package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

func (Photo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("photos").Unique(),
	}
}

func (Photo) Hooks() []ent.Hook {
	return []ent.Hook{
		UploadBlob(),
	}
}
