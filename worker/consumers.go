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
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

			w.logger.Info("Received a message", zap.String("MediaUUID", body.MediaUUID.String()), zap.String("OriginalFilePath", body.OriginalFile.Filepath))

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
			fmt.Println(111, strings.Join(args, " "))
			cmd := exec.Command(w.ffmpegConfig.FfmpegBinPath, args...)
			err = cmd.Run()
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
				SetResolutionWidth(uint16(body.VideoWidth)).
				SetResolutionHeight(uint16(body.VideoHeight)).
				SetDurationSeconds(body.OriginalFile.DurationSeconds).
				SetFramerate(body.OriginalFile.FrameRate).
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
