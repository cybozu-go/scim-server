package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type GroupMember struct {
	ent.Schema
}

func (GroupMember) Fields() []ent.Field {
	return []ent.Field{
		field.String("value"),
		field.String(`type`).Optional(),
		field.String(`ref`).Optional(),
	}
}

func (GroupMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From(`user`, User.Type).Ref(`groups`).Unique(),
		edge.From(`group`, Group.Type).Ref(`members`).Unique(),
	}
}
