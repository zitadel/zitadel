package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LockoutPolicy struct {
	models.ObjectRoot

	Default             bool
	MaxPasswordAttempts uint64
	ShowLockOutFailures bool
}
