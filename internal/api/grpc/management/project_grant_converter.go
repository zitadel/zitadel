package management

import (
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	proj_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ListProjectGrantsRequestToModel(req *mgmt_pb.ListProjectGrantsRequest) (*proj_model.ProjectGrantViewSearchRequest, error) {
	queries := proj_grpc.ProjectGrantQueriesToModel(req.Queries)
	queries = append(queries, &proj_model.ProjectGrantViewSearchQuery{
		Key:    proj_model.GrantedProjectSearchKeyProjectID,
		Method: domain.SearchMethodEquals,
		Value:  req.ProjectId,
	})
	return &proj_model.ProjectGrantViewSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
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
	queries := member_grpc.MemberQueriesToProjectGrantMember(req.Queries)
	queries = append(queries, &proj_model.ProjectGrantMemberSearchQuery{
		Key:    proj_model.ProjectGrantMemberSearchKeyProjectID,
		Method: domain.SearchMethodEquals,
		Value:  req.ProjectId,
	})
	return &proj_model.ProjectGrantMemberSearchRequest{
		Offset: req.MetaData.Offset,
		Limit:  uint64(req.MetaData.Limit),
		Asc:    req.MetaData.Asc,
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
