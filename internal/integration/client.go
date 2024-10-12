package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	feature_v2beta "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	object_v3alpha "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	oidc_pb_v2beta "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	org_v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	user_v3alpha "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
	userschema_v3alpha "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
	webkey_v3alpha "github.com/zitadel/zitadel/pkg/grpc/resources/webkey/v3alpha"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	session_v2beta "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	settings_v2beta "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

type Client struct {
	CC             *grpc.ClientConn
	Admin          admin.AdminServiceClient
	Mgmt           mgmt.ManagementServiceClient
	Auth           auth.AuthServiceClient
	UserV2beta     user_v2beta.UserServiceClient
	UserV2         user_v2.UserServiceClient
	SessionV2beta  session_v2beta.SessionServiceClient
	SessionV2      session.SessionServiceClient
	SettingsV2beta settings_v2beta.SettingsServiceClient
	SettingsV2     settings.SettingsServiceClient
	OIDCv2beta     oidc_pb_v2beta.OIDCServiceClient
	OIDCv2         oidc_pb.OIDCServiceClient
	OrgV2beta      org_v2beta.OrganizationServiceClient
	OrgV2          org.OrganizationServiceClient
	ActionV3Alpha  action.ZITADELActionsClient
	FeatureV2beta  feature_v2beta.FeatureServiceClient
	FeatureV2      feature.FeatureServiceClient
	UserSchemaV3   userschema_v3alpha.ZITADELUserSchemasClient
	WebKeyV3Alpha  webkey_v3alpha.ZITADELWebKeysClient
	IDPv2          idp_pb.IdentityProviderServiceClient
	UserV3Alpha    user_v3alpha.ZITADELUsersClient
}

func newClient(ctx context.Context, target string) (*Client, error) {
	cc, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := &Client{
		CC:             cc,
		Admin:          admin.NewAdminServiceClient(cc),
		Mgmt:           mgmt.NewManagementServiceClient(cc),
		Auth:           auth.NewAuthServiceClient(cc),
		UserV2beta:     user_v2beta.NewUserServiceClient(cc),
		UserV2:         user_v2.NewUserServiceClient(cc),
		SessionV2beta:  session_v2beta.NewSessionServiceClient(cc),
		SessionV2:      session.NewSessionServiceClient(cc),
		SettingsV2beta: settings_v2beta.NewSettingsServiceClient(cc),
		SettingsV2:     settings.NewSettingsServiceClient(cc),
		OIDCv2beta:     oidc_pb_v2beta.NewOIDCServiceClient(cc),
		OIDCv2:         oidc_pb.NewOIDCServiceClient(cc),
		OrgV2beta:      org_v2beta.NewOrganizationServiceClient(cc),
		OrgV2:          org.NewOrganizationServiceClient(cc),
		ActionV3Alpha:  action.NewZITADELActionsClient(cc),
		FeatureV2beta:  feature_v2beta.NewFeatureServiceClient(cc),
		FeatureV2:      feature.NewFeatureServiceClient(cc),
		UserSchemaV3:   userschema_v3alpha.NewZITADELUserSchemasClient(cc),
		WebKeyV3Alpha:  webkey_v3alpha.NewZITADELWebKeysClient(cc),
		IDPv2:          idp_pb.NewIdentityProviderServiceClient(cc),
		UserV3Alpha:    user_v3alpha.NewZITADELUsersClient(cc),
	}
	return client, client.pollHealth(ctx)
}

// pollHealth waits until a healthy status is reported.
func (c *Client) pollHealth(ctx context.Context) (err error) {
	for {
		err = func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			_, err := c.Admin.Healthz(ctx, &admin.HealthzRequest{})
			return err
		}(ctx)
		if err == nil {
			return nil
		}
		logging.WithError(err).Debug("poll healthz")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			continue
		}
	}
}

