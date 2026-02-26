//go:build integration

package org_test

import (
	"context"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX               context.Context
	Instance          *integration.Instance
	Client            org.OrganizationServiceClient
	User              *user.AddHumanUserResponse
	OtherOrganization *org.AddOrganizationResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OrgV2

		CTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		User = Instance.CreateHumanUser(CTX)
		OtherOrganization = Instance.CreateOrganization(CTX, integration.OrganizationName(), integration.Email())
		return m.Run()
	}())
}

func TestServer_AddOrganization(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)
	userId := "userID"

	tests := []struct {
		name    string
		ctx     context.Context
		req     *org.AddOrganizationRequest
		want    *org.AddOrganizationResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &org.AddOrganizationRequest{
				Name:   "name",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name:   "",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid admin type",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "no admin, custom org ID",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name:  integration.OrganizationName(),
				OrgId: gu.Ptr("custom-org-ID"),
			},
			want: &org.AddOrganizationResponse{
				OrganizationId: "custom-org-ID",
				CreatedAdmins:  []*org.AddOrganizationResponse_CreatedAdmin{},
			},
		},
		{
			name: "no admin, custom organization ID",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name:           integration.OrganizationName(),
				OrganizationId: gu.Ptr("custom-organization-ID"),
			},
			want: &org.AddOrganizationResponse{
				OrganizationId: "custom-organization-ID",
				CreatedAdmins:  []*org.AddOrganizationResponse_CreatedAdmin{},
			},
		},
		{
			name: "no admin, custom organization ID (precedence over org ID)",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name:           integration.OrganizationName(),
				OrganizationId: gu.Ptr("custom-organization-ID2"),
				OrgId:          gu.Ptr("custom-org-ID"), // will be ignored in favor of OrganizationId
			},
			want: &org.AddOrganizationResponse{
				OrganizationId: "custom-organization-ID2",
				CreatedAdmins:  []*org.AddOrganizationResponse_CreatedAdmin{},
			},
		},
		{
			name: "admin with init with userID passed for Human admin",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{
						UserType: &org.AddOrganizationRequest_Admin_Human{
							Human: &user.AddHumanUserRequest{
								UserId: &userId,
								Profile: &user.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user.SetHumanEmail{
									Email: integration.Email(),
									Verification: &user.SetHumanEmail_ReturnCode{
										ReturnCode: &user.ReturnEmailVerificationCode{},
									},
								},
							},
						},
					},
				},
			},
			want: &org.AddOrganizationResponse{
				OrganizationId: integration.NotEmpty,
				CreatedAdmins: []*org.AddOrganizationResponse_CreatedAdmin{
					{
						UserId:    userId,
						EmailCode: gu.Ptr(integration.NotEmpty),
						PhoneCode: nil,
					},
				},
			},
		},
		{
			name: "existing user and new human with idp",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{
						UserType: &org.AddOrganizationRequest_Admin_UserId{UserId: User.GetUserId()},
					},
					{
						UserType: &org.AddOrganizationRequest_Admin_Human{
							Human: &user.AddHumanUserRequest{
								Profile: &user.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user.SetHumanEmail{
									Email: integration.Email(),
									Verification: &user.SetHumanEmail_IsVerified{
										IsVerified: true,
									},
								},
								IdpLinks: []*user.IDPLink{
									{
										IdpId:    idpResp.Id,
										UserId:   userId,
										UserName: "username",
									},
								},
							},
						},
					},
				},
			},
			want: &org.AddOrganizationResponse{
				CreatedAdmins: []*org.AddOrganizationResponse_CreatedAdmin{
					// a single admin is expected, because the first provided already exists
					{
						UserId: integration.NotEmpty,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOrganization(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check details
			assert.NotZero(t, got.GetDetails().GetSequence())
			gotCD := got.GetDetails().GetChangeDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			assert.NotEmpty(t, got.GetDetails().GetResourceOwner())

			// organization id must be the same as the resourceOwner
			assert.Equal(t, got.GetDetails().GetResourceOwner(), got.GetOrganizationId())

			// check the admins
			require.Len(t, got.GetCreatedAdmins(), len(tt.want.GetCreatedAdmins()))
			for i, admin := range tt.want.GetCreatedAdmins() {
				gotAdmin := got.GetCreatedAdmins()[i]
				assertCreatedAdmin(t, admin, gotAdmin)
			}
		})
	}
}

