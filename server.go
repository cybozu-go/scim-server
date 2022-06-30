//go:generate ./tools/genent.sh
//go:generate go generate ./ent

package server

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"entgo.io/ent/dialect"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim-server/ent/user"
	"github.com/cybozu-go/scim/filter"
	"github.com/cybozu-go/scim/resource"
	"github.com/cybozu-go/scim/schema"
	"github.com/google/uuid"
	"github.com/lestrrat-go/rungroup"
	"golang.org/x/text/secure/precis"

	// default driver
	_ "github.com/mattn/go-sqlite3"
)

// TODO: remove these when they have been incorporated
var _ = groupPresencePredicate
var _ = (&Backend{}).createIMS
var _ = (&Backend{}).createEntitlements
var _ = (&Backend{}).createPhotos

var entTrace bool

func init() {
	if v, err := strconv.ParseBool(os.Getenv(`SCIM_ENT_TRACE`)); err == nil {
		entTrace = v
	}
}

func splitScimField(s string) (string, string, error) {
	i := strings.IndexByte(s, '.')
	if i == -1 {
		return s, "", nil
	}

	if i == len(s)-1 {
		return "", "", fmt.Errorf(`invalid field name specification`)
	}

	return s[:i], s[i+1:], nil
}

type Backend struct {
	db       *ent.Client
	spc      *resource.ServiceProviderConfig
	rts      []*resource.ResourceType
	etagSalt []byte
}

func New(connspec string, spc *resource.ServiceProviderConfig) (*Backend, error) {
	client, err := ent.Open(dialect.SQLite, connspec)
	if err != nil {
		return nil, fmt.Errorf(`failed to open database: %w`, err)
	}

	if entTrace {
		client = client.Debug()
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf(`failed to create schema resources: %w`, err)
	}

	var b resource.Builder

	rts := []*resource.ResourceType{
		b.ResourceType().
			ID("User").
			Name("User").
			Endpoint("/Users").
			Description("User Account").
			Schema(resource.UserSchemaURI).
			SchemaExtensions(
				b.SchemaExtension().
					Schema(
						resource.EnterpriseUserSchemaURI).
					MustBuild(),
			).
			MustBuild(),
		b.ResourceType().
			ID("Group").
			Name("Group").
			Endpoint("/Groups").
			Description("Group").
			Schema(resource.GroupSchemaURI).
			MustBuild(),
	}

	salt := make([]byte, 0, 256)
	_, _ = rand.Read(salt)

	return &Backend{
		db:       client,
		spc:      spc,
		rts:      rts,
		etagSalt: salt,
	}, nil
}

func (b *Backend) Close() error {
	return b.db.Close()
}

var chars []byte
var maxchars *big.Int

func init() {
	charmap := make(map[byte]struct{})
	for i := 0x21; i < 0x7E; i++ {
		charmap[byte(i)] = struct{}{}
	}
	delete(charmap, 'I')
	delete(charmap, 'O')
	delete(charmap, '\\')
	delete(charmap, 'l')
	delete(charmap, 'o')
	for c := range charmap {
		chars = append(chars, c)
	}
	maxchars = big.NewInt(int64(len(chars)))
}

func randomString(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		bn, err := rand.Int(rand.Reader, maxchars)
		if err != nil {
			panic(err)
		}
		b.WriteByte(chars[int(bn.Int64())])
	}
	return b.String()
}

func (b *Backend) createAddresses(in *resource.User, h hash.Hash) ([]*ent.Address, error) {
	list := make([]*ent.Address, len(in.Addresses()))
	for i, v := range in.Addresses() {
		addressCreateCall := b.db.Address.Create()
		if v.HasCountry() {
			addressCreateCall.SetCountry(v.Country())
			fmt.Fprint(h, v.Country())
		}

		if v.HasFormatted() {
			addressCreateCall.SetFormatted(v.Formatted())
			fmt.Fprint(h, v.Formatted())
		}

		if v.HasLocality() {
			addressCreateCall.SetLocality(v.Locality())
			fmt.Fprint(h, v.Locality())
		}

		if v.HasPostalCode() {
			addressCreateCall.SetPostalCode(v.PostalCode())
			fmt.Fprint(h, v.PostalCode())
		}

		if v.HasRegion() {
			addressCreateCall.SetRegion(v.Region())
			fmt.Fprint(h, v.Region())
		}

		if v.HasStreetAddress() {
			addressCreateCall.SetStreetAddress(v.StreetAddress())
			fmt.Fprint(h, v.StreetAddress())
		}

		address, err := addressCreateCall.Save(context.TODO())
		if err != nil {
			return nil, fmt.Errorf(`failed to save address: %w`, err)
		}

		list[i] = address
	}

	return list, nil
}

