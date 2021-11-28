package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func ListMyLinkedIDPsRequestToModel(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*query.LinkedIDPsSearchQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	q, err := query.NewLinkedIDPsUserIDSearchQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &query.LinkedIDPsSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{q},
	}, nil
}

func RemoveMyLinkedIDPRequestToDomain(ctx context.Context, req *auth_pb.RemoveMyLinkedIDPRequest) *domain.UserIDPLink {
	return &domain.UserIDPLink{
		ObjectRoot:     ctxToObjectRoot(ctx),
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}
