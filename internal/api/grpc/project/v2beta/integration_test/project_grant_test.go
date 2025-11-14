//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/integration"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_CreateProjectGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *project.CreateProjectGrantRequest)
		req     *project.CreateProjectGrantRequest
		want
		wantErr bool
	}{
		{
			name:    "empty projectID",
			ctx:     iamOwnerCtx,
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "empty granted organization",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectGrantRequest{
				ProjectId:             "something",
				GrantedOrganizationId: "",
			},
			wantErr: true,
		},
		{
			name: "project not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				request.ProjectId = "something"
				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "org not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				request.GrantedOrganizationId = "something"
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "already existing, error",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "same organization, error",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				request.GrantedOrganizationId = orgResp.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "empty, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
			},
			req: &project.CreateProjectGrantRequest{},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "with roles, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				roles := []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
				request.RoleKeys = roles

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "with roles, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				roles := []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()
				for _, role := range roles {
					instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), role, role, "")
				}

				request.RoleKeys = roles

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
			},
			req: &project.CreateProjectGrantRequest{},
			want: want{
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}

			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.CreateProjectGrant(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateProjectGrantResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func TestServer_CreateProjectGrant_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userResp := instance.CreateMachineUser(iamOwnerCtx)
	patResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userResp.GetUserId())
	projectResp := createProject(iamOwnerCtx, instance, t, instance.DefaultOrg.GetId(), false, false)
	instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patResp.Token)

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *project.CreateProjectGrantRequest)
		req     *project.CreateProjectGrantRequest
		want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "no permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "project owner, other project",
			ctx:  projectOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "project owner, ok",
			ctx:  projectOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
			},
			req: &project.CreateProjectGrantRequest{},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "organization owner, other org",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req:     &project.CreateProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
			},
			req: &project.CreateProjectGrantRequest{},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "instance owner",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.CreateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				grantedOrgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

				request.ProjectId = projectResp.GetId()
				request.GrantedOrganizationId = grantedOrgResp.GetOrganizationId()
			},
			req: &project.CreateProjectGrantRequest{},
			want: want{
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}

			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.CreateProjectGrant(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateProjectGrantResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func assertCreateProjectGrantResponse(t *testing.T, creationDate, changeDate time.Time, expectedCreationDate bool, actualResp *project.CreateProjectGrantResponse) {
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

func TestServer_UpdateProjectGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.UpdateProjectGrantRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.UpdateProjectGrantRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "not existing",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				request.ProjectId = "notexisting"
				request.GrantedOrganizationId = "notexisting"
				return
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectGrantRequest{
					RoleKeys: []string{"notexisting"},
				},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectGrantRequest{},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "change roles, ok",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				roles := []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
				for _, role := range roles {
					instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), role, role, "")
				}

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId(), roles...)
				request.RoleKeys = roles
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectGrantRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "change roles, not existing",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				request.RoleKeys = []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectGrantRequest{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			if tt.prepare != nil {
				tt.prepare(tt.args.req)
			}

			got, err := instance.Client.Projectv2Beta.UpdateProjectGrant(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateProjectGrantResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_UpdateProjectGrant_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userResp := instance.CreateMachineUser(iamOwnerCtx)
	patResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userResp.GetUserId())
	projectID := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false).GetId()
	instance.CreateProjectGrant(iamOwnerCtx, t, projectID, orgResp.GetOrganizationId())
	instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectID, orgResp.GetOrganizationId(), userResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patResp.Token)

	type args struct {
		ctx context.Context
		req *project.UpdateProjectGrantRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.UpdateProjectGrantRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: CTX,
				req: &project.UpdateProjectGrantRequest{
					RoleKeys: []string{"nopermission"},
				},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.UpdateProjectGrantRequest{
					RoleKeys: []string{"nopermission"},
				},
			},
			wantErr: true,
		},
		{
			name: "project grant owner, no permission",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				roles := []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
				request.ProjectId = projectID
				request.GrantedOrganizationId = orgResp.GetOrganizationId()

				for _, role := range roles {
					instance.AddProjectRole(iamOwnerCtx, t, projectID, role, role, "")
				}

				request.RoleKeys = roles
			},
			args: args{
				ctx: projectGrantOwnerCtx,
				req: &project.UpdateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.UpdateProjectGrantRequest{
					RoleKeys: []string{"nopermission"},
				},
			},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				roles := []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
				for _, role := range roles {
					instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), role, role, "")
				}

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId(), roles...)
				request.RoleKeys = roles
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.UpdateProjectGrantRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner, ok",
			prepare: func(request *project.UpdateProjectGrantRequest) {
				roles := []string{integration.RoleKey(), integration.RoleKey(), integration.RoleKey()}
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()
				for _, role := range roles {
					instance.AddProjectRole(iamOwnerCtx, t, projectResp.GetId(), role, role, "")
				}

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId(), roles...)
				request.RoleKeys = roles
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectGrantRequest{},
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
			if tt.prepare != nil {
				tt.prepare(tt.args.req)
			}

			got, err := instance.Client.Projectv2Beta.UpdateProjectGrant(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateProjectGrantResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertUpdateProjectGrantResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.UpdateProjectGrantResponse) {
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

func TestServer_DeleteProjectGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time)
		req              *project.DeleteProjectGrantRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "empty project id",
			ctx:  iamOwnerCtx,
			req: &project.DeleteProjectGrantRequest{
				ProjectId: "",
			},
			wantErr: true,
		},
		{
			name: "empty grantedorg id",
			ctx:  iamOwnerCtx,
			req: &project.DeleteProjectGrantRequest{
				ProjectId:             "notempty",
				GrantedOrganizationId: "",
			},
			wantErr: true,
		},
		{
			name: "delete, not existing",
			ctx:  iamOwnerCtx,
			req: &project.DeleteProjectGrantRequest{
				ProjectId:             "notexisting",
				GrantedOrganizationId: "notexisting",
			},
			wantDeletionDate: false,
		},
		{
			name: "delete",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectGrantRequest{},
			wantDeletionDate: true,
		},
		{
			name: "delete deactivated",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectGrantRequest{},
			wantDeletionDate: true,
		},
		{
			name: "delete, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeleteProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return creationDate, time.Now().UTC()
			},
			req:              &project.DeleteProjectGrantRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := instance.Client.Projectv2Beta.DeleteProjectGrant(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteProjectGrantResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func TestServer_DeleteProjectGrant_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time)
		req              *project.DeleteProjectGrantRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return time.Time{}, time.Time{}
			},
			req:     &project.DeleteProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "no permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return time.Time{}, time.Time{}
			},
			req:     &project.DeleteProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return creationDate, time.Time{}
			},
			req:     &project.DeleteProjectGrantRequest{},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectGrantRequest{},
			wantDeletionDate: true,
		},
		{
			name: "organization owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectGrantRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectGrantRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := instance.Client.Projectv2Beta.DeleteProjectGrant(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteProjectGrantResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func assertDeleteProjectGrantResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *project.DeleteProjectGrantResponse) {
	if expectedDeletionDate {
		if !deletionDate.IsZero() {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, deletionDate)
		} else {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.DeletionDate)
	}
}

