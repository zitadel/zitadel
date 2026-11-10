package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestOTPEmailChallengeCommand_Validate(t *testing.T) {
	t.Parallel()
	otpEnabledAt := time.Now().Add(-30 * time.Minute)

	tests := []struct {
		name                  string
		sessionID             string
		instanceID            string
		secretGeneratorConfig *crypto.GeneratorConfig
		otpAlgorithm          crypto.EncryptionAlgorithm
		challengeTypeOTPEmail *domain.ChallengeTypeOTPEmail
		userRepo              func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo           func(ctrl *gomock.Controller) domain.SessionRepository
		wantErr               error
	}{
		{
			name:                  "no request otp email challenge",
			challengeTypeOTPEmail: nil,
			wantErr:               nil,
		},
		{
			name:                  "no session id",
			sessionID:             "",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			wantErr:               zerrors.ThrowPreconditionFailed(nil, "DOM-BQ5UgK", "Errors.Missing.SessionID"),
		},
		{
			name:                  "no instance id",
			sessionID:             "session-1",
			instanceID:            "",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			wantErr:               zerrors.ThrowPreconditionFailed(nil, "DOM-kDnkDn", "Errors.Missing.InstanceID"),
		},
		{
			name:                  "no default secret generator config",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			wantErr:               zerrors.ThrowInternal(nil, "DOM-nnB9MS", "missing default secret generator config"),
		},
		{
			name:                  "no otp algorithm",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-kuG75Q", "missing MFA encryption algorithm"),
		},
		{
			name:       "failed to render url template",
			sessionID:  "session-1",
			instanceID: "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{
						URLTemplate: "http://{{.Invalid",
					},
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			wantErr:               zerrors.ThrowInvalidArgument(nil, "DOM-wkDwQM", "Errors.Invalid.URLTemplate"),
		},
		{
			name:       "failed to get session",
			sessionID:  "session-1",
			instanceID: "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{
						URLTemplate: "https://example.com",
					},
				},
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "session-1"),
						),
					)).
					Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-JArUai", "failed fetching Session"),
		},
		{
			name:                  "session not found",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "session-1"),
						),
					)).
					Times(1).
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-JArUai", "Session not found"),
		},
		{
			name:                  "missing user id in session",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							repo.PrimaryKeyCondition("instance-1", "session-1"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "",
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-wG2XoJ", "Errors.Missing.Session.UserID"),
		},
		{
			name:                  "failed to get user",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-56MWkg", "failed fetching User"),
		},
		{
			name:                  "user not found",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-56MWkg", "User not found"),
		},
		{
			name:                  "human user not set",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-1",
						OrganizationID: "org-id",
						ID:             "123",
						Username:       "username",
						Human:          nil,
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2d", "Errors.NotFound.User.Human.Email"),
		},
		{
			name:                  "email not set",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-1",
						OrganizationID: "org-id",
						ID:             "123",
						Username:       "username",
						Human: &domain.HumanUser{
							FirstName: "first",
							LastName:  "last",
						},
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2d", "Errors.NotFound.User.Human.Email"),
		},
		{
			name:                  "email address not set",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-1",
						OrganizationID: "org-id",
						ID:             "123",
						Username:       "username",
						Human: &domain.HumanUser{
							FirstName: "first",
							LastName:  "last",
							Phone: &domain.HumanPhone{
								Number: "09080706050",
							},
							Email: domain.HumanEmail{
								Address: "",
							},
						},
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2d", "Errors.NotFound.User.Human.Email"),
		},
		{
			name:                  "email otp not enabled",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-1",
						OrganizationID: "org-id",
						ID:             "123",
						Username:       "username",
						Human: &domain.HumanUser{
							FirstName: "first",
							LastName:  "last",
							Phone: &domain.HumanPhone{
								Number: "09080706050",
							},
							Email: domain.HumanEmail{
								Address: "email@example.com",
								OTP:     domain.OTP{},
							},
						},
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4q", "Errors.User.MFA.OTP.NotReady"),
		},
		{
			name:                  "valid OTP Email challenge request",
			sessionID:             "session-1",
			instanceID:            "instance-1",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewUserRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.PrimaryKeyCondition("instance-1", "user-1"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-1",
						OrganizationID: "org-id",
						ID:             "123",
						Username:       "username",
						Human: &domain.HumanUser{
							FirstName: "first",
							LastName:  "last",
							Email: domain.HumanEmail{
								Address: "email@example.com",
								OTP: domain.OTP{
									EnabledAt: otpEnabledAt,
								},
							},
						},
					}, nil)
				return repo
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPEmailChallengeCommand(
				tt.challengeTypeOTPEmail,
				tt.sessionID,
				tt.instanceID,
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				nil,
			)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			err := cmd.Validate(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestOTPEmailChallengeCommand_Execute(t *testing.T) {
	t.Parallel()
	codeErr := errors.New("failed to create code")
	challengedAt := time.Now()
	defaultExpiry := 10 * time.Minute
	expiry := 30 * time.Minute
	otpEnabledAt := time.Now().Add(-30 * time.Minute)

	tests := []struct {
		name                        string
		challengeTypeOTPEmail       *domain.ChallengeTypeOTPEmail
		sessionRepo                 func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo                    func(ctrl *gomock.Controller) domain.UserRepository
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		secretGeneratorConfig       *crypto.GeneratorConfig
		otpAlgorithm                crypto.EncryptionAlgorithm
		newEmailCodeFn              func(g crypto.Generator) (*crypto.CryptoValue, string, error)
		wantErr                     error
		wantOTPEmailChallenge       *string
	}{
		{
			name:                  "no otp email challenge request",
			challengeTypeOTPEmail: nil,
			wantErr:               nil,
		},
		{
			name:                  "failed to retrieve secret generator config from settings DB",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{},
			otpAlgorithm:          crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: getUser(otpEnabledAt),
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewSecretGeneratorSettingsRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.UniqueCondition("instance-1", nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
							),
						),
					).Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-x7Yd3E", "failed fetching SecretGeneratorSettings"),
		},
		{
			name:                  "failed to generate code",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo:                    getUser(otpEnabledAt),
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return nil, "", codeErr
			},
			wantErr: codeErr,
		},
		{
			name: "failed to update session",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					ReturnCode: true,
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: expiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challengeOTPEmail, assert.AnError, 0)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-YfQIA3", "failed updating Session"),
		},
		{
			name: "failed to update session - no rows updated",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{
						URLTemplate: "http://example.com",
					},
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
				Length: uint(6),
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					URLTemplate:      "http://example.com",
					CodeReturned:     false,
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challengeOTPEmail, nil, 0)
				return repo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-YfQIA3", "Session not found"),
		},
		{
			name: "failed to update session - more than 1 row updated",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					ReturnCode: true,
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					LastChallengedAt: challengedAt,
				})
				repo.EXPECT().
					Update(
						gomock.Any(),
						gomock.Any(),
						repo.PrimaryKeyCondition("instance-1", "session-1"),
						challengeOTPEmail,
					).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-YfQIA3", "unexpected number of rows updated"),
		},
		{
			name: "update session succeeded - delivery type return code",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					ReturnCode: true,
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
			wantOTPEmailChallenge: gu.Ptr("code"), // OTPEmailChallenge is only set when the delivery type is return code
		},
		{
			name: "update session succeeded - delivery type send code",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{
						URLTemplate: "http://example.com",
					},
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					URLTemplate:      "http://example.com",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
		},
		{
			name:                  "update session succeeded - delivery type default",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					URLTemplate:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
		},
		{
			name:                  "update session succeeded - with default config when setting state is not active",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry:              defaultExpiry,
				Length:              uint(8),
				IncludeLowerLetters: false,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStatePreview, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           defaultExpiry,
					CodeReturned:     false,
					URLTemplate:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
		},
		{
			name:                  "update session succeeded - with default config when OTP email settings not found",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry:              defaultExpiry,
				Length:              uint(8),
				IncludeLowerLetters: false,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, 2),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           defaultExpiry,
					CodeReturned:     false,
					URLTemplate:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPEmailChallengeCommand(
				tt.challengeTypeOTPEmail,
				"session-1",
				"instance-1",
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.newEmailCodeFn,
			)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.secretGeneratorSettingsRepo != nil {
				domain.WithSecretGeneratorSettingsRepo(tt.secretGeneratorSettingsRepo(ctrl))(opts)
			}
			// to fetch/validate session and user before calling Execute
			err := cmd.Validate(ctx, opts)
			require.NoError(t, err)

			err = cmd.Execute(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantOTPEmailChallenge, cmd.GetOTPEmailChallenge())
		})
	}
}

