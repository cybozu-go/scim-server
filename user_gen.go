package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/google/uuid"
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
			q.WithAddresses()
		case resource.UserDisplayNameKey:
			selectNames = append(selectNames, user.FieldDisplayName)
		case resource.UserEmailsKey:
			q.WithEmails()
		case resource.UserEntitlementsKey:
			q.WithEntitlements()
		case resource.UserExternalIDKey:
			selectNames = append(selectNames, user.FieldExternalID)
		case resource.UserGroupsKey:
			q.WithGroups()
		case resource.UserIDKey:
			selectNames = append(selectNames, user.FieldID)
		case resource.UserIMSKey:
			q.WithImses()
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
			q.WithPhoneNumbers()
		case resource.UserPhotosKey:
			q.WithPhotos()
		case resource.UserPreferredLanguageKey:
			selectNames = append(selectNames, user.FieldPreferredLanguage)
		case resource.UserProfileURLKey:
			selectNames = append(selectNames, user.FieldProfileURL)
		case resource.UserRolesKey:
			q.WithRoles()
		case resource.UserTimezoneKey:
			selectNames = append(selectNames, user.FieldTimezone)
		case resource.UserTitleKey:
			selectNames = append(selectNames, user.FieldTitle)
		case resource.UserUserNameKey:
			selectNames = append(selectNames, user.FieldUserName)
		case resource.UserUserTypeKey:
			selectNames = append(selectNames, user.FieldUserType)
		case resource.UserX509CertificatesKey:
			q.WithX509Certificates()
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

	if el := len(in.Edges.Addresses); el > 0 {
		list := make([]*resource.Address, 0, el)
		for _, ine := range in.Edges.Addresses {
			r, err := AddressResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build addresses information for User")
			}
			list = append(list, r)
		}
		builder.Addresses(list...)
	}

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

	if el := len(in.Edges.Entitlements); el > 0 {
		list := make([]*resource.Entitlement, 0, el)
		for _, ine := range in.Edges.Entitlements {
			r, err := EntitlementResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build entitlements information for User")
			}
			list = append(list, r)
		}
		builder.Entitlements(list...)
	}

	if el := len(in.Edges.Imses); el > 0 {
		list := make([]*resource.IMS, 0, el)
		for _, ine := range in.Edges.Imses {
			r, err := IMSResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build ims information for User")
			}
			list = append(list, r)
		}
		builder.IMS(list...)
	}

	if el := in.Edges.Name; el != nil {
		r, err := NamesResourceFromEnt(el)
		if err != nil {
			return nil, fmt.Errorf("failed to convert names to SCIM resource: %w", err)
		}
		builder.Name(r)
	}

	if el := len(in.Edges.PhoneNumbers); el > 0 {
		list := make([]*resource.PhoneNumber, 0, el)
		for _, ine := range in.Edges.PhoneNumbers {
			r, err := PhoneNumberResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build phoneNumbers information for User")
			}
			list = append(list, r)
		}
		builder.PhoneNumbers(list...)
	}

	if el := len(in.Edges.Photos); el > 0 {
		list := make([]*resource.Photo, 0, el)
		for _, ine := range in.Edges.Photos {
			r, err := PhotoResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build photos information for User")
			}
			list = append(list, r)
		}
		builder.Photos(list...)
	}

	if el := len(in.Edges.Roles); el > 0 {
		list := make([]*resource.Role, 0, el)
		for _, ine := range in.Edges.Roles {
			r, err := RoleResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build roles information for User")
			}
			list = append(list, r)
		}
		builder.Roles(list...)
	}

	if el := len(in.Edges.X509Certificates); el > 0 {
		list := make([]*resource.X509Certificate, 0, el)
		for _, ine := range in.Edges.X509Certificates {
			r, err := X509CertificateResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build x509Certificates information for User")
			}
			list = append(list, r)
		}
		builder.X509Certificates(list...)
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
	if err := userResourceFromEntHelper(in, builder); err != nil {
		return nil, err
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

func (b *Backend) ReplaceUser(id string, in *resource.User) (*resource.User, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	h := sha256.New()
	fmt.Fprint(h, b.etagSalt)

	u, err := b.db.User.Query().
		Select("id").
		Where(user.IDEQ(parsedUUID)).
		Only(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	replaceUserCall := u.Update().
		ClearAddresses().
		ClearEmails().
		ClearEntitlements().
		ClearGroups().
		ClearImses().
		ClearName().
		ClearPhoneNumbers().
		ClearPhotos().
		ClearRoles().
		ClearX509Certificates()

	if in.HasActive() {
		replaceUserCall.SetActive(in.Active())
		fmt.Fprint(h, in.Active())
	}

	var addresses []*ent.Address
	if in.HasAddresses() {
		created, err := b.createAddresses(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create address: %w", err)
		}
		replaceUserCall.AddAddresses(created...)
		addresses = created
	}

	if in.HasDisplayName() {
		replaceUserCall.SetDisplayName(in.DisplayName())
		fmt.Fprint(h, in.DisplayName())
	}

	var emails []*ent.Email
	if in.HasEmails() {
		created, err := b.createEmails(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create email: %w", err)
		}
		replaceUserCall.AddEmails(created...)
		emails = created
	}

	var entitlements []*ent.Entitlement
	if in.HasEntitlements() {
		created, err := b.createEntitlements(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create entitlement: %w", err)
		}
		replaceUserCall.AddEntitlements(created...)
		entitlements = created
	}

	if in.HasExternalID() {
		replaceUserCall.SetExternalID(in.ExternalID())
		fmt.Fprint(h, in.ExternalID())
	}

	var ims []*ent.IMS
	if in.HasIMS() {
		created, err := b.createIMS(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create ims: %w", err)
		}
		replaceUserCall.AddImses(created...)
		ims = created
	}

	if in.HasLocale() {
		replaceUserCall.SetLocale(in.Locale())
		fmt.Fprint(h, in.Locale())
	}

	var name *ent.Names
	if in.HasName() {
		created, err := b.createName(in.Name(), h)
		if err != nil {
			return nil, fmt.Errorf("failed to create name: %w", err)
		}
		replaceUserCall.SetName(created)
		name = created
	}

	if in.HasNickName() {
		replaceUserCall.SetNickName(in.NickName())
		fmt.Fprint(h, in.NickName())
	}

	if in.HasPassword() {
		replaceUserCall.SetPassword(in.Password())
		fmt.Fprint(h, in.Password())
	}

	var phoneNumbers []*ent.PhoneNumber
	if in.HasPhoneNumbers() {
		created, err := b.createPhoneNumbers(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create phoneNumber: %w", err)
		}
		replaceUserCall.AddPhoneNumbers(created...)
		phoneNumbers = created
	}

	var photos []*ent.Photo
	if in.HasPhotos() {
		created, err := b.createPhotos(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create photo: %w", err)
		}
		replaceUserCall.AddPhotos(created...)
		photos = created
	}

	if in.HasPreferredLanguage() {
		replaceUserCall.SetPreferredLanguage(in.PreferredLanguage())
		fmt.Fprint(h, in.PreferredLanguage())
	}

	if in.HasProfileURL() {
		replaceUserCall.SetProfileURL(in.ProfileURL())
		fmt.Fprint(h, in.ProfileURL())
	}

	var roles []*ent.Role
	if in.HasRoles() {
		created, err := b.createRoles(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create role: %w", err)
		}
		replaceUserCall.AddRoles(created...)
		roles = created
	}

	if in.HasTimezone() {
		replaceUserCall.SetTimezone(in.Timezone())
		fmt.Fprint(h, in.Timezone())
	}

	if in.HasTitle() {
		replaceUserCall.SetTitle(in.Title())
		fmt.Fprint(h, in.Title())
	}

	if in.HasUserName() {
		replaceUserCall.SetUserName(in.UserName())
		fmt.Fprint(h, in.UserName())
	}

	if in.HasUserType() {
		replaceUserCall.SetUserType(in.UserType())
		fmt.Fprint(h, in.UserType())
	}

	var x509Certificates []*ent.X509Certificate
	if in.HasX509Certificates() {
		created, err := b.createX509Certificates(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create x509Certificate: %w", err)
		}
		replaceUserCall.AddX509Certificates(created...)
		x509Certificates = created
	}

	replaceUserCall.SetEtag(fmt.Sprintf("W/%q", base64.RawStdEncoding.EncodeToString(h.Sum(nil))))

	u2, err := replaceUserCall.
		Save(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	u2.Edges.Addresses = addresses
	u2.Edges.Emails = emails
	u2.Edges.Entitlements = entitlements
	u2.Edges.Imses = ims
	u2.Edges.Name = name
	u2.Edges.PhoneNumbers = phoneNumbers
	u2.Edges.Photos = photos
	u2.Edges.Roles = roles
	u2.Edges.X509Certificates = x509Certificates

	return UserResourceFromEnt(u2)
}

func (b *Backend) CreateUser(in *resource.User) (*resource.User, error) {
	password, err := b.generatePassword(in)
	if err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	h := sha256.New()
	fmt.Fprint(h, b.etagSalt)

	createUserCall := b.db.User.Create().
		SetUserName(in.UserName()).
		SetPassword(password)
	fmt.Fprint(h, in.UserName())

	if in.HasActive() {
		createUserCall.SetActive(in.Active())
		fmt.Fprint(h, in.Active())
	}

	var addresses []*ent.Address
	if in.HasAddresses() {
		created, err := b.createAddresses(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddAddresses(created...)
		addresses = created
	}

	if in.HasDisplayName() {
		createUserCall.SetDisplayName(in.DisplayName())
		fmt.Fprint(h, in.DisplayName())
	}

	var emails []*ent.Email
	if in.HasEmails() {
		created, err := b.createEmails(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddEmails(created...)
		emails = created
	}

	var entitlements []*ent.Entitlement
	if in.HasEntitlements() {
		created, err := b.createEntitlements(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddEntitlements(created...)
		entitlements = created
	}

	if in.HasExternalID() {
		createUserCall.SetExternalID(in.ExternalID())
		fmt.Fprint(h, in.ExternalID())
	}

	var ims []*ent.IMS
	if in.HasIMS() {
		created, err := b.createIMS(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddImses(created...)
		ims = created
	}

	if in.HasLocale() {
		createUserCall.SetLocale(in.Locale())
		fmt.Fprint(h, in.Locale())
	}

	var name *ent.Names
	if in.HasName() {
		created, err := b.createName(in.Name(), h)
		if err != nil {
			return nil, fmt.Errorf("failed to create name: %w", err)
		}
		createUserCall.SetName(created)
		name = created
	}

	if in.HasNickName() {
		createUserCall.SetNickName(in.NickName())
		fmt.Fprint(h, in.NickName())
	}

	var phoneNumbers []*ent.PhoneNumber
	if in.HasPhoneNumbers() {
		created, err := b.createPhoneNumbers(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddPhoneNumbers(created...)
		phoneNumbers = created
	}

	var photos []*ent.Photo
	if in.HasPhotos() {
		created, err := b.createPhotos(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddPhotos(created...)
		photos = created
	}

	if in.HasPreferredLanguage() {
		createUserCall.SetPreferredLanguage(in.PreferredLanguage())
		fmt.Fprint(h, in.PreferredLanguage())
	}

	if in.HasProfileURL() {
		createUserCall.SetProfileURL(in.ProfileURL())
		fmt.Fprint(h, in.ProfileURL())
	}

	var roles []*ent.Role
	if in.HasRoles() {
		created, err := b.createRoles(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddRoles(created...)
		roles = created
	}

	if in.HasTimezone() {
		createUserCall.SetTimezone(in.Timezone())
		fmt.Fprint(h, in.Timezone())
	}

	if in.HasTitle() {
		createUserCall.SetTitle(in.Title())
		fmt.Fprint(h, in.Title())
	}

	if in.HasUserType() {
		createUserCall.SetUserType(in.UserType())
		fmt.Fprint(h, in.UserType())
	}

	var x509Certificates []*ent.X509Certificate
	if in.HasX509Certificates() {
		created, err := b.createX509Certificates(in, h)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		createUserCall.AddX509Certificates(created...)
		x509Certificates = created
	}

	createUserCall.SetEtag(fmt.Sprintf("W/%q", base64.RawStdEncoding.EncodeToString(h.Sum(nil))))

	u, err := createUserCall.
		Save(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	u.Edges.Addresses = addresses
	u.Edges.Emails = emails
	u.Edges.Entitlements = entitlements
	u.Edges.Imses = ims
	u.Edges.Name = name
	u.Edges.PhoneNumbers = phoneNumbers
	u.Edges.Photos = photos
	u.Edges.Roles = roles
	u.Edges.X509Certificates = x509Certificates

	return UserResourceFromEnt(u)
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

func (b *Backend) createX509Certificates(in *resource.User, h hash.Hash) ([]*ent.X509Certificate, error) {
	list := make([]*ent.X509Certificate, len(in.X509Certificates()))
	inbound := in.X509Certificates()
	sort.Slice(inbound, func(i, j int) bool {
		return inbound[i].Value() <= inbound[j].Value()
	})

	var hasPrimary bool
	for i, v := range inbound {
		createCall := b.db.X509Certificate.Create()
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
				return nil, fmt.Errorf("invalid user.x509Certificates: multiple x509Certificates have been set to primary")
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
