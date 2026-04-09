package domain_test

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCheckUserCommand_Validate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tt := []struct {
		testName      string
		cmd           *domain.CheckUserCommand
		expectedError error
	}{
		{
			testName:      "when userID and loginName are nil should return error",
			cmd:           domain.NewCheckUserCommand(domainmock.InitCheckUserParent(ctrl), nil, nil),
			expectedError: &zerrors.ZitadelError{Kind: zerrors.KindInvalidArgument},
		},
		{
			testName:      "userID is set should return no error",
			cmd:           domain.NewCheckUserCommand(domainmock.InitCheckUserParent(ctrl), gu.Ptr("user-id"), nil),
			expectedError: nil,
		},
		{
			testName:      "loginName is set should return no error",
			cmd:           domain.NewCheckUserCommand(domainmock.InitCheckUserParent(ctrl), nil, gu.Ptr("login-name")),
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			err := tc.cmd.Validate(t.Context(), new(domain.InvokeOpts))

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestCheckUserCommand_Execute(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	fetchUserResult := &domain.User{
		ID:             "user-123",
		OrganizationID: "org-1",
		InstanceID:     "instance-1",
		Username:       "username",
		State:          domain.UserStateActive,
		Human: &domain.HumanUser{
			PreferredLanguage: language.Nepali,
		},
	}
	fetchSessionResult := &domain.Session{
		ID:         "session-1",
		InstanceID: "instance-1",
	}

	tt := []struct {
		testName         string
		cmd              *domain.CheckUserCommand
		queryExecutor    database.QueryExecutor
		expectedError    error
		shouldHaveResult bool
	}{
		{
			testName: "when isolation fails should return error",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl),
				gu.Ptr("user-123"), nil,
			),
			expectedError:    &zerrors.ZitadelError{Kind: zerrors.KindInternal, ID: "some-id"},
			shouldHaveResult: false,
			queryExecutor: dbmock.NewMockTransaction(ctrl).
				AddExpectation(func(recorder *dbmock.MockTransactionMockRecorder) {
					recorder.Begin(gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "some-id", "some message"))
				}),
		},
		{
			testName: "when user is not found should return no error",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(nil, zerrors.ThrowNotFound(nil, "DOM-9s8f1", "user not found"))
					}),
				gu.Ptr("user-123"), nil,
			),
			expectedError:    nil,
			shouldHaveResult: false,
		},
		{
			testName: "when user retrieval fails should return error",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "DOM-9s8f1", "something went wrong"))
					}),
				gu.Ptr("user-123"), nil,
			),
			expectedError: &zerrors.ZitadelError{Kind: zerrors.KindInternal, ID: "DOM-9s8f1"},
		},
		{
			testName: "when session retrieval fails should return error",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(fetchUserResult, nil)
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(nil, zerrors.ThrowNotFound(nil, "nf", "not found"))
					}),
				gu.Ptr("user-123"), nil,
			),
			expectedError: zerrors.ThrowNotFound(nil, "nf", "not found"),
		},
		{
			testName: "when user change is attempted should return invalid argument error",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(fetchUserResult, nil)
						session := *fetchSessionResult
						session.UserID = "user-456"
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(&session, nil)
					}),
				nil, gu.Ptr("different-user"),
			),
			expectedError: zerrors.ThrowInvalidArgument(nil, "DOM-78g1TV", "user change not possible"),
		},
		{
			testName: "when user is not active should return error",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						user := *fetchUserResult
						user.State = domain.UserStateInactive
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(&user, nil)
					}),
				nil, gu.Ptr("different-user"),
			),
			expectedError: zerrors.ThrowPreconditionFailed(nil, "DOM-vgDIu9", "Errors.User.NotActive"),
		},
		{
			testName: "when session has no user and user is set should execute successfully",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(fetchUserResult, nil)
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(fetchSessionResult, nil)
					}),
				gu.Ptr("user-123"), nil,
			),
			shouldHaveResult: true,
		},
		{
			testName: "when session user matches fetched user should execute successfully",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(fetchUserResult, nil)
						session := *fetchSessionResult
						session.UserID = "user-123"
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(&session, nil)
					}),
				nil, gu.Ptr("different-user"),
			),
			shouldHaveResult: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "", "")
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.queryExecutor != nil {
				domain.WithQueryExecutor(tc.queryExecutor)(opts)
			}

			// Test
			err := tc.cmd.Execute(ctx, opts)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.shouldHaveResult, tc.cmd.Result() != nil)
		})
	}
}

