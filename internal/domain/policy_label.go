package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	State   LabelPolicyState
	Default bool

	PrimaryColor   string
	SecondaryColor string
	WarnColor      string
	LogoURL        string
	IconURL        string

	PrimaryColorDark   string
	SecondaryColorDark string
	WarnColorDark      string
	LogoDarkURL        string
	IconDarkURL        string

	HideLoginNameSuffix bool
	ErrorMsgPopup       bool
	DisableWatermark    bool
}

func (p *LabelPolicy) IsValid() bool {
	return p.PrimaryColor != "" && p.SecondaryColor != "" && p.WarnColor != ""
}

type LabelPolicyState int32

const (
	LabelPolicyStateUnspecified LabelPolicyState = iota
	LabelPolicyStateActive
	LabelPolicyStateRemoved
	LabelPolicyStatePreview

	labelPolicyStateCount
)

func (f LabelPolicyState) Valid() bool {
	return f >= 0 && f < labelPolicyStateCount
}

func (s LabelPolicyState) Exists() bool {
	return s != LabelPolicyStateUnspecified && s != LabelPolicyStateRemoved
}
