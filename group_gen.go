package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"

	"entgo.io/ent/dialect/sql"
	"github.com/cybozu-go/scim-server/ent"
	"github.com/cybozu-go/scim-server/ent/group"
	"github.com/cybozu-go/scim-server/ent/groupmember"
	"github.com/cybozu-go/scim-server/ent/predicate"
	"github.com/cybozu-go/scim/filter"
	"github.com/cybozu-go/scim/resource"
	"github.com/google/uuid"
)

func groupLoadEntFields(q *ent.GroupQuery, scimFields, excludedFields []string) {
	fields := make(map[string]struct{})
	if len(scimFields) == 0 {
		scimFields = []string{resource.GroupDisplayNameKey, resource.GroupExternalIDKey, resource.GroupIDKey, resource.GroupMembersKey}
	}

	for _, name := range scimFields {
		fields[name] = struct{}{}
	}

	for _, name := range excludedFields {
		delete(fields, name)
	}
	selectNames := make([]string, 0, len(fields))
	for f := range fields {
		switch f {
		case resource.GroupDisplayNameKey:
			selectNames = append(selectNames, group.FieldDisplayName)
		case resource.GroupExternalIDKey:
			selectNames = append(selectNames, group.FieldExternalID)
		case resource.GroupIDKey:
			selectNames = append(selectNames, group.FieldID)
		case resource.GroupMembersKey:
		case resource.GroupMetaKey:
		}
	}
	selectNames = append(selectNames, group.FieldEtag)
	q.Select(selectNames...)
}

func groupLocation(id string) string {
	return "https://foobar.com/scim/v2/Groups/" + id
}

func GroupResourceFromEnt(in *ent.Group) (*resource.Group, error) {
	var b resource.Builder

	builder := b.Group()

	meta, err := b.Meta().
		ResourceType("Group").
		Location(groupLocation(in.ID.String())).
		Version(in.Etag).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build meta information for Group")
	}

	builder.
		Meta(meta)

	if el := len(in.Edges.Members); el > 0 {
		list := make([]*resource.GroupMember, 0, el)
		for _, ine := range in.Edges.Members {
			r, err := GroupMemberResourceFromEnt(ine)
			if err != nil {
				return nil, fmt.Errorf("failed to build members information for Group")
			}
			list = append(list, r)
		}
		builder.Members(list...)
	}
	if !reflect.ValueOf(in.DisplayName).IsZero() {
		builder.DisplayName(in.DisplayName)
	}
	if !reflect.ValueOf(in.ExternalID).IsZero() {
		builder.ExternalID(in.ExternalID)
	}
	builder.ID(in.ID.String())
	if err := groupResourceFromEntHelper(in, builder); err != nil {
		return nil, err
	}
	return builder.Build()
}

func GroupEntFieldFromSCIM(s string) string {
	switch s {
	case resource.GroupDisplayNameKey:
		return group.FieldDisplayName
	case resource.GroupExternalIDKey:
		return group.FieldExternalID
	case resource.GroupIDKey:
		return group.FieldID
	default:
		return s
	}
}

