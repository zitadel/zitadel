//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_AddProjectRole(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.ProjectName(), integration.Email())
	alreadyExistingProject := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
	alreadyExistingProjectRoleName := integration.RoleDisplayName()
	instance.AddProjectRole(iamOwnerCtx, t, alreadyExistingProject.GetId(), alreadyExistingProjectRoleName, alreadyExistingProjectRoleName, "")

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(t *testing.T, request *project.AddProjectRoleRequest)
		req     *project.AddProjectRoleRequest
		want
		wantErr bool
	}{
		{
			name: "empty key",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     "",
				DisplayName: integration.ProjectName(),
			},
			wantErr: true,
		},
		{
			name: "empty displayname",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: "",
			},
			wantErr: true,
		},
		{
			name: "already existing, error",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				request.ProjectId = alreadyExistingProject.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     alreadyExistingProjectRoleName,
				DisplayName: alreadyExistingProjectRoleName,
			},
			wantErr: true,
		},
		{
			name: "empty, ok",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: integration.RoleDisplayName(),
			},
			want: want{
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(t, tt.req)
			}

			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.AddProjectRole(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertAddProjectRoleResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func TestServer_AddProjectRole_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	alreadyExistingProject := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
	alreadyExistingProjectRoleName := integration.RoleKey()
	instance.AddProjectRole(iamOwnerCtx, t, alreadyExistingProject.GetId(), alreadyExistingProjectRoleName, alreadyExistingProjectRoleName, "")

	type want struct {
		creationDate bool
	}

	type test struct {
		name    string
		ctx     context.Context
		prepare func(t *testing.T, request *project.AddProjectRoleRequest)
		req     *project.AddProjectRoleRequest
		want
		wantErr bool
	}
	tests := []*test{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: integration.RoleDisplayName(),
			},
			wantErr: true,
		},
		{
			name: "no permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: integration.RoleDisplayName(),
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: integration.RoleDisplayName(),
			},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: integration.RoleDisplayName(),
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.AddProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
			},
			req: &project.AddProjectRoleRequest{
				RoleKey:     integration.RoleKey(),
				DisplayName: integration.RoleDisplayName(),
			},
			want: want{
				creationDate: true,
			},
		},
		func() *test {
			out := test{
				name: "add project role as a added project admin, ok",
				req: &project.AddProjectRoleRequest{
					RoleKey:     integration.RoleKey(),
					DisplayName: integration.RoleDisplayName(),
				},
				want: want{
					creationDate: true,
				},
			}

			out.prepare = func(t *testing.T, request *project.AddProjectRoleRequest) {
				// create project
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, integration.ProjectName(), false, false)
				// create user
				userID := instance.CreateHumanUser(iamOwnerCtx).GetUserId()
				loginCTX := instance.WithAuthorization(iamOwnerCtx, integration.UserTypeLogin)
				instance.RegisterUserPasskey(iamOwnerCtx, userID)
				_, token, _, _ := instance.CreateVerifiedWebAuthNSession(t, loginCTX, userID)
				// assign user as project admin
				_, err := instance.Client.Mgmt.AddProjectMember(iamOwnerCtx, &management.AddProjectMemberRequest{
					ProjectId: projectResp.GetId(),
					UserId:    userID,
					Roles:     []string{"PROJECT_OWNER_GLOBAL"},
				})
				assert.NoError(t, err)

				// set context
				out.ctx = integration.WithAuthorizationToken(context.Background(), token)

				request.ProjectId = projectResp.GetId()
			}

			return &out
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(t, tt.req)
			}

			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.AddProjectRole(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertAddProjectRoleResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func assertAddProjectRoleResponse(t *testing.T, creationDate, changeDate time.Time, expectedCreationDate bool, actualResp *project.AddProjectRoleResponse) {
	if expectedCreationDate {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.CreationDate)
	}
}

