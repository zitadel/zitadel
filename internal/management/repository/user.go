package repository

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.UserChanges, error)
}
