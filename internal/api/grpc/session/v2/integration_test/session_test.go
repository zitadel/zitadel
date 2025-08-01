//go:build integration

package session_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/sink"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func verifyCurrentSession(t *testing.T, id, token string, sequence uint64, window time.Duration, metadata map[string][]byte, userAgent *session.UserAgent, expirationWindow time.Duration, userID string, factors ...wantFactor) *session.Session {
	t.Helper()
	require.NotEmpty(t, id)
	require.NotEmpty(t, token)

	resp, err := Client.GetSession(CTX, &session.GetSessionRequest{
		SessionId:    id,
		SessionToken: &token,
	})
	require.NoError(t, err)
	s := resp.GetSession()
	want := &session.Session{
		Id:        id,
		Sequence:  sequence,
		Metadata:  metadata,
		UserAgent: userAgent,
	}
	verifySession(t, s, want, window, expirationWindow, userID, factors...)
	return s
}

func verifySession(t assert.TestingT, s *session.Session, want *session.Session, window time.Duration, expirationWindow time.Duration, userID string, factors ...wantFactor) {
	assert.Equal(t, want.Id, s.GetId())
	assert.WithinRange(t, s.GetCreationDate().AsTime(), time.Now().Add(-window), time.Now().Add(window))
	assert.WithinRange(t, s.GetChangeDate().AsTime(), time.Now().Add(-window), time.Now().Add(window))
	assert.Equal(t, want.Sequence, s.GetSequence())
	assert.Equal(t, want.Metadata, s.GetMetadata())

	if !proto.Equal(want.UserAgent, s.GetUserAgent()) {
		t.Errorf("user agent =\n%v\nwant\n%v", s.GetUserAgent(), want.UserAgent)
	}
	if expirationWindow == 0 {
		assert.Nil(t, s.GetExpirationDate())
	} else {
		assert.WithinRange(t, s.GetExpirationDate().AsTime(), time.Now().Add(-expirationWindow), time.Now().Add(expirationWindow))
	}

	verifyFactors(t, s.GetFactors(), window, userID, factors)
}

type wantFactor int

const (
	wantUserFactor wantFactor = iota
	wantPasswordFactor
	wantWebAuthNFactor
	wantWebAuthNFactorUserVerified
	wantTOTPFactor
	wantIntentFactor
	wantOTPSMSFactor
	wantOTPEmailFactor
)

