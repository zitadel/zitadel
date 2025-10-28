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
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func TestListInstancesCommand_Sorting(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name    string
		request *instance.ListInstancesRequest

		expectedSortingDirection database.OrderDirection
		expectedOrderBy          database.Columns
	}{
		{
			name: "sort by creation date DESC",
			request: &instance.ListInstancesRequest{
				SortingColumn: instance.FieldName_FIELD_NAME_CREATION_DATE.Enum(),
			},

			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("instances", "created_at")},
		},
		{
			name: "sort by creation date ASC",
			request: &instance.ListInstancesRequest{
				SortingColumn: instance.FieldName_FIELD_NAME_CREATION_DATE.Enum(),
				Pagination:    &filter.PaginationRequest{Asc: true},
			},

			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          database.Columns{database.NewColumn("instances", "created_at")},
		},
		{
			name: "sort by id DESC",
			request: &instance.ListInstancesRequest{
				SortingColumn: instance.FieldName_FIELD_NAME_ID.Enum(),
			},
			expectedSortingDirection: database.OrderDirectionDesc,
			expectedOrderBy:          database.Columns{database.NewColumn("instances", "id")},
		},
		{
			name: "sort by id ASC",
			request: &instance.ListInstancesRequest{
				SortingColumn: instance.FieldName_FIELD_NAME_ID.Enum(),
				Pagination:    &filter.PaginationRequest{Asc: true},
			},
			expectedSortingDirection: database.OrderDirectionAsc,
			expectedOrderBy:          database.Columns{database.NewColumn("instances", "id")},
		},
		{
			name: "unspecified field",
			request: &instance.ListInstancesRequest{
				SortingColumn: instance.FieldName_FIELD_NAME_UNSPECIFIED.Enum(),
			},
			expectedSortingDirection: database.OrderDirectionAsc,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			mockCtrl := gomock.NewController(t)
			mockRepo := domainmock.NewInstanceRepo(mockCtrl)

			cmd := domain.NewListInstancesCommand(tc.request)
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

func TestListInstancesCommand_Execute(t *testing.T) {
	t.Parallel()
	listErr := errors.New("list mock error")

	tt := []struct {
		testName string

		repos func(ctrl *gomock.Controller, queryParams ...database.QueryOption) (domain.InstanceRepository, domain.InstanceDomainRepository)

		queryParams  []database.QueryOption
		inputRequest *instance.ListInstancesRequest

		expectedInstances []*domain.Instance
		expectedError     error
	}{
		{
			testName: "when listing instances fails should return error",
			repos: func(ctrl *gomock.Controller, _ ...database.QueryOption) (domain.InstanceRepository, domain.InstanceDomainRepository) {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)
				domainRepo := domainmock.NewInstancesDomainRepo(ctrl)
				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				domainExistsORedConditions := instanceRepo.ExistsDomain(
					database.Or([]database.Condition{
						domainRepo.DomainCondition(database.TextOperationEqual, "domain1.example.com"),
						domainRepo.DomainCondition(database.TextOperationEqual, "domain2.example.com"),
					}...),
				)

				idConditions := database.Or(
					[]database.Condition{
						instanceRepo.IDCondition("instance-1"),
						instanceRepo.IDCondition("instance-2"),
					}...,
				)

				instanceRepo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(database.And(
								idConditions, domainExistsORedConditions,
							)),
						),
						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionDesc, instanceRepo.NameColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Times(1).
					Return(nil, listErr)

				return instanceRepo, domainRepo
			},
			inputRequest: &instance.ListInstancesRequest{
				Pagination:    &filter.PaginationRequest{Limit: 2, Offset: 1, Asc: false},
				SortingColumn: instance.FieldName_FIELD_NAME_NAME.Enum(),
				Queries: []*instance.Query{
					{Query: &instance.Query_IdQuery{IdQuery: &instance.IdsQuery{
						Ids: []string{"instance-1", "instance-2"},
					}}},
					{Query: &instance.Query_DomainQuery{DomainQuery: &instance.DomainsQuery{
						Domains: []string{"domain1.example.com", "domain2.example.com"},
					}}},
				},
			},
			expectedError: listErr,
		},
		{
			testName: "when listing instances succeeds should save into result and return nil",
			repos: func(ctrl *gomock.Controller, _ ...database.QueryOption) (domain.InstanceRepository, domain.InstanceDomainRepository) {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)
				domainRepo := domainmock.NewInstancesDomainRepo(ctrl)
				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				domainExistsORedConditions := instanceRepo.ExistsDomain(
					database.Or([]database.Condition{
						domainRepo.DomainCondition(database.TextOperationEqual, "domain1.example.com"),
						domainRepo.DomainCondition(database.TextOperationEqual, "domain2.example.com"),
					}...),
				)

				idConditions := database.Or(
					[]database.Condition{
						instanceRepo.IDCondition("instance-1"),
						instanceRepo.IDCondition("instance-2"),
					}...,
				)

				instanceRepo.EXPECT().
					List(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(database.And(
								idConditions, domainExistsORedConditions,
							)),
						),
						dbmock.QueryOptions(database.WithOrderBy(database.OrderDirectionDesc, instanceRepo.NameColumn())),
						dbmock.QueryOptions(database.WithLimit(2)),
						dbmock.QueryOptions(database.WithOffset(1)),
					).
					Times(1).
					Return([]*domain.Instance{
						{
							ID:      "instance-2",
							Name:    "Instance Two",
							Domains: []*domain.InstanceDomain{{Domain: "domain1.example.com"}},
						},
						{
							ID:      "instance-1",
							Name:    "Instance One",
							Domains: []*domain.InstanceDomain{{Domain: "domain2.example.com"}},
						},
					}, nil)

				return instanceRepo, domainRepo
			},
			inputRequest: &instance.ListInstancesRequest{
				Pagination:    &filter.PaginationRequest{Limit: 2, Offset: 1, Asc: false},
				SortingColumn: instance.FieldName_FIELD_NAME_NAME.Enum(),
				Queries: []*instance.Query{
					{Query: &instance.Query_IdQuery{IdQuery: &instance.IdsQuery{
						Ids: []string{"instance-1", "instance-2"},
					}}},
					{Query: &instance.Query_DomainQuery{DomainQuery: &instance.DomainsQuery{
						Domains: []string{"domain1.example.com", "domain2.example.com"},
					}}},
				},
			},
			expectedInstances: []*domain.Instance{
				{
					ID:      "instance-2",
					Name:    "Instance Two",
					Domains: []*domain.InstanceDomain{{Domain: "domain1.example.com"}},
				},
				{
					ID:      "instance-1",
					Name:    "Instance One",
					Domains: []*domain.InstanceDomain{{Domain: "domain2.example.com"}},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctx := authz.NewMockContext("instance-1", "org-1", "")
			ctrl := gomock.NewController(t)
			cmd := domain.NewListInstancesCommand(tc.inputRequest)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)
			if tc.repos != nil {
				inst, instDomain := tc.repos(ctrl, tc.queryParams...)
				domain.WithInstanceRepo(inst)(opts)
				domain.WithInstanceDomainRepo(instDomain)(opts)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.ElementsMatch(t, tc.expectedInstances, cmd.Result())

		})
	}
}

