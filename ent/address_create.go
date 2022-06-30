// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/address"
)

// AddressCreate is the builder for creating a Address entity.
type AddressCreate struct {
	config
	mutation *AddressMutation
	hooks    []Hook
}

// SetCountry sets the "country" field.
func (ac *AddressCreate) SetCountry(s string) *AddressCreate {
	ac.mutation.SetCountry(s)
	return ac
}

// SetNillableCountry sets the "country" field if the given value is not nil.
func (ac *AddressCreate) SetNillableCountry(s *string) *AddressCreate {
	if s != nil {
		ac.SetCountry(*s)
	}
	return ac
}

// SetFormatted sets the "formatted" field.
func (ac *AddressCreate) SetFormatted(s string) *AddressCreate {
	ac.mutation.SetFormatted(s)
	return ac
}

// SetNillableFormatted sets the "formatted" field if the given value is not nil.
func (ac *AddressCreate) SetNillableFormatted(s *string) *AddressCreate {
	if s != nil {
		ac.SetFormatted(*s)
	}
	return ac
}

// SetLocality sets the "locality" field.
func (ac *AddressCreate) SetLocality(s string) *AddressCreate {
	ac.mutation.SetLocality(s)
	return ac
}

// SetNillableLocality sets the "locality" field if the given value is not nil.
func (ac *AddressCreate) SetNillableLocality(s *string) *AddressCreate {
	if s != nil {
		ac.SetLocality(*s)
	}
	return ac
}

// SetPostalCode sets the "postalCode" field.
func (ac *AddressCreate) SetPostalCode(s string) *AddressCreate {
	ac.mutation.SetPostalCode(s)
	return ac
}

// SetNillablePostalCode sets the "postalCode" field if the given value is not nil.
func (ac *AddressCreate) SetNillablePostalCode(s *string) *AddressCreate {
	if s != nil {
		ac.SetPostalCode(*s)
	}
	return ac
}

// SetRegion sets the "region" field.
func (ac *AddressCreate) SetRegion(s string) *AddressCreate {
	ac.mutation.SetRegion(s)
	return ac
}

// SetNillableRegion sets the "region" field if the given value is not nil.
func (ac *AddressCreate) SetNillableRegion(s *string) *AddressCreate {
	if s != nil {
		ac.SetRegion(*s)
	}
	return ac
}

// SetStreetAddress sets the "streetAddress" field.
func (ac *AddressCreate) SetStreetAddress(s string) *AddressCreate {
	ac.mutation.SetStreetAddress(s)
	return ac
}

// SetNillableStreetAddress sets the "streetAddress" field if the given value is not nil.
func (ac *AddressCreate) SetNillableStreetAddress(s *string) *AddressCreate {
	if s != nil {
		ac.SetStreetAddress(*s)
	}
	return ac
}

// Mutation returns the AddressMutation object of the builder.
func (ac *AddressCreate) Mutation() *AddressMutation {
	return ac.mutation
}

// Save creates the Address in the database.
func (ac *AddressCreate) Save(ctx context.Context) (*Address, error) {
	var (
		err  error
		node *Address
	)
	if len(ac.hooks) == 0 {
		if err = ac.check(); err != nil {
			return nil, err
		}
		node, err = ac.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AddressMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ac.check(); err != nil {
				return nil, err
			}
			ac.mutation = mutation
			if node, err = ac.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ac.hooks) - 1; i >= 0; i-- {
			if ac.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ac.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ac.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AddressCreate) SaveX(ctx context.Context) *Address {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *AddressCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *AddressCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *AddressCreate) check() error {
	return nil
}

func (ac *AddressCreate) sqlSave(ctx context.Context) (*Address, error) {
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (ac *AddressCreate) createSpec() (*Address, *sqlgraph.CreateSpec) {
	var (
		_node = &Address{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: address.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: address.FieldID,
			},
		}
	)
	if value, ok := ac.mutation.Country(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: address.FieldCountry,
		})
		_node.Country = value
	}
	if value, ok := ac.mutation.Formatted(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: address.FieldFormatted,
		})
		_node.Formatted = value
	}
	if value, ok := ac.mutation.Locality(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: address.FieldLocality,
		})
		_node.Locality = value
	}
	if value, ok := ac.mutation.PostalCode(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: address.FieldPostalCode,
		})
		_node.PostalCode = value
	}
	if value, ok := ac.mutation.Region(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: address.FieldRegion,
		})
		_node.Region = value
	}
	if value, ok := ac.mutation.StreetAddress(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: address.FieldStreetAddress,
		})
		_node.StreetAddress = value
	}
	return _node, _spec
}

// AddressCreateBulk is the builder for creating many Address entities in bulk.
type AddressCreateBulk struct {
	config
	builders []*AddressCreate
}

// Save creates the Address entities in the database.
func (acb *AddressCreateBulk) Save(ctx context.Context) ([]*Address, error) {
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Address, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AddressMutation)
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
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
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
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *AddressCreateBulk) SaveX(ctx context.Context) []*Address {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *AddressCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *AddressCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}
