package model

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type LabelPolicy struct {
	models.ObjectRoot

	State               PolicyState
	Default             bool
	PrimaryColor        string
	BackgroundColor     string
	FontColor           string
	WarnColor           string
	PrimaryColorDark    string
	BackgroundColorDark string
	FontColorDark       string
	WarnColorDark       string
	HideLoginNameSuffix bool
}

func (p *LabelPolicy) IsValid() bool {
	return p.ObjectRoot.AggregateID != ""
}
