//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	internal_permission_v2beta "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_CreateProject(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	alreadyExistingProjectName := integration.ProjectName()
	instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), alreadyExistingProjectName, false, false)

	type want struct {
		id           bool
		creationDate bool
	}
	tests := []struct {
		name string
		ctx  context.Context
		req  *project.CreateProjectRequest
		want
		wantErr bool
	}{
		{
			name: "empty name",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty organization",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: "",
			},
			wantErr: true,
		},
		{
			name: "already existing, error",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectRequest{
				Name:           alreadyExistingProjectName,
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "empty, ok",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			want: want{
				id:           true,
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.CreateProject(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateProjectResponse(t, creationDate, changeDate, tt.want.creationDate, tt.want.id, got)
		})
	}
}

func TestServer_CreateProject_Permission(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type want struct {
		id           bool
		creationDate bool
	}
	tests := []struct {
		name string
		ctx  context.Context
		req  *project.CreateProjectRequest
		want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "missing permission, other organization",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "with ORG_PROJECT_CREATOR permission, same organization, ok",
			ctx:  integration.WithAuthorizationToken(CTX, getOrgProjectCreatorToken(t, iamOwnerCtx, orgResp.GetOrganizationId(), orgResp.GetOrganizationId())),
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			want: want{
				id:           true,
				creationDate: true,
			},
		},
		{
			name: "with ORG_PROJECT_CREATOR permission, other organization, ok",
			ctx:  integration.WithAuthorizationToken(CTX, getOrgProjectCreatorToken(t, iamOwnerCtx, orgResp.GetOrganizationId(), instance.DefaultOrg.GetId())),
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: instance.DefaultOrg.GetId(),
			},
			want: want{
				id:           true,
				creationDate: true,
			},
		},
		{
			name: "organization owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: instance.DefaultOrg.GetId(),
			},
			want: want{
				id:           true,
				creationDate: true,
			},
		},
		{
			name: "instance owner, ok",
			ctx:  iamOwnerCtx,
			req: &project.CreateProjectRequest{
				Name:           integration.ProjectName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			want: want{
				id:           true,
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			creationDate := time.Now().UTC()
			got, err := instance.Client.Projectv2Beta.CreateProject(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateProjectResponse(t, creationDate, changeDate, tt.want.creationDate, tt.want.id, got)
		})
	}
}

func getOrgProjectCreatorToken(t *testing.T, ctx context.Context, orgId1, orgId2 string) string {
	// create a machine user in Org 1
	userResp := instance.CreateUserTypeMachine(ctx, orgId1)

	// assign ORG_PROJECT_CREATOR role in Org 2
	_, err := instance.Client.InternalPermissionv2Beta.CreateAdministrator(ctx, &internal_permission_v2beta.CreateAdministratorRequest{
		Resource: &internal_permission_v2beta.ResourceType{
			Resource: &internal_permission_v2beta.ResourceType_OrganizationId{OrganizationId: orgId2},
		},
		UserId: userResp.GetId(),
		Roles:  []string{domain.RoleOrgProjectCreator},
	})
	require.NoError(t, err)

	return instance.CreatePersonalAccessToken(ctx, userResp.GetId()).Token
}

func assertCreateProjectResponse(t *testing.T, creationDate, changeDate time.Time, expectedCreationDate, expectedID bool, actualResp *project.CreateProjectResponse) {
	if expectedCreationDate {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.CreationDate)
	}

	if expectedID {
		assert.NotEmpty(t, actualResp.GetId())
	} else {
		assert.Nil(t, actualResp.Id)
	}
}

