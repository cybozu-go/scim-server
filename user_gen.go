package server

import (
	"context"
	"fmt"
	"hash"
	"reflect"
	"sort"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim-server/ent/phonenumber"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/role"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim/resource"
)

func userLoadEntFields(q *ent.UserQuery, scimFields, excludedFields []string) {
	fields := make(map[string]struct{})
	if len(scimFields) == 0 {
		scimFields = []string{resource.UserActiveKey, resource.UserAddressesKey, resource.UserDisplayNameKey, resource.UserEmailsKey, resource.UserEntitlementsKey, resource.UserExternalIDKey, resource.UserGroupsKey, resource.UserIDKey, resource.UserIMSKey, resource.UserLocaleKey, resource.UserNameKey, resource.UserNickNameKey, resource.UserPasswordKey, resource.UserPhoneNumbersKey, resource.UserPhotosKey, resource.UserPreferredLanguageKey, resource.UserProfileURLKey, resource.UserRolesKey, resource.UserTimezoneKey, resource.UserTitleKey, resource.UserUserNameKey, resource.UserUserTypeKey, resource.UserX509CertificatesKey}
	}

	for _, name := range scimFields {
		fields[name] = struct{}{}
	}

	for _, name := range excludedFields {
		delete(fields, name)
	}
	selectNames := make([]string, 0, len(fields))
	for f := range fields {
		switch f {
		case resource.UserActiveKey:
			selectNames = append(selectNames, user.FieldActive)
		case resource.UserAddressesKey:
		case resource.UserDisplayNameKey:
			selectNames = append(selectNames, user.FieldDisplayName)
		case resource.UserEmailsKey:
			q.WithEmails()
		case resource.UserEntitlementsKey:
		case resource.UserExternalIDKey:
			selectNames = append(selectNames, user.FieldExternalID)
		case resource.UserGroupsKey:
		case resource.UserIDKey:
			selectNames = append(selectNames, user.FieldID)
		case resource.UserIMSKey:
		case resource.UserLocaleKey:
			selectNames = append(selectNames, user.FieldLocale)
		case resource.UserMetaKey:
		case resource.UserNameKey:
			q.WithName()
		case resource.UserNickNameKey:
			selectNames = append(selectNames, user.FieldNickName)
		case resource.UserPasswordKey:
			selectNames = append(selectNames, user.FieldPassword)
		case resource.UserPhoneNumbersKey:
		case resource.UserPhotosKey:
		case resource.UserPreferredLanguageKey:
			selectNames = append(selectNames, user.FieldPreferredLanguage)
		case resource.UserProfileURLKey:
			selectNames = append(selectNames, user.FieldProfileURL)
		case resource.UserRolesKey:
		case resource.UserTimezoneKey:
			selectNames = append(selectNames, user.FieldTimezone)
		case resource.UserTitleKey:
			selectNames = append(selectNames, user.FieldTitle)
		case resource.UserUserNameKey:
			selectNames = append(selectNames, user.FieldUserName)
		case resource.UserUserTypeKey:
			selectNames = append(selectNames, user.FieldUserType)
		case resource.UserX509CertificatesKey:
		}
	}
	selectNames = append(selectNames, user.FieldEtag)
	q.Select(selectNames...)
}

func userLocation(id string) string {
	return "https://foobar.com/scim/v2/Users/" + id
}

func UserResourceFromEnt(in *ent.User) (*resource.User, error) {
	var b resource.Builder

	builder := b.User()

	meta, err := b.Meta().
		ResourceType("User").
		Location(userLocation(in.ID.String())).
		Version(in.Etag).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build meta information for User")
	}

	builder.
		Meta(meta)

	if el := len(in.Edges.Emails); el > 0 {
		list := make([]*resource.Email, 0, el)
		for _, ine := range in.Edges.Emails {
			r, err := EmailResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build emails information for User")
			}
			list = append(list, r)
		}
		builder.Emails(list...)
	}

	if el := len(in.Edges.Name); el > 0 {
		list := make([]*resource.Names, 0, el)
		for _, ine := range in.Edges.Name {
			r, err := NamesResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build name information for User")
			}
			list = append(list, r)
		}
		builder.Name(list[0])
	}
	if !reflect.ValueOf(in.Active).IsZero() {
		builder.Active(in.Active)
	}
	if !reflect.ValueOf(in.DisplayName).IsZero() {
		builder.DisplayName(in.DisplayName)
	}
	if !reflect.ValueOf(in.ExternalID).IsZero() {
		builder.ExternalID(in.ExternalID)
	}
	builder.ID(in.ID.String())
	if !reflect.ValueOf(in.Locale).IsZero() {
		builder.Locale(in.Locale)
	}
	if !reflect.ValueOf(in.PreferredLanguage).IsZero() {
		builder.PreferredLanguage(in.PreferredLanguage)
	}
	if !reflect.ValueOf(in.Timezone).IsZero() {
		builder.Timezone(in.Timezone)
	}
	if !reflect.ValueOf(in.UserName).IsZero() {
		builder.UserName(in.UserName)
	}
	if !reflect.ValueOf(in.UserType).IsZero() {
		builder.UserType(in.UserType)
	}
	return builder.Build()
}

