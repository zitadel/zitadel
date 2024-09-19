package auth

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/query"
	auth_pb "github.com/zitadel/zitadel/v2/pkg/grpc/auth"
)

func ListMyLinkedIDPsRequestToQuery(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*query.IDPUserLinksSearchQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	q, err := query.NewIDPUserLinksUserIDSearchQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &query.IDPUserLinksSearchQuery{
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