func assertCreatedAdmin(t *testing.T, expected, got *org.AddOrganizationResponse_CreatedAdmin) {
	if expected.GetUserId() != "" {
		assert.NotEmpty(t, got.GetUserId())
	} else {
		assert.Empty(t, got.GetUserId())
	}
	if expected.GetEmailCode() != "" {
		assert.NotEmpty(t, got.GetEmailCode())
	} else {
		assert.Empty(t, got.GetEmailCode())
	}
	if expected.GetPhoneCode() != "" {
		assert.NotEmpty(t, got.GetPhoneCode())
	} else {
		assert.Empty(t, got.GetPhoneCode())
	}
}

func TestServer_UpdateOrganization(t *testing.T) {
	orgs, orgsName, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].OrganizationId
	orgName := orgsName[0]

	tests := []struct {
		name    string
		ctx     context.Context
		req     *org.UpdateOrganizationRequest
		want    *org.UpdateOrganizationResponse
		wantErr bool
	}{
		{
			name: "update org with new name",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.UpdateOrganizationRequest{
				OrganizationId: orgId,
				Name:           "new org name",
			},
		},
		{
			name: "update org with same name",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.UpdateOrganizationRequest{
				OrganizationId: orgId,
				Name:           orgName,
			},
		},
		{
			name: "update org with non existent org id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.UpdateOrganizationRequest{
				OrganizationId: "non existant org id",
				// Name: "",
			},
			wantErr: true,
		},
		{
			name: "update org with no id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.UpdateOrganizationRequest{
				OrganizationId: "",
				Name:           orgName,
			},
			wantErr: true,
		},
		{
			name: "no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &org.UpdateOrganizationRequest{
				OrganizationId: OtherOrganization.GetOrganizationId(),
				Name:           integration.OrganizationName(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateOrganization(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check details
			gotCD := got.GetChangeDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
		})
	}
}

func TestServer_DeleteOrganization(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		createOrgFunc func() string
		req           *org.DeleteOrganizationRequest
		want          *org.DeleteOrganizationResponse
		dontCheckTime bool
		err           error
	}{
		{
			name: "delete org no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			createOrgFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				return orgs[0].OrganizationId
			},
			req: &org.DeleteOrganizationRequest{},
			err: errors.New("membership not found"),
		},
		{
			name: "delete org happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			createOrgFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				return orgs[0].OrganizationId
			},
			req: &org.DeleteOrganizationRequest{},
		},
		{
			name: "delete already deleted org",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			createOrgFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				// delete org
				_, err := Client.DeleteOrganization(CTX, &org.DeleteOrganizationRequest{OrganizationId: orgs[0].OrganizationId})
				require.NoError(t, err)

				return orgs[0].OrganizationId
			},
			req:           &org.DeleteOrganizationRequest{},
			dontCheckTime: true,
		},
		{
			name: "delete non existent org",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.DeleteOrganizationRequest{
				OrganizationId: "non existent org id",
			},
			dontCheckTime: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createOrgFunc != nil {
				tt.req.OrganizationId = tt.createOrgFunc()
			}

			got, err := Client.DeleteOrganization(tt.ctx, tt.req)
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			// check details
			gotCD := got.GetDeletionDate().AsTime()
			if !tt.dontCheckTime {
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			}
		})
	}
}

func TestServer_DeactivateReactivateNonExistentOrganization(t *testing.T) {
	ctx := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	// deactivate non existent organization
	_, err := Client.DeactivateOrganization(ctx, &org.DeactivateOrganizationRequest{
		OrganizationId: "non existent organization",
	})
	require.Contains(t, err.Error(), "Organisation not found")

	// reactivate non existent organization
	_, err = Client.ActivateOrganization(ctx, &org.ActivateOrganizationRequest{
		OrganizationId: "non existent organization",
	})
	require.Contains(t, err.Error(), "Organisation not found")
}

