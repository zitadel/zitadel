package model

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	State          PolicyState
	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
