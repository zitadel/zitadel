package repository

import (
	"context"

	"github.com/caos/zitadel/internal/user_agent/model"
)

type UserRepository interface {
	userAuthorizationRepository
}

type userAuthorizationRepository interface {
	VerifyPassword(ctx context.Context, authRequestID, password string, info *model.BrowserInfo) (*model.AuthSession, error)
	VerifyMfa(ctx context.Context, agentID, authRequestID string, mfa interface{}, info *model.BrowserInfo) (*model.AuthSession, error)
}
