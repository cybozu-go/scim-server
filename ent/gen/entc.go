//go:build ignore

package main

import (
	"fmt"
	"os"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/cybozu-go/scim-server/helper"
	"gocloud.dev/blob"
)

func main() {
	var pr helper.PhotoURLFunc = func(string, string) (string, error) { return "dummy", nil }
	// &helper.NilPhotoURL{}
	err := entc.Generate("./schema",
		&gen.Config{},
		entc.Extensions(&ETag{}),
		entc.Dependency(
			// object that is responsible for
			entc.DependencyName(`PhotoURL`),
			entc.DependencyType(pr),
		),
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
