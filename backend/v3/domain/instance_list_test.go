package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func TestListInstancesCommand_ResultToGRPC(t *testing.T) {
	t.Parallel()
	now := time.Now()

	tt := []struct {
		testName       string
		inputResult    []*domain.Instance
		expectedResult []*instance.Instance
	}{
		{
			testName:       "empty result",
			inputResult:    []*domain.Instance{},
			expectedResult: []*instance.Instance{},
		},
		{
			testName: "single instance without domains",
			inputResult: []*domain.Instance{
				{
					ID:        "instance1",
					Name:      "test-instance",
					CreatedAt: now,
					UpdatedAt: now,
					Domains:   nil,
				},
			},
			expectedResult: []*instance.Instance{
				{
					Id:           "instance1",
					Name:         "test-instance",
					CreationDate: timestamppb.New(now),
					ChangeDate:   timestamppb.New(now),
					State:        instance.State_STATE_RUNNING,
					Domains:      []*instance.Domain{},
				},
			},
		},
		{
			testName: "multiple instances with domains",
			inputResult: []*domain.Instance{
				{
					ID:        "instance1",
					Name:      "test-instance-1",
					CreatedAt: now,
					UpdatedAt: now,
					Domains: []*domain.InstanceDomain{
						{
							InstanceID:  "instance1",
							Domain:      "domain1.com",
							CreatedAt:   now,
							IsPrimary:   gu.Ptr(true),
							IsGenerated: gu.Ptr(false),
						},
					},
				},
				{
					ID:        "instance2",
					Name:      "test-instance-2",
					CreatedAt: now,
					UpdatedAt: now,
					Domains:   nil,
				},
			},
			expectedResult: []*instance.Instance{
				{
					Id:           "instance1",
					Name:         "test-instance-1",
					CreationDate: timestamppb.New(now),
					ChangeDate:   timestamppb.New(now),
					State:        instance.State_STATE_RUNNING,
					Domains: []*instance.Domain{
						{
							InstanceId:   "instance1",
							Domain:       "domain1.com",
							CreationDate: timestamppb.New(now),
							Primary:      true,
							Generated:    false,
						},
					},
				},
				{
					Id:           "instance2",
					Name:         "test-instance-2",
					CreationDate: timestamppb.New(now),
					ChangeDate:   timestamppb.New(now),
					State:        instance.State_STATE_RUNNING,
					Domains:      []*instance.Domain{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.ListInstancesCommand{
				Result: tc.inputResult,
			}

			// Test
			result := cmd.ResultToGRPC()

			// Verify
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

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
	txInitErr := errors.New("tx init error")
	listErr := errors.New("list mock error")

	tt := []struct {
		testName string

		queryExecutor func(ctrl *gomock.Controller) database.QueryExecutor
		repos         func(ctrl *gomock.Controller, queryParams ...database.QueryOption) (domain.InstanceRepository, domain.InstanceDomainRepository)

		queryParams  []database.QueryOption
		inputRequest *instance.ListInstancesRequest

		expectedInstances []*domain.Instance
		expectedError     error
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
			opts := &domain.CommandOpts{
				DB: new(noopdb.Pool),
			}
			if tc.repos != nil {
				inst, domain := tc.repos(ctrl, tc.queryParams...)
				opts.SetInstanceRepo(inst)
				opts.SetInstanceDomainRepo(domain)
			}
			if tc.queryExecutor != nil {
				opts.DB = tc.queryExecutor(ctrl)
			}

			// Test
			err := cmd.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.ElementsMatch(t, tc.expectedInstances, cmd.Result)

		})
	}
}
