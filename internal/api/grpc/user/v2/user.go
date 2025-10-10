package user

import (
	"context"
	"io"

	"connectrpc.com/connect"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddHumanUser(ctx context.Context, req *connect.Request[user.AddHumanUserRequest]) (_ *connect.Response[user.AddHumanUserResponse], err error) {
	human, err := AddUserRequestToAddHuman(req.Msg)
	if err != nil {
		return nil, err
	}
	orgID := authz.GetCtxData(ctx).OrgID
	if err = s.command.AddUserHuman(ctx, orgID, human, false, s.userCodeAlg); err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.AddHumanUserResponse{
		UserId:    human.ID,
		Details:   object.DomainToDetailsPb(human.Details),
		EmailCode: human.EmailCode,
		PhoneCode: human.PhoneCode,
	}), nil
}

func AddUserRequestToAddHuman(req *user.AddHumanUserRequest) (*command.AddHuman, error) {
	username := req.GetUsername()
	if username == "" {
		username = req.GetEmail().GetEmail()
	}
	passwordChangeRequired := req.GetPassword().GetChangeRequired() || req.GetHashedPassword().GetChangeRequired()
	metadata := make([]*command.AddMetadataEntry, len(req.Metadata))
	for i, metadataEntry := range req.Metadata {
		metadata[i] = &command.AddMetadataEntry{
			Key:   metadataEntry.GetKey(),
			Value: metadataEntry.GetValue(),
		}
	}
	links := make([]*command.AddLink, len(req.GetIdpLinks()))
	for i, link := range req.GetIdpLinks() {
		links[i] = &command.AddLink{
			IDPID:         link.GetIdpId(),
			IDPExternalID: link.GetUserId(),
			DisplayName:   link.GetUserName(),
		}
	}
	email, err := addUserRequestEmailToCommand(req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &command.AddHuman{
		ID:          req.GetUserId(),
		Username:    username,
		FirstName:   req.GetProfile().GetGivenName(),
		LastName:    req.GetProfile().GetFamilyName(),
		NickName:    req.GetProfile().GetNickName(),
		DisplayName: req.GetProfile().GetDisplayName(),
		Email:       email,
		Phone: command.Phone{
			Number:     domain.PhoneNumber(req.GetPhone().GetPhone()),
			Verified:   req.GetPhone().GetIsVerified(),
			ReturnCode: req.GetPhone().GetReturnCode() != nil,
		},
		PreferredLanguage:      language.Make(req.GetProfile().GetPreferredLanguage()),
		Gender:                 genderToDomain(req.GetProfile().GetGender()),
		Password:               req.GetPassword().GetPassword(),
		EncodedPasswordHash:    req.GetHashedPassword().GetHash(),
		PasswordChangeRequired: passwordChangeRequired,
		Passwordless:           false,
		Register:               false,
		Metadata:               metadata,
		Links:                  links,
		TOTPSecret:             req.GetTotpSecret(),
	}, nil
}

func addUserRequestEmailToCommand(email *user.SetHumanEmail) (command.Email, error) {
	address := domain.EmailAddress(email.GetEmail())
	switch v := email.GetVerification().(type) {
	case *user.SetHumanEmail_ReturnCode:
		return command.Email{Address: address, ReturnCode: true}, nil
	case *user.SetHumanEmail_SendCode:
		urlTemplate := v.SendCode.GetUrlTemplate()
		// test the template execution so the async notification will not fail because of it and the user won't realize
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, "userID", "code", "orgID"); err != nil {
			return command.Email{}, err
		}
		return command.Email{Address: address, URLTemplate: urlTemplate}, nil
	case *user.SetHumanEmail_IsVerified:
		return command.Email{Address: address, Verified: v.IsVerified, NoEmailVerification: true}, nil
	default:
		return command.Email{Address: address}, nil
	}
}

func genderToDomain(gender user.Gender) domain.Gender {
	switch gender {
	case user.Gender_GENDER_UNSPECIFIED:
		return domain.GenderUnspecified
	case user.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	case user.Gender_GENDER_MALE:
		return domain.GenderMale
	case user.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	default:
		return domain.GenderUnspecified
	}
}

func (s *Server) UpdateHumanUser(ctx context.Context, req *connect.Request[user.UpdateHumanUserRequest]) (_ *connect.Response[user.UpdateHumanUserResponse], err error) {
	human, err := updateHumanUserRequestToChangeHuman(req.Msg)
	if err != nil {
		return nil, err
	}
	err = s.command.ChangeUserHuman(ctx, human, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.UpdateHumanUserResponse{
		Details:   object.DomainToDetailsPb(human.Details),
		EmailCode: human.EmailCode,
		PhoneCode: human.PhoneCode,
	}), nil
}

