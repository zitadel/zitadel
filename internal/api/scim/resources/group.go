package resources

import (
	"context"
	"slices"
	"strconv"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	scim_schemas "github.com/zitadel/zitadel/internal/api/scim/schemas"
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
	return g.Resource.Schemas
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
	if err := validateMemberTypes(group.Members); err != nil {
		return nil, err
	}
	details, err := h.command.CreateGroup(ctx, &command.CreateGroup{
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		Name: group.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	if len(group.Members) > 0 {
		userIDs := make([]string, len(group.Members))
		for i, member := range group.Members {
			userIDs[i] = member.Value
		}
		if _, err = h.command.AddUsersToGroup(ctx, details.ID, userIDs); err != nil {
			return nil, err
		}
	}
	group.ID = details.ID
	group.Resource = buildResource[*ScimGroup](ctx, h, details)
	return group, nil
}

func (h *GroupsHandler) Replace(ctx context.Context, id string, group *ScimGroup) (*ScimGroup, error) {
	if group.DisplayName == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-GRP2r", "Errors.Group.Invalid")
	}
	if err := validateMemberTypes(group.Members); err != nil {
		return nil, err
	}
	existing, err := h.getOrgGroup(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.Name != group.DisplayName {
		if _, err = h.command.UpdateGroup(ctx, &command.UpdateGroup{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   id,
				ResourceOwner: existing.ResourceOwner,
			},
			Name: gu.Ptr(group.DisplayName),
		}); err != nil {
			return nil, err
		}
	}
	if err = h.replaceMembers(ctx, id, group.Members); err != nil {
		return nil, err
	}
	return h.Get(ctx, id)
}

func (h *GroupsHandler) Update(ctx context.Context, id string, operations patch.OperationCollection) error {
	return zerrors.ThrowUnimplemented(nil, "SCIM-GRP3p", "PATCH is not supported for groups, use PUT instead")
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
	return h.mapToScimGroup(ctx, group, members), nil
}

func (h *GroupsHandler) List(ctx context.Context, request *ListRequest) (*ListResponse[*ScimGroup], error) {
	if request.Filter != nil {
		return nil, zerrors.ThrowUnimplemented(nil, "SCIM-GRP4l", "Filtering groups is not supported")
	}
	sr, err := request.toSearchRequest(query.GroupColumnName, filter.FieldPathMapping{})
	if err != nil {
		return nil, err
	}
	orgIDQuery, err := query.NewGroupOrganizationIdSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	groups, err := h.query.SearchGroups(ctx, &query.GroupSearchQuery{
		SearchRequest: sr,
		Queries:       []query.SearchQuery{orgIDQuery},
	}, nil)
	if err != nil {
		return nil, err
	}
	scimGroups := make([]*ScimGroup, len(groups.Groups))
	for i, group := range groups.Groups {
		// members are omitted in listings, use a get on the single group to resolve them
		scimGroups[i] = h.mapToScimGroup(ctx, group, nil)
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

func (h *GroupsHandler) replaceMembers(ctx context.Context, groupID string, members []*ScimGroupMember) error {
	existingMembers, err := h.groupMembers(ctx, groupID)
	if err != nil {
		return err
	}
	existingIDs := make([]string, len(existingMembers))
	for i, member := range existingMembers {
		existingIDs[i] = member.UserID
	}
	requestedIDs := make([]string, len(members))
	for i, member := range members {
		requestedIDs[i] = member.Value
	}

	toAdd := make([]string, 0, len(requestedIDs))
	for _, id := range requestedIDs {
		if !slices.Contains(existingIDs, id) {
			toAdd = append(toAdd, id)
		}
	}
	toRemove := make([]string, 0, len(existingIDs))
	for _, id := range existingIDs {
		if !slices.Contains(requestedIDs, id) {
			toRemove = append(toRemove, id)
		}
	}

	if len(toAdd) > 0 {
		if _, err = h.command.AddUsersToGroup(ctx, groupID, toAdd); err != nil {
			return err
		}
	}
	if len(toRemove) > 0 {
		if _, err = h.command.RemoveUsersFromGroup(ctx, groupID, toRemove); err != nil {
			return err
		}
	}
	return nil
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
