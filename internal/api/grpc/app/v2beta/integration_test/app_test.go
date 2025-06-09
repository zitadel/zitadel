//go:build integration

package instance_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func TestCreateApplication(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	orgNotInCtx := instance.CreateOrganization(iamOwnerCtx, gofakeit.Name(), gofakeit.Email())
	p := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, gofakeit.AppName(), false, false)
	pNotInCtx := instance.CreateProject(iamOwnerCtx, t, orgNotInCtx.GetOrganizationId(), gofakeit.AppName(), false, false)

	tt := []struct {
		testName        string
		creationRequest *app.CreateApplicationRequest

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		{
			testName: "when project for API app creation is not found should return failed precondition error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Name:      "App Name",
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateAPIApp request is valid should create app and return no error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      "App Name",
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &app.CreateApplicationResponse_ApiResponse{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			res, err := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				resType := fmt.Sprintf("%T", res.GetCreationResponseType())
				assert.Equal(t, tc.expectedResponseType, resType)
				assert.NotEmpty(t, res.GetAppId())
				assert.NotEmpty(t, res.GetCreationDate())
			}

		})
	}
}
