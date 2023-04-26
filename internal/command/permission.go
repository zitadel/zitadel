package command

import (
	"context"
)

type permissionCheck func(ctx context.Context, permission, orgID, resourceID string, allowSelf bool) (err error)

const (
	permissionUserWrite = "user.write"
)
