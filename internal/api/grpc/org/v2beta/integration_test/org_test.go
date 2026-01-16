//go:build integration

package org_test

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2beta"
	v2beta_object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

var (
	CTX               context.Context
	Instance          *integration.Instance
	Client            v2beta_org.OrganizationServiceClient
	User              *user.AddHumanUserResponse
	OtherOrganization *org.AddOrganizationResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OrgV2beta

		CTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		User = Instance.CreateHumanUser(CTX)
		OtherOrganization = Instance.CreateOrganization(CTX, integration.OrganizationName(), integration.Email())
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
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	instOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)

	relationInst := integration.NewInstance(ctxWithSysAuthZ)
	instOwnerRelationCtx := relationInst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	integration.EnsureInstanceFeature(t, ctxWithSysAuthZ, relationInst,
		&feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)},
		func(tCollection *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.True(tCollection, got.EnableRelationalTables.GetEnabled())
		},
	)

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		relationInst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	type instanceAndCtx struct {
		testType string
		instance *integration.Instance
		ctx      context.Context
	}
	testedInstances := []instanceAndCtx{
		{testType: "eventstore", instance: inst, ctx: instOwnerCtx},
		{testType: "relations", instance: inst, ctx: instOwnerRelationCtx},
	}

	cases := len(testedInstances)
	orgs, orgNames, _ := createOrgs(CTX, t, Client, 2*cases)

	for i, stateCase := range testedInstances {
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
					Id:   orgs[0+(i*cases)].GetId(),
					Name: fmt.Sprintf("new org name %d", i),
				},
			},
			{
				name: "update org with same name",
				ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				req: &v2beta_org.UpdateOrganizationRequest{
					Id:   orgs[1+(i*cases)].GetId(),
					Name: orgNames[1+(i*cases)],
				},
				wantErr: true,
			},
			{
				name: "update org with non existent org id",
				ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				req: &v2beta_org.UpdateOrganizationRequest{
					Id:   "non existent org id",
					Name: "new name",
				},
				wantErr: true,
			},
			{
				name: "update org with no id",
				ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				req: &v2beta_org.UpdateOrganizationRequest{
					Id:   " ",
					Name: "new name",
				},
				wantErr: true,
			},
			{
				name: "no permission",
				ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &v2beta_org.UpdateOrganizationRequest{
					Id:   orgs[1+(i*cases)].GetId(),
					Name: integration.OrganizationName(),
				},
				wantErr: stateCase.instance.ID() == inst.ID(), //TODO: implement for relational tables
			},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s - %s", stateCase.testType, tt.name), func(t1 *testing.T) {
				got, err := Client.UpdateOrganization(tt.ctx, tt.req)
				if tt.wantErr {
					require.Error(t1, err)
					return
				}
				require.NoError(t1, err)

				// check details
				gotCD := got.GetChangeDate().AsTime()
				now := time.Now()
				assert.WithinRange(t1, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			})
		}
	}
}

