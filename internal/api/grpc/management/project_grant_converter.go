package management

import (
	member_grpc "github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/query"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
	proj_pb "github.com/caos/zitadel/pkg/grpc/project"
)

func listProjectGrantsRequestToModel(req *mgmt_pb.ListProjectGrantsRequest) (*query.ProjectGrantSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := ProjectGrantQueriesToModel(req)
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

func ProjectGrantQueriesToModel(req *mgmt_pb.ListProjectGrantsRequest) (_ []query.SearchQuery, err error) {
	queries := make([]query.SearchQuery, 0, len(req.Queries)+1)
	for _, query := range req.Queries {
		q, err := ProjectGrantQueryToModel(query)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}
	projectIDQuery, err := query.NewProjectGrantProjectIDSearchQuery(query.TextEquals, req.ProjectId)
	if err != nil {
		return nil, err
	}
	queries = append(queries, projectIDQuery)

	return queries, nil
}

func ProjectGrantQueryToModel(apiQuery *proj_pb.ProjectGrantQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *proj_pb.ProjectGrantQuery_ProjectNameQuery:
		return query.NewProjectGrantProjectNameSearchQuery(object.TextMethodToQuery(q.ProjectNameQuery.Method), q.ProjectNameQuery.Name)
	case *proj_pb.ProjectGrantQuery_RoleKeyQuery:
		return query.NewProjectGrantRoleKeySearchQuery(q.RoleKeyQuery.RoleKey)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "PROJECT-M099f", "List.Query.Invalid")
	}
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
