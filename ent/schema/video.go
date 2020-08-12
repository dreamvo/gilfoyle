package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"time"
)

// VideoStatusProcessing relates to the initial state of a video
const VideoStatusProcessing = "processing"

// VideoStatusReady relates to the final state of a video
const VideoStatusReady = "ready"

// Video holds the schema definition for the Video entity.
type Video struct {
	ent.Schema
}

func (Video) Config() ent.Config {
	return ent.Config{
		Table: "video",
	}
}

func (Video) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.Fields("uuid").
			Unique(),
	}
}

// Fields of the Video.
func (Video) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").Unique().MinLen(36).MaxLen(36),
		field.String("title").NotEmpty().MinLen(1).MaxLen(255),
		field.Enum("status").Values(VideoStatusProcessing, VideoStatusReady),
		field.Time("created_at").Default(func() time.Time {
			return time.Now()
		}),
		field.Time("updated_at").Default(func() time.Time {
			return time.Now()
		}),
	}
}

// Edges of the Video.
func (Video) Edges() []ent.Edge {
	return nil
}
