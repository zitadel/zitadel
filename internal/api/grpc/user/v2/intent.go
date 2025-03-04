package user

import (
	"context"
	"encoding/json"
	"errors"

	oidc_pkg "github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/apple"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) StartIdentityProviderIntent(ctx context.Context, req *user.StartIdentityProviderIntentRequest) (_ *user.StartIdentityProviderIntentResponse, err error) {
	switch t := req.GetContent().(type) {
	case *user.StartIdentityProviderIntentRequest_Urls:
		return s.startIDPIntent(ctx, req.GetIdpId(), t.Urls)
	case *user.StartIdentityProviderIntentRequest_Ldap:
		return s.startLDAPIntent(ctx, req.GetIdpId(), t.Ldap)
	default:
		return nil, zerrors.ThrowUnimplementedf(nil, "USERv2-S2g21", "type oneOf %T in method StartIdentityProviderIntent not implemented", t)
	}
}

func (s *Server) startIDPIntent(ctx context.Context, idpID string, urls *user.RedirectURLs) (*user.StartIdentityProviderIntentResponse, error) {
	state, session, err := s.command.AuthFromProvider(ctx, idpID, s.idpCallback(ctx), s.samlRootURL(ctx, idpID))
	if err != nil {
		return nil, err
	}
	_, details, err := s.command.CreateIntent(ctx, state, idpID, urls.GetSuccessUrl(), urls.GetFailureUrl(), authz.GetInstance(ctx).InstanceID(), session.PersistentParameters())
	if err != nil {
		return nil, err
	}
	content, redirect := session.GetAuth(ctx)
	if redirect {
		return &user.StartIdentityProviderIntentResponse{
			Details:  object.DomainToDetailsPb(details),
			NextStep: &user.StartIdentityProviderIntentResponse_AuthUrl{AuthUrl: content},
		}, nil
	}
	return &user.StartIdentityProviderIntentResponse{
		Details: object.DomainToDetailsPb(details),
		NextStep: &user.StartIdentityProviderIntentResponse_PostForm{
			PostForm: []byte(content),
		},
	}, nil
}

func (s *Server) startLDAPIntent(ctx context.Context, idpID string, ldapCredentials *user.LDAPCredentials) (*user.StartIdentityProviderIntentResponse, error) {
	intentWriteModel, details, err := s.command.CreateIntent(ctx, "", idpID, "", "", authz.GetInstance(ctx).InstanceID(), nil)
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
	links, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: queries}, nil)
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
	return externalUser, userID, attributes, nil
}

func (s *Server) RetrieveIdentityProviderIntent(ctx context.Context, req *user.RetrieveIdentityProviderIntentRequest) (_ *user.RetrieveIdentityProviderIntentResponse, err error) {
	intent, err := s.command.GetIntentWriteModel(ctx, req.GetIdpIntentId(), "")
	if err != nil {
		return nil, err
	}
	if err := s.checkIntentToken(req.GetIdpIntentToken(), intent.AggregateID); err != nil {
		return nil, err
	}
	if intent.State != domain.IDPIntentStateSucceeded {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDP-nme4gszsvx", "Errors.Intent.NotSucceeded")
	}
	idpIntent, err := idpIntentToIDPIntentPb(intent, s.idpAlg)
	if err != nil {
		return nil, err
	}
	if idpIntent.UserId == "" {
		provider, err := s.command.GetProvider(ctx, idpIntent.IdpInformation.IdpId, "", "")
		if err != nil && !errors.Is(err, oidc_pkg.ErrDiscoveryFailed) {
			return nil, err
		}
		var idpUser idp.User
		switch p := provider.(type) {
		case *apple.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &apple.User{})
		case *oauth.Provider:
			idpUser, err = unmarshalRawIdpUser(intent.IDPUser, p.User())
		case *oidc.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &oidc.User{UserInfo: &oidc_pkg.UserInfo{}})
		case *jwt.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &jwt.User{})
		case *azuread.Provider:
			idpUser, err = unmarshalRawIdpUser(intent.IDPUser, p.User())
		case *github.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &github.User{})
		case *gitlab.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &oidc.User{UserInfo: &oidc_pkg.UserInfo{}})
		case *google.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &oidc.User{UserInfo: &oidc_pkg.UserInfo{}})
		case *saml.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &saml.UserMapper{})
		case *ldap.Provider:
			idpUser, err = unmarshalIdpUser(intent.IDPUser, &ldap.User{})
		default:
			return nil, zerrors.ThrowInvalidArgument(nil, "IDP-7rPBbls4Zn", "Errors.ExternalIDP.IDPTypeNotImplemented")
		}
		if err != nil {
			return nil, err
		}
		idpIntent.AddHumanUser = idpUserToAddHumanUser(idpUser, idpIntent.IdpInformation.IdpId)
	}
	return idpIntent, nil
}

