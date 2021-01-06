package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type LabelPolicyWriteModel struct {
	eventstore.WriteModel

	PrimaryColor   string
	SecondaryColor string
	IsActive       bool
}

func (wm *LabelPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.LabelPolicyAddedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
			wm.IsActive = true
		case *policy.LabelPolicyChangedEvent:
			if e.PrimaryColor != nil {
				wm.PrimaryColor = *e.PrimaryColor
			}
			if e.SecondaryColor != nil {
				wm.SecondaryColor = *e.SecondaryColor
			}
		case *policy.LabelPolicyRemovedEvent:
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
