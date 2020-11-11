package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type ReadModel struct {
	eventstore.ReadModel

	SetUpStarted Step
	SetUpDone    Step

	Members MembersReadModel

	GlobalOrgID string
	ProjectID   string

	DefaultLoginPolicy              policy.LoginPolicyReadModel
	DefaultLabelPolicy              policy.LabelPolicyReadModel
	DefaultOrgIAMPolicy             policy.OrgIAMPolicyReadModel
	DefaultPasswordComplexityPolicy policy.PasswordComplexityPolicyReadModel
	DefaultPasswordAgePolicy        policy.PasswordAgePolicyReadModel
	DefaultPasswordLockoutPolicy    policy.PasswordLockoutPolicyReadModel
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) (err error) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *member.AddedEvent, *member.ChangedEvent, *member.RemovedEvent:
			rm.Members.AppendEvents(event)
		case *policy.LabelPolicyAddedEvent, *policy.LabelPolicyChangedEvent:
			rm.DefaultLabelPolicy.AppendEvents(event)
		case *policy.LoginPolicyAddedEvent, *policy.LoginPolicyChangedEvent:
			rm.DefaultLoginPolicy.AppendEvents(event)
		case *policy.OrgIAMPolicyAddedEvent:
			rm.DefaultOrgIAMPolicy.AppendEvents(event)
		case *policy.PasswordComplexityPolicyAddedEvent, *policy.PasswordComplexityPolicyChangedEvent:
			rm.DefaultPasswordComplexityPolicy.AppendEvents(event)
		case *policy.PasswordAgePolicyAddedEvent, *policy.PasswordAgePolicyChangedEvent:
			rm.DefaultPasswordAgePolicy.AppendEvents(event)
		case *policy.PasswordLockoutPolicyAddedEvent, *policy.PasswordLockoutPolicyChangedEvent:
			rm.DefaultPasswordLockoutPolicy.AppendEvents(event)
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
	for _, reduce := range []func() error{
		rm.Members.Reduce,
		rm.Members.Reduce,
		rm.DefaultLoginPolicy.Reduce,
		rm.DefaultLabelPolicy.Reduce,
		rm.DefaultOrgIAMPolicy.Reduce,
		rm.DefaultPasswordComplexityPolicy.Reduce,
		rm.DefaultPasswordAgePolicy.Reduce,
		rm.DefaultPasswordLockoutPolicy.Reduce,
		rm.ReadModel.Reduce,
	} {
		if err = reduce(); err != nil {
			return err
		}
	}

	return nil
}
