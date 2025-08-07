package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgSecondFactorWriteModel struct {
	SecondFactorWriteModel
}

func NewOrgSecondFactorWriteModel(orgID string, factorType domain.SecondFactorType) *OrgSecondFactorWriteModel {
	return &OrgSecondFactorWriteModel{
		SecondFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			MFAType: factorType,
		},
	}
}

func (wm *OrgSecondFactorWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LoginPolicySecondFactorAddedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.SecondFactorAddedEvent)
			}
		case *org.LoginPolicySecondFactorRemovedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.SecondFactorRemovedEvent)
			}
		case *org.LoginPolicyRemovedEvent:
			wm.WriteModel.AppendEvents(&e.LoginPolicyRemovedEvent)
		}
	}
}

func (wm *OrgSecondFactorWriteModel) Reduce() error {
	return wm.SecondFactorWriteModel.Reduce()
}

func (wm *OrgSecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		EventTypes(
			org.LoginPolicySecondFactorAddedEventType,
			org.LoginPolicySecondFactorRemovedEventType,
			org.LoginPolicyRemovedEventType).
		Builder()
}

type OrgMultiFactorWriteModel struct {
	MultiFactorWriteModel
}

func NewOrgMultiFactorWriteModel(orgID string, factorType domain.MultiFactorType) *OrgMultiFactorWriteModel {
	return &OrgMultiFactorWriteModel{
		MultiFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			MFAType: factorType,
		},
	}
}

func (wm *OrgMultiFactorWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LoginPolicyMultiFactorAddedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.MultiFactorAddedEvent)
			}
		case *org.LoginPolicyMultiFactorRemovedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.MultiFactorRemovedEvent)
			}
		case *org.LoginPolicyRemovedEvent:
			wm.WriteModel.AppendEvents(&e.LoginPolicyRemovedEvent)
		}
	}
}

func (wm *OrgMultiFactorWriteModel) Reduce() error {
	return wm.MultiFactorWriteModel.Reduce()
}

func (wm *OrgMultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		EventTypes(
			org.LoginPolicyMultiFactorAddedEventType,
			org.LoginPolicyMultiFactorRemovedEventType,
			org.LoginPolicyRemovedEventType).
		Builder()
}

func NewOrgAuthFactorsAllowedWriteModel(ctx context.Context, orgID string) *OrgAuthFactorsAllowedWriteModel {
	return &OrgAuthFactorsAllowedWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    authz.GetInstance(ctx).InstanceID(),
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
		SecondFactors: map[domain.SecondFactorType]*factorState{},
		MultiFactors:  map[domain.MultiFactorType]*factorState{},
	}
}

type OrgAuthFactorsAllowedWriteModel struct {
	eventstore.WriteModel
	SecondFactors map[domain.SecondFactorType]*factorState
	MultiFactors  map[domain.MultiFactorType]*factorState
}

type factorState struct {
	IAM domain.FactorState
	Org domain.FactorState
}

func (wm *OrgAuthFactorsAllowedWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.LoginPolicySecondFactorAddedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].IAM = domain.FactorStateActive
		case *instance.LoginPolicySecondFactorRemovedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].IAM = domain.FactorStateRemoved
		case *org.LoginPolicySecondFactorAddedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].Org = domain.FactorStateActive
		case *org.LoginPolicySecondFactorRemovedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].Org = domain.FactorStateRemoved
		case *instance.LoginPolicyMultiFactorAddedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].IAM = domain.FactorStateActive
		case *instance.LoginPolicyMultiFactorRemovedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].IAM = domain.FactorStateRemoved
		case *org.LoginPolicyMultiFactorAddedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].Org = domain.FactorStateActive
		case *org.LoginPolicyMultiFactorRemovedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].Org = domain.FactorStateRemoved
		case *org.LoginPolicyRemovedEvent:
			for factorType := range wm.SecondFactors {
				wm.SecondFactors[factorType].Org = domain.FactorStateRemoved
			}
			for factorType := range wm.MultiFactors {
				wm.MultiFactors[factorType].Org = domain.FactorStateRemoved
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgAuthFactorsAllowedWriteModel) ensureSecondFactor(secondFactor domain.SecondFactorType) {
	_, ok := wm.SecondFactors[secondFactor]
	if !ok {
		wm.SecondFactors[secondFactor] = &factorState{}
	}
}

func (wm *OrgAuthFactorsAllowedWriteModel) ensureMultiFactor(multiFactor domain.MultiFactorType) {
	_, ok := wm.MultiFactors[multiFactor]
	if !ok {
		wm.MultiFactors[multiFactor] = &factorState{}
	}
}

func (wm *OrgAuthFactorsAllowedWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.InstanceID).
		EventTypes(
			instance.LoginPolicySecondFactorAddedEventType,
			instance.LoginPolicySecondFactorRemovedEventType,
			instance.LoginPolicyMultiFactorAddedEventType,
			instance.LoginPolicyMultiFactorRemovedEventType,
		).
		Or().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		EventTypes(
			org.LoginPolicySecondFactorAddedEventType,
			org.LoginPolicySecondFactorRemovedEventType,
			org.LoginPolicyMultiFactorAddedEventType,
			org.LoginPolicyMultiFactorRemovedEventType,
			org.LoginPolicyRemovedEventType,
		).
		Builder()
}

func (wm *OrgAuthFactorsAllowedWriteModel) ToSecondFactorWriteModel(factor domain.SecondFactorType) *OrgSecondFactorWriteModel {
	orgSecondFactorWriteModel := NewOrgSecondFactorWriteModel(wm.AggregateID, factor)
	orgSecondFactorWriteModel.ProcessedSequence = wm.ProcessedSequence
	orgSecondFactorWriteModel.State = wm.SecondFactors[factor].Org
	return orgSecondFactorWriteModel
}

func (wm *OrgAuthFactorsAllowedWriteModel) ToMultiFactorWriteModel(factor domain.MultiFactorType) *OrgMultiFactorWriteModel {
	orgMultiFactorWriteModel := NewOrgMultiFactorWriteModel(wm.AggregateID, factor)
	orgMultiFactorWriteModel.ProcessedSequence = wm.ProcessedSequence
	orgMultiFactorWriteModel.State = wm.MultiFactors[factor].Org
	return orgMultiFactorWriteModel
}
