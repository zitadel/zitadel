//go:build integration

package instance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateInstance(t *testing.T) {
	tt := []struct {
		testName          string
		inputRequest      *instance.CreateInstanceRequest
		inputContext      context.Context
		expectedErrorMsg  string
		expectedErrorCode codes.Code
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.CreateInstanceRequest{
				InstanceName: " ",
				FirstOrgName: " ",
				CustomDomain: " ",
				Owner: &instance.CreateInstanceRequest_Machine_{
					Machine: &instance.CreateInstanceRequest_Machine{
						UserName:            "owner",
						Name:                "owner",
						PersonalAccessToken: &instance.CreateInstanceRequest_PersonalAccessToken{},
					},
				},
				DefaultLanguage: "",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when invalid input should return no admin error",
			inputRequest: &instance.CreateInstanceRequest{
				InstanceName:    " ",
				FirstOrgName:    " ",
				CustomDomain:    " ",
				DefaultLanguage: "",
				Owner: &instance.CreateInstanceRequest_Human_{
					Human: &instance.CreateInstanceRequest_Human{
						Email: &instance.CreateInstanceRequest_Email{
							Email: " user@example.com ",
						},
					},
				},
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Given name in profile is empty (USER-UCej2)",
		},
		{
			testName: "when input and context are valid should return expected response and no error",
			inputRequest: &instance.CreateInstanceRequest{
				InstanceName: " ",
				FirstOrgName: " ",
				CustomDomain: " ",
				Owner: &instance.CreateInstanceRequest_Machine_{
					Machine: &instance.CreateInstanceRequest_Machine{
						UserName:            "owner",
						Name:                "owner",
						PersonalAccessToken: &instance.CreateInstanceRequest_PersonalAccessToken{},
					},
				},
				DefaultLanguage: "",
			},
			inputContext: CTXWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			cli, err := integration.NewDefaultClient(CTXWithSysAuthZ)
			require.Nil(t, err)

			// Test
			res, err := cli.InstanceV2Beta.CreateInstance(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())
			if tc.expectedErrorMsg == "" {
				require.NotEmpty(t, res.GetInstanceId())
				assert.NotEmpty(t, res.GetPat())

				require.NotNil(t, res.GetDetails())
				assert.NotEmpty(t, res.GetDetails().GetResourceOwner())
			}
		})
	}
}

func TestDeleteInstace(t *testing.T) {
	// Given
	inst := integration.NewInstance(CTXWithSysAuthZ)

	tt := []struct {
		testName           string
		inputRequest       *instance.DeleteInstanceRequest
		inputContext       context.Context
		expectedErrorMsg   string
		expectedErrorCode  codes.Code
		expectedInstanceID string
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: " ",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: " ",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "instance_id must not be empty (instance_id)",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID() + "invalid",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "Instance not found (COMMA-AE3GS)",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID(),
			},
			inputContext:       CTXWithSysAuthZ,
			expectedInstanceID: inst.ID(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.DeleteInstance(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())
			if tc.expectedErrorMsg == "" {
				require.NotNil(t, res)
				require.NotNil(t, res.GetDetails())
				assert.Equal(t, tc.expectedInstanceID, res.GetDetails().GetResourceOwner())
				assert.NotEmpty(t, res.GetDetails().GetChangeDate())
			}
		})
	}
}
