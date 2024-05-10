package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	LabelPolicyKeyAggregateID  = "aggregate_id"
	LabelPolicyKeyState        = "label_policy_state"
	LabelPolicyKeyInstanceID   = "instance_id"
	LabelPolicyKeyOwnerRemoved = "owner_removed"

	LabelPolicyKeyCreationDate        = "creation_date"
	LabelPolicyKeyChangeDate          = "change_date"
	LabelPolicyKeyPrimaryColor        = "primary_color"
	LabelPolicyKeyBackgroundColor     = "background_color"
	LabelPolicyKeyWarnColor           = "warn_color"
	LabelPolicyKeyFontColor           = "font_color"
	LabelPolicyKeyPrimaryColorDark    = "primary_color_dark"
	LabelPolicyKeyBackgroundColorDark = "background_color_dark"
	LabelPolicyKeyWarnColorDark       = "warn_color_dark"
	LabelPolicyKeyFontColorDark       = "font_color_dark"
	LabelPolicyKeyLogoURL             = "logo_url"
	LabelPolicyKeyIconURL             = "icon_url"
	LabelPolicyKeyLogoDarkURL         = "logo_dark_url"
	LabelPolicyKeyIconDarkURL         = "icon_dark_url"
	LabelPolicyKeyFontURL             = "font_url"
	LabelPolicyKeyHideLoginNameSuffix = "hide_login_name_suffix"
	LabelPolicyKeyErrorMsgPopup       = "err_msg_popup"
	LabelPolicyKeyDisableWatermark    = "disable_watermark"
	LabelPolicyKeySequence            = "sequence"
)

type LabelPolicyView struct {
	AggregateID string `json:"-" gorm:"column:aggregate_id;primary_key"`
	State       int32  `json:"-" gorm:"column:label_policy_state;primary_key"`
	InstanceID  string `json:"instanceID" gorm:"column:instance_id;primary_key"`

	Default bool `json:"-" gorm:"-"`

	CreationDate        repository.Field[time.Time] `json:"-" gorm:"column:creation_date"`
	ChangeDate          repository.Field[time.Time] `json:"-" gorm:"column:change_date"`
	PrimaryColor        repository.Field[string]    `json:"primaryColor" gorm:"column:primary_color"`
	BackgroundColor     repository.Field[string]    `json:"backgroundColor" gorm:"column:background_color"`
	WarnColor           repository.Field[string]    `json:"warnColor" gorm:"column:warn_color"`
	FontColor           repository.Field[string]    `json:"fontColor" gorm:"column:font_color"`
	PrimaryColorDark    repository.Field[string]    `json:"primaryColorDark" gorm:"column:primary_color_dark"`
	BackgroundColorDark repository.Field[string]    `json:"backgroundColorDark" gorm:"column:background_color_dark"`
	WarnColorDark       repository.Field[string]    `json:"warnColorDark" gorm:"column:warn_color_dark"`
	FontColorDark       repository.Field[string]    `json:"fontColorDark" gorm:"column:font_color_dark"`
	LogoURL             repository.Field[string]    `json:"-" gorm:"column:logo_url"`
	IconURL             repository.Field[string]    `json:"-" gorm:"column:icon_url"`
	LogoDarkURL         repository.Field[string]    `json:"-" gorm:"column:logo_dark_url"`
	IconDarkURL         repository.Field[string]    `json:"-" gorm:"column:icon_dark_url"`
	FontURL             repository.Field[string]    `json:"-" gorm:"column:font_url"`
	HideLoginNameSuffix repository.Field[bool]      `json:"hideLoginNameSuffix" gorm:"column:hide_login_name_suffix"`
	ErrorMsgPopup       repository.Field[bool]      `json:"errorMsgPopup" gorm:"column:err_msg_popup"`
	DisableWatermark    repository.Field[bool]      `json:"disableWatermark" gorm:"column:disable_watermark"`
	Sequence            repository.Field[uint64]    `json:"-" gorm:"column:sequence"`
}

type AssetView struct {
	AssetURL string `json:"storeKey"`
}

func (p *LabelPolicyView) ToDomain() *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  p.AggregateID,
			CreationDate: p.CreationDate.Value(),
			ChangeDate:   p.ChangeDate.Value(),
			Sequence:     p.Sequence.Value(),
		},
		Default:         p.Default,
		PrimaryColor:    p.PrimaryColor.Value(),
		BackgroundColor: p.BackgroundColor.Value(),
		WarnColor:       p.WarnColor.Value(),
		FontColor:       p.FontColor.Value(),
		LogoURL:         p.LogoURL.Value(),
		IconURL:         p.IconURL.Value(),

		PrimaryColorDark:    p.PrimaryColorDark.Value(),
		BackgroundColorDark: p.BackgroundColorDark.Value(),
		WarnColorDark:       p.WarnColorDark.Value(),
		FontColorDark:       p.FontColorDark.Value(),
		LogoDarkURL:         p.LogoDarkURL.Value(),
		IconDarkURL:         p.IconDarkURL.Value(),
		Font:                p.FontURL.Value(),

		HideLoginNameSuffix: p.HideLoginNameSuffix.Value(),
		ErrorMsgPopup:       p.ErrorMsgPopup.Value(),
		DisableWatermark:    p.DisableWatermark.Value(),
	}
}

