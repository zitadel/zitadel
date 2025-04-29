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
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)

					request.Id = resp.GetId()

					response.Project.Id = resp.GetId()
					response.Project.OrganizationId = orgID
					response.Project.Name = name
					response.Project.CreationDate = resp.GetCreationDate()
					response.Project.ChangeDate = resp.GetCreationDate()
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
			name: "get, full, ok",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.GetProjectRequest, response *project.GetProjectResponse) {
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, name, true, true)

					request.Id = resp.GetId()

					response.Project.Id = resp.GetId()
					response.Project.OrganizationId = orgID
					response.Project.Name = name
					response.Project.CreationDate = resp.GetCreationDate()
					response.Project.ChangeDate = resp.GetCreationDate()
				},
				req: &project.GetProjectRequest{},
			},
			want: &project.GetProjectResponse{
				Project: &project.Project{
					State:                  1,
					ProjectRoleAssertion:   false,
					ProjectAccessRequired:  true,
					AuthorizationRequired:  true,
					PrivateLabelingSetting: project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
				},
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
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
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
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{resp.GetId()},
						},
					}

					response.Projects[0].Id = resp.GetId()
					response.Projects[0].Name = name
					response.Projects[0].OrganizationId = orgID
					response.Projects[0].CreationDate = resp.GetCreationDate()
					response.Projects[0].ChangeDate = resp.GetCreationDate()
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
			name: "list single name",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *project.ListProjectsRequest, response *project.ListProjectsResponse) {
					name := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					resp := instance.CreateProject(iamOwnerCtx, t, orgID, name, false, false)
					request.Filters[0].Filter = &project.ProjectSearchFilter_ProjectNameFilter{
						ProjectNameFilter: &project.ProjectNameFilter{
							ProjectName: name,
						},
					}

					response.Projects[0].Id = resp.GetId()
					response.Projects[0].Name = name
					response.Projects[0].OrganizationId = orgID
					response.Projects[0].CreationDate = resp.GetCreationDate()
					response.Projects[0].ChangeDate = resp.GetCreationDate()
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
					name1 := gofakeit.Name()
					name2 := gofakeit.Name()
					name3 := gofakeit.Name()
					orgID := instance.DefaultOrg.GetId()
					resp1 := instance.CreateProject(iamOwnerCtx, t, orgID, name1, false, false)
					resp2 := instance.CreateProject(iamOwnerCtx, t, orgID, name2, true, false)
					resp3 := instance.CreateProject(iamOwnerCtx, t, orgID, name3, false, true)
					request.Filters[0].Filter = &project.ProjectSearchFilter_InProjectIdsFilter{
						InProjectIdsFilter: &project.InProjectIDsFilter{
							ProjectIds: []string{resp1.GetId(), resp2.GetId(), resp3.GetId()},
						},
					}

					response.Projects[2].Id = resp1.GetId()
					response.Projects[2].Name = name1
					response.Projects[2].OrganizationId = orgID
					response.Projects[2].CreationDate = resp1.GetCreationDate()
					response.Projects[2].ChangeDate = resp1.GetCreationDate()

					response.Projects[1].Id = resp2.GetId()
					response.Projects[1].Name = name2
					response.Projects[1].OrganizationId = orgID
					response.Projects[1].CreationDate = resp2.GetCreationDate()
					response.Projects[1].ChangeDate = resp2.GetCreationDate()

					response.Projects[0].Id = resp3.GetId()
					response.Projects[0].Name = name3
					response.Projects[0].OrganizationId = orgID
					response.Projects[0].CreationDate = resp3.GetCreationDate()
					response.Projects[0].ChangeDate = resp3.GetCreationDate()
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
					{
						State:                  1,
						ProjectRoleAssertion:   false,
						ProjectAccessRequired:  false,
						AuthorizationRequired:  true,
						PrivateLabelingSetting: project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
					},
					{
						State:                  1,
						ProjectRoleAssertion:   false,
						ProjectAccessRequired:  true,
						AuthorizationRequired:  false,
						PrivateLabelingSetting: project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_UNSPECIFIED,
					},
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
					grantedOrgName := gofakeit.AppName()
					grantedOrg := instance.CreateOrganization(iamOwnerCtx, grantedOrgName, gofakeit.Email())
					request.Filters[1].Filter = &project.ProjectGrantSearchFilter_ProjectGrantResourceOwnerFilter{
						ProjectGrantResourceOwnerFilter: &project.ProjectGrantResourceOwnerFilter{
							ProjectGrantResourceOwner: grantedOrg.GetOrganizationId(),
						},
					}
					projectGrantResp := instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())

					response.ProjectGrants[0].OrganizationId = orgID
					response.ProjectGrants[0].CreationDate = projectGrantResp.GetCreationDate()
					response.ProjectGrants[0].ChangeDate = projectGrantResp.GetCreationDate()
					response.ProjectGrants[0].GrantedOrganizationId = grantedOrg.GetOrganizationId()
					response.ProjectGrants[0].GrantedOrganizationName = grantedOrgName
					response.ProjectGrants[0].ProjectId = projectResp.GetId()
					response.ProjectGrants[0].ProjectName = name
				},
				req: &project.ListProjectGrantsRequest{
					Filters: []*project.ProjectGrantSearchFilter{{}, {}},
				},
			},
			want: &project.ListProjectGrantsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				ProjectGrants: []*project.ProjectGrant{
					{
						State: 1,
					},
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
					grantedOrgName1 := gofakeit.AppName()
					grantedOrg1 := instance.CreateOrganization(iamOwnerCtx, grantedOrgName1, gofakeit.Email())
					projectGrantResp1 := instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg1.GetOrganizationId())

					response.ProjectGrants[0].OrganizationId = orgID
					response.ProjectGrants[0].CreationDate = projectGrantResp1.GetCreationDate()
					response.ProjectGrants[0].ChangeDate = projectGrantResp1.GetCreationDate()
					response.ProjectGrants[0].GrantedOrganizationId = grantedOrg1.GetOrganizationId()
					response.ProjectGrants[0].GrantedOrganizationName = grantedOrgName1
					response.ProjectGrants[0].ProjectId = projectResp.GetId()
					response.ProjectGrants[0].ProjectName = name

					grantedOrgName2 := gofakeit.AppName()
					grantedOrg2 := instance.CreateOrganization(iamOwnerCtx, grantedOrgName2, gofakeit.Email())
					projectGrantResp2 := instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg2.GetOrganizationId())

					response.ProjectGrants[1].OrganizationId = orgID
					response.ProjectGrants[1].CreationDate = projectGrantResp2.GetCreationDate()
					response.ProjectGrants[1].ChangeDate = projectGrantResp2.GetCreationDate()
					response.ProjectGrants[1].GrantedOrganizationId = grantedOrg2.GetOrganizationId()
					response.ProjectGrants[1].GrantedOrganizationName = grantedOrgName2
					response.ProjectGrants[1].ProjectId = projectResp.GetId()
					response.ProjectGrants[1].ProjectName = name

					grantedOrgName3 := gofakeit.AppName()
					grantedOrg3 := instance.CreateOrganization(iamOwnerCtx, grantedOrgName3, gofakeit.Email())
					projectGrantResp3 := instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg3.GetOrganizationId())

					response.ProjectGrants[2].OrganizationId = orgID
					response.ProjectGrants[2].CreationDate = projectGrantResp3.GetCreationDate()
					response.ProjectGrants[2].ChangeDate = projectGrantResp3.GetCreationDate()
					response.ProjectGrants[2].GrantedOrganizationId = grantedOrg3.GetOrganizationId()
					response.ProjectGrants[2].GrantedOrganizationName = grantedOrgName3
					response.ProjectGrants[2].ProjectId = projectResp.GetId()
					response.ProjectGrants[2].ProjectName = name
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
					{State: 1}, {State: 1}, {State: 1},
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
