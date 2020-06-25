package auth

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/auth/grpc"
)

func (s *Server) GetMyUserSessions(ctx context.Context, _ *empty.Empty) (_ *grpc.UserSessionViews, err error) {
	userSessions, err := s.repo.GetMyUserSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &grpc.UserSessionViews{UserSessions: userSessionViewsFromModel(userSessions)}, nil
}
