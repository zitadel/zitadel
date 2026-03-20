package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestDeleteSessionCommand_Validate(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())

	type args struct {
		sessionID      string
		sessionToken   string
		mustCheckToken bool
	}
	type fields struct {
		sessionRepo     func(ctrl *gomock.Controller) domain.SessionRepository
		permissionCheck func(ctrl *gomock.Controller) domain.PermissionChecker
	}
	tt := []struct {
		name          string
		args          args
		fields        fields
		expectedError error
	}{
		{
			name: "when sessionID is empty should return invalid argument error",
			args: args{
				sessionID: "",
			},
			fields:        fields{},
			expectedError: zerrors.ThrowInvalidArgument(nil, domain.ErrIDMissing, "Errors.IDMissing"),
		},
		{
			name: "when sessionID is provided must validate successfully",
			args: args{
				sessionID: "session",
			},
			fields:        fields{},
			expectedError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			d := domain.NewDeleteSessionCommand(tc.args.sessionID, tc.args.sessionToken, tc.args.mustCheckToken)
			ctrl := gomock.NewController(t)
			opts := domain.DefaultOpts(nil)
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.fields.sessionRepo != nil {
				domain.WithSessionRepo(tc.fields.sessionRepo(ctrl))(opts)
			}
			if tc.fields.permissionCheck != nil {
				domain.WithPermissionChecker(tc.fields.permissionCheck(ctrl))(opts)
			}

			// Test
			err := d.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeleteSessionCommand_Execute(t *testing.T) {
	t.Parallel()

	userID := gofakeit.UUID()
	ctx := authz.NewMockContext("inst-1", "org-1", userID)
	deleteErr := errors.New("delete error")

	type args struct {
		mockTx                func(ctrl *gomock.Controller) database.QueryExecutor
		sessionRepo           func(ctrl *gomock.Controller) domain.SessionRepository
		sessionTokenDecryptor func(t *testing.T) domain.SessionTokenDecryptor
	}
	type fields struct {
		sessionID           string
		sessionToken        string
		mustCheckPermission bool
	}
	type res struct {
		error     error
		deletedAt bool
		events    int
	}
	tt := []struct {
		testName string

		args   args
		fields fields
		res    res
	}{
		{
			testName: "when delete session fails should return error",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							nil,
						).
						Times(1).
						Return(int64(0), time.Time{}, deleteErr)
					return repo
				},
			},
			fields: fields{
				sessionID: "session-1",
			},
			res: res{
				error: deleteErr,
			},
		},
		{
			testName: "when no token is provided but must check permission, should use userID check",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							database.Or(
								repo.UserIDCondition(userID),
								database.RaisePermissionDeniedException(),
							),
						).
						Times(1).
						Return(int64(0), time.Time{}, nil)
					return repo
				},
			},
			fields: fields{
				sessionID:           "session-1",
				mustCheckPermission: true,
			},
			res: res{
				error: nil,
			},
		},
		{
			testName: "when token is provided and must check permission, should use tokenID check",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							database.Or(
								repo.TokenIDCondition("token-1"),
								database.RaisePermissionDeniedException(),
							),
						).
						Times(1).
						Return(int64(0), time.Time{}, nil)
					return repo
				},
				sessionTokenDecryptor: mockSessionTokenDecryptor("session-token-1", "session-1", "token-1"),
			},
			fields: fields{
				sessionID:           "session-1",
				sessionToken:        "session-token-1",
				mustCheckPermission: true,
			},
			res: res{
				error: nil,
			},
		},
		{
			testName: "when delete session without permission should return error",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							database.Or(
								repo.UserIDCondition(userID),
								database.RaisePermissionDeniedException(),
							),
						).
						Times(1).
						Return(int64(0), time.Time{}, database.NewPermissionError(nil))
					return repo
				},
			},
			fields: fields{
				sessionID:           "session-1",
				mustCheckPermission: true,
			},
			res: res{
				error: new(database.PermissionError),
			},
		},
		{
			testName: "when more than one row deleted must return internal error",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							nil,
						).
						Times(1).
						Return(int64(2), time.Time{}, nil)
					return repo
				},
			},
			fields: fields{
				sessionID: "session-1",
			},
			res: res{
				error: zerrors.ThrowInternalf(nil, domain.ErrMoreThanOneRowAffected, "expected 1 session to be deleted, got %d", 2),
			},
		},
		{
			testName: "when no rows deleted should execute successfully",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							nil,
						).
						Times(1).
						Return(int64(0), time.Time{}, nil)
					return repo
				},
			},
			fields: fields{
				sessionID: "session-1",
			},
			res: res{
				error: nil,
			},
		},
		{
			testName: "when one row deleted should execute successfully and set DeletedAt",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							nil,
						).
						Times(1).
						Return(int64(1), time.Now(), nil)
					return repo
				},
			},
			fields: fields{
				sessionID: "session-1",
			},
			res: res{
				deletedAt: true,
				events:    1,
			},
		},
		{
			testName: "when session already deleted should execute successfully, set DeletedAt but no event",
			args: args{
				sessionRepo: func(ctrl *gomock.Controller) domain.SessionRepository {
					repo := domainmock.NewSessionRepo(ctrl)
					repo.EXPECT().
						Delete(gomock.Any(), gomock.Any(),
							repo.PrimaryKeyCondition("inst-1", "session-1"),
							nil,
						).
						Times(1).
						Return(int64(0), time.Now(), nil)
					return repo
				},
			},
			fields: fields{
				sessionID: "session-1",
			},
			res: res{
				deletedAt: true,
				events:    0,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewDeleteSessionCommand(tc.fields.sessionID, tc.fields.sessionToken, tc.fields.mustCheckPermission)
			ctrl := gomock.NewController(t)
			opts := domain.DefaultOpts(nil)
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.args.mockTx != nil {
				domain.WithQueryExecutor(tc.args.mockTx(ctrl))(opts)
			}
			if tc.args.sessionRepo != nil {
				domain.WithSessionRepo(tc.args.sessionRepo(ctrl))(opts)
			}
			if tc.args.sessionTokenDecryptor != nil {
				domain.WithSessionTokenDecryptor(tc.args.sessionTokenDecryptor(t))(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.res.error)
			assert.Equal(t, tc.res.deletedAt, !cmd.DeletedAt.IsZero())

			// Test events
			cmds, err := cmd.Events(ctx, opts)
			assert.NoError(t, err)
			assert.Len(t, cmds, tc.res.events)
		})
	}
}

func mockSessionTokenDecryptor(expectedSessionToken, returnSessionID, returnSessionToken string) func(t *testing.T) domain.SessionTokenDecryptor {
	return func(t *testing.T) domain.SessionTokenDecryptor {
		return func(ctx context.Context, sessionToken string) (sessionID, tokenID string, err error) {
			if sessionToken != expectedSessionToken {
				return "", "", zerrors.ThrowInvalidArgumentf(nil, "SESS-S3gq1", "Errors.Session.TokenInvalid")
			}
			return returnSessionID, returnSessionToken, nil
		}
	}
}
