// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/address"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim-server/ent/entitlement"
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim-server/ent/names"
	"github.com/cybozu-go/scim-server/ent/phonenumber"
	"github.com/cybozu-go/scim-server/ent/photo"
	"github.com/cybozu-go/scim-server/ent/role"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim-server/ent/x509certificate"
	"github.com/google/uuid"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	mutation *UserMutation
	hooks    []Hook
}

// SetActive sets the "active" field.
func (uc *UserCreate) SetActive(b bool) *UserCreate {
	uc.mutation.SetActive(b)
	return uc
}

// SetNillableActive sets the "active" field if the given value is not nil.
func (uc *UserCreate) SetNillableActive(b *bool) *UserCreate {
	if b != nil {
		uc.SetActive(*b)
	}
	return uc
}

// SetDisplayName sets the "displayName" field.
func (uc *UserCreate) SetDisplayName(s string) *UserCreate {
	uc.mutation.SetDisplayName(s)
	return uc
}

// SetNillableDisplayName sets the "displayName" field if the given value is not nil.
func (uc *UserCreate) SetNillableDisplayName(s *string) *UserCreate {
	if s != nil {
		uc.SetDisplayName(*s)
	}
	return uc
}

// SetExternalID sets the "externalID" field.
func (uc *UserCreate) SetExternalID(s string) *UserCreate {
	uc.mutation.SetExternalID(s)
	return uc
}

// SetNillableExternalID sets the "externalID" field if the given value is not nil.
func (uc *UserCreate) SetNillableExternalID(s *string) *UserCreate {
	if s != nil {
		uc.SetExternalID(*s)
	}
	return uc
}

// SetLocale sets the "locale" field.
func (uc *UserCreate) SetLocale(s string) *UserCreate {
	uc.mutation.SetLocale(s)
	return uc
}

// SetNillableLocale sets the "locale" field if the given value is not nil.
func (uc *UserCreate) SetNillableLocale(s *string) *UserCreate {
	if s != nil {
		uc.SetLocale(*s)
	}
	return uc
}

// SetNickName sets the "nickName" field.
func (uc *UserCreate) SetNickName(s string) *UserCreate {
	uc.mutation.SetNickName(s)
	return uc
}

// SetNillableNickName sets the "nickName" field if the given value is not nil.
func (uc *UserCreate) SetNillableNickName(s *string) *UserCreate {
	if s != nil {
		uc.SetNickName(*s)
	}
	return uc
}

// SetPassword sets the "password" field.
func (uc *UserCreate) SetPassword(s string) *UserCreate {
	uc.mutation.SetPassword(s)
	return uc
}

// SetNillablePassword sets the "password" field if the given value is not nil.
func (uc *UserCreate) SetNillablePassword(s *string) *UserCreate {
	if s != nil {
		uc.SetPassword(*s)
	}
	return uc
}

// SetPreferredLanguage sets the "preferredLanguage" field.
func (uc *UserCreate) SetPreferredLanguage(s string) *UserCreate {
	uc.mutation.SetPreferredLanguage(s)
	return uc
}

// SetNillablePreferredLanguage sets the "preferredLanguage" field if the given value is not nil.
func (uc *UserCreate) SetNillablePreferredLanguage(s *string) *UserCreate {
	if s != nil {
		uc.SetPreferredLanguage(*s)
	}
	return uc
}

// SetProfileURL sets the "profileURL" field.
func (uc *UserCreate) SetProfileURL(s string) *UserCreate {
	uc.mutation.SetProfileURL(s)
	return uc
}

// SetNillableProfileURL sets the "profileURL" field if the given value is not nil.
func (uc *UserCreate) SetNillableProfileURL(s *string) *UserCreate {
	if s != nil {
		uc.SetProfileURL(*s)
	}
	return uc
}

// SetTimezone sets the "timezone" field.
func (uc *UserCreate) SetTimezone(s string) *UserCreate {
	uc.mutation.SetTimezone(s)
	return uc
}

