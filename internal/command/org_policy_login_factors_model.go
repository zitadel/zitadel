package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
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

func (wm *OrgSecondFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
			org.LoginPolicySecondFactorRemovedEventType).
		SearchQueryBuilder()
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

func (wm *OrgMultiFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
			org.LoginPolicyMultiFactorRemovedEventType).
		SearchQueryBuilder()
}

func NewOrgAuthFactorsAllowedWriteModel(orgID string) *OrgAuthFactorsAllowedWriteModel {
	return &OrgAuthFactorsAllowedWriteModel{
		WriteModel: eventstore.WriteModel{
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
		case *iam.LoginPolicySecondFactorAddedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].IAM = domain.FactorStateActive
		case *iam.LoginPolicySecondFactorRemovedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].IAM = domain.FactorStateRemoved
		case *org.LoginPolicySecondFactorAddedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].Org = domain.FactorStateActive
		case *org.LoginPolicySecondFactorRemovedEvent:
			wm.ensureSecondFactor(e.MFAType)
			wm.SecondFactors[e.MFAType].Org = domain.FactorStateRemoved
		case *iam.LoginPolicyMultiFactorAddedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].IAM = domain.FactorStateActive
		case *iam.LoginPolicyMultiFactorRemovedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].IAM = domain.FactorStateRemoved
		case *org.LoginPolicyMultiFactorAddedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].Org = domain.FactorStateActive
		case *org.LoginPolicyMultiFactorRemovedEvent:
			wm.ensureMultiFactor(e.MFAType)
			wm.MultiFactors[e.MFAType].Org = domain.FactorStateRemoved
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
		AggregateTypes(iam.AggregateType).
		AggregateIDs(domain.IAMID).
		EventTypes(
			iam.LoginPolicySecondFactorAddedEventType,
			iam.LoginPolicySecondFactorRemovedEventType,
			iam.LoginPolicyMultiFactorAddedEventType,
			iam.LoginPolicyMultiFactorRemovedEventType,
		).
		Or().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		EventTypes(
			org.LoginPolicySecondFactorAddedEventType,
			org.LoginPolicySecondFactorRemovedEventType,
			org.LoginPolicyMultiFactorAddedEventType,
			org.LoginPolicyMultiFactorRemovedEventType,
		).
		SearchQueryBuilder()
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
