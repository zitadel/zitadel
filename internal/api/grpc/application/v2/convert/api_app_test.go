package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func TestCreateAPIApplicationRequestToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		appName   string
		projectID string
		appID     string
		req       *application.CreateAPIApplicationRequest
		want      *domain.APIApp
	}{
		{
			name:      "basic auth method",
			appName:   "my-application",
			projectID: "proj-1",
			appID:     "someID",
			req: &application.CreateAPIApplicationRequest{
				AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-1"},
				AppName:        "my-application",
				AuthMethodType: domain.APIAuthMethodTypeBasic,
				AppID:          "someID",
			},
		},
		{
			name:      "private key jwt",
			appName:   "jwt-application",
			projectID: "proj-2",
			req: &application.CreateAPIApplicationRequest{
				AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-2"},
				AppName:        "jwt-application",
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
		req       *application.UpdateAPIApplicationConfigurationRequest
		want      *domain.APIApp
	}{
		{
			name:      "basic auth method",
			appID:     "application-1",
			projectID: "proj-1",
			req: &application.UpdateAPIApplicationConfigurationRequest{
				AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-1"},
				AppID:          "application-1",
				AuthMethodType: domain.APIAuthMethodTypeBasic,
			},
		},
		{
			name:      "private key jwt",
			appID:     "application-2",
			projectID: "proj-2",
			req: &application.UpdateAPIApplicationConfigurationRequest{
				AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
			},
			want: &domain.APIApp{
				ObjectRoot:     models.ObjectRoot{AggregateID: "proj-2"},
				AppID:          "application-2",
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
		expectedResult application.APIAuthMethodType
	}{
		{
			name:           "basic auth method",
			methodType:     domain.APIAuthMethodTypeBasic,
			expectedResult: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
		},
		{
			name:           "private key jwt",
			methodType:     domain.APIAuthMethodTypePrivateKeyJWT,
			expectedResult: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
		},
		{
			name:           "unknown auth method defaults to basic",
			expectedResult: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
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
			inputAppID:            "application-1",
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
			testName:              "only application ID",
			inputOrgID:            " ",
			inputProjectID:        "",
			inputAppID:            "application-1",
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
			inputAppID:            " application-1 ",
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