func (i *Instance) CreateHumanUser(ctx context.Context) *user_v2.AddHumanUserResponse {
	resp, err := i.Client.UserV2.AddHumanUser(ctx, &user_v2.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: i.DefaultOrg.GetId(),
			},
		},
		Profile: &user_v2.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user_v2.Gender_GENDER_MALE),
		},
		Email: &user_v2.SetHumanEmail{
			Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
			Verification: &user_v2.SetHumanEmail_ReturnCode{
				ReturnCode: &user_v2.ReturnEmailVerificationCode{},
			},
		},
		Phone: &user_v2.SetHumanPhone{
			Phone: "+41791234567",
			Verification: &user_v2.SetHumanPhone_ReturnCode{
				ReturnCode: &user_v2.ReturnPhoneVerificationCode{},
			},
		},
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) CreateHumanUserNoPhone(ctx context.Context) *user_v2.AddHumanUserResponse {
	resp, err := i.Client.UserV2.AddHumanUser(ctx, &user_v2.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: i.DefaultOrg.GetId(),
			},
		},
		Profile: &user_v2.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user_v2.Gender_GENDER_MALE),
		},
		Email: &user_v2.SetHumanEmail{
			Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
			Verification: &user_v2.SetHumanEmail_ReturnCode{
				ReturnCode: &user_v2.ReturnEmailVerificationCode{},
			},
		},
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) CreateHumanUserWithTOTP(ctx context.Context, secret string) *user_v2.AddHumanUserResponse {
	resp, err := i.Client.UserV2.AddHumanUser(ctx, &user_v2.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: i.DefaultOrg.GetId(),
			},
		},
		Profile: &user_v2.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user_v2.Gender_GENDER_MALE),
		},
		Email: &user_v2.SetHumanEmail{
			Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
			Verification: &user_v2.SetHumanEmail_ReturnCode{
				ReturnCode: &user_v2.ReturnEmailVerificationCode{},
			},
		},
		Phone: &user_v2.SetHumanPhone{
			Phone: "+41791234567",
			Verification: &user_v2.SetHumanPhone_ReturnCode{
				ReturnCode: &user_v2.ReturnPhoneVerificationCode{},
			},
		},
		TotpSecret: gu.Ptr(secret),
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) CreateOrganization(ctx context.Context, name, adminEmail string) *org.AddOrganizationResponse {
	resp, err := i.Client.OrgV2.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: name,
		Admins: []*org.AddOrganizationRequest_Admin{
			{
				UserType: &org.AddOrganizationRequest_Admin_Human{
					Human: &user_v2.AddHumanUserRequest{
						Profile: &user_v2.SetHumanProfile{
							GivenName:  "firstname",
							FamilyName: "lastname",
						},
						Email: &user_v2.SetHumanEmail{
							Email: adminEmail,
							Verification: &user_v2.SetHumanEmail_ReturnCode{
								ReturnCode: &user_v2.ReturnEmailVerificationCode{},
							},
						},
					},
				},
			},
		},
	})
	logging.OnError(err).Panic("create org")
	return resp
}

func (i *Instance) DeactivateOrganization(ctx context.Context, orgID string) *mgmt.DeactivateOrgResponse {
	resp, err := i.Client.Mgmt.DeactivateOrg(
		SetOrgID(ctx, orgID),
		&mgmt.DeactivateOrgRequest{},
	)
	logging.OnError(err).Fatal("deactivate org")
	return resp
}

func SetOrgID(ctx context.Context, orgID string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.AppendToOutgoingContext(ctx, "x-zitadel-orgid", orgID)
	}
	md.Set("x-zitadel-orgid", orgID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (i *Instance) CreateOrganizationWithUserID(ctx context.Context, name, userID string) *org.AddOrganizationResponse {
	resp, err := i.Client.OrgV2.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: name,
		Admins: []*org.AddOrganizationRequest_Admin{
			{
				UserType: &org.AddOrganizationRequest_Admin_UserId{
					UserId: userID,
				},
			},
		},
	})
	logging.OnError(err).Fatal("create org")
	return resp
}

