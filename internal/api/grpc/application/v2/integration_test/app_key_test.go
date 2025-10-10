//go:build integration

package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func TestCreateApplicationKey(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)
	createdApp := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	t.Parallel()

	tt := []struct {
		testName        string
		creationRequest *application.CreateApplicationKeyRequest
		inputCtx        context.Context

		expectedErrorType codes.Code
	}{
		{
			testName: "when application id is not found should return failed precondition",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				ApplicationId:  integration.ID(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateAPIApp request is valid should create application and return no error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				ApplicationId:  createdApp.GetApplicationId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
		},

		// LoginUser
		{
			testName: "when user has no project.application.write permission for application key generation should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &application.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				ApplicationId:  createdApp.GetApplicationId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// OrgOwner
		{
			testName: "when user is OrgOwner application key request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &application.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				ApplicationId:  createdApp.GetApplicationId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
		},

		// ProjectOwner
		{
			testName: "when user is ProjectOwner application key request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &application.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				ApplicationId:  createdApp.GetApplicationId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.ApplicationV2.CreateApplicationKey(tc.inputCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetKeyId())
				assert.NotZero(t, res.GetCreationDate())
				assert.NotZero(t, res.GetKeyDetails())
			}
		})
	}
}

func TestDeleteApplicationKey(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)
	createdApp := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	t.Parallel()

	tt := []struct {
		testName        string
		deletionRequest func(ttt *testing.T) *application.DeleteApplicationKeyRequest
		inputCtx        context.Context

		expectedErrorType codes.Code
	}{
		{
			testName: "when application key ID is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			deletionRequest: func(ttt *testing.T) *application.DeleteApplicationKeyRequest {
				return &application.DeleteApplicationKeyRequest{
					KeyId:         integration.ID(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetApplicationId(),
				}
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when valid application key ID should delete successfully",
			inputCtx: IAMOwnerCtx,
			deletionRequest: func(ttt *testing.T) *application.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetApplicationId(), time.Now().AddDate(0, 0, 1))

				return &application.DeleteApplicationKeyRequest{
					KeyId:         createdAppKey.GetKeyId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetApplicationId(),
				}
			},
		},

		// LoginUser
		{
			testName: "when user has no project.application.write permission for application key deletion should return permission error",
			inputCtx: LoginUserCtx,
			deletionRequest: func(ttt *testing.T) *application.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetApplicationId(), time.Now().AddDate(0, 0, 1))

				return &application.DeleteApplicationKeyRequest{
					KeyId:         createdAppKey.GetKeyId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetApplicationId(),
				}
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// ProjectOwner
		{
			testName: "when user is OrgOwner API request should succeed",
			inputCtx: projectOwnerCtx,
			deletionRequest: func(ttt *testing.T) *application.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetApplicationId(), time.Now().AddDate(0, 0, 1))

				return &application.DeleteApplicationKeyRequest{
					KeyId:         createdAppKey.GetKeyId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetApplicationId(),
				}
			},
		},

		// OrganizationOwner
		{
			testName: "when user is OrgOwner application key deletion request should succeed",
			inputCtx: OrgOwnerCtx,
			deletionRequest: func(ttt *testing.T) *application.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetApplicationId(), time.Now().AddDate(0, 0, 1))

				return &application.DeleteApplicationKeyRequest{
					KeyId:         createdAppKey.GetKeyId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetApplicationId(),
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			deletionReq := tc.deletionRequest(t)

			// When
			res, err := instance.Client.ApplicationV2.DeleteApplicationKey(tc.inputCtx, deletionReq)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotEmpty(t, res.GetDeletionDate())
			}
		})
	}
}
