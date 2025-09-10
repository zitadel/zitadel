//go:build integration

package authorization_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_ListAuthorizations(t *testing.T) {
	iamOwnerCtx := Instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeIAMOwner)
	projectOwnerResp := Instance.CreateMachineUser(iamOwnerCtx)
	projectOwnerPatResp := Instance.CreatePersonalAccessToken(iamOwnerCtx, projectOwnerResp.GetUserId())
	projectResp := createProject(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), false, false)
	Instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), projectOwnerResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectOwnerPatResp.Token)

	projectGrantOwnerResp := Instance.CreateMachineUser(iamOwnerCtx)
	projectGrantOwnerPatResp := Instance.CreatePersonalAccessToken(iamOwnerCtx, projectGrantOwnerResp.GetUserId())
	grantedProjectResp := createGrantedProject(iamOwnerCtx, Instance, t, projectResp)
	Instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), grantedProjectResp.GetGrantedOrganizationId(), projectGrantOwnerResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectGrantOwnerPatResp.Token)

	type args struct {
		ctx context.Context
		dep func(*authorization.ListAuthorizationsRequest, *authorization.ListAuthorizationsResponse)
		req *authorization.ListAuthorizationsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *authorization.ListAuthorizationsResponse
		wantErr bool
	}{
		{
			name: "list by user id, unauthenticated",
			args: args{
				ctx: EmptyCTX,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: Instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeNoPermission),
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{Filter: &authorization.AuthorizationsSearchFilter_UserId{
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
			name: "list single id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_AuthorizationIds{
						AuthorizationIds: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single project id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_ProjectId{
						ProjectId: &filter.IDFilter{
							Id: resp.GetProjectId(),
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single project name",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_ProjectName{
						ProjectName: &authorization.ProjectNameQuery{
							Name:   resp.GetProjectName(),
							Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), true)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single grant id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), true)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_ProjectGrantId{
						ProjectGrantId: &filter.IDFilter{
							Id: resp.GetProjectGrantId(),
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project and project grant, multiple",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[5] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
					response.Authorizations[4] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
					response.Authorizations[3] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), false)
					response.Authorizations[2] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), true)
					response.Authorizations[1] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), true)
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), true)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  6,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{}, {}, {}, {}, {}, {},
				},
			},
		},
		{
			name: "list single id, project and project grant, org owner",
			args: args{
				ctx: Instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeOrgOwner),
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}

					response.Authorizations[1] = createAuthorizationForProject(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), projectResp.GetName(), projectResp.GetId())
					response.Authorizations[0] = createAuthorizationWithProjectGrant(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), grantedProjectResp.GetName(), grantedProjectResp.GetId())
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{}, {},
				},
			},
		},
		{
			name: "list single id, project and project grant, project owner",
			args: args{
				ctx: projectOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}

					response.Authorizations[0] = createAuthorizationForProject(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), projectResp.GetName(), projectResp.GetId())
					createAuthorizationWithProjectGrant(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), grantedProjectResp.GetName(), grantedProjectResp.GetId())
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project and project grant, project grant owner",
			args: args{
				ctx: projectGrantOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := Instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}

					createAuthorizationForProject(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), projectResp.GetName(), projectResp.GetId())
					response.Authorizations[0] = createAuthorizationForProjectGrant(iamOwnerCtx, Instance, t, Instance.DefaultOrg.GetId(), userResp.GetId(), grantedProjectResp.GetName(), grantedProjectResp.GetId(), grantedProjectResp.GetGrantedOrganizationId())
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Instance.Client.AuthorizationV2Beta.ListAuthorizations(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Authorizations, len(tt.want.Authorizations)) {
					for i := range tt.want.Authorizations {
						assert.EqualExportedValues(ttt, tt.want.Authorizations[i], got.Authorizations[i])
					}
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

func createAuthorization(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, userID string, grant bool) *authorization.Authorization {
	projectName := integration.ProjectName()
	projectResp := instance.CreateProject(ctx, t, orgID, projectName, false, false)

	if grant {
		return createAuthorizationWithProjectGrant(ctx, instance, t, orgID, userID, projectName, projectResp.GetId())
	}
	return createAuthorizationForProject(ctx, instance, t, orgID, userID, projectName, projectResp.GetId())
}

func createAuthorizationForProject(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, userID, projectName, projectID string) *authorization.Authorization {
	userResp, err := instance.Client.UserV2.GetUserByID(ctx, &user.GetUserByIDRequest{UserId: userID})
	require.NoError(t, err)

	userGrantResp := instance.CreateProjectUserGrant(t, ctx, orgID, projectID, userID)
	return &authorization.Authorization{
		Id:                    userGrantResp.GetUserGrantId(),
		ProjectId:             projectID,
		ProjectName:           projectName,
		ProjectOrganizationId: orgID,
		OrganizationId:        orgID,
		CreationDate:          userGrantResp.Details.GetCreationDate(),
		ChangeDate:            userGrantResp.Details.GetCreationDate(),
		State:                 1,
		User: &authorization.User{
			Id:                 userID,
			PreferredLoginName: userResp.User.GetPreferredLoginName(),
			DisplayName:        userResp.User.GetHuman().GetProfile().GetDisplayName(),
			AvatarUrl:          userResp.User.GetHuman().GetProfile().GetAvatarUrl(),
			OrganizationId:     userResp.GetUser().GetDetails().GetResourceOwner(),
		},
	}
}

func createAuthorizationWithProjectGrant(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, userID, projectName, projectID string) *authorization.Authorization {
	grantedOrgName := integration.OrganizationName()
	grantedOrg := instance.CreateOrganization(ctx, grantedOrgName, integration.Email())
	instance.CreateProjectGrant(ctx, t, projectID, grantedOrg.GetOrganizationId())

	return createAuthorizationForProjectGrant(ctx, instance, t, orgID, userID, projectName, projectID, grantedOrg.GetOrganizationId())
}

func createAuthorizationForProjectGrant(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, userID, projectName, projectID, grantedOrgID string) *authorization.Authorization {
	userResp, err := instance.Client.UserV2.GetUserByID(ctx, &user.GetUserByIDRequest{UserId: userID})
	require.NoError(t, err)

	userGrantResp := instance.CreateProjectGrantUserGrant(ctx, orgID, projectID, grantedOrgID, userID)
	return &authorization.Authorization{
		Id:                    userGrantResp.GetUserGrantId(),
		ProjectId:             projectID,
		ProjectName:           projectName,
		ProjectOrganizationId: orgID,
		ProjectGrantId:        gu.Ptr(grantedOrgID),
		GrantedOrganizationId: gu.Ptr(grantedOrgID),
		OrganizationId:        orgID,
		CreationDate:          userGrantResp.Details.GetCreationDate(),
		ChangeDate:            userGrantResp.Details.GetCreationDate(),
		State:                 1,
		User: &authorization.User{
			Id:                 userID,
			PreferredLoginName: userResp.User.GetPreferredLoginName(),
			DisplayName:        userResp.User.GetHuman().GetProfile().GetDisplayName(),
			AvatarUrl:          userResp.User.GetHuman().GetProfile().GetAvatarUrl(),
			OrganizationId:     userResp.GetUser().GetDetails().GetResourceOwner(),
		},
	}
}

func createProject(ctx context.Context, instance *integration.Instance, t *testing.T, orgID string, projectRoleCheck, hasProjectCheck bool) *project.Project {
	name := integration.ProjectName()
	resp := instance.CreateProject(ctx, t, orgID, name, projectRoleCheck, hasProjectCheck)
	return &project.Project{
		Id:                     resp.GetId(),
		Name:                   name,
		OrganizationId:         orgID,
		CreationDate:           resp.GetCreationDate(),
		ChangeDate:             resp.GetCreationDate(),
		State:                  1,
		ProjectRoleAssertion:   false,
		ProjectAccessRequired:  hasProjectCheck,
		AuthorizationRequired:  projectRoleCheck,
		PrivateLabelingSetting: project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
	}
}

func createGrantedProject(ctx context.Context, instance *integration.Instance, t *testing.T, projectToGrant *project.Project) *project.Project {
	grantedOrgName := integration.OrganizationName()
	grantedOrg := instance.CreateOrganization(ctx, grantedOrgName, integration.Email())
	projectGrantResp := instance.CreateProjectGrant(ctx, t, projectToGrant.GetId(), grantedOrg.GetOrganizationId())

	return &project.Project{
		Id:                      projectToGrant.GetId(),
		Name:                    projectToGrant.GetName(),
		OrganizationId:          projectToGrant.GetOrganizationId(),
		CreationDate:            projectGrantResp.GetCreationDate(),
		ChangeDate:              projectGrantResp.GetCreationDate(),
		State:                   1,
		ProjectRoleAssertion:    false,
		ProjectAccessRequired:   projectToGrant.GetProjectAccessRequired(),
		AuthorizationRequired:   projectToGrant.GetAuthorizationRequired(),
		PrivateLabelingSetting:  project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
		GrantedOrganizationId:   gu.Ptr(grantedOrg.GetOrganizationId()),
		GrantedOrganizationName: gu.Ptr(grantedOrgName),
		GrantedState:            1,
	}
}

func TestServer_ListAuthorizations_PermissionsV2(t *testing.T) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// ensureFeaturePermissionV2Enabled(t, InstancePermissionV2)
	iamOwnerCtx := InstancePermissionV2.WithAuthorizationToken(EmptyCTX, integration.UserTypeIAMOwner)

	projectOwnerResp := InstancePermissionV2.CreateMachineUser(iamOwnerCtx)
	projectOwnerPatResp := InstancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, projectOwnerResp.GetUserId())
	projectResp := createProject(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), false, false)
	InstancePermissionV2.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), projectOwnerResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectOwnerPatResp.Token)

	projectGrantOwnerResp := InstancePermissionV2.CreateMachineUser(iamOwnerCtx)
	projectGrantOwnerPatResp := InstancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, projectGrantOwnerResp.GetUserId())
	grantedProjectResp := createGrantedProject(iamOwnerCtx, InstancePermissionV2, t, projectResp)
	InstancePermissionV2.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), grantedProjectResp.GetGrantedOrganizationId(), projectGrantOwnerResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(EmptyCTX, projectGrantOwnerPatResp.Token)

	type args struct {
		ctx context.Context
		dep func(*authorization.ListAuthorizationsRequest, *authorization.ListAuthorizationsResponse)
		req *authorization.ListAuthorizationsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *authorization.ListAuthorizationsResponse
		wantErr bool
	}{
		{
			name: "list by user id, unauthenticated",
			args: args{
				ctx: EmptyCTX,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: InstancePermissionV2.WithAuthorizationToken(EmptyCTX, integration.UserTypeNoPermission),
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{Filter: &authorization.AuthorizationsSearchFilter_UserId{
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
			name: "list single id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_AuthorizationIds{
						AuthorizationIds: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single project id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_ProjectId{
						ProjectId: &filter.IDFilter{
							Id: resp.GetProjectId(),
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single project name",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_ProjectName{
						ProjectName: &authorization.ProjectNameQuery{
							Name:   resp.GetProjectName(),
							Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), true)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single grant id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					resp := createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), true)
					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_ProjectGrantId{
						ProjectGrantId: &filter.IDFilter{
							Id: resp.GetProjectGrantId(),
						},
					}
					response.Authorizations[0] = resp
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project and project grant, multiple",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}
					response.Authorizations[5] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
					response.Authorizations[4] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
					response.Authorizations[3] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), false)
					response.Authorizations[2] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), true)
					response.Authorizations[1] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), true)
					response.Authorizations[0] = createAuthorization(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), true)
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  6,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{}, {}, {}, {}, {}, {},
				},
			},
		},
		{
			name: "list single id, project and project grant, org owner",
			args: args{
				ctx: InstancePermissionV2.WithAuthorizationToken(EmptyCTX, integration.UserTypeOrgOwner),
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}

					response.Authorizations[1] = createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), projectResp.GetName(), projectResp.GetId())
					response.Authorizations[0] = createAuthorizationWithProjectGrant(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), grantedProjectResp.GetName(), grantedProjectResp.GetId())
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{}, {},
				},
			},
		},
		{
			name: "list single id, project and project grant, project owner",
			args: args{
				ctx: projectOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}

					response.Authorizations[0] = createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), projectResp.GetName(), projectResp.GetId())
					createAuthorizationWithProjectGrant(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), grantedProjectResp.GetName(), grantedProjectResp.GetId())
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
		{
			name: "list single id, project and project grant, project grant owner",
			args: args{
				ctx: projectGrantOwnerCtx,
				dep: func(request *authorization.ListAuthorizationsRequest, response *authorization.ListAuthorizationsResponse) {
					userResp := InstancePermissionV2.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

					request.Filters[0].Filter = &authorization.AuthorizationsSearchFilter_UserId{
						UserId: &filter.IDFilter{
							Id: userResp.GetId(),
						},
					}

					createAuthorizationForProject(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), projectResp.GetName(), projectResp.GetId())
					response.Authorizations[0] = createAuthorizationForProjectGrant(iamOwnerCtx, InstancePermissionV2, t, InstancePermissionV2.DefaultOrg.GetId(), userResp.GetId(), grantedProjectResp.GetName(), grantedProjectResp.GetId(), grantedProjectResp.GetGrantedOrganizationId())
				},
				req: &authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{{}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Authorizations: []*authorization.Authorization{
					{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := InstancePermissionV2.Client.AuthorizationV2Beta.ListAuthorizations(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Authorizations, len(tt.want.Authorizations)) {
					for i := range tt.want.Authorizations {
						assert.EqualExportedValues(ttt, tt.want.Authorizations[i], got.Authorizations[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}
