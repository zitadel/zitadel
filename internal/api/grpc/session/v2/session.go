package session

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) GetSession(ctx context.Context, req *session.GetSessionRequest) (*session.GetSessionResponse, error) {
	return &session.GetSessionResponse{
		Session: &session.Session{
			Id:   req.Id,
			User: &user.User{Id: authz.GetCtxData(ctx).UserID},
		},
	}, nil
}
