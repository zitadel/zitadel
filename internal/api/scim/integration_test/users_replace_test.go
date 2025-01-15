//go:build integration

package integration_test

import (
	"context"
	_ "embed"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/scim"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"golang.org/x/text/language"
	"net/http"
	"path"
	"testing"
	"time"
)

var (
	//go:embed testdata/users_replace_test_minimal_with_external_id.json
	minimalUserWithExternalIDJson []byte

	//go:embed testdata/users_replace_test_minimal.json
	minimalUserReplaceJson []byte

	//go:embed testdata/users_replace_test_full.json
	fullUserReplaceJson []byte
)

func TestReplaceUser(t *testing.T) {
	tests := []struct {
		name          string
		body          []byte
		ctx           context.Context
		want          *resources.ScimUser
		wantErr       bool
		scimErrorType string
		errorStatus   int
		zitadelErrID  string
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
				ProfileUrl:  integration.Must(schemas.ParseHTTPURL("http://login.example.com/bjensen-updated")),
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
						Value: *integration.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/F-updated")),
						Type:  "photo-updated",
					},
					{
						Value: *integration.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/T-updated")),
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
				Active:            gu.Ptr(false),
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
			createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
			require.NoError(t, err)

			defer func() {
				_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
				assert.NoError(t, err)
			}()

			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			replacedUser, err := Instance.Client.SCIM.Users.Replace(ctx, Instance.DefaultOrg.Id, createdUser.ID, tt.body)
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

			if !integration.PartiallyDeepEqual(tt.want, replacedUser) {
				t.Errorf("ReplaceUser() got = %#v, want %#v", replacedUser, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// ensure the user is really stored and not just returned to the caller
				fetchedUser, err := Instance.Client.SCIM.Users.Get(CTX, Instance.DefaultOrg.Id, replacedUser.ID)
				require.NoError(ttt, err)
				if !integration.PartiallyDeepEqual(tt.want, fetchedUser) {
					ttt.Errorf("GetUser() got = %#v, want %#v", fetchedUser, tt.want)
				}
			}, retryDuration, tick)
		})
	}

}

func TestReplaceUser_removeOldMetadata(t *testing.T) {
	// ensure old metadata is removed correctly
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
	require.NoError(t, err)

	_, err = Instance.Client.SCIM.Users.Replace(CTX, Instance.DefaultOrg.Id, createdUser.ID, minimalUserJson)
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		md, err := Instance.Client.Mgmt.ListUserMetadata(CTX, &management.ListUserMetadataRequest{
			Id: createdUser.ID,
		})
		require.NoError(tt, err)
		require.Equal(tt, 0, len(md.Result))
	}, retryDuration, tick)

	_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
	require.NoError(t, err)
}

func TestReplaceUser_scopedExternalID(t *testing.T) {
	// create user without provisioning domain set
	createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
	require.NoError(t, err)

	// set provisioning domain of service user
	_, err = Instance.Client.Mgmt.SetUserMetadata(CTX, &management.SetUserMetadataRequest{
		Id:    Instance.Users.Get(integration.UserTypeOrgOwner).ID,
		Key:   "urn:zitadel:scim:provisioning_domain",
		Value: []byte("fooBazz"),
	})
	require.NoError(t, err)

	// replace the user with provisioning domain set
	_, err = Instance.Client.SCIM.Users.Replace(CTX, Instance.DefaultOrg.Id, createdUser.ID, minimalUserWithExternalIDJson)
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		md, err := Instance.Client.Mgmt.ListUserMetadata(CTX, &management.ListUserMetadataRequest{
			Id: createdUser.ID,
		})
		require.NoError(tt, err)

		mdMap := make(map[string]string)
		for i := range md.Result {
			mdMap[md.Result[i].Key] = string(md.Result[i].Value)
		}

		// both external IDs should be present on the user
		integration.AssertMapContains(tt, mdMap, "urn:zitadel:scim:externalId", "701984")
		integration.AssertMapContains(tt, mdMap, "urn:zitadel:scim:fooBazz:externalId", "replaced-external-id")
	}, retryDuration, tick)

	_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: createdUser.ID})
	require.NoError(t, err)

	_, err = Instance.Client.Mgmt.RemoveUserMetadata(CTX, &management.RemoveUserMetadataRequest{
		Id:  Instance.Users.Get(integration.UserTypeOrgOwner).ID,
		Key: "urn:zitadel:scim:provisioning_domain",
	})
	require.NoError(t, err)
}