func (s *Server) LockUser(ctx context.Context, req *connect.Request[user.LockUserRequest]) (_ *connect.Response[user.LockUserResponse], err error) {
	details, err := s.command.LockUserV2(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.LockUserResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) UnlockUser(ctx context.Context, req *connect.Request[user.UnlockUserRequest]) (_ *connect.Response[user.UnlockUserResponse], err error) {
	details, err := s.command.UnlockUserV2(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.UnlockUserResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *connect.Request[user.DeactivateUserRequest]) (_ *connect.Response[user.DeactivateUserResponse], err error) {
	details, err := s.command.DeactivateUserV2(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.DeactivateUserResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) ReactivateUser(ctx context.Context, req *connect.Request[user.ReactivateUserRequest]) (_ *connect.Response[user.ReactivateUserResponse], err error) {
	details, err := s.command.ReactivateUserV2(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ReactivateUserResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func ifNotNilPtr[v, p any](value *v, conv func(v) p) *p {
	var pNil *p
	if value == nil {
		return pNil
	}
	pVal := conv(*value)
	return &pVal
}

func (s *Server) DeleteUser(ctx context.Context, req *connect.Request[user.DeleteUserRequest]) (_ *connect.Response[user.DeleteUserResponse], err error) {
	memberships, grants, err := s.removeUserDependencies(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveUserV2(ctx, req.Msg.GetUserId(), "", memberships, grants...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.DeleteUserResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) removeUserDependencies(ctx context.Context, userID string) ([]*command.CascadingMembership, []string, error) {
	userGrantUserQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{userGrantUserQuery},
	}, true, nil)
	if err != nil {
		return nil, nil, err
	}
	membershipsUserQuery, err := query.NewMembershipUserIDQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	memberships, err := s.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{membershipsUserQuery},
	}, false)
	if err != nil {
		return nil, nil, err
	}
	return cascadingMemberships(memberships.Memberships), userGrantsToIDs(grants.UserGrants), nil
}

func cascadingMemberships(memberships []*query.Membership) []*command.CascadingMembership {
	cascades := make([]*command.CascadingMembership, len(memberships))
	for i, membership := range memberships {
		cascades[i] = &command.CascadingMembership{
			UserID:        membership.UserID,
			ResourceOwner: membership.ResourceOwner,
			IAM:           cascadingIAMMembership(membership.IAM),
			Org:           cascadingOrgMembership(membership.Org),
			Project:       cascadingProjectMembership(membership.Project),
			ProjectGrant:  cascadingProjectGrantMembership(membership.ProjectGrant),
		}
	}
	return cascades
}

func cascadingIAMMembership(membership *query.IAMMembership) *command.CascadingIAMMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingIAMMembership{IAMID: membership.IAMID}
}
func cascadingOrgMembership(membership *query.OrgMembership) *command.CascadingOrgMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingOrgMembership{OrgID: membership.OrgID}
}
func cascadingProjectMembership(membership *query.ProjectMembership) *command.CascadingProjectMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectMembership{ProjectID: membership.ProjectID}
}
func cascadingProjectGrantMembership(membership *query.ProjectGrantMembership) *command.CascadingProjectGrantMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectGrantMembership{ProjectID: membership.ProjectID, GrantID: membership.GrantID}
}

func userGrantsToIDs(userGrants []*query.UserGrant) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}

func (s *Server) ListAuthenticationMethodTypes(ctx context.Context, req *connect.Request[user.ListAuthenticationMethodTypesRequest]) (*connect.Response[user.ListAuthenticationMethodTypesResponse], error) {
	authMethods, err := s.query.ListUserAuthMethodTypes(ctx, req.Msg.GetUserId(), true, req.Msg.GetDomainQuery().GetIncludeWithoutDomain(), req.Msg.GetDomainQuery().GetDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ListAuthenticationMethodTypesResponse{
		Details:         object.ToListDetails(authMethods.SearchResponse),
		AuthMethodTypes: authMethodTypesToPb(authMethods.AuthMethodTypes),
	}), nil
}

func (s *Server) ListAuthenticationFactors(ctx context.Context, req *connect.Request[user.ListAuthenticationFactorsRequest]) (*connect.Response[user.ListAuthenticationFactorsResponse], error) {
	query := new(query.UserAuthMethodSearchQueries)

	if err := query.AppendUserIDQuery(req.Msg.GetUserId()); err != nil {
		return nil, err
	}

	authMethodsType := []domain.UserAuthMethodType{domain.UserAuthMethodTypeU2F, domain.UserAuthMethodTypeTOTP, domain.UserAuthMethodTypeOTPSMS, domain.UserAuthMethodTypeOTPEmail}
	if len(req.Msg.GetAuthFactors()) > 0 {
		authMethodsType = object.AuthFactorsToPb(req.Msg.GetAuthFactors())
	}
	if err := query.AppendAuthMethodsQuery(authMethodsType...); err != nil {
		return nil, err
	}

	states := []domain.MFAState{domain.MFAStateReady}
	if len(req.Msg.GetStates()) > 0 {
		states = object.AuthFactorStatesToPb(req.Msg.GetStates())
	}
	if err := query.AppendStatesQuery(states...); err != nil {
		return nil, err
	}

	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&user.ListAuthenticationFactorsResponse{
		Result: object.AuthMethodsToPb(authMethods),
	}), nil
}

