package multi_factors

import "github.com/caos/zitadel/internal/eventstore/v2"

type MultiFactoryWriteModel struct {
	eventstore.WriteModel
	MFAType MultiFactorType
}

func (wm *MultiFactoryWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *MultiFactorAddedEvent:
			wm.MFAType = e.MFAType
		case *MultiFactorRemovedEvent:
			wm.MFAType = e.MFAType
		}
	}
	return wm.WriteModel.Reduce()
}
