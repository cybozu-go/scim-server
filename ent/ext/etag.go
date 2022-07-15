package ext

import (
	_ "embed"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

//go:embed etag.tmpl
var etagTemplate string

type ETag struct {
	entc.DefaultExtension
}

func (*ETag) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(
			gen.NewTemplate(`etag`).Parse(etagTemplate),
		),
	}
}
