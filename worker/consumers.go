package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func setMediaStatusAck(w *Worker, d amqp.Delivery, uuid uuid.UUID, status media.Status) error {
	_, err := w.dbClient.Media.UpdateOneID(uuid).SetStatus(status).Save(context.Background())
	if err != nil {
		return err
	}
	err = d.Nack(false, false)
	if err != nil {
		return err
	}
	return nil
}

func videoTranscodingConsumer(w *Worker, msgs <-chan amqp.Delivery) {
	ctx := context.Background()
	for d := range msgs {
		func() {
			var body VideoTranscodingParams

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				w.logger.Error("Unmarshal error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			w.logger.Info("Received a message", zap.String("MediaUUID", body.MediaUUID.String()), zap.String("OriginalFilePath", body.OriginalFilePath))

			m, err := w.dbClient.Media.Query().
				Where(media.ID(body.MediaUUID)).
				WithMediaFiles().
				Only(context.Background())
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				return
			}

			r, err := w.storage.Open(ctx, fmt.Sprintf("%s/%s", body.MediaUUID.String(), transcoding.OriginalFileName))
			if err != nil {
				w.logger.Error("Storage error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}
			defer func() { _ = r.Close() }()

			ctxWithTimeout, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			probe, err := ffprobe.ProbeReader(ctxWithTimeout, r)
			if err != nil {
				w.logger.Error("ffprobe error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			srcTmpPath := filepath.Join(os.TempDir(), uuid.New().String())
			dstTmpPath := filepath.Join(os.TempDir(), uuid.New().String())

			err = os.MkdirAll(dstTmpPath, 0755)
			if err != nil {
				w.logger.Error("Error creating destination directory", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			out, err := os.Create(srcTmpPath)
			if err != nil {
				w.logger.Error("File create error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}
			defer func() { _ = out.Close() }()

			_, err = io.Copy(out, r)
			if err != nil {
				w.logger.Error("File copy error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			originalFormat := strings.Split(probe.Format.FormatName, ",")[0]

			args := []string{
				"-i",
				srcTmpPath,
				"-vf",
				fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease", body.VideoWidth, body.VideoHeight),
				"-c:a",
				body.AudioCodec,
				"-ar",
				fmt.Sprintf("%d", body.AudioRate),
				"-b:a",
				fmt.Sprintf("%d", body.AudioBitrate),
				"-c:v",
				body.VideoCodec,
				"-r",
				fmt.Sprintf("%d", body.FrameRate),
				"-f",
				originalFormat,
				"-profile:v",
				"main",
				"-crf",
				fmt.Sprintf("%d", body.Crf),
				"-g",
				fmt.Sprintf("%d", body.KeyframeInterval),
				"-b:v",
				fmt.Sprintf("%d", body.VideoBitRate),
				"-maxrate",
				fmt.Sprintf("%d", body.VideoMaxBitRate),
				"-bufsize",
				fmt.Sprintf("%d", body.BufferSize),
				"-hls_time",
				fmt.Sprintf("%d", body.HlsSegmentDuration),
				"-hls_playlist_type",
				body.HlsPlaylistType,
				"-hls_segment_filename",
				dstTmpPath + "/%03d.ts",
				dstTmpPath + "/" + transcoding.HLSPlaylistFilename,
			}
			err = exec.Command(w.ffmpegConfig.FfmpegBinPath, args...).Run()
			if err != nil {
				w.logger.Error("Command execution error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			files, err := filepath.Glob(fmt.Sprintf("%s/*", dstTmpPath))
			if err != nil {
				w.logger.Error("Directory walk error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			for _, file := range files {
				f, err := os.Open(file)
				if err != nil {
					w.logger.Error("File open error", zap.Error(err))
					_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
					return
				}
				defer func() { _ = f.Close() }()

				filename := strings.Replace(file, dstTmpPath, "", 1)

				// Save the file
				err = w.storage.Save(ctx, f, fmt.Sprintf("%s/%s/%s", body.MediaUUID.String(), body.RenditionName, filename))
				if err != nil {
					w.logger.Error("Storage error", zap.Error(err))
					_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
					return
				}
			}

			_, err = w.dbClient.MediaFile.Create().
				SetMedia(m).
				SetRenditionName(body.RenditionName).
				SetFormat("hls").
				SetVideoBitrate(int64(body.VideoBitRate)).
				SetScaledWidth(int16(body.VideoWidth)).
				SetDurationSeconds(m.Edges.MediaFiles[0].DurationSeconds).
				SetFramerate(m.Edges.MediaFiles[0].Framerate).
				SetMediaType(schema.MediaFileTypeVideo).
				Save(ctx)
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			err = d.Ack(false)
			if err != nil {
				w.logger.Error("Error trying to send ack", zap.Error(err))
				_ = setMediaStatusAck(w, d, body.MediaUUID, media.StatusErrored)
			}
		}()
	}
}
