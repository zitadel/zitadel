package grpc

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

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

func userCreateToModel(u *CreateUserRequest) *usr_model.User {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-cK5k2").OnError(err).Debug("language malformed")

	user := &usr_model.User{
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
		user.Password = &usr_model.Password{SecretString: u.Password}
	}
	if u.Phone != "" {
		user.Phone = &usr_model.Phone{PhoneNumber: u.Phone, IsPhoneVerified: u.IsPhoneVerified}
	}
	return user
}

func passwordRequestToModel(r *PasswordRequest) *usr_model.Password {
	return &usr_model.Password{
		ObjectRoot:   models.ObjectRoot{AggregateID: r.Id},
		SecretString: r.Password,
	}
}

func userSearchRequestsToModel(project *UserSearchRequest) *usr_model.UserSearchRequest {
	return &usr_model.UserSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: userSearchQueriesToModel(project.Queries),
	}
}

func userSearchQueriesToModel(queries []*UserSearchQuery) []*usr_model.UserSearchQuery {
	converted := make([]*usr_model.UserSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userSearchQueryToModel(q)
	}
	return converted
}

func userSearchQueryToModel(query *UserSearchQuery) *usr_model.UserSearchQuery {
	return &usr_model.UserSearchQuery{
		Key:    userSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userSearchKeyToModel(key UserSearchKey) usr_model.UserSearchKey {
	switch key {
	case UserSearchKey_USERSEARCHKEY_USER_NAME:
		return usr_model.USERSEARCHKEY_USER_NAME
	case UserSearchKey_USERSEARCHKEY_FIRST_NAME:
		return usr_model.USERSEARCHKEY_FIRST_NAME
	case UserSearchKey_USERSEARCHKEY_LAST_NAME:
		return usr_model.USERSEARCHKEY_LAST_NAME
	case UserSearchKey_USERSEARCHKEY_NICK_NAME:
		return usr_model.USERSEARCHKEY_NICK_NAME
	case UserSearchKey_USERSEARCHKEY_DISPLAY_NAME:
		return usr_model.USERSEARCHKEY_DISPLAY_NAME
	case UserSearchKey_USERSEARCHKEY_EMAIL:
		return usr_model.USERSEARCHKEY_EMAIL
	case UserSearchKey_USERSEARCHKEY_STATE:
		return usr_model.USERSEARCHKEY_STATE
	default:
		return usr_model.USERSEARCHKEY_UNSPECIFIED
	}
}

func profileFromModel(profile *usr_model.Profile) *UserProfile {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-dkso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-ski8d").OnError(err).Debug("unable to parse timestamp")

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
	logging.Log("GRPC-d8k2s").OnError(err).Debug("language malformed")

	return &usr_model.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: u.Id},
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
		PreferredLanguage: preferredLanguage,
		Gender:            genderToModel(u.Gender),
	}
}

func emailFromModel(email *usr_model.Email) *UserEmail {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-d9ow2").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-s0dkw").OnError(err).Debug("unable to parse timestamp")

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
	logging.Log("GRPC-ps9ws").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-09ewq").OnError(err).Debug("unable to parse timestamp")

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
	logging.Log("GRPC-ud8w7").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-si9ws").OnError(err).Debug("unable to parse timestamp")

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

func userSearchResponseFromModel(response *usr_model.UserSearchResponse) *UserSearchResponse {
	return &UserSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      userViewsFromModel(response.Result),
	}
}

func userViewsFromModel(users []*usr_model.UserView) []*UserView {
	converted := make([]*UserView, len(users))
	for i, user := range users {
		converted[i] = userViewFromModel(user)
	}
	return converted
}

func userViewFromModel(user *usr_model.UserView) *UserView {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-dl9we").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-lpsg5").OnError(err).Debug("unable to parse timestamp")

	lastLogin, err := ptypes.TimestampProto(user.LastLogin)
	logging.Log("GRPC-dksi3").OnError(err).Debug("unable to parse timestamp")

	passwordChanged, err := ptypes.TimestampProto(user.PasswordChanged)
	logging.Log("GRPC-dl9ws").OnError(err).Debug("unable to parse timestamp")

	return &UserView{
		Id:              user.ID,
		State:           userStateFromModel(user.State),
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		LastLogin:       lastLogin,
		PasswordChanged: passwordChanged,
		Sequence:        user.Sequence,
		ResourceOwner:   user.ResourceOwner,
		UserName:        user.UserName,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		NickName:        user.NickName,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		Phone:           user.Phone,
		IsPhoneVerified: user.IsPhoneVerified,
		Country:         user.Country,
		Locality:        user.Locality,
		Region:          user.Region,
		PostalCode:      user.PostalCode,
		StreetAddress:   user.StreetAddress,
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
	case usr_model.USERSTATE_INITIAL:
		return UserState_USERSTATE_INITIAL
	case usr_model.USERSTATE_SUSPEND:
		return UserState_USERSTATE_SUSPEND
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

func mfaTypeFromModel(mfatype usr_model.MFAType) MfaType {
	switch mfatype {
	case usr_model.MFATYPE_OTP:
		return MfaType_MFATYPE_OTP
	case usr_model.MFATYPE_SMS:
		return MfaType_MFATYPE_SMS
	default:
		return MfaType_MFATYPE_UNSPECIFIED
	}
}

func mfaStateFromModel(state usr_model.MfaState) MFAState {
	switch state {
	case usr_model.MFASTATE_READY:
		return MFAState_MFASTATE_READY
	case usr_model.MFASTATE_NOTREADY:
		return MFAState_MFASTATE_NOT_READY
	default:
		return MFAState_MFASTATE_UNSPECIFIED
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
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Data:       data,
		}
	}

	return result
}
