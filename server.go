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

func (b *Backend) createAddress(_ context.Context, in ...*resource.Address) ([]*ent.AddressCreate, error) {
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

func (b *Backend) createName(ctx context.Context, v *resource.Names) (*ent.Names, error) {
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

	name, err := nameCreateCall.Save(ctx)
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

func (b *Backend) RetrieveUser(ctx context.Context, id string, fields []string, excludedFields []string) (*resource.User, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ID: %w`, err)
	}

	userQuery := b.db.User.Query().
		Unique(false).
		Where(user.IDEQ(parsedUUID))

	userLoadEntFields(userQuery, fields, excludedFields)

	u, err := userQuery.
		Only(ctx)
	if err != nil {
		return nil, resource.NewErrorBuilder().
			Status(http.StatusNotFound).
			Detail(fmt.Sprintf(`failed to retrieve user: %s`, err)).
			ScimType(resource.ErrUnknown).
			MustBuild()
	}

	return UserResourceFromEnt(u)
}

func (b *Backend) DeleteUser(ctx context.Context, id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf(`failed to parse ID: %w`, err)
	}

	if err := b.db.User.DeleteOneID(parsedUUID).Exec(ctx); err != nil {
		return fmt.Errorf(`failed to delete user: %w`, err)
	}

	return nil
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

func (b *Backend) Search(ctx context.Context, in *resource.SearchRequest) (*resource.ListResponse, error) {
	return b.search(ctx, in, true, true)
}

func (b *Backend) SearchUser(ctx context.Context, in *resource.SearchRequest) (*resource.ListResponse, error) {
	return b.search(ctx, in, true, false)
}

func (b *Backend) SearchGroup(ctx context.Context, in *resource.SearchRequest) (*resource.ListResponse, error) {
	return b.search(ctx, in, false, true)
}

func (b *Backend) search(ctx context.Context, in *resource.SearchRequest, searchUser, searchGroup bool) (*resource.ListResponse, error) {
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

	if err := <-g.Run(ctx); err != nil {
		return nil, err
	}

	var builder resource.Builder
	return builder.ListResponse().
		TotalResults(len(list)).
		Resources(list...).
		Build()
}

func (b *Backend) RetrieveGroup(ctx context.Context, id string, fields []string, excludedFields []string) (*resource.Group, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ID: %w`, err)
	}

	groupQuery := b.db.Group.Query().
		WithMembers().
		Where(group.IDEQ(parsedUUID))

	groupLoadEntFields(groupQuery, fields, excludedFields)

	g, err := groupQuery.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf(`failed to retrieve group %s: %w`, id, err)
	}

	return GroupResourceFromEnt(g)
}

func (b *Backend) DeleteGroup(ctx context.Context, id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf(`failed to parse ID: %w`, err)
	}

	if err := b.db.Group.DeleteOneID(parsedUUID).Exec(ctx); err != nil {
		return fmt.Errorf(`failed to delete group: %w`, err)
	}

	return nil
}

func (b *Backend) RetrieveServiceProviderConfig(_ context.Context) (*resource.ServiceProviderConfig, error) {
	// TODO: meta?
	return b.spc, nil
}

func (b *Backend) RetrieveResourceTypes(_ context.Context) ([]*resource.ResourceType, error) {
	// TODO: meta?
	return b.rts, nil
}

func (b *Backend) ListSchemas(_ context.Context) (*resource.ListResponse, error) {
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

func (b *Backend) RetrieveSchema(_ context.Context, id string) (*resource.Schema, error) {
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

func (b *Backend) PatchUser(ctx context.Context, id string, r *resource.PatchRequest) (*resource.User, error) {
	tx, err := b.db.Tx(ctx)
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
		Only(ctx)
	if err != nil {
		return nil, rollbackTx(tx, fmt.Errorf(`failed to retrieve user: %w`, err))
	}

	retrieve := true
	for _, op := range r.Operations() {
		switch op.Op() {
		case resource.PatchAdd:
			if err := b.patchAddUser(ctx, u, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		case resource.PatchRemove:
			if err := b.patchRemoveUser(ctx, u, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		default:
			return nil, rollbackTx(tx, err)
		}
	}

	var u2 *resource.User

	if retrieve {
		// This is silly, but we're going to have to re-load the object
		u2, err = b.RetrieveUser(ctx, id, nil, nil)
		if err != nil {
			return nil, rollbackTx(tx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`failed to commit transaction: %w`, err)
	}

	return u2, nil
}

func (b *Backend) PatchGroup(ctx context.Context, id string, r *resource.PatchRequest) (*resource.Group, error) {
	tx, err := b.db.Tx(ctx)
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
		Only(ctx)
	if err != nil {
		return nil, rollbackTx(tx, fmt.Errorf(`failed to retrieve group: %w`, err))
	}

	retrieve := true
	for _, op := range r.Operations() {
		switch op.Op() {
		case resource.PatchAdd:
			if err := b.patchAddGroup(ctx, g, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		case resource.PatchRemove:
			if err := b.patchRemoveGroup(ctx, g, op); err != nil {
				return nil, rollbackTx(tx, err)
			}
		default:
			return nil, rollbackTx(tx, err)
		}
	}

	var g2 *resource.Group

	if retrieve {
		// This is silly, but we're going to have to re-load the object
		g2, err = b.RetrieveGroup(ctx, id, nil, nil)
		if err != nil {
			return nil, rollbackTx(tx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(`failed to commit transaction: %w`, err)
	}

	return g2, nil
}
