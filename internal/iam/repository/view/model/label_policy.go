package model

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/domain"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	LabelPolicyKeyAggregateID = "aggregate_id"
	LabelPolicyKeyState       = "label_policy_state"
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

	Sequence uint64 `json:"-" gorm:"column:sequence"`
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

func LabelPolicyViewToModel(policy *LabelPolicyView) *model.LabelPolicyView {
	return &model.LabelPolicyView{
		AggregateID:  policy.AggregateID,
		Sequence:     policy.Sequence,
		CreationDate: policy.CreationDate,
		ChangeDate:   policy.ChangeDate,

		PrimaryColor:    policy.PrimaryColor,
		BackgroundColor: policy.BackgroundColor,
		WarnColor:       policy.WarnColor,
		FontColor:       policy.FontColor,
		LogoURL:         policy.LogoURL,
		IconURL:         policy.IconURL,

		PrimaryColorDark:    policy.PrimaryColorDark,
		BackgroundColorDark: policy.BackgroundColorDark,
		WarnColorDark:       policy.WarnColorDark,
		FontColorDark:       policy.FontColorDark,
		LogoDarkURL:         policy.LogoDarkURL,
		IconDarkURL:         policy.IconDarkURL,

		FontURL: policy.FontURL,

		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		ErrorMsgPopup:       policy.ErrorMsgPopup,
		DisableWatermark:    policy.DisableWatermark,
		Default:             policy.Default,
	}
}

func (i *LabelPolicyView) AppendEvent(event *models.Event) (err error) {
	asset := &AssetView{}
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.LabelPolicyAdded, org_es_model.LabelPolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		i.State = int32(domain.LabelPolicyStatePreview)
		err = i.SetData(event)
	case es_model.LabelPolicyChanged, org_es_model.LabelPolicyChanged:
		err = i.SetData(event)
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyLogoAdded, org_es_model.LabelPolicyLogoAdded:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.LogoURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyLogoRemoved, org_es_model.LabelPolicyLogoRemoved:
		i.LogoURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyIconAdded, org_es_model.LabelPolicyIconAdded:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.IconURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyIconRemoved, org_es_model.LabelPolicyIconRemoved:
		i.IconURL = ""
	case es_model.LabelPolicyLogoDarkAdded, org_es_model.LabelPolicyLogoDarkAdded:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.LogoDarkURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyLogoDarkRemoved, org_es_model.LabelPolicyLogoDarkRemoved:
		i.LogoDarkURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyIconDarkAdded, org_es_model.LabelPolicyIconDarkAdded:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.IconDarkURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyIconDarkRemoved, org_es_model.LabelPolicyIconDarkRemoved:
		i.IconDarkURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyFontAdded, org_es_model.LabelPolicyFontAdded:
		err = asset.SetData(event)
		if err != nil {
			return err
		}
		i.FontURL = asset.AssetURL
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyFontRemoved, org_es_model.LabelPolicyFontRemoved:
		i.FontURL = ""
		i.State = int32(domain.LabelPolicyStatePreview)
	case es_model.LabelPolicyActivated, org_es_model.LabelPolicyActivated:
		i.State = int32(domain.LabelPolicyStateActive)
	}
	return err
}

func (r *LabelPolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *LabelPolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("MODEL-Flp9C").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}

func (r *AssetView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("MODEL-Ms8f2").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
