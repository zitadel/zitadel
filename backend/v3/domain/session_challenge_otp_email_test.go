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
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestOTPEmailChallengeCommand_Validate(t *testing.T) {
	t.Parallel()
	otpEnabledAt := time.Now().Add(-30 * time.Minute)

	tests := []struct {
		name                string
		sessionID           string
		instanceID          string
		requestChallengeOTP *session_grpc.RequestChallenges_OTPEmail
		userRepo            func(ctrl *gomock.Controller) domain.UserRepository
		sessionRepo         func(ctrl *gomock.Controller) domain.SessionRepository
		wantErr             error
		wantUser            *domain.User
		wantSession         *domain.Session
	}{
		{
			name:                "no request otp email challenge",
			requestChallengeOTP: nil,
			wantErr:             nil,
		},
		{
			name:                "no session id",
			sessionID:           "",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-BQ5UgK", "Errors.Missing.SessionID"),
		},
		{
			name:                "no instance id",
			sessionID:           "session-id",
			instanceID:          "",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			wantErr:             zerrors.ThrowPreconditionFailed(nil, "DOM-kDnkDn", "Errors.Missing.InstanceID"),
		},
		{
			name:                "failed to get session",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-JArUai", "Errors.Get.Session"),
		},
		{
			name:                "session not found",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-JArUai", "Errors.NotFound.Session"),
		},
		{
			name:                "missing user id in session",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
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
			name:                "failed to get user",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-56MWkg", "Errors.Get.User"),
		},
		{
			name:                "user not found",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
					Return(nil, new(database.NoRowFoundError))
				return repo
			},
			wantErr: zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-56MWkg", "Errors.NotFound.User"),
		},
		{
			name:                "human user not set",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-id",
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
			name:                "email not set",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
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
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-7hG2d", "Errors.NotFound.User.Human.Email"),
		},
		{
			name:                "email address not set",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
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
			name:                "email otp not enabled",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID: "user-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
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
							Email: domain.HumanEmail{
								Address: "email@example.com",
								OTP:     domain.OTP{},
							},
						},
					}, nil)
				return repo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-9kL4q", "Errors.OTPEmail.NotEnabled"),
		},
		{
			name:                "valid OTP Email challenge request",
			sessionID:           "session-id",
			instanceID:          "instance-id",
			requestChallengeOTP: &session_grpc.RequestChallenges_OTPEmail{},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(), dbmock.QueryOptions(
						database.WithCondition(
							getSessionIDCondition(repo, "session-id"),
						),
					)).
					Times(1).
					Return(&domain.Session{
						UserID:     "user-id",
						InstanceID: "instance-id",
						ID:         "session-id",
					}, nil)
				return repo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				repo := domainmock.NewMockUserRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								getUserIDCondition(repo, "user-id"),
							),
						)).
					Times(1).
					Return(&domain.User{
						InstanceID:     "instance-id",
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
			wantUser: &domain.User{
				InstanceID:     "instance-id",
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
			cmd := domain.NewOTPEmailChallengeCommand(
				tt.sessionID,
				tt.instanceID,
				tt.requestChallengeOTP,
				nil,
				nil,
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
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantUser, cmd.User)
			assert.Equal(t, tt.wantSession, cmd.Session)
		})
	}
}

