package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type LabelPolicyWriteModel struct {
	eventstore.WriteModel

	PrimaryColor    string
	BackgroundColor string
	WarnColor       string
	FontColor       string
	LogoKey         string
	IconKey         string

	PrimaryColorDark    string
	BackgroundColorDark string
	WarnColorDark       string
	FontColorDark       string
	LogoDarkKey         string
	IconDarkKey         string

	FontKey string

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
			wm.BackgroundColor = e.BackgroundColor
			wm.WarnColor = e.WarnColor
			wm.FontColor = e.FontColor
			wm.PrimaryColorDark = e.PrimaryColorDark
			wm.BackgroundColorDark = e.BackgroundColorDark
			wm.WarnColorDark = e.WarnColorDark
			wm.FontColorDark = e.FontColorDark
			wm.HideLoginNameSuffix = e.HideLoginNameSuffix
			wm.ErrorMsgPopup = e.ErrorMsgPopup
			wm.DisableWatermark = e.DisableWatermark
			wm.State = domain.PolicyStateActive
		case *policy.LabelPolicyChangedEvent:
			if e.PrimaryColor != nil {
				wm.PrimaryColor = *e.PrimaryColor
			}
			if e.BackgroundColor != nil {
				wm.BackgroundColor = *e.BackgroundColor
			}
			if e.WarnColor != nil {
				wm.WarnColor = *e.WarnColor
			}
			if e.FontColor != nil {
				wm.FontColor = *e.FontColor
			}
			if e.PrimaryColorDark != nil {
				wm.PrimaryColorDark = *e.PrimaryColorDark
			}
			if e.BackgroundColorDark != nil {
				wm.BackgroundColorDark = *e.BackgroundColorDark
			}
			if e.WarnColorDark != nil {
				wm.WarnColorDark = *e.WarnColorDark
			}
			if e.FontColorDark != nil {
				wm.FontColorDark = *e.FontColorDark
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
		case *policy.LabelPolicyLogoAddedEvent:
			wm.LogoKey = e.StoreKey
		case *policy.LabelPolicyLogoRemovedEvent:
			wm.LogoKey = ""
		case *policy.LabelPolicyLogoDarkAddedEvent:
			wm.LogoDarkKey = e.StoreKey
		case *policy.LabelPolicyLogoDarkRemovedEvent:
			wm.LogoDarkKey = ""
		case *policy.LabelPolicyIconAddedEvent:
			wm.IconKey = e.StoreKey
		case *policy.LabelPolicyIconRemovedEvent:
			wm.IconKey = ""
		case *policy.LabelPolicyIconDarkAddedEvent:
			wm.IconDarkKey = e.StoreKey
		case *policy.LabelPolicyIconDarkRemovedEvent:
			wm.IconDarkKey = ""
		case *policy.LabelPolicyFontAddedEvent:
			wm.FontKey = e.StoreKey
		case *policy.LabelPolicyFontRemovedEvent:
			wm.FontKey = ""
		case *policy.LabelPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
