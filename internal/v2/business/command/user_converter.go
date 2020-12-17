package command

import (
	"github.com/caos/zitadel/internal/user/model"
)

func writeModelToUser(wm *UserWriteModel) *model.User {
	return &model.User{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		UserName:   wm.UserName,
		State:      model.UserState(wm.UserState),
	}
}

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

func writeModelToAddress(wm *HumanAddressWriteModel) *model.Address {
	return &model.Address{
		ObjectRoot:    writeModelToObjectRoot(wm.WriteModel),
		Country:       wm.Country,
		Locality:      wm.Locality,
		PostalCode:    wm.PostalCode,
		Region:        wm.Region,
		StreetAddress: wm.StreetAddress,
	}
}
