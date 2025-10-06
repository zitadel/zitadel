package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func TestGetInstanceCommand_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		inputInstanceID string
		ctx             context.Context
		expectedError   error
	}{
		{
			name:            "empty instance id, error",
			inputInstanceID: "",
			ctx:             context.Background(),
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-32a0o2", "invalid instance ID"),
		},
		{
			name:            "whitespace instance id, error",
			inputInstanceID: "   ",
			ctx:             context.Background(),
			expectedError:   zerrors.ThrowInvalidArgument(nil, "DOM-32a0o2", "invalid instance ID"),
		},
		{
			name:            "instance id mismatch, error",
			inputInstanceID: "instance1",
			ctx:             authz.NewMockContext("instance2", "", ""),
			expectedError:   zerrors.ThrowPermissionDenied(nil, "DOM-n0SvVB", "input instance ID doesn't match context instance"),
		},
		{
			name:            "valid instance id, success",
			inputInstanceID: "instance1",
			ctx:             authz.NewMockContext("instance1", "", ""),
			expectedError:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewGetInstanceCommand(tt.inputInstanceID)

			// Test
			err := cmd.Validate(tt.ctx, nil)

			// Verify
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
func TestGetInstanceCommand_Events(t *testing.T) {
	t.Parallel()

	// Given
	cmd := domain.NewGetInstanceCommand("instance1")

	// Test
	events, err := cmd.Events(context.Background(), nil)

	// Verify
	assert.NoError(t, err)
	assert.Empty(t, events)
}

func TestGetInstanceCommand_Execute(t *testing.T) {
	t.Parallel()

	ctx := authz.NewMockContext("inst-1", "org-1", gofakeit.UUID())
	txInitErr := errors.New("tx init error")
	getErr := errors.New("get error")

	tt := []struct {
		testName string

		mockTx       func(ctrl *gomock.Controller) database.QueryExecutor
		instanceRepo func(ctrl *gomock.Controller) domain.InstanceRepository

		inputInstanceID string

		expectedError    error
		expectedInstance *domain.Instance
	}{
		{
			testName: "when EnsureTx fails should return error",
			mockTx: func(ctrl *gomock.Controller) database.QueryExecutor {
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
			testName: "when retrieving instance fails should return error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(nil, getErr)
				return instanceRepo
			},
			inputInstanceID: "instance-1",
			expectedError:   getErr,
		},
		{
			testName: "when retrieving instance succeeds should set instance on struct and return no error",
			instanceRepo: func(ctrl *gomock.Controller) domain.InstanceRepository {
				instanceRepo := domainmock.NewInstanceRepo(ctrl)

				instanceRepo.EXPECT().
					LoadDomains().
					Times(1).
					Return(instanceRepo)

				instanceRepo.EXPECT().
					Get(
						gomock.Any(),
						gomock.Any(),
						dbmock.QueryOptions(database.WithCondition(
							instanceRepo.IDCondition("instance-1"),
						)),
					).
					Times(1).
					Return(&domain.Instance{
						ID:   "instance-1",
						Name: "My instance 1",
						Domains: []*domain.InstanceDomain{
							{Domain: "d1.example.com"},
							{Domain: "d2.example.com"},
							{Domain: "d3.example.com"},
						},
					}, nil)

				return instanceRepo
			},
			inputInstanceID: "instance-1",
			expectedInstance: &domain.Instance{
				ID:   "instance-1",
				Name: "My instance 1",
				Domains: []*domain.InstanceDomain{
					{Domain: "d1.example.com"},
					{Domain: "d2.example.com"},
					{Domain: "d3.example.com"},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			g := domain.NewGetInstanceCommand(tc.inputInstanceID)
			ctrl := gomock.NewController(t)
			opts := &domain.CommandOpts{DB: new(noopdb.Pool)}

			if tc.mockTx != nil {
				opts.DB = tc.mockTx(ctrl)
			}
			if tc.instanceRepo != nil {
				opts.SetInstanceRepo(tc.instanceRepo(ctrl))
			}

			// Test
			err := g.Execute(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedInstance, g.ReturnedInstance)
		})
	}
}

func TestGetInstanceCommand_ResultToGRPC(t *testing.T) {
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

	cmd := domain.GetInstanceCommand{ReturnedInstance: inputInstance}

	// Test
	res := cmd.ResultToGRPC()

	// Verify
	assert.Equal(t, expectedInstance, res)
}
