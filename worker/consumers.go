package worker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/mediafile"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func setMediaStatusNack(w *Worker, d amqp.Delivery, uuid uuid.UUID, status media.Status, errMessage error) error {
	_, err := w.dbClient.Media.
		UpdateOneID(uuid).
		SetStatus(status).
		SetMessage(errMessage.Error()).
		Save(context.Background())
	if err != nil {
		return err
	}
	err = d.Nack(false, false)
	if err != nil {
		return err
	}
	return nil
}

func setRenditionStatusNack(w *Worker, d amqp.Delivery, uuid uuid.UUID, status mediafile.Status, errMessage error) error {
	_, err := w.dbClient.MediaFile.
		UpdateOneID(uuid).
		SetStatus(status).
		SetMessage(errMessage.Error()).
		Save(context.Background())
	if err != nil {
		return err
	}
	err = d.Nack(false, false)
	if err != nil {
		return err
	}
	return nil
}

func encodingEntrypointConsumer(w *Worker, d amqp.Delivery) {
	ctx := context.Background()
	var body EncodingEntrypointParams

	err := json.Unmarshal(d.Body, &body)
	if err != nil {
		w.logger.Error("Unmarshal error", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}

	w.logger.Info("Media encoding entrypoint", zap.String("MediaUUID", body.MediaUUID.String()))

	m, err := w.dbClient.Media.Get(ctx, body.MediaUUID)
	if err != nil {
		w.logger.Error("Database error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	file, err := w.storage.Open(ctx, path.Join(m.ID.String(), m.OriginalFilename))
	if err != nil {
		w.logger.Error("Storage error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	ctxWithTimeout, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	// Run original file analysis
	data, err := ffprobe.ProbeReader(ctxWithTimeout, file)
	if err != nil {
		w.logger.Error("ffprobe command failed", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	videoStreams := data.StreamType(ffprobe.StreamVideo)
	if len(videoStreams) == 0 {
		err := errors.New("uploaded media should have at least 1 video stream")
		w.logger.Error("File input error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	videoStream := data.FirstVideoStream()
	nbFrames, _ := strconv.ParseFloat(videoStream.NbFrames, 64)
	framerate := nbFrames / data.Format.Duration().Seconds()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		w.logger.Error("Failed to generate checksum", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	probe, err := w.dbClient.Probe.
		Create().
		SetMedia(m).
		SetFilename(m.OriginalFilename).
		//SetAspectRatio(videoStream.DisplayAspectRatio).
		SetFilesize(0).
		SetVideoBitrate(0).
		SetAudioBitrate(0).
		SetWidth(videoStream.Width).
		SetHeight(videoStream.Height).
		SetDurationSeconds(data.Format.Duration().Seconds()).
		SetChecksumSha256(fmt.Sprintf("%x", h.Sum(nil))).
		SetFramerate(framerate).
		SetFormat(data.Format.FormatName).
		SetNbStreams(data.Format.NBStreams).
		Save(ctx)
	if err != nil {
		w.logger.Error("Database error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	ch, err := w.Client.Channel()
	if err != nil {
		w.logger.Error("Worker channel", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	// Schedule encoding jobs
	for _, rendition := range w.config.Settings.Encoding.Renditions {
		// Ignore resolutions higher than original
		if rendition.Width > probe.Width && rendition.Height > probe.Height {
			continue
		}

		// Calculate resolution with original aspect ratio
		var width = rendition.Width
		var height = math.Round((float64(probe.Height) / float64(probe.Width)) * float64(rendition.Width))

		mediaFile, err := w.dbClient.MediaFile.
			Create().
			SetMedia(m).
			SetRenditionName(rendition.Name).
			SetMediaType(schema.MediaFileTypeVideo).
			SetFormat("hls").
			SetStatus(mediafile.StatusProcessing).
			SetEntryFile(transcoding.HLSPlaylistFilename).
			SetMimetype(transcoding.HLSPlaylistMimeType).
			SetTargetBandwidth(uint64(rendition.VideoBitrate + rendition.AudioBitrate)).
			SetVideoBitrate(int64(rendition.VideoBitrate)).
			SetAudioBitrate(int64(rendition.AudioBitrate)).
			SetVideoCodec(rendition.VideoCodec).
			SetAudioCodec(rendition.AudioCodec).
			SetResolutionWidth(uint16(width)).
			SetResolutionHeight(uint16(height)).
			SetDurationSeconds(probe.DurationSeconds).
			SetFramerate(uint8(rendition.Framerate)).
			Save(ctx)
		if err != nil {
			w.logger.Error("Database error", zap.Error(err))
			_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
			return
		}

		err = HlsVideoEncodingProducer(ch, HlsVideoEncodingParams{
			MediaFileUUID:      mediaFile.ID,
			KeyframeInterval:   48,
			HlsSegmentDuration: 4,
			HlsPlaylistType:    "vod",
		})
		if err != nil {
			w.logger.Error("HlsVideoEncoding job creation", zap.Error(err))
			_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
			return
		}
	}

	err = EncodingFinalizerProducer(ch, EncodingFinalizerParams{MediaUUID: body.MediaUUID})
	if err != nil {
		w.logger.Error("EncodingFinalizer job creation", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
	}

	err = d.Ack(false)
	if err != nil {
		w.logger.Error("Error trying to send ack", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
	}
}

func hlsVideoEncodingConsumer(w *Worker, d amqp.Delivery) {
	var body HlsVideoEncodingParams
	ctx := context.Background()

	err := json.Unmarshal(d.Body, &body)
	if err != nil {
		w.logger.Error("Unmarshal error", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}

	w.logger.Info("Received HLS video encoding message", zap.String("MediaFileUUID", body.MediaFileUUID.String()))

	m, err := w.dbClient.MediaFile.
		Query().
		Where(mediafile.ID(body.MediaFileUUID)).
		WithMedia().
		Only(context.Background())
	if err != nil {
		w.logger.Error("Database error", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}

	r, err := w.storage.Open(ctx, path.Join(m.Edges.Media.ID.String(), m.Edges.Media.OriginalFilename))
	if err != nil {
		w.logger.Error("Storage error", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}
	defer func() { _ = r.Close() }()

	srcTmpPath := filepath.Join(os.TempDir(), uuid.New().String())
	dstTmpPath := filepath.Join(os.TempDir(), uuid.New().String())

	err = os.MkdirAll(dstTmpPath, 0755)
	if err != nil {
		w.logger.Error("Error creating destination directory", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}

	out, err := os.Create(srcTmpPath)
	if err != nil {
		w.logger.Error("File create error", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, r)
	if err != nil {
		w.logger.Error("File copy error", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}

	VideoFilter := fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease", m.ResolutionWidth, m.ResolutionHeight)
	HlsSegmentFilename := path.Join(dstTmpPath, "/%03d.ts")
	VideoProfile := "main"
	preset := "medium"
	overwrite := true
	AudioBitrate := int(m.AudioBitrate)
	VideoBitrate := int(m.VideoBitrate)
	Framerate := int(m.Framerate)

	p := w.transcoder.
		Process().
		SetInput(srcTmpPath).
		SetOutput(path.Join(dstTmpPath, transcoding.HLSPlaylistFilename)).
		WithOptions(transcoding.ProcessOptions{
			AudioCodec:         &m.AudioCodec,
			VideoCodec:         &m.VideoCodec,
			AudioBitrate:       &AudioBitrate,
			VideoBitRate:       &VideoBitrate,
			FrameRate:          &Framerate,
			HlsSegmentDuration: &body.HlsSegmentDuration,
			HlsPlaylistType:    &body.HlsPlaylistType,
			HlsSegmentFilename: &HlsSegmentFilename,
			VideoProfile:       &VideoProfile,
			KeyframeInterval:   &body.KeyframeInterval,
			VideoFilter:        &VideoFilter,
			Overwrite:          &overwrite,
			Preset:             &preset,
		})
	err = w.transcoder.Run(p)
	if err != nil {
		w.logger.Error("Command execution error", zap.Error(err), zap.String("arguments", strings.Join(p.GetStrArguments(), " ")))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}

	err = filepath.Walk(dstTmpPath, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		filename := strings.Replace(filepath, dstTmpPath, "", 1)

		// Save the file
		err = w.storage.Save(ctx, f, path.Join(m.Edges.Media.ID.String(), m.RenditionName, filename))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		w.logger.Error("Storage error (walk)", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}

	_, err = w.dbClient.MediaFile.
		UpdateOne(m).
		SetStatus(mediafile.StatusReady).
		SetMessage("Encoding job succeeded").
		Save(ctx)
	if err != nil {
		w.logger.Error("Database error", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
		return
	}

	err = d.Ack(false)
	if err != nil {
		w.logger.Error("Error trying to send ack", zap.Error(err))
		_ = setRenditionStatusNack(w, d, body.MediaFileUUID, mediafile.StatusErrored, err)
	}
}

func encodingFinalizerConsumer(w *Worker, d amqp.Delivery) {
	ctx := context.Background()
	var body EncodingFinalizerParams

	err := json.Unmarshal(d.Body, &body)
	if err != nil {
		w.logger.Error("Unmarshal error", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}

	w.logger.Info("Received encoding finalizer message", zap.String("MediaUUID", body.MediaUUID.String()))

	m, err := w.dbClient.Media.Query().Where(media.ID(body.MediaUUID)).WithMediaFiles().Only(ctx)
	if err != nil {
		w.logger.Error("Database error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	// Create master playlist
	masterPlaylistContent := transcoding.CreateMasterPlaylist(m.Edges.MediaFiles)
	err = w.storage.Save(ctx, strings.NewReader(masterPlaylistContent), path.Join(m.ID.String(), transcoding.HLSMasterPlaylistFilename))
	if err != nil {
		w.logger.Error("Storage error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	mediaStatus := media.StatusErrored
	mediaMessage := "All encoding jobs failed"

	if len(m.Edges.MediaFiles) == 0 {
		mediaMessage = "Media doesn't have any rendition"
	}

	for _, r := range m.Edges.MediaFiles {
		// If at least one rendition is still in
		// processing state then requeue the job in 15 sec.
		if r.Status == mediafile.StatusProcessing {
			go func(w *Worker, d amqp.Delivery, body EncodingFinalizerParams) {
				time.Sleep(15 * time.Second)

				ch, err := w.Client.Channel()
				if err != nil {
					w.logger.Error("Worker channel", zap.Error(err))
					_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
					return
				}

				err = EncodingFinalizerProducer(ch, body)
				if err != nil {
					w.logger.Error("Error creating encoding finalizer job", zap.Error(err))
					_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
				}
			}(w, d, body)

			mediaStatus = media.StatusProcessing
			mediaMessage = "Media is not yet available for streaming"
		}

		// If at least one media file is ready
		// the media is now available for streaming
		// so we set the media status to ready.
		if r.Status == mediafile.StatusReady {
			mediaStatus = media.StatusReady
			mediaMessage = "One or more rendition succeeded. Media is available for streaming"
		}
	}

	_, err = w.dbClient.Media.UpdateOne(m).SetStatus(mediaStatus).SetMessage(mediaMessage).Save(ctx)
	if err != nil {
		w.logger.Error("Database error", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
		return
	}

	err = d.Ack(false)
	if err != nil {
		w.logger.Error("Error trying to send ack", zap.Error(err))
		_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored, err)
	}
}