func TestServer_UpdateProject(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.UpdateProjectRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.UpdateProjectRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "not existing",
			prepare: func(request *project.UpdateProjectRequest) {
				request.Id = "notexisting"
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *project.UpdateProjectRequest) {
				name := integration.ProjectName()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), name, false, false).GetId()
				request.Id = projectID
				request.Name = gu.Ptr(name)
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRequest{},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "change name, ok",
			prepare: func(request *project.UpdateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "change full, ok",
			prepare: func(request *project.UpdateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRequest{
					Name:                   gu.Ptr(integration.ProjectName()),
					ProjectRoleAssertion:   gu.Ptr(true),
					ProjectRoleCheck:       gu.Ptr(true),
					HasProjectCheck:        gu.Ptr(true),
					PrivateLabelingSetting: gu.Ptr(project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY),
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
			t.Parallel()
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.UpdateProject(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateProjectResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_UpdateProject_Permission(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userResp := instance.CreateMachineUser(iamOwnerCtx)
	patResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userResp.GetUserId())
	projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
	instance.CreateProjectMembership(t, iamOwnerCtx, projectID, userResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patResp.Token)

	type args struct {
		ctx context.Context
		req *project.UpdateProjectRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.UpdateProjectRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			prepare: func(request *project.UpdateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: CTX,
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			prepare: func(request *project.UpdateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			wantErr: true,
		},
		{
			name: "project owner, no permission",
			prepare: func(request *project.UpdateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: projectOwnerCtx,
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			wantErr: true,
		},
		{
			name: " roject owner, ok",
			prepare: func(request *project.UpdateProjectRequest) {
				request.Id = projectID
			},
			args: args{
				ctx: projectOwnerCtx,
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "missing permission, other organization",
			prepare: func(request *project.UpdateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			wantErr: true,
		},
		{
			name: "organization owner, ok",
			prepare: func(request *project.UpdateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner, ok",
			prepare: func(request *project.UpdateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.UpdateProjectRequest{
					Name: gu.Ptr(integration.ProjectName()),
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
			t.Parallel()
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.UpdateProject(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateProjectResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertUpdateProjectResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.UpdateProjectResponse) {
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

func TestServer_DeleteProject(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *project.DeleteProjectRequest) (time.Time, time.Time)
		req              *project.DeleteProjectRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "empty id",
			ctx:  iamOwnerCtx,
			req: &project.DeleteProjectRequest{
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "delete, not existing",
			ctx:  iamOwnerCtx,
			req: &project.DeleteProjectRequest{
				Id: "notexisting",
			},
			wantDeletionDate: false,
		},
		{
			name: "delete",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectRequest{},
			wantDeletionDate: true,
		},
		{
			name: "delete, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				instance.DeleteProject(iamOwnerCtx, t, projectID)
				return creationDate, time.Now().UTC()
			},
			req:              &project.DeleteProjectRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := instance.Client.Projectv2Beta.DeleteProject(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteProjectResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func TestServer_DeleteProject_Permission(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userResp := instance.CreateMachineUser(iamOwnerCtx)
	patResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userResp.GetUserId())
	projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
	instance.CreateProjectMembership(t, iamOwnerCtx, projectID, userResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patResp.Token)

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *project.DeleteProjectRequest) (time.Time, time.Time)
		req              *project.DeleteProjectRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:     &project.DeleteProjectRequest{},
			wantErr: true,
		},
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:     &project.DeleteProjectRequest{},
			wantErr: true,
		},
		{
			name: "project owner, no permission",
			ctx:  projectOwnerCtx,
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:     &project.DeleteProjectRequest{},
			wantErr: true,
		},
		{
			name: "project owner, ok",
			ctx:  projectOwnerCtx,
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectRequest{},
			wantDeletionDate: true,
		},
		{
			name: "organization owner, other org",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:     &project.DeleteProjectRequest{},
			wantErr: true,
		},
		{
			name: "organization owner",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectRequest{},
			wantDeletionDate: true,
		},
		{
			name: "instance owner",
			ctx:  iamOwnerCtx,
			prepare: func(request *project.DeleteProjectRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				projectID := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				return creationDate, time.Time{}
			},
			req:              &project.DeleteProjectRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := instance.Client.Projectv2Beta.DeleteProject(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteProjectResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func assertDeleteProjectResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *project.DeleteProjectResponse) {
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

func TestServer_DeactivateProject(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.DeactivateProjectRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.DeactivateProjectRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "not existing",
			prepare: func(request *project.DeactivateProjectRequest) {
				request.Id = "notexisting"
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				instance.DeactivateProject(iamOwnerCtx, t, projectID)
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "change, ok",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.DeactivateProject(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertDeactivateProjectResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_DeactivateProject_Permission(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.DeactivateProjectRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.DeactivateProjectRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: CTX,
				req: &project.DeactivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.DeactivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.DeactivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.DeactivateProjectRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner",
			prepare: func(request *project.DeactivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.DeactivateProjectRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.DeactivateProject(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertDeactivateProjectResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertDeactivateProjectResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.DeactivateProjectResponse) {
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

func TestServer_ActivateProject(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.ActivateProjectRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.ActivateProjectRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *project.ActivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.ActivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *project.ActivateProjectRequest) {
				request.Id = "notexisting"
				return
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *project.ActivateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "change, ok",
			prepare: func(request *project.ActivateProjectRequest) {
				projectID := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false).GetId()
				request.Id = projectID
				instance.DeactivateProject(iamOwnerCtx, t, projectID)
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.ActivateProject(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertActivateProjectResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_ActivateProject_Permission(t *testing.T) {
	t.Parallel()
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	type args struct {
		ctx context.Context
		req *project.ActivateProjectRequest
	}
	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		prepare func(request *project.ActivateProjectRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			prepare: func(request *project.ActivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
				instance.DeactivateProject(iamOwnerCtx, t, projectResp.GetId())
			},
			args: args{
				ctx: CTX,
				req: &project.ActivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			prepare: func(request *project.ActivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
				instance.DeactivateProject(iamOwnerCtx, t, projectResp.GetId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
				req: &project.ActivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner, other org",
			prepare: func(request *project.ActivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
				instance.DeactivateProject(iamOwnerCtx, t, projectResp.GetId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.ActivateProjectRequest{},
			},
			wantErr: true,
		},
		{
			name: "organization owner",
			prepare: func(request *project.ActivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
				instance.DeactivateProject(iamOwnerCtx, t, projectResp.GetId())
			},
			args: args{
				ctx: instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
				req: &project.ActivateProjectRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance owner",
			prepare: func(request *project.ActivateProjectRequest) {
				projectResp := instance.CreateProject(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.ProjectName(), false, false)
				request.Id = projectResp.GetId()
				instance.DeactivateProject(iamOwnerCtx, t, projectResp.GetId())
			},
			args: args{
				ctx: iamOwnerCtx,
				req: &project.ActivateProjectRequest{},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.Projectv2Beta.ActivateProject(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertActivateProjectResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func assertActivateProjectResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *project.ActivateProjectResponse) {
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
