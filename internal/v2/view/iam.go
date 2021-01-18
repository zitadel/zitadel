package view

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAM struct {
	eventstore.ReadModel

	SetUpStarted domain.Step
	SetUpDone    domain.Step

	GlobalOrgID string
	ProjectID   string

	// TODO: how to implement queries?
}

func (rm *IAM) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
}

//Reduce implements eventstore.IAMMemberReadModel
//
func (rm *IAM) Reduce() (err error) {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *iam.ProjectSetEvent:
			rm.ProjectID = e.ProjectID
		case *iam.GlobalOrgSetEvent:
			rm.GlobalOrgID = e.OrgID
		case *iam.SetupStepEvent:
			if e.Done {
				rm.SetUpDone = e.Step
			} else {
				rm.SetUpStarted = e.Step
			}
		}
	}
	return rm.ReadModel.Reduce()
	//execute all queries
}
