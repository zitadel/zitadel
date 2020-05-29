package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) SearchUserGrant(ctx context.Context, grantSearch *UserGrantSearchRequest) (*UserGrantSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-s8iSf", "Not implemented")
}

func (s *Server) SearchMyProjectOrgs(ctx context.Context, request *MyProjectOrgSearchRequest) (*MyProjectOrgSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-8kdRf", "Not implemented")
}
