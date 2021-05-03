package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	Default bool

	PrimaryColor   string
	SecondaryColor string
	WarnColor      string

	PrimaryColorDark   string
	SecondaryColorDark string
	WarnColorDark      string

	HideLoginNameSuffix bool
	ErrorMsgPopup       bool
	DisableWatermark    bool
}

func (p *LabelPolicy) IsValid() bool {
	return p.PrimaryColor != "" && p.SecondaryColor != "" && p.WarnColor != ""
}
