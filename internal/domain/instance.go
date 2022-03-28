package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	IAMID = "IAM"
)

type Instance struct {
	models.ObjectRoot

	Name                            string
	GeneratedDomain                 *InstanceDomain
	State                           InstanceState
	GlobalOrgID                     string
	IAMProjectID                    string
	SetUpDone                       Step
	SetUpStarted                    Step
	Members                         []*Member
	IDPs                            []*IDPConfig
	DefaultLoginPolicy              *LoginPolicy
	DefaultLabelPolicy              *LabelPolicy
	DefaultDomainPolicy             *DomainPolicy
	DefaultPasswordComplexityPolicy *PasswordComplexityPolicy
	DefaultPasswordAgePolicy        *PasswordAgePolicy
	DefaultPasswordLockoutPolicy    *LockoutPolicy
}

func (i *IAM) IsValid() bool {
	return i.Name != ""
}

func (i *IAM) AddGeneratedDomain(iamDomain string) {
	i.GeneratedDomain = &InstanceDomain{Domain: NewIAMDomainName(i.Name, iamDomain), Generated: true}
}

type InstanceState int32

const (
	InstanceStateUnspecified InstanceState = iota
	InstanceStateActive
	InstanceStateRemoved

	instanceStateCount
)

func (f InstanceState) Valid() bool {
	return f >= 0 && f < instanceStateCount
}
