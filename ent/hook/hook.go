// Code generated by ent, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/cybozu-go/scim-server/ent"
)

// The AddressFunc type is an adapter to allow the use of ordinary
// function as Address mutator.
type AddressFunc func(context.Context, *ent.AddressMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f AddressFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.AddressMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.AddressMutation", m)
	}
	return f(ctx, mv)
}

// The EmailFunc type is an adapter to allow the use of ordinary
// function as Email mutator.
type EmailFunc func(context.Context, *ent.EmailMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EmailFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EmailMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EmailMutation", m)
	}
	return f(ctx, mv)
}

// The EntitlementFunc type is an adapter to allow the use of ordinary
// function as Entitlement mutator.
type EntitlementFunc func(context.Context, *ent.EntitlementMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EntitlementFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EntitlementMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EntitlementMutation", m)
	}
	return f(ctx, mv)
}

// The GroupFunc type is an adapter to allow the use of ordinary
// function as Group mutator.
type GroupFunc func(context.Context, *ent.GroupMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f GroupFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.GroupMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.GroupMutation", m)
	}
	return f(ctx, mv)
}

// The IMSFunc type is an adapter to allow the use of ordinary
// function as IMS mutator.
type IMSFunc func(context.Context, *ent.IMSMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f IMSFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.IMSMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.IMSMutation", m)
	}
	return f(ctx, mv)
}

// The MemberFunc type is an adapter to allow the use of ordinary
// function as Member mutator.
type MemberFunc func(context.Context, *ent.MemberMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f MemberFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.MemberMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.MemberMutation", m)
	}
	return f(ctx, mv)
}

// The NamesFunc type is an adapter to allow the use of ordinary
// function as Names mutator.
type NamesFunc func(context.Context, *ent.NamesMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f NamesFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.NamesMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.NamesMutation", m)
	}
	return f(ctx, mv)
}

// The PhoneNumberFunc type is an adapter to allow the use of ordinary
// function as PhoneNumber mutator.
type PhoneNumberFunc func(context.Context, *ent.PhoneNumberMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PhoneNumberFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.PhoneNumberMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PhoneNumberMutation", m)
	}
	return f(ctx, mv)
}

// The PhotoFunc type is an adapter to allow the use of ordinary
// function as Photo mutator.
type PhotoFunc func(context.Context, *ent.PhotoMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PhotoFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.PhotoMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PhotoMutation", m)
	}
	return f(ctx, mv)
}

// The RoleFunc type is an adapter to allow the use of ordinary
// function as Role mutator.
type RoleFunc func(context.Context, *ent.RoleMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f RoleFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.RoleMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.RoleMutation", m)
	}
	return f(ctx, mv)
}

// The UserFunc type is an adapter to allow the use of ordinary
// function as User mutator.
type UserFunc func(context.Context, *ent.UserMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f UserFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.UserMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.UserMutation", m)
	}
	return f(ctx, mv)
}

// The X509CertificateFunc type is an adapter to allow the use of ordinary
// function as X509Certificate mutator.
type X509CertificateFunc func(context.Context, *ent.X509CertificateMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f X509CertificateFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.X509CertificateMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.X509CertificateMutation", m)
	}
	return f(ctx, mv)
}

// Condition is a hook condition function.
type Condition func(context.Context, ent.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op ent.Op) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
//
func If(hk ent.Hook, cond Condition) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
//
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, ent.Update|ent.UpdateOne)
//
func Unless(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) ent.Hook {
	return func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
//
func Reject(op ent.Op) ent.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []ent.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...ent.Hook) Chain {
	return Chain{append([]ent.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() ent.Hook {
	return func(mutator ent.Mutator) ent.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...ent.Hook) Chain {
	newHooks := make([]ent.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
