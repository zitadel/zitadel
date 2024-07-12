package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	crewjam_saml "github.com/crewjam/saml"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/repository/idp"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	feature_v2beta "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	oidc_pb_v2beta "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	organisation "github.com/zitadel/zitadel/pkg/grpc/org/v2"
	org_v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	session_v2beta "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	settings_v2beta "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	schema "github.com/zitadel/zitadel/pkg/grpc/user/schema/v3alpha"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

type Client struct {
	CC             *grpc.ClientConn
	Admin          admin.AdminServiceClient
	Mgmt           mgmt.ManagementServiceClient
	Auth           auth.AuthServiceClient
	UserV2beta     user_v2beta.UserServiceClient
	UserV2         user.UserServiceClient
	SessionV2beta  session_v2beta.SessionServiceClient
	SessionV2      session.SessionServiceClient
	SettingsV2beta settings_v2beta.SettingsServiceClient
	SettingsV2     settings.SettingsServiceClient
	OIDCv2beta     oidc_pb_v2beta.OIDCServiceClient
	OIDCv2         oidc_pb.OIDCServiceClient
	OrgV2beta      org_v2beta.OrganizationServiceClient
	OrgV2          organisation.OrganizationServiceClient
	System         system.SystemServiceClient
	ActionV3       action.ActionServiceClient
	FeatureV2beta  feature_v2beta.FeatureServiceClient
	FeatureV2      feature.FeatureServiceClient
	UserSchemaV3   schema.UserSchemaServiceClient
}

func newClient(cc *grpc.ClientConn) Client {
	return Client{
		CC:             cc,
		Admin:          admin.NewAdminServiceClient(cc),
		Mgmt:           mgmt.NewManagementServiceClient(cc),
		Auth:           auth.NewAuthServiceClient(cc),
		UserV2beta:     user_v2beta.NewUserServiceClient(cc),
		UserV2:         user.NewUserServiceClient(cc),
		SessionV2beta:  session_v2beta.NewSessionServiceClient(cc),
		SessionV2:      session.NewSessionServiceClient(cc),
		SettingsV2beta: settings_v2beta.NewSettingsServiceClient(cc),
		SettingsV2:     settings.NewSettingsServiceClient(cc),
		OIDCv2beta:     oidc_pb_v2beta.NewOIDCServiceClient(cc),
		OIDCv2:         oidc_pb.NewOIDCServiceClient(cc),
		OrgV2beta:      org_v2beta.NewOrganizationServiceClient(cc),
		OrgV2:          organisation.NewOrganizationServiceClient(cc),
		System:         system.NewSystemServiceClient(cc),
		ActionV3:       action.NewActionServiceClient(cc),
		FeatureV2beta:  feature_v2beta.NewFeatureServiceClient(cc),
		FeatureV2:      feature.NewFeatureServiceClient(cc),
		UserSchemaV3:   schema.NewUserSchemaServiceClient(cc),
	}
}

