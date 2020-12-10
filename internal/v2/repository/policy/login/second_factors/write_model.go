package second_factors

import "github.com/caos/zitadel/internal/eventstore/v2"

type SecondFactoryWriteModel struct {
	eventstore.WriteModel
	MFAType SecondFactorType
}

func (wm *SecondFactoryWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *SecondFactorAddedEvent:
			wm.MFAType = e.MFAType
		case *SecondFactorRemovedEvent:
			wm.MFAType = e.MFAType
		}
	}
	return wm.WriteModel.Reduce()
}
