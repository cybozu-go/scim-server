package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cybozu-go/scim/resource"
	"github.com/cybozu-go/scim/schema"
	"github.com/goccy/go-yaml"
	"github.com/lestrrat-go/codegen"
	"github.com/lestrrat-go/runcmd"
	"github.com/lestrrat-go/xstrings"
)

var objectMap map[string]*codegen.Object
var groupEdges = []string{
	`Members`,
}
var userEdges = []string{
	`Emails`,
	`Roles`,
	`PhoneNumbers`,
	`Entitlements`,
	`IMS`,
	`Photos`,
	`Addresses`,
	`X509Certificates`,
}

var groupEdgeMap map[string]struct{}
var userEdgeMap map[string]struct{}
var groupSchema *resource.Schema
var userSchema *resource.Schema

func main() {
	s, ok := schema.GetByResourceType(`User`)
	if !ok {
		panic("could not find User resource schema")
	}
	userSchema = s
	s, ok = schema.GetByResourceType(`Group`)
	if !ok {
		panic("could not find Group resource schema")
	}
	groupSchema = s

	userEdgeMap = make(map[string]struct{})
	for _, e := range userEdges {
		userEdgeMap[e] = struct{}{}
	}
	groupEdgeMap = make(map[string]struct{})
	for _, e := range groupEdges {
		groupEdgeMap[e] = struct{}{}
	}

	objectMap = make(map[string]*codegen.Object)
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func isGroupEdge(field codegen.Field) bool {
	_, ok := groupEdgeMap[field.Name(true)]
	return ok
}

func isUserEdge(field codegen.Field) bool {
	_, ok := userEdgeMap[field.Name(true)]
	return ok
}

func isEdge(object *codegen.Object, field codegen.Field) bool {
	switch object.Name(true) {
	case `Group`:
		return isGroupEdge(field)
	case `User`:
		return isUserEdge(field)
	}
	return false
}

func isMutable(object *codegen.Object, field codegen.Field) bool {
	return false
}

func entName(field codegen.Field, exported bool) string {
	n := field.Name(false)
	if v := field.String(`ent_name`); v != "" {
		n = v
	}

	if exported {
		n = xstrings.Camel(n)

		// cheat
		if strings.HasSuffix(n, `Id`) {
			n = strings.TrimSuffix(n, `Id`) + `ID`
		}
	}
	return n
}

func addMethod(field codegen.Field) string {
	queryMethod := fmt.Sprintf("Add%s", field.Name(true))
	return queryMethod
}

func queryMethod(field codegen.Field) string {
	queryMethod := fmt.Sprintf("Query%s", field.Name(true))
	return queryMethod
}

func clearMethod(field codegen.Field) string {
	clearMethod := fmt.Sprintf(`Clear%s`, field.Name(true))
	return clearMethod
}

func edgeName(field codegen.Field) string {
	n := field.Name(true)
	return n
}

func resourceName(src interface{}) string {
	typ := scimResourceName(src)
	if typ == `GroupMember` {
		typ = `Member`
	}
	return typ
}

func scimResourceName(src interface{}) string {
	switch src := src.(type) {
	case *codegen.Object:
		return src.Name(true)
	case codegen.Field:
		return strings.TrimPrefix(strings.TrimPrefix(src.Type(), `[]`), `*`)
	default:
		panic(fmt.Sprintf(`invalid object %T`, src))
	}
}

func yaml2json(fn string) ([]byte, error) {
	in, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf(`failed to open %q: %w`, fn, err)
	}
	defer in.Close()

	var v interface{}
	if err := yaml.NewDecoder(in).Decode(&v); err != nil {
		return nil, fmt.Errorf(`failed to decode %q: %w`, fn, err)
	}

	return json.Marshal(v)
}

type GenEnt struct {
	dir      string
	cloneURL string
	cloneDir string
}

func (g *GenEnt) RemoveCloneDir() {
	if g.cloneDir != "" {
		_ = os.RemoveAll(g.cloneDir)
	}
}

func _main() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cloneDir = flag.String("clone-dir", "", "")
	flag.Parse()

	var g = GenEnt{
		cloneURL: `https://github.com/cybozu-go/scim.git`,
		cloneDir: *cloneDir,
	}

	if g.cloneDir == "" {
		if err := g.cloneSCIM(ctx); err != nil {
			return fmt.Errorf(`failed to clone cybozu-go/scim: %w`, err)
		}
		defer g.RemoveCloneDir()
	}

	var objectsFile = flag.String("objects", filepath.Join(g.cloneDir, `tools`, `cmd`, `genresources`, "objects.yml"), "")
	flag.Parse()
	jsonSrc, err := yaml2json(*objectsFile)
	if err != nil {
		return err
	}

	var def struct {
		Common  codegen.FieldList
		Objects []*codegen.Object `json:"objects"`
	}
	if err := json.NewDecoder(bytes.NewReader(jsonSrc)).Decode(&def); err != nil {
		return fmt.Errorf(`failed to decode %q: %w`, *objectsFile, err)
	}

	for _, object := range def.Objects {
		// Each object needs a common set of fields.
		if !object.Bool(`skipCommonFields`) {
			for _, commonField := range def.Common {
				object.AddField(commonField)
			}
		}
		if object.String(`schema`) != "" {
			// TODO: we needed codegen.FieldBulder
			var fl codegen.FieldList
			if err := json.Unmarshal([]byte(`[{"name":"schemas","type":"schemas"}]`), &fl); err != nil {
				return fmt.Errorf(`failed to unmarshal schemas field: %w`, err)
			}
			object.AddField(fl[0])
		}

		object.Organize()

		objectMap[object.Name(true)] = object
	}

	for _, object := range def.Objects {
		if err := generateEnt(object); err != nil {
			return fmt.Errorf(`failed to generate ent adapter: %w`, err)
		}
	}

	return nil
}

func (g *GenEnt) cloneSCIM(ctx context.Context) error {
	dir, err := os.MkdirTemp("", "scim-server-*")
	if err != nil {
		return fmt.Errorf(`failed to create temporary directory: %w`, err)
	}
	g.dir = dir

	if err := runcmd.Run(runcmd.Context(ctx).WithDir(dir), "git", "clone", g.cloneURL); err != nil {
		return fmt.Errorf(`failed to run git clone command: %w`, err)
	}

	g.cloneDir = filepath.Join(dir, `scim`)
	return nil
}

