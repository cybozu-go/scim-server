package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", GroupMember.Type),
	}
}
