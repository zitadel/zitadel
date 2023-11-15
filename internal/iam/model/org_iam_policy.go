package model

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore/v1/models"
)

type DomainPolicy struct {
	models.ObjectRoot

	State                 PolicyState
	UserLoginMustBeDomain bool
	Default               bool
}
