// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/member"
	"github.com/google/uuid"
)

// MemberCreate is the builder for creating a Member entity.
type MemberCreate struct {
	config
	mutation *MemberMutation
	hooks    []Hook
}

// SetValue sets the "value" field.
func (mc *MemberCreate) SetValue(s string) *MemberCreate {
	mc.mutation.SetValue(s)
	return mc
}

// SetType sets the "type" field.
func (mc *MemberCreate) SetType(s string) *MemberCreate {
	mc.mutation.SetType(s)
	return mc
}

// SetRef sets the "ref" field.
func (mc *MemberCreate) SetRef(s string) *MemberCreate {
	mc.mutation.SetRef(s)
	return mc
}

// SetNillableRef sets the "ref" field if the given value is not nil.
func (mc *MemberCreate) SetNillableRef(s *string) *MemberCreate {
	if s != nil {
		mc.SetRef(*s)
	}
	return mc
}

// SetGroupID sets the "group" edge to the Group entity by ID.
func (mc *MemberCreate) SetGroupID(id uuid.UUID) *MemberCreate {
	mc.mutation.SetGroupID(id)
	return mc
}

// SetNillableGroupID sets the "group" edge to the Group entity by ID if the given value is not nil.
func (mc *MemberCreate) SetNillableGroupID(id *uuid.UUID) *MemberCreate {
	if id != nil {
		mc = mc.SetGroupID(*id)
	}
	return mc
}

// SetGroup sets the "group" edge to the Group entity.
func (mc *MemberCreate) SetGroup(g *Group) *MemberCreate {
	return mc.SetGroupID(g.ID)
}

// Mutation returns the MemberMutation object of the builder.
func (mc *MemberCreate) Mutation() *MemberMutation {
	return mc.mutation
}

// Save creates the Member in the database.
func (mc *MemberCreate) Save(ctx context.Context) (*Member, error) {
	var (
		err  error
		node *Member
	)
	if len(mc.hooks) == 0 {
		if err = mc.check(); err != nil {
			return nil, err
		}
		node, err = mc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MemberMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = mc.check(); err != nil {
				return nil, err
			}
			mc.mutation = mutation
			if node, err = mc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(mc.hooks) - 1; i >= 0; i-- {
			if mc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, mc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Member)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from MemberMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MemberCreate) SaveX(ctx context.Context) *Member {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MemberCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MemberCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MemberCreate) check() error {
	if _, ok := mc.mutation.Value(); !ok {
		return &ValidationError{Name: "value", err: errors.New(`ent: missing required field "Member.value"`)}
	}
	if _, ok := mc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Member.type"`)}
	}
	return nil
}

func (mc *MemberCreate) sqlSave(ctx context.Context) (*Member, error) {
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (mc *MemberCreate) createSpec() (*Member, *sqlgraph.CreateSpec) {
	var (
		_node = &Member{config: mc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: member.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: member.FieldID,
			},
		}
	)
	if value, ok := mc.mutation.Value(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: member.FieldValue,
		})
		_node.Value = value
	}
	if value, ok := mc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: member.FieldType,
		})
		_node.Type = value
	}
	if value, ok := mc.mutation.Ref(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: member.FieldRef,
		})
		_node.Ref = value
	}
	if nodes := mc.mutation.GroupIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   member.GroupTable,
			Columns: []string{member.GroupColumn},
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
		_node.group_members = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MemberCreateBulk is the builder for creating many Member entities in bulk.
type MemberCreateBulk struct {
	config
	builders []*MemberCreate
}

// Save creates the Member entities in the database.
func (mcb *MemberCreateBulk) Save(ctx context.Context) ([]*Member, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Member, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MemberMutation)
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
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MemberCreateBulk) SaveX(ctx context.Context) []*Member {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MemberCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MemberCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}