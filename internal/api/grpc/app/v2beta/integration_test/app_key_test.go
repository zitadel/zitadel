//go:build integration

package instance_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func TestCreateApplicationKey(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)
	createdApp := createAPIApp(t, p.GetId())

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
				AppId:          gofakeit.UUID(),
				Type:           app.ApplicationKeyType_APPLICATION_KEY_TYPE_JSON,
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
				Type:           app.ApplicationKeyType_APPLICATION_KEY_TYPE_JSON,
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
				Type:           app.ApplicationKeyType_APPLICATION_KEY_TYPE_JSON,
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
				Type:           app.ApplicationKeyType_APPLICATION_KEY_TYPE_JSON,
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
				Type:           app.ApplicationKeyType_APPLICATION_KEY_TYPE_JSON,
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