func UserEntFieldFromSCIM(s string) string {
	switch s {
	case resource.UserActiveKey:
		return user.FieldActive
	case resource.UserDisplayNameKey:
		return user.FieldDisplayName
	case resource.UserExternalIDKey:
		return user.FieldExternalID
	case resource.UserIDKey:
		return user.FieldID
	case resource.UserLocaleKey:
		return user.FieldLocale
	case resource.UserNickNameKey:
		return user.FieldNickName
	case resource.UserPasswordKey:
		return user.FieldPassword
	case resource.UserPreferredLanguageKey:
		return user.FieldPreferredLanguage
	case resource.UserProfileURLKey:
		return user.FieldProfileURL
	case resource.UserTimezoneKey:
		return user.FieldTimezone
	case resource.UserTitleKey:
		return user.FieldTitle
	case resource.UserUserNameKey:
		return user.FieldUserName
	case resource.UserUserTypeKey:
		return user.FieldUserType
	default:
		return s
	}
}

func userStartsWithPredicate(q *ent.UserQuery, scimField string, val interface{}) (predicate.User, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.UserDisplayNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserEmailsKey:
		switch subfield {
		case resource.EmailDisplayKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.DisplayHasPrefix(val.(string))), nil
		case resource.EmailTypeKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.TypeHasPrefix(val.(string))), nil
		case resource.EmailValueKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.ValueHasPrefix(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserExternalIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserLocaleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserNickNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPasswordKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPhoneNumbersKey:
		switch subfield {
		case resource.PhoneNumberDisplayKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.DisplayHasPrefix(val.(string))), nil
		case resource.PhoneNumberTypeKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.TypeHasPrefix(val.(string))), nil
		case resource.PhoneNumberValueKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.ValueHasPrefix(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserPreferredLanguageKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserProfileURLKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserRolesKey:
		switch subfield {
		case resource.RoleDisplayKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.DisplayHasPrefix(val.(string))), nil
		case resource.RoleTypeKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.TypeHasPrefix(val.(string))), nil
		case resource.RoleValueKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.ValueHasPrefix(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserTimezoneKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserTitleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserTypeKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func userEndsWithPredicate(q *ent.UserQuery, scimField string, val interface{}) (predicate.User, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.UserDisplayNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserEmailsKey:
		switch subfield {
		case resource.EmailDisplayKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.DisplayHasSuffix(val.(string))), nil
		case resource.EmailTypeKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.TypeHasSuffix(val.(string))), nil
		case resource.EmailValueKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.ValueHasSuffix(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserExternalIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserLocaleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserNickNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPasswordKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPhoneNumbersKey:
		switch subfield {
		case resource.PhoneNumberDisplayKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.DisplayHasSuffix(val.(string))), nil
		case resource.PhoneNumberTypeKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.TypeHasSuffix(val.(string))), nil
		case resource.PhoneNumberValueKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.ValueHasSuffix(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserPreferredLanguageKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserProfileURLKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserRolesKey:
		switch subfield {
		case resource.RoleDisplayKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.DisplayHasSuffix(val.(string))), nil
		case resource.RoleTypeKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.TypeHasSuffix(val.(string))), nil
		case resource.RoleValueKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.ValueHasSuffix(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserTimezoneKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserTitleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserTypeKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func userContainsPredicate(q *ent.UserQuery, scimField string, val interface{}) (predicate.User, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.UserDisplayNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserEmailsKey:
		switch subfield {
		case resource.EmailDisplayKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.DisplayContains(val.(string))), nil
		case resource.EmailTypeKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.TypeContains(val.(string))), nil
		case resource.EmailValueKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.ValueContains(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserExternalIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserLocaleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserNickNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPasswordKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPhoneNumbersKey:
		switch subfield {
		case resource.PhoneNumberDisplayKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.DisplayContains(val.(string))), nil
		case resource.PhoneNumberTypeKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.TypeContains(val.(string))), nil
		case resource.PhoneNumberValueKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.ValueContains(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserPreferredLanguageKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserProfileURLKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserRolesKey:
		switch subfield {
		case resource.RoleDisplayKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.DisplayContains(val.(string))), nil
		case resource.RoleTypeKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.TypeContains(val.(string))), nil
		case resource.RoleValueKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.ValueContains(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserTimezoneKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserTitleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserTypeKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func userEqualsPredicate(q *ent.UserQuery, scimField string, val interface{}) (predicate.User, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.UserDisplayNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserEmailsKey:
		switch subfield {
		case resource.EmailDisplayKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.DisplayEQ(val.(string))), nil
		case resource.EmailPrimaryKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.PrimaryEQ(val.(bool))), nil
		case resource.EmailTypeKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.TypeEQ(val.(string))), nil
		case resource.EmailValueKey:
			//nolint:forcetypeassert
			return user.HasEmailsWith(email.ValueEQ(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserExternalIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserIDKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserLocaleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserNickNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPasswordKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserPhoneNumbersKey:
		switch subfield {
		case resource.PhoneNumberDisplayKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.DisplayEQ(val.(string))), nil
		case resource.PhoneNumberPrimaryKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.PrimaryEQ(val.(bool))), nil
		case resource.PhoneNumberTypeKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.TypeEQ(val.(string))), nil
		case resource.PhoneNumberValueKey:
			//nolint:forcetypeassert
			return user.HasPhoneNumbersWith(phonenumber.ValueEQ(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserPreferredLanguageKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserProfileURLKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserRolesKey:
		switch subfield {
		case resource.RoleDisplayKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.DisplayEQ(val.(string))), nil
		case resource.RolePrimaryKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.PrimaryEQ(val.(bool))), nil
		case resource.RoleTypeKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.TypeEQ(val.(string))), nil
		case resource.RoleValueKey:
			//nolint:forcetypeassert
			return user.HasRolesWith(role.ValueEQ(val.(string))), nil
		default:
			return nil, fmt.Errorf("invalid filter specification: invalid subfield for %q", field)
		}
	case resource.UserTimezoneKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserTitleKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserNameKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.UserUserTypeKey:
		entFieldName := UserEntFieldFromSCIM(scimField)
		return predicate.User(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func userPresencePredicate(scimField string) predicate.User {
	switch scimField {
	case resource.UserDisplayNameKey:
		return user.And(user.DisplayNameNotNil(), user.DisplayNameNEQ(""))
	case resource.UserExternalIDKey:
		return user.And(user.ExternalIDNotNil(), user.ExternalIDNEQ(""))
	case resource.UserLocaleKey:
		return user.And(user.LocaleNotNil(), user.LocaleNEQ(""))
	case resource.UserNickNameKey:
		return user.And(user.NickNameNotNil(), user.NickNameNEQ(""))
	case resource.UserPasswordKey:
		return user.And(user.PasswordNotNil(), user.PasswordNEQ(""))
	case resource.UserPreferredLanguageKey:
		return user.And(user.PreferredLanguageNotNil(), user.PreferredLanguageNEQ(""))
	case resource.UserProfileURLKey:
		return user.And(user.ProfileURLNotNil(), user.ProfileURLNEQ(""))
	case resource.UserTimezoneKey:
		return user.And(user.TimezoneNotNil(), user.TimezoneNEQ(""))
	case resource.UserTitleKey:
		return user.And(user.TitleNotNil(), user.TitleNEQ(""))
	case resource.UserUserTypeKey:
		return user.And(user.UserTypeNotNil(), user.UserTypeNEQ(""))
	default:
		return nil
	}
}

func (b *Backend) createEmails(in *resource.User, h hash.Hash) ([]*ent.Email, error) {
	list := make([]*ent.Email, len(in.Emails()))
	inbound := in.Emails()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.Email.Create()
		createCall.SetValue(v.Value())
		fmt.Fprint(h, v.Value())

		if v.HasDisplay() {
			createCall.SetDisplay(v.Display())
			fmt.Fprint(h, v.Display())
		}

		if v.HasType() {
			createCall.SetType(v.Type())
			fmt.Fprint(h, v.Type())
		}

		if sv := v.Primary(); sv {
			if hasPrimary {
				return nil, fmt.Errorf("invalid user.emails: multiple emails have been set to primary")
			}
			createCall.SetPrimary(true)
			fmt.Fprint(h, []byte{1})
			hasPrimary = true
		} else {
			fmt.Fprint(h, []byte{0})
		}

		r, err := createCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to save email %d: %w", i, err)
		}

		list[i] = r
	}
	return list, nil
}

