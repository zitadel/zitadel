package label

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type LabelPolicyWriteModel struct {
	eventstore.WriteModel

	PrimaryColor   string
	SecondaryColor string
}

func (wm *LabelPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
		case *LabelPolicyChangedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
		}
	}
	return wm.WriteModel.Reduce()
}
