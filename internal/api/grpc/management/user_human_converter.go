package management

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
)

func humanFromDomain(user *domain.Human) *management.HumanResponse {
	human := &management.HumanResponse{
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage.String(),
		Gender:            genderFromDomain(user.Gender),
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

func humanFromModel(user *usr_model.Human) *management.HumanResponse {
	human := &management.HumanResponse{
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage.String(),
		//TODO: User Converter
		Gender: management.Gender(user.Gender),
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
		//TODO: User converter
		Gender:          management.Gender(user.Gender),
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		Phone:           user.Phone,
		IsPhoneVerified: user.IsPhoneVerified,
		Country:         user.Country,
		Locality:        user.Locality,
		PostalCode:      user.PostalCode,
		Region:          user.Region,
		StreetAddress:   user.StreetAddress,
		PasswordChanged: passwordChanged,
	}
}

func humanCreateToDomain(u *management.CreateHumanRequest) *domain.Human {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-cK5k2").OnError(err).Debug("language malformed")

	human := &domain.Human{
		Profile: &domain.Profile{
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			NickName:          u.NickName,
			PreferredLanguage: preferredLanguage,
			Gender:            genderToDomain(u.Gender),
		},
		Email: &domain.Email{
			EmailAddress:    u.Email,
			IsEmailVerified: u.IsEmailVerified,
		},
		Address: &domain.Address{
			Country:       u.Country,
			Locality:      u.Locality,
			PostalCode:    u.PostalCode,
			Region:        u.Region,
			StreetAddress: u.StreetAddress,
		},
	}
	if u.Password != "" {
		human.Password = &domain.Password{SecretString: u.Password}
	}
	if u.Phone != "" {
		human.Phone = &domain.Phone{PhoneNumber: u.Phone, IsPhoneVerified: u.IsPhoneVerified}
	}
	return human
}
