package server

import (
	"fmt"
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/names"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
	"github.com/cybozu-go/scim/resource"
)

func NamesResourceFromEnt(in *ent.Names) (*resource.Names, error) {
	var b resource.Builder

	builder := b.Names()
	if !reflect.ValueOf(in.FamilyName).IsZero() {
		builder.FamilyName(in.FamilyName)
	}
	if !reflect.ValueOf(in.Formatted).IsZero() {
		builder.Formatted(in.Formatted)
	}
	if !reflect.ValueOf(in.GivenName).IsZero() {
		builder.GivenName(in.GivenName)
	}
	if !reflect.ValueOf(in.HonorificPrefix).IsZero() {
		builder.HonorificPrefix(in.HonorificPrefix)
	}
	if !reflect.ValueOf(in.HonorificSuffix).IsZero() {
		builder.HonorificSuffix(in.HonorificSuffix)
	}
	if !reflect.ValueOf(in.MiddleName).IsZero() {
		builder.MiddleName(in.MiddleName)
	}
	return builder.Build()
}

func NamesEntFieldFromSCIM(s string) string {
	switch s {
	case resource.NamesFamilyNameKey:
		return names.FieldFamilyName
	case resource.NamesFormattedKey:
		return names.FieldFormatted
	case resource.NamesGivenNameKey:
		return names.FieldGivenName
	case resource.NamesHonorificPrefixKey:
		return names.FieldHonorificPrefix
	case resource.NamesHonorificSuffixKey:
		return names.FieldHonorificSuffix
	case resource.NamesMiddleNameKey:
		return names.FieldMiddleName
	default:
		return s
	}
}

type NamesPredicateBuilder struct {
	predicates []predicate.Names
}

func (b *NamesPredicateBuilder) Build(expr filter.Expr) ([]predicate.Names, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *NamesPredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *NamesPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.Names{names.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.Names{names.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *NamesPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
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
		case resource.NamesFamilyNameKey:
			b.predicates = append(b.predicates, names.FamilyName(srhe))
		case resource.NamesFormattedKey:
			b.predicates = append(b.predicates, names.Formatted(srhe))
		case resource.NamesGivenNameKey:
			b.predicates = append(b.predicates, names.GivenName(srhe))
		case resource.NamesHonorificPrefixKey:
			b.predicates = append(b.predicates, names.HonorificPrefix(srhe))
		case resource.NamesHonorificSuffixKey:
			b.predicates = append(b.predicates, names.HonorificSuffix(srhe))
		case resource.NamesMiddleNameKey:
			b.predicates = append(b.predicates, names.MiddleName(srhe))
		default:
			return fmt.Errorf("invalid field name for Names: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
