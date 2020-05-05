package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
)

func userFromModel(user *usr_model.User) *User {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-ao4OI").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-LLsq4").OnError(err).Debug("unable to parse timestamp")

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

func passwordRequestToModel(r *PasswordRequest) *usr_model.Password {
	return &usr_model.Password{
		ObjectRoot:   models.ObjectRoot{AggregateID: r.Id},
		SecretString: r.Password,
	}
}

func profileFromModel(profile *usr_model.Profile) *UserProfile {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-56t5s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-K58ds").OnError(err).Debug("unable to parse timestamp")

	return &UserProfile{
		Id:                profile.AggregateID,
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          profile.Sequence,
		UserName:          profile.UserName,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		DisplayName:       profile.DisplayName,
		NickName:          profile.NickName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		Gender:            genderFromModel(profile.Gender),
	}
}

func updateProfileToModel(u *UpdateUserProfileRequest) *usr_model.Profile {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-lk73L").OnError(err).Debug("language malformed")

	return &usr_model.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: u.Id},
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
		DisplayName:       u.DisplayName,
		PreferredLanguage: preferredLanguage,
		Gender:            genderToModel(u.Gender),
	}
}

func emailFromModel(email *usr_model.Email) *UserEmail {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-sdoi3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-klJK3").OnError(err).Debug("unable to parse timestamp")

	return &UserEmail{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func updateEmailToModel(e *UpdateUserEmailRequest) *usr_model.Email {
	return &usr_model.Email{
		ObjectRoot:      models.ObjectRoot{AggregateID: e.Id},
		EmailAddress:    e.Email,
		IsEmailVerified: e.IsEmailVerified,
	}
}

func phoneFromModel(phone *usr_model.Phone) *UserPhone {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-kjn5J").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-LKA9S").OnError(err).Debug("unable to parse timestamp")

	return &UserPhone{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func updatePhoneToModel(e *UpdateUserPhoneRequest) *usr_model.Phone {
	return &usr_model.Phone{
		ObjectRoot:      models.ObjectRoot{AggregateID: e.Id},
		PhoneNumber:     e.Phone,
		IsPhoneVerified: e.IsPhoneVerified,
	}
}

func addressFromModel(address *usr_model.Address) *UserAddress {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-65FRs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-aslk4").OnError(err).Debug("unable to parse timestamp")

	return &UserAddress{
		Id:            address.AggregateID,
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Sequence:      address.Sequence,
		Country:       address.Country,
		StreetAddress: address.StreetAddress,
		Region:        address.Region,
		PostalCode:    address.PostalCode,
		Locality:      address.Locality,
	}
}

func updateAddressToModel(address *UpdateUserAddressRequest) *usr_model.Address {
	return &usr_model.Address{
		ObjectRoot:    models.ObjectRoot{AggregateID: address.Id},
		Country:       address.Country,
		StreetAddress: address.StreetAddress,
		Region:        address.Region,
		PostalCode:    address.PostalCode,
		Locality:      address.Locality,
	}
}

func otpFromModel(otp *usr_model.OTP) *MfaOtpResponse {
	return &MfaOtpResponse{
		UserId: otp.AggregateID,
		Url:    otp.Url,
		Secret: otp.SecretString,
		State:  mfaStateFromModel(otp.State),
	}
}

func notifyTypeToModel(state NotificationType) usr_model.NotificationType {
	switch state {
	case NotificationType_NOTIFICATIONTYPE_EMAIL:
		return usr_model.NOTIFICATIONTYPE_EMAIL
	case NotificationType_NOTIFICATIONTYPE_SMS:
		return usr_model.NOTIFICATIONTYPE_SMS
	default:
		return usr_model.NOTIFICATIONTYPE_EMAIL
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

func mfaStateFromModel(state usr_model.MfaState) MFAState {
	switch state {
	case usr_model.MFASTATE_READY:
		return MFAState_MFASTATE_NOT_READY
	case usr_model.MFASTATE_NOTREADY:
		return MFAState_MFASTATE_NOT_READY
	default:
		return MFAState_MFASTATE_UNSPECIFIED
	}
}
