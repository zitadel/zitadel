//go:build integration

package management_test

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_ListAppKeys(t *testing.T) {
	// create project
	prjName := integration.ProjectName()
	createPrjRes, err := Instance.Client.Projectv2Beta.CreateProject(IAMOwnerCTX, &project.CreateProjectRequest{
		Name:           prjName,
		OrganizationId: Instance.DefaultOrg.Id,
	})
	require.NoError(t, err)
	prjId := createPrjRes.Id

	// add app to project
	createAppjRes, err := Instance.Client.AppV2Beta.CreateApplication(IAMOwnerCTX, &app.CreateApplicationRequest{
		ProjectId: prjId,
		Name:      integration.ApplicationName(),
		CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
			ApiRequest: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
			},
		},
	})
	require.NoError(t, err)
	appId := createAppjRes.AppId

	type test struct {
		name               string
		expectedKeyIdsFunc func() []string
	}

	tests := []test{
		{
			name: "happy path",
			expectedKeyIdsFunc: func() []string {
				// add other app to project
				createOtherAppjRes, err := Instance.Client.AppV2Beta.CreateApplication(IAMOwnerCTX, &app.CreateApplicationRequest{
					ProjectId: prjId,
					Name:      integration.ApplicationName(),
					CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
						ApiRequest: &app.CreateAPIApplicationRequest{
							AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
						},
					},
				})
				require.NoError(t, err)
				otherAppId := createOtherAppjRes.AppId
				// add other project key ids - These SHOULD NOT be returned when calling ListAppKeys()
				for range 5 {
					_, err := Instance.Client.AppV2Beta.CreateApplicationKey(IAMOwnerCTX, &app.CreateApplicationKeyRequest{
						AppId:          otherAppId,
						ProjectId:      prjId,
						ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
					})
					require.NoError(t, err)
				}

				// create app keys we expect to be rturned form ListAppKeys()
				keyIDs := make([]string, 5)
				for i := range len(keyIDs) {
					res, err := Instance.Client.AppV2Beta.CreateApplicationKey(IAMOwnerCTX, &app.CreateApplicationKeyRequest{
						AppId:          appId,
						ProjectId:      prjId,
						ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 1).UTC()),
					})
					require.NoError(t, err)
					keyIDs[i] = res.Id
				}
				return keyIDs
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedKeyIds := tt.expectedKeyIdsFunc()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ct *assert.CollectT) {
				res, err := Client.ListAppKeys(IAMOwnerCTX, &management.ListAppKeysRequest{
					AppId:     appId,
					ProjectId: prjId,
				})
				require.NoError(ct, err)
				assert.Equal(ct, len(expectedKeyIds), len(res.GetResult()))

				for _, key := range res.GetResult() {
					assert.True(ct, slices.Contains(expectedKeyIds, key.Id))
				}
			}, retryDuration, tick)
		})
	}
}
