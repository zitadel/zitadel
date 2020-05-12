package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Iam struct {
	es_models.ObjectRoot
	GlobalOrgID  string
	IamProjectID string
	SetUpDone    bool
	SetUpStarted bool
	Members      []*IamMember
}

func (iam *Iam) GetMember(userID string) (int, *IamMember) {
	for i, m := range iam.Members {
		if m.UserID == userID {
			return i, m
		}
	}
	return -1, nil
}
