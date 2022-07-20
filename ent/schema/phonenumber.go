package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

func (PhoneNumber) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("phone_numbers").Unique(),
	}
}
