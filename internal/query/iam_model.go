package query

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/policy"
)

type ReadModel struct {
	eventstore.ReadModel

	SetUpStarted domain.Step
	SetUpDone    domain.Step

	Members IAMMembersReadModel
	IDPs    IAMIDPConfigsReadModel

	GlobalOrgID string
	ProjectID   string

	DefaultLoginPolicy              IAMLoginPolicyReadModel
	DefaultLabelPolicy              IAMLabelPolicyReadModel
	DefaultOrgIAMPolicy             IAMOrgIAMPolicyReadModel
	DefaultPasswordComplexityPolicy IAMPasswordComplexityPolicyReadModel
	DefaultPasswordAgePolicy        IAMPasswordAgePolicyReadModel
	DefaultPasswordLockoutPolicy    IAMPasswordLockoutPolicyReadModel
}

func NewReadModel(id string) *ReadModel {
	return &ReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: id,
		},
	}
}

func (rm *ReadModel) IDPByID(idpID string) *IAMIDPConfigReadModel {
	_, config := rm.IDPs.ConfigByID(idpID)
	if config == nil {
		return nil
	}
	return &IAMIDPConfigReadModel{IDPConfigReadModel: *config}
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *member.MemberAddedEvent,
			*member.MemberChangedEvent,
			*member.MemberRemovedEvent:

			rm.Members.AppendEvents(event)
		case *iam.IDPConfigAddedEvent,
			*iam.IDPConfigChangedEvent,
			*iam.IDPConfigDeactivatedEvent,
			*iam.IDPConfigReactivatedEvent,
			*iam.IDPConfigRemovedEvent,
			*iam.IDPOIDCConfigAddedEvent,
			*iam.IDPOIDCConfigChangedEvent:

			rm.IDPs.AppendEvents(event)
		case *policy.LabelPolicyAddedEvent,
			*policy.LabelPolicyChangedEvent:

			rm.DefaultLabelPolicy.AppendEvents(event)
		case *policy.LoginPolicyAddedEvent,
			*policy.LoginPolicyChangedEvent:

			rm.DefaultLoginPolicy.AppendEvents(event)
		case *policy.OrgIAMPolicyAddedEvent:
			rm.DefaultOrgIAMPolicy.AppendEvents(event)
		case *policy.PasswordComplexityPolicyAddedEvent,
			*policy.PasswordComplexityPolicyChangedEvent:

			rm.DefaultPasswordComplexityPolicy.AppendEvents(event)
		case *policy.PasswordAgePolicyAddedEvent,
			*policy.PasswordAgePolicyChangedEvent:

			rm.DefaultPasswordAgePolicy.AppendEvents(event)
		case *policy.LockoutPolicyAddedEvent,
			*policy.LockoutPolicyChangedEvent:

			rm.DefaultPasswordLockoutPolicy.AppendEvents(event)
		}
	}
}

func (rm *ReadModel) Reduce() (err error) {
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
	for _, reduce := range []func() error{
		rm.Members.Reduce,
		rm.IDPs.Reduce,
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

func (rm *ReadModel) AppendAndReduce(events ...eventstore.EventReader) error {
	rm.AppendEvents(events...)
	return rm.Reduce()
}

func (rm *ReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(rm.AggregateID).
		Builder()
}
