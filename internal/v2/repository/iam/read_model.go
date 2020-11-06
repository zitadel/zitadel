package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

type ReadModel struct {
	eventstore.ReadModel

	SetUpStarted Step
	SetUpDone    Step

	Members member.MembersReadModel

	GlobalOrgID string
	ProjectID   string
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) (err error) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *member.MemberAddedEvent, *member.MemberChangedEvent, *member.MemberRemovedEvent:
			err = rm.Members.AppendEvents(events...)
		}

	}
	return err
}

func (rm *ReadModel) Reduce() (err error) {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *ProjectSetEvent:
			rm.ProjectID = e.ProjectID
		case *GlobalOrgSetEvent:
			rm.GlobalOrgID = e.OrgID
		case *SetupStepEvent:
			if e.Done {
				rm.SetUpDone = e.Step
			} else {
				rm.SetUpStarted = e.Step
			}
		}
	}
	err = rm.Members.Reduce()
	if err != nil {
		return err
	}
	return rm.ReadModel.Reduce()
}
