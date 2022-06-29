package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/role"
	"github.com/cybozu-go/scim/resource"
)

func RoleResourceFromEnt(in *ent.Role) (*resource.Role, error) {
	var b resource.Builder

	builder := b.Role()
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

func RoleEntFieldFromSCIM(s string) string {
	switch s {
	case resource.RoleDisplayKey:
		return role.FieldDisplay
	case resource.RolePrimaryKey:
		return role.FieldPrimary
	case resource.RoleTypeKey:
		return role.FieldType
	case resource.RoleValueKey:
		return role.FieldValue
	default:
		return s
	}
}
