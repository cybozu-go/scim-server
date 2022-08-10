package server_test

import (
	"context"
	"testing"

	server "github.com/cybozu-go/scim-server"
	"github.com/cybozu-go/scim-server/ent"
	_ "github.com/cybozu-go/scim-server/ent/runtime"
	"github.com/cybozu-go/scim-server/helper"
	"github.com/cybozu-go/scim/test"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/memblob"
)

func TestSample(t *testing.T) {
	ctx := context.TODO()

	bucket, err := blob.OpenBucket(ctx, "mem://photos/")
	require.NoError(t, err, `blob.OpenBucket should succeed`)

	s, err := server.New("file:ent?mode=memory&cache=shared&_fk=1",
		ent.Bucket(bucket),
		ent.PhotoURL(helper.PhotoURLFunc(func(uid, path string) (string, error) {
			return "https://sample/foo.png", nil
		})),
	)
	require.NoError(t, err, `server.New should succeed`)

	test.RunConformanceTests(t, "Sample backend", s)
}
