package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_label "github.com/caos/zitadel/internal/v2/repository/iam/policy/label"
	iam_login "github.com/caos/zitadel/internal/v2/repository/iam/policy/login"
	iam_org_iam "github.com/caos/zitadel/internal/v2/repository/iam/policy/org_iam"
	iam_password_age "github.com/caos/zitadel/internal/v2/repository/iam/policy/password_age"
	iam_password_complexity "github.com/caos/zitadel/internal/v2/repository/iam/policy/password_complexity"
	iam_password_lockout "github.com/caos/zitadel/internal/v2/repository/iam/policy/password_lockout"
	"github.com/caos/zitadel/internal/v2/repository/member"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
	"github.com/caos/zitadel/internal/v2/repository/policy/org_iam"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_age"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_lockout"
)

type ReadModel struct {
	eventstore.ReadModel

	SetUpStarted Step
	SetUpDone    Step

	Members MembersReadModel
	IDPs    IDPConfigsReadModel

	GlobalOrgID string
	ProjectID   string

	DefaultLoginPolicy              iam_login.ReadModel
	DefaultLabelPolicy              iam_label.ReadModel
	DefaultOrgIAMPolicy             iam_org_iam.ReadModel
	DefaultPasswordComplexityPolicy iam_password_complexity.ReadModel
	DefaultPasswordAgePolicy        iam_password_age.ReadModel
	DefaultPasswordLockoutPolicy    iam_password_lockout.ReadModel
}

func NewReadModel(id string) *ReadModel {
	return &ReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: id,
		},
	}
}

func (rm *ReadModel) IDPByID(idpID string) *IDPConfigReadModel {
	_, config := rm.IDPs.ConfigByID(idpID)
	if config == nil {
		return nil
	}
	return &IDPConfigReadModel{ConfigReadModel: *config}
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *member.AddedEvent,
			*member.ChangedEvent,
			*member.RemovedEvent:

			rm.Members.AppendEvents(event)
		case *IDPConfigAddedEvent,
			*IDPConfigChangedEvent,
			*IDPConfigDeactivatedEvent,
			*IDPConfigReactivatedEvent,
			*IDPConfigRemovedEvent,
			*IDPOIDCConfigAddedEvent,
			*IDPOIDCConfigChangedEvent:

			rm.IDPs.AppendEvents(event)
		case *label.AddedEvent,
			*label.ChangedEvent:

			rm.DefaultLabelPolicy.AppendEvents(event)
		case *login.AddedEvent,
			*login.ChangedEvent:

			rm.DefaultLoginPolicy.AppendEvents(event)
		case *org_iam.AddedEvent:
			rm.DefaultOrgIAMPolicy.AppendEvents(event)
		case *password_complexity.AddedEvent,
			*password_complexity.ChangedEvent:

			rm.DefaultPasswordComplexityPolicy.AppendEvents(event)
		case *password_age.AddedEvent,
			*password_age.ChangedEvent:

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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).AggregateIDs(rm.AggregateID)
}
