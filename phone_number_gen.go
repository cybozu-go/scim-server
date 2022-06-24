package server

import (
	"fmt"
	"reflect"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/phonenumber"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/resource"
)

func PhoneNumberResourceFromEnt(in *ent.PhoneNumber) (*resource.PhoneNumber, error) {
	var b resource.Builder

	builder := b.PhoneNumber()
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

func PhoneNumberEntFieldFromSCIM(s string) string {
	switch s {
	case resource.PhoneNumberDisplayKey:
		return phonenumber.FieldDisplay
	case resource.PhoneNumberPrimaryKey:
		return phonenumber.FieldPrimary
	case resource.PhoneNumberTypeKey:
		return phonenumber.FieldType
	case resource.PhoneNumberValueKey:
		return phonenumber.FieldValue
	default:
		return s
	}
}

func phoneNumberStartsWithPredicate(q *ent.PhoneNumberQuery, scimField string, val interface{}) (predicate.PhoneNumber, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhoneNumberDisplayKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberTypeKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberValueKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func phoneNumberEndsWithPredicate(q *ent.PhoneNumberQuery, scimField string, val interface{}) (predicate.PhoneNumber, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhoneNumberDisplayKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberTypeKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberValueKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func phoneNumberContainsPredicate(q *ent.PhoneNumberQuery, scimField string, val interface{}) (predicate.PhoneNumber, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhoneNumberDisplayKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberTypeKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberValueKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func phoneNumberEqualsPredicate(q *ent.PhoneNumberQuery, scimField string, val interface{}) (predicate.PhoneNumber, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhoneNumberDisplayKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberTypeKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhoneNumberValueKey:
		entFieldName := PhoneNumberEntFieldFromSCIM(scimField)
		return predicate.PhoneNumber(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func phoneNumberPresencePredicate(scimField string) predicate.PhoneNumber {
	switch scimField {
	case resource.PhoneNumberDisplayKey:
		return phonenumber.And(phonenumber.DisplayNotNil(), phonenumber.DisplayNEQ(""))
	case resource.PhoneNumberTypeKey:
		return phonenumber.And(phonenumber.TypeNotNil(), phonenumber.TypeNEQ(""))
	case resource.PhoneNumberValueKey:
		return phonenumber.And(phonenumber.ValueNotNil(), phonenumber.ValueNEQ(""))
	default:
		return nil
	}
}
