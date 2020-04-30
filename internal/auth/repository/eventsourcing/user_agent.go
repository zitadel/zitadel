package eventsourcing

import (
	"context"

	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	user_agent_model "github.com/caos/zitadel/internal/user_agent/model"
	user_agent_event "github.com/caos/zitadel/internal/user_agent/repository/eventsourcing"
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

func (repo *UserAgentRepo) GetAuthSession(ctx context.Context, agentID, userSessionID, authSessionID string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
	//return repo.UserAgentEvents.Auth(ctx, session)
}

//func (repo *UserAgentRepo) SelectUser(ctx context.Context, agentID, authSessionID, userSessionID string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
//	return repo.UserAgentEvents.AuthSessionSetUserSession(ctx, agentID, userSessionID, authSessionID)
//}

//func (repo *UserAgentRepo) VerifyUser(ctx context.Context, agentID, authSessionID, userName string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
//	user, err := repo.UserByUsername(ctx, userName)
//	if err != nil {
//		return nil, err
//	}
//	return repo.UserAgentEvents.AddUserSession(ctx, agentID, authSessionID, user_agent_model.NewUserSession(agentID, "", user.AggregateID))
//}
//func (repo *UserAgentRepo) VerifyPassword(ctx context.Context, agentID, userSessionID, authSessionID, password string, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
//	authSession, err := repo.UserAgentEvents.AuthSessionByIDs(ctx, agentID, userSessionID, authSessionID)
//	if err != nil {
//		return nil, err
//	}
//	if err := repo.UserEvents.VerifyPassword(ctx, authSession.UserSession.UserID, password); err == nil {
//		return repo.UserAgentEvents.PasswordCheckSucceeded(ctx, agentID, authSession.UserSession.SessionID, authSessionID)
//	}
//	return repo.UserAgentEvents.PasswordCheckFailed(ctx, agentID, authSession.UserSession.SessionID, authSessionID)
//}
//func (repo *UserAgentRepo) VerifyMfa(ctx context.Context, agentID, userSessionID, authSessionID string, mfa int32, info *user_agent_model.BrowserInfo) (*user_agent_model.AuthSession, error) {
//	authSession, err := repo.UserAgentEvents.AuthSessionByIDs(ctx, agentID, userSessionID, authSessionID)
//	if err != nil {
//		return nil, err
//	}
//	if err := repo.UserEvents.CheckMfaOTP(ctx, authSession.UserSession.UserID, mfa); err == nil {
//		return repo.UserAgentEvents.MfaCheckSucceeded(ctx, agentID, authSession.UserSession.SessionID, authSessionID, mfa)
//	}
//	return repo.UserAgentEvents.MfaCheckFailed(ctx, agentID, authSession.UserSession.SessionID, authSessionID, mfa)
//}

func (repo *UserAgentRepo) TerminateUserSession(ctx context.Context, agentID, sessionID string) error {
	return repo.UserAgentEvents.TerminateUserSession(ctx, agentID, sessionID)
}
func (repo *UserAgentRepo) CreateToken(ctx context.Context, agentID, authSessionID string) (*user_agent_model.Token, error) {
	return repo.UserAgentEvents.CreateToken(ctx, agentID, "", authSessionID)
}
