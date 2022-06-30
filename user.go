package server

import (
	"fmt"

	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim/resource"
)

func userResourceFromEntHelper(in *ent.User, builder *resource.UserBuilder) error {
	groups := make([]*resource.GroupMember, len(in.Edges.Groups))
	for _, g := range in.Edges.Groups {
		id := g.ID.String()
		gm, err := resource.NewGroupMemberBuilder().
			Value(id).
			Ref(groupLocation(id)).
			Build()
		if err != nil {
			return fmt.Errorf(`failed to convert internal group data to SCIM resource: %w`, err)
		}
		groups = append(groups, gm)
	}
	builder.Groups(groups...)
	return nil
}
