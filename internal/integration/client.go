package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	authorization_v2 "github.com/zitadel/zitadel/pkg/grpc/authorization/v2"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	internal_permission_v2 "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	project_v2 "github.com/zitadel/zitadel/pkg/grpc/project/v2"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
	webkey_v2 "github.com/zitadel/zitadel/pkg/grpc/webkey/v2"
)

type Client struct {
	CC                   *grpc.ClientConn
	Admin                admin.AdminServiceClient
	Mgmt                 mgmt.ManagementServiceClient
	Auth                 auth.AuthServiceClient
	UserV2               user_v2.UserServiceClient
	SessionV2            session.SessionServiceClient
	SettingsV2           settings.SettingsServiceClient
	OIDCv2               oidc_pb.OIDCServiceClient
	OrgV2                org.OrganizationServiceClient
	ActionV2             action.ActionServiceClient
	FeatureV2            feature.FeatureServiceClient
	WebKeyV2             webkey_v2.WebKeyServiceClient
	IDPv2                idp_pb.IdentityProviderServiceClient
	SAMLv2               saml_pb.SAMLServiceClient
	SCIM                 *scim.Client
	ProjectV2            project_v2.ProjectServiceClient
	InstanceV2           instance_v2.InstanceServiceClient
	ApplicationV2        application.ApplicationServiceClient
	InternalPermissionV2 internal_permission_v2.InternalPermissionServiceClient
	AuthorizationV2      authorization_v2.AuthorizationServiceClient
	GroupV2              group_v2.GroupServiceClient
}

func NewDefaultClient(ctx context.Context) (*Client, error) {
	return newClient(ctx, loadedConfig.Host())
}

func newClient(ctx context.Context, target string) (*Client, error) {
	cc, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := &Client{
		CC:                   cc,
		Admin:                admin.NewAdminServiceClient(cc),
		Mgmt:                 mgmt.NewManagementServiceClient(cc),
		Auth:                 auth.NewAuthServiceClient(cc),
		UserV2:               user_v2.NewUserServiceClient(cc),
		SessionV2:            session.NewSessionServiceClient(cc),
		SettingsV2:           settings.NewSettingsServiceClient(cc),
		OIDCv2:               oidc_pb.NewOIDCServiceClient(cc),
		OrgV2:                org.NewOrganizationServiceClient(cc),
		ActionV2:             action.NewActionServiceClient(cc),
		FeatureV2:            feature.NewFeatureServiceClient(cc),
		WebKeyV2:             webkey_v2.NewWebKeyServiceClient(cc),
		IDPv2:                idp_pb.NewIdentityProviderServiceClient(cc),
		SAMLv2:               saml_pb.NewSAMLServiceClient(cc),
		SCIM:                 scim.NewScimClient(target),
		ProjectV2:            project_v2.NewProjectServiceClient(cc),
		InstanceV2:           instance_v2.NewInstanceServiceClient(cc),
		ApplicationV2:        application.NewApplicationServiceClient(cc),
		InternalPermissionV2: internal_permission_v2.NewInternalPermissionServiceClient(cc),
		AuthorizationV2:      authorization_v2.NewAuthorizationServiceClient(cc),
		GroupV2:              group_v2.NewGroupServiceClient(cc),
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

// Deprecated: use CreateUserTypeHuman instead
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
	i.TriggerUserByID(ctx, resp.GetUserId())
	return resp
}

// Deprecated: user CreateUserTypeHuman instead
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
	i.TriggerUserByID(ctx, resp.GetUserId())
	return resp
}

// Deprecated: user CreateUserTypeHuman instead
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
	i.TriggerUserByID(ctx, resp.GetUserId())
	return resp
}

func (i *Instance) SetUserMetadata(ctx context.Context, id, key, value string) *user_v2.SetUserMetadataResponse {
	resp, err := i.Client.UserV2.SetUserMetadata(ctx, &user_v2.SetUserMetadataRequest{
		UserId: id,
		Metadata: []*user_v2.Metadata{{
			Key:   key,
			Value: []byte(base64.StdEncoding.EncodeToString([]byte(value))),
		},
		},
	})
	logging.OnError(err).Panic("set user metadata")
	return resp
}