func verifyFactors(t assert.TestingT, factors *session.Factors, window time.Duration, userID string, want []wantFactor) {
	for _, w := range want {
		switch w {
		case wantUserFactor:
			uf := factors.GetUser()
			assert.NotNil(t, uf)
			assert.WithinRange(t, uf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
			assert.Equal(t, userID, uf.GetId())
		case wantPasswordFactor:
			pf := factors.GetPassword()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		case wantWebAuthNFactor:
			pf := factors.GetWebAuthN()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
			assert.False(t, pf.GetUserVerified())
		case wantWebAuthNFactorUserVerified:
			pf := factors.GetWebAuthN()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
			assert.True(t, pf.GetUserVerified())
		case wantTOTPFactor:
			pf := factors.GetTotp()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		case wantIntentFactor:
			pf := factors.GetIntent()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		case wantOTPSMSFactor:
			pf := factors.GetOtpSms()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		case wantOTPEmailFactor:
			pf := factors.GetOtpEmail()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		}
	}
}

func TestServer_CreateSession(t *testing.T) {
	tests := []struct {
		name                 string
		req                  *session.CreateSessionRequest
		want                 *session.CreateSessionResponse
		wantErr              bool
		wantFactors          []wantFactor
		wantUserAgent        *session.UserAgent
		wantExpirationWindow time.Duration
	}{
		{
			name: "empty session",
			req: &session.CreateSessionRequest{
				Metadata: map[string][]byte{"foo": []byte("bar")},
			},
			want: &session.CreateSessionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "full session",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					User: &session.CheckUser{
						Search: &session.CheckUser_UserId{
							UserId: User.GetUserId(),
						},
					},
				},
				Metadata: map[string][]byte{"foo": []byte("bar")},
				UserAgent: &session.UserAgent{
					FingerprintId: gu.Ptr("fingerPrintID"),
					Ip:            gu.Ptr("1.2.3.4"),
					Description:   gu.Ptr("Description"),
					Header: map[string]*session.UserAgent_HeaderValues{
						"foo": {Values: []string{"foo", "bar"}},
					},
				},
				Lifetime: durationpb.New(5 * time.Minute),
			},
			want: &session.CreateSessionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "negative lifetime",
			req: &session.CreateSessionRequest{
				Metadata: map[string][]byte{"foo": []byte("bar")},
				Lifetime: durationpb.New(-5 * time.Minute),
			},
			wantErr: true,
		},
		{
			name: "deactivated user",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					User: &session.CheckUser{
						Search: &session.CheckUser_UserId{
							UserId: DeactivatedUser.GetUserId(),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "locked user",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					User: &session.CheckUser{
						Search: &session.CheckUser_UserId{
							UserId: LockedUser.GetUserId(),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "password without user error",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					Password: &session.CheckPassword{
						Password: "Difficult",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "passkey without user error",
			req: &session.CreateSessionRequest{
				Challenges: &session.RequestChallenges{
					WebAuthN: &session.RequestChallenges_WebAuthN{
						Domain:                      Instance.Domain,
						UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "passkey without domain (not registered) error",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					User: &session.CheckUser{
						Search: &session.CheckUser_UserId{
							UserId: User.GetUserId(),
						},
					},
				},
				Challenges: &session.RequestChallenges{
					WebAuthN: &session.RequestChallenges_WebAuthN{
						UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateSession(LoginCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_CreateSession_lock_user(t *testing.T) {
	// create a separate org so we don't interfere with any other test
	org := Instance.CreateOrganization(IAMOwnerCTX,
		fmt.Sprintf("TestServer_CreateSession_lock_user_%s", gofakeit.AppName()),
		gofakeit.Email(),
	)
	userID := org.CreatedAdmins[0].GetUserId()
	Instance.SetUserPassword(IAMOwnerCTX, userID, integration.UserPassword, false)

	// enable password lockout
	maxAttempts := 2
	ctxOrg := metadata.AppendToOutgoingContext(IAMOwnerCTX, "x-zitadel-orgid", org.GetOrganizationId())
	_, err := Instance.Client.Mgmt.AddCustomLockoutPolicy(ctxOrg, &mgmt.AddCustomLockoutPolicyRequest{
		MaxPasswordAttempts: uint32(maxAttempts),
	})
	require.NoError(t, err)

	for i := 0; i <= maxAttempts; i++ {
		_, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
			Checks: &session.Checks{
				User: &session.CheckUser{
					Search: &session.CheckUser_UserId{
						UserId: userID,
					},
				},
				Password: &session.CheckPassword{
					Password: "invalid",
				},
			},
		})
		assert.Error(t, err)
		statusCode := status.Code(err)
		expectedCode := codes.InvalidArgument
		// as soon as we hit the limit the user is locked and following request will
		// already deny any check with a precondition failed since the user is locked
		if i >= maxAttempts {
			expectedCode = codes.FailedPrecondition
		}
		assert.Equal(t, expectedCode, statusCode)
	}
}

func TestServer_CreateSession_webauthn(t *testing.T) {
	// create new session with user and request the webauthn challenge
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
		Challenges: &session.RequestChallenges{
			WebAuthN: &session.RequestChallenges_WebAuthN{
				Domain:                      Instance.Domain,
				UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())

	assertionData, err := Instance.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), true)
	require.NoError(t, err)

	// update the session with webauthn assertion data
	updateResp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			WebAuthN: &session.CheckWebAuthN{
				CredentialAssertionData: assertionData,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), updateResp.GetSessionToken(), updateResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantWebAuthNFactorUserVerified)
}

func TestServer_CreateSession_successfulIntent(t *testing.T) {
	idpID := Instance.AddGenericOAuthProvider(IAMOwnerCTX, gofakeit.AppName()).GetId()
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())

	intentID, token, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, "id", User.GetUserId(), time.Now().Add(time.Hour))
	require.NoError(t, err)
	updateResp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: token,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), updateResp.GetSessionToken(), updateResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantIntentFactor)
}

func TestServer_CreateSession_successfulIntent_instant(t *testing.T) {
	idpID := Instance.AddGenericOAuthProvider(IAMOwnerCTX, gofakeit.AppName()).GetId()

	intentID, token, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, "id", User.GetUserId(), time.Now().Add(time.Hour))
	require.NoError(t, err)
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: token,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantIntentFactor)
}

func TestServer_CreateSession_successfulIntentUnknownUserID(t *testing.T) {
	idpID := Instance.AddGenericOAuthProvider(IAMOwnerCTX, gofakeit.AppName()).GetId()

	// successful intent without known / linked user
	idpUserID := "id"
	intentID, token, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, idpUserID, "", time.Now().Add(time.Hour))

	// link the user (with info from intent)
	Instance.CreateUserIDPlink(CTX, User.GetUserId(), idpUserID, idpID, User.GetUserId())

	// session with intent check must now succeed
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: token,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantIntentFactor)
}

