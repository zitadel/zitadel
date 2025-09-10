//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

func TestServer_ListAdministrators(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	projectName := integration.ProjectName()
	projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), projectName, false, false)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

	userProjectResp := instance.CreateMachineUser(iamOwnerCtx)
	instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	userProjectGrantResp := instance.CreateMachineUser(iamOwnerCtx)
	instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
	patProjectGrantResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectGrantResp.GetUserId())
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
					admin := createInstanceAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
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
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin := createInstanceAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
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
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin := createInstanceAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
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
					admin := createInstanceAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
			name: "list single id, instance",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin := createInstanceAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					admin2 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					admin3 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					admin3 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin2 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin3 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					admin2 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					admin3 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
					response.Administrators[1] = admin3
					response.Administrators[2] = admin2
					response.Administrators[3] = admin1
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
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin1 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
					response.Administrators[1] = admin3
					response.Administrators[2] = admin2
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
					response.Administrators[1] = admin3
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instance, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instance, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instance.Client.InternalPermissionv2Beta.ListAdministrators(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Administrators, len(tt.want.Administrators)) {
					for i := range tt.want.Administrators {
						assert.EqualExportedValues(ttt, tt.want.Administrators[i], got.Administrators[i])
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

func createOrganizationAdministrator(ctx context.Context, instance *integration.Instance, t *testing.T) *internal_permission.Administrator {
	email := integration.Email()
	userResp := instance.CreateUserTypeHuman(ctx, email)
	memberResp := instance.CreateOrgMembership(t, ctx, instance.DefaultOrg.Id, userResp.GetId())
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
				Id:   instance.DefaultOrg.GetId(),
				Name: instance.DefaultOrg.GetName(),
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
				Id:                    grantedOrgID,
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
	iamOwnerCtx := instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	projectName := integration.ProjectName()
	projectResp := instancePermissionV2.CreateProject(iamOwnerCtx, t, instancePermissionV2.DefaultOrg.GetId(), projectName, false, false)
	orgResp := instancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	instancePermissionV2.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

	userProjectResp := instancePermissionV2.CreateMachineUser(iamOwnerCtx)
	instancePermissionV2.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instancePermissionV2.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	userProjectGrantResp := instancePermissionV2.CreateMachineUser(iamOwnerCtx)
	instancePermissionV2.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
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
					admin := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
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
					admin := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
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
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
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
					admin := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
			name: "list single id, instance",
			args: args{
				ctx: iamOwnerCtx,
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin2 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin3 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin3 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin2 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin3 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin
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
					admin1 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					admin2 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					admin3 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin3
					response.Administrators[1] = admin2
					response.Administrators[2] = admin1
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
					response.Administrators[1] = admin3
					response.Administrators[2] = admin2
					response.Administrators[3] = admin1
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
				ctx: instancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				dep: func(request *internal_permission.ListAdministratorsRequest, response *internal_permission.ListAdministratorsResponse) {
					admin1 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
					response.Administrators[1] = admin3
					response.Administrators[2] = admin2
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
					response.Administrators[1] = admin3
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
					admin1 := createInstanceAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin2 := createOrganizationAdministrator(iamOwnerCtx, instancePermissionV2, t)
					admin3 := createProjectAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName)
					admin4 := createProjectGrantAdministrator(iamOwnerCtx, instancePermissionV2, t, instancePermissionV2.DefaultOrg.GetId(), projectResp.GetId(), projectName, orgResp.GetOrganizationId())
					request.Filters[0].Filter = &internal_permission.AdministratorSearchFilter_InUserIdsFilter{
						InUserIdsFilter: &filter.InIDsFilter{
							Ids: []string{admin1.GetUser().GetId(), admin2.GetUser().GetId(), admin3.GetUser().GetId(), admin4.GetUser().GetId()},
						},
					}
					response.Administrators[0] = admin4
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := instancePermissionV2.Client.InternalPermissionv2Beta.ListAdministrators(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)

				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Administrators, len(tt.want.Administrators)) {
					for i := range tt.want.Administrators {
						assert.EqualExportedValues(ttt, tt.want.Administrators[i], got.Administrators[i])
					}
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}
