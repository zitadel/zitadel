package grpc

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func userFromModel(user *usr_model.User) *grpc.User {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-ckoe3d").OnError(err).Debug("unable to parse timestamp")

	converted := &grpc.User{
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

func userCreateToModel(u *grpc.CreateUserRequest) *usr_model.User {
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

func passwordRequestToModel(r *grpc.PasswordRequest) *usr_model.Password {
	return &usr_model.Password{
		ObjectRoot:   models.ObjectRoot{AggregateID: r.Id},
		SecretString: r.Password,
	}
}

func userSearchRequestsToModel(project *grpc.UserSearchRequest) *usr_model.UserSearchRequest {
	return &usr_model.UserSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: userSearchQueriesToModel(project.Queries),
	}
}

func userSearchQueriesToModel(queries []*grpc.UserSearchQuery) []*usr_model.UserSearchQuery {
	converted := make([]*usr_model.UserSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userSearchQueryToModel(q)
	}
	return converted
}

func userSearchQueryToModel(query *grpc.UserSearchQuery) *usr_model.UserSearchQuery {
	return &usr_model.UserSearchQuery{
		Key:    userSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userSearchKeyToModel(key grpc.UserSearchKey) usr_model.UserSearchKey {
	switch key {
	case grpc.UserSearchKey_USERSEARCHKEY_USER_NAME:
		return usr_model.UserSearchKeyUserName
	case grpc.UserSearchKey_USERSEARCHKEY_FIRST_NAME:
		return usr_model.UserSearchKeyFirstName
	case grpc.UserSearchKey_USERSEARCHKEY_LAST_NAME:
		return usr_model.UserSearchKeyLastName
	case grpc.UserSearchKey_USERSEARCHKEY_NICK_NAME:
		return usr_model.UserSearchKeyNickName
	case grpc.UserSearchKey_USERSEARCHKEY_DISPLAY_NAME:
		return usr_model.UserSearchKeyDisplayName
	case grpc.UserSearchKey_USERSEARCHKEY_EMAIL:
		return usr_model.UserSearchKeyEmail
	case grpc.UserSearchKey_USERSEARCHKEY_STATE:
		return usr_model.UserSearchKeyState
	default:
		return usr_model.UserSearchKeyUnspecified
	}
}

func profileFromModel(profile *usr_model.Profile) *grpc.UserProfile {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-dkso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-ski8d").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserProfile{
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

func profileViewFromModel(profile *usr_model.Profile) *grpc.UserProfileView {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-sk8sk").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-s30Ks'").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserProfileView{
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

func updateProfileToModel(u *grpc.UpdateUserProfileRequest) *usr_model.Profile {
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

func emailFromModel(email *usr_model.Email) *grpc.UserEmail {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-d9ow2").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-s0dkw").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserEmail{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func emailViewFromModel(email *usr_model.Email) *grpc.UserEmailView {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-sKefs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-0isjD").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserEmailView{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func updateEmailToModel(e *grpc.UpdateUserEmailRequest) *usr_model.Email {
	return &usr_model.Email{
		ObjectRoot:      models.ObjectRoot{AggregateID: e.Id},
		EmailAddress:    e.Email,
		IsEmailVerified: e.IsEmailVerified,
	}
}

func phoneFromModel(phone *usr_model.Phone) *grpc.UserPhone {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-ps9ws").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-09ewq").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserPhone{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func phoneViewFromModel(phone *usr_model.Phone) *grpc.UserPhoneView {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-6gSj").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-lKs8f").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserPhoneView{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}
func updatePhoneToModel(e *grpc.UpdateUserPhoneRequest) *usr_model.Phone {
	return &usr_model.Phone{
		ObjectRoot:      models.ObjectRoot{AggregateID: e.Id},
		PhoneNumber:     e.Phone,
		IsPhoneVerified: e.IsPhoneVerified,
	}
}

func addressFromModel(address *usr_model.Address) *grpc.UserAddress {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-ud8w7").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-si9ws").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserAddress{
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

func addressViewFromModel(address *usr_model.Address) *grpc.UserAddressView {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-67stC").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-0jSfs").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserAddressView{
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

func updateAddressToModel(address *grpc.UpdateUserAddressRequest) *usr_model.Address {
	return &usr_model.Address{
		ObjectRoot:    models.ObjectRoot{AggregateID: address.Id},
		Country:       address.Country,
		StreetAddress: address.StreetAddress,
		Region:        address.Region,
		PostalCode:    address.PostalCode,
		Locality:      address.Locality,
	}
}

func userSearchResponseFromModel(response *usr_model.UserSearchResponse) *grpc.UserSearchResponse {
	return &grpc.UserSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      userViewsFromModel(response.Result),
	}
}

func userViewsFromModel(users []*usr_model.UserView) []*grpc.UserView {
	converted := make([]*grpc.UserView, len(users))
	for i, user := range users {
		converted[i] = userViewFromModel(user)
	}
	return converted
}

func userViewFromModel(user *usr_model.UserView) *grpc.UserView {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-dl9we").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-lpsg5").OnError(err).Debug("unable to parse timestamp")

	lastLogin, err := ptypes.TimestampProto(user.LastLogin)
	logging.Log("GRPC-dksi3").OnError(err).Debug("unable to parse timestamp")

	passwordChanged, err := ptypes.TimestampProto(user.PasswordChanged)
	logging.Log("GRPC-dl9ws").OnError(err).Debug("unable to parse timestamp")

	return &grpc.UserView{
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

func mfasFromModel(mfas []*usr_model.MultiFactor) []*grpc.MultiFactor {
	converted := make([]*grpc.MultiFactor, len(mfas))
	for i, mfa := range mfas {
		converted[i] = mfaFromModel(mfa)
	}
	return converted
}

func mfaFromModel(mfa *usr_model.MultiFactor) *grpc.MultiFactor {
	return &grpc.MultiFactor{
		State: mfaStateFromModel(mfa.State),
		Type:  mfaTypeFromModel(mfa.Type),
	}
}

func notifyTypeToModel(state grpc.NotificationType) usr_model.NotificationType {
	switch state {
	case grpc.NotificationType_NOTIFICATIONTYPE_EMAIL:
		return usr_model.NotificationTypeEmail
	case grpc.NotificationType_NOTIFICATIONTYPE_SMS:
		return usr_model.NotificationTypeSms
	default:
		return usr_model.NotificationTypeEmail
	}
}

func userStateFromModel(state usr_model.UserState) grpc.UserState {
	switch state {
	case usr_model.UserStateActive:
		return grpc.UserState_USERSTATE_ACTIVE
	case usr_model.UserStateInactive:
		return grpc.UserState_USERSTATE_INACTIVE
	case usr_model.UserStateLocked:
		return grpc.UserState_USERSTATE_LOCKED
	case usr_model.UserStateInitial:
		return grpc.UserState_USERSTATE_INITIAL
	case usr_model.UserStateSuspend:
		return grpc.UserState_USERSTATE_SUSPEND
	default:
		return grpc.UserState_USERSTATE_UNSPECIFIED
	}
}

func genderFromModel(gender usr_model.Gender) grpc.Gender {
	switch gender {
	case usr_model.GenderFemale:
		return grpc.Gender_GENDER_FEMALE
	case usr_model.GenderMale:
		return grpc.Gender_GENDER_MALE
	case usr_model.GenderDiverse:
		return grpc.Gender_GENDER_DIVERSE
	default:
		return grpc.Gender_GENDER_UNSPECIFIED
	}
}

func genderToModel(gender grpc.Gender) usr_model.Gender {
	switch gender {
	case grpc.Gender_GENDER_FEMALE:
		return usr_model.GenderFemale
	case grpc.Gender_GENDER_MALE:
		return usr_model.GenderMale
	case grpc.Gender_GENDER_DIVERSE:
		return usr_model.GenderDiverse
	default:
		return usr_model.GenderUnspecified
	}
}

func mfaTypeFromModel(mfatype usr_model.MfaType) grpc.MfaType {
	switch mfatype {
	case usr_model.MfaTypeOTP:
		return grpc.MfaType_MFATYPE_OTP
	case usr_model.MfaTypeSMS:
		return grpc.MfaType_MFATYPE_SMS
	default:
		return grpc.MfaType_MFATYPE_UNSPECIFIED
	}
}

func mfaStateFromModel(state usr_model.MfaState) grpc.MFAState {
	switch state {
	case usr_model.MfaStateReady:
		return grpc.MFAState_MFASTATE_READY
	case usr_model.MfaStateNotReady:
		return grpc.MFAState_MFASTATE_NOT_READY
	default:
		return grpc.MFAState_MFASTATE_UNSPECIFIED
	}
}

func userChangesToResponse(response *usr_model.UserChanges, offset uint64, limit uint64) (_ *grpc.Changes) {
	return &grpc.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: userChangesToMgtAPI(response),
	}
}

func userChangesToMgtAPI(changes *usr_model.UserChanges) (_ []*grpc.Change) {
	result := make([]*grpc.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &grpc.Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Data:       data,
		}
	}

	return result
}
