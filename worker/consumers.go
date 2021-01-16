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
	"path/filepath"
	"strings"
)

func setMediaStatusNack(w *Worker, d amqp.Delivery, uuid uuid.UUID, status media.Status) error {
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
	for d := range msgs {
		func() {
			var body VideoTranscodingParams
			ctx := context.Background()

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				w.logger.Error("Unmarshal error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			w.logger.Info("Received a message", zap.String("MediaUUID", body.MediaUUID.String()))

			m, err := w.dbClient.Media.
				Query().
				Where(media.ID(body.MediaUUID)).
				Only(context.Background())
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				return
			}

			r, err := w.storage.Open(ctx, fmt.Sprintf("%s/%s", body.MediaUUID.String(), transcoding.OriginalFileName))
			if err != nil {
				w.logger.Error("Storage error", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}
			defer func() { _ = r.Close() }()

			srcTmpPath := filepath.Join(os.TempDir(), uuid.New().String())
			dstTmpPath := filepath.Join(os.TempDir(), uuid.New().String())

			err = os.MkdirAll(dstTmpPath, 0755)
			if err != nil {
				w.logger.Error("Error creating destination directory", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			out, err := os.Create(srcTmpPath)
			if err != nil {
				w.logger.Error("File create error", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}
			defer func() { _ = out.Close() }()

			_, err = io.Copy(out, r)
			if err != nil {
				w.logger.Error("File copy error", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			HlsSegmentFilename := dstTmpPath + "/%03d.ts"
			VideoProfile := "main"
			VideoFilter := fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease", body.VideoWidth, body.VideoHeight)

			p := w.transcoder.
				Process().
				Input(srcTmpPath).
				Output(fmt.Sprintf("%s/%s", dstTmpPath, transcoding.HLSPlaylistFilename)).
				WithOptions(transcoding.ProcessOptions{
					AudioCodec:         &body.AudioCodec,
					VideoCodec:         &body.VideoCodec,
					AudioRate:          &body.AudioRate,
					AudioBitrate:       &body.AudioBitrate,
					VideoBitRate:       &body.VideoBitRate,
					VideoMaxBitRate:    &body.VideoMaxBitRate,
					FrameRate:          &body.FrameRate,
					BufferSize:         &body.BufferSize,
					HlsSegmentDuration: &body.HlsSegmentDuration,
					HlsPlaylistType:    &body.HlsPlaylistType,
					HlsSegmentFilename: &HlsSegmentFilename,
					VideoProfile:       &VideoProfile,
					Crf:                &body.Crf,
					KeyframeInterval:   &body.KeyframeInterval,
					VideoFilter:        &VideoFilter,
				})
			err = w.transcoder.Run(p)
			if err != nil {
				w.logger.Error("Command execution error", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			err = filepath.Walk(dstTmpPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				f, err := os.Open(path)
				if err != nil {
					return err
				}
				defer func() { _ = f.Close() }()

				filename := strings.Replace(path, dstTmpPath, "", 1)

				// Save the file
				err = w.storage.Save(ctx, f, fmt.Sprintf("%s/%s/%s", body.MediaUUID.String(), body.RenditionName, filename))
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				w.logger.Error("Storage error (walk)", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
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
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			err = d.Ack(false)
			if err != nil {
				w.logger.Error("Error trying to send ack", zap.Error(err))
				_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
			}
		}()
	}
}
