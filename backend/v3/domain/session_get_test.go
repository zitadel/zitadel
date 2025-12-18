package domain_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestGetSessionQuery_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		sessionID    string
		instanceID   string
		sessionToken string
		wantErr      error
	}{
		{
			name:       "missing session ID",
			sessionID:  "     ",
			instanceID: "instance-1",
			wantErr:    zerrors.ThrowPreconditionFailed(nil, "QUERY-CtWgrV", "Errors.Missing.SessionID"),
		},
		{
			name:       "missing instance ID",
			sessionID:  "session-1",
			instanceID: "",
			wantErr:    zerrors.ThrowPreconditionFailed(nil, "QUERY-3n9fGv", "Errors.Missing.InstanceID"),
		},
		{
			name:       "valid request",
			sessionID:  "session-1",
			instanceID: "instance-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			query := domain.NewGetSessionQuery(tt.sessionID, tt.instanceID, tt.sessionToken, nil)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			err := query.Validate(context.Background(), opts)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetSessionQuery_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ctx                  context.Context
		name                 string
		sessionToken         string
		sessionTokenVerifier func(ctx context.Context, sessionToken string, sessionID string, tokenID string) (err error)
		sessionRepo          func(ctrl *gomock.Controller) domain.SessionRepository
		permissionChecker    func(ctrl *gomock.Controller) domain.PermissionChecker
		wantErr              error
		wantSession          *domain.Session
	}{
		{
			name:        "failed to get session",
			ctx:         context.Background(),
			sessionRepo: getSessionRepo(nil, assert.AnError),
			wantErr:     zerrors.ThrowInternal(assert.AnError, "DOM-QiiiFY", "Errors.Get.Session"),
		},
		{
			name:        "session not found",
			ctx:         context.Background(),
			sessionRepo: getSessionRepo(nil, new(database.NoRowFoundError)),
			wantErr:     zerrors.ThrowNotFound(new(database.NoRowFoundError), "DOM-QiiiFY", "Errors.NotFound.Session"),
		},
		{
			name: "skip permission check - user retrieving their own session",
			ctx:  authz.NewMockContext("instance-1", "", "user-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				Factors: []domain.SessionFactor{
					&domain.SessionFactorUser{
						UserID: "user-1",
					},
				},
			}, nil),
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				Factors: []domain.SessionFactor{
					&domain.SessionFactorUser{
						UserID: "user-1",
					},
				},
			},
		},
		{
			name: "skip permission check - same user agent",
			ctx:  authz.NewMockContextWithAgent("instance-1", "", "user-1", "fingerprint-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserAgent: &domain.SessionUserAgent{
					InstanceID:    "instance-1",
					FingerprintID: gu.Ptr("fingerprint-1"),
				},
			}, nil),
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserAgent: &domain.SessionUserAgent{
					InstanceID:    "instance-1",
					FingerprintID: gu.Ptr("fingerprint-1"),
				},
			},
		},
		{
			name: "user with no permissions - permission denied",
			ctx:  authz.NewMockContextWithPermissions("instance-1", "", "user-1", nil),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
			}, nil),
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permissionChecker := domainmock.NewMockPermissionChecker(ctrl)
				permissionChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.SessionReadPermission).
					Times(1).
					Return(assert.AnError)
				return permissionChecker
			},
			wantErr: zerrors.ThrowPermissionDenied(assert.AnError, "QUERY-5RJSUU", "Errors.PermissionDenied"),
		},
		{
			name: "user with sufficient permissions - success",
			ctx:  authz.NewMockContextWithPermissions("instance-1", "", "user-1", nil),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
			}, nil),
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permissionChecker := domainmock.NewMockPermissionChecker(ctrl)
				permissionChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.SessionReadPermission).
					Times(1).
					Return(nil)
				return permissionChecker
			},
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
			},
		},
		{
			name:         "session token verification failed",
			ctx:          context.Background(),
			sessionToken: "random",
			sessionTokenVerifier: func(ctx context.Context, sessionToken string, sessionID string, tokenID string) error {
				return assert.AnError
			},
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
			}, nil),
			wantErr: zerrors.ThrowPermissionDenied(assert.AnError, "QUERY-M3f4fS", "Errors.PermissionDenied"),
		},
		{
			name:         "session token verification succeeded",
			ctx:          context.Background(),
			sessionToken: "random",
			sessionTokenVerifier: func(ctx context.Context, sessionToken string, sessionID string, tokenID string) error {
				return nil
			},
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
			}, nil),
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			query := domain.NewGetSessionQuery("session-1", "instance-1", tt.sessionToken, tt.sessionTokenVerifier)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			ctrl := gomock.NewController(t)
			if tt.permissionChecker != nil {
				opts.Permissions = tt.permissionChecker(ctrl)
			}
			if tt.sessionRepo != nil {
				domain.WithSessionRepo(tt.sessionRepo(ctrl))(opts)
			}

			err := query.Execute(tt.ctx, opts)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantSession, query.Result())
		})
	}
}

func getSessionRepo(session *domain.Session, err error) func(ctrl *gomock.Controller) domain.SessionRepository {
	return func(ctrl *gomock.Controller) domain.SessionRepository {
		repo := domainmock.NewMockSessionRepository(ctrl)
		repo.EXPECT().
			Get(gomock.Any(),
				gomock.Any(),
				dbmock.QueryOptions(
					database.WithCondition(
						database.And(
							getSessionInstanceIDCondition(repo, "instance-1"),
							getSessionIDCondition(repo, "session-1"),
						),
					),
				),
			).
			Times(1).
			Return(session, err)
		return repo
	}
}
