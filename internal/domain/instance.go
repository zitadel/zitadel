package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	IAMID = "IAM"
)

type Instance struct {
	models.ObjectRoot

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
