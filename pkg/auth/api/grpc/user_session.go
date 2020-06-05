package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMyUserSessions(ctx context.Context, _ *empty.Empty) (_ *UserSessionViews, err error) {
	userSessions, err := s.repo.GetMyUserSessions(ctx)
	if err != nil {
		return nil, err
	}
	return &UserSessionViews{UserSessions: userSessionViewsFromModel(userSessions)}, nil
}
