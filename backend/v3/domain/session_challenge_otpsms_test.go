package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestOTPSMSChallengeCommand_Validate(t *testing.T) {
	t.Parallel()
	otpEnabledAt := time.Now().Add(-30 * time.Minute)

	tests := []struct {
		name                string
		sessionID           string
		instanceID          string
		requestChallengeOTP *session_grpc.RequestChallenges_OTPSMS
		userRepo            func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo         func(ctrl *gomock.Controller) domain.SessionRepository
		smsProvider         func(ctx context.Context, instanceID string) (string, error)
		wantErr             error
		wantUser            *domain.User
		wantSession         *domain.Session
	}{
		{
			name:                "no request otpsms challenge",
			requestChallengeOTP: nil,
			wantErr:             nil,
		},
		{
			name:                "no session id",
			sessionID:           "",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-3XpM6A", "session id missing"),
		},
		{
			name:                "no instance id",
			sessionID:           "session-id",
			instanceID:          "",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-jNNJ9f", "instance id missing"),
		},
		{
			name:                "no sms provider",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider:         nil,
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-1aGWeE", "sms provider not configured"),
		},
		{
			name:                "failed to get session",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-2aGWWE", "failed fetching session"),
		},
		{
			name:                "session not found",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-2aGWWE", "session not found"),
		},
		{
			name:                "missing user id in session",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "",
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-Vi16Fs", "missing user id in session"),
		},
		{
			name:                "failed to get user",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					AnyTimes().
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-3aGHDs", "failed fetching user"),
		},
		{
			name:                "user not found",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					AnyTimes().
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-3aGHDs", "user not found"),
		},
		{
			name:                "human user not set",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					AnyTimes().
					Return(&domain.User{
						InstanceID:     "instance-id",
						OrganizationID: "org-id",
						ID:             "123",
						Username:       "username",
						Human:          nil,
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2w", "user phone not configured"),
		},
		{
			name:                "phone not set",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					AnyTimes().
					Return(&domain.User{
						InstanceID:     "instance-id",
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2w", "user phone not configured"),
		},
		{
			name:                "phone otp not enabled",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					AnyTimes().
					Return(&domain.User{
						InstanceID:     "instance-id",
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4m", "phone OTP not enabled"),
		},
		{
			name:                "valid OTP SMS challenge request",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPSMS{},
			smsProvider: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					AnyTimes().
					Return(&domain.Session{
						UserID:     "user-id",
						InstanceID: "instance-id",
						ID:         "session-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				humanRepo := domainmock.NewMockHumanUserRepository(ctrl)
				repo.EXPECT().
					LoadPasskeys().
					AnyTimes().
					Return(repo)
				repo.EXPECT().
					Human().
					AnyTimes().
					Return(humanRepo)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					AnyTimes().
					Return(&domain.User{
						InstanceID:     "instance-id",
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
						},
					}, nil)
				return repo
			},
			wantUser: &domain.User{
				InstanceID:     "instance-id",
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
				},
			},
			wantSession: &domain.Session{
				UserID:     "user-id",
				InstanceID: "instance-id",
				ID:         "session-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-1", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPSMSChallengeCommand(
				tt.requestChallengeOTP,
				tt.sessionID,
				tt.instanceID,
				nil,
				nil,
				tt.smsProvider,
				nil,
			)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			err := cmd.Validate(ctx, opts)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantUser, cmd.User)
				assert.Equal(t, tt.wantSession, cmd.Session)
			}
		})
	}
}

func TestOTPSMSChallengeCommand_Events(t *testing.T) {
	t.Parallel()
	challengedAt := time.Now()
	expiry := 30 * time.Minute
	code := &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "alg",
		KeyID:      "encKey",
		Crypted:    []byte("crypted"),
	}

	ctx := authz.NewMockContext("instance-id", "", "")
	tests := []struct {
		name                   string
		requestChallengeOTPSMS *session_grpc.RequestChallenges_OTPSMS
		challengeOTPSMS        *domain.SessionChallengeOTPSMS
		wantErr                error
		wantEvent              eventstore.Command
	}{
		{
			name:                   "no otp sms challenge request",
			requestChallengeOTPSMS: nil,
		},
		{
			name: "valid otp sms challenge request with return code set to false",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: false,
			},
			challengeOTPSMS: &domain.SessionChallengeOTPSMS{
				LastChallengedAt: challengedAt,
				Code:             code,
				Expiry:           expiry,
				GeneratorID:      "generator-id",
			},
			wantEvent: session.NewOTPSMSChallengedEvent(ctx,
				&session.NewAggregate("session-id", "instance-id").Aggregate,
				code,
				expiry,
				false,
				"generator-id",
			),
		},
		{
			name: "valid otp sms challenge request with return code set to true",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			challengeOTPSMS: &domain.SessionChallengeOTPSMS{
				LastChallengedAt: challengedAt,
				Code:             code,
				Expiry:           expiry,
				GeneratorID:      "generator-id",
				CodeReturned:     true,
			},
			wantEvent: session.NewOTPSMSChallengedEvent(ctx,
				&session.NewAggregate("session-id", "instance-id").Aggregate,
				code,
				expiry,
				true,
				"generator-id",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cmd := domain.NewOTPSMSChallengeCommand(
				tt.requestChallengeOTPSMS,
				"session-id",
				"instance-id",
				nil,
				nil,
				nil,
				nil,
			)
			cmd.ChallengeOTPSMS = tt.challengeOTPSMS
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			events, err := cmd.Events(ctx, opts)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantEvent != nil {
				require.Len(t, events, 1)
				assert.Equal(t, tt.wantEvent, events[0])
			}
		})
	}
}

func TestOTPSMSChallengeCommand_Execute(t *testing.T) {
	t.Parallel()
	smsProviderErr := errors.New("failed to get active sms provider")
	codeErr := errors.New("failed to create code")
	defaultExpiry := 10 * time.Minute
	expiry := 30 * time.Minute

	tests := []struct {
		name                        string
		requestChallengeOTPSMS      *session_grpc.RequestChallenges_OTPSMS
		sessionRepo                 func(ctrl *gomock.Controller) domain.SessionRepository
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		user                        *domain.User
		session                     *domain.Session
		secretGeneratorConfig       *crypto.GeneratorConfig
		otpAlgorithm                crypto.EncryptionAlgorithm
		smsProviderFn               func(ctx context.Context, instanceID string) (string, error)
		newPhoneCodeFn              func(g crypto.Generator) (*crypto.CryptoValue, string, error)
		wantErr                     error
	}{
		{
			name:                   "no otp sms challenge request",
			requestChallengeOTPSMS: nil,
			wantErr:                nil,
		},
		{
			name: "failed to get active sms provider",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: false,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", smsProviderErr
			},
			wantErr: smsProviderErr,
		},
		{
			name: "failed to retrieve secret generator config from settings DB",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: false,
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil // no external sms provider
			},
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewMockSecretGeneratorSettingsRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								database.And(
									getSettingsInstanceIDCondition(repo, "instance-id"),
									database.NewTextCondition(
										getSettingsTypeColumn(repo),
										database.TextOperationEqual,
										domain.SettingTypeSecretGenerator.String(),
									),
								),
							),
						),
					).AnyTimes().
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-kAcM0U", "failed to get OTP SMS secret generator config"),
		},
		{
			name:                   "failed to generate code",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil // no external sms provider
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive),
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return nil, "", codeErr
			},
			wantErr: codeErr,
		},
		{
			name: "update successful - external sms provider",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "external-sms-provider-id", nil // with an external sms provider
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedChallengeOTPSMS := &domain.SessionChallengeOTPSMS{
					Code:         nil,
					Expiry:       0,
					CodeReturned: true,
					GeneratorID:  "external-sms-provider-id",
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertOTPSMSChallengeChange(t, expectedChallengeOTPSMS))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(1), nil)
				return repo
			},
		},
		{
			name: "update successful - internal sms provider",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedChallengeOTPSMS := &domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       expiry,
					CodeReturned: true,
					GeneratorID:  "",
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertOTPSMSChallengeChange(t, expectedChallengeOTPSMS))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(1), nil)
				return repo
			},
		},
		{
			name: "update successful - internal sms provider - with default config",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStatePreview),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedChallengeOTPSMS := &domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       defaultExpiry,
					CodeReturned: true,
					GeneratorID:  "",
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertOTPSMSChallengeChange(t, expectedChallengeOTPSMS))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(1), nil)
				return repo
			},
		},
		{
			name: "failed to update session",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedChallengeOTPSMS := &domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       expiry,
					CodeReturned: true,
					GeneratorID:  "",
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertOTPSMSChallengeChange(t, expectedChallengeOTPSMS))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(0), assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-AigB0Z", "failed updating session"),
		},
		{
			name: "failed to update session - no rows updated",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedChallengeOTPSMS := &domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       expiry,
					CodeReturned: true,
					GeneratorID:  "",
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertOTPSMSChallengeChange(t, expectedChallengeOTPSMS))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(0), nil)
				return repo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-AigB0Z", "session not found"),
		},
		{
			name: "failed to update session - more than 1 row updated",
			requestChallengeOTPSMS: &session_grpc.RequestChallenges_OTPSMS{
				ReturnCode: true,
			},
			smsProviderFn: func(ctx context.Context, instanceID string) (string, error) {
				return "", nil
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive),
			newPhoneCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				expectedChallengeOTPSMS := &domain.SessionChallengeOTPSMS{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       expiry,
					CodeReturned: true,
					GeneratorID:  "",
				}
				repo.EXPECT().
					SetChallenge(gomock.Any()).
					AnyTimes().
					DoAndReturn(assertOTPSMSChallengeChange(t, expectedChallengeOTPSMS))
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, gomock.Any()).
					AnyTimes().
					Return(int64(2), nil)
				return repo
			},
			wantErr: zerrors.ThrowInternal(domain.NewMultipleObjectsUpdatedError(1, 2), "DOM-AigB0Z", "unexpected number of rows updated"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-id", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPSMSChallengeCommand(tt.requestChallengeOTPSMS,
				"session-id",
				"instance-id",
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.smsProviderFn,
				tt.newPhoneCodeFn,
			)
			cmd.User = tt.user
			cmd.Session = tt.session
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}
			if tt.secretGeneratorSettingsRepo != nil {
				domain.WithSecretGeneratorSettingsRepo(tt.secretGeneratorSettingsRepo(ctrl))(opts)
			}
			err := cmd.Execute(ctx, opts)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func assertOTPSMSChallengeChange(t *testing.T, expectedChallengeOTPSMS *domain.SessionChallengeOTPSMS) func(challenge *domain.SessionChallengeOTPSMS) database.Change {
	return func(challenge *domain.SessionChallengeOTPSMS) database.Change {
		assert.Equal(t, expectedChallengeOTPSMS.Code, challenge.Code)
		assert.Equal(t, expectedChallengeOTPSMS.Expiry, challenge.Expiry)
		assert.Equal(t, expectedChallengeOTPSMS.CodeReturned, challenge.CodeReturned)
		assert.Equal(t, expectedChallengeOTPSMS.GeneratorID, challenge.GeneratorID)

		changes := []database.Change{
			database.NewChange(
				database.NewColumn("zitadel.sessions", "otp_sms_challenge_expiry"), challenge.Expiry,
			),
			database.NewChange(
				database.NewColumn("zitadel.sessions", "otp_sms_challenge_code_returned"), challenge.CodeReturned,
			),
			database.NewChange(
				database.NewColumn("zitadel.sessions", "otp_sms_challenge_generator_id"), challenge.GeneratorID,
			),
		}
		if challenge.Code != nil { // is nil in the case of an external sms provider
			changes = append(changes, database.NewChange(
				database.NewColumn("zitadel.sessions", "otp_sms_challenge_code"), challenge.Code.Crypted,
			))
		}
		return database.NewChanges(changes...)
	}
}
