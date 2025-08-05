//go:build !integration

package sink

import "github.com/zitadel/zitadel/internal/command"

// StartServer and its returned close function are a no-op
// when the `integration` build tag is disabled.
func StartServer(cmd *command.Commands) (close func()) {
	return func() {}
}
