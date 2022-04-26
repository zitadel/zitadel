package user

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/pkg/grpc/user"
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
		AvatarUrl:   session.AvatarURL,
		Details: object.ToViewDetailsPb(
			session.Sequence,
			session.CreationDate,
			session.ChangeDate,
			session.ResourceOwner,
		),
	}
}

func SessionStateToPb(state domain.UserSessionState) user.SessionState {
	switch state {
	case domain.UserSessionStateActive:
		return user.SessionState_SESSION_STATE_ACTIVE
	case domain.UserSessionStateTerminated:
		return user.SessionState_SESSION_STATE_TERMINATED
	default:
		return user.SessionState_SESSION_STATE_UNSPECIFIED
	}
}
