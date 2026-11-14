package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	scim_schemas "github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GroupsHandler struct {
	command *command.Commands
	query   *query.Queries
	schema  *scim_schemas.ResourceSchema
}

type ScimGroup struct {
	*scim_schemas.Resource `scim:"ignoreInSchema"`
	ID                     string             `json:"id" scim:"ignoreInSchema"`
	ExternalID             string             `json:"externalId,omitempty"`
	DisplayName            string             `json:"displayName,omitempty" scim:"required"`
	Members                []*ScimGroupMember `json:"members,omitempty"`
}

type ScimGroupMember struct {
	Value   string `json:"value"`
	Display string `json:"display,omitempty"`
	Ref     string `json:"$ref,omitempty" scim:"ignoreInSchema"`
	Type    string `json:"type,omitempty"`
}

func NewGroupsHandler(
	command *command.Commands,
	query *query.Queries,
) ResourceHandler[*ScimGroup] {
	return &GroupsHandler{
		command: command,
		query:   query,
		schema: scim_schemas.BuildSchema(scim_schemas.SchemaBuilderArgs{
			ID:           scim_schemas.IdGroup,
			Name:         scim_schemas.GroupResourceType,
			EndpointName: scim_schemas.GroupsResourceType,
			Description:  "Group",
			Resource:     new(ScimGroup),
		}),
	}
}

func (g *ScimGroup) GetResource() *scim_schemas.Resource {
	return g.Resource
}

func (g *ScimGroup) GetSchemas() []scim_schemas.ScimSchemaType {
	if g.Resource == nil {
		return nil
	}
	return g.Schemas
}

func (h *GroupsHandler) Schema() *scim_schemas.ResourceSchema {
	return h.schema
}

func (h *GroupsHandler) NewResource() *ScimGroup {
	return new(ScimGroup)
}

// validateMemberTypes makes the flat group model explicit: RFC 7643 allows
// members of type "Group" (nested groups), which ZITADEL does not support.
// Rejecting them loudly beats silently misreading a group ID as a user ID.
func validateMemberTypes(members []*ScimGroupMember) error {
	for _, member := range members {
		if member == nil || strings.TrimSpace(member.Value) == "" {
			return invalidGroupPatch("Each member requires a user ID")
		}
		switch member.Type {
		case "", "User":
		case "Group":
			return zerrors.ThrowInvalidArgument(nil, "SCIM-GRP6n", "Nested groups are not supported, members must be of type User")
		default:
			return zerrors.ThrowInvalidArgumentf(nil, "SCIM-GRP7t", "Invalid member type %q, supported values: User", member.Type)
		}
	}
	return nil
}

func (h *GroupsHandler) Create(ctx context.Context, group *ScimGroup) (*ScimGroup, error) {
	if group.DisplayName == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-GRP1m", "Errors.Group.Invalid")
	}
	ids, err := groupMemberIDs(group.Members)
	if err != nil {
		return nil, err
	}

	details, err := h.command.CreateGroup(ctx, &command.CreateGroup{
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		Name:       group.DisplayName,
		UserIDs:    ids,
		ExternalID: &command.GroupExternalID{Scope: string(metadata.ScopeKey(ctx, metadata.KeyExternalId)), Value: group.ExternalID},
	})
	if err != nil {
		return nil, err
	}

	group.ID = details.ID
	group.Resource = buildResource[*ScimGroup](ctx, h, details)
	return group, nil
}

func (h *GroupsHandler) Replace(ctx context.Context, id string, group *ScimGroup) (*ScimGroup, error) {
	if group.DisplayName == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-GRP2r", "Errors.Group.Invalid")
	}

	ids, err := groupMemberIDs(group.Members)
	if err != nil {
		return nil, err
	}
	details, err := h.command.MutateGroup(ctx, id, authz.GetCtxData(ctx).OrgID, []command.GroupMutation{{Name: &group.DisplayName, Members: &ids, ExternalID: &command.GroupExternalID{Scope: string(metadata.ScopeKey(ctx, metadata.KeyExternalId)), Value: group.ExternalID}}})
	if err != nil {
		return nil, err
	}
	group.ID = id
	group.Resource = buildResource[*ScimGroup](ctx, h, details)
	return group, nil
}

