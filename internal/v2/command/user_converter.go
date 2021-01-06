package command

import (
	"github.com/caos/zitadel/internal/v2/domain"
)

func writeModelToUser(wm *UserWriteModel) *domain.User {
	return &domain.User{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		UserName:   wm.UserName,
		State:      wm.UserState,
	}
}

func writeModelToHuman(wm *HumanWriteModel) *domain.Human {
	return &domain.Human{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Profile: &domain.Profile{
			FirstName:         wm.FirstName,
			LastName:          wm.LastName,
			NickName:          wm.NickName,
			DisplayName:       wm.DisplayName,
			PreferredLanguage: wm.PreferredLanguage,
			Gender:            wm.Gender,
		},
		Email: &domain.Email{
			EmailAddress:    wm.Email,
			IsEmailVerified: wm.IsEmailVerified,
		},
		Address: &domain.Address{
			Country:       wm.Country,
			Locality:      wm.Locality,
			PostalCode:    wm.PostalCode,
			Region:        wm.Region,
			StreetAddress: wm.StreetAddress,
		},
	}
}

func writeModelToProfile(wm *HumanProfileWriteModel) *domain.Profile {
	return &domain.Profile{
		ObjectRoot:        writeModelToObjectRoot(wm.WriteModel),
		FirstName:         wm.FirstName,
		LastName:          wm.LastName,
		NickName:          wm.NickName,
		DisplayName:       wm.DisplayName,
		PreferredLanguage: wm.PreferredLanguage,
		Gender:            wm.Gender,
	}
}

func writeModelToEmail(wm *HumanEmailWriteModel) *domain.Email {
	return &domain.Email{
		ObjectRoot:      writeModelToObjectRoot(wm.WriteModel),
		EmailAddress:    wm.Email,
		IsEmailVerified: wm.IsEmailVerified,
	}
}

func writeModelToAddress(wm *HumanAddressWriteModel) *domain.Address {
	return &domain.Address{
		ObjectRoot:    writeModelToObjectRoot(wm.WriteModel),
		Country:       wm.Country,
		Locality:      wm.Locality,
		PostalCode:    wm.PostalCode,
		Region:        wm.Region,
		StreetAddress: wm.StreetAddress,
	}
}

func writeModelToMachine(wm *MachineWriteModel) *domain.Machine {
	return &domain.Machine{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		Name:        wm.Name,
		Description: wm.Description,
	}
}
