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
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
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

func (s *Server) StartIdentityProviderIntent(ctx context.Context, req *user.StartIdentityProviderIntentRequest) (_ *user.StartIdentityProviderIntentResponse, err error) {
	id, details, err := s.command.CreateIntent(ctx, req.GetIdpId(), req.GetSuccessUrl(), req.GetFailureUrl(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	authURL, err := s.command.AuthURLFromProvider(ctx, req.GetIdpId(), id, s.idpCallback(ctx))
	if err != nil {
		return nil, err
	}
	return &user.StartIdentityProviderIntentResponse{
		Details:  object.DomainToDetailsPb(details),
		NextStep: &user.StartIdentityProviderIntentResponse_AuthUrl{AuthUrl: authURL},
	}, nil
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
	rawInformation := new(structpb.Struct)
	err = rawInformation.UnmarshalJSON(intent.IDPUser)
	if err != nil {
		return nil, err
	}

	return &user.RetrieveIdentityProviderIntentResponse{
		Details: &object_pb.Details{
			Sequence:      intent.ProcessedSequence,
			ChangeDate:    timestamppb.New(intent.ChangeDate),
			ResourceOwner: intent.ResourceOwner,
		},
		IdpInformation: &user.IDPInformation{
			Access: &user.IDPInformation_Oauth{
				Oauth: &user.IDPOAuthAccessInformation{
					AccessToken: accessToken,
					IdToken:     idToken,
				},
			},
			IdpId:          intent.IDPID,
			UserId:         intent.IDPUserID,
			UserName:       intent.IDPUserName,
			RawInformation: rawInformation,
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
