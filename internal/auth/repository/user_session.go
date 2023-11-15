package repository

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/user/model"
)

type UserSessionRepository interface {
	GetMyUserSessions(ctx context.Context) ([]*model.UserSessionView, error)
	ActiveUserSessionCount() int64
}