func (t *Tester) UseIsolatedInstance(tt *testing.T, iamOwnerCtx, systemCtx context.Context) (primaryDomain, instanceId, adminID string, authenticatedIamOwnerCtx context.Context) {
	primaryDomain = RandString(5) + ".integration.localhost"
	instance, err := t.Client.System.CreateInstance(systemCtx, &system.CreateInstanceRequest{
		InstanceName: "testinstance",
		CustomDomain: primaryDomain,
		Owner: &system.CreateInstanceRequest_Machine_{
			Machine: &system.CreateInstanceRequest_Machine{
				UserName:            "owner",
				Name:                "owner",
				PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
			},
		},
	})
	require.NoError(tt, err)
	t.createClientConn(iamOwnerCtx, fmt.Sprintf("%s:%d", primaryDomain, t.Config.Port))
	instanceId = instance.GetInstanceId()
	owner, err := t.Queries.GetUserByLoginName(authz.WithInstanceID(iamOwnerCtx, instanceId), true, "owner@"+primaryDomain)
	require.NoError(tt, err)
	t.Users.Set(instanceId, IAMOwner, &User{
		User:  owner,
		Token: instance.GetPat(),
	})
	newCtx := t.WithInstanceAuthorization(iamOwnerCtx, IAMOwner, instanceId)
	var adminUser *mgmt.ImportHumanUserResponse
	// the following serves two purposes:
	// 1. it ensures that the instance is ready to be used
	// 2. it enables a normal login with the default admin user credentials
	require.EventuallyWithT(tt, func(collectT *assert.CollectT) {
		var importErr error
		adminUser, importErr = t.Client.Mgmt.ImportHumanUser(newCtx, &mgmt.ImportHumanUserRequest{
			UserName: "zitadel-admin@zitadel.localhost",
			Email: &mgmt.ImportHumanUserRequest_Email{
				Email:           "zitadel-admin@zitadel.localhost",
				IsEmailVerified: true,
			},
			Password: "Password1!",
			Profile: &mgmt.ImportHumanUserRequest_Profile{
				FirstName: "hodor",
				LastName:  "hodor",
				NickName:  "hodor",
			},
		})
		assert.NoError(collectT, importErr)
	}, 2*time.Minute, 100*time.Millisecond, "instance not ready")
	return primaryDomain, instanceId, adminUser.GetUserId(), newCtx
}

func (s *Tester) CreateHumanUser(ctx context.Context) *user.AddHumanUserResponse {
	resp, err := s.Client.UserV2.AddHumanUser(ctx, &user.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: s.Organisation.ID,
			},
		},
		Profile: &user.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user.Gender_GENDER_MALE),
		},
		Email: &user.SetHumanEmail{
			Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
			Verification: &user.SetHumanEmail_ReturnCode{
				ReturnCode: &user.ReturnEmailVerificationCode{},
			},
		},
		Phone: &user.SetHumanPhone{
			Phone: "+41791234567",
			Verification: &user.SetHumanPhone_ReturnCode{
				ReturnCode: &user.ReturnPhoneVerificationCode{},
			},
		},
	})
	logging.OnError(err).Fatal("create human user")
	return resp
}

func (s *Tester) CreateHumanUserWithTOTP(ctx context.Context, secret string) *user.AddHumanUserResponse {
	resp, err := s.Client.UserV2.AddHumanUser(ctx, &user.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: s.Organisation.ID,
			},
		},
		Profile: &user.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user.Gender_GENDER_MALE),
		},
		Email: &user.SetHumanEmail{
			Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
			Verification: &user.SetHumanEmail_ReturnCode{
				ReturnCode: &user.ReturnEmailVerificationCode{},
			},
		},
		Phone: &user.SetHumanPhone{
			Phone: "+41791234567",
			Verification: &user.SetHumanPhone_ReturnCode{
				ReturnCode: &user.ReturnPhoneVerificationCode{},
			},
		},
		TotpSecret: gu.Ptr(secret),
	})
	logging.OnError(err).Fatal("create human user")
	return resp
}

func (s *Tester) CreateOrganization(ctx context.Context, name, adminEmail string) *org.AddOrganizationResponse {
	resp, err := s.Client.OrgV2.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: name,
		Admins: []*org.AddOrganizationRequest_Admin{
			{
				UserType: &org.AddOrganizationRequest_Admin_Human{
					Human: &user.AddHumanUserRequest{
						Profile: &user.SetHumanProfile{
							GivenName:  "firstname",
							FamilyName: "lastname",
						},
						Email: &user.SetHumanEmail{
							Email: adminEmail,
							Verification: &user.SetHumanEmail_ReturnCode{
								ReturnCode: &user.ReturnEmailVerificationCode{},
							},
						},
					},
				},
			},
		},
	})
	logging.OnError(err).Fatal("create org")
	return resp
}

