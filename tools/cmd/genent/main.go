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

func main() {
	objects = make(map[string]*codegen.Object)
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
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
	case `User`, `Group`, `Email`, `Names`, `Role`, `Photo`, `IMS`, `PhoneNumber`:
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

		switch field.Name(true) {
		case `Entitlements`, `Roles`:
			if err := generateSimpleEdge(field); err != nil {
				return fmt.Errorf(`failed to generate edge %q: %w`, field.Name(true), err)
			}
		default:
		}
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

	fn := filepath.Join(`ent`, `schema`, xstrings.Snake(object.Name(false))+`_gen.go`)
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
			// Special case
			var ft = field.Type()
			if strings.HasPrefix(ft, `[]`) || strings.HasPrefix(ft, `*`) {
				// TODO: later
				switch field.Name(false) {
				case `emails`, `name`:
					o.L(`q.With%s()`, field.Name(true))
				}
				continue
			} else {
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
		if field.Name(false) == "schemas" {
			continue
		}

		switch field.Name(false) {
		case `emails`, `name`:
		default:
			continue
		}

		// This section is just really really confusing because not all
		// stored data map 1-to-1 to the SCIM resource (for example,
		// Group.members can't be expressed 1-to-1 in a straight forward
		// manner).
		// I think it's better if we give people an escape hatch, so we're
		// going to inject a call to a helper of your choice at the end.
		rsname := strings.TrimSuffix(field.Name(true), "s")
		if rsname == "Name" {
			rsname = "Names"
		}
		o.LL(`if el := len(in.Edges.%s); el > 0 {`, field.Name(true))
		o.L(`list := make([]*resource.%s, 0, el)`, rsname)
		o.L(`for _, ine := range in.Edges.%s {`, field.Name(true))
		o.L(`r, err := %sResourceFromEnt(ine)`, rsname)
		o.L(`if err != nil {`)
		o.L(`return nil, fmt.Errorf("failed to build %s information for %s")`, field.Name(false), object.Name(true))
		o.L(`}`)
		o.L(`list = append(list, r)`)
		o.L(`}`)

		if strings.HasPrefix(field.Type(), "*") {
			o.L(`builder.%s(list[0])`, field.Name(true))
		} else {
			o.L(`builder.%s(list...)`, field.Name(true))
		}
		o.L(`}`)
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

	fn := fmt.Sprintf(`%s_gen.go`, relationFilename(object.Name(false)))
	if err := o.WriteFile(fn, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return fmt.Errorf(`failed to write to %s: %w`, fn, err)
	}
	return nil
}

func generateSimpleEdge(field codegen.Field) error {

	// If this is a simple one-to-many, we can just generate the schema
	var o2buf bytes.Buffer
	o2 := codegen.NewOutput(&o2buf)
	structName := singularName(field.Name(true))
	fmt.Printf("  ⌛ Generating smple edge adapters for %s...\n", structName)
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

	fn := fmt.Sprintf(`ent/schema/%s_gen.go`, relationFilename(field.Name(false)))
	if err := o2.WriteFile(fn, codegen.WithFormatCode(true)); err != nil {
		if cfe, ok := err.(codegen.CodeFormatError); ok {
			fmt.Fprint(os.Stderr, cfe.Source())
		}
		return fmt.Errorf(`failed to write to %s: %w`, fn, err)
	}
	return nil
}
