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

// MediaStatusAwaitingUpload relates to the pending state of a media
const MediaStatusAwaitingUpload = "AwaitingUpload"

// MediaStatusErrored relates to the errored state of a media
const MediaStatusErrored = "Errored"

// MediaStatusProcessing relates to the initial state of a media
const MediaStatusProcessing = "Processing"

// MediaStatusPending relates to the pending state of a media
const MediaStatusPending = "Pending"

// MediaStatusReady relates to the final state of a media
const MediaStatusReady = "Ready"

// Media holds the schema definition for the Media entity.
type Media struct {
	ent.Schema
}

func (Media) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "media"},
	}
}

func (Media) Index() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

// Fields of the Media.
func (Media) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(func() uuid.UUID {
			return uuid.New()
		}).Immutable(),
		field.String("title").NotEmpty().MinLen(1).MaxLen(255),
		field.String("original_filename").Optional().MaxLen(150).Default(""),
		field.Enum("status").Values(MediaStatusAwaitingUpload, MediaStatusProcessing, MediaStatusReady, MediaStatusErrored),
		field.String("message").Optional().MaxLen(255).Default(""),
		field.Bool("playable").Default(false).Optional(),
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

// Edges of the Media.
func (Media) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("media_files", MediaFile.Type).
			StorageKey(edge.Column("media")).
			StructTag(`json:"media_files,omitempty"`),
		edge.To("probe", Probe.Type).
			StorageKey(edge.Column("media")).
			StructTag(`json:"probe,omitempty"`).
			Unique(),
		edge.To("events", MediaEvents.Type).
			StorageKey(edge.Column("media")).
			StructTag(`json:"events,omitempty"`),
	}
}