func (h *GroupsHandler) Update(ctx context.Context, id string, operations patch.OperationCollection) error {
	mutations, err := groupMutations(operations)
	if err != nil {
		return err
	}
	for i := range mutations {
		if mutations[i].ExternalID != nil {
			mutations[i].ExternalID.Scope = string(metadata.ScopeKey(ctx, metadata.KeyExternalId))
		}
	}
	_, err = h.command.MutateGroup(ctx, id, authz.GetCtxData(ctx).OrgID, mutations)
	return err
}

func (h *GroupsHandler) Delete(ctx context.Context, id string) error {
	if _, err := h.getOrgGroup(ctx, id); err != nil {
		return err
	}
	cascadingGrantIDs, err := h.groupGrantIDs(ctx, id)
	if err != nil {
		return err
	}
	_, err = h.command.DeleteGroup(ctx, id, cascadingGrantIDs...)
	return err
}

func (h *GroupsHandler) Get(ctx context.Context, id string) (*ScimGroup, error) {
	group, err := h.getOrgGroup(ctx, id)
	if err != nil {
		return nil, err
	}
	members, err := h.groupMembers(ctx, id)
	if err != nil {
		return nil, err
	}
	bindings, err := h.query.GroupExternalIDs(ctx, authz.GetCtxData(ctx).OrgID, string(metadata.ScopeKey(ctx, metadata.KeyExternalId)), []string{id})
	if err != nil {
		return nil, err
	}
	result := h.mapToScimGroup(ctx, group, members)
	result.ExternalID = bindings[id]
	return result, nil
}

func (h *GroupsHandler) List(ctx context.Context, request *ListRequest) (*ListResponse[*ScimGroup], error) {
	mapping := filter.FieldPathMapping{
		"displayname": {Column: query.GroupColumnName, FieldType: filter.FieldTypeString, CaseInsensitive: true},
		"externalid":  {FieldType: filter.FieldTypeCustom, BuildMappedQuery: buildGroupExternalIDQuery},
		"id":          {Column: query.GroupColumnID, FieldType: filter.FieldTypeString},
	}
	sr, err := request.toSearchRequest(query.GroupColumnID, mapping)
	if err != nil {
		return nil, err
	}
	orgIDQuery, err := query.NewGroupOrganizationIdSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries := []query.SearchQuery{orgIDQuery}
	if request.Filter != nil {
		fq, err := request.Filter.BuildQuery(ctx, h.schema.ID, mapping)
		if err != nil {
			return nil, err
		}
		queries = append(queries, fq)
	}
	groups, err := h.query.SearchGroups(ctx, &query.GroupSearchQuery{
		SearchRequest: sr,
		Queries:       queries,
	}, nil)
	if err != nil {
		return nil, err
	}
	groupIDs := make([]string, len(groups.Groups))
	for i, group := range groups.Groups {
		groupIDs[i] = group.ID
	}
	bindings, err := h.query.GroupExternalIDs(ctx, authz.GetCtxData(ctx).OrgID, string(metadata.ScopeKey(ctx, metadata.KeyExternalId)), groupIDs)
	if err != nil {
		return nil, err
	}
	scimGroups := make([]*ScimGroup, len(groups.Groups))
	for i, group := range groups.Groups {
		var members []*query.GroupUser
		if !excludeGroupMembers(request.ExcludedAttributes) {
			members, err = h.groupMembers(ctx, group.ID)
			if err != nil {
				return nil, err
			}
		}
		scimGroups[i] = h.mapToScimGroup(ctx, group, members)
		scimGroups[i].ExternalID = bindings[group.ID]
	}
	return NewListResponse(groups.Count, sr, scimGroups), nil
}

