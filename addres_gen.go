package server

import (
	"fmt"
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/address"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
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

type AddressPredicateBuilder struct {
	predicates []predicate.Address
}

func (b *AddressPredicateBuilder) Build(expr filter.Expr) ([]predicate.Address, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *AddressPredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *AddressPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.Address{address.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.Address{address.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *AddressPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
	lhe, err := exprAttr(expr.LHE())
	slhe, ok := lhe.(string)
	if err != nil || !ok {
		return fmt.Errorf("left hand side of CompareExpr is not valid")
	}

	rhe, err := exprAttr(expr.RHE())
	if err != nil {
		return fmt.Errorf("right hand side of CompareExpr is not valid: %w", err)
	}

	// convert rhe to string so it can be passed to regexp.QuoteMeta
	srhe := fmt.Sprintf("%v", rhe)

	switch expr.Operator() {
	case filter.EqualOp:
		switch slhe {
		case resource.AddressCountryKey:
			b.predicates = append(b.predicates, address.Country(srhe))
		case resource.AddressFormattedKey:
			b.predicates = append(b.predicates, address.Formatted(srhe))
		case resource.AddressLocalityKey:
			b.predicates = append(b.predicates, address.Locality(srhe))
		case resource.AddressPostalCodeKey:
			b.predicates = append(b.predicates, address.PostalCode(srhe))
		case resource.AddressRegionKey:
			b.predicates = append(b.predicates, address.Region(srhe))
		case resource.AddressStreetAddressKey:
			b.predicates = append(b.predicates, address.StreetAddress(srhe))
		default:
			return fmt.Errorf("invalid field name for Address: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
