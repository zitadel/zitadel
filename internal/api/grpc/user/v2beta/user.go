package user

import (
	"context"
	"errors"
	"io"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
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
	var urlTemplate string
	if req.GetEmail().GetSendCode() != nil {
		urlTemplate = req.GetEmail().GetSendCode().GetUrlTemplate()
		// test the template execution so the async notification will not fail because of it and the user won't realize
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, req.GetUserId(), "code", "orgID"); err != nil {
			return nil, err
		}
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
	return &command.AddHuman{
		ID:          req.GetUserId(),
		Username:    username,
		FirstName:   req.GetProfile().GetGivenName(),
		LastName:    req.GetProfile().GetFamilyName(),
		NickName:    req.GetProfile().GetNickName(),
		DisplayName: req.GetProfile().GetDisplayName(),
		Email: command.Email{
			Address:     domain.EmailAddress(req.GetEmail().GetEmail()),
			Verified:    req.GetEmail().GetIsVerified(),
			ReturnCode:  req.GetEmail().GetReturnCode() != nil,
			URLTemplate: urlTemplate,
		},
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
	human, err := UpdateUserRequestToChangeHuman(req.Msg)
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

func (s *Server) AddIDPLink(ctx context.Context, req *connect.Request[user.AddIDPLinkRequest]) (_ *connect.Response[user.AddIDPLinkResponse], err error) {
	details, err := s.command.AddUserIDPLink(ctx, req.Msg.GetUserId(), "", &command.AddLink{
		IDPID:         req.Msg.GetIdpLink().GetIdpId(),
		DisplayName:   req.Msg.GetIdpLink().GetUserName(),
		IDPExternalID: req.Msg.GetIdpLink().GetUserId(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.AddIDPLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
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

func (s *Server) StartIdentityProviderIntent(ctx context.Context, req *connect.Request[user.StartIdentityProviderIntentRequest]) (_ *connect.Response[user.StartIdentityProviderIntentResponse], err error) {
	switch t := req.Msg.GetContent().(type) {
	case *user.StartIdentityProviderIntentRequest_Urls:
		return s.startIDPIntent(ctx, req.Msg.GetIdpId(), t.Urls)
	case *user.StartIdentityProviderIntentRequest_Ldap:
		return s.startLDAPIntent(ctx, req.Msg.GetIdpId(), t.Ldap)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "USERv2-S2g21", "type oneOf %T in method StartIdentityProviderIntent not implemented", t)
	}
}

func (s *Server) startIDPIntent(ctx context.Context, idpID string, urls *user.RedirectURLs) (*connect.Response[user.StartIdentityProviderIntentResponse], error) {
	state, session, err := s.command.AuthFromProvider(ctx, idpID, s.idpCallback(ctx), s.samlRootURL(ctx, idpID))
	if err != nil {
		return nil, err
	}
	_, details, err := s.command.CreateIntent(ctx, state, idpID, urls.GetSuccessUrl(), urls.GetFailureUrl(), authz.GetInstance(ctx).InstanceID(), session.PersistentParameters())
	if err != nil {
		return nil, err
	}
	auth, err := session.GetAuth(ctx)
	if err != nil {
		return nil, err
	}
	switch a := auth.(type) {
	case *idp.RedirectAuth:
		return connect.NewResponse(&user.StartIdentityProviderIntentResponse{
			Details:  object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderIntentResponse_AuthUrl{AuthUrl: a.RedirectURL},
		}), nil
	case *idp.FormAuth:
		return connect.NewResponse(&user.StartIdentityProviderIntentResponse{
			Details: object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderIntentResponse_FormData{
				FormData: &user.FormData{
					Url:    a.URL,
					Fields: a.Fields,
				},
			},
		}), nil
	}
	return nil, zerrors.ThrowInvalidArgumentf(nil, "USERv2-3g2j3", "type oneOf %T in method StartIdentityProviderIntent not implemented", auth)
}

func (s *Server) startLDAPIntent(ctx context.Context, idpID string, ldapCredentials *user.LDAPCredentials) (*connect.Response[user.StartIdentityProviderIntentResponse], error) {
	intentWriteModel, details, err := s.command.CreateIntent(ctx, "", idpID, "", "", authz.GetInstance(ctx).InstanceID(), nil)
	if err != nil {
		return nil, err
	}
	externalUser, userID, session, err := s.ldapLogin(ctx, intentWriteModel.IDPID, ldapCredentials.GetUsername(), ldapCredentials.GetPassword())
	if err != nil {
		if err := s.command.FailIDPIntent(ctx, intentWriteModel, err.Error()); err != nil {
			return nil, err
		}
		return nil, err
	}
	token, err := s.command.SucceedLDAPIDPIntent(ctx, intentWriteModel, externalUser, userID, session)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.StartIdentityProviderIntentResponse{
		Details: object.DomainToDetailsPb(details),
		NextStep: &user.StartIdentityProviderIntentResponse_IdpIntent{
			IdpIntent: &user.IDPIntent{
				IdpIntentId:    intentWriteModel.AggregateID,
				IdpIntentToken: token,
				UserId:         userID,
			},
		},
	}), nil
}

func (s *Server) checkLinkedExternalUser(ctx context.Context, idpID, externalUserID string) (string, error) {
	idQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(idpID)
	if err != nil {
		return "", err
	}
	externalIDQuery, err := query.NewIDPUserLinksExternalIDSearchQuery(externalUserID)
	if err != nil {
		return "", err
	}
	queries := []query.SearchQuery{
		idQuery, externalIDQuery,
	}
	links, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: queries}, nil)
	if err != nil {
		return "", err
	}
	if len(links.Links) == 1 {
		return links.Links[0].UserID, nil
	}
	return "", nil
}

func (s *Server) ldapLogin(ctx context.Context, idpID, username, password string) (idp.User, string, *ldap.Session, error) {
	provider, err := s.command.GetProvider(ctx, idpID, "", "")
	if err != nil {
		return nil, "", nil, err
	}
	ldapProvider, ok := provider.(*ldap.Provider)
	if !ok {
		return nil, "", nil, zerrors.ThrowInvalidArgument(nil, "IDP-9a02j2n2bh", "Errors.ExternalIDP.IDPTypeNotImplemented")
	}
	session := ldapProvider.GetSession(username, password)
	externalUser, err := session.FetchUser(ctx)
	if errors.Is(err, ldap.ErrFailedLogin) || errors.Is(err, ldap.ErrNoSingleUser) {
		return nil, "", nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-nzun2i", "Errors.User.ExternalIDP.LoginFailed")
	}
	if err != nil {
		return nil, "", nil, err
	}
	userID, err := s.checkLinkedExternalUser(ctx, idpID, externalUser.GetID())
	if err != nil {
		return nil, "", nil, err
	}

	attributes := make(map[string][]string, 0)
	for _, item := range session.Entry.Attributes {
		attributes[item.Name] = item.Values
	}
	return externalUser, userID, session, nil
}

func (s *Server) RetrieveIdentityProviderIntent(ctx context.Context, req *connect.Request[user.RetrieveIdentityProviderIntentRequest]) (_ *connect.Response[user.RetrieveIdentityProviderIntentResponse], err error) {
	intent, err := s.command.GetIntentWriteModel(ctx, req.Msg.GetIdpIntentId(), "")
	if err != nil {
		return nil, err
	}
	if err := s.checkIntentToken(req.Msg.GetIdpIntentToken(), intent.AggregateID); err != nil {
		return nil, err
	}
	if intent.State != domain.IDPIntentStateSucceeded {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-nme4gszsvx", "Errors.Intent.NotSucceeded")
	}
	if time.Now().After(intent.ExpiresAt()) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-Afb2s", "Errors.Intent.Expired")
	}
	return idpIntentToIDPIntentPb(intent, s.idpAlg)
}

func idpIntentToIDPIntentPb(intent *command.IDPIntentWriteModel, alg crypto.EncryptionAlgorithm) (_ *connect.Response[user.RetrieveIdentityProviderIntentResponse], err error) {
	rawInformation := new(structpb.Struct)
	err = rawInformation.UnmarshalJSON(intent.IDPUser)
	if err != nil {
		return nil, err
	}
	information := &user.RetrieveIdentityProviderIntentResponse{
		Details: intentToDetailsPb(intent),
		IdpInformation: &user.IDPInformation{
			IdpId:          intent.IDPID,
			UserId:         intent.IDPUserID,
			UserName:       intent.IDPUserName,
			RawInformation: rawInformation,
		},
		UserId: intent.UserID,
	}
	if intent.IDPIDToken != "" || intent.IDPAccessToken != nil {
		information.IdpInformation.Access, err = idpOAuthTokensToPb(intent.IDPIDToken, intent.IDPAccessToken, alg)
		if err != nil {
			return nil, err
		}
	}

	if intent.IDPEntryAttributes != nil {
		access, err := IDPEntryAttributesToPb(intent.IDPEntryAttributes)
		if err != nil {
			return nil, err
		}
		information.IdpInformation.Access = access
	}

	if intent.Assertion != nil {
		assertion, err := crypto.Decrypt(intent.Assertion, alg)
		if err != nil {
			return nil, err
		}
		information.IdpInformation.Access = IDPSAMLResponseToPb(assertion)
	}

	return connect.NewResponse(information), nil
}

func idpOAuthTokensToPb(idpIDToken string, idpAccessToken *crypto.CryptoValue, alg crypto.EncryptionAlgorithm) (_ *user.IDPInformation_Oauth, err error) {
	var idToken *string
	if idpIDToken != "" {
		idToken = &idpIDToken
	}
	var accessToken string
	if idpAccessToken != nil {
		accessToken, err = crypto.DecryptString(idpAccessToken, alg)
		if err != nil {
			return nil, err
		}
	}
	return &user.IDPInformation_Oauth{
		Oauth: &user.IDPOAuthAccessInformation{
			AccessToken: accessToken,
			IdToken:     idToken,
		},
	}, nil
}

func intentToDetailsPb(intent *command.IDPIntentWriteModel) *object_pb.Details {
	return &object_pb.Details{
		Sequence:      intent.ProcessedSequence,
		ChangeDate:    timestamppb.New(intent.ChangeDate),
		ResourceOwner: intent.ResourceOwner,
	}
}

func IDPEntryAttributesToPb(entryAttributes map[string][]string) (*user.IDPInformation_Ldap, error) {
	values := make(map[string]interface{}, 0)
	for k, v := range entryAttributes {
		intValues := make([]interface{}, len(v))
		for i, value := range v {
			intValues[i] = value
		}
		values[k] = intValues
	}
	attributes, err := structpb.NewStruct(values)
	if err != nil {
		return nil, err
	}
	return &user.IDPInformation_Ldap{
		Ldap: &user.IDPLDAPAccessInformation{
			Attributes: attributes,
		},
	}, nil
}

func IDPSAMLResponseToPb(assertion []byte) *user.IDPInformation_Saml {
	return &user.IDPInformation_Saml{
		Saml: &user.IDPSAMLAccessInformation{
			Assertion: assertion,
		},
	}
}

func (s *Server) checkIntentToken(token string, intentID string) error {
	return crypto.CheckToken(s.idpAlg, token, intentID)
}

func (s *Server) ListAuthenticationMethodTypes(ctx context.Context, req *connect.Request[user.ListAuthenticationMethodTypesRequest]) (*connect.Response[user.ListAuthenticationMethodTypesResponse], error) {
	authMethods, err := s.query.ListUserAuthMethodTypes(ctx, req.Msg.GetUserId(), true, false, "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ListAuthenticationMethodTypesResponse{
		Details:         object.ToListDetails(authMethods.SearchResponse),
		AuthMethodTypes: authMethodTypesToPb(authMethods.AuthMethodTypes),
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
	case domain.UserAuthMethodTypeUnspecified, domain.UserAuthMethodTypeOTP, domain.UserAuthMethodTypePrivateKey:
		// Handle all remaining cases so the linter succeeds
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED
	default:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED
	}
}
