package auth

import "github.com/caos/zitadel/internal/api/grpc/server/middleware"

func (c *ListMyUserChangesResponse) Localizers() []middleware.Localizer {
	if c == nil {
		return nil
	}
	localizers := make([]middleware.Localizer, len(c.Result))
	for i, change := range c.Result {
		localizers[i] = change.EventType
	}
	return localizers
}
