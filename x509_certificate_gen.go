package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/x509certificate"
	"github.com/cybozu-go/scim/resource"
)

func X509CertificateResourceFromEnt(in *ent.X509Certificate) (*resource.X509Certificate, error) {
	var b resource.Builder

	builder := b.X509Certificate()
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

func X509CertificateEntFieldFromSCIM(s string) string {
	switch s {
	case resource.X509CertificateDisplayKey:
		return x509certificate.FieldDisplay
	case resource.X509CertificatePrimaryKey:
		return x509certificate.FieldPrimary
	case resource.X509CertificateTypeKey:
		return x509certificate.FieldType
	case resource.X509CertificateValueKey:
		return x509certificate.FieldValue
	default:
		return s
	}
}
