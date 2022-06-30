package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

func (Names) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("name").
			Unique(),
		// We would like to set `Required()` here, but for implementation
		// reasons, we keep it short of making it required
	}
}
