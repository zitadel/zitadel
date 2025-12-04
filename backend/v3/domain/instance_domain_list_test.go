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
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestListInstanceDomainsQuery_Validate(t *testing.T) {
	t.Parallel()
	permissionErr := errors.New("permission error")

	tt := []struct {
		name              string
		request           *instance.ListCustomDomainsRequest
		permissionChecker func(ctrl *gomock.Controller) domain.PermissionChecker
		expectedError     error
	}{
		{
			name:    "when missing domain read permission should return permission denied (empty instance id)",
			request: &instance.ListCustomDomainsRequest{InstanceId: " "},
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				checker := domainmock.NewMockPermissionChecker(ctrl)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainReadPermission).
					Times(1).
					Return(permissionErr)
				return checker
			},
			expectedError: zerrors.ThrowPermissionDenied(permissionErr, "DOM-RyCEyr", "permission denied"),
		},
		{
			name:    "when valid permissions should return no error (ctx instance same as input)",
			request: &instance.ListCustomDomainsRequest{InstanceId: "instance-1"},
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				checker := domainmock.NewMockPermissionChecker(ctrl)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainReadPermission).
					Times(1).
					Return(nil)
				return checker
			},
		},
		{
			name: "when input instance doesn't match context should check instance and return error",
			request: &instance.ListCustomDomainsRequest{
				InstanceId: "different-instance",
			},
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				checker := domainmock.NewMockPermissionChecker(ctrl)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(permissionErr)
				return checker
			},
			expectedError: zerrors.ThrowPermissionDenied(permissionErr, "DOM-yN7oCp", "permission denied"),
		},
		{
			name: "when input instance doesn't match context should check instance and succeed then check domain and fail",
			request: &instance.ListCustomDomainsRequest{
				InstanceId: "different-instance",
			},
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				checker := domainmock.NewMockPermissionChecker(ctrl)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(nil)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainReadPermission).
					Times(1).
					Return(permissionErr)
				return checker
			},
			expectedError: zerrors.ThrowPermissionDenied(permissionErr, "DOM-RyCEyr", "permission denied"),
		},
		{
			name: "when input instance doesn't match context should check instance and succeed then check domain and succeed",
			request: &instance.ListCustomDomainsRequest{
				InstanceId: "different-instance",
			},
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				checker := domainmock.NewMockPermissionChecker(ctrl)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(nil)
				checker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.DomainReadPermission).
					Times(1).
					Return(nil)
				return checker
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ctx := authz.NewMockContext("instance-1", "org-1", "")

			query := domain.NewListInstanceDomainsQuery(tc.request)
			err := query.Validate(ctx, &domain.InvokeOpts{
				Permissions: tc.permissionChecker(ctrl),
			})

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListInstanceDomainsQuery_Execute(t *testing.T) {
	t.Parallel()
	listErr := errors.New("list mock error")

	tt := []struct {
		name            string
		request         *instance.ListCustomDomainsRequest
		repo            func(ctrl *gomock.Controller) domain.InstanceDomainRepository
		expectedError   error
		expectedDomains []*domain.InstanceDomain
	}{
		{
			name: "when parsing conditions fails should return invalid argument error",
			request: &instance.ListCustomDomainsRequest{
				Pagination: &filter.PaginationRequest{Limit: 2, Offset: 1, Asc: false},
				Filters: []*instance.CustomDomainFilter{
					{Filter: &instance.CustomDomainFilter_DomainFilter{
						DomainFilter: &instance.DomainFilter{
							Domain: "test.domain",
							Method: 99,
						},
					}},
				},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
			},
			repo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				return repo
			},
			expectedError: zerrors.ThrowInvalidArgument(nil, "OBJ-iBRBVe", "invalid text query method"),
		},
		{
			name: "when listing domains fails should return error",
			request: &instance.ListCustomDomainsRequest{
				Pagination: &filter.PaginationRequest{Limit: 2, Offset: 1, Asc: false},
				Filters: []*instance.CustomDomainFilter{
					{Filter: &instance.CustomDomainFilter_DomainFilter{
						DomainFilter: &instance.DomainFilter{
							Domain: "test.domain",
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					}},
				},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
			},
			repo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(database.And(
								repo.DomainCondition(database.TextOperationContains, "test.domain"),
								repo.TypeCondition(domain.DomainTypeCustom),
							)),
						),
						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionDesc, repo.DomainColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Return(nil, listErr)
				return repo
			},
			expectedError: zerrors.ThrowInternal(listErr, "DOM-ubaPNU", "failed fetching instance domains"),
		},
		{
			name: "when listing domains succeeds should return domains",
			request: &instance.ListCustomDomainsRequest{
				Filters: []*instance.CustomDomainFilter{
					{Filter: &instance.CustomDomainFilter_DomainFilter{
						DomainFilter: &instance.DomainFilter{
							Domain: "test.domain",
							Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					}},
					{Filter: &instance.CustomDomainFilter_GeneratedFilter{
						GeneratedFilter: true,
					}},
					{Filter: &instance.CustomDomainFilter_PrimaryFilter{
						PrimaryFilter: true,
					}},
				},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
			},
			repo: func(ctrl *gomock.Controller) domain.InstanceDomainRepository {
				repo := domainmock.NewInstancesDomainRepo(ctrl)
				domains := []*domain.InstanceDomain{
					{Domain: "test1.domain"},
					{Domain: "test2.domain"},
				}
				repo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(database.And(
								repo.DomainCondition(database.TextOperationContains, "test.domain"),
								repo.IsGeneratedCondition(true),
								repo.IsPrimaryCondition(true),
								repo.TypeCondition(domain.DomainTypeCustom),
							)),
						),
						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionDesc, repo.DomainColumn())),
						dbmock.QueryOptions(database.WithLimit(0)),
						dbmock.QueryOptions(database.WithOffset(0)),
					).
					Return(domains, nil)
				return repo
			},
			expectedDomains: []*domain.InstanceDomain{
				{Domain: "test1.domain"},
				{Domain: "test2.domain"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.repo != nil {
				instDomain := tc.repo(ctrl)
				domain.WithInstanceDomainRepo(instDomain)(opts)
			}
			query := domain.NewListInstanceDomainsQuery(tc.request)
			err := query.Execute(t.Context(), opts)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDomains, query.Result())
			}
		})
	}
}