func (b *Backend) createEntitlements(in *resource.User, h hash.Hash) ([]*ent.Entitlement, error) {
	list := make([]*ent.Entitlement, len(in.Entitlements()))
	inbound := in.Entitlements()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.Entitlement.Create()
		createCall.SetValue(v.Value())
		fmt.Fprint(h, v.Value())

		if v.HasDisplay() {
			createCall.SetDisplay(v.Display())
			fmt.Fprint(h, v.Display())
		}

		if v.HasType() {
			createCall.SetType(v.Type())
			fmt.Fprint(h, v.Type())
		}

		if sv := v.Primary(); sv {
			if hasPrimary {
				return nil, fmt.Errorf("invalid user.entitlements: multiple entitlements have been set to primary")
			}
			createCall.SetPrimary(true)
			fmt.Fprint(h, []byte{1})
			hasPrimary = true
		} else {
			fmt.Fprint(h, []byte{0})
		}

		r, err := createCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to save email %d: %w", i, err)
		}

		list[i] = r
	}
	return list, nil
}

func (b *Backend) createIMS(in *resource.User, h hash.Hash) ([]*ent.IMS, error) {
	list := make([]*ent.IMS, len(in.IMS()))
	inbound := in.IMS()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.IMS.Create()
		createCall.SetValue(v.Value())
		fmt.Fprint(h, v.Value())

		if v.HasDisplay() {
			createCall.SetDisplay(v.Display())
			fmt.Fprint(h, v.Display())
		}

		if v.HasType() {
			createCall.SetType(v.Type())
			fmt.Fprint(h, v.Type())
		}

		if sv := v.Primary(); sv {
			if hasPrimary {
				return nil, fmt.Errorf("invalid user.ims: multiple ims have been set to primary")
			}
			createCall.SetPrimary(true)
			fmt.Fprint(h, []byte{1})
			hasPrimary = true
		} else {
			fmt.Fprint(h, []byte{0})
		}

		r, err := createCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to save email %d: %w", i, err)
		}

		list[i] = r
	}
	return list, nil
}