func TestOTPEmailChallengeCommand_Events(t *testing.T) {
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
		name                     string
		requestChallengeOTPEmail *session_grpc.RequestChallenges_OTPEmail
		challengeOTPEmail        *domain.SessionChallengeOTPEmail
		wantErr                  error
		wantEvent                eventstore.Command
	}{
		{
			name:                     "no request otp email challenge",
			requestChallengeOTPEmail: nil,
		},
		{
			name:                     "valid OTP email challenge request - no delivery type",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{},
			challengeOTPEmail: &domain.SessionChallengeOTPEmail{
				LastChallengedAt: challengedAt,
				Code:             code,
				Expiry:           expiry,
			},
			wantEvent: session.NewOTPEmailChallengedEvent(
				ctx,
				&session.NewAggregate("session-id", "instance-id").Aggregate,
				code,
				expiry,
				false,
				"",
			),
		},
		{
			name: "valid OTP email challenge request - delivery type send code",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_SendCode_{
					SendCode: &session_grpc.RequestChallenges_OTPEmail_SendCode{
						UrlTemplate: gu.Ptr("http://example.com"),
					},
				},
			},
			challengeOTPEmail: &domain.SessionChallengeOTPEmail{
				LastChallengedAt: challengedAt,
				Code:             code,
				Expiry:           expiry,
				CodeReturned:     false,
				URLTmpl:          "https://example.com",
			},
			wantEvent: session.NewOTPEmailChallengedEvent(
				ctx,
				&session.NewAggregate("session-id", "instance-id").Aggregate,
				code,
				expiry,
				false,
				"https://example.com",
			),
		},
		{
			name: "valid OTP email challenge request - delivery type return code",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_ReturnCode_{
					ReturnCode: &session_grpc.RequestChallenges_OTPEmail_ReturnCode{},
				},
			},
			challengeOTPEmail: &domain.SessionChallengeOTPEmail{
				LastChallengedAt: challengedAt,
				Code:             code,
				Expiry:           expiry,
				CodeReturned:     true,
				URLTmpl:          "",
			},
			wantEvent: session.NewOTPEmailChallengedEvent(
				ctx,
				&session.NewAggregate("session-id", "instance-id").Aggregate,
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
				"session-id",
				"instance-id",
				tt.requestChallengeOTPEmail,
				nil,
				nil,
				nil,
			)
			cmd.ChallengeOTPEmail = tt.challengeOTPEmail
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			events, err := cmd.Events(ctx, opts)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantEvent != nil {
				require.Len(t, events, 1)
				assert.Equal(t, tt.wantEvent, events[0])
			}
		})
	}
}

