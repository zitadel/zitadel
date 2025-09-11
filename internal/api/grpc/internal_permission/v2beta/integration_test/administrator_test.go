//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/integration"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

func TestServer_CreateAdministrator(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *internal_permission.CreateAdministratorRequest)
		req     *internal_permission.CreateAdministratorRequest
		want
		wantErr bool
	}{
		{
			name: "empty user ID",
			ctx:  iamOwnerCtx,
			req: &internal_permission.CreateAdministratorRequest{
				UserId: "",
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "empty roles",
			ctx:  iamOwnerCtx,
			req: &internal_permission.CreateAdministratorRequest{
				UserId: "notexisting",
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{},
			},
			wantErr: true,
		},
		{
			name: "empty resource",
			ctx:  iamOwnerCtx,
			req: &internal_permission.CreateAdministratorRequest{
				UserId: "notexisting",
				Roles:  []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "already existing, error",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "instance, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "instance, not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
		{
			name: "org, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: "notexisting",
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "org, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, no existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
		{
			name: "project, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: "notexisting",
					},
				},
				Roles: []string{"PROJECT_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "project, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
			},
			req: &internal_permission.CreateAdministratorRequest{
				Roles: []string{"PROJECT_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project, not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
			},
			req: &internal_permission.CreateAdministratorRequest{
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
		{
			name: "project grant, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: "notexisting",
						},
					},
				}
			},
			req: &internal_permission.CreateAdministratorRequest{
				Roles: []string{"PROJECT_GRANT_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "project grant, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
			},
			req: &internal_permission.CreateAdministratorRequest{
				Roles: []string{"PROJECT_GRANT_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project grant, not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
			},
			req: &internal_permission.CreateAdministratorRequest{
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}
			creationDate := time.Now().UTC()
			got, err := instance.Client.InternalPermissionv2Beta.CreateAdministrator(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateAdministratorResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func TestServer_CreateAdministrator_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userProjectResp := instance.CreateMachineUser(iamOwnerCtx)
	projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
	instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
	userProjectGrantResp := instance.CreateMachineUser(iamOwnerCtx)
	instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
	patProjectGrantResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectGrantResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectGrantResp.Token)

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *internal_permission.CreateAdministratorRequest)
		req     *internal_permission.CreateAdministratorRequest
		want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "instance, missing permission, org owner",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "instance, instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, org owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "project, project owner, ok",
			ctx:  projectOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				},
				Roles: []string{"PROJECT_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				},
				Roles: []string{"PROJECT_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "project grant, project grant owner, ok",
			ctx:  projectGrantOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
				Roles: []string{"PROJECT_GRANT_OWNER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project grant, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())

				request.UserId = userResp.GetId()
			},
			req: &internal_permission.CreateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
				Roles: []string{"PROJECT_GRANT_OWNER"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}
			creationDate := time.Now().UTC()
			got, err := instance.Client.InternalPermissionv2Beta.CreateAdministrator(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateAdministratorResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func assertCreateAdministratorResponse(t *testing.T, creationDate, changeDate time.Time, expectedCreationDate bool, actualResp *internal_permission.CreateAdministratorResponse) {
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

func TestServer_UpdateAdministrator(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type want struct {
		change     bool
		changeDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *internal_permission.UpdateAdministratorRequest)
		req     *internal_permission.UpdateAdministratorRequest
		want
		wantErr bool
	}{
		{
			name: "empty user ID",
			ctx:  iamOwnerCtx,
			req: &internal_permission.UpdateAdministratorRequest{
				UserId: "",
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "empty roles",
			ctx:  iamOwnerCtx,
			req: &internal_permission.UpdateAdministratorRequest{
				UserId: "notexisting",
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{},
			},
			wantErr: true,
		},
		{
			name: "empty resource",
			ctx:  iamOwnerCtx,
			req: &internal_permission.UpdateAdministratorRequest{
				UserId: "notexisting",
				Roles:  []string{"IAM_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "instance, no change",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER"},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "instance, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER_VIEWER"},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "instance, not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
		{
			name: "org, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: "notexisting",
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "org, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER_VIEWER"},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "org, no change",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER"},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "org, no existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
		{
			name: "project, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: "notexisting",
					},
				},
				Roles: []string{"PROJECT_OWNER"},
			},
			wantErr: true,
		},
		{
			name: "project, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"PROJECT_OWNER_VIEWER"},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "project, no change",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"PROJECT_OWNER"},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "project, not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
		{
			name: "project grant, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: "notexisting",
						},
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"PROJECT_GRANT_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "project grant, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"PROJECT_GRANT_OWNER_VIEWER"},
			},
			want: want{
				change:     true,
				changeDate: true,
			},
		},
		{
			name: "project grant, no change",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"PROJECT_GRANT_OWNER"},
			},
			want: want{
				change:     false,
				changeDate: true,
			},
		},
		{
			name: "project grant, not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Roles: []string{"notexisting"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}
			got, err := instance.Client.InternalPermissionv2Beta.UpdateAdministrator(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateAdministratorResponse(t, creationDate, changeDate, tt.want.changeDate, got)
		})
	}
}

func TestServer_UpdateAdministrator_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userProjectResp := instance.CreateMachineUser(iamOwnerCtx)
	projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
	instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
	userProjectGrantResp := instance.CreateMachineUser(iamOwnerCtx)
	instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
	patProjectGrantResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectGrantResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectGrantResp.Token)

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *internal_permission.UpdateAdministratorRequest)
		req     *internal_permission.UpdateAdministratorRequest
		want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "instance, missing permission, org owner",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "instance, instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
				Roles: []string{"IAM_OWNER_VIEWER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER_VIEWER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, org owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER_VIEWER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
				Roles: []string{"ORG_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "project, project owner, ok",
			ctx:  projectOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				},
				Roles: []string{"PROJECT_OWNER_VIEWER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				},
				Roles: []string{"PROJECT_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "project grant, project grant owner, ok",
			ctx:  projectGrantOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
				Roles: []string{"PROJECT_GRANT_OWNER_VIEWER"},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project grant, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
				Roles: []string{"PROJECT_GRANT_OWNER_VIEWER"},
			},
			wantErr: true,
		},
		{
			name: "project grant, project owner, error",
			ctx:  projectOwnerCtx,
			prepare: func(request *internal_permission.UpdateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.UpdateAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
				Roles: []string{"PROJECT_GRANT_OWNER_VIEWER"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}
			creationDate := time.Now().UTC()
			got, err := instance.Client.InternalPermissionv2Beta.UpdateAdministrator(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertUpdateAdministratorResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func assertUpdateAdministratorResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *internal_permission.UpdateAdministratorResponse) {
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

func TestServer_DeleteAdministrator(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time)
		req              *internal_permission.DeleteAdministratorRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "empty user ID",
			ctx:  iamOwnerCtx,
			req: &internal_permission.DeleteAdministratorRequest{
				UserId: "",
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "instance, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
				return creationDate, time.Now().UTC()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantDeletionDate: false,
		},
		{
			name: "instance, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				instance.DeleteInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
				return creationDate, time.Now().UTC()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantDeletionDate: true,
		},
		{
			name: "instance, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
				return creationDate, time.Time{}
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantDeletionDate: true,
		},
		{
			name: "org, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
				return creationDate, time.Now().UTC()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: "notexisting",
					},
				},
			},
			wantDeletionDate: false,
		},
		{
			name: "org, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
				return creationDate, time.Time{}
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
			},
			wantDeletionDate: true,
		},
		{
			name: "org, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				instance.DeleteOrgMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
				return creationDate, time.Now().UTC()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
			},
			wantDeletionDate: true,
		},
		{
			name: "project, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				request.UserId = userResp.GetId()
				return creationDate, time.Now().UTC()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: "notexisting",
					},
				},
			},
			wantDeletionDate: false,
		},
		{
			name: "project, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
				return creationDate, time.Time{}
			},
			req:              &internal_permission.DeleteAdministratorRequest{},
			wantDeletionDate: true,
		},
		{
			name: "project, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())
				instance.DeleteProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				}
				return creationDate, time.Now().UTC()
			},
			req:              &internal_permission.DeleteAdministratorRequest{},
			wantDeletionDate: true,
		},
		{
			name: "project grant, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: "notexisting",
						},
					},
				}
				return creationDate, time.Time{}
			},
			req:              &internal_permission.DeleteAdministratorRequest{},
			wantDeletionDate: false,
		},
		{
			name: "project grant, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
				return creationDate, time.Time{}
			},
			req:              &internal_permission.DeleteAdministratorRequest{},
			wantDeletionDate: true,
		},
		{
			name: "project grant, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
				instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				instance.DeleteProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())

				request.UserId = userResp.GetId()
				request.Resource = &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				}
				return creationDate, time.Now().UTC()
			},
			req:              &internal_permission.DeleteAdministratorRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := instance.Client.InternalPermissionv2Beta.DeleteAdministrator(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteAdministratorResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func TestServer_DeleteAdministrator_Permission(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	userProjectResp := instance.CreateMachineUser(iamOwnerCtx)
	projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
	instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userProjectResp.GetUserId())
	patProjectResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectResp.Token)

	instance.CreateProjectGrant(iamOwnerCtx, t, projectResp.GetId(), orgResp.GetOrganizationId())
	userProjectGrantResp := instance.CreateMachineUser(iamOwnerCtx)
	instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userProjectGrantResp.GetUserId())
	patProjectGrantResp := instance.CreatePersonalAccessToken(iamOwnerCtx, userProjectGrantResp.GetUserId())
	projectGrantOwnerCtx := integration.WithAuthorizationToken(CTX, patProjectGrantResp.Token)

	type want struct {
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(request *internal_permission.DeleteAdministratorRequest)
		req     *internal_permission.DeleteAdministratorRequest
		want
		wantErr bool
	}{
		{
			name: "unauthenticated",
			ctx:  CTX,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "instance, missing permission, org owner",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "instance, instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateInstanceMembership(t, iamOwnerCtx, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_Instance{
						Instance: true,
					},
				},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, instance owner, ok",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, org owner, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "org, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateOrgMembership(t, iamOwnerCtx, instance.DefaultOrg.Id, userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_OrganizationId{
						OrganizationId: instance.DefaultOrg.GetId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "project, project owner, ok",
			ctx:  projectOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectMembership(t, iamOwnerCtx, projectResp.GetId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectId{
						ProjectId: projectResp.GetId(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "project grant, project grant owner, ok",
			ctx:  projectGrantOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
			},
			want: want{
				creationDate: true,
			},
		},
		{
			name: "project grant, missing permission",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "project grant, project owner, error",
			ctx:  projectOwnerCtx,
			prepare: func(request *internal_permission.DeleteAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx, integration.Email())
				instance.CreateProjectGrantMembership(t, iamOwnerCtx, projectResp.GetId(), orgResp.GetOrganizationId(), userResp.GetId())
				request.UserId = userResp.GetId()
			},
			req: &internal_permission.DeleteAdministratorRequest{
				Resource: &internal_permission.ResourceType{
					Resource: &internal_permission.ResourceType_ProjectGrant_{
						ProjectGrant: &internal_permission.ResourceType_ProjectGrant{
							ProjectId:      projectResp.GetId(),
							ProjectGrantId: orgResp.GetOrganizationId(),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(tt.req)
			}
			creationDate := time.Now().UTC()
			got, err := instance.Client.InternalPermissionv2Beta.DeleteAdministrator(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteAdministratorResponse(t, creationDate, changeDate, tt.want.creationDate, got)
		})
	}
}

func assertDeleteAdministratorResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *internal_permission.DeleteAdministratorResponse) {
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
