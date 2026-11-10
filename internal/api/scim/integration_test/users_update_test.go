//go:build integration

package integration_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
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

var (
	//go:embed testdata/users_update_test_full.json
	fullUserUpdateJson []byte

	minimalUserUpdateJson = simpleReplacePatchBody("nickname", "\"foo\"")
)

func init() {
	fullUserUpdateJson = removeComments(fullUserUpdateJson)
}

func TestUpdateUser(t *testing.T) {
	type testCase struct {
		name          string
		body          []byte
		ctx           context.Context
		orgID         string
		userID        string
		want          *resources.ScimUser
		wantErr       bool
		scimErrorType string
		errorStatus   int
	}
	tests := []struct {
		name  string
		setup func(t *testing.T) testCase
	}{
		{
			name: "not authenticated",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:      created.ID,
					ctx:         context.Background(),
					body:        minimalUserUpdateJson,
					wantErr:     true,
					errorStatus: http.StatusUnauthorized,
				}
			},
		},
		{
			name: "no permissions",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:      created.ID,
					ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
					body:        minimalUserUpdateJson,
					wantErr:     true,
					errorStatus: http.StatusNotFound,
				}
			},
		},
		{
			name: "other org",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:      created.ID,
					orgID:       SecondaryOrganization.OrganizationId,
					body:        minimalUserUpdateJson,
					wantErr:     true,
					errorStatus: http.StatusNotFound,
				}
			},
		},
		{
			name: "other org with permissions",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:      created.ID,
					ctx:         Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner),
					orgID:       SecondaryOrganization.OrganizationId,
					body:        minimalUserUpdateJson,
					wantErr:     true,
					errorStatus: http.StatusNotFound,
				}
			},
		},
		{
			name: "invalid patch json",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:        created.ID,
					body:          simpleReplacePatchBody("nickname", "10"),
					wantErr:       true,
					scimErrorType: "invalidValue",
				}
			},
		},
		{
			name: "password complexity violation",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:        created.ID,
					body:          simpleReplacePatchBody("password", `"fooBar"`),
					wantErr:       true,
					scimErrorType: "invalidValue",
				}
			},
		},
		{
			name: "invalid profile url",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:        created.ID,
					body:          simpleReplacePatchBody("profileUrl", `"ftp://example.com/profiles"`),
					wantErr:       true,
					scimErrorType: "invalidValue",
				}
			},
		},
		{
			name: "invalid time zone",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:        created.ID,
					body:          simpleReplacePatchBody("timezone", `"foobar"`),
					wantErr:       true,
					scimErrorType: "invalidValue",
				}
			},
		},
		{
			name: "invalid locale",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID:        created.ID,
					body:          simpleReplacePatchBody("locale", `"foobar"`),
					wantErr:       true,
					scimErrorType: "invalidValue",
				}
			},
		},
		{
			name: "unknown user id",
			setup: func(t *testing.T) testCase {
				return testCase{
					body:        simpleReplacePatchBody("nickname", `"foo"`),
					userID:      "fooBar",
					wantErr:     true,
					errorStatus: http.StatusNotFound,
				}
			},
		},
		{
			name: "full",
			setup: func(t *testing.T) testCase {
				username := integration.Username()
				created, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, withUsername(fullUserJson, username))
				require.NoError(t, err)
				return testCase{
					userID: created.ID,
					body:   withUsername(fullUserUpdateJson, username),
					want: &resources.ScimUser{
						ExternalID: "fooBAR",
						UserName:   username + "@example.com",
						Name: &resources.ScimUserName{
							Formatted:       "replaced-display-name",
							FamilyName:      "added-family-name",
							GivenName:       "added-given-name",
							MiddleName:      "added-middle-name-2",
							HonorificPrefix: "added-honorific-prefix",
							HonorificSuffix: "replaced-honorific-suffix",
						},
						DisplayName: "replaced-display-name",
						NickName:    "",
						ProfileUrl:  test.Must(schemas.ParseHTTPURL("http://login.example.com/bjensen")),
						Emails: []*resources.ScimEmail{
							{
								Value: username + "@example.com",
								Type:  "work",
							},
							{
								Value: username + "+1@example.com",
								Type:  "home",
							},
							{
								Value:   username + "+2@example.com",
								Primary: true,
								Type:    "home",
							},
						},
						Addresses: []*resources.ScimAddress{
							{
								Type:          "replaced-work",
								StreetAddress: "replaced-100 Universal City Plaza",
								Locality:      "replaced-Hollywood",
								Region:        "replaced-CA",
								PostalCode:    "replaced-91608",
								Country:       "replaced-USA",
								Formatted:     "replaced-100 Universal City Plaza\nHollywood, CA 91608 USA",
								Primary:       true,
							},
						},
						PhoneNumbers: []*resources.ScimPhoneNumber{
							{
								Value:   "+41711234567",
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
								Type:  "",
							},
						},
						Photos: []*resources.ScimPhoto{
							{
								Value: *test.Must(schemas.ParseHTTPURL("https://photos.example.com/profilephoto/72930000000Ccne/F")),
								Type:  "photo",
							},
						},
						Roles: nil,
						Entitlements: []*resources.ScimEntitlement{
							{
								Value:   "my-entitlement-1",
								Display: "added-entitlement-1",
								Type:    "added-entitlement-1",
								Primary: false,
							},
							{
								Value:   "my-entitlement-2",
								Display: "Entitlement 2",
								Type:    "secondary-entitlement",
								Primary: false,
							},
							{
								Value:   "added-entitlement-1",
								Primary: false,
							},
							{
								Value:   "added-entitlement-2",
								Primary: false,
							},
							{
								Value:   "added-entitlement-3",
								Primary: true,
							},
						},
						Title:             "Tour Guide",
						PreferredLanguage: language.MustParse("en"),
						Locale:            "en-US",
						Timezone:          "America/Los_Angeles",
						Active:            schemas.NewRelaxedBool(true),
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttt := tt.setup(t)
			if ttt.orgID == "" {
				ttt.orgID = Instance.DefaultOrg.Id
			}
			if ttt.ctx == nil {
				ttt.ctx = CTX
			}
			err := Instance.Client.SCIM.Users.Update(ttt.ctx, ttt.orgID, ttt.userID, ttt.body)
			if ttt.wantErr {
				require.Error(t, err)
				statusCode := ttt.errorStatus
				if statusCode == 0 {
					statusCode = http.StatusBadRequest
				}
				scimErr := scim.RequireScimError(t, statusCode, err)
				assert.Equal(t, ttt.scimErrorType, scimErr.Error.ScimType)
				return
			}

			require.NoError(t, err)
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(collect *assert.CollectT) {
				fetchedUser, err := Instance.Client.SCIM.Users.Get(ttt.ctx, ttt.orgID, ttt.userID)
				if !assert.NoError(collect, err) {
					return
				}
				fetchedUser.Resource = nil
				fetchedUser.ID = ""
				fetched, err := json.Marshal(fetchedUser)
				require.NoError(collect, err)
				want, err := json.Marshal(ttt.want)
				require.NoError(collect, err)
				assert.JSONEq(collect, string(want), string(fetched))
			}, retryDuration, tick)
		})
	}
}

func simpleReplacePatchBody(path, value string) []byte {
	return []byte(fmt.Sprintf(
		`{
		  "schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
		  "Operations": [
			{
			  "op": "replace",
			  "path": "%s",
			  "value": %s
			}
		  ]
		}`,
		path,
		value,
	))
}
