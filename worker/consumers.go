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
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

func encodingEntrypointConsumer(w *Worker, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		func() {
			ctx := context.Background()
			var body MediaEncodingEntrypoint

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				w.logger.Error("Unmarshal error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			w.logger.Info("Starting media encoding", zap.String("MediaUUID", body.MediaUUID.String()))

			m, err := w.dbClient.Media.Get(ctx, body.MediaUUID)
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				_ = d.Nack(false, true)
				return
			}

			file, err := w.storage.Open(ctx, path.Join(m.ID.String(), m.OriginalFilename))
			if err != nil {
				w.logger.Error("Storage error", zap.Error(err))
				_ = d.Nack(false, true)
				return
			}

			ctxWithTimeout, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			// Run original file analysis
			data, err := ffprobe.ProbeReader(ctxWithTimeout, file)
			if err != nil {
				w.logger.Error("FFprobe error", zap.Error(err))
				_ = d.Nack(false, true)
				return
			}

			videoStreams := data.StreamType(ffprobe.StreamVideo)
			if len(videoStreams) != 1 {
				w.logger.Error("File input error", zap.Error(errors.New("uploaded media must have only 1 video stream")))
				_ = d.Nack(false, false)
				return
			}

			audioStreams := data.StreamType(ffprobe.StreamAudio)
			if len(audioStreams) != 1 {
				w.logger.Error("File input error", zap.Error(errors.New("uploaded media must have only 1 audio stream")))
				_ = d.Nack(false, false)
				return
			}

			videoStream := data.FirstVideoStream()
			framerate := int(transcoding.ParseFrameRates(videoStream.RFrameRate))

			mime, err := mimetype.DetectReader(file)
			if err != nil {
				w.logger.Error("Mimetype detection", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				w.logger.Error("File to bytes conversion error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			checksum := sha256.Sum256(fileBytes)

			probe, err := w.dbClient.Probe.
				Create().
				SetMedia(m).
				SetFilename(m.OriginalFilename).
				SetAspectRatio(videoStream.DisplayAspectRatio).
				SetFilesize(0).
				SetVideoBitrate(0).
				SetAudioBitrate(0).
				SetWidth(videoStream.Width).
				SetHeight(videoStream.Height).
				SetDurationSeconds(data.Format.Duration().Seconds()).
				SetChecksumSha256(string(checksum[:])).
				SetMimetype(mime.String()).
				SetFramerate(framerate).
				SetFormat(data.Format.FormatName).
				SetNbStreams(data.Format.NBStreams).
				Save(ctx)
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			// Schedule encoding jobs
			for _, rendition := range w.config.Settings.Encoding.Renditions {
				// Ignore resolutions higher than original
				if rendition.Width > videoStream.Width && rendition.Height > videoStream.Height {
					continue
				}

				// Calculate resolution with original aspect ratio
				var width = uint16((videoStream.Width / videoStream.Height) * rendition.Height)
				var height = uint16((videoStream.Height / videoStream.Width) * rendition.Width)

				mediaFile, err := w.dbClient.MediaFile.
					Create().
					SetMedia(m).
					SetRenditionName(rendition.Name).
					SetMediaType(schema.MediaFileTypeVideo).
					SetFormat("hls").
					SetStatus(mediafile.StatusProcessing).
					SetEntryFile("index.m3u8").
					SetMimetype("application/x-mpegURL").
					SetTargetBandwidth(uint64(rendition.VideoBitrate + rendition.AudioBitrate)).
					SetVideoBitrate(int64(rendition.VideoBitrate)).
					SetAudioBitrate(int64(rendition.AudioBitrate)).
					SetVideoCodec(rendition.VideoCodec).
					SetAudioCodec(rendition.AudioCodec).
					SetResolutionWidth(width).
					SetResolutionHeight(height).
					SetDurationSeconds(probe.DurationSeconds).
					SetFramerate(uint8(rendition.Framerate)).
					Save(ctx)
				if err != nil {
					w.logger.Error("Database error", zap.Error(err))
					_ = d.Nack(false, false)
					return
				}

				ch, err := w.Client.Channel()
				if err != nil {
					w.logger.Error("Worker channel", zap.Error(err))
					_ = d.Nack(false, false)
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
					_ = d.Nack(false, false)
					return
				}
			}
		}()
	}
}

func hlsVideoEncodingConsumer(w *Worker, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		func() {
			var body HlsVideoEncodingParams
			ctx := context.Background()

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				w.logger.Error("Unmarshal error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			w.logger.Info("Received video transcoding message", zap.String("MediaFileUUID", body.MediaFileUUID.String()))

			m, err := w.dbClient.MediaFile.
				Query().
				Where(mediafile.ID(body.MediaFileUUID)).
				WithMedia().
				Only(context.Background())
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				return
			}

			r, err := w.storage.Open(ctx, path.Join(m.Edges.Media.String(), m.Edges.Media.OriginalFilename))
			if err != nil {
				w.logger.Error("Storage error", zap.Error(err))
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}
			defer func() { _ = r.Close() }()

			srcTmpPath := filepath.Join(os.TempDir(), uuid.New().String())
			dstTmpPath := filepath.Join(os.TempDir(), uuid.New().String())

			err = os.MkdirAll(dstTmpPath, 0755)
			if err != nil {
				w.logger.Error("Error creating destination directory", zap.Error(err))
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			out, err := os.Create(srcTmpPath)
			if err != nil {
				w.logger.Error("File create error", zap.Error(err))
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}
			defer func() { _ = out.Close() }()

			_, err = io.Copy(out, r)
			if err != nil {
				w.logger.Error("File copy error", zap.Error(err))
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			VideoFilter := fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease", m.ResolutionWidth, m.ResolutionHeight)
			HlsSegmentFilename := dstTmpPath + "/%03d.ts"
			VideoProfile := "main"
			preset := "medium"
			overwrite := true
			AudioBitrate := int(m.AudioBitrate)
			VideoBitrate := int(m.VideoBitrate)
			Framerate := int(m.Framerate)

			p := w.transcoder.
				Process().
				SetInput(srcTmpPath).
				SetOutput(fmt.Sprintf("%s/%s", dstTmpPath, transcoding.HLSPlaylistFilename)).
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
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
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
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			_, err = w.dbClient.MediaFile.UpdateOne(m).
				SetStatus(mediafile.StatusReady).
				Save(ctx)
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
				return
			}

			err = d.Ack(false)
			if err != nil {
				w.logger.Error("Error trying to send ack", zap.Error(err))
				//_ = setMediaStatusNack(w, d, body.MediaUUID, media.StatusErrored)
			}
		}()
	}
}

func mediaEncodingCallbackConsumer(w *Worker, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		func() {
			ctx := context.Background()
			var body MediaEncodingCallbackParams

			err := json.Unmarshal(d.Body, &body)
			if err != nil {
				w.logger.Error("Unmarshal error", zap.Error(err))
				_ = d.Nack(false, false)
				return
			}

			w.logger.Info("Received media callback message", zap.String("MediaUUID", body.MediaUUID.String()), zap.Int("MediaFilesCount", body.MediaFilesCount))

			m, err := w.dbClient.Media.Query().Where(media.ID(body.MediaUUID)).WithMediaFiles().Only(ctx)
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				_ = d.Nack(false, true)
				return
			}

			if len(m.Edges.MediaFiles) != body.MediaFilesCount {
				time.Sleep(2 * time.Second)
				_ = d.Nack(false, m.Status == media.StatusProcessing)
				return
			}

			masterPlaylist := transcoding.CreateMasterPlaylist(m.Edges.MediaFiles)

			err = w.storage.Save(ctx, strings.NewReader(masterPlaylist), path.Join(m.ID.String(), transcoding.HLSMasterPlaylistFilename))
			if err != nil {
				w.logger.Error("Storage error", zap.Error(err))
				_ = d.Nack(false, true)
				return
			}

			_, err = w.dbClient.Media.UpdateOne(m).SetStatus(media.StatusReady).Save(ctx)
			if err != nil {
				w.logger.Error("Database error", zap.Error(err))
				_ = d.Nack(false, true)
				return
			}

			err = d.Ack(false)
			if err != nil {
				w.logger.Error("Error trying to send ack", zap.Error(err))
			}
		}()
	}
}
