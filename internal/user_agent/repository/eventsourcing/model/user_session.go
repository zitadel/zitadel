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
	//AuthSessions         []*AuthSession
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

//func (s *UserSession) Changes(changed *UserSession) map[string]interface{} {
//	changes := make(map[string]interface{}, 1)
//	if changed.Name != "" && p.Name != changed.Name {
//		changes["name"] = changed.Name
//	}
//	return changes
//}

func (s *UserSession) setData(event *es_models.Event) error {
	s.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, s); err != nil {
		logging.Log("MODEL-s231F").WithError(err).Debug("could not unmarshal event data")
		return err
	}
	return nil
}

func (a *UserAgent) appendUserSessionAddedEvent(event *es_models.Event) error {
	s := new(UserSession)
	if err := s.setData(event); err != nil {
		return err
	}
	s.State = int32(model.UserSessionStateActive)
	a.UserSessions = append(a.UserSessions, s)
	return nil
}

func (a *UserAgent) appendUserSessionTerminatedEvent(event *es_models.Event) error {
	id, err := userSessionIDFromEvent(event)
	if err != nil {
		return err
	}
	if _, s := GetUserSession(a.UserSessions, id.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateTerminated)
	}
	return nil
}

func (a *UserAgent) appendUserNameCheckSucceededEvent(event *es_models.Event) error {
	id, err := userSessionIDFromEvent(event)
	if err != nil {
		return err
	}
	if _, s := GetUserSession(a.UserSessions, id.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.PasswordFailureCount = 0
		s.PasswordVerified = true
	}
	return nil
}

func (a *UserAgent) appendUserNameCheckFailedEvent(event *es_models.Event) error {
	id, err := userSessionIDFromEvent(event)
	if err != nil {
		return err
	}
	if _, s := GetUserSession(a.UserSessions, id.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.PasswordFailureCount++
		s.PasswordVerified = false
	}
	return nil
}

func (a *UserAgent) appendPasswordCheckSucceededEvent(event *es_models.Event) error {
	id, err := userSessionIDFromEvent(event)
	if err != nil {
		return err
	}
	if _, s := GetUserSession(a.UserSessions, id.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.PasswordFailureCount = 0
		s.PasswordVerified = true
	}
	return nil
}

func (a *UserAgent) appendPasswordCheckFailedEvent(event *es_models.Event) error {
	id, err := userSessionIDFromEvent(event)
	if err != nil {
		return err
	}
	if _, s := GetUserSession(a.UserSessions, id.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.PasswordFailureCount++
		s.PasswordVerified = false
	}
	return nil
}

func (a *UserAgent) appendMfaCheckSucceededEvent(event *es_models.Event) error {
	mfaSession := new(MfaUserSession)
	if err := json.Unmarshal(event.Data, mfaSession); err != nil {
		logging.Log("MODEL-s2gyx").WithError(err).Debug("could not unmarshal event data")
		return err
	}
	if _, s := GetUserSession(a.UserSessions, mfaSession.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.MfaFailureCount = 0
		s.MfaVerified = true
		s.Mfa = mfaSession.MfaType
	}
	return nil
}

func (a *UserAgent) appendMfaCheckFailedEvent(event *es_models.Event) error {
	mfaSession := new(MfaUserSession)
	if err := json.Unmarshal(event.Data, mfaSession); err != nil {
		logging.Log("MODEL-s2gyx").WithError(err).Debug("could not unmarshal event data")
		return err
	}
	if _, s := GetUserSession(a.UserSessions, mfaSession.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.MfaFailureCount++
		s.MfaVerified = false
		s.Mfa = mfaSession.MfaType
	}
	return nil
}

func (a *UserAgent) appendReAuthRequestedEvent(event *es_models.Event) error {
	id, err := userSessionIDFromEvent(event)
	if err != nil {
		return err
	}
	if _, s := GetUserSession(a.UserSessions, id.UserSessionID); s != nil {
		s.State = int32(model.UserSessionStateActive)
		s.PasswordVerified = false
		s.PasswordFailureCount = 0
		s.MfaVerified = false
		s.MfaFailureCount = 0
	}
	return nil
}

func userSessionIDFromEvent(event *es_models.Event) (*UserSessionID, error) {
	id := new(UserSessionID)
	if err := json.Unmarshal(event.Data, id); err != nil {
		logging.Log("MODEL-s231F").WithError(err).Debug("could not unmarshal event data")
		return nil, err
	}
	return id, nil
}