func (s *Tester) CreateHumanUserVerified(ctx context.Context, org, email string) *user.AddHumanUserResponse {
	resp, err := s.Client.UserV2.AddHumanUser(ctx, &user.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: org,
			},
		},
		Profile: &user.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			NickName:          gu.Ptr("Mickey"),
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user.Gender_GENDER_MALE),
		},
		Email: &user.SetHumanEmail{
			Email: email,
			Verification: &user.SetHumanEmail_IsVerified{
				IsVerified: true,
			},
		},
		Phone: &user.SetHumanPhone{
			Phone: "+41791234567",
			Verification: &user.SetHumanPhone_IsVerified{
				IsVerified: true,
			},
		},
	})
	logging.OnError(err).Fatal("create human user")
	return resp
}

func (s *Tester) CreateMachineUser(ctx context.Context) *mgmt.AddMachineUserResponse {
	resp, err := s.Client.Mgmt.AddMachineUser(ctx, &mgmt.AddMachineUserRequest{
		UserName:        fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
		Name:            "Mickey",
		Description:     "Mickey Mouse",
		AccessTokenType: user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
	})
	logging.OnError(err).Fatal("create human user")
	return resp
}

func (s *Tester) CreateUserIDPlink(ctx context.Context, userID, externalID, idpID, username string) *user.AddIDPLinkResponse {
	resp, err := s.Client.UserV2.AddIDPLink(
		ctx,
		&user.AddIDPLinkRequest{
			UserId: userID,
			IdpLink: &user.IDPLink{
				IdpId:    idpID,
				UserId:   externalID,
				UserName: username,
			},
		},
	)
	logging.OnError(err).Fatal("create human user link")
	return resp
}

