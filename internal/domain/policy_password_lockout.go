package domain

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type LockoutPolicy struct {
	models.ObjectRoot

	Default             bool
	MaxPasswordAttempts uint64
	ShowLockOutFailures bool
}
