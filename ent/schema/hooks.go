package schema

// NOTE: This file has a dependency on generated files in ../ent.
// This means that if we generated code that does not compile
// in ../ent, we now run the risk of not being able to compile
// this hook either.
//
// It also won't work if the hook depends on some code that has
// not yet been generated in the ../ent directory
//
// In those cases, remove this from the files to be compiled

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent"
	gen "github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/hook"
	"github.com/lestrrat-go/dataurl"
)

func (Photo) Hooks() []ent.Hook {
	return []ent.Hook{
		UploadBlob(),
	}
}

func UploadBlob() ent.Hook {
	h := func(next ent.Mutator) ent.Mutator {
		return hook.PhotoFunc(func(ctx context.Context, m *gen.PhotoMutation) (ent.Value, error) {
			value, ok := m.Value()
			if !ok {
				return next.Mutate(ctx, m)
			}

			id, ok := m.ID()
			if !ok {
				return next.Mutate(ctx, m)
			}

			userID, ok := m.UserID()
			if !ok {
				return next.Mutate(ctx, m)
			}

			if !strings.HasPrefix(value, `data:`) {
				return next.Mutate(ctx, m)
			}

			parsed, err := dataurl.Parse([]byte(value))
			if err != nil {
				return nil, fmt.Errorf(`failed to parse photo.value: %w`, err)
			}

			// We only support image/jpeg, image/png
			var suffix string
			switch parsed.MediaType.Type {
			case `image/jpeg`:
				suffix = `.jpg`
			case `image/png`:
				suffix = `.png`
			default:
				return nil, fmt.Errorf(`media type %q not supported for photo`, parsed.MediaType.Type)
			}

			if err := m.Bucket.WriteAll(ctx, id.String()+suffix, []byte(parsed.Data), nil); err != nil {
				return nil, fmt.Errorf(`failed to store blob`)
			}

			u, err := m.PhotoURL.Make(userID.String(), id.String()+suffix)
			if err != nil {
				// TODO: if this happens, we need to delete the old one?
				return nil, fmt.Errorf(`failed to create URL for new object: %w`, err)
			}
			m.SetValue(u)
			return next.Mutate(ctx, m)
		})
	}
	return hook.On(h, ent.OpCreate|ent.OpUpdate)
}
