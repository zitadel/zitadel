package domain

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
