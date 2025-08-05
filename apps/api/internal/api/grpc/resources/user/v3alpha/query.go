package user

import (
	"context"

	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) SearchUsers(ctx context.Context, _ *user.SearchUsersRequest) (_ *user.SearchUsersResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	return &user.SearchUsersResponse{}, nil
}
