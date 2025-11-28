package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
)

func TestDeleteSessionCommand_Validate(t *testing.T) {
	t.Parallel()
	ctx := authz.NewMockContext("inst-1", "org-default", gofakeit.UUID())
	getErr := errors.New("get error")

	tt := []struct {
		testName             string
		orgRepo              func(ctrl *gomock.Controller) domain.OrganizationRepository
		projectRepo          func(ctrl *gomock.Controller) domain.ProjectRepository
		inputOrganizationID  string
		inputSessionVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
		expectedError        error
	}{
		{
			testName:             "",
			orgRepo:              nil,
			projectRepo:          nil,
			inputOrganizationID:  "",
			inputSessionVerifier: mockSessionVerification(nil),
			expectedError:        nil,
		}, {
			testName:             "",
			orgRepo:              nil,
			projectRepo:          nil,
			inputOrganizationID:  "",
			inputSessionVerifier: mockSessionVerification(getErr),
			expectedError:        nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			d := domain.NewDeleteSessionCommand(tc.inputOrganizationID)
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.orgRepo != nil {
				domain.WithOrganizationRepo(tc.orgRepo(ctrl))(opts)
			}
			if tc.projectRepo != nil {
				domain.WithProjectRepo(tc.projectRepo(ctrl))(opts)
			}

			// Test
			err := d.Validate(ctx, opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func mockSessionVerification(err error) func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		return err
	}
}