func generateEnt(object *codegen.Object) error {
	// for the time being, only generate for hardcoded objects.
	// later, move this definition to objects.yml
	switch object.Name(true) {
	case `User`, `Group`, `Email`, `Names`, `Role`, `Photo`, `IMS`, `PhoneNumber`, `Address`, `Entitlement`, `X509Certificate`, `GroupMember`:
	default:
		return nil
	}

	fmt.Printf("  ⌛ Generating ent adapters for %s...\n", object.Name(true))

	// TODO: check if this is legit
	if object.Name(true) != `GroupMember` {
		if err := generateSchema(object); err != nil {
			return fmt.Errorf(`failed to generate schema: %w`, err)
		}
	}

	if err := generateUtilities(object); err != nil {
		return fmt.Errorf(`failed to generate utilities: %w`, err)
	}

	return nil
}

func singularName(s string) string {
	if s == "ims" {
		return s
	}

	s2 := strings.Replace(s, `ddresses`, `ddress`, 1)
	if s != s2 {
		s = s2
	} else {
		if s[len(s)-1] == 's' {
			s = s[:len(s)-1]
		}
	}
	return s
}

func relationFilename(s string) string {
	if s == "groupMember" || s == "member" {
		return "member"
	}
	s = xstrings.Snake(s)
	s = strings.Replace(s, `im_s`, `ims`, 1)
	s = strings.Replace(s, `x_509`, `x509`, 1)
	return singularName(s)
}

func packageName(s string) string {
	if s == "IMS" {
		return "ims"
	}
	if strings.Contains(s, "509") {
		return strings.ToLower(s)
	}
	if s == "groupMember" || s == "member" {
		return "member"
	}
	return strings.ToLower(s)
}

func generateSchema(object *codegen.Object) error {
	var buf bytes.Buffer
	o := codegen.NewOutput(&buf)

	o.L(`package schema`)

	pkgs := []string{
		"entgo.io/ent/schema/edge",
		"entgo.io/ent/schema/field",
	}
	o.L(`import (`)
	for _, pkg := range pkgs {
		o.L(`%q`, pkg)
	}
	o.L(`)`)

	o.LL(`type %s struct {`, object.Name(true))
	o.L(`ent.Schema`)
	o.L(`}`)

	o.LL(`func (%s) Fields() []ent.Field {`, object.Name(true))
	o.L(`return []ent.Field{`)
	o.L(`field.UUID("id", uuid.UUID{}).Default(uuid.New),`)

	for _, field := range object.Fields() {
		switch field.Name(true) {
		case "Schemas", `Meta`:
			continue
		}

		ft := field.Type()
		if strings.HasPrefix(ft, `[]`) {
			continue
		}

		var entType = field.String(`ent_type`)
		var entMethod string
		if strings.HasPrefix(ft, `*`) {
			entMethod = `UUID`
			entType = `uuid.UUID{}`
		} else {
			entMethod = xstrings.Camel(ft)
		}

		if v := field.String(`ent_build_method`); v != "" {
			entMethod = v
		}

		var entName = entName(field, false)
		var entDefault = field.String(`ent_default`)

		if entType != "" {
			o.L(`field.%s(%q, %s)`, entMethod, entName, entType)
		} else {
			o.L(`field.%s(%q)`, entMethod, entName)
		}

		if entDefault != "" {
			o.R(`.Default(%s)`, entDefault)
		}
		if !field.IsRequired() {
			o.R(`.Optional()`)
		}

		if field.Bool(`ent_unique`) {
			o.R(`.Unique()`)
		}
		if field.Bool(`ent_notempty`) {
			o.R(`.NotEmpty()`)
		}
		if field.Bool(`ent_sensitive`) {
			o.R(`.Sensitive()`)
		}
		o.R(`,`)
	}

	// For Users and Groups, we need to store/create ETags
	switch object.Name(true) {
	case `User`, `Group`:
		o.L(`field.String("etag").Optional(),`)
	default:
	}
	o.L(`}`)
	o.L(`}`)

	fn := filepath.Join(`ent`, `schema`, packageName(object.Name(false))+`_gen.go`)
	if err := o.WriteFile(fn, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return fmt.Errorf(`failed to write to %s: %w`, fn, err)
	}
	return nil
}

func generateLoadEntFields(dst io.Writer, object *codegen.Object) error {
	o := codegen.NewOutput(dst)

	o.LL(`func %sLoadEntFields(q *ent.%sQuery, scimFields, excludedFields []string) {`, object.Name(false), object.Name(true))
	o.L(`fields := make(map[string]struct{})`)
	o.L(`if len(scimFields) == 0 {`)
	o.L(`scimFields = []string {`)

	for i, field := range object.Fields() {
		switch field.Name(false) {
		case "schemas", "meta": // These are handled separately
			continue
		}
		if field.Bool(`skipCommonFields`) {
			switch field.Name(false) {
			case "id", "externalID": // these are only required when they are imported
				continue
			}
		}

		// Theoretically, there cold be any number of fields that
		// have the "returned" field set to `never` or `request`, but
		// in practice only password is set to never, and
		// there are no fields set to request (TODO: check again)
		if i > 0 {
			o.R(`,`)
		}
		o.R(`resource.%s%sKey`, object.Name(true), field.Name(true))
	}
	o.R(`}`)
	o.L(`}`)
	o.LL(`for _, name := range scimFields {`)
	// Theoretically we need to prevent the user from deleting
	// fields set to "always", but only "id" has this in practice
	o.L(`fields[name] = struct{}{}`)
	o.L(`}`)

	o.LL(`for _, name := range excludedFields {`)
	o.L(`delete(fields, name)`)
	o.L(`}`)

	o.L(`selectNames := make([]string, 0, len(fields))`)
	o.L(`for f := range fields {`)
	o.L(`switch f {`)
	for _, field := range object.Fields() {
		if object.Name(true) == `User` && field.Name(true) == `Groups` {
			continue
		}
		if field.Name(false) == "schemas" {
			continue
		}
		if field.Bool(`skipCommonFields`) {
			switch field.Name(false) {
			case "id", "externalID", "meta":
				continue
			}
		}

		o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
		if isEdge(object, field) || (object.Name(true) == `User` && field.Name(true) == `Name`) {
			o.L(`q.With%s()`, edgeName(field))
			continue
		}

		switch field.Name(true) {
		case `Meta`, `Members`:
		default:
			// Otherwise, accumulate in the list of names
			o.L(`selectNames = append(selectNames, %s.Field%s)`, object.Name(false), entName(field, true))
		}
	}
	o.L(`}`)
	o.L(`}`)
	// there are some fields that MUST exist
	switch object.Name(true) {
	case `User`, `Group`:
		o.L(`selectNames = append(selectNames, %s.FieldEtag)`, object.Name(false))
	}
	o.L(`q.Select(selectNames...)`)
	o.L(`}`)
	return nil
}

