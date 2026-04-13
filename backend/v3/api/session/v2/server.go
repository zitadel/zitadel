package v2

import (
	"github.com/zitadel/zitadel/pkg/grpc/session/v2/sessionconnect"
)

var defaultServer = new(server)

type server struct {
	sessionconnect.UnimplementedSessionServiceHandler
}

var _ sessionconnect.SessionServiceHandler = (*server)(nil)
