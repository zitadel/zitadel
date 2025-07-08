package management

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	proj_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	proj_pb "github.com/zitadel/zitadel/pkg/grpc/project"
)

func ProjectCreateToCommand(req *mgmt_pb.AddProjectRequest, projectID string, resourceOwner string) *command.AddProject {
	return &command.AddProject{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.ProjectRoleCheck,
		HasProjectCheck:        req.HasProjectCheck,
		PrivateLabelingSetting: privateLabelingSettingToDomain(req.PrivateLabelingSetting),
	}
}

func ProjectUpdateToCommand(req *mgmt_pb.UpdateProjectRequest, resourceOwner string) *command.ChangeProject {
	return &command.ChangeProject{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.Id,
			ResourceOwner: resourceOwner,
		},
		Name:                   gu.Ptr(req.Name),
		ProjectRoleAssertion:   gu.Ptr(req.ProjectRoleAssertion),
		ProjectRoleCheck:       gu.Ptr(req.ProjectRoleCheck),
		HasProjectCheck:        gu.Ptr(req.HasProjectCheck),
		PrivateLabelingSetting: gu.Ptr(privateLabelingSettingToDomain(req.PrivateLabelingSetting)),
	}
}

func privateLabelingSettingToDomain(setting proj_pb.PrivateLabelingSetting) domain.PrivateLabelingSetting {
	switch setting {
	case proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY:
		return domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy
	case proj_pb.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY:
		return domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy
	default:
		return domain.PrivateLabelingSettingUnspecified
	}
}

func AddProjectRoleRequestToCommand(req *mgmt_pb.AddProjectRoleRequest, resourceOwner string) *command.AddProjectRole {
	return &command.AddProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.ProjectId,
			ResourceOwner: resourceOwner,
		},
		Key:         req.RoleKey,
		DisplayName: req.DisplayName,
		Group:       req.Group,
	}
}

func BulkAddProjectRolesRequestToCommand(req *mgmt_pb.BulkAddProjectRolesRequest, resourceOwner string) []*command.AddProjectRole {
	roles := make([]*command.AddProjectRole, len(req.Roles))
	for i, role := range req.Roles {
		roles[i] = &command.AddProjectRole{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   req.ProjectId,
				ResourceOwner: resourceOwner,
			},
			Key:         role.Key,
			DisplayName: role.DisplayName,
			Group:       role.Group,
		}
	}
	return roles
}

func UpdateProjectRoleRequestToCommand(req *mgmt_pb.UpdateProjectRoleRequest, resourceOwner string) *command.ChangeProjectRole {
	return &command.ChangeProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.ProjectId,
			ResourceOwner: resourceOwner,
		},
		Key:         req.RoleKey,
		DisplayName: req.DisplayName,
		Group:       req.Group,
	}
}

func ProjectGrantsToIDs(projectGrants *query.ProjectGrants) []string {
	converted := make([]string, len(projectGrants.ProjectGrants))
	for i, grant := range projectGrants.ProjectGrants {
		converted[i] = grant.GrantID
	}
	return converted
}

func AddProjectMemberRequestToCommand(req *mgmt_pb.AddProjectMemberRequest, orgID string) *command.AddProjectMember {
	return &command.AddProjectMember{
		ResourceOwner: orgID,
		ProjectID:     req.ProjectId,
		UserID:        req.UserId,
		Roles:         req.Roles,
	}
}

func UpdateProjectMemberRequestToCommand(req *mgmt_pb.UpdateProjectMemberRequest, orgID string) *command.ChangeProjectMember {
	return &command.ChangeProjectMember{
		ResourceOwner: orgID,
		ProjectID:     req.ProjectId,
		UserID:        req.UserId,
		Roles:         req.Roles,
	}
}

func listProjectRequestToModel(req *mgmt_pb.ListProjectsRequest) (*query.ProjectSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := proj_grpc.ProjectQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.ProjectSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func listGrantedProjectsRequestToModel(req *mgmt_pb.ListGrantedProjectsRequest) (*query.ProjectGrantSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := proj_grpc.ProjectQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.ProjectGrantSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func listProjectRolesRequestToModel(req *mgmt_pb.ListProjectRolesRequest) (*query.ProjectRoleSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := proj_grpc.RoleQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.ProjectRoleSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func listGrantedProjectRolesRequestToModel(req *mgmt_pb.ListGrantedProjectRolesRequest) (*query.ProjectRoleSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := proj_grpc.RoleQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.ProjectRoleSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func ListProjectMembersRequestToModel(ctx context.Context, req *mgmt_pb.ListProjectMembersRequest) (*query.ProjectMembersQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := query.NewMemberResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, ownerQuery)
	return &query.ProjectMembersQuery{
		MembersQuery: query.MembersQuery{
			SearchRequest: query.SearchRequest{
				Offset: offset,
				Limit:  limit,
				Asc:    asc,
				//SortingColumn: //TODO: sorting
			},
			Queries: queries,
		},
		ProjectID: req.ProjectId,
	}, nil
}
