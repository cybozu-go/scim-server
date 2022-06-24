package server

import (
	"fmt"
	"reflect"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim-server/ent/predicate"
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

func imSStartsWithPredicate(q *ent.IMSQuery, scimField string, val interface{}) (predicate.IMS, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.IMSDisplayKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSTypeKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSValueKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func imSEndsWithPredicate(q *ent.IMSQuery, scimField string, val interface{}) (predicate.IMS, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.IMSDisplayKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSTypeKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSValueKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func imSContainsPredicate(q *ent.IMSQuery, scimField string, val interface{}) (predicate.IMS, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.IMSDisplayKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSTypeKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSValueKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func imSEqualsPredicate(q *ent.IMSQuery, scimField string, val interface{}) (predicate.IMS, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.IMSDisplayKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSTypeKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.IMSValueKey:
		entFieldName := IMSEntFieldFromSCIM(scimField)
		return predicate.IMS(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func imSPresencePredicate(scimField string) predicate.IMS {
	switch scimField {
	case resource.IMSDisplayKey:
		return ims.And(ims.DisplayNotNil(), ims.DisplayNEQ(""))
	case resource.IMSTypeKey:
		return ims.And(ims.TypeNotNil(), ims.TypeNEQ(""))
	case resource.IMSValueKey:
		return ims.And(ims.ValueNotNil(), ims.ValueNEQ(""))
	default:
		return nil
	}
}
