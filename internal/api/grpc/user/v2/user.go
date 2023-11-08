package user

import (
	"context"
	errs "errors"
	"io"

	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/query"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func (s *Server) AddHumanUser(ctx context.Context, req *user.AddHumanUserRequest) (_ *user.AddHumanUserResponse, err error) {
	human, err := AddUserRequestToAddHuman(req)
	if err != nil {
		return nil, err
	}
	orgID := authz.GetCtxData(ctx).OrgID
	err = s.command.AddHuman(ctx, orgID, human, false)
	if err != nil {
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

func (s *Server) AddIDPLink(ctx context.Context, req *user.AddIDPLinkRequest) (_ *user.AddIDPLinkResponse, err error) {
	orgID := authz.GetCtxData(ctx).OrgID
	details, err := s.command.AddUserIDPLink(ctx, req.UserId, orgID, &command.AddLink{
		IDPID:         req.GetIdpLink().GetIdpId(),
		DisplayName:   req.GetIdpLink().GetUserName(),
		IDPExternalID: req.GetIdpLink().GetUserId(),
	})
	if err != nil {
		return nil, err
	}
	return &user.AddIDPLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) StartIdentityProviderIntent(ctx context.Context, req *user.StartIdentityProviderIntentRequest) (_ *user.StartIdentityProviderIntentResponse, err error) {
	switch t := req.GetContent().(type) {
	case *user.StartIdentityProviderIntentRequest_Urls:
		return s.startIDPIntent(ctx, req.GetIdpId(), t.Urls)
	case *user.StartIdentityProviderIntentRequest_Ldap:
		return s.startLDAPIntent(ctx, req.GetIdpId(), t.Ldap)
	default:
		return nil, errors.ThrowUnimplementedf(nil, "USERv2-S2g21", "type oneOf %T in method StartIdentityProviderIntent not implemented", t)
	}
}

func (s *Server) startIDPIntent(ctx context.Context, idpID string, urls *user.RedirectURLs) (*user.StartIdentityProviderIntentResponse, error) {
	intentWriteModel, details, err := s.command.CreateIntent(ctx, idpID, urls.GetSuccessUrl(), urls.GetFailureUrl(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	content, redirect, err := s.command.AuthFromProvider(ctx, idpID, intentWriteModel.AggregateID, s.idpCallback(ctx), s.samlRootURL(ctx, idpID))
	if err != nil {
		return nil, err
	}
	if redirect {
		return &user.StartIdentityProviderIntentResponse{
			Details:  object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderIntentResponse_AuthUrl{AuthUrl: content},
		}, nil
	} else {
		return &user.StartIdentityProviderIntentResponse{
			Details: object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderIntentResponse_PostForm{
				PostForm: []byte(content),
			},
		}, nil
	}
}

func (s *Server) startLDAPIntent(ctx context.Context, idpID string, ldapCredentials *user.LDAPCredentials) (*user.StartIdentityProviderIntentResponse, error) {
	intentWriteModel, details, err := s.command.CreateIntent(ctx, idpID, "", "", authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	externalUser, userID, attributes, err := s.ldapLogin(ctx, intentWriteModel.IDPID, ldapCredentials.GetUsername(), ldapCredentials.GetPassword())
	if err != nil {
		if err := s.command.FailIDPIntent(ctx, intentWriteModel, err.Error()); err != nil {
			return nil, err
		}
		return nil, err
	}
	token, err := s.command.SucceedLDAPIDPIntent(ctx, intentWriteModel, externalUser, userID, attributes)
	if err != nil {
		return nil, err
	}
	return &user.StartIdentityProviderIntentResponse{
		Details: object.DomainToDetailsPb(details),
		NextStep: &user.StartIdentityProviderIntentResponse_IdpIntent{
			IdpIntent: &user.IDPIntent{
				IdpIntentId:    intentWriteModel.AggregateID,
				IdpIntentToken: token,
				UserId:         userID,
			},
		},
	}, nil
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
	links, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: queries}, false)
	if err != nil {
		return "", err
	}
	if len(links.Links) == 1 {
		return links.Links[0].UserID, nil
	}
	return "", nil
}

func (s *Server) ldapLogin(ctx context.Context, idpID, username, password string) (idp.User, string, map[string][]string, error) {
	provider, err := s.command.GetProvider(ctx, idpID, "", "")
	if err != nil {
		return nil, "", nil, err
	}
	ldapProvider, ok := provider.(*ldap.Provider)
	if !ok {
		return nil, "", nil, errors.ThrowInvalidArgument(nil, "IDP-9a02j2n2bh", "Errors.ExternalIDP.IDPTypeNotImplemented")
	}
	session := ldapProvider.GetSession(username, password)
	externalUser, err := session.FetchUser(ctx)
	if errs.Is(err, ldap.ErrFailedLogin) || errs.Is(err, ldap.ErrNoSingleUser) {
		return nil, "", nil, errors.ThrowInvalidArgument(nil, "COMMAND-nzun2i", "Errors.User.ExternalIDP.LoginFailed")
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
	return externalUser, userID, attributes, nil
}

func (s *Server) RetrieveIdentityProviderIntent(ctx context.Context, req *user.RetrieveIdentityProviderIntentRequest) (_ *user.RetrieveIdentityProviderIntentResponse, err error) {
	intent, err := s.command.GetIntentWriteModel(ctx, req.GetIdpIntentId(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	if err := s.checkIntentToken(req.GetIdpIntentToken(), intent.AggregateID); err != nil {
		return nil, err
	}
	if intent.State != domain.IDPIntentStateSucceeded {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-Hk38e", "Errors.Intent.NotSucceeded")
	}
	return idpIntentToIDPIntentPb(intent, s.idpAlg)
}

func idpIntentToIDPIntentPb(intent *command.IDPIntentWriteModel, alg crypto.EncryptionAlgorithm) (_ *user.RetrieveIdentityProviderIntentResponse, err error) {
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

	return information, nil
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

func (s *Server) ListAuthenticationMethodTypes(ctx context.Context, req *user.ListAuthenticationMethodTypesRequest) (*user.ListAuthenticationMethodTypesResponse, error) {
	authMethods, err := s.query.ListActiveUserAuthMethodTypes(ctx, req.GetUserId(), false)
	if err != nil {
		return nil, err
	}
	return &user.ListAuthenticationMethodTypesResponse{
		Details:         object.ToListDetails(authMethods.SearchResponse),
		AuthMethodTypes: authMethodTypesToPb(authMethods.AuthMethodTypes),
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
