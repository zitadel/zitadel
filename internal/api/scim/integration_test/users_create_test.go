//go:build integration

package integration_test

import (
	"context"
	_ "embed"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"google.golang.org/grpc/codes"
	"net/http"
	"path"
	"testing"
)

var (
	//go:embed testdata/users_create_test_minimal.json
	minimalUserJson []byte

	//go:embed testdata/users_create_test_full.json
	fullUserJson []byte

	//go:embed testdata/users_create_test_missing_username.json
	missingUserNameUserJson []byte

	//go:embed testdata/users_create_test_missing_name.json
	missingNameUserJson []byte

	//go:embed testdata/users_create_test_missing_email.json
	missingEmailUserJson []byte

	//go:embed testdata/users_create_test_invalid_password.json
	invalidPasswordUserJson []byte

	//go:embed testdata/users_create_test_invalid_profile_url.json
	invalidProfileUrlUserJson []byte

	//go:embed testdata/users_create_test_invalid_locale.json
	invalidLocaleUserJson []byte

	//go:embed testdata/users_create_test_invalid_timezone.json
	invalidTimeZoneUserJson []byte
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		body          []byte
		ctx           context.Context
		wantErr       bool
		scimErrorType string
		errorStatus   int
		zitadelErrID  string
	}{
		{
			name: "minimal user",
			body: minimalUserJson,
		},
		{
			name: "full user",
			body: fullUserJson,
		},
		{
			name:          "missing userName",
			wantErr:       true,
			scimErrorType: "invalidValue",
			body:          missingUserNameUserJson,
		},
		{
			// this is an expected schema violation
			name:          "missing name",
			wantErr:       true,
			scimErrorType: "invalidValue",
			body:          missingNameUserJson,
		},
		{
			name:          "missing email",
			wantErr:       true,
			scimErrorType: "invalidValue",
			body:          missingEmailUserJson,
		},
		{
			name:          "password complexity violation",
			wantErr:       true,
			scimErrorType: "invalidValue",
			body:          invalidPasswordUserJson,
		},
		{
			name:          "invalid profile url",
			wantErr:       true,
			scimErrorType: "invalidValue",
			zitadelErrID:  "SCIM-htturl1",
			body:          invalidProfileUrlUserJson,
		},
		{
			name:          "invalid time zone",
			wantErr:       true,
			scimErrorType: "invalidValue",
			body:          invalidTimeZoneUserJson,
		},
		{
			name:          "invalid locale",
			wantErr:       true,
			scimErrorType: "invalidValue",
			body:          invalidLocaleUserJson,
		},
		{
			name:        "not authenticated",
			body:        minimalUserJson,
			ctx:         context.Background(),
			wantErr:     true,
			errorStatus: http.StatusUnauthorized,
		},
		{
			name:        "no permissions",
			body:        minimalUserJson,
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			wantErr:     true,
			errorStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			createdUser, err := Instance.Client.SCIM.Users.Create(ctx, Instance.DefaultOrg.Id, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				statusCode := tt.errorStatus
				if statusCode == 0 {
					statusCode = http.StatusBadRequest
				}
				scimErr := scim.RequireScimError(t, statusCode, err)
				assert.Equal(t, tt.scimErrorType, scimErr.Error.ScimType)
				if tt.zitadelErrID != "" {
					assert.Equal(t, tt.zitadelErrID, scimErr.Error.ZitadelDetail.ID)
				}

				return
			}

			assert.NotEmpty(t, createdUser.ID)
			assert.EqualValues(t, []schemas.ScimSchemaType{"urn:ietf:params:scim:schemas:core:2.0:User"}, createdUser.Resource.Schemas)
			assert.Equal(t, schemas.ScimResourceTypeSingular("User"), createdUser.Resource.Meta.ResourceType)
			assert.Equal(t, "http://"+Instance.Host()+path.Join(schemas.HandlerPrefix, Instance.DefaultOrg.Id, "Users", createdUser.ID), createdUser.Resource.Meta.Location)
			assert.Nil(t, createdUser.Password)

			_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
			assert.NoError(t, err)
		})
	}
}

