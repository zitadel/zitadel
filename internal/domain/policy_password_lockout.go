package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	Default             bool
	MaxAttempts         uint64
	ShowLockOutFailures bool
}
