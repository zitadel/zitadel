package management

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func userMachineFromDomain(machine *domain.Machine) *management.UserResponse {
	creationDate, err := ptypes.TimestampProto(machine.CreationDate)
	logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(machine.ChangeDate)
	logging.Log("GRPC-ckoe3d").OnError(err).Debug("unable to parse timestamp")

	userResp := &management.UserResponse{
		Id:           machine.AggregateID,
		State:        userStateFromDomain(machine.GetState()),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     machine.Sequence,
		UserName:     machine.GetUsername(),
	}
	userResp.User = &management.UserResponse_Machine{Machine: machineFromDomain(machine)}
	return userResp
}

func userHumanFromDomain(human *domain.Human) *management.UserResponse {
	creationDate, err := ptypes.TimestampProto(human.CreationDate)
	logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(human.ChangeDate)
	logging.Log("GRPC-ckoe3d").OnError(err).Debug("unable to parse timestamp")

	userResp := &management.UserResponse{
		Id:           human.AggregateID,
		State:        userStateFromDomain(human.GetState()),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     human.Sequence,
		UserName:     human.GetUsername(),
	}
	userResp.User = &management.UserResponse_Human{Human: humanFromDomain(human)}
	return userResp
}

func userCreateToDomain(user *management.CreateUserRequest) (*domain.Human, *domain.Machine) {
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

func passwordRequestToModel(r *management.PasswordRequest) *usr_model.Password {
	return &usr_model.Password{
		ObjectRoot:   models.ObjectRoot{AggregateID: r.Id},
		SecretString: r.Password,
	}
}

func externalIDPSearchRequestToModel(request *management.ExternalIDPSearchRequest) *usr_model.ExternalIDPSearchRequest {
	return &usr_model.ExternalIDPSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: []*usr_model.ExternalIDPSearchQuery{{Key: usr_model.ExternalIDPSearchKeyUserID, Method: model.SearchMethodEquals, Value: request.UserId}},
	}
}

func externalIDPRemoveToDomain(ctxData authz.CtxData, idp *management.ExternalIDPRemoveRequest) *domain.ExternalIDP {
	return &domain.ExternalIDP{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   idp.UserId,
			ResourceOwner: ctxData.ResourceOwner,
		},
		IDPConfigID:    idp.IdpConfigId,
		ExternalUserID: idp.ExternalUserId,
	}
}

func externalIDPSearchResponseFromModel(response *usr_model.ExternalIDPSearchResponse) *management.ExternalIDPSearchResponse {
	viewTimestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-3h8is").OnError(err).Debug("unable to parse timestamp")

	return &management.ExternalIDPSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     viewTimestamp,
		Result:            externalIDPViewsFromModel(response.Result),
	}
}

func externalIDPViewsFromModel(externalIDPs []*usr_model.ExternalIDPView) []*management.ExternalIDPView {
	converted := make([]*management.ExternalIDPView, len(externalIDPs))
	for i, externalIDP := range externalIDPs {
		converted[i] = externalIDPViewFromModel(externalIDP)
	}
	return converted
}

func externalIDPViewFromModel(externalIDP *usr_model.ExternalIDPView) *management.ExternalIDPView {
	creationDate, err := ptypes.TimestampProto(externalIDP.CreationDate)
	logging.Log("GRPC-Fdu8s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(externalIDP.ChangeDate)
	logging.Log("GRPC-Was7u").OnError(err).Debug("unable to parse timestamp")

	return &management.ExternalIDPView{
		UserId:                  externalIDP.UserID,
		IdpConfigId:             externalIDP.IDPConfigID,
		ExternalUserId:          externalIDP.ExternalUserID,
		ExternalUserDisplayName: externalIDP.UserDisplayName,
		IdpName:                 externalIDP.IDPName,
		CreationDate:            creationDate,
		ChangeDate:              changeDate,
	}
}

func userSearchRequestsToModel(project *management.UserSearchRequest) *usr_model.UserSearchRequest {
	return &usr_model.UserSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: userSearchQueriesToModel(project.Queries),
	}
}

func userSearchQueriesToModel(queries []*management.UserSearchQuery) []*usr_model.UserSearchQuery {
	converted := make([]*usr_model.UserSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userSearchQueryToModel(q)
	}
	return converted
}

