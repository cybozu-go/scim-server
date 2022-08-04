// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/ims"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/google/uuid"
)

// IMSCreate is the builder for creating a IMS entity.
type IMSCreate struct {
	config
	mutation *IMSMutation
	hooks    []Hook
}

// SetDisplay sets the "display" field.
func (ic *IMSCreate) SetDisplay(s string) *IMSCreate {
	ic.mutation.SetDisplay(s)
	return ic
}

// SetNillableDisplay sets the "display" field if the given value is not nil.
func (ic *IMSCreate) SetNillableDisplay(s *string) *IMSCreate {
	if s != nil {
		ic.SetDisplay(*s)
	}
	return ic
}

// SetPrimary sets the "primary" field.
func (ic *IMSCreate) SetPrimary(b bool) *IMSCreate {
	ic.mutation.SetPrimary(b)
	return ic
}

// SetNillablePrimary sets the "primary" field if the given value is not nil.
func (ic *IMSCreate) SetNillablePrimary(b *bool) *IMSCreate {
	if b != nil {
		ic.SetPrimary(*b)
	}
	return ic
}

// SetType sets the "type" field.
func (ic *IMSCreate) SetType(s string) *IMSCreate {
	ic.mutation.SetType(s)
	return ic
}

// SetNillableType sets the "type" field if the given value is not nil.
func (ic *IMSCreate) SetNillableType(s *string) *IMSCreate {
	if s != nil {
		ic.SetType(*s)
	}
	return ic
}

// SetValue sets the "value" field.
func (ic *IMSCreate) SetValue(s string) *IMSCreate {
	ic.mutation.SetValue(s)
	return ic
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (ic *IMSCreate) SetNillableValue(s *string) *IMSCreate {
	if s != nil {
		ic.SetValue(*s)
	}
	return ic
}

// SetID sets the "id" field.
func (ic *IMSCreate) SetID(u uuid.UUID) *IMSCreate {
	ic.mutation.SetID(u)
	return ic
}

// SetNillableID sets the "id" field if the given value is not nil.
func (ic *IMSCreate) SetNillableID(u *uuid.UUID) *IMSCreate {
	if u != nil {
		ic.SetID(*u)
	}
	return ic
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ic *IMSCreate) SetUserID(id uuid.UUID) *IMSCreate {
	ic.mutation.SetUserID(id)
	return ic
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ic *IMSCreate) SetNillableUserID(id *uuid.UUID) *IMSCreate {
	if id != nil {
		ic = ic.SetUserID(*id)
	}
	return ic
}

// SetUser sets the "user" edge to the User entity.
func (ic *IMSCreate) SetUser(u *User) *IMSCreate {
	return ic.SetUserID(u.ID)
}

// Mutation returns the IMSMutation object of the builder.
func (ic *IMSCreate) Mutation() *IMSMutation {
	return ic.mutation
}

// Save creates the IMS in the database.
func (ic *IMSCreate) Save(ctx context.Context) (*IMS, error) {
	var (
		err  error
		node *IMS
	)
	ic.defaults()
	if len(ic.hooks) == 0 {
		if err = ic.check(); err != nil {
			return nil, err
		}
		node, err = ic.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*IMSMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ic.check(); err != nil {
				return nil, err
			}
			ic.mutation = mutation
			if node, err = ic.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ic.hooks) - 1; i >= 0; i-- {
			if ic.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ic.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, ic.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*IMS)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from IMSMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ic *IMSCreate) SaveX(ctx context.Context) *IMS {
	v, err := ic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ic *IMSCreate) Exec(ctx context.Context) error {
	_, err := ic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ic *IMSCreate) ExecX(ctx context.Context) {
	if err := ic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ic *IMSCreate) defaults() {
	if _, ok := ic.mutation.ID(); !ok {
		v := ims.DefaultID()
		ic.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ic *IMSCreate) check() error {
	return nil
}

func (ic *IMSCreate) sqlSave(ctx context.Context) (*IMS, error) {
	_node, _spec := ic.createSpec()
	if err := sqlgraph.CreateNode(ctx, ic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
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

func (ic *IMSCreate) createSpec() (*IMS, *sqlgraph.CreateSpec) {
	var (
		_node = &IMS{config: ic.config}
		_spec = &sqlgraph.CreateSpec{
			Table: ims.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: ims.FieldID,
			},
		}
	)
	if id, ok := ic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := ic.mutation.Display(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ims.FieldDisplay,
		})
		_node.Display = value
	}
	if value, ok := ic.mutation.Primary(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: ims.FieldPrimary,
		})
		_node.Primary = value
	}
	if value, ok := ic.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ims.FieldType,
		})
		_node.Type = value
	}
	if value, ok := ic.mutation.Value(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: ims.FieldValue,
		})
		_node.Value = value
	}
	if nodes := ic.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   ims.UserTable,
			Columns: []string{ims.UserColumn},
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
		_node.user_ims = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// IMSCreateBulk is the builder for creating many IMS entities in bulk.
type IMSCreateBulk struct {
	config
	builders []*IMSCreate
}

// Save creates the IMS entities in the database.
func (icb *IMSCreateBulk) Save(ctx context.Context) ([]*IMS, error) {
	specs := make([]*sqlgraph.CreateSpec, len(icb.builders))
	nodes := make([]*IMS, len(icb.builders))
	mutators := make([]Mutator, len(icb.builders))
	for i := range icb.builders {
		func(i int, root context.Context) {
			builder := icb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*IMSMutation)
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
					_, err = mutators[i+1].Mutate(root, icb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, icb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
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
		if _, err := mutators[0].Mutate(ctx, icb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (icb *IMSCreateBulk) SaveX(ctx context.Context) []*IMS {
	v, err := icb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (icb *IMSCreateBulk) Exec(ctx context.Context) error {
	_, err := icb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (icb *IMSCreateBulk) ExecX(ctx context.Context) {
	if err := icb.Exec(ctx); err != nil {
		panic(err)
	}
}
