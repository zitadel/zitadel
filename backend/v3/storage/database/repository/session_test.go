package repository_test

import (
	"net"
	"net/http"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/crypto"
	zdomain "github.com/zitadel/zitadel/internal/domain"
)

func TestSession_Create(t *testing.T) {
	beforeCreate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	sessionRepo := repository.SessionRepository()

	instanceRepo := repository.InstanceRepository()
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "managementConsoleClient",
		ConsoleAppID:    "managementConsoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)
	tests := []struct {
		name    string
		session *domain.Session
		err     error
	}{
		{
			name: "create valid session",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				CreatedAt:  beforeCreate,
				UpdatedAt:  beforeCreate,
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
				UserAgent: &domain.SessionUserAgent{
					FingerprintID: gu.Ptr(gofakeit.Name()),
					Description:   gu.Ptr(gofakeit.Name()),
					IP:            net.IPv4(127, 0, 0, 1),
					Header:        http.Header{"User-Agent": []string{"user-agent"}},
				},
			},
			err: nil,
		},
		{
			name: "create valid session no dates",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
				UserAgent: &domain.SessionUserAgent{
					FingerprintID: gu.Ptr(gofakeit.Name()),
					Description:   gu.Ptr(gofakeit.Name()),
					IP:            net.IPv4(127, 0, 0, 1),
					Header:        http.Header{"User-Agent": []string{"user-agent"}},
				},
			},
			err: nil,
		},
		{
			name: "create session user agent fingerprint only",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
				UserAgent: &domain.SessionUserAgent{
					FingerprintID: gu.Ptr(gofakeit.Name()),
				},
			},
			err: nil,
		},
		{
			name: "create session empty fingerprint",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
				UserAgent: &domain.SessionUserAgent{
					FingerprintID: gu.Ptr(""),
				},
			},
			err: new(database.CheckError),
		},
		{
			name: "create session without user agent",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
			},
			err: nil,
		},
		{
			name: "create session without lifetime",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				CreatorID:  gofakeit.Name(),
			},
			err: nil,
		},
		{
			name: "create session without creatorID",
			session: &domain.Session{
				InstanceID: instanceId,
				ID:         gofakeit.Name(),
				Lifetime:   time.Hour * 24,
			},
			err: nil,
		},
		{
			name: "create session without id",
			session: &domain.Session{
				InstanceID: instanceId,
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
			},
			err: new(database.CheckError),
		},
		{
			name: "create session without instanceID",
			session: &domain.Session{
				ID:        gofakeit.Name(),
				Lifetime:  time.Hour * 24,
				CreatorID: gofakeit.Name(),
			},
			err: new(database.ForeignKeyError),
		},
		{
			name: "create session without existent instanceID",
			session: &domain.Session{
				InstanceID: gofakeit.Name(),
				ID:         gofakeit.Name(),
				Lifetime:   time.Hour * 24,
				CreatorID:  gofakeit.Name(),
			},
			err: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				err = savepoint.Rollback(t.Context())
				if err != nil {
					t.Logf("error during rollback: %v", err)
				}
			}()

			err = sessionRepo.Create(t.Context(), savepoint, tt.session)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check session values
			createdSession, err := sessionRepo.Get(t.Context(), savepoint,
				database.WithCondition(sessionRepo.PrimaryKeyCondition(tt.session.InstanceID, tt.session.ID)),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.session.ID, createdSession.ID)
			assert.Equal(t, tt.session.InstanceID, createdSession.InstanceID)
			assert.Equal(t, tt.session.Lifetime, createdSession.Lifetime)
			assert.Equal(t, createdSession.UpdatedAt.Add(tt.session.Lifetime), createdSession.Expiration)
			assert.Equal(t, tt.session.CreatorID, createdSession.CreatorID)
			assert.Equal(t, tt.session.UserAgent, createdSession.UserAgent)
			assert.WithinRange(t, createdSession.CreatedAt, beforeCreate.Add(-time.Second), afterCreate.Add(time.Second))
			assert.WithinRange(t, createdSession.UpdatedAt, beforeCreate.Add(-time.Second), afterCreate.Add(time.Second))
		})
	}
}

