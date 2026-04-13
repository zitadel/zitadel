package v2

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func GetSession(ctx context.Context, request *connect.Request[session.GetSessionRequest]) (*connect.Response[session.GetSessionResponse], error) {
	return defaultServer.GetSession(ctx, request)
}

// GetSession implements [sessionconnect.SessionServiceHandler].
func (s *server) GetSession(ctx context.Context, request *connect.Request[session.GetSessionRequest]) (*connect.Response[session.GetSessionResponse], error) {
	panic("unimplemented")
}
