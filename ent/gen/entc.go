//go:build ignore

package main

import (
	"fmt"
	"os"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"gocloud.dev/blob"
)

func main() {
	err := entc.Generate("./schema",
		&gen.Config{},
		entc.Extensions(&ETag{}),
		entc.Dependency(
			entc.DependencyName("Bucket"),
			entc.DependencyType(&blob.Bucket{}),
		),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate ent code: %s\n", err)
		os.Exit(1)
	}
}

type ETag struct {
	entc.DefaultExtension
}

func (*ETag) Templates() []*gen.Template {
	return []*gen.Template{
		gen.MustParse(
			gen.NewTemplate(`etag`).ParseDir(`./tmpl/`),
		),
	}
}
