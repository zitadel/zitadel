package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/model"
)

type UserSession struct {
	es_models.ObjectRoot

	UserID string
	//UserName string
	SessionID            string
	State                int32
	PasswordVerified     bool
	PasswordFailureCount uint16
	Mfa                  int32
	MfaVerified          bool
	MfaFailureCount      uint16
	AuthTime             time.Time
}

type UserSessionID struct {
	es_models.ObjectRoot
	UserSessionID string `json:"userSessionID"`
}

type MfaUserSession struct {
	es_models.ObjectRoot
	UserSessionID string `json:"userSessionID"`
	MfaType       int32  `json:"mfaType"`
}

type SetUserSession struct {
	es_models.ObjectRoot
	UserSessionID string `json:"userSessionID"`
	AuthSessionID string `json:"authSessionID"`
}

func UserSessionsFromModel(sessions []*model.UserSession) []*UserSession {
	convertedSessions := make([]*UserSession, len(sessions))
	for i, session := range sessions {
		convertedSessions[i] = UserSessionFromModel(session)
	}
	return convertedSessions
}

func UserSessionsToModel(sessions []*UserSession) []*model.UserSession {
	convertedSessions := make([]*model.UserSession, len(sessions))
	for i, session := range sessions {
		convertedSessions[i] = UserSessionToModel(session)
	}
	return convertedSessions
}

func UserSessionFromModel(userSession *model.UserSession) *UserSession {
	return &UserSession{
		ObjectRoot:           userSession.ObjectRoot,
		UserID:               userSession.UserID,
		SessionID:            userSession.SessionID,
		State:                int32(userSession.State),
		PasswordVerified:     userSession.PasswordVerified,
		PasswordFailureCount: userSession.PasswordFailureCount,
		Mfa:                  int32(userSession.Mfa),
		MfaVerified:          userSession.MfaVerified,
		MfaFailureCount:      userSession.MfaFailureCount,
		AuthTime:             userSession.AuthTime,
	}
}

func UserSessionToModel(userSession *UserSession) *model.UserSession {
	return &model.UserSession{
		ObjectRoot:           userSession.ObjectRoot,
		UserID:               userSession.UserID,
		SessionID:            userSession.SessionID,
		State:                model.UserSessionState(userSession.State),
		PasswordVerified:     userSession.PasswordVerified,
		PasswordFailureCount: userSession.PasswordFailureCount,
		Mfa:                  model.MfaType(userSession.Mfa),
		MfaVerified:          userSession.MfaVerified,
		MfaFailureCount:      userSession.MfaFailureCount,
		AuthTime:             userSession.AuthTime,
	}
}

func GetUserSession(sessions []*UserSession, id string) (int, *UserSession) {
	for i, s := range sessions {
		if s.SessionID == id {
			return i, s
		}
	}
	return -1, nil
}

func (s *UserSession) Changes(changed *UserSession) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Name != "" && p.Name != changed.Name {
		changes["name"] = changed.Name
	}
	return changes
}

func (s *UserSession) getData(event *es_models.Event) error {
	s.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, s); err != nil {
		logging.Log("MODEL-s231F").WithError(err).Debug("could not unmarshal event data")
		return err
	}
	return nil
}
