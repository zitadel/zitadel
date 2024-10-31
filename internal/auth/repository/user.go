package repository

import (
	"context"
)

type UserRepository interface {
	UserSessionsByAgentID(ctx context.Context, agentID string) (sessions map[string]string, err error)
	UserAgentIDBySessionID(ctx context.Context, sessionID string) (string, error)
	ActiveUserIDsBySessionID(ctx context.Context, sessionID string) (userAgentID string, sessions map[string]string, err error)
}
