package domain

import "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type NotificationPolicy struct {
	models.ObjectRoot

	State   PolicyState
	Default bool

	PasswordChange bool
}
