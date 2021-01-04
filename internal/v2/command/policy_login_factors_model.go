package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type SecondFactorWriteModel struct {
	eventstore.WriteModel
	MFAType  domain.SecondFactorType
	IsActive bool
}

func (wm *SecondFactorWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.SecondFactorAddedEvent:
			wm.MFAType = e.MFAType
			wm.IsActive = true
		case *policy.SecondFactorRemovedEvent:
			wm.MFAType = e.MFAType
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}

type MultiFactoryWriteModel struct {
	eventstore.WriteModel
	MFAType  domain.MultiFactorType
	IsActive bool
}

func (wm *MultiFactoryWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.MultiFactorAddedEvent:
			wm.MFAType = e.MFAType
			wm.IsActive = true
		case *policy.MultiFactorRemovedEvent:
			wm.MFAType = e.MFAType
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