func TestSession_Update(t *testing.T) {
	beforeUpdate := time.Now()
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	sessionRepo := repository.SessionRepository()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "managementConsoleClient",
		ConsoleAppID:    "managementConsoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	orgID := createOrganization(t, tx, instanceId)
	userID := createHumanUser(t, tx, instanceId, orgID)

	testNow := time.Now()
	tests := []struct {
		name         string
		testFunc     func(t *testing.T) *domain.Session
		update       []database.Change
		rowsAffected int64
		err          error
	}{
		{
			name: "update no changes",
			testFunc: func(t *testing.T) *domain.Session {
				return &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
			},
			update:       []database.Change{},
			rowsAffected: 0,
			err:          database.ErrNoChanges,
		},
		{
			name: "update non existent session",
			testFunc: func(t *testing.T) *domain.Session {
				return &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
			},
			update: []database.Change{
				sessionRepo.SetLifetime(time.Hour * 48),
			},
			rowsAffected: 0,
		},
		{
			name: "update session updated at",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.UpdatedAt = testNow
				return session
			},
			update: []database.Change{
				sessionRepo.SetUpdatedAt(testNow),
			},
			rowsAffected: 1,
		},
		{
			name: "update session lifetime",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Lifetime = time.Hour * 48
				return session
			},
			update: []database.Change{
				sessionRepo.SetLifetime(time.Hour * 48),
			},
			rowsAffected: 1,
		},
		{
			name: "update session token",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.TokenID = "new-token"
				return session
			},
			update: []database.Change{
				sessionRepo.SetToken("new-token"),
			},
			rowsAffected: 1,
		},
		{
			name: "update session passkey challenge",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Challenges = []domain.SessionChallenge{
					&domain.SessionChallengePasskey{
						LastChallengedAt:     testNow,
						Challenge:            "challenge",
						AllowedCredentialIDs: [][]byte{[]byte("allowed-id-1"), []byte("allowed-id-2")},
						UserVerification:     zdomain.UserVerificationRequirementRequired,
						RPID:                 "rp-id",
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetChallenge(&domain.SessionChallengePasskey{
					LastChallengedAt:     testNow,
					Challenge:            "challenge",
					AllowedCredentialIDs: [][]byte{[]byte("allowed-id-1"), []byte("allowed-id-2")},
					UserVerification:     zdomain.UserVerificationRequirementRequired,
					RPID:                 "rp-id",
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session otp sms challenge",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Challenges = []domain.SessionChallenge{
					&domain.SessionChallengeOTPSMS{
						LastChallengedAt: testNow,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "key-id",
							Crypted:    []byte("code"),
						},
						Expiry:            time.Minute * 10,
						CodeReturned:      false,
						GeneratorID:       "",
						TriggeredAtOrigin: "https://example.com",
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetChallenge(&domain.SessionChallengeOTPSMS{
					LastChallengedAt: testNow,
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "key-id",
						Crypted:    []byte("code"),
					},
					Expiry:            time.Minute * 10,
					CodeReturned:      false,
					GeneratorID:       "",
					TriggeredAtOrigin: "https://example.com",
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session otp email challenge",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Challenges = []domain.SessionChallenge{
					&domain.SessionChallengeOTPEmail{
						LastChallengedAt: testNow,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "key-id",
							Crypted:    []byte("code"),
						},
						Expiry:            time.Minute * 10,
						CodeReturned:      false,
						TriggeredAtOrigin: "https://example.com",
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetChallenge(&domain.SessionChallengeOTPEmail{
					LastChallengedAt: testNow,
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "key-id",
						Crypted:    []byte("code"),
					},
					Expiry:            time.Minute * 10,
					CodeReturned:      false,
					TriggeredAtOrigin: "https://example.com",
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session user factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorUser{
						UserID:         userID,
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorUser{
					UserID:         userID,
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session password factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorPassword{
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorPassword{
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session identity provider intent factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorIdentityProviderIntent{
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorIdentityProviderIntent{
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session passkey factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorPasskey{
						LastVerifiedAt: testNow,
						UserVerified:   true,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorPasskey{
					LastVerifiedAt: testNow,
					UserVerified:   true,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session totp factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorTOTP{
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorTOTP{
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session otp sms factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorOTPSMS{
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorOTPSMS{
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session otp email factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorOTPEmail{
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorOTPEmail{
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session recovery code factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorRecoveryCode{
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
					LastVerifiedAt: testNow,
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session metadata",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.Metadata = []*domain.SessionMetadata{
					{
						Metadata: domain.Metadata{
							InstanceID: instanceId,
							Key:        "key1",
							Value:      []byte("value1"),
						},
						SessionID: session.ID,
					},
					{
						Metadata: domain.Metadata{
							InstanceID: instanceId,
							Key:        "key2",
							Value:      []byte("value2"),
						},
						SessionID: session.ID,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetMetadata([]*domain.SessionMetadata{
					{
						Metadata: domain.Metadata{
							Key:   "key1",
							Value: []byte("value1"),
						},
					},
					{
						Metadata: domain.Metadata{
							Key:   "key2",
							Value: []byte("value2"),
						},
					},
				}),
			},
			rowsAffected: 1,
		},
		{
			name: "update session multiple fields",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				// update with updated values
				session.UpdatedAt = testNow
				session.Lifetime = time.Hour * 48
				session.TokenID = "new-token"
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorUser{
						UserID:         userID,
						LastVerifiedAt: testNow,
					},
					&domain.SessionFactorPassword{
						LastVerifiedAt: testNow,
					},
				}
				session.Challenges = []domain.SessionChallenge{
					&domain.SessionChallengeOTPSMS{
						LastChallengedAt: testNow,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "key-id",
							Crypted:    []byte("code"),
						},
						Expiry:            time.Minute * 10,
						CodeReturned:      false,
						GeneratorID:       "",
						TriggeredAtOrigin: "https://example.com",
					},
				}
				session.Metadata = []*domain.SessionMetadata{
					{
						Metadata: domain.Metadata{
							InstanceID: instanceId,
							Key:        "key1",
							Value:      []byte("value1"),
						},
						SessionID: session.ID,
					},
					{
						Metadata: domain.Metadata{
							InstanceID: instanceId,
							Key:        "key2",
							Value:      []byte("value2"),
						},
						SessionID: session.ID,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.SetUpdatedAt(testNow),
				sessionRepo.SetLifetime(time.Hour * 48),
				sessionRepo.SetFactor(&domain.SessionFactorUser{
					UserID:         userID,
					LastVerifiedAt: testNow,
				}),
				sessionRepo.SetFactor(&domain.SessionFactorPassword{
					LastVerifiedAt: testNow,
				}),
				sessionRepo.SetChallenge(&domain.SessionChallengeOTPSMS{
					LastChallengedAt: testNow,
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "key-id",
						Crypted:    []byte("code"),
					},
					Expiry:            time.Minute * 10,
					CodeReturned:      false,
					GeneratorID:       "",
					TriggeredAtOrigin: "https://example.com",
				}),
				sessionRepo.SetMetadata([]*domain.SessionMetadata{
					{
						Metadata: domain.Metadata{
							Key:   "key1",
							Value: []byte("value1"),
						},
					},
					{
						Metadata: domain.Metadata{
							Key:   "key2",
							Value: []byte("value2"),
						},
					},
				}),
				sessionRepo.SetToken("new-token"),
			},
			rowsAffected: 1,
		},
		{
			name: "update session clear factor",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)

				updatedRows, err := sessionRepo.Update(t.Context(), tx,
					sessionRepo.PrimaryKeyCondition(session.InstanceID, session.ID),
					sessionRepo.SetFactor(&domain.SessionFactorUser{
						UserID:         userID,
						LastVerifiedAt: testNow,
					}),
					sessionRepo.SetFactor(&domain.SessionFactorPassword{
						LastVerifiedAt: testNow,
					}),
				)
				require.NoError(t, err)
				require.Equal(t, int64(1), updatedRows)

				// update with updated values
				session.Factors = []domain.SessionFactor{
					&domain.SessionFactorUser{
						UserID:         userID,
						LastVerifiedAt: testNow,
					},
				}
				return session
			},
			update: []database.Change{
				sessionRepo.ClearFactor(domain.SessionFactorTypePassword),
			},
			rowsAffected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdSession := tt.testFunc(t)

			// update org
			rowsAffected, err := sessionRepo.Update(t.Context(), tx,
				sessionRepo.PrimaryKeyCondition(createdSession.InstanceID, createdSession.ID),
				tt.update...,
			)
			afterUpdate := time.Now()
			require.ErrorIs(t, err, tt.err)

			assert.Equal(t, tt.rowsAffected, rowsAffected)

			if rowsAffected == 0 {
				return
			}

			// check session values
			session, err := sessionRepo.Get(t.Context(), tx,
				database.WithCondition(
					sessionRepo.PrimaryKeyCondition(createdSession.InstanceID, createdSession.ID),
				),
			)
			require.NoError(t, err)

			assert.Equal(t, createdSession.ID, session.ID)
			assert.Equal(t, createdSession.Lifetime, session.Lifetime)
			assert.Equal(t, createdSession.TokenID, session.TokenID)
			assert.WithinRange(t, session.Expiration, beforeUpdate.Add(createdSession.Lifetime), afterUpdate.Add(createdSession.Lifetime))
			assert.Equal(t, createdSession.CreatorID, session.CreatorID)
			assert.Equal(t, createdSession.UserAgent, session.UserAgent)
			assert.WithinRange(t, session.UpdatedAt, beforeUpdate, afterUpdate)
			assert.EqualExportedValues(t, createdSession.Challenges, session.Challenges)
			slices.SortFunc(session.Factors, func(a, b domain.SessionFactor) int { // sort to make comparison possible
				return int(a.SessionFactorType()) - int(b.SessionFactorType())
			})
			assert.EqualExportedValues(t, createdSession.Factors, session.Factors)
			assert.EqualExportedValues(t, createdSession.Metadata, session.Metadata)
		})
	}
}

func TestSession_Delete(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	sessionRepo := repository.SessionRepository()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "managementConsoleClient",
		ConsoleAppID:    "managementConsoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	type test struct {
		name            string
		testFunc        func(t *testing.T) *domain.Session
		noOfDeletedRows int64
	}
	tests := []test{
		{
			name: "delete existent session",
			testFunc: func(t *testing.T) *domain.Session {
				session := &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
				err := sessionRepo.Create(t.Context(), tx, session)
				require.NoError(t, err)
				return session
			},
			noOfDeletedRows: 1,
		},
		{
			name: "delete non existent session",
			testFunc: func(t *testing.T) *domain.Session {
				return &domain.Session{
					InstanceID: instanceId,
					ID:         gofakeit.Name(),
					Lifetime:   time.Hour * 24,
					CreatorID:  gofakeit.Name(),
				}
			},
			noOfDeletedRows: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := tt.testFunc(t)

			// delete session
			deletedRows, err := sessionRepo.Delete(t.Context(), tx,
				sessionRepo.PrimaryKeyCondition(session.InstanceID, session.ID),
			)
			require.NoError(t, err)
			assert.Equal(t, tt.noOfDeletedRows, deletedRows)

			// verify session deletion
			deletedSession, err := sessionRepo.Get(t.Context(), tx,
				database.WithCondition(
					sessionRepo.PrimaryKeyCondition(session.InstanceID, session.ID),
				),
			)
			require.Error(t, err, new(database.NoRowFoundError))
			assert.Nil(t, deletedSession)
		})
	}
}

func TestSession_Get(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	sessionRepo := repository.SessionRepository()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "managementConsoleClient",
		ConsoleAppID:    "managementConsoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	session := &domain.Session{
		InstanceID: instanceId,
		ID:         gofakeit.Name(),
		Lifetime:   time.Hour * 24,
		CreatorID:  gofakeit.Name(),
	}
	err = sessionRepo.Create(t.Context(), tx, session)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		err       error
	}{
		{
			name:      "get existent session by primary key",
			condition: sessionRepo.PrimaryKeyCondition(instanceId, session.ID),
			err:       nil,
		},
		{
			name:      "get without instanceID condition",
			condition: sessionRepo.IDCondition(session.ID),
			err:       new(database.MissingConditionError),
		},
		{
			name:      "get non existent session",
			condition: sessionRepo.PrimaryKeyCondition(instanceId, gofakeit.Name()),
			err:       new(database.NoRowFoundError),
		},
		// other conditions are tested in the List tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// get session
			got, err := sessionRepo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}

			// check session values
			assert.Equal(t, session.ID, got.ID)
			assert.Equal(t, session.InstanceID, got.InstanceID)
			assert.Equal(t, session.Lifetime, got.Lifetime)
			assert.Equal(t, session.CreatorID, got.CreatorID)
		})
	}
}

func TestSession_List(t *testing.T) {
	tx, err := pool.Begin(t.Context(), nil)
	require.NoError(t, err)
	defer func() {
		err := tx.Rollback(t.Context())
		if err != nil {
			t.Logf("error during rollback: %v", err)
		}
	}()

	instanceRepo := repository.InstanceRepository()
	sessionRepo := repository.SessionRepository()

	// create instance
	instanceId := gofakeit.Name()
	instance := domain.Instance{
		ID:              instanceId,
		Name:            gofakeit.Name(),
		DefaultOrgID:    "defaultOrgId",
		IAMProjectID:    "iamProject",
		ConsoleClientID: "managementConsoleClient",
		ConsoleAppID:    "managementConsoleApp",
		DefaultLanguage: "defaultLanguage",
	}
	err = instanceRepo.Create(t.Context(), tx, &instance)
	require.NoError(t, err)

	orgID := createOrganization(t, tx, instanceId)
	userIDs := []string{createHumanUser(t, tx, instanceId, orgID), createHumanUser(t, tx, instanceId, orgID)}

	// create sessions
	sessions := make([]*domain.Session, 5)
	for i := range 5 {
		session := &domain.Session{
			InstanceID: instanceId,
			ID:         strconv.Itoa(i),
			Lifetime:   time.Hour * 24,
			CreatorID:  gofakeit.Name(),
			CreatedAt:  time.Now(),
			UserAgent: &domain.SessionUserAgent{
				FingerprintID: gu.Ptr(gofakeit.UUID()),
				Description:   gu.Ptr(gofakeit.Name()),
				IP:            net.ParseIP(gofakeit.IPv4Address()),
				Header:        http.Header{"user-agent": []string{gofakeit.Name()}},
			},
		}
		err = sessionRepo.Create(t.Context(), tx, session)
		require.NoError(t, err)
		session.Expiration = session.CreatedAt.Add(session.Lifetime)

		changes := make([]database.Change, 0, 4)
		updated := time.Now()
		changes = append(changes, sessionRepo.SetUpdatedAt(updated))
		session.UpdatedAt = updated
		metadata := []*domain.SessionMetadata{
			{
				Metadata: domain.Metadata{
					InstanceID: instanceId,
					Key:        "key" + strconv.Itoa(i),
					Value:      []byte{uint8(i)},
				},
				SessionID: session.ID,
			},
		}
		session.Metadata = metadata
		changes = append(changes, sessionRepo.SetMetadata(metadata))
		userFactor := &domain.SessionFactorUser{
			UserID:         userIDs[i%2],
			LastVerifiedAt: time.Now(),
		}
		session.Factors.AppendTo(userFactor)
		changes = append(changes, sessionRepo.SetFactor(userFactor))
		if i < 3 {
			passwordFactor := &domain.SessionFactorPassword{
				LastVerifiedAt: time.Now(),
			}
			session.Factors.AppendTo(passwordFactor)
			changes = append(changes, sessionRepo.SetFactor(passwordFactor))
		}
		_, err = sessionRepo.Update(t.Context(), tx, sessionRepo.PrimaryKeyCondition(instanceId, session.ID),
			changes...,
		)
		require.NoError(t, err)
		session.UserID = userFactor.UserID
		sessions[i] = session
	}

	tests := []struct {
		name     string
		options  []database.QueryOption
		expected []*domain.Session
		err      error
	}{
		{
			name: "list sessions by instanceID",
			options: []database.QueryOption{
				database.WithCondition(
					sessionRepo.InstanceIDCondition(instanceId),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: sessions,
			err:      nil,
		},
		{
			name: "list single session by id",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.IDCondition(sessions[0].ID),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
			},
			expected: []*domain.Session{sessions[0]},
			err:      nil,
		},
		{
			name: "list sessions with non matching condition",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.IDCondition("non-existent-id"),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
			},
			expected: []*domain.Session{},
			err:      nil,
		},
		{
			name:     "list sessions without conditions",
			options:  nil,
			expected: nil,
			err:      new(database.MissingConditionError),
		},
		{
			name: "list sessions by user agent fingerprintID",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.UserAgentIDCondition(gu.Value(sessions[1].UserAgent.FingerprintID)),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
			},
			expected: []*domain.Session{sessions[1]},
			err:      nil,
		},
		{
			name: "list sessions by userID",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.UserIDCondition(sessions[0].Factors.GetUserFactor().UserID),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
				database.WithLimit(2),
			},
			expected: []*domain.Session{sessions[0], sessions[2]},
			err:      nil,
		},
		{
			name: "list sessions by creatorID",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.CreatorIDCondition(sessions[3].CreatorID),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
			},
			expected: []*domain.Session{sessions[3]},
		},
		{
			name: "list sessions by expiration",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.ExpirationCondition(database.NumberOperationGreaterThanOrEqual, sessions[3].Expiration),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: []*domain.Session{sessions[3], sessions[4]},
		},
		{
			name: "list sessions by created at",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.CreatedAtCondition(database.NumberOperationLessThan, sessions[3].CreatedAt),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: []*domain.Session{sessions[0], sessions[1], sessions[2]},
		},
		{
			name: "list sessions by factor type",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.ExistsFactor(sessionRepo.FactorConditions().FactorTypeCondition(domain.SessionFactorTypePassword)),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: []*domain.Session{sessions[0], sessions[1], sessions[2]},
		},
		{
			name: "list sessions by factor verification time",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.ExistsFactor(sessionRepo.FactorConditions().LastVerifiedBeforeCondition(sessions[2].Factors.GetUserFactor().LastVerifiedAt)),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: []*domain.Session{sessions[0], sessions[1]},
		},
		{
			name: "list sessions by metadata key",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.ExistsMetadata(sessionRepo.MetadataConditions().KeyCondition(database.TextOperationEqual, "key3")),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: []*domain.Session{sessions[3]},
		},
		{
			name: "list sessions by metadata value",
			options: []database.QueryOption{
				database.WithCondition(
					database.And(
						sessionRepo.ExistsMetadata(sessionRepo.MetadataConditions().ValueCondition(database.BytesOperationEqual, []byte{3})),
						sessionRepo.InstanceIDCondition(instanceId),
					),
				),
				database.WithOrderByAscending(sessionRepo.IDColumn()),
			},
			expected: []*domain.Session{sessions[3]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// list sessions
			got, err := sessionRepo.List(t.Context(), tx, tt.options...)
			assert.ErrorIs(t, err, tt.err)
			if err != nil {
				return
			}
			assert.Len(t, got, len(tt.expected))
			for _, session := range got {
				slices.SortFunc(session.Factors, func(a, b domain.SessionFactor) int { // sort to make comparison possible
					return int(a.SessionFactorType()) - int(b.SessionFactorType())
				})
			}
			assert.EqualExportedValues(t, tt.expected, got)
		})
	}
}
