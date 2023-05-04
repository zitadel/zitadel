package session

import (
	"context"
)

type permissionCheck func(ctx context.Context, permission, orgID, resourceID string) (err error)

const (
	permissionSessionRead = "session.read"
)
