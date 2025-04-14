package user

import (
	"context"
	"io"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) AddHumanUser(ctx context.Context, req *user.AddHumanUserRequest) (_ *user.AddHumanUserResponse, err error) {
	human, err := AddUserRequestToAddHuman(req)
	if err != nil {
		return nil, err
	}
	orgID := authz.GetCtxData(ctx).OrgID
	if err = s.command.AddUserHuman(ctx, orgID, human, false, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.AddHumanUserResponse{
		UserId:    human.ID,
		Details:   object.DomainToDetailsPb(human.Details),
		EmailCode: human.EmailCode,
		PhoneCode: human.PhoneCode,
	}, nil
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

func (s *Server) UpdateHumanUser(ctx context.Context, req *user.UpdateHumanUserRequest) (_ *user.UpdateHumanUserResponse, err error) {
	human, err := UpdateUserRequestToChangeHuman(req)
	if err != nil {
		return nil, err
	}
	err = s.command.ChangeUserHuman(ctx, human, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	return &user.UpdateHumanUserResponse{
		Details:   object.DomainToDetailsPb(human.Details),
		EmailCode: human.EmailCode,
		PhoneCode: human.PhoneCode,
	}, nil
}

func (s *Server) LockUser(ctx context.Context, req *user.LockUserRequest) (_ *user.LockUserResponse, err error) {
	details, err := s.command.LockUserV2(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &user.LockUserResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UnlockUser(ctx context.Context, req *user.UnlockUserRequest) (_ *user.UnlockUserResponse, err error) {
	details, err := s.command.UnlockUserV2(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &user.UnlockUserResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *user.DeactivateUserRequest) (_ *user.DeactivateUserResponse, err error) {
	details, err := s.command.DeactivateUserV2(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &user.DeactivateUserResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateUser(ctx context.Context, req *user.ReactivateUserRequest) (_ *user.ReactivateUserResponse, err error) {
	details, err := s.command.ReactivateUserV2(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &user.ReactivateUserResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func ifNotNilPtr[v, p any](value *v, conv func(v) p) *p {
	var pNil *p
	if value == nil {
		return pNil
	}
	pVal := conv(*value)
	return &pVal
}

func UpdateUserRequestToChangeHuman(req *user.UpdateHumanUserRequest) (*command.ChangeHuman, error) {
	email, err := SetHumanEmailToEmail(req.Email, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &command.ChangeHuman{
		ID:       req.GetUserId(),
		Username: req.Username,
		Profile:  SetHumanProfileToProfile(req.Profile),
		Email:    email,
		Phone:    SetHumanPhoneToPhone(req.Phone),
		Password: SetHumanPasswordToPassword(req.Password),
	}, nil
}

func SetHumanProfileToProfile(profile *user.SetHumanProfile) *command.Profile {
	if profile == nil {
		return nil
	}
	var firstName *string
	if profile.GivenName != "" {
		firstName = &profile.GivenName
	}
	var lastName *string
	if profile.FamilyName != "" {
		lastName = &profile.FamilyName
	}
	return &command.Profile{
		FirstName:         firstName,
		LastName:          lastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: ifNotNilPtr(profile.PreferredLanguage, language.Make),
		Gender:            ifNotNilPtr(profile.Gender, genderToDomain),
	}
}

func SetHumanEmailToEmail(email *user.SetHumanEmail, userID string) (*command.Email, error) {
	if email == nil {
		return nil, nil
	}
	var urlTemplate string
	if email.GetSendCode() != nil && email.GetSendCode().UrlTemplate != nil {
		urlTemplate = *email.GetSendCode().UrlTemplate
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, userID, "code", "orgID"); err != nil {
			return nil, err
		}
	}
	return &command.Email{
		Address:     domain.EmailAddress(email.Email),
		Verified:    email.GetIsVerified(),
		ReturnCode:  email.GetReturnCode() != nil,
		URLTemplate: urlTemplate,
	}, nil
}

func SetHumanPhoneToPhone(phone *user.SetHumanPhone) *command.Phone {
	if phone == nil {
		return nil
	}
	return &command.Phone{
		Number:     domain.PhoneNumber(phone.GetPhone()),
		Verified:   phone.GetIsVerified(),
		ReturnCode: phone.GetReturnCode() != nil,
	}
}

func SetHumanPasswordToPassword(password *user.SetPassword) *command.Password {
	if password == nil {
		return nil
	}
	return &command.Password{
		PasswordCode:        password.GetVerificationCode(),
		OldPassword:         password.GetCurrentPassword(),
		Password:            password.GetPassword().GetPassword(),
		EncodedPasswordHash: password.GetHashedPassword().GetHash(),
		ChangeRequired:      password.GetPassword().GetChangeRequired() || password.GetHashedPassword().GetChangeRequired(),
	}
}

func (s *Server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (_ *user.DeleteUserResponse, err error) {
	memberships, grants, err := s.removeUserDependencies(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveUserV2(ctx, req.UserId, "", memberships, grants...)
	if err != nil {
		return nil, err
	}
	return &user.DeleteUserResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) removeUserDependencies(ctx context.Context, userID string) ([]*command.CascadingMembership, []string, error) {
	userGrantUserQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{userGrantUserQuery},
	}, true)
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

func (s *Server) ListAuthenticationMethodTypes(ctx context.Context, req *user.ListAuthenticationMethodTypesRequest) (*user.ListAuthenticationMethodTypesResponse, error) {
	authMethods, err := s.query.ListUserAuthMethodTypes(ctx, req.GetUserId(), true, req.GetDomainQuery().GetIncludeWithoutDomain(), req.GetDomainQuery().GetDomain())
	if err != nil {
		return nil, err
	}
	return &user.ListAuthenticationMethodTypesResponse{
		Details:         object.ToListDetails(authMethods.SearchResponse),
		AuthMethodTypes: authMethodTypesToPb(authMethods.AuthMethodTypes),
	}, nil
}

func (s *Server) ListAuthenticationFactors(ctx context.Context, req *user.ListAuthenticationFactorsRequest) (*user.ListAuthenticationFactorsResponse, error) {
	query := new(query.UserAuthMethodSearchQueries)

	if err := query.AppendUserIDQuery(req.UserId); err != nil {
		return nil, err
	}

	authMethodsType := []domain.UserAuthMethodType{domain.UserAuthMethodTypeU2F, domain.UserAuthMethodTypeTOTP, domain.UserAuthMethodTypeOTPSMS, domain.UserAuthMethodTypeOTPEmail}
	if len(req.GetAuthFactors()) > 0 {
		authMethodsType = object.AuthFactorsToPb(req.GetAuthFactors())
	}
	if err := query.AppendAuthMethodsQuery(authMethodsType...); err != nil {
		return nil, err
	}

	states := []domain.MFAState{domain.MFAStateReady}
	if len(req.GetStates()) > 0 {
		states = object.AuthFactorStatesToPb(req.GetStates())
	}
	if err := query.AppendStatesQuery(states...); err != nil {
		return nil, err
	}

	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return &user.ListAuthenticationFactorsResponse{
		Result: object.AuthMethodsToPb(authMethods),
	}, nil
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

func (s *Server) CreateInviteCode(ctx context.Context, req *user.CreateInviteCodeRequest) (*user.CreateInviteCodeResponse, error) {
	invite, err := createInviteCodeRequestToCommand(req)
	if err != nil {
		return nil, err
	}
	details, code, err := s.command.CreateInviteCode(ctx, invite)
	if err != nil {
		return nil, err
	}
	return &user.CreateInviteCodeResponse{
		Details:    object.DomainToDetailsPb(details),
		InviteCode: code,
	}, nil
}

func (s *Server) ResendInviteCode(ctx context.Context, req *user.ResendInviteCodeRequest) (*user.ResendInviteCodeResponse, error) {
	details, err := s.command.ResendInviteCode(ctx, req.GetUserId(), "", "")
	if err != nil {
		return nil, err
	}
	return &user.ResendInviteCodeResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) VerifyInviteCode(ctx context.Context, req *user.VerifyInviteCodeRequest) (*user.VerifyInviteCodeResponse, error) {
	details, err := s.command.VerifyInviteCode(ctx, req.GetUserId(), req.GetVerificationCode())
	if err != nil {
		return nil, err
	}
	return &user.VerifyInviteCodeResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
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

func (s *Server) HumanMFAInitSkipped(ctx context.Context, req *user.HumanMFAInitSkippedRequest) (_ *user.HumanMFAInitSkippedResponse, err error) {
	details, err := s.command.HumanMFAInitSkippedV2(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &user.HumanMFAInitSkippedResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}
