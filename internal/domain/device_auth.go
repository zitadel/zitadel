package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type DeviceAuthorization struct {
	models.ObjectRoot

	ClientID   string
	DeviceCode string
	UserCode   string
	Expires    time.Time
	Scopes     []string
	Subject    string
	State      DeviceAuthState
}

type DeviceAuthState uint

const (
	DeviceAuthStateInitiated DeviceAuthState = iota
	DeviceAuthStateApproved
	DeviceAuthStateUserDenied
	DeviceAuthStateCompleted
	DeviceAuthStateRemoved
)

func (s DeviceAuthState) Exists() bool {
	return s < DeviceAuthStateRemoved
}
func (s DeviceAuthState) Done() bool {
	return s == DeviceAuthStateApproved
}
func (s DeviceAuthState) Denied() bool {
	return s == DeviceAuthStateUserDenied || s == DeviceAuthStateRemoved
}