func (i *Instance) DeleteUserMetadata(ctx context.Context, id, key string) *user_v2.DeleteUserMetadataResponse {
	resp, err := i.Client.UserV2.DeleteUserMetadata(ctx, &user_v2.DeleteUserMetadataRequest{
		UserId: id,
		Keys:   []string{key},
	})
	logging.OnError(err).Panic("delete user metadata")
	return resp
}

func (i *Instance) CreateUserTypeHuman(ctx context.Context, email string) *user_v2.CreateUserResponse {
	resp, err := i.Client.UserV2.CreateUser(ctx, &user_v2.CreateUserRequest{
		OrganizationId: i.DefaultOrg.GetId(),
		UserType: &user_v2.CreateUserRequest_Human_{
			Human: &user_v2.CreateUserRequest_Human{
				Profile: &user_v2.SetHumanProfile{
					GivenName:  "Mickey",
					FamilyName: "Mouse",
				},
				Email: &user_v2.SetHumanEmail{
					Email: email,
					Verification: &user_v2.SetHumanEmail_ReturnCode{
						ReturnCode: &user_v2.ReturnEmailVerificationCode{},
					},
				},
			},
		},
	})
	logging.OnError(err).Panic("create human user")
	i.TriggerUserByID(ctx, resp.GetId())
	return resp
}

func (i *Instance) CreateUserTypeMachine(ctx context.Context, orgId string) *user_v2.CreateUserResponse {
	if orgId == "" {
		orgId = i.DefaultOrg.GetId()
	}
	resp, err := i.Client.UserV2.CreateUser(ctx, &user_v2.CreateUserRequest{
		OrganizationId: orgId,
		UserType: &user_v2.CreateUserRequest_Machine_{
			Machine: &user_v2.CreateUserRequest_Machine{
				Name: "machine",
			},
		},
	})
	logging.OnError(err).Panic("create service account")
	i.TriggerUserByID(ctx, resp.GetId())
	return resp
}

func (i *Instance) CreatePersonalAccessToken(ctx context.Context, userID string) *user_v2.AddPersonalAccessTokenResponse {
	resp, err := i.Client.UserV2.AddPersonalAccessToken(ctx, &user_v2.AddPersonalAccessTokenRequest{
		UserId:         userID,
		ExpirationDate: timestamppb.New(time.Now().Add(30 * time.Minute)),
	})
	logging.OnError(err).Panic("create pat")
	return resp
}

// TriggerUserByID makes sure the user projection gets triggered after creation.
func (i *Instance) TriggerUserByID(ctx context.Context, users ...string) {
	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, user := range users {
		go func(user string) {
			defer wg.Done()
			_, err := i.Client.UserV2.GetUserByID(ctx, &user_v2.GetUserByIDRequest{
				UserId: user,
			})
			logging.OnError(err).Warn("get user by ID for trigger failed")
		}(user)
	}
	wg.Wait()
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

	users := make([]string, len(resp.GetCreatedAdmins()))
	for i, admin := range resp.GetCreatedAdmins() {
		users[i] = admin.GetUserId()
	}
	i.TriggerUserByID(ctx, users...)

	return resp
}

func (i *Instance) SetOrganizationMetadata(ctx context.Context, id, key, value string) *org.SetOrganizationMetadataResponse {
	resp, err := i.Client.OrgV2.SetOrganizationMetadata(ctx, &org.SetOrganizationMetadataRequest{
		OrganizationId: id,
		Metadata: []*org.Metadata{{
			Key:   key,
			Value: []byte(base64.StdEncoding.EncodeToString([]byte(value))),
		},
		},
	})
	logging.OnError(err).Panic("set organization metadata")
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

func (i *Instance) CreateOrganizationWithCustomOrgID(ctx context.Context, name, orgID string) *org.AddOrganizationResponse {
	resp, err := i.Client.OrgV2.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name:  name,
		OrgId: gu.Ptr(orgID),
	})
	logging.OnError(err).Fatal("create org")
	return resp
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

