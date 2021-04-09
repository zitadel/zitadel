package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type SecondFactorWriteModel struct {
	eventstore.WriteModel
	MFAType domain.SecondFactorType
	State   domain.FactorState
}

func (wm *SecondFactorWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.SecondFactorAddedEvent:
			wm.MFAType = e.MFAType
			wm.State = domain.FactorStateActive
		case *policy.SecondFactorRemovedEvent:
			wm.MFAType = e.MFAType
			wm.State = domain.FactorStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

type MultiFactorWriteModel struct {
	eventstore.WriteModel
	MFAType domain.MultiFactorType
	State   domain.FactorState
}

func (wm *MultiFactorWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.MultiFactorAddedEvent:
			wm.MFAType = e.MFAType
			wm.State = domain.FactorStateActive
		case *policy.MultiFactorRemovedEvent:
			wm.MFAType = e.MFAType
			wm.State = domain.FactorStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

type AuthFactorsWriteModel struct {
	eventstore.WriteModel
	SecondFactors map[domain.SecondFactorType]domain.FactorState
	MultiFactors  map[domain.MultiFactorType]domain.FactorState
}

func (wm *AuthFactorsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.SecondFactorAddedEvent:
			wm.SecondFactors[e.MFAType] = domain.FactorStateActive
		case *policy.SecondFactorRemovedEvent:
			wm.SecondFactors[e.MFAType] = domain.FactorStateRemoved
		case *policy.MultiFactorAddedEvent:
			wm.MultiFactors[e.MFAType] = domain.FactorStateActive
		case *policy.MultiFactorRemovedEvent:
			wm.MultiFactors[e.MFAType] = domain.FactorStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
