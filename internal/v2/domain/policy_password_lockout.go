package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	MaxAttempts         uint64
	ShowLockOutFailures bool
}
