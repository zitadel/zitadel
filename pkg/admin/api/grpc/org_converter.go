package grpc

import (
	"github.com/caos/logging"
	admin_model "github.com/caos/zitadel/internal/admin/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
)

func setUpRequestToModel(setUp *OrgSetUpRequest) *admin_model.SetupOrg {
	return &admin_model.SetupOrg{
		Org:  orgCreateRequestToModel(setUp.Org),
		User: userCreateRequestToModel(setUp.User),
	}
}

func orgCreateRequestToModel(org *CreateOrgRequest) *org_model.Org {
	return &org_model.Org{
		Domain: org.Domain,
		Name:   org.Name,
	}
}

func userCreateRequestToModel(user *CreateUserRequest) *usr_model.User {
	preferredLanguage, err := language.Parse(user.PreferredLanguage)
	logging.Log("GRPC-30hwz").OnError(err).Debug("unable to parse language")

	return &usr_model.User{
		Profile: &usr_model.Profile{
			UserName:          user.UserName,
			DisplayName:       user.DisplayName,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			NickName:          user.NickName,
			PreferredLanguage: preferredLanguage,
			Gender:            genderToModel(user.Gender),
		},
		Password: &usr_model.Password{
			SecretString: user.Password,
		},
		Email: &usr_model.Email{
			EmailAddress:    user.Email,
			IsEmailVerified: user.IsEmailVerified,
		},
		Phone: &usr_model.Phone{
			IsPhoneVerified: user.IsPhoneVerified,
			PhoneNumber:     user.Phone,
		},
		Address: &usr_model.Address{
			Country:       user.Country,
			Locality:      user.Locality,
			PostalCode:    user.PostalCode,
			Region:        user.Region,
			StreetAddress: user.StreetAddress,
		},
	}
}

func setUpOrgResponseFromModel(setUp *admin_model.SetupOrg) *OrgSetUpResponse {
	return &OrgSetUpResponse{
		Org:  orgFromModel(setUp.Org),
		User: userFromModel(setUp.User),
	}
}

func orgsFromModel(orgs []*org_model.Org) []*Org {
	result := make([]*Org, len(orgs))
	for i, org := range orgs {
		result[i] = orgFromModel(org)
	}

	return result
}

func orgFromModel(org *org_model.Org) *Org {
	creationDate, err := ptypes.TimestampProto(org.CreationDate)
	logging.Log("GRPC-GTHsZ").OnError(err).Debug("unable to get timestamp from time")

	changeDate, err := ptypes.TimestampProto(org.ChangeDate)
	logging.Log("GRPC-dVnoj").OnError(err).Debug("unable to get timestamp from time")

	return &Org{
		Domain:       org.Domain,
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Id:           org.AggregateID,
		Name:         org.Name,
		State:        orgStateFromModel(org.State),
	}
}

func userFromModel(user *usr_model.User) *User {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-ckoe3d").OnError(err).Debug("unable to parse timestamp")

	converted := &User{
		Id:                user.AggregateID,
		State:             userStateFromModel(user.State),
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          user.Sequence,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage.String(),
		Gender:            genderFromModel(user.Gender),
	}
	if user.Email != nil {
		converted.Email = user.EmailAddress
		converted.IsEmailVerified = user.IsEmailVerified
	}
	if user.Phone != nil {
		converted.Phone = user.PhoneNumber
		converted.IsPhoneVerified = user.IsPhoneVerified
	}
	if user.Address != nil {
		converted.Country = user.Country
		converted.Locality = user.Locality
		converted.PostalCode = user.PostalCode
		converted.Region = user.Region
		converted.StreetAddress = user.StreetAddress
	}
	return converted
}

func orgStateFromModel(state org_model.OrgState) OrgState {
	switch state {
	case org_model.ORGSTATE_ACTIVE:
		return OrgState_ORGSTATE_ACTIVE
	case org_model.ORGSTATE_INACTIVE:
		return OrgState_ORGSTATE_INACTIVE
	default:
		return OrgState_ORGSTATE_UNSPECIFIED
	}
}

func genderFromModel(gender usr_model.Gender) Gender {
	switch gender {
	case usr_model.GENDER_FEMALE:
		return Gender_GENDER_FEMALE
	case usr_model.GENDER_MALE:
		return Gender_GENDER_MALE
	case usr_model.GENDER_DIVERSE:
		return Gender_GENDER_DIVERSE
	default:
		return Gender_GENDER_UNSPECIFIED
	}
}

func genderToModel(gender Gender) usr_model.Gender {
	switch gender {
	case Gender_GENDER_FEMALE:
		return usr_model.GENDER_FEMALE
	case Gender_GENDER_MALE:
		return usr_model.GENDER_MALE
	case Gender_GENDER_DIVERSE:
		return usr_model.GENDER_DIVERSE
	default:
		return usr_model.GENDER_UNDEFINED
	}
}

func userStateFromModel(state usr_model.UserState) UserState {
	switch state {
	case usr_model.USERSTATE_ACTIVE:
		return UserState_USERSTATE_ACTIVE
	case usr_model.USERSTATE_INACTIVE:
		return UserState_USERSTATE_INACTIVE
	case usr_model.USERSTATE_LOCKED:
		return UserState_USERSTATE_LOCKED
	default:
		return UserState_USERSTATE_UNSPECIFIED
	}
}