func TestServer_ListOrganizations(t *testing.T) {
	t.Cleanup(func() {
		_, err := Instance.Client.FeatureV2.ResetInstanceFeatures(CTX, &feature.ResetInstanceFeaturesRequest{})
		require.NoError(t, err)
	})

	testStartTimestamp := time.Now()
	listOrgInstance := integration.NewInstance(CTX)
	listOrgIAmOwnerCtx := listOrgInstance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	listOrgClient := listOrgInstance.Client.OrgV2beta

	noOfOrgs := 3
	orgs, orgsName, orgsDomain := createOrgs(listOrgIAmOwnerCtx, t, listOrgClient, noOfOrgs)

	// deactivate org[1]
	_, err := listOrgClient.DeactivateOrganization(listOrgIAmOwnerCtx, &v2beta_org.DeactivateOrganizationRequest{
		Id: orgs[1].Id,
	})
	require.NoError(t, err)

	relTableState := integration.RelationalTablesEnableMatrix()

	tests := []struct {
		name  string
		ctx   context.Context
		query []*v2beta_org.OrganizationSearchFilter
		want  *v2beta_org.ListOrganizationsResponse
		err   error
	}{
		// TODO: re-enable when permission model is implemented in relational tables
		//{
		//	name: "list organizations, without required permissions",
		//	ctx:  listOrgInstance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
		//	want: &v2beta_org.ListOrganizationsResponse{
		//		Pagination: &filter.PaginationResponse{
		//			TotalResult: 4,
		//		},
		//	},
		//},
		{
			name: "list organizations happy path, no filter",
			ctx:  listOrgIAmOwnerCtx,
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 4,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   listOrgInstance.DefaultOrg.Id,
						Name: listOrgInstance.DefaultOrg.Name,
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 3,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   listOrgInstance.DefaultOrg.Id,
						Name: listOrgInstance.DefaultOrg.Name,
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 0,
				},
				Organizations: nil,
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
							Domain: orgsDomain[1][1 : len(orgsDomain[1])-2],
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
						},
					},
				},
			},
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
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
							Domain: strings.ToUpper(orgsDomain[1][1 : len(orgsDomain[1])-2]),
							Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
						},
					},
				},
			},
			want: &v2beta_org.ListOrganizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult: 1,
				},
				Organizations: []*v2beta_org.Organization{
					{
						Id:   orgs[1].Id,
						Name: orgsName[1],
					},
				},
			},
		},
	}

	for _, stateCase := range relTableState {
		integration.EnsureInstanceFeature(t, listOrgIAmOwnerCtx, listOrgInstance, stateCase.FeatureSet, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.Equal(tCollect, stateCase.FeatureSet.GetEnableRelationalTables(), got.EnableRelationalTables.GetEnabled())
		})

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s - %s", stateCase.State, tt.name), func(t *testing.T) {
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 20*time.Second)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					got, err := listOrgClient.ListOrganizations(tt.ctx, &v2beta_org.ListOrganizationsRequest{
						Filter: tt.query,
						Pagination: &filter.PaginationRequest{
							Asc: true,
						},
						SortingColumn: v2beta_org.OrgFieldName_ORG_FIELD_NAME_CREATION_DATE,
					})
					if tt.err != nil {
						require.ErrorContains(ttt, err, tt.err.Error())
						return
					}
					require.NoError(ttt, err)

					require.Equal(ttt, tt.want.GetPagination(), got.GetPagination())

					require.Len(ttt, got.Organizations, len(tt.want.Organizations))

					for i, got := range got.Organizations {
						// created/chagned date
						gotCD := got.GetCreationDate().AsTime()
						now := time.Now()
						assert.WithinRange(ttt, gotCD, testStartTimestamp, now.Add(time.Minute))
						gotCD = got.GetChangedDate().AsTime()
						assert.WithinRange(ttt, gotCD, testStartTimestamp, now.Add(time.Minute))

						assert.Equal(ttt, tt.want.Organizations[i].Id, got.Id)
						assert.Equal(ttt, tt.want.Organizations[i].Name, got.Name)
					}
				}, retryDuration, tick, "timeout waiting for expected organizations being created")
			})
		}
	}
}

