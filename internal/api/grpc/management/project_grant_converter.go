package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	proj_pb "github.com/zitadel/zitadel/pkg/grpc/project"
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
	projectIDQuery, err := query.NewProjectGrantProjectIDSearchQuery(req.ProjectId)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-M099f", "List.Query.Invalid")
	}
}
func listAllProjectGrantsRequestToModel(req *mgmt_pb.ListAllProjectGrantsRequest) (*query.ProjectGrantSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := AllProjectGrantQueriesToModel(req)
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

func AllProjectGrantQueriesToModel(req *mgmt_pb.ListAllProjectGrantsRequest) (_ []query.SearchQuery, err error) {
	queries := make([]query.SearchQuery, 0, len(req.Queries))
	for _, query := range req.Queries {
		q, err := AllProjectGrantQueryToModel(query)
		if err != nil {
			return nil, err
		}
		queries = append(queries, q)
	}
	return queries, nil
}

func AllProjectGrantQueryToModel(apiQuery *proj_pb.AllProjectGrantQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *proj_pb.AllProjectGrantQuery_ProjectNameQuery:
		return query.NewProjectGrantProjectNameSearchQuery(object.TextMethodToQuery(q.ProjectNameQuery.Method), q.ProjectNameQuery.Name)
	case *proj_pb.AllProjectGrantQuery_RoleKeyQuery:
		return query.NewProjectGrantRoleKeySearchQuery(q.RoleKeyQuery.RoleKey)
	case *proj_pb.AllProjectGrantQuery_ProjectIdQuery:
		return query.NewProjectGrantProjectIDSearchQuery(q.ProjectIdQuery.ProjectId)
	case *proj_pb.AllProjectGrantQuery_GrantedOrgIdQuery:
		return query.NewProjectGrantGrantedOrgIDSearchQuery(q.GrantedOrgIdQuery.GrantedOrgId)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-M099f", "List.Query.Invalid")
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

func ListProjectGrantMembersRequestToModel(ctx context.Context, req *mgmt_pb.ListProjectGrantMembersRequest) (*query.ProjectGrantMembersQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.ProjectGrantMembersQuery{
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
		GrantID:   req.GrantId,
		OrgID:     authz.GetCtxData(ctx).OrgID,
	}, nil
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
