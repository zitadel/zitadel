//go:build integration

package project_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/integration"
	internal_permission "github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2beta"
)

func TestServer_CreateAdministrator(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

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
			name: "not existing roles",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			name: "already existing, error",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			name: "org, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			name: "project, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.AppName(), false, false)

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
			name: "project grant, not existing",
			ctx:  iamOwnerCtx,
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.AppName(), false, false)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
				projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.AppName(), false, false)
				orgResp := instance.CreateOrganization(iamOwnerCtx, gofakeit.Company(), gofakeit.Email())
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
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, gofakeit.AppName(), gofakeit.Email())

	userProjectResp := instance.CreateMachineUser(iamOwnerCtx)
	projectResp := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.AppName(), false, false)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeOrgOwner),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)
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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)

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
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)

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
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)

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
			ctx:  instance.WithAuthorization(CTX, integration.UserTypeNoPermission),
			prepare: func(request *internal_permission.CreateAdministratorRequest) {
				userResp := instance.CreateUserTypeHuman(iamOwnerCtx)

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