func (b *Backend) createName(v *resource.Names, h hash.Hash) (*ent.Names, error) {
	nameCreateCall := b.db.Names.Create()
	if v.HasFamilyName() {
		nameCreateCall.SetFamilyName(v.FamilyName())
		fmt.Fprint(h, v.FamilyName())
	}

	if v.HasFormatted() {
		nameCreateCall.SetFormatted(v.Formatted())
		fmt.Fprint(h, v.Formatted())
	}

	if v.HasGivenName() {
		nameCreateCall.SetGivenName(v.GivenName())
		fmt.Fprint(h, v.GivenName())
	}
	if v.HasHonorificPrefix() {
		nameCreateCall.SetHonorificPrefix(v.HonorificPrefix())
		fmt.Fprint(h, v.HonorificPrefix())
	}
	if v.HasHonorificSuffix() {
		nameCreateCall.SetHonorificSuffix(v.HonorificSuffix())
		fmt.Fprint(h, v.HonorificSuffix())
	}
	if v.HasMiddleName() {
		nameCreateCall.SetMiddleName(v.MiddleName())
		fmt.Fprint(h, v.MiddleName())
	}

	name, err := nameCreateCall.Save(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to save name: %w`, err)
	}

	return name, nil
}

func (b *Backend) generatePassword(in *resource.User) (string, error) {
	password := in.Password()
	if password == "" {
		password = randomString(25)
	} else {
		norm, err := precis.OpaqueString.String(password)
		if err != nil {
			return "", fmt.Errorf(`failed to normalize password: %w`, err)
		}
		password = norm
	}
	return password, nil
}

func (b *Backend) RetrieveUser(id string, fields []string, excludedFields []string) (*resource.User, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ID: %w`, err)
	}

	userQuery := b.db.User.Query().
		Where(user.IDEQ(parsedUUID))

	userLoadEntFields(userQuery, fields, excludedFields)

	u, err := userQuery.
		Only(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to retrieve user: %w`, err)
	}

	return UserResourceFromEnt(u)
}

func (b *Backend) DeleteUser(id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf(`failed to parse ID: %w`, err)
	}

	if err := b.db.User.DeleteOneID(parsedUUID).Exec(context.TODO()); err != nil {
		return fmt.Errorf(`failed to delete user: %w`, err)
	}

	return nil
}

func (b *Backend) memberIDs(members []*resource.GroupMember) ([]uuid.UUID, []uuid.UUID, error) {
	var userMembers []uuid.UUID
	var groupMembers []uuid.UUID

	for _, member := range members {
		asUUID, err := uuid.Parse(member.Value())
		if err != nil {
			return nil, nil, fmt.Errorf(`expected "value" to contain a valid UUID: %w`, err)
		}

		if strings.Contains(member.Ref(), `/Users/`) {
			userMembers = append(userMembers, asUUID)
		} else if strings.Contains(member.Ref(), `/Groups/`) {
			groupMembers = append(groupMembers, asUUID)
		} else {
			return nil, nil, fmt.Errorf(`$ref is required in group "members" attribute when creating Groups`)
		}
	}

	sort.Slice(userMembers, func(i, j int) bool {
		return userMembers[i].String() < userMembers[j].String()
	})

	sort.Slice(groupMembers, func(i, j int) bool {
		return groupMembers[i].String() < groupMembers[j].String()
	})

	return userMembers, groupMembers, nil
}

func (b *Backend) CreateGroup(in *resource.Group) (*resource.Group, error) {
	userMembers, groupMembers, err := b.memberIDs(in.Members())
	if err != nil {
		return nil, err
	}

	createGroupCall := b.db.Group.Create().
		SetDisplayName(in.DisplayName())

	h := sha256.New()
	fmt.Fprint(h, b.etagSalt)

	if len(userMembers) > 0 {
		createGroupCall.AddUserIDs(userMembers...)
	}

	if len(groupMembers) > 0 {
		createGroupCall.AddChildIDs(groupMembers...)
	}

	createGroupCall.SetEtag(fmt.Sprintf(`W/%q`, base64.RawStdEncoding.EncodeToString(h.Sum(nil))))
	g, err := createGroupCall.
		Save(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to save data: %w`, err)
	}

	// Unfortunately we're going to have to load the actual members here
	// because that's how we transform the data

	children, err := b.db.Group.Query().Where(group.HasParentWith(group.ID(g.ID))).All(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to load children for group`)
	}
	g.Edges.Children = children

	users, err := b.db.User.Query().Where(user.HasGroupsWith(group.ID(g.ID))).All(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to load member users for group`)
	}
	g.Edges.Users = users

	return GroupResourceFromEnt(g)
}