func TestServer_ActivateOrganization(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		testFunc func() string
		err      error
	}{
		{
			name: "Activate, happy path",
			ctx:  CTX,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				// 2. deactivate organization once
				deactivate_res, err := Client.DeactivateOrganization(CTX, &org.DeactivateOrganizationRequest{
					OrganizationId: orgId,
				})
				require.NoError(t, err)
				gotCD := deactivate_res.GetChangeDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 3. check organization state is deactivated
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					listOrgRes, err := Client.ListOrganizations(CTX, &org.ListOrganizationsRequest{
						Queries: []*org.SearchQuery{
							{
								Query: &org.SearchQuery_IdQuery{
									IdQuery: &org.OrganizationIDQuery{
										Id: orgId,
									},
								},
							},
						},
					})
					require.NoError(ttt, err)
					if assert.GreaterOrEqual(ttt, len(listOrgRes.Result), 1) {
						require.Equal(ttt, org.OrganizationState_ORGANIZATION_STATE_INACTIVE, listOrgRes.Result[0].State)
					}
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
		},
		{
			name: "Activate, no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId
				return orgId
			},
			err: errors.New("membership not found"),
		},
		{
			name: "Activate, not existing",
			ctx:  CTX,
			testFunc: func() string {
				return "non-existing-org-id"
			},
			err: errors.New("Organisation not found"),
		},
		{
			name: "Activate, already activated",
			ctx:  CTX,
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId
				return orgId
			},
			err: errors.New("Organisation is already active"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var orgId string
			if tt.testFunc != nil {
				orgId = tt.testFunc()
			}
			_, err := Client.ActivateOrganization(tt.ctx, &org.ActivateOrganizationRequest{
				OrganizationId: orgId,
			})
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_DeactivateOrganization(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		testFunc func() string
		err      error
	}{
		{
			name: "Deactivate, happy path",
			ctx:  CTX,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				return orgId
			},
		},
		{
			name: "Deactivate, no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId
				return orgId
			},
			err: errors.New("membership not found"),
		},
		{
			name: "Deactivate, not existing",
			ctx:  CTX,
			testFunc: func() string {
				return "non-existing-org-id"
			},
			err: errors.New("Organisation not found"),
		},
		{
			name: "Deactivate, already deactivated",
			ctx:  CTX,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				// 2. deactivate organization once
				deactivate_res, err := Client.DeactivateOrganization(CTX, &org.DeactivateOrganizationRequest{
					OrganizationId: orgId,
				})
				require.NoError(t, err)
				gotCD := deactivate_res.GetChangeDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 3. check organization state is deactivated
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					listOrgRes, err := Client.ListOrganizations(CTX, &org.ListOrganizationsRequest{
						Queries: []*org.SearchQuery{
							{
								Query: &org.SearchQuery_IdQuery{
									IdQuery: &org.OrganizationIDQuery{
										Id: orgId,
									},
								},
							},
						},
					})
					require.NoError(ttt, err)
					require.Equal(ttt, org.OrganizationState_ORGANIZATION_STATE_INACTIVE, listOrgRes.Result[0].State)
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
			err: errors.New("Organisation is already deactivated"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var orgId string
			orgId = tt.testFunc()
			_, err := Client.DeactivateOrganization(tt.ctx, &org.DeactivateOrganizationRequest{
				OrganizationId: orgId,
			})
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_AddOrganizationDomain(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		domain   string
		testFunc func() string
		err      error
	}{
		{
			name:   "add org domain, happy path",
			ctx:    CTX,
			domain: integration.DomainName(),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId
				return orgId
			},
		},
		{
			name:   "no permission",
			ctx:    Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			domain: integration.DomainName(),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId
				return orgId
			},
			err: errors.New("membership not found"),
		},
		{
			name:   "add org domain, twice",
			ctx:    CTX,
			domain: integration.DomainName(),
			testFunc: func() string {
				t.Helper()
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				domain := integration.DomainName()
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// check domain added
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					queryRes, err := Client.ListOrganizationDomains(CTX, &org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)
					found := false
					for _, res := range queryRes.Domains {
						if res.Domain == domain {
							found = true
						}
					}
					require.True(ttt, found, "unable to find added domain")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
		},
		{
			name:   "add org domain to non existent org",
			ctx:    CTX,
			domain: integration.DomainName(),
			testFunc: func() string {
				return "non-existing-org-id"
			},
			err: errors.New("Organisation not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgId := tt.testFunc()
			addOrgDomainRes, err := Client.AddOrganizationDomain(tt.ctx, &org.AddOrganizationDomainRequest{
				OrganizationId: orgId,
				Domain:         tt.domain,
			})
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
			} else {
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			}
		})
	}
}

