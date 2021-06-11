package model

import (
	"golang.org/x/text/language"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Profile struct {
	es_models.ObjectRoot

	FirstName          string
	LastName           string
	NickName           string
	DisplayName        string
	PreferredLanguage  language.Tag
	Gender             Gender
	PreferredLoginName string
	LoginNames         []string
	Avatar             string
}

func (p *Profile) IsValid() bool {
	return p.FirstName != "" && p.LastName != ""
}

func (p *Profile) SetNamesAsDisplayname() {
	p.DisplayName = p.FirstName + " " + p.LastName
}
