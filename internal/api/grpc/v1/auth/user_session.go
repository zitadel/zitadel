package auth

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyUserSessions(ctx context.Context, _ *empty.Empty) (_ *auth.UserSessionViews, err error) {
	userSessions, err := s.repo.GetMyUserSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &auth.UserSessionViews{UserSessions: userSessionViewsFromModel(userSessions)}, nil
}
