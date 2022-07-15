package server

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/entitlement"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
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

type EntitlementPredicateBuilder struct {
	predicates []predicate.Entitlement
}

func (b *EntitlementPredicateBuilder) Build(expr filter.Expr) ([]predicate.Entitlement, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *EntitlementPredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *EntitlementPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.Entitlement{entitlement.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.Entitlement{entitlement.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *EntitlementPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
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
		case resource.EntitlementDisplayKey:
			b.predicates = append(b.predicates, entitlement.Display(srhe))
		case resource.EntitlementPrimaryKey:
			v, err := strconv.ParseBool(srhe)
			if err != nil {
				return fmt.Errorf("failed to parse boolean expression")
			}
			b.predicates = append(b.predicates, entitlement.Primary(v))
		case resource.EntitlementTypeKey:
			b.predicates = append(b.predicates, entitlement.Type(srhe))
		case resource.EntitlementValueKey:
			b.predicates = append(b.predicates, entitlement.Value(srhe))
		default:
			return fmt.Errorf("invalid field name for Entitlement: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