func (i *Instance) CreateHumanUserVerified(ctx context.Context, org, email string) *user_v2.AddHumanUserResponse {
	resp, err := i.Client.UserV2.AddHumanUser(ctx, &user_v2.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: org,
			},
		},
		Profile: &user_v2.SetHumanProfile{
			GivenName:         "Mickey",
			FamilyName:        "Mouse",
			NickName:          gu.Ptr("Mickey"),
			PreferredLanguage: gu.Ptr("nl"),
			Gender:            gu.Ptr(user_v2.Gender_GENDER_MALE),
		},
		Email: &user_v2.SetHumanEmail{
			Email: email,
			Verification: &user_v2.SetHumanEmail_IsVerified{
				IsVerified: true,
			},
		},
		Phone: &user_v2.SetHumanPhone{
			Phone: "+41791234567",
			Verification: &user_v2.SetHumanPhone_IsVerified{
				IsVerified: true,
			},
		},
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) CreateMachineUser(ctx context.Context) *mgmt.AddMachineUserResponse {
	resp, err := i.Client.Mgmt.AddMachineUser(ctx, &mgmt.AddMachineUserRequest{
		UserName:        fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
		Name:            "Mickey",
		Description:     "Mickey Mouse",
		AccessTokenType: user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER,
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) CreateUserIDPlink(ctx context.Context, userID, externalID, idpID, username string) (*user_v2.AddIDPLinkResponse, error) {
	return i.Client.UserV2.AddIDPLink(
		ctx,
		&user_v2.AddIDPLinkRequest{
			UserId: userID,
			IdpLink: &user_v2.IDPLink{
				IdpId:    idpID,
				UserId:   externalID,
				UserName: username,
			},
		},
	)
}

func (i *Instance) RegisterUserPasskey(ctx context.Context, userID string) {
	reg, err := i.Client.UserV2.CreatePasskeyRegistrationLink(ctx, &user_v2.CreatePasskeyRegistrationLinkRequest{
		UserId: userID,
		Medium: &user_v2.CreatePasskeyRegistrationLinkRequest_ReturnCode{},
	})
	logging.OnError(err).Panic("create user passkey")

	pkr, err := i.Client.UserV2.RegisterPasskey(ctx, &user_v2.RegisterPasskeyRequest{
		UserId: userID,
		Code:   reg.GetCode(),
		Domain: i.Domain,
	})
	logging.OnError(err).Panic("create user passkey")
	attestationResponse, err := i.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	logging.OnError(err).Panic("create user passkey")

	_, err = i.Client.UserV2.VerifyPasskeyRegistration(ctx, &user_v2.VerifyPasskeyRegistrationRequest{
		UserId:              userID,
		PasskeyId:           pkr.GetPasskeyId(),
		PublicKeyCredential: attestationResponse,
		PasskeyName:         "nice name",
	})
	logging.OnError(err).Panic("create user passkey")
}

func (i *Instance) RegisterUserU2F(ctx context.Context, userID string) {
	pkr, err := i.Client.UserV2.RegisterU2F(ctx, &user_v2.RegisterU2FRequest{
		UserId: userID,
		Domain: i.Domain,
	})
	logging.OnError(err).Panic("create user u2f")
	attestationResponse, err := i.WebAuthN.CreateAttestationResponse(pkr.GetPublicKeyCredentialCreationOptions())
	logging.OnError(err).Panic("create user u2f")

	_, err = i.Client.UserV2.VerifyU2FRegistration(ctx, &user_v2.VerifyU2FRegistrationRequest{
		UserId:              userID,
		U2FId:               pkr.GetU2FId(),
		PublicKeyCredential: attestationResponse,
		TokenName:           "nice name",
	})
	logging.OnError(err).Panic("create user u2f")
}

func (i *Instance) SetUserPassword(ctx context.Context, userID, password string, changeRequired bool) *object.Details {
	resp, err := i.Client.UserV2.SetPassword(ctx, &user_v2.SetPasswordRequest{
		UserId: userID,
		NewPassword: &user_v2.Password{
			Password:       password,
			ChangeRequired: changeRequired,
		},
	})
	logging.OnError(err).Panic("set user password")
	return resp.GetDetails()
}

func (i *Instance) AddGenericOAuthProvider(ctx context.Context, name string) *admin.AddGenericOAuthProviderResponse {
	resp, err := i.Client.Admin.AddGenericOAuthProvider(ctx, &admin.AddGenericOAuthProviderRequest{
		Name:                  name,
		ClientId:              "clientID",
		ClientSecret:          "clientSecret",
		AuthorizationEndpoint: "https://example.com/oauth/v2/authorize",
		TokenEndpoint:         "https://example.com/oauth/v2/token",
		UserEndpoint:          "https://api.example.com/user",
		Scopes:                []string{"openid", "profile", "email"},
		IdAttribute:           "id",
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
			AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
		},
	})
	logging.OnError(err).Panic("create generic OAuth idp")

	mustAwait(func() error {
		_, err := i.Client.Admin.GetProviderByID(ctx, &admin.GetProviderByIDRequest{
			Id: resp.GetId(),
		})
		return err
	})

	return resp
}

func (i *Instance) AddOrgGenericOAuthProvider(ctx context.Context, name string) *mgmt.AddGenericOAuthProviderResponse {
	resp, err := i.Client.Mgmt.AddGenericOAuthProvider(ctx, &mgmt.AddGenericOAuthProviderRequest{
		Name:                  name,
		ClientId:              "clientID",
		ClientSecret:          "clientSecret",
		AuthorizationEndpoint: "https://example.com/oauth/v2/authorize",
		TokenEndpoint:         "https://example.com/oauth/v2/token",
		UserEndpoint:          "https://api.example.com/user",
		Scopes:                []string{"openid", "profile", "email"},
		IdAttribute:           "id",
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
			AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
		},
	})
	logging.OnError(err).Panic("create generic OAuth idp")
	/*
		mustAwait(func() error {
			_, err := i.Client.Mgmt.GetProviderByID(ctx, &mgmt.GetProviderByIDRequest{
				Id: resp.GetId(),
			})
			return err
		})
	*/
	return resp
}

