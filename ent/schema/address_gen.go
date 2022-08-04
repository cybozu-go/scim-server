package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Address struct {
	ent.Schema
}

func (Address) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("country").Optional(),
		field.String("formatted").Optional(),
		field.String("locality").Optional(),
		field.String("postalCode").Optional(),
		field.String("region").Optional(),
		field.String("streetAddress").Optional(),
	}
}