// XXX passing these boolean variables is so ugly
func (b *Backend) buildWhere(src string, buildUsers, buildGroups bool) ([]predicate.User, []predicate.Group, error) {
	expr, err := filter.Parse(src)
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to parse filter: %w`, err)
	}

	var v filterVisitor

	v.uq = b.db.User.Query()
	v.gq = b.db.Group.Query()

	// XXX while /.search (at the root level) allows querying for
	// all resources (well, User and Group only, really), /Users/.search
	// and /Group/.search restrict the search domain to either User or Group
	// only. In this case we need to limit the predicates that we generate

	// we do this by explicitly initializing the storage space
	// (the []predicate.* fields) with a non-nill value
	if buildUsers {
		v.users = []predicate.User{}
	}
	if buildGroups {
		v.groups = []predicate.Group{}
	}

	if err := v.visit(expr); err != nil {
		return nil, nil, fmt.Errorf(`failed to parse filter expression: %w`, err)
	}

	return v.users, v.groups, nil
}

type filterVisitor struct {
	uq     *ent.UserQuery
	gq     *ent.GroupQuery
	users  []predicate.User
	groups []predicate.Group
}

func (v *filterVisitor) visit(expr filter.Expr) error {
	switch expr := expr.(type) {
	case filter.PresenceExpr:
		return v.visitPresenceExpr(expr)
	case filter.CompareExpr:
		return v.visitCompareExpr(expr)
	case filter.RegexExpr:
		return v.visitRegexExpr(expr)
	case filter.LogExpr: // RENAME ME TO LogicalStatement
		return v.visitLogExpr(expr)
	case filter.ParenExpr:
		return v.visitParenExpr(expr)
	case filter.ValuePath:
		return v.visitValuePath(expr)
	default:
		return fmt.Errorf(`unhandled statement type: %T`, expr)
	}
}

func exprAttr(expr interface{}) (interface{}, error) {
	switch v := expr.(type) {
	case string:
		return v, nil
	case interface{ Lit() string }: // IdentifierExpr, AttrValueExpr
		return v.Lit(), nil
	case filter.BoolExpr:
		return v.Lit(), nil
	case filter.NumberExpr:
		return v.Lit(), nil
	default:
		return nil, fmt.Errorf(`unhandled type: %T`, v)
	}
}

func (v *filterVisitor) visitPresenceExpr(expr filter.PresenceExpr) error {
	attr, err := exprAttr(expr.Attr())
	sattr, ok := attr.(string)
	if err != nil || !ok {
		if err == nil && !ok {
			err = fmt.Errorf(`expected string, got %T`, attr)
		}
		return fmt.Errorf(`left hand side of PresenceExpr is not valid: %w`, err)
	}

	switch expr.Operator() {
	case filter.PresenceOp:
		if v.users != nil {
			if pred := userPresencePredicate(sattr); pred != nil {
				v.users = append(v.users, pred)
			}
		}
		return nil
	default:
		return fmt.Errorf(`unhandled attr operator %q`, expr.Operator())
	}
}

func (v *filterVisitor) visitRegexExpr(expr filter.RegexExpr) error {
	lhe, err := exprAttr(expr.LHE())
	slhe, ok := lhe.(string)
	if err != nil || !ok {
		return fmt.Errorf(`left hand side of RegexExpr is not valid`)
	}

	rhe, err := exprAttr(expr.Value())
	if err != nil {
		return fmt.Errorf(`right hand side of RegexExpr is not valid: %w`, err)
	}
	// convert rhe to string so it can be passed to regexp.QuoteMeta
	srhe := fmt.Sprintf(`%v`, rhe)

	switch expr.Operator() {
	case filter.ContainsOp:
		if v.users != nil {
			pred, err := userContainsPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupContainsPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	case filter.StartsWithOp:
		if v.users != nil {
			pred, err := userStartsWithPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupStartsWithPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	case filter.EndsWithOp:
		if v.users != nil {
			pred, err := userEndsWithPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupEndsWithPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	default:
		return fmt.Errorf(`unhandled regexp operator %q`, expr.Operator())
	}
}

func (v *filterVisitor) visitCompareExpr(expr filter.CompareExpr) error {
	lhe, err := exprAttr(expr.LHE())
	slhe, ok := lhe.(string)
	if err != nil || !ok {
		return fmt.Errorf(`left hand side of CompareExpr is not valid`)
	}

	rhe, err := exprAttr(expr.RHE())
	if err != nil {
		return fmt.Errorf(`right hand side of CompareExpr is not valid: %w`, err)
	}
	// convert rhe to string so it can be passed to regexp.QuoteMeta
	srhe := fmt.Sprintf(`%v`, rhe)

	switch expr.Operator() {
	case filter.EqualOp:
		if v.users != nil {
			pred, err := userEqualsPredicate(v.uq, slhe, srhe)
			if err != nil {
				return err
			}
			v.users = append(v.users, pred)
		}
		if v.groups != nil {
			pred, err := groupEqualsPredicate(v.gq, slhe, srhe)
			if err != nil {
				return err
			}
			v.groups = append(v.groups, pred)
		}
		return nil
	default:
		panic(expr.Operator())
	}
}

func (v *filterVisitor) visitLogExpr(expr filter.LogExpr) error {
	if err := v.visit(expr.LHE()); err != nil {
		return fmt.Errorf(`failed to parse left hand side of %q statement: %w`, expr.Operator(), err)
	}
	if err := v.visit(expr.RHS()); err != nil {
		return fmt.Errorf(`failed to parse right hand side of %q statement: %w`, expr.Operator(), err)
	}

	switch expr.Operator() {
	case "and":
		if v.users != nil {
			v.users = []predicate.User{user.And(v.users...)}
		}
		if v.groups != nil {
			v.groups = []predicate.Group{group.And(v.groups...)}
		}
	case "or":
		if v.users != nil {
			v.users = []predicate.User{user.Or(v.users...)}
		}
		if v.groups != nil {
			v.groups = []predicate.Group{group.Or(v.groups...)}
		}
	default:
		return fmt.Errorf(`unhandled logical statement operator %q`, expr.Operator())
	}
	return nil
}
func (v *filterVisitor) visitParenExpr(expr filter.Expr) error {
	return fmt.Errorf(`unimplemented`)
}

func (v *filterVisitor) visitValuePath(expr filter.Expr) error {
	return fmt.Errorf(`unimplemented`)
}

func (b *Backend) Search(in *resource.SearchRequest) (*resource.ListResponse, error) {
	return b.search(in, true, true)
}

func (b *Backend) SearchUser(in *resource.SearchRequest) (*resource.ListResponse, error) {
	return b.search(in, true, false)
}

func (b *Backend) SearchGroup(in *resource.SearchRequest) (*resource.ListResponse, error) {
	return b.search(in, false, true)
}

func (b *Backend) search(in *resource.SearchRequest, searchUser, searchGroup bool) (*resource.ListResponse, error) {
	userWhere, groupWhere, err := b.buildWhere(in.Filter(), searchUser, searchGroup)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse filter: %w`, err)
	}

	if in.Filter() != "" && (len(userWhere) == 0 && len(groupWhere) == 0) {
		var builder resource.Builder
		return builder.ListResponse().
			TotalResults(0).
			Build()
	}

	var listMu sync.Mutex
	var list []interface{}

	var g rungroup.Group
	if searchUser {
		_ = g.Add(rungroup.ActorFunc(func(ctx context.Context) error {
			users, err := b.db.User.Query().Where(userWhere...).
				All(ctx)
			if err != nil {
				return fmt.Errorf(`failed to execute query: %w`, err)
			}

			for _, user := range users {
				r, err := UserResourceFromEnt(user)
				if err != nil {
					return fmt.Errorf(`failed to convert internal data to SCIM resource: %w`, err)
				}
				listMu.Lock()
				list = append(list, r)
				listMu.Unlock()
			}
			return nil
		}))
	}

	if searchGroup {
		_ = g.Add(rungroup.ActorFunc(func(ctx context.Context) error {
			groups, err := b.db.Group.Query().Where(groupWhere...).
				All(ctx)
			if err != nil {
				return fmt.Errorf(`failed to execute query: %w`, err)
			}

			for _, group := range groups {
				r, err := GroupResourceFromEnt(group)
				if err != nil {
					return fmt.Errorf(`failed to convert internal data to SCIM resource: %w`, err)
				}
				listMu.Lock()
				list = append(list, r)
				listMu.Unlock()
			}
			return nil
		}))
	}

	if err := <-g.Run(context.TODO()); err != nil {
		return nil, err
	}

	var builder resource.Builder
	return builder.ListResponse().
		TotalResults(len(list)).
		Resources(list...).
		Build()
}

