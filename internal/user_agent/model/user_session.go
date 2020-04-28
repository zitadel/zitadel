package model

import (
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type UserSession struct {
	es_models.ObjectRoot

	UserID string
	//UserName string
	SessionID            string
	State                UserSessionState
	PasswordVerified     bool
	PasswordFailureCount uint16
	Mfa                  MfaType
	MfaVerified          bool
	MfaFailureCount      uint16
	AuthTime             time.Time
	AuthSessions         []*AuthSession
}

type UserSessionState int32

const (
	UserSessionStateActive UserSessionState = iota
	UserSessionStateTerminated
)

type MfaType int32

const (
	MfaTypeNone MfaType = iota
	MfaTypeOTP
	MFaTypeSMS
)

func NewUserSession(agentID, sessionID string, userID string) *UserSession {
	return &UserSession{
		ObjectRoot: es_models.ObjectRoot{AggregateID: agentID},
		UserID:     userID,
		SessionID:  sessionID,
		State:      UserSessionStateActive,
	}
}

func (u *UserSession) IsValid() bool {
	return true
}
