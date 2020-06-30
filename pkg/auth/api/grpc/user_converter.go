package grpc

import (
	"context"
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func userViewFromModel(user *usr_model.UserView) *UserView {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-sd32g").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-FJKq1").OnError(err).Debug("unable to parse timestamp")

	lastLogin, err := ptypes.TimestampProto(user.LastLogin)
	logging.Log("GRPC-Gteh2").OnError(err).Debug("unable to parse timestamp")

	passwordChanged, err := ptypes.TimestampProto(user.PasswordChanged)
	logging.Log("GRPC-fgQFT").OnError(err).Debug("unable to parse timestamp")

	return &UserView{
		Id:                 user.ID,
		State:              userStateFromModel(user.State),
		CreationDate:       creationDate,
		ChangeDate:         changeDate,
		LastLogin:          lastLogin,
		PasswordChanged:    passwordChanged,
		UserName:           user.UserName,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		DisplayName:        user.DisplayName,
		NickName:           user.NickName,
		PreferredLanguage:  user.PreferredLanguage,
		Gender:             genderFromModel(user.Gender),
		Email:              user.Email,
		IsEmailVerified:    user.IsEmailVerified,
		Phone:              user.Phone,
		IsPhoneVerified:    user.IsPhoneVerified,
		Country:            user.Country,
		Locality:           user.Locality,
		PostalCode:         user.PostalCode,
		Region:             user.Region,
		StreetAddress:      user.StreetAddress,
		Sequence:           user.Sequence,
		ResourceOwner:      user.ResourceOwner,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
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

func profileViewFromModel(profile *usr_model.Profile) *UserProfileView {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-s9iKs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-9sujE").OnError(err).Debug("unable to parse timestamp")

	return &UserProfileView{
		Id:                 profile.AggregateID,
		CreationDate:       creationDate,
		ChangeDate:         changeDate,
		Sequence:           profile.Sequence,
		UserName:           profile.UserName,
		FirstName:          profile.FirstName,
		LastName:           profile.LastName,
		DisplayName:        profile.DisplayName,
		NickName:           profile.NickName,
		PreferredLanguage:  profile.PreferredLanguage.String(),
		Gender:             genderFromModel(profile.Gender),
		LoginNames:         profile.LoginNames,
		PreferredLoginName: profile.PreferredLoginName,
	}
}

func updateProfileToModel(ctx context.Context, u *UpdateUserProfileRequest) *usr_model.Profile {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-lk73L").OnError(err).Debug("language malformed")

	return &usr_model.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: auth.GetCtxData(ctx).UserID},
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
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

func emailViewFromModel(email *usr_model.Email) *UserEmailView {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-LSp8s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-6szJe").OnError(err).Debug("unable to parse timestamp")

	return &UserEmailView{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func updateEmailToModel(ctx context.Context, e *UpdateUserEmailRequest) *usr_model.Email {
	return &usr_model.Email{
		ObjectRoot:   models.ObjectRoot{AggregateID: auth.GetCtxData(ctx).UserID},
		EmailAddress: e.Email,
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

func phoneViewFromModel(phone *usr_model.Phone) *UserPhoneView {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-s5zJS").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-s9kLe").OnError(err).Debug("unable to parse timestamp")

	return &UserPhoneView{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func updatePhoneToModel(ctx context.Context, e *UpdateUserPhoneRequest) *usr_model.Phone {
	return &usr_model.Phone{
		ObjectRoot:  models.ObjectRoot{AggregateID: auth.GetCtxData(ctx).UserID},
		PhoneNumber: e.Phone,
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

func addressViewFromModel(address *usr_model.Address) *UserAddressView {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-sk4fS").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-9siEs").OnError(err).Debug("unable to parse timestamp")

	return &UserAddressView{
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

func updateAddressToModel(ctx context.Context, address *UpdateUserAddressRequest) *usr_model.Address {
	return &usr_model.Address{
		ObjectRoot:    models.ObjectRoot{AggregateID: auth.GetCtxData(ctx).UserID},
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

func userStateFromModel(state usr_model.UserState) UserState {
	switch state {
	case usr_model.UserStateActive:
		return UserState_USERSTATE_ACTIVE
	case usr_model.UserStateInactive:
		return UserState_USERSTATE_INACTIVE
	case usr_model.UserStateLocked:
		return UserState_USERSTATE_LOCKED
	case usr_model.UserStateInitial:
		return UserState_USERSTATE_INITIAL
	case usr_model.UserStateSuspend:
		return UserState_USERSTATE_SUSPEND
	default:
		return UserState_USERSTATE_UNSPECIFIED
	}
}

func genderFromModel(gender usr_model.Gender) Gender {
	switch gender {
	case usr_model.GenderFemale:
		return Gender_GENDER_FEMALE
	case usr_model.GenderMale:
		return Gender_GENDER_MALE
	case usr_model.GenderDiverse:
		return Gender_GENDER_DIVERSE
	default:
		return Gender_GENDER_UNSPECIFIED
	}
}

func genderToModel(gender Gender) usr_model.Gender {
	switch gender {
	case Gender_GENDER_FEMALE:
		return usr_model.GenderFemale
	case Gender_GENDER_MALE:
		return usr_model.GenderMale
	case Gender_GENDER_DIVERSE:
		return usr_model.GenderDiverse
	default:
		return usr_model.GenderUnspecified
	}
}

func mfaStateFromModel(state usr_model.MfaState) MFAState {
	switch state {
	case usr_model.MfaStateReady:
		return MFAState_MFASTATE_READY
	case usr_model.MfaStateNotReady:
		return MFAState_MFASTATE_NOT_READY
	default:
		return MFAState_MFASTATE_UNSPECIFIED
	}
}

func mfasFromModel(mfas []*usr_model.MultiFactor) []*MultiFactor {
	converted := make([]*MultiFactor, len(mfas))
	for i, mfa := range mfas {
		converted[i] = mfaFromModel(mfa)
	}
	return converted
}

func mfaFromModel(mfa *usr_model.MultiFactor) *MultiFactor {
	return &MultiFactor{
		State: mfaStateFromModel(mfa.State),
		Type:  mfaTypeFromModel(mfa.Type),
	}
}

func mfaTypeFromModel(mfatype usr_model.MfaType) MfaType {
	switch mfatype {
	case usr_model.MfaTypeOTP:
		return MfaType_MFATYPE_OTP
	case usr_model.MfaTypeSMS:
		return MfaType_MFATYPE_SMS
	default:
		return MfaType_MFATYPE_UNSPECIFIED
	}
}

func userChangesToResponse(response *usr_model.UserChanges, offset uint64, limit uint64) (_ *Changes) {
	return &Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: userChangesToMgtAPI(response),
	}
}

func userChangesToMgtAPI(changes *usr_model.UserChanges) (_ []*Change) {
	result := make([]*Change, len(changes.Changes))

	for i, change := range changes.Changes {
		var data *structpb.Struct
		changedData, err := json.Marshal(change.Data)
		if err == nil {
			data = new(structpb.Struct)
			err = protojson.Unmarshal(changedData, data)
			logging.Log("GRPC-a7F54").OnError(err).Debug("unable to marshal changed data to struct")
		}

		result[i] = &Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Data:       data,
			EditorId:   change.ModifierId,
			Editor:     change.ModifierName,
		}
	}

	return result
}
