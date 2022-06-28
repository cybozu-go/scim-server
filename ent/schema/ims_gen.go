package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type IMS struct {
	ent.Schema
}

func (IMS) Fields() []ent.Field {
	return []ent.Field{
		field.String("display").Optional(),
		field.Bool("primary").Optional(),
		field.String("type").Optional(),
		field.String("value").Optional(),
	}
}
