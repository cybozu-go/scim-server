package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim/resource"
)

func IMSResourceFromEnt(in *ent.IMS) (*resource.IMS, error) {
	var b resource.Builder

	builder := b.IMS()
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

func IMSEntFieldFromSCIM(s string) string {
	switch s {
	case resource.IMSDisplayKey:
		return ims.FieldDisplay
	case resource.IMSPrimaryKey:
		return ims.FieldPrimary
	case resource.IMSTypeKey:
		return ims.FieldType
	case resource.IMSValueKey:
		return ims.FieldValue
	default:
		return s
	}
}
