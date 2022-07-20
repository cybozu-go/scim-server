package server

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/x509certificate"
	"github.com/cybozu-go/scim/filter"
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

type X509CertificatePredicateBuilder struct {
	predicates []predicate.X509Certificate
}

func (b *X509CertificatePredicateBuilder) Build(expr filter.Expr) ([]predicate.X509Certificate, error) {
	b.predicates = nil
	if err := b.visit(expr); err != nil {
		return nil, err
	}
	return b.predicates, nil
}

func (b *X509CertificatePredicateBuilder) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.CompareExpr:
		return b.visitCompareExpr(expr)
	case filter.LogExpr:
		return b.visitLogExpr(expr)
	default:
		return fmt.Errorf("unhandled expression type %T", expr)
	}
}

func (b *X509CertificatePredicateBuilder) visitLogExpr(expr filter.LogExpr) error {
	if err := b.visit(expr.LHE()); err != nil {
		return fmt.Errorf("failed to parse left hand side of %q statement: %w", expr.Operator(), err)
	}
	if err := b.visit(expr.RHS()); err != nil {
		return fmt.Errorf("failed to parse right hand side of %q statement: %w", expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		b.predicates = []predicate.X509Certificate{x509certificate.And(b.predicates...)}
	case "or":
		b.predicates = []predicate.X509Certificate{x509certificate.Or(b.predicates...)}
	default:
		return fmt.Errorf("unhandled logical operator %q", expr.Operator())
	}
	return nil
}

func (b *X509CertificatePredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {
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
		case resource.X509CertificateDisplayKey:
			b.predicates = append(b.predicates, x509certificate.Display(srhe))
		case resource.X509CertificatePrimaryKey:
			v, err := strconv.ParseBool(srhe)
			if err != nil {
				return fmt.Errorf("failed to parse boolean expression")
			}
			b.predicates = append(b.predicates, x509certificate.Primary(v))
		case resource.X509CertificateTypeKey:
			b.predicates = append(b.predicates, x509certificate.Type(srhe))
		case resource.X509CertificateValueKey:
			b.predicates = append(b.predicates, x509certificate.Value(srhe))
		default:
			return fmt.Errorf("invalid field name for X509Certificate: %q", slhe)
		}
	default:
		return fmt.Errorf("invalid operator: %q", expr.Operator())
	}
	return nil
}
