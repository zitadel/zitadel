package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	State          PolicyState
	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
