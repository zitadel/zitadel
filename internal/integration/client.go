package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

type Client struct {
	CC        *grpc.ClientConn
	Admin     admin.AdminServiceClient
	Mgmt      mgmt.ManagementServiceClient
	Auth      auth.AuthServiceClient
	UserV2    user.UserServiceClient
	SessionV2 session.SessionServiceClient
	OIDCv2    oidc_pb.OIDCServiceClient
	System    system.SystemServiceClient
}

func newClient(cc *grpc.ClientConn) Client {
	return Client{
		CC:        cc,
		Admin:     admin.NewAdminServiceClient(cc),
		Mgmt:      mgmt.NewManagementServiceClient(cc),
		Auth:      auth.NewAuthServiceClient(cc),
		UserV2:    user.NewUserServiceClient(cc),
		SessionV2: session.NewSessionServiceClient(cc),
		OIDCv2:    oidc_pb.NewOIDCServiceClient(cc),
		System:    system.NewSystemServiceClient(cc),
	}
}

func (t *Tester) UseIsolatedInstance(iamOwnerCtx, systemCtx context.Context) (primaryDomain, instanceId string, authenticatedIamOwnerCtx context.Context) {
	primaryDomain = randString(5) + ".integration"
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
	if err != nil {
		panic(err)
	}
	t.createClientConn(iamOwnerCtx, grpc.WithAuthority(primaryDomain))
	instanceId = instance.GetInstanceId()
	t.Users.Set(instanceId, IAMOwner, &User{
		Token: instance.GetPat(),
	})
	return primaryDomain, instanceId, t.WithInstanceAuthorization(iamOwnerCtx, IAMOwner, instanceId)
}

func (s *Tester) CreateHumanUser(ctx context.Context) *user.AddHumanUserResponse {
	resp, err := s.Client.UserV2.AddHumanUser(ctx, &user.AddHumanUserRequest{
		Organisation: &object.Organisation{
			Org: &object.Organisation_OrgId{
				OrgId: s.Organisation.ID,
			},
		},
		Profile: &user.SetHumanProfile{
			FirstName: "Mickey",
			LastName:  "Mouse",
		},
		Email: &user.SetHumanEmail{
			Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
			Verification: &user.SetHumanEmail_ReturnCode{
				ReturnCode: &user.ReturnEmailVerificationCode{},
			},
		},
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

func (s *Tester) SetUserPassword(ctx context.Context, userID, password string) {
	_, err := s.Client.UserV2.SetPassword(ctx, &user.SetPasswordRequest{
		UserId:      userID,
		NewPassword: &user.Password{Password: password},
	})
	logging.OnError(err).Fatal("set user password")
}

func (s *Tester) AddGenericOAuthProvider(t *testing.T) string {
	ctx := authz.WithInstance(context.Background(), s.Instance)
	id, _, err := s.Commands.AddOrgGenericOAuthProvider(ctx, s.Organisation.ID, command.GenericOAuthProvider{
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

func (s *Tester) CreateIntent(t *testing.T, idpID string) string {
	ctx := authz.WithInstance(context.Background(), s.Instance)
	id, _, err := s.Commands.CreateIntent(ctx, idpID, "https://example.com/success", "https://example.com/failure", s.Organisation.ID)
	require.NoError(t, err)
	return id
}

func (s *Tester) CreateSuccessfulIntent(t *testing.T, idpID, userID, idpUserID string) (string, string, time.Time, uint64) {
	ctx := authz.WithInstance(context.Background(), s.Instance)
	intentID := s.CreateIntent(t, idpID)
	writeModel, err := s.Commands.GetIntentWriteModel(ctx, intentID, s.Organisation.ID)
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

func (s *Tester) CreatePasskeySession(t *testing.T, ctx context.Context, userID string) (id, token string, start, change time.Time) {
	createResp, err := s.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: userID},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
		Domain: s.Config.ExternalDomain,
	})
	require.NoError(t, err)

	assertion, err := s.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := s.Client.SessionV2.SetSession(ctx, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
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
		Domain: s.Config.ExternalDomain,
	})
	require.NoError(t, err)
	return createResp.GetSessionId(), createResp.GetSessionToken(),
		createResp.GetDetails().GetChangeDate().AsTime(), createResp.GetDetails().GetChangeDate().AsTime()
}
