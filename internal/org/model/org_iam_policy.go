package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	Description           string
	State                 PolicyState
	UserLoginMustBeDomain bool
	Default               bool
	IamDomain             string
}

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateRemoved
)
