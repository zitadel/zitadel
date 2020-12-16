package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	iam_password_complexity "github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	iam_password_lockout "github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
	"github.com/caos/zitadel/internal/v2/repository/member"
	"github.com/caos/zitadel/internal/v2/repository/policy"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_lockout"
)

type ReadModel struct {
	eventstore.ReadModel

	SetUpStarted iam.Step
	SetUpDone    iam.Step

	Members IAMMembersReadModel
	IDPs    iam.IDPConfigsReadModel

	GlobalOrgID string
	ProjectID   string

	DefaultLoginPolicy              IAMLoginPolicyReadModel
	DefaultLabelPolicy              IAMLabelPolicyReadModel
	DefaultOrgIAMPolicy             IAMOrgIAMPolicyReadModel
	DefaultPasswordComplexityPolicy iam_password_complexity.ReadModel
	DefaultPasswordAgePolicy        IAMPasswordAgePolicyReadModel
	DefaultPasswordLockoutPolicy    iam_password_lockout.ReadModel
}

func NewReadModel(id string) *ReadModel {
	return &ReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: id,
		},
	}
}

func (rm *ReadModel) IDPByID(idpID string) *iam.IDPConfigReadModel {
	_, config := rm.IDPs.ConfigByID(idpID)
	if config == nil {
		return nil
	}
	return &iam.IDPConfigReadModel{ConfigReadModel: *config}
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *member.MemberAddedEvent,
			*member.ChangedEvent,
			*member.RemovedEvent:

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
		case *password_complexity.AddedEvent,
			*password_complexity.ChangedEvent:

			rm.DefaultPasswordComplexityPolicy.AppendEvents(event)
		case *policy.PassowordAgePolicyAddedEvent,
			*policy.PasswordAgePolicyChangedEvent:

			rm.DefaultPasswordAgePolicy.AppendEvents(event)
		case *password_lockout.AddedEvent,
			*password_lockout.ChangedEvent:

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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).AggregateIDs(rm.AggregateID)
}

func AggregateFromReadModel(rm *ReadModel) *iam.Aggregate {
	return &iam.Aggregate{
		Aggregate: *eventstore.NewAggregate(
			rm.AggregateID,
			iam.AggregateType,
			rm.ResourceOwner,
			iam.AggregateVersion,
			rm.ProcessedSequence,
		),
	}
}
