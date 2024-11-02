package domain

import (
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type LockoutPolicy struct {
	models.ObjectRoot

	Default             bool
	MaxPasswordAttempts uint64
	MaxOTPAttempts      uint64
	ShowLockOutFailures bool
}
