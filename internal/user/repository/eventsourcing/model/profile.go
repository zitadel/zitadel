package model

import (
	"encoding/json"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"golang.org/x/text/language"
)

type Profile struct {
	es_models.ObjectRoot

	FirstName         string      `json:"firstName,omitempty"`
	LastName          string      `json:"lastName,omitempty"`
	NickName          string      `json:"nickName,omitempty"`
	DisplayName       string      `json:"displayName,omitempty"`
	PreferredLanguage LanguageTag `json:"preferredLanguage,omitempty"`
	Gender            int32       `json:"gender,omitempty"`
}

func (p *Profile) Changes(changed *Profile) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.FirstName != "" && p.FirstName != changed.FirstName {
		changes["firstName"] = changed.FirstName
	}
	if changed.LastName != "" && p.LastName != changed.LastName {
		changes["lastName"] = changed.LastName
	}
	if changed.NickName != p.NickName {
		changes["nickName"] = changed.NickName
	}
	if changed.DisplayName != "" && p.DisplayName != changed.DisplayName {
		changes["displayName"] = changed.DisplayName
	}
	if language.Tag(changed.PreferredLanguage) != language.Und && changed.PreferredLanguage != p.PreferredLanguage {
		changes["preferredLanguage"] = changed.PreferredLanguage
	}
	if changed.Gender != p.Gender {
		changes["gender"] = changed.Gender
	}
	return changes
}

func ProfileFromModel(profile *model.Profile) *Profile {
	return &Profile{
		ObjectRoot:        profile.ObjectRoot,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: LanguageTag(profile.PreferredLanguage),
		Gender:            int32(profile.Gender),
	}
}

func ProfileToModel(profile *Profile) *model.Profile {
	return &model.Profile{
		ObjectRoot:        profile.ObjectRoot,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: language.Tag(profile.PreferredLanguage),
		Gender:            model.Gender(profile.Gender),
	}
}

type LanguageTag language.Tag

func (t *LanguageTag) UnmarshalJSON(data []byte) error {
	var tag string
	err := json.Unmarshal(data, &tag)
	if err != nil {
		return err
	}
	*t = LanguageTag(language.Make(tag))
	return nil
}

func (t LanguageTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(language.Tag(t))
}

func (t *LanguageTag) MarshalBinary() ([]byte, error) {
	if t == nil {
		return nil, nil
	}

	return []byte(language.Tag(*t).String()), nil
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (t *LanguageTag) UnmarshalBinary(data []byte) error {
	*t = LanguageTag(language.Make(string(data)))
	return nil
}
