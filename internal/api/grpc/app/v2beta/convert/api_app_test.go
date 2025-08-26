package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func TestCreateAPIApplicationRequestToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		appName   string
		projectID string
		appID     string
		req       *app.CreateAPIApplicationRequest
		want      *domain.APIApp
	}{
		{
			name:      "basic auth method",
			appName:   "my-app",
			projectID: "proj-1",
			appID:     "someID",
			req: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-1"},
				AppName:        "my-app",
				AuthMethodType: domain.APIAuthMethodTypeBasic,
				AppID:          "someID",
			},
		},
		{
			name:      "private key jwt",
			appName:   "jwt-app",
			projectID: "proj-2",
			req: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-2"},
				AppName:        "jwt-app",
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got := CreateAPIApplicationRequestToDomain(tt.appName, tt.projectID, tt.appID, tt.req)

			// Then
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateAPIApplicationConfigurationRequestToDomain(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		appID     string
		projectID string
		req       *app.UpdateAPIApplicationConfigurationRequest
		want      *domain.APIApp
	}{
		{
			name:      "basic auth method",
			appID:     "app-1",
			projectID: "proj-1",
			req: &app.UpdateAPIApplicationConfigurationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-1"},
				AppID:          "app-1",
				AuthMethodType: domain.APIAuthMethodTypeBasic,
			},
		},
		{
			name:      "private key jwt",
			appID:     "app-2",
			projectID: "proj-2",
			req: &app.UpdateAPIApplicationConfigurationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-2"},
				AppID:          "app-2",
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			got := UpdateAPIApplicationConfigurationRequestToDomain(tt.appID, tt.projectID, tt.req)

			// Then
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_apiAuthMethodTypeToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name           string
		methodType     domain.APIAuthMethodType
		expectedResult app.APIAuthMethodType
	}{
		{
			name:           "basic auth method",
			methodType:     domain.APIAuthMethodTypeBasic,
			expectedResult: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
		},
		{
			name:           "private key jwt",
			methodType:     domain.APIAuthMethodTypePrivateKeyJWT,
			expectedResult: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
		},
		{
			name:           "unknown auth method defaults to basic",
			expectedResult: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			res := apiAuthMethodTypeToPb(tc.methodType)

			// Then
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
func TestGetApplicationKeyQueriesRequestToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName       string
		inputOrgID     string
		inputProjectID string
		inputAppID     string

		expectedQueriesLength int
	}{
		{
			testName:              "all IDs provided",
			inputOrgID:            "org-1",
			inputProjectID:        "proj-1",
			inputAppID:            "app-1",
			expectedQueriesLength: 3,
		},
		{
			testName:              "only org ID",
			inputOrgID:            "org-1",
			inputProjectID:        " ",
			inputAppID:            "",
			expectedQueriesLength: 1,
		},
		{
			testName:              "only project ID",
			inputOrgID:            "",
			inputProjectID:        "proj-1",
			inputAppID:            " ",
			expectedQueriesLength: 1,
		},
		{
			testName:              "only app ID",
			inputOrgID:            " ",
			inputProjectID:        "",
			inputAppID:            "app-1",
			expectedQueriesLength: 1,
		},
		{
			testName:              "empty IDs",
			inputOrgID:            " ",
			inputProjectID:        " ",
			inputAppID:            " ",
			expectedQueriesLength: 0,
		},
		{
			testName:              "with spaces",
			inputOrgID:            " org-1 ",
			inputProjectID:        " proj-1 ",
			inputAppID:            " app-1 ",
			expectedQueriesLength: 3,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			got, err := GetApplicationKeyQueriesRequestToDomain(tc.inputOrgID, tc.inputProjectID, tc.inputAppID)

			// Then
			require.NoError(t, err)

			assert.Len(t, got, tc.expectedQueriesLength)
		})
	}
}
