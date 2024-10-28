package repository

import (
	"context"
)

type UserRepository interface {
	UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error)
	UserAgentIDBySessionID(ctx context.Context, sessionID string) (string, error)
	ActiveUserIDsBySessionID(ctx context.Context, sessionID string) (userAgentID string, userIDs []string, err error)
	UserSessionsByAgentID(ctx context.Context, agentID string) ([]string, error)
}
