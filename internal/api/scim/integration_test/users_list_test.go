//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var totalCountOfHumanUsers = 13

/*
	func TestListUser(t *testing.T) {
		createdUserIDs := createUsers(t, CTX, Instance.DefaultOrg.Id)
		// secondary organization with same set of users,
		// these should never be modified.
		// This allows testing list requests without filters.
		iamOwnerCtx := Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
		secondaryOrg := Instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
		secondaryOrgCreatedUserIDs := createUsers(t, iamOwnerCtx, secondaryOrg.OrganizationId)

		testsInitializedUtc := time.Now().UTC()

		// Wait one second to ensure a change in the least significant value of the timestamp.
		time.Sleep(time.Second)

		tests := []struct {
			name        string
			ctx         context.Context
			orgID       string
			req         *scim.ListRequest
			prepare     func(require.TestingT) *scim.ListRequest
			wantErr     bool
			errorStatus int
			errorType   string
			assert      func(assert.TestingT, *scim.ListResponse[*resources.ScimUser])
			cleanup     func(require.TestingT)
		}{
			{
				name:        "not authenticated",
				ctx:         context.Background(),
				req:         new(scim.ListRequest),
				wantErr:     true,
				errorStatus: http.StatusUnauthorized,
			},
			{
				name:        "no permissions",
				ctx:         Instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
				req:         new(scim.ListRequest),
				wantErr:     true,
				errorStatus: http.StatusNotFound,
			},
			{
				name: "unknown sort order",
				req: &scim.ListRequest{
					SortBy:    gu.Ptr("id"),
					SortOrder: gu.Ptr(scim.ListRequestSortOrder("fooBar")),
				},
				wantErr:   true,
				errorType: "invalidValue",
			},
			{
				name: "unknown sort field",
				req: &scim.ListRequest{
					SortBy: gu.Ptr("fooBar"),
				},
				wantErr:   true,
				errorType: "invalidValue",
			},
			{
				name: "custom sort field",
				req: &scim.ListRequest{
					SortBy: gu.Ptr("externalid"),
				},
				wantErr:   true,
				errorType: "invalidValue",
			},
			{
				name: "unknown filter field",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`fooBar eq "10"`),
				},
				wantErr:   true,
				errorType: "invalidFilter",
			},
			{
				name: "invalid filter",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`fooBarBaz`),
				},
				wantErr:   true,
				errorType: "invalidFilter",
			},
			{
				name: "list users without filter",
				// use other org, modifications of users happens only on primary org
				orgID: secondaryOrg.OrganizationId,
				ctx:   iamOwnerCtx,
				req:   new(scim.ListRequest),
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 100, resp.ItemsPerPage)
					assert.Equal(t, totalCountOfHumanUsers, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, totalCountOfHumanUsers)
				},
			},
			{
				name: "list paged sorted users without filter",
				// use other org, modifications of users happens only on primary org
				orgID: secondaryOrg.OrganizationId,
				ctx:   iamOwnerCtx,
				req: &scim.ListRequest{
					Count:      gu.Ptr(2),
					StartIndex: gu.Ptr(5),
					SortOrder:  gu.Ptr(scim.ListRequestSortOrderAsc),
					SortBy:     gu.Ptr("username"),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					// sort the created users with usernames instead of creation date
					sortedResources := sortScimUserByUsername(resp.Resources)

					assert.Equal(t, 2, resp.ItemsPerPage)
					assert.Equal(t, totalCountOfHumanUsers, resp.TotalResults)
					assert.Equal(t, 5, resp.StartIndex)
					assert.Len(t, sortedResources, 2)
					assert.True(t, strings.HasPrefix(sortedResources[0].UserName, "scim-username-1: "), "got %q", resp.Resources[0].UserName)
					assert.True(t, strings.HasPrefix(sortedResources[1].UserName, "scim-username-2: "), "got %q", resp.Resources[1].UserName)
				},
			},
			{
				name: "list users with simple filter",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`username sw "scim-username-1"`),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 100, resp.ItemsPerPage)
					assert.Equal(t, 2, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, 2)
					for _, resource := range resp.Resources {
						assert.True(t, strings.HasPrefix(resource.UserName, "scim-username-1"))
					}
				},
			},
			{
				name: "list paged sorted users with filter",
				req: &scim.ListRequest{
					Count:      gu.Ptr(5),
					StartIndex: gu.Ptr(1),
					SortOrder:  gu.Ptr(scim.ListRequestSortOrderAsc),
					SortBy:     gu.Ptr("username"),
					Filter:     gu.Ptr(`emails sw "scim-email-1" and emails ew "@example.com"`),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					// sort the created users with usernames instead of creation date
					sortedResources := sortScimUserByUsername(resp.Resources)

					assert.Equal(t, 5, resp.ItemsPerPage)
					assert.Equal(t, 2, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, sortedResources, 2)
					for _, resource := range sortedResources {
						assert.True(t, strings.HasPrefix(resource.UserName, "scim-username-1"))
						assert.Len(t, resource.Emails, 1)
						assert.True(t, strings.HasPrefix(resource.Emails[0].Value, "scim-email-1"))
						assert.True(t, strings.HasSuffix(resource.Emails[0].Value, "@example.com"))
					}
				},
			},
			{
				name: "list paged sorted users with filter as post",
				req: &scim.ListRequest{
					Schemas:    []schemas.ScimSchemaType{schemas.IdSearchRequest},
					Count:      gu.Ptr(5),
					StartIndex: gu.Ptr(1),
					SortOrder:  gu.Ptr(scim.ListRequestSortOrderAsc),
					SortBy:     gu.Ptr("username"),
					Filter:     gu.Ptr(`emails sw "scim-email-1" and emails ew "@example.com"`),
					SendAsPost: true,
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					// sort the created users with usernames instead of creation date
					sortedResources := sortScimUserByUsername(resp.Resources)

					assert.Equal(t, 5, resp.ItemsPerPage)
					assert.Equal(t, 2, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, sortedResources, 2)
					for _, resource := range sortedResources {
						assert.True(t, strings.HasPrefix(resource.UserName, "scim-username-1"))
						assert.Len(t, resource.Emails, 1)
						assert.True(t, strings.HasPrefix(resource.Emails[0].Value, "scim-email-1"))
						assert.True(t, strings.HasSuffix(resource.Emails[0].Value, "@example.com"))
					}
				},
			},
			{
				name: "count users without filter",
				// use other org, modifications of users happens only on primary org
				orgID: secondaryOrg.OrganizationId,
				ctx:   iamOwnerCtx,
				prepare: func(t require.TestingT) *scim.ListRequest {
					return &scim.ListRequest{
						Count: gu.Ptr(0),
					}
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 0, resp.ItemsPerPage)
					assert.Equal(t, totalCountOfHumanUsers, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, 0)
				},
			},
			{
				name: "list users with active filter",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`active eq false`),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 100, resp.ItemsPerPage)
					assert.Equal(t, 1, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, 1)
					assert.True(t, strings.HasPrefix(resp.Resources[0].UserName, "scim-username-0"))
					assert.False(t, resp.Resources[0].Active.Bool())
				},
			},
			{
				name: "list users with externalid filter",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`externalid eq "701984"`),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 100, resp.ItemsPerPage)
					assert.Equal(t, 1, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, 1)
					assert.Equal(t, resp.Resources[0].ExternalID, "701984")
				},
			},
			{
				name: "list users with externalid filter invalid operator",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`externalid pr`),
				},
				wantErr:   true,
				errorType: "invalidFilter",
			},
			{
				name: "list users with externalid complex filter",
				req: &scim.ListRequest{
					Filter: gu.Ptr(`externalid eq "701984" and username eq "bjensen@example.com"`),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 100, resp.ItemsPerPage)
					assert.Equal(t, 1, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, 1)
					assert.Equal(t, resp.Resources[0].UserName, "bjensen@example.com")
					assert.Equal(t, resp.Resources[0].ExternalID, "701984")
				},
			},
			{
				name: "count users with filter",
				req: &scim.ListRequest{
					Count:  gu.Ptr(0),
					Filter: gu.Ptr(`emails sw "scim-email-1" and emails ew "@example.com"`),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Equal(t, 0, resp.ItemsPerPage)
					assert.Equal(t, 2, resp.TotalResults)
					assert.Equal(t, 1, resp.StartIndex)
					assert.Len(t, resp.Resources, 0)
				},
			},
			{
				name: "list users with modification date filter",
				prepare: func(t require.TestingT) *scim.ListRequest {
					userID := createdUserIDs[len(createdUserIDs)-1] // use the last entry, as we use the others for other assertions
					_, err := Instance.Client.UserV2.UpdateHumanUser(CTX, &user_v2.UpdateHumanUserRequest{
						UserId: userID,

						Profile: &user_v2.SetHumanProfile{
							GivenName:  "scim-user-given-name-modified-0: " + integration.FirstName(),
							FamilyName: "scim-user-family-name-modified-0: " + integration.LastName(),
						},
					})
					require.NoError(t, err)

					return &scim.ListRequest{
						// filter by id too to exclude other random users
						Filter: gu.Ptr(fmt.Sprintf(`id eq "%s" and meta.LASTMODIFIED gt "%s"`, userID, testsInitializedUtc.Format(time.RFC3339))),
					}
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Len(t, resp.Resources, 1)
					assert.Equal(t, resp.Resources[0].ID, createdUserIDs[len(createdUserIDs)-1])
					assert.True(t, strings.HasPrefix(resp.Resources[0].Name.FamilyName, "scim-user-family-name-modified-0:"))
					assert.True(t, strings.HasPrefix(resp.Resources[0].Name.GivenName, "scim-user-given-name-modified-0:"))
				},
			},
			{
				name: "list users with creation date filter",
				prepare: func(t require.TestingT) *scim.ListRequest {
					resp := createHumanUser(t, CTX, Instance.DefaultOrg.Id, 100)
					return &scim.ListRequest{
						Filter: gu.Ptr(fmt.Sprintf(`id eq "%s" and meta.created gt "%s"`, resp.UserId, testsInitializedUtc.Format(time.RFC3339))),
					}
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Len(t, resp.Resources, 1)
					assert.True(t, strings.HasPrefix(resp.Resources[0].UserName, "scim-username-100:"))
				},
			},
			{
				name: "validate returned objects",
				req: &scim.ListRequest{
					Filter: gu.Ptr(fmt.Sprintf(`id eq "%s"`, createdUserIDs[0])),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Len(t, resp.Resources, 1)
					if !test.PartiallyDeepEqual(fullUser, resp.Resources[0]) {
						t.Errorf("got = %#v, want %#v", resp.Resources[0], fullUser)
					}
				},
			},
			{
				name: "do not return user of other org",
				req: &scim.ListRequest{
					Filter: gu.Ptr(fmt.Sprintf(`id eq "%s"`, secondaryOrgCreatedUserIDs[0])),
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Len(t, resp.Resources, 0)
				},
			},
			{
				name: "do not count user of other org",
				prepare: func(t require.TestingT) *scim.ListRequest {
					iamOwnerCtx := Instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
					org := Instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp := createHumanUser(t, iamOwnerCtx, org.OrganizationId, 102)

					return &scim.ListRequest{
						Count:  gu.Ptr(0),
						Filter: gu.Ptr(fmt.Sprintf(`id eq "%s"`, resp.UserId)),
					}
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Len(t, resp.Resources, 0)
				},
			},
			{
				name: "scoped externalID",
				prepare: func(t require.TestingT) *scim.ListRequest {
					resp := createHumanUser(t, CTX, Instance.DefaultOrg.Id, 102)

					// set provisioning domain of service user
					setProvisioningDomain(t, Instance.Users.Get(integration.UserTypeOrgOwner).ID, "fooBar")

					// set externalID for provisioning domain
					setAndEnsureMetadata(t, resp.UserId, "urn:zitadel:scim:fooBar:externalId", "100-scopedExternalId")
					return &scim.ListRequest{
						Filter: gu.Ptr(fmt.Sprintf(`id eq "%s"`, resp.UserId)),
					}
				},
				assert: func(t assert.TestingT, resp *scim.ListResponse[*resources.ScimUser]) {
					assert.Len(t, resp.Resources, 1)
					assert.Equal(t, resp.Resources[0].ExternalID, "100-scopedExternalId")
				},
				cleanup: func(t require.TestingT) {
					// delete provisioning domain of service user
					removeProvisioningDomain(t, Instance.Users.Get(integration.UserTypeOrgOwner).ID)
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.ctx == nil {
					tt.ctx = CTX
				}

				if tt.prepare != nil {
					tt.req = tt.prepare(t)
				}

				if tt.orgID == "" {
					tt.orgID = Instance.DefaultOrg.Id
				}

				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					listResp, err := Instance.Client.SCIM.Users.List(tt.ctx, tt.orgID, tt.req)
					if tt.wantErr {
						statusCode := tt.errorStatus
						if statusCode == 0 {
							statusCode = http.StatusBadRequest
						}

						scimErr := scim.RequireScimError(ttt, statusCode, err)
						if tt.errorType != "" {
							assert.Equal(t, tt.errorType, scimErr.Error.ScimType)
						}
						return
					}

					require.NoError(t, err)
					assert.EqualValues(ttt, []schemas.ScimSchemaType{"urn:ietf:params:scim:api:messages:2.0:ListResponse"}, listResp.Schemas)
					if tt.assert != nil {
						tt.assert(ttt, listResp)
					}
				}, retryDuration, tick)

				if tt.cleanup != nil {
					tt.cleanup(t)
				}
			})
		}
	}
*/
func sortScimUserByUsername(users []*resources.ScimUser) []*resources.ScimUser {
	sortedResources := users
	slices.SortFunc(sortedResources, func(a, b *resources.ScimUser) int {
		return strings.Compare(a.UserName, b.UserName)
	})
	return sortedResources
}

