package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/zitadel/logging"
	"google.golang.org/grpc"

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
