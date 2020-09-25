package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type IAM struct {
	es_models.ObjectRoot
	GlobalOrgID                     string
	IAMProjectID                    string
	SetUpDone                       bool
	SetUpStarted                    bool
	Members                         []*IAMMember
	IDPs                            []*IDPConfig
	DefaultLoginPolicy              *LoginPolicy
	DefaultPasswordComplexityPolicy *PasswordComplexityPolicy
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
