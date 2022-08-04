package schema

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent"
	gen "github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/hook"
	"github.com/lestrrat-go/dataurl"
)

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

			m.SetValue(fmt.Sprintf(`https://random/%s/%s`, userID.String(), id.String()+suffix))
			return next.Mutate(ctx, m)
		})
	}
	return hook.On(h, ent.OpCreate|ent.OpUpdate)
}
