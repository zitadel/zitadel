package management

import (
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/api/grpc/object"
	proj_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ListProjectGrantsRequestToModel(req *mgmt_pb.ListProjectGrantsRequest) (*proj_model.ProjectGrantViewSearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries := proj_grpc.ProjectGrantQueriesToModel(req.Queries)
	queries = append(queries, &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyProjectID,
		Method: domain.SearchMethodEquals,
		Value:  req.ProjectId,
	})
	return &proj_model.ProjectGrantViewSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func AddProjectGrantRequestToDomain(req *mgmt_pb.AddProjectGrantRequest) *domain.ProjectGrant {
	return &domain.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantedOrgID: req.GrantedOrgId,
		RoleKeys:     req.RoleKeys,
	}
}

func UpdateProjectGrantRequestToDomain(req *mgmt_pb.UpdateProjectGrantRequest) *domain.ProjectGrant {
	return &domain.ProjectGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantID:  req.GrantId,
		RoleKeys: req.RoleKeys,
	}
}

func ListProjectGrantMembersRequestToModel(req *mgmt_pb.ListProjectGrantMembersRequest) *proj_model.ProjectGrantMemberSearchRequest {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries := member_grpc.MemberQueriesToProjectGrantMember(req.Queries)
	queries = append(queries,
		&proj_model.ProjectGrantMemberSearchQuery{
			Key:    proj_model.ProjectGrantMemberSearchKeyProjectID,
			Method: domain.SearchMethodEquals,
			Value:  req.ProjectId,
		},
		&proj_model.ProjectGrantMemberSearchQuery{
			Key:    proj_model.ProjectGrantMemberSearchKeyGrantID,
			Method: domain.SearchMethodEquals,
			Value:  req.GrantId,
		})
	return &proj_model.ProjectGrantMemberSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}
}

func AddProjectGrantMemberRequestToDomain(req *mgmt_pb.AddProjectGrantMemberRequest) *domain.ProjectGrantMember {
	return &domain.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantID: req.GrantId,
		UserID:  req.UserId,
		Roles:   req.Roles,
	}
}

func UpdateProjectGrantMemberRequestToDomain(req *mgmt_pb.UpdateProjectGrantMemberRequest) *domain.ProjectGrantMember {
	return &domain.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		GrantID: req.GrantId,
		UserID:  req.UserId,
		Roles:   req.Roles,
	}
}
