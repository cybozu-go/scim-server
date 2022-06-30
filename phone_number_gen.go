package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/phone_number"
	"github.com/cybozu-go/scim/resource"
)

func PhoneNumberResourceFromEnt(in *ent.PhoneNumber) (*resource.PhoneNumber, error) {
	var b resource.Builder

	builder := b.PhoneNumber()
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

func PhoneNumberEntFieldFromSCIM(s string) string {
	switch s {
	case resource.PhoneNumberDisplayKey:
		return phone_number.FieldDisplay
	case resource.PhoneNumberPrimaryKey:
		return phone_number.FieldPrimary
	case resource.PhoneNumberTypeKey:
		return phone_number.FieldType
	case resource.PhoneNumberValueKey:
		return phone_number.FieldValue
	default:
		return s
	}
}
