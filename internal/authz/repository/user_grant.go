package repository

import "context"

type UserGrantRepository interface {
	SearchMyZitadelPermissions(ctx context.Context) ([]string, error)
}