func TestServer_DeleteOrganization(t *testing.T) {
	t.Cleanup(func() {
		_, err := Instance.Client.FeatureV2.ResetInstanceFeatures(CTX, &feature.ResetInstanceFeaturesRequest{})
		require.NoError(t, err)
	})

	relTableState := integration.RelationalTablesEnableMatrix()
	var orgs []*v2beta_org.CreateOrganizationResponse
	orgsNumPerCase := 3
	orgs, _, _ = createOrgs(CTX, t, Client, orgsNumPerCase*len(relTableState))
	require.NotNil(t, orgs)
	require.NotEmpty(t, orgs)

	for i, stateCase := range relTableState {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		require.EventuallyWithT(t, func(ttt *assert.CollectT) {
			deleteRes, err := Client.DeleteOrganization(CTX, &v2beta_org.DeleteOrganizationRequest{Id: orgs[2+(orgsNumPerCase*i)].GetId()})
			assert.Nil(ttt, err)
			assert.NotNil(ttt, deleteRes)
			assert.NotZero(ttt, deleteRes.GetDeletionDate())
		}, retryDuration, tick)

		tests := []struct {
			name          string
			ctx           context.Context
			req           *v2beta_org.DeleteOrganizationRequest
			want          *v2beta_org.DeleteOrganizationResponse
			dontCheckTime bool
			err           error
		}{
			{
				name: "delete org no permission",
				ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &v2beta_org.DeleteOrganizationRequest{
					Id: orgs[0+(orgsNumPerCase*i)].GetId(),
				},
				err: errors.New("membership not found"),
			},
			{
				name: "delete org happy path",
				ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				req: &v2beta_org.DeleteOrganizationRequest{
					Id: orgs[1+(orgsNumPerCase*i)].GetId(),
				},
			},
			{
				name:          "delete already deleted org",
				ctx:           Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				req:           &v2beta_org.DeleteOrganizationRequest{Id: orgs[2+(orgsNumPerCase*i)].GetId()},
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
		integration.EnsureInstanceFeature(t, CTX, Instance, stateCase.FeatureSet, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.Equal(tCollect, stateCase.FeatureSet.GetEnableRelationalTables(), got.EnableRelationalTables.GetEnabled())
		})
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s - %s", stateCase.State, tt.name), func(t *testing.T) {
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
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	type TestCase struct {
		name      string
		inst      *integration.Instance
		instOwner context.Context
	}

	cases := []TestCase{
		func() TestCase {
			inst := integration.NewInstance(ctxWithSysAuthZ)
			instOwner := inst.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
			return TestCase{
				name:      "eventstore",
				inst:      inst,
				instOwner: instOwner,
			}
		}(),
		func() TestCase {
			inst := integration.NewInstance(ctxWithSysAuthZ)
			instOwner := inst.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
			_, err := inst.Client.FeatureV2.SetInstanceFeatures(instOwner, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(false)})
			require.NoError(t, err)
			return TestCase{
				name:      "relational",
				inst:      inst,
				instOwner: instOwner,
			}
		}(),
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(tt *testing.T) {
			inst := testCase.inst
			client := inst.Client.OrgV2beta
			instOwner := testCase.instOwner

			tt.Cleanup(func() {
				_, err := inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
				assert.NoError(tt, err)
			})

			tt.Run("Happy path", func(ttt *testing.T) {
				// given
				orgs, _, _ := createOrgs(instOwner, ttt, client, 1)
				organisation := orgs[0]

				_, err := client.DeactivateOrganization(instOwner, &v2beta_org.DeactivateOrganizationRequest{Id: organisation.Id})
				require.NoError(ttt, err)

				// when
				_, err = client.ActivateOrganization(instOwner, &v2beta_org.ActivateOrganizationRequest{Id: organisation.Id})

				// then
				assert.NoError(ttt, err)
			})

			tt.Run("Unhappy: no permission", func(ttt *testing.T) {
				// given
				orgs, _, _ := createOrgs(instOwner, ttt, client, 1)
				organisation := orgs[0]
				usersWithoutPermissions := []integration.UserType{
					integration.UserTypeOrgOwner,
				}

				_, err := client.DeactivateOrganization(instOwner, &v2beta_org.DeactivateOrganizationRequest{Id: organisation.Id})
				assert.NoError(ttt, err)

				for _, userType := range usersWithoutPermissions {
					ttt.Run(userType.String(), func(tttt *testing.T) {
						u := inst.WithAuthorizationToken(ctx, userType)

						// when
						_, err = client.ActivateOrganization(u, &v2beta_org.ActivateOrganizationRequest{Id: organisation.Id})

						// then
						assert.ErrorContains(tttt, err, "membership not found")
					})
				}
			})

			tt.Run("Unhappy: unknown org", func(ttt *testing.T) {
				// when
				_, err := client.ActivateOrganization(instOwner, &v2beta_org.ActivateOrganizationRequest{Id: "does not exist"})

				// then
				assert.ErrorContains(ttt, err, "Organisation not found")
			})

			tt.Run("Unhappy: already activated", func(ttt *testing.T) {
				// given
				orgs, _, _ := createOrgs(instOwner, ttt, inst.Client.OrgV2beta, 1)
				organisation := orgs[0]

				// when
				_, err := client.ActivateOrganization(instOwner, &v2beta_org.ActivateOrganizationRequest{Id: organisation.Id})

				// then
				assert.ErrorContains(ttt, err, "Organisation is already active")
			})
		})
	}
}