func (i *Instance) SetOrganizationSettings(ctx context.Context, t *testing.T, orgID string, organizationScopedUsernames bool) *settings.SetOrganizationSettingsResponse {
	resp, err := i.Client.SettingsV2.SetOrganizationSettings(ctx,
		&settings.SetOrganizationSettingsRequest{
			OrganizationId:              orgID,
			OrganizationScopedUsernames: gu.Ptr(organizationScopedUsernames),
		},
	)
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteOrganizationSettings(ctx context.Context, t *testing.T, orgID string) *settings.DeleteOrganizationSettingsResponse {
	resp, err := i.Client.SettingsV2.DeleteOrganizationSettings(ctx,
		&settings.DeleteOrganizationSettingsRequest{
			OrganizationId: orgID,
		},
	)
	require.NoError(t, err)
	return resp
}

func (i *Instance) CreateHumanUserVerified(ctx context.Context, org, email, phone string) *user_v2.AddHumanUserResponse {
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
			Phone: phone,
			Verification: &user_v2.SetHumanPhone_IsVerified{
				IsVerified: true,
			},
		},
	})
	logging.OnError(err).Panic("create human user")
	i.TriggerUserByID(ctx, resp.GetUserId())
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
	i.TriggerUserByID(ctx, resp.GetUserId())
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

func (i *Instance) RegisterUserPasskey(ctx context.Context, userID string) string {
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
	return pkr.GetPasskeyId()
}

func (i *Instance) RegisterUserU2F(ctx context.Context, userID string) string {
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
	return pkr.GetU2FId()
}

func (i *Instance) RegisterUserOTPSMS(ctx context.Context, userID string) {
	_, err := i.Client.UserV2.AddOTPSMS(ctx, &user_v2.AddOTPSMSRequest{
		UserId: userID,
	})
	logging.OnError(err).Panic("create user sms")
}

