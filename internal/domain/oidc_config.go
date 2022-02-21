package domain

import (
	"time"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type OIDCConfig struct {
	models.ObjectRoot

	State   OIDCConfigState
	Default bool

	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
}

type OIDCConfigState int32

const (
	OIDCConfigStateUnspecified OIDCConfigState = iota
	OIDCConfigStateActive
	OIDCConfigStateRemoved

	oidcConfigStateCount
)

func (c OIDCConfig) IsValid() error {

	return nil
}

func (c OIDCConfigState) Valid() bool {
	return c >= 0 && c < oidcConfigStateCount
}

func (s OIDCConfigState) Exists() bool {
	return s != OIDCConfigStateUnspecified && s != OIDCConfigStateRemoved
}
