package domain_test

import (
	"errors"
	"testing"
	"time"

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

func TestListSessionsQuery_Validate(t *testing.T) {
	t.Parallel()

	// Given
	q := domain.NewListSessionsQuery(nil)

	// Test
	err := q.Validate(t.Context(), nil)

	// Verify
	assert.NoError(t, err)

}

func TestListSessionsQuery_Sorting(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name                     string
		request                  *domain.ListSessionsRequest
		expectedSortingDirection database.OrderDirection
		expectedOrderBy          database.Columns
	}{
		{
			name: "sort by creation date DESC",
			request: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnCreationDate,
			},
			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("sessions", "created_at")},
		},
		{
			name: "sort by creation date ASC",
			request: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnCreationDate,
				Ascending:  true,
			},
			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          database.Columns{database.NewColumn("sessions", "created_at")},
		},
		{
			name: "unspecified sorting column — no ordering",
			request: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnUnspecified,
			},
			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			ctrl := gomock.NewController(t)
			mockRepo := domainmock.NewSessionRepo(ctrl)

			q := domain.NewListSessionsQuery(tc.request)
			opts := &database.QueryOpts{}

			// Test
			queryOpt := q.Sorting(mockRepo)
			queryOpt(opts)

			// Verify
			assert.Equal(t, tc.expectedOrderBy, opts.OrderBy)
			assert.Equal(t, tc.expectedSortingDirection, opts.Ordering)
		})
	}
}

func TestListSessionsQuery_Conditions(t *testing.T) {
	t.Parallel()

	const (
		userID  = "user-1"
		agentID = "agent-fp-1"
	)

	now := time.Now()

	defaultConds := func(repo *domainmock.SessionRepo) []database.Condition {
		return []database.Condition{
			repo.InstanceIDCondition("inst-1"),
			database.Or(
				repo.UserIDCondition(userID),
				repo.UserAgentIDCondition(agentID),
				repo.CreatorIDCondition(userID),
				database.PermissionCheck(domain.SessionReadPermission, false),
			),
		}
	}
	appendToDefaultConditions := func(repo *domainmock.SessionRepo, conds database.Condition) database.Condition {
		return database.And(
			append(defaultConds(repo), conds)...,
		)
	}

	tt := []struct {
		name          string
		request       *domain.ListSessionsRequest
		expectedCond  func(repo *domainmock.SessionRepo) database.Condition
		expectedError error
	}{
		{
			name:    "no filters returns permission and instance id conditions",
			request: &domain.ListSessionsRequest{},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(defaultConds(repo)...)
			},
		},
		{
			name: "SessionIDsFilter empty returns empty OR",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionIDsFilter{IDs: []string{}},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return appendToDefaultConditions(repo, database.Or())
			},
		},
		{
			name: "SessionIDsFilter with multiple IDs returns OR",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionIDsFilter{IDs: []string{"s-1", "s-2"}},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return appendToDefaultConditions(repo, database.Or(repo.IDCondition("s-1"), repo.IDCondition("s-2")))
			},
		},
		{
			name: "SessionUserIDFilter",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionUserIDFilter{UserID: "target-user"},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo, repo.UserIDCondition("target-user")),
				)
			},
		},
		{
			name: "SessionCreationDateFilter",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionCreationDateFilter{Op: database.NumberOperationGreaterThan, Date: now},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo, repo.CreatedAtCondition(database.NumberOperationGreaterThan, now)),
				)
			},
		},
		{
			name: "SessionCreatorFilter",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionCreatorFilter{ID: "other-creator"},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo, repo.CreatorIDCondition("other-creator")),
				)
			},
		},
		{
			name: "SessionUserAgentFilter",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionUserAgentFilter{FingerprintID: "explicit-fp"},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo, repo.UserAgentIDCondition("explicit-fp")),
				)
			},
		},
		{
			name: "SessionExpirationDateFilter with GREATER_THAN includes null expiration",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionExpirationDateFilter{Op: database.NumberOperationGreaterThan, Date: now},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo,
						database.Or(repo.ExpirationCondition(database.NumberOperationGreaterThan, now), database.IsNull(repo.ExpirationColumn())),
					),
				)
			},
		},
		{
			name: "SessionExpirationDateFilter with GREATER_OR_EQUALS includes null expiration",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionExpirationDateFilter{Op: database.NumberOperationGreaterThanOrEqual, Date: now},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo,
						database.Or(repo.ExpirationCondition(database.NumberOperationGreaterThanOrEqual, now), database.IsNull(repo.ExpirationColumn())),
					),
				)

			},
		},
		{
			name: "SessionExpirationDateFilter with LESS_THAN does not include null expiration",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionExpirationDateFilter{Op: database.NumberOperationLessThan, Date: now},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo, repo.ExpirationCondition(database.NumberOperationLessThan, now)),
				)
			},
		},
		{
			name: "SessionExpirationDateFilter with LESS_OR_EQUALS does not include null expiration",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionExpirationDateFilter{Op: database.NumberOperationLessThanOrEqual, Date: now},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo, repo.ExpirationCondition(database.NumberOperationLessThanOrEqual, now)),
				)
			},
		},
		{
			name: "multiple filters combined with AND",
			request: &domain.ListSessionsRequest{
				Filters: []domain.SessionFilter{
					domain.SessionUserIDFilter{UserID: "target-user"},
					domain.SessionCreatorFilter{ID: "creator-1"},
				},
			},
			expectedCond: func(repo *domainmock.SessionRepo) database.Condition {
				return database.And(
					appendToDefaultConditions(repo,
						database.And(repo.UserIDCondition("target-user"), repo.CreatorIDCondition("creator-1")),
					),
				)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := authz.NewMockContextWithAgent("inst-1", "org-1", userID, agentID)
			ctrl := gomock.NewController(t)
			sessionRepo := domainmock.NewSessionRepo(ctrl)

			q := domain.NewListSessionsQuery(tc.request)
			cond, err := q.Conditions(ctx, sessionRepo)

			assert.ErrorIs(t, err, tc.expectedError)
			expected := tc.expectedCond(sessionRepo)
			assert.Equal(t, expected.String(), cond.String())
		})
	}
}

