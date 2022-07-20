// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/email"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/google/uuid"
)

// EmailUpdate is the builder for updating Email entities.
type EmailUpdate struct {
	config
	hooks    []Hook
	mutation *EmailMutation
}

// Where appends a list predicates to the EmailUpdate builder.
func (eu *EmailUpdate) Where(ps ...predicate.Email) *EmailUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetDisplay sets the "display" field.
func (eu *EmailUpdate) SetDisplay(s string) *EmailUpdate {
	eu.mutation.SetDisplay(s)
	return eu
}

// SetNillableDisplay sets the "display" field if the given value is not nil.
func (eu *EmailUpdate) SetNillableDisplay(s *string) *EmailUpdate {
	if s != nil {
		eu.SetDisplay(*s)
	}
	return eu
}

// ClearDisplay clears the value of the "display" field.
func (eu *EmailUpdate) ClearDisplay() *EmailUpdate {
	eu.mutation.ClearDisplay()
	return eu
}

// SetPrimary sets the "primary" field.
func (eu *EmailUpdate) SetPrimary(b bool) *EmailUpdate {
	eu.mutation.SetPrimary(b)
	return eu
}

// SetNillablePrimary sets the "primary" field if the given value is not nil.
func (eu *EmailUpdate) SetNillablePrimary(b *bool) *EmailUpdate {
	if b != nil {
		eu.SetPrimary(*b)
	}
	return eu
}

// ClearPrimary clears the value of the "primary" field.
func (eu *EmailUpdate) ClearPrimary() *EmailUpdate {
	eu.mutation.ClearPrimary()
	return eu
}