func TestServer_CreateSession_startedIntentFalseToken(t *testing.T) {
	idpID := Instance.AddGenericOAuthProvider(IAMOwnerCTX, gofakeit.AppName()).GetId()

	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())

	intent := Instance.CreateIntent(CTX, idpID)
	_, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intent.GetIdpIntent().GetIdpIntentId(),
				IdpIntentToken: "false",
			},
		},
	})
	require.Error(t, err)
}

func TestServer_CreateSession_reuseIntent(t *testing.T) {
	idpID := Instance.AddGenericOAuthProvider(IAMOwnerCTX, gofakeit.AppName()).GetId()
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())

	intentID, token, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, "id", User.GetUserId(), time.Now().Add(time.Hour))
	require.NoError(t, err)
	updateResp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: token,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), updateResp.GetSessionToken(), updateResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantIntentFactor)

	// the reuse of the intent token is not allowed, not even on the same session
	session2, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: token,
			},
		},
	})
	require.Error(t, err)
	_ = session2
}

func TestServer_CreateSession_expiredIntent(t *testing.T) {
	idpID := Instance.AddGenericOAuthProvider(IAMOwnerCTX, gofakeit.AppName()).GetId()
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())

	intentID, token, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, "id", User.GetUserId(), time.Now().Add(time.Second))
	require.NoError(t, err)

	// wait for the intent to expire
	time.Sleep(2 * time.Second)

	_, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Checks: &session.Checks{
			IdpIntent: &session.CheckIDPIntent{
				IdpIntentId:    intentID,
				IdpIntentToken: token,
			},
		},
	})
	require.Error(t, err)
}

