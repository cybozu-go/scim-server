package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim/resource"
)

func EmailResourceFromEnt(in *ent.Email) (*resource.Email, error) {
	var b resource.Builder

	builder := b.Email()
	if !reflect.ValueOf(in.Display).IsZero() {
		builder.Display(in.Display)
	}
	if !reflect.ValueOf(in.Primary).IsZero() {
		builder.Primary(in.Primary)
	}
	if !reflect.ValueOf(in.Type).IsZero() {
		builder.Type(in.Type)
	}
	if !reflect.ValueOf(in.Value).IsZero() {
		builder.Value(in.Value)
	}
	return builder.Build()
}

func EmailEntFieldFromSCIM(s string) string {
	switch s {
	case resource.EmailDisplayKey:
		return email.FieldDisplay
	case resource.EmailPrimaryKey:
		return email.FieldPrimary
	case resource.EmailTypeKey:
		return email.FieldType
	case resource.EmailValueKey:
		return email.FieldValue
	default:
		return s
	}
}
