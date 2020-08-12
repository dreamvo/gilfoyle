package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

const VideoStatusCreated = "created"
const VideoStatusReady = "ready"

// Video holds the schema definition for the Video entity.
type Video struct {
	ent.Schema
}

// Fields of the Video.
func (Video) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").Unique(),
		field.String("title").NotEmpty(),
		field.Enum("status").Values(VideoStatusCreated, VideoStatusReady),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

// Edges of the Video.
func (Video) Edges() []ent.Edge {
	return nil
}