func TestCreateUser_duplicate(t *testing.T) {
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, minimalUserJson)
	require.NoError(t, err)

	_, err = Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, minimalUserJson)
	scimErr := scim.RequireScimError(t, http.StatusConflict, err)
	assert.Equal(t, "User already exists", scimErr.Error.Detail)

	_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
	require.NoError(t, err)
}

func TestCreateUser_metadata(t *testing.T) {
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
	require.NoError(t, err)

	md, err := Instance.Client.Mgmt.ListUserMetadata(CTX, &management.ListUserMetadataRequest{
		Id: createdUser.ID,
	})
	require.NoError(t, err)

	mdMap := make(map[string]string)
	for i := range md.Result {
		mdMap[md.Result[i].Key] = string(md.Result[i].Value)
	}

	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:name.honorificPrefix", "Ms.")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:timezone", "America/Los_Angeles")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:photos", `[{"value":"https://photos.example.com/profilephoto/72930000000Ccne/F","type":"photo"},{"value":"https://photos.example.com/profilephoto/72930000000Ccne/T","type":"thumbnail"}]`)
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:addresses", `[{"type":"work","streetAddress":"100 Universal City Plaza","locality":"Hollywood","region":"CA","postalCode":"91608","country":"USA","formatted":"100 Universal City Plaza\nHollywood, CA 91608 USA","primary":true},{"type":"home","streetAddress":"456 Hollywood Blvd","locality":"Hollywood","region":"CA","postalCode":"91608","country":"USA","formatted":"456 Hollywood Blvd\nHollywood, CA 91608 USA"}]`)
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:entitlements", `[{"value":"my-entitlement-1","display":"Entitlement 1","type":"main-entitlement","primary":true},{"value":"my-entitlement-2","display":"Entitlement 2","type":"secondary-entitlement"}]`)
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:externalId", "701984")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:name.middleName", "Jane")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:name.honorificSuffix", "III")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:profileURL", "http://login.example.com/bjensen")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:title", "Tour Guide")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:locale", "en-US")
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:ims", `[{"value":"someaimhandle","type":"aim"},{"value":"twitterhandle","type":"X"}]`)
	integration.AssertMapContains(t, mdMap, "urn:zitadel:scim:roles", `[{"value":"my-role-1","display":"Rolle 1","type":"main-role","primary":true},{"value":"my-role-2","display":"Rolle 2","type":"secondary-role"}]`)

	_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
	require.NoError(t, err)
}

func TestCreateUser_scopedExternalID(t *testing.T) {
	_, err := Instance.Client.Mgmt.SetUserMetadata(CTX, &management.SetUserMetadataRequest{
		Id:    Instance.Users.Get(integration.UserTypeOrgOwner).ID,
		Key:   "urn:zitadel:scim:provisioning_domain",
		Value: []byte("fooBar"),
	})
	require.NoError(t, err)

	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
	require.NoError(t, err)

	// unscoped externalID should not exist
	_, err = Instance.Client.Mgmt.GetUserMetadata(CTX, &management.GetUserMetadataRequest{
		Id:  createdUser.ID,
		Key: "urn:zitadel:scim:externalId",
	})
	integration.AssertGrpcStatus(t, codes.NotFound, err)

	// scoped externalID should exist
	md, err := Instance.Client.Mgmt.GetUserMetadata(CTX, &management.GetUserMetadataRequest{
		Id:  createdUser.ID,
		Key: "urn:zitadel:scim:fooBar:externalId",
	})
	require.NoError(t, err)
	assert.Equal(t, "701984", string(md.Metadata.Value))

	_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
	require.NoError(t, err)
}

func TestCreateUser_anotherOrg(t *testing.T) {
	org := Instance.CreateOrganization(Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner), gofakeit.Name(), gofakeit.Email())
	_, err := Instance.Client.SCIM.Users.Create(CTX, org.OrganizationId, fullUserJson)
	scim.RequireScimError(t, http.StatusNotFound, err)
}
