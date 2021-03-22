package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	auth_req_model "github.com/caos/zitadel/internal/auth_request/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/user"
)

func UserSessionsToPb(sessions []*user_model.UserSessionView) []*user.Session {
	s := make([]*user.Session, len(sessions))
	for i, session := range sessions {
		s[i] = UserSessionToPb(session)
	}
	return s
}

func UserSessionToPb(session *user_model.UserSessionView) *user.Session {
	return &user.Session{
		// SessionId: session.,//TOOD: not return from be
		AgentId:     session.UserAgentID,
		UserId:      session.UserID,
		UserName:    session.UserName,
		LoginName:   session.LoginName,
		DisplayName: session.DisplayName,
		AuthState:   SessionStateToPb(session.State),
		Details: object.ToViewDetailsPb(
			session.Sequence,
			session.CreationDate,
			session.ChangeDate,
			session.ResourceOwner,
		),
	}
}

func SessionStateToPb(state auth_req_model.UserSessionState) user.SessionState {
	switch state {
	case auth_req_model.UserSessionStateActive:
		return user.SessionState_SESSION_STATE_ACTIVE
	case auth_req_model.UserSessionStateTerminated:
		return user.SessionState_SESSION_STATE_TERMINATED
	default:
		return user.SessionState_SESSION_STATE_UNSPECIFIED
	}
}