func TestServer_AddOrganizationDomain_ClaimDomain(t *testing.T) {
	domain := integration.DomainName()

	// create an organization, ensure it has globally unique usernames
	// and create a user with a loginname that matches the domain later on
	organization, err := Client.AddOrganization(CTX, &org.AddOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	_, err = Instance.Client.Admin.AddCustomDomainPolicy(CTX, &admin.AddCustomDomainPolicyRequest{
		OrgId:                 organization.GetOrganizationId(),
		UserLoginMustBeDomain: false,
	})
	require.NoError(t, err)
	username := integration.Username() + "@" + domain
	ownUser := Instance.CreateHumanUserVerified(CTX, organization.GetOrganizationId(), username, "")

	// create another organization, ensure it has globally unique usernames
	// and create a user with a loginname that matches the domain later on
	otherOrg, err := Client.AddOrganization(CTX, &org.AddOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	_, err = Instance.Client.Admin.AddCustomDomainPolicy(CTX, &admin.AddCustomDomainPolicyRequest{
		OrgId:                 otherOrg.GetOrganizationId(),
		UserLoginMustBeDomain: false,
	})
	require.NoError(t, err)

	otherUsername := integration.Username() + "@" + domain
	otherUser := Instance.CreateHumanUserVerified(CTX, otherOrg.GetOrganizationId(), otherUsername, "")

	// if we add the domain now to the first organization, it should be claimed on the second organization, resp. its user(s)
	_, err = Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
		OrganizationId: organization.GetOrganizationId(),
		Domain:         domain,
	})
	require.NoError(t, err)

	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		// check both users: the first one must be untouched, the second one must be updated
		users, err := Instance.Client.UserV2.ListUsers(CTX, &user.ListUsersRequest{
			Queries: []*user.SearchQuery{
				{
					Query: &user.SearchQuery_InUserIdsQuery{
						InUserIdsQuery: &user.InUserIDQuery{UserIds: []string{ownUser.GetUserId(), otherUser.GetUserId()}},
					},
				},
			},
		})
		require.NoError(collect, err)
		require.Len(collect, users.GetResult(), 2)

		for _, u := range users.GetResult() {
			if u.GetUserId() == ownUser.GetUserId() {
				assert.Equal(collect, username, u.GetPreferredLoginName())
				continue
			}
			if u.GetUserId() == otherUser.GetUserId() {
				assert.NotEqual(collect, otherUsername, u.GetPreferredLoginName())
				assert.Contains(collect, u.GetPreferredLoginName(), "@temporary.")
			}
		}
	}, 5*time.Second, time.Second, "user not updated in time")
}