func generateResourceFromEnt(dst io.Writer, object *codegen.Object) error {
	o := codegen.NewOutput(dst)

	scimRsname := scimResourceName(object)
	entRsname := resourceName(object)

	o.LL(`func %[1]sResourceFromEnt(in *ent.%[2]s) (*resource.%[1]s, error) {`, scimRsname, entRsname)
	o.L(`var b resource.Builder`)

	o.LL(`builder := b.%s()`, scimRsname)

	if !object.Bool(`skipCommonFields`) {
		o.LL(`meta, err := b.Meta().`)
		o.L(`ResourceType(%q).`, object.Name(true))
		o.L(`Location(%sLocation(in.ID.String())).`, object.Name(false))
		o.L(`Version(in.Etag).`)
		o.L(`Build()`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to build meta information for %s")`, object.Name(true))
		o.L(`}`)
		o.LL(`builder.`)
		o.L(`Meta(meta)`)
	}

	for _, field := range object.Fields() {
		if ft := field.Type(); !strings.HasPrefix(ft, `[]`) && !strings.HasPrefix(ft, `*`) {
			continue
		}

		// This section is just really really confusing because not all
		// stored data map 1-to-1 to the SCIM resource (for example,
		// Group.members can't be expressed 1-to-1 in a straight forward
		// manner).
		// I think it's better if we give people an escape hatch, so we're
		// going to inject a call to a helper of your choice at the end.
		edgeName := edgeName(field)
		scimRsname := scimResourceName(field)
		entRsname := resourceName(field)

		switch scimRsname {
		case `Names`:
			entRsname = `Names`
			o.LL(`if el := in.Edges.Name; el != nil {`)
			o.L(`r, err := NamesResourceFromEnt(el)`)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to convert names to SCIM resource: %%w", err)`)
			o.L(`}`)
			o.L(`builder.%s(r)`, field.Name(true))
			o.L(`}`)
		case "Group", `Meta`:
			// no op. done in helper
		default:
			if isEdge(object, field) {
				o.LL(`if el := len(in.Edges.%s); el > 0 {`, edgeName)
				o.L(`list := make([]*resource.%s, 0, el)`, scimRsname)
				o.L(`for _, ine := range in.Edges.%s {`, edgeName)
				o.L(`r, err := %sResourceFromEnt(ine)`, entRsname)
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to build %s information for %s")`, field.Name(false), object.Name(true))
				o.L(`}`)
				o.L(`list = append(list, r)`)
				o.L(`}`)

				o.L(`builder.%s(list...)`, field.Name(true))
				o.L(`}`)
			}
		}
	}

	if object.Name(true) == `User` {
		// User has a special case... namely in.Groups, which is not technicaly
		// and edge, bt something that is dynamically computed outside of the
		// scope of entgo
		o.LL(`if gl := len(in.Groups); gl > 0 {`)
		o.L(`memberships := make([]*resource.GroupMember, gl)`)
		o.L(`for i, m := range in.Groups {`)
		o.L(`var gmb resource.GroupMemberBuilder`)
		o.L(`gm, err := gmb.Value(m.Value).`)
		// o.L(`Set("display", m.Display).`) // TODO
		o.L(`Ref(m.Ref).`)
		o.L(`Type(m.Type).`) // direct or indirect
		o.L(`Build()`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to compute \"groups\" field: %%w", err)`)
		o.L(`}`) // end of if err != nil
		o.L(`memberships[i] = gm`)
		o.L(`}`) // end for
		o.L(`builder.Groups(memberships...)`)
		o.L(`}`) // end of if gl := len...
	}

	for _, field := range object.Fields() {
		switch field.Name(true) {
		// FIXME: do't hard codethis
		case "Password":
			continue
		case "ID":
			o.L(`builder.%[1]s(in.%[1]s.String())`, field.Name(true))
		case "Schemas", "Meta", "Members", "Addresses", "Emails", "Entitlements", "IMS", "NickName", "Name", "Groups", "PhoneNumbers", "ProfileURL", "Roles", "X509Certificates", "Photos":
			// TODO: FIXME
		default:
			o.L(`if !reflect.ValueOf(in.%s).IsZero() {`, field.Name(true))
			o.L(`builder.%[1]s(in.%[1]s)`, field.Name(true))
			o.L(`}`)
		}
	}
	if h := object.String(`ent_conversion_helper`); h != "" {
		o.L(`if err := %s(in, builder); err != nil {`, h)
		o.L(`return nil, err`)
		o.L(`}`)
	}
	o.L(`return builder.Build()`)
	o.L(`}`)
	return nil
}

func generateEntFieldFromSCIM(dst io.Writer, object *codegen.Object) error {
	o := codegen.NewOutput(dst)

	o.LL(`func %sEntFieldFromSCIM(s string) string {`, object.Name(true))
	o.L(`switch s {`)
	for _, field := range object.Fields() {
		if strings.HasPrefix(field.Type(), `[]`) || strings.HasPrefix(field.Type(), `*`) {
			continue
		}
		switch field.Name(false) {
		case `schemas`:
			continue
		default:
		}
		o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
		o.L(`return %s.Field%s`, packageName(object.Name(false)), entName(field, true))
	}
	o.L(`default:`)
	o.L(`return s`)
	o.L(`}`)
	o.L(`}`)
	return nil
}

func generatePredicate(dst io.Writer, object *codegen.Object, scimMethod, entMethod string) error {
	o := codegen.NewOutput(dst)
	o.LL(`func %[1]s%[2]sPredicate(q *ent.%[3]sQuery, scimField string, val interface{}) (predicate.%[3]s, error) {`, object.Name(false), scimMethod, object.Name(true))
	o.L(`_ = q`) // in case the predicate doesn't actually need to use the query object
	// The scim field may either be a flat (simple) field or a nested field.
	o.L(`field, subfield, err := splitScimField(scimField)`)
	o.L(`if err != nil {`)
	o.L(`return nil, err`)
	o.L(`}`)
	o.L(`_ = subfield // TODO: remove later`)

	o.L(`switch field {`)
	for _, field := range object.Fields() {
		switch field.Name(false) {
		case `schemas`:
			continue
		default:
		}

		switch field.Type() {
		// predicates against a list actually means "... if any of the values match"
		// so things like `roles.value eq "foo"` means `if any of the role.value is equal to "foo"`
		case "[]*Role", "[]*Email", "[]*PhoneNumber":
			o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
			// It's going to be a relation, so add a query that goes into the separate entity
			o.L(`switch subfield {`)

			// TODO don't hardcode
			// We know at this point that this type is something like []*Foo, so extract the Foo
			// and get the object definition
			subObjectName := strings.TrimPrefix(field.Type(), `[]*`)
			subObject, ok := objectMap[subObjectName]
			if !ok {
				return fmt.Errorf(`could not find object %q`, subObjectName)
			}
			for _, subField := range subObject.Fields() {
				switch entMethod {
				case `EQ`:
					// anything goes
				default:
					// only strings allowed
					if subField.Type() != `string` {
						continue
					}
				}

				o.L(`case resource.%s%sKey:`, subObjectName, subField.Name(true))
				o.L(`//nolint:forcetypeassert`)
				o.L(`return %s.Has%sWith(%s.%s%s(val.(%s))), nil`, object.Name(false), field.Name(true), strings.ToLower(singularName(field.Name(false))), subField.Name(true), entMethod, subField.Type())
			}
			o.L(`default:`)
			o.L(`return nil, fmt.Errorf("invalid filter specification: invalid subfield for %%q", field)`)
			o.L(`}`)
		case "string":
			o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
			// We can't just use ${Field}HasPrefix here, because we're going to
			// receive the field name as a parameter
			o.L(`entFieldName := %sEntFieldFromSCIM(scimField)`, object.Name(true))
			o.L(`return predicate.%[1]s(func(s *sql.Selector) {`, object.Name(true))
			o.L(`//nolint:forcetypeassert`)
			o.L(`s.Where(sql.%s(s.C(entFieldName), val.(%s)))`, entMethod, field.Type())
			o.L(`}), nil`)
		}
	}
	o.L(`default:`)
	o.L(`return nil, fmt.Errorf("invalid filter field specification")`)
	o.L(`}`)
	o.L(`}`)
	return nil
}