func (b *Backend) RetrieveGroup(id string, fields []string, excludedFields []string) (*resource.Group, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ID: %w`, err)
	}

	groupQuery := b.db.Group.Query().
		WithUsers().
		WithChildren().
		Where(group.IDEQ(parsedUUID))

	groupLoadEntFields(groupQuery, fields, excludedFields)

	g, err := groupQuery.
		Only(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to retrieve group %s: %w`, id, err)
	}

	return GroupResourceFromEnt(g)
}

func (b *Backend) ReplaceGroup(id string, in *resource.Group) (*resource.Group, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ID: %w`, err)
	}

	// TODO: is it possible to just grab the ID or check existence?
	g, err := b.db.Group.Query().
		Where(group.IDEQ(parsedUUID)).
		Only(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to retrieve group for replace: %w`, err)
	}

	replaceGroupCall := g.Update().
		ClearUsers().
		ClearChildren()

	// optional fields
	if in.HasDisplayName() {
		replaceGroupCall.SetDisplayName(in.DisplayName())
	}

	userMembers, groupMembers, err := b.memberIDs(in.Members())
	if err != nil {
		return nil, err
	}

	if len(userMembers) > 0 {
		replaceGroupCall.AddUserIDs(userMembers...)
	}

	if len(groupMembers) > 0 {
		replaceGroupCall.AddChildIDs(groupMembers...)
	}

	if _, err := replaceGroupCall.Save(context.TODO()); err != nil {
		return nil, fmt.Errorf(`failed to update group: %w`, err)
	}

	// Okay, I'm sure we can get the edges (users+children -> members)
	// somehow without re-fetching the data, but we're going to punt this
	return b.RetrieveGroup(id, nil, nil)
}

func (b *Backend) DeleteGroup(id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf(`failed to parse ID: %w`, err)
	}

	if err := b.db.Group.DeleteOneID(parsedUUID).Exec(context.TODO()); err != nil {
		return fmt.Errorf(`failed to delete group: %w`, err)
	}

	return nil
}

func (b *Backend) RetrieveServiceProviderConfig() (*resource.ServiceProviderConfig, error) {
	// TODO: meta?
	return b.spc, nil
}

func (b *Backend) RetrieveResourceTypes() ([]*resource.ResourceType, error) {
	// TODO: meta?
	return b.rts, nil
}

func (b *Backend) ListSchemas() (*resource.ListResponse, error) {
	schemas := schema.All()

	// Need to convert this to interface{}
	list := make([]interface{}, len(schemas))
	for i, s := range schemas {
		list[i] = s
	}

	var builder resource.Builder
	return builder.ListResponse().
		TotalResults(len(list)).
		Resources(list...).
		Build()
}

func (b *Backend) RetrieveSchema(id string) (*resource.Schema, error) {
	s, ok := schema.Get(id)
	if !ok {
		return nil, fmt.Errorf(`schema %q not found`, id)
	}
	return s, nil
}
