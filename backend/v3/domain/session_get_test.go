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
		userRepo             func(ctrl *gomock.Controller) domain.UserRepository
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
			name: "skip permission check - requestor retrieving their own session - successful",
			ctx:  authz.NewMockContext("instance-1", "", "user-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			},
		},
		{
			name: "skip permission check - requestor retrieving session of a different user, but same user agent - successful",
			ctx:  authz.NewMockContextWithAgent("instance-1", "", "requestor-1", "fingerprint-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
				UserAgent: &domain.SessionUserAgent{
					InstanceID:    "instance-1",
					FingerprintID: gu.Ptr("fingerprint-1"),
				},
			}, nil),
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
				UserAgent: &domain.SessionUserAgent{
					InstanceID:    "instance-1",
					FingerprintID: gu.Ptr("fingerprint-1"),
				},
			},
		},
		{
			name: "check requestor's permission on the session user's organization - failed to get user",
			ctx:  authz.NewMockContext("instance-1", "", "requestor-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			userRepo: getUserRepo(nil, assert.AnError),
			wantErr:  zerrors.ThrowInternal(assert.AnError, "QUERY-4n9fGv", "Errors.Get.User"),
		},
		{
			name: "check requestor's permission on the session user's organization - user not found",
			ctx:  authz.NewMockContext("instance-1", "", "requestor-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			userRepo: getUserRepo(nil, new(database.NoRowFoundError)),
			wantErr:  zerrors.ThrowNotFound(new(database.NoRowFoundError), "QUERY-4n9fGv", "Errors.NotFound.User"),
		},
		{
			name: "requestor with no permissions on the session user's organization - permission denied",
			ctx:  authz.NewMockContext("instance-1", "", "requestor-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			userRepo: getUserRepo(&domain.User{
				ID:             "user-1",
				InstanceID:     "instance-1",
				OrganizationID: "org-1",
			}, nil),
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permissionChecker := domainmock.NewMockPermissionChecker(ctrl)
				permissionChecker.EXPECT().
					CheckOrganizationPermission(gomock.Any(), domain.SessionReadPermission, "org-1").
					Times(1).
					Return(assert.AnError)
				return permissionChecker
			},
			wantErr: zerrors.ThrowPermissionDenied(assert.AnError, "QUERY-6l8oWp", "Errors.PermissionDenied"),
		},
		{
			name: "requestor with permissions on the session user's organization - successful",
			ctx:  authz.NewMockContext("instance-1", "", "requestor-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			userRepo: getUserRepo(&domain.User{
				ID:             "user-1",
				InstanceID:     "instance-1",
				OrganizationID: "org-1",
			}, nil),
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permissionChecker := domainmock.NewMockPermissionChecker(ctrl)
				permissionChecker.EXPECT().
					CheckOrganizationPermission(gomock.Any(), domain.SessionReadPermission, "org-1").
					Times(1).
					Return(nil)
				return permissionChecker
			},
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			},
		},
		{
			name: "requestor with no permissions on the instance-level - permission denied",
			ctx:  authz.NewMockContext("instance-1", "", "requestor-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			userRepo: getUserRepo(&domain.User{
				ID:             "user-1",
				InstanceID:     "instance-1",
				OrganizationID: "",
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
			name: "requestor with sufficient permissions on the instance-level - successful",
			ctx:  authz.NewMockContext("instance-1", "", "requestor-1"),
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				UserID:     "user-1",
			}, nil),
			userRepo: getUserRepo(&domain.User{
				ID:             "user-1",
				InstanceID:     "instance-1",
				OrganizationID: "",
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
				UserID:     "user-1",
			},
		},
		{
			name:         "session token verification failed",
			ctx:          context.Background(),
			sessionToken: "session_token",
			sessionTokenVerifier: func(ctx context.Context, sessionToken string, sessionID string, tokenID string) error {
				return assert.AnError
			},
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				Token:      "sessiontoken-1",
			}, nil),
			wantErr: zerrors.ThrowPermissionDenied(assert.AnError, "QUERY-M3f4fS", "Errors.PermissionDenied"),
		},
		{
			name:         "session token verification succeeded",
			ctx:          context.Background(),
			sessionToken: "session_token",
			sessionTokenVerifier: func(ctx context.Context, sessionToken string, sessionID string, tokenID string) error {
				return nil
			},
			sessionRepo: getSessionRepo(&domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				Token:      "sessiontoken-1",
			}, nil),
			wantSession: &domain.Session{
				ID:         "session-1",
				InstanceID: "instance-1",
				Token:      "sessiontoken-1",
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
			if tt.userRepo != nil {
				domain.WithUserRepo(tt.userRepo(ctrl))(opts)
			}

			err := query.Execute(tt.ctx, opts)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantSession, query.Result())
		})
	}
}

func getUserRepo(user *domain.User, err error) func(ctrl *gomock.Controller) domain.UserRepository {
	return func(ctrl *gomock.Controller) domain.UserRepository {
		repo := domainmock.NewMockUserRepository(ctrl)
		repo.EXPECT().
			Get(gomock.Any(),
				gomock.Any(),
				dbmock.QueryOptions(
					database.WithCondition(
						getUserIDCondition(repo, "user-1"),
					),
				),
			).
			Times(1).
			Return(user, err)
		return repo
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
