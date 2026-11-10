//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2"
)

func TestServer_ListAdministrators(t *testing.T) {
	// instance administrators with IAM_OWNER role for tests and 3 instance administrators to query
	iamOwnerCtx := instanceQuery.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	instanceAdministrator1 := createInstanceAdministrator(iamOwnerCtx, instanceQuery, t)
	instanceAdministrator2 := createInstanceAdministrator(iamOwnerCtx, instanceQuery, t)
	instanceAdministrator3 := createInstanceAdministrator(iamOwnerCtx, instanceQuery, t)

	// organization setup with 3 organization administrators to query
	organizationName := integration.OrganizationName()
	organizationResp := instanceQuery.CreateOrganization(iamOwnerCtx, organizationName, integration.Email())
	organizationAdmin1 := createOrganizationAdministrator(iamOwnerCtx, instanceQuery, organizationResp.GetOrganizationId(), organizationName, t)
	organizationAdmin2 := createOrganizationAdministrator(iamOwnerCtx, instanceQuery, organizationResp.GetOrganizationId(), organizationName, t)
	organizationAdmin3 := createOrganizationAdministrator(iamOwnerCtx, instanceQuery, organizationResp.GetOrganizationId(), organizationName, t)

	// Additionally, a service account as a project grant owner to test permissions
	userOrganizationResp := instanceQuery.CreateMachineUser(iamOwnerCtx)
	instanceQuery.CreateOrgMembership(t, iamOwnerCtx, organizationResp.GetOrganizationId(), userOrganizationResp.GetUserId())
	patOrganizationResp := instanceQuery.CreatePersonalAccessToken(iamOwnerCtx, userOrganizationResp.GetUserId())
	organizationOwnerCtx := integration.WithAuthorizationToken(CTX, patOrganizationResp.Token)

	// project setup on the same organization with 3 project administrators to query
	projectName := integration.ProjectName()
	projectResp := instanceQuery.CreateProject(iamOwnerCtx, t, organizationResp.GetOrganizationId(), projectName, false, false)
	projectAdministrator1 := createProjectAdministrator(iamOwnerCtx, instanceQuery, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName)
	projectAdministrator2 := createProjectAdministrator(iamOwnerCtx, instanceQuery, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName)
	projectAdministrator3 := createProjectAdministrator(iamOwnerCtx, instanceQuery, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName)

	// Additionally, a service account as a project owner to test permissions
	userProjectResp := instanceQuery.CreateMachineUser(iamOwnerCtx)
	instanceQuery.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instanceQuery.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	// project grant setup on the same project with 3 project grant administrators to query
	grantedOrganizationName := integration.OrganizationName()
	grantedOrganizationResp := instanceQuery.CreateOrganization(iamOwnerCtx, grantedOrganizationName, integration.Email())
	instanceQuery.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrganizationResp.GetOrganizationId())
	projectGrantAdministrator1 := createProjectGrantAdministrator(iamOwnerCtx, instanceQuery, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName, grantedOrganizationResp.GetOrganizationId())
	projectGrantAdministrator2 := createProjectGrantAdministrator(iamOwnerCtx, instanceQuery, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName, grantedOrganizationResp.GetOrganizationId())
	projectGrantAdministrator3 := createProjectGrantAdministrator(iamOwnerCtx, instanceQuery, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName, grantedOrganizationResp.GetOrganizationId())

	// Additionally, a service account as a project grant owner to test permissions
	userProjectGrantResp := instanceQuery.CreateMachineUser(iamOwnerCtx)
	instanceQuery.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), grantedOrganizationResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
	patProjectGrantResp := instanceQuery.CreatePersonalAccessToken(iamOwnerCtx, userProjectGrantResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectGrantResp.Token)

	type args struct {
		ctx context.Context
		dep func(*internal_permission.ListAdministratorsRequest, *internal_permission.ListAdministratorsResponse)
		req *internal_permission.ListAdministratorsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *internal_permission.ListAdministratorsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instanceQuery.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{},
			},
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: organizationOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{
						{
							Filter: &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
								InUserIdsFilter: &filter.InIDsFilter{Ids: []string{"notexisting"}},
							},
						},
					},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id, instance",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = instanceAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, instance",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), instanceAdministrator2.GetUser().GetId(), instanceAdministrator3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = instanceAdministrator3
					response.Administrators[1] = instanceAdministrator2
					response.Administrators[2] = instanceAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list single id, org",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{organizationAdmin1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = organizationAdmin1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, org",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{organizationAdmin1.GetUser().GetId(), organizationAdmin2.GetUser().GetId(), organizationAdmin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = organizationAdmin3
					response.Administrators[1] = organizationAdmin2
					response.Administrators[2] = organizationAdmin1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list single id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectAdministrator1.GetUser().GetId(), projectAdministrator2.GetUser().GetId(), projectAdministrator3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectAdministrator3
					response.Administrators[1] = projectAdministrator2
					response.Administrators[2] = projectAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list single id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectGrantAdministrator1.GetUser().GetId(), projectGrantAdministrator2.GetUser().GetId(), projectGrantAdministrator3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator3
					response.Administrators[1] = projectGrantAdministrator2
					response.Administrators[2] = projectGrantAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, instance owner",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
					response.Administrators[1] = projectAdministrator1
					response.Administrators[2] = organizationAdmin1
					response.Administrators[3] = instanceAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, org owner",
			args: args{
				ctx: organizationOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
					response.Administrators[1] = projectAdministrator1
					response.Administrators[2] = organizationAdmin1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, project owner",
			args: args{
				ctx: projectOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
					response.Administrators[1] = projectAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {},
				},
			},
		},
		{
			name: "list multiple id, project grant owner",
			args: args{
				ctx: projectGrantOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list, resource filter, instance admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_Instance{Instance: true},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2, Offset: 2}, // skip the first two, which were created
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5, // 3 test admins + 2 during the integration test setup
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{instanceAdministrator1, instanceAdministrator2},
			},
		},
		{
			name: "list, resource filter, organization admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_OrganizationId{OrganizationId: organizationResp.GetOrganizationId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2, Offset: 1}, // skip the first as this is the instance admin creating the organization
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5, // 3 test admins + 1 for testing the permissions + 1 for the instance admin creating the org
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{organizationAdmin1, organizationAdmin2},
			},
		},
		{
			name: "list, resource filter, project admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_ProjectId{ProjectId: projectResp.GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4, // 3 test admins + 1 for testing the permissions
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{projectAdministrator1, projectAdministrator2},
			},
		},
		{
			name: "list, resource filter, project grant admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_ProjectGrant_{
								ProjectGrant: &internal_permission.ResourceFilter_ProjectGrant{
									ProjectId:      projectResp.GetId(),
									OrganizationId: grantedOrganizationResp.GetOrganizationId(),
								},
							},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4, // 3 test admins + 1 for testing the permissions
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{projectGrantAdministrator1, projectGrantAdministrator2},
			},
		},
		{
			name: "list, role filter",
			args: args{
				ctx: iamOwnerCtx,
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{
						{
							Filter: &internal_permission.AdministratorSearchFilter_Role{
								Role: &internal_permission.RoleFilter{
									RoleKey: "PROJECT_OWNER",
								},
							},
						},
					},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4, // 3 project admins with PROJECT_OWNER role and the project creator
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{projectAdministrator1, projectAdministrator2},
			},
		},
		{
			name: "list, and with not filter",
			args: args{
				ctx: iamOwnerCtx,
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{
						{
							Filter: &internal_permission.AdministratorSearchFilter_And{
								And: &internal_permission.AndFilter{
									Queries: []*internal_permission.AdministratorSearchFilter{
										{
											Filter: &internal_permission.AdministratorSearchFilter_Role{
												Role: &internal_permission.RoleFilter{
													RoleKey: "PROJECT_OWNER",
												},
											},
										},
										{
											Filter: &internal_permission.AdministratorSearchFilter_Not{
												Not: &internal_permission.NotFilter{
													Query: &internal_permission.AdministratorSearchFilter{
														Filter: &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
															InUserIdsFilter: &filter.InIDsFilter{
																Ids: []string{projectAdministrator1.GetUser().GetId(), projectAdministrator2.GetUser().GetId()},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 1},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  2, // 1 project admin and the project creator
					AppliedLimit: 1,
				},
				Administrators: []*internal_permission.Administrator{projectAdministrator3},
			},
		},
		{
			name: "list, or filter",
			args: args{
				ctx: iamOwnerCtx,
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{
						{
							Filter: &internal_permission.AdministratorSearchFilter_Or{
								Or: &internal_permission.OrFilter{
									Queries: []*internal_permission.AdministratorSearchFilter{
										{
											Filter: &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
												InUserIdsFilter: &filter.InIDsFilter{
													Ids: []string{projectGrantAdministrator1.GetUser().GetId()},
												},
											},
										},
										{
											Filter: &internal_permission.AdministratorSearchFilter_Role{
												Role: &internal_permission.RoleFilter{
													RoleKey: "PROJECT_OWNER",
												},
											},
										},
									},
								},
							},
						},
					},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5, // 3 project grant admins, the project creator and the project grant admin
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{projectAdministrator1, projectAdministrator2},
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
				got, listErr := instanceQuery.Client.InternalPermissionV2.ListAdministrators(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Administrators, len(tt.want.Administrators)) {
					for i := range tt.want.Administrators {
						// need to set the project grant ID as it is generated
						if grant := got.Administrators[i].GetProjectGrant(); grant != nil {
							tt.want.Administrators[i].GetProjectGrant().Id = grant.Id
						}
						assert.EqualExportedValues(ttt, tt.want.Administrators[i], got.Administrators[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	t.Helper()
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func createInstanceAdministrator(ctx context.Context, instance *integration.Instance, t *testing.T) *internal_permission.Administrator {
	email := integration.Email()
	userResp := instance.CreateUserTypeHuman(ctx, email)
	memberResp := instance.CreateInstanceMembership(t, ctx, userResp.GetId())
	return &internal_permission.Administrator{
		CreationDate: memberResp.GetCreationDate(),
		ChangeDate:   memberResp.GetCreationDate(),
		User: &internal_permission.User{
			Id:                 userResp.GetId(),
			PreferredLoginName: email,
			DisplayName:        "Mickey Mouse",
			OrganizationId:     instance.DefaultOrg.GetId(),
		},
		Resource: &internal_permission.Administrator_Instance{
			Instance: true,
		},
		Roles: []string{"IAM_OWNER"},
	}
}

func createOrganizationAdministrator(ctx context.Context, instance *integration.Instance, organizationID, organizationName string, t *testing.T) *internal_permission.Administrator {
	email := integration.Email()

	userResp := instance.CreateUserTypeHuman(ctx, email)
	memberResp := instance.CreateOrgMembership(t, ctx, organizationID, userResp.GetId(), "ORG_OWNER")
	return &internal_permission.Administrator{
		CreationDate: memberResp.GetCreationDate(),
		ChangeDate:   memberResp.GetCreationDate(),
		User: &internal_permission.User{
			Id:                 userResp.GetId(),
			PreferredLoginName: email,
			DisplayName:        "Mickey Mouse",
			OrganizationId:     instance.DefaultOrg.GetId(),
		},
		Resource: &internal_permission.Administrator_Organization{
			Organization: &internal_permission.Organization{
				Id:   organizationID,
				Name: organizationName,
			},
		},
		Roles: []string{"ORG_OWNER"},
	}
}

func createProjectAdministrator(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, projectID, projectName string) *internal_permission.Administrator {
	email := integration.Email()
	userResp := instance.CreateUserTypeHuman(ctx, email)
	memberResp := instance.CreateProjectMembership(t, ctx, projectID, userResp.GetId())
	return &internal_permission.Administrator{
		CreationDate: memberResp.GetCreationDate(),
		ChangeDate:   memberResp.GetCreationDate(),
		User: &internal_permission.User{
			Id:                 userResp.GetId(),
			PreferredLoginName: email,
			DisplayName:        "Mickey Mouse",
			OrganizationId:     instance.DefaultOrg.GetId(),
		},
		Resource: &internal_permission.Administrator_Project{
			Project: &internal_permission.Project{
				Id:             projectID,
				Name:           projectName,
				OrganizationId: orgID,
			},
		},
		Roles: []string{"PROJECT_OWNER"},
	}
}

func createProjectGrantAdministrator(ctx context.Context, instance *integration.Instance, t *testing.T, orgID, projectID, projectName, grantedOrgID string) *internal_permission.Administrator {
	email := integration.Email()
	userResp := instance.CreateUserTypeHuman(ctx, email)
	memberResp := instance.CreateProjectGrantMembership(t, ctx, projectID, grantedOrgID, userResp.GetId())
	return &internal_permission.Administrator{
		CreationDate: memberResp.GetCreationDate(),
		ChangeDate:   memberResp.GetCreationDate(),
		User: &internal_permission.User{
			Id:                 userResp.GetId(),
			PreferredLoginName: email,
			DisplayName:        "Mickey Mouse",
			OrganizationId:     instance.DefaultOrg.GetId(),
		},
		Resource: &internal_permission.Administrator_ProjectGrant{
			ProjectGrant: &internal_permission.ProjectGrant{
				// left empty as generated
				Id:                    "",
				ProjectId:             projectID,
				ProjectName:           projectName,
				OrganizationId:        orgID,
				GrantedOrganizationId: grantedOrgID,
			},
		},
		Roles: []string{"PROJECT_GRANT_OWNER"},
	}
}

func TestServer_ListAdministrators_PermissionV2(t *testing.T) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// ensureFeaturePermissionV2Enabled(t, instancePermissionV2)

	// instance administrators with IAM_OWNER role for tests and 3 instance administrators to query
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	instanceAdministrator1 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
	instanceAdministrator2 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
	instanceAdministrator3 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)

	// organization setup with 3 organization administrators to query
	organizationName := integration.OrganizationName()
	organizationResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, organizationName, integration.Email())
	organizationAdmin1 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, organizationResp.GetOrganizationId(), organizationName, t)
	organizationAdmin2 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, organizationResp.GetOrganizationId(), organizationName, t)
	organizationAdmin3 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, organizationResp.GetOrganizationId(), organizationName, t)

	// Additionally, a service account as a project grant owner to test permissions
	userOrganizationResp := instancePermissionV2.CreateMachineUser(iamOwnerCtx)
	instancePermissionV2.CreateOrgMembership(t, iamOwnerCtx, organizationResp.GetOrganizationId(), userOrganizationResp.GetUserId())
	patOrganizationResp := instancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, userOrganizationResp.GetUserId())
	organizationOwnerCtx := integration.WithAuthorizationToken(CTX, patOrganizationResp.Token)

	// project setup on the same organization with 3 project administrators to query
	projectName := integration.ProjectName()
	projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, organizationResp.GetOrganizationId(), projectName, false, false)
	projectAdministrator1 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName)
	projectAdministrator2 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName)
	projectAdministrator3 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName)

	// Additionally, a service account as a project owner to test permissions
	userProjectResp := instancePermissionV2.CreateMachineUser(iamOwnerCtx)
	instancePermissionV2.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	// project grant setup on the same project with 3 project grant administrators to query
	grantedOrganizationName := integration.OrganizationName()
	grantedOrganizationResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, grantedOrganizationName, integration.Email())
	instancePermissionV2.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrganizationResp.GetOrganizationId())
	projectGrantAdministrator1 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName, grantedOrganizationResp.GetOrganizationId())
	projectGrantAdministrator2 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName, grantedOrganizationResp.GetOrganizationId())
	projectGrantAdministrator3 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, organizationResp.GetOrganizationId(), projectResp.GetId(), projectName, grantedOrganizationResp.GetOrganizationId())

	// Additionally, a service account as a project grant owner to test permissions
	userProjectGrantResp := instancePermissionV2.CreateMachineUser(iamOwnerCtx)
	instancePermissionV2.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), grantedOrganizationResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
	patProjectGrantResp := instancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, userProjectGrantResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectGrantResp.Token)
	type args struct {
		ctx context.Context
		dep func(*internal_permission.ListAdministratorsRequest, *internal_permission.ListAdministratorsResponse)
		req *internal_permission.ListAdministratorsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *internal_permission.ListAdministratorsResponse
		wantErr bool
	}{
		{
			name: "list by id, unauthenticated",
			args: args{
				ctx: CTX,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			wantErr: true,
		},
		{
			name: "list by id, no permission",
			args: args{
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{},
			},
		},
		{
			name: "list by id, missing permission",
			args: args{
				ctx: organizationOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{},
			},
		},
		{
			name: "list, not found",
			args: args{
				ctx: iamOwnerCtx,
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{
						{
							Filter: &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
								InUserIdsFilter: &filter.InIDsFilter{Ids: []string{"notexisting"}},
							},
						},
					},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list single id, instance",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = instanceAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, instance",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), instanceAdministrator2.GetUser().GetId(), instanceAdministrator3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = instanceAdministrator3
					response.Administrators[1] = instanceAdministrator2
					response.Administrators[2] = instanceAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list single id, org",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{organizationAdmin1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = organizationAdmin1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, org",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{organizationAdmin1.GetUser().GetId(), organizationAdmin2.GetUser().GetId(), organizationAdmin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = organizationAdmin3
					response.Administrators[1] = organizationAdmin2
					response.Administrators[2] = organizationAdmin1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list single id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, project",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectAdministrator1.GetUser().GetId(), projectAdministrator2.GetUser().GetId(), projectAdministrator3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectAdministrator3
					response.Administrators[1] = projectAdministrator2
					response.Administrators[2] = projectAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list single id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{},
				},
			},
		},
		{
			name: "list multiple id, project grant",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{projectGrantAdministrator1.GetUser().GetId(), projectGrantAdministrator2.GetUser().GetId(), projectGrantAdministrator3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator3
					response.Administrators[1] = projectGrantAdministrator2
					response.Administrators[2] = projectGrantAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, instance owner",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
					response.Administrators[1] = projectAdministrator1
					response.Administrators[2] = organizationAdmin1
					response.Administrators[3] = instanceAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, org owner",
			args: args{
				ctx: organizationOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
					response.Administrators[1] = projectAdministrator1
					response.Administrators[2] = organizationAdmin1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {}, {},
				},
			},
		},
		{
			name: "list multiple id, project owner",
			args: args{
				ctx: projectOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
					response.Administrators[1] = projectAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{
					{}, {},
				},
			},
		},
		{
			name: "list multiple id, project grant owner",
			args: args{
				ctx: projectGrantOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{instanceAdministrator1.GetUser().GetId(), organizationAdmin1.GetUser().GetId(), projectAdministrator1.GetUser().GetId(), projectGrantAdministrator1.GetUser().GetId()},
						},
					}
					response.Administrators[0] = projectGrantAdministrator1
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters: []*internal_permission.AdministratorSearchFilter{{}},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
				Administrators: []*internal_permission.Administrator{{}},
			},
		},
		{
			name: "list, resource filter, instance admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_Instance{Instance: true},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2, Offset: 2}, // skip the first two, which were created
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5, // 3 test admins + 2 during the integration test setup
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{instanceAdministrator1, instanceAdministrator2},
			},
		},
		{
			name: "list, resource filter, organization admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_OrganizationId{OrganizationId: organizationResp.GetOrganizationId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2, Offset: 1}, // skip the first as this is the instance admin creating the organization
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  5, // 3 test admins + 1 for testing the permissions + 1 for the instance admin creating the org
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{organizationAdmin1, organizationAdmin2},
			},
		},
		{
			name: "list, resource filter, project admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_ProjectId{ProjectId: projectResp.GetId()},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4, // 3 test admins + 1 for testing the permissions
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{projectAdministrator1, projectAdministrator2},
			},
		},
		{
			name: "list, resource filter, project grant admin",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_Resource{
						Resource: &internal_permission.ResourceFilter{
							Resource: &internal_permission.ResourceFilter_ProjectGrant_{
								ProjectGrant: &internal_permission.ResourceFilter_ProjectGrant{
									ProjectId:      projectResp.GetId(),
									OrganizationId: grantedOrganizationResp.GetOrganizationId(),
								},
							},
						},
					}
				},
				req: &internal_permission.ListAdministratorsRequest{
					Filters:    []*internal_permission.AdministratorSearchFilter{{}},
					Pagination: &filter.PaginationRequest{Asc: true, Limit: 2},
				},
			},
			want: &internal_permission.ListAdministratorsResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  4, // 3 test admins + 1 for testing the permissions
					AppliedLimit: 2,
				},
				Administrators: []*internal_permission.Administrator{projectGrantAdministrator1, projectGrantAdministrator2},
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
				got, listErr := instancePermissionV2.Client.InternalPermissionV2.ListAdministrators(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Administrators, len(tt.want.Administrators)) {
					for i := range tt.want.Administrators {
						// set as project grant id is generated
						if grant := got.Administrators[i].GetProjectGrant(); grant != nil {
							tt.want.Administrators[i].GetProjectGrant().Id = grant.Id
						}
						assert.EqualExportedValues(ttt, tt.want.Administrators[i], got.Administrators[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}
