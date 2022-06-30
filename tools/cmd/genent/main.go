package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/lestrrat-go/codegen"
	"github.com/lestrrat-go/runcmd"
	"github.com/lestrrat-go/xstrings"
)

var objects map[string]*codegen.Object
var userEdges = []string{
	`Name`,
	`Emails`,
	`Roles`,
	`Groups`,
	`PhoneNumbers`,
	`Entitlements`,
	`IMS`,
	`Photos`,
	`Addresses`,
	`X509Certificates`,
}

var userEdgeMap map[string]struct{}

func main() {
	userEdgeMap = make(map[string]struct{})
	for _, e := range userEdges {
		userEdgeMap[e] = struct{}{}
	}

	objects = make(map[string]*codegen.Object)
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func isUserEdge(field codegen.Field) bool {
	_, ok := userEdgeMap[field.Name(true)]
	return ok
}

func isEdge(object *codegen.Object, field codegen.Field) bool {
	switch object.Name(true) {
	case `User`:
		return isUserEdge(field)
	}
	return false
}

func edgeName(field codegen.Field) string {
	n := field.Name(true)
	if n == "IMS" {
		return "Imses"
	}
	return n
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

		objects[object.Name(true)] = object
	}

	for _, object := range def.Objects {
		if err := generateEnt(object); err != nil {
			return fmt.Errorf(`failed to generate ent adapter: %s`, err)
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
	case `User`, `Group`, `Email`, `Names`, `Role`, `Photo`, `IMS`, `PhoneNumber`, `Address`, `Entitlement`, `X509Certificate`:
	default:
		return nil
	}

	fmt.Printf("  ⌛ Generating ent adapters for %s...\n", object.Name(true))

	if err := generateSchema(object); err != nil {
		return fmt.Errorf(`failed to generate schema: %w`, err)
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
	return strings.ToLower(s)
}

func generateSchema(object *codegen.Object) error {
	var buf bytes.Buffer
	o := codegen.NewOutput(&buf)

	o.L(`package schema`)

	o.LL(`type %s struct {`, object.Name(true))
	o.L(`ent.Schema`)
	o.L(`}`)

	o.LL(`func (%s) Fields() []ent.Field {`, object.Name(true))
	o.L(`return []ent.Field{`)
	for _, field := range object.Fields() {
		if field.Name(false) == "schemas" {
			continue
		}

		ft := field.Type()
		if strings.HasPrefix(ft, `[]`) || strings.HasPrefix(ft, `*`) {
			continue
		}

		var entMethod = xstrings.Camel(ft)
		if v := field.String(`ent_build_method`); v != "" {
			entMethod = v
		}

		var entName = field.Name(false)
		if v := field.String(`ent_name`); v != "" {
			entName = v
		}
		var entType = field.String(`ent_type`)
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
		o.L(`field.String("etag").NotEmpty(),`)
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

func generateUtilities(object *codegen.Object) error {
	var buf bytes.Buffer
	o := codegen.NewOutput(&buf)

	o.L(`package server`)

	o.LL(`import (`)
	o.L(`"github.com/cybozu-go/scim/resource"`)
	o.L(`"github.com/cybozu-go/scim-server/ent"`)
	o.L(`"github.com/cybozu-go/scim-server/ent/predicate"`)
	o.L(`"github.com/cybozu-go/scim-server/ent/%s"`, packageName(object.Name(false)))
	o.L(`)`)

	if object.String(`schema`) != "" {
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
			if isEdge(object, field) {
				o.L(`q.With%s()`, edgeName(field))
				continue
			}

			switch field.Name(true) {
			case `Meta`, `Members`:
			default:
				// Otherwise, accumulate in the list of names
				o.L(`selectNames = append(selectNames, %s.Field%s)`, object.Name(false), field.Name(true))
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
	}

	// TODO: prefix is hard coded, need to fix
	if !object.Bool(`skipCommonFields`) {
		o.LL(`func %sLocation(id string) string {`, object.Name(false))
		o.L(`return %q+id` /* TODO: FIXME */, fmt.Sprintf(`https://foobar.com/scim/v2/%ss/`, object.Name(true)))
		o.L(`}`)
	}

	o.LL(`func %[1]sResourceFromEnt(in *ent.%[1]s) (*resource.%[1]s, error) {`, object.Name(true))
	o.L(`var b resource.Builder`)

	o.LL(`builder := b.%s()`, object.Name(true))

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
		edgeName := field.Name(true)
		rsname := singularName(field.Name(true))
		if rsname == "Member" || rsname == "Meta" {
			continue
		}

		switch rsname {
		case "Name":
			rsname = "Names"
			o.LL(`if el := in.Edges.Name; el != nil {`)
			o.L(`r, err := NamesResourceFromEnt(el)`)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to convert names to SCIM resource: %%w", err)`)
			o.L(`}`)
			o.L(`builder.%s(r)`, field.Name(true))
			o.L(`}`)
		case "Group":
			// no op. done in helper
		case "IMS":
			edgeName = "Imses"
			fallthrough
		default:
			o.LL(`if el := len(in.Edges.%s); el > 0 {`, edgeName)
			o.L(`list := make([]*resource.%s, 0, el)`, rsname)
			o.L(`for _, ine := range in.Edges.%s {`, edgeName)
			o.L(`r, err := %sResourceFromEnt(ine)`, rsname)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to build %s information for %s")`, field.Name(false), object.Name(true))
			o.L(`}`)
			o.L(`list = append(list, r)`)
			o.L(`}`)

			o.L(`builder.%s(list...)`, field.Name(true))
			o.L(`}`)
		}
	}

	for _, field := range object.Fields() {
		switch field.Name(true) {
		// FIXME: do't hard codethis
		case "Password":
			continue
		case "ID":
			o.L(`builder.%[1]s(in.%[1]s.String())`, field.Name(true))
		case "Schemas", "Meta", "Members", "Addresses", "Emails", "Entitlements", "IMS", "NickName", "Name", "Groups", "PhoneNumbers", "ProfileURL", "Title", "Roles", "X509Certificates", "Photos":
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
		o.L(`return %s.Field%s`, packageName(object.Name(false)), field.Name(true))
	}
	o.L(`default:`)
	o.L(`return s`)
	o.L(`}`)
	o.L(`}`)

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
			o.LL(`func %[1]s%[2]sPredicate(q *ent.%[3]sQuery, scimField string, val interface{}) (predicate.%[3]s, error) {`, object.Name(false), pred.Name, object.Name(true))
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
					subObject, ok := objects[subObjectName]
					if !ok {
						return fmt.Errorf(`could not find object %q`, subObjectName)
					}
					for _, subField := range subObject.Fields() {
						switch pred.Method {
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
						o.L(`return %s.Has%sWith(%s.%s%s(val.(%s))), nil`, object.Name(false), field.Name(true), strings.ToLower(singularName(field.Name(false))), subField.Name(true), pred.Method, subField.Type())
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
					o.L(`s.Where(sql.%s(s.C(entFieldName), val.(%s)))`, pred.Method, field.Type())
					o.L(`}), nil`)
				}
			}
			o.L(`default:`)
			o.L(`return nil, fmt.Errorf("invalid filter field specification")`)
			o.L(`}`)
			o.L(`}`)
		}

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
	}

	if object.Name(true) == `User` {
		o.LL(`func (b *Backend) ReplaceUser(id string, in *resource.User) (*resource.User, error) {`)
		o.L(`parsedUUID, err := uuid.Parse(id)`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to parse ID: %%w", err)`)
		o.L(`}`)

		o.LL(`h := sha256.New()`)
		o.L(`fmt.Fprint(h, b.etagSalt)`)

		o.LL(`u, err := b.db.User.Query().`)
		o.L(`Select("id").`)
		o.L(`Where(user.IDEQ(parsedUUID)).`)
		o.L(`Only(context.TODO())`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to retrieve user: %%w", err)`)
		o.L(`}`)

		o.LL(`replaceUserCall := u.Update().`)
		edgeFields := make([]codegen.Field, 0, len(object.Fields()))
		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `Meta`:
				continue
			default:
			}
			if ft := field.Type(); strings.HasPrefix(ft, `[]`) || strings.HasPrefix(ft, `*`) {
				edgeFields = append(edgeFields, field)
			}
		}

		for i, field := range edgeFields {
			if i > 0 {
				o.R(`.`)
			}
			edgeName := field.Name(true)
			if edgeName == "IMS" {
				edgeName = "Imses"
			}
			o.L(`Clear%s()`, edgeName)
		}

		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `ID`, `Groups`, `Meta`, `Schemas`: // can't change this
			case `Name`:
				o.LL(`var name *ent.Names`)
				o.L(`if in.HasName() {`)
				o.L(`created, err := b.createName(in.Name(), h)`)
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to create name: %%w", err)`)
				o.L(`}`)
				o.L(`replaceUserCall.SetName(created)`)
				o.L(`name = created`)
				o.L(`}`)
			default:
				if isUserEdge(field) {
					o.LL(`var %s []*ent.%s`, field.Name(false), singularName(field.Name(true)))
					o.L(`if in.Has%s() {`, field.Name(true))
					o.L(`created, err := b.create%s(in, h)`, field.Name(true))
					o.L(`if err != nil {`)
					o.L(`return nil, fmt.Errorf("failed to create %s: %%w", err)`, singularName(field.Name(false)))
					o.L(`}`)
					o.L(`replaceUserCall.Add%s(created...)`, edgeName(field))
					o.L(`%s = created`, field.Name(false))
					o.L(`}`)
				} else {
					o.LL(`if in.Has%s() {`, field.Name(true))
					o.L(`replaceUserCall.Set%[1]s(in.%[1]s())`, field.Name(true))
					o.L(`fmt.Fprint(h, in.%s())`, field.Name(true))
					o.L(`}`)
				}
			}
		}
		o.LL(`replaceUserCall.SetEtag(fmt.Sprintf("W/%%q", base64.RawStdEncoding.EncodeToString(h.Sum(nil))))`)
		o.LL(`u2, err := replaceUserCall.`)
		o.L(`Save(context.TODO())`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to save user: %%w", err)`)
		o.L(`}`)

		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `Groups`:
				continue
			}
			if isUserEdge(field) {
				o.L(`u2.Edges.%s = %s`, edgeName(field), field.Name(false))
			}
		}
		o.LL(`return UserResourceFromEnt(u2)`)
		o.L(`}`)

		o.LL(`func (b *Backend) CreateUser(in *resource.User) (*resource.User, error) {`)
		o.L(`password, err := b.generatePassword(in)`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to process password: %%w", err)`)
		o.L(`}`)

		o.LL(`h := sha256.New()`)
		o.L(`fmt.Fprint(h, b.etagSalt)`)
		o.LL(`createUserCall := b.db.User.Create().`)

		requiredFields := make([]codegen.Field, 0, len(object.Fields()))
		for _, field := range object.Fields() {
			if field.Name(true) == `ID` {
				continue
			}

			if field.IsRequired() {
				requiredFields = append(requiredFields, field)
			}
		}

		for i, field := range requiredFields {
			if i > 0 {
				o.R(`.`)
			}
			o.L(`Set%[1]s(in.%[1]s())`, field.Name(true))
		}
		if len(requiredFields) > 0 {
			o.R(`.`)
		}
		o.L(`SetPassword(password)`)
		for _, field := range requiredFields {
			o.L(`fmt.Fprint(h, in.%s())`, field.Name(true))
		}

		for _, field := range object.Fields() {
			if field.IsRequired() {
				continue
			}
			switch field.Name(true) {
			case `Password`, `Meta`, `Schemas`, `Groups`:
			case `Name`:
				o.LL(`var name *ent.Names`)
				o.L(`if in.HasName() {`)
				o.L(`created, err := b.createName(in.Name(), h)`)
				o.L(`if err != nil {`)
				o.L(`return nil, fmt.Errorf("failed to create name: %%w", err)`)
				o.L(`}`)
				o.L(`createUserCall.SetName(created)`)
				o.L(`name = created`)
				o.L(`}`)
			default:
				if isUserEdge(field) {
					// TODO: add `X509Certificates` later
					o.LL(`var %s []*ent.%s`, field.Name(false), singularName(field.Name(true)))
					o.L(`if in.Has%s() {`, field.Name(true))
					o.L(`created, err := b.create%s(in, h)`, field.Name(true))
					o.L(`if err != nil {`)
					o.L(`return nil, fmt.Errorf("failed to create roles: %%w", err)`)
					o.L(`}`)
					o.L(`createUserCall.Add%s(created...)`, edgeName(field))
					o.L(`%s = created`, field.Name(false))
					o.L(`}`)
				} else {
					o.LL(`if in.Has%s() {`, field.Name(true))
					o.L(`createUserCall.Set%[1]s(in.%[1]s())`, field.Name(true))
					o.L(`fmt.Fprint(h, in.%s())`, field.Name(true))
					o.L(`}`)
				}
			}
		}

		o.LL(`createUserCall.SetEtag(fmt.Sprintf("W/%%q", base64.RawStdEncoding.EncodeToString(h.Sum(nil))))`)
		o.LL(`u, err := createUserCall.`)
		o.L(`Save(context.TODO())`)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to save user: %%w", err)`)
		o.L(`}`)

		for _, field := range object.Fields() {
			switch field.Name(true) {
			case `Groups`:
				continue
			}
			if isUserEdge(field) {
				o.L(`u.Edges.%s = %s`, edgeName(field), field.Name(false))
			}
		}
		o.LL(`return UserResourceFromEnt(u)`)
		o.L(`}`)

		// Email, Roles, PhoneNumbers, Certficates, Entitlement, IMS, Photo
		// all require the same type of helpers
		for _, field := range object.Fields() {
			if !isUserEdge(field) {
				continue
			}
			switch field.Name(true) {
			case `Name`, `Addresses`, `Groups`:
				continue
			}

			typ := singularName(field.Name(true))
			o.LL(`func (b *Backend) create%s(in *resource.User, h hash.Hash) ([]*ent.%s, error) {`, field.Name(true), typ)
			o.L(`list := make([]*ent.%s, len(in.%s()))`, singularName(field.Name(true)), field.Name(true))
			o.L(`inbound := in.%s()`, field.Name(true))
			o.L(`sort.Slice(inbound, func(i, j int) bool {`)
			o.L(`return inbound[i].Value() <= inbound[j].Value()`)
			o.L(`})`)

			o.LL(`var hasPrimary bool`)
			o.L(`for i, v := range inbound {`)
			o.L(`createCall := b.db.%s.Create()`, singularName(field.Name(true)))
			o.L(`createCall.SetValue(v.Value())`)
			o.L(`fmt.Fprint(h, v.Value())`)
			o.LL(`if v.HasDisplay() {`)
			o.L(`createCall.SetDisplay(v.Display())`)
			o.L(`fmt.Fprint(h, v.Display())`)
			o.L(`}`)
			o.LL(`if v.HasType() {`)
			o.L(`createCall.SetType(v.Type())`)
			o.L(`fmt.Fprint(h, v.Type())`)
			o.L(`}`)

			o.LL(`if sv := v.Primary(); sv {`)
			o.L(`if hasPrimary {`)
			o.L(`return nil, fmt.Errorf("invalid user.%[1]s: multiple %[1]s have been set to primary")`, field.JSON())
			o.L(`}`)
			o.L(`createCall.SetPrimary(true)`)
			o.L(`fmt.Fprint(h, []byte{1})`)
			o.L(`hasPrimary = true`)
			o.L(`} else {`)
			o.L(`fmt.Fprint(h, []byte{0})`)
			o.L(`}`)

			o.LL(`r, err := createCall.Save(context.TODO())`)
			o.L(`if err != nil {`)
			o.L(`return nil, fmt.Errorf("failed to save email %%d: %%w", i, err)`)
			o.L(`}`)

			o.LL(`list[i] = r`)
			o.L(`}`)
			o.L(`return list, nil`)
			o.L(`}`)
		}
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

/*
func generateSimpleEdge(field codegen.Field) error {
	// If this is a simple one-to-many, we can just generate the schema
	var o2buf bytes.Buffer
	o2 := codegen.NewOutput(&o2buf)
	structName := singularName(field.Name(true))
	fmt.Printf("  ⌛ Generating simple edge adapters for %s...\n", structName)
	o2.L(`package schema`)
	o2.LL(`type %s struct {`, structName)
	o2.L(`ent.Schema`)
	o2.L(`}`)
	o2.LL(`func (%s) Fields() []ent.Field {`, structName)
	o2.L(`return []ent.Field{`)
	o2.L(`field.String("value"),`)
	o2.L(`field.String("display"),`)
	o2.L(`field.String("type"),`)
	o2.L(`field.Bool("primary"),`)
	o2.L(`}`)
	o2.L(`}`)
	o2.LL(`func (%s) Edges() []ent.Edge {`, structName)
	o2.L(`return []ent.Edge{`)
	o2.L(`edge.To("user", User.Type).Unique(),`)
	o2.L(`}`)
	o2.L(`}`)

	fn := fmt.Sprintf(`ent/schema/%s_gen.go`, packageName(field.Name(false)))
	if err := o2.WriteFile(fn, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return fmt.Errorf(`failed to write to %s: %w`, fn, err)
	}
	return nil
}*/
