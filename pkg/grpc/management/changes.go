package management

import (
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

func (c *Changes) Localizers() []middleware.Localizer {
	if c == nil {
		return nil
	}
	localizers := make([]middleware.Localizer, len(c.Changes))
	for i, change := range c.Changes {
		localizers[i] = change.EventType
	}
	return localizers
}
