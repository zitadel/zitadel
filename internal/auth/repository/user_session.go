package repository

import (
	"context"
	"github.com/caos/zitadel/internal/user/model"
)

type UserSessionRepository interface {
	GetMyUserSessions(ctx context.Context) ([]*model.UserSessionView, error)
	ActiveUserSessionCount() int64
}
