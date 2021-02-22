package auth

import (
	"context"
	"encoding/json"
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func userViewFromModel(user *usr_model.UserView) *auth.UserView {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-sd32g").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-FJKq1").OnError(err).Debug("unable to parse timestamp")

	lastLogin, err := ptypes.TimestampProto(user.LastLogin)
	logging.Log("GRPC-Gteh2").OnError(err).Debug("unable to parse timestamp")

	userView := &auth.UserView{
		Id:                 user.ID,
		State:              userStateFromModel(user.State),
		CreationDate:       creationDate,
		ChangeDate:         changeDate,
		LastLogin:          lastLogin,
		UserName:           user.UserName,
		Sequence:           user.Sequence,
		ResourceOwner:      user.ResourceOwner,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
	}

	if user.HumanView != nil {
		userView.User = &auth.UserView_Human{Human: humanViewFromModel(user.HumanView)}
	}
	if user.MachineView != nil {
		userView.User = &auth.UserView_Machine{Machine: machineViewFromModel(user.MachineView)}

	}

	return userView
}

func profileFromDomain(profile *domain.Profile) *auth.UserProfile {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-56t5s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-K58ds").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserProfile{
		Id:                profile.AggregateID,
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          profile.Sequence,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		DisplayName:       profile.DisplayName,
		NickName:          profile.NickName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		Gender:            genderFromDomain(profile.Gender),
	}
}

func profileViewFromModel(profile *usr_model.Profile) *auth.UserProfileView {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-s9iKs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-9sujE").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserProfileView{
		Id:                profile.AggregateID,
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          profile.Sequence,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		DisplayName:       profile.DisplayName,
		NickName:          profile.NickName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		//TODO: Use converter
		Gender:             auth.Gender(profile.Gender),
		LoginNames:         profile.LoginNames,
		PreferredLoginName: profile.PreferredLoginName,
	}
}

func updateProfileToDomain(ctx context.Context, u *auth.UpdateUserProfileRequest) *domain.Profile {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-lk73L").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("language malformed")

	return &domain.Profile{
		ObjectRoot:        ctxToObjectRoot(ctx),
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
		PreferredLanguage: preferredLanguage,
		Gender:            genderToDomain(u.Gender),
	}
}

