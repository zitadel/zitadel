package eventsourcing

import (
	"context"

	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	user_agent_model "github.com/caos/zitadel/internal/user_agent/model"
	user_agent_event "github.com/caos/zitadel/internal/user_agent/repository/eventsourcing"
	"github.com/caos/zitadel/internal/user_agent/repository/eventsourcing/model"
)

type UserAgentRepo struct {
	UserAgentEvents *user_agent_event.UserAgentEventstore
	UserEvents      *user_event.UserEventstore
	//view      *view.View
}

func (repo *UserAgentRepo) UserAgentByID(ctx context.Context, id string) (*user_agent_model.UserAgent, error) {
	return repo.UserAgentEvents.UserAgentByID(ctx, id)
}

func (repo *UserAgentRepo) CreateUserAgent(ctx context.Context, info *user_agent_model.BrowserInfo) (*user_agent_model.UserAgent, error) {
	agent := user_agent_model.NewUserAgent(info.UserAgent, info.AcceptLanguage, info.RemoteIP)
	return repo.UserAgentEvents.CreateUserAgent(ctx, agent)
}

func (repo *UserAgentRepo) RevokeUserAgent(ctx context.Context, id string) (*user_agent_model.UserAgent, error) {
	return repo.UserAgentEvents.RevokeUserAgent(ctx, id)
}

func (repo *UserAgentRepo) CreateAuthSession(ctx context.Context, session *user_agent_model.AuthSession) (*user_agent_model.AuthSession, error) {
	return repo.UserAgentEvents.AuthSessionAdded(ctx, session)
}

func (repo *UserAgentRepo) GetAuthSession(ctx context.Context, id, agentID string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
	//return repo.UserAgentEvents.Auth(ctx, session)
}
func (repo *UserAgentRepo) GetAuthSessionByTokenID(ctx context.Context, tokenID string) (*user_agent_model.AuthSession, error) { //view?
}

func (repo *UserAgentRepo) SelectUser(ctx context.Context, agentID, authSessionID, userSessionID string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
	return repo.UserAgentEvents.AuthSessionSetUserSession(ctx, agentID, userSessionID, authSessionID)
}

func (repo *UserAgentRepo) VerifyUser(ctx context.Context, agentID, authSessionID, userName string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
	//return repo.UserAgentEvents.Usern(ctx, agentID)
}
func (repo *UserAgentRepo) VerifyPassword(ctx context.Context, agentID, authSessionID, password string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {

	return repo.UserAgentEvents.PasswordCheckSucceeded(ctx, agentID, user)
}
func (repo *UserAgentRepo) VerifyMfa(ctx context.Context, agentID, authSessionID string, mfa interface{}, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error)
func (repo *UserAgentRepo) GetUserSessions(ctx context.Context, agentID string) ([]*user_agent_model.UserSession, error)
func (repo *UserAgentRepo) GetUserSessionByID(ctx context.Context, agentID, sessionID string) (*user_agent_model.UserSession, error)
func (repo *UserAgentRepo) GetUserSessionByUserID(ctx context.Context) (*user_agent_model.UserAgent, error) //my? view?
func (repo *UserAgentRepo) TerminateUserSession(ctx context.Context, agentID, sessionID string) error
func (repo *UserAgentRepo) CreateToken(ctx context.Context, agentID, authSessionID string) (*user_agent_model.Token, error)
