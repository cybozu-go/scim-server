// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/cybozu-go/scim-server/ent/address"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim-server/ent/entitlement"
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim-server/ent/names"
	"github.com/cybozu-go/scim-server/ent/phonenumber"
	"github.com/cybozu-go/scim-server/ent/photo"
	"github.com/cybozu-go/scim-server/ent/role"
	"github.com/cybozu-go/scim-server/ent/schema"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim-server/ent/x509certificate"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	addressFields := schema.Address{}.Fields()
	_ = addressFields
	// addressDescID is the schema descriptor for id field.
	addressDescID := addressFields[0].Descriptor()
	// address.DefaultID holds the default value on creation for the id field.
	address.DefaultID = addressDescID.Default.(func() uuid.UUID)
	emailFields := schema.Email{}.Fields()
	_ = emailFields
	// emailDescID is the schema descriptor for id field.
	emailDescID := emailFields[0].Descriptor()
	// email.DefaultID holds the default value on creation for the id field.
	email.DefaultID = emailDescID.Default.(func() uuid.UUID)
	entitlementFields := schema.Entitlement{}.Fields()
	_ = entitlementFields
	// entitlementDescID is the schema descriptor for id field.
	entitlementDescID := entitlementFields[0].Descriptor()
	// entitlement.DefaultID holds the default value on creation for the id field.
	entitlement.DefaultID = entitlementDescID.Default.(func() uuid.UUID)
	groupFields := schema.Group{}.Fields()
	_ = groupFields
	// groupDescID is the schema descriptor for id field.
	groupDescID := groupFields[3].Descriptor()
	// group.DefaultID holds the default value on creation for the id field.
	group.DefaultID = groupDescID.Default.(func() uuid.UUID)
	imsFields := schema.IMS{}.Fields()
	_ = imsFields
	// imsDescID is the schema descriptor for id field.
	imsDescID := imsFields[0].Descriptor()
	// ims.DefaultID holds the default value on creation for the id field.
	ims.DefaultID = imsDescID.Default.(func() uuid.UUID)
	namesFields := schema.Names{}.Fields()
	_ = namesFields
	// namesDescID is the schema descriptor for id field.
	namesDescID := namesFields[0].Descriptor()
	// names.DefaultID holds the default value on creation for the id field.
	names.DefaultID = namesDescID.Default.(func() uuid.UUID)
	phonenumberFields := schema.PhoneNumber{}.Fields()
	_ = phonenumberFields
	// phonenumberDescID is the schema descriptor for id field.
	phonenumberDescID := phonenumberFields[0].Descriptor()
	// phonenumber.DefaultID holds the default value on creation for the id field.
	phonenumber.DefaultID = phonenumberDescID.Default.(func() uuid.UUID)
	photoFields := schema.Photo{}.Fields()
	_ = photoFields
	// photoDescID is the schema descriptor for id field.
	photoDescID := photoFields[0].Descriptor()
	// photo.DefaultID holds the default value on creation for the id field.
	photo.DefaultID = photoDescID.Default.(func() uuid.UUID)
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescID is the schema descriptor for id field.
	roleDescID := roleFields[0].Descriptor()
	// role.DefaultID holds the default value on creation for the id field.
	role.DefaultID = roleDescID.Default.(func() uuid.UUID)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[8].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = userDescPassword.Validators[0].(func(string) error)
	// userDescUserName is the schema descriptor for userName field.
	userDescUserName := userFields[13].Descriptor()
	// user.UserNameValidator is a validator for the "userName" field. It is called by the builders before save.
	user.UserNameValidator = userDescUserName.Validators[0].(func(string) error)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[4].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
	x509certificateFields := schema.X509Certificate{}.Fields()
	_ = x509certificateFields
	// x509certificateDescID is the schema descriptor for id field.
	x509certificateDescID := x509certificateFields[0].Descriptor()
	// x509certificate.DefaultID holds the default value on creation for the id field.
	x509certificate.DefaultID = x509certificateDescID.Default.(func() uuid.UUID)
}