func (i *Instance) AddSAMLProvider(ctx context.Context) string {
	resp, err := i.Client.Admin.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name: "saml-idp",
		Metadata: &admin.AddSAMLProviderRequest_MetadataXml{
			MetadataXml: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
		},
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	logging.OnError(err).Panic("create saml idp")
	return resp.GetId()
}

func (i *Instance) AddSAMLRedirectProvider(ctx context.Context, transientMappingAttributeName string) string {
	resp, err := i.Client.Admin.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name:    "saml-idp-redirect",
		Binding: idp.SAMLBinding_SAML_BINDING_REDIRECT,
		Metadata: &admin.AddSAMLProviderRequest_MetadataXml{
			MetadataXml: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
		},
		TransientMappingAttributeName: &transientMappingAttributeName,
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	logging.OnError(err).Panic("create saml idp")
	return resp.GetId()
}

func (i *Instance) AddSAMLPostProvider(ctx context.Context) string {
	resp, err := i.Client.Admin.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name:    "saml-idp-post",
		Binding: idp.SAMLBinding_SAML_BINDING_POST,
		Metadata: &admin.AddSAMLProviderRequest_MetadataXml{
			MetadataXml: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
		},
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	logging.OnError(err).Panic("create saml idp")
	return resp.GetId()
}

