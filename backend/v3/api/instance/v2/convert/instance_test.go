package convert

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

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

	expectedInstance := &instance.Instance{
		Id:           "instance-1",
		ChangeDate:   timestamppb.New(now),
		CreationDate: timestamppb.New(yesterday),
		State:        instance.State_STATE_RUNNING,
		Name:         "Instance One",
		Version:      "",
		Domains: []*instance.Domain{
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

func TestDomainInstanceListModelToGRPCResponse(t *testing.T) {
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
			// Test
			result := DomainInstanceListModelToGRPCResponse(tc.inputResult)

			// Verify
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
