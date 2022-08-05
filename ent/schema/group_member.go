package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Member struct {
	ent.Schema
}

func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.String(`value`),
		field.String(`display`).Optional(),
		field.String(`type`),
		field.String(`ref`).Optional(),
	}
}

func (Member) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From(`group`, Group.Type).Ref(`members`).Unique(),
	}
}
