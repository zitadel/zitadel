package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	IAMID = "IAM"
)

type IAM struct {
	models.ObjectRoot

	GlobalOrgID                     string
	IAMProjectID                    string
	SetUpDone                       Step
	SetUpStarted                    Step
	Members                         []*Member
	IDPs                            []*IDPConfig
	DefaultLoginPolicy              *LoginPolicy
	DefaultLabelPolicy              *LabelPolicy
	DefaultOrgIAMPolicy             *OrgIAMPolicy
	DefaultPasswordComplexityPolicy *PasswordComplexityPolicy
	DefaultPasswordAgePolicy        *PasswordAgePolicy
	DefaultPasswordLockoutPolicy    *LockoutPolicy
}

type InstanceDomain struct {
	models.ObjectRoot

	Domain string
}

func (i *InstanceDomain) IsValid() bool {
	return i.Domain != ""
}

type InstanceDomainState int32

const (
	InstanceDomainStateUnspecified InstanceDomainState = iota
	InstanceDomainStateActive
	InstanceDomainStateRemoved

	instanceDomainStateCount
)

func (f InstanceDomainState) Valid() bool {
	return f >= 0 && f < instanceDomainStateCount
}
