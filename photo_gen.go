package server

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/photo"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
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

type PhotoPredicateBuilder struct {
	predicates []predicate.Photo
}

func (b *PhotoPredicateBuilder) Build(expr filter.Expr) ([]predicate.Photo, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *PhotoPredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *PhotoPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.Photo{photo.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.Photo{photo.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *PhotoPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
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
		case resource.PhotoDisplayKey:
			b.predicates = append(b.predicates, photo.Display(srhe))
		case resource.PhotoPrimaryKey:
			v, err := strconv.ParseBool(srhe)
			if err != nil {
				return fmt.Errorf("failed to parse boolean expression")
			}
			b.predicates = append(b.predicates, photo.Primary(v))
		case resource.PhotoTypeKey:
			b.predicates = append(b.predicates, photo.Type(srhe))
		case resource.PhotoValueKey:
			b.predicates = append(b.predicates, photo.Value(srhe))
		default:
			return fmt.Errorf("invalid field name for Photo: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
