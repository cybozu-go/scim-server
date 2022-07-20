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
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/groupmember"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/google/uuid"
)

// GroupMemberQuery is the builder for querying GroupMember entities.
type GroupMemberQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.GroupMember
	// eager-loading edges.
	withUser  *UserQuery
	withGroup *GroupQuery
	withFKs   bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GroupMemberQuery builder.
func (gmq *GroupMemberQuery) Where(ps ...predicate.GroupMember) *GroupMemberQuery {
	gmq.predicates = append(gmq.predicates, ps...)
	return gmq
}

// Limit adds a limit step to the query.
func (gmq *GroupMemberQuery) Limit(limit int) *GroupMemberQuery {
	gmq.limit = &limit
	return gmq
}

// Offset adds an offset step to the query.
func (gmq *GroupMemberQuery) Offset(offset int) *GroupMemberQuery {
	gmq.offset = &offset
	return gmq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (gmq *GroupMemberQuery) Unique(unique bool) *GroupMemberQuery {
	gmq.unique = &unique
	return gmq
}

// Order adds an order step to the query.
func (gmq *GroupMemberQuery) Order(o ...OrderFunc) *GroupMemberQuery {
	gmq.order = append(gmq.order, o...)
	return gmq
}

// QueryUser chains the current query on the "user" edge.
func (gmq *GroupMemberQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: gmq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gmq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(groupmember.Table, groupmember.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, groupmember.UserTable, groupmember.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(gmq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryGroup chains the current query on the "group" edge.
func (gmq *GroupMemberQuery) QueryGroup() *GroupQuery {
	query := &GroupQuery{config: gmq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gmq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(groupmember.Table, groupmember.FieldID, selector),
			sqlgraph.To(group.Table, group.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, groupmember.GroupTable, groupmember.GroupColumn),
		)
		fromU = sqlgraph.SetNeighbors(gmq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first GroupMember entity from the query.
// Returns a *NotFoundError when no GroupMember was found.
func (gmq *GroupMemberQuery) First(ctx context.Context) (*GroupMember, error) {
	nodes, err := gmq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{groupmember.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (gmq *GroupMemberQuery) FirstX(ctx context.Context) *GroupMember {
	node, err := gmq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first GroupMember ID from the query.
// Returns a *NotFoundError when no GroupMember ID was found.
func (gmq *GroupMemberQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = gmq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{groupmember.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (gmq *GroupMemberQuery) FirstIDX(ctx context.Context) int {
	id, err := gmq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single GroupMember entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one GroupMember entity is found.
// Returns a *NotFoundError when no GroupMember entities are found.
func (gmq *GroupMemberQuery) Only(ctx context.Context) (*GroupMember, error) {
	nodes, err := gmq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{groupmember.Label}
	default:
		return nil, &NotSingularError{groupmember.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (gmq *GroupMemberQuery) OnlyX(ctx context.Context) *GroupMember {
	node, err := gmq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only GroupMember ID in the query.
// Returns a *NotSingularError when more than one GroupMember ID is found.
// Returns a *NotFoundError when no entities are found.
func (gmq *GroupMemberQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = gmq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = &NotSingularError{groupmember.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (gmq *GroupMemberQuery) OnlyIDX(ctx context.Context) int {
	id, err := gmq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of GroupMembers.
func (gmq *GroupMemberQuery) All(ctx context.Context) ([]*GroupMember, error) {
	if err := gmq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return gmq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (gmq *GroupMemberQuery) AllX(ctx context.Context) []*GroupMember {
	nodes, err := gmq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of GroupMember IDs.
func (gmq *GroupMemberQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := gmq.Select(groupmember.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (gmq *GroupMemberQuery) IDsX(ctx context.Context) []int {
	ids, err := gmq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (gmq *GroupMemberQuery) Count(ctx context.Context) (int, error) {
	if err := gmq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return gmq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (gmq *GroupMemberQuery) CountX(ctx context.Context) int {
	count, err := gmq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (gmq *GroupMemberQuery) Exist(ctx context.Context) (bool, error) {
	if err := gmq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return gmq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (gmq *GroupMemberQuery) ExistX(ctx context.Context) bool {
	exist, err := gmq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GroupMemberQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (gmq *GroupMemberQuery) Clone() *GroupMemberQuery {
	if gmq == nil {
		return nil
	}
	return &GroupMemberQuery{
		config:     gmq.config,
		limit:      gmq.limit,
		offset:     gmq.offset,
		order:      append([]OrderFunc{}, gmq.order...),
		predicates: append([]predicate.GroupMember{}, gmq.predicates...),
		withUser:   gmq.withUser.Clone(),
		withGroup:  gmq.withGroup.Clone(),
		// clone intermediate query.
		sql:    gmq.sql.Clone(),
		path:   gmq.path,
		unique: gmq.unique,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (gmq *GroupMemberQuery) WithUser(opts ...func(*UserQuery)) *GroupMemberQuery {
	query := &UserQuery{config: gmq.config}
	for _, opt := range opts {
		opt(query)
	}
	gmq.withUser = query
	return gmq
}

// WithGroup tells the query-builder to eager-load the nodes that are connected to
// the "group" edge. The optional arguments are used to configure the query builder of the edge.
func (gmq *GroupMemberQuery) WithGroup(opts ...func(*GroupQuery)) *GroupMemberQuery {
	query := &GroupQuery{config: gmq.config}
	for _, opt := range opts {
		opt(query)
	}
	gmq.withGroup = query
	return gmq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Value string `json:"value,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.GroupMember.Query().
//		GroupBy(groupmember.FieldValue).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (gmq *GroupMemberQuery) GroupBy(field string, fields ...string) *GroupMemberGroupBy {
	group := &GroupMemberGroupBy{config: gmq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := gmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return gmq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Value string `json:"value,omitempty"`
//	}
//
//	client.GroupMember.Query().
//		Select(groupmember.FieldValue).
//		Scan(ctx, &v)
//
func (gmq *GroupMemberQuery) Select(fields ...string) *GroupMemberSelect {
	gmq.fields = append(gmq.fields, fields...)
	return &GroupMemberSelect{GroupMemberQuery: gmq}
}

func (gmq *GroupMemberQuery) prepareQuery(ctx context.Context) error {
	for _, f := range gmq.fields {
		if !groupmember.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if gmq.path != nil {
		prev, err := gmq.path(ctx)
		if err != nil {
			return err
		}
		gmq.sql = prev
	}
	return nil
}

func (gmq *GroupMemberQuery) sqlAll(ctx context.Context) ([]*GroupMember, error) {
	var (
		nodes       = []*GroupMember{}
		withFKs     = gmq.withFKs
		_spec       = gmq.querySpec()
		loadedTypes = [2]bool{
			gmq.withUser != nil,
			gmq.withGroup != nil,
		}
	)
	if gmq.withUser != nil || gmq.withGroup != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, groupmember.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &GroupMember{config: gmq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, gmq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := gmq.withUser; query != nil {
		ids := make([]uuid.UUID, 0, len(nodes))
		nodeids := make(map[uuid.UUID][]*GroupMember)
		for i := range nodes {
			if nodes[i].user_groups == nil {
				continue
			}
			fk := *nodes[i].user_groups
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(user.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "user_groups" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.User = n
			}
		}
	}

	if query := gmq.withGroup; query != nil {
		ids := make([]uuid.UUID, 0, len(nodes))
		nodeids := make(map[uuid.UUID][]*GroupMember)
		for i := range nodes {
			if nodes[i].group_members == nil {
				continue
			}
			fk := *nodes[i].group_members
			if _, ok := nodeids[fk]; !ok {
				ids = append(ids, fk)
			}
			nodeids[fk] = append(nodeids[fk], nodes[i])
		}
		query.Where(group.IDIn(ids...))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			nodes, ok := nodeids[n.ID]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "group_members" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.Group = n
			}
		}
	}

	return nodes, nil
}

func (gmq *GroupMemberQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := gmq.querySpec()
	_spec.Node.Columns = gmq.fields
	if len(gmq.fields) > 0 {
		_spec.Unique = gmq.unique != nil && *gmq.unique
	}
	return sqlgraph.CountNodes(ctx, gmq.driver, _spec)
}

func (gmq *GroupMemberQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := gmq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (gmq *GroupMemberQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   groupmember.Table,
			Columns: groupmember.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: groupmember.FieldID,
			},
		},
		From:   gmq.sql,
		Unique: true,
	}
	if unique := gmq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := gmq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, groupmember.FieldID)
		for i := range fields {
			if fields[i] != groupmember.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := gmq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := gmq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := gmq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := gmq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (gmq *GroupMemberQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(gmq.driver.Dialect())
	t1 := builder.Table(groupmember.Table)
	columns := gmq.fields
	if len(columns) == 0 {
		columns = groupmember.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if gmq.sql != nil {
		selector = gmq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if gmq.unique != nil && *gmq.unique {
		selector.Distinct()
	}
	for _, p := range gmq.predicates {
		p(selector)
	}
	for _, p := range gmq.order {
		p(selector)
	}
	if offset := gmq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := gmq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// GroupMemberGroupBy is the group-by builder for GroupMember entities.
type GroupMemberGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (gmgb *GroupMemberGroupBy) Aggregate(fns ...AggregateFunc) *GroupMemberGroupBy {
	gmgb.fns = append(gmgb.fns, fns...)
	return gmgb
}

// Scan applies the group-by query and scans the result into the given value.
func (gmgb *GroupMemberGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := gmgb.path(ctx)
	if err != nil {
		return err
	}
	gmgb.sql = query
	return gmgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := gmgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(gmgb.fields) > 1 {
		return nil, errors.New("ent: GroupMemberGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := gmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) StringsX(ctx context.Context) []string {
	v, err := gmgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = gmgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) StringX(ctx context.Context) string {
	v, err := gmgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(gmgb.fields) > 1 {
		return nil, errors.New("ent: GroupMemberGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := gmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) IntsX(ctx context.Context) []int {
	v, err := gmgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = gmgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) IntX(ctx context.Context) int {
	v, err := gmgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(gmgb.fields) > 1 {
		return nil, errors.New("ent: GroupMemberGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := gmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := gmgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = gmgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) Float64X(ctx context.Context) float64 {
	v, err := gmgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(gmgb.fields) > 1 {
		return nil, errors.New("ent: GroupMemberGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := gmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := gmgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (gmgb *GroupMemberGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = gmgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (gmgb *GroupMemberGroupBy) BoolX(ctx context.Context) bool {
	v, err := gmgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (gmgb *GroupMemberGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range gmgb.fields {
		if !groupmember.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := gmgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gmgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (gmgb *GroupMemberGroupBy) sqlQuery() *sql.Selector {
	selector := gmgb.sql.Select()
	aggregation := make([]string, 0, len(gmgb.fns))
	for _, fn := range gmgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(gmgb.fields)+len(gmgb.fns))
		for _, f := range gmgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(gmgb.fields...)...)
}

// GroupMemberSelect is the builder for selecting fields of GroupMember entities.
type GroupMemberSelect struct {
	*GroupMemberQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (gms *GroupMemberSelect) Scan(ctx context.Context, v interface{}) error {
	if err := gms.prepareQuery(ctx); err != nil {
		return err
	}
	gms.sql = gms.GroupMemberQuery.sqlQuery(ctx)
	return gms.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (gms *GroupMemberSelect) ScanX(ctx context.Context, v interface{}) {
	if err := gms.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Strings(ctx context.Context) ([]string, error) {
	if len(gms.fields) > 1 {
		return nil, errors.New("ent: GroupMemberSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := gms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (gms *GroupMemberSelect) StringsX(ctx context.Context) []string {
	v, err := gms.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = gms.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (gms *GroupMemberSelect) StringX(ctx context.Context) string {
	v, err := gms.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Ints(ctx context.Context) ([]int, error) {
	if len(gms.fields) > 1 {
		return nil, errors.New("ent: GroupMemberSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := gms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (gms *GroupMemberSelect) IntsX(ctx context.Context) []int {
	v, err := gms.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = gms.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (gms *GroupMemberSelect) IntX(ctx context.Context) int {
	v, err := gms.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(gms.fields) > 1 {
		return nil, errors.New("ent: GroupMemberSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := gms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (gms *GroupMemberSelect) Float64sX(ctx context.Context) []float64 {
	v, err := gms.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = gms.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (gms *GroupMemberSelect) Float64X(ctx context.Context) float64 {
	v, err := gms.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(gms.fields) > 1 {
		return nil, errors.New("ent: GroupMemberSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := gms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (gms *GroupMemberSelect) BoolsX(ctx context.Context) []bool {
	v, err := gms.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (gms *GroupMemberSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = gms.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{groupmember.Label}
	default:
		err = fmt.Errorf("ent: GroupMemberSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (gms *GroupMemberSelect) BoolX(ctx context.Context) bool {
	v, err := gms.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (gms *GroupMemberSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := gms.sql.Query()
	if err := gms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
