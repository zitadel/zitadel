package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	LabelPolicyKeyAggregateID  = "aggregate_id"
	LabelPolicyKeyState        = "label_policy_state"
	LabelPolicyKeyInstanceID   = "instance_id"
	LabelPolicyKeyOwnerRemoved = "owner_removed"
)

type LabelPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	State        int32     `json:"-" gorm:"column:label_policy_state;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`

	PrimaryColor        string `json:"primaryColor" gorm:"column:primary_color"`
	BackgroundColor     string `json:"backgroundColor" gorm:"column:background_color"`
	WarnColor           string `json:"warnColor" gorm:"column:warn_color"`
	FontColor           string `json:"fontColor" gorm:"column:font_color"`
	PrimaryColorDark    string `json:"primaryColorDark" gorm:"column:primary_color_dark"`
	BackgroundColorDark string `json:"backgroundColorDark" gorm:"column:background_color_dark"`
	WarnColorDark       string `json:"warnColorDark" gorm:"column:warn_color_dark"`
	FontColorDark       string `json:"fontColorDark" gorm:"column:font_color_dark"`
	LogoURL             string `json:"-" gorm:"column:logo_url"`
	IconURL             string `json:"-" gorm:"column:icon_url"`
	LogoDarkURL         string `json:"-" gorm:"column:logo_dark_url"`
	IconDarkURL         string `json:"-" gorm:"column:icon_dark_url"`
	FontURL             string `json:"-" gorm:"column:font_url"`
	HideLoginNameSuffix bool   `json:"hideLoginNameSuffix" gorm:"column:hide_login_name_suffix"`
	ErrorMsgPopup       bool   `json:"errorMsgPopup" gorm:"column:err_msg_popup"`
	DisableWatermark    bool   `json:"disableWatermark" gorm:"column:disable_watermark"`
	Default             bool   `json:"-" gorm:"-"`

	Sequence   uint64 `json:"-" gorm:"column:sequence"`
	InstanceID string `json:"instanceID" gorm:"column:instance_id;primary_key"`
}

type AssetView struct {
	AssetURL string `json:"storeKey"`
}

func (p *LabelPolicyView) ToDomain() *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  p.AggregateID,
			CreationDate: p.CreationDate,
			ChangeDate:   p.ChangeDate,
			Sequence:     p.Sequence,
		},
		Default:         p.Default,
		PrimaryColor:    p.PrimaryColor,
		BackgroundColor: p.BackgroundColor,
		WarnColor:       p.WarnColor,
		FontColor:       p.FontColor,
		LogoURL:         p.LogoURL,
		IconURL:         p.IconURL,

		PrimaryColorDark:    p.PrimaryColorDark,
		BackgroundColorDark: p.BackgroundColorDark,
		WarnColorDark:       p.WarnColorDark,
		FontColorDark:       p.FontColorDark,
		LogoDarkURL:         p.LogoDarkURL,
		IconDarkURL:         p.IconDarkURL,
		Font:                p.FontURL,

		HideLoginNameSuffix: p.HideLoginNameSuffix,
		ErrorMsgPopup:       p.ErrorMsgPopup,
		DisableWatermark:    p.DisableWatermark,
	}
}

func (i *LabelPolicyView) AppendEvent(event eventstore.Event) (err error) {
	asset := &AssetView{}
	i.Sequence = event.Sequence()
	i.ChangeDate = event.CreatedAt()
	switch event.Type() {
	case instance.LabelPolicyAddedEventType,
		org.LabelPolicyAddedEventType:
		i.setRootData(event)
		i.CreationDate = event.CreatedAt()
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
		i.LogoURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyLogoRemovedEventType,
		org.LabelPolicyLogoRemovedEventType:
		i.LogoURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconAddedEventType,
		org.LabelPolicyIconAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.IconURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconRemovedEventType,
		org.LabelPolicyIconRemovedEventType:
		i.IconURL = ""
	case instance.LabelPolicyLogoDarkAddedEventType,
		org.LabelPolicyLogoDarkAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.LogoDarkURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyLogoDarkRemovedEventType,
		org.LabelPolicyLogoDarkRemovedEventType:
		i.LogoDarkURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconDarkAddedEventType,
		org.LabelPolicyIconDarkAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.IconDarkURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyIconDarkRemovedEventType,
		org.LabelPolicyIconDarkRemovedEventType:
		i.IconDarkURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyFontAddedEventType,
		org.LabelPolicyFontAddedEventType:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.FontURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyFontRemovedEventType,
		org.LabelPolicyFontRemovedEventType:
		i.FontURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case instance.LabelPolicyActivatedEventType,
		org.LabelPolicyActivatedEventType:
		i.State = int32(domain.LabelPolicyStateActive)
	case instance.LabelPolicyAssetsRemovedEventType,
		org.LabelPolicyAssetsRemovedEventType:
		i.LogoURL = ""
		i.IconURL = ""
		i.LogoDarkURL = ""
		i.IconDarkURL = ""
		i.FontURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	}
	return err
}

func (r *LabelPolicyView) setRootData(event eventstore.Event) {
	r.AggregateID = event.Aggregate().ID
	r.InstanceID = event.Aggregate().InstanceID
}

func (r *LabelPolicyView) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(r); err != nil {
		logging.Log("MODEL-Flp9C").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}

func (r *AssetView) SetData(event eventstore.Event) error {
	if err := event.Unmarshal(r); err != nil {
		logging.Log("MODEL-Ms8f2").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
