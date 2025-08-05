package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/user/model"
)

type UserSessionRepository interface {
	GetMyUserSessions(ctx context.Context) ([]*model.UserSessionView, error)
}
