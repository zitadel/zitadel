package convert

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	filter_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	filter_v2beta "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func TestDomainInstanceModelToGRPCBetaResponse(t *testing.T) {
	// Given
	t.Parallel()

	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)

	inputInstance := &domain.Instance{
		ID:        "instance-1",
		Name:      "Instance One",
		CreatedAt: yesterday,
		UpdatedAt: now,
		Domains: []*domain.InstanceDomain{
			{
				InstanceID: "instance-1",
				Domain:     "d1.example.com",
				IsPrimary:  gu.Ptr(true),
				CreatedAt:  yesterday,
			},
			{
				InstanceID:  "instance-1",
				Domain:      "d2.example.com",
				IsGenerated: gu.Ptr(true),
				CreatedAt:   yesterday,
			},
			{
				InstanceID:  "instance-1",
				Domain:      "d3.example.com",
				IsPrimary:   gu.Ptr(true),
				IsGenerated: gu.Ptr(false),
				CreatedAt:   yesterday,
			},
			{
				InstanceID:  "instance-1",
				Domain:      "d4.example.com",
				IsPrimary:   gu.Ptr(false),
				IsGenerated: gu.Ptr(true),
				CreatedAt:   yesterday,
			},
		},
	}

	expectedInstance := &instance_v2beta.Instance{
		Id:           "instance-1",
		ChangeDate:   timestamppb.New(now),
		CreationDate: timestamppb.New(yesterday),
		Name:         "Instance One",
		Version:      "",
		Domains: []*instance_v2beta.Domain{
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d1.example.com",
				Primary:      true,
				Generated:    false,
			},
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d2.example.com",
				Primary:      false,
				Generated:    true,
			},
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d3.example.com",
				Primary:      true,
				Generated:    false,
			},
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d4.example.com",
				Primary:      false,
				Generated:    true,
			},
		},
	}

	// Test
	res := DomainInstanceModelToGRPCBetaResponse(inputInstance)

	// Verify
	assert.Equal(t, expectedInstance, res)
}

func TestDomainInstanceModelToGRPCResponse(t *testing.T) {
	// Given
	t.Parallel()

	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)

	inputInstance := &domain.Instance{
		ID:        "instance-1",
		Name:      "Instance One",
		CreatedAt: yesterday,
		UpdatedAt: now,
		Domains: []*domain.InstanceDomain{
			{
				InstanceID: "instance-1",
				Domain:     "d1.example.com",
				IsPrimary:  gu.Ptr(true),
				CreatedAt:  yesterday,
			},
			{
				InstanceID:  "instance-1",
				Domain:      "d2.example.com",
				IsGenerated: gu.Ptr(true),
				CreatedAt:   yesterday,
			},
			{
				InstanceID:  "instance-1",
				Domain:      "d3.example.com",
				IsPrimary:   gu.Ptr(true),
				IsGenerated: gu.Ptr(false),
				CreatedAt:   yesterday,
			},
			{
				InstanceID:  "instance-1",
				Domain:      "d4.example.com",
				IsPrimary:   gu.Ptr(false),
				IsGenerated: gu.Ptr(true),
				CreatedAt:   yesterday,
			},
		},
	}

	expectedInstance := &instance_v2.Instance{
		Id:           "instance-1",
		ChangeDate:   timestamppb.New(now),
		CreationDate: timestamppb.New(yesterday),
		State:        instance_v2.State_STATE_RUNNING,
		Name:         "Instance One",
		Version:      "",
		CustomDomains: []*instance_v2.CustomDomain{
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d1.example.com",
				Primary:      true,
				Generated:    false,
			},
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d2.example.com",
				Primary:      false,
				Generated:    true,
			},
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d3.example.com",
				Primary:      true,
				Generated:    false,
			},
			{
				InstanceId:   "instance-1",
				CreationDate: timestamppb.New(yesterday),
				Domain:       "d4.example.com",
				Primary:      false,
				Generated:    true,
			},
		},
	}

	// Test
	res := DomainInstanceModelToGRPCResponse(inputInstance)

	// Verify
	assert.Equal(t, expectedInstance, res)
}

