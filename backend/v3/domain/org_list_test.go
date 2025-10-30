package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
)

func TestListOrgsCommand_sorting(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name    string
		request *org.ListOrganizationsRequest

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
			l := &domain.ListOrgsQuery{
				Request: tc.request,
			}
			opts := &database.QueryOpts{}

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
			l := &domain.ListOrgsQuery{
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
	listErr := errors.New("list mock error")

	tt := []struct {
		testName string

		repos func(ctrl *gomock.Controller, queryParams ...database.QueryOption) (domain.OrganizationRepository, domain.OrganizationDomainRepository)

		queryParams  []database.QueryOption
		inputRequest *org.ListOrganizationsRequest

		expectedOrganizations []*domain.Organization
		expectedError         error
	}{
		{
			testName: "when condition parsing fails should return error",
			repos: func(ctrl *gomock.Controller, _ ...database.QueryOption) (domain.OrganizationRepository, domain.OrganizationDomainRepository) {
				orgRepo := domainmock.NewOrgRepo(ctrl)
				domainRepo := domainmock.NewOrgDomainRepo(ctrl)
				orgRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(orgRepo)

				return orgRepo, domainRepo
			},
			inputRequest: &org.ListOrganizationsRequest{
				Queries: []*org.SearchQuery{{Query: &org.SearchQuery_DomainQuery{DomainQuery: &org.OrganizationDomainQuery{
					Domain: "some domain",
					Method: object.TextQueryMethod(99)}},
				},
				},
			},
			expectedError: zerrors.ThrowInvalidArgument(nil, "OBJ-iBRBVe", "invalid text query method"),
		},
		{
			testName: "when listing orgs fails should return error",
			repos: func(ctrl *gomock.Controller, queryParams ...database.QueryOption) (domain.OrganizationRepository, domain.OrganizationDomainRepository) {
				orgRepo := domainmock.NewOrgRepo(ctrl)
				domainRepo := domainmock.NewOrgDomainRepo(ctrl)
				orgRepo.SetExistsDomain(database.Exists("domains", domainRepo.DomainCondition(database.TextOperationEqual, "some domain")))

				orgRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(orgRepo)

				orgRepo.EXPECT().
					List(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								database.And(
									database.And(
										orgRepo.InstanceIDCondition("instance-1"),
										orgRepo.ExistsDomain(domainRepo.DomainCondition(database.TextOperationEqual, "some domain")),
									),
									orgRepo.InstanceIDCondition("instance-1"),
								),
							),
						),

						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionAsc, orgRepo.NameColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Times(1).
					Return(nil, listErr)

				return orgRepo, domainRepo
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
			repos: func(ctrl *gomock.Controller, queryParams ...database.QueryOption) (domain.OrganizationRepository, domain.OrganizationDomainRepository) {
				orgRepo := domainmock.NewOrgRepo(ctrl)
				domainRepo := domainmock.NewOrgDomainRepo(ctrl)

				orgRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(orgRepo)

				orgRepo.EXPECT().
					List(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(database.And(
							orgRepo.IDCondition("org-1"),
							orgRepo.IDCondition("org-2"),
							orgRepo.NameCondition(database.TextOperationEqual, "Named Org"),
							orgRepo.StateCondition(domain.OrgStateActive),
							orgRepo.InstanceIDCondition("instance-1"),
						))),

						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionDesc, orgRepo.NameColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Times(1).
					Return([]*domain.Organization{
						{ID: "org-1"},
						{ID: "org-2"},
					}, nil)

				return orgRepo, domainRepo
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
			cmd := &domain.ListOrgsQuery{
				Request: tc.inputRequest,
			}
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.repos != nil {
				orgRepo, domainRepo := tc.repos(ctrl, tc.queryParams...)
				domain.WithOrganizationRepo(orgRepo)(opts)
				domain.WithOrganizationDomainRepo(domainRepo)(opts)
			}

			// Test
			err := opts.Invoke(ctx, cmd)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.ElementsMatch(t, tc.expectedOrganizations, cmd.Result())
		})
	}
}
