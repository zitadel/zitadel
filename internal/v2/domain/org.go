package domain

import (
	"strings"

	"github.com/caos/zitadel/internal/eventstore/models"
)

type Org struct {
	models.ObjectRoot

	State OrgState
	Name  string

	Domains                  []*OrgDomain
	Members                  []*Member
	OrgIamPolicy             *OrgIAMPolicy
	LoginPolicy              *LoginPolicy
	LabelPolicy              *LabelPolicy
	PasswordComplexityPolicy *PasswordComplexityPolicy
	PasswordAgePolicy        *PasswordAgePolicy
	PasswordLockoutPolicy    *PasswordLockoutPolicy
	IDPs                     []*IDPConfig
}

func (o *Org) IsValid() bool {
	return o.Name != ""
}

func (o *Org) AddIAMDomain(iamDomain string) {
	o.Domains = append(o.Domains, &OrgDomain{Domain: o.nameForDomain(iamDomain), Verified: true, Primary: true})
}

func (o *Org) nameForDomain(iamDomain string) string {
	return strings.ToLower(strings.ReplaceAll(o.Name, " ", "-") + "." + iamDomain)
}

type OrgState int32

const (
	OrgStateUnspecified OrgState = iota
	OrgStateActive
	OrgStateInactive
	OrgStateRemoved
)