func (s *Tester) RegisterUserPasskey(ctx context.Context, userID string) {
	reg, err := s.Client.UserV2.CreatePasskeyRegistrationLink(ctx, &user.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	logging.OnError(err).Fatal("create user passkey")

	pkr, err := s.Client.UserV2.RegisterPasskey(ctx, &user.RegisterPasskeyRequest{
		UserId: userID,
		Code:   reg.GetCode(),
		Domain: s.Config.ExternalDomain,
	})
	logging.OnError(err).Fatal("create user passkey")
	attestationResponse, err := s.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	logging.OnError(err).Fatal("create user passkey")

	_, err = s.Client.UserV2.VerifyPasskeyRegistration(ctx, &user.VerifyPasskeyRegistrationRequest{
		UserId:              userID,
		PasskeyId:           pkr.GetPasskeyId(),
		PublicKeyCredential: attestationResponse,
		PasskeyName:         "nice name",
	})
	logging.OnError(err).Fatal("create user passkey")
}

func (s *Tester) RegisterUserU2F(ctx context.Context, userID string) {
	pkr, err := s.Client.UserV2.RegisterU2F(ctx, &user.RegisterU2FRequest{
		UserId: userID,
		Domain: s.Config.ExternalDomain,
	})
	logging.OnError(err).Fatal("create user u2f")
	attestationResponse, err := s.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	logging.OnError(err).Fatal("create user u2f")

	_, err = s.Client.UserV2.VerifyU2FRegistration(ctx, &user.VerifyU2FRegistrationRequest{
		UserId:              userID,
		U2FId:               pkr.GetU2FId(),
		PublicKeyCredential: attestationResponse,
		TokenName:           "nice name",
	})
	logging.OnError(err).Fatal("create user u2f")
}

func (s *Tester) SetUserPassword(ctx context.Context, userID, password string, changeRequired bool) *object.Details {
	resp, err := s.Client.UserV2.SetPassword(ctx, &user.SetPasswordRequest{
		UserId: userID,
		NewPassword: &user.Password{
			Password:       password,
			ChangeRequired: changeRequired,
		},
	})
	logging.OnError(err).Fatal("set user password")
	return resp.GetDetails()
}

func (s *Tester) AddGenericOAuthProvider(t *testing.T, ctx context.Context) string {
	ctx = authz.WithInstance(ctx, s.Instance)
	id, _, err := s.Commands.AddInstanceGenericOAuthProvider(ctx, command.GenericOAuthProvider{
		Name:                  "idp",
		ClientID:              "clientID",
		ClientSecret:          "clientSecret",
		AuthorizationEndpoint: "https://example.com/oauth/v2/authorize",
		TokenEndpoint:         "https://example.com/oauth/v2/token",
		UserEndpoint:          "https://api.example.com/user",
		Scopes:                []string{"openid", "profile", "email"},
		IDAttribute:           "id",
		IDPOptions: idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	require.NoError(t, err)
	return id
}

func (s *Tester) AddOrgGenericOAuthProvider(t *testing.T, ctx context.Context, orgID string) string {
	ctx = authz.WithInstance(ctx, s.Instance)
	id, _, err := s.Commands.AddOrgGenericOAuthProvider(ctx, orgID,
		command.GenericOAuthProvider{
			Name:                  "idp",
			ClientID:              "clientID",
			ClientSecret:          "clientSecret",
			AuthorizationEndpoint: "https://example.com/oauth/v2/authorize",
			TokenEndpoint:         "https://example.com/oauth/v2/token",
			UserEndpoint:          "https://api.example.com/user",
			Scopes:                []string{"openid", "profile", "email"},
			IDAttribute:           "id",
			IDPOptions: idp.Options{
				IsLinkingAllowed:  true,
				IsCreationAllowed: true,
				IsAutoCreation:    true,
				IsAutoUpdate:      true,
			},
		})
	require.NoError(t, err)
	return id
}

func (s *Tester) AddSAMLProvider(t *testing.T, ctx context.Context) string {
	ctx = authz.WithInstance(ctx, s.Instance)
	id, _, err := s.Server.Commands.AddInstanceSAMLProvider(ctx, command.SAMLProvider{
		Name:     "saml-idp",
		Metadata: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
		IDPOptions: idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	require.NoError(t, err)
	return id
}

func (s *Tester) AddSAMLRedirectProvider(t *testing.T, ctx context.Context, transientMappingAttributeName string) string {
	ctx = authz.WithInstance(ctx, s.Instance)
	id, _, err := s.Server.Commands.AddInstanceSAMLProvider(ctx, command.SAMLProvider{
		Name:                          "saml-idp-redirect",
		Binding:                       "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
		Metadata:                      []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
		TransientMappingAttributeName: transientMappingAttributeName,
		IDPOptions: idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	require.NoError(t, err)
	return id
}

func (s *Tester) AddSAMLPostProvider(t *testing.T, ctx context.Context) string {
	ctx = authz.WithInstance(ctx, s.Instance)
	id, _, err := s.Server.Commands.AddInstanceSAMLProvider(ctx, command.SAMLProvider{
		Name:     "saml-idp-post",
		Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
		Metadata: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
		IDPOptions: idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	require.NoError(t, err)
	return id
}

func (s *Tester) CreateIntent(t *testing.T, ctx context.Context, idpID string) string {
	ctx = authz.WithInstance(context.WithoutCancel(ctx), s.Instance)
	writeModel, _, err := s.Commands.CreateIntent(ctx, idpID, "https://example.com/success", "https://example.com/failure", s.Instance.InstanceID())
	require.NoError(t, err)
	return writeModel.AggregateID
}

func (s *Tester) CreateSuccessfulOAuthIntent(t *testing.T, ctx context.Context, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
	ctx = authz.WithInstance(context.WithoutCancel(ctx), s.Instance)
	intentID := s.CreateIntent(t, ctx, idpID)
	writeModel, err := s.Commands.GetIntentWriteModel(ctx, intentID, s.Instance.InstanceID())
	require.NoError(t, err)
	idpUser := openid.NewUser(
		&oidc.UserInfo{
			Subject: idpUserID,
			UserInfoProfile: oidc.UserInfoProfile{
				PreferredUsername: "username",
			},
		},
	)
	idpSession := &openid.Session{
		Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
			Token: &oauth2.Token{
				AccessToken: "accessToken",
			},
			IDToken: "idToken",
		},
	}
	token, err := s.Commands.SucceedIDPIntent(ctx, writeModel, idpUser, idpSession, userID)
	require.NoError(t, err)
	return intentID, token, writeModel.ChangeDate, writeModel.ProcessedSequence
}

func (s *Tester) CreateSuccessfulLDAPIntent(t *testing.T, ctx context.Context, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
	ctx = authz.WithInstance(context.WithoutCancel(ctx), s.Instance)
	intentID := s.CreateIntent(t, ctx, idpID)
	writeModel, err := s.Commands.GetIntentWriteModel(ctx, intentID, s.Instance.InstanceID())
	require.NoError(t, err)
	username := "username"
	lang := language.Make("en")
	idpUser := ldap.NewUser(
		idpUserID,
		"",
		"",
		"",
		"",
		username,
		"",
		false,
		"",
		false,
		lang,
		"",
		"",
	)
	attributes := map[string][]string{"id": {idpUserID}, "username": {username}, "language": {lang.String()}}
	token, err := s.Commands.SucceedLDAPIDPIntent(ctx, writeModel, idpUser, userID, attributes)
	require.NoError(t, err)
	return intentID, token, writeModel.ChangeDate, writeModel.ProcessedSequence
}

func (s *Tester) CreateSuccessfulSAMLIntent(t *testing.T, ctx context.Context, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
	ctx = authz.WithInstance(context.WithoutCancel(ctx), s.Instance)
	intentID := s.CreateIntent(t, ctx, idpID)
	writeModel, err := s.Server.Commands.GetIntentWriteModel(ctx, intentID, s.Instance.InstanceID())
	require.NoError(t, err)

	idpUser := &saml.UserMapper{
		ID:         idpUserID,
		Attributes: map[string][]string{"attribute1": {"value1"}},
	}
	assertion := &crewjam_saml.Assertion{ID: "id"}

	token, err := s.Server.Commands.SucceedSAMLIDPIntent(ctx, writeModel, idpUser, userID, assertion)
	require.NoError(t, err)
	return intentID, token, writeModel.ChangeDate, writeModel.ProcessedSequence
}

func (s *Tester) CreateVerifiedWebAuthNSession(t *testing.T, ctx context.Context, userID string) (id, token string, start, change time.Time) {
	return s.CreateVerifiedWebAuthNSessionWithLifetime(t, ctx, userID, 0)
}

func (s *Tester) CreateVerifiedWebAuthNSessionWithLifetime(t *testing.T, ctx context.Context, userID string, lifetime time.Duration) (id, token string, start, change time.Time) {
	var sessionLifetime *durationpb.Duration
	if lifetime > 0 {
		sessionLifetime = durationpb.New(lifetime)
	}
	createResp, err := s.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: userID},
			},
		},
		Challenges: &session.RequestChallenges{
			WebAuthN: &session.RequestChallenges_WebAuthN{
				Domain:                      s.Config.ExternalDomain,
				UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
			},
		},
		Lifetime: sessionLifetime,
	})
	require.NoError(t, err)

	assertion, err := s.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), true)
	require.NoError(t, err)

	updateResp, err := s.Client.SessionV2.SetSession(ctx, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			WebAuthN: &session.CheckWebAuthN{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)
	return createResp.GetSessionId(), updateResp.GetSessionToken(),
		createResp.GetDetails().GetChangeDate().AsTime(), updateResp.GetDetails().GetChangeDate().AsTime()
}

func (s *Tester) CreatePasswordSession(t *testing.T, ctx context.Context, userID, password string) (id, token string, start, change time.Time) {
	createResp, err := s.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: userID},
			},
			Password: &session.CheckPassword{
				Password: password,
			},
		},
	})
	require.NoError(t, err)
	return createResp.GetSessionId(), createResp.GetSessionToken(),
		createResp.GetDetails().GetChangeDate().AsTime(), createResp.GetDetails().GetChangeDate().AsTime()
}

