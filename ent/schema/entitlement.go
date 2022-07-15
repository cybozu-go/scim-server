package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

func (Entitlement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("entitlements").Unique(),
	}
}