func TestServer_DeleteOrganizationDomain(t *testing.T) {
	domain := integration.DomainName()
	tests := []struct {
		name     string
		ctx      context.Context
		domain   string
		testFunc func() string
		err      error
	}{
		{
			name:   "delete org domain, happy path",
			ctx:    CTX,
			domain: domain,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// check domain added
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					queryRes, err := Client.ListOrganizationDomains(CTX, &org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)

					found := slices.ContainsFunc(queryRes.Domains, func(d *org.Domain) bool { return d.GetDomain() == domain })
					require.True(ttt, found, "unable to find added domain")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
		},
		{
			name:   "delete org domain, twice",
			ctx:    CTX,
			domain: integration.DomainName(),
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				domain := integration.DomainName()
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// check domain added
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					queryRes, err := Client.ListOrganizationDomains(CTX, &org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)
					found := false
					for _, res := range queryRes.Domains {
						if res.Domain == domain {
							found = true
						}
					}
					require.True(ttt, found, "unable to find added domain")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				_, err = Client.DeleteOrganizationDomain(CTX, &org.DeleteOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)

				return orgId
			},
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name:   "delete org domain to non existent org",
			ctx:    CTX,
			domain: integration.DomainName(),
			testFunc: func() string {
				return "non-existing-org-id"
			},
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name:   "delete org domain no permission",
			ctx:    Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			domain: domain,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].OrganizationId

				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// check domain added
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					queryRes, err := Client.ListOrganizationDomains(CTX, &org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)

					found := slices.ContainsFunc(queryRes.Domains, func(d *org.Domain) bool { return d.GetDomain() == domain })
					require.True(ttt, found, "unable to find added domain")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
			err: errors.New("membership not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgId := tt.testFunc()

			_, err := Client.DeleteOrganizationDomain(tt.ctx, &org.DeleteOrganizationDomainRequest{
				OrganizationId: orgId,
				Domain:         tt.domain,
			})

			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServer_ValidateOrganizationDomain(t *testing.T) {
	t.Cleanup(func() {
		_, err := Instance.Client.Admin.UpdateDomainPolicy(CTX, &admin.UpdateDomainPolicyRequest{
			ValidateOrgDomains: false,
		})
		require.NoError(t, err)
	})

	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].OrganizationId

	_, err := Instance.Client.Admin.UpdateDomainPolicy(CTX, &admin.UpdateDomainPolicyRequest{
		ValidateOrgDomains: true,
	})
	if err != nil && !strings.Contains(err.Error(), "Organisation is already deactivated") {
		require.NoError(t, err)
	}

	domain := integration.DomainName()
	_, err = Client.AddOrganizationDomain(CTX, &org.AddOrganizationDomainRequest{
		OrganizationId: orgId,
		Domain:         domain,
	})
	require.NoError(t, err)

	tests := []struct {
		name string
		ctx  context.Context
		req  *org.GenerateOrganizationDomainValidationRequest
		err  error
	}{
		{
			name: "validate org http happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         domain,
				Type:           org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
		},
		{
			name: "validate org http non existent org id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: "non existent org id",
				Domain:         domain,
				Type:           org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name: "validate org dns happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         domain,
				Type:           org.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS,
			},
		},
		{
			name: "validate org dns non existent org id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: "non existent org id",
				Domain:         domain,
				Type:           org.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS,
			},
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name: "validate org non existent domain",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         "non existent domain",
				Type:           org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name: "validate without permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         domain,
				Type:           org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
			err: errors.New("membership not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.GenerateOrganizationDomainValidation(tt.ctx, tt.req)
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			require.NotEmpty(t, got.Token)
			require.Contains(t, got.Url, domain)
		})
	}
}

