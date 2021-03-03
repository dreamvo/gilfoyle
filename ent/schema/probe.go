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



// Probe holds the schema definition for the Probe entity.
type Probe struct {
	ent.Schema
}

func (Probe) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "meida_probe"},
	}
}

func (Probe) Index() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

// Fields of the Probe.
func (Probe) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(func() uuid.UUID {
			return uuid.New()
		}).Immutable(),
		field.String("filename").MinLen(1).MaxLen(255),
		field.String("mimetype").MinLen(1).MaxLen(255),
		field.Int("filesize").Default(0).Min(0),
		field.String("checksum_sha256").MinLen(64).MaxLen(64),
		field.String("aspect_ratio").MinLen(3).MaxLen(5).Default("16:9"),
		field.Int("width").Min(1),
		field.Int("height").Min(1),
		field.Float("duration_seconds").Min(0).Default(0),
		field.Int("video_bitrate").Min(0).Default(0),
		field.Int("audio_bitrate").Min(0).Default(0),
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

// Edges of the Probe.
func (Probe) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("media", Media.Type).
			Ref("probe").
			StructTag(`json"media,omitempty"`).
			Required().
			Unique(),
	}
}
