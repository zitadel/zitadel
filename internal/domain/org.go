package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Org struct {
	models.ObjectRoot

	State OrgState
	Name  string

	PrimaryDomain            string
	Domains                  []*OrgDomain
	Members                  []*Member
	OrgIamPolicy             *OrgIAMPolicy
	LoginPolicy              *LoginPolicy
	LabelPolicy              *LabelPolicy
	PasswordComplexityPolicy *PasswordComplexityPolicy
	PasswordAgePolicy        *PasswordAgePolicy
	PasswordLockoutPolicy    *LockoutPolicy
	IDPs                     []*IDPConfig
}

func (o *Org) IsValid() bool {
	return o.Name != ""
}

func (o *Org) AddIAMDomain(iamDomain string) {
	o.Domains = append(o.Domains, &OrgDomain{Domain: NewIAMDomainName(o.Name, iamDomain), Verified: true, Primary: true})
}

type OrgState int32

const (
	OrgStateUnspecified OrgState = iota
	OrgStateActive
	OrgStateInactive
	OrgStateRemoved
)
