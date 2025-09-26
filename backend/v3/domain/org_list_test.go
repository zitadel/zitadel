package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func TestListOrgsCommand_sorting(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name    string
		request *org.ListOrganizationsRequest

		orgRepo func(*gomock.Controller) domain.OrganizationRepository

		expectedSortingDirection database.OrderDirection
		expectedOrderBy          database.Columns
	}{
		{
			name: "sorting by name desc",
			request: &org.ListOrganizationsRequest{
				SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME,
				Query: &object.ListQuery{
					Asc: false,
				},
			},

			orgRepo: func(c *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(c)

				repo.EXPECT().
					NameColumn().
					Times(1).
					Return(database.NewColumn("organizations", "name"))

				return repo
			},

			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("organizations", "name")},
		},
		{
			name: "sorting by name asc",
			request: &org.ListOrganizationsRequest{
				SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME,
				Query: &object.ListQuery{
					Asc: true,
				},
			},

			orgRepo: func(c *gomock.Controller) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(c)

				repo.EXPECT().
					NameColumn().
					Times(1).
					Return(database.NewColumn("organizations", "name"))

				return repo
			},

			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          database.Columns{database.NewColumn("organizations", "name")},
		},
		{
			name: "unspecified sorting",
			request: &org.ListOrganizationsRequest{
				SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_UNSPECIFIED,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			ctrl := gomock.NewController(t)
			orgRepo := domainmock.NewOrgRepo(ctrl)
			l := &domain.ListOrgsCommand{
				Request: tc.request,
			}
			opts := &database.QueryOpts{
				OrderBy: tc.expectedOrderBy,
			}

			// Test
			gotFunc := l.Sorting(orgRepo)
			gotFunc(opts)

			// Verify
			assert.Equal(t, tc.expectedOrderBy, opts.OrderBy)
			assert.Equal(t, tc.expectedSortingDirection, opts.Ordering)
		})
	}
}

func TestListOrgsCommand_pagination(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name       string
		request    *org.ListOrganizationsRequest
		wantLimit  uint64
		wantOffset uint32
	}{
		{
			name: "pagination with limit and offset",
			request: &org.ListOrganizationsRequest{
				Query: &object.ListQuery{
					Limit:  10,
					Offset: 5,
				},
			},
			wantLimit:  10,
			wantOffset: 5,
		},
		{
			name: "pagination with zero values",
			request: &org.ListOrganizationsRequest{
				Query: &object.ListQuery{
					Limit:  0,
					Offset: 0,
				},
			},
			wantLimit:  0,
			wantOffset: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			l := &domain.ListOrgsCommand{
				Request: tc.request,
			}
			opts := &database.QueryOpts{}

			// Test
			limitFunc, offsetFunc := l.Pagination()
			limitFunc(opts)
			offsetFunc(opts)

			// Verify
			assert.EqualValues(t, tc.wantLimit, opts.Limit)
			assert.EqualValues(t, tc.wantOffset, opts.Offset)
		})
	}
}

