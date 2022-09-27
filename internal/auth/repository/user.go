package repository

import (
	"context"
)

type UserRepository interface {
	UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error)
}