func generatePresencePredicate(dst io.Writer, object *codegen.Object) error {
	o := codegen.NewOutput(dst)

	o.LL(`func %sPresencePredicate(scimField string) predicate.%s {`, object.Name(false), object.Name(true))
	o.L(`switch scimField {`)
	for _, field := range object.Fields() {
		switch field.Name(false) {
		case `schemas`:
			continue
		default:
		}
		if field.Type() != "string" {
			continue
		}
		if field.IsRequired() {
			continue
		}
		o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
		o.L(`return %[1]s.And(%[1]s.%[2]sNotNil(), %[1]s.%[2]sNEQ(""))`, packageName(object.Name(false)), field.Name(true))
	}
	o.L(`default:`)
	o.L(`return nil`)
	o.L(`}`)
	o.L(`}`)
	return nil
}

func generatePredicateBuilder(dst io.Writer, object *codegen.Object) error {
	o := codegen.NewOutput(dst)

	scimRsname := scimResourceName(object)
	entRsname := resourceName(object)

	// visits filter.Expr to build predicates for PATCH operations against GroupMembers
	// is only intended to parse the query portion (e.g. `members[HERE]`)
	o.LL(`type %sPredicateBuilder struct {`, entRsname)
	o.L(`predicates []predicate.%s`, entRsname)
	o.L(`}`)

	o.LL(`func (b *%[1]sPredicateBuilder) Build(expr filter.Expr) ([]predicate.%[1]s, error) {`, entRsname)
	o.L(`b.predicates = nil`)
	o.L(`if err := b.visit(expr); err != nil {`)
	o.L(`return nil, err`)
	o.L(`}`)
	o.L(`return b.predicates, nil`)
	o.L(`}`)

	o.LL(`func (b *%sPredicateBuilder) visit(expr filter.Expr) error {`, entRsname)
	o.L(`switch expr := expr.(type) {`)
	o.L(`case filter.CompareExpr:`)
	o.L(`return b.visitCompareExpr(expr)`)
	o.L(`case filter.LogExpr:`)
	o.L(`return b.visitLogExpr(expr)`)
	o.L(`default:`)
	o.L(`return fmt.Errorf("unhandled expression type %%T", expr)`)
	o.L(`}`)
	o.L(`}`)

	o.LL(`func (b *%sPredicateBuilder) visitLogExpr(expr filter.LogExpr) error {`, entRsname)
	o.L(`if err := b.visit(expr.LHE()); err != nil {`)
	o.L(`return fmt.Errorf("failed to parse left hand side of %%q statement: %%w", expr.Operator(), err)`)
	o.L(`}`)
	o.L(`if err := b.visit(expr.RHS()); err != nil {`)
	o.L(`return fmt.Errorf("failed to parse right hand side of %%q statement: %%w", expr.Operator(), err)`)
	o.L(`}`)
	o.LL(`switch expr.Operator() {`)
	o.L(`case "and":`)
	o.L(`b.predicates = []predicate.%s{%s.And(b.predicates...)}`, entRsname, packageName(entRsname))
	o.L(`case "or":`)
	o.L(`b.predicates = []predicate.%s{%s.Or(b.predicates...)}`, entRsname, packageName(entRsname))
	o.L(`default:`)
	o.L(`return fmt.Errorf("unhandled logical operator %%q", expr.Operator())`)
	o.L(`}`)
	o.L(`return nil`)
	o.L(`}`)

	o.LL(`func (b *%sPredicateBuilder) visitCompareExpr(expr filter.CompareExpr) error {`, entRsname)
	o.L(`lhe, err := exprAttr(expr.LHE())`)
	o.L(`slhe, ok := lhe.(string)`)
	o.L(`if err != nil || !ok {`)
	o.L(`return fmt.Errorf("left hand side of CompareExpr is not valid")`)
	o.L(`}`)
	o.LL(`rhe, err := exprAttr(expr.RHE())`)
	o.L(`if err != nil {`)
	o.L(`return fmt.Errorf("right hand side of CompareExpr is not valid: %%w", err)`)
	o.L(`}`)
	o.LL(`// convert rhe to string so it can be passed to regexp.QuoteMeta`)
	o.L(`srhe := fmt.Sprintf("%%v", rhe)`)
	o.LL(`switch expr.Operator() {`)
	o.L(`case filter.EqualOp:`)
	o.L(`switch slhe {`)
	for _, field := range object.Fields() {
		switch field.Name(true) {
		case `Schemas`, `Meta`:
			continue
		}

		o.L(`case resource.%s%sKey:`, scimRsname, field.Name(true))
		var localVar string
		if field.Type() == `bool` {
			localVar = `v`
			o.L(`v, err := strconv.ParseBool(srhe)`)
			o.L(`if err != nil {`)
			o.L(`return fmt.Errorf("failed to parse boolean expression")`)
			o.L(`}`)
		} else if field.Name(false) == `id` {
			localVar = `v`
			o.L(`v, err := uuid.Parse(srhe)`)
			o.L(`if err != nil {`)
			o.L(`return fmt.Errorf("failed to parse UUID")`)
			o.L(`}`)
		} else {
			localVar = `srhe`
		}

		o.L(`b.predicates = append(b.predicates, %s.%s(%s))`, packageName(entRsname), field.Name(true), localVar)
	}
	o.L(`default:`)
	o.L(`return fmt.Errorf("invalid field name for %s: %%q", slhe)`, object.Name(true))
	o.L(`}`)
	o.L(`default:`)
	o.L(`return fmt.Errorf("invalid operator: %%q", expr.Operator())`)
	o.L(`}`)
	o.L(`return nil`)
	o.L(`}`)
	return nil
}

