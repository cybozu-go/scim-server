// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/phonenumber"
	"github.com/cybozu-go/scim-server/ent/predicate"
)

// PhoneNumberUpdate is the builder for updating PhoneNumber entities.
type PhoneNumberUpdate struct {
	config
	hooks    []Hook
	mutation *PhoneNumberMutation
}

// Where appends a list predicates to the PhoneNumberUpdate builder.
func (pnu *PhoneNumberUpdate) Where(ps ...predicate.PhoneNumber) *PhoneNumberUpdate {
	pnu.mutation.Where(ps...)
	return pnu
}

// SetDisplay sets the "display" field.
func (pnu *PhoneNumberUpdate) SetDisplay(s string) *PhoneNumberUpdate {
	pnu.mutation.SetDisplay(s)
	return pnu
}

// SetNillableDisplay sets the "display" field if the given value is not nil.
func (pnu *PhoneNumberUpdate) SetNillableDisplay(s *string) *PhoneNumberUpdate {
	if s != nil {
		pnu.SetDisplay(*s)
	}
	return pnu
}

// ClearDisplay clears the value of the "display" field.
func (pnu *PhoneNumberUpdate) ClearDisplay() *PhoneNumberUpdate {
	pnu.mutation.ClearDisplay()
	return pnu
}

// SetPrimary sets the "primary" field.
func (pnu *PhoneNumberUpdate) SetPrimary(b bool) *PhoneNumberUpdate {
	pnu.mutation.SetPrimary(b)
	return pnu
}

// SetNillablePrimary sets the "primary" field if the given value is not nil.
func (pnu *PhoneNumberUpdate) SetNillablePrimary(b *bool) *PhoneNumberUpdate {
	if b != nil {
		pnu.SetPrimary(*b)
	}
	return pnu
}

// ClearPrimary clears the value of the "primary" field.
func (pnu *PhoneNumberUpdate) ClearPrimary() *PhoneNumberUpdate {
	pnu.mutation.ClearPrimary()
	return pnu
}

// SetType sets the "type" field.
func (pnu *PhoneNumberUpdate) SetType(s string) *PhoneNumberUpdate {
	pnu.mutation.SetType(s)
	return pnu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (pnu *PhoneNumberUpdate) SetNillableType(s *string) *PhoneNumberUpdate {
	if s != nil {
		pnu.SetType(*s)
	}
	return pnu
}

// ClearType clears the value of the "type" field.
func (pnu *PhoneNumberUpdate) ClearType() *PhoneNumberUpdate {
	pnu.mutation.ClearType()
	return pnu
}

// SetValue sets the "value" field.
func (pnu *PhoneNumberUpdate) SetValue(s string) *PhoneNumberUpdate {
	pnu.mutation.SetValue(s)
	return pnu
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (pnu *PhoneNumberUpdate) SetNillableValue(s *string) *PhoneNumberUpdate {
	if s != nil {
		pnu.SetValue(*s)
	}
	return pnu
}

// ClearValue clears the value of the "value" field.
func (pnu *PhoneNumberUpdate) ClearValue() *PhoneNumberUpdate {
	pnu.mutation.ClearValue()
	return pnu
}

// Mutation returns the PhoneNumberMutation object of the builder.
func (pnu *PhoneNumberUpdate) Mutation() *PhoneNumberMutation {
	return pnu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pnu *PhoneNumberUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(pnu.hooks) == 0 {
		affected, err = pnu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PhoneNumberMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pnu.mutation = mutation
			affected, err = pnu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pnu.hooks) - 1; i >= 0; i-- {
			if pnu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pnu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pnu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pnu *PhoneNumberUpdate) SaveX(ctx context.Context) int {
	affected, err := pnu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pnu *PhoneNumberUpdate) Exec(ctx context.Context) error {
	_, err := pnu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pnu *PhoneNumberUpdate) ExecX(ctx context.Context) {
	if err := pnu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pnu *PhoneNumberUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   phonenumber.Table,
			Columns: phonenumber.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: phonenumber.FieldID,
			},
		},
	}
	if ps := pnu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pnu.mutation.Display(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: phonenumber.FieldDisplay,
		})
	}
	if pnu.mutation.DisplayCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: phonenumber.FieldDisplay,
		})
	}
	if value, ok := pnu.mutation.Primary(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: phonenumber.FieldPrimary,
		})
	}
	if pnu.mutation.PrimaryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: phonenumber.FieldPrimary,
		})
	}
	if value, ok := pnu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: phonenumber.FieldType,
		})
	}
	if pnu.mutation.TypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: phonenumber.FieldType,
		})
	}
	if value, ok := pnu.mutation.Value(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: phonenumber.FieldValue,
		})
	}
	if pnu.mutation.ValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: phonenumber.FieldValue,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pnu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{phonenumber.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// PhoneNumberUpdateOne is the builder for updating a single PhoneNumber entity.
type PhoneNumberUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PhoneNumberMutation
}

// SetDisplay sets the "display" field.
func (pnuo *PhoneNumberUpdateOne) SetDisplay(s string) *PhoneNumberUpdateOne {
	pnuo.mutation.SetDisplay(s)
	return pnuo
}