func registerTOTP(ctx context.Context, t *testing.T, userID string) (secret string) {
	resp, err := Instance.Client.UserV2.RegisterTOTP(ctx, &user.RegisterTOTPRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	secret = resp.GetSecret()
	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	_, err = Instance.Client.UserV2.VerifyTOTPRegistration(ctx, &user.VerifyTOTPRegistrationRequest{
		UserId: userID,
		Code:   code,
	})
	require.NoError(t, err)
	return secret
}

func registerOTPSMS(ctx context.Context, t *testing.T, userID string) {
	_, err := Instance.Client.UserV2.AddOTPSMS(ctx, &user.AddOTPSMSRequest{
		UserId: userID,
	})
	require.NoError(t, err)
}

func registerOTPEmail(ctx context.Context, t *testing.T, userID string) {
	_, err := Instance.Client.UserV2.AddOTPEmail(ctx, &user.AddOTPEmailRequest{
		UserId: userID,
	})
	require.NoError(t, err)
}

func TestServer_SetSession_flow_totp(t *testing.T) {
	userExisting := createFullUser(CTX)

	// create new, empty session
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{})
	require.NoError(t, err)
	sessionToken := createResp.GetSessionToken()
	verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, "")

	t.Run("check user", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				User: &session.CheckUser{
					Search: &session.CheckUser_UserId{
						UserId: userExisting.GetUserId(),
					},
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, userExisting.GetUserId(), wantUserFactor)
	})

	t.Run("check webauthn, user verified (passkey)", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Challenges: &session.RequestChallenges{
				WebAuthN: &session.RequestChallenges_WebAuthN{
					Domain:                      Instance.Domain,
					UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
				},
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, userExisting.GetUserId())
		sessionToken = resp.GetSessionToken()

		assertionData, err := Instance.WebAuthN.CreateAssertionResponse(resp.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), true)
		require.NoError(t, err)

		resp, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				WebAuthN: &session.CheckWebAuthN{
					CredentialAssertionData: assertionData,
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, userExisting.GetUserId(), wantUserFactor, wantWebAuthNFactorUserVerified)
	})

	userAuthCtx := integration.WithAuthorizationToken(CTX, sessionToken)
	Instance.RegisterUserU2F(userAuthCtx, userExisting.GetUserId())
	totpSecret := registerTOTP(userAuthCtx, t, userExisting.GetUserId())
	registerOTPSMS(userAuthCtx, t, userExisting.GetUserId())
	registerOTPEmail(userAuthCtx, t, userExisting.GetUserId())

	t.Run("check TOTP", func(t *testing.T) {
		code, err := totp.GenerateCode(totpSecret, time.Now())
		require.NoError(t, err)
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				Totp: &session.CheckTOTP{
					Code: code,
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, userExisting.GetUserId(), wantUserFactor, wantTOTPFactor)
	})

	userImport := Instance.CreateHumanUserWithTOTP(CTX, totpSecret)
	createRespImport, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{})
	require.NoError(t, err)
	sessionTokenImport := createRespImport.GetSessionToken()
	verifyCurrentSession(t, createRespImport.GetSessionId(), sessionTokenImport, createRespImport.GetDetails().GetSequence(), time.Minute, nil, nil, 0, "")

	t.Run("check user", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createRespImport.GetSessionId(),
			Checks: &session.Checks{
				User: &session.CheckUser{
					Search: &session.CheckUser_UserId{
						UserId: userImport.GetUserId(),
					},
				},
			},
		})
		require.NoError(t, err)
		sessionTokenImport = resp.GetSessionToken()
		verifyCurrentSession(t, createRespImport.GetSessionId(), sessionTokenImport, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, userImport.GetUserId(), wantUserFactor)
	})
	t.Run("check TOTP", func(t *testing.T) {
		code, err := totp.GenerateCode(totpSecret, time.Now())
		require.NoError(t, err)
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createRespImport.GetSessionId(),
			Checks: &session.Checks{
				Totp: &session.CheckTOTP{
					Code: code,
				},
			},
		})
		require.NoError(t, err)
		sessionTokenImport = resp.GetSessionToken()
		verifyCurrentSession(t, createRespImport.GetSessionId(), sessionTokenImport, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, userImport.GetUserId(), wantUserFactor, wantTOTPFactor)
	})
}

