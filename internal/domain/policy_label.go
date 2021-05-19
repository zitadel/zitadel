package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	State   LabelPolicyState
	Default bool

	PrimaryColor    string
	BackgroundColor string
	WarnColor       string
	FontColor       string
	LogoURL         string
	IconURL         string

	PrimaryColorDark    string
	BackgroundColorDark string
	WarnColorDark       string
	FontColorDark       string
	LogoDarkURL         string
	IconDarkURL         string

	Font string

	HideLoginNameSuffix bool
	ErrorMsgPopup       bool
	DisableWatermark    bool
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