func (h *GroupsHandler) getOrgGroup(ctx context.Context, id string) (*query.Group, error) {
	group, err := h.query.GetGroupByID(ctx, id, nil)
	if err != nil {
		return nil, err
	}
	if group.ResourceOwner != authz.GetCtxData(ctx).OrgID {
		return nil, zerrors.ThrowNotFound(nil, "SCIM-GRP5o", "Errors.Group.NotFound")
	}
	return group, nil
}

func (h *GroupsHandler) groupMembers(ctx context.Context, groupID string) ([]*query.GroupUser, error) {
	groupIDQuery, err := query.NewGroupUsersGroupIDsSearchQuery([]string{groupID})
	if err != nil {
		return nil, err
	}
	members, err := h.query.SearchGroupUsers(ctx, &query.GroupUsersSearchQuery{
		Queries: []query.SearchQuery{groupIDQuery},
	}, nil)
	if err != nil {
		return nil, err
	}
	return members.GroupUsers, nil
}

func (h *GroupsHandler) groupGrantIDs(ctx context.Context, groupID string) ([]string, error) {
	groupIDsQuery, err := query.NewGroupGrantGroupIDsSearchQuery([]string{groupID})
	if err != nil {
		return nil, err
	}
	grants, err := h.query.SearchGroupGrants(ctx, &query.GroupGrantsSearchQuery{
		Queries: []query.SearchQuery{groupIDsQuery},
	}, nil)
	if err != nil {
		return nil, err
	}
	grantIDs := make([]string, len(grants.GroupGrants))
	for i, grant := range grants.GroupGrants {
		grantIDs[i] = grant.ID
	}
	return grantIDs, nil
}

func (h *GroupsHandler) mapToScimGroup(ctx context.Context, group *query.Group, members []*query.GroupUser) *ScimGroup {
	scimGroup := &ScimGroup{
		Resource: &scim_schemas.Resource{
			ID:      group.ID,
			Schemas: []scim_schemas.ScimSchemaType{scim_schemas.IdGroup},
			Meta: &scim_schemas.ResourceMeta{
				ResourceType: scim_schemas.GroupResourceType,
				Created:      gu.Ptr(group.CreationDate.UTC()),
				LastModified: gu.Ptr(group.ChangeDate.UTC()),
				Version:      strconv.FormatUint(group.Sequence, 10),
				Location:     scim_schemas.BuildLocationForResource(ctx, h.schema.PluralName, group.ID),
			},
		},
		ID:          group.ID,
		DisplayName: group.Name,
	}
	for _, member := range members {
		scimGroup.Members = append(scimGroup.Members, &ScimGroupMember{
			Value:   member.UserID,
			Display: member.DisplayName,
			Ref:     scim_schemas.BuildLocationForResource(ctx, scim_schemas.UsersResourceType, member.UserID),
			Type:    "User",
		})
	}
	return scimGroup
}

func excludeGroupMembers(attributes string) bool {
	for _, attribute := range strings.Split(attributes, ",") {
		if strings.EqualFold(strings.TrimSpace(attribute), "members") || strings.EqualFold(strings.TrimSpace(attribute), string(scim_schemas.IdGroup)+":members") {
			return true
		}
	}
	return false
}

func buildGroupExternalIDQuery(ctx context.Context, value *filter.CompValue, op *filter.CompareOp) (query.SearchQuery, error) {
	if !op.Equal || value.StringValue == nil {
		return nil, serrors.ThrowInvalidFilter(zerrors.ThrowInvalidArgument(nil, "SCIM-GEXT1", "externalId supports equality with a string"))
	}
	return query.NewGroupExternalIDSearchQuery(string(metadata.ScopeKey(ctx, metadata.KeyExternalId)), *value.StringValue)
}
