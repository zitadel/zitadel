package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	State                 PolicyState
	UserLoginMustBeDomain bool
	Default               bool
}
