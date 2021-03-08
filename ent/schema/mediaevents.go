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

// MediaEvents holds the schema definition for the MediaEvents entity.
type MediaEvents struct {
	ent.Schema
}

func (MediaEvents) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "media_events"},
	}
}

func (MediaEvents) Index() []ent.Index {
	return []ent.Index{
		index.Fields("id").
			Unique(),
	}
}

// Fields of the MediaEvents.
func (MediaEvents) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(func() uuid.UUID {
			return uuid.New()
		}).Immutable(),
		field.Enum("level").Values("Normal", "Warning", "Error").Default("Normal"),
		field.String("reason").MinLen(1).MaxLen(255).Immutable(),
		field.String("message").MinLen(1).MaxLen(255).Immutable(),
		field.Time("created_at").Default(func() time.Time {
			return time.Now()
		}).Immutable(),
	}
}

// Edges of the MediaEvents.
func (MediaEvents) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("media", Media.Type).
			Ref("events").
			StructTag(`json"media,omitempty"`).
			Required().
			Unique(),
	}
}