func TestListInstancesCommand_Validate(t *testing.T) {
	t.Parallel()
	permissionErr := errors.New("permission error")
	tt := []struct {
		name              string
		permissionChecker func(ctrl *gomock.Controller) domain.PermissionChecker
		expectedError     error
	}{
		{
			name: "when user is missing permission should return permission denied",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(permissionErr)

				return permChecker
			},
			expectedError: zerrors.ThrowPermissionDenied(permissionErr, "DOM-cuT6Ws", "permission denied"),
		},
		{
			name: "when valid permission should return no error",
			permissionChecker: func(ctrl *gomock.Controller) domain.PermissionChecker {
				permChecker := domainmock.NewMockPermissionChecker(ctrl)

				permChecker.EXPECT().
					CheckInstancePermission(gomock.Any(), domain.InstanceReadPermission).
					Times(1).
					Return(nil)

				return permChecker
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Given
			l := &domain.ListInstancesQuery{}

			cmdOpts := &domain.InvokeOpts{}
			if tc.permissionChecker != nil {
				ctrl := gomock.NewController(t)
				cmdOpts.Permissions = tc.permissionChecker(ctrl)
			}

			// Test
			err := l.Validate(context.Background(), cmdOpts)

			// Verify
			assert.Equal(t, tc.expectedError, err)

		})
	}
}