func userSearchQueryToModel(query *management.UserSearchQuery) *usr_model.UserSearchQuery {
	return &usr_model.UserSearchQuery{
		Key:    userSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userSearchKeyToModel(key management.UserSearchKey) usr_model.UserSearchKey {
	switch key {
	case management.UserSearchKey_USERSEARCHKEY_USER_NAME:
		return usr_model.UserSearchKeyUserName
	case management.UserSearchKey_USERSEARCHKEY_FIRST_NAME:
		return usr_model.UserSearchKeyFirstName
	case management.UserSearchKey_USERSEARCHKEY_LAST_NAME:
		return usr_model.UserSearchKeyLastName
	case management.UserSearchKey_USERSEARCHKEY_NICK_NAME:
		return usr_model.UserSearchKeyNickName
	case management.UserSearchKey_USERSEARCHKEY_DISPLAY_NAME:
		return usr_model.UserSearchKeyDisplayName
	case management.UserSearchKey_USERSEARCHKEY_EMAIL:
		return usr_model.UserSearchKeyEmail
	case management.UserSearchKey_USERSEARCHKEY_STATE:
		return usr_model.UserSearchKeyState
	case management.UserSearchKey_USERSEARCHKEY_TYPE:
		return usr_model.UserSearchKeyType
	default:
		return usr_model.UserSearchKeyUnspecified
	}
}

func userMembershipSearchRequestsToModel(request *management.UserMembershipSearchRequest) *usr_model.UserMembershipSearchRequest {
	return &usr_model.UserMembershipSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: userMembershipSearchQueriesToModel(request.Queries),
	}
}

func userMembershipSearchQueriesToModel(queries []*management.UserMembershipSearchQuery) []*usr_model.UserMembershipSearchQuery {
	converted := make([]*usr_model.UserMembershipSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userMembershipSearchQueryToModel(q)
	}
	return converted
}

func userMembershipSearchQueryToModel(query *management.UserMembershipSearchQuery) *usr_model.UserMembershipSearchQuery {
	return &usr_model.UserMembershipSearchQuery{
		Key:    userMembershipSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userMembershipSearchKeyToModel(key management.UserMembershipSearchKey) usr_model.UserMembershipSearchKey {
	switch key {
	case management.UserMembershipSearchKey_USERMEMBERSHIPSEARCHKEY_TYPE:
		return usr_model.UserMembershipSearchKeyMemberType
	case management.UserMembershipSearchKey_USERMEMBERSHIPSEARCHKEY_OBJECT_ID:
		return usr_model.UserMembershipSearchKeyObjectID
	default:
		return usr_model.UserMembershipSearchKeyUnspecified
	}
}

func profileFromDomain(profile *domain.Profile) *management.UserProfile {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-dkso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-ski8d").OnError(err).Debug("unable to parse timestamp")

	return &management.UserProfile{
		Id:                profile.AggregateID,
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          profile.Sequence,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		DisplayName:       profile.DisplayName,
		NickName:          profile.NickName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		Gender:            management.Gender(profile.Gender),
	}
}

func profileViewFromModel(profile *usr_model.Profile) *management.UserProfileView {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-sk8sk").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-s30Ks'").OnError(err).Debug("unable to parse timestamp")

	return &management.UserProfileView{
		Id:                 profile.AggregateID,
		CreationDate:       creationDate,
		ChangeDate:         changeDate,
		Sequence:           profile.Sequence,
		FirstName:          profile.FirstName,
		LastName:           profile.LastName,
		DisplayName:        profile.DisplayName,
		NickName:           profile.NickName,
		PreferredLanguage:  profile.PreferredLanguage.String(),
		Gender:             management.Gender(profile.Gender),
		LoginNames:         profile.LoginNames,
		PreferredLoginName: profile.PreferredLoginName,
	}
}

func updateProfileToDomain(u *management.UpdateUserProfileRequest) *domain.Profile {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-d8k2s").OnError(err).Debug("language malformed")

	return &domain.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: u.Id},
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
		PreferredLanguage: preferredLanguage,
		Gender:            genderToDomain(u.Gender),
	}
}

