package server

import (
	"fmt"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim/filter"
)

type filterVisitor struct {
	uq     *ent.UserQuery
	gq     *ent.GroupQuery
	users  []predicate.User
	groups []predicate.Group
}

// visit visits the filter expression AST and collects ent-predicates.
//
// during the traversal, we build predicates for multiple resources
// at the same time... this makes things more complicated, but the
// thing is, the search endpoint allows users to query both User and Group
// at the same time -- for example:
//
//   displayName eq "foo"
//
// There's nothing stopping the user from mixing this with resource-specific queries,
// (although that would affect the results): for example:
//
//   displayName eq "foo" or roles.value eq "bar"
//
// This would presumably match groups with displayName foo, and users wit displayName
// foo or with roles.value of bar.
//
// This is why we build predicates for multiple resources at the same time.
//
// However, in order to re-use the same mechanic in resource-specific search endpoints
// such as /Users/.search and /Group/.search, we also make it possible to "toggle"
// building the predicates by checking if the accumulator (v.users and v.groups) are
// nil or not.
//
// If the accumulator is nil, we assume that the caller is not interested in accumulating
// predicates for that specific resource. For example, if v.users = nil and v.groups is NOT nil,
// then only the predicates for the Group resources are accumulated
func (v *filterVisitor) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.PresenceExpr:
		return v.visitPresenceExpr(expr)
	case filter.CompareExpr:
		return v.visitCompareExpr(expr)
	case filter.RegexExpr:
		return v.visitRegexExpr(expr)
	case filter.LogExpr: // RENAME ME TO LogicalStatement
		return v.visitLogExpr(expr)
	case filter.ParenExpr:
		return v.visitParenExpr(expr)
	case filter.ValuePath:
		return v.visitValuePath(expr)
	default:
		return fmt.Errorf(`unhandled statement type: %T`, expr)
	}
}

func exprAttr(expr interface{}) (interface{}, error) {
	switch v := expr.(type) {
	case string:
		return v, nil
	case interface{ Lit() string }: // IdentifierExpr, AttrValueExpr
		return v.Lit(), nil
	case filter.BoolExpr:
		return v.Lit(), nil
	case filter.NumberExpr:
		return v.Lit(), nil
	default:
		return nil, fmt.Errorf(`unhandled type: %T`, v)
	}
}

func (v *filterVisitor) visitPresenceExpr(expr filter.PresenceExpr) error {
	attr, err := exprAttr(expr.Attr())
	sattr, ok := attr.(string)
	if err != nil || !ok {
		if err == nil && !ok {
			err = fmt.Errorf(`expected string, got %T`, attr)
		}
		return fmt.Errorf(`left hand side of PresenceExpr is not valid: %w`, err)
	}

	switch expr.Operator() {
	case filter.PresenceOp:
		if v.users != nil {
			if pred := userPresencePredicate(sattr); pred != nil {
				v.users = append(v.users, pred)
			}
		}
		return nil
	default:
		return fmt.Errorf(`unhandled attr operator %q`, expr.Operator())
	}
}

func (v *filterVisitor) visitRegexExpr(expr filter.RegexExpr) error {
	lhe, err := exprAttr(expr.LHE())
	slhe, ok := lhe.(string)
	if err != nil || !ok {
		return fmt.Errorf(`left hand side of RegexExpr is not valid`)
	}

	rhe, err := exprAttr(expr.Value())
	if err != nil {
		return fmt.Errorf(`right hand side of RegexExpr is not valid: %w`, err)
	}
	// convert rhe to string so it can be passed to regexp.QuoteMeta
	srhe := fmt.Sprintf(`%v`, rhe)

	switch expr.Operator() {
	case filter.ContainsOp:
		if v.users != nil {
			pred, err := userContainsPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupContainsPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	case filter.StartsWithOp:
		if v.users != nil {
			pred, err := userStartsWithPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupStartsWithPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	case filter.EndsWithOp:
		if v.users != nil {
			pred, err := userEndsWithPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupEndsWithPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	default:
		return fmt.Errorf(`unhandled regexp operator %q`, expr.Operator())
	}
}

func (v *filterVisitor) visitCompareExpr(expr filter.CompareExpr) error {
	lhe, err := exprAttr(expr.LHE())
	slhe, ok := lhe.(string)
	if err != nil || !ok {
		return fmt.Errorf(`left hand side of CompareExpr is not valid`)
	}

	rhe, err := exprAttr(expr.RHE())
	if err != nil {
		return fmt.Errorf(`right hand side of CompareExpr is not valid: %w`, err)
	}
	// convert rhe to string so it can be passed to regexp.QuoteMeta
	srhe := fmt.Sprintf(`%v`, rhe)

	switch expr.Operator() {
	case filter.EqualOp:
		if v.users != nil {
			pred, err := userEqualsPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupEqualsPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	default:
		panic(expr.Operator())
	}
}

func (v *filterVisitor) visitLogExpr(expr filter.LogExpr) error {
	if err := v.visit(expr.LHE()); err != nil {
		return fmt.Errorf(`failed to parse left hand side of %q statement: %w`, expr.Operator(), err)
	}
	if err := v.visit(expr.RHS()); err != nil {
		return fmt.Errorf(`failed to parse right hand side of %q statement: %w`, expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		if v.users != nil {
			v.users = []predicate.User{user.And(v.users...)}
		}
		if v.groups != nil {
			v.groups = []predicate.Group{group.And(v.groups...)}
		}
	case "or":
		if v.users != nil {
			v.users = []predicate.User{user.Or(v.users...)}
		}
		if v.groups != nil {
			v.groups = []predicate.Group{group.Or(v.groups...)}
		}
	default:
		return fmt.Errorf(`unhandled logical statement operator %q`, expr.Operator())
	}
	return nil
}

func (v *filterVisitor) visitParenExpr(expr filter.Expr) error {
	return fmt.Errorf(`unimplemented`)
}

func (v *filterVisitor) visitValuePath(expr filter.ValuePath) error {
	return fmt.Errorf(`unimplemented`)
}
