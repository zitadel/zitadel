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
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

type Client struct {
	CC        *grpc.ClientConn
	Admin     admin.AdminServiceClient
	UserV2    user.UserServiceClient
	SessionV2 session.SessionServiceClient
}

func newClient(cc *grpc.ClientConn) Client {
	return Client{
		CC:        cc,
		Admin:     admin.NewAdminServiceClient(cc),
		UserV2:    user.NewUserServiceClient(cc),
		SessionV2: session.NewSessionServiceClient(cc),
	}
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
				IdpId:         idpID,
				IdpExternalId: externalID,
				DisplayName:   username,
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

func (s *Tester) AddGenericOAuthProvider(t *testing.T) string {
	ctx := authz.WithInstance(context.Background(), s.Instance)
	id, _, err := s.Commands.AddOrgGenericOAuthProvider(ctx, s.Organisation.ID, command.GenericOAuthProvider{
		"idp",
		"clientID",
		"clientSecret",
		"https://example.com/oauth/v2/authorize",
		"https://example.com/oauth/v2/token",
		"https://api.example.com/user",
		[]string{"openid", "profile", "email"},
		"id",
		idp.Options{
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

func (s *Tester) CreateSuccessfulIntent(t *testing.T, idpID, userID string) (string, string, time.Time, uint64) {
	ctx := authz.WithInstance(context.Background(), s.Instance)
	intentID := s.CreateIntent(t, idpID)
	writeModel, err := s.Commands.GetIntentWriteModel(ctx, intentID, s.Organisation.ID)
	require.NoError(t, err)
	idpUser := &oauth.UserMapper{
		RawInfo: map[string]interface{}{
			"id": "id",
		},
	}
	idpSession := &oauth.Session{
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
