package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
