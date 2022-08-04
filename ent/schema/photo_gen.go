package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Photo struct {
	ent.Schema
}

func (Photo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("display").Optional(),
		field.Bool("primary").Optional(),
		field.String("type").Optional(),
		field.String("value").Optional(),
	}
}
