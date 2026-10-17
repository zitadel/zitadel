package command

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/action/otp"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/user"
)

var (
	hookTestDefaultGenerators = &SecretGenerators{
		OTPSMS: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
		OTPEmail: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
	}

	hookTestAuthRequest = &domain.AuthRequest{
		ID:      "authRequestID",
		AgentID: "userAgentID",
		BrowserInfo: &domain.BrowserInfo{
			UserAgent:      "user-agent",
			AcceptLanguage: "en",
			RemoteIP:       net.IP{192, 0, 2, 1},
		},
	}

	hookTestAuthRequestInfo = &user.AuthRequestInfo{
		ID:          "authRequestID",
		UserAgentID: "userAgentID",
		BrowserInfo: &user.BrowserInfo{
			UserAgent:      "user-agent",
			AcceptLanguage: "en",
			RemoteIP:       net.IP{192, 0, 2, 1},
		},
	}
)

func flagOnCtx() context.Context {
	return authz.NewMockContext("inst1", "org1", "user1",
		authz.WithMockFeatures(feature.Features{AllowOTPCodeOverride: true}))
}

func TestCommandSide_HumanSendOTPSMS_hookReturnsCode(t *testing.T) {
	ctx := flagOnCtx()

	r := &Commands{
		eventstore: expectEventstore(
			expectFilter(
				eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
			),
			expectFilter(), // lastActivated SMS config — none
			expectFilter(), // getSMSConfig — none
			expectFilter(), // secret generator config — fall back to default
			expectPush(
				user.NewHumanOTPSMSCodeAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate,
					&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("A7F2B9"),
					},
					5*time.Minute,
					hookTestAuthRequestInfo,
					"",
				),
			),
		)(t),
		defaultSecretGenerators: hookTestDefaultGenerators,
		userEncryption:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
		preOTPSMSCodeHook: func(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPSMSCodeResponse, error) {
			return &otp.PreOTPSMSCodeResponse{
				Code:   gu.Ptr("A7F2B9"),
				Expiry: gu.Ptr(otp.Duration(5 * time.Minute)),
			}, nil
		},
	}
	assert.NoError(t, r.HumanSendOTPSMS(ctx, "user1", "org1", hookTestAuthRequest))
}

func TestCommandSide_HumanSendOTPSMS_flagOffHookNotInvoked(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1") // no features

	hookCalled := false
	r := &Commands{
		eventstore: expectEventstore(
			expectFilter(
				eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
			),
			expectFilter(), // lastActivated SMS config — none
			expectFilter(), // getSMSConfig — none
			expectPush(
				user.NewHumanOTPSMSCodeAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate,
					&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("12345678"),
					},
					time.Hour,
					hookTestAuthRequestInfo,
					"",
				),
			),
		)(t),
		defaultSecretGenerators:     hookTestDefaultGenerators,
		newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
		preOTPSMSCodeHook: func(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPSMSCodeResponse, error) {
			hookCalled = true
			return nil, nil
		},
	}
	assert.NoError(t, r.HumanSendOTPSMS(ctx, "user1", "org1", hookTestAuthRequest))
	assert.False(t, hookCalled, "preOTPSMSCodeHook must not be invoked when AllowOTPCodeOverride is off")
}

func TestCommandSide_HumanSendOTPSMS_flagOnNilHookFallsThrough(t *testing.T) {
	ctx := flagOnCtx()

	r := &Commands{
		eventstore: expectEventstore(
			expectFilter(
				eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
			),
			expectFilter(), // lastActivated SMS config — none
			expectFilter(), // getSMSConfig — none
			expectPush(
				user.NewHumanOTPSMSCodeAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate,
					&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("12345678"),
					},
					time.Hour,
					hookTestAuthRequestInfo,
					"",
				),
			),
		)(t),
		defaultSecretGenerators:     hookTestDefaultGenerators,
		newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
		// preOTPSMSCodeHook left nil on purpose.
	}
	assert.NoError(t, r.HumanSendOTPSMS(ctx, "user1", "org1", hookTestAuthRequest))
}

func TestCommandSide_HumanSendOTPSMS_hookError(t *testing.T) {
	ctx := flagOnCtx()

	hookErr := errors.New("hook target unreachable")
	r := &Commands{
		eventstore: expectEventstore(
			expectFilter(
				eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
			),
			expectFilter(), // lastActivated SMS config — none
			expectFilter(), // getSMSConfig — none
			expectFilter(), // secret generator config — fall back to default
		)(t),
		defaultSecretGenerators: hookTestDefaultGenerators,
		userEncryption:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
		preOTPSMSCodeHook: func(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPSMSCodeResponse, error) {
			return nil, hookErr
		},
	}
	err := r.HumanSendOTPSMS(ctx, "user1", "org1", hookTestAuthRequest)
	assert.ErrorIs(t, err, hookErr)
}

func TestCommandSide_HumanSendOTPEmail_hookReturnsCode(t *testing.T) {
	ctx := flagOnCtx()

	r := &Commands{
		eventstore: expectEventstore(
			expectFilter(
				eventFromEventPusher(user.NewHumanEmailVerifiedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
				eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
			),
			expectFilter(), // secret generator config — fall back to default
			expectPush(
				user.NewHumanOTPEmailCodeAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate,
					&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("B8E3C1"),
					},
					10*time.Minute,
					hookTestAuthRequestInfo,
					"",
				),
			),
		)(t),
		defaultSecretGenerators: hookTestDefaultGenerators,
		userEncryption:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
		preOTPEmailCodeHook: func(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPEmailCodeResponse, error) {
			return &otp.PreOTPEmailCodeResponse{
				Code:   gu.Ptr("B8E3C1"),
				Expiry: gu.Ptr(otp.Duration(10 * time.Minute)),
			}, nil
		},
	}
	assert.NoError(t, r.HumanSendOTPEmail(ctx, "user1", "org1", hookTestAuthRequest))
}

func TestCommandSide_HumanSendOTPEmail_flagOffHookNotInvoked(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1") // no features

	hookCalled := false
	r := &Commands{
		eventstore: expectEventstore(
			expectFilter(
				eventFromEventPusher(user.NewHumanEmailVerifiedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
				eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate)),
			),
			expectPush(
				user.NewHumanOTPEmailCodeAddedEvent(ctx,
					&user.NewAggregate("user1", "org1").Aggregate,
					&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("12345678"),
					},
					time.Hour,
					hookTestAuthRequestInfo,
					"",
				),
			),
		)(t),
		defaultSecretGenerators:     hookTestDefaultGenerators,
		newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
		preOTPEmailCodeHook: func(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPEmailCodeResponse, error) {
			hookCalled = true
			return nil, nil
		},
	}
	assert.NoError(t, r.HumanSendOTPEmail(ctx, "user1", "org1", hookTestAuthRequest))
	assert.False(t, hookCalled, "preOTPEmailCodeHook must not be invoked when AllowOTPCodeOverride is off")
}
