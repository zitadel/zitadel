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

func (p *Profile) Validate(allowedLanguages []language.Tag, allowUndefinedLanguage bool) error {
	if p == nil {
		return errors.ThrowInvalidArgument(nil, "PROFILE-GPY3p", "Errors.User.Profile.Empty")
	}
	if p.FirstName == "" {
		return errors.ThrowInvalidArgument(nil, "PROFILE-RF5z2", "Errors.User.Profile.FirstNameEmpty")
	}
	if p.LastName == "" {
		return errors.ThrowInvalidArgument(nil, "PROFILE-DSUkN", "Errors.User.Profile.LastNameEmpty")
	}
	if err := LanguagesAreSupported(p.PreferredLanguage); err != nil {
		return err
	}
	if err := LanguageIsAllowed(allowUndefinedLanguage, allowedLanguages, p.PreferredLanguage); err != nil {
		return err
	}
	return nil
}

func AvatarURL(prefix, resourceOwner, key string) string {
	if prefix == "" || resourceOwner == "" || key == "" {
		return ""
	}
	return prefix + "/" + resourceOwner + "/" + key
}
