package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