// SetType sets the "type" field.
func (eu *EmailUpdate) SetType(s string) *EmailUpdate {
	eu.mutation.SetType(s)
	return eu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (eu *EmailUpdate) SetNillableType(s *string) *EmailUpdate {
	if s != nil {
		eu.SetType(*s)
	}
	return eu
}

// ClearType clears the value of the "type" field.
func (eu *EmailUpdate) ClearType() *EmailUpdate {
	eu.mutation.ClearType()
	return eu
}

// SetValue sets the "value" field.
func (eu *EmailUpdate) SetValue(s string) *EmailUpdate {
	eu.mutation.SetValue(s)
	return eu
}

// SetUserID sets the "user" edge to the User entity by ID.
func (eu *EmailUpdate) SetUserID(id uuid.UUID) *EmailUpdate {
	eu.mutation.SetUserID(id)
	return eu
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (eu *EmailUpdate) SetNillableUserID(id *uuid.UUID) *EmailUpdate {
	if id != nil {
		eu = eu.SetUserID(*id)
	}
	return eu
}

// SetUser sets the "user" edge to the User entity.
func (eu *EmailUpdate) SetUser(u *User) *EmailUpdate {
	return eu.SetUserID(u.ID)
}

// Mutation returns the EmailMutation object of the builder.
func (eu *EmailUpdate) Mutation() *EmailMutation {
	return eu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (eu *EmailUpdate) ClearUser() *EmailUpdate {
	eu.mutation.ClearUser()
	return eu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EmailUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(eu.hooks) == 0 {
		affected, err = eu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EmailMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			eu.mutation = mutation
			affected, err = eu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(eu.hooks) - 1; i >= 0; i-- {
			if eu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = eu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, eu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EmailUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EmailUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EmailUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (eu *EmailUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   email.Table,
			Columns: email.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: email.FieldID,
			},
		},
	}
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.Display(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: email.FieldDisplay,
		})
	}
	if eu.mutation.DisplayCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: email.FieldDisplay,
		})
	}
	if value, ok := eu.mutation.Primary(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: email.FieldPrimary,
		})
	}
	if eu.mutation.PrimaryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: email.FieldPrimary,
		})
	}
	if value, ok := eu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: email.FieldType,
		})
	}
	if eu.mutation.TypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: email.FieldType,
		})
	}
	if value, ok := eu.mutation.Value(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: email.FieldValue,
		})
	}
	if eu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.UserTable,
			Columns: []string{email.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.UserTable,
			Columns: []string{email.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{email.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// EmailUpdateOne is the builder for updating a single Email entity.
type EmailUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *EmailMutation
}

// SetDisplay sets the "display" field.
func (euo *EmailUpdateOne) SetDisplay(s string) *EmailUpdateOne {
	euo.mutation.SetDisplay(s)
	return euo
}

// SetNillableDisplay sets the "display" field if the given value is not nil.
func (euo *EmailUpdateOne) SetNillableDisplay(s *string) *EmailUpdateOne {
	if s != nil {
		euo.SetDisplay(*s)
	}
	return euo
}

// ClearDisplay clears the value of the "display" field.
func (euo *EmailUpdateOne) ClearDisplay() *EmailUpdateOne {
	euo.mutation.ClearDisplay()
	return euo
}

// SetPrimary sets the "primary" field.
func (euo *EmailUpdateOne) SetPrimary(b bool) *EmailUpdateOne {
	euo.mutation.SetPrimary(b)
	return euo
}

// SetNillablePrimary sets the "primary" field if the given value is not nil.
func (euo *EmailUpdateOne) SetNillablePrimary(b *bool) *EmailUpdateOne {
	if b != nil {
		euo.SetPrimary(*b)
	}
	return euo
}

// ClearPrimary clears the value of the "primary" field.
func (euo *EmailUpdateOne) ClearPrimary() *EmailUpdateOne {
	euo.mutation.ClearPrimary()
	return euo
}

// SetType sets the "type" field.
func (euo *EmailUpdateOne) SetType(s string) *EmailUpdateOne {
	euo.mutation.SetType(s)
	return euo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (euo *EmailUpdateOne) SetNillableType(s *string) *EmailUpdateOne {
	if s != nil {
		euo.SetType(*s)
	}
	return euo
}

// ClearType clears the value of the "type" field.
func (euo *EmailUpdateOne) ClearType() *EmailUpdateOne {
	euo.mutation.ClearType()
	return euo
}

// SetValue sets the "value" field.
func (euo *EmailUpdateOne) SetValue(s string) *EmailUpdateOne {
	euo.mutation.SetValue(s)
	return euo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (euo *EmailUpdateOne) SetUserID(id uuid.UUID) *EmailUpdateOne {
	euo.mutation.SetUserID(id)
	return euo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (euo *EmailUpdateOne) SetNillableUserID(id *uuid.UUID) *EmailUpdateOne {
	if id != nil {
		euo = euo.SetUserID(*id)
	}
	return euo
}

// SetUser sets the "user" edge to the User entity.
func (euo *EmailUpdateOne) SetUser(u *User) *EmailUpdateOne {
	return euo.SetUserID(u.ID)
}

// Mutation returns the EmailMutation object of the builder.
func (euo *EmailUpdateOne) Mutation() *EmailMutation {
	return euo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (euo *EmailUpdateOne) ClearUser() *EmailUpdateOne {
	euo.mutation.ClearUser()
	return euo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EmailUpdateOne) Select(field string, fields ...string) *EmailUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Email entity.
func (euo *EmailUpdateOne) Save(ctx context.Context) (*Email, error) {
	var (
		err  error
		node *Email
	)
	if len(euo.hooks) == 0 {
		node, err = euo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*EmailMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			euo.mutation = mutation
			node, err = euo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(euo.hooks) - 1; i >= 0; i-- {
			if euo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = euo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, euo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EmailUpdateOne) SaveX(ctx context.Context) *Email {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EmailUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EmailUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (euo *EmailUpdateOne) sqlSave(ctx context.Context) (_node *Email, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   email.Table,
			Columns: email.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: email.FieldID,
			},
		},
	}
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Email.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, email.FieldID)
		for _, f := range fields {
			if !email.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != email.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.Display(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: email.FieldDisplay,
		})
	}
	if euo.mutation.DisplayCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: email.FieldDisplay,
		})
	}
	if value, ok := euo.mutation.Primary(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: email.FieldPrimary,
		})
	}
	if euo.mutation.PrimaryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: email.FieldPrimary,
		})
	}
	if value, ok := euo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: email.FieldType,
		})
	}
	if euo.mutation.TypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: email.FieldType,
		})
	}
	if value, ok := euo.mutation.Value(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: email.FieldValue,
		})
	}
	if euo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.UserTable,
			Columns: []string{email.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   email.UserTable,
			Columns: []string{email.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Email{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{email.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
