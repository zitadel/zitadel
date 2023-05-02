package command

import (
	"context"
)

type permissionCheck func(ctx context.Context, permission, orgID, resourceID string) (err error)

const (
	permissionUserWrite     = "user.write"
	permissionSessionWrite  = "session.write"
	permissionSessionDelete = "session.delete"
)
