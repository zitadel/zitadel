package domain

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"golang.org/x/text/language"
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
}

func (p *Profile) IsValid() bool {
	return p.FirstName != "" && p.LastName != ""
}
