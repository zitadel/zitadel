package management

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
)

func humanFromModel(user *usr_model.Human) *management.HumanResponse {
	human := &management.HumanResponse{
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

func humanViewFromModel(user *usr_model.HumanView) *management.HumanView {
	passwordChanged, err := ptypes.TimestampProto(user.PasswordChanged)
	logging.Log("MANAG-h4ByY").OnError(err).Debug("unable to parse date")

	return &management.HumanView{
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage,
		Gender:            genderFromModel(user.Gender),
		Email:             user.Email,
		IsEmailVerified:   user.IsEmailVerified,
		Phone:             user.Phone,
		IsPhoneVerified:   user.IsPhoneVerified,
		Country:           user.Country,
		Locality:          user.Locality,
		PostalCode:        user.PostalCode,
		Region:            user.Region,
		StreetAddress:     user.StreetAddress,
		PasswordChanged:   passwordChanged,
	}
}

func humanCreateToModel(u *management.CreateHumanRequest) *usr_model.Human {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-cK5k2").OnError(err).Debug("language malformed")

	human := &usr_model.Human{
		Profile: &usr_model.Profile{
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
