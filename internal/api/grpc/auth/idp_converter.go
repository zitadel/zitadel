package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func ListMyLinkedIDPsRequestToQuery(ctx context.Context, req *auth_pb.ListMyLinkedIDPsRequest) (*query.IDPUserLinksSearchQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	q := make([]query.SearchQuery, 0, 2)
	userIDQuery, err := query.NewIDPUserLinksUserIDSearchQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	q = append(q, userIDQuery)
	if !req.GetWithLinkedLoginPolicyOnly() {
		withPolicyQuery, err := query.NewIDPUserLinksWithLoginPolicyOnlySearchQuery()
		if err != nil {
			return nil, err
		}
		q = append(q, withPolicyQuery)
	}
	return &query.IDPUserLinksSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: q,
	}, nil
}

func RemoveMyLinkedIDPRequestToDomain(ctx context.Context, req *auth_pb.RemoveMyLinkedIDPRequest) *domain.UserIDPLink {
	return &domain.UserIDPLink{
		ObjectRoot:     ctxToObjectRoot(ctx),
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}
