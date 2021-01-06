package domain

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Org struct {
	models.ObjectRoot

	State OrgState
	Name  string

	Domains                  []*OrgDomain
	Members                  []*OrgMember
	OrgIamPolicy             *OrgIAMPolicy
	LoginPolicy              *LoginPolicy
	LabelPolicy              *LabelPolicy
	PasswordComplexityPolicy *PasswordComplexityPolicy
	PasswordAgePolicy        *PasswordAgePolicy
	PasswordLockoutPolicy    *PasswordLockoutPolicy
	IDPs                     []*IDPConfig
}

type OrgState int32

const (
	OrgStateActive OrgState = iota
	OrgStateInactive
)
