//go:generate ./tools/genent.sh
//go:generate env GOWORK=off go generate ./ent

package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"os"
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

func New(connspec string, options ...ent.Option) (*Backend, error) {
	var b resource.Builder
	spc, err := b.ServiceProviderConfig().
		AuthenticationSchemes(
			b.AuthenticationScheme().
				Name("OAuth Bearer Token").
				Description("Authentication scheme using the OAuth Bearer Token Standard").
				SpecURI("http://www.rfc-editor.org/info/rfc6750").
				DocumentationURI("http://example.com/help/oauth.html").
				Type(resource.OAuthBearerToken).
				MustBuild(),
		).
		Bulk(b.BulkSupport().
			Supported(false).
			MaxOperations(0).
			MaxPayloadSize(0).
			MustBuild(),
		).
		ETag(b.GenericSupport().
			Supported(true).
			MustBuild(),
		).
		Filter(b.FilterSupport().
			Supported(true).
			MaxResults(200). // TODO: arbitrary value used
			MustBuild(),
		).
		Sort(b.GenericSupport().
			Supported(false).
			MustBuild(),
		).
		// Notes on PATCH support.
		//
		// * fully qualified field names are currently not supported.
		//   Copied from RFC7644:
		//
		//   The attribute notation rules described in Section 3.10 apply for
		//   describing attribute paths.  For all operations, the value of the
		//   "schemas" attribute on the SCIM service provider's representation of
		//   the resource SHALL be assumed by default.  If one of the PATCH
		//   operations modifies the "schemas" attribute, subsequent operations
		//   SHALL assume the modified state of the "schemas" attribute.  Clients
		//   MAY implicitly modify the "schemas" attribute by adding (or
		//   replacing) an attribute with its fully qualified name, including
		//   schema URN.  For example, adding the attribute "urn:ietf:params:scim:
		//   schemas:extension:enterprise:2.0:User:employeeNumber" automatically
		//   adds the value
		//   "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User" to the
		//   resource's "schemas" attribute.
		Patch(b.GenericSupport().
			Supported(true).
			MustBuild(),
		).
		ChangePassword(b.GenericSupport().
			Supported(false).
			MustBuild(),
		).
		Build()
	if err != nil {
		return nil, fmt.Errorf(`failed to setup ServiceProviderConfig: %w`, err)
	}

	client, err := ent.Open(dialect.SQLite, connspec, options...)
	if err != nil {
		return nil, fmt.Errorf(`failed to open database: %w`, err)
	}

	if entTrace {
		client = client.Debug()
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf(`failed to create schema resources: %w`, err)
	}

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

func (b *Backend) createAddress(in ...*resource.Address) ([]*ent.AddressCreate, error) {
	list := make([]*ent.AddressCreate, len(in))
	for i, v := range in {
		addressCreateCall := b.db.Address.Create()
		if v.HasCountry() {
			addressCreateCall.SetCountry(v.Country())
		}

		if v.HasFormatted() {
			addressCreateCall.SetFormatted(v.Formatted())
		}

		if v.HasLocality() {
			addressCreateCall.SetLocality(v.Locality())
		}

		if v.HasPostalCode() {
			addressCreateCall.SetPostalCode(v.PostalCode())
		}

		if v.HasRegion() {
			addressCreateCall.SetRegion(v.Region())
		}

		if v.HasStreetAddress() {
			addressCreateCall.SetStreetAddress(v.StreetAddress())
		}

		list[i] = addressCreateCall
	}

	return list, nil
}

func (b *Backend) createName(v *resource.Names) (*ent.Names, error) {
	nameCreateCall := b.db.Names.Create()
	if v.HasFamilyName() {
		nameCreateCall.SetFamilyName(v.FamilyName())
	}
	if v.HasFormatted() {
		nameCreateCall.SetFormatted(v.Formatted())
	}
	if v.HasGivenName() {
		nameCreateCall.SetGivenName(v.GivenName())
	}
	if v.HasHonorificPrefix() {
		nameCreateCall.SetHonorificPrefix(v.HonorificPrefix())
	}
	if v.HasHonorificSuffix() {
		nameCreateCall.SetHonorificSuffix(v.HonorificSuffix())
	}
	if v.HasMiddleName() {
		nameCreateCall.SetMiddleName(v.MiddleName())
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
		Unique(false).
		Where(user.IDEQ(parsedUUID))

	userLoadEntFields(userQuery, fields, excludedFields)

	u, err := userQuery.
		Only(context.TODO())
	if err != nil {
		return nil, resource.NewErrorBuilder().
			Status(http.StatusNotFound).
			Detail(fmt.Sprintf(`failed to retrieve user: %s`, err)).
			ScimType(resource.ErrUnknown).
			MustBuild()
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

/*
func (b *Backend) createGroupMembers(members ...*resource.GroupMember) ([]*ent.GroupMember, error) {
	ctx := context.TODO()

	list := make([]*ent.GroupMember, len(members))
	for i, in := range members {
		createCall := b.db.GroupMember.Create()

		if !in.HasValue() {
			return nil, fmt.Errorf(`members.value is required`)
		}

		if !in.HasRef() {
			return nil, fmt.Errorf(`members.$ref is required`)
		}

		createCall.SetValue(in.Value())
		ref := in.Ref()
		if in.HasType() {
			createCall.SetType(in.Type())
		} else {
			switch {
			case strings.Contains(ref, `/Users/`):
				createCall.SetType(`User`)
			case strings.Contains(ref, `/Groups/`):
				createCall.SetType(`Group`)
			default:
				return nil, fmt.Errorf(`failed to determine if the resource is a group or a user`)
			}
		}
		createCall.SetRef(ref)

		gm, err := createCall.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf(`failed to save object: %w`, err)
		}
		list[i] = gm
	}
	return list, nil
}
*/
/*
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
}*/

/*
func (b *Backend) CreateGroup(in *resource.Group) (*resource.Group, error) {
	createGroupCall := b.db.Group.Create().
		SetDisplayName(in.DisplayName())

	h := sha256.New()
	fmt.Fprint(h, b.etagSalt)

	members, err := b.createGroupMembers(in.Members()...)
	if err != nil {
		return nil, fmt.Errorf(`failed to create group members: %w`, err)
	}

	createGroupCall.AddMembers(members...)
	createGroupCall.SetEtag(fmt.Sprintf(`W/%q`, base64.RawStdEncoding.EncodeToString(h.Sum(nil))))
	g, err := createGroupCall.
		Save(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to save data: %w`, err)
	}

	// Unfortunately we're going to have to load the actual members here
	// because that's how we transform the data
	g.Edges.Members = members

	return GroupResourceFromEnt(g)
}*/

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
		WithMembers().
		Where(group.IDEQ(parsedUUID))

	groupLoadEntFields(groupQuery, fields, excludedFields)

	g, err := groupQuery.
		Only(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to retrieve group %s: %w`, id, err)
	}

	return GroupResourceFromEnt(g)
}

/*
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
		ClearMembers()

	// optional fields
	if in.HasDisplayName() {
		replaceGroupCall.SetDisplayName(in.DisplayName())
	}

	members, err := b.createGroupMembers(in.Members()...)
	if err != nil {
		return nil, err
	}

	if len(members) > 0 {
		replaceGroupCall.AddMembers(members...)
	}

	if _, err := replaceGroupCall.Save(context.TODO()); err != nil {
		return nil, fmt.Errorf(`failed to update group: %w`, err)
	}

	// Okay, I'm sure we can get the edges (users+children -> members)
	// somehow without re-fetching the data, but we're going to punt this
	return b.RetrieveGroup(id, nil, nil)
}*/

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

func rollbackTx(tx *ent.Tx, oerr error) error {
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf(`failed to rollback transaction: %s (original error = %w)`, err, oerr)
	}
	return oerr
}

/*
   o  If the target location does not exist, the attribute and value are
      added.

   o  If the target location specifies a complex attribute, a set of
      sub-attributes SHALL be specified in the "value" parameter.

   o  If the target location specifies a multi-valued attribute, a new
      value is added to the attribute.

   o  If the target location specifies a single-valued attribute, the
      existing value is replaced.

   o  If the target location specifies an attribute that does not exist
      (has no value), the attribute is added with the new value.

   o  If the target location exists, the value is replaced.

   o  If the target location already contains the value specified, no
      changes SHOULD be made to the resource, and a success response
      SHOULD be returned.  Unless other operations change the resource,
      this operation SHALL NOT change the modify timestamp of the
      resource.
*/

/*
func (m *multiValueMutator) Add(src json.RawMessage) error {
	// Note: in this mode, the parent must exist, because multi-valued
	// elements must be associated with a parent resource
	switch m.field {
	case resource.UserEmailKey:
		var child resource.Email
		if err := json.Unmarshal(src, &child); err != nil {

		}
	}
	return nil
}*/

// The PATCH operation takes a path specifier that might take one of the following
// forms:
//
//  1. "" (empty)
//  2. attr
//  3. attr.subAttr
//  4. attr[...]
//  5. attr[...].subAttr
//
// # Empty path
//
// The first one is special. It is effectively the same as
// Replace, Add, Remove operations. For example,
//
//   PATCH /Users/...
//   {"op":"replace", "path": "", "value": {...}}
//
// Would effectively be the same as PUT /Users, and
//
//   PATCH /Users/...
//   {"op":"remove", "path":"", ...}
//
// Would effectively be the same as DELETE /Users/..., so we can
// just delegate tehese operations to the equivalent operations
//
// The one exception is "add" operation with a path of "", because
// while "add" for a singular value implies to only add if the value
// is not previously present, by nature of the PATCH operation
// we already know the ID of the resource being patched. Therefore
// this case will result in an error
//
// # Everything else
//
// For all of the rest of cases, we know that you will have a top-level
// object (User/Group, or a sub attribute), and its sub attribute to
// perform operations on.
//

// when we use filter.Parse against a PATCH path, we would have
// one of the four forms:
//  1. attr
//  2. attr.subAttr
//  3. attr[...]
//  4. attr[...].subAttr
// Note that this isn't strictly "correct" for a filter query
//
// Cases 1 and 2 are handled by the previous block, and now
// we need to build a where clause
/*
	var buildUsers, buildGroups bool
	switch parent.(type) {
	case *ent.User:
		buildUsers = true
	case *ent.Group:
		buildGroups = true
	default:
		return nil, fmt.Errorf(`invalid parent type %T`, parent)
	}
	uw, gw, err := b.buildWhere(path, buildUsers, buildGroups)
	if err != nil {
		return nil, fmt.Errorf(`failed to build where: %w`, err)
	}

	log.Printf("uw = %#v", uw)
	log.Printf("gw = %#v", gw)

	switch parent := parent.(type) {
	case *ent.Group:
		// whoa, can this be multiple elements?
		attr, err := parent.Query().Where(gw...).Only(ctx)
		return &attrValueMutator{
			backend: b,
			parent:  parent,
			attr:    attr,
		}, nil
	}
*/

func exprStr(expr filter.Expr) (string, error) {
	v, err := exprAttr(expr)
	if err != nil {
		return "", err
	}
	sv, ok := v.(string)
	if err != nil || !ok {
		return "", fmt.Errorf(`expected string, got %T`, v)
	}
	return sv, nil
}

/*
func (b *Backend) patchAddGroup(parent *ent.Group, op *resource.PatchOperation) error {
	ctx := context.TODO()

	root, err := filter.Parse(op.Path())
	if err != nil {
		return fmt.Errorf(`failed to parse PATH path %q`, op.Path())
	}

	expr, ok := root.(filter.ValuePath)
	if !ok {
		return fmt.Errorf(`root element should be a valuePath (got %T)`, root)
	}

	attr, err := exprAttr(expr.ParentAttr())
	sattr, ok := attr.(string)
	if err != nil || !ok {
		return fmt.Errorf(`invalid attribute specification`)
	}

	switch sattr {
	case resource.GroupMembersKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			// Adding to a subAttr of a multi-value element doesn't make
			// sense
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf(`patch add operation on sub attribute of multi-value item members with unspecified element is not possible`)
			}

			// if we're adding to the list, we need the entire thing
			var in resource.GroupMember
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf(`failed to decode patch add value: %w`, err)
			}

			created, err := b.createGroupMembers(&in)
			if err != nil {
				return fmt.Errorf(`failed to create GroupMember: %w`, err)
			}
			if _, err := parent.Update().AddMembers(created...).Save(ctx); err != nil {
				return fmt.Errorf(`failed to save object: %w`, err)
			}
		} else {
			// TODO: this looks fishy, as all fields in the members
			// sub attribute are immutable

			var pb GroupMemberPredicateBuilder
			// so we have a subExpr, that must mean we must have have a subAttr
			// "path": "members[value eq \"...\"].value"  // OK
			// "path": "members[value eq \"...\"]"        // NOT OK
			// also, the query must resolve to a single member element
			// Load attr with the given conditions
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf(`failed to parse valuePath expression: %w`, err)
			}
			members, err := parent.QueryMembers().
				Where(predicates...).
				All(ctx)

			if len(members) > 0 {
				return fmt.Errorf(`query must resolve to one element, got multiple`)
			}

			member := members[0]
			// we must have subAttr
			subAttr, err := exprAttr(expr.SubAttr())
			sSubAttr, ok := subAttr.(string)
			if err != nil || !ok {
				return fmt.Errorf(`query must have a sub attribute`)
			}
			switch sSubAttr {
			case resource.GroupMemberValueKey:
				var s string
				if err := json.Unmarshal(op.Value(), &s); err != nil {
					return fmt.Errorf(`failed to decode value: %w`, err)
				}
				if _, err := member.Update().SetValue(s).Save(ctx); err != nil {
					return fmt.Errorf(`failed to save object: %w`, err)
				}
				return nil
			}
		}
	}
	return nil
}*/

/*
func (b *Backend) patchRemoveGroup(parent *ent.Group, op *resource.PatchOperation) error {
	ctx := context.TODO()

	root, err := filter.Parse(op.Path())
	if err != nil {
		return fmt.Errorf(`failed to parse PATH path %q`, op.Path())
	}

	expr, ok := root.(filter.ValuePath)
	if !ok {
		return fmt.Errorf(`root element should be a valuePath (got %T)`, root)
	}

	sattr, err := exprStr(expr.ParentAttr())
	if err != nil {
		return fmt.Errorf(`invalid attribute specification`)
	}

	switch sattr {
	case resource.GroupDisplayNameKey:
		if subExpr := expr.SubExpr(); subExpr != nil {
			return fmt.Errorf(`patch remove operation on displayName cannot have a query`)
		}
		if subAttr := expr.SubAttr(); subAttr != nil {
			return fmt.Errorf(`patch remove operation on displayName cannot have a sub attribute`)
		}

		if _, err := parent.Update().ClearDisplayName().Save(ctx); err != nil {
			return fmt.Errorf(`failed to save object: %w`, err)
		}
	case resource.GroupMembersKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			// Removing a subAttr of a multi-value element doesn't make sense
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf(`patch remove operation on sub attribute of multi-value item members with unspecified element is not possible`)
			}

			if _, err := parent.Update().ClearMembers().Save(ctx); err != nil {
				return fmt.Errorf(`failed to save object: %w`, err)
			}
		} else {
			var pb GroupMemberPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf(`failed to parse valuePath expression: %w`, err)
			}
			members, err := parent.QueryMembers().
				Where(predicates...).
				All(ctx)

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf(`invalid sub attribute specified`)
				}
				switch subAttr {
				case resource.GroupMemberRefKey:
					return fmt.Errorf(`$ref is not mutable`)
				case resource.GroupMemberTypeKey:
					return fmt.Errorf(`type is not mutable`)
				case resource.GroupMemberValueKey:
					return fmt.Errorf(`value is not mutable`)
				default:
					return fmt.Errorf(`unknown sub attribute specified`)
				}
			}

			ids := make([]int, len(members))
			for i, member := range members {
				ids[i] = member.ID
			}

			b.db.GroupMember.Delete().
				Where(groupmember.IDIn(ids...)).
				Exec(ctx)
		}
	}
	return nil
}
*/

func (b *Backend) PatchUser(id string, r *resource.PatchRequest) (*resource.User, error) {
	tx, err := b.db.Tx(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to start transaction`)
	}
	old := b.db
	b.db = tx.Client()
	defer func() {
		b.db = old
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, rollbackTx(tx, fmt.Errorf(`failed to parse ID: %w`, err))
	}

	userQuery := b.db.User.Query().
		Where(user.IDEQ(parsedUUID))

	u, err := userQuery.
		Only(context.TODO())
	if err != nil {
		return nil, rollbackTx(tx, fmt.Errorf(`failed to retrieve user: %w`, err))
	}

	retrieve := true
	for _, op := range r.Operations() {
		switch op.Op() {
		case resource.PatchAdd:
			if err := b.patchAddUser(u, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		case resource.PatchRemove:
			if err := b.patchRemoveUser(u, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		default:
			return nil, rollbackTx(tx, err)
		}
	}

	var u2 *resource.User

	if retrieve {
		// This is silly, but we're going to have to re-load the object
		u2, err = b.RetrieveUser(id, nil, nil)
		if err != nil {
			return nil, rollbackTx(tx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`failed to commit transaction: %w`, err)
	}

	return u2, nil
}

func (b *Backend) PatchGroup(id string, r *resource.PatchRequest) (*resource.Group, error) {
	tx, err := b.db.Tx(context.TODO())
	if err != nil {
		return nil, fmt.Errorf(`failed to start transaction`)
	}
	old := b.db
	b.db = tx.Client()
	defer func() {
		b.db = old
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, rollbackTx(tx, fmt.Errorf(`failed to parse ID: %w`, err))
	}

	groupQuery := b.db.Group.Query().
		Where(group.IDEQ(parsedUUID))

	g, err := groupQuery.
		Only(context.TODO())
	if err != nil {
		return nil, rollbackTx(tx, fmt.Errorf(`failed to retrieve group: %w`, err))
	}

	retrieve := true
	for _, op := range r.Operations() {
		switch op.Op() {
		case resource.PatchAdd:
			if err := b.patchAddGroup(g, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		case resource.PatchRemove:
			if err := b.patchRemoveGroup(g, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		default:
			return nil, rollbackTx(tx, err)
		}
	}

	var g2 *resource.Group

	if retrieve {
		// This is silly, but we're going to have to re-load the object
		g2, err = b.RetrieveGroup(id, nil, nil)
		if err != nil {
			return nil, rollbackTx(tx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`failed to commit transaction: %w`, err)
	}

	return g2, nil
}