func authMethodTypesToPb(methodTypes []domain.UserAuthMethodType) []user.AuthenticationMethodType {
	methods := make([]user.AuthenticationMethodType, len(methodTypes))
	for i, method := range methodTypes {
		methods[i] = authMethodTypeToPb(method)
	}
	return methods
}

func authMethodTypeToPb(methodType domain.UserAuthMethodType) user.AuthenticationMethodType {
	switch methodType {
	case domain.UserAuthMethodTypeTOTP:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_TOTP
	case domain.UserAuthMethodTypeU2F:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_U2F
	case domain.UserAuthMethodTypePasswordless:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY
	case domain.UserAuthMethodTypePassword:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSWORD
	case domain.UserAuthMethodTypeIDP:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP
	case domain.UserAuthMethodTypeOTPSMS:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_OTP_SMS
	case domain.UserAuthMethodTypeOTPEmail:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_OTP_EMAIL
	case domain.UserAuthMethodTypeUnspecified:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED
	default:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED
	}
}

func (s *Server) CreateInviteCode(ctx context.Context, req *connect.Request[user.CreateInviteCodeRequest]) (*connect.Response[user.CreateInviteCodeResponse], error) {
	invite, err := createInviteCodeRequestToCommand(req.Msg)
	if err != nil {
		return nil, err
	}
	details, code, err := s.command.CreateInviteCode(ctx, invite)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.CreateInviteCodeResponse{
		Details:    object.DomainToDetailsPb(details),
		InviteCode: code,
	}), nil
}

func (s *Server) ResendInviteCode(ctx context.Context, req *connect.Request[user.ResendInviteCodeRequest]) (*connect.Response[user.ResendInviteCodeResponse], error) {
	details, err := s.command.ResendInviteCode(ctx, req.Msg.GetUserId(), "", "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ResendInviteCodeResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) VerifyInviteCode(ctx context.Context, req *connect.Request[user.VerifyInviteCodeRequest]) (*connect.Response[user.VerifyInviteCodeResponse], error) {
	details, err := s.command.VerifyInviteCode(ctx, req.Msg.GetUserId(), req.Msg.GetVerificationCode())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.VerifyInviteCodeResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func createInviteCodeRequestToCommand(req *user.CreateInviteCodeRequest) (*command.CreateUserInvite, error) {
	switch v := req.GetVerification().(type) {
	case *user.CreateInviteCodeRequest_SendCode:
		urlTemplate := v.SendCode.GetUrlTemplate()
		// test the template execution so the async notification will not fail because of it and the user won't realize
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, req.GetUserId(), "code", "orgID"); err != nil {
			return nil, err
		}
		return &command.CreateUserInvite{UserID: req.GetUserId(), URLTemplate: urlTemplate, ApplicationName: v.SendCode.GetApplicationName()}, nil
	case *user.CreateInviteCodeRequest_ReturnCode:
		return &command.CreateUserInvite{UserID: req.GetUserId(), ReturnCode: true}, nil
	default:
		return &command.CreateUserInvite{UserID: req.GetUserId()}, nil
	}
}

func (s *Server) HumanMFAInitSkipped(ctx context.Context, req *connect.Request[user.HumanMFAInitSkippedRequest]) (_ *connect.Response[user.HumanMFAInitSkippedResponse], err error) {
	details, err := s.command.HumanMFAInitSkippedV2(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.HumanMFAInitSkippedResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) CreateUser(ctx context.Context, req *connect.Request[user.CreateUserRequest]) (*connect.Response[user.CreateUserResponse], error) {
	switch userType := req.Msg.GetUserType().(type) {
	case *user.CreateUserRequest_Human_:
		return s.createUserTypeHuman(ctx, userType.Human, req.Msg.GetOrganizationId(), req.Msg.Username, req.Msg.UserId)
	case *user.CreateUserRequest_Machine_:
		return s.createUserTypeMachine(ctx, userType.Machine, req.Msg.GetOrganizationId(), req.Msg.GetUsername(), req.Msg.GetUserId())
	default:
		return nil, zerrors.ThrowInternal(nil, "", "user type is not implemented")
	}
}

func (s *Server) UpdateUser(ctx context.Context, req *connect.Request[user.UpdateUserRequest]) (*connect.Response[user.UpdateUserResponse], error) {
	switch userType := req.Msg.GetUserType().(type) {
	case *user.UpdateUserRequest_Human_:
		return s.updateUserTypeHuman(ctx, userType.Human, req.Msg.GetUserId(), req.Msg.Username)
	case *user.UpdateUserRequest_Machine_:
		return s.updateUserTypeMachine(ctx, userType.Machine, req.Msg.GetUserId(), req.Msg.Username)
	default:
		return nil, zerrors.ThrowUnimplemented(nil, "", "user type is not implemented")
	}
}
