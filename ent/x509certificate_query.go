// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim-server/ent/x509certificate"
	"github.com/google/uuid"
)

// X509CertificateQuery is the builder for querying X509Certificate entities.
type X509CertificateQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.X509Certificate
	// eager-loading edges.
	withUser *UserQuery
	withFKs  bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the X509CertificateQuery builder.
func (xq *X509CertificateQuery) Where(ps ...predicate.X509Certificate) *X509CertificateQuery {
	xq.predicates = append(xq.predicates, ps...)
	return xq
}

// Limit adds a limit step to the query.
func (xq *X509CertificateQuery) Limit(limit int) *X509CertificateQuery {
	xq.limit = &limit
	return xq
}

// Offset adds an offset step to the query.
func (xq *X509CertificateQuery) Offset(offset int) *X509CertificateQuery {
	xq.offset = &offset
	return xq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (xq *X509CertificateQuery) Unique(unique bool) *X509CertificateQuery {
	xq.unique = &unique
	return xq
}

// Order adds an order step to the query.
func (xq *X509CertificateQuery) Order(o ...OrderFunc) *X509CertificateQuery {
	xq.order = append(xq.order, o...)
	return xq
}

// QueryUser chains the current query on the "user" edge.
func (xq *X509CertificateQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: xq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := xq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := xq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(x509certificate.Table, x509certificate.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, x509certificate.UserTable, x509certificate.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(xq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first X509Certificate entity from the query.
// Returns a *NotFoundError when no X509Certificate was found.
func (xq *X509CertificateQuery) First(ctx context.Context) (*X509Certificate, error) {
	nodes, err := xq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{x509certificate.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (xq *X509CertificateQuery) FirstX(ctx context.Context) *X509Certificate {
	node, err := xq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first X509Certificate ID from the query.
// Returns a *NotFoundError when no X509Certificate ID was found.
func (xq *X509CertificateQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = xq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{x509certificate.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (xq *X509CertificateQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := xq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single X509Certificate entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one X509Certificate entity is found.
// Returns a *NotFoundError when no X509Certificate entities are found.
func (xq *X509CertificateQuery) Only(ctx context.Context) (*X509Certificate, error) {
	nodes, err := xq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{x509certificate.Label}
	default:
		return nil, &NotSingularError{x509certificate.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (xq *X509CertificateQuery) OnlyX(ctx context.Context) *X509Certificate {
	node, err := xq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only X509Certificate ID in the query.
// Returns a *NotSingularError when more than one X509Certificate ID is found.
// Returns a *NotFoundError when no entities are found.
func (xq *X509CertificateQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = xq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{x509certificate.Label}
	default:
		err = &NotSingularError{x509certificate.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (xq *X509CertificateQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := xq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of X509Certificates.
func (xq *X509CertificateQuery) All(ctx context.Context) ([]*X509Certificate, error) {
	if err := xq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return xq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (xq *X509CertificateQuery) AllX(ctx context.Context) []*X509Certificate {
	nodes, err := xq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of X509Certificate IDs.
func (xq *X509CertificateQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := xq.Select(x509certificate.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (xq *X509CertificateQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := xq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (xq *X509CertificateQuery) Count(ctx context.Context) (int, error) {
	if err := xq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return xq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (xq *X509CertificateQuery) CountX(ctx context.Context) int {
	count, err := xq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (xq *X509CertificateQuery) Exist(ctx context.Context) (bool, error) {
	if err := xq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return xq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (xq *X509CertificateQuery) ExistX(ctx context.Context) bool {
	exist, err := xq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the X509CertificateQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (xq *X509CertificateQuery) Clone() *X509CertificateQuery {
	if xq == nil {
		return nil
	}
	return &X509CertificateQuery{
		config:     xq.config,
		limit:      xq.limit,
		offset:     xq.offset,
		order:      append([]OrderFunc{}, xq.order...),
		predicates: append([]predicate.X509Certificate{}, xq.predicates...),
		withUser:   xq.withUser.Clone(),
		// clone intermediate query.
		sql:    xq.sql.Clone(),
		path:   xq.path,
		unique: xq.unique,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (xq *X509CertificateQuery) WithUser(opts ...func(*UserQuery)) *X509CertificateQuery {
	query := &UserQuery{config: xq.config}
	for _, opt := range opts {
		opt(query)
	}
	xq.withUser = query
	return xq
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
//	client.X509Certificate.Query().
//		GroupBy(x509certificate.FieldDisplay).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (xq *X509CertificateQuery) GroupBy(field string, fields ...string) *X509CertificateGroupBy {
	grbuild := &X509CertificateGroupBy{config: xq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := xq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return xq.sqlQuery(ctx), nil
	}
	grbuild.label = x509certificate.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
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
//	client.X509Certificate.Query().
//		Select(x509certificate.FieldDisplay).
//		Scan(ctx, &v)
//
func (xq *X509CertificateQuery) Select(fields ...string) *X509CertificateSelect {
	xq.fields = append(xq.fields, fields...)
	selbuild := &X509CertificateSelect{X509CertificateQuery: xq}
	selbuild.label = x509certificate.Label
	selbuild.flds, selbuild.scan = &xq.fields, selbuild.Scan
	return selbuild
}

func (xq *X509CertificateQuery) prepareQuery(ctx context.Context) error {
	for _, f := range xq.fields {
		if !x509certificate.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if xq.path != nil {
		prev, err := xq.path(ctx)
		if err != nil {
			return err
		}
		xq.sql = prev
	}
	return nil
}

func (xq *X509CertificateQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*X509Certificate, error) {
	var (
		nodes       = []*X509Certificate{}
		withFKs     = xq.withFKs
		_spec       = xq.querySpec()
		loadedTypes = [1]bool{
			xq.withUser != nil,
		}
	)
	if xq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, x509certificate.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*X509Certificate).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &X509Certificate{config: xq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, xq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := xq.withUser; query != nil {
		ids := make([]uuid.UUID, 0, len(nodes))
		nodeids := make(map[uuid.UUID][]*X509Certificate)
		for i := range nodes {
			if nodes[i].user_x509_certificates == nil {
				continue
			}
			fk := *nodes[i].user_x509_certificates
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
				return nil, fmt.Errorf(`unexpected foreign-key "user_x509_certificates" returned %v`, n.ID)
			}
			for i := range nodes {
				nodes[i].Edges.User = n
			}
		}
	}

	return nodes, nil
}

func (xq *X509CertificateQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := xq.querySpec()
	_spec.Node.Columns = xq.fields
	if len(xq.fields) > 0 {
		_spec.Unique = xq.unique != nil && *xq.unique
	}
	return sqlgraph.CountNodes(ctx, xq.driver, _spec)
}

func (xq *X509CertificateQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := xq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (xq *X509CertificateQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   x509certificate.Table,
			Columns: x509certificate.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: x509certificate.FieldID,
			},
		},
		From:   xq.sql,
		Unique: true,
	}
	if unique := xq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := xq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, x509certificate.FieldID)
		for i := range fields {
			if fields[i] != x509certificate.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := xq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := xq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := xq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := xq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (xq *X509CertificateQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(xq.driver.Dialect())
	t1 := builder.Table(x509certificate.Table)
	columns := xq.fields
	if len(columns) == 0 {
		columns = x509certificate.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if xq.sql != nil {
		selector = xq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if xq.unique != nil && *xq.unique {
		selector.Distinct()
	}
	for _, p := range xq.predicates {
		p(selector)
	}
	for _, p := range xq.order {
		p(selector)
	}
	if offset := xq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := xq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// X509CertificateGroupBy is the group-by builder for X509Certificate entities.
type X509CertificateGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (xgb *X509CertificateGroupBy) Aggregate(fns ...AggregateFunc) *X509CertificateGroupBy {
	xgb.fns = append(xgb.fns, fns...)
	return xgb
}

// Scan applies the group-by query and scans the result into the given value.
func (xgb *X509CertificateGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := xgb.path(ctx)
	if err != nil {
		return err
	}
	xgb.sql = query
	return xgb.sqlScan(ctx, v)
}

func (xgb *X509CertificateGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range xgb.fields {
		if !x509certificate.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := xgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := xgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (xgb *X509CertificateGroupBy) sqlQuery() *sql.Selector {
	selector := xgb.sql.Select()
	aggregation := make([]string, 0, len(xgb.fns))
	for _, fn := range xgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(xgb.fields)+len(xgb.fns))
		for _, f := range xgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(xgb.fields...)...)
}

// X509CertificateSelect is the builder for selecting fields of X509Certificate entities.
type X509CertificateSelect struct {
	*X509CertificateQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (xs *X509CertificateSelect) Scan(ctx context.Context, v interface{}) error {
	if err := xs.prepareQuery(ctx); err != nil {
		return err
	}
	xs.sql = xs.X509CertificateQuery.sqlQuery(ctx)
	return xs.sqlScan(ctx, v)
}

func (xs *X509CertificateSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := xs.sql.Query()
	if err := xs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
