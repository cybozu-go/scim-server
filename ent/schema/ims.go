package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

func (IMS) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("IMS").Unique(),
	}
}