func (i *Instance) RegisterUserOTPEmail(ctx context.Context, userID string) {
	_, err := i.Client.UserV2.AddOTPEmail(ctx, &user_v2.AddOTPEmailRequest{
		UserId: userID,
	})
	logging.OnError(err).Panic("create user email")
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

func (i *Instance) CreateProject(ctx context.Context, t *testing.T, orgID, name string, projectRoleCheck, hasProjectCheck bool) *project_v2.CreateProjectResponse {
	if orgID == "" {
		orgID = i.DefaultOrg.GetId()
	}

	resp, err := i.Client.ProjectV2.CreateProject(ctx, &project_v2.CreateProjectRequest{
		OrganizationId:        orgID,
		Name:                  name,
		AuthorizationRequired: projectRoleCheck,
		ProjectAccessRequired: hasProjectCheck,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProject(ctx context.Context, t *testing.T, projectID string) *project_v2.DeleteProjectResponse {
	resp, err := i.Client.ProjectV2.DeleteProject(ctx, &project_v2.DeleteProjectRequest{
		ProjectId: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeactivateProject(ctx context.Context, t *testing.T, projectID string) *project_v2.DeactivateProjectResponse {
	resp, err := i.Client.ProjectV2.DeactivateProject(ctx, &project_v2.DeactivateProjectRequest{
		ProjectId: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) ActivateProject(ctx context.Context, t *testing.T, projectID string) *project_v2.ActivateProjectResponse {
	resp, err := i.Client.ProjectV2.ActivateProject(ctx, &project_v2.ActivateProjectRequest{
		ProjectId: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) AddProjectRole(ctx context.Context, t *testing.T, projectID, roleKey, displayName, group string) *project_v2.AddProjectRoleResponse {
	var groupP *string
	if group != "" {
		groupP = &group
	}

	resp, err := i.Client.ProjectV2.AddProjectRole(ctx, &project_v2.AddProjectRoleRequest{
		ProjectId:   projectID,
		RoleKey:     roleKey,
		DisplayName: displayName,
		Group:       groupP,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) RemoveProjectRole(ctx context.Context, t *testing.T, projectID, roleKey string) *project_v2.RemoveProjectRoleResponse {
	resp, err := i.Client.ProjectV2.RemoveProjectRole(ctx, &project_v2.RemoveProjectRoleRequest{
		ProjectId: projectID,
		RoleKey:   roleKey,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) AddProviderToDefaultLoginPolicy(ctx context.Context, id string) {
	_, err := i.Client.Admin.AddIDPToLoginPolicy(ctx, &admin.AddIDPToLoginPolicyRequest{
		IdpId: id,
	})
	logging.OnError(err).Panic("add provider to default login policy")
}

func (i *Instance) AddAzureADProvider(ctx context.Context, name string) *admin.AddAzureADProviderResponse {
	resp, err := i.Client.Admin.AddAzureADProvider(ctx, &admin.AddAzureADProviderRequest{
		Name:          name,
		ClientId:      "clientID",
		ClientSecret:  "clientSecret",
		Tenant:        nil,
		EmailVerified: false,
		Scopes:        []string{"openid", "profile", "email"},
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
			AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
		},
	})
	logging.OnError(err).Panic("create Azure AD idp")

	mustAwait(func() error {
		_, err := i.Client.Admin.GetProviderByID(ctx, &admin.GetProviderByIDRequest{
			Id: resp.GetId(),
		})
		return err
	})

	return resp
}

func (i *Instance) AddGenericOAuthProvider(ctx context.Context, name string) *admin.AddGenericOAuthProviderResponse {
	return i.AddGenericOAuthProviderWithOptions(ctx, name, true, true, true, idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME)
}

func (i *Instance) AddGenericOAuthProviderWithOptions(ctx context.Context, name string, isLinkingAllowed, isCreationAllowed, isAutoCreation bool, autoLinking idp.AutoLinkingOption) *admin.AddGenericOAuthProviderResponse {
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
			IsLinkingAllowed:  isLinkingAllowed,
			IsCreationAllowed: isCreationAllowed,
			IsAutoCreation:    isAutoCreation,
			IsAutoUpdate:      true,
			AutoLinking:       autoLinking,
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
	return resp
}

func (i *Instance) AddGenericOIDCProvider(ctx context.Context, name string) *admin.AddGenericOIDCProviderResponse {
	resp, err := i.Client.Admin.AddGenericOIDCProvider(ctx, &admin.AddGenericOIDCProviderRequest{
		Name:         name,
		Issuer:       "https://example.com",
		ClientId:     "clientID",
		ClientSecret: "clientSecret",
		Scopes:       []string{"openid", "profile", "email"},
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
			AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
		},
		IsIdTokenMapping: false,
	})
	logging.OnError(err).Panic("create generic oidc idp")
	return resp
}

func (i *Instance) AddSAMLProvider(ctx context.Context) string {
	cert := i.SAMLCertificateString()
	resp, err := i.Client.Admin.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name: "saml-idp",
		Metadata: &admin.AddSAMLProviderRequest_MetadataXml{
			MetadataXml: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">" + cert + "</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">" + cert + "</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
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
	cert := i.SAMLCertificateString()
	resp, err := i.Client.Admin.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name:    "saml-idp-redirect",
		Binding: idp.SAMLBinding_SAML_BINDING_REDIRECT,
		Metadata: &admin.AddSAMLProviderRequest_MetadataXml{
			MetadataXml: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">" + cert + "</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">" + cert + "</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
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
	cert := i.SAMLCertificateString()
	resp, err := i.Client.Admin.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name:    "saml-idp-post",
		Binding: idp.SAMLBinding_SAML_BINDING_POST,
		Metadata: &admin.AddSAMLProviderRequest_MetadataXml{
			MetadataXml: []byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-09-16T09:00:32.986Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">" + cert + "</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">" + cert + "</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
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

var samlPrivateKeyOnce sync.Map  // map[*Instance]*sync.Once
var samlCertificateOnce sync.Map // map[*Instance]*sync.Once

func (i *Instance) SAMLPrivateKey() *rsa.PrivateKey {
	if i.SAML.PrivateKey != nil {
		return i.SAML.PrivateKey
	}

	onceIface, _ := samlPrivateKeyOnce.LoadOrStore(i, &sync.Once{})
	once := onceIface.(*sync.Once)

	once.Do(func() {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		logging.OnError(err).Panic("generate saml private key")
		i.SAML.PrivateKey = privateKey
	})

	return i.SAML.PrivateKey
}

func (i *Instance) SAMLCertificateString() string {
	return base64.StdEncoding.EncodeToString(i.SAMLCertificate().Raw)
}
func (i *Instance) SAMLCertificate() *x509.Certificate {
	if i.SAML.Certificate != nil {
		return i.SAML.Certificate
	}

	onceIface, _ := samlCertificateOnce.LoadOrStore(i, &sync.Once{})
	once := onceIface.(*sync.Once)

	once.Do(func() {
		template := &x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{"Example Co"},
				CommonName:   "www.example.com",
			},
			Issuer: pkix.Name{
				Organization: []string{"Example Co"},
				CommonName:   "www.example.com",
			},
			NotBefore:             time.Now().Add(-time.Hour),
			NotAfter:              time.Now().Add(365 * 24 * time.Hour),
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			BasicConstraintsValid: true,
			IsCA:                  true,
		}
		privateKey := i.SAMLPrivateKey()
		cert, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
		logging.OnError(err).Panic("create saml certificate")
		parsedCert, err := x509.ParseCertificate(cert)
		logging.OnError(err).Panic("parse saml certificate")
		i.SAML.Certificate = parsedCert
	})
	return i.SAML.Certificate
}

func (i *Instance) AddLDAPProvider(ctx context.Context) string {
	resp, err := i.Client.Admin.AddLDAPProvider(ctx, &admin.AddLDAPProviderRequest{
		Name:              "ldap-idp-post",
		Servers:           []string{"https://localhost:8000"},
		StartTls:          false,
		BaseDn:            "baseDn",
		BindDn:            "admin",
		BindPassword:      "admin",
		UserBase:          "dn",
		UserObjectClasses: []string{"user"},
		UserFilters:       []string{"(objectclass=*)"},
		Timeout:           durationpb.New(10 * time.Second),
		Attributes: &idp.LDAPAttributes{
			IdAttribute: "id",
		},
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	logging.OnError(err).Panic("create ldap idp")
	return resp.GetId()
}

func (i *Instance) AddJWTProvider(ctx context.Context) string {
	resp, err := i.Client.Admin.AddJWTProvider(ctx, &admin.AddJWTProviderRequest{
		Name:         "jwt-idp",
		Issuer:       "https://example.com",
		JwtEndpoint:  "https://example.com/jwt",
		KeysEndpoint: "https://example.com/keys",
		HeaderName:   "Authorization",
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	logging.OnError(err).Panic("create jwt idp")
	return resp.GetId()
}

func (i *Instance) CreateIntent(ctx context.Context, idpID string) *user_v2.StartIdentityProviderIntentResponse {
	resp, err := i.Client.UserV2.StartIdentityProviderIntent(ctx, &user_v2.StartIdentityProviderIntentRequest{
		IdpId: idpID,
		Content: &user_v2.StartIdentityProviderIntentRequest_Urls{
			Urls: &user_v2.RedirectURLs{
				SuccessUrl: "https://example.com/success",
				FailureUrl: "https://example.com/failure",
			},
		},
	})
	logging.OnError(err).Fatal("create generic OAuth idp")
	return resp
}

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

func (i *Instance) CreateIntentSession(t *testing.T, ctx context.Context, userID, intentID, intentToken string) (id, token string, start, change time.Time) {
	createResp, err := i.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: userID},
			},
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: intentToken,
			},
		},
	})
	require.NoError(t, err)
	return createResp.GetSessionId(), createResp.GetSessionToken(),
		createResp.GetDetails().GetChangeDate().AsTime(), createResp.GetDetails().GetChangeDate().AsTime()
}

func (i *Instance) CreateProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string, roles ...string) *project_v2.CreateProjectGrantResponse {
	resp, err := i.Client.ProjectV2.CreateProjectGrant(ctx, &project_v2.CreateProjectGrantRequest{
		GrantedOrganizationId: grantedOrgID,
		ProjectId:             projectID,
		RoleKeys:              roles,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string) *project_v2.DeleteProjectGrantResponse {
	resp, err := i.Client.ProjectV2.DeleteProjectGrant(ctx, &project_v2.DeleteProjectGrantRequest{
		GrantedOrganizationId: grantedOrgID,
		ProjectId:             projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeactivateProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string) *project_v2.DeactivateProjectGrantResponse {
	resp, err := i.Client.ProjectV2.DeactivateProjectGrant(ctx, &project_v2.DeactivateProjectGrantRequest{
		ProjectId:             projectID,
		GrantedOrganizationId: grantedOrgID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) ActivateProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string) *project_v2.ActivateProjectGrantResponse {
	resp, err := i.Client.ProjectV2.ActivateProjectGrant(ctx, &project_v2.ActivateProjectGrantRequest{
		ProjectId:             projectID,
		GrantedOrganizationId: grantedOrgID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) CreateAuthorizationProject(t *testing.T, ctx context.Context, projectID, userID, organizationID string, roles ...string) *authorization_v2.CreateAuthorizationResponse {
	resp, err := i.Client.AuthorizationV2.CreateAuthorization(ctx, &authorization_v2.CreateAuthorizationRequest{
		UserId:         userID,
		ProjectId:      projectID,
		OrganizationId: organizationID,
		RoleKeys:       roles,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) CreateAuthorizationProjectGrant(t *testing.T, ctx context.Context, projectID, organizationID, userID string, roles ...string) *authorization_v2.CreateAuthorizationResponse {
	resp, err := i.Client.AuthorizationV2.CreateAuthorization(ctx, &authorization_v2.CreateAuthorizationRequest{
		UserId:         userID,
		ProjectId:      projectID,
		OrganizationId: organizationID,
		RoleKeys:       roles,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) CreateProjectUserGrant(t *testing.T, ctx context.Context, orgID, projectID, userID string) *mgmt.AddUserGrantResponse {
	//nolint:staticcheck
	resp, err := i.Client.Mgmt.AddUserGrant(SetOrgID(ctx, orgID), &mgmt.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) CreateProjectGrantUserGrant(ctx context.Context, orgID, projectID, projectGrantID, userID string) *mgmt.AddUserGrantResponse {
	//nolint:staticcheck
	resp, err := i.Client.Mgmt.AddUserGrant(SetOrgID(ctx, orgID), &mgmt.AddUserGrantRequest{
		UserId:         userID,
		ProjectId:      projectID,
		ProjectGrantId: projectGrantID,
	})
	logging.OnError(err).Panic("create project grant user grant")
	return resp
}

func (i *Instance) CreateInstanceMembership(t *testing.T, ctx context.Context, userID string) *internal_permission_v2.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_Instance{Instance: true},
		},
		UserId: userID,
		Roles:  []string{domain.RoleIAMOwner},
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteInstanceMembership(t *testing.T, ctx context.Context, userID string) {
	_, err := i.Client.Admin.RemoveIAMMember(ctx, &admin.RemoveIAMMemberRequest{
		UserId: userID,
	})
	require.NoError(t, err)
}

// CreateOrgMembership creates an org membership with the given roles. If no roles are provided, the user will be assigned the "org.owner" role.
func (i *Instance) CreateOrgMembership(t *testing.T, ctx context.Context, orgID, userID string, roles ...string) *internal_permission_v2.CreateAdministratorResponse {
	if len(roles) == 0 {
		roles = []string{domain.RoleOrgOwner}
	}
	resp, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_OrganizationId{OrganizationId: orgID},
		},
		UserId: userID,
		Roles:  roles,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteOrgMembership(t *testing.T, ctx context.Context, userID string) {
	_, err := i.Client.Mgmt.RemoveOrgMember(ctx, &mgmt.RemoveOrgMemberRequest{
		UserId: userID,
	})
	require.NoError(t, err)
}

func (i *Instance) CreateProjectMembership(t *testing.T, ctx context.Context, projectID, userID string) *internal_permission_v2.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_ProjectId{ProjectId: projectID},
		},
		UserId: userID,
		Roles:  []string{domain.RoleProjectOwner},
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProjectMembership(t *testing.T, ctx context.Context, projectID, userID string) {
	_, err := i.Client.InternalPermissionV2.DeleteAdministrator(ctx, &internal_permission_v2.DeleteAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{Resource: &internal_permission_v2.ResourceType_ProjectId{ProjectId: projectID}},
		UserId:   userID,
	})
	require.NoError(t, err)
}

func (i *Instance) CreateProjectGrantMembership(t *testing.T, ctx context.Context, projectID, orgID, userID string) *internal_permission_v2.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionV2.CreateAdministrator(ctx, &internal_permission_v2.CreateAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_ProjectGrant_{ProjectGrant: &internal_permission_v2.ResourceType_ProjectGrant{
				ProjectId:      projectID,
				OrganizationId: orgID,
			}},
		},
		UserId: userID,
		Roles:  []string{domain.RoleProjectGrantOwner},
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProjectGrantMembership(t *testing.T, ctx context.Context, projectID, orgID, userID string) {
	_, err := i.Client.InternalPermissionV2.DeleteAdministrator(ctx, &internal_permission_v2.DeleteAdministratorRequest{
		Resource: &internal_permission_v2.ResourceType{
			Resource: &internal_permission_v2.ResourceType_ProjectGrant_{ProjectGrant: &internal_permission_v2.ResourceType_ProjectGrant{
				ProjectId:      projectID,
				OrganizationId: orgID,
			}},
		},
		UserId: userID,
	})
	require.NoError(t, err)
}

func (i *Instance) CreateTargetWithoutPayloadType(ctx context.Context, t *testing.T, name, endpoint string, ty target_domain.TargetType, interrupt bool) *action.CreateTargetResponse {
	return i.CreateTarget(ctx, t, name, endpoint, ty, interrupt, action.PayloadType_PAYLOAD_TYPE_UNSPECIFIED)
}

func (i *Instance) CreateTarget(ctx context.Context, t *testing.T, name, endpoint string, ty target_domain.TargetType, interrupt bool, payloadType action.PayloadType) *action.CreateTargetResponse {
	if name == "" {
		name = TargetName()
	}
	req := &action.CreateTargetRequest{
		Name:        name,
		Endpoint:    endpoint,
		Timeout:     durationpb.New(5 * time.Second),
		PayloadType: payloadType,
	}
	switch ty {
	case target_domain.TargetTypeWebhook:
		req.TargetType = &action.CreateTargetRequest_RestWebhook{
			RestWebhook: &action.RESTWebhook{
				InterruptOnError: interrupt,
			},
		}
	case target_domain.TargetTypeCall:
		req.TargetType = &action.CreateTargetRequest_RestCall{
			RestCall: &action.RESTCall{
				InterruptOnError: interrupt,
			},
		}
	case target_domain.TargetTypeAsync:
		req.TargetType = &action.CreateTargetRequest_RestAsync{
			RestAsync: &action.RESTAsync{},
		}
	}
	target, err := i.Client.ActionV2.CreateTarget(ctx, req)
	require.NoError(t, err)
	return target
}

func (i *Instance) DeleteTarget(ctx context.Context, t *testing.T, id string) {
	_, err := i.Client.ActionV2.DeleteTarget(ctx, &action.DeleteTargetRequest{
		Id: id,
	})
	require.NoError(t, err)
}

func (i *Instance) DeleteExecution(ctx context.Context, t *testing.T, cond *action.Condition) {
	_, err := i.Client.ActionV2.SetExecution(ctx, &action.SetExecutionRequest{
		Condition: cond,
	})
	require.NoError(t, err)
}

func (i *Instance) SetExecution(ctx context.Context, t *testing.T, cond *action.Condition, targets []string) *action.SetExecutionResponse {
	target, err := i.Client.ActionV2.SetExecution(ctx, &action.SetExecutionRequest{
		Condition: cond,
		Targets:   targets,
	})
	require.NoError(t, err)
	return target
}

func (i *Instance) CreateInviteCode(ctx context.Context, userID string) *user_v2.CreateInviteCodeResponse {
	user, err := i.Client.UserV2.CreateInviteCode(ctx, &user_v2.CreateInviteCodeRequest{
		UserId:       userID,
		Verification: &user_v2.CreateInviteCodeRequest_ReturnCode{ReturnCode: &user_v2.ReturnInviteCode{}},
	})
	logging.OnError(err).Fatal("create invite code")
	return user
}

func (i *Instance) CreateGroup(ctx context.Context, t *testing.T, orgID, name string) *group_v2.CreateGroupResponse {
	resp, err := i.Client.GroupV2.CreateGroup(ctx, &group_v2.CreateGroupRequest{
		OrganizationId: orgID,
		Name:           name,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteGroup(ctx context.Context, t *testing.T, id string) *group_v2.DeleteGroupResponse {
	resp, err := i.Client.GroupV2.DeleteGroup(ctx, &group_v2.DeleteGroupRequest{
		Id: id,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) AddUsersToGroup(ctx context.Context, t *testing.T, groupID string, userIDs []string) *group_v2.AddUsersToGroupResponse {
	resp, err := i.Client.GroupV2.AddUsersToGroup(ctx, &group_v2.AddUsersToGroupRequest{
		Id:      groupID,
		UserIds: userIDs,
	})
	require.NoError(t, err)
	return resp
}
