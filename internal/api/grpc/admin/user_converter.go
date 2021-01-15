package admin

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
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

func userCreateRequestToModel(user *admin.CreateUserRequest) *usr_model.User {
	var human *usr_model.Human
	var machine *usr_model.Machine

	if h := user.GetHuman(); h != nil {
		human = humanCreateToModel(h)
	}
	if m := user.GetMachine(); m != nil {
		machine = machineCreateToModel(m)
	}

	return &usr_model.User{
		UserName: user.UserName,
		Human:    human,
		Machine:  machine,
	}
}

func humanCreateToModel(u *admin.CreateHumanRequest) *usr_model.Human {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-1ouQc").OnError(err).Debug("language malformed")

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

func machineCreateToModel(machine *admin.CreateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func userFromModel(user *usr_model.User) *admin.UserResponse {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-yo0FW").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-jxoQr").OnError(err).Debug("unable to parse timestamp")

	userResp := &admin.UserResponse{
		Id:           user.AggregateID,
		State:        userStateFromModel(user.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     user.Sequence,
		UserName:     user.UserName,
	}

	if user.Machine != nil {
		userResp.User = &admin.UserResponse_Machine{Machine: machineFromModel(user.Machine)}
	}
	if user.Human != nil {
		userResp.User = &admin.UserResponse_Human{Human: humanFromModel(user.Human)}
	}

	return userResp
}

func machineFromModel(account *usr_model.Machine) *admin.MachineResponse {
	return &admin.MachineResponse{
		Name:        account.Name,
		Description: account.Description,
	}
}

func humanFromModel(user *usr_model.Human) *admin.HumanResponse {
	human := &admin.HumanResponse{
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
