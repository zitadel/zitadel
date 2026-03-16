package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestRecoveryCodeCheckCommand_Validate(t *testing.T) {
	t.Parallel()

	getError := errors.New("get error")

	tests := []struct {
		name        string
		sessionID   string
		instanceID  string
		check       *domain.CheckRecoveryCode
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo    func(ctrl *gomock.Controller) domain.UserRepository
		wantErr     error
	}{
		{
			name: "no recovery code check to be done",
		},
		{
			name:    "empty recovery code",
			check:   &domain.CheckRecoveryCode{},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-cEKxoG", "Errors.user.MFA.RecoveryCodes.Empty"),
		},
		{
			name:    "no session id",
			check:   &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-hsbVyd", "Errors.session.IDMissing"),
		},
		{
			name:       "no instance id",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			wantErr:    zerrors.ThrowPreconditionFailed(nil, "DOM-lGIe1v", "Errors.Instance.IDMissing"),
		},
		{
			name:       "failed to get session - session not found",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, &database.NoRowFoundError{})
				return sessionRepo
			},
			wantErr: zerrors.ThrowNotFound(&database.NoRowFoundError{}, "DOM-Ot3qO6", "Errors.session.NotFound"),
		},
		{
			name:       "failed to get session - database error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, getError)
				return sessionRepo
			},
			wantErr: zerrors.ThrowInternal(getError, "DOM-2sF2kF", "Errors.Internal"),
		},
		{
			name:       "missing user id in session",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{}, nil)
				return sessionRepo
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOM-EaLqwq", "Errors.user.UserIDMissing"),
		},
		{
			name:       "failed to get user - user not found",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, &database.NoRowFoundError{})
				return userRepo
			},
			wantErr: zerrors.ThrowNotFound(&database.NoRowFoundError{}, "DOM-Ot3qO6", "Errors.user.NotFound"),
		},
		{
			name:       "failed to get user - database error",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, getError)
				return userRepo
			},
			wantErr: zerrors.ThrowInternal(getError, "DOM-7sWTNf", "Errors.Internal"),
		},
		{
			name:       "user is locked",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.User{State: domain.UserStateLocked}, nil)
				return userRepo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-47H1Ii", "Errors.user.Locked"),
		},
		{
			name:       "user does not have recovery codes - not a human user",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(
						&domain.User{
							InstanceID:     "instance-1",
							OrganizationID: "org-1",
							State:          domain.UserStateActive,
						}, nil)
				return userRepo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-tzN2a1", "Errors.user.MFA.RecoveryCodes.NotReady"),
		},
		{
			name:       "user does not have recovery codes",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(
						&domain.User{
							InstanceID:     "instance-1",
							OrganizationID: "org-1",
							State:          domain.UserStateActive,
							Human:          &domain.HumanUser{},
						}, nil)
				return userRepo
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "DOM-tzN2a1", "Errors.user.MFA.RecoveryCodes.NotReady"),
		},
		{
			name:       "valid reqquest",
			sessionID:  "session-1",
			instanceID: "instance-1",
			check:      &domain.CheckRecoveryCode{RecoveryCode: "test-code"},
			sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
				sessionRepo := domainmock.NewSessionRepo(ctrl)
				sessionRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&domain.Session{UserID: "user-1"}, nil)
				return sessionRepo
			},
			userRepo: func(ctrl *gomock.Controller) domain.UserRepository {
				userRepo := domainmock.NewUserRepo(ctrl)
				userRepo.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(
						&domain.User{
							InstanceID:     "instance-1",
							OrganizationID: "org-1",
							State:          domain.UserStateActive,
							Human: &domain.HumanUser{
								RecoveryCodes: &domain.HumanRecoveryCodes{
									Codes: []string{"code1", "code2", "code3"},
								},
							},
						}, nil)
				return userRepo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

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

			cmd := domain.NewRecoveryCodeCheckCommand(tt.sessionID, tt.instanceID, tt.check, nil)
			got := cmd.Validate(context.Background(), opts)
			assert.ErrorIs(t, tt.wantErr, got)
		})
	}
}

func TestRecoveryCodeCheckCommand_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		sessionID   string
		instanceID  string
		check       *domain.CheckRecoveryCode
		hasher      *crypto.Hasher
		sessionRepo func(ctrl *gomock.Controller) domain.SessionRepository
		userRepo    func(ctrl *gomock.Controller) domain.UserRepository
		wantErr     error
	}{
		{
			name: "no recovery code check to be done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

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

			cmd := domain.NewRecoveryCodeCheckCommand(tt.sessionID, tt.instanceID, tt.check, tt.hasher)
			got := cmd.Execute(context.Background(), opts)
			assert.ErrorIs(t, tt.wantErr, got)
		})
	}
}