/*
func (s *Instance) CreateIntent(t *testing.T, ctx context.Context, idpID string) string {
	resp, err := i.Client.UserV2.StartIdentityProviderIntent(ctx, &user.StartIdentityProviderIntentRequest{
		IdpId: idpID,
		Content: &user.StartIdentityProviderIntentRequest_Urls{
			Urls: &user.RedirectURLs{
				SuccessUrl: "https://example.com/success",
				FailureUrl: "https://example.com/failure",
			},
			AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
		},
	})
	logging.OnError(err).Fatal("create generic OAuth idp")
	return resp
}

func (i *Instance) CreateIntent(t *testing.T, ctx context.Context, idpID string) string {
	ctx = authz.WithInstance(context.WithoutCancel(ctx), s.Instance)
	writeModel, _, err := s.Commands.CreateIntent(ctx, idpID, "https://example.com/success", "https://example.com/failure", s.Instance.InstanceID())
	require.NoError(t, err)
	return writeModel.AggregateID
}

func (i *Instance) CreateSuccessfulOAuthIntent(t *testing.T, ctx context.Context, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
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

func (s *Instance) CreateSuccessfulLDAPIntent(t *testing.T, ctx context.Context, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
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

func (s *Instance) CreateSuccessfulSAMLIntent(t *testing.T, ctx context.Context, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
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
*/

func (i *Instance) CreateVerifiedWebAuthNSession(t *testing.T, ctx context.Context, userID string) (id, token string, start, change time.Time) {
	return i.CreateVerifiedWebAuthNSessionWithLifetime(t, ctx, userID, 0)
}