func TestListOrgsCommand_Execute(t *testing.T) {
	t.Parallel()
	txInitErr := errors.New("tx init error")
	listErr := errors.New("list mock error")

	tt := []struct {
		testName string

		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		orgRepo       func(ctrl *gomock.Controller, queryParams ...database.QueryOption) func(client database.QueryExecutor) domain.OrganizationRepository

		queryParams  []database.QueryOption
		inputRequest *org.ListOrganizationsRequest

		expectedOrganizations []*domain.Organization
		expectedError         error
	}{
		{
			testName: "when EnsureTx fails should return error",
			queryExecutor: func(ctrl *gomock.Controller) database.QueryExecutor {
				mockDB := dbmock.NewMockPool(ctrl)
				mockDB.EXPECT().
					Begin(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, txInitErr)
				return mockDB
			},
			expectedError: txInitErr,
		},
		{
			testName: "when condition parsing fails should return error",
			orgRepo: func(ctrl *gomock.Controller, _ ...database.QueryOption) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)

				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputRequest: &org.ListOrganizationsRequest{
				Queries: []*org.SearchQuery{{Query: &org.SearchQuery_DomainQuery{DomainQuery: &org.OrganizationDomainQuery{
					Domain: "some domain",
					Method: object.TextQueryMethod(99)}},
				},
				},
			},
			expectedError: domain.NewUnexpectedTextQueryOperationError("DOM-iBRBVe", object.TextQueryMethod(99)),
		},
		{
			testName: "when listing orgs fails should return error",
			orgRepo: func(ctrl *gomock.Controller, queryParams ...database.QueryOption) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				domainRepo := domainmock.NewOrgDomainRepo(ctrl)
				repo.SetDomains(domainRepo)
				repo.SetExistsDomain(database.Exists("domains", domainRepo.DomainCondition(database.TextOperationEqual, "some domain")))

				repo.EXPECT().
					List(
						gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								database.And(
									repo.InstanceIDCondition("instance-1"),
									repo.ExistsDomain(domainRepo.DomainCondition(database.TextOperationEqual, "some domain")),
								),
							),
						),
						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionAsc, repo.NameColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Times(1).
					Return(nil, listErr)

				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputRequest: &org.ListOrganizationsRequest{
				Queries: []*org.SearchQuery{
					{
						Query: &org.SearchQuery_DomainQuery{
							DomainQuery: &org.OrganizationDomainQuery{
								Domain: "some domain",
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
							},
						},
					},
				},
				SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME,
				Query: &object.ListQuery{
					Asc:    true,
					Offset: 1,
					Limit:  2,
				},
			},

			expectedError: listErr,
		},
		{
			testName: "when listing orgs succeeds should return expected organizations",
			orgRepo: func(ctrl *gomock.Controller, queryParams ...database.QueryOption) func(client database.QueryExecutor) domain.OrganizationRepository {
				repo := domainmock.NewOrgRepo(ctrl)
				domainRepo := domainmock.NewOrgDomainRepo(ctrl)
				repo.SetDomains(domainRepo)

				repo.EXPECT().
					List(
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-1"))),
						dbmock.QueryOptions(database.WithCondition(repo.IDCondition("org-2"))),
						dbmock.QueryOptions(database.WithCondition(repo.NameCondition(database.TextOperationEqual, "Named Org"))),
						dbmock.QueryOptions(database.WithCondition(repo.StateCondition(domain.OrgStateActive))),

						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionDesc, repo.NameColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Times(1).
					Return([]*domain.Organization{
						{ID: "org-1"},
						{ID: "org-2"},
					}, nil)

				return func(_ database.QueryExecutor) domain.OrganizationRepository { return repo }
			},
			inputRequest: &org.ListOrganizationsRequest{
				Queries: []*org.SearchQuery{
					{
						Query: &org.SearchQuery_DefaultQuery{
							DefaultQuery: &org.DefaultOrganizationQuery{},
						},
					},
					{
						Query: &org.SearchQuery_IdQuery{
							IdQuery: &org.OrganizationIDQuery{
								Id: "org-2",
							},
						},
					},
					{
						Query: &org.SearchQuery_NameQuery{
							NameQuery: &org.OrganizationNameQuery{
								Name:   "Named Org",
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
							},
						},
					},
					{
						Query: &org.SearchQuery_StateQuery{
							StateQuery: &org.OrganizationStateQuery{
								State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
							},
						},
					},
				},
				SortingColumn: org.OrganizationFieldName_ORGANIZATION_FIELD_NAME_NAME,
				Query: &object.ListQuery{
					Asc:    false,
					Offset: 1,
					Limit:  2,
				},
			},

			expectedOrganizations: []*domain.Organization{
				{ID: "org-1"},
				{ID: "org-2"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "org-1", "")
			ctrl := gomock.NewController(t)
			cmd := &domain.ListOrgsCommand{
				Request: tc.inputRequest,
			}
			opts := &domain.CommandOpts{
				DB: new(noopdb.Pool),
			}
			if tc.orgRepo != nil {
				opts.SetOrgRepo(tc.orgRepo(ctrl, tc.queryParams...))
			}
			if tc.queryExecutor != nil {
				opts.DB = tc.queryExecutor(ctrl)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.ElementsMatch(t, tc.expectedOrganizations, cmd.Result)
		})
	}
}

func TestListOrgsCommand_ResultToGRPC(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)

	tt := []struct {
		name string
		orgs []*domain.Organization
		want []*org.Organization
	}{
		{
			name: "empty result",
			orgs: nil,
			want: []*org.Organization{},
		},
		{
			name: "multiple organizations",
			orgs: []*domain.Organization{
				{
					ID:        "org-1",
					Name:      "org 1",
					State:     domain.OrgStateActive,
					CreatedAt: yesterday,
					UpdatedAt: now,
					Domains: []*domain.OrganizationDomain{
						{Domain: "wrong selected domain"},
						{Domain: "domain.example.com", IsPrimary: true},
					},
				},
				{
					ID:        "org-2",
					Name:      "org 2",
					State:     domain.OrgStateInactive,
					CreatedAt: yesterday,
					UpdatedAt: now,
					Domains: []*domain.OrganizationDomain{
						{Domain: "wrong selected domain 2"},
						{Domain: "domain2.example.com", IsPrimary: true},
					},
				},
			},
			want: []*org.Organization{
				{
					Id:    "org-1",
					Name:  "org 1",
					State: org.OrganizationState_ORGANIZATION_STATE_ACTIVE,
					Details: &object.Details{
						ChangeDate:   timestamppb.New(now),
						CreationDate: timestamppb.New(yesterday),
					},
					PrimaryDomain: "domain.example.com",
				},
				{
					Id:    "org-2",
					Name:  "org 2",
					State: org.OrganizationState_ORGANIZATION_STATE_INACTIVE,
					Details: &object.Details{
						ChangeDate:   timestamppb.New(now),
						CreationDate: timestamppb.New(yesterday),
					},
					PrimaryDomain: "domain2.example.com",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cmd := &domain.ListOrgsCommand{
				Result: tc.orgs,
			}

			got := cmd.ResultToGRPC()
			assert.Equal(t, tc.want, got)
		})
	}
}
