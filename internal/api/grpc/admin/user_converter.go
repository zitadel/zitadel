package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"golang.org/x/text/language"
)

func userCreateRequestToDomain(user *admin.CreateUserRequest) (*domain.Human, *domain.Machine) {
	if h := user.GetHuman(); h != nil {
		human := humanCreateToDomain(h)
		human.Username = user.UserName
		return human, nil
	}
	if m := user.GetMachine(); m != nil {
		machine := machineCreateToDomain(m)
		machine.Username = user.UserName
		return nil, machine
	}
	return nil, nil
}

func humanCreateToDomain(u *admin.CreateHumanRequest) *domain.Human {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-1ouQc").OnError(err).Debug("language malformed")

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

func genderToDomain(gender admin.Gender) domain.Gender {
	switch gender {
	case admin.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	case admin.Gender_GENDER_MALE:
		return domain.GenderMale
	case admin.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	default:
		return domain.GenderUnspecified
	}
}

func machineCreateToDomain(machine *admin.CreateMachineRequest) *domain.Machine {
	return &domain.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func externalIDPViewsToDomain(idps []*usr_model.ExternalIDPView) []*domain.ExternalIDP {
	externalIDPs := make([]*domain.ExternalIDP, len(idps))
	for i, idp := range idps {
		externalIDPs[i] = &domain.ExternalIDP{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   idp.UserID,
				ResourceOwner: idp.ResourceOwner,
			},
			IDPConfigID:    idp.IDPConfigID,
			ExternalUserID: idp.ExternalUserID,
			DisplayName:    idp.UserDisplayName,
		}
	}
	return externalIDPs
}
