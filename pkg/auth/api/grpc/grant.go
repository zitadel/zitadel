package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchMyProjectOrgs(ctx context.Context, request *MyProjectOrgSearchRequest) (*MyProjectOrgSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-8kdRf", "Not implemented")
}

func (s *Server) IsIamAdmin(ctx context.Context, _ *empty.Empty) (*IsAdminResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-9odFv", "Not implemented")
}
