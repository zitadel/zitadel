//go:build integration

package integration_test

import (
	"context"
	_ "embed"
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
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
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
	fullUserCreated, err := Instance.Client.SCIM.Users.Create(CTX, Instance.DefaultOrg.Id, fullUserJson)
	require.NoError(t, err)

	defer func() {
		_, err = Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: fullUserCreated.ID})
		require.NoError(t, err)
	}()

	tests := []struct {
		name          string
		body          []byte
		ctx           context.Context
		orgID         string
		userID        string
		want          *resources.ScimUser
		wantErr       bool
		scimErrorType string
		errorStatus   int
	}{
		{
			name:        "not authenticated",
			ctx:         context.Background(),
			body:        minimalUserUpdateJson,
			wantErr:     true,
			errorStatus: http.StatusUnauthorized,
		},
		{
			name:        "no permissions",
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			body:        minimalUserUpdateJson,
			wantErr:     true,
			errorStatus: http.StatusNotFound,
		},
		{
			name:        "other org",
			orgID:       SecondaryOrganization.OrganizationId,
			body:        minimalUserUpdateJson,
			wantErr:     true,
			errorStatus: http.StatusNotFound,
		},
		{
			name:        "other org with permissions",
			ctx:         Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner),
			orgID:       SecondaryOrganization.OrganizationId,
			body:        minimalUserUpdateJson,
			wantErr:     true,
			errorStatus: http.StatusNotFound,
		},
		{
			name:          "invalid patch json",
			body:          simpleReplacePatchBody("nickname", "10"),
			wantErr:       true,
			scimErrorType: "invalidValue",
		},
		{
			name:          "password complexity violation",
			body:          simpleReplacePatchBody("password", `"fooBar"`),
			wantErr:       true,
			scimErrorType: "invalidValue",
		},
		{
			name:          "invalid profile url",
			body:          simpleReplacePatchBody("profileUrl", `"ftp://example.com/profiles"`),
			wantErr:       true,
			scimErrorType: "invalidValue",
		},
		{
			name:          "invalid time zone",
			body:          simpleReplacePatchBody("timezone", `"foobar"`),
			wantErr:       true,
			scimErrorType: "invalidValue",
		},
		{
			name:          "invalid locale",
			body:          simpleReplacePatchBody("locale", `"foobar"`),
			wantErr:       true,
			scimErrorType: "invalidValue",
		},
		{
			name:        "unknown user id",
			body:        simpleReplacePatchBody("nickname", `"foo"`),
			userID:      "fooBar",
			wantErr:     true,
			errorStatus: http.StatusNotFound,
		},
		{
			name: "full",
			body: fullUserUpdateJson,
			want: &resources.ScimUser{
				ExternalID: "fooBAR",
				UserName:   "bjensen@example.com",
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
						Value:   "babs@example.com",
						Primary: true,
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
				PreferredLanguage: language.MustParse("en-US"),
				Locale:            "en-US",
				Timezone:          "America/Los_Angeles",
				Active:            schemas.NewRelaxedBool(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ctx == nil {
				tt.ctx = CTX
			}

			if tt.orgID == "" {
				tt.orgID = Instance.DefaultOrg.Id
			}

			if tt.userID == "" {
				tt.userID = fullUserCreated.ID
			}

			err := Instance.Client.SCIM.Users.Update(tt.ctx, tt.orgID, tt.userID, tt.body)

			if tt.wantErr {
				require.Error(t, err)

				statusCode := tt.errorStatus
				if statusCode == 0 {
					statusCode = http.StatusBadRequest
				}

				scimErr := scim.RequireScimError(t, statusCode, err)
				assert.Equal(t, tt.scimErrorType, scimErr.Error.ScimType)
				return
			}

			require.NoError(t, err)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				fetchedUser, err := Instance.Client.SCIM.Users.Get(tt.ctx, tt.orgID, fullUserCreated.ID)
				require.NoError(ttt, err)

				fetchedUser.Resource = nil
				fetchedUser.ID = ""
				if tt.want != nil && !test.PartiallyDeepEqual(tt.want, fetchedUser) {
					ttt.Errorf("got = %#v, want = %#v", fetchedUser, tt.want)
				}
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
