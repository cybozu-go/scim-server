package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/photo"
	"github.com/cybozu-go/scim/resource"
)

func PhotoResourceFromEnt(in *ent.Photo) (*resource.Photo, error) {
	var b resource.Builder

	builder := b.Photo()
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

func PhotoEntFieldFromSCIM(s string) string {
	switch s {
	case resource.PhotoDisplayKey:
		return photo.FieldDisplay
	case resource.PhotoPrimaryKey:
		return photo.FieldPrimary
	case resource.PhotoTypeKey:
		return photo.FieldType
	case resource.PhotoValueKey:
		return photo.FieldValue
	default:
		return s
	}
}
