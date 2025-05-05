//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_GetProject(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

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
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeNoPermission),
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := createProject(iamOwnerCtx, t, orgID, false, false)

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
					resp := createProject(iamOwnerCtx, t, orgID, false, false)

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
					resp := createProject(iamOwnerCtx, t, orgID, true, true)

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
					assert.Error(ttt, err, "Error: "+err.Error())
					return
				}
				assert.NoError(ttt, err)
				assert.EqualExportedValues(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected project result")
		})
	}
}

func TestServer_ListProjects(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
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
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					orgID := instance.DefaultOrg.GetId()
					resp := createProject(iamOwnerCtx, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{resp.GetId()},
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
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ListProjectsRequest{
					Filters: []*project.ProjectSearchFilter{
						{Filter: &project.ProjectSearchFilter_InProjectIdsFilter{
							InProjectIdsFilter: &project.InProjectIDsFilter{
								ProjectIds: []string{"notfound"},
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
					response.Projects[0] = createProject(iamOwnerCtx, t, orgID, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{response.Projects[0].GetId()},
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
					response.Projects[0] = createProject(iamOwnerCtx, t, orgID, false, false)
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
					orgID := instance.DefaultOrg.GetId()
					response.Projects[2] = createProject(iamOwnerCtx, t, orgID, false, false)
					response.Projects[1] = createProject(iamOwnerCtx, t, orgID, true, false)
					response.Projects[0] = createProject(iamOwnerCtx, t, orgID, false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{response.Projects[0].GetId(), response.Projects[1].GetId(), response.Projects[2].GetId()},
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
					require.Error(ttt, listErr, "Error: "+listErr.Error())
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

func createProject(ctx context.Context, t *testing.T, orgID string, projectRoleCheck, hasProjectCheck bool) *project.Project {
	name := gofakeit.Name()
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

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func TestServer_ListProjectGrants(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
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
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.Name(), false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{projectResp.GetId()},
						},
					}
					grantedOrg := instance.CreateOrganization(iamOwnerCtx, gofakeit.AppName(), gofakeit.Email())
					request.Filters[1].Filter = &project.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter{
						ProjectGrantResourceOwnerFilter: &project.ProjectGrantResourceOwnerFilter{
							ProjectGrantResourceOwner: grantedOrg.GetOrganizationId(),
						},
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
							InProjectIdsFilter: &project.InProjectIDsFilter{
								ProjectIds: []string{"notfound"},
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
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{projectResp.GetId()},
						},
					}

					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, t, orgID, projectResp.GetId(), name)
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
			name: "list multiple id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectGrantsRequest, response *project.ListProjectGrantsResponse) {
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectGrantSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{projectResp.GetId()},
						},
					}

					response.ProjectGrants[2] = createProjectGrant(iamOwnerCtx, t, orgID, projectResp.GetId(), name)
					response.ProjectGrants[1] = createProjectGrant(iamOwnerCtx, t, orgID, projectResp.GetId(), name)
					response.ProjectGrants[0] = createProjectGrant(iamOwnerCtx, t, orgID, projectResp.GetId(), name)
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

func createProjectGrant(ctx context.Context, t *testing.T, orgID, projectID, projectName string) *project.ProjectGrant {
	grantedOrgName := gofakeit.AppName()
	grantedOrg := instance.CreateOrganization(ctx, grantedOrgName, gofakeit.Email())
	projectGrantResp := instance.CreateProjectGrant(ctx, t, projectID, grantedOrg.GetOrganizationId())

	return &project.ProjectGrant{
		OrganizationId:          orgID,
		CreationDate:            projectGrantResp.GetCreationDate(),
		ChangeDate:              projectGrantResp.GetCreationDate(),
		GrantedOrganizationId:   grantedOrg.GetOrganizationId(),
		GrantedOrganizationName: grantedOrgName,
		ProjectId:               projectID,
		ProjectName:             projectName,
		State:                   1,
	}
}

func TestServer_ListProjectRoles(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
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
			name: "list by id, missing permission",
			args: args{
				ctx: instance.WithAuthorization(context.Background(), integration.UserTypeNoPermission),
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.Name(), false, false)

					request.ProjectId = projectResp.GetId()
					addProjectRole(iamOwnerCtx, t, projectResp.GetId())
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
			name: "list single id",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectRolesRequest, response *project.ListProjectRolesResponse) {
					orgID := instance.DefaultOrg.GetId()
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, gofakeit.Name(), false, false)

					request.ProjectId = projectResp.GetId()
					response.ProjectRoles[0] = addProjectRole(iamOwnerCtx, t, projectResp.GetId())
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
					projectResp := instance.CreateProject(iamOwnerCtx, t, orgID, gofakeit.Name(), false, false)

					request.ProjectId = projectResp.GetId()
					response.ProjectRoles[2] = addProjectRole(iamOwnerCtx, t, projectResp.GetId())
					response.ProjectRoles[1] = addProjectRole(iamOwnerCtx, t, projectResp.GetId())
					response.ProjectRoles[0] = addProjectRole(iamOwnerCtx, t, projectResp.GetId())
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

func addProjectRole(ctx context.Context, t *testing.T, projectID string) *project.ProjectRole {
	name := gofakeit.Animal()
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