func TestServer_SetSession_flow(t *testing.T) {
	// create new, empty session
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{})
	require.NoError(t, err)
	sessionToken := createResp.GetSessionToken()
	verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, createResp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())

	t.Run("check user", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				User: &session.CheckUser{
					Search: &session.CheckUser_UserId{
						UserId: User.GetUserId(),
					},
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor)
	})

	t.Run("check webauthn, user verified (passkey)", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Challenges: &session.RequestChallenges{
				WebAuthN: &session.RequestChallenges_WebAuthN{
					Domain:                      Instance.Domain,
					UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
				},
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())
		sessionToken = resp.GetSessionToken()

		assertionData, err := Instance.WebAuthN.CreateAssertionResponse(resp.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), true)
		require.NoError(t, err)

		resp, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				WebAuthN: &session.CheckWebAuthN{
					CredentialAssertionData: assertionData,
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantWebAuthNFactorUserVerified)
	})

	userAuthCtx := integration.WithAuthorizationToken(CTX, sessionToken)
	Instance.RegisterUserU2F(userAuthCtx, User.GetUserId())
	totpSecret := registerTOTP(userAuthCtx, t, User.GetUserId())
	registerOTPSMS(userAuthCtx, t, User.GetUserId())
	registerOTPEmail(userAuthCtx, t, User.GetUserId())

	t.Run("check webauthn, user not verified (U2F)", func(t *testing.T) {

		for _, userVerificationRequirement := range []session.UserVerificationRequirement{
			session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED,
			session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_DISCOURAGED,
		} {
			t.Run(userVerificationRequirement.String(), func(t *testing.T) {
				resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
					SessionId: createResp.GetSessionId(),
					Challenges: &session.RequestChallenges{
						WebAuthN: &session.RequestChallenges_WebAuthN{
							Domain:                      Instance.Domain,
							UserVerificationRequirement: userVerificationRequirement,
						},
					},
				})
				require.NoError(t, err)
				verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())
				sessionToken = resp.GetSessionToken()

				assertionData, err := Instance.WebAuthN.CreateAssertionResponse(resp.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), false)
				require.NoError(t, err)

				resp, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
					SessionId: createResp.GetSessionId(),
					Checks: &session.Checks{
						WebAuthN: &session.CheckWebAuthN{
							CredentialAssertionData: assertionData,
						},
					},
				})
				require.NoError(t, err)
				sessionToken = resp.GetSessionToken()
				verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantWebAuthNFactor)
			})
		}
	})

	t.Run("check TOTP", func(t *testing.T) {
		code, err := totp.GenerateCode(totpSecret, time.Now())
		require.NoError(t, err)
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				Totp: &session.CheckTOTP{
					Code: code,
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantWebAuthNFactor, wantTOTPFactor)
	})

	t.Run("check OTP SMS", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Challenges: &session.RequestChallenges{
				OtpSms: &session.RequestChallenges_OTPSMS{ReturnCode: true},
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())
		sessionToken = resp.GetSessionToken()

		otp := resp.GetChallenges().GetOtpSms()
		require.NotEmpty(t, otp)

		resp, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				OtpSms: &session.CheckOTP{
					Code: otp,
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantWebAuthNFactor, wantOTPSMSFactor)
	})

	t.Run("check OTP Email", func(t *testing.T) {
		resp, err := Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Challenges: &session.RequestChallenges{
				OtpEmail: &session.RequestChallenges_OTPEmail{
					DeliveryType: &session.RequestChallenges_OTPEmail_ReturnCode_{},
				},
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId())
		sessionToken = resp.GetSessionToken()

		otp := resp.GetChallenges().GetOtpEmail()
		require.NotEmpty(t, otp)

		resp, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
			SessionId: createResp.GetSessionId(),
			Checks: &session.Checks{
				OtpEmail: &session.CheckOTP{
					Code: otp,
				},
			},
		})
		require.NoError(t, err)
		sessionToken = resp.GetSessionToken()
		verifyCurrentSession(t, createResp.GetSessionId(), sessionToken, resp.GetDetails().GetSequence(), time.Minute, nil, nil, 0, User.GetUserId(), wantUserFactor, wantWebAuthNFactor, wantOTPEmailFactor)
	})
}

func TestServer_SetSession_expired(t *testing.T) {
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Lifetime: durationpb.New(20 * time.Second),
	})
	require.NoError(t, err)

	// test session token works
	_, err = Instance.Client.SessionV2.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Lifetime:  durationpb.New(20 * time.Second),
	})
	require.NoError(t, err)

	// ensure session expires and does not work anymore
	time.Sleep(20 * time.Second)
	_, err = Client.SetSession(LoginCTX, &session.SetSessionRequest{
		SessionId: createResp.GetSessionId(),
		Lifetime:  durationpb.New(20 * time.Second),
	})
	require.Error(t, err)
}

func TestServer_DeleteSession_token(t *testing.T) {
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{})
	require.NoError(t, err)

	_, err = Client.DeleteSession(CTX, &session.DeleteSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: gu.Ptr("invalid"),
	})
	require.Error(t, err)

	_, err = Client.DeleteSession(CTX, &session.DeleteSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: gu.Ptr(createResp.GetSessionToken()),
	})
	require.NoError(t, err)
}

