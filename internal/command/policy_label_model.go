package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type LabelPolicyWriteModel struct {
	eventstore.WriteModel

	PrimaryColor   string
	SecondaryColor string
	WarnColor      string

	PrimaryColorDark   string
	SecondaryColorDark string
	WarnColorDark      string

	HideLoginNameSuffix bool
	ErrorMsgPopup       bool
	DisableWatermark    bool

	State domain.PolicyState
}

func (wm *LabelPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.LabelPolicyAddedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
			wm.WarnColor = e.WarnColor
			wm.PrimaryColorDark = e.PrimaryColorDark
			wm.SecondaryColorDark = e.SecondaryColorDark
			wm.WarnColorDark = e.WarnColorDark
			wm.HideLoginNameSuffix = e.HideLoginNameSuffix
			wm.ErrorMsgPopup = e.ErrorMsgPopup
			wm.DisableWatermark = e.DisableWatermark
			wm.State = domain.PolicyStateActive
		case *policy.LabelPolicyChangedEvent:
			if e.PrimaryColor != nil {
				wm.PrimaryColor = *e.PrimaryColor
			}
			if e.SecondaryColor != nil {
				wm.SecondaryColor = *e.SecondaryColor
			}
			if e.WarnColor != nil {
				wm.WarnColor = *e.WarnColor
			}
			if e.PrimaryColorDark != nil {
				wm.PrimaryColorDark = *e.PrimaryColorDark
			}
			if e.SecondaryColorDark != nil {
				wm.SecondaryColorDark = *e.SecondaryColorDark
			}
			if e.WarnColorDark != nil {
				wm.WarnColorDark = *e.WarnColorDark
			}
			if e.HideLoginNameSuffix != nil {
				wm.HideLoginNameSuffix = *e.HideLoginNameSuffix
			}
			if e.ErrorMsgPopup != nil {
				wm.ErrorMsgPopup = *e.ErrorMsgPopup
			}
			if e.DisableWatermark != nil {
				wm.DisableWatermark = *e.DisableWatermark
			}
		case *policy.LabelPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
