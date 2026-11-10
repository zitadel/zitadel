package domain_test

import (
	"context"
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

func TestOTPSMSChallengeCommand_Validate(t *testing.T) {
	t.Parallel()
	otpEnabledAt := time.Now().Add(-30 * time.Minute)

	tests := []struct {
		name                  string
		sessionID             string
		instanceID            string
		secretGeneratorConfig *crypto.GeneratorConfig
		otpAlgorithm          crypto.EncryptionAlgorithm
		challengeTypeOTPSMS   *domain.ChallengeTypeOTPSMS
		userRepo              func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo           func(ctrl *gomock.Controller) domain.SessionRepository
		smsProvider           func(ctx context.Context, instanceID string) (string, error)
		wantErr               error
	}{
		{
			name:                "no request otpsms challenge",
			challengeTypeOTPSMS: nil,
			wantErr:             nil,
		},
		{
			name:                "no session id",
			sessionID:           "",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-3XpM6A", "Errors.Missing.SessionID"),
		},
		{
			name:                "no instance id",
			sessionID:           "session-1",
			instanceID:          "",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-jNNJ9f", "Errors.Missing.InstanceID"),
		},
		{
			name:                "no sms provider",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			smsProvider:         nil,
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-1aGWeE", "missing sms provider"),
		},
		{
			name:                "no default secret generator config",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-IDcOzP", "missing default secret generator config"),
		},
		{
			name:                "no otp algorithm",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-VeBPmV", "missing MFA encryption algorithm"),
		},
		{
			name:                "failed to get session",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-2aGWWE", "failed fetching Session"),
		},
		{
			name:                "session not found",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-2aGWWE", "Session not found"),
		},
		{
			name:                "missing user id in session",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-Vi16Fs", "Errors.Missing.Session.UserID"),
		},
		{
			name:                "failed to get user",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-3aGHDs", "failed fetching User"),
		},
		{
			name:                "user not found",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-3aGHDs", "User not found"),
		},
		{
			name:                "human user not set",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2w", "Errors.NotFound.User.Human.Phone"),
		},
		{
			name:                "phone not set",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2w", "Errors.NotFound.User.Human.Phone"),
		},
		{
			name:                "phone otp not enabled",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
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
						},
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4m", "Errors.User.MFA.OTP.NotReady"),
		},
		{
			name:                "valid OTP SMS challenge request",
			sessionID:           "session-1",
			instanceID:          "instance-1",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: 10 * time.Minute,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: getUser(otpEnabledAt),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPSMSChallengeCommand(
				tt.challengeTypeOTPSMS,
				tt.sessionID,
				tt.instanceID,
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.smsProvider,
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

func TestOTPSMSChallengeCommand_Execute(t *testing.T) {
	t.Parallel()
	challengedAt := time.Now()
	otpEnabledAt := time.Now().Add(-30 * time.Minute)
	smsProviderErr := errors.New("failed to get active sms provider")
	codeErr := errors.New("failed to create code")
	defaultExpiry := 10 * time.Minute
	expiry := 30 * time.Minute

	tests := []struct {
		name                        string
		challengeTypeOTPSMS         *domain.ChallengeTypeOTPSMS
		sessionRepo                 func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo                    func(ctrl *gomock.Controller) domain.UserRepository
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		secretGeneratorConfig       *crypto.GeneratorConfig
		otpAlgorithm                crypto.EncryptionAlgorithm
		smsProviderFn               func(ctx context.Context, instanceID string) (string, error)
		newPhoneCodeFn              func(g crypto.Generator) (*crypto.CryptoValue, string, error)
		wantErr                     error
		wantOTPSMSChallenge         *string
	}{
		{
			name:                "no otp sms challenge request",
			challengeTypeOTPSMS: nil,
			wantErr:             nil,
		},
		{
			name: "failed to get active sms provider",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: false,
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo: getUser(otpEnabledAt),
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", smsProviderErr
			},
			wantErr: smsProviderErr,
		},
		{
			name: "failed to retrieve secret generator config from settings DB",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: false,
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil // no external sms provider
			},
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
			name:                "failed to generate code",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil // no external sms provider
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)
				return repo
			},
			userRepo:                    getUser(otpEnabledAt),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return nil, "", codeErr
			},
			wantErr: codeErr,
		},
		{
			name:                "update successful - external sms provider",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "external-sms-provider-id", nil // with an external sms provider
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			userRepo:     getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code:             nil,
					Expiry:           0,
					CodeReturned:     false,
					GeneratorID:      "external-sms-provider-id",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
		},
		{
			name: "update successful - external sms provider with return code true",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "external-sms-provider-id", nil // with an external sms provider
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code:             nil,
					Expiry:           0,
					CodeReturned:     true,
					GeneratorID:      "external-sms-provider-id",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantOTPSMSChallenge: gu.Ptr(""), // no code is generated for external sms provider
		},
		{
			name: "update successful - internal sms provider",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantOTPSMSChallenge: gu.Ptr("code"),
		},
		{
			name: "update successful - internal sms provider - with return code false",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: false,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantOTPSMSChallenge: nil,
		},
		{
			name: "update successful - internal sms provider - with default config",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStatePreview, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           defaultExpiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantOTPSMSChallenge: gu.Ptr("code"),
		},
		{
			name: "update successful - internal sms provider - missing SMS OTP config - with default config",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, 2),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           defaultExpiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantOTPSMSChallenge: gu.Ptr("code"),
		},
		{
			name: "failed to update session",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challengeOTPSMSChange, assert.AnError, 0)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-AigB0Z", "failed updating Session"),
		},
		{
			name: "failed to update session - no rows updated",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challengeOTPSMSChange, nil, 0)
				return repo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-AigB0Z", "Session not found"),
		},
		{
			name: "failed to update session - more than 1 row updated",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionFailedExpectation(repo, challengeOTPSMSChange, nil, 2)
				return repo
			},
			wantErr: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-AigB0Z", "unexpected number of rows updated"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPSMSChallengeCommand(tt.challengeTypeOTPSMS,
				"session-1",
				"instance-1",
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.smsProviderFn,
				tt.newPhoneCodeFn,
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
			assert.NoError(t, err)

			err = cmd.Execute(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantOTPSMSChallenge, cmd.GetOTPSMSChallenge())
		})
	}
}

