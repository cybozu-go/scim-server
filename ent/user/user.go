// Code generated by entc, DO NOT EDIT.

package user

import (
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldActive holds the string denoting the active field in the database.
	FieldActive = "active"
	// FieldDisplayName holds the string denoting the displayname field in the database.
	FieldDisplayName = "display_name"
	// FieldExternalID holds the string denoting the externalid field in the database.
	FieldExternalID = "external_id"
	// FieldLocale holds the string denoting the locale field in the database.
	FieldLocale = "locale"
	// FieldNickName holds the string denoting the nickname field in the database.
	FieldNickName = "nick_name"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldPreferredLanguage holds the string denoting the preferredlanguage field in the database.
	FieldPreferredLanguage = "preferred_language"
	// FieldProfileURL holds the string denoting the profileurl field in the database.
	FieldProfileURL = "profile_url"
	// FieldTimezone holds the string denoting the timezone field in the database.
	FieldTimezone = "timezone"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldUserName holds the string denoting the username field in the database.
	FieldUserName = "user_name"
	// FieldUserType holds the string denoting the usertype field in the database.
	FieldUserType = "user_type"
	// EdgeGroups holds the string denoting the groups edge name in mutations.
	EdgeGroups = "groups"
	// EdgeEmails holds the string denoting the emails edge name in mutations.
	EdgeEmails = "emails"
	// EdgeName holds the string denoting the name edge name in mutations.
	EdgeName = "name"
	// EdgeEntitlements holds the string denoting the entitlements edge name in mutations.
	EdgeEntitlements = "entitlements"
	// EdgeRoles holds the string denoting the roles edge name in mutations.
	EdgeRoles = "roles"
	// EdgeImses holds the string denoting the imses edge name in mutations.
	EdgeImses = "imses"
	// EdgePhoneNumbers holds the string denoting the phone_numbers edge name in mutations.
	EdgePhoneNumbers = "phone_numbers"
	// EdgePhotos holds the string denoting the photos edge name in mutations.
	EdgePhotos = "photos"
	// Table holds the table name of the user in the database.
	Table = "users"
	// GroupsTable is the table that holds the groups relation/edge.
	GroupsTable = "groups"
	// GroupsInverseTable is the table name for the Group entity.
	// It exists in this package in order to avoid circular dependency with the "group" package.
	GroupsInverseTable = "groups"
	// GroupsColumn is the table column denoting the groups relation/edge.
	GroupsColumn = "user_groups"
	// EmailsTable is the table that holds the emails relation/edge.
	EmailsTable = "emails"
	// EmailsInverseTable is the table name for the Email entity.
	// It exists in this package in order to avoid circular dependency with the "email" package.
	EmailsInverseTable = "emails"
	// EmailsColumn is the table column denoting the emails relation/edge.
	EmailsColumn = "user_emails"
	// NameTable is the table that holds the name relation/edge.
	NameTable = "names"
	// NameInverseTable is the table name for the Names entity.
	// It exists in this package in order to avoid circular dependency with the "names" package.
	NameInverseTable = "names"
	// NameColumn is the table column denoting the name relation/edge.
	NameColumn = "user_name"
	// EntitlementsTable is the table that holds the entitlements relation/edge.
	EntitlementsTable = "entitlements"
	// EntitlementsInverseTable is the table name for the Entitlement entity.
	// It exists in this package in order to avoid circular dependency with the "entitlement" package.
	EntitlementsInverseTable = "entitlements"
	// EntitlementsColumn is the table column denoting the entitlements relation/edge.
	EntitlementsColumn = "user_entitlements"
	// RolesTable is the table that holds the roles relation/edge.
	RolesTable = "roles"
	// RolesInverseTable is the table name for the Role entity.
	// It exists in this package in order to avoid circular dependency with the "role" package.
	RolesInverseTable = "roles"
	// RolesColumn is the table column denoting the roles relation/edge.
	RolesColumn = "user_roles"
	// ImsesTable is the table that holds the imses relation/edge.
	ImsesTable = "im_ss"
	// ImsesInverseTable is the table name for the IMS entity.
	// It exists in this package in order to avoid circular dependency with the "ims" package.
	ImsesInverseTable = "im_ss"
	// ImsesColumn is the table column denoting the imses relation/edge.
	ImsesColumn = "user_imses"
	// PhoneNumbersTable is the table that holds the phone_numbers relation/edge.
	PhoneNumbersTable = "phone_numbers"
	// PhoneNumbersInverseTable is the table name for the PhoneNumber entity.
	// It exists in this package in order to avoid circular dependency with the "phonenumber" package.
	PhoneNumbersInverseTable = "phone_numbers"
	// PhoneNumbersColumn is the table column denoting the phone_numbers relation/edge.
	PhoneNumbersColumn = "user_phone_numbers"
	// PhotosTable is the table that holds the photos relation/edge.
	PhotosTable = "photos"
	// PhotosInverseTable is the table name for the Photo entity.
	// It exists in this package in order to avoid circular dependency with the "photo" package.
	PhotosInverseTable = "photos"
	// PhotosColumn is the table column denoting the photos relation/edge.
	PhotosColumn = "user_photos"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldActive,
	FieldDisplayName,
	FieldExternalID,
	FieldLocale,
	FieldNickName,
	FieldPassword,
	FieldPreferredLanguage,
	FieldProfileURL,
	FieldTimezone,
	FieldTitle,
	FieldUserName,
	FieldUserType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "users"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"group_users",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	PasswordValidator func(string) error
	// UserNameValidator is a validator for the "userName" field. It is called by the builders before save.
	UserNameValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)
