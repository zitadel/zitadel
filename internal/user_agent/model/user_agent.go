package model

import (
	"net"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type UserAgent struct {
	es_models.ObjectRoot

	UserAgent      string
	AcceptLanguage string
	RemoteIP       net.IP
	State          UserAgentState
	UserSessions   []*UserSession
	//LastUsedUserSession ?
}

type UserAgentState int32

const (
	UserAgentStateActive UserAgentState = iota
	UserAgentStateInctive
)

func NewUserAgent(id, userAgent, acceptLanguage string, remoteIP net.IP) *UserAgent {
	return &UserAgent{
		ObjectRoot:     es_models.ObjectRoot{ID: id},
		UserAgent:      userAgent,
		AcceptLanguage: acceptLanguage,
		RemoteIP:       remoteIP,
		State:          UserAgentStateActive,
	}
}

func (u *UserAgent) IsActive() bool {
	return u.State == UserAgentStateActive
}

func (u *UserAgent) IsValid() bool {
	return u.UserAgent != "" &&
		u.AcceptLanguage != "" &&
		u.RemoteIP != nil && !u.RemoteIP.IsUnspecified()
}
