package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	State                 PolicyState
	UserLoginMustBeDomain bool
	IamDomain             string
}
