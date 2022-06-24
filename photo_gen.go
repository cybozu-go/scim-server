package server

import (
	"fmt"
	"reflect"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/photo"
	"github.com/cybozu-go/scim-server/ent/predicate"
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

func photoStartsWithPredicate(q *ent.PhotoQuery, scimField string, val interface{}) (predicate.Photo, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhotoDisplayKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoTypeKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoValueKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func photoEndsWithPredicate(q *ent.PhotoQuery, scimField string, val interface{}) (predicate.Photo, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhotoDisplayKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoTypeKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoValueKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func photoContainsPredicate(q *ent.PhotoQuery, scimField string, val interface{}) (predicate.Photo, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhotoDisplayKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoTypeKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoValueKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func photoEqualsPredicate(q *ent.PhotoQuery, scimField string, val interface{}) (predicate.Photo, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.PhotoDisplayKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoTypeKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.PhotoValueKey:
		entFieldName := PhotoEntFieldFromSCIM(scimField)
		return predicate.Photo(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func photoPresencePredicate(scimField string) predicate.Photo {
	switch scimField {
	case resource.PhotoDisplayKey:
		return photo.And(photo.DisplayNotNil(), photo.DisplayNEQ(""))
	case resource.PhotoTypeKey:
		return photo.And(photo.TypeNotNil(), photo.TypeNEQ(""))
	case resource.PhotoValueKey:
		return photo.And(photo.ValueNotNil(), photo.ValueNEQ(""))
	default:
		return nil
	}
}