func createUsers(t *testing.T, ctx context.Context, orgID string) []string {
	count := totalCountOfHumanUsers - 1 // zitadel admin is always created by default
	createdUserIDs := make([]string, 0, count)

	// create the full scim user if on primary org
	if orgID == Instance.DefaultOrg.Id {
		fullUserCreatedResp, err := Instance.Client.SCIM.Users.Create(ctx, orgID, withUsername(fullUserJson, integration.Username()))
		require.NoError(t, err)
		createdUserIDs = append(createdUserIDs, fullUserCreatedResp.ID)
		count--
	}

	// set the first user inactive
	resp := createHumanUser(t, ctx, orgID, 0)
	_, err := Instance.Client.UserV2.DeactivateUser(ctx, &user_v2.DeactivateUserRequest{
		UserId: resp.UserId,
	})
	require.NoError(t, err)
	createdUserIDs = append(createdUserIDs, resp.UserId)

	for i := 1; i < count; i++ {
		resp = createHumanUser(t, ctx, orgID, i)
		createdUserIDs = append(createdUserIDs, resp.UserId)
	}

	return createdUserIDs
}

func createHumanUser(t require.TestingT, ctx context.Context, orgID string, i int) *user_v2.AddHumanUserResponse {
	// create remaining minimal users with faker data
	// no need to clean these up as identification attributes change each time
	resp, err := Instance.Client.UserV2.AddHumanUser(ctx, &user_v2.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{
				OrgId: orgID,
			},
		},
		Username: gu.Ptr(fmt.Sprintf("scim-username-%d: %s", i, integration.Username())),
		Profile: &user_v2.SetHumanProfile{
			GivenName:         fmt.Sprintf("scim-givenname-%d: %s", i, integration.FirstName()),
			FamilyName:        fmt.Sprintf("scim-familyname-%d: %s", i, integration.LastName()),
			PreferredLanguage: gu.Ptr("en-US"),
			Gender:            gu.Ptr(user_v2.Gender_GENDER_MALE),
		},
		Email: &user_v2.SetHumanEmail{
			Email: fmt.Sprintf("scim-email-%d-%d@example.com", i, integration.Number()),
		},
	})
	require.NoError(t, err)
	return resp
}
