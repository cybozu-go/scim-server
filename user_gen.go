package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/address"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim-server/ent/entitlement"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim-server/ent/phonenumber"
	"github.com/cybozu-go/scim-server/ent/photo"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/role"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim-server/ent/x509certificate"
	"github.com/cybozu-go/scim/filter"
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
		case resource.UserIDKey:
			selectNames = append(selectNames, user.FieldID)
		case resource.UserIMSKey:
			q.WithIMS()
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

	if el := len(in.Edges.IMS); el > 0 {
		list := make([]*resource.IMS, 0, el)
		for _, ine := range in.Edges.IMS {
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

	if gl := len(in.Groups); gl > 0 {
		memberships := make([]*resource.GroupMember, gl)
		for i, m := range in.Groups {
			var gmb resource.GroupMemberBuilder
			gm, err := gmb.Value(m.Value).
				Ref(m.Ref).
				Type(m.Type).
				Build()
			if err != nil {
				return nil, fmt.Errorf("failed to compute \"groups\" field: %w", err)
			}
			memberships[i] = gm
		}
		builder.Groups(memberships...)
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
	if !reflect.ValueOf(in.Title).IsZero() {
		builder.Title(in.Title)
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

func (b *Backend) existsUserAddress(ctx context.Context, parent *ent.User, in *resource.Address) bool {
	queryCall := parent.QueryAddresses()
	var predicates []predicate.Address
	if in.HasCountry() {
		predicates = append(predicates, address.Country(in.Country()))
	}
	if in.HasFormatted() {
		predicates = append(predicates, address.Formatted(in.Formatted()))
	}
	if in.HasLocality() {
		predicates = append(predicates, address.Locality(in.Locality()))
	}
	if in.HasPostalCode() {
		predicates = append(predicates, address.PostalCode(in.PostalCode()))
	}
	if in.HasRegion() {
		predicates = append(predicates, address.Region(in.Region()))
	}
	if in.HasStreetAddress() {
		predicates = append(predicates, address.StreetAddress(in.StreetAddress()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserEmail(ctx context.Context, parent *ent.User, in *resource.Email) bool {
	queryCall := parent.QueryEmails()
	var predicates []predicate.Email
	if in.HasDisplay() {
		predicates = append(predicates, email.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, email.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, email.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, email.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserEntitlement(ctx context.Context, parent *ent.User, in *resource.Entitlement) bool {
	queryCall := parent.QueryEntitlements()
	var predicates []predicate.Entitlement
	if in.HasDisplay() {
		predicates = append(predicates, entitlement.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, entitlement.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, entitlement.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, entitlement.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserIMS(ctx context.Context, parent *ent.User, in *resource.IMS) bool {
	queryCall := parent.QueryIMS()
	var predicates []predicate.IMS
	if in.HasDisplay() {
		predicates = append(predicates, ims.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, ims.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, ims.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, ims.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserPhoneNumber(ctx context.Context, parent *ent.User, in *resource.PhoneNumber) bool {
	queryCall := parent.QueryPhoneNumbers()
	var predicates []predicate.PhoneNumber
	if in.HasDisplay() {
		predicates = append(predicates, phonenumber.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, phonenumber.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, phonenumber.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, phonenumber.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserPhoto(ctx context.Context, parent *ent.User, in *resource.Photo) bool {
	queryCall := parent.QueryPhotos()
	var predicates []predicate.Photo
	if in.HasDisplay() {
		predicates = append(predicates, photo.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, photo.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, photo.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, photo.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserRole(ctx context.Context, parent *ent.User, in *resource.Role) bool {
	queryCall := parent.QueryRoles()
	var predicates []predicate.Role
	if in.HasDisplay() {
		predicates = append(predicates, role.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, role.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, role.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, role.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) existsUserX509Certificate(ctx context.Context, parent *ent.User, in *resource.X509Certificate) bool {
	queryCall := parent.QueryX509Certificates()
	var predicates []predicate.X509Certificate
	if in.HasDisplay() {
		predicates = append(predicates, x509certificate.Display(in.Display()))
	}
	if in.HasPrimary() {
		predicates = append(predicates, x509certificate.Primary(in.Primary()))
	}
	if in.HasType() {
		predicates = append(predicates, x509certificate.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, x509certificate.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) createEmail(ctx context.Context, resources ...*resource.Email) ([]*ent.EmailCreate, error) {
	list := make([]*ent.EmailCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.Email.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) createEntitlement(ctx context.Context, resources ...*resource.Entitlement) ([]*ent.EntitlementCreate, error) {
	list := make([]*ent.EntitlementCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.Entitlement.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) createIMS(ctx context.Context, resources ...*resource.IMS) ([]*ent.IMSCreate, error) {
	list := make([]*ent.IMSCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.IMS.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) createPhoneNumber(ctx context.Context, resources ...*resource.PhoneNumber) ([]*ent.PhoneNumberCreate, error) {
	list := make([]*ent.PhoneNumberCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.PhoneNumber.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) createPhoto(ctx context.Context, resources ...*resource.Photo) ([]*ent.PhotoCreate, error) {
	list := make([]*ent.PhotoCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.Photo.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) createRole(ctx context.Context, resources ...*resource.Role) ([]*ent.RoleCreate, error) {
	list := make([]*ent.RoleCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.Role.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) createX509Certificate(ctx context.Context, resources ...*resource.X509Certificate) ([]*ent.X509CertificateCreate, error) {
	list := make([]*ent.X509CertificateCreate, len(resources))
	for i, in := range resources {
		createCall := b.db.X509Certificate.Create()
		if in.HasDisplay() {
			createCall.SetDisplay(in.Display())
		}
		if in.HasPrimary() {
			createCall.SetPrimary(in.Primary())
		}
		if in.HasType() {
			createCall.SetType(in.Type())
		}
		if in.HasValue() {
			createCall.SetValue(in.Value())
		}
		list[i] = createCall
	}
	return list, nil
}

func (b *Backend) CreateUser(ctx context.Context, in *resource.User) (*resource.User, error) {

	createCall := b.db.User.Create()
	password, err := b.generatePassword(in)
	if err != nil {
		return nil, fmt.Errorf("failed to process password: %w", err)
	}
	createCall.SetPassword(password)
	if !in.HasUserName() {
		return nil, fmt.Errorf("required field userName not found")
	}
	createCall.SetUserName(in.UserName())
	if in.HasActive() {
		createCall.SetActive(in.Active())
	}
	var addressCreateCalls []*ent.AddressCreate
	if in.HasAddresses() {
		calls, err := b.createAddress(ctx, in.Addresses()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create addresses: %w", err)
		}
		addressCreateCalls = calls
	}
	if in.HasDisplayName() {
		createCall.SetDisplayName(in.DisplayName())
	}
	var emailCreateCalls []*ent.EmailCreate
	if in.HasEmails() {
		calls, err := b.createEmail(ctx, in.Emails()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create emails: %w", err)
		}
		emailCreateCalls = calls
	}
	var entitlementCreateCalls []*ent.EntitlementCreate
	if in.HasEntitlements() {
		calls, err := b.createEntitlement(ctx, in.Entitlements()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create entitlements: %w", err)
		}
		entitlementCreateCalls = calls
	}
	if in.HasExternalID() {
		createCall.SetExternalID(in.ExternalID())
	}
	var imsCreateCalls []*ent.IMSCreate
	if in.HasIMS() {
		calls, err := b.createIMS(ctx, in.IMS()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create ims: %w", err)
		}
		imsCreateCalls = calls
	}
	if in.HasLocale() {
		createCall.SetLocale(in.Locale())
	}
	if in.HasName() {
		created, err := b.createName(ctx, in.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to create name: %w", err)
		}
		createCall.SetName(created)
	}
	if in.HasNickName() {
		createCall.SetNickName(in.NickName())
	}
	var phoneNumberCreateCalls []*ent.PhoneNumberCreate
	if in.HasPhoneNumbers() {
		calls, err := b.createPhoneNumber(ctx, in.PhoneNumbers()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create phoneNumbers: %w", err)
		}
		phoneNumberCreateCalls = calls
	}
	var photoCreateCalls []*ent.PhotoCreate
	if in.HasPhotos() {
		calls, err := b.createPhoto(ctx, in.Photos()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create photos: %w", err)
		}
		photoCreateCalls = calls
	}
	if in.HasPreferredLanguage() {
		createCall.SetPreferredLanguage(in.PreferredLanguage())
	}
	if in.HasProfileURL() {
		createCall.SetProfileURL(in.ProfileURL())
	}
	var roleCreateCalls []*ent.RoleCreate
	if in.HasRoles() {
		calls, err := b.createRole(ctx, in.Roles()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		roleCreateCalls = calls
	}
	if in.HasTimezone() {
		createCall.SetTimezone(in.Timezone())
	}
	if in.HasTitle() {
		createCall.SetTitle(in.Title())
	}
	if in.HasUserType() {
		createCall.SetUserType(in.UserType())
	}
	var x509CertificateCreateCalls []*ent.X509CertificateCreate
	if in.HasX509Certificates() {
		calls, err := b.createX509Certificate(ctx, in.X509Certificates()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create x509Certificates: %w", err)
		}
		x509CertificateCreateCalls = calls
	}

	rs, err := createCall.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save object: %w", err)
	}

	for _, call := range addressCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create addresses: %w", err)
		}
		rs.Edges.Addresses = append(rs.Edges.Addresses, created)
	}

	for _, call := range emailCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create emails: %w", err)
		}
		rs.Edges.Emails = append(rs.Edges.Emails, created)
	}

	for _, call := range entitlementCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create entitlements: %w", err)
		}
		rs.Edges.Entitlements = append(rs.Edges.Entitlements, created)
	}

	for _, call := range imsCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create ims: %w", err)
		}
		rs.Edges.IMS = append(rs.Edges.IMS, created)
	}

	for _, call := range phoneNumberCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create phoneNumbers: %w", err)
		}
		rs.Edges.PhoneNumbers = append(rs.Edges.PhoneNumbers, created)
	}

	for _, call := range photoCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create photos: %w", err)
		}
		rs.Edges.Photos = append(rs.Edges.Photos, created)
	}

	for _, call := range roleCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		rs.Edges.Roles = append(rs.Edges.Roles, created)
	}

	for _, call := range x509CertificateCreateCalls {
		call.SetUserID(rs.ID)
		created, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create x509Certificates: %w", err)
		}
		rs.Edges.X509Certificates = append(rs.Edges.X509Certificates, created)
	}

	h := sha256.New()
	if err := rs.ComputeETag(h); err != nil {
		return nil, fmt.Errorf("failed to compute etag: %w", err)
	}
	etag := fmt.Sprintf("W/%x", h.Sum(nil))

	if _, err := rs.Update().SetEtag(etag).Save(ctx); err != nil {
		return nil, fmt.Errorf("failed to save etag: %w", err)
	}
	rs.Etag = etag
	return UserResourceFromEnt(rs)
}

func (b *Backend) ReplaceUser(ctx context.Context, id string, in *resource.User) (*resource.User, error) {

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	r, err := b.db.User.Query().Where(user.ID(parsedUUID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve resource for replacing: %w", err)
	}

	replaceCall := r.Update()

	replaceCall.ClearActive()
	if in.HasActive() {
		replaceCall.SetActive(in.Active())
	}

	replaceCall.ClearAddresses()
	var addressesCreateCalls []*ent.AddressCreate
	if in.HasAddresses() {
		calls, err := b.createAddress(ctx, in.Addresses()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create addresses: %w", err)
		}
		addressesCreateCalls = calls
	}

	replaceCall.ClearDisplayName()
	if in.HasDisplayName() {
		replaceCall.SetDisplayName(in.DisplayName())
	}

	replaceCall.ClearEmails()
	var emailsCreateCalls []*ent.EmailCreate
	if in.HasEmails() {
		calls, err := b.createEmail(ctx, in.Emails()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create emails: %w", err)
		}
		emailsCreateCalls = calls
	}

	replaceCall.ClearEntitlements()
	var entitlementsCreateCalls []*ent.EntitlementCreate
	if in.HasEntitlements() {
		calls, err := b.createEntitlement(ctx, in.Entitlements()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create entitlements: %w", err)
		}
		entitlementsCreateCalls = calls
	}

	replaceCall.ClearExternalID()
	if in.HasExternalID() {
		replaceCall.SetExternalID(in.ExternalID())
	}

	replaceCall.ClearIMS()
	var imsCreateCalls []*ent.IMSCreate
	if in.HasIMS() {
		calls, err := b.createIMS(ctx, in.IMS()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create ims: %w", err)
		}
		imsCreateCalls = calls
	}

	replaceCall.ClearLocale()
	if in.HasLocale() {
		replaceCall.SetLocale(in.Locale())
	}

	replaceCall.ClearName()
	if in.HasName() {
		created, err := b.createName(ctx, in.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to create name: %w", err)
		}
		replaceCall.SetName(created)
	}

	replaceCall.ClearNickName()
	if in.HasNickName() {
		replaceCall.SetNickName(in.NickName())
	}

	replaceCall.ClearPassword()
	if in.HasPassword() {
		replaceCall.SetPassword(in.Password())
	}

	replaceCall.ClearPhoneNumbers()
	var phoneNumbersCreateCalls []*ent.PhoneNumberCreate
	if in.HasPhoneNumbers() {
		calls, err := b.createPhoneNumber(ctx, in.PhoneNumbers()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create phoneNumbers: %w", err)
		}
		phoneNumbersCreateCalls = calls
	}

	replaceCall.ClearPhotos()
	var photosCreateCalls []*ent.PhotoCreate
	if in.HasPhotos() {
		calls, err := b.createPhoto(ctx, in.Photos()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create photos: %w", err)
		}
		photosCreateCalls = calls
	}

	replaceCall.ClearPreferredLanguage()
	if in.HasPreferredLanguage() {
		replaceCall.SetPreferredLanguage(in.PreferredLanguage())
	}

	replaceCall.ClearProfileURL()
	if in.HasProfileURL() {
		replaceCall.SetProfileURL(in.ProfileURL())
	}

	replaceCall.ClearRoles()
	var rolesCreateCalls []*ent.RoleCreate
	if in.HasRoles() {
		calls, err := b.createRole(ctx, in.Roles()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create roles: %w", err)
		}
		rolesCreateCalls = calls
	}

	replaceCall.ClearTimezone()
	if in.HasTimezone() {
		replaceCall.SetTimezone(in.Timezone())
	}

	replaceCall.ClearTitle()
	if in.HasTitle() {
		replaceCall.SetTitle(in.Title())
	}
	if in.HasUserName() {
		replaceCall.SetUserName(in.UserName())
	}

	replaceCall.ClearUserType()
	if in.HasUserType() {
		replaceCall.SetUserType(in.UserType())
	}

	replaceCall.ClearX509Certificates()
	var x509CertificatesCreateCalls []*ent.X509CertificateCreate
	if in.HasX509Certificates() {
		calls, err := b.createX509Certificate(ctx, in.X509Certificates()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create x509Certificates: %w", err)
		}
		x509CertificatesCreateCalls = calls
	}

	if _, err := replaceCall.Save(ctx); err != nil {
		return nil, fmt.Errorf("failed to save object: %w", err)
	}
	for _, call := range addressesCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save Addresses: %w", err)
		}
	}
	for _, call := range emailsCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save Emails: %w", err)
		}
	}
	for _, call := range entitlementsCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save Entitlements: %w", err)
		}
	}
	for _, call := range imsCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save IMS: %w", err)
		}
	}
	for _, call := range phoneNumbersCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save PhoneNumbers: %w", err)
		}
	}
	for _, call := range photosCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save Photos: %w", err)
		}
	}
	for _, call := range rolesCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save Roles: %w", err)
		}
	}
	for _, call := range x509CertificatesCreateCalls {
		call.SetUserID(parsedUUID)
		_, err := call.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to save X509Certificates: %w", err)
		}
	}

	r2, err := b.db.User.Query().Where(user.ID(parsedUUID)).
		WithAddresses().
		WithEmails().
		WithEntitlements().
		WithIMS().
		WithPhoneNumbers().
		WithPhotos().
		WithRoles().
		WithX509Certificates().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data")
	}

	h := sha256.New()
	if err := r2.ComputeETag(h); err != nil {
		return nil, fmt.Errorf("failed to compute etag: %w", err)
	}
	etag := fmt.Sprintf("W/%x", h.Sum(nil))

	if _, err := r2.Update().SetEtag(etag).Save(ctx); err != nil {
		return nil, fmt.Errorf("failed to save etag: %w", err)
	}
	r2.Etag = etag

	return UserResourceFromEnt(r2)
}

func (b *Backend) patchAddUser(ctx context.Context, parent *ent.User, op *resource.PatchOperation) error {

	root, err := filter.Parse(op.Path(), filter.WithPatchExpression(true))
	if err != nil {
		return fmt.Errorf("failed to parse PATH path %q", op.Path())
	}

	expr, ok := root.(filter.ValuePath)
	if !ok {
		return fmt.Errorf("root element should be a valuePath (got %T)", root)
	}

	sattr, err := exprStr(expr.ParentAttr())
	if err != nil {
		return fmt.Errorf("invalid attribute specification: %w", err)
	}

	switch sattr {
	case resource.UserAddressesKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item addresses with unspecified element is not possible")
			}

			var in resource.Address
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserAddress(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createAddress(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create Address: %w", err)
			}
			list := make([]*ent.Address, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create Address: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb AddressPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryAddresses().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.AddressCountryKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetCountry(v)
			case resource.AddressFormattedKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetFormatted(v)
			case resource.AddressLocalityKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetLocality(v)
			case resource.AddressPostalCodeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPostalCode(v)
			case resource.AddressRegionKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetRegion(v)
			case resource.AddressStreetAddressKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetStreetAddress(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserDisplayNameKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element displayName")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element displayName")
		}

		if _, err := parent.Update().SetDisplayName(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserEmailsKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item emails with unspecified element is not possible")
			}

			var in resource.Email
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserEmail(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createEmail(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create Email: %w", err)
			}
			list := make([]*ent.Email, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create Email: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb EmailPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryEmails().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.EmailDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.EmailPrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.EmailTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.EmailValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserEntitlementsKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item entitlements with unspecified element is not possible")
			}

			var in resource.Entitlement
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserEntitlement(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createEntitlement(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create Entitlement: %w", err)
			}
			list := make([]*ent.Entitlement, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create Entitlement: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb EntitlementPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryEntitlements().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.EntitlementDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.EntitlementPrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.EntitlementTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.EntitlementValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserExternalIDKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element externalId")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element externalId")
		}

		if _, err := parent.Update().SetExternalID(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserGroupsKey:
		return fmt.Errorf("cannot create group memberships through User resource")
	case resource.UserIMSKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item ims with unspecified element is not possible")
			}

			var in resource.IMS
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserIMS(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createIMS(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create IMS: %w", err)
			}
			list := make([]*ent.IMS, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create IMS: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb IMSPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryIMS().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.IMSDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.IMSPrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.IMSTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.IMSValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserLocaleKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element locale")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element locale")
		}

		if _, err := parent.Update().SetLocale(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserNickNameKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element nickName")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element nickName")
		}

		if _, err := parent.Update().SetNickName(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserPasswordKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element password")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element password")
		}

		if _, err := parent.Update().SetPassword(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserPhoneNumbersKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item phoneNumbers with unspecified element is not possible")
			}

			var in resource.PhoneNumber
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserPhoneNumber(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createPhoneNumber(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create PhoneNumber: %w", err)
			}
			list := make([]*ent.PhoneNumber, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create PhoneNumber: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb PhoneNumberPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryPhoneNumbers().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.PhoneNumberDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.PhoneNumberPrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.PhoneNumberTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.PhoneNumberValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserPhotosKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item photos with unspecified element is not possible")
			}

			var in resource.Photo
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserPhoto(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createPhoto(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create Photo: %w", err)
			}
			list := make([]*ent.Photo, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create Photo: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb PhotoPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryPhotos().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.PhotoDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.PhotoPrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.PhotoTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.PhotoValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserPreferredLanguageKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element preferredLanguage")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element preferredLanguage")
		}

		if _, err := parent.Update().SetPreferredLanguage(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserProfileURLKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element profileUrl")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element profileUrl")
		}

		if _, err := parent.Update().SetProfileURL(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserRolesKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item roles with unspecified element is not possible")
			}

			var in resource.Role
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserRole(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createRole(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create Role: %w", err)
			}
			list := make([]*ent.Role, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create Role: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb RolePredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryRoles().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.RoleDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.RolePrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.RoleTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.RoleValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	case resource.UserTimezoneKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element timezone")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element timezone")
		}

		if _, err := parent.Update().SetTimezone(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserTitleKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element title")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element title")
		}

		if _, err := parent.Update().SetTitle(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserUserNameKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element userName")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element userName")
		}

		if _, err := parent.Update().SetUserName(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserUserTypeKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element userType")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element userType")
		}

		if _, err := parent.Update().SetUserType(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserX509CertificatesKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item x509Certificates with unspecified element is not possible")
			}

			var in resource.X509Certificate
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsUserX509Certificate(ctx, parent, &in) {
				return nil
			}

			calls, err := b.createX509Certificate(ctx, &in)
			if err != nil {
				return fmt.Errorf("failed to create X509Certificate: %w", err)
			}
			list := make([]*ent.X509Certificate, len(calls))
			for i, call := range calls {
				call.SetUserID(parent.ID)
				created, err := call.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create X509Certificate: %w", err)
				}
				list[i] = created
			}
		} else {
			var pb X509CertificatePredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryX509Certificates().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.X509CertificateDisplayKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetDisplay(v)
			case resource.X509CertificatePrimaryKey:
				var v bool
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetPrimary(v)
			case resource.X509CertificateTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.X509CertificateValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (b *Backend) patchRemoveUser(ctx context.Context, parent *ent.User, op *resource.PatchOperation) error {
	if op.Path() == "" {
		return resource.NewErrorBuilder().
			Status(http.StatusBadRequest).
			ScimType(resource.ErrNoTarget).
			Detail("empty path").
			MustBuild()
	}

	root, err := filter.Parse(op.Path(), filter.WithPatchExpression(true))
	if err != nil {
		return fmt.Errorf("failed to parse path %q", op.Path())
	}

	expr, ok := root.(filter.ValuePath)
	if !ok {
		return fmt.Errorf("root element should be a valuePath (got %T)", root)
	}

	sattr, err := exprStr(expr.ParentAttr())
	if err != nil {
		return fmt.Errorf("invalid attribute specification: %w", err)
	}
	switch sattr {
	case resource.UserActiveKey:
	case resource.UserAddressesKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item addresses without a query is not possible")
			}
			if _, err := b.db.Address.Delete().Where(address.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from addresses: %w", err)
			}
			if _, err := parent.Update().ClearAddresses().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to addresses: %w", err)
			}
		} else {
			var pb AddressPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryAddresses().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.AddressCountryKey:
					return fmt.Errorf("country is not mutable")
				case resource.AddressFormattedKey:
					return fmt.Errorf("formatted is not mutable")
				case resource.AddressLocalityKey:
					return fmt.Errorf("locality is not mutable")
				case resource.AddressPostalCodeKey:
					return fmt.Errorf("postalCode is not mutable")
				case resource.AddressRegionKey:
					return fmt.Errorf("region is not mutable")
				case resource.AddressStreetAddressKey:
					return fmt.Errorf("streetAddress is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.Address.Delete().Where(address.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserDisplayNameKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on displayName cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on displayName cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearDisplayName().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserEmailsKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item emails without a query is not possible")
			}
			if _, err := b.db.Email.Delete().Where(email.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from emails: %w", err)
			}
			if _, err := parent.Update().ClearEmails().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to emails: %w", err)
			}
		} else {
			var pb EmailPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryEmails().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.EmailDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.EmailPrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.EmailTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.EmailValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.Email.Delete().Where(email.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserEntitlementsKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item entitlements without a query is not possible")
			}
			if _, err := b.db.Entitlement.Delete().Where(entitlement.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from entitlements: %w", err)
			}
			if _, err := parent.Update().ClearEntitlements().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to entitlements: %w", err)
			}
		} else {
			var pb EntitlementPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryEntitlements().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.EntitlementDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.EntitlementPrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.EntitlementTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.EntitlementValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.Entitlement.Delete().Where(entitlement.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserExternalIDKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on externalId cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on externalId cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearExternalID().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserGroupsKey:
		return fmt.Errorf("cannot delete group memberships through User resource")
	case resource.UserIMSKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item ims without a query is not possible")
			}
			if _, err := b.db.IMS.Delete().Where(ims.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from ims: %w", err)
			}
			if _, err := parent.Update().ClearIMS().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to ims: %w", err)
			}
		} else {
			var pb IMSPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryIMS().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.IMSDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.IMSPrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.IMSTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.IMSValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.IMS.Delete().Where(ims.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserLocaleKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on locale cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on locale cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearLocale().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserNameKey:
	case resource.UserNickNameKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on nickName cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on nickName cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearNickName().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserPasswordKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on password cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on password cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearPassword().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserPhoneNumbersKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item phoneNumbers without a query is not possible")
			}
			if _, err := b.db.PhoneNumber.Delete().Where(phonenumber.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from phoneNumbers: %w", err)
			}
			if _, err := parent.Update().ClearPhoneNumbers().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to phoneNumbers: %w", err)
			}
		} else {
			var pb PhoneNumberPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryPhoneNumbers().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.PhoneNumberDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.PhoneNumberPrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.PhoneNumberTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.PhoneNumberValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.PhoneNumber.Delete().Where(phonenumber.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserPhotosKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item photos without a query is not possible")
			}
			if _, err := b.db.Photo.Delete().Where(photo.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from photos: %w", err)
			}
			if _, err := parent.Update().ClearPhotos().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to photos: %w", err)
			}
		} else {
			var pb PhotoPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryPhotos().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.PhotoDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.PhotoPrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.PhotoTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.PhotoValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.Photo.Delete().Where(photo.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserPreferredLanguageKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on preferredLanguage cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on preferredLanguage cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearPreferredLanguage().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserProfileURLKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on profileUrl cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on profileUrl cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearProfileURL().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserRolesKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item roles without a query is not possible")
			}
			if _, err := b.db.Role.Delete().Where(role.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from roles: %w", err)
			}
			if _, err := parent.Update().ClearRoles().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to roles: %w", err)
			}
		} else {
			var pb RolePredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryRoles().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.RoleDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.RolePrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.RoleTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.RoleValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.Role.Delete().Where(role.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	case resource.UserTimezoneKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on timezone cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on timezone cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearTimezone().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserTitleKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on title cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on title cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearTitle().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserUserTypeKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on userType cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on userType cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearUserType().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.UserX509CertificatesKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item x509Certificates without a query is not possible")
			}
			if _, err := b.db.X509Certificate.Delete().Where(x509certificate.HasUserWith(user.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from x509Certificates: %w", err)
			}
			if _, err := parent.Update().ClearX509Certificates().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to x509Certificates: %w", err)
			}
		} else {
			var pb X509CertificatePredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryX509Certificates().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.X509CertificateDisplayKey:
					return fmt.Errorf("display is not mutable")
				case resource.X509CertificatePrimaryKey:
					return fmt.Errorf("primary is not mutable")
				case resource.X509CertificateTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.X509CertificateValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]uuid.UUID, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.X509Certificate.Delete().Where(x509certificate.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	}
	return nil
}
