package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type UserRepository interface {
	UserSessionsByAgentID(ctx context.Context, agentID string) (sessions []command.HumanSignOutSession, err error)
	UserAgentIDBySessionID(ctx context.Context, sessionID string) (string, error)
	UserSessionByID(ctx context.Context, sessionID string) (*model.UserSessionView, error)
	ActiveUserSessionsBySessionID(ctx context.Context, sessionID string) (userAgentID string, sessions []command.HumanSignOutSession, err error)
}