func TestServer_DeactivateOrganization(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	type TestCase struct {
		name      string
		inst      *integration.Instance
		instOwner context.Context
	}

	cases := []TestCase{
		func() TestCase {
			inst := integration.NewInstance(ctxWithSysAuthZ)
			instOwner := inst.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
			return TestCase{
				name:      "eventstore",
				inst:      inst,
				instOwner: instOwner,
			}
		}(),
		func() TestCase {
			inst := integration.NewInstance(ctxWithSysAuthZ)
			instOwner := inst.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
			_, err := inst.Client.FeatureV2.SetInstanceFeatures(instOwner, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(false)})
			require.NoError(t, err)
			return TestCase{
				name:      "relational",
				inst:      inst,
				instOwner: instOwner,
			}
		}(),
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(tt *testing.T) {
			inst := testCase.inst
			client := inst.Client.OrgV2beta
			instOwner := testCase.instOwner

			tt.Cleanup(func() {
				_, err := inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
				require.NoError(tt, err)
			})

			tt.Run("Happy path", func(ttt *testing.T) {
				// given
				orgs, _, _ := createOrgs(instOwner, ttt, client, 1)
				organisation := orgs[0]

				// when
				_, err := client.DeactivateOrganization(instOwner, &v2beta_org.DeactivateOrganizationRequest{Id: organisation.Id})

				// then
				assert.NoError(ttt, err)
			})

			tt.Run("Unhappy: no permission", func(ttt *testing.T) {
				// given
				orgs, _, _ := createOrgs(instOwner, ttt, client, 1)
				organisation := orgs[0]
				usersWithoutPermissions := []integration.UserType{
					integration.UserTypeOrgOwner,
				}

				for _, userType := range usersWithoutPermissions {
					ttt.Run(userType.String(), func(tttt *testing.T) {
						u := inst.WithAuthorizationToken(ctx, userType)

						// when
						_, err := client.DeactivateOrganization(u, &v2beta_org.DeactivateOrganizationRequest{Id: organisation.Id})

						// then
						assert.ErrorContains(tttt, err, "membership not found")
					})
				}
			})

			tt.Run("Unhappy: unknown org", func(ttt *testing.T) {
				// when
				_, err := client.DeactivateOrganization(instOwner, &v2beta_org.DeactivateOrganizationRequest{Id: "does not exist"})

				// then
				assert.ErrorContains(ttt, err, "Organisation not found")
			})

			tt.Run("Unhappy: already deactivated", func(ttt *testing.T) {
				// given
				orgs, _, _ := createOrgs(instOwner, ttt, inst.Client.OrgV2beta, 1)
				organisation := orgs[0]
				_, err := client.DeactivateOrganization(instOwner, &v2beta_org.DeactivateOrganizationRequest{Id: organisation.Id})
				require.NoError(ttt, err)

				// when
				_, err = client.DeactivateOrganization(instOwner, &v2beta_org.DeactivateOrganizationRequest{Id: organisation.Id})

				// then
				assert.ErrorContains(ttt, err, "Organisation is already deactivated")
			})
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
				orgId := orgs[0].Id
				return orgId
			},
		},
		{
			name:   "no permission",
			ctx:    Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			domain: integration.DomainName(),
			testFunc: func() string {
				orgs, _, _ := createOrgs(CTX, t, Client, 1)
				orgId := orgs[0].Id
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
			addOrgDomainRes, err := Client.AddOrganizationDomain(tt.ctx, &v2beta_org.AddOrganizationDomainRequest{
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

	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].Id

	var primaryDomain string
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Second)
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		organizations, err := Client.ListOrganizations(CTX, &v2beta_org.ListOrganizationsRequest{
			Filter: []*v2beta_org.OrganizationSearchFilter{
				{Filter: &v2beta_org.OrganizationSearchFilter_IdFilter{
					IdFilter: &v2beta_org.OrgIDFilter{Id: orgId},
				}},
			},
		})
		require.NoError(t, err)
		require.Len(t, organizations.GetOrganizations(), 1)
		primaryDomain = organizations.GetOrganizations()[0].GetPrimaryDomain()
	}, retryDuration, tick, "could not find primary domain")

	_, err := Client.AddOrganizationDomain(CTX, &v2beta_org.AddOrganizationDomainRequest{
		OrganizationId: orgId,
		Domain:         domain,
	})
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		request *v2beta_org.ListOrganizationDomainsRequest
	}
	type want struct {
		response *v2beta_org.ListOrganizationDomainsResponse
		err      bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "non existing organization",
			args: args{
				ctx:     CTX,
				request: &v2beta_org.ListOrganizationDomainsRequest{OrganizationId: "not-existing"},
			},
			want: want{
				response: &v2beta_org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 0,
					},
					Domains: nil,
				},
			},
		},
		{
			name: "no permission (different organization), error",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				request: &v2beta_org.ListOrganizationDomainsRequest{
					OrganizationId: orgId,
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 0,
					},
					Domains: nil,
				},
			},
		},
		{
			name: "list org domain, all domains",
			args: args{
				ctx: CTX,
				request: &v2beta_org.ListOrganizationDomainsRequest{
					OrganizationId: orgId,
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 2,
					},
					Domains: []*v2beta_org.Domain{
						{
							OrganizationId: orgId,
							DomainName:     domain,
							IsVerified:     true,
							IsPrimary:      false,
							ValidationType: 0,
						},
						{
							OrganizationId: orgId,
							DomainName:     primaryDomain,
							IsVerified:     true,
							IsPrimary:      true,
							ValidationType: 0,
						},
					},
				},
			},
		},
		{
			name: "list specific domain",
			args: args{
				ctx: CTX,
				request: &v2beta_org.ListOrganizationDomainsRequest{
					OrganizationId: orgId,
					Filters: []*v2beta_org.DomainSearchFilter{
						{Filter: &v2beta_org.DomainSearchFilter_DomainNameFilter{DomainNameFilter: &v2beta_org.DomainNameFilter{Name: domain}}},
					},
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationDomainsResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 1,
					},
					Domains: []*v2beta_org.Domain{
						{
							OrganizationId: orgId,
							DomainName:     domain,
							IsVerified:     true,
							IsPrimary:      false,
							ValidationType: 0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, 10*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				queryRes, err := Client.ListOrganizationDomains(tt.args.ctx, tt.args.request)
				if tt.want.err {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				assert.Len(ttt, queryRes.Domains, int(tt.want.response.GetPagination().GetTotalResult()))
				assert.EqualExportedValues(ttt, tt.want.response.GetPagination(), queryRes.GetPagination())
				assert.ElementsMatch(ttt, tt.want.response.GetDomains(), queryRes.GetDomains())
			}, retryDuration, tick, "timeout waiting for adding domain")
		})
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
			ctx:    CTX,
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
			ctx:    CTX,
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
			ctx:    CTX,
			domain: integration.DomainName(),
			testFunc: func() string {
				return "non-existing-org-id"
			},
			// BUG:
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name:   "delete org domain no permission",
			ctx:    Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
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
			err: errors.New("membership not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgId := tt.testFunc()

			_, err := Client.DeleteOrganizationDomain(tt.ctx, &v2beta_org.DeleteOrganizationDomainRequest{
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
			name: "validate org non existent domain",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         "non existent domain",
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
			},
			err: errors.New("Domain doesn't exist on organization"),
		},
		{
			name: "validate without permission",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &v2beta_org.GenerateOrganizationDomainValidationRequest{
				OrganizationId: orgId,
				Domain:         domain,
				Type:           v2beta_org.DomainValidationType_DOMAIN_VALIDATION_TYPE_HTTP,
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
	setRespoonse, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
		OrganizationId: orgId,
		Metadata: []*v2beta_org.Metadata{
			{
				Key:   "key1",
				Value: []byte("value1"),
			},
			{
				Key:   "key2",
				Value: []byte("value2"),
			},
			{
				Key:   "key2.1",
				Value: []byte("value3"),
			},
			{
				Key:   "key2.2",
				Value: []byte("value4"),
			},
		},
	})
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		request *v2beta_org.ListOrganizationMetadataRequest
	}
	type want struct {
		response *v2beta_org.ListOrganizationMetadataResponse
		err      error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "list org metadata happy path",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				request: &v2beta_org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult: 4,
					},
					Metadata: []*metadata.Metadata{
						{
							Key:          "key1",
							Value:        []byte("value1"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2",
							Value:        []byte("value2"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2.1",
							Value:        []byte("value3"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2.2",
							Value:        []byte("value4"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
					},
				},
			},
		},
		{
			name: "list org metadata filter key",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				request: &v2beta_org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
					Pagination: &filter.PaginationRequest{
						Offset: 1,
						Limit:  2,
					},
					Filter: []*metadata.MetadataQuery{
						{
							Query: &metadata.MetadataQuery_KeyQuery{
								KeyQuery: &metadata.MetadataKeyQuery{
									Key:    "key2",
									Method: v2beta_object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH,
								},
							},
						},
					},
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{
						TotalResult:  3,
						AppliedLimit: 2,
					},
					Metadata: []*metadata.Metadata{
						{
							Key:          "key2.1",
							Value:        []byte("value3"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
						{
							Key:          "key2.2",
							Value:        []byte("value4"),
							CreationDate: setRespoonse.GetSetDate(),
							ChangeDate:   setRespoonse.GetSetDate(),
						},
					},
				},
			},
		},
		{
			name: "list org metadata for non existent org",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
				request: &v2beta_org.ListOrganizationMetadataRequest{
					OrganizationId: "non existent orgid",
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{},
				},
			},
		},
		{
			name: "list org metadata without permission (other organization)",
			args: args{
				ctx: Instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				request: &v2beta_org.ListOrganizationMetadataRequest{
					OrganizationId: orgId,
				},
			},
			want: want{
				response: &v2beta_org.ListOrganizationMetadataResponse{
					Pagination: &filter.PaginationResponse{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 1*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListOrganizationMetadata(tt.args.ctx, tt.args.request)
				require.NoError(ttt, err)

				assert.EqualExportedValues(ttt, tt.want.response, got)
			}, retryDuration, tick, "timeout waiting for expected organizations being created")
		})
	}
}

func TestServer_DeleteOrganizationMetadata(t *testing.T) {
	orgs, _, _ := createOrgs(CTX, t, Client, 1)
	orgId := orgs[0].Id

	_, err := Client.SetOrganizationMetadata(CTX, &v2beta_org.SetOrganizationMetadataRequest{
		OrganizationId: orgId,
		Metadata: []*v2beta_org.Metadata{
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
		listOrgMetadataRes, err := Client.ListOrganizationMetadata(CTX, &v2beta_org.ListOrganizationMetadataRequest{
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
			_, err := Client.DeleteOrganizationMetadata(tt.ctx, &v2beta_org.DeleteOrganizationMetadataRequest{
				OrganizationId: tt.orgId,
				Keys:           keys,
			})
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 10*time.Minute)
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

	for i := range noOfOrgs {
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
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
