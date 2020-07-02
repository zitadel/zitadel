package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type OrgIamPolicy struct {
	models.ObjectRoot

	Description           string
	State                 PolicyState
	UserLoginMustBeDomain bool
	Default               bool
}

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateRemoved
)