func (i *Instance) CreateVerifiedWebAuthNSessionWithLifetime(t *testing.T, ctx context.Context, userID string, lifetime time.Duration) (id, token string, start, change time.Time) {
	var sessionLifetime *durationpb.Duration
	if lifetime > 0 {
		sessionLifetime = durationpb.New(lifetime)
	}
	createResp, err := i.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: userID},
			},
		},
		Challenges: &session.RequestChallenges{
			WebAuthN: &session.RequestChallenges_WebAuthN{
				Domain:                      i.Domain,
				UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
			},
		},
		Lifetime: sessionLifetime,
	})
	require.NoError(t, err)

	assertion, err := i.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), true)
	require.NoError(t, err)

	updateResp, err := i.Client.SessionV2.SetSession(ctx, &session.SetSessionRequest{
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

func (i *Instance) CreatePasswordSession(t *testing.T, ctx context.Context, userID, password string) (id, token string, start, change time.Time) {
	createResp, err := i.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
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

func (i *Instance) CreateProjectUserGrant(t *testing.T, ctx context.Context, projectID, userID string) string {
	resp, err := i.Client.Mgmt.AddUserGrant(ctx, &mgmt.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
	require.NoError(t, err)
	return resp.GetUserGrantId()
}

func (i *Instance) CreateOrgMembership(t *testing.T, ctx context.Context, userID string) {
	_, err := i.Client.Mgmt.AddOrgMember(ctx, &mgmt.AddOrgMemberRequest{
		UserId: userID,
		Roles:  []string{domain.RoleOrgOwner},
	})
	require.NoError(t, err)
}

func (i *Instance) CreateProjectMembership(t *testing.T, ctx context.Context, projectID, userID string) {
	_, err := i.Client.Mgmt.AddProjectMember(ctx, &mgmt.AddProjectMemberRequest{
		ProjectId: projectID,
		UserId:    userID,
		Roles:     []string{domain.RoleProjectOwner},
	})
	require.NoError(t, err)
}

func (i *Instance) CreateTarget(ctx context.Context, t *testing.T, name, endpoint string, ty domain.TargetType, interrupt bool) *action.CreateTargetResponse {
	if name == "" {
		name = gofakeit.Name()
	}
	reqTarget := &action.Target{
		Name:     name,
		Endpoint: endpoint,
		Timeout:  durationpb.New(10 * time.Second),
	}
	switch ty {
	case domain.TargetTypeWebhook:
		reqTarget.TargetType = &action.Target_RestWebhook{
			RestWebhook: &action.SetRESTWebhook{
				InterruptOnError: interrupt,
			},
		}
	case domain.TargetTypeCall:
		reqTarget.TargetType = &action.Target_RestCall{
			RestCall: &action.SetRESTCall{
				InterruptOnError: interrupt,
			},
		}
	case domain.TargetTypeAsync:
		reqTarget.TargetType = &action.Target_RestAsync{
			RestAsync: &action.SetRESTAsync{},
		}
	}
	target, err := i.Client.ActionV3Alpha.CreateTarget(ctx, &action.CreateTargetRequest{Target: reqTarget})
	require.NoError(t, err)
	return target
}

func (i *Instance) DeleteExecution(ctx context.Context, t *testing.T, cond *action.Condition) {
	_, err := i.Client.ActionV3Alpha.SetExecution(ctx, &action.SetExecutionRequest{
		Condition: cond,
	})
	require.NoError(t, err)
}

func (i *Instance) SetExecution(ctx context.Context, t *testing.T, cond *action.Condition, targets []*action.ExecutionTargetType) *action.SetExecutionResponse {
	target, err := i.Client.ActionV3Alpha.SetExecution(ctx, &action.SetExecutionRequest{
		Condition: cond,
		Execution: &action.Execution{
			Targets: targets,
		},
	})
	require.NoError(t, err)
	return target
}

func (i *Instance) CreateUserSchemaEmpty(ctx context.Context) *userschema_v3alpha.CreateUserSchemaResponse {
	return i.CreateUserSchemaEmptyWithType(ctx, fmt.Sprint(time.Now().UnixNano()+1))
}

func (i *Instance) CreateUserSchema(ctx context.Context, schemaData []byte) *userschema_v3alpha.CreateUserSchemaResponse {
	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON(schemaData)
	logging.OnError(err).Fatal("create userschema unmarshal")
	schema, err := i.Client.UserSchemaV3.CreateUserSchema(ctx, &userschema_v3alpha.CreateUserSchemaRequest{
		UserSchema: &userschema_v3alpha.UserSchema{
			Type: fmt.Sprint(time.Now().UnixNano() + 1),
			DataType: &userschema_v3alpha.UserSchema_Schema{
				Schema: userSchema,
			},
		},
	})
	logging.OnError(err).Fatal("create userschema")
	return schema
}

func (i *Instance) CreateUserSchemaEmptyWithType(ctx context.Context, schemaType string) *userschema_v3alpha.CreateUserSchemaResponse {
	userSchema := new(structpb.Struct)
	err := userSchema.UnmarshalJSON([]byte(`{
		"$schema": "urn:zitadel:schema:v1",
		"type": "object",
		"properties": {}
	}`))
	logging.OnError(err).Fatal("create userschema unmarshal")
	schema, err := i.Client.UserSchemaV3.CreateUserSchema(ctx, &userschema_v3alpha.CreateUserSchemaRequest{
		UserSchema: &userschema_v3alpha.UserSchema{
			Type: schemaType,
			DataType: &userschema_v3alpha.UserSchema_Schema{
				Schema: userSchema,
			},
		},
	})
	logging.OnError(err).Fatal("create userschema")
	return schema
}

func (i *Instance) CreateSchemaUser(ctx context.Context, orgID string, schemaID string, data []byte) *user_v3alpha.CreateUserResponse {
	userData := new(structpb.Struct)
	err := userData.UnmarshalJSON(data)
	logging.OnError(err).Fatal("create user unmarshal")
	user, err := i.Client.UserV3Alpha.CreateUser(ctx, &user_v3alpha.CreateUserRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		User: &user_v3alpha.CreateUser{
			SchemaId: schemaID,
			Data:     userData,
		},
	})
	logging.OnError(err).Fatal("create user")
	return user
}

func (i *Instance) AddAuthenticatorUsername(ctx context.Context, orgID string, userID string, username string, isOrgSpecific bool) *user_v3alpha.AddUsernameResponse {
	user, err := i.Client.UserV3Alpha.AddUsername(ctx, &user_v3alpha.AddUsernameRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userID,
		Username: &user_v3alpha.SetUsername{
			Username:               username,
			IsOrganizationSpecific: isOrgSpecific,
		},
	})
	logging.OnError(err).Fatal("create username")
	return user
}

func (i *Instance) RemoveAuthenticatorUsername(ctx context.Context, orgID string, userid, id string) *user_v3alpha.RemoveUsernameResponse {
	user, err := i.Client.UserV3Alpha.RemoveUsername(ctx, &user_v3alpha.RemoveUsernameRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userid,
		UsernameId:   id,
	})
	logging.OnError(err).Fatal("remove username")
	return user
}

func (i *Instance) SetAuthenticatorPassword(ctx context.Context, orgID string, userID string, password string) *user_v3alpha.SetPasswordResponse {
	user, err := i.Client.UserV3Alpha.SetPassword(ctx, &user_v3alpha.SetPasswordRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userID,
		NewPassword: &user_v3alpha.SetPassword{
			Type: &user_v3alpha.SetPassword_Password{Password: password},
		},
	})
	logging.OnError(err).Fatal("create password")
	return user
}

