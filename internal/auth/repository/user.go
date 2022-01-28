package repository

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	myUserRepo

	UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error)
}

type myUserRepo interface {
	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.UserChanges, error)
}
