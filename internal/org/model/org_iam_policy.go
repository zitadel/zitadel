package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgIamPolicy struct {
	models.ObjectRoot

	Description           string
	State                 PolicyState
	UserLoginMustBeDomain bool
}

type PolicyState int32

const (
	POLICYSTATE_ACTIVE PolicyState = iota
	POLICYSTATE_REMOVED
)
