//go:build integration

package events_test

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/crypto"
	zdomain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/sink"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_SessionReduces(t *testing.T) {
	instanceID := Instance.ID()
	sessionRepo := repository.SessionRepository()

	email := integration.Email()
	testUser := Instance.CreateUserTypeHuman(CTX, email)
	_, err := UserClient.VerifyEmail(CTX, &user.VerifyEmailRequest{
		UserId:           testUser.GetId(),
		VerificationCode: testUser.GetEmailCode(),
	})
	require.NoError(t, err)
	_, err = UserClient.SetPassword(CTX, &user.SetPasswordRequest{
		UserId: testUser.GetId(),
		NewPassword: &user.Password{
			Password: integration.UserPassword,
		},
		Verification: nil,
	})
	require.NoError(t, err)
	totpSecret := registerTOTP(CTX, t, testUser.GetId())
	Instance.RegisterUserPasskey(CTX, testUser.GetId())
	registerOTPEmail(CTX, t, testUser.GetId())
	_, err = UserClient.SetPhone(CTX, &user.SetPhoneRequest{
		UserId: testUser.GetId(),
		Phone:  integration.Phone(),
		Verification: &user.SetPhoneRequest_IsVerified{
			IsVerified: true,
		},
	})
	require.NoError(t, err)
	registerOTPSMS(CTX, t, testUser.GetId())

	idpID := Instance.AddGenericOAuthProvider(CTX, integration.IDPName()).GetId()
	idpUserID := integration.ID()
	_, err = UserClient.AddIDPLink(CTX, &user.AddIDPLinkRequest{
		UserId: testUser.GetId(),
		IdpLink: &user.IDPLink{
			IdpId:    idpID,
			UserId:   idpUserID,
			UserName: integration.Username(),
		},
	})
	require.NoError(t, err)

	recoveryCodes := createRecoveryCodes(CTX, t, testUser.GetId())

	lifetime := time.Hour
	creatorID := "tester"
	userAgent := &domain.SessionUserAgent{
		FingerprintID: gu.Ptr(integration.ID()),
		Description:   gu.Ptr("description"),
		IP:            net.IPv4(127, 0, 0, 1),
		Header: http.Header{
			"User-Agent": []string{"ZITADEL-Integration-Test"},
		},
	}
	createdSession, err := SessionClient.CreateSession(CTX, &session.CreateSessionRequest{
		Checks:     nil,
		Metadata:   nil,
		Challenges: nil,
		UserAgent: &session.UserAgent{
			FingerprintId: userAgent.FingerprintID,
			Ip:            gu.Ptr(userAgent.IP.String()),
			Description:   userAgent.Description,
			Header: map[string]*session.UserAgent_HeaderValues{
				"User-Agent": {Values: []string{"ZITADEL-Integration-Test"}},
			},
		},
		Lifetime: durationpb.New(lifetime),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Second)
	var updatedSession *session.SetSessionResponse

	t.Run("create session reduces", func(t *testing.T) {
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.NoError(collect, err)
			assertSessionsEqual(collect, &domain.Session{
				InstanceID: instanceID,
				ID:         createdSession.GetSessionId(),
				TokenID:    createdSession.GetSessionToken(),
				Lifetime:   lifetime,
				Expiration: createdSession.GetDetails().GetChangeDate().AsTime().Add(lifetime),
				UserID:     "",
				CreatorID:  creatorID,
				CreatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				UpdatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				Factors:    nil,
				Challenges: nil,
				Metadata:   nil,
				UserAgent:  userAgent,
			}, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	metadata := []*domain.SessionMetadata{
		{
			Metadata: domain.Metadata{
				InstanceID: instanceID,
				Key:        "key1",
				Value:      []byte("value1"),
			},
			SessionID: createdSession.GetSessionId(),
		},
	}

	t.Run("metadata set reduces", func(t *testing.T) {
		updatedSession, err = SessionClient.SetSession(CTX, &session.SetSessionRequest{
			SessionId: createdSession.GetSessionId(),
			Checks:    nil,
			Metadata:  map[string][]byte{"key1": []byte("value1")},
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.NoError(collect, err)
			assertSessionsEqual(collect, &domain.Session{
				InstanceID: instanceID,
				ID:         createdSession.GetSessionId(),
				TokenID:    createdSession.GetSessionToken(),
				Lifetime:   lifetime,
				Expiration: createdSession.GetDetails().GetChangeDate().AsTime().Add(lifetime),
				UserID:     "",
				CreatorID:  creatorID,
				CreatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				UpdatedAt:  updatedSession.GetDetails().GetChangeDate().AsTime(),
				Factors:    nil,
				Challenges: nil,
				Metadata:   metadata,
				UserAgent:  userAgent,
			}, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	var lifetimeUpdated time.Time

	t.Run("lifetime set reduces", func(t *testing.T) {
		updatedSession, err = SessionClient.SetSession(CTX, &session.SetSessionRequest{
			SessionId: createdSession.GetSessionId(),
			Lifetime:  durationpb.New(lifetime),
		})
		require.NoError(t, err)
		lifetimeUpdated = updatedSession.GetDetails().GetChangeDate().AsTime()

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.NoError(collect, err)
			assertSessionsEqual(collect, &domain.Session{
				InstanceID: instanceID,
				ID:         createdSession.GetSessionId(),
				TokenID:    createdSession.GetSessionToken(),
				Lifetime:   lifetime,
				Expiration: lifetimeUpdated.Add(lifetime),
				UserID:     "",
				CreatorID:  creatorID,
				CreatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				UpdatedAt:  updatedSession.GetDetails().GetChangeDate().AsTime(),
				Factors:    nil,
				Challenges: nil,
				Metadata:   metadata,
				UserAgent:  userAgent,
			}, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	var firstFactorCheckTime time.Time

	t.Run("checks set reduces", func(t *testing.T) {
		_, err := url.Parse(Instance.CreateIntent(CTX, idpID).GetAuthUrl())
		require.NoError(t, err)
		expiry := time.Now().Add(1 * time.Hour)
		intentID, intentToken, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), idpID, idpUserID, testUser.GetId(), expiry)
		require.NoError(t, err)
		updatedSession, err = SessionClient.SetSession(CTX, &session.SetSessionRequest{
			SessionId: createdSession.GetSessionId(),
			Checks: &session.Checks{
				User: &session.CheckUser{
					Search: &session.CheckUser_LoginName{
						LoginName: email,
					},
				},
				Password: &session.CheckPassword{
					Password: integration.UserPassword,
				},
				IdpIntent: &session.CheckIDPIntent{
					IdpIntentId:    intentID,
					IdpIntentToken: intentToken,
				},
				Totp: &session.CheckTOTP{
					Code: func() string {
						code, err := totp.GenerateCode(totpSecret, time.Now())
						require.NoError(t, err)
						return code
					}(),
				},
				RecoveryCode: &session.CheckRecoveryCode{
					Code: recoveryCodes[0],
				},
			},
			Challenges: &session.RequestChallenges{
				WebAuthN: &session.RequestChallenges_WebAuthN{
					Domain:                      Instance.Domain,
					UserVerificationRequirement: session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED,
				},
				OtpEmail: &session.RequestChallenges_OTPEmail{
					DeliveryType: &session.RequestChallenges_OTPEmail_ReturnCode_{
						ReturnCode: &session.RequestChallenges_OTPEmail_ReturnCode{},
					},
				},
				OtpSms: &session.RequestChallenges_OTPSMS{
					ReturnCode: true,
				},
			},
		})
		require.NoError(t, err)
		firstFactorCheckTime = updatedSession.GetDetails().GetChangeDate().AsTime()

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.NoError(collect, err)
			assertSessionsEqual(collect, &domain.Session{
				InstanceID: instanceID,
				ID:         createdSession.GetSessionId(),
				TokenID:    createdSession.GetSessionToken(),
				Lifetime:   lifetime,
				Expiration: lifetimeUpdated.Add(lifetime),
				UserID:     testUser.GetId(),
				CreatorID:  creatorID,
				CreatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				UpdatedAt:  updatedSession.GetDetails().GetChangeDate().AsTime(),
				Factors: domain.SessionFactors{
					&domain.SessionFactorUser{
						UserID:         testUser.GetId(),
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorPassword{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorIdentityProviderIntent{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorTOTP{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorRecoveryCode{
						LastVerifiedAt: firstFactorCheckTime,
					},
				},
				Challenges: domain.SessionChallenges{
					&domain.SessionChallengePasskey{
						LastChallengedAt:     updatedSession.GetDetails().GetChangeDate().AsTime(),
						Challenge:            "challenge", // placeholder
						AllowedCredentialIDs: [][]byte{Instance.WebAuthN.KeyID()},
						UserVerification:     zdomain.UserVerificationRequirementRequired,
						RPID:                 Instance.Domain,
					},
					&domain.SessionChallengeOTPEmail{
						LastChallengedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "",
							KeyID:      "",
							Crypted:    nil,
						},
						Expiry:            0,
						CodeReturned:      true,
						URLTemplate:       "",
						TriggeredAtOrigin: "",
					},
					&domain.SessionChallengeOTPSMS{
						LastChallengedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "",
							KeyID:      "",
							Crypted:    nil,
						},
						Expiry:            0,
						CodeReturned:      true,
						GeneratorID:       "",
						TriggeredAtOrigin: "",
					},
				},
				Metadata:  metadata,
				UserAgent: userAgent,
			}, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	t.Run("checks with challenges set reduces", func(t *testing.T) {
		assertionData, err := Instance.WebAuthN.CreateAssertionResponse(updatedSession.GetChallenges().GetWebAuthN().GetPublicKeyCredentialRequestOptions(), true)
		require.NoError(t, err)

		updatedSession, err = SessionClient.SetSession(CTX, &session.SetSessionRequest{
			SessionId: createdSession.GetSessionId(),
			Checks: &session.Checks{
				WebAuthN: &session.CheckWebAuthN{
					CredentialAssertionData: assertionData,
				},
				OtpEmail: &session.CheckOTP{
					Code: updatedSession.GetChallenges().GetOtpEmail(),
				},
				OtpSms: &session.CheckOTP{
					Code: updatedSession.GetChallenges().GetOtpSms(),
				},
			},
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.NoError(collect, err)
			assertSessionsEqual(collect, &domain.Session{
				InstanceID: instanceID,
				ID:         createdSession.GetSessionId(),
				TokenID:    createdSession.GetSessionToken(),
				Lifetime:   lifetime,
				Expiration: lifetimeUpdated.Add(lifetime),
				UserID:     testUser.GetId(),
				CreatorID:  creatorID,
				CreatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				UpdatedAt:  updatedSession.GetDetails().GetChangeDate().AsTime(),
				Factors: domain.SessionFactors{
					&domain.SessionFactorUser{
						UserID:         testUser.GetId(),
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorPassword{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorIdentityProviderIntent{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorPasskey{
						LastVerifiedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
						UserVerified:   true,
					},
					&domain.SessionFactorTOTP{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorOTPEmail{
						LastVerifiedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
					},
					&domain.SessionFactorOTPSMS{
						LastVerifiedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
					},
					&domain.SessionFactorRecoveryCode{
						LastVerifiedAt: firstFactorCheckTime,
					},
				},
				Challenges: domain.SessionChallenges{},
				Metadata:   metadata,
				UserAgent:  userAgent,
			}, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	t.Run("password changed reduces", func(t *testing.T) {
		passwordSet, err := UserClient.SetPassword(CTX, &user.SetPasswordRequest{
			UserId: testUser.GetId(),
			NewPassword: &user.Password{
				Password: integration.UserPassword + integration.RandString(5),
			},
			Verification: &user.SetPasswordRequest_CurrentPassword{
				CurrentPassword: integration.UserPassword,
			},
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.NoError(collect, err)
			assertSessionsEqual(collect, &domain.Session{
				InstanceID: instanceID,
				ID:         createdSession.GetSessionId(),
				TokenID:    createdSession.GetSessionToken(),
				Lifetime:   lifetime,
				Expiration: lifetimeUpdated.Add(lifetime),
				UserID:     testUser.GetId(),
				CreatorID:  creatorID,
				CreatedAt:  createdSession.GetDetails().GetChangeDate().AsTime(),
				UpdatedAt:  passwordSet.GetDetails().GetChangeDate().AsTime(),
				Factors: domain.SessionFactors{
					&domain.SessionFactorUser{
						UserID:         testUser.GetId(),
						LastVerifiedAt: firstFactorCheckTime,
					},
					// password factor removed
					&domain.SessionFactorIdentityProviderIntent{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorPasskey{
						LastVerifiedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
						UserVerified:   true,
					},
					&domain.SessionFactorTOTP{
						LastVerifiedAt: firstFactorCheckTime,
					},
					&domain.SessionFactorOTPEmail{
						LastVerifiedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
					},
					&domain.SessionFactorOTPSMS{
						LastVerifiedAt: updatedSession.GetDetails().GetChangeDate().AsTime(),
					},
					&domain.SessionFactorRecoveryCode{
						LastVerifiedAt: firstFactorCheckTime,
					},
				},
				Challenges: domain.SessionChallenges{},
				Metadata:   metadata,
				UserAgent:  userAgent,
			}, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	t.Run("terminated reduces", func(t *testing.T) {
		_, err := SessionClient.DeleteSession(CTX, &session.DeleteSessionRequest{
			SessionId:    createdSession.GetSessionId(),
			SessionToken: gu.Ptr(updatedSession.GetSessionToken()),
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbSession, err := sessionRepo.Get(CTX, pool, database.WithCondition(
				sessionRepo.PrimaryKeyCondition(instanceID, createdSession.GetSessionId()),
			))
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
			assert.Nil(collect, dbSession)
		}, retryDuration, tick, "session not found within %v: %v", retryDuration, err)
	})

	// TODO (@grvijayan): add tests for session reducers based on user deactivated, locked, deleted events
}

func assertSessionsEqual(t *assert.CollectT, expected, actual *domain.Session) {
	t.Helper()
	assert.Equal(t, expected.InstanceID, actual.InstanceID)
	assert.Equal(t, expected.ID, actual.ID)
	assert.NotNil(t, actual.TokenID)
	assert.Equal(t, expected.Lifetime, actual.Lifetime)
	assert.Equal(t, expected.Expiration.UTC(), actual.Expiration.UTC())
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.CreatorID, actual.CreatorID)
	assert.NotNil(t, actual.CreatedAt)
	assert.NotNil(t, actual.UpdatedAt)
	assertSessionFactorsEqual(t, expected.Factors, actual.Factors)
	assertSessionChallengesEqual(t, expected.Challenges, actual.Challenges)
	assertSessionMetadataEqual(t, expected.Metadata, actual.Metadata)
	assert.Equal(t, expected.UserAgent, actual.UserAgent)
}

func assertSessionMetadataEqual(t *assert.CollectT, expected, actual []*domain.SessionMetadata) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for _, exp := range expected {
		found := false
		for _, act := range actual {
			if exp.InstanceID == act.InstanceID &&
				exp.SessionID == act.SessionID &&
				exp.Key == act.Key &&
				bytes.Equal(exp.Value, act.Value) {
				found = true
				break
			}
		}
		assert.Truef(t, found, "expected metadata not found: %+v", exp)
	}
}

func assertSessionFactorsEqual(t *assert.CollectT, expected, actual domain.SessionFactors) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for _, exp := range expected {
		found := false
		for _, act := range actual {
			if exp.SessionFactorType() == act.SessionFactorType() {
				switch expTyped := exp.(type) {
				case *domain.SessionFactorUser:
					actTyped := act.(*domain.SessionFactorUser)
					assert.Equal(t, expTyped.UserID, actTyped.UserID)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				case *domain.SessionFactorPassword:
					actTyped := act.(*domain.SessionFactorPassword)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				case *domain.SessionFactorIdentityProviderIntent:
					actTyped := act.(*domain.SessionFactorIdentityProviderIntent)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				case *domain.SessionFactorPasskey:
					actTyped := act.(*domain.SessionFactorPasskey)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
					assert.Equal(t, expTyped.UserVerified, actTyped.UserVerified)
				case *domain.SessionFactorTOTP:
					actTyped := act.(*domain.SessionFactorTOTP)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				case *domain.SessionFactorOTPEmail:
					actTyped := act.(*domain.SessionFactorOTPEmail)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				case *domain.SessionFactorOTPSMS:
					actTyped := act.(*domain.SessionFactorOTPSMS)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				case *domain.SessionFactorRecoveryCode:
					actTyped := act.(*domain.SessionFactorRecoveryCode)
					assert.Equal(t, expTyped.LastVerifiedAt.UTC(), actTyped.LastVerifiedAt.UTC())
				}
				found = true
				break
			}
		}
		assert.Truef(t, found, "expected factor not found: %+v", exp)
	}
}

func assertSessionChallengesEqual(t *assert.CollectT, expected, actual domain.SessionChallenges) {
	t.Helper()
	assert.Len(t, actual, len(expected))
	for _, exp := range expected {
		found := false
		for _, act := range actual {
			if exp.SessionChallengeType() == act.SessionChallengeType() {
				switch expTyped := exp.(type) {
				case *domain.SessionChallengePasskey:
					actTyped := act.(*domain.SessionChallengePasskey)
					assert.Equal(t, expTyped.LastChallengedAt.UTC(), actTyped.LastChallengedAt.UTC())
					assert.NotEmpty(t, actTyped.Challenge)
					assert.Equal(t, expTyped.AllowedCredentialIDs, actTyped.AllowedCredentialIDs)
					assert.Equal(t, expTyped.RPID, actTyped.RPID)
					assert.Equal(t, expTyped.UserVerification, actTyped.UserVerification)
				case *domain.SessionChallengeOTPEmail:
					actTyped := act.(*domain.SessionChallengeOTPEmail)
					assert.Equal(t, expTyped.LastChallengedAt.UTC(), actTyped.LastChallengedAt.UTC())
					assert.Equal(t, expTyped.Code.CryptoType, actTyped.Code.CryptoType)
				case *domain.SessionChallengeOTPSMS:
					actTyped := act.(*domain.SessionChallengeOTPSMS)
					assert.Equal(t, expTyped.LastChallengedAt.UTC(), actTyped.LastChallengedAt.UTC())
					assert.Equal(t, expTyped.Code.CryptoType, actTyped.Code.CryptoType)
				}
				found = true
				break
			}
		}
		assert.Truef(t, found, "expected challenge not found: %+v", exp)
	}
}

func registerTOTP(ctx context.Context, t *testing.T, userID string) (secret string) {
	resp, err := UserClient.RegisterTOTP(ctx, &user.RegisterTOTPRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	secret = resp.GetSecret()
	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)

	_, err = UserClient.VerifyTOTPRegistration(ctx, &user.VerifyTOTPRegistrationRequest{
		UserId: userID,
		Code:   code,
	})
	require.NoError(t, err)
	return secret
}

func registerOTPSMS(ctx context.Context, t *testing.T, userID string) {
	_, err := UserClient.AddOTPSMS(ctx, &user.AddOTPSMSRequest{
		UserId: userID,
	})
	require.NoError(t, err)
}

func registerOTPEmail(ctx context.Context, t *testing.T, userID string) {
	_, err := UserClient.AddOTPEmail(ctx, &user.AddOTPEmailRequest{
		UserId: userID,
	})
	require.NoError(t, err)
}

func createRecoveryCodes(ctx context.Context, t *testing.T, userID string) (secret []string) {
	resp, err := UserClient.GenerateRecoveryCodes(ctx, &user.GenerateRecoveryCodesRequest{
		UserId: userID,
		Count:  5,
	})
	require.NoError(t, err)
	return resp.GetRecoveryCodes()
}
