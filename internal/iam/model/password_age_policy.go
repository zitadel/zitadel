package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	State          PolicyState
	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