type rawUserMapper struct {
	RawInfo map[string]interface{}
}

func unmarshalRawIdpUser(idpUserData []byte, idpUser idp.User) (idp.User, error) {
	userMapper := &rawUserMapper{}
	if err := json.Unmarshal(idpUserData, userMapper); err != nil {
		return nil, err
	}
	idpUserData, err := json.Marshal(userMapper.RawInfo)
	if err != nil {
		return nil, err
	}
	return unmarshalIdpUser(idpUserData, idpUser)
}

func unmarshalIdpUser(idpUserData []byte, idpUser idp.User) (idp.User, error) {
	if err := json.Unmarshal(idpUserData, idpUser); err != nil {
		return nil, err
	}
	return idpUser, nil
}

func idpIntentToIDPIntentPb(intent *command.IDPIntentWriteModel, alg crypto.EncryptionAlgorithm) (_ *user.RetrieveIdentityProviderIntentResponse, err error) {
	rawInformation := new(structpb.Struct)
	err = rawInformation.UnmarshalJSON(intent.IDPUser)
	if err != nil {
		return nil, err
	}
	information := &user.RetrieveIdentityProviderIntentResponse{
		IdpInformation: &user.IDPInformation{
			IdpId:          intent.IDPID,
			UserId:         intent.IDPUserID,
			UserName:       intent.IDPUserName,
			RawInformation: rawInformation,
		},
		UserId: intent.UserID,
	}
	information.Details = intentToDetailsPb(intent)
	// OAuth / OIDC
	if intent.IDPIDToken != "" || intent.IDPAccessToken != nil {
		information.IdpInformation.Access, err = idpOAuthTokensToPb(intent.IDPIDToken, intent.IDPAccessToken, alg)
		if err != nil {
			return nil, err
		}
	}
	// LDAP
	if intent.IDPEntryAttributes != nil {
		access, err := IDPEntryAttributesToPb(intent.IDPEntryAttributes)
		if err != nil {
			return nil, err
		}
		information.IdpInformation.Access = access
	}
	// SAML
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

func idpUserToAddHumanUser(idpUser idp.User, idpID string) *user.AddHumanUserRequest {
	addHumanUser := &user.AddHumanUserRequest{
		Profile: &user.SetHumanProfile{
			GivenName:  idpUser.GetFirstName(),
			FamilyName: idpUser.GetLastName(),
		},
		Email: &user.SetHumanEmail{
			Email:        string(idpUser.GetEmail()),
			Verification: &user.SetHumanEmail_SendCode{},
		},
		Metadata: make([]*user.SetMetadataEntry, 0),
		IdpLinks: []*user.IDPLink{
			{
				IdpId:    idpID,
				UserId:   idpUser.GetID(),
				UserName: idpUser.GetPreferredUsername(),
			},
		},
	}
	if username := idpUser.GetPreferredUsername(); username != "" {
		addHumanUser.Username = &username
	}
	if nickName := idpUser.GetNickname(); nickName != "" {
		addHumanUser.Profile.NickName = &nickName
	}
	if displayName := idpUser.GetDisplayName(); displayName != "" {
		addHumanUser.Profile.DisplayName = &displayName
	}
	if lang := idpUser.GetPreferredLanguage().String(); lang != "" {
		addHumanUser.Profile.PreferredLanguage = &lang
	}
	if isEmailVerified := idpUser.IsEmailVerified(); isEmailVerified {
		addHumanUser.Email.Verification = &user.SetHumanEmail_IsVerified{IsVerified: isEmailVerified}
	}
	if phone := idpUser.GetPhone(); phone != "" {
		addHumanUser.Phone = &user.SetHumanPhone{
			Phone:        string(phone),
			Verification: &user.SetHumanPhone_SendCode{},
		}
		if isPhoneVerified := idpUser.IsPhoneVerified(); isPhoneVerified {
			addHumanUser.Phone.Verification = &user.SetHumanPhone_IsVerified{IsVerified: isPhoneVerified}
		}
	}
	return addHumanUser
}