func TestOTPEmailChallengeCommand_Execute(t *testing.T) {
	t.Parallel()
	codeErr := errors.New("failed to create code")
	defaultExpiry := 10 * time.Minute
	expiry := 30 * time.Minute

	tests := []struct {
		name                        string
		requestChallengeOTPEmail    *session_grpc.RequestChallenges_OTPEmail
		sessionRepo                 func(ctrl *gomock.Controller) domain.SessionRepository
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		user                        *domain.User
		session                     *domain.Session
		secretGeneratorConfig       *crypto.GeneratorConfig
		otpAlgorithm                crypto.EncryptionAlgorithm
		newEmailCodeFn              func(g crypto.Generator) (*crypto.CryptoValue, string, error)
		wantErr                     error
		wantOTPEmailChallenge       string
		wantChallengeOTPEmail       *domain.SessionChallengeOTPEmail
	}{
		{
			name:                     "no otp email challenge request",
			requestChallengeOTPEmail: nil,
			wantErr:                  nil,
		},
		{
			name:                     "failed to retrieve secret generator config from settings DB",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{},
			secretGeneratorConfig:    &crypto.GeneratorConfig{},
			otpAlgorithm:             crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewMockSecretGeneratorSettingsRepository(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								database.And(
									getSettingsInstanceIDCondition(repo, "instance-id"),
									getSettingsTypeCondition(repo, domain.SettingTypeSecretGenerator),
								),
							),
						),
					).Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-x7Yd3E", "Errors.Get.SecretGeneratorSettings"),
		},
		{
			name:                        "failed to generate code",
			requestChallengeOTPEmail:    &session_grpc.RequestChallenges_OTPEmail{},
			secretGeneratorConfig:       &crypto.GeneratorConfig{},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return nil, "", codeErr
			},
			wantErr: codeErr,
		},
		{
			name: "failed to render url template",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_SendCode_{
					SendCode: &session_grpc.RequestChallenges_OTPEmail_SendCode{
						UrlTemplate: gu.Ptr("http://{{.Invalid"),
					},
				},
			},
			secretGeneratorConfig:       &crypto.GeneratorConfig{},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-wkDwQM", "Errors.Invalid.URLTemplate"),
		},
		{
			name: "failed to update session",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_ReturnCode_{},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: expiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(0), assert.AnError)
				return repo
			},
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-YfQIA3", "Errors.Update.Session"),
		},
		{
			name: "failed to update session - no rows updated",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_SendCode_{
					SendCode: &session_grpc.RequestChallenges_OTPEmail_SendCode{
						UrlTemplate: gu.Ptr("http://example.com"),
					},
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
				Length: uint(6),
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:            expiry,
					URLTmpl:           "http://example.com",
					LastChallengedAt:  time.Now(),
					TriggeredAtOrigin: "://",
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(0), nil)
				return repo
			},
			wantErr: zerrors.ThrowNotFound(nil, "DOM-YfQIA3", "Errors.NotFound.Session"),
		},
		{
			name: "failed to update session - more than 1 row updated",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_ReturnCode_{},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
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
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					URLTmpl:          "",
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(2), nil)
				return repo
			},
			wantErr: zerrors.ThrowInternal(nil, "DOM-YfQIA3", "Errors.Update.MultipleRows"),
		},
		{
			name: "update session succeeded - delivery type return code",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_ReturnCode_{},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     true,
					URLTmpl:          "",
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			wantOTPEmailChallenge: "code", // OTPEmailChallenge is only set when the delivery type is return code
			wantChallengeOTPEmail: &domain.SessionChallengeOTPEmail{
				Code: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				Expiry:            defaultExpiry,
				CodeReturned:      true,
				URLTmpl:           "",
				LastChallengedAt:  time.Now(),
				TriggeredAtOrigin: "://",
			},
		},
		{
			name: "update session succeeded - delivery type send code",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{
				DeliveryType: &session_grpc.RequestChallenges_OTPEmail_SendCode_{
					SendCode: &session_grpc.RequestChallenges_OTPEmail_SendCode{
						UrlTemplate: gu.Ptr("http://example.com"),
					},
				},
			},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					URLTmpl:          "http://example.com",
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			wantChallengeOTPEmail: &domain.SessionChallengeOTPEmail{
				Code: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				Expiry:            defaultExpiry,
				CodeReturned:      false,
				URLTmpl:           "http://example.com",
				LastChallengedAt:  time.Now(),
				TriggeredAtOrigin: "://",
			},
		},
		{
			name:                     "update session succeeded - delivery type default",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry: defaultExpiry,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           expiry,
					CodeReturned:     false,
					URLTmpl:          "",
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			wantChallengeOTPEmail: &domain.SessionChallengeOTPEmail{
				Code: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				Expiry:            defaultExpiry,
				CodeReturned:      false,
				URLTmpl:           "",
				LastChallengedAt:  time.Now(),
				TriggeredAtOrigin: "://",
			},
		},
		{
			name:                     "update session succeeded - with default config when setting state is not active",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{},
			secretGeneratorConfig: &crypto.GeneratorConfig{
				Expiry:              defaultExpiry,
				Length:              uint(8),
				IncludeLowerLetters: false,
				IncludeUpperLetters: true,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
			otpAlgorithm:                crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStatePreview, domain.OTPEmailRequestType),
			newEmailCodeFn: func(g crypto.Generator) (*crypto.CryptoValue, string, error) {
				return &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				}, "code", nil
			},
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           defaultExpiry,
					CodeReturned:     false,
					URLTmpl:          "",
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			wantChallengeOTPEmail: &domain.SessionChallengeOTPEmail{
				Code: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				Expiry:            defaultExpiry,
				CodeReturned:      false,
				URLTmpl:           "",
				LastChallengedAt:  time.Now(),
				TriggeredAtOrigin: "://",
			},
		},
		{
			name:                     "update session succeeded - with default config when OTP email settings not found",
			requestChallengeOTPEmail: &session_grpc.RequestChallenges_OTPEmail{},
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
			session: &domain.Session{
				ID: "session-id",
			},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				repo := domainmock.NewMockSessionRepository(ctrl)
				challengeOTPEmail := getOTPEmailChallengeChange(repo, &domain.SessionChallengeOTPEmail{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:           defaultExpiry,
					CodeReturned:     false,
					URLTmpl:          "",
					LastChallengedAt: time.Now(),
				})
				idCondition := getSessionIDCondition(repo, "session-id")
				repo.EXPECT().
					Update(gomock.Any(), gomock.Any(), idCondition, challengeOTPEmail).
					Times(1).
					Return(int64(1), nil)
				return repo
			},
			wantChallengeOTPEmail: &domain.SessionChallengeOTPEmail{
				Code: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("code"),
				},
				Expiry:            defaultExpiry,
				CodeReturned:      false,
				URLTmpl:           "",
				LastChallengedAt:  time.Now(),
				TriggeredAtOrigin: "://",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := authz.NewMockContext("instance-id", "", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewOTPEmailChallengeCommand(
				"session-id",
				"instance-id",
				tt.requestChallengeOTPEmail,
				tt.secretGeneratorConfig,
				tt.otpAlgorithm,
				tt.newEmailCodeFn,
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
			after := time.Now()
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantOTPEmailChallenge, gu.Value(cmd.OTPEmailChallenge))
			assertChallengeOTPEmailEqual(t, tt.wantChallengeOTPEmail, cmd.ChallengeOTPEmail, after)
		})
	}
}

