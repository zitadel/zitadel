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
)

func TestGetUser(t *testing.T) {
	type testCase struct {
		ctx         context.Context
		orgID       string
		userID      string
		want        *resources.ScimUser
		wantErr     bool
		errorStatus int
	}
	tests := []struct {
		name  string
		setup func(t *testing.T) testCase
	}{
		{
			name: "not authenticated",
			setup: func(t *testing.T) testCase {
				return testCase{
					ctx:         context.Background(),
					errorStatus: http.StatusUnauthorized,
					wantErr:     true,
				}
			},
		},
		{
			name: "no permissions",
			setup: func(t *testing.T) testCase {
				return testCase{
					ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
					errorStatus: http.StatusNotFound,
					wantErr:     true,
				}
			},
		},
		{
			name: "another org",
			setup: func(t *testing.T) testCase {
				return testCase{
					orgID:       SecondaryOrganization.OrganizationId,
					errorStatus: http.StatusNotFound,
					wantErr:     true,
				}
			},
		},
		{
			name: "another org with permissions",
			setup: func(t *testing.T) testCase {
				return testCase{
					ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
					orgID:       SecondaryOrganization.OrganizationId,
					errorStatus: http.StatusNotFound,
					wantErr:     true,
				}
			},
		},
		{
			name: "unknown user id",
			setup: func(t *testing.T) testCase {
				return testCase{
					userID:      "unknown",
					errorStatus: http.StatusNotFound,
					wantErr:     true,
				}
			},
		},
		{
			name: "created via grpc",
			setup: func(t *testing.T) testCase {
				return testCase{
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
				}
			},
		},
		{
			name: "created via scim",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID: createdUser.ID,
					want: &resources.ScimUser{
						ExternalID: "701984",
						UserName:   username + "@example.com",
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
								Value:   username + "@example.com",
								Primary: true,
								Type:    "work",
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
				}
			},
		},
		{
			name: "scoped externalID",
			setup: func(t *testing.T) testCase {
				createdUser, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, integration.Username()))
				require.NoError(t, err)
				callingUserId, callingUserPat, err := Instance.CreateMachineUserPATWithMembership(CTX, "ORG_OWNER")
				require.NoError(t, err)
				setProvisioningDomain(t, callingUserId, "fooBar")
				setAndEnsureMetadata(t, createdUser.ID, "urn:zitadel:scim:fooBar:externalId", "100-scopedExternalId")
				return testCase{
					ctx:    integration.WithAuthorizationToken(CTX, callingUserPat),
					userID: createdUser.ID,
					want: &resources.ScimUser{
						ExternalID: "100-scopedExternalId",
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttt := tt.setup(t)
			if ttt.userID == "" {
				ttt.userID = Instance.CreateHumanUser(CTX).UserId
			}
			if ttt.ctx == nil {
				ttt.ctx = CTX
			}
			if ttt.orgID == "" {
				ttt.orgID = Instance.DefaultOrg.Id
			}
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(collect *assert.CollectT) {
				fetchedUser, err := Instance.Client.SCIM.Users.Get(ttt.ctx, ttt.orgID, ttt.userID)
				if ttt.wantErr {
					statusCode := ttt.errorStatus
					if statusCode == 0 {
						statusCode = http.StatusBadRequest
					}
					scim.RequireScimError(collect, statusCode, err)
					return
				}
				if !assert.NoError(collect, err) {
					scim.RequireScimError(collect, http.StatusNotFound, err)
					return
				}
				assert.Equal(collect, ttt.userID, fetchedUser.ID)
				assert.EqualValues(collect, []schemas.ScimSchemaType{"urn:ietf:params:scim:schemas:core:2.0:User"}, fetchedUser.Schemas)
				assert.Equal(collect, schemas.ScimResourceTypeSingular("User"), fetchedUser.Resource.Meta.ResourceType)
				assert.Equal(collect, "http://"+Instance.Host()+path.Join(schemas.HandlerPrefix, ttt.orgID, "Users", fetchedUser.ID), fetchedUser.Resource.Meta.Location)
				assert.Nil(collect, fetchedUser.Password)
				if !test.PartiallyDeepEqual(ttt.want, fetchedUser) {
					collect.Errorf("GetUser() got = %#v, want %#v", fetchedUser, ttt.want)
				}
			}, retryDuration, tick)
		})
	}
}
