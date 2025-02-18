//go:build integration

package integration_test

import (
	"context"
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
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name        string
		orgID       string
		buildUserID func() string
		cleanup     func(userID string)
		ctx         context.Context
		want        *resources.ScimUser
		wantErr     bool
		errorStatus int
	}{
		{
			name:        "not authenticated",
			ctx:         context.Background(),
			errorStatus: http.StatusUnauthorized,
			wantErr:     true,
		},
		{
			name:        "no permissions",
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			errorStatus: http.StatusNotFound,
			wantErr:     true,
		},
		{
			name:        "another org",
			orgID:       SecondaryOrganization.OrganizationId,
			errorStatus: http.StatusNotFound,
			wantErr:     true,
		},
		{
			name:        "another org with permissions",
			orgID:       SecondaryOrganization.OrganizationId,
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			errorStatus: http.StatusNotFound,
			wantErr:     true,
		},
		{
			name: "unknown user id",
			buildUserID: func() string {
				return "unknown"
			},
			errorStatus: http.StatusNotFound,
			wantErr:     true,
		},
		{
			name: "created via grpc",
			want: &resources.ScimUser{
				Name: &resources.ScimUserName{
					FamilyName: "Mouse",
					GivenName:  "Mickey",
				},
				PreferredLanguage: language.MustParse("nl"),
				PhoneNumbers: []*resources.ScimPhoneNumber{
					{
						Value:   "+41791234567",
						Primary: true,
					},
				},
			},
		},
		{
			name: "created via scim",
			buildUserID: func() string {
				createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
				require.NoError(t, err)
				return createdUser.ID
			},
			cleanup: func(userID string) {
				_, err := Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: userID})
				require.NoError(t, err)
			},
			want: &resources.ScimUser{
				ExternalID: "701984",
				UserName:   "bjensen@example.com",
				Name: &resources.ScimUserName{
					Formatted:       "Babs Jensen", // DisplayName takes precedence
					FamilyName:      "Jensen",
					GivenName:       "Barbara",
					MiddleName:      "Jane",
					HonorificPrefix: "Ms.",
					HonorificSuffix: "III",
				},
				DisplayName:       "Babs Jensen",
				NickName:          "Babs",
				ProfileUrl:        test.Must(schemas.ParseHTTPURL("http://login.example.com/bjensen")),
				Title:             "Tour Guide",
				PreferredLanguage: language.Make("en-US"),
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            schemas.NewRelaxedBool(true),
				Emails: []*resources.ScimEmail{
					{
						Value:   "bjensen@example.com",
						Primary: true,
					},
				},
				PhoneNumbers: []*resources.ScimPhoneNumber{
					{
						Value:   "+415555555555",
						Primary: true,
					},
				},
				Ims: []*resources.ScimIms{
					{
						Value: "someaimhandle",
						Type:  "aim",
					},
					{
						Value: "twitterhandle",
						Type:  "X",
					},
				},
				Addresses: []*resources.ScimAddress{
					{
						Type:          "work",
						StreetAddress: "100 Universal City Plaza",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "100 Universal City Plaza\nHollywood, CA 91608 USA",
						Primary:       true,
					},
					{
						Type:          "home",
						StreetAddress: "456 Hollywood Blvd",
						Locality:      "Hollywood",
						Region:        "CA",
						PostalCode:    "91608",
						Country:       "USA",
						Formatted:     "456 Hollywood Blvd\nHollywood, CA 91608 USA",
					},
				},
				Photos: []*resources.ScimPhoto{
					{
						Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/F")),
						Type:  "photo",
					},
					{
						Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/T")),
						Type:  "thumbnail",
					},
				},
				Roles: []*resources.ScimRole{
					{
						Value:   "my-role-1",
						Display: "Rolle 1",
						Type:    "main-role",
						Primary: true,
					},
					{
						Value:   "my-role-2",
						Display: "Rolle 2",
						Type:    "secondary-role",
						Primary: false,
					},
				},
				Entitlements: []*resources.ScimEntitlement{
					{
						Value:   "my-entitlement-1",
						Display: "Entitlement 1",
						Type:    "main-entitlement",
						Primary: true,
					},
					{
						Value:   "my-entitlement-2",
						Display: "Entitlement 2",
						Type:    "secondary-entitlement",
						Primary: false,
					},
				},
			},
		},
		{
			name: "scoped externalID",
			buildUserID: func() string {
				// create user without provisioning domain
				createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
				require.NoError(t, err)

				// set provisioning domain of service user
				setProvisioningDomain(t, Instance.Users.Get(integration.UserTypeOrgOwner).ID, "fooBar")

				// set externalID for provisioning domain
				setAndEnsureMetadata(t, createdUser.ID, "urn:zitadel:scim:fooBar:externalId", "100-scopedExternalId")
				return createdUser.ID
			},
			cleanup: func(userID string) {
				_, err := Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: userID})
				require.NoError(t, err)

				removeProvisioningDomain(t, Instance.Users.Get(integration.UserTypeOrgOwner).ID)
			},
			want: &resources.ScimUser{
				ExternalID: "100-scopedExternalId",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = CTX
			}

			var userID string
			if tt.buildUserID != nil {
				userID = tt.buildUserID()
			} else {
				createUserResp := Instance.CreateHumanUser(CTX)
				userID = createUserResp.UserId
			}

			orgID := tt.orgID
			if orgID == "" {
				orgID = Instance.DefaultOrg.Id
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			var fetchedUser *resources.ScimUser
			var err error
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				fetchedUser, err = Instance.Client.SCIM.Users.Get(ctx, orgID, userID)
				if tt.wantErr {
					statusCode := tt.errorStatus
					if statusCode == 0 {
						statusCode = http.StatusBadRequest
					}

					scim.RequireScimError(ttt, statusCode, err)
					return
				}

				assert.Equal(ttt, userID, fetchedUser.ID)
				assert.EqualValues(ttt, []schemas.ScimSchemaType{"urn:ietf:params:scim:schemas:core:2.0:User"}, fetchedUser.Schemas)
				assert.Equal(ttt, schemas.ScimResourceTypeSingular("User"), fetchedUser.Resource.Meta.ResourceType)
				assert.Equal(ttt, "http://"+Instance.Host()+path.Join(schemas.HandlerPrefix, orgID, "Users", fetchedUser.ID), fetchedUser.Resource.Meta.Location)
				assert.Nil(ttt, fetchedUser.Password)
				if !test.PartiallyDeepEqual(tt.want, fetchedUser) {
					ttt.Errorf("GetUser() got = %#v, want %#v", fetchedUser, tt.want)
				}
			}, retryDuration, tick)

			if tt.cleanup != nil {
				tt.cleanup(fetchedUser.ID)
			}
		})
	}
}
