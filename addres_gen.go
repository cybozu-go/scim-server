package server

import (
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/address"
	"github.com/cybozu-go/scim/resource"
)

func AddressResourceFromEnt(in *ent.Address) (*resource.Address, error) {
	var b resource.Builder

	builder := b.Address()
	if !reflect.ValueOf(in.Country).IsZero() {
		builder.Country(in.Country)
	}
	if !reflect.ValueOf(in.Formatted).IsZero() {
		builder.Formatted(in.Formatted)
	}
	if !reflect.ValueOf(in.Locality).IsZero() {
		builder.Locality(in.Locality)
	}
	if !reflect.ValueOf(in.PostalCode).IsZero() {
		builder.PostalCode(in.PostalCode)
	}
	if !reflect.ValueOf(in.Region).IsZero() {
		builder.Region(in.Region)
	}
	if !reflect.ValueOf(in.StreetAddress).IsZero() {
		builder.StreetAddress(in.StreetAddress)
	}
	return builder.Build()
}

func AddressEntFieldFromSCIM(s string) string {
	switch s {
	case resource.AddressCountryKey:
		return address.FieldCountry
	case resource.AddressFormattedKey:
		return address.FieldFormatted
	case resource.AddressLocalityKey:
		return address.FieldLocality
	case resource.AddressPostalCodeKey:
		return address.FieldPostalCode
	case resource.AddressRegionKey:
		return address.FieldRegion
	case resource.AddressStreetAddressKey:
		return address.FieldStreetAddress
	default:
		return s
	}
}
