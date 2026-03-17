package start

import "context"

// StartZitadel is an exported wrapper around startZitadel, allowing
// in-process ZITADEL startup from integration test orchestrators.
func StartZitadel(ctx context.Context, config *Config, masterKey string, server chan<- *Server) error {
	return startZitadel(ctx, config, masterKey, server)
}
