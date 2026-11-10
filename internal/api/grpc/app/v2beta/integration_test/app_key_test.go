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
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func TestCreateApplicationKey(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)
	createdApp := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	t.Parallel()

	tt := []struct {
		testName        string
		creationRequest *app.CreateApplicationKeyRequest
		inputCtx        context.Context

		expectedErrorType codes.Code
	}{
		{
			testName: "when app id is not found should return failed precondition",
			inputCtx: IAMOwnerCtx,
			creationRequest: &app.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				AppId:          integration.ID(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateAPIApp request is valid should create app and return no error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &app.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				AppId:          createdApp.GetAppId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
		},

		// LoginUser
		{
			testName: "when user has no project.app.write permission for app key generation should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &app.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				AppId:          createdApp.GetAppId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// OrgOwner
		{
			testName: "when user is OrgOwner app key request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &app.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				AppId:          createdApp.GetAppId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
		},

		// ProjectOwner
		{
			testName: "when user is ProjectOwner app key request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &app.CreateApplicationKeyRequest{
				ProjectId:      p.GetId(),
				AppId:          createdApp.GetAppId(),
				ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.AppV2Beta.CreateApplicationKey(tc.inputCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetId())
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
		deletionRequest func(ttt *testing.T) *app.DeleteApplicationKeyRequest
		inputCtx        context.Context

		expectedErrorType codes.Code
	}{
		{
			testName: "when app key ID is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			deletionRequest: func(ttt *testing.T) *app.DeleteApplicationKeyRequest {
				return &app.DeleteApplicationKeyRequest{
					Id:            integration.ID(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetAppId(),
				}
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when valid app key ID should delete successfully",
			inputCtx: IAMOwnerCtx,
			deletionRequest: func(ttt *testing.T) *app.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetAppId(), time.Now().AddDate(0, 0, 1))

				return &app.DeleteApplicationKeyRequest{
					Id:            createdAppKey.GetId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetAppId(),
				}
			},
		},

		// LoginUser
		{
			testName: "when user has no project.app.write permission for app key deletion should return permission error",
			inputCtx: LoginUserCtx,
			deletionRequest: func(ttt *testing.T) *app.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetAppId(), time.Now().AddDate(0, 0, 1))

				return &app.DeleteApplicationKeyRequest{
					Id:            createdAppKey.GetId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetAppId(),
				}
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// ProjectOwner
		{
			testName: "when user is OrgOwner API request should succeed",
			inputCtx: projectOwnerCtx,
			deletionRequest: func(ttt *testing.T) *app.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetAppId(), time.Now().AddDate(0, 0, 1))

				return &app.DeleteApplicationKeyRequest{
					Id:            createdAppKey.GetId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetAppId(),
				}
			},
		},

		// OrganizationOwner
		{
			testName: "when user is OrgOwner app key deletion request should succeed",
			inputCtx: OrgOwnerCtx,
			deletionRequest: func(ttt *testing.T) *app.DeleteApplicationKeyRequest {
				createdAppKey := createAppKey(ttt, IAMOwnerCtx, instance, p.GetId(), createdApp.GetAppId(), time.Now().AddDate(0, 0, 1))

				return &app.DeleteApplicationKeyRequest{
					Id:            createdAppKey.GetId(),
					ProjectId:     p.GetId(),
					ApplicationId: createdApp.GetAppId(),
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
			res, err := instance.Client.AppV2Beta.DeleteApplicationKey(tc.inputCtx, deletionReq)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotEmpty(t, res.GetDeletionDate())
			}
		})
	}
}
