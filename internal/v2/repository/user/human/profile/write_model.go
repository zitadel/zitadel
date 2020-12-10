package profile

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/user/human"
	"golang.org/x/text/language"
)

type WriteModel struct {
	eventstore.WriteModel

	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            human.Gender
}
