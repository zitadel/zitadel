package authorization

import (
	"context"

	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
)

func (s *Server) ListAuthorizations(ctx context.Context, request *authorization.ListAuthorizationsRequest) (*authorization.ListAuthorizationsResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Server) GetAuthorization(ctx context.Context, request *authorization.GetAuthorizationRequest) (*authorization.GetAuthorizationResponse, error) {
	// TODO implement me
	panic("implement me")
}
