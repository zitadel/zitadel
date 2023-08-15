package user

import (
	"context"
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
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) AddHumanUser(ctx context.Context, req *user.AddHumanUserRequest) (_ *user.AddHumanUserResponse, err error) {
	human, err := addUserRequestToAddHuman(req)
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
	}, nil
}

func addUserRequestToAddHuman(req *user.AddHumanUserRequest) (*command.AddHuman, error) {
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
		FirstName:   req.GetProfile().GetFirstName(),
		LastName:    req.GetProfile().GetLastName(),
		NickName:    req.GetProfile().GetNickName(),
		DisplayName: req.GetProfile().GetDisplayName(),
		Email: command.Email{
			Address:     domain.EmailAddress(req.GetEmail().GetEmail()),
			Verified:    req.GetEmail().GetIsVerified(),
			ReturnCode:  req.GetEmail().GetReturnCode() != nil,
			URLTemplate: urlTemplate,
		},
		PreferredLanguage:      language.Make(req.GetProfile().GetPreferredLanguage()),
		Gender:                 genderToDomain(req.GetProfile().GetGender()),
		Phone:                  command.Phone{}, // TODO: add as soon as possible
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
	details, err := s.command.AddUserIDPLink(ctx, req.UserId, orgID, &domain.UserIDPLink{
		IDPConfigID:    req.GetIdpLink().GetIdpId(),
		ExternalUserID: req.GetIdpLink().GetUserId(),
		DisplayName:    req.GetIdpLink().GetUserName(),
	})
	if err != nil {
		return nil, err
	}
	return &user.AddIDPLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) StartIdentityProviderFlow(ctx context.Context, req *user.StartIdentityProviderFlowRequest) (_ *user.StartIdentityProviderFlowResponse, err error) {
	if req.GetLdap() != nil {
		intentWriteModel, details, err := s.command.CreateIntent(ctx, req.GetIdpId(), "", "", authz.GetCtxData(ctx).OrgID)
		if err != nil {
			return nil, err
		}
		externalUser, userID, attributes, err := s.ldapLogin(ctx, intentWriteModel.IDPID, req.GetLdap().GetUsername(), req.GetLdap().GetPassword())
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
		return &user.StartIdentityProviderFlowResponse{
			Details:  object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderFlowResponse_Intent{Intent: &user.Intent{IntentId: intentWriteModel.AggregateID, Token: token}},
		}, nil
	} else {
		intentWriteModel, details, err := s.command.CreateIntent(ctx, req.GetIdpId(), req.GetUrls().GetSuccessUrl(), req.GetUrls().GetFailureUrl(), authz.GetCtxData(ctx).OrgID)
		if err != nil {
			return nil, err
		}
		authURL, err := s.command.AuthURLFromProvider(ctx, req.GetIdpId(), intentWriteModel.AggregateID, s.idpCallback(ctx))
		if err != nil {
			return nil, err
		}
		return &user.StartIdentityProviderFlowResponse{
			Details:  object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderFlowResponse_AuthUrl{AuthUrl: authURL},
		}, nil
	}
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
	provider, err := s.command.GetProvider(ctx, idpID, "")
	if err != nil {
		return nil, "", nil, err
	}
	ldapProvider, ok := provider.(ldap.ProviderInterface)
	if !ok {
		return nil, "", nil, errors.ThrowInvalidArgument(nil, "IDP-9a02j2n2bh", "Errors.ExternalIDP.IDPTypeNotImplemented")
	}

	externalUser, idpSession, err := s.command.LoginWithLDAP(ctx, ldapProvider, username, password)
	if err != nil {
		return nil, "", nil, err
	}

	userID, err := s.checkLinkedExternalUser(ctx, idpID, externalUser.GetID())
	if err != nil {
		return nil, "", nil, err
	}

	ldapSession := idpSession.(*ldap.Session)
	attributes := make(map[string][]string, 0)
	for _, item := range ldapSession.Entry.Attributes {
		attributes[item.Name] = item.Values
	}
	return externalUser, userID, attributes, nil
}

func (s *Server) RetrieveIdentityProviderInformation(ctx context.Context, req *user.RetrieveIdentityProviderInformationRequest) (_ *user.RetrieveIdentityProviderInformationResponse, err error) {
	intent, err := s.command.GetIntentWriteModel(ctx, req.GetIntentId(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	if err := s.checkIntentToken(req.GetToken(), intent.AggregateID); err != nil {
		return nil, err
	}
	if intent.State != domain.IDPIntentStateSucceeded {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-Hk38e", "Errors.Intent.NotSucceeded")
	}
	return intentToIDPInformationPb(intent, s.idpAlg)
}

func intentToIDPInformationPb(intent *command.IDPIntentWriteModel, alg crypto.EncryptionAlgorithm) (_ *user.RetrieveIdentityProviderInformationResponse, err error) {

	rawInformation := new(structpb.Struct)
	err = rawInformation.UnmarshalJSON(intent.IDPUser)
	if err != nil {
		return nil, err
	}

	information := &user.RetrieveIdentityProviderInformationResponse{
		Details: intentToDetailsPb(intent),
		IdpInformation: &user.IDPInformation{
			IdpId:          intent.IDPID,
			UserId:         intent.IDPUserID,
			UserName:       intent.IDPUserName,
			RawInformation: rawInformation,
		},
	}

	if intent.IDPIDToken != "" || intent.IDPAccessToken != nil {
		var idToken *string
		if intent.IDPIDToken != "" {
			idToken = &intent.IDPIDToken
		}
		var accessToken string
		if intent.IDPAccessToken != nil {
			accessToken, err = crypto.DecryptString(intent.IDPAccessToken, alg)
			if err != nil {
				return nil, err
			}
		}
		information.IdpInformation.Access = &user.IDPInformation_Oauth{
			Oauth: &user.IDPOAuthAccessInformation{
				AccessToken: accessToken,
				IdToken:     idToken,
			},
		}
	}

	if intent.IDPEntryAttributes != nil {
		access, err := IDPEntryAttributesToPb(intent.IDPEntryAttributes)
		if err != nil {
			return nil, err
		}
		information.IdpInformation.Access = access
	}

	return information, nil
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
	case domain.UserAuthMethodTypeOTP:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_TOTP
	case domain.UserAuthMethodTypeU2F:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_U2F
	case domain.UserAuthMethodTypePasswordless:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY
	case domain.UserAuthMethodTypePassword:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSWORD
	case domain.UserAuthMethodTypeIDP:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP
	case domain.UserAuthMethodTypeUnspecified:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED
	default:
		return user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_UNSPECIFIED
	}
}
