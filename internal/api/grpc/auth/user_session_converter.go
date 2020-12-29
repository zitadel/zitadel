package auth

import (
	auth_req_model "github.com/caos/zitadel/internal/auth_request/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func userSessionViewsFromModel(userSessions []*usr_model.UserSessionView) []*auth.UserSessionView {
	converted := make([]*auth.UserSessionView, len(userSessions))
	for i, s := range userSessions {
		converted[i] = userSessionViewFromModel(s)
	}
	return converted
}

func userSessionViewFromModel(userSession *usr_model.UserSessionView) *auth.UserSessionView {
	return &auth.UserSessionView{
		Sequence:    userSession.Sequence,
		AgentId:     userSession.UserAgentID,
		UserId:      userSession.UserID,
		UserName:    userSession.UserName,
		LoginName:   userSession.LoginName,
		DisplayName: userSession.DisplayName,
		AuthState:   userSessionStateFromModel(userSession.State),
	}
}

func userSessionStateFromModel(state auth_req_model.UserSessionState) auth.UserSessionState {
	switch state {
	case auth_req_model.UserSessionStateActive:
		return auth.UserSessionState_USERSESSIONSTATE_ACTIVE
	case auth_req_model.UserSessionStateTerminated:
		return auth.UserSessionState_USERSESSIONSTATE_TERMINATED
	case auth_req_model.UserSessionStateInitiated:
		return auth.UserSessionState_USERSESSIONSTATE_INITIATED
	default:
		return auth.UserSessionState_USERSESSIONSTATE_UNSPECIFIED
	}
}
