package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/floostack/transcoder/ffmpeg"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
)

func setMediaStatus(w *Worker, uuid uuid.UUID, status media.Status) error {
	_, err := w.dbClient.Media.UpdateOneID(uuid).SetStatus(status).Save(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func videoTranscodingConsumer(w *Worker, msg <-chan amqp.Delivery) {
	ctx := context.Background()
	for d := range msg {
		var body VideoTranscodingParams

		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			w.logger.Error("Unmarshal error", zap.Error(err))
			return
		}

		w.logger.Info("Received a message", zap.String("MediaUUID", body.MediaUUID.String()), zap.String("SourceFilePath", body.SourceFilePath))

		m, err := w.dbClient.Media.Query().
			Where(media.ID(body.MediaUUID)).
			WithMediaFiles().
			Only(context.Background())
		if err != nil {
			w.logger.Error("Database error", zap.Error(err))
			return
		}

		// Create a new MediaFile for this Media
		r, err := w.storage.Open(ctx, fmt.Sprintf("%s/%s", body.MediaUUID.String(), transcoding.OriginalFileName))
		if err != nil {
			w.logger.Error("Storage error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}
		defer func() { _ = r.Close() }()

		srcTmpPath := filepath.Join(os.TempDir(), uuid.New().String())
		dstTmpPath := filepath.Join(os.TempDir(), uuid.New().String())

		out, err := os.Create(srcTmpPath)
		if err != nil {
			w.logger.Error("File create error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}
		defer func() { _ = out.Close() }()

		_, err = io.Copy(out, r)
		if err != nil {
			w.logger.Error("File copy error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}

		overwrite := true
		opts := ffmpeg.Options{
			Overwrite:          &overwrite,
			VideoFilter:        nil,
			AudioCodec:         nil,
			AudioRate:          nil,
			VideoCodec:         nil,
			VideoProfile:       nil,
			Crf:                nil,
			KeyframeInterval:   nil,
			HlsSegmentDuration: nil,
			HlsPlaylistType:    nil,
			VideoBitRate:       nil,
			VideoMaxBitRate:    nil,
			BufferSize:         nil,
			AudioBitrate:       nil,
			HlsSegmentFilename: nil,
		}

		progress, err := w.transcoder.
			Input(srcTmpPath).
			Output(dstTmpPath).
			WithOptions(opts).
			Start(opts)
		if err != nil {
			w.logger.Error("ffmpeg error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}

		for msg := range progress {
			w.logger.Info("Video transcoding progress", zap.String("MediaUUID", body.MediaUUID.String()), zap.String("CurrentFramesProcessed", fmt.Sprintf("%+v", msg.GetFramesProcessed())), zap.Float64("progress", msg.GetProgress()))
		}

		f, err := os.Open(dstTmpPath)
		if err != nil {
			w.logger.Error("File open error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}
		defer func() { _ = f.Close() }()

		// Save the created files
		err = w.storage.Save(ctx, f, fmt.Sprintf("%s/%s", body.MediaUUID.String(), "medium.m3u8"))
		if err != nil {
			w.logger.Error("Storage error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}

		_, err = w.dbClient.MediaFile.Create().
			SetMedia(m).
			SetVideoBitrate(100).
			SetEncoderPreset(schema.MediaFileEncoderPresetMedium).
			SetDurationSeconds(100).
			SetScaledWidth(1280).
			SetFramerate(transcoding.ParseFrameRates("25/1")).
			SetMediaType(schema.MediaFileTypeVideo).
			Save(ctx)
		if err != nil {
			w.logger.Error("Database error", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
			return
		}

		err = d.Ack(false)
		if err != nil {
			w.logger.Error("Error trying to send ack", zap.Error(err))
			_ = setMediaStatus(w, body.MediaUUID, media.StatusErrored)
		}
	}
}
