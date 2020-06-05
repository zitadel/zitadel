package grpc

import (
	auth_req_model "github.com/caos/zitadel/internal/auth_request/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

func userSessionViewsFromModel(userSessions []*usr_model.UserSessionView) []*UserSessionView {
	converted := make([]*UserSessionView, len(userSessions))
	for i, s := range userSessions {
		converted[i] = userSessionViewFromModel(s)
	}
	return converted
}

func userSessionViewFromModel(userSession *usr_model.UserSessionView) *UserSessionView {
	return &UserSessionView{
		Sequence:  userSession.Sequence,
		AgentId:   userSession.UserAgentID,
		UserId:    userSession.UserID,
		UserName:  userSession.UserName,
		AuthState: userSessionStateFromModel(userSession.State),
	}
}

func userSessionStateFromModel(state auth_req_model.UserSessionState) UserSessionState {
	switch state {
	case auth_req_model.UserSessionStateActive:
		return UserSessionState_USERSESSIONSTATE_ACTIVE
	case auth_req_model.UserSessionStateTerminated:
		return UserSessionState_USERSESSIONSTATE_TERMINATED
	default:
		return UserSessionState_USERSESSIONSTATE_UNSPECIFIED
	}
}