func (s *Tester) CreateProjectUserGrant(t *testing.T, ctx context.Context, projectID, userID string) string {
	resp, err := s.Client.Mgmt.AddUserGrant(ctx, &mgmt.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
	require.NoError(t, err)
	return resp.GetUserGrantId()
}

func (s *Tester) CreateOrgMembership(t *testing.T, ctx context.Context, userID string) {
	_, err := s.Client.Mgmt.AddOrgMember(ctx, &mgmt.AddOrgMemberRequest{
		UserId: userID,
		Roles:  []string{domain.RoleOrgOwner},
	})
	require.NoError(t, err)
}

func (s *Tester) CreateProjectMembership(t *testing.T, ctx context.Context, projectID, userID string) {
	_, err := s.Client.Mgmt.AddProjectMember(ctx, &mgmt.AddProjectMemberRequest{
		ProjectId: projectID,
		UserId:    userID,
		Roles:     []string{domain.RoleProjectOwner},
	})
	require.NoError(t, err)
}

func (s *Tester) CreateTarget(ctx context.Context, t *testing.T, name, endpoint string, ty domain.TargetType, interrupt bool) *action.CreateTargetResponse {
	nameSet := fmt.Sprint(time.Now().UnixNano() + 1)
	if name != "" {
		nameSet = name
	}
	req := &action.CreateTargetRequest{
		Name:     nameSet,
		Endpoint: endpoint,
		Timeout:  durationpb.New(10 * time.Second),
	}
	switch ty {
	case domain.TargetTypeWebhook:
		req.TargetType = &action.CreateTargetRequest_RestWebhook{
			RestWebhook: &action.SetRESTWebhook{
				InterruptOnError: interrupt,
			},
		}
	case domain.TargetTypeCall:
		req.TargetType = &action.CreateTargetRequest_RestCall{
			RestCall: &action.SetRESTCall{
				InterruptOnError: interrupt,
			},
		}
	case domain.TargetTypeAsync:
		req.TargetType = &action.CreateTargetRequest_RestAsync{
			RestAsync: &action.SetRESTAsync{},
		}
	}
	target, err := s.Client.ActionV3.CreateTarget(ctx, req)
	require.NoError(t, err)
	return target
}

func (s *Tester) SetExecution(ctx context.Context, t *testing.T, cond *action.Condition, targets []*action.ExecutionTargetType) *action.SetExecutionResponse {
	target, err := s.Client.ActionV3.SetExecution(ctx, &action.SetExecutionRequest{
		Condition: cond,
		Targets:   targets,
	})
	require.NoError(t, err)
	return target
}

func (s *Tester) DeleteExecution(ctx context.Context, t *testing.T, cond *action.Condition) {
	_, err := s.Client.ActionV3.DeleteExecution(ctx, &action.DeleteExecutionRequest{
		Condition: cond,
	})
	require.NoError(t, err)
}

func (s *Tester) CreateUserSchema(ctx context.Context, t *testing.T) *schema.CreateUserSchemaResponse {
	return s.CreateUserSchemaWithType(ctx, t, fmt.Sprint(time.Now().UnixNano()+1))
}

func (s *Tester) CreateUserSchemaWithType(ctx context.Context, t *testing.T, schemaType string) *schema.CreateUserSchemaResponse {
	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON([]byte(`{
		"$schema": "urn:zitadel:schema:v1",
		"type": "object",
		"properties": {}
	}`))
	require.NoError(t, err)
	target, err := s.Client.UserSchemaV3.CreateUserSchema(ctx, &schema.CreateUserSchemaRequest{
		Type: schemaType,
		DataType: &schema.CreateUserSchemaRequest_Schema{
			Schema: userSchema,
		},
	})
	require.NoError(t, err)
	return target
}
