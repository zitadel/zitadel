package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) GetApplicationByID(ctx context.Context, request *ApplicationID) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-Rfh8e", "Not implemented")
}

func (s *Server) SearchApplications(ctx context.Context, appSearch *ApplicationSearchRequest) (*ApplicationSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ju8Rd", "Not implemented")
}

func (s *Server) AuthorizeApplication(ctx context.Context, auth *ApplicationAuthorizeRequest) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-lo8ws", "Not implemented")
}
