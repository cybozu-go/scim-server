// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/entitlement"
	"github.com/cybozu-go/scim-server/ent/predicate"
)

// EntitlementQuery is the builder for querying Entitlement entities.
type EntitlementQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.Entitlement
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the EntitlementQuery builder.
func (eq *EntitlementQuery) Where(ps ...predicate.Entitlement) *EntitlementQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit adds a limit step to the query.
func (eq *EntitlementQuery) Limit(limit int) *EntitlementQuery {
	eq.limit = &limit
	return eq
}

// Offset adds an offset step to the query.
func (eq *EntitlementQuery) Offset(offset int) *EntitlementQuery {
	eq.offset = &offset
	return eq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (eq *EntitlementQuery) Unique(unique bool) *EntitlementQuery {
	eq.unique = &unique
	return eq
}

// Order adds an order step to the query.
func (eq *EntitlementQuery) Order(o ...OrderFunc) *EntitlementQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// First returns the first Entitlement entity from the query.
// Returns a *NotFoundError when no Entitlement was found.
func (eq *EntitlementQuery) First(ctx context.Context) (*Entitlement, error) {
	nodes, err := eq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{entitlement.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *EntitlementQuery) FirstX(ctx context.Context) *Entitlement {
	node, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Entitlement ID from the query.
// Returns a *NotFoundError when no Entitlement ID was found.
func (eq *EntitlementQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{entitlement.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (eq *EntitlementQuery) FirstIDX(ctx context.Context) int {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Entitlement entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Entitlement entity is found.
// Returns a *NotFoundError when no Entitlement entities are found.
func (eq *EntitlementQuery) Only(ctx context.Context) (*Entitlement, error) {
	nodes, err := eq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{entitlement.Label}
	default:
		return nil, &NotSingularError{entitlement.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *EntitlementQuery) OnlyX(ctx context.Context) *Entitlement {
	node, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Entitlement ID in the query.
// Returns a *NotSingularError when more than one Entitlement ID is found.
// Returns a *NotFoundError when no entities are found.
func (eq *EntitlementQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = &NotSingularError{entitlement.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (eq *EntitlementQuery) OnlyIDX(ctx context.Context) int {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Entitlements.
func (eq *EntitlementQuery) All(ctx context.Context) ([]*Entitlement, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return eq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (eq *EntitlementQuery) AllX(ctx context.Context) []*Entitlement {
	nodes, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Entitlement IDs.
func (eq *EntitlementQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := eq.Select(entitlement.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *EntitlementQuery) IDsX(ctx context.Context) []int {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *EntitlementQuery) Count(ctx context.Context) (int, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return eq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (eq *EntitlementQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *EntitlementQuery) Exist(ctx context.Context) (bool, error) {
	if err := eq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return eq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *EntitlementQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the EntitlementQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *EntitlementQuery) Clone() *EntitlementQuery {
	if eq == nil {
		return nil
	}
	return &EntitlementQuery{
		config:     eq.config,
		limit:      eq.limit,
		offset:     eq.offset,
		order:      append([]OrderFunc{}, eq.order...),
		predicates: append([]predicate.Entitlement{}, eq.predicates...),
		// clone intermediate query.
		sql:    eq.sql.Clone(),
		path:   eq.path,
		unique: eq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Display string `json:"display,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Entitlement.Query().
//		GroupBy(entitlement.FieldDisplay).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (eq *EntitlementQuery) GroupBy(field string, fields ...string) *EntitlementGroupBy {
	group := &EntitlementGroupBy{config: eq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return eq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Display string `json:"display,omitempty"`
//	}
//
//	client.Entitlement.Query().
//		Select(entitlement.FieldDisplay).
//		Scan(ctx, &v)
//
func (eq *EntitlementQuery) Select(fields ...string) *EntitlementSelect {
	eq.fields = append(eq.fields, fields...)
	return &EntitlementSelect{EntitlementQuery: eq}
}

func (eq *EntitlementQuery) prepareQuery(ctx context.Context) error {
	for _, f := range eq.fields {
		if !entitlement.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if eq.path != nil {
		prev, err := eq.path(ctx)
		if err != nil {
			return err
		}
		eq.sql = prev
	}
	return nil
}

func (eq *EntitlementQuery) sqlAll(ctx context.Context) ([]*Entitlement, error) {
	var (
		nodes   = []*Entitlement{}
		withFKs = eq.withFKs
		_spec   = eq.querySpec()
	)
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, entitlement.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &Entitlement{config: eq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, eq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (eq *EntitlementQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eq.querySpec()
	_spec.Node.Columns = eq.fields
	if len(eq.fields) > 0 {
		_spec.Unique = eq.unique != nil && *eq.unique
	}
	return sqlgraph.CountNodes(ctx, eq.driver, _spec)
}

func (eq *EntitlementQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := eq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (eq *EntitlementQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   entitlement.Table,
			Columns: entitlement.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: entitlement.FieldID,
			},
		},
		From:   eq.sql,
		Unique: true,
	}
	if unique := eq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := eq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, entitlement.FieldID)
		for i := range fields {
			if fields[i] != entitlement.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := eq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eq *EntitlementQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(entitlement.Table)
	columns := eq.fields
	if len(columns) == 0 {
		columns = entitlement.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if eq.unique != nil && *eq.unique {
		selector.Distinct()
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// EntitlementGroupBy is the group-by builder for Entitlement entities.
type EntitlementGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *EntitlementGroupBy) Aggregate(fns ...AggregateFunc) *EntitlementGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the group-by query and scans the result into the given value.
func (egb *EntitlementGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := egb.path(ctx)
	if err != nil {
		return err
	}
	egb.sql = query
	return egb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (egb *EntitlementGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := egb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EntitlementGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (egb *EntitlementGroupBy) StringsX(ctx context.Context) []string {
	v, err := egb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = egb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (egb *EntitlementGroupBy) StringX(ctx context.Context) string {
	v, err := egb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EntitlementGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (egb *EntitlementGroupBy) IntsX(ctx context.Context) []int {
	v, err := egb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = egb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (egb *EntitlementGroupBy) IntX(ctx context.Context) int {
	v, err := egb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EntitlementGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (egb *EntitlementGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := egb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = egb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (egb *EntitlementGroupBy) Float64X(ctx context.Context) float64 {
	v, err := egb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(egb.fields) > 1 {
		return nil, errors.New("ent: EntitlementGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := egb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (egb *EntitlementGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := egb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (egb *EntitlementGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = egb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (egb *EntitlementGroupBy) BoolX(ctx context.Context) bool {
	v, err := egb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (egb *EntitlementGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range egb.fields {
		if !entitlement.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := egb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := egb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (egb *EntitlementGroupBy) sqlQuery() *sql.Selector {
	selector := egb.sql.Select()
	aggregation := make([]string, 0, len(egb.fns))
	for _, fn := range egb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(egb.fields)+len(egb.fns))
		for _, f := range egb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(egb.fields...)...)
}

// EntitlementSelect is the builder for selecting fields of Entitlement entities.
type EntitlementSelect struct {
	*EntitlementQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (es *EntitlementSelect) Scan(ctx context.Context, v interface{}) error {
	if err := es.prepareQuery(ctx); err != nil {
		return err
	}
	es.sql = es.EntitlementQuery.sqlQuery(ctx)
	return es.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (es *EntitlementSelect) ScanX(ctx context.Context, v interface{}) {
	if err := es.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Strings(ctx context.Context) ([]string, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EntitlementSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (es *EntitlementSelect) StringsX(ctx context.Context) []string {
	v, err := es.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = es.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (es *EntitlementSelect) StringX(ctx context.Context) string {
	v, err := es.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Ints(ctx context.Context) ([]int, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EntitlementSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (es *EntitlementSelect) IntsX(ctx context.Context) []int {
	v, err := es.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = es.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (es *EntitlementSelect) IntX(ctx context.Context) int {
	v, err := es.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EntitlementSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (es *EntitlementSelect) Float64sX(ctx context.Context) []float64 {
	v, err := es.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = es.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (es *EntitlementSelect) Float64X(ctx context.Context) float64 {
	v, err := es.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(es.fields) > 1 {
		return nil, errors.New("ent: EntitlementSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := es.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (es *EntitlementSelect) BoolsX(ctx context.Context) []bool {
	v, err := es.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (es *EntitlementSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = es.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{entitlement.Label}
	default:
		err = fmt.Errorf("ent: EntitlementSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (es *EntitlementSelect) BoolX(ctx context.Context) bool {
	v, err := es.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (es *EntitlementSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := es.sql.Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