func (b *Backend) createPhoneNumbers(in *resource.User, h hash.Hash) ([]*ent.PhoneNumber, error) {
	list := make([]*ent.PhoneNumber, len(in.PhoneNumbers()))
	inbound := in.PhoneNumbers()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.PhoneNumber.Create()
		createCall.SetValue(v.Value())
		fmt.Fprint(h, v.Value())

		if v.HasDisplay() {
			createCall.SetDisplay(v.Display())
			fmt.Fprint(h, v.Display())
		}

		if v.HasType() {
			createCall.SetType(v.Type())
			fmt.Fprint(h, v.Type())
		}

		if sv := v.Primary(); sv {
			if hasPrimary {
				return nil, fmt.Errorf("invalid user.phoneNumbers: multiple phoneNumbers have been set to primary")
			}
			createCall.SetPrimary(true)
			fmt.Fprint(h, []byte{1})
			hasPrimary = true
		} else {
			fmt.Fprint(h, []byte{0})
		}

		r, err := createCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to save email %d: %w", i, err)
		}

		list[i] = r
	}
	return list, nil
}

func (b *Backend) createPhotos(in *resource.User, h hash.Hash) ([]*ent.Photo, error) {
	list := make([]*ent.Photo, len(in.Photos()))
	inbound := in.Photos()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.Photo.Create()
		createCall.SetValue(v.Value())
		fmt.Fprint(h, v.Value())

		if v.HasDisplay() {
			createCall.SetDisplay(v.Display())
			fmt.Fprint(h, v.Display())
		}

		if v.HasType() {
			createCall.SetType(v.Type())
			fmt.Fprint(h, v.Type())
		}

		if sv := v.Primary(); sv {
			if hasPrimary {
				return nil, fmt.Errorf("invalid user.photos: multiple photos have been set to primary")
			}
			createCall.SetPrimary(true)
			fmt.Fprint(h, []byte{1})
			hasPrimary = true
		} else {
			fmt.Fprint(h, []byte{0})
		}

		r, err := createCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to save email %d: %w", i, err)
		}

		list[i] = r
	}
	return list, nil
}

func (b *Backend) createRoles(in *resource.User, h hash.Hash) ([]*ent.Role, error) {
	list := make([]*ent.Role, len(in.Roles()))
	inbound := in.Roles()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.Role.Create()
		createCall.SetValue(v.Value())
		fmt.Fprint(h, v.Value())

		if v.HasDisplay() {
			createCall.SetDisplay(v.Display())
			fmt.Fprint(h, v.Display())
		}

		if v.HasType() {
			createCall.SetType(v.Type())
			fmt.Fprint(h, v.Type())
		}

		if sv := v.Primary(); sv {
			if hasPrimary {
				return nil, fmt.Errorf("invalid user.roles: multiple roles have been set to primary")
			}
			createCall.SetPrimary(true)
			fmt.Fprint(h, []byte{1})
			hasPrimary = true
		} else {
			fmt.Fprint(h, []byte{0})
		}

		r, err := createCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to save email %d: %w", i, err)
		}

		list[i] = r
	}
	return list, nil
}