func TestOTPEmailChallengeCommand_Events(t *testing.T) {
	t.Parallel()
	challengedAt := time.Now()
	expiry := 30 * time.Minute
	defaultExpiry := 10 * time.Minute
	otpEnabledAt := time.Now().Add(-30 * time.Minute)
	code := &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}
	ctx := authz.NewMockContext("instance-1", "", "")
	tests := []struct {
		name                        string
		challengeTypeOTPEmail       *domain.ChallengeTypeOTPEmail
		sessionRepo                 func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo                    func(ctrl *gomock.Controller) domain.UserRepository
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		secretGeneratorConfig       *crypto.GeneratorConfig
		otpAlgorithm                crypto.EncryptionAlgorithm
		newEmailCodeFn              func(g crypto.Generator) (*crypto.CryptoValue, string, error)
		wantErr                     error
		wantEvent                   eventstore.Command
	}{
		{
			name:                  "no request otp email challenge",
			challengeTypeOTPEmail: nil,
		},
		{
			name:                  "valid OTP email challenge request - no delivery type",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					URLTemplate:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
			wantEvent: session.NewOTPEmailChallengedEvent(
				ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				code,
				expiry,
				false,
				"",
			),
		},
		{
			name: "valid OTP email challenge request - delivery type send code",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					SendCode: &domain.SendCode{
						URLTemplate: "https://example.com",
					},
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					URLTemplate:      "https://example.com",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
			wantEvent: session.NewOTPEmailChallengedEvent(
				ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				code,
				expiry,
				false,
				"https://example.com",
			),
		},
		{
			name: "valid OTP email challenge request - delivery type return code",
			challengeTypeOTPEmail: &domain.ChallengeTypeOTPEmail{
				DeliveryType: domain.DeliveryType{
					ReturnCode: true,
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				challengeOTPEmail := repo.SetChallenge(&domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					URLTemplate:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPEmail)
				return repo
			},
			wantEvent: session.NewOTPEmailChallengedEvent(
				ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				code,
				expiry,
				true,
				"",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := domain.NewOTPEmailChallengeCommand(
				tt.challengeTypeOTPEmail,
				"session-1",
				"instance-1",
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.newEmailCodeFn,
			)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.secretGeneratorSettingsRepo != nil {
				domain.WithSecretGeneratorSettingsRepo(tt.secretGeneratorSettingsRepo(ctrl))(opts)
			}

			// to fetch/validate session and user before calling Execute
			err := cmd.Validate(ctx, opts)
			require.NoError(t, err)

			// to set the challengeOTPSMS in the command for event generation
			err = cmd.Execute(ctx, opts)
			require.NoError(t, err)
			events, err := cmd.Events(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantEvent != nil {
				require.Len(t, events, 1)
				assert.Equal(t, tt.wantEvent, events[0])
			}
		})
	}
}
