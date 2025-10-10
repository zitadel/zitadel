//go:build integration

package integration_test

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/internal/test"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

var (
	//go:embed testdata/users_replace_test_minimal_with_external_id.json
	minimalUserWithExternalIDJson []byte

	//go:embed testdata/users_replace_test_minimal_with_email_type.json
	minimalUserWithEmailTypeReplaceJson []byte

	//go:embed testdata/users_replace_test_minimal.json
	minimalUserReplaceJson []byte

	//go:embed testdata/users_replace_test_full.json
	fullUserReplaceJson []byte
)

func TestReplaceUser(t *testing.T) {
	tests := []struct {
		name             string
		body             []byte
		ctx              context.Context
		createUserOrgID  string
		replaceUserOrgID string
		want             *resources.ScimUser
		wantErr          bool
		scimErrorType    string
		errorStatus      int
		zitadelErrID     string
	}{
		{
			name: "minimal user",
			body: minimalUserReplaceJson,
			want: &resources.ScimUser{
				UserName: "acmeUser1-minimal-replaced",
				Name: &resources.ScimUserName{
					FamilyName: "Ross-replaced",
					GivenName:  "Bethany-replaced",
				},
				Emails: []*resources.ScimEmail{
					{
						Value:   "user1-minimal-replaced@example.com",
						Primary: true,
					},
				},
			},
		},
		{
			name: "full user",
			body: fullUserReplaceJson,
			want: &resources.ScimUser{
				ExternalID: "701984-updated",
				UserName:   "bjensen-replaced-full@example.com",
				Name: &resources.ScimUserName{
					Formatted:       "Babs Jensen-updated", // display name takes precedence
					FamilyName:      "Jensen-updated",
					GivenName:       "Barbara-updated",
					MiddleName:      "Jane-updated",
					HonorificPrefix: "Ms.-updated",
					HonorificSuffix: "III",
				},
				DisplayName: "Babs Jensen-updated",
				NickName:    "Babs-updated",
				ProfileUrl:  test.Must(schemas.ParseHTTPURL("http://login.example.com/bjensen-updated")),
				Emails: []*resources.ScimEmail{
					{
						Value:   "bjensen-replaced-full@example.com",
						Primary: true,
					},
				},
				Addresses: []*resources.ScimAddress{
					{
						Type:          "work-updated",
						StreetAddress: "100 Universal City Plaza-updated",
						Locality:      "Hollywood-updated",
						Region:        "CA-updated",
						PostalCode:    "91608-updated",
						Country:       "USA-updated",
						Formatted:     "100 Universal City Plaza\nHollywood, CA 91608 USA-updated",
						Primary:       true,
					},
					{
						Type:          "home-updated",
						StreetAddress: "456 Hollywood Blvd-updated",
						Locality:      "Hollywood-updated",
						Region:        "CA-updated",
						PostalCode:    "91608-updated",
						Country:       "USA-updated",
						Formatted:     "456 Hollywood Blvd\nHollywood, CA 91608 USA-updated",
					},
				},
				PhoneNumbers: []*resources.ScimPhoneNumber{
					{
						Value:   "+4155555555558732833",
						Primary: true,
					},
				},
				Ims: []*resources.ScimIms{
					{
						Value: "someaimhandle-updated",
						Type:  "aim-updated",
					},
					{
						Value: "twitterhandle-updated",
						Type:  "X-updated",
					},
				},
				Photos: []*resources.ScimPhoto{
					{
						Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/F-updated")),
						Type:  "photo-updated",
					},
					{
						Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/T-updated")),
						Type:  "thumbnail-updated",
					},
				},
				Roles: []*resources.ScimRole{
					{
						Value:   "my-role-1-updated",
						Display: "Rolle 1-updated",
						Type:    "main-role-updated",
						Primary: true,
					},
					{
						Value:   "my-role-2-updated",
						Display: "Rolle 2-updated",
						Type:    "secondary-role-updated",
						Primary: false,
					},
				},
				Entitlements: []*resources.ScimEntitlement{
					{
						Value:   "my-entitlement-1-updated",
						Display: "Entitlement 1-updated",
						Type:    "main-entitlement-updated",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2-updated",
						Display: "Entitlement 2-updated",
						Type:    "secondary-entitlement-updated",
						Primary: false,
					},
				},
				Title:             "Tour Guide-updated",
				PreferredLanguage: language.MustParse("en-CH"),
				Locale:            "en-CH",
				Timezone:          "Europe/Zurich",
				Active:            schemas.NewRelaxedBool(false),
			},
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
			body:        withUsername(minimalUserJson, integration.Username()),
			ctx:         context.Background(),
			wantErr:     true,
			errorStatus: http.StatusUnauthorized,
		},
		{
			name:        "no permissions",
			body:        withUsername(minimalUserJson, integration.Username()),
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			wantErr:     true,
			errorStatus: http.StatusNotFound,
		},
		{
			name:             "another org",
			body:             withUsername(minimalUserJson, integration.Username()),
			replaceUserOrgID: SecondaryOrganization.OrganizationId,
			wantErr:          true,
			errorStatus:      http.StatusNotFound,
		},
		{
			name:             "another org with permissions",
			body:             withUsername(minimalUserJson, integration.Username()),
			replaceUserOrgID: SecondaryOrganization.OrganizationId,
			ctx:              Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner),
			wantErr:          true,
			errorStatus:      http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// use iam owner => we don't want to test permissions of the create endpoint.
			createdUser, err := Instance.Client.SCIM.Users.Create(Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner), Instance.DefaultOrg.Id, withUsername(fullUserJson, integration.Username()))
			require.NoError(t, err)

			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			replaceUserOrgID := tt.replaceUserOrgID
			if replaceUserOrgID == "" {
				replaceUserOrgID = Instance.DefaultOrg.Id
			}

			replacedUser, err := Instance.Client.SCIM.Users.Replace(ctx, replaceUserOrgID, createdUser.ID, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceUser() error = %v, wantErr %v", err, tt.wantErr)
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

			assert.NotEmpty(t, replacedUser.ID)
			assert.EqualValues(t, []schemas.ScimSchemaType{"urn:ietf:params:scim:schemas:core:2.0:User"}, replacedUser.Resource.Schemas)
			assert.Equal(t, schemas.ScimResourceTypeSingular("User"), replacedUser.Resource.Meta.ResourceType)
			assert.Equal(t, "http://"+Instance.Host()+path.Join(schemas.HandlerPrefix, Instance.DefaultOrg.Id, "Users", createdUser.ID), replacedUser.Resource.Meta.Location)
			assert.Nil(t, createdUser.Password)

			if !test.PartiallyDeepEqual(tt.want, replacedUser) {
				t.Errorf("ReplaceUser() got = %#v, want %#v", replacedUser, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// ensure the user is really stored and not just returned to the caller
				fetchedUser, err := Instance.Client.SCIM.Users.Get(CTX, Instance.DefaultOrg.Id, replacedUser.ID)
				require.NoError(ttt, err)
				if !test.PartiallyDeepEqual(tt.want, fetchedUser) {
					ttt.Errorf("GetUser() got = %#v, want %#v", fetchedUser, tt.want)
				}
			}, retryDuration, tick)
		})
	}

}

