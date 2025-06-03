//go:build integration

package authorization_test

import (
	"context"
	"github.com/muhlemmer/gu"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
)

func TestServer_CreateAuthorization(t *testing.T) {
	type args struct {
		prepare func(*testing.T, *authorization.CreateAuthorizationRequest) context.Context
	}
	tests := []struct {
		name    string
		skip    string
		args    args
		wantErr bool
	}{
		{
			name: "add authorization, project owned, PROJECT_OWNER, ok",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					request.ProjectId = Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, project owned, PROJECT_OWNER, no org id, ok",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.ProjectId = Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, project owned, ORG_OWNER, ok",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					request.ProjectId = Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateOrgMembership(t, IAMCTX, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, project owned, no permission, error",
			args: args{

				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					request.ProjectId = Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, role does not exist, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					request.ProjectId = Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, project does not exist, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					request.ProjectId = gofakeit.AppName()
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, org does not exist, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = gu.Ptr(gofakeit.AppName())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, project granted, ok",
			skip: "fails because of a bug in CreateProjectGrant on a foreign org",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					foreignOrg := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					Instance.CreateProjectGrant(IAMCTX, t, request.ProjectId, Instance.DefaultOrg.Id, request.RoleKeys...)
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, role key not granted, error",
			skip: "fails because of a bug in CreateProjectGrant on a foreign org",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					foreignOrg := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					// TODO: This fails because of a bug in CreateProjectGrant
					Instance.CreateProjectGrant(IAMCTX, t, request.ProjectId, Instance.DefaultOrg.Id)
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, grant does not exist, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					foreignOrg := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, PROJECT_OWNER on wrong org, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					foreignOrg := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrg.OrganizationId}), request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip != "" {
				t.Skip(tt.skip)
			}
			now := time.Now()
			req := &authorization.CreateAuthorizationRequest{}
			ctx := tt.args.prepare(t, req)
			got, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(ctx, req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.Id, "id is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_UpdateAuthorization(t *testing.T) {
	type args struct {
		prepare func(*testing.T, *authorization.UpdateAuthorizationRequest) context.Context
	}
	tests := []struct {
		name                         string
		skip                         string
		args                         args
		wantErr                      bool
		wantChangedDateDuringPrepare bool
	}{
		{
			name: "update authorization, owned project, ok",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					projectId := Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "update authorization, owned project, role not found, error",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					projectId := Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1, projectRole2, gofakeit.AppName()}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "update authorization, owned project, unchanged, ok, changed date is creation date",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					projectId := Instance.CreateProject(IAMCTX, t, Instance.DefaultOrg.Id, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1, projectRole2}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantChangedDateDuringPrepare: true,
		},
		{
			name: "update authorization, granted project, ok",
			skip: "fails because of a bug in CreateProjectGrant on a foreign org",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					// TODO: This fails because of a bug in CreateProjectGrant
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "update authorization, granted project, role not granted, error",
			skip: "fails because of a bug in CreateProjectGrant on a foreign org",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					// TODO: This fails because of a bug in CreateProjectGrant
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2, projectRole3},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "update authorization, granted project, grant removed, error",
			skip: "fails because of a bug in CreateProjectGrant on a foreign org",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					// TODO: This fails because of a bug in CreateProjectGrant
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					_, err = Instance.Client.Projectv2Beta.DeleteProjectGrant(IAMCTX, &project.DeleteProjectGrantRequest{
						ProjectId:             projectId,
						GrantedOrganizationId: Instance.DefaultOrg.Id,
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.Mgmt.AddPersonalAccessToken(IAMCTX, &management.AddPersonalAccessTokenRequest{UserId: callingUser.UserId})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip != "" {
				t.Skip(tt.skip)
			}
			now := time.Now()
			req := &authorization.UpdateAuthorizationRequest{}
			ctx := tt.args.prepare(t, req)
			afterPrepare := time.Now()
			got, err := Instance.Client.AuthorizationV2Beta.UpdateAuthorization(ctx, req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.ChangeDate, "change date is empty")
			changeDate := got.ChangeDate.AsTime()
			assert.Greater(t, changeDate, now, "change date is before the test started")
			if tt.wantChangedDateDuringPrepare {
				assert.Less(t, changeDate, afterPrepare, "change date is after prepare finished")
			} else {
				assert.Less(t, changeDate, time.Now(), "change date is in the future")
			}
		})
	}
}

/*
func TestServer_RemoveAuthorization(t *testing.T) {
	resp := Instance.CreateUserTypeMachine(IamCTX)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		req     *authorization.RemoveAuthorizationRequest
		prepare func(request *authorization.RemoveAuthorizationRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove authorization, user not existing",
			args: args{
				&authorization.RemoveAuthorizationRequest{
					UserId: "notexisting",
				},
				func(request *authorization.RemoveAuthorizationRequest) error {
					authorization, err := Client.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					request.TokenId = authorization.GetTokenId()
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "remove authorization, not existing",
			args: args{
				&authorization.RemoveAuthorizationRequest{
					TokenId: "notexisting",
				},
				func(request *authorization.RemoveAuthorizationRequest) error {
					request.UserId = userId
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "remove authorization, ok",
			args: args{
				&authorization.RemoveAuthorizationRequest{},
				func(request *authorization.RemoveAuthorizationRequest) error {
					authorization, err := Client.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					request.TokenId = authorization.GetTokenId()
					request.UserId = userId
					return err
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)
			got, err := Client.RemoveAuthorization(IAMCTX, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			deletionDate := got.DeletionDate.AsTime()
			assert.Greater(t, deletionDate, now, "creation date is before the test started")
			assert.Less(t, deletionDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_RemoveAuthorization_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, fmt.Sprintf("RemoveAuthorization-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Client.CreateUser(IamCTX, &authorization.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &authorization.CreateUserRequest_Machine_{
			Machine: &authorization.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
			},
		},
	})
	request := &authorization.RemoveAuthorizationRequest{
		UserId: otherOrgUser.GetId(),
	}
	prepare := func(request *authorization.RemoveAuthorizationRequest) error {
		authorization, err := Client.CreateAuthorization(IamCTX, &authorization.CreateAuthorizationRequest{
			ExpirationDate: timestamppb.New(time.Now().Add(time.Hour * 24)),
			UserId:         otherOrgUser.GetId(),
		})
		request.TokenId = authorization.GetTokenId()
		return err
	}
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *authorization.RemoveAuthorizationRequest
		prepare func(request *authorization.RemoveAuthorizationRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "system, ok",
			args: args{SystemCTX, request, prepare},
		},
		{
			name: "instance, ok",
			args: args{IamCTX, request, prepare},
		},
		{
			name:    "org, error",
			args:    args{IAMCTX, request, prepare},
			wantErr: true,
		},
		{
			name:    "user, error",
			args:    args{UserCTX, request, prepare},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			require.NoError(t, tt.args.prepare(tt.args.req))
			got, err := Client.RemoveAuthorization(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.DeletionDate, "client authorization is empty")
			creationDate := got.DeletionDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_ListAuthorizations(t *testing.T) {
	type args struct {
		ctx context.Context
		req *authorization.ListAuthorizationsRequest
	}
	type testCase struct {
		name string
		args args
		want *authorization.ListAuthorizationsResponse
	}
	OrgCTX := IAMCTX
	otherOrg := Instance.CreateOrganization(SystemCTX, fmt.Sprintf("ListAuthorizations-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Client.CreateUser(SystemCTX, &authorization.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &authorization.CreateUserRequest_Machine_{
			Machine: &authorization.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
			},
		},
	})
	require.NoError(t, err)
	otherOrgUserId := otherOrgUser.GetId()
	otherUserId := Instance.CreateUserTypeMachine(SystemCTX).GetId()
	onlySinceTestStartFilter := &authorization.AuthorizationsSearchFilter{Filter: &authorization.AuthorizationsSearchFilter_CreatedDateFilter{CreatedDateFilter: &filter.TimestampFilter{
		Timestamp: timestamppb.Now(),
		Method:    filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS,
	}}}
	myOrgId := Instance.DefaultOrg.GetId()
	myUserId := Instance.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Hour).Add(time.Hour * 24)
	myDataPoint := setupPATDataPoint(t, myUserId, myOrgId, expiresInADay)
	otherUserDataPoint := setupPATDataPoint(t, otherUserId, myOrgId, expiresInADay)
	otherOrgDataPointExpiringSoon := setupPATDataPoint(t, otherOrgUserId, otherOrg.OrganizationId, time.Now().Truncate(time.Hour).Add(time.Hour))
	otherOrgDataPointExpiringLate := setupPATDataPoint(t, otherOrgUserId, otherOrg.OrganizationId, expiresInADay.Add(time.Hour*24*30))
	sortingColumnExpirationDate := authorization.AuthorizationFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_EXPIRATION_DATE
	awaitAuthorizations(t,
		onlySinceTestStartFilter,
		otherOrgDataPointExpiringSoon.GetId(),
		otherOrgDataPointExpiringLate.GetId(),
		otherUserDataPoint.GetId(),
		myDataPoint.GetId(),
	)
	tests := []testCase{
		{
			name: "list all, instance",
			args: args{
				IamCTX,
				&authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					otherOrgDataPointExpiringLate,
					otherOrgDataPointExpiringSoon,
					otherUserDataPoint,
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, org",
			args: args{
				OrgCTX,
				&authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					otherUserDataPoint,
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, user",
			args: args{
				UserCTX,
				&authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list by id",
			args: args{
				IamCTX,
				&authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &authorization.AuthorizationsSearchFilter_TokenIdFilter{
								TokenIdFilter: &filter.IDFilter{Id: otherOrgDataPointExpiringSoon.Id},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					otherOrgDataPointExpiringSoon,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all from other org",
			args: args{
				IamCTX,
				&authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &authorization.AuthorizationsSearchFilter_OrganizationIdFilter{
								OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId},
							},
						}},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					otherOrgDataPointExpiringLate,
					otherOrgDataPointExpiringSoon,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "sort by next expiration dates",
			args: args{
				IamCTX,
				&authorization.ListAuthorizationsRequest{
					Pagination: &filter.PaginationRequest{
						Asc: true,
					},
					SortingColumn: &sortingColumnExpirationDate,
					Filters: []*authorization.AuthorizationsSearchFilter{
						onlySinceTestStartFilter,
						{Filter: &authorization.AuthorizationsSearchFilter_OrganizationIdFilter{OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId}}},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					otherOrgDataPointExpiringSoon,
					otherOrgDataPointExpiringLate,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "get page",
			args: args{
				IamCTX,
				&authorization.ListAuthorizationsRequest{
					Pagination: &filter.PaginationRequest{
						Offset: 2,
						Limit:  2,
						Asc:    true,
					},
					Filters: []*authorization.AuthorizationsSearchFilter{
						onlySinceTestStartFilter,
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{
					otherOrgDataPointExpiringSoon,
					otherOrgDataPointExpiringLate,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 2,
				},
			},
		},
		{
			name: "empty list",
			args: args{
				UserCTX,
				&authorization.ListAuthorizationsRequest{
					Filters: []*authorization.AuthorizationsSearchFilter{
						{
							Filter: &authorization.AuthorizationsSearchFilter_TokenIdFilter{
								TokenIdFilter: &filter.IDFilter{Id: otherUserDataPoint.Id},
							},
						},
					},
				},
			},
			want: &authorization.ListAuthorizationsResponse{
				Result: []*authorization.Authorization{},
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
	}
	t.Run("with permission flag v2", func(t *testing.T) {
		setPermissionCheckV2Flag(t, true)
		defer setPermissionCheckV2Flag(t, false)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Client.ListAuthorizations(tt.args.ctx, tt.args.req)
				require.NoError(t, err)
				assert.Len(t, got.Result, len(tt.want.Result))
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					t.Errorf("ListAuthorizations() mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})
	t.Run("without permission flag v2", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Client.ListAuthorizations(tt.args.ctx, tt.args.req)
				require.NoError(t, err)
				assert.Len(t, got.Result, len(tt.want.Result))
				// ignore the total result, as this is a known bug with the in-memory permission checks.
				// The command can't know how many keys exist in the system if the SQL statement has a limit.
				// This is fixed, once the in-memory permission checks are removed with https://github.com/zitadel/zitadel/issues/9188
				tt.want.Pagination.TotalResult = got.Pagination.TotalResult
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					t.Errorf("ListAuthorizations() mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})
}

func setupPATDataPoint(t *testing.T, userId, orgId string, expirationDate time.Time) *authorization.Authorization {
	expirationDatePb := timestamppb.New(expirationDate)
	newAuthorization, err := Client.CreateAuthorization(SystemCTX, &authorization.CreateAuthorizationRequest{
		UserId:         userId,
		ExpirationDate: expirationDatePb,
	})
	require.NoError(t, err)
	return &authorization.Authorization{
		CreationDate:   newAuthorization.CreationDate,
		ChangeDate:     newAuthorization.CreationDate,
		Id:             newAuthorization.GetTokenId(),
		UserId:         userId,
		OrganizationId: orgId,
		ExpirationDate: expirationDatePb,
	}
}

func awaitAuthorizations(t *testing.T, sinceTestStartFilter *authorization.AuthorizationsSearchFilter, authorizationIds ...string) {
	sortingColumn := authorization.AuthorizationFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_ID
	slices.Sort(authorizationIds)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		result, err := Client.ListAuthorizations(SystemCTX, &authorization.ListAuthorizationsRequest{
			Filters:       []*authorization.AuthorizationsSearchFilter{sinceTestStartFilter},
			SortingColumn: &sortingColumn,
			Pagination: &filter.PaginationRequest{
				Asc: true,
			},
		})
		require.NoError(t, err)
		if !assert.Len(collect, result.Result, len(authorizationIds)) {
			return
		}
		for i := range authorizationIds {
			authorizationId := authorizationIds[i]
			require.Equal(collect, authorizationId, result.Result[i].GetId())
		}
	}, 5*time.Second, time.Second, "authorization not created in time")
}
*/
