package server

import (
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim/resource"
)

func MemberResourceFromEnt(in *ent.Member) (*resource.GroupMember, error) {
	return resource.NewGroupMemberBuilder().
		Value(in.Value).
		Ref(userLocation(in.Value)).
		Type(in.Type).
		Build()
}

func groupResourceFromEntHelper(in *ent.Group, builder *resource.GroupBuilder) error {
	/*
		members := make([]*resource.GroupMember, 0, len(in.Edges.Users)+len(in.Edges.Children))
		for _, u := range in.Edges.Users {
			id := u.ID.String()
			gm, err := resource.NewGroupMemberBuilder().
				Value(id).
				Ref(userLocation(id)).
				Build()
			if err != nil {
				return fmt.Errorf(`failed to convert internal member data to SCIM resource: %w`, err)
			}
			members = append(members, gm)
		}
		for _, subg := range in.Edges.Children {
			id := subg.ID.String()
			gm, err := resource.NewGroupMemberBuilder().
				Value(id).
				Ref(groupLocation(id)).
				Build()
			if err != nil {
				return fmt.Errorf(`failed to convert interna member data to SCIM resource: %w`, err)
			}
			members = append(members, gm)
		}

		builder.Members(members...)
		return nil*/
	return nil
}
