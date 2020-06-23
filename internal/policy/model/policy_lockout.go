package model

import "github.com/caos/zitadel/internal/eventstore/models"

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	Description         string
	State               PolicyState
	MaxAttempts         uint64
	ShowLockOutFailures bool
}