func TestCheckUserCommand_Events(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	fetchUserResult := &domain.User{
		ID:             "user-123",
		OrganizationID: "org-1",
		InstanceID:     "instance-1",
		Username:       "username",
		State:          domain.UserStateActive,
	}
	fetchSessionResult := &domain.Session{
		ID:         "session-1",
		InstanceID: "instance-1",
	}

	type expectedEventFields struct {
		UserID            string
		UserResourceOwner string
		PreferredLanguage *language.Tag
		AggregateID       string
		ResourceOwner     string
	}
	tt := []struct {
		testName      string
		cmd           *domain.CheckUserCommand
		event         *expectedEventFields
		expectedError error
	}{
		{
			testName: "when user is machine should return event without language",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						user := *fetchUserResult
						user.Machine = &domain.MachineUser{}
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(&user, nil)
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(fetchSessionResult, nil).Times(2)
					}),
				gu.Ptr("user-123"), nil,
			),
			event: &expectedEventFields{
				UserID:            "user-123",
				UserResourceOwner: "org-1",
				AggregateID:       "session-1",
				ResourceOwner:     "instance-1",
			},
			expectedError: nil,
		},
		{
			testName: "when user is human without preferred language should return event without language",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						user := *fetchUserResult
						user.Machine = nil
						user.Human = &domain.HumanUser{}
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(&user, nil)
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(fetchSessionResult, nil).Times(2)
					}),
				gu.Ptr("user-123"), nil,
			),
			event: &expectedEventFields{
				UserID:            "user-123",
				UserResourceOwner: "org-1",
				AggregateID:       "session-1",
				ResourceOwner:     "instance-1",
			},
			expectedError: nil,
		},
		{
			testName: "when user is human with preferred language should return event with language",
			cmd: domain.NewCheckUserCommand(
				domainmock.InitCheckUserParent(ctrl).
					AddExpectation(func(recorder *domainmock.MockCheckUserParentMockRecorder) {
						user := *fetchUserResult
						user.Human = &domain.HumanUser{PreferredLanguage: language.Afrikaans}
						recorder.FetchUser(gomock.Any(), gomock.Any()).Return(&user, nil)
						recorder.FetchSession(gomock.Any(), gomock.Any()).Return(fetchSessionResult, nil).Times(2)
					}),
				gu.Ptr("user-123"), nil,
			),
			event: &expectedEventFields{
				UserID:            "user-123",
				UserResourceOwner: "org-1",
				AggregateID:       "session-1",
				ResourceOwner:     "instance-1",
				PreferredLanguage: &language.Afrikaans,
			},
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			require.NoError(t, opts.Invoke(t.Context(), tc.cmd))

			// Test
			events, err := tc.cmd.Events(t.Context(), opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			if tc.event == nil {
				assert.Empty(t, events)
				return
			}
			require.Len(t, events, 1)
			userCheckedEvent, ok := events[0].(*session.UserCheckedEvent)
			require.True(t, ok)

			assert.Equal(t, tc.event.UserID, userCheckedEvent.UserID)
			assert.Equal(t, tc.event.UserResourceOwner, userCheckedEvent.UserResourceOwner)
			assert.NotZero(t, userCheckedEvent.CheckedAt)
			assert.Equal(t, tc.event.PreferredLanguage, userCheckedEvent.PreferredLanguage)
			assert.Equal(t, tc.event.AggregateID, userCheckedEvent.Aggregate().ID)
			assert.Equal(t, tc.event.ResourceOwner, userCheckedEvent.Aggregate().ResourceOwner)
		})
	}
}