// SetNillableTimezone sets the "timezone" field if the given value is not nil.
func (uc *UserCreate) SetNillableTimezone(s *string) *UserCreate {
	if s != nil {
		uc.SetTimezone(*s)
	}
	return uc
}

// SetTitle sets the "title" field.
func (uc *UserCreate) SetTitle(s string) *UserCreate {
	uc.mutation.SetTitle(s)
	return uc
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (uc *UserCreate) SetNillableTitle(s *string) *UserCreate {
	if s != nil {
		uc.SetTitle(*s)
	}
	return uc
}

// SetUserName sets the "userName" field.
func (uc *UserCreate) SetUserName(s string) *UserCreate {
	uc.mutation.SetUserName(s)
	return uc
}

// SetUserType sets the "userType" field.
func (uc *UserCreate) SetUserType(s string) *UserCreate {
	uc.mutation.SetUserType(s)
	return uc
}

// SetNillableUserType sets the "userType" field if the given value is not nil.
func (uc *UserCreate) SetNillableUserType(s *string) *UserCreate {
	if s != nil {
		uc.SetUserType(*s)
	}
	return uc
}

// SetEtag sets the "etag" field.
func (uc *UserCreate) SetEtag(s string) *UserCreate {
	uc.mutation.SetEtag(s)
	return uc
}

// SetID sets the "id" field.
func (uc *UserCreate) SetID(u uuid.UUID) *UserCreate {
	uc.mutation.SetID(u)
	return uc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (uc *UserCreate) SetNillableID(u *uuid.UUID) *UserCreate {
	if u != nil {
		uc.SetID(*u)
	}
	return uc
}

// AddAddressIDs adds the "addresses" edge to the Address entity by IDs.
func (uc *UserCreate) AddAddressIDs(ids ...int) *UserCreate {
	uc.mutation.AddAddressIDs(ids...)
	return uc
}

// AddAddresses adds the "addresses" edges to the Address entity.
func (uc *UserCreate) AddAddresses(a ...*Address) *UserCreate {
	ids := make([]int, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return uc.AddAddressIDs(ids...)
}

// AddGroupIDs adds the "groups" edge to the Group entity by IDs.
func (uc *UserCreate) AddGroupIDs(ids ...uuid.UUID) *UserCreate {
	uc.mutation.AddGroupIDs(ids...)
	return uc
}

// AddGroups adds the "groups" edges to the Group entity.
func (uc *UserCreate) AddGroups(g ...*Group) *UserCreate {
	ids := make([]uuid.UUID, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uc.AddGroupIDs(ids...)
}

// AddEmailIDs adds the "emails" edge to the Email entity by IDs.
func (uc *UserCreate) AddEmailIDs(ids ...int) *UserCreate {
	uc.mutation.AddEmailIDs(ids...)
	return uc
}

// AddEmails adds the "emails" edges to the Email entity.
func (uc *UserCreate) AddEmails(e ...*Email) *UserCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return uc.AddEmailIDs(ids...)
}

// SetNameID sets the "name" edge to the Names entity by ID.
func (uc *UserCreate) SetNameID(id int) *UserCreate {
	uc.mutation.SetNameID(id)
	return uc
}

// SetNillableNameID sets the "name" edge to the Names entity by ID if the given value is not nil.
func (uc *UserCreate) SetNillableNameID(id *int) *UserCreate {
	if id != nil {
		uc = uc.SetNameID(*id)
	}
	return uc
}

// SetName sets the "name" edge to the Names entity.
func (uc *UserCreate) SetName(n *Names) *UserCreate {
	return uc.SetNameID(n.ID)
}

// AddEntitlementIDs adds the "entitlements" edge to the Entitlement entity by IDs.
func (uc *UserCreate) AddEntitlementIDs(ids ...int) *UserCreate {
	uc.mutation.AddEntitlementIDs(ids...)
	return uc
}

// AddEntitlements adds the "entitlements" edges to the Entitlement entity.
func (uc *UserCreate) AddEntitlements(e ...*Entitlement) *UserCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return uc.AddEntitlementIDs(ids...)
}

// AddRoleIDs adds the "roles" edge to the Role entity by IDs.
func (uc *UserCreate) AddRoleIDs(ids ...int) *UserCreate {
	uc.mutation.AddRoleIDs(ids...)
	return uc
}

// AddRoles adds the "roles" edges to the Role entity.
func (uc *UserCreate) AddRoles(r ...*Role) *UserCreate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return uc.AddRoleIDs(ids...)
}

// AddImseIDs adds the "imses" edge to the IMS entity by IDs.
func (uc *UserCreate) AddImseIDs(ids ...int) *UserCreate {
	uc.mutation.AddImseIDs(ids...)
	return uc
}

// AddImses adds the "imses" edges to the IMS entity.
func (uc *UserCreate) AddImses(i ...*IMS) *UserCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return uc.AddImseIDs(ids...)
}