func TestOTPSMSChallengeCommand_Events(t *testing.T) {
	t.Parallel()
	challengedAt := time.Now()
	otpEnabledAt := time.Now().Add(-30 * time.Minute)
	defaultExpiry := 10 * time.Minute
	expiry := 30 * time.Minute
	code := &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    []byte("code"),
	}

	ctx := authz.NewMockContext("instance-1", "", "")
	tests := []struct {
		name                        string
		challengeTypeOTPSMS         *domain.ChallengeTypeOTPSMS
		sessionRepo                 func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo                    func(ctrl *gomock.Controller) domain.UserRepository
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		secretGeneratorConfig       *crypto.GeneratorConfig
		otpAlgorithm                crypto.EncryptionAlgorithm
		smsProviderFn               func(ctx context.Context, instanceID string) (string, error)
		newPhoneCodeFn              func(g crypto.Generator) (*crypto.CryptoValue, string, error)
		wantErr                     error
		wantEvent                   eventstore.Command
	}{
		{
			name:                "no otp sms challenge request",
			challengeTypeOTPSMS: nil,
		},
		{
			name: "valid otp sms challenge request - internal provider - with return code set to false",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: false,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantEvent: session.NewOTPSMSChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				code,
				expiry,
				false,
				"",
			),
		},
		{
			name: "valid otp sms challenge request - internal provider - with return code set to true",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					GeneratorID:      "",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantEvent: session.NewOTPSMSChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				code,
				expiry,
				true,
				"",
			),
		},
		{
			name:                "update successful - external sms provider - return code false",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "external-sms-provider-id", nil // with an external sms provider
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code:             nil,
					Expiry:           0,
					CodeReturned:     false,
					GeneratorID:      "external-sms-provider-id",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantEvent: session.NewOTPSMSChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				nil,
				0,
				false,
				"external-sms-provider-id",
			),
		},
		{
			name: "update successful - external sms provider - return code true",
			challengeTypeOTPSMS: &domain.ChallengeTypeOTPSMS{
				ReturnCode: true,
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "external-sms-provider-id", nil // with an external sms provider
			},
			userRepo: getUser(otpEnabledAt),
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewSessionRepo(ctrl)
				getSessionSucceededExpectation(repo)

				challengeOTPSMSChange := repo.SetChallenge(&domain.SessionChallengeOTPSMS{
					Code:             nil,
					Expiry:           0,
					CodeReturned:     true,
					GeneratorID:      "external-sms-provider-id",
					LastChallengedAt: challengedAt,
				})
				updateSessionSucceededExpectation(repo, challengeOTPSMSChange)
				return repo
			},
			wantEvent: session.NewOTPSMSChallengedEvent(ctx,
				&session.NewAggregate("session-1", "instance-1").Aggregate,
				nil,
				0,
				true,
				"external-sms-provider-id",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPSMSChallengeCommand(tt.challengeTypeOTPSMS,
				"session-1",
				"instance-1",
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.smsProviderFn,
				tt.newPhoneCodeFn,
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
			assert.NoError(t, err)

			// to set the challengeOTPSMS in the command for event generation
			err = cmd.Execute(ctx, opts)
			assert.NoError(t, err)

			events, err := cmd.Events(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantEvent != nil {
				require.Len(t, events, 1)
				assert.Equal(t, tt.wantEvent, events[0])
			}
		})
	}
}

