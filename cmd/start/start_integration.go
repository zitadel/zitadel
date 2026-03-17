package start

import (
	"context"
)

// StartZitadel wraps the private [startZitadel] function so that it can
// be called from orchestrator tests in the tests/integration/ package.
// The server channel receives a [*Server] once the startup sequence
// completes; callers should read from it to obtain the running server
// reference.
func StartZitadel(ctx context.Context, config *Config, masterKey string, server chan<- *Server) error {
	return startZitadel(ctx, config, masterKey, server)
}