func TestServer_SetOrganizationMetadata(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].OrganizationId

	tests := []struct {
		name      string
		ctx       context.Context
		setupFunc func()
		orgId     string
		key       string
		value     string
		err       error
	}{
		{
			name:  "no permission",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			orgId: orgId,
			key:   "key1",
			value: "value1",
			err:   errors.New("membership not found"),
		},
		{
			name:  "set org metadata",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId: orgId,
			key:   "key1",
			value: "value1",
		},
		{
			name:  "set org metadata on non existant org",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId: "non existant orgid",
			key:   "key2",
			value: "value2",
			err:   errors.New("Organisation not found"),
		},
		{
			name: "update org metadata",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*org.Metadata{
						{
							Key:   "key3",
							Value: []byte("value3"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			key:   "key4",
			value: "value4",
		},
		{
			name: "update org metadata with same value",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*org.Metadata{
						{
							Key:   "key5",
							Value: []byte("value5"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			key:   "key5",
			value: "value5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}
			got, err := Client.SetOrganizationMetadata(tt.ctx, &org.SetOrganizationMetadataRequest{
				OrganizationId: tt.orgId,
				Metadata: []*org.Metadata{
					{
						Key:   tt.key,
						Value: []byte(tt.value),
					},
				},
			})
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			// check details
			gotCD := got.GetSetDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// check metadata
				listMetadataRes, err := Client.ListOrganizationMetadata(tt.ctx, &org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
				})
				require.NoError(ttt, err)
				foundMetadata := false
				foundMetadataKeyCount := 0
				for _, res := range listMetadataRes.Metadata {
					if res.Key == tt.key {
						foundMetadataKeyCount += 1
					}
					if res.Key == tt.key &&
						string(res.Value) == tt.value {
						foundMetadata = true
					}
				}
				require.True(ttt, foundMetadata, "unable to find added metadata")
				require.Equal(ttt, 1, foundMetadataKeyCount, "same metadata key found multiple times")
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func TestServer_DeleteOrganizationMetadata(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].OrganizationId

	_, err := Client.SetOrganizationMetadata(CTX, &org.SetOrganizationMetadataRequest{
		OrganizationId: orgId,
		Metadata: []*org.Metadata{
			{
				Key:   "key1",
				Value: []byte("value1"),
			},
			{
				Key:   "key2",
				Value: []byte("value2"),
			},
			{
				Key:   "key3",
				Value: []byte("value3"),
			}, {

				Key:   "key4",
				Value: []byte("value4"),
			},
		},
	})
	require.NoError(t, err)

	// check metadata exists
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 1*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		listOrgMetadataRes, err := Client.ListOrganizationMetadata(CTX, &org.ListOrganizationMetadataRequest{
			OrganizationId: orgId,
		})
		require.NoError(ttt, err)
		require.Len(ttt, listOrgMetadataRes.GetMetadata(), 4)
	}, retryDuration, tick, "timeout waiting for expected organizations being created")

	tests := []struct {
		name             string
		ctx              context.Context
		setupFunc        func()
		orgId            string
		metadataToDelete []struct {
			key   string
			value string
		}
		err error
	}{
		{
			name:  "delete org metadata happy path",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId: orgId,
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key1",
					value: "value1",
				},
			},
		},
		{
			name:  "delete multiple org metadata happy path",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId: orgId,
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key2",
					value: "value2",
				},
				{
					key:   "key3",
					value: "value3",
				},
			},
		},
		{
			name:  "delete org metadata that does not exist",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId: orgId,
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key5",
					value: "value5",
				},
			},
			err: errors.New("One or more keys do not exist"),
		},
		{
			name:  "delete org metadata for org that does not exist",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId: "non existant org id",
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key4",
					value: "value4",
				},
			},
			err: errors.New("Organisation not found"),
		},
		{
			name:  "delete org metadata without permission",
			ctx:   Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			orgId: orgId,
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key4",
					value: "value4",
				},
			},
			err: errors.New("membership not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			keys := make([]string, len(tt.metadataToDelete))
			for i, kvp := range tt.metadataToDelete {
				keys[i] = kvp.key
			}

			// run delete
			_, err := Client.DeleteOrganizationMetadata(tt.ctx, &org.DeleteOrganizationMetadataRequest{
				OrganizationId: tt.orgId,
				Keys:           keys,
			})
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// check metadata was definitely deleted
				listOrgMetadataRes, err := Client.ListOrganizationMetadata(tt.ctx, &org.ListOrganizationMetadataRequest{
					OrganizationId: tt.orgId,
				})
				require.NoError(ttt, err)
				foundMetadataCount := 0
				for _, kv := range tt.metadataToDelete {
					for _, res := range listOrgMetadataRes.Metadata {
						if res.Key == kv.key &&
							string(res.Value) == kv.value {
							foundMetadataCount += 1
						}
					}
				}
				require.Equal(ttt, foundMetadataCount, 0)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func createOrgs(ctx context.Context, t *testing.T, client org.OrganizationServiceClient, noOfOrgs int) ([]*org.AddOrganizationResponse, []string, []string) {
	var err error
	orgs := make([]*org.AddOrganizationResponse, noOfOrgs)
	orgNames := make([]string, noOfOrgs)
	orgDomains := make([]string, noOfOrgs)

	for i := range noOfOrgs {
		orgName := integration.OrganizationName()
		orgNames[i] = orgName
		orgs[i], err = client.AddOrganization(ctx,
			&org.AddOrganizationRequest{
				Name: orgName,
			},
		)
		require.NoError(t, err)
	}

	for i := range noOfOrgs {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			listOrgRes, err := client.ListOrganizations(ctx, &org.ListOrganizationsRequest{
				Queries: []*org.SearchQuery{
					{
						Query: &org.SearchQuery_IdQuery{
							IdQuery: &org.OrganizationIDQuery{
								Id: orgs[i].GetOrganizationId(),
							},
						},
					},
				},
			})
			require.NoError(collect, err)
			require.Len(collect, listOrgRes.Result, 1)

			orgDomains[i] = listOrgRes.Result[0].PrimaryDomain
		}, retryDuration, tick, "timeout waiting for org creation")
	}

	return orgs, orgNames, orgDomains
}
