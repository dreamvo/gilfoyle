package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/entsql"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

const (
	// MediaFileKindAudio relates to the audio type
	MediaFileTypeAudio = "audio"

	// MediaFileTypeVideo relates to the video type
	MediaFileTypeVideo = "video"
)

const (
	MediaFileEncoderPresetSource    = "source"
	MediaFileEncoderPresetUltraFast = "ultrafast"
	MediaFileEncoderPresetVeryFast  = "veryfast"
	MediaFileEncoderPresetFast      = "fast"
	MediaFileEncoderPresetMedium    = "medium"
	MediaFileEncoderPresetSlow      = "slow"
	MediaFileEncoderPresetVerySlow  = "veryslow"
)

// MediaFile holds the schema definition for the MediaFile entity.
type MediaFile struct {
	ent.Schema
}

func (MediaFile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "media_file"},
	}
}

func (MediaFile) Index() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

// Fields of the MediaFile.
func (MediaFile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(func() uuid.UUID {
			return uuid.New()
		}),
		field.Int64("video_bitrate").Min(0),
		field.Int16("scaled_width").Min(144).Max(2160),
		field.Enum("encoder_preset").Values(
			MediaFileEncoderPresetSource,
			MediaFileEncoderPresetUltraFast,
			MediaFileEncoderPresetVeryFast,
			MediaFileEncoderPresetFast,
			MediaFileEncoderPresetMedium,
			MediaFileEncoderPresetSlow,
			MediaFileEncoderPresetVerySlow,
		),
		field.Int8("framerate").Min(12).Max(60),
		field.Float("duration_seconds").Min(0),
		field.Enum("media_type").Values(MediaFileTypeAudio, MediaFileTypeVideo),
		field.Time("created_at").Default(func() time.Time {
			return time.Now()
		}),
		field.Time("updated_at").
			Default(func() time.Time {
				return time.Now()
			}).
			UpdateDefault(func() time.Time {
				return time.Now()
			}),
	}
}

// Edges of the MediaFile.
func (MediaFile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("media", Media.Type).
			Ref("media_files").
			Required().
			Unique(),
	}
}
