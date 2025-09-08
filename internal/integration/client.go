package integration

import (
	"context"
	"encoding/base64"
	"fmt"
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
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
	action_v2beta "github.com/zitadel/zitadel/pkg/grpc/action/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	feature_v2beta "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	internal_permission_v2beta "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	object_v3alpha "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	oidc_pb_v2beta "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	org_v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	project_v2beta "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	user_v3alpha "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
	userschema_v3alpha "github.com/zitadel/zitadel/pkg/grpc/resources/userschema/v3alpha"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	session_v2beta "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	settings_v2beta "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
	webkey_v2 "github.com/zitadel/zitadel/pkg/grpc/webkey/v2"
	webkey_v2beta "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
)

type Client struct {
	CC                       *grpc.ClientConn
	Admin                    admin.AdminServiceClient
	Mgmt                     mgmt.ManagementServiceClient
	Auth                     auth.AuthServiceClient
	UserV2beta               user_v2beta.UserServiceClient
	UserV2                   user_v2.UserServiceClient
	SessionV2beta            session_v2beta.SessionServiceClient
	SessionV2                session.SessionServiceClient
	SettingsV2beta           settings_v2beta.SettingsServiceClient
	SettingsV2               settings.SettingsServiceClient
	OIDCv2beta               oidc_pb_v2beta.OIDCServiceClient
	OIDCv2                   oidc_pb.OIDCServiceClient
	OrgV2beta                org_v2beta.OrganizationServiceClient
	OrgV2                    org.OrganizationServiceClient
	ActionV2beta             action_v2beta.ActionServiceClient
	ActionV2                 action.ActionServiceClient
	FeatureV2beta            feature_v2beta.FeatureServiceClient
	FeatureV2                feature.FeatureServiceClient
	UserSchemaV3             userschema_v3alpha.ZITADELUserSchemasClient
	WebKeyV2                 webkey_v2.WebKeyServiceClient
	WebKeyV2Beta             webkey_v2beta.WebKeyServiceClient
	IDPv2                    idp_pb.IdentityProviderServiceClient
	UserV3Alpha              user_v3alpha.ZITADELUsersClient
	SAMLv2                   saml_pb.SAMLServiceClient
	SCIM                     *scim.Client
	Projectv2Beta            project_v2beta.ProjectServiceClient
	InstanceV2Beta           instance.InstanceServiceClient
	AppV2Beta                app.AppServiceClient
	InternalPermissionv2Beta internal_permission_v2beta.InternalPermissionServiceClient
	AuthorizationV2Beta      authorization.AuthorizationServiceClient
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
		CC:                       cc,
		Admin:                    admin.NewAdminServiceClient(cc),
		Mgmt:                     mgmt.NewManagementServiceClient(cc),
		Auth:                     auth.NewAuthServiceClient(cc),
		UserV2beta:               user_v2beta.NewUserServiceClient(cc),
		UserV2:                   user_v2.NewUserServiceClient(cc),
		SessionV2beta:            session_v2beta.NewSessionServiceClient(cc),
		SessionV2:                session.NewSessionServiceClient(cc),
		SettingsV2beta:           settings_v2beta.NewSettingsServiceClient(cc),
		SettingsV2:               settings.NewSettingsServiceClient(cc),
		OIDCv2beta:               oidc_pb_v2beta.NewOIDCServiceClient(cc),
		OIDCv2:                   oidc_pb.NewOIDCServiceClient(cc),
		OrgV2beta:                org_v2beta.NewOrganizationServiceClient(cc),
		OrgV2:                    org.NewOrganizationServiceClient(cc),
		ActionV2beta:             action_v2beta.NewActionServiceClient(cc),
		ActionV2:                 action.NewActionServiceClient(cc),
		FeatureV2beta:            feature_v2beta.NewFeatureServiceClient(cc),
		FeatureV2:                feature.NewFeatureServiceClient(cc),
		UserSchemaV3:             userschema_v3alpha.NewZITADELUserSchemasClient(cc),
		WebKeyV2:                 webkey_v2.NewWebKeyServiceClient(cc),
		WebKeyV2Beta:             webkey_v2beta.NewWebKeyServiceClient(cc),
		IDPv2:                    idp_pb.NewIdentityProviderServiceClient(cc),
		UserV3Alpha:              user_v3alpha.NewZITADELUsersClient(cc),
		SAMLv2:                   saml_pb.NewSAMLServiceClient(cc),
		SCIM:                     scim.NewScimClient(target),
		Projectv2Beta:            project_v2beta.NewProjectServiceClient(cc),
		InstanceV2Beta:           instance.NewInstanceServiceClient(cc),
		AppV2Beta:                app.NewAppServiceClient(cc),
		InternalPermissionv2Beta: internal_permission_v2beta.NewInternalPermissionServiceClient(cc),
		AuthorizationV2Beta:      authorization.NewAuthorizationServiceClient(cc),
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
	logging.OnError(err).Panic("create machine user")
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

func (i *Instance) SetOrganizationSettings(ctx context.Context, t *testing.T, orgID string, organizationScopedUsernames bool) *settings_v2beta.SetOrganizationSettingsResponse {
	resp, err := i.Client.SettingsV2beta.SetOrganizationSettings(ctx,
		&settings_v2beta.SetOrganizationSettingsRequest{
			OrganizationId:              orgID,
			OrganizationScopedUsernames: gu.Ptr(organizationScopedUsernames),
		},
	)
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteOrganizationSettings(ctx context.Context, t *testing.T, orgID string) *settings_v2beta.DeleteOrganizationSettingsResponse {
	resp, err := i.Client.SettingsV2beta.DeleteOrganizationSettings(ctx,
		&settings_v2beta.DeleteOrganizationSettingsRequest{
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

func (i *Instance) CreateProject(ctx context.Context, t *testing.T, orgID, name string, projectRoleCheck, hasProjectCheck bool) *project_v2beta.CreateProjectResponse {
	if orgID == "" {
		orgID = i.DefaultOrg.GetId()
	}

	resp, err := i.Client.Projectv2Beta.CreateProject(ctx, &project_v2beta.CreateProjectRequest{
		OrganizationId:        orgID,
		Name:                  name,
		AuthorizationRequired: projectRoleCheck,
		ProjectAccessRequired: hasProjectCheck,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProject(ctx context.Context, t *testing.T, projectID string) *project_v2beta.DeleteProjectResponse {
	resp, err := i.Client.Projectv2Beta.DeleteProject(ctx, &project_v2beta.DeleteProjectRequest{
		Id: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeactivateProject(ctx context.Context, t *testing.T, projectID string) *project_v2beta.DeactivateProjectResponse {
	resp, err := i.Client.Projectv2Beta.DeactivateProject(ctx, &project_v2beta.DeactivateProjectRequest{
		Id: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) ActivateProject(ctx context.Context, t *testing.T, projectID string) *project_v2beta.ActivateProjectResponse {
	resp, err := i.Client.Projectv2Beta.ActivateProject(ctx, &project_v2beta.ActivateProjectRequest{
		Id: projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) AddProjectRole(ctx context.Context, t *testing.T, projectID, roleKey, displayName, group string) *project_v2beta.AddProjectRoleResponse {
	var groupP *string
	if group != "" {
		groupP = &group
	}

	resp, err := i.Client.Projectv2Beta.AddProjectRole(ctx, &project_v2beta.AddProjectRoleRequest{
		ProjectId:   projectID,
		RoleKey:     roleKey,
		DisplayName: displayName,
		Group:       groupP,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) RemoveProjectRole(ctx context.Context, t *testing.T, projectID, roleKey string) *project_v2beta.RemoveProjectRoleResponse {
	resp, err := i.Client.Projectv2Beta.RemoveProjectRole(ctx, &project_v2beta.RemoveProjectRoleRequest{
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

func (i *Instance) CreateProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string, roles ...string) *project_v2beta.CreateProjectGrantResponse {
	resp, err := i.Client.Projectv2Beta.CreateProjectGrant(ctx, &project_v2beta.CreateProjectGrantRequest{
		GrantedOrganizationId: grantedOrgID,
		ProjectId:             projectID,
		RoleKeys:              roles,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string) *project_v2beta.DeleteProjectGrantResponse {
	resp, err := i.Client.Projectv2Beta.DeleteProjectGrant(ctx, &project_v2beta.DeleteProjectGrantRequest{
		GrantedOrganizationId: grantedOrgID,
		ProjectId:             projectID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeactivateProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string) *project_v2beta.DeactivateProjectGrantResponse {
	resp, err := i.Client.Projectv2Beta.DeactivateProjectGrant(ctx, &project_v2beta.DeactivateProjectGrantRequest{
		ProjectId:             projectID,
		GrantedOrganizationId: grantedOrgID,
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) ActivateProjectGrant(ctx context.Context, t *testing.T, projectID, grantedOrgID string) *project_v2beta.ActivateProjectGrantResponse {
	resp, err := i.Client.Projectv2Beta.ActivateProjectGrant(ctx, &project_v2beta.ActivateProjectGrantRequest{
		ProjectId:             projectID,
		GrantedOrganizationId: grantedOrgID,
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

func (i *Instance) CreateInstanceMembership(t *testing.T, ctx context.Context, userID string) *internal_permission_v2beta.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionv2Beta.CreateAdministrator(ctx, &internal_permission_v2beta.CreateAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{
			Resource: &internal_permission_v2beta.ResourceType_Instance{Instance: true},
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

func (i *Instance) CreateOrgMembership(t *testing.T, ctx context.Context, orgID, userID string) *internal_permission_v2beta.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionv2Beta.CreateAdministrator(ctx, &internal_permission_v2beta.CreateAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{
			Resource: &internal_permission_v2beta.ResourceType_OrganizationId{OrganizationId: orgID},
		},
		UserId: userID,
		Roles:  []string{domain.RoleOrgOwner},
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

func (i *Instance) CreateProjectMembership(t *testing.T, ctx context.Context, projectID, userID string) *internal_permission_v2beta.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionv2Beta.CreateAdministrator(ctx, &internal_permission_v2beta.CreateAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{
			Resource: &internal_permission_v2beta.ResourceType_ProjectId{ProjectId: projectID},
		},
		UserId: userID,
		Roles:  []string{domain.RoleProjectOwner},
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProjectMembership(t *testing.T, ctx context.Context, projectID, userID string) {
	_, err := i.Client.InternalPermissionv2Beta.DeleteAdministrator(ctx, &internal_permission_v2beta.DeleteAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{Resource: &internal_permission_v2beta.ResourceType_ProjectId{ProjectId: projectID}},
		UserId:   userID,
	})
	require.NoError(t, err)
}

func (i *Instance) CreateProjectGrantMembership(t *testing.T, ctx context.Context, projectID, grantID, userID string) *internal_permission_v2beta.CreateAdministratorResponse {
	resp, err := i.Client.InternalPermissionv2Beta.CreateAdministrator(ctx, &internal_permission_v2beta.CreateAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{
			Resource: &internal_permission_v2beta.ResourceType_ProjectGrant_{ProjectGrant: &internal_permission_v2beta.ResourceType_ProjectGrant{
				ProjectId:      projectID,
				ProjectGrantId: grantID,
			}},
		},
		UserId: userID,
		Roles:  []string{domain.RoleProjectGrantOwner},
	})
	require.NoError(t, err)
	return resp
}

func (i *Instance) DeleteProjectGrantMembership(t *testing.T, ctx context.Context, projectID, grantID, userID string) {
	_, err := i.Client.InternalPermissionv2Beta.DeleteAdministrator(ctx, &internal_permission_v2beta.DeleteAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{
			Resource: &internal_permission_v2beta.ResourceType_ProjectGrant_{ProjectGrant: &internal_permission_v2beta.ResourceType_ProjectGrant{
				ProjectId:      projectID,
				ProjectGrantId: grantID,
			}},
		},
		UserId: userID,
	})
	require.NoError(t, err)
}

func (i *Instance) CreateTarget(ctx context.Context, t *testing.T, name, endpoint string, ty target_domain.TargetType, interrupt bool) *action.CreateTargetResponse {
	if name == "" {
		name = TargetName()
	}
	req := &action.CreateTargetRequest{
		Name:     name,
		Endpoint: endpoint,
		Timeout:  durationpb.New(5 * time.Second),
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