func (i *Instance) RequestAuthenticatorPasswordReset(ctx context.Context, orgID string, userID string) *user_v3alpha.RequestPasswordResetResponse {
	user, err := i.Client.UserV3Alpha.RequestPasswordReset(ctx, &user_v3alpha.RequestPasswordResetRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userID,
		Medium:       &user_v3alpha.RequestPasswordResetRequest_ReturnCode{},
	})
	logging.OnError(err).Fatal("reset password")
	return user
}

func (i *Instance) RemoveAuthenticatorPassword(ctx context.Context, orgID string, id string) *user_v3alpha.RemovePasswordResponse {
	user, err := i.Client.UserV3Alpha.RemovePassword(ctx, &user_v3alpha.RemovePasswordRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           id,
	})
	logging.OnError(err).Fatal("remove username")
	return user
}

func (i *Instance) AddAuthenticatorPublicKey(ctx context.Context, orgID string, userID string, pk []byte) *user_v3alpha.AddPublicKeyResponse {
	user, err := i.Client.UserV3Alpha.AddPublicKey(ctx, &user_v3alpha.AddPublicKeyRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userID,
		PublicKey: &user_v3alpha.SetPublicKey{
			Type: &user_v3alpha.SetPublicKey_PublicKey{PublicKey: &user_v3alpha.ProvidedPublicKey{
				PublicKey: pk,
			}},
		},
	})
	logging.OnError(err).Fatal("create pk")
	return user
}

func (i *Instance) RemoveAuthenticatorPublicKey(ctx context.Context, orgID string, userid, id string) *user_v3alpha.RemovePublicKeyResponse {
	user, err := i.Client.UserV3Alpha.RemovePublicKey(ctx, &user_v3alpha.RemovePublicKeyRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userid,
		PublicKeyId:  id,
	})
	logging.OnError(err).Fatal("remove pk")
	return user
}

func (i *Instance) AddAuthenticatorPersonalAccessToken(ctx context.Context, orgID string, userID string) *user_v3alpha.AddPersonalAccessTokenResponse {
	user, err := i.Client.UserV3Alpha.AddPersonalAccessToken(ctx, &user_v3alpha.AddPersonalAccessTokenRequest{
		Organization:        &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:                  userID,
		PersonalAccessToken: &user_v3alpha.SetPersonalAccessToken{},
	})
	logging.OnError(err).Fatal("create pat")
	return user
}