func TestServer_DeactivateProjectGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.DeactivateProjectGrantRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.DeactivateProjectGrantRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.DeactivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				request.ProjectId = "notexisting"
				request.GrantedOrganizationId = "notexisting"
				return
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectGrantRequest{},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "change, ok",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectGrantRequest{},
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
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.DeactivateProjectGrant(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertDeactivateProjectGrantResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_DeactivateProjectGrant_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.DeactivateProjectGrantRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.DeactivateProjectGrantRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: CTX,
				req: &project.DeactivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.DeactivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.DeactivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.DeactivateProjectGrantRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner, ok",
			prepare: func(request *project.DeactivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectGrantRequest{},
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
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.DeactivateProjectGrant(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertDeactivateProjectGrantResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertDeactivateProjectGrantResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.DeactivateProjectGrantResponse) {
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

func TestServer_ActivateProjectGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.ActivateProjectGrantRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.ActivateProjectGrantRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "not existing",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				request.ProjectId = "notexisting"
				request.GrantedOrganizationId = "notexisting"
				return
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectGrantRequest{},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "change, ok",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectGrantRequest{},
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
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.ActivateProjectGrant(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertActivateProjectGrantResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_ActivateProjectGrant_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.ActivateProjectGrantRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.ActivateProjectGrantRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: CTX,
				req: &project.ActivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "no permission",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.ActivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.ActivateProjectGrantRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.ActivateProjectGrantRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner, ok",
			prepare: func(request *project.ActivateProjectGrantRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.ProjectId = projectResp.GetId()

				grantedOrg := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				request.GrantedOrganizationId = grantedOrg.GetOrganizationId()

				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
				instance.DeactivateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), grantedOrg.GetOrganizationId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectGrantRequest{},
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
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.ActivateProjectGrant(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertActivateProjectGrantResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertActivateProjectGrantResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.ActivateProjectGrantResponse) {
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
