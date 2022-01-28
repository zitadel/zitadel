package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyRefreshTokens(ctx context.Context, req *auth.ListMyRefreshTokensRequest) (*auth.ListMyRefreshTokensResponse, error) {
	res, err := s.repo.SearchMyRefreshTokens(ctx, authz.GetCtxData(ctx).UserID, ListMyRefreshTokensRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &auth.ListMyRefreshTokensResponse{
		Result: user_grpc.RefreshTokensToPb(res.Result),
		Details: object.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}

func (s *Server) RevokeMyRefreshToken(ctx context.Context, req *auth.RevokeMyRefreshTokenRequest) (*auth.RevokeMyRefreshTokenResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	details, err := s.command.RevokeRefreshToken(ctx, ctxData.UserID, ctxData.ResourceOwner, req.Id)
	if err != nil {
		return nil, err
	}
	return &auth.RevokeMyRefreshTokenResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RevokeAllMyRefreshTokens(ctx context.Context, _ *auth.RevokeAllMyRefreshTokensRequest) (*auth.RevokeAllMyRefreshTokensResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	res, err := s.repo.SearchMyRefreshTokens(ctx, ctxData.UserID, ListMyRefreshTokensRequestToModel(nil))
	if err != nil {
		return nil, err
	}
	tokenIDs := make([]string, len(res.Result))
	for i, view := range res.Result {
		tokenIDs[i] = view.ID
	}
	err = s.command.RevokeRefreshTokens(ctx, ctxData.UserID, ctxData.ResourceOwner, tokenIDs)
	if err != nil {
		return nil, err
	}
	return &auth.RevokeAllMyRefreshTokensResponse{}, nil
}

func ListMyRefreshTokensRequestToModel(_ *auth.ListMyRefreshTokensRequest) *model.RefreshTokenSearchRequest {
	return &model.RefreshTokenSearchRequest{} //add sorting, queries, ... when possible
}