func (i *Instance) RemoveAuthenticatorPersonalAccessToken(ctx context.Context, orgID string, userid, id string) *user_v3alpha.RemovePersonalAccessTokenResponse {
	user, err := i.Client.UserV3Alpha.RemovePersonalAccessToken(ctx, &user_v3alpha.RemovePersonalAccessTokenRequest{
		Organization:          &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:                    userid,
		PersonalAccessTokenId: id,
	})
	logging.OnError(err).Fatal("remove pat")
	return user
}

func (i *Instance) UpdateSchemaUserEmail(ctx context.Context, orgID string, userID string, email string) *user_v3alpha.SetContactEmailResponse {
	user, err := i.Client.UserV3Alpha.SetContactEmail(ctx, &user_v3alpha.SetContactEmailRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userID,
		Email: &user_v3alpha.SetEmail{
			Address:      email,
			Verification: &user_v3alpha.SetEmail_ReturnCode{},
		},
	})
	logging.OnError(err).Fatal("create user")
	return user
}

func (i *Instance) UpdateSchemaUserPhone(ctx context.Context, orgID string, userID string, phone string) *user_v3alpha.SetContactPhoneResponse {
	user, err := i.Client.UserV3Alpha.SetContactPhone(ctx, &user_v3alpha.SetContactPhoneRequest{
		Organization: &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}},
		Id:           userID,
		Phone: &user_v3alpha.SetPhone{
			Number:       phone,
			Verification: &user_v3alpha.SetPhone_ReturnCode{},
		},
	})
	logging.OnError(err).Fatal("create user")
	return user
}

func (i *Instance) CreateInviteCode(ctx context.Context, userID string) *user_v2.CreateInviteCodeResponse {
	user, err := i.Client.UserV2.CreateInviteCode(ctx, &user_v2.CreateInviteCodeRequest{
		UserId:       userID,
		Verification: &user_v2.CreateInviteCodeRequest_ReturnCode{ReturnCode: &user_v2.ReturnInviteCode{}},
	})
	logging.OnError(err).Fatal("create invite code")
	return user
}

func (i *Instance) LockSchemaUser(ctx context.Context, orgID string, userID string) *user_v3alpha.LockUserResponse {
	var org *object_v3alpha.Organization
	if orgID != "" {
		org = &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}}
	}
	user, err := i.Client.UserV3Alpha.LockUser(ctx, &user_v3alpha.LockUserRequest{
		Organization: org,
		Id:           userID,
	})
	logging.OnError(err).Fatal("lock user")
	return user
}

func (i *Instance) UnlockSchemaUser(ctx context.Context, orgID string, userID string) *user_v3alpha.UnlockUserResponse {
	var org *object_v3alpha.Organization
	if orgID != "" {
		org = &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}}
	}
	user, err := i.Client.UserV3Alpha.UnlockUser(ctx, &user_v3alpha.UnlockUserRequest{
		Organization: org,
		Id:           userID,
	})
	logging.OnError(err).Fatal("unlock user")
	return user
}

func (i *Instance) DeactivateSchemaUser(ctx context.Context, orgID string, userID string) *user_v3alpha.DeactivateUserResponse {
	var org *object_v3alpha.Organization
	if orgID != "" {
		org = &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}}
	}
	user, err := i.Client.UserV3Alpha.DeactivateUser(ctx, &user_v3alpha.DeactivateUserRequest{
		Organization: org,
		Id:           userID,
	})
	logging.OnError(err).Fatal("deactivate user")
	return user
}

func (i *Instance) ActivateSchemaUser(ctx context.Context, orgID string, userID string) *user_v3alpha.ActivateUserResponse {
	var org *object_v3alpha.Organization
	if orgID != "" {
		org = &object_v3alpha.Organization{Property: &object_v3alpha.Organization_OrgId{OrgId: orgID}}
	}
	user, err := i.Client.UserV3Alpha.ActivateUser(ctx, &user_v3alpha.ActivateUserRequest{
		Organization: org,
		Id:           userID,
	})
	logging.OnError(err).Fatal("reactivate user")
	return user
}
