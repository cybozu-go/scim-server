package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/entitlement"
	"github.com/cybozu-go/scim/resource"
)

func EntitlementResourceFromEnt(in *ent.Entitlement) (*resource.Entitlement, error) {
	var b resource.Builder

	builder := b.Entitlement()
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

func EntitlementEntFieldFromSCIM(s string) string {
	switch s {
	case resource.EntitlementDisplayKey:
		return entitlement.FieldDisplay
	case resource.EntitlementPrimaryKey:
		return entitlement.FieldPrimary
	case resource.EntitlementTypeKey:
		return entitlement.FieldType
	case resource.EntitlementValueKey:
		return entitlement.FieldValue
	default:
		return s
	}
}
