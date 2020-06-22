package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	Description    string
	State          PolicyState
	MaxAgeDays     uint64
	ExpireWarnDays uint64
}