func groupStartsWithPredicate(q *ent.GroupQuery, scimField string, val interface{}) (predicate.Group, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.GroupDisplayNameKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupExternalIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasPrefix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func groupEndsWithPredicate(q *ent.GroupQuery, scimField string, val interface{}) (predicate.Group, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.GroupDisplayNameKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupExternalIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.HasSuffix(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func groupContainsPredicate(q *ent.GroupQuery, scimField string, val interface{}) (predicate.Group, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.GroupDisplayNameKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupExternalIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.Contains(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func groupEqualsPredicate(q *ent.GroupQuery, scimField string, val interface{}) (predicate.Group, error) {
	_ = q
	field, subfield, err := splitScimField(scimField)
	if err != nil {
		return nil, err
	}
	_ = subfield // TODO: remove later
	switch field {
	case resource.GroupDisplayNameKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupExternalIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	case resource.GroupIDKey:
		entFieldName := GroupEntFieldFromSCIM(scimField)
		return predicate.Group(func(s *sql.Selector) {
			//nolint:forcetypeassert
			s.Where(sql.EQ(s.C(entFieldName), val.(string)))
		}), nil
	default:
		return nil, fmt.Errorf("invalid filter field specification")
	}
}

func groupPresencePredicate(scimField string) predicate.Group {
	switch scimField {
	case resource.GroupDisplayNameKey:
		return group.And(group.DisplayNameNotNil(), group.DisplayNameNEQ(""))
	case resource.GroupExternalIDKey:
		return group.And(group.ExternalIDNotNil(), group.ExternalIDNEQ(""))
	default:
		return nil
	}
}

func (b *Backend) existsGroupGroupMember(parent *ent.Group, in *resource.GroupMember) bool {
	ctx := context.TODO()
	queryCall := parent.QueryMembers()
	var predicates []predicate.GroupMember
	if in.HasRef() {
		predicates = append(predicates, groupmember.Ref(in.Ref()))
	}
	if in.HasType() {
		predicates = append(predicates, groupmember.Type(in.Type()))
	}
	if in.HasValue() {
		predicates = append(predicates, groupmember.Value(in.Value()))
	}

	v, err := queryCall.Where(predicates...).Exist(ctx)
	if err != nil {
		return false
	}
	return v
}

func (b *Backend) CreateGroup(in *resource.Group) (*resource.Group, error) {
	ctx := context.TODO()

	createCall := b.db.Group.Create()
	if in.HasDisplayName() {
		createCall.SetDisplayName(in.DisplayName())
	}
	if in.HasExternalID() {
		createCall.SetExternalID(in.ExternalID())
	}
	var members []*ent.GroupMember
	if in.HasMembers() {
		created, err := b.createGroupMember(in.Members()...)
		if err != nil {
			return nil, fmt.Errorf("failed to create members: %w", err)
		}
		createCall.AddMembers(created...)
		members = created
	}

	rs, err := createCall.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save object: %w", err)
	}
	rs.Edges.Members = members

	h := sha256.New()
	if err := rs.ComputeETag(h); err != nil {
		return nil, fmt.Errorf("failed to compute etag: %w", err)
	}
	etag := fmt.Sprintf("W/%x", h.Sum(nil))

	if _, err := rs.Update().SetEtag(etag).Save(ctx); err != nil {
		return nil, fmt.Errorf("failed to save etag: %w", err)
	}
	rs.Etag = etag
	return GroupResourceFromEnt(rs)
}

func (b *Backend) ReplaceGroup(id string, in *resource.Group) (*resource.Group, error) {
	ctx := context.TODO()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	r, err := b.db.Group.Query().Where(group.ID(parsedUUID)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve resource for replacing: %w", err)
	}

	replaceCall := r.Update()

	if in.HasDisplayName() {
		replaceCall.ClearDisplayName()
		replaceCall.SetDisplayName(in.DisplayName())
	}

	if in.HasExternalID() {
		replaceCall.ClearExternalID()
		replaceCall.SetExternalID(in.ExternalID())
	}

	if in.HasMembers() {
		replaceCall.ClearMembers()
		created, err := b.createGroupMember(in.Members()...)
		if err != nil {
			return nil, err
		}
		replaceCall.AddMembers(created...)
	}
	if _, err := replaceCall.Save(ctx); err != nil {
		return nil, fmt.Errorf("failed to save object: %w", err)
	}

	r2, err := b.db.Group.Query().Where(group.ID(parsedUUID)).
		WithMembers().
		Only(ctx)

	h := sha256.New()
	if err := r2.ComputeETag(h); err != nil {
		return nil, fmt.Errorf("failed to compute etag: %w", err)
	}
	etag := fmt.Sprintf("W/%x", h.Sum(nil))

	if _, err := r2.Update().SetEtag(etag).Save(ctx); err != nil {
		return nil, fmt.Errorf("failed to save etag: %w", err)
	}
	r2.Etag = etag

	return GroupResourceFromEnt(r2)
}

func (b *Backend) patchAddGroup(parent *ent.Group, op *resource.PatchOperation) error {
	ctx := context.TODO()
	root, err := filter.Parse(op.Path())
	if err != nil {
		return fmt.Errorf("failed to parse PATH path %q", op.Path())
	}

	expr, ok := root.(filter.ValuePath)
	if !ok {
		return fmt.Errorf("root element should be a valuePath (got %T)", root)
	}

	sattr, err := exprStr(expr.ParentAttr())
	if err != nil {
		return fmt.Errorf("invalid attribute specification: %w", err)
	}

	switch sattr {
	case resource.GroupDisplayNameKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element displayName")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element displayName")
		}

		if _, err := parent.Update().SetDisplayName(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.GroupExternalIDKey:
		subExpr := expr.SubExpr()
		if subExpr != nil {
			return fmt.Errorf("subexpr on string element is unimplmented")
		}

		if expr.SubAttr() != nil {
			return fmt.Errorf("invalid sub attrribute on string element externalId")
		}

		var v string
		if err := json.Unmarshal(op.Value(), &v); err != nil {
			return fmt.Errorf("invalid value for string element externalId")
		}

		if _, err := parent.Update().SetExternalID(v).Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.GroupMembersKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch add operation on sub attribute of multi-value item members with unspecified element is not possible")
			}

			var in resource.GroupMember
			if err := json.Unmarshal(op.Value(), &in); err != nil {
				return fmt.Errorf("failed to decode patch add value: %w", err)
			}

			if b.existsGroupGroupMember(parent, &in) {
				return nil
			}

			created, err := b.createGroupMember(&in)
			if err != nil {
				return fmt.Errorf("failed to create GroupMember: %w", err)
			}

			if _, err := parent.Update().AddMembers(created...).Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
		} else {
			var pb GroupMemberPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}
			list, err := parent.QueryMembers().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to look up value: %w", err)
			}

			if len(list) > 0 {
				return fmt.Errorf("query must resolve to one element, got multiple")
			}

			item := list[0]
			sSubAttr, err := exprStr(expr.SubAttr())
			if err != nil {
				return fmt.Errorf("query must have a sub attribute")
			}

			updateCall := item.Update()

			switch sSubAttr {
			case resource.GroupMemberRefKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetRef(v)
			case resource.GroupMemberTypeKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetType(v)
			case resource.GroupMemberValueKey:
				var v string
				if err := json.Unmarshal(op.Value(), &v); err != nil {
					return fmt.Errorf("failed to decode value: %w", err)
				}
				updateCall.SetValue(v)
			}

			if _, err := updateCall.Save(ctx); err != nil {
				return fmt.Errorf("failed to save object: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (b *Backend) patchRemoveGroup(parent *ent.Group, op *resource.PatchOperation) error {
	ctx := context.TODO()

	root, err := filter.Parse(op.Path())
	if err != nil {
		return fmt.Errorf("failed to parse path %q", op.Path())
	}

	expr, ok := root.(filter.ValuePath)
	if !ok {
		return fmt.Errorf("root element should be a valuePath (got %T)", root)
	}

	sattr, err := exprStr(expr.ParentAttr())
	if err != nil {
		return fmt.Errorf("invalid attribute specification: %w", err)
	}
	switch sattr {
	case resource.GroupDisplayNameKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on displayName cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on displayName cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearDisplayName().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.GroupExternalIDKey:
		if subexpr := expr.SubExpr(); subexpr != nil {
			return fmt.Errorf("patch remove operation on externalId cannot have a sub attribute query")
		}

		if subattr := expr.SubAttr(); subattr != nil {
			return fmt.Errorf("patch remove operation on externalId cannot have a sub attribute")
		}

		if _, err := parent.Update().ClearExternalID().Save(ctx); err != nil {
			return fmt.Errorf("failed to save object: %w", err)
		}
	case resource.GroupMembersKey:
		subExpr := expr.SubExpr()
		if subExpr == nil {
			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				return fmt.Errorf("patch remove operation on su attribute of multi-valued item members without a query is not possible")
			}
			if _, err := b.db.GroupMember.Delete().Where(groupmember.HasGroupWith(group.ID(parent.ID))).Exec(ctx); err != nil {
				return fmt.Errorf("failed to remove elements from members: %w", err)
			}
			if _, err := parent.Update().ClearMembers().Save(ctx); err != nil {
				return fmt.Errorf("failed to remove references to members: %w", err)
			}
		} else {
			var pb GroupMemberPredicateBuilder
			predicates, err := pb.Build(subExpr)
			if err != nil {
				return fmt.Errorf("failed to parse valuePath expression: %w", err)
			}

			list, err := parent.QueryMembers().
				Where(predicates...).
				All(ctx)
			if err != nil {
				return fmt.Errorf("failed to query context object: %w", err)
			}

			if subAttrExpr := expr.SubAttr(); subAttrExpr != nil {
				subAttr, err := exprStr(subAttrExpr)
				if err != nil {
					return fmt.Errorf("invalid sub attribute specified")
				}
				switch subAttr {
				case resource.GroupMemberRefKey:
					return fmt.Errorf("$ref is not mutable")
				case resource.GroupMemberTypeKey:
					return fmt.Errorf("type is not mutable")
				case resource.GroupMemberValueKey:
					return fmt.Errorf("value is not mutable")
				default:
					return fmt.Errorf("unknown sub attribute specified")
				}
			}

			ids := make([]int, len(list))
			for i, elem := range list {
				ids[i] = elem.ID
			}
			if _, err := b.db.GroupMember.Delete().Where(groupmember.IDIn(ids...)).Exec(ctx); err != nil {
				return fmt.Errorf("failed to delete object: %w", err)
			}
		}
	}
	return nil
}
