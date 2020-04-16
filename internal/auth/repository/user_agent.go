package repository

import (
	"context"

	"github.com/caos/zitadel/internal/user_agent/model"
)

type UserAgentRepository interface {
	UserAgentByID(ctx context.Context, id string) (*model.UserAgent, error)
	CreateUserAgent(ctx context.Context, info *model.BrowserInfo) (*model.UserAgent, error)
	RevokeUserAgent(ctx context.Context, id string) (*model.UserAgent, error)
	CreateAuthSession(ctx context.Context, session *model.AuthSession) (*model.AuthSession, error)
	GetAuthSession(ctx context.Context, id, agentID string, info *model.BrowserInfo) (*model.AuthSession, error)
	GetAuthSessionByTokenID(ctx context.Context, tokenID string) (*model.AuthSession, error) //view?
	SelectUser(ctx context.Context, agentID, authSessionID, userSessionID string, info *model.BrowserInfo) (*model.AuthSession, error)
	VerifyUser(ctx context.Context, agentID, authSessionID, userName string, info *model.BrowserInfo) (*model.AuthSession, error)
	VerifyPassword(ctx context.Context, agentID, authSessionID, password string, info *model.BrowserInfo) (*model.AuthSession, error)
	VerifyMfa(ctx context.Context, agentID, authSessionID string, mfa interface{}, info *model.BrowserInfo) (*model.AuthSession, error)
	GetUserSessions(ctx context.Context, agentID string) ([]*model.UserSession, error)
	GetUserSessionByID(ctx context.Context, agentID, sessionID string) (*model.UserSession, error)
	GetUserSessionByUserID(ctx context.Context) (*model.UserAgent, error) //my? view?
	TerminateUserSession(ctx context.Context, agentID, sessionID string) error
	CreateToken(ctx context.Context, agentID, authSessionID string) (*model.Token, error)
}
