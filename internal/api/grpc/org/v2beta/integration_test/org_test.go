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
	v2beta_object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

var (
	CTX      context.Context
	Instance *integration.Instance
	Client   v2beta_org.OrganizationServiceClient
	User     *user.AddHumanUserResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OrgV2beta

		CTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		User = Instance.CreateHumanUser(CTX)
		return m.Run()
	}())
}

func TestServer_CreateOrganization(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)

	type test struct {
		name     string
		ctx      context.Context
		req      *v2beta_org.CreateOrganizationRequest
		id       string
		testFunc func(ctx context.Context, t *testing.T)
		want     *v2beta_org.CreateOrganizationResponse
		wantErr  bool
	}

	tests := []test{
		{
			name: "missing permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &v2beta_org.CreateOrganizationRequest{
				Name:   "name",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name:   "",
				Admins: nil,
			},
			wantErr: true,
		},
		func() test {
			orgName := integration.OrganizationName()
			return test{
				name: "adding org with same name twice",
				ctx:  CTX,
				req: &v2beta_org.CreateOrganizationRequest{
					Name:   orgName,
					Admins: nil,
				},
				testFunc: func(ctx context.Context, t *testing.T) {
					// create org initially
					_, err := Client.CreateOrganization(ctx, &v2beta_org.CreateOrganizationRequest{
						Name: orgName,
					})
					require.NoError(t, err)
				},
				wantErr: true,
			}
		}(),
		{
			name: "invalid admin type",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "existing user as admin",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_UserId{UserId: User.GetUserId()},
					},
				},
			},
			want: &v2beta_org.CreateOrganizationResponse{
				OrganizationAdmins: []*v2beta_org.OrganizationAdmin{
					{
						OrganizationAdmin: &v2beta_org.OrganizationAdmin_AssignedAdmin{
							AssignedAdmin: &v2beta_org.AssignedAdmin{
								UserId: User.GetUserId(),
							},
						},
					},
				},
			},
		},
		{
			name: "admin with init",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_Human{
							Human: &user_v2beta.AddHumanUserRequest{
								Profile: &user_v2beta.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user_v2beta.SetHumanEmail{
									Email: integration.Email(),
									Verification: &user_v2beta.SetHumanEmail_ReturnCode{
										ReturnCode: &user_v2beta.ReturnEmailVerificationCode{},
									},
								},
							},
						},
					},
				},
			},
			want: &v2beta_org.CreateOrganizationResponse{
				Id: integration.NotEmpty,
				OrganizationAdmins: []*v2beta_org.OrganizationAdmin{
					{
						OrganizationAdmin: &v2beta_org.OrganizationAdmin_CreatedAdmin{
							CreatedAdmin: &v2beta_org.CreatedAdmin{
								UserId:    integration.NotEmpty,
								EmailCode: gu.Ptr(integration.NotEmpty),
								PhoneCode: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "existing user and new human with idp",
			ctx:  CTX,
			req: &v2beta_org.CreateOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_UserId{UserId: User.GetUserId()},
					},
					{
						UserType: &v2beta_org.CreateOrganizationRequest_Admin_Human{
							Human: &user_v2beta.AddHumanUserRequest{
								Profile: &user_v2beta.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user_v2beta.SetHumanEmail{
									Email: integration.Email(),
									Verification: &user_v2beta.SetHumanEmail_IsVerified{
										IsVerified: true,
									},
								},
								IdpLinks: []*user_v2beta.IDPLink{
									{
										IdpId:    idpResp.Id,
										UserId:   "userID",
										UserName: "username",
									},
								},
							},
						},
					},
				},
			},
			want: &v2beta_org.CreateOrganizationResponse{
				// OrganizationId: integration.NotEmpty,
				OrganizationAdmins: []*v2beta_org.OrganizationAdmin{
					{
						OrganizationAdmin: &v2beta_org.OrganizationAdmin_AssignedAdmin{
							AssignedAdmin: &v2beta_org.AssignedAdmin{
								UserId: User.GetUserId(),
							},
						},
					},
					{
						OrganizationAdmin: &v2beta_org.OrganizationAdmin_CreatedAdmin{
							CreatedAdmin: &v2beta_org.CreatedAdmin{
								UserId: integration.NotEmpty,
							},
						},
					},
				},
			},
		},
		{
			name: "create with ID",
			ctx:  CTX,
			id:   "custom_id",
			req: &v2beta_org.CreateOrganizationRequest{
				Name: integration.OrganizationName(),
				Id:   gu.Ptr("custom_id"),
			},
			want: &v2beta_org.CreateOrganizationResponse{
				Id: "custom_id",
			},
		},
		func() test {
			orgID := integration.OrganizationName()
			return test{
				name: "adding org with same ID twice",
				ctx:  CTX,
				req: &v2beta_org.CreateOrganizationRequest{
					Id:     &orgID,
					Name:   integration.OrganizationName(),
					Admins: nil,
				},
				testFunc: func(ctx context.Context, t *testing.T) {
					// create org initially
					_, err := Client.CreateOrganization(ctx, &v2beta_org.CreateOrganizationRequest{
						Id:   &orgID,
						Name: integration.OrganizationName(),
					})
					require.NoError(t, err)
				},
				wantErr: true,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.testFunc != nil {
				tt.testFunc(tt.ctx, t)
			}

			got, err := Client.CreateOrganization(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.id != "" {
				require.Equal(t, tt.id, got.Id)
			}

			// check details
			gotCD := got.GetCreationDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

			// check the admins
			require.Equal(t, len(tt.want.GetOrganizationAdmins()), len(got.GetOrganizationAdmins()))
			for i, admin := range tt.want.GetOrganizationAdmins() {
				gotAdmin := got.GetOrganizationAdmins()[i].OrganizationAdmin
				switch admin := admin.OrganizationAdmin.(type) {
				case *v2beta_org.OrganizationAdmin_CreatedAdmin:
					assertCreatedAdmin(t, admin.CreatedAdmin, gotAdmin.(*v2beta_org.OrganizationAdmin_CreatedAdmin).CreatedAdmin)
				case *v2beta_org.OrganizationAdmin_AssignedAdmin:
					assert.Equal(t, admin.AssignedAdmin.GetUserId(), gotAdmin.(*v2beta_org.OrganizationAdmin_AssignedAdmin).AssignedAdmin.GetUserId())
				}
			}
		})
	}
}