func (i *LabelPolicyView) AppendEvent(event eventstore.Event) (err error) {
	asset := &AssetView{}
	i.Sequence.Set(event.Sequence())
	i.ChangeDate.Set(event.CreatedAt())
	switch event.Type() {
	case instance.LabelPolicyAddedEventType,
		org.LabelPolicyAddedEventType:
		i.setRootData(event)
		i.CreationDate.Set(event.CreatedAt())
		i.State = int32(domain.LabelPolicyStatePreview)
		err = i.SetData(event)
	case instance.LabelPolicyChangedEventType,
		org.LabelPolicyChangedEventType:
		err = i.SetData(event)
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyLogoAddedEventType,
		org.LabelPolicyLogoAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.LogoURL.Set(asset.AssetURL)
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyLogoRemovedEventType,
		org.LabelPolicyLogoRemovedEventType:
		i.LogoURL.Set("")
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconAddedEventType,
		org.LabelPolicyIconAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.IconURL.Set(asset.AssetURL)
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconRemovedEventType,
		org.LabelPolicyIconRemovedEventType:
		i.IconURL.Set("")
	case instance.LabelPolicyLogoDarkAddedEventType,
		org.LabelPolicyLogoDarkAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.LogoDarkURL.Set(asset.AssetURL)
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyLogoDarkRemovedEventType,
		org.LabelPolicyLogoDarkRemovedEventType:
		i.LogoDarkURL.Set("")
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconDarkAddedEventType,
		org.LabelPolicyIconDarkAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.IconDarkURL.Set(asset.AssetURL)
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconDarkRemovedEventType,
		org.LabelPolicyIconDarkRemovedEventType:
		i.IconDarkURL.Set("")
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyFontAddedEventType,
		org.LabelPolicyFontAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.FontURL.Set(asset.AssetURL)
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyFontRemovedEventType,
		org.LabelPolicyFontRemovedEventType:
		i.FontURL.Set("")
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyActivatedEventType,
		org.LabelPolicyActivatedEventType:
		i.State = int32(domain.LabelPolicyStateActive)
	case instance.LabelPolicyAssetsRemovedEventType,
		org.LabelPolicyAssetsRemovedEventType:
		i.LogoURL.Set("")
		i.IconURL.Set("")
		i.LogoDarkURL.Set("")
		i.IconDarkURL.Set("")
		i.FontURL.Set("")
		i.State = int32(domain.LabelPolicyStatePreview)
	}
	return err
}

func (r *LabelPolicyView) setRootData(event eventstore.Event) {
	r.AggregateID = event.Aggregate().ID
	r.InstanceID = event.Aggregate().InstanceID
}

func (v *LabelPolicyView) PKColumns() []handler.Column {
	return []handler.Column{
		handler.NewCol("aggregate_id", v.AggregateID),
		handler.NewCol("label_policy_state", v.State),
		handler.NewCol("instance_id", v.InstanceID),
	}
}

func (v *LabelPolicyView) PKConditions() []handler.Condition {
	return []handler.Condition{
		handler.NewCond("aggregate_id", v.AggregateID),
		handler.NewCond("label_policy_state", v.State),
		handler.NewCond("instance_id", v.InstanceID),
	}
}

func (r *LabelPolicyView) Changes() []handler.Column {
	changes := make([]handler.Column, 0, 18)

	if r.CreationDate.DidChange() {
		changes = append(changes, handler.NewCol("creation_date", r.CreationDate.Value()))
	}
	if r.ChangeDate.DidChange() {
		changes = append(changes, handler.NewCol("change_date", r.ChangeDate.Value()))
	}
	if r.PrimaryColor.DidChange() {
		changes = append(changes, handler.NewCol("primary_color", r.PrimaryColor.Value()))
	}
	if r.BackgroundColor.DidChange() {
		changes = append(changes, handler.NewCol("background_color", r.BackgroundColor.Value()))
	}
	if r.WarnColor.DidChange() {
		changes = append(changes, handler.NewCol("warn_color", r.WarnColor.Value()))
	}
	if r.FontColor.DidChange() {
		changes = append(changes, handler.NewCol("font_color", r.FontColor.Value()))
	}
	if r.PrimaryColorDark.DidChange() {
		changes = append(changes, handler.NewCol("primary_color_dark", r.PrimaryColorDark.Value()))
	}
	if r.BackgroundColorDark.DidChange() {
		changes = append(changes, handler.NewCol("background_color_dark", r.BackgroundColorDark.Value()))
	}
	if r.WarnColorDark.DidChange() {
		changes = append(changes, handler.NewCol("warn_color_dark", r.WarnColorDark.Value()))
	}
	if r.FontColorDark.DidChange() {
		changes = append(changes, handler.NewCol("font_color_dark", r.FontColorDark.Value()))
	}
	if r.LogoURL.DidChange() {
		changes = append(changes, handler.NewCol("logo_url", r.LogoURL.Value()))
	}
	if r.IconURL.DidChange() {
		changes = append(changes, handler.NewCol("icon_url", r.IconURL.Value()))
	}
	if r.LogoDarkURL.DidChange() {
		changes = append(changes, handler.NewCol("logo_dark_url", r.LogoDarkURL.Value()))
	}
	if r.IconDarkURL.DidChange() {
		changes = append(changes, handler.NewCol("icon_dark_url", r.IconDarkURL.Value()))
	}
	if r.FontURL.DidChange() {
		changes = append(changes, handler.NewCol("font_url", r.FontURL.Value()))
	}
	if r.HideLoginNameSuffix.DidChange() {
		changes = append(changes, handler.NewCol("hide_login_name_suffix", r.HideLoginNameSuffix.Value()))
	}
	if r.ErrorMsgPopup.DidChange() {
		changes = append(changes, handler.NewCol("err_msg_popup", r.ErrorMsgPopup.Value()))
	}
	if r.DisableWatermark.DidChange() {
		changes = append(changes, handler.NewCol("disable_watermark", r.DisableWatermark.Value()))
	}
	if r.Sequence.DidChange() {
		changes = append(changes, handler.NewCol("sequence", r.Sequence.Value()))
	}

	return changes
}

func (r *LabelPolicyView) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(r); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}

func (r *AssetView) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(r); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
