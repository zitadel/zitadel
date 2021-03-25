package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type LabelPolicyWriteModel struct {
	eventstore.WriteModel

	PrimaryColor        string
	SecondaryColor      string
	HideLoginNameSuffix bool

	State domain.PolicyState
}

func (wm *LabelPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.LabelPolicyAddedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
			wm.HideLoginNameSuffix = e.HideLoginNameSuffix
			wm.State = domain.PolicyStateActive
		case *policy.LabelPolicyChangedEvent:
			if e.PrimaryColor != nil {
				wm.PrimaryColor = *e.PrimaryColor
			}
			if e.SecondaryColor != nil {
				wm.SecondaryColor = *e.SecondaryColor
			}
			if e.HideLoginNameSuffix != nil {
				wm.HideLoginNameSuffix = *e.HideLoginNameSuffix
			}
		case *policy.LabelPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