func TestReplaceUser_removeOldMetadata(t *testing.T) {
	// ensure old metadata is removed correctly
	username := integration.Username()
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
	require.NoError(t, err)

	_, err = Instance.Client.SCIM.Users.Replace(CTX, Instance.DefaultOrg.Id, createdUser.ID, withUsername(minimalUserJson, username))
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		md, err := Instance.Client.Mgmt.ListUserMetadata(CTX, &management.ListUserMetadataRequest{
			Id: createdUser.ID,
		})
		require.NoError(tt, err)
		require.Equal(tt, 1, len(md.Result))

		mdMap := make(map[string]string)
		for i := range md.Result {
			mdMap[md.Result[i].Key] = string(md.Result[i].Value)
		}
		test.AssertMapContains(tt, mdMap, "urn:zitadel:scim:emails", fmt.Sprintf("[{\"value\":\"%s@example.com\",\"primary\":true}]", username))
	}, retryDuration, tick)
}

func TestReplaceUser_emailType(t *testing.T) {
	// ensure old metadata is removed correctly
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, integration.Username()))
	require.NoError(t, err)

	replacedUsername := integration.Username()
	_, err = Instance.Client.SCIM.Users.Replace(CTX, Instance.DefaultOrg.Id, createdUser.ID, withUsername(minimalUserWithEmailTypeReplaceJson, replacedUsername))
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		md, err := Instance.Client.Mgmt.ListUserMetadata(CTX, &management.ListUserMetadataRequest{
			Id: createdUser.ID,
		})
		require.NoError(tt, err)
		require.Equal(tt, 1, len(md.Result))

		mdMap := make(map[string]string)
		for i := range md.Result {
			mdMap[md.Result[i].Key] = string(md.Result[i].Value)
		}

		test.AssertMapContains(tt, mdMap, "urn:zitadel:scim:emails", fmt.Sprintf("[{\"value\":\"%s@example.com\",\"primary\":true,\"type\":\"work\"}]", replacedUsername))
	}, retryDuration, tick)
}

func TestReplaceUser_scopedExternalID(t *testing.T) {
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, integration.Username()))
	require.NoError(t, err)
	callingUserId, callingUserPat, err := Instance.CreateMachineUserPATWithMembership(CTX, "ORG_OWNER")
	require.NoError(t, err)
	ctx := integration.WithAuthorizationToken(CTX, callingUserPat)
	// set provisioning domain of service user
	setProvisioningDomain(t, callingUserId, "fooBazz")

	// replace the user with provisioning domain set
	_, err = Instance.Client.SCIM.Users.Replace(ctx, Instance.DefaultOrg.Id, createdUser.ID, minimalUserWithExternalIDJson)
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		md, err := Instance.Client.Mgmt.ListUserMetadata(ctx, &management.ListUserMetadataRequest{
			Id: createdUser.ID,
		})
		require.NoError(tt, err)

		mdMap := make(map[string]string)
		for i := range md.Result {
			mdMap[md.Result[i].Key] = string(md.Result[i].Value)
		}

		// both external IDs should be present on the user
		test.AssertMapContains(tt, mdMap, "urn:zitadel:scim:externalId", "701984")
		test.AssertMapContains(tt, mdMap, "urn:zitadel:scim:fooBazz:externalId", "replaced-external-id")
	}, retryDuration, tick)
}
