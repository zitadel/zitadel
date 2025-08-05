//go:build integration

package integration_test

import (
	_ "embed"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/service_provider_config_expected.json
	expectedProviderConfigJson []byte

	//go:embed testdata/service_provider_config_expected_schemas.json
	expectedSchemasJson []byte

	//go:embed testdata/service_provider_config_expected_resource_types.json
	expectedResourceTypesJson []byte

	//go:embed testdata/service_provider_config_expected_resource_type_user.json
	expectedResourceTypeUserJson []byte

	//go:embed testdata/service_provider_config_expected_user_schema.json
	expectedUserSchemaJson []byte
)

func TestServiceProviderConfig(t *testing.T) {
	resp, err := Instance.Client.SCIM.GetServiceProviderConfig(CTX, Instance.DefaultOrg.Id)
	assert.NoError(t, err)
	assertJsonEqual(t, expectedProviderConfigJson, resp)
}

func TestResourceTypes(t *testing.T) {
	resp, err := Instance.Client.SCIM.GetResourceTypes(CTX, Instance.DefaultOrg.Id)
	assert.NoError(t, err)
	assertJsonEqual(t, expectedResourceTypesJson, resp)
}

func TestResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceName string
		want         []byte
		wantErr      bool
	}{
		{
			name:         "user",
			resourceName: "User",
			want:         expectedResourceTypeUserJson,
		},
		{
			name:         "not found",
			resourceName: "foobar",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Instance.Client.SCIM.GetResourceType(CTX, Instance.DefaultOrg.Id, tt.resourceName)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assertJsonEqual(t, tt.want, resp)
		})
	}
}

func TestSchemas(t *testing.T) {
	resp, err := Instance.Client.SCIM.GetSchemas(CTX, Instance.DefaultOrg.Id)
	assert.NoError(t, err)
	assertJsonEqual(t, expectedSchemasJson, resp)
}

func TestSchema(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    []byte
		wantErr bool
	}{
		{
			name: "user",
			id:   "urn:ietf:params:scim:schemas:core:2.0:User",
			want: expectedUserSchemaJson,
		},
		{
			name:    "not found",
			id:      "foobar",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Instance.Client.SCIM.GetSchema(CTX, Instance.DefaultOrg.Id, tt.id)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assertJsonEqual(t, tt.want, resp)
		})
	}
}

func assertJsonEqual(t *testing.T, expected, actual []byte) {
	t.Helper()

	// replace dynamic data json
	expectedJson := strings.Replace(string(expected), "{domain}", Instance.Domain, 1)
	expectedJson = strings.Replace(expectedJson, "{orgId}", Instance.DefaultOrg.Id, 1)
	assert.Equal(t, normalizeJson(t, []byte(expectedJson)), normalizeJson(t, actual))
}

func normalizeJson(t *testing.T, content []byte) string {
	t.Helper()

	raw := new(json.RawMessage)
	err := json.Unmarshal(content, raw)
	require.NoError(t, err)

	content, err = json.MarshalIndent(raw, "", "  ")
	require.NoError(t, err)
	return string(content) // use string for easier assertion diffs
}
