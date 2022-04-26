package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type OIDCSettings struct {
	models.ObjectRoot

	State   OIDCSettingsState
	Default bool

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
}

type OIDCSettingsState int32

const (
	OIDCSettingsStateUnspecified OIDCSettingsState = iota
	OIDCSettingsStateActive
	OIDCSettingsStateRemoved

	oidcSettingsStateCount
)

func (c OIDCSettingsState) Valid() bool {
	return c >= 0 && c < oidcSettingsStateCount
}

func (s OIDCSettingsState) Exists() bool {
	return s != OIDCSettingsStateUnspecified && s != OIDCSettingsStateRemoved
}
