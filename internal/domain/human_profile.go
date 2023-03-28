package domain

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
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

func (p *Profile) Validate() error {
	if p == nil {
		return errors.ThrowInvalidArgument(nil, "PROFILE-GPY3p", "Errors.User.Profile.Empty")
	}
	if p.FirstName == "" {
		return errors.ThrowInvalidArgument(nil, "PROFILE-RF5z2", "Errors.User.Profile.FirstNameEmpty")
	}
	if p.LastName == "" {
		return errors.ThrowInvalidArgument(nil, "PROFILE-DSUkN", "Errors.User.Profile.LastNameEmpty")
	}
	return nil
}

func AvatarURL(prefix, resourceOwner, key string) string {
	if prefix == "" || resourceOwner == "" || key == "" {
		return ""
	}
	return prefix + "/" + resourceOwner + "/" + key
}
