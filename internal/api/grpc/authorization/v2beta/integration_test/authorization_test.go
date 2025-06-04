//go:build integration

package authorization_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/integration"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_CreateAuthorization(t *testing.T) {
	type args struct {
		prepare func(*testing.T, *authorization.CreateAuthorizationRequest) context.Context
	}
	tests := []struct {
		name    string
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, project granted, ok",
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, role key not granted, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					request.OrganizationId = &Instance.DefaultOrg.Id
					foreignOrg := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, gofakeit.AppName(), false, false).Id
					request.RoleKeys = []string{gofakeit.AppName()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], gofakeit.AppName(), "")
					Instance.CreateProjectGrant(IAMCTX, t, request.ProjectId, Instance.DefaultOrg.Id)
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantChangedDateDuringPrepare: true,
		},
		{
			name: "update authorization, granted project, ok",
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
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, Instance.DefaultOrg.Id, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "update authorization, granted project, role not granted, error",
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
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectId, projectRole3, projectRole3}
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, Instance.DefaultOrg.Id, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "update authorization, granted project, grant removed, error",
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
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, Instance.DefaultOrg.Id, callingUser.UserId)
					_, err = Instance.Client.Projectv2Beta.DeleteProjectGrant(IAMCTX, &project.DeleteProjectGrantRequest{
						ProjectId:             projectId,
						GrantedOrganizationId: Instance.DefaultOrg.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestServer_DeleteAuthorization(t *testing.T) {
	type args struct {
		prepare func(*testing.T, *authorization.DeleteAuthorizationRequest) context.Context
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		wantDeletionDateDuringPrepare bool
	}{
		{
			name: "delete authorization, project owned by calling users org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		}, {
			name: "delete authorization, owned project, user membership on project owning org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "delete authorization, granted project, user membership on project owning org, error",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "delete authorization, granted project, user membership on project granted org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, Instance.DefaultOrg.Id, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "delete authorization, already deleted, ok, deletion date is creation date",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
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
					_, err = Instance.Client.AuthorizationV2Beta.DeleteAuthorization(IAMCTX, &authorization.DeleteAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantDeletionDateDuringPrepare: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			req := &authorization.DeleteAuthorizationRequest{}
			ctx := tt.args.prepare(t, req)
			afterPrepare := time.Now()
			got, err := Instance.Client.AuthorizationV2Beta.DeleteAuthorization(ctx, req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.DeletionDate, "deletion date is empty")
			changeDate := got.DeletionDate.AsTime()
			assert.Greater(t, changeDate, now, "deletion date is before the test started")
			if tt.wantDeletionDateDuringPrepare {
				assert.Less(t, changeDate, afterPrepare, "deletion date is after prepare finished")
			} else {
				assert.Less(t, changeDate, time.Now(), "deletion date is in the future")
			}
		})
	}
}

func TestServer_DeactivateAuthorization(t *testing.T) {
	type args struct {
		prepare func(*testing.T, *authorization.DeactivateAuthorizationRequest) context.Context
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		wantDeletionDateDuringPrepare bool
	}{
		{
			name: "deactivate authorization, project owned by calling users org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		}, {
			name: "deactivate authorization, owned project, user membership on project owning org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "deactivate authorization, granted project, user membership on project owning org, error",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "deactivate authorization, granted project, user membership on project granted org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, Instance.DefaultOrg.Id, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "deactivate authorization, already inactive, ok, change date is creation date",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
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
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantDeletionDateDuringPrepare: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			req := &authorization.DeactivateAuthorizationRequest{}
			ctx := tt.args.prepare(t, req)
			afterPrepare := time.Now()
			got, err := Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(ctx, req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.ChangeDate, "change date is empty")
			changeDate := got.ChangeDate.AsTime()
			assert.Greater(t, changeDate, now, "change date is before the test started")
			if tt.wantDeletionDateDuringPrepare {
				assert.Less(t, changeDate, afterPrepare, "change date is after prepare finished")
			} else {
				assert.Less(t, changeDate, time.Now(), "change date is in the future")
			}
		})
	}
}

func TestServer_ActivateAuthorization(t *testing.T) {
	type args struct {
		prepare func(*testing.T, *authorization.ActivateAuthorizationRequest) context.Context
	}
	tests := []struct {
		name                          string
		args                          args
		wantErr                       bool
		wantDeletionDateDuringPrepare bool
	}{
		{
			name: "activate authorization, project owned by calling users org, ok",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		}, {
			name: "activate authorization, owned project, user membership on project owning org, ok",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.UserId)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "activate authorization, granted project, user membership on project owning org, error",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.UserId)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "activate authorization, granted project, user membership on project granted org, ok",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
					foreignOrgId := Instance.CreateOrganization(IAMCTX, gofakeit.AppName(), gofakeit.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, gofakeit.AppName(), false, false).Id
					projectRole1 := gofakeit.AppName()
					projectRole2 := gofakeit.AppName()
					projectRole3 := gofakeit.AppName()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, Instance.DefaultOrg.Id, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &Instance.DefaultOrg.Id,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, Instance.DefaultOrg.Id, callingUser.UserId)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "activate authorization, already active, ok, change date is creation date",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
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
					callingUser := Instance.CreateMachineUser(IAMCTX)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.UserId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.UserId, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantDeletionDateDuringPrepare: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			req := &authorization.ActivateAuthorizationRequest{}
			ctx := tt.args.prepare(t, req)
			afterPrepare := time.Now()
			got, err := Instance.Client.AuthorizationV2Beta.ActivateAuthorization(ctx, req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.ChangeDate, "change date is empty")
			changeDate := got.ChangeDate.AsTime()
			assert.Greater(t, changeDate, now, "change date is before the test started")
			if tt.wantDeletionDateDuringPrepare {
				assert.Less(t, changeDate, afterPrepare, "change date is after prepare finished")
			} else {
				assert.Less(t, changeDate, time.Now(), "change date is in the future")
			}
		})
	}
}