func assertChallengeOTPEmailEqual(t *testing.T, want, got *domain.SessionChallengeOTPEmail, upperThreshold time.Time) {
	if want == nil {
		assert.Nil(t, got)
		return
	}
	assert.WithinRange(t, got.LastChallengedAt, want.LastChallengedAt, upperThreshold)
	assert.Equal(t, want.Code, got.Code)
	assert.Equal(t, want.Expiry, want.Expiry)
	assert.Equal(t, want.CodeReturned, got.CodeReturned)
	assert.Equal(t, want.URLTmpl, got.URLTmpl)
	assert.Equal(t, want.TriggeredAtOrigin, got.TriggeredAtOrigin)
}

func secretGeneratorSettingsRepo(state domain.SettingState, otpType domain.OTPRequestType) func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
	var attrs domain.SecretGeneratorSettingsAttributes
	switch otpType {
	case domain.OTPEmailRequestType:
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
	case domain.OTPSMSRequestType:
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
		repo := domainmock.NewMockSecretGeneratorSettingsRepository(ctrl)
		repo.EXPECT().
			Get(gomock.Any(), gomock.Any(),
				dbmock.QueryOptions(
					database.WithCondition(
						database.And(
							getSettingsInstanceIDCondition(repo, "instance-id"),
							getSettingsTypeCondition(repo, domain.SettingTypeSecretGenerator),
						),
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

func getSettingsInstanceIDCondition(repo *domainmock.MockSecretGeneratorSettingsRepository, instanceID string) database.Condition {
	instanceIDCondition := getTextCondition("zitadel.settings", "instance_id", instanceID)

	repo.EXPECT().
		InstanceIDCondition(instanceID).
		Times(1).
		Return(instanceIDCondition)
	return instanceIDCondition
}

func getSettingsTypeCondition(repo *domainmock.MockSecretGeneratorSettingsRepository, settingType domain.SettingType) database.Condition {
	typeCondition := getTextCondition("zitadel.settings", "type", settingType.String())

	repo.EXPECT().
		TypeCondition(settingType).
		Times(1).
		Return(typeCondition)
	return typeCondition
}

func getOTPEmailChallengeChange(repo *domainmock.MockSessionRepository, challenge *domain.SessionChallengeOTPEmail) database.Change {
	changes := database.NewChanges(
		database.NewChange(
			database.NewColumn("zitadel.sessions", "otp_email_challenge_code_crypted"), challenge.Code.Crypted,
		),
		database.NewChange(
			database.NewColumn("zitadel.sessions", "otp_email_challenge_expiry"), challenge.Expiry,
		),
		database.NewChange(
			database.NewColumn("zitadel.sessions", "otp_email_challenge_code_returned"), challenge.CodeReturned,
		),
		database.NewChange(
			database.NewColumn("zitadel.sessions", "otp_email_challenge_url_template"), challenge.URLTmpl,
		),
	)
	repo.EXPECT().
		SetChallenge(challenge).
		AnyTimes().
		Return(changes)
	return changes
}
