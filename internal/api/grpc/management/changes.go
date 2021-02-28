package management

import (
	"github.com/caos/zitadel/internal/api/grpc/server/middleware"
)

func (c *ListUserChangesResponse) Localizers() []middleware.Localizer {
	if c == nil {
		return nil
	}
	localizers := make([]middleware.Localizer, len(c.Result))
	for i, change := range c.Result {
		localizers[i] = change.EventType
	}
	return localizers
}

func (c *ListOrgChangesResponse) Localizers() []middleware.Localizer {
	if c == nil {
		return nil
	}
	localizers := make([]middleware.Localizer, len(c.Result))
	for i, change := range c.Result {
		localizers[i] = change.EventType
	}
	return localizers
}

func (c *ListProjectChangesResponse) Localizers() []middleware.Localizer {
	if c == nil {
		return nil
	}
	localizers := make([]middleware.Localizer, len(c.Result))
	for i, change := range c.Result {
		localizers[i] = change.EventType
	}
	return localizers
}

func (c *ListAppChangesResponse) Localizers() []middleware.Localizer {
	if c == nil {
		return nil
	}
	localizers := make([]middleware.Localizer, len(c.Result))
	for i, change := range c.Result {
		localizers[i] = change.EventType
	}
	return localizers
}
