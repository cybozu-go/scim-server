package server

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
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

type IMSPredicateBuilder struct {
	predicates []predicate.IMS
}

func (b *IMSPredicateBuilder) Build(expr filter.Expr) ([]predicate.IMS, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *IMSPredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *IMSPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.IMS{ims.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.IMS{ims.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *IMSPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
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
		case resource.IMSDisplayKey:
			b.predicates = append(b.predicates, ims.Display(srhe))
		case resource.IMSPrimaryKey:
			v, err := strconv.ParseBool(srhe)
			if err != nil {
				return fmt.Errorf("failed to parse boolean expression")
			}
			b.predicates = append(b.predicates, ims.Primary(v))
		case resource.IMSTypeKey:
			b.predicates = append(b.predicates, ims.Type(srhe))
		case resource.IMSValueKey:
			b.predicates = append(b.predicates, ims.Value(srhe))
		default:
			return fmt.Errorf("invalid field name for IMS: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
