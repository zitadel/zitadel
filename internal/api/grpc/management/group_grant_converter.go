package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	group_grpc "github.com/zitadel/zitadel/internal/api/grpc/group"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/group"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func ListGroupGrantsRequestToQuery(ctx context.Context, req *mgmt_pb.ListGroupGrantRequest) (*query.GroupGrantsQueries, error) {
	queries, err := group_grpc.GroupGrantQueriesToQuery(ctx, req.Queries)
	if err != nil {
		return nil, err
	}

	if shouldAppendGroupGrantOwnerQuery(req.Queries) {
		ownerQuery, err := query.NewUserGrantResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
		if err != nil {
			return nil, err
		}
		queries = append(queries, ownerQuery)
	}

	offset, limit, asc := object.ListQueryToModel(req.Query)
	request := &query.GroupGrantsQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}

	return request, nil
}

func shouldAppendGroupGrantOwnerQuery(queries []*group.GroupGrantQuery) bool {
	for _, query := range queries {
		if _, ok := query.Query.(*group.GroupGrantQuery_WithGrantedQuery); ok {
			return false
		}
	}
	return true
}

func AddGroupGrantRequestToDomain(req *mgmt_pb.AddGroupGrantRequest) *domain.GroupGrant {
	return &domain.GroupGrant{
		GroupID:        req.GroupId,
		ProjectID:      req.ProjectId,
		ProjectGrantID: req.ProjectGrantId,
		RoleKeys:       req.RoleKeys,
	}
}

func UpdateGroupGrantRequestToDomain(req *mgmt_pb.UpdateGroupGrantRequest) *domain.GroupGrant {
	return &domain.GroupGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GrantId,
		},
		GroupID:  req.GroupId,
		RoleKeys: req.RoleKeys,
	}

}
