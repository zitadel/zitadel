package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	State               PolicyState
	MaxAttempts         uint64
	ShowLockOutFailures bool
}