func TestServer_DeleteSession_own_session(t *testing.T) {
	// create two users for the test and a session each to get tokens for authorization
	user1 := Instance.CreateHumanUser(CTX)
	Instance.SetUserPassword(CTX, user1.GetUserId(), integration.UserPassword, false)
	_, token1, _, _ := Instance.CreatePasswordSession(t, LoginCTX, user1.GetUserId(), integration.UserPassword)

	user2 := Instance.CreateHumanUser(CTX)
	Instance.SetUserPassword(CTX, user2.GetUserId(), integration.UserPassword, false)
	_, token2, _, _ := Instance.CreatePasswordSession(t, LoginCTX, user2.GetUserId(), integration.UserPassword)

	// create a new session for the first user
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: user1.GetUserId(),
				},
			},
		},
	})
	require.NoError(t, err)

	// delete the new (user1) session must not be possible with user (has no permission)
	_, err = Client.DeleteSession(integration.WithAuthorizationToken(context.Background(), token2), &session.DeleteSessionRequest{
		SessionId: createResp.GetSessionId(),
	})
	require.Error(t, err)

	// delete the new (user1) session by themselves
	_, err = Client.DeleteSession(integration.WithAuthorizationToken(context.Background(), token1), &session.DeleteSessionRequest{
		SessionId: createResp.GetSessionId(),
	})
	require.NoError(t, err)
}

func TestServer_DeleteSession_with_permission(t *testing.T) {
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
	})
	require.NoError(t, err)

	// delete the new session by ORG_OWNER
	_, err = Client.DeleteSession(Instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner), &session.DeleteSessionRequest{
		SessionId: createResp.GetSessionId(),
	})
	require.NoError(t, err)
}

func Test_ZITADEL_API_missing_authentication(t *testing.T) {
	// create new, empty session
	createResp, err := Client.CreateSession(LoginCTX, &session.CreateSessionRequest{})
	require.NoError(t, err)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("Bearer %s", createResp.GetSessionToken()))
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		sessionResp, err := Client.GetSession(ctx, &session.GetSessionRequest{SessionId: createResp.GetSessionId()})
		if !assert.Error(tt, err) {
			return
		}
		assert.Nil(tt, sessionResp)
	}, retryDuration, tick)
}

func Test_ZITADEL_API_success(t *testing.T) {
	id, token, _, _ := Instance.CreateVerifiedWebAuthNSession(t, LoginCTX, User.GetUserId())
	ctx := integration.WithAuthorizationToken(context.Background(), token)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		sessionResp, err := Client.GetSession(ctx, &session.GetSessionRequest{SessionId: id})
		if !assert.NoError(tt, err) {
			return
		}
		webAuthN := sessionResp.GetSession().GetFactors().GetWebAuthN()
		assert.NotNil(tt, id, webAuthN.GetVerifiedAt().AsTime())
		assert.True(tt, webAuthN.GetUserVerified())
	}, retryDuration, tick)
}

func Test_ZITADEL_API_session_not_found(t *testing.T) {
	id, token, _, _ := Instance.CreateVerifiedWebAuthNSession(t, LoginCTX, User.GetUserId())

	// test session token works
	ctx := integration.WithAuthorizationToken(context.Background(), token)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		_, err := Client.GetSession(ctx, &session.GetSessionRequest{SessionId: id})
		if !assert.NoError(tt, err) {
			return
		}
	}, retryDuration, tick)

	//terminate the session and test it does not work anymore
	_, err := Client.DeleteSession(CTX, &session.DeleteSessionRequest{
		SessionId:    id,
		SessionToken: gu.Ptr(token),
	})
	require.NoError(t, err)

	ctx = integration.WithAuthorizationToken(context.Background(), token)
	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		_, err := Client.GetSession(ctx, &session.GetSessionRequest{SessionId: id})
		if !assert.Error(tt, err) {
			return
		}
	}, retryDuration, tick)
}

func Test_ZITADEL_API_session_expired(t *testing.T) {
	id, token, _, _ := Instance.CreateVerifiedWebAuthNSessionWithLifetime(t, LoginCTX, User.GetUserId(), 20*time.Second)

	// test session token works
	ctx := integration.WithAuthorizationToken(context.Background(), token)
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		_, err := Client.GetSession(ctx, &session.GetSessionRequest{SessionId: id})
		if !assert.NoError(tt, err) {
			return
		}
	}, retryDuration, tick)

	// ensure session expires and does not work anymore
	time.Sleep(20 * time.Second)
	sessionResp, err := Client.GetSession(ctx, &session.GetSessionRequest{SessionId: id})
	require.Error(t, err)
	require.Nil(t, sessionResp)
}
