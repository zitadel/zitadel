package session

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	user "github.com/zitadel/zitadel/internal/api/grpc/user/v2alpha"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
)

func (s *Server) CreateSession(ctx context.Context, req *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	return &session.CreateSessionResponse{
		SessionId: "hodor",
	}, nil
}

func (s *Server) GetSession(ctx context.Context, req *session.GetSessionRequest) (*session.GetSessionResponse, error) {
	return &session.GetSessionResponse{
		Session: &session.Session{
			Id:   req.SessionId,
			User: &user.User{Id: authz.GetCtxData(ctx).UserID},
		},
	}, nil
}
