package auth

import (
	auth_req_model "github.com/caos/zitadel/internal/auth_request/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/auth/grpc"
)

func userSessionViewsFromModel(userSessions []*usr_model.UserSessionView) []*grpc.UserSessionView {
	converted := make([]*grpc.UserSessionView, len(userSessions))
	for i, s := range userSessions {
		converted[i] = userSessionViewFromModel(s)
	}
	return converted
}

func userSessionViewFromModel(userSession *usr_model.UserSessionView) *grpc.UserSessionView {
	return &grpc.UserSessionView{
		Sequence:    userSession.Sequence,
		AgentId:     userSession.UserAgentID,
		UserId:      userSession.UserID,
		UserName:    userSession.UserName,
		LoginName:   userSession.LoginName,
		DisplayName: userSession.DisplayName,
		AuthState:   userSessionStateFromModel(userSession.State),
	}
}

func userSessionStateFromModel(state auth_req_model.UserSessionState) grpc.UserSessionState {
	switch state {
	case auth_req_model.UserSessionStateActive:
		return grpc.UserSessionState_USERSESSIONSTATE_ACTIVE
	case auth_req_model.UserSessionStateTerminated:
		return grpc.UserSessionState_USERSESSIONSTATE_TERMINATED
	default:
		return grpc.UserSessionState_USERSESSIONSTATE_UNSPECIFIED
	}
}
