//go:build integration

package authorization_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/authorization/v2"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_ListAuthorizations(t *testing.T) {
	iamOwnerCtx := InstanceQuery.WithAuthorizationToken(EmptyCTX, integration.UserTypeIAMOwner)

	organizationName := integration.OrganizationName()
	organization := InstanceQuery.CreateOrganization(iamOwnerCtx, organizationName, integration.Email())

	projectName1 := integration.ProjectName()
	project1 := InstanceQuery.CreateProject(iamOwnerCtx, t, organization.GetOrganizationId(), projectName1, false, false)

	projectName2 := integration.ProjectName()
	project2 := InstanceQuery.CreateProject(iamOwnerCtx, t, organization.GetOrganizationId(), projectName2, false, false)

	roles := []*authorization.Role{
		{Key: "role1", DisplayName: "Role 1", Group: "group1"},
		{Key: "role2", DisplayName: "Role 2", Group: "group2"},
		{Key: "role3", DisplayName: "Role 3", Group: "group2"}}
	for _, role := range roles {
		InstanceQuery.AddProjectRole(iamOwnerCtx, t, project1.GetId(), role.Key, role.DisplayName, role.Group)
		InstanceQuery.AddProjectRole(iamOwnerCtx, t, project2.GetId(), role.Key, role.DisplayName, role.Group)
	}

	grantedOrganizationName := integration.OrganizationName()
	grantedOrganization := InstanceQuery.CreateOrganization(iamOwnerCtx, grantedOrganizationName, integration.Email())

	grantedRoles := roles[0:2]
	grantedRoleKeys := make([]string, len(grantedRoles))
	for i, role := range grantedRoles {
		grantedRoleKeys[i] = role.Key
	}

	InstanceQuery.CreateProjectGrant(iamOwnerCtx, t, project2.GetId(), grantedOrganization.GetOrganizationId(), grantedRoleKeys...)

	user1 := InstanceQuery.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
	user2 := InstanceQuery.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
	user3 := InstanceQuery.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

	authorizationUser1Project1 := createAuthorizationForProject(iamOwnerCtx, InstanceQuery, t, organization.GetOrganizationId(), organizationName, user1.GetId(), project1.GetId(), projectName1, roles[0:1])
	authorizationUser2Project1 := createAuthorizationForProject(iamOwnerCtx, InstanceQuery, t, organization.GetOrganizationId(), organizationName, user2.GetId(), project1.GetId(), projectName1, roles[0:1])
	authorizationUser2Project2 := createAuthorizationForProject(iamOwnerCtx, InstanceQuery, t, organization.GetOrganizationId(), organizationName, user2.GetId(), project2.GetId(), projectName2, roles[1:2])
	authorizationUser3Project2 := createAuthorizationForProject(iamOwnerCtx, InstanceQuery, t, organization.GetOrganizationId(), organizationName, user3.GetId(), project2.GetId(), projectName2, roles[2:3])
	authorizationUser3ProjectGrant := createAuthorizationForProjectGrant(iamOwnerCtx, InstanceQuery, t, organization.GetOrganizationId(), user3.GetId(), project2.GetId(), projectName2, grantedOrganization.GetOrganizationId(), grantedOrganizationName, grantedRoles)

	organizationOwnerResp := InstanceQuery.CreateMachineUser(iamOwnerCtx)
	organizationOwnerPatResp := InstanceQuery.CreatePersonalAccessToken(iamOwnerCtx, organizationOwnerResp.GetUserId())
	InstanceQuery.CreateOrgMembership(t, iamOwnerCtx, organization.GetOrganizationId(), organizationOwnerResp.GetUserId())
	organizationOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, organizationOwnerPatResp.Token)

	projectOwnerResp := InstanceQuery.CreateMachineUser(iamOwnerCtx)
	projectOwnerPatResp := InstanceQuery.CreatePersonalAccessToken(iamOwnerCtx, projectOwnerResp.GetUserId())
	InstanceQuery.CreateProjectMembership(t, iamOwnerCtx, project2.GetId(), projectOwnerResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectOwnerPatResp.Token)

	projectGrantOwnerResp := InstanceQuery.CreateMachineUser(iamOwnerCtx)
	projectGrantOwnerPatResp := InstanceQuery.CreatePersonalAccessToken(iamOwnerCtx, projectGrantOwnerResp.GetUserId())
	InstanceQuery.CreateProjectGrantMembership(t, iamOwnerCtx, project2.GetId(), grantedOrganization.GetOrganizationId(), projectGrantOwnerResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectGrantOwnerPatResp.Token)

	type args struct {
		ctx context.Context
		req *authorization.ListAuthorizationsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *authorization.ListAuthorizationsResponse
		wantErr bool
	}{
		{
			name: "unauthenticated",
			args: args{
				ctx: EmptyCTX,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user1.GetId(),
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			args: args{
				ctx: InstanceQuery.WithAuthorizationToken(EmptyCTX, integration.UserTypeNoPermission),
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user1.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: nil,
			},
		},
		{
			name: "not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: "notexisting",
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "limit and pagination",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Pagination: &filter.PaginationRequest{
						Limit:  3,
						Offset: 1,
						Asc:    true,
					},
					SortingColumn: authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_CREATED_DATE,
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5,
					AppliedLimit: 3,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser2Project1,
					authorizationUser2Project2,
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "organization owner permission, only expect non granted",
			args: args{
				ctx: organizationOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // count issue with permission check
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "project owner permission, only expect own project non granted",
			args: args{
				ctx: projectOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // count issue with permission check
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "project grant owner permission, only expect granted",
			args: args{
				ctx: projectGrantOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // count issue with permission check
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
				},
			},
		},
		{
			name: "filter by single authorizationID, expect one result",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_AuthorizationIds{
								AuthorizationIds: &filter.InIDsFilter{
									Ids: []string{authorizationUser1Project1.GetId()},
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter by organization ID (expect only own / non granted)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_OrganizationId{
								OrganizationId: &filter.IDFilter{
									Id: organization.GetOrganizationId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3Project2,
					authorizationUser2Project2,
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter by organization ID (only project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_OrganizationId{
								OrganizationId: &filter.IDFilter{
									Id: grantedOrganization.GetOrganizationId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
				},
			},
		},
		{
			name: "filter by userID",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user1.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "list by userID (where project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "filter by projectID (no project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_ProjectId{
								ProjectId: &filter.IDFilter{
									Id: project1.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter by projectID (project grant exists), expect own and granted",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_ProjectId{
								ProjectId: &filter.IDFilter{
									Id: project2.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
					authorizationUser3Project2,
					authorizationUser2Project2,
				},
			},
		},
		{
			name: "filter by project name (no project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_ProjectName{
								ProjectName: &authorization.ProjectNameQuery{
									Name:   projectName1,
									Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter role key, expect all with that role, own and granted",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_RoleKey{
								RoleKey: &authorization.RoleKeyQuery{
									Key:    roles[0].Key,
									Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := InstanceQuery.Client.AuthorizationV2.ListAuthorizations(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Authorizations, len(tt.want.Authorizations)) {
					assert.EqualExportedValues(ttt, tt.want.Authorizations, got.Authorizations)
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func createAuthorizationForProject(ctx context.Context, instance *integration.Instance, t *testing.T, organizationID, organizationName, userID, projectID, projectName string, roles []*authorization.Role) *authorization.Authorization {
	userResp, err := instance.Client.UserV2.GetUserByID(ctx, &user.GetUserByIDRequest{UserId: userID})
	require.NoError(t, err)

	roleKeys := make([]string, len(roles))
	for i, role := range roles {
		roleKeys[i] = role.Key
	}
	authResp := instance.CreateAuthorizationProject(t, ctx, projectID, userID, organizationID, roleKeys...)
	return &authorization.Authorization{
		Id:           authResp.GetId(),
		CreationDate: authResp.GetCreationDate(),
		ChangeDate:   authResp.GetCreationDate(),
		Project: &authorization.Project{
			Id:             projectID,
			Name:           projectName,
			OrganizationId: organizationID,
		},
		Organization: &authorization.Organization{
			Id:   organizationID,
			Name: organizationName,
		},
		User: &authorization.User{
			Id:                 userID,
			PreferredLoginName: userResp.User.GetPreferredLoginName(),
			DisplayName:        userResp.User.GetHuman().GetProfile().GetDisplayName(),
			AvatarUrl:          userResp.User.GetHuman().GetProfile().GetAvatarUrl(),
			OrganizationId:     userResp.GetUser().GetDetails().GetResourceOwner(),
		},
		State: 1,
		Roles: roles,
	}
}

func createAuthorizationForProjectGrant(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, userID, projectID, projectName, grantedOrganizationID, grantedOrganizationName string, roles []*authorization.Role) *authorization.Authorization {
	userResp, err := instance.Client.UserV2.GetUserByID(ctx, &user.GetUserByIDRequest{UserId: userID})
	require.NoError(t, err)

	roleKeys := make([]string, len(roles))
	for i, role := range roles {
		roleKeys[i] = role.Key
	}
	authResp := instance.CreateAuthorizationProjectGrant(t, ctx, projectID, grantedOrganizationID, userID, roleKeys...)
	return &authorization.Authorization{
		Id:           authResp.GetId(),
		CreationDate: authResp.GetCreationDate(),
		ChangeDate:   authResp.GetCreationDate(),
		Project: &authorization.Project{
			Id:             projectID,
			Name:           projectName,
			OrganizationId: orgID,
		},
		Organization: &authorization.Organization{
			Id:   grantedOrganizationID,
			Name: grantedOrganizationName,
		},
		User: &authorization.User{
			Id:                 userID,
			PreferredLoginName: userResp.User.GetPreferredLoginName(),
			DisplayName:        userResp.User.GetHuman().GetProfile().GetDisplayName(),
			AvatarUrl:          userResp.User.GetHuman().GetProfile().GetAvatarUrl(),
			OrganizationId:     userResp.GetUser().GetDetails().GetResourceOwner(),
		},
		State: 1,
		Roles: roles,
	}
}

func TestServer_ListAuthorizations_PermissionsV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, InstancePermissionV2)
	iamOwnerCtx := InstancePermissionV2.WithAuthorizationToken(EmptyCTX, integration.UserTypeIAMOwner)

	organizationName := integration.OrganizationName()
	organization := InstancePermissionV2.CreateOrganization(iamOwnerCtx, organizationName, integration.Email())

	projectName1 := integration.ProjectName()
	project1 := InstancePermissionV2.CreateProject(iamOwnerCtx, t, organization.GetOrganizationId(), projectName1, false, false)

	projectName2 := integration.ProjectName()
	project2 := InstancePermissionV2.CreateProject(iamOwnerCtx, t, organization.GetOrganizationId(), projectName2, false, false)

	roles := []*authorization.Role{
		{Key: "role1", DisplayName: "Role 1", Group: "group1"},
		{Key: "role2", DisplayName: "Role 2", Group: "group2"},
		{Key: "role3", DisplayName: "Role 3", Group: "group2"}}
	for _, role := range roles {
		InstancePermissionV2.AddProjectRole(iamOwnerCtx, t, project1.GetId(), role.Key, role.DisplayName, role.Group)
		InstancePermissionV2.AddProjectRole(iamOwnerCtx, t, project2.GetId(), role.Key, role.DisplayName, role.Group)
	}

	grantedOrganizationName := integration.OrganizationName()
	grantedOrganization := InstancePermissionV2.CreateOrganization(iamOwnerCtx, grantedOrganizationName, integration.Email())

	grantedRoles := roles[0:2]
	grantedRoleKeys := make([]string, len(grantedRoles))
	for i, role := range grantedRoles {
		grantedRoleKeys[i] = role.Key
	}

	InstancePermissionV2.CreateProjectGrant(iamOwnerCtx, t, project2.GetId(), grantedOrganization.GetOrganizationId(), grantedRoleKeys...)

	user1 := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
	user2 := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
	user3 := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

	authorizationUser1Project1 := createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, organization.GetOrganizationId(), organizationName, user1.GetId(), project1.GetId(), projectName1, roles[0:1])
	authorizationUser2Project1 := createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, organization.GetOrganizationId(), organizationName, user2.GetId(), project1.GetId(), projectName1, roles[0:1])
	authorizationUser2Project2 := createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, organization.GetOrganizationId(), organizationName, user2.GetId(), project2.GetId(), projectName2, roles[1:2])
	authorizationUser3Project2 := createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, organization.GetOrganizationId(), organizationName, user3.GetId(), project2.GetId(), projectName2, roles[2:3])
	authorizationUser3ProjectGrant := createAuthorizationForProjectGrant(iamOwnerCtx, InstancePermissionV2, t, organization.GetOrganizationId(), user3.GetId(), project2.GetId(), projectName2, grantedOrganization.GetOrganizationId(), grantedOrganizationName, grantedRoles)

	organizationOwnerResp := InstancePermissionV2.CreateMachineUser(iamOwnerCtx)
	organizationOwnerPatResp := InstancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, organizationOwnerResp.GetUserId())
	InstancePermissionV2.CreateOrgMembership(t, iamOwnerCtx, organization.GetOrganizationId(), organizationOwnerResp.GetUserId())
	organizationOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, organizationOwnerPatResp.Token)

	projectOwnerResp := InstancePermissionV2.CreateMachineUser(iamOwnerCtx)
	projectOwnerPatResp := InstancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, projectOwnerResp.GetUserId())
	InstancePermissionV2.CreateProjectMembership(t, iamOwnerCtx, project2.GetId(), projectOwnerResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectOwnerPatResp.Token)

	projectGrantOwnerResp := InstancePermissionV2.CreateMachineUser(iamOwnerCtx)
	projectGrantOwnerPatResp := InstancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, projectGrantOwnerResp.GetUserId())
	InstancePermissionV2.CreateProjectGrantMembership(t, iamOwnerCtx, project2.GetId(), grantedOrganization.GetOrganizationId(), projectGrantOwnerResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectGrantOwnerPatResp.Token)

	type args struct {
		ctx context.Context
		req *authorization.ListAuthorizationsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *authorization.ListAuthorizationsResponse
		wantErr bool
	}{
		{
			name: "unauthenticated",
			args: args{
				ctx: EmptyCTX,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user1.GetId(),
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			args: args{
				ctx: InstancePermissionV2.WithAuthorizationToken(EmptyCTX, integration.UserTypeNoPermission),
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user1.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: nil,
			},
		},
		{
			name: "not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: "notexisting",
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "limit and pagination",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Pagination: &filter.PaginationRequest{
						Limit:  3,
						Offset: 1,
						Asc:    true,
					},
					SortingColumn: authorization.AuthorizationFieldName_AUTHORIZATION_FIELD_NAME_CREATED_DATE,
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5,
					AppliedLimit: 3,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser2Project1,
					authorizationUser2Project2,
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "organization owner permission, only expect non granted",
			args: args{
				ctx: organizationOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // count issue with permission check
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "project owner permission, only expect own project non granted",
			args: args{
				ctx: projectOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // count issue with permission check
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "project grant owner permission, only expect granted",
			args: args{
				ctx: projectGrantOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // count issue with permission check
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
				},
			},
		},
		{
			name: "filter by single authorizationID, expect one result",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_AuthorizationIds{
								AuthorizationIds: &filter.InIDsFilter{
									Ids: []string{authorizationUser1Project1.GetId()},
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter by organization ID (expect only own / non granted)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_OrganizationId{
								OrganizationId: &filter.IDFilter{
									Id: organization.GetOrganizationId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3Project2,
					authorizationUser2Project2,
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter by organization ID (only project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_OrganizationId{
								OrganizationId: &filter.IDFilter{
									Id: grantedOrganization.GetOrganizationId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
				},
			},
		},
		{
			name: "filter by userID",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user1.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "list by userID (where project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_UserId{
								UserId: &filter.IDFilter{
									Id: user3.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
					authorizationUser3Project2,
				},
			},
		},
		{
			name: "filter by projectID (no project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_ProjectId{
								ProjectId: &filter.IDFilter{
									Id: project1.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter by projectID (project grant exists), expect own and granted",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_ProjectId{
								ProjectId: &filter.IDFilter{
									Id: project2.GetId(),
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
					authorizationUser3Project2,
					authorizationUser2Project2,
				},
			},
		},
		{
			name: "filter by project name (no project grant exists)",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_ProjectName{
								ProjectName: &authorization.ProjectNameQuery{
									Name:   projectName1,
									Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
		{
			name: "filter role key, expect all with that role, own and granted",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_RoleKey{
								RoleKey: &authorization.RoleKeyQuery{
									Key:    roles[0].Key,
									Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					authorizationUser3ProjectGrant,
					authorizationUser2Project1,
					authorizationUser1Project1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := InstancePermissionV2.Client.AuthorizationV2.ListAuthorizations(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Authorizations, len(tt.want.Authorizations)) {
					assert.EqualExportedValues(ttt, tt.want.Authorizations, got.Authorizations)
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}