// AddPhoneNumberIDs adds the "phone_numbers" edge to the PhoneNumber entity by IDs.
func (uc *UserCreate) AddPhoneNumberIDs(ids ...int) *UserCreate {
	uc.mutation.AddPhoneNumberIDs(ids...)
	return uc
}

// AddPhoneNumbers adds the "phone_numbers" edges to the PhoneNumber entity.
func (uc *UserCreate) AddPhoneNumbers(p ...*PhoneNumber) *UserCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uc.AddPhoneNumberIDs(ids...)
}

// AddPhotoIDs adds the "photos" edge to the Photo entity by IDs.
func (uc *UserCreate) AddPhotoIDs(ids ...int) *UserCreate {
	uc.mutation.AddPhotoIDs(ids...)
	return uc
}

// AddPhotos adds the "photos" edges to the Photo entity.
func (uc *UserCreate) AddPhotos(p ...*Photo) *UserCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return uc.AddPhotoIDs(ids...)
}

// AddX509CertificateIDs adds the "x509Certificates" edge to the X509Certificate entity by IDs.
func (uc *UserCreate) AddX509CertificateIDs(ids ...int) *UserCreate {
	uc.mutation.AddX509CertificateIDs(ids...)
	return uc
}

// AddX509Certificates adds the "x509Certificates" edges to the X509Certificate entity.
func (uc *UserCreate) AddX509Certificates(x ...*X509Certificate) *UserCreate {
	ids := make([]int, len(x))
	for i := range x {
		ids[i] = x[i].ID
	}
	return uc.AddX509CertificateIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uc *UserCreate) Mutation() *UserMutation {
	return uc.mutation
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	var (
		err  error
		node *User
	)
	uc.defaults()
	if len(uc.hooks) == 0 {
		if err = uc.check(); err != nil {
			return nil, err
		}
		node, err = uc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = uc.check(); err != nil {
				return nil, err
			}
			uc.mutation = mutation
			if node, err = uc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(uc.hooks) - 1; i >= 0; i-- {
			if uc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = uc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UserCreate) SaveX(ctx context.Context) *User {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (uc *UserCreate) Exec(ctx context.Context) error {
	_, err := uc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uc *UserCreate) ExecX(ctx context.Context) {
	if err := uc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uc *UserCreate) defaults() {
	if _, ok := uc.mutation.ID(); !ok {
		v := user.DefaultID()
		uc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uc *UserCreate) check() error {
	if v, ok := uc.mutation.Password(); ok {
		if err := user.PasswordValidator(v); err != nil {
			return &ValidationError{Name: "password", err: fmt.Errorf(`ent: validator failed for field "User.password": %w`, err)}
		}
	}
	if _, ok := uc.mutation.UserName(); !ok {
		return &ValidationError{Name: "userName", err: errors.New(`ent: missing required field "User.userName"`)}
	}
	if v, ok := uc.mutation.UserName(); ok {
		if err := user.UserNameValidator(v); err != nil {
			return &ValidationError{Name: "userName", err: fmt.Errorf(`ent: validator failed for field "User.userName": %w`, err)}
		}
	}
	if _, ok := uc.mutation.Etag(); !ok {
		return &ValidationError{Name: "etag", err: errors.New(`ent: missing required field "User.etag"`)}
	}
	if v, ok := uc.mutation.Etag(); ok {
		if err := user.EtagValidator(v); err != nil {
			return &ValidationError{Name: "etag", err: fmt.Errorf(`ent: validator failed for field "User.etag": %w`, err)}
		}
	}
	return nil
}

func (uc *UserCreate) sqlSave(ctx context.Context) (*User, error) {
	_node, _spec := uc.createSpec()
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (uc *UserCreate) createSpec() (*User, *sqlgraph.CreateSpec) {
	var (
		_node = &User{config: uc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: user.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: user.FieldID,
			},
		}
	)
	if id, ok := uc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := uc.mutation.Active(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: user.FieldActive,
		})
		_node.Active = value
	}
	if value, ok := uc.mutation.DisplayName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldDisplayName,
		})
		_node.DisplayName = value
	}
	if value, ok := uc.mutation.ExternalID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldExternalID,
		})
		_node.ExternalID = value
	}
	if value, ok := uc.mutation.Locale(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldLocale,
		})
		_node.Locale = value
	}
	if value, ok := uc.mutation.NickName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldNickName,
		})
		_node.NickName = value
	}
	if value, ok := uc.mutation.Password(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPassword,
		})
		_node.Password = value
	}
	if value, ok := uc.mutation.PreferredLanguage(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPreferredLanguage,
		})
		_node.PreferredLanguage = value
	}
	if value, ok := uc.mutation.ProfileURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldProfileURL,
		})
		_node.ProfileURL = value
	}
	if value, ok := uc.mutation.Timezone(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldTimezone,
		})
		_node.Timezone = value
	}
	if value, ok := uc.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldTitle,
		})
		_node.Title = value
	}
	if value, ok := uc.mutation.UserName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldUserName,
		})
		_node.UserName = value
	}
	if value, ok := uc.mutation.UserType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldUserType,
		})
		_node.UserType = value
	}
	if value, ok := uc.mutation.Etag(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEtag,
		})
		_node.Etag = value
	}
	if nodes := uc.mutation.AddressesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.AddressesTable,
			Columns: []string{user.AddressesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: address.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   user.GroupsTable,
			Columns: user.GroupsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: group.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.EmailsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.EmailsTable,
			Columns: []string{user.EmailsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: email.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.NameIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.NameTable,
			Columns: []string{user.NameColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: names.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.EntitlementsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.EntitlementsTable,
			Columns: []string{user.EntitlementsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: entitlement.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.RolesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.RolesTable,
			Columns: []string{user.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: role.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.ImsesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ImsesTable,
			Columns: []string{user.ImsesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: ims.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.PhoneNumbersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.PhoneNumbersTable,
			Columns: []string{user.PhoneNumbersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: phonenumber.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.PhotosIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.PhotosTable,
			Columns: []string{user.PhotosColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: photo.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.X509CertificatesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.X509CertificatesTable,
			Columns: []string{user.X509CertificatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: x509certificate.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// UserCreateBulk is the builder for creating many User entities in bulk.
type UserCreateBulk struct {
	config
	builders []*UserCreate
}

// Save creates the User entities in the database.
func (ucb *UserCreateBulk) Save(ctx context.Context) ([]*User, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ucb.builders))
	nodes := make([]*User, len(ucb.builders))
	mutators := make([]Mutator, len(ucb.builders))
	for i := range ucb.builders {
		func(i int, root context.Context) {
			builder := ucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*UserMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ucb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ucb *UserCreateBulk) SaveX(ctx context.Context) []*User {
	v, err := ucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ucb *UserCreateBulk) Exec(ctx context.Context) error {
	_, err := ucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ucb *UserCreateBulk) ExecX(ctx context.Context) {
	if err := ucb.Exec(ctx); err != nil {
		panic(err)
	}
}
