package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	proj_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	proj_pb "github.com/zitadel/zitadel/pkg/grpc/project"
)

func ProjectCreateToDomain(req *mgmt_pb.AddProjectRequest) *domain.Project {
	return &domain.Project{
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.ProjectRoleCheck,
		HasProjectCheck:        req.HasProjectCheck,
		PrivateLabelingSetting: privateLabelingSettingToDomain(req.PrivateLabelingSetting),
	}
}

func ProjectUpdateToDomain(req *mgmt_pb.UpdateProjectRequest) *domain.Project {
	return &domain.Project{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Id,
		},
		Name:                   req.Name,
		ProjectRoleAssertion:   req.ProjectRoleAssertion,
		ProjectRoleCheck:       req.ProjectRoleCheck,
		HasProjectCheck:        req.HasProjectCheck,
		PrivateLabelingSetting: privateLabelingSettingToDomain(req.PrivateLabelingSetting),
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

func AddProjectRoleRequestToDomain(req *mgmt_pb.AddProjectRoleRequest) *domain.ProjectRole {
	return &domain.ProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		Key:         req.RoleKey,
		DisplayName: req.DisplayName,
		Group:       req.Group,
	}
}

func BulkAddProjectRolesRequestToDomain(req *mgmt_pb.BulkAddProjectRolesRequest) []*domain.ProjectRole {
	roles := make([]*domain.ProjectRole, len(req.Roles))
	for i, role := range req.Roles {
		roles[i] = &domain.ProjectRole{
			ObjectRoot: models.ObjectRoot{
				AggregateID: req.ProjectId,
			},
			Key:         role.Key,
			DisplayName: role.DisplayName,
			Group:       role.Group,
		}
	}
	return roles
}

func UpdateProjectRoleRequestToDomain(req *mgmt_pb.UpdateProjectRoleRequest) *domain.ProjectRole {
	return &domain.ProjectRole{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
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

func AddProjectMemberRequestToDomain(req *mgmt_pb.AddProjectMemberRequest) *domain.Member {
	return domain.NewMember(req.ProjectId, req.UserId, req.Roles...)
}

func UpdateProjectMemberRequestToDomain(req *mgmt_pb.UpdateProjectMemberRequest) *domain.Member {
	return domain.NewMember(req.ProjectId, req.UserId, req.Roles...)
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