func generateExistsPredicate(dst io.Writer, object *codegen.Object) error {
	o := codegen.NewOutput(dst)

	for _, field := range object.Fields() {
		if !strings.HasPrefix(field.Type(), `[]`) {
			continue
		}

		scimRsname := scimResourceName(field)
		rsname := resourceName(field)

		if object.Name(true) == `User` && scimRsname == `GroupMember` {
			continue
		}

		o.LL(`func (b *Backend) exists%[2]s%[1]s(ctx context.Context, parent *ent.%[2]s, in *resource.%[3]s) bool {`, rsname, object.Name(true), scimRsname)
		o.L(`queryCall := parent.Query%s()`, field.Name(true))

		subObject, ok := objectMap[scimRsname]
		if !ok {
			return fmt.Errorf(`could not locate object %q`, scimRsname)
		}

		o.L(`var predicates []predicate.%s`, rsname)
		for _, subField := range subObject.Fields() {
			o.L(`if in.Has%s() {`, subField.Name(true))
			o.L(`predicates = append(predicates, %[1]s.%[2]s(in.%[2]s()))`, packageName(rsname), subField.Name(true))
			o.L(`}`)
		}
		o.LL(`v, err := queryCall.Where(predicates...).Exist(ctx)`)
		o.L(`if err != nil {`)
		o.L(`return false`)
		o.L(`}`)
		o.L(`return v`)
		o.L(`}`)
	}
	return nil
}