func TestServer_UpdateProjectRole(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.UpdateProjectRoleRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(t *testing.T, request *project.UpdateProjectRoleRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				request.RoleKey = "notexisting"
				return
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				request.DisplayName = gu.Ptr(roleName)
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRoleRequest{},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "change display name, ok",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr(integration.RoleKey()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "change full, ok",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr(integration.RoleKey()),
					Group:       gu.Ptr(integration.RoleKey()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			tt.prepare(t, tt.args.req)

			got, err := instance.Client.Projectv2Beta.UpdateProjectRole(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateProjectRoleResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_UpdateProjectRole_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.UpdateProjectRoleRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	type test struct {
		name    string
		prepare func(t *testing.T, request *project.UpdateProjectRoleRequest)
		args    args
		want    want
		wantErr bool
	}
	tests := []*test{
		{
			name: "unauthenicated",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: CTX,
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr(integration.RoleKey()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner, ok",
			prepare: func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRoleRequest{
					DisplayName: gu.Ptr(integration.RoleKey()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		func() *test {
			out := test{
				name: "change project role as a added project admin, ok",
				args: args{
					req: &project.UpdateProjectRoleRequest{
						DisplayName: gu.Ptr(integration.RoleKey()),
						Group:       gu.Ptr(integration.RoleKey()),
					},
				},
				want: want{
					change:     true,
					changeDate: true,
				},
			}

			out.prepare = func(t *testing.T, request *project.UpdateProjectRoleRequest) {
				// create project
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, integration.ProjectName(), false, false)
				// create user
				userID := instance.CreateHumanUser(iamOwnerCtx).GetUserId()
				loginCTX := instance.WithAuthorization(iamOwnerCtx, integration.UserTypeLogin)
				instance.RegisterUserPasskey(iamOwnerCtx, userID)
				_, token, _, _ := instance.CreateVerifiedWebAuthNSession(t, loginCTX, userID)
				// assign user as project admin
				_, err := instance.Client.Mgmt.AddProjectMember(iamOwnerCtx, &management.AddProjectMemberRequest{
					ProjectId: projectResp.GetId(),
					UserId:    userID,
					Roles:     []string{"PROJECT_OWNER_GLOBAL"},
				})
				assert.NoError(t, err)

				// set context
				out.args.ctx = integration.WithAuthorizationToken(context.Background(), token)

				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
			}

			return &out
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			tt.prepare(t, tt.args.req)

			got, err := instance.Client.Projectv2Beta.UpdateProjectRole(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateProjectRoleResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertUpdateProjectRoleResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.UpdateProjectRoleResponse) {
	if expectedChangeDate {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetChangeDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetChangeDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.ChangeDate)
	}
}

func TestServer_DeleteProjectRole(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time)
		req              *project.RemoveProjectRoleRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "empty id",
			ctx:  iamOwnerCtx,
			req: &project.RemoveProjectRoleRequest{
				ProjectId: "",
				RoleKey:   "notexisting",
			},
			wantErr: true,
		},
		{
			name: "delete, not existing",
			ctx:  iamOwnerCtx,
			req: &project.RemoveProjectRoleRequest{
				ProjectId: "notexisting",
				RoleKey:   "notexisting",
			},
			wantDeletionDate: false,
		},
		{
			name: "delete",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			},
			req:              &project.RemoveProjectRoleRequest{},
			wantDeletionDate: true,
		},
		{
			name: "delete, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				instance.RemoveProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName)
				return creationDate, time.Now().UTC()
			},
			req:              &project.RemoveProjectRoleRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(t, tt.req)
			}
			got, err := instance.Client.Projectv2Beta.RemoveProjectRole(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertRemoveProjectRoleResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func TestServer_DeleteProjectRole_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type test struct {
		name             string
		ctx              context.Context
		prepare          func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time)
		req              *project.RemoveProjectRoleRequest
		wantDeletionDate bool
		wantErr          bool
	}
	tests := []*test{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			},
			req:     &project.RemoveProjectRoleRequest{},
			wantErr: true,
		},
		{
			name: "no permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			},
			req:     &project.RemoveProjectRoleRequest{},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			},
			req:     &project.RemoveProjectRoleRequest{},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			},
			req:              &project.RemoveProjectRoleRequest{},
			wantDeletionDate: true,
		},
		{
			name: "instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			},
			req:              &project.RemoveProjectRoleRequest{},
			wantDeletionDate: true,
		},
		func() *test {
			out := test{
				name:             "delete project role as a added project admin, ok",
				req:              &project.RemoveProjectRoleRequest{},
				wantDeletionDate: true,
			}

			out.prepare = func(t *testing.T, request *project.RemoveProjectRoleRequest) (time.Time, time.Time) {
				// create project
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, integration.ProjectName(), false, false)
				// create user
				userID := instance.CreateHumanUser(iamOwnerCtx).GetUserId()
				loginCTX := instance.WithAuthorization(iamOwnerCtx, integration.UserTypeLogin)
				instance.RegisterUserPasskey(iamOwnerCtx, userID)
				_, token, _, _ := instance.CreateVerifiedWebAuthNSession(t, loginCTX, userID)
				// assign user as project admin
				_, err := instance.Client.Mgmt.AddProjectMember(iamOwnerCtx, &management.AddProjectMemberRequest{
					ProjectId: projectResp.GetId(),
					UserId:    userID,
					Roles:     []string{"PROJECT_OWNER_GLOBAL"},
				})
				assert.NoError(t, err)

				// set context
				out.ctx = integration.WithAuthorizationToken(context.Background(), token)

				roleName := integration.RoleKey()
				instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), roleName, roleName, "")
				request.ProjectId = projectResp.GetId()
				request.RoleKey = roleName
				return creationDate, time.Time{}
			}

			return &out
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(t, tt.req)
			}
			got, err := instance.Client.Projectv2Beta.RemoveProjectRole(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertRemoveProjectRoleResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func assertRemoveProjectRoleResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *project.RemoveProjectRoleResponse) {
	if expectedDeletionDate {
		if !deletionDate.IsZero() {
			assert.WithinRange(t, actualResp.GetRemovalDate().AsTime(), creationDate, deletionDate)
		} else {
			assert.WithinRange(t, actualResp.GetRemovalDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.RemovalDate)
	}
}
