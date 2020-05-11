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
}