// SetNillableDisplay sets the "display" field if the given value is not nil.
func (pnuo *PhoneNumberUpdateOne) SetNillableDisplay(s *string) *PhoneNumberUpdateOne {
	if s != nil {
		pnuo.SetDisplay(*s)
	}
	return pnuo
}

// ClearDisplay clears the value of the "display" field.
func (pnuo *PhoneNumberUpdateOne) ClearDisplay() *PhoneNumberUpdateOne {
	pnuo.mutation.ClearDisplay()
	return pnuo
}

// SetPrimary sets the "primary" field.
func (pnuo *PhoneNumberUpdateOne) SetPrimary(b bool) *PhoneNumberUpdateOne {
	pnuo.mutation.SetPrimary(b)
	return pnuo
}

// SetNillablePrimary sets the "primary" field if the given value is not nil.
func (pnuo *PhoneNumberUpdateOne) SetNillablePrimary(b *bool) *PhoneNumberUpdateOne {
	if b != nil {
		pnuo.SetPrimary(*b)
	}
	return pnuo
}

// ClearPrimary clears the value of the "primary" field.
func (pnuo *PhoneNumberUpdateOne) ClearPrimary() *PhoneNumberUpdateOne {
	pnuo.mutation.ClearPrimary()
	return pnuo
}

// SetType sets the "type" field.
func (pnuo *PhoneNumberUpdateOne) SetType(s string) *PhoneNumberUpdateOne {
	pnuo.mutation.SetType(s)
	return pnuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (pnuo *PhoneNumberUpdateOne) SetNillableType(s *string) *PhoneNumberUpdateOne {
	if s != nil {
		pnuo.SetType(*s)
	}
	return pnuo
}

// ClearType clears the value of the "type" field.
func (pnuo *PhoneNumberUpdateOne) ClearType() *PhoneNumberUpdateOne {
	pnuo.mutation.ClearType()
	return pnuo
}

// SetValue sets the "value" field.
func (pnuo *PhoneNumberUpdateOne) SetValue(s string) *PhoneNumberUpdateOne {
	pnuo.mutation.SetValue(s)
	return pnuo
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (pnuo *PhoneNumberUpdateOne) SetNillableValue(s *string) *PhoneNumberUpdateOne {
	if s != nil {
		pnuo.SetValue(*s)
	}
	return pnuo
}

// ClearValue clears the value of the "value" field.
func (pnuo *PhoneNumberUpdateOne) ClearValue() *PhoneNumberUpdateOne {
	pnuo.mutation.ClearValue()
	return pnuo
}

// Mutation returns the PhoneNumberMutation object of the builder.
func (pnuo *PhoneNumberUpdateOne) Mutation() *PhoneNumberMutation {
	return pnuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pnuo *PhoneNumberUpdateOne) Select(field string, fields ...string) *PhoneNumberUpdateOne {
	pnuo.fields = append([]string{field}, fields...)
	return pnuo
}

// Save executes the query and returns the updated PhoneNumber entity.
func (pnuo *PhoneNumberUpdateOne) Save(ctx context.Context) (*PhoneNumber, error) {
	var (
		err  error
		node *PhoneNumber
	)
	if len(pnuo.hooks) == 0 {
		node, err = pnuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PhoneNumberMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pnuo.mutation = mutation
			node, err = pnuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pnuo.hooks) - 1; i >= 0; i-- {
			if pnuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pnuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pnuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (pnuo *PhoneNumberUpdateOne) SaveX(ctx context.Context) *PhoneNumber {
	node, err := pnuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pnuo *PhoneNumberUpdateOne) Exec(ctx context.Context) error {
	_, err := pnuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pnuo *PhoneNumberUpdateOne) ExecX(ctx context.Context) {
	if err := pnuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pnuo *PhoneNumberUpdateOne) sqlSave(ctx context.Context) (_node *PhoneNumber, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   phonenumber.Table,
			Columns: phonenumber.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: phonenumber.FieldID,
			},
		},
	}
	id, ok := pnuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PhoneNumber.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pnuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, phonenumber.FieldID)
		for _, f := range fields {
			if !phonenumber.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != phonenumber.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pnuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pnuo.mutation.Display(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: phonenumber.FieldDisplay,
		})
	}
	if pnuo.mutation.DisplayCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: phonenumber.FieldDisplay,
		})
	}
	if value, ok := pnuo.mutation.Primary(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: phonenumber.FieldPrimary,
		})
	}
	if pnuo.mutation.PrimaryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Column: phonenumber.FieldPrimary,
		})
	}
	if value, ok := pnuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: phonenumber.FieldType,
		})
	}
	if pnuo.mutation.TypeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: phonenumber.FieldType,
		})
	}
	if value, ok := pnuo.mutation.Value(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: phonenumber.FieldValue,
		})
	}
	if pnuo.mutation.ValueCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: phonenumber.FieldValue,
		})
	}
	_node = &PhoneNumber{config: pnuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pnuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{phonenumber.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
