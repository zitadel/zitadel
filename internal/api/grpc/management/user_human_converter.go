package management

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"golang.org/x/text/language"
)

func humanFromModel(user *usr_model.User) *management.HumanResponse {
	human := &management.HumanResponse{
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage.String(),
		Gender:            genderFromModel(user.Gender),
	}

	if user.Email != nil {
		human.Email = user.EmailAddress
		human.IsEmailVerified = user.IsEmailVerified
	}
	if user.Phone != nil {
		human.Phone = user.PhoneNumber
		human.IsPhoneVerified = user.IsPhoneVerified
	}
	if user.Address != nil {
		human.Country = user.Country
		human.Locality = user.Locality
		human.PostalCode = user.PostalCode
		human.Region = user.Region
		human.StreetAddress = user.StreetAddress
	}
	return human
}

func humanCreateToModel(u *management.CreateHumanRequest) *usr_model.Human {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-cK5k2").OnError(err).Debug("language malformed")

	human := &usr_model.Human{
		Profile: &usr_model.Profile{
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			NickName:          u.NickName,
			PreferredLanguage: preferredLanguage,
			Gender:            genderToModel(u.Gender),
		},
		Email: &usr_model.Email{
			EmailAddress:    u.Email,
			IsEmailVerified: u.IsEmailVerified,
		},
		Address: &usr_model.Address{
			Country:       u.Country,
			Locality:      u.Locality,
			PostalCode:    u.PostalCode,
			Region:        u.Region,
			StreetAddress: u.StreetAddress,
		},
	}
	if u.Password != "" {
		human.Password = &usr_model.Password{SecretString: u.Password}
	}
	if u.Phone != "" {
		human.Phone = &usr_model.Phone{PhoneNumber: u.Phone, IsPhoneVerified: u.IsPhoneVerified}
	}
	return human
}
