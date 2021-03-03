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
		}).Immutable(),
		field.String("rendition_name").MinLen(1).MaxLen(100),
		field.String("format").MinLen(1),
		field.Uint64("target_bandwidth").Default(800000),
		field.Int64("video_bitrate").Min(0),
		field.Int64("audio_bitrate").Min(0),
		field.String("video_codec").MinLen(1).Default("h264"),
		field.String("audio_codec").MinLen(1).Default("aac"),
		field.Uint16("resolution_width").Min(8).Max(2160),
		field.Uint16("resolution_height").Min(8).Max(2160),
		field.Uint8("framerate").Min(8).Max(120),
		field.Float("duration_seconds").Min(0),
		field.Enum("media_type").Values(MediaFileTypeAudio, MediaFileTypeVideo),
		field.Enum("status").Values(MediaStatusProcessing, MediaStatusReady, MediaStatusErrored),
		field.String("message").Optional().MaxLen(255).Default(""),
		field.String("entry_file").MaxLen(255).Default("index.m3u8"),
		field.String("mimetype").MaxLen(255).Default("application/x-mpegURL"),
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
			StructTag(`json"media,omitempty"`).
			Required().
			Unique(),
	}
}
