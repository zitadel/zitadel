package connect_middleware

import (
	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

func LogHandler(ignoredPrefix ...string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return logging.NewConnectInterceptor(next, ignoredPrefix...)
	}
}
