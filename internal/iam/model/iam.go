package model

import (
	"github.com/caos/zitadel/internal/domain"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
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
	Step9
	Step10
	//StepCount marks the the length of possible steps (StepCount-1 == last possible step)
	StepCount
)

type IAM struct {
	es_models.ObjectRoot
	GlobalOrgID  string
	IAMProjectID string
	SetUpDone    domain.Step
	SetUpStarted domain.Step
	Members      []*IAMMember
}