func TestListInstanceDomainsQuery_Sorting(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name    string
		request *instance.ListCustomDomainsRequest

		expectedSortingDirection database.OrderDirection
		expectedOrderBy          database.Columns
	}{
		{
			name: "sort by creation date DESC",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
			},
			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("instance_domains", "created_at")},
		},
		{
			name: "sort by creation date ASC",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Pagination:    &filter.PaginationRequest{Asc: true},
			},
			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          database.Columns{database.NewColumn("instance_domains", "created_at")},
		},
		{
			name: "sort by domain DESC",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
			},
			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("instance_domains", "domain")},
		},
		{
			name: "sort by domain ASC",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
				Pagination:    &filter.PaginationRequest{Asc: true},
			},
			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          database.Columns{database.NewColumn("instance_domains", "domain")},
		},
		{
			name: "sort by generated DESC",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED,
			},
			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("instance_domains", "is_generated")},
		},
		{
			name: "sort by primary DESC",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY,
			},
			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("instance_domains", "is_primary")},
		},
		{
			name: "unspecified field",
			request: &instance.ListCustomDomainsRequest{
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED,
			},
			expectedSortingDirection: database.OrderDirectionAsc,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			mockCtrl := gomock.NewController(t)
			mockRepo := domainmock.NewInstancesDomainRepo(mockCtrl)

			cmd := domain.NewListInstanceDomainsQuery(tc.request)
			opts := &database.QueryOpts{}

			// Test
			queryOpt := cmd.Sorting(mockRepo)
			queryOpt(opts)

			// Verify
			assert.Equal(t, tc.expectedOrderBy, opts.OrderBy)
			assert.Equal(t, tc.expectedSortingDirection, opts.Ordering)
		})
	}
}
