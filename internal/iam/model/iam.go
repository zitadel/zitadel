package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/v2/business/domain"
)

type Step int

const (
	Step1 Step = iota + 1
	Step2
	Step3
	Step4
	Step5
	Step6
	Step7
	Step8
	//StepCount marks the the length of possible steps (StepCount-1 == last possible step)
	StepCount
)

type IAM struct {
	es_models.ObjectRoot
	GlobalOrgID                     string
	IAMProjectID                    string
	SetUpDone                       domain.Step
	SetUpStarted                    domain.Step
	Members                         []*IAMMember
	IDPs                            []*IDPConfig
	DefaultLoginPolicy              *LoginPolicy
	DefaultLabelPolicy              *LabelPolicy
	DefaultOrgIAMPolicy             *OrgIAMPolicy
	DefaultPasswordComplexityPolicy *PasswordComplexityPolicy
	DefaultPasswordAgePolicy        *PasswordAgePolicy
	DefaultPasswordLockoutPolicy    *PasswordLockoutPolicy
}

func (iam *IAM) GetMember(userID string) (int, *IAMMember) {
	for i, m := range iam.Members {
		if m.UserID == userID {
			return i, m
		}
	}
	return -1, nil
}

func (iam *IAM) GetIDP(idpID string) (int, *IDPConfig) {
	for i, idp := range iam.IDPs {
		if idp.IDPConfigID == idpID {
			return i, idp
		}
	}
	return -1, nil
}
