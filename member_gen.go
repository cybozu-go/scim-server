package server

import (
	"fmt"
	"reflect"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/member"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
	"github.com/cybozu-go/scim/resource"
)

func GroupMemberResourceFromEnt(in *ent.Member) (*resource.GroupMember, error) {
	var b resource.Builder

	builder := b.GroupMember()
	if !reflect.ValueOf(in.Ref).IsZero() {
		builder.Ref(in.Ref)
	}
	if !reflect.ValueOf(in.Type).IsZero() {
		builder.Type(in.Type)
	}
	if !reflect.ValueOf(in.Value).IsZero() {
		builder.Value(in.Value)
	}
	return builder.Build()
}

func GroupMemberEntFieldFromSCIM(s string) string {
	switch s {
	case resource.GroupMemberRefKey:
		return member.FieldRef
	case resource.GroupMemberTypeKey:
		return member.FieldType
	case resource.GroupMemberValueKey:
		return member.FieldValue
	default:
		return s
	}
}

type MemberPredicateBuilder struct {
	predicates []predicate.Member
}

func (b *MemberPredicateBuilder) Build(expr filter.Expr) ([]predicate.Member, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *MemberPredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *MemberPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.Member{member.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.Member{member.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *MemberPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
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
		case resource.GroupMemberRefKey:
			b.predicates = append(b.predicates, member.Ref(srhe))
		case resource.GroupMemberTypeKey:
			b.predicates = append(b.predicates, member.Type(srhe))
		case resource.GroupMemberValueKey:
			b.predicates = append(b.predicates, member.Value(srhe))
		default:
			return fmt.Errorf("invalid field name for GroupMember: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
