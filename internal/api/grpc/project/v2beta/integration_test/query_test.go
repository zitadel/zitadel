//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_GetProject(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(*project.GetProjectRequest, *project.GetProjectResponse)
		req *project.GetProjectRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.GetProjectResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := createProject(iamOwnerCtx, instance, t, orgID, false, false)

					request.Id = resp.GetId()
				},
				req: &project.GetProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "missing permission, other org owner",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)

					request.Id = resp.GetId()
				},
				req: &project.GetProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.GetProjectRequest{Id: "notexisting"},
			},
			wantErr: true,
		},
		{
			name: "get, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := createProject(iamOwnerCtx, instance, t, orgID, false, false)

					request.Id = resp.GetId()
					response.Project = resp
				},
				req: &project.GetProjectRequest{},
			},
			want: &project.GetProjectResponse{
				Project: &project.Project{
					State:                  1,
					ProjectRoleAssertion:   false,
					ProjectAccessRequired:  false,
					AuthorizationRequired:  false,
					PrivateLabelingSetting: project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
				},
			},
		},
		{
			name: "get, ok, org owner",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := createProject(iamOwnerCtx, instance, t, orgID, false, false)

					request.Id = resp.GetId()
					response.Project = resp
				},
				req: &project.GetProjectRequest{},
			},
			want: &project.GetProjectResponse{
				Project: &project.Project{},
			},
		},
		{
			name: "get, full, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := createProject(iamOwnerCtx, instance, t, orgID, true, true)

					request.Id = resp.GetId()
					response.Project = resp
				},
				req: &project.GetProjectRequest{},
			},
			want: &project.GetProjectResponse{
				Project: &project.Project{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, 2*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.Projectv2Beta.GetProject(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)
				assert.EqualExportedValues(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected project result")
		})
	}
}