func secretGeneratorSettingsRepo(state domain.SettingState, otpType domain.OTPType) func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
	var attrs domain.SecretGeneratorSettingsAttributes
	switch otpType {
	case domain.OTPTypeEmail:
		attrs = domain.SecretGeneratorSettingsAttributes{
			OTPEmail: &domain.OTPEmailAttributes{
				SecretGeneratorAttrsWithExpiry: domain.SecretGeneratorAttrsWithExpiry{
					Expiry: gu.Ptr(30 * time.Minute),
					SecretGeneratorAttrs: domain.SecretGeneratorAttrs{
						Length:              gu.Ptr(uint(8)),
						IncludeLowerLetters: gu.Ptr(true),
						IncludeUpperLetters: gu.Ptr(false),
						IncludeDigits:       gu.Ptr(true),
						IncludeSymbols:      gu.Ptr(false),
					},
				},
			},
		}
	case domain.OTPTypeSMS:
		attrs = domain.SecretGeneratorSettingsAttributes{
			OTPSMS: &domain.OTPSMSAttributes{
				SecretGeneratorAttrsWithExpiry: domain.SecretGeneratorAttrsWithExpiry{
					Expiry: gu.Ptr(30 * time.Minute),
					SecretGeneratorAttrs: domain.SecretGeneratorAttrs{
						Length:              gu.Ptr(uint(6)),
						IncludeLowerLetters: gu.Ptr(true),
						IncludeUpperLetters: gu.Ptr(false),
						IncludeDigits:       gu.Ptr(true),
						IncludeSymbols:      gu.Ptr(false),
					},
				},
			},
		}
	}
	return func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
		repo := domainmock.NewSecretGeneratorSettingsRepo(ctrl)
		repo.EXPECT().
			Get(gomock.Any(), gomock.Any(),
				dbmock.QueryOptions(
					database.WithCondition(
						repo.UniqueCondition("instance-1", nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
					),
				),
			).Times(1).
			Return(&domain.SecretGeneratorSettings{
				Settings: domain.Settings{
					State: state,
				},
				SecretGeneratorSettingsAttributes: attrs,
			}, nil)
		return repo
	}
}

func getUser(otpEnabledAt time.Time) func(ctrl *gomock.Controller) domain.UserRepository {
	return func(ctrl *gomock.Controller) domain.UserRepository {
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
						OTP: domain.OTP{
							EnabledAt: otpEnabledAt,
						},
					},
					Email: domain.HumanEmail{
						Address: "testuser@example.com",
						OTP: domain.OTP{
							EnabledAt: otpEnabledAt,
						},
					},
				},
			}, nil)
		return repo
	}
}
