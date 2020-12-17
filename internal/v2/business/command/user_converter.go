package command

import (
	"github.com/caos/zitadel/internal/user/model"
)

func writeModelToProfile(wm *HumanProfileWriteModel) *model.Profile {
	return &model.Profile{
		ObjectRoot:        writeModelToObjectRoot(wm.WriteModel),
		FirstName:         wm.FirstName,
		LastName:          wm.LastName,
		NickName:          wm.NickName,
		DisplayName:       wm.DisplayName,
		PreferredLanguage: wm.PreferredLanguage,
		Gender:            model.Gender(wm.Gender),
	}
}

func writeModelToEmail(wm *HumanEmailWriteModel) *model.Email {
	return &model.Email{
		ObjectRoot:      writeModelToObjectRoot(wm.WriteModel),
		EmailAddress:    wm.Email,
		IsEmailVerified: wm.IsEmailVerified,
	}
}