func TestServer_UpdateOrganization(t *testing.T) {
	orgs, orgsName, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].Id
	orgName := orgsName[0]

	tests := []struct {
		name    string
		ctx     context.Context
		req     *v2beta_org.UpdateOrganizationRequest
		want    *v2beta_org.UpdateOrganizationResponse
		wantErr bool
	}{
		{
			name: "update org with new name",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id:   orgId,
				Name: "new org name",
			},
		},
		{
			name: "update org with same name",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id:   orgId,
				Name: orgName,
			},
		},
		{
			name: "update org with non existent org id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id: "non existant org id",
				// Name: "",
			},
			wantErr: true,
		},
		{
			name: "update org with no id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.UpdateOrganizationRequest{
				Id:   "",
				Name: orgName,
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

func TestServer_ListOrganizations(t *testing.T) {
	testStartTimestamp := time.Now()
	ListOrgIinstance := integration.NewInstance(CTX)
	listOrgIAmOwnerCtx := ListOrgIinstance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	listOrgClient := ListOrgIinstance.Client.OrgV2beta

	noOfOrgs := 3
	orgs, orgsName, orgsDomain := createOrgs(listOrgIAmOwnerCtx, t, listOrgClient, noOfOrgs)

	// deactivat org[1]
	_, err := listOrgClient.DeactivateOrganization(listOrgIAmOwnerCtx, &v2beta_org.DeactivateOrganizationRequest{
		Id: orgs[1].Id,
	})
	require.NoError(t, err)

	tests := []struct {
		name  string
		ctx   context.Context
		query []*v2beta_org.OrganizationSearchFilter
		want  []*v2beta_org.Organization
		err   error
	}{
		{
			name: "list organizations, without required permissions",
			ctx:  ListOrgIinstance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			err:  errors.New("membership not found"),
		},
		{
			name: "list organizations happy path, no filter",
			ctx:  listOrgIAmOwnerCtx,
			want: []*v2beta_org.Organization{
				{
					// default org
					Name: "testinstance",
				},
				{
					Id:   orgs[0].Id,
					Name: orgsName[0],
				},
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
				{
					Id:   orgs[2].Id,
					Name: orgsName[2],
				},
			},
		},
		{
			name: "list organizations by id happy path",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_IdFilter{
						IdFilter: &v2beta_org.OrgIDFilter{
							Id: orgs[1].Id,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations by state active",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_StateFilter{
						StateFilter: &v2beta_org.OrgStateFilter{
							State: v2beta_org.OrgState_ORG_STATE_ACTIVE,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					// default org
					Name: "testinstance",
				},
				{
					Id:   orgs[0].Id,
					Name: orgsName[0],
				},
				{
					Id:   orgs[2].Id,
					Name: orgsName[2],
				},
			},
		},
		{
			name: "list organizations by state inactive",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_StateFilter{
						StateFilter: &v2beta_org.OrgStateFilter{
							State: v2beta_org.OrgState_ORG_STATE_INACTIVE,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations by id bad id",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_IdFilter{
						IdFilter: &v2beta_org.OrgIDFilter{
							Id: "bad id",
						},
					},
				},
			},
		},
		{
			name: "list organizations specify org name equals",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_NameFilter{
						NameFilter: &v2beta_org.OrgNameFilter{
							Name:   orgsName[1],
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations specify org name contains",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_NameFilter{
						NameFilter: &v2beta_org.OrgNameFilter{
							Name: func() string {
								return orgsName[1][1 : len(orgsName[1])-2]
							}(),
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations specify org name contains IGNORE CASE",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_NameFilter{
						NameFilter: &v2beta_org.OrgNameFilter{
							Name: func() string {
								return strings.ToUpper(orgsName[1][1 : len(orgsName[1])-2])
							}(),
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations specify domain name equals",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_DomainFilter{
						DomainFilter: &v2beta_org.OrgDomainFilter{
							Domain: orgsDomain[1],
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations specify domain name contains",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_DomainFilter{
						DomainFilter: &v2beta_org.OrgDomainFilter{
							Domain: func() string {
								domain := strings.ToLower(strings.ReplaceAll(orgsName[1][1:len(orgsName[1])-2], " ", "-"))
								return domain
							}(),
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
		{
			name: "list organizations specify org name contains IGNORE CASE",
			ctx:  listOrgIAmOwnerCtx,
			query: []*v2beta_org.OrganizationSearchFilter{
				{
					Filter: &v2beta_org.OrganizationSearchFilter_DomainFilter{
						DomainFilter: &v2beta_org.OrgDomainFilter{
							Domain: func() string {
								domain := strings.ToUpper(strings.ReplaceAll(orgsName[1][1:len(orgsName[1])-2], " ", "-"))
								return domain
							}(),
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
						},
					},
				},
			},
			want: []*v2beta_org.Organization{
				{
					Id:   orgs[1].Id,
					Name: orgsName[1],
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := listOrgClient.ListOrganizations(tt.ctx, &v2beta_org.ListOrganizationsRequest{
					Filter: tt.query,
				})
				if tt.err != nil {
					require.ErrorContains(ttt, err, tt.err.Error())
					return
				}
				require.NoError(ttt, err)

				require.Equal(ttt, uint64(len(tt.want)), got.Pagination.GetTotalResult())

				foundOrgs := 0
				for _, got := range got.Organizations {
					for _, org := range tt.want {

						// created/chagned date
						gotCD := got.GetCreationDate().AsTime()
						now := time.Now()
						assert.WithinRange(ttt, gotCD, testStartTimestamp, now.Add(time.Minute))
						gotCD = got.GetChangedDate().AsTime()
						assert.WithinRange(ttt, gotCD, testStartTimestamp, now.Add(time.Minute))

						// default org
						if org.Name == got.Name && got.Name == "testinstance" {
							foundOrgs += 1
							continue
						}

						if org.Name == got.Name &&
							org.Id == got.Id {
							foundOrgs += 1
						}
					}
				}
				require.Equal(ttt, len(tt.want), foundOrgs)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func TestServer_DeleteOrganization(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		createOrgFunc func() string
		req           *v2beta_org.DeleteOrganizationRequest
		want          *v2beta_org.DeleteOrganizationResponse
		dontCheckTime bool
		err           error
	}{
		{
			name: "delete org no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			createOrgFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				return orgs[0].Id
			},
			req: &v2beta_org.DeleteOrganizationRequest{},
			err: errors.New("membership not found"),
		},
		{
			name: "delete org happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			createOrgFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				return orgs[0].Id
			},
			req: &v2beta_org.DeleteOrganizationRequest{},
		},
		{
			name: "delete already deleted org",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			createOrgFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				// delete org
				_, err := Client.DeleteOrganization(CTX, &v2beta_org.DeleteOrganizationRequest{Id: orgs[0].Id})
				require.NoError(t, err)

				return orgs[0].Id
			},
			req:           &v2beta_org.DeleteOrganizationRequest{},
			dontCheckTime: true,
		},
		{
			name: "delete non existent org",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.DeleteOrganizationRequest{
				Id: "non existent org id",
			},
			dontCheckTime: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createOrgFunc != nil {
				tt.req.Id = tt.createOrgFunc()
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
	_, err := Client.DeactivateOrganization(ctx, &v2beta_org.DeactivateOrganizationRequest{
		Id: "non existent organization",
	})
	require.Contains(t, err.Error(), "Organisation not found")

	// reactivate non existent organization
	_, err = Client.ActivateOrganization(ctx, &v2beta_org.ActivateOrganizationRequest{
		Id: "non existent organization",
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
				orgId := orgs[0].Id

				// 2. deactivate organization once
				deactivate_res, err := Client.DeactivateOrganization(CTX, &v2beta_org.DeactivateOrganizationRequest{
					Id: orgId,
				})
				require.NoError(t, err)
				gotCD := deactivate_res.GetChangeDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 3. check organization state is deactivated
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					listOrgRes, err := Client.ListOrganizations(CTX, &v2beta_org.ListOrganizationsRequest{
						Filter: []*v2beta_org.OrganizationSearchFilter{
							{
								Filter: &v2beta_org.OrganizationSearchFilter_IdFilter{
									IdFilter: &v2beta_org.OrgIDFilter{
										Id: orgId,
									},
								},
							},
						},
					})
					require.NoError(ttt, err)
					if assert.GreaterOrEqual(ttt, len(listOrgRes.Organizations), 1) {
						require.Equal(ttt, v2beta_org.OrgState_ORG_STATE_INACTIVE, listOrgRes.Organizations[0].State)
					}
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
		},
		{
			name: "Activate, no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id
				return orgId
			},
			// BUG: this needs changing
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
				orgId := orgs[0].Id
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
			_, err := Client.ActivateOrganization(tt.ctx, &v2beta_org.ActivateOrganizationRequest{
				Id: orgId,
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
				orgId := orgs[0].Id

				return orgId
			},
		},
		{
			name: "Deactivate, no permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id
				return orgId
			},
			// BUG: this needs changing
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
				orgId := orgs[0].Id

				// 2. deactivate organization once
				deactivate_res, err := Client.DeactivateOrganization(CTX, &v2beta_org.DeactivateOrganizationRequest{
					Id: orgId,
				})
				require.NoError(t, err)
				gotCD := deactivate_res.GetChangeDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 3. check organization state is deactivated
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					listOrgRes, err := Client.ListOrganizations(CTX, &v2beta_org.ListOrganizationsRequest{
						Filter: []*v2beta_org.OrganizationSearchFilter{
							{
								Filter: &v2beta_org.OrganizationSearchFilter_IdFilter{
									IdFilter: &v2beta_org.OrgIDFilter{
										Id: orgId,
									},
								},
							},
						},
					})
					require.NoError(ttt, err)
					require.Equal(ttt, v2beta_org.OrgState_ORG_STATE_INACTIVE, listOrgRes.Organizations[0].State)
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
			_, err := Client.DeactivateOrganization(tt.ctx, &v2beta_org.DeactivateOrganizationRequest{
				Id: orgId,
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
			domain: integration.DomainName(),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id
				return orgId
			},
		},
		{
			name:   "add org domain, twice",
			domain: integration.DomainName(),
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id

				domain := integration.DomainName()
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
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
					queryRes, err := Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)
					found := false
					for _, res := range queryRes.Domains {
						if res.DomainName == domain {
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
			domain: integration.DomainName(),
			testFunc: func() string {
				return "non-existing-org-id"
			},
			// BUG: should return a error
			err: nil,
		},
	}

	for _, tt := range tests {
		var orgId string
		t.Run(tt.name, func(t *testing.T) {
			orgId = tt.testFunc()
		})
		addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
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
	}
}

func TestServer_AddOrganizationDomain_ClaimDomain(t *testing.T) {
	domain := integration.DomainName()

	// create an organization, ensure it has globally unique usernames
	// and create a user with a loginname that matches the domain later on
	organization, err := Client.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	_, err = Instance.Client.Admin.AddCustomDomainPolicy(CTX, &admin.AddCustomDomainPolicyRequest{
		OrgId:                 organization.GetId(),
		UserLoginMustBeDomain: false,
	})
	require.NoError(t, err)
	username := integration.Username() + "@" + domain
	ownUser := Instance.CreateHumanUserVerified(CTX, organization.GetId(), username, "")

	// create another organization, ensure it has globally unique usernames
	// and create a user with a loginname that matches the domain later on
	otherOrg, err := Client.CreateOrganization(CTX, &v2beta_org.CreateOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	_, err = Instance.Client.Admin.AddCustomDomainPolicy(CTX, &admin.AddCustomDomainPolicyRequest{
		OrgId:                 otherOrg.GetId(),
		UserLoginMustBeDomain: false,
	})
	require.NoError(t, err)

	otherUsername := integration.Username() + "@" + domain
	otherUser := Instance.CreateHumanUserVerified(CTX, otherOrg.GetId(), otherUsername, "")

	// if we add the domain now to the first organization, it should be claimed on the second organization, resp. its user(s)
	_, err = Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
		OrganizationId: organization.GetId(),
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

func TestServer_ListOrganizationDomains(t *testing.T) {
	domain := integration.DomainName()
	tests := []struct {
		name     string
		ctx      context.Context
		domain   string
		testFunc func() string
		err      error
	}{
		{
			name:   "list org domain, happy path",
			domain: domain,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				return orgId
			},
		},
	}

	for _, tt := range tests {
		var orgId string
		t.Run(tt.name, func(t *testing.T) {
			orgId = tt.testFunc()
		})

		var err error
		var queryRes *v2beta_org.ListOrganizationDomainsResponse

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
		require.EventuallyWithT(t, func(ttt *assert.CollectT) {
			queryRes, err = Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
				OrganizationId: orgId,
			})
			require.NoError(ttt, err)
			found := false
			for _, res := range queryRes.Domains {
				if res.DomainName == tt.domain {
					found = true
				}
			}
			require.True(ttt, found, "unable to find added domain")
		}, retryDuration, tick, "timeout waiting for adding domain")

	}
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
			domain: domain,
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id

				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
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
					queryRes, err := Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)

					found := slices.ContainsFunc(queryRes.Domains, func(d *v2beta_org.Domain) bool { return d.GetDomainName() == domain })
					require.True(ttt, found, "unable to find added domain")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				return orgId
			},
		},
		{
			name:   "delete org domain, twice",
			domain: integration.DomainName(),
			testFunc: func() string {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id

				domain := integration.DomainName()
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
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
					queryRes, err := Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)
					found := false
					for _, res := range queryRes.Domains {
						if res.DomainName == domain {
							found = true
						}
					}
					require.True(ttt, found, "unable to find added domain")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				_, err = Client.DeleteOrganizationDomain(CTX, &v2beta_org.DeleteOrganizationDomainRequest{
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
			domain: integration.DomainName(),
			testFunc: func() string {
				return "non-existing-org-id"
			},
			// BUG:
			err: errors.New("Domain doesn't exist on organization"),
		},
	}

	for _, tt := range tests {
		var orgId string
		t.Run(tt.name, func(t *testing.T) {
			orgId = tt.testFunc()
		})

		_, err := Client.DeleteOrganizationDomain(CTX, &v2beta_org.DeleteOrganizationDomainRequest{
			OrganizationId: orgId,
			Domain:         tt.domain,
		})

		if tt.err != nil {
			require.Contains(t, err.Error(), tt.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}

func TestServer_AddListDeleteOrganizationDomain(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func()
	}{
		{
			name: "add org domain, re-add org domain",
			testFunc: func() {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id

				domain := integration.DomainName()
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 3. re-add domain
				_, err = Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				// TODO remove error for adding already existing domain
				// require.NoError(t, err)
				require.Contains(t, err.Error(), "Errors.Already.Exists")
				// check details
				// gotCD = addOrgDomainRes.GetDetails().GetChangeDate().AsTime()
				// now = time.Now()
				// assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 4. check domain is added
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
				require.EventuallyWithT(t, func(collect *assert.CollectT) {
					queryRes, err := Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(collect, err)
					found := false
					for _, res := range queryRes.Domains {
						if res.DomainName == domain {
							found = true
						}
					}
					require.True(collect, found, "unable to find added domain")
				}, retryDuration, tick)
			},
		},
		{
			name: "add org domain, delete org domain, re-delete org domain",
			testFunc: func() {
				// 1. create organization
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id

				domain := integration.DomainName()
				// 2. add domain
				addOrgDomainRes, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD := addOrgDomainRes.GetCreationDate().AsTime()
				now := time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 2. delete organisation domain
				deleteOrgDomainRes, err := Client.DeleteOrganizationDomain(CTX, &v2beta_org.DeleteOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				require.NoError(t, err)
				// check details
				gotCD = deleteOrgDomainRes.GetDeletionDate().AsTime()
				now = time.Now()
				assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					// 3. check organization domain deleted
					queryRes, err := Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
						OrganizationId: orgId,
					})
					require.NoError(ttt, err)
					found := false
					for _, res := range queryRes.Domains {
						if res.DomainName == domain {
							found = true
						}
					}
					require.False(ttt, found, "deleted domain found")
				}, retryDuration, tick, "timeout waiting for expected organizations being created")

				// 4. redelete organisation domain
				_, err = Client.DeleteOrganizationDomain(CTX, &v2beta_org.DeleteOrganizationDomainRequest{
					OrganizationId: orgId,
					Domain:         domain,
				})
				// TODO remove error for deleting org domain already deleted
				// require.NoError(t, err)
				require.Contains(t, err.Error(), "Domain doesn't exist on organization")
				// check details
				// gotCD = deleteOrgDomainRes.GetDetails().GetChangeDate().AsTime()
				// now = time.Now()
				// assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

				// 5. check organization domain deleted
				queryRes, err := Client.ListOrganizationDomains(CTX, &v2beta_org.ListOrganizationDomainsRequest{
					OrganizationId: orgId,
				})
				require.NoError(t, err)
				found := false
				for _, res := range queryRes.Domains {
					if res.DomainName == domain {
						found = true
					}
				}
				require.False(t, found, "deleted domain found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc()
		})
	}
}

func TestServer_ValidateOrganizationDomain(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].Id

	_, err := Instance.Client.Admin.UpdateDomainPolicy(CTX, &admin.UpdateDomainPolicyRequest{
		ValidateOrgDomains: true,
	})
	if err != nil && !strings.Contains(err.Error(), "Organisation is already deactivated") {
		require.NoError(t, err)
	}

	domain := integration.DomainName()
	_, err = Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
		OrganizationId: orgId,
		Domain:         domain,
	})
	require.NoError(t, err)

	tests := []struct {
		name string
		ctx  context.Context
		req  *v2beta_org.GenerateOrganizationDomainValidationRequest
		err  error
	}{
		{
			name: "validate org http happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         domain,
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
		},
		{
			name: "validate org http non existnetn org id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: "non existent org id",
				Domain:         domain,
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
			// BUG: this should be 'organization does not exist'
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name: "validate org dns happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         domain,
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS,
			},
		},
		{
			name: "validate org dns non existnetn org id",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: "non existent org id",
				Domain:         domain,
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_DNS,
			},
			// BUG: this should be 'organization does not exist'
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name: "validate org non existnetn domain",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         "non existent domain",
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
			err: errors.New("Domain doesn't exist on organization"),
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
	orgId := orgs[0].Id

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
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
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
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
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
			got, err := Client.SetOrganizationMetadata(tt.ctx, &v2beta_org.SetOrganizationMetadataRequest{
				OrganizationId: tt.orgId,
				Metadata: []*v2beta_org.Metadata{
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
				listMetadataRes, err := Client.ListOrganizationMetadata(tt.ctx, &v2beta_org.ListOrganizationMetadataRequest{
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

func TestServer_ListOrganizationMetadata(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].Id

	tests := []struct {
		name          string
		ctx           context.Context
		setupFunc     func()
		orgId         string
		keyValuePairs []struct {
			key   string
			value string
		}
	}{
		{
			name: "list org metadata happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key1",
							Value: []byte("value1"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			keyValuePairs: []struct{ key, value string }{
				{
					key:   "key1",
					value: "value1",
				},
			},
		},
		{
			name: "list multiple org metadata happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key2",
							Value: []byte("value2"),
						},
						{
							Key:   "key3",
							Value: []byte("value3"),
						},
						{
							Key:   "key4",
							Value: []byte("value4"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			keyValuePairs: []struct{ key, value string }{
				{
					key:   "key2",
					value: "value2",
				},
				{
					key:   "key3",
					value: "value3",
				},
				{
					key:   "key4",
					value: "value4",
				},
			},
		},
		{
			name:          "list org metadata for non existent org",
			ctx:           Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			orgId:         "non existent orgid",
			keyValuePairs: []struct{ key, value string }{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListOrganizationMetadata(tt.ctx, &v2beta_org.ListOrganizationMetadataRequest{
					OrganizationId: tt.orgId,
				})
				require.NoError(ttt, err)

				foundMetadataCount := 0
				for _, kv := range tt.keyValuePairs {
					for _, res := range got.Metadata {
						if res.Key == kv.key &&
							string(res.Value) == kv.value {
							foundMetadataCount += 1
						}
					}
				}
				require.Len(ttt, tt.keyValuePairs, foundMetadataCount)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func TestServer_DeleteOrganizationMetadata(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].Id

	tests := []struct {
		name             string
		ctx              context.Context
		setupFunc        func()
		orgId            string
		metadataToDelete []struct {
			key   string
			value string
		}
		metadataToRemain []struct {
			key   string
			value string
		}
		err error
	}{
		{
			name: "delete org metadata happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key1",
							Value: []byte("value1"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key1",
					value: "value1",
				},
			},
		},
		{
			name: "delete multiple org metadata happy path",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key2",
							Value: []byte("value2"),
						},
						{
							Key:   "key3",
							Value: []byte("value3"),
						},
					},
				})
				require.NoError(t, err)
			},
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
			name: "delete some org metadata but not all",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key4",
							Value: []byte("value4"),
						},
						// key5 should not be deleted
						{
							Key:   "key5",
							Value: []byte("value5"),
						},
						{
							Key:   "key6",
							Value: []byte("value6"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			metadataToDelete: []struct{ key, value string }{
				{
					key:   "key4",
					value: "value4",
				},
				{
					key:   "key6",
					value: "value6",
				},
			},
			metadataToRemain: []struct{ key, value string }{
				{
					key:   "key5",
					value: "value5",
				},
			},
		},
		{
			name: "delete org metadata that does not exist",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key88",
							Value: []byte("value74"),
						},
						{
							Key:   "key5888",
							Value: []byte("value8885"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: orgId,
			// TODO: this error message needs to be either removed or changed
			err: errors.New("Metadata list is empty"),
		},
		{
			name: "delete org metadata for org that does not exist",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			setupFunc: func() {
				_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
					OrganizationId: orgId,
					Metadata: []*v2beta_org.Metadata{
						{
							Key:   "key88",
							Value: []byte("value74"),
						},
						{
							Key:   "key5888",
							Value: []byte("value8885"),
						},
					},
				})
				require.NoError(t, err)
			},
			orgId: "non existant org id",
			// TODO: this error message needs to be either removed or changed
			err: errors.New("Metadata list is empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			// check metadata exists
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				listOrgMetadataRes, err := Client.ListOrganizationMetadata(tt.ctx, &v2beta_org.ListOrganizationMetadataRequest{
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
				require.Equal(ttt, len(tt.metadataToDelete), foundMetadataCount)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")

			keys := make([]string, len(tt.metadataToDelete))
			for i, kvp := range tt.metadataToDelete {
				keys[i] = kvp.key
			}

			// run delete
			_, err := Client.DeleteOrganizationMetadata(tt.ctx, &v2beta_org.DeleteOrganizationMetadataRequest{
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
				listOrgMetadataRes, err := Client.ListOrganizationMetadata(tt.ctx, &v2beta_org.ListOrganizationMetadataRequest{
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

			// check metadata that should not be delted was not deleted
			listOrgMetadataRes, err := Client.ListOrganizationMetadata(tt.ctx, &v2beta_org.ListOrganizationMetadataRequest{
				OrganizationId: tt.orgId,
			})
			require.NoError(t, err)
			foundMetadataCount := 0
			for _, kv := range tt.metadataToRemain {
				for _, res := range listOrgMetadataRes.Metadata {
					if res.Key == kv.key &&
						string(res.Value) == kv.value {
						foundMetadataCount += 1
					}
				}
			}
			require.Equal(t, len(tt.metadataToRemain), foundMetadataCount)
		})
	}
}

func createOrgs(ctx context.Context, t *testing.T, client v2beta_org.OrganizationServiceClient, noOfOrgs int) ([]*v2beta_org.CreateOrganizationResponse, []string, []string) {
	var err error
	orgs := make([]*v2beta_org.CreateOrganizationResponse, noOfOrgs)
	orgNames := make([]string, noOfOrgs)
	orgDomains := make([]string, noOfOrgs)

	for i := range noOfOrgs {
		orgName := integration.OrganizationName()
		orgNames[i] = orgName
		orgs[i], err = client.CreateOrganization(ctx,
			&v2beta_org.CreateOrganizationRequest{
				Name: orgName,
			},
		)
		require.NoError(t, err)
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	for i := range noOfOrgs {
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			listOrgRes, err := client.ListOrganizations(ctx, &v2beta_org.ListOrganizationsRequest{
				Filter: []*v2beta_org.OrganizationSearchFilter{
					{
						Filter: &v2beta_org.OrganizationSearchFilter_IdFilter{
							IdFilter: &v2beta_org.OrgIDFilter{
								Id: orgs[i].Id,
							},
						},
					},
				},
			})
			require.NoError(collect, err)
			require.Len(collect, listOrgRes.Organizations, 1)

			orgDomains[i] = listOrgRes.Organizations[0].PrimaryDomain
		}, retryDuration, tick, "timeout waiting for org creation")
	}

	return orgs, orgNames, orgDomains
}

func assertCreatedAdmin(t *testing.T, expected, got *v2beta_org.CreatedAdmin) {
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
