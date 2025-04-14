//go:build integration

package session_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX             context.Context
	IAMOwnerCTX     context.Context
	UserCTX         context.Context
	LoginCTX        context.Context
	Instance        *integration.Instance
	Client          session.SessionServiceClient
	User            *user.AddHumanUserResponse
	DeactivatedUser *user.AddHumanUserResponse
	LockedUser      *user.AddHumanUserResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.SessionV2

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		IAMOwnerCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		UserCTX = Instance.WithAuthorization(ctx, integration.UserTypeNoPermission)
		LoginCTX = Instance.WithAuthorization(ctx, integration.UserTypeLogin)
		User = createFullUser(CTX)
		DeactivatedUser = createDeactivatedUser(CTX)
		LockedUser = createLockedUser(CTX)
		return m.Run()
	}())
}

func createFullUser(ctx context.Context) *user.AddHumanUserResponse {
	userResp := Instance.CreateHumanUser(ctx)
	Instance.Client.UserV2.VerifyEmail(ctx, &user.VerifyEmailRequest{
		UserId:           userResp.GetUserId(),
		VerificationCode: userResp.GetEmailCode(),
	})
	Instance.Client.UserV2.VerifyPhone(ctx, &user.VerifyPhoneRequest{
		UserId:           userResp.GetUserId(),
		VerificationCode: userResp.GetPhoneCode(),
	})
	Instance.SetUserPassword(ctx, userResp.GetUserId(), integration.UserPassword, false)
	Instance.RegisterUserPasskey(ctx, userResp.GetUserId())
	return userResp
}

func createDeactivatedUser(ctx context.Context) *user.AddHumanUserResponse {
	userResp := Instance.CreateHumanUser(ctx)
	_, err := Instance.Client.UserV2.DeactivateUser(ctx, &user.DeactivateUserRequest{UserId: userResp.GetUserId()})
	logging.OnError(err).Fatal("deactivate human user")
	return userResp
}

func createLockedUser(ctx context.Context) *user.AddHumanUserResponse {
	userResp := Instance.CreateHumanUser(ctx)
	_, err := Instance.Client.UserV2.LockUser(ctx, &user.LockUserRequest{UserId: userResp.GetUserId()})
	logging.OnError(err).Fatal("lock human user")
	return userResp
}