func emailFromDomain(email *domain.Email) *management.UserEmail {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-d9ow2").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-s0dkw").OnError(err).Debug("unable to parse timestamp")

	return &management.UserEmail{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func emailViewFromModel(email *usr_model.Email) *management.UserEmailView {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-sKefs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-0isjD").OnError(err).Debug("unable to parse timestamp")

	return &management.UserEmailView{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func updateEmailToDomain(e *management.UpdateUserEmailRequest) *domain.Email {
	return &domain.Email{
		ObjectRoot:      models.ObjectRoot{AggregateID: e.Id},
		EmailAddress:    e.Email,
		IsEmailVerified: e.IsEmailVerified,
	}
}

func phoneFromDomain(phone *domain.Phone) *management.UserPhone {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-ps9ws").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-09ewq").OnError(err).Debug("unable to parse timestamp")

	return &management.UserPhone{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func phoneViewFromModel(phone *usr_model.Phone) *management.UserPhoneView {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-6gSj").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-lKs8f").OnError(err).Debug("unable to parse timestamp")

	return &management.UserPhoneView{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}
func updatePhoneToDomain(e *management.UpdateUserPhoneRequest) *domain.Phone {
	return &domain.Phone{
		ObjectRoot:      models.ObjectRoot{AggregateID: e.Id},
		PhoneNumber:     e.Phone,
		IsPhoneVerified: e.IsPhoneVerified,
	}
}

func addressFromDomain(address *domain.Address) *management.UserAddress {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-ud8w7").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-si9ws").OnError(err).Debug("unable to parse timestamp")

	return &management.UserAddress{
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

func addressViewFromModel(address *usr_model.Address) *management.UserAddressView {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-67stC").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-0jSfs").OnError(err).Debug("unable to parse timestamp")

	return &management.UserAddressView{
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

func updateAddressToDomain(ctxData authz.CtxData, address *management.UpdateUserAddressRequest) *domain.Address {
	return &domain.Address{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   address.Id,
			ResourceOwner: ctxData.OrgID,
		},
		Country:       address.Country,
		StreetAddress: address.StreetAddress,
		Region:        address.Region,
		PostalCode:    address.PostalCode,
		Locality:      address.Locality,
	}
}

func userSearchResponseFromModel(response *usr_model.UserSearchResponse) *management.UserSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-aBezr").OnError(err).Debug("unable to parse timestamp")
	return &management.UserSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            userViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func userViewsFromModel(users []*usr_model.UserView) []*management.UserView {
	converted := make([]*management.UserView, len(users))
	for i, user := range users {
		converted[i] = userViewFromModel(user)
	}
	return converted
}

func userViewFromModel(user *usr_model.UserView) *management.UserView {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-dl9we").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-lpsg5").OnError(err).Debug("unable to parse timestamp")

	lastLogin, err := ptypes.TimestampProto(user.LastLogin)
	logging.Log("GRPC-dksi3").OnError(err).Debug("unable to parse timestamp")

	userView := &management.UserView{
		Id:                 user.ID,
		State:              management.UserState(user.State),
		CreationDate:       creationDate,
		ChangeDate:         changeDate,
		LastLogin:          lastLogin,
		Sequence:           user.Sequence,
		ResourceOwner:      user.ResourceOwner,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		UserName:           user.UserName,
	}
	if user.HumanView != nil {
		userView.User = &management.UserView_Human{Human: humanViewFromModel(user.HumanView)}
	}
	if user.MachineView != nil {
		userView.User = &management.UserView_Machine{Machine: machineViewFromModel(user.MachineView)}

	}
	return userView
}

func userMembershipSearchResponseFromModel(response *usr_model.UserMembershipSearchResponse) *management.UserMembershipSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-Hs8jd").OnError(err).Debug("unable to parse timestamp")
	return &management.UserMembershipSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            userMembershipViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func userMembershipViewsFromModel(memberships []*usr_model.UserMembershipView) []*management.UserMembershipView {
	converted := make([]*management.UserMembershipView, len(memberships))
	for i, membership := range memberships {
		converted[i] = userMembershipViewFromModel(membership)
	}
	return converted
}

func userMembershipViewFromModel(membership *usr_model.UserMembershipView) *management.UserMembershipView {
	creationDate, err := ptypes.TimestampProto(membership.CreationDate)
	logging.Log("GRPC-Msnu8").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(membership.ChangeDate)
	logging.Log("GRPC-Slco9").OnError(err).Debug("unable to parse timestamp")

	return &management.UserMembershipView{
		UserId:        membership.UserID,
		AggregateId:   membership.AggregateID,
		ObjectId:      membership.ObjectID,
		MemberType:    memberTypeFromModel(membership.MemberType),
		DisplayName:   membership.DisplayName,
		Roles:         membership.Roles,
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Sequence:      membership.Sequence,
		ResourceOwner: membership.ResourceOwner,
	}
}

func mfasFromModel(mfas []*usr_model.MultiFactor) []*management.UserMultiFactor {
	converted := make([]*management.UserMultiFactor, len(mfas))
	for i, mfa := range mfas {
		converted[i] = mfaFromModel(mfa)
	}
	return converted
}

func mfaFromModel(mfa *usr_model.MultiFactor) *management.UserMultiFactor {
	return &management.UserMultiFactor{
		State:     mfaStateFromModel(mfa.State),
		Type:      mfaTypeFromModel(mfa.Type),
		Attribute: mfa.Attribute,
		Id:        mfa.ID,
	}
}

func notifyTypeToDomain(state management.NotificationType) domain.NotificationType {
	switch state {
	case management.NotificationType_NOTIFICATIONTYPE_EMAIL:
		return domain.NotificationTypeEmail
	case management.NotificationType_NOTIFICATIONTYPE_SMS:
		return domain.NotificationTypeSms
	default:
		return domain.NotificationTypeEmail
	}
}

func userStateFromDomain(state domain.UserState) management.UserState {
	switch state {
	case domain.UserStateActive:
		return management.UserState_USERSTATE_ACTIVE
	case domain.UserStateInactive:
		return management.UserState_USERSTATE_INACTIVE
	case domain.UserStateLocked:
		return management.UserState_USERSTATE_LOCKED
	case domain.UserStateInitial:
		return management.UserState_USERSTATE_INITIAL
	case domain.UserStateSuspend:
		return management.UserState_USERSTATE_SUSPEND
	default:
		return management.UserState_USERSTATE_UNSPECIFIED
	}
}

func genderFromDomain(gender domain.Gender) management.Gender {
	switch gender {
	case domain.GenderFemale:
		return management.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return management.Gender_GENDER_MALE
	case domain.GenderDiverse:
		return management.Gender_GENDER_DIVERSE
	default:
		return management.Gender_GENDER_UNSPECIFIED
	}
}

func memberTypeFromModel(memberType usr_model.MemberType) management.MemberType {
	switch memberType {
	case usr_model.MemberTypeOrganisation:
		return management.MemberType_MEMBERTYPE_ORGANISATION
	case usr_model.MemberTypeProject:
		return management.MemberType_MEMBERTYPE_PROJECT
	case usr_model.MemberTypeProjectGrant:
		return management.MemberType_MEMBERTYPE_PROJECT_GRANT
	default:
		return management.MemberType_MEMBERTYPE_UNSPECIFIED
	}
}

func genderToDomain(gender management.Gender) domain.Gender {
	switch gender {
	case management.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	case management.Gender_GENDER_MALE:
		return domain.GenderMale
	case management.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	default:
		return domain.GenderUnspecified
	}
}

func mfaTypeFromModel(mfatype usr_model.MFAType) management.MfaType {
	switch mfatype {
	case usr_model.MFATypeOTP:
		return management.MfaType_MFATYPE_OTP
	case usr_model.MFATypeU2F:
		return management.MfaType_MFATYPE_U2F
	default:
		return management.MfaType_MFATYPE_UNSPECIFIED
	}
}

func mfaStateFromModel(state usr_model.MFAState) management.MFAState {
	switch state {
	case usr_model.MFAStateReady:
		return management.MFAState_MFASTATE_READY
	case usr_model.MFAStateNotReady:
		return management.MFAState_MFASTATE_NOT_READY
	default:
		return management.MFAState_MFASTATE_UNSPECIFIED
	}
}

func userChangesToResponse(response *usr_model.UserChanges, offset uint64, limit uint64) (_ *management.Changes) {
	return &management.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: userChangesToMgtAPI(response),
	}
}

func userChangesToMgtAPI(changes *usr_model.UserChanges) (_ []*management.Change) {
	result := make([]*management.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		var data *structpb.Struct
		changedData, err := json.Marshal(change.Data)
		if err == nil {
			data = new(structpb.Struct)
			err = protojson.Unmarshal(changedData, data)
			logging.Log("GRPC-a7F54").OnError(err).Debug("unable to marshal changed data to struct")
		}

		result[i] = &management.Change{
			ChangeDate: change.ChangeDate,
			EventType:  message.NewLocalizedEventType(change.EventType),
			Sequence:   change.Sequence,
			Data:       data,
			EditorId:   change.ModifierID,
			Editor:     change.ModifierName,
		}
	}

	return result
}

func webAuthNTokensFromModel(tokens []*usr_model.WebAuthNToken) *management.WebAuthNTokens {
	result := make([]*management.WebAuthNToken, len(tokens))
	for i, token := range tokens {
		result[i] = webAuthNTokenFromModel(token)
	}
	return &management.WebAuthNTokens{Tokens: result}
}

func webAuthNTokenFromModel(token *usr_model.WebAuthNToken) *management.WebAuthNToken {
	return &management.WebAuthNToken{
		Id:    token.WebAuthNTokenID,
		Name:  token.WebAuthNTokenName,
		State: mfaStateFromModel(token.State),
	}
}