func TestServer_ListProjects(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	userResp := instance.CreateMachineUser(iamOwnerCtx)
	patResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userResp.GetUserId())
	projectResp := createProject(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), false, false)
	instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetUserId())
	grantedProjectResp := createGrantedProject(iamOwnerCtx, instance, t, projectResp)
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patResp.Token)

	type args struct {
		ctx context.Context
		dep func(*project.ListProjectsRequest, *project.ListProjectsResponse)
		req *project.ListProjectsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.ListProjectsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, integration.ProjectName(), false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, integration.ProjectName(), false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{
						{Filter: &project.ProjectSearchFilter_InProjectIdsFilter{
							InProjectIdsFilter: &filter.InIDsFilter{
								Ids: []string{"notfound"},
							},
						},
						},
					},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					response.Projects[0] = createProject(iamOwnerCtx, instance, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{response.Projects[0].GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list single name",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					response.Projects[0] = createProject(iamOwnerCtx, instance, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectNameFilter{
						ProjectNameFilter: &project.ProjectNameFilter{
							ProjectName: response.Projects[0].Name,
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{
						State:                  1,
						ProjectRoleAssertion:   false,
						ProjectAccessRequired:  false,
						AuthorizationRequired:  false,
						PrivateLabelingSetting: project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
					},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					response.Projects[2] = createProject(iamOwnerCtx, instance, t, orgID, false, false)
					response.Projects[1] = createProject(iamOwnerCtx, instance, t, orgID, true, false)
					response.Projects[0] = createProject(iamOwnerCtx, instance, t, orgID, false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{response.Projects[0].GetId(), response.Projects[1].GetId(), response.Projects[2].GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
					{},
				},
			},
		},
		{
			name: "list multiple id, limited permissions",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp1 := createProject(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), false, false)
					resp2 := createProject(iamOwnerCtx, instance, t, orgID, true, false)
					resp3 := createProject(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp1.GetId(), resp2.GetId(), resp3.GetId()},
						},
					}

					response.Projects[0] = resp2
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list multiple id, limited permissions, project owner",
			args: args{
				ctx: projectOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp1 := createProject(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), false, false)
					resp2 := createProject(iamOwnerCtx, instance, t, orgID, true, false)
					resp3 := createProject(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp1.GetId(), resp2.GetId(), resp3.GetId(), projectResp.GetId()}},
					}
					response.Projects[0] = grantedProjectResp
					response.Projects[1] = projectResp
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
				},
			},
		},
		{
			name: "list project and granted projects",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					projectResp := createProject(iamOwnerCtx, instance, t, orgID, true, true)
					response.Projects[3] = projectResp
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					response.Projects[2] = createGrantedProject(iamOwnerCtx, instance, t, projectResp)
					response.Projects[1] = createGrantedProject(iamOwnerCtx, instance, t, projectResp)
					response.Projects[0] = createGrantedProject(iamOwnerCtx, instance, t, projectResp)
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
					{},
					{},
				},
			},
		},
		{
			name: "list project and granted projects, organization",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					projectResp := createProject(iamOwnerCtx, instance, t, orgID, true, true)

					grantedProjectResp := createGrantedProject(iamOwnerCtx, instance, t, projectResp)
					response.Projects[1] = grantedProjectResp
					response.Projects[0] = createProject(iamOwnerCtx, instance, t, *grantedProjectResp.GrantedOrganizationId, true, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectOrganizationIdFilter{
						ProjectOrganizationIdFilter: &filter.IDFilter{Id: *grantedProjectResp.GrantedOrganizationId},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
				},
			},
		},
		{
			name: "list project and granted projects, project resourceowner",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					projectResp := createProject(iamOwnerCtx, instance, t, orgID, true, true)

					grantedProjectResp := createGrantedProject(iamOwnerCtx, instance, t, projectResp)
					response.Projects[0] = createProject(iamOwnerCtx, instance, t, *grantedProjectResp.GrantedOrganizationId, true, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectResourceOwnerFilter{
						ProjectResourceOwnerFilter: &filter.IDFilter{Id: *grantedProjectResp.GrantedOrganizationId},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list granted project, project id",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()

					orgName := integration.OrganizationName()
					projectName := integration.ProjectName()
					orgResp := instance.CreateOrganization(iamOwnerCtx, orgName, integration.Email())
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), projectName, true, true)
					projectGrantResp := instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgID)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{Ids: []string{projectResp.GetId()}},
					}
					response.Projects[0] = &project.Project{
						Id:                      projectResp.GetId(),
						Name:                    projectName,
						OrganizationId:          orgResp.GetOrganizationId(),
						CreationDate:            projectGrantResp.GetCreationDate(),
						ChangeDate:              projectGrantResp.GetCreationDate(),
						State:                   1,
						ProjectRoleAssertion:    false,
						ProjectAccessRequired:   true,
						AuthorizationRequired:   true,
						PrivateLabelingSetting:  project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
						GrantedOrganizationId:   gu.Ptr(orgID),
						GrantedOrganizationName: gu.Ptr(instance.DefaultOrg.GetName()),
						GrantedState:            1,
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
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
				got, listErr := instance.Client.Projectv2Beta.ListProjects(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Projects, len(tt.want.Projects)) {
					for i := range tt.want.Projects {
						assert.EqualExportedValues(ttt, tt.want.Projects[i], got.Projects[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ListProjects_PermissionV2(t *testing.T) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgID := instancePermissionV2.DefaultOrg.GetId()

	type args struct {
		ctx context.Context
		dep func(*project.ListProjectsRequest, *project.ListProjectsResponse)
		req *project.ListProjectsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.ListProjectsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					resp := createProject(iamOwnerCtx, instancePermissionV2, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					resp := createProject(iamOwnerCtx, instancePermissionV2, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp := createProject(iamOwnerCtx, instancePermissionV2, t, orgResp.GetOrganizationId(), false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp.GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{
						{Filter: &project.ProjectSearchFilter_InProjectIdsFilter{
							InProjectIdsFilter: &filter.InIDsFilter{
								Ids: []string{"notfound"},
							},
						},
						},
					},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					response.Projects[0] = createProject(iamOwnerCtx, instancePermissionV2, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{response.Projects[0].GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list single name",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					response.Projects[0] = createProject(iamOwnerCtx, instancePermissionV2, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectNameFilter{
						ProjectNameFilter: &project.ProjectNameFilter{
							ProjectName: response.Projects[0].Name,
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					response.Projects[2] = createProject(iamOwnerCtx, instancePermissionV2, t, orgID, false, false)
					response.Projects[1] = createProject(iamOwnerCtx, instancePermissionV2, t, orgID, true, false)
					response.Projects[0] = createProject(iamOwnerCtx, instancePermissionV2, t, orgID, false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{response.Projects[0].GetId(), response.Projects[1].GetId(), response.Projects[2].GetId()},
						},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
					{},
				},
			},
		},
		{
			name: "list project and granted projects",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					projectResp := createProject(iamOwnerCtx, instancePermissionV2, t, orgID, true, true)
					response.Projects[3] = projectResp
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					response.Projects[2] = createGrantedProject(iamOwnerCtx, instancePermissionV2, t, projectResp)
					response.Projects[1] = createGrantedProject(iamOwnerCtx, instancePermissionV2, t, projectResp)
					response.Projects[0] = createGrantedProject(iamOwnerCtx, instancePermissionV2, t, projectResp)
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
					{},
					{},
				},
			},
		},
		{
			name: "list project and granted projects, organization",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					projectResp := createProject(iamOwnerCtx, instancePermissionV2, t, orgID, true, true)

					grantedProjectResp := createGrantedProject(iamOwnerCtx, instancePermissionV2, t, projectResp)
					response.Projects[1] = grantedProjectResp
					response.Projects[0] = createProject(iamOwnerCtx, instancePermissionV2, t, *grantedProjectResp.GrantedOrganizationId, true, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectOrganizationIdFilter{
						ProjectOrganizationIdFilter: &filter.IDFilter{Id: *grantedProjectResp.GrantedOrganizationId},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
					{},
				},
			},
		},
		{
			name: "list project and granted projects, project resourceowner",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					projectResp := createProject(iamOwnerCtx, instancePermissionV2, t, orgID, true, true)

					grantedProjectResp := createGrantedProject(iamOwnerCtx, instancePermissionV2, t, projectResp)
					response.Projects[0] = createProject(iamOwnerCtx, instancePermissionV2, t, *grantedProjectResp.GrantedOrganizationId, true, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectResourceOwnerFilter{
						ProjectResourceOwnerFilter: &filter.IDFilter{Id: *grantedProjectResp.GrantedOrganizationId},
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list multiple id, limited permissions",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					resp1 := createProject(iamOwnerCtx, instancePermissionV2, t, orgResp.GetOrganizationId(), false, false)
					resp2 := createProject(iamOwnerCtx, instancePermissionV2, t, orgID, true, false)
					resp3 := createProject(iamOwnerCtx, instancePermissionV2, t, orgResp.GetOrganizationId(), false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{resp1.GetId(), resp2.GetId(), resp3.GetId()},
						},
					}

					response.Projects[0] = resp2
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{
					{},
				},
			},
		},
		{
			name: "list granted project, project id",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instancePermissionV2.DefaultOrg.GetId()

					orgName := integration.OrganizationName()
					projectName := integration.ProjectName()
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, orgName, integration.Email())
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), projectName, true, true)
					projectGrantResp := instancePermissionV2.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgID)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{Ids: []string{projectResp.GetId()}},
					}

					response.Projects[0] = &project.Project{
						Id:                      projectResp.GetId(),
						Name:                    projectName,
						OrganizationId:          orgResp.GetOrganizationId(),
						CreationDate:            projectGrantResp.GetCreationDate(),
						ChangeDate:              projectGrantResp.GetCreationDate(),
						State:                   1,
						ProjectRoleAssertion:    false,
						ProjectAccessRequired:   true,
						AuthorizationRequired:   true,
						PrivateLabelingSetting:  project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
						GrantedOrganizationId:   gu.Ptr(orgID),
						GrantedOrganizationName: gu.Ptr(instancePermissionV2.DefaultOrg.GetName()),
						GrantedState:            1,
					}
				},
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{{}},
				},
			},
			want: &project.ListProjectsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
				Projects: []*project.Project{{}},
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
				got, listErr := instancePermissionV2.Client.Projectv2Beta.ListProjects(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Projects, len(tt.want.Projects)) {
					for i := range tt.want.Projects {
						assert.EqualExportedValues(ttt, tt.want.Projects[i], got.Projects[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
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

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func TestServer_ListProjectGrants(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	userResp := instance.CreateMachineUser(iamOwnerCtx)
	patResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userResp.GetUserId())
	projectResp := createProject(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), false, false)
	projectGrantResp := createProjectGrant(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectResp.GetName())
	instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), projectGrantResp.GetGrantedOrganizationId(), userResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patResp.Token)

	type args struct {
		ctx context.Context
		dep func(*project.ListProjectGrantsRequest, *project.ListProjectGrantsResponse)
		req *project.ListProjectGrantsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.ListProjectGrantsResponse
		wantErr bool
	}{{
		name: "list by id, unauthenticated",
		args: args{
			ctx: CTX,
			dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
					InProjectIdsFilter: &filter.InIDsFilter{
						Ids: []string{projectResp.GetId()},
					},
				}
				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.Filters[1].Filter = &project.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter{
					ProjectGrantResourceOwnerFilter: &filter.IDFilter{Id: grantedOrg.GetOrganizationId()},
				}

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			req: &project.ListProjectGrantsRequest{
				Filters: []*project.ProjectGrantSearchFilter{{}, {}},
			},
		},
		wantErr: true,
	},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					request.Filters[1].Filter = &project.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter{
						ProjectGrantResourceOwnerFilter: &filter.IDFilter{Id: grantedOrg.GetOrganizationId()},
					}

					instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}, {}},
				},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{
						{Filter: &project.ProjectGrantSearchFilter_InProjectIdsFilter{
							InProjectIdsFilter: &filter.InIDsFilter{
								Ids: []string{"notfound"},
							},
						},
						},
					},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list by id",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}

					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instance, t, orgID, projectResp.GetId(), name)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{},
				},
			},
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}

					createProjectGrant(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), projectResp.GetId(), name)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}

					response.ProjectGrants[2] = createProjectGrant(iamOwnerCtx, instance, t, orgID, projectResp.GetId(), name)
					response.ProjectGrants[1] = createProjectGrant(iamOwnerCtx, instance, t, orgID, projectResp.GetId(), name)
					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instance, t, orgID, projectResp.GetId(), name)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, limited permissions",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name1 := integration.ProjectName()
					name2 := integration.ProjectName()
					name3 := integration.ProjectName()
					orgID := instance.DefaultOrg.GetId()
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					project1Resp := instance.CreateProject(iamOwnerCtx, t, orgID, name1, false, false)
					project2Resp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name2, false, false)
					project3Resp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name3, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{project1Resp.GetId(), project2Resp.GetId(), project3Resp.GetId()},
						},
					}

					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instance, t, orgID, project1Resp.GetId(), name1)
					createProjectGrant(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), project2Resp.GetId(), name2)
					createProjectGrant(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), project3Resp.GetId(), name3)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{},
				},
			},
		},
		{
			name: "list multiple id, limited permissions, project grant owner",
			args: args{
				ctx: projectGrantOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name1 := integration.ProjectName()
					name2 := integration.ProjectName()
					name3 := integration.ProjectName()
					orgID := instance.DefaultOrg.GetId()
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					project1Resp := instance.CreateProject(iamOwnerCtx, t, orgID, name1, false, false)
					project2Resp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name2, false, false)
					project3Resp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name3, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{project1Resp.GetId(), project2Resp.GetId(), project3Resp.GetId(), projectResp.GetId()}},
					}

					createProjectGrant(iamOwnerCtx, instance, t, orgID, project1Resp.GetId(), name1)
					createProjectGrant(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), project2Resp.GetId(), name2)
					createProjectGrant(iamOwnerCtx, instance, t, orgResp.GetOrganizationId(), project3Resp.GetId(), name3)
					response.ProjectGrants[0] = projectGrantResp
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{},
				},
			},
		},
		{
			name: "list single id with role",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					projectRoleResp := addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instance, t, orgID, projectResp.GetId(), name, projectRoleResp.Key)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{},
				},
			},
		},
		{
			name: "list single id with removed role",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					projectRoleResp := addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instance, t, orgID, projectResp.GetId(), name, projectRoleResp.Key)

					removeRoleResp := instance.RemoveProjectRole(iamOwnerCtx, t, projectResp.GetId(), projectRoleResp.Key)
					response.ProjectGrants[0].GrantedRoleKeys = nil
					response.ProjectGrants[0].ChangeDate = removeRoleResp.GetRemovalDate()
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
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
				got, listErr := instance.Client.Projectv2Beta.ListProjectGrants(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.ProjectGrants, len(tt.want.ProjectGrants)) {
					for i := range tt.want.ProjectGrants {
						assert.EqualExportedValues(ttt, tt.want.ProjectGrants[i], got.ProjectGrants[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ListProjectGrants_PermissionV2(t *testing.T) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(*project.ListProjectGrantsRequest, *project.ListProjectGrantsResponse)
		req *project.ListProjectGrantsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.ListProjectGrantsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, instancePermissionV2.DefaultOrg.GetId(), integration.ProjectName(), false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					grantedOrg := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					request.Filters[1].Filter = &project.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter{
						ProjectGrantResourceOwnerFilter: &filter.IDFilter{Id: grantedOrg.GetOrganizationId()},
					}

					instancePermissionV2.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}, {}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, instancePermissionV2.DefaultOrg.GetId(), integration.ProjectName(), false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}
					grantedOrg := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					request.Filters[1].Filter = &project.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter{
						ProjectGrantResourceOwnerFilter: &filter.IDFilter{Id: grantedOrg.GetOrganizationId()},
					}

					instancePermissionV2.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}, {}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgID := instancePermissionV2.DefaultOrg.GetId()
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}

					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgID, projectResp.GetId(), name)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{},
				},
			},
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}

					createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgResp.GetOrganizationId(), projectResp.GetId(), name)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := integration.ProjectName()
					orgID := instancePermissionV2.DefaultOrg.GetId()
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectResp.GetId()},
						},
					}

					response.ProjectGrants[2] = createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgID, projectResp.GetId(), name)
					response.ProjectGrants[1] = createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgID, projectResp.GetId(), name)
					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgID, projectResp.GetId(), name)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, limited permissions",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name1 := integration.ProjectName()
					name2 := integration.ProjectName()
					name3 := integration.ProjectName()
					orgID := instancePermissionV2.DefaultOrg.GetId()
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					project1Resp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgID, name1, false, false)
					project2Resp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name2, false, false)
					project3Resp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name3, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &filter.InIDsFilter{
							Ids: []string{project1Resp.GetId(), project2Resp.GetId(), project3Resp.GetId()},
						},
					}

					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgID, project1Resp.GetId(), name1)
					createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgResp.GetOrganizationId(), project2Resp.GetId(), name2)
					createProjectGrant(iamOwnerCtx, instancePermissionV2, t, orgResp.GetOrganizationId(), project3Resp.GetId(), name3)
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
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
				got, listErr := instancePermissionV2.Client.Projectv2Beta.ListProjectGrants(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.ProjectGrants, len(tt.want.ProjectGrants)) {
					for i := range tt.want.ProjectGrants {
						assert.EqualExportedValues(ttt, tt.want.ProjectGrants[i], got.ProjectGrants[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func createProjectGrant(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, projectID, projectName string, roles ...string) *project.ProjectGrant {
	grantedOrgName := integration.ProjectName()
	grantedOrg := instance.CreateOrganization(ctx, grantedOrgName, integration.Email())
	projectGrantResp := instance.CreateProjectGrant(ctx, t, projectID, grantedOrg.GetOrganizationId(), roles...)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		projects, err := instance.Client.Projectv2Beta.ListProjects(ctx, &project.ListProjectsRequest{
			Filters: []*project.ProjectSearchFilter{{
				Filter: &project.ProjectSearchFilter_InProjectIdsFilter{
					InProjectIdsFilter: &filter.InIDsFilter{Ids: []string{projectID}},
				},
			}, {
				Filter: &project.ProjectSearchFilter_ProjectOrganizationIdFilter{
					ProjectOrganizationIdFilter: &filter.IDFilter{Id: grantedOrg.GetOrganizationId()},
				},
			}},
		})
		require.NoError(collect, err)
		require.Len(collect, projects.Projects, 1)
	}, retryDuration, tick)

	return &project.ProjectGrant{
		OrganizationId:          orgID,
		CreationDate:            projectGrantResp.GetCreationDate(),
		ChangeDate:              projectGrantResp.GetCreationDate(),
		GrantedOrganizationId:   grantedOrg.GetOrganizationId(),
		GrantedOrganizationName: grantedOrgName,
		ProjectId:               projectID,
		ProjectName:             projectName,
		State:                   1,
		GrantedRoleKeys:         roles,
	}
}

func TestServer_ListProjectRoles(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		dep func(*project.ListProjectRolesRequest, *project.ListProjectRolesResponse)
		req *project.ListProjectRolesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.ListProjectRolesResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{
					Filters: []*project.ProjectRoleSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{
					Filters: []*project.ProjectRoleSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ListProjectRolesRequest{
					ProjectId: "notfound",
				},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id, missing permission",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectRoles: []*project.ProjectRole{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					response.ProjectRoles[0] = addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectRoles: []*project.ProjectRole{
					{},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					response.ProjectRoles[2] = addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
					response.ProjectRoles[1] = addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
					response.ProjectRoles[0] = addProjectRole(iamOwnerCtx, instance, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				ProjectRoles: []*project.ProjectRole{
					{}, {}, {},
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
				got, listErr := instance.Client.Projectv2Beta.ListProjectRoles(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.ProjectRoles, len(tt.want.ProjectRoles)) {
					for i := range tt.want.ProjectRoles {
						assert.EqualExportedValues(ttt, tt.want.ProjectRoles[i], got.ProjectRoles[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func TestServer_ListProjectRoles_PermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(*project.ListProjectRolesRequest, *project.ListProjectRolesResponse)
		req *project.ListProjectRolesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *project.ListProjectRolesResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, instancePermissionV2.DefaultOrg.GetId(), integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{
					Filters: []*project.ProjectRoleSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, instancePermissionV2.DefaultOrg.GetId(), integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{
					Filters: []*project.ProjectRoleSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ListProjectRolesRequest{
					ProjectId: "notfound",
				},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id, missing permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
				ProjectRoles: []*project.ProjectRole{},
			},
		},
		{
			name: "list single id",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgID := instancePermissionV2.DefaultOrg.GetId()
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgID, integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					response.ProjectRoles[0] = addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectRoles: []*project.ProjectRole{
					{},
				},
			},
		},
		{
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgID := instancePermissionV2.DefaultOrg.GetId()
					projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, orgID, integration.ProjectName(), false, false)

					request.ProjectId = projectResp.GetId()
					response.ProjectRoles[2] = addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
					response.ProjectRoles[1] = addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
					response.ProjectRoles[0] = addProjectRole(iamOwnerCtx, instancePermissionV2, t, projectResp.GetId())
				},
				req: &project.ListProjectRolesRequest{},
			},
			want: &project.ListProjectRolesResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				ProjectRoles: []*project.ProjectRole{
					{}, {}, {},
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
				got, listErr := instancePermissionV2.Client.Projectv2Beta.ListProjectRoles(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr, "Error: "+listErr.Error())
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.ProjectRoles, len(tt.want.ProjectRoles)) {
					for i := range tt.want.ProjectRoles {
						assert.EqualExportedValues(ttt, tt.want.ProjectRoles[i], got.ProjectRoles[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func addProjectRole(ctx context.Context, instance *integration.Instance, t *testing.T, projectID string) *project.ProjectRole {
	name := integration.RoleKey()
	projectRoleResp := instance.AddProjectRole(ctx, t, projectID, name, name, name)

	return &project.ProjectRole{
		ProjectId:    projectID,
		CreationDate: projectRoleResp.GetCreationDate(),
		ChangeDate:   projectRoleResp.GetCreationDate(),
		Key:          name,
		DisplayName:  name,
		Group:        name,
	}
}
