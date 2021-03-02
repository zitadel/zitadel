package management

import (
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	proj_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ProjectCreateToDomain(req *mgmt_pb.AddProjectRequest) *domain.Project {
	return &domain.Project{
		Name:                 req.Name,
		ProjectRoleAssertion: req.ProjectRoleAssertion,
		ProjectRoleCheck:     req.ProjectRoleCheck,
	}
}

func ProjectUpdateToDomain(req *mgmt_pb.UpdateProjectRequest) *domain.Project {
	return &domain.Project{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Id,
		},
		Name:                 req.Name,
		ProjectRoleAssertion: req.ProjectRoleAssertion,
		ProjectRoleCheck:     req.ProjectRoleCheck,
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

func ProjectGrantsToIDs(projectGrants []*proj_model.ProjectGrantView) []string {
	converted := make([]string, len(projectGrants))
	for i, grant := range projectGrants {
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

func ListProjectsRequestToModel(req *mgmt_pb.ListProjectsRequest) (*proj_model.ProjectViewSearchRequest, error) {
	queries, err := proj_grpc.ProjectQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectViewSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func ListGrantedProjectsRequestToModel(req *mgmt_pb.ListGrantedProjectsRequest) (*proj_model.ProjectGrantViewSearchRequest, error) {
	queries, err := proj_grpc.GrantedProjectQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectGrantViewSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}
func ListProjectRolesRequestToModel(req *mgmt_pb.ListProjectRolesRequest) (*proj_model.ProjectRoleSearchRequest, error) {
	queries, err := proj_grpc.RoleQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &proj_model.ProjectRoleSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func ListProjectMembersRequestToModel(req *mgmt_pb.ListProjectMembersRequest) (*proj_model.ProjectMemberSearchRequest, error) {
	queries := member_grpc.MemberQueriesToProjectMember(req.Queries)
	queries = append(queries, &proj_model.ProjectMemberSearchQuery{
		Key:    proj_model.ProjectMemberSearchKeyProjectID,
		Method: domain.SearchMethodEquals,
		Value:  req.ProjectId,
	})
	return &proj_model.ProjectMemberSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}