func TestDomainInstanceListModelToGRPCBetaResponse(t *testing.T) {
	t.Parallel()
	now := time.Now()

	tt := []struct {
		testName       string
		inputResult    []*domain.Instance
		expectedResult []*instance_v2beta.Instance
	}{
		{
			testName:       "empty result",
			inputResult:    []*domain.Instance{},
			expectedResult: []*instance_v2beta.Instance{},
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
			expectedResult: []*instance_v2beta.Instance{
				{
					Id:           "instance1",
					Name:         "test-instance",
					CreationDate: timestamppb.New(now),
					ChangeDate:   timestamppb.New(now),
					Domains:      []*instance_v2beta.Domain{},
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
			expectedResult: []*instance_v2beta.Instance{
				{
					Id:           "instance1",
					Name:         "test-instance-1",
					CreationDate: timestamppb.New(now),
					ChangeDate:   timestamppb.New(now),
					Domains: []*instance_v2beta.Domain{
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
					Domains:      []*instance_v2beta.Domain{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Test
			result := DomainInstanceListModelToGRPCBetaResponse(tc.inputResult)

			// Verify
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDomainInstanceListModelToGRPCResponse(t *testing.T) {
	t.Parallel()
	now := time.Now()

	tt := []struct {
		testName       string
		inputResult    []*domain.Instance
		expectedResult []*instance_v2.Instance
	}{
		{
			testName:       "empty result",
			inputResult:    []*domain.Instance{},
			expectedResult: []*instance_v2.Instance{},
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
			expectedResult: []*instance_v2.Instance{
				{
					Id:            "instance1",
					Name:          "test-instance",
					CreationDate:  timestamppb.New(now),
					ChangeDate:    timestamppb.New(now),
					State:         instance_v2.State_STATE_RUNNING,
					CustomDomains: []*instance_v2.CustomDomain{},
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
			expectedResult: []*instance_v2.Instance{
				{
					Id:           "instance1",
					Name:         "test-instance-1",
					CreationDate: timestamppb.New(now),
					ChangeDate:   timestamppb.New(now),
					State:        instance_v2.State_STATE_RUNNING,
					CustomDomains: []*instance_v2.CustomDomain{
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
					Id:            "instance2",
					Name:          "test-instance-2",
					CreationDate:  timestamppb.New(now),
					ChangeDate:    timestamppb.New(now),
					State:         instance_v2.State_STATE_RUNNING,
					CustomDomains: []*instance_v2.CustomDomain{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Test
			result := DomainInstanceListModelToGRPCResponse(tc.inputResult)

			// Verify
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestListInstancesBetaRequestToV2Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		request *instance_v2beta.ListInstancesRequest
		want    *instance_v2.ListInstancesRequest
	}{
		{
			name: "empty request",
			request: &instance_v2beta.ListInstancesRequest{
				Pagination: &filter_v2beta.PaginationRequest{},
			},
			want: &instance_v2.ListInstancesRequest{
				Pagination:    &filter_v2.PaginationRequest{},
				SortingColumn: instance_v2.FieldName_FIELD_NAME_UNSPECIFIED,
				Filters:       []*instance_v2.Filter{},
			},
		},
		{
			name: "request with all fields",
			request: &instance_v2beta.ListInstancesRequest{
				Pagination: &filter_v2beta.PaginationRequest{
					Offset: 10,
					Limit:  20,
					Asc:    true,
				},
				SortingColumn: gu.Ptr(instance_v2beta.FieldName_FIELD_NAME_NAME),
				Queries: []*instance_v2beta.Query{
					{
						Query: &instance_v2beta.Query_DomainQuery{
							DomainQuery: &instance_v2beta.DomainsQuery{
								Domains: []string{"domain1.com", "domain2.com"},
							},
						},
					},
					{
						Query: &instance_v2beta.Query_IdQuery{
							IdQuery: &instance_v2beta.IdsQuery{
								Ids: []string{"id1", "id2"},
							},
						},
					},
				},
			},
			want: &instance_v2.ListInstancesRequest{
				Pagination: &filter_v2.PaginationRequest{
					Offset: 10,
					Limit:  20,
					Asc:    true,
				},
				SortingColumn: instance_v2.FieldName_FIELD_NAME_NAME,
				Filters: []*instance_v2.Filter{
					{
						Filter: &instance_v2.Filter_CustomDomainsFilter{
							CustomDomainsFilter: &instance_v2.CustomDomainsFilter{
								Domains: []string{"domain1.com", "domain2.com"},
							},
						},
					},
					{
						Filter: &instance_v2.Filter_InIdsFilter{
							InIdsFilter: &filter_v2.InIDsFilter{
								Ids: []string{"id1", "id2"},
							},
						},
					},
				},
			},
		},
		{
			name: "request with different sorting field",
			request: &instance_v2beta.ListInstancesRequest{
				Pagination:    &filter_v2beta.PaginationRequest{},
				SortingColumn: gu.Ptr(instance_v2beta.FieldName_FIELD_NAME_CREATION_DATE),
			},
			want: &instance_v2.ListInstancesRequest{
				Pagination:    &filter_v2.PaginationRequest{},
				SortingColumn: instance_v2.FieldName_FIELD_NAME_CREATION_DATE,
				Filters:       []*instance_v2.Filter{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ListInstancesBetaRequestToV2Request(tc.request)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestListInstancesBetaSortingColToV2Request(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldName *instance_v2beta.FieldName
		want      instance_v2.FieldName
	}{
		{
			name:      "nil field name",
			fieldName: nil,
			want:      instance_v2.FieldName_FIELD_NAME_UNSPECIFIED,
		},
		{
			name:      "creation date field",
			fieldName: gu.Ptr(instance_v2beta.FieldName_FIELD_NAME_CREATION_DATE),
			want:      instance_v2.FieldName_FIELD_NAME_CREATION_DATE,
		},
		{
			name:      "id field",
			fieldName: gu.Ptr(instance_v2beta.FieldName_FIELD_NAME_ID),
			want:      instance_v2.FieldName_FIELD_NAME_ID,
		},
		{
			name:      "name field",
			fieldName: gu.Ptr(instance_v2beta.FieldName_FIELD_NAME_NAME),
			want:      instance_v2.FieldName_FIELD_NAME_NAME,
		},
		{
			name:      "unspecified field",
			fieldName: gu.Ptr(instance_v2beta.FieldName_FIELD_NAME_UNSPECIFIED),
			want:      instance_v2.FieldName_FIELD_NAME_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := listInstancesBetaSortingColToV2Request(tt.fieldName)
			assert.Equal(t, tt.want, got)
		})
	}
}
