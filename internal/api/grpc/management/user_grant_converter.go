package management

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/v2/internal/api/grpc/user"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/v2/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/v2/pkg/grpc/management"
	"github.com/zitadel/zitadel/v2/pkg/grpc/user"
)

func ListUserGrantsRequestToQuery(ctx context.Context, req *mgmt_pb.ListUserGrantRequest) (*query.UserGrantsQueries, error) {
	queries, err := user_grpc.UserGrantQueriesToQuery(ctx, req.Queries)
	if err != nil {
		return nil, err
	}

	if shouldAppendUserGrantOwnerQuery(req.Queries) {
		ownerQuery, err := query.NewUserGrantResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
		if err != nil {
			return nil, err
		}
		queries = append(queries, ownerQuery)
	}

	offset, limit, asc := object.ListQueryToModel(req.Query)
	request := &query.UserGrantsQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}

	return request, nil
}

func shouldAppendUserGrantOwnerQuery(queries []*user.UserGrantQuery) bool {
	for _, query := range queries {
		if _, ok := query.Query.(*user.UserGrantQuery_WithGrantedQuery); ok {
			return false
		}
	}
	return true
}

func AddUserGrantRequestToDomain(req *mgmt_pb.AddUserGrantRequest) *domain.UserGrant {
	return &domain.UserGrant{
		UserID:         req.UserId,
		ProjectID:      req.ProjectId,
		ProjectGrantID: req.ProjectGrantId,
		RoleKeys:       req.RoleKeys,
	}
}

func UpdateUserGrantRequestToDomain(req *mgmt_pb.UpdateUserGrantRequest) *domain.UserGrant {
	return &domain.UserGrant{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.GrantId,
		},
		UserID:   req.UserId,
		RoleKeys: req.RoleKeys,
	}

}