func TestListSessionsQuery_Execute(t *testing.T) {
	t.Parallel()

	const (
		instanceID = "inst-1"
		orgID      = "org-1"
		userID     = "user-1"
		agentID    = "agent-fp-1"
	)

	listErr := errors.New("list mock error")

	tt := []struct {
		name             string
		request          *domain.ListSessionsRequest
		setupMock        func(sessionRepo *domainmock.SessionRepo)
		expectedSessions []*domain.Session
		expectedError    error
	}{
		{
			name:    "when List fails should return internal error",
			request: &domain.ListSessionsRequest{},
			setupMock: func(sessionRepo *domainmock.SessionRepo) {
				instanceCond := sessionRepo.InstanceIDCondition(instanceID)
				permCond := database.Or(
					sessionRepo.UserIDCondition(userID),
					sessionRepo.UserAgentIDCondition(agentID),
					sessionRepo.CreatorIDCondition(userID),
					database.PermissionCheck(domain.SessionReadPermission, false),
				)
				sessionRepo.EXPECT().
					LoadUserData().
					Times(1).
					Return(sessionRepo)
				sessionRepo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(database.And(instanceCond, permCond))),
						dbmock.QueryOptions(func(*database.QueryOpts) {}),
						dbmock.QueryOptions(database.WithLimit(0)),
						dbmock.QueryOptions(database.WithOffset(0)),
					).
					Times(1).
					Return(nil, listErr)
			},
			expectedError: zerrors.ThrowInternal(listErr, "DOM-Yx8q2r", "Errors.Session.List"),
		},
		{
			name: "when List succeeds should return sessions",
			request: &domain.ListSessionsRequest{
				SortColumn: domain.SessionSortColumnCreationDate,
				Ascending:  true,
				Limit:      10,
				Offset:     5,
			},
			setupMock: func(sessionRepo *domainmock.SessionRepo) {
				instanceCond := sessionRepo.InstanceIDCondition(instanceID)
				permCond := database.Or(
					sessionRepo.UserIDCondition(userID),
					sessionRepo.UserAgentIDCondition(agentID),
					sessionRepo.CreatorIDCondition(userID),
					database.PermissionCheck(domain.SessionReadPermission, false),
				)
				sessionRepo.EXPECT().
					LoadUserData().
					Times(1).
					Return(sessionRepo)
				sessionRepo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(database.And(instanceCond, permCond))),
						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionAsc, sessionRepo.CreatedAtColumn())),
						dbmock.QueryOptions(database.WithLimit(10)),
						dbmock.QueryOptions(database.WithOffset(5)),
					).
					Times(1).
					Return([]*domain.Session{{ID: "session-1", UserID: userID}}, nil)
			},
			expectedSessions: []*domain.Session{{ID: "session-1", UserID: userID}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := authz.NewMockContextWithAgent(instanceID, orgID, userID, agentID)
			ctrl := gomock.NewController(t)
			sessionRepo := domainmock.NewSessionRepo(ctrl)

			if tc.setupMock != nil {
				tc.setupMock(sessionRepo)
			}

			opts := domain.DefaultOpts(nil)
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			domain.WithSessionRepo(sessionRepo)(opts)

			q := domain.NewListSessionsQuery(tc.request)
			err := q.Execute(ctx, opts)

			assert.ErrorIs(t, err, tc.expectedError)
			assert.ElementsMatch(t, tc.expectedSessions, q.Result())
		})
	}
}