func generateUtilities(object *codegen.Object) error {
	var buf bytes.Buffer
	o := codegen.NewOutput(&buf)

	lcObject := object.Name(false)

	o.L(`package server`)

	pkgs := []string{
		`github.com/cybozu-go/scim/resource`,
		`github.com/cybozu-go/scim-server/ent`,
	}
	for _, pkg := range []string{`predicate`, packageName(lcObject), `address`, `email`, `entitlement`, `group`, `ims`, `member`, `phonenumber`, `photo`, `user`, `role`, `x509certificate`} {
		pkgs = append(pkgs, fmt.Sprintf(`github.com/cybozu-go/scim-server/ent/%s`, pkg))
	}
	o.LL(`import (`)
	for _, pkg := range pkgs {
		o.L(`%q`, pkg)
	}
	o.L(`)`)

	if object.String(`schema`) != "" {
		if err := generateLoadEntFields(&buf, object); err != nil {
			return fmt.Errorf(`failed to generate XXXLoadEntFields utility: %w`, err)
		}
	}

	// TODO: prefix is hard coded, need to fix
	if !object.Bool(`skipCommonFields`) {
		o.LL(`func %sLocation(id string) string {`, object.Name(false))
		o.L(`return %q+id` /* TODO: FIXME */, fmt.Sprintf(`https://foobar.com/scim/v2/%ss/`, object.Name(true)))
		o.L(`}`)
	}

	if err := generateResourceFromEnt(&buf, object); err != nil {
		return fmt.Errorf(`failed to generate XXXResourceFromEnt utility: %w`, err)
	}

	if err := generateEntFieldFromSCIM(&buf, object); err != nil {
		return fmt.Errorf(`failed to generate XXXEntFieldFromSCIM utility: %w`, err)
	}

	switch object.Name(true) {
	case `User`, `Group`:
		for _, pred := range []struct {
			Name   string
			Method string
		}{
			{Name: `StartsWith`, Method: `HasPrefix`},
			{Name: `EndsWith`, Method: `HasSuffix`},
			{Name: `Contains`, Method: `Contains`},
			{Name: `Equals`, Method: `EQ`},
		} {
			if err := generatePredicate(&buf, object, pred.Name, pred.Method); err != nil {
				return fmt.Errorf(`failed to generate predicate for %s: %w`, pred.Name, err)
			}
		}

		if err := generatePresencePredicate(&buf, object); err != nil {
			return fmt.Errorf(`failed to generate predicate for presence: %w`, err)
		}
	default:
		if err := generatePredicateBuilder(&buf, object); err != nil {
			return fmt.Errorf(`failed to generate predicate builder: %w`, err)
		}
	}

	if err := generateExistsPredicate(&buf, object); err != nil {
		return fmt.Errorf(`failed to generate exists predicate: %w`, err)
	}

	if object.Name(true) == `User` || object.Name(true) == `Group` {
		for _, field := range object.Fields() {
			if !strings.HasPrefix(field.Type(), `[]`) {
				continue
			}

			// These are defined elsewhere
			switch field.Name(true) {
			// Groups cannot created/modified via User resource
			case `Addresses`, `Groups`:
				continue
			}
			scimRsname := scimResourceName(field)
			rsname := resourceName(field)

			o.LL(`func (b *Backend) create%[1]s(ctx context.Context, resources ...*resource.%[2]s) ([]*ent.%[1]sCreate, error) {`, rsname, scimRsname)
			o.L(`list := make([]*ent.%sCreate, len(resources))`, rsname)
			o.L(`for i, in := range resources {`)
			o.L(`createCall := b.db.%s.Create()`, rsname)
			var fields []string
			if rsname == `Member` {
				fields = []string{`Value`, `Type`, `Ref`}
			} else {
				fields = []string{"Display", "Primary", "Type", "Value"}
			}
			for _, subf := range fields {
				o.L(`if in.Has%s() {`, subf)
				o.L(`createCall.Set%[1]s(in.%[1]s())`, subf)

				if subf == `Type` && rsname == `Member` {
					o.L(`} else {`)
					o.LL(`parsedUUID, err := uuid.Parse(in.Value())`)
					o.L(`if err != nil {`)
					o.L(`return nil, fmt.Errorf("failed to parse ID in \"value\" field: %%w", err)`)
					o.L(`}`)
					o.L(`if ok, _ := b.db.User.Query().Where(user.ID(parsedUUID)).Exist(ctx); ok {`)
					o.L(`createCall.SetType("User")`)
					o.L(`} else if ok, _ := b.db.Group.Query().Where(group.ID(parsedUUID)).Exist(ctx); ok {`)
					o.L(`createCall.SetType("Group")`)
					o.L(`} else {`)
					o.L(`return nil, fmt.Errorf("could not determine resource type (User/Group) from provided ID")`)
					o.L(`}`)
				}
				o.L(`}`)

			}
			o.L(`list[i] = createCall`)
			o.L(`}`)
			o.L(`return list, nil`)
			o.L(`}`)
		}
	}

	var resourceSchema *resource.Schema
	if object.Name(true) == `User` {
		resourceSchema = userSchema
	} else {
		resourceSchema = groupSchema
	}
	switch object.Name(true) {
	case `User`, `Group`:
		var required []codegen.Field
		var optional []codegen.Field
		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `ID`, `Schemas`, `Meta`, `Password`:
				continue
			}

			if field.IsRequired() {
				required = append(required, field)
			} else {
				optional = append(optional, field)
			}
		}
		o.LL(`func (b *Backend) Create%[1]s(ctx context.Context, in *resource.%[1]s) (*resource.%[1]s, error) {`, object.Name(true))
		o.LL(`createCall := b.db.%s.Create()`, object.Name(true))

		if object.Name(true) == `User` {
			o.L(`password, err := b.generatePassword(in)`)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to process password: %%w", err)`)
			o.L(`}`)
			o.L(`createCall.SetPassword(password)`)
		}
		for _, field := range required {
			o.L(`if !in.Has%s() {`, field.Name(true))
			o.L(`return nil, fmt.Errorf("required field %s not found")`, field.JSON())
			o.L(`}`)
			o.L(`createCall.Set%[1]s(in.%[1]s())`, field.Name(true))
		}

		for _, field := range optional {
			if isEdge(object, field) {
				o.L(`var %[1]sCreateCalls []*ent.%[2]sCreate`, singularName(field.Name(false)), resourceName(field))
				o.L(`if in.Has%s() {`, field.Name(true))
				o.L(`calls, err := b.create%s(ctx, in.%s()...)`, resourceName(field), field.Name(true))
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to create %s: %%w", err)`, field.JSON())
				o.L(`}`)
				o.L(`%sCreateCalls = calls`, singularName(field.Name(false)))
				o.L(`}`)
			} else if object.Name(true) == `User` && field.Name(true) == `Groups` {
				continue
			} else if field.Name(true) == `Name` {
				o.L(`if in.Has%s() {`, field.Name(true))
				o.L(`created, err := b.create%[1]s(ctx, in.%[1]s())`, field.Name(true))
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to create %s: %%w", err)`, field.JSON())
				o.L(`}`)
				o.L(`createCall.Set%s(created)`, field.Name(true))
				o.L(`}`)
			} else {
				o.L(`if in.Has%s() {`, field.Name(true))
				o.L(`createCall.Set%[1]s(in.%[1]s())`, field.Name(true))
				o.L(`}`)
			}
		}

		o.LL(`rs, err := createCall.Save(ctx)`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to save object: %%w", err)`)
		o.L(`}`)

		for _, field := range optional {
			if !strings.HasPrefix(field.Type(), `[]`) || field.Name(true) == `Groups` {
				continue
			}
			o.LL(`for _, call := range %sCreateCalls {`, singularName(field.Name(false)))
			o.L(`call.Set%sID(rs.ID)`, object.Name(true))
			o.L(`created, err := call.Save(ctx)`)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to create %s: %%w", err)`, field.JSON())
			o.L(`}`)
			o.L(`rs.Edges.%[1]s = append(rs.Edges.%[1]s, created)`, edgeName(field))
			o.L(`}`)
		}

		o.LL(`h := sha256.New()`)
		o.L(`if err := rs.ComputeETag(h); err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to compute etag: %%w", err)`)
		o.L(`}`)
		o.L(`etag := fmt.Sprintf("W/%%x", h.Sum(nil))`)
		o.LL(`if _, err := rs.Update().SetEtag(etag).Save(ctx); err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to save etag: %%w", err)`)
		o.L(`}`)
		o.L(`rs.Etag = etag`)
		o.L(`return %sResourceFromEnt(rs)`, object.Name(true))
		o.L(`}`)

		o.LL(`func (b *Backend) Replace%[1]s(ctx context.Context, id string, in *resource.%[1]s) (*resource.%[1]s, error) {`, object.Name(true))
		o.LL(`parsedUUID, err := uuid.Parse(id)`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to parse ID: %%w", err)`)
		o.L(`}`)
		o.LL(`r, err := b.db.%s.Query().Where(%s.ID(parsedUUID)).Only(ctx)`, object.Name(true), packageName(object.Name(false)))
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to retrieve resource for replacing: %%w", err)`)
		o.L(`}`)

		// TODO: THIS IS NOT THE RIGHT IMPLEMENTATION
		o.LL(`replaceCall := r.Update()`)
		for _, field := range object.Fields() {
			attr, ok := resourceSchema.AttributeByName(field.JSON())
			if !ok {
				continue
			}

			if !attr.CanWrite() {
				continue
			}

			// UserName cannot be empty
			if field.Name(true) != `UserName` {
				o.LL(`replaceCall.Clear%s()`, field.Name(true))
			}

			if isEdge(object, field) {
				entRsname := resourceName(field)

				o.L(`var %sCreateCalls []*ent.%sCreate`, field.Name(false), entRsname)
				o.L(`if in.Has%s() {`, field.Name(true))
				o.L(`calls, err := b.create%s(ctx, in.%s()...)`, entRsname, field.Name(true))
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to create %s: %%w", err)`, field.JSON())
				o.L(`}`)
				o.L(`%sCreateCalls = calls`, field.Name(false))
				o.L(`}`)
			} else if field.Name(true) == `Name` {
				o.L(`if in.Has%s() {`, field.Name(true))
				o.L(`created, err := b.createName(ctx, in.Name())`)
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to create name: %%w", err)`)
				o.L(`}`)
				o.L(`replaceCall.SetName(created)`)
				o.L(`}`)
			} else {
				o.L(`if in.Has%s() {`, field.Name(true))
				o.L(`replaceCall.Set%[1]s(in.%[1]s())`, field.Name(true))
				o.L(`}`)
			}
		}
		o.LL(`if _, err := replaceCall.Save(ctx); err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to save object: %%w", err)`)
		o.L(`}`)

		// Do the Edge updates AFTER we saved the results
		for _, field := range object.Fields() {
			if !isEdge(object, field) {
				continue
			}
			if !strings.HasPrefix(field.Type(), `[]`) {
				continue
			}

			o.L(`for _, call := range %sCreateCalls {`, field.Name(false))
			o.L(`call.Set%sID(parsedUUID)`, object.Name(true))
			o.L(`_, err := call.Save(ctx)`)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to save %s: %%w", err)`, field.Name(true))
			o.L(`}`)
			o.L(`}`)
		}

		o.LL(`r2, err := b.db.%s.Query().Where(%s.ID(parsedUUID)).`, object.Name(true), packageName(object.Name(false)))
		for _, field := range object.Fields() {
			if !isEdge(object, field) || field.Name(true) == `Names` {
				continue
			}
			o.L(`With%s().`, field.Name(true))
		}
		o.L(`Only(ctx)`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to retrieve data")`)
		o.L(`}`)
		o.LL(`h := sha256.New()`)
		o.L(`if err := r2.ComputeETag(h); err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to compute etag: %%w", err)`)
		o.L(`}`)
		o.L(`etag := fmt.Sprintf("W/%%x", h.Sum(nil))`)
		o.LL(`if _, err := r2.Update().SetEtag(etag).Save(ctx); err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to save etag: %%w", err)`)
		o.L(`}`)
		o.L(`r2.Etag = etag`)

		o.LL(`return %sResourceFromEnt(r2)`, object.Name(true))
		o.L(`}`)

		o.LL(`func (b *Backend) patchAdd%[1]s(ctx context.Context, parent *ent.%[1]s, op *resource.PatchOperation) error {`, object.Name(true))
		o.LL(`root, err := filter.Parse(op.Path(), filter.WithPatchExpression(true))`)
		o.L(`if err != nil {`)
		o.L(`return fmt.Errorf("failed to parse PATH path %%q", op.Path())`)
		o.L(`}`)
		o.LL(`expr, ok := root.(filter.ValuePath)`)
		o.L(`if !ok {`)
		o.L(`return fmt.Errorf("root element should be a valuePath (got %%T)", root)`)
		o.L(`}`)
		o.LL(`sattr, err := exprStr(expr.ParentAttr())`)
		o.L(`if err != nil {`)
		o.L(`return fmt.Errorf("invalid attribute specification: %%w", err)`)
		o.L(`}`)

		o.LL(`switch sattr {`)
		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `ID`, `Schema`, `Meta`:
				continue
			}

			if object.Name(true) == `User` && field.Name(true) == `Groups` {
				o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
				o.L(`return fmt.Errorf("cannot create group memberships through User resource")`)
			} else if field.Type() == `string` {
				o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
				o.L(`subExpr := expr.SubExpr()`) //
				o.L(`if subExpr != nil {`)
				o.L(`return fmt.Errorf("subexpr on string element is unimplmented")`)
				o.L(`}`)

				o.LL(`if expr.SubAttr() != nil {`)
				o.L(`return fmt.Errorf("invalid sub attrribute on string element %s")`, field.JSON())
				o.L(`}`)

				o.LL(`var v string`)
				o.L(`if err := json.Unmarshal(op.Value(), &v); err != nil {`)
				o.L(`return fmt.Errorf("invalid value for string element %s")`, field.JSON())
				o.L(`}`)
				o.LL(`if _, err := parent.Update().Set%s(v).Save(ctx); err != nil {`, field.Name(true))
				o.L(`return fmt.Errorf("failed to save object: %%w", err)`)
				o.L(`}`)
			} else if strings.HasPrefix(field.Type(), `[]`) {
				o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))
				o.L(`subExpr := expr.SubExpr()`)
				// there's  no query
				o.L(`if subExpr == nil {`)
				o.L(`if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {`)
				// Adding to a subAttr of a multi-value element doesn't make sense
				o.L(`return fmt.Errorf("patch add operation on sub attribute of multi-value item %s with unspecified element is not possible")`, field.JSON())
				o.L(`}`)

				// if we're adding to the list, we need the entire thing
				scimRsname := scimResourceName(field)
				rsname := resourceName(field)
				o.LL(`var in resource.%s`, scimRsname)
				o.L(`if err := json.Unmarshal(op.Value(), &in); err != nil {`)
				o.L(`return fmt.Errorf("failed to decode patch add value: %%w", err)`)
				o.L(`}`)

				o.LL(`if b.exists%s%s(ctx, parent, &in) {`, object.Name(true), rsname)
				o.L(`return nil`)
				o.L(`}`)

				o.LL(`calls, err := b.create%s(ctx, &in)`, rsname)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("failed to create %s: %%w", err)`, rsname)
				o.L(`}`)

				o.L(`list := make([]*ent.%s, len(calls))`, rsname)
				o.L(`for i, call := range calls {`)
				o.L(`call.Set%sID(parent.ID)`, object.Name(true))
				o.L(`created, err := call.Save(ctx)`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("failed to create %s: %%w", err)`, rsname)
				o.L(`}`)
				o.L(`list[i] = created`)
				o.L(`}`)

				/*
					o.LL(`if _, err := parent.Update().%s(list...).Save(ctx); err != nil {`, addMethod(field))
					o.L(`return fmt.Errorf("failed to save object: %%w", err)`)
					o.L(`}`)
				*/
				o.L(`} else {`)
				o.L(`var pb %sPredicateBuilder`, rsname)
				// so we have a subExpr, that must mean we must have have a subAttr
				// "path": "members[value eq \"...\"].value"  // OK
				// "path": "members[value eq \"...\"]"        // NOT OK
				// also, the query must resolve to a single member element
				// Load attr with the given conditions
				o.L(`predicates, err := pb.Build(subExpr)`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("failed to parse valuePath expression: %%w", err)`)
				o.L(`}`)
				o.L(`list, err := parent.Query%s().`, field.Name(true))
				o.L(`Where(predicates...).`)
				o.L(`All(ctx)`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("failed to look up value: %%w", err)`)
				o.L(`}`)

				o.LL(`if len(list) > 0 {`)
				o.L(`return fmt.Errorf("query must resolve to one element, got multiple")`)
				o.L(`}`)

				o.LL(`item := list[0]`)
				// we must have subAttr
				o.L(`sSubAttr, err := exprStr(expr.SubAttr())`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("query must have a sub attribute")`)
				o.L(`}`)

				o.LL(`updateCall := item.Update()`)
				o.LL(`switch sSubAttr {`)
				subObject, ok := objectMap[scimRsname]
				if !ok {
					return fmt.Errorf(`could not find object for %q`, scimRsname)
				}
				for _, subField := range subObject.Fields() {
					// TODO check for mutability
					o.L(`case resource.%s%sKey:`, scimRsname, subField.Name(true))
					o.L(`var v %s`, subField.Type())
					o.L(`if err := json.Unmarshal(op.Value(), &v); err != nil {`)
					o.L(`return fmt.Errorf("failed to decode value: %%w", err)`)
					o.L(`}`)
					o.L(`updateCall.Set%s(v)`, subField.Name(true))
				}
				o.L(`}`) // switch sSubAttr
				o.LL(`if _, err := updateCall.Save(ctx); err != nil {`)
				o.L(`return fmt.Errorf("failed to save object: %%w", err)`)
				o.L(`}`)
				o.L(`return nil`)
				o.L(`}`) // else
			}
		}
		o.L(`}`) // switch sattr
		o.L(`return nil`)
		o.L(`}`) // patchAdd%[1]s

		o.LL(`func (b *Backend) patchRemove%[1]s(ctx context.Context, parent *ent.%[1]s, op *resource.PatchOperation) error {`, object.Name(true))
		o.L(`if op.Path() == "" {`)
		o.L(`return resource.NewErrorBuilder().`)
		o.L(`Status(http.StatusBadRequest).`)
		o.L(`ScimType(resource.ErrNoTarget).`)
		o.L(`Detail("empty path").`)
		o.L(`MustBuild()`)
		o.L(`}`)
		o.LL(`root, err := filter.Parse(op.Path(), filter.WithPatchExpression(true))`)
		o.L(`if err != nil {`)
		o.L(`return fmt.Errorf("failed to parse path %%q", op.Path())`)
		o.L(`}`)
		o.LL(`expr, ok := root.(filter.ValuePath)`)
		o.L(`if !ok {`)
		o.L(`return fmt.Errorf("root element should be a valuePath (got %%T)", root)`)
		o.L(`}`)
		o.LL(`sattr, err := exprStr(expr.ParentAttr())`)
		o.L(`if err != nil {`)
		o.L(`return fmt.Errorf("invalid attribute specification: %%w", err)`)
		o.L(`}`)

		o.L(`switch sattr {`)
		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `ID`, `Meta`, `Schemas`, `UserName`:
				continue
			}
			o.L(`case resource.%s%sKey:`, object.Name(true), field.Name(true))

			if object.Name(true) == `User` && field.Name(true) == `Groups` {
				o.L(`return fmt.Errorf("cannot delete group memberships through User resource")`)
			} else if field.Type() == `string` {
				o.L(`if subexpr := expr.SubExpr(); subexpr != nil {`)
				o.L(`return fmt.Errorf("patch remove operation on %s cannot have a sub attribute query")`, field.JSON())
				o.L(`}`)
				o.LL(`if subattr := expr.SubAttr(); subattr != nil {`)
				o.L(`return fmt.Errorf("patch remove operation on %s cannot have a sub attribute")`, field.JSON())
				o.L(`}`)
				o.LL(`if _, err := parent.Update().Clear%s().Save(ctx); err != nil {`, field.Name(true))
				o.L(`return fmt.Errorf("failed to save object: %%w", err)`)
				o.L(`}`)
			} else if strings.HasPrefix(field.Type(), `[]`) {
				scimRsname := scimResourceName(field)
				rsname := resourceName(field)

				o.L(`subExpr := expr.SubExpr()`)

				// This means no query, so we can't specify which item in the multi-value element we're dealing with.
				o.L(`if subExpr == nil {`)
				// This means we have `attr.subAttr` to remove.
				// In this case removing a subAttr of a multi-value element doesn't make sense
				o.L(`if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {`)
				o.L(`return fmt.Errorf("patch remove operation on su attribute of multi-valued item %s without a query is not possible")`, field.JSON())
				o.L(`}`)
				// This means we have `attr` to remove. clear the entire thing
				o.L(`if _, err := b.db.%s.Delete().Where(%s.Has%sWith(%s.ID(parent.ID))).Exec(ctx); err != nil {`, rsname, packageName(rsname), object.Name(true), packageName(singularName(object.Name(false))))
				o.L(`return fmt.Errorf("failed to remove elements from %s: %%w", err)`, field.JSON())
				o.L(`}`)
				o.L(`if _, err := parent.Update().%s().Save(ctx); err != nil {`, clearMethod(field))
				o.L(`return fmt.Errorf("failed to remove references to %s: %%w", err)`, field.JSON())
				o.L(`}`)
				o.L(`} else {`) // subExpr == nil
				o.L(`var pb %sPredicateBuilder`, rsname)
				o.L(`predicates, err := pb.Build(subExpr)`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("failed to parse valuePath expression: %%w", err)`)
				o.L(`}`)
				o.LL(`list, err := parent.%s().`, queryMethod(field))
				o.L(`Where(predicates...).`)
				o.L(`All(ctx)`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("failed to query context object: %%w", err)`)
				o.L(`}`)
				o.LL(`if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {`)
				o.L(`subAttr, err := exprStr(subAttrExpr)`)
				o.L(`if err != nil {`)
				o.L(`return fmt.Errorf("invalid sub attribute specified")`)
				o.L(`}`)
				o.L(`switch subAttr {`)

				subObject, ok := objectMap[scimRsname]
				if !ok {
					return fmt.Errorf(`failed to find object %q`, scimRsname)
				}
				for _, subField := range subObject.Fields() {
					o.L(`case resource.%s%sKey:`, subObject.Name(true), subField.Name(true))
					if isMutable(subObject, subField) {

					} else {
						o.L(`return fmt.Errorf("%s is not mutable")`, subField.JSON())
					}
				}
				o.L(`default:`)
				o.L(`return fmt.Errorf("unknown sub attribute specified")`)
				o.L(`}`)
				o.L(`}`)
				if subObject.Name(true) == `GroupMember` {
					o.LL(`ids := make([]int, len(list))`)
				} else {
					o.LL(`ids := make([]uuid.UUID, len(list))`)
				}
				o.L(`for i, elem := range list {`)
				o.L(`ids[i] = elem.ID`)
				o.L(`}`)

				o.L(`if _, err := b.db.%s.Delete().Where(%s.IDIn(ids...)).Exec(ctx); err != nil {`, rsname, packageName(rsname))
				o.L(`return fmt.Errorf("failed to delete object: %%w", err)`)
				o.L(`}`)
				o.L(`}`) // subExpr == nil
			}
		}
		o.L(`}`) // switch sattr
		o.L(`return nil`)
		o.L(`}`) // func
	}

	fn := fmt.Sprintf(`%s_gen.go`, relationFilename(object.Name(false)))
	if err := o.WriteFile(fn, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return fmt.Errorf(`failed to write to %s: %w`, fn, err)
	}
	return nil
}