func emailFromDomain(email *domain.Email) *auth.UserEmail {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-sdoi3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-klJK3").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserEmail{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func emailViewFromModel(email *usr_model.Email) *auth.UserEmailView {
	creationDate, err := ptypes.TimestampProto(email.CreationDate)
	logging.Log("GRPC-LSp8s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(email.ChangeDate)
	logging.Log("GRPC-6szJe").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserEmailView{
		Id:              email.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        email.Sequence,
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func updateEmailToDomain(ctx context.Context, e *auth.UpdateUserEmailRequest) *domain.Email {
	return &domain.Email{
		ObjectRoot:   ctxToObjectRoot(ctx),
		EmailAddress: e.Email,
	}
}

func phoneFromDomain(phone *domain.Phone) *auth.UserPhone {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-kjn5J").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-LKA9S").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserPhone{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func phoneViewFromModel(phone *usr_model.Phone) *auth.UserPhoneView {
	creationDate, err := ptypes.TimestampProto(phone.CreationDate)
	logging.Log("GRPC-s5zJS").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(phone.ChangeDate)
	logging.Log("GRPC-s9kLe").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserPhoneView{
		Id:              phone.AggregateID,
		CreationDate:    creationDate,
		ChangeDate:      changeDate,
		Sequence:        phone.Sequence,
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func updatePhoneToDomain(ctx context.Context, e *auth.UpdateUserPhoneRequest) *domain.Phone {
	return &domain.Phone{
		ObjectRoot:  ctxToObjectRoot(ctx),
		PhoneNumber: e.Phone,
	}
}

func addressFromDomain(address *domain.Address) *auth.UserAddress {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-65FRs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-aslk4").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserAddress{
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

func addressViewFromModel(address *usr_model.Address) *auth.UserAddressView {
	creationDate, err := ptypes.TimestampProto(address.CreationDate)
	logging.Log("GRPC-sk4fS").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(address.ChangeDate)
	logging.Log("GRPC-9siEs").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserAddressView{
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

func updateAddressToDomain(ctx context.Context, address *auth.UpdateUserAddressRequest) *domain.Address {
	return &domain.Address{
		ObjectRoot:    ctxToObjectRoot(ctx),
		Country:       address.Country,
		StreetAddress: address.StreetAddress,
		Region:        address.Region,
		PostalCode:    address.PostalCode,
		Locality:      address.Locality,
	}
}

func externalIDPSearchRequestToModel(request *auth.ExternalIDPSearchRequest) *usr_model.ExternalIDPSearchRequest {
	return &usr_model.ExternalIDPSearchRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}
}

func externalIDPRemoveToDomain(ctx context.Context, idp *auth.ExternalIDPRemoveRequest) *domain.ExternalIDP {
	return &domain.ExternalIDP{
		ObjectRoot:     ctxToObjectRoot(ctx),
		IDPConfigID:    idp.IdpConfigId,
		ExternalUserID: idp.ExternalUserId,
	}
}

func externalIDPResponseFromModel(idp *usr_model.ExternalIDP) *auth.ExternalIDPResponse {
	return &auth.ExternalIDPResponse{
		IdpConfigId: idp.IDPConfigID,
		UserId:      idp.UserID,
		DisplayName: idp.DisplayName,
	}
}

func externalIDPSearchResponseFromModel(response *usr_model.ExternalIDPSearchResponse) *auth.ExternalIDPSearchResponse {
	viewTimestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-3h8is").OnError(err).Debug("unable to parse timestamp")

	return &auth.ExternalIDPSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     viewTimestamp,
		Result:            externalIDPViewsFromModel(response.Result),
	}
}

func externalIDPViewsFromModel(externalIDPs []*usr_model.ExternalIDPView) []*auth.ExternalIDPView {
	converted := make([]*auth.ExternalIDPView, len(externalIDPs))
	for i, externalIDP := range externalIDPs {
		converted[i] = externalIDPViewFromModel(externalIDP)
	}
	return converted
}

func externalIDPViewFromModel(externalIDP *usr_model.ExternalIDPView) *auth.ExternalIDPView {
	creationDate, err := ptypes.TimestampProto(externalIDP.CreationDate)
	logging.Log("GRPC-Sj8dw").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(externalIDP.ChangeDate)
	logging.Log("GRPC-Nf8ue").OnError(err).Debug("unable to parse timestamp")

	return &auth.ExternalIDPView{
		UserId:                  externalIDP.UserID,
		IdpConfigId:             externalIDP.IDPConfigID,
		ExternalUserId:          externalIDP.ExternalUserID,
		ExternalUserDisplayName: externalIDP.UserDisplayName,
		IdpName:                 externalIDP.IDPName,
		CreationDate:            creationDate,
		ChangeDate:              changeDate,
	}
}

func otpFromDomain(otp *domain.OTP) *auth.MfaOtpResponse {
	return &auth.MfaOtpResponse{
		UserId: otp.AggregateID,
		Url:    otp.Url,
		Secret: otp.SecretString,
		State:  mfaStateFromDomain(otp.State),
	}
}

func userStateFromModel(state usr_model.UserState) auth.UserState {
	switch state {
	case usr_model.UserStateActive:
		return auth.UserState_USERSTATE_ACTIVE
	case usr_model.UserStateInactive:
		return auth.UserState_USERSTATE_INACTIVE
	case usr_model.UserStateLocked:
		return auth.UserState_USERSTATE_LOCKED
	case usr_model.UserStateInitial:
		return auth.UserState_USERSTATE_INITIAL
	case usr_model.UserStateSuspend:
		return auth.UserState_USERSTATE_SUSPEND
	default:
		return auth.UserState_USERSTATE_UNSPECIFIED
	}
}

func genderFromDomain(gender domain.Gender) auth.Gender {
	switch gender {
	case domain.GenderFemale:
		return auth.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return auth.Gender_GENDER_MALE
	case domain.GenderDiverse:
		return auth.Gender_GENDER_DIVERSE
	default:
		return auth.Gender_GENDER_UNSPECIFIED
	}
}

func genderToDomain(gender auth.Gender) domain.Gender {
	switch gender {
	case auth.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	case auth.Gender_GENDER_MALE:
		return domain.GenderMale
	case auth.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	default:
		return domain.GenderUnspecified
	}
}

func mfaStateFromDomain(state domain.MFAState) auth.MFAState {
	switch state {
	case domain.MFAStateReady:
		return auth.MFAState_MFASTATE_READY
	case domain.MFAStateNotReady:
		return auth.MFAState_MFASTATE_NOT_READY
	default:
		return auth.MFAState_MFASTATE_UNSPECIFIED
	}
}

func mfasFromModel(mfas []*usr_model.MultiFactor) []*auth.MultiFactor {
	converted := make([]*auth.MultiFactor, len(mfas))
	for i, mfa := range mfas {
		converted[i] = mfaFromModel(mfa)
	}
	return converted
}

func mfaFromModel(mfa *usr_model.MultiFactor) *auth.MultiFactor {
	return &auth.MultiFactor{
		State:     auth.MFAState(mfa.State),
		Type:      mfaTypeFromModel(mfa.Type),
		Attribute: mfa.Attribute,
		Id:        mfa.ID,
	}
}

func mfaTypeFromModel(mfaType usr_model.MFAType) auth.MfaType {
	switch mfaType {
	case usr_model.MFATypeOTP:
		return auth.MfaType_MFATYPE_OTP
	case usr_model.MFATypeU2F:
		return auth.MfaType_MFATYPE_U2F
	default:
		return auth.MfaType_MFATYPE_UNSPECIFIED
	}
}

func userChangesToResponse(response *usr_model.UserChanges, offset uint64, limit uint64) (_ *auth.Changes) {
	return &auth.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: userChangesToAPI(response),
	}
}

func userChangesToAPI(changes *usr_model.UserChanges) (_ []*auth.Change) {
	result := make([]*auth.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		var data *structpb.Struct
		changedData, err := json.Marshal(change.Data)
		if err == nil {
			data = new(structpb.Struct)
			err = protojson.Unmarshal(changedData, data)
			logging.Log("GRPC-0kRsY").OnError(err).Debug("unable to marshal changed data to struct")
		}
		result[i] = &auth.Change{
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

func verifyWebAuthNFromDomain(u2f *domain.WebAuthNToken) *auth.WebAuthNResponse {
	return &auth.WebAuthNResponse{
		Id:        u2f.WebAuthNTokenID,
		PublicKey: u2f.CredentialCreationData,
		State:     mfaStateFromDomain(u2f.State),
	}
}

func webAuthNTokensFromModel(tokens []*usr_model.WebAuthNView) *auth.WebAuthNTokens {
	result := make([]*auth.WebAuthNToken, len(tokens))
	for i, token := range tokens {
		result[i] = webAuthNTokenFromModel(token)
	}
	return &auth.WebAuthNTokens{Tokens: result}
}

func webAuthNTokenFromModel(token *usr_model.WebAuthNView) *auth.WebAuthNToken {
	return &auth.WebAuthNToken{
		Id:    token.TokenID,
		Name:  token.Name,
		State: auth.MFAState(token.State),
	}
}

func ctxToObjectRoot(ctx context.Context) models.ObjectRoot {
	ctxData := authz.GetCtxData(ctx)
	return models.ObjectRoot{
		AggregateID:   ctxData.UserID,
		ResourceOwner: ctxData.ResourceOwner,
	}
}

func userMembershipSearchResponseFromModel(response *usr_model.UserMembershipSearchResponse) *auth.UserMembershipSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-Hs8jd").OnError(err).Debug("unable to parse timestamp")
	return &auth.UserMembershipSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            userMembershipViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func userMembershipViewsFromModel(memberships []*usr_model.UserMembershipView) []*auth.UserMembershipView {
	converted := make([]*auth.UserMembershipView, len(memberships))
	for i, membership := range memberships {
		converted[i] = userMembershipViewFromModel(membership)
	}
	return converted
}

func userMembershipViewFromModel(membership *usr_model.UserMembershipView) *auth.UserMembershipView {
	creationDate, err := ptypes.TimestampProto(membership.CreationDate)
	logging.Log("GRPC-Msnu8").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(membership.ChangeDate)
	logging.Log("GRPC-Slco9").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserMembershipView{
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

func userMembershipSearchRequestsToModel(request *auth.UserMembershipSearchRequest) *usr_model.UserMembershipSearchRequest {
	return &usr_model.UserMembershipSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: userMembershipSearchQueriesToModel(request.Queries),
	}
}

func userMembershipSearchQueriesToModel(queries []*auth.UserMembershipSearchQuery) []*usr_model.UserMembershipSearchQuery {
	converted := make([]*usr_model.UserMembershipSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userMembershipSearchQueryToModel(q)
	}
	return converted
}

func userMembershipSearchQueryToModel(query *auth.UserMembershipSearchQuery) *usr_model.UserMembershipSearchQuery {
	return &usr_model.UserMembershipSearchQuery{
		Key:    userMembershipSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userMembershipSearchKeyToModel(key auth.UserMembershipSearchKey) usr_model.UserMembershipSearchKey {
	switch key {
	case auth.UserMembershipSearchKey_USERMEMBERSHIPSEARCHKEY_TYPE:
		return usr_model.UserMembershipSearchKeyMemberType
	case auth.UserMembershipSearchKey_USERMEMBERSHIPSEARCHKEY_OBJECT_ID:
		return usr_model.UserMembershipSearchKeyObjectID
	default:
		return usr_model.UserMembershipSearchKeyUnspecified
	}
}

func memberTypeFromModel(memberType usr_model.MemberType) auth.MemberType {
	switch memberType {
	case usr_model.MemberTypeOrganisation:
		return auth.MemberType_MEMBERTYPE_ORGANISATION
	case usr_model.MemberTypeProject:
		return auth.MemberType_MEMBERTYPE_PROJECT
	case usr_model.MemberTypeProjectGrant:
		return auth.MemberType_MEMBERTYPE_PROJECT_GRANT
	default:
		return auth.MemberType_MEMBERTYPE_UNSPECIFIED
	}
}
