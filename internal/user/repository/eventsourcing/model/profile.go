package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"golang.org/x/text/language"
)

type Profile struct {
	es_models.ObjectRoot

	UserName          string       `json:"userName,omitempty"`
	FirstName         string       `json:"firstName,omitempty"`
	LastName          string       `json:"lastName,omitempty"`
	NickName          string       `json:"nickName,omitempty"`
	DisplayName       string       `json:"displayName,omitempty"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitempty"`
	Gender            int32        `json:"gender,omitempty"`

	isUserNameUnique bool `json:"-"`
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
	if changed.DisplayName != p.DisplayName {
		changes["displayName"] = changed.DisplayName
	}
	if p.PreferredLanguage != language.Und && changed.PreferredLanguage != p.PreferredLanguage {
		changes["preferredLanguage"] = changed.PreferredLanguage
	}
	if p.Gender > 0 && changed.Gender != p.Gender {
		changes["gender"] = changed.Gender
	}
	return changes
}

func ProfileFromModel(project *model.Profile) *Profile {
	return &Profile{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  project.ObjectRoot.AggregateID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		UserName:          project.UserName,
		FirstName:         project.FirstName,
		LastName:          project.LastName,
		NickName:          project.NickName,
		DisplayName:       project.DisplayName,
		PreferredLanguage: project.PreferredLanguage,
		Gender:            int32(project.Gender),
	}
}

func ProfileToModel(project *Profile) *model.Profile {
	return &model.Profile{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  project.ObjectRoot.AggregateID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		UserName:          project.UserName,
		FirstName:         project.FirstName,
		LastName:          project.LastName,
		NickName:          project.NickName,
		DisplayName:       project.DisplayName,
		PreferredLanguage: project.PreferredLanguage,
		Gender:            model.Gender(project.Gender),
	}
}
