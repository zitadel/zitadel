//go:build integration

package authorization_test

import (
	"context"
	"testing"
	"time"

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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					request.ProjectId = Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, project owned, PROJECT_OWNER, no org id, ok",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.ProjectId = Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, project owned, ORG_OWNER, ok",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					request.ProjectId = Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateOrgMembership(t, IAMCTX, selfOrgId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "add authorization, project owned, no permission, error",
			args: args{

				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					request.ProjectId = Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					request.ProjectId = Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					request.ProjectId = "notexists"
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = gu.Ptr("notexists")
					request.ProjectId = Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)

				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, project owner, project granted, no permission",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					foreignOrg := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					Instance.CreateProjectGrant(IAMCTX, t, request.ProjectId, selfOrgId, request.RoleKeys...)
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
			wantErr: true,
		},
		{
			name: "add authorization, role key not granted, error",
			args: args{
				func(t *testing.T, request *authorization.CreateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					foreignOrg := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					Instance.CreateProjectGrant(IAMCTX, t, request.ProjectId, selfOrgId)
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					foreignOrg := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email())
					projectID := Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, integration.ProjectName(), false, false).Id
					request.ProjectId = projectID
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, projectID, request.RoleKeys[0], integration.RoleDisplayName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectID, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					request.OrganizationId = &selfOrgId
					foreignOrg := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email())
					request.ProjectId = Instance.CreateProject(IAMCTX, t, foreignOrg.OrganizationId, integration.ProjectName(), false, false).Id
					request.RoleKeys = []string{integration.RoleKey()}
					Instance.AddProjectRole(IAMCTX, t, request.ProjectId, request.RoleKeys[0], integration.RoleDisplayName(), "")
					request.UserId = Instance.Users.Get(integration.UserTypeIAMOwner).ID
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrg.OrganizationId}), request.ProjectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
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
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "update authorization, owned project, role not found, error",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1, projectRole2, "rolenotfound"}
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
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
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)

					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1}
					token := createUserWithProjectGrantMembership(IAMCTX, t, Instance, projectId, selfOrgId)
					return integration.WithAuthorizationToken(EmptyCTX, token)

				},
			},
		},
		{
			name: "update authorization, granted project, role not granted, error",
			args: args{
				func(t *testing.T, request *authorization.UpdateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectId, projectRole3, projectRole3}
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, selfOrgId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					request.RoleKeys = []string{projectRole1}
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectGrantMembership(t, IAMCTX, projectId, selfOrgId, callingUser.Id)
					_, err = Instance.Client.Projectv2Beta.DeleteProjectGrant(IAMCTX, &project.DeleteProjectGrantRequest{
						ProjectId:             projectId,
						GrantedOrganizationId: selfOrgId,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		}, {
			name: "delete authorization, owned project, user membership on project owning org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "delete authorization, granted project, user membership on project owning org, error",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)

					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					return integration.WithAuthorizationToken(EmptyCTX, createUserWithProjectGrantMembership(IAMCTX, t, Instance, projectId, selfOrgId))

				},
			},
		},
		{
			name: "delete authorization, already deleted, ok, deletion date is creation date",
			args: args{
				func(t *testing.T, request *authorization.DeleteAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
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
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		}, {
			name: "deactivate authorization, owned project, user membership on project owning org, ok",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "deactivate authorization, granted project, user membership on project owning org, error",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)

					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					return integration.WithAuthorizationToken(EmptyCTX, createUserWithProjectGrantMembership(IAMCTX, t, Instance, projectId, selfOrgId))
				},
			},
		},
		{
			name: "deactivate authorization, already inactive, ok, change date is creation date",
			args: args{
				func(t *testing.T, request *authorization.DeactivateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
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
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		}, {
			name: "activate authorization, owned project, user membership on project owning org, ok",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.Id)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, token.Token)
				},
			},
		},
		{
			name: "activate authorization, granted project, user membership on project owning org, error",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, authz.SetCtxData(IAMCTX, authz.CtxData{OrgID: foreignOrgId}), projectId, callingUser.Id)
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					foreignOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, foreignOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					projectRole3 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole3, projectRole3, "")
					Instance.CreateProjectGrant(IAMCTX, t, projectId, selfOrgId, projectRole1, projectRole2)
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:         Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId:      projectId,
						OrganizationId: &selfOrgId,
						RoleKeys:       []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					_, err = Instance.Client.AuthorizationV2Beta.DeactivateAuthorization(IAMCTX, &authorization.DeactivateAuthorizationRequest{
						Id: preparedAuthorization.Id,
					})
					require.NoError(t, err)
					return integration.WithAuthorizationToken(EmptyCTX, createUserWithProjectGrantMembership(IAMCTX, t, Instance, projectId, selfOrgId))
				},
			},
		},
		{
			name: "activate authorization, already active, ok, change date is creation date",
			args: args{
				func(t *testing.T, request *authorization.ActivateAuthorizationRequest) context.Context {
					selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
					projectId := Instance.CreateProject(IAMCTX, t, selfOrgId, integration.ProjectName(), false, false).Id
					projectRole1 := integration.RoleKey()
					projectRole2 := integration.RoleKey()
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole1, projectRole1, "")
					Instance.AddProjectRole(IAMCTX, t, projectId, projectRole2, projectRole2, "")
					preparedAuthorization, err := Instance.Client.AuthorizationV2Beta.CreateAuthorization(IAMCTX, &authorization.CreateAuthorizationRequest{
						UserId:    Instance.Users.Get(integration.UserTypeIAMOwner).ID,
						ProjectId: projectId,
						RoleKeys:  []string{projectRole1, projectRole2},
					})
					require.NoError(t, err)
					request.Id = preparedAuthorization.Id
					callingUser := Instance.CreateUserTypeMachine(IAMCTX, selfOrgId)
					Instance.CreateProjectMembership(t, IAMCTX, projectId, callingUser.Id)
					token, err := Instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
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

func createUserWithProjectGrantMembership(ctx context.Context, t *testing.T, instance *integration.Instance, projectID, grantID string) string {
	selfOrgId := Instance.CreateOrganization(IAMCTX, integration.OrganizationName(), integration.Email()).OrganizationId
	callingUser := instance.CreateUserTypeMachine(ctx, selfOrgId)
	instance.CreateProjectGrantMembership(t, ctx, projectID, grantID, callingUser.Id)
	token, err := instance.Client.UserV2.AddPersonalAccessToken(IAMCTX, &user.AddPersonalAccessTokenRequest{UserId: callingUser.Id, ExpirationDate: timestamppb.New(time.Now().Add(24 * time.Hour))})
	require.NoError(t, err)
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 10*time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		got, err := instance.Client.AuthorizationV2Beta.ListAuthorizations(ctx, &authorization.ListAuthorizationsRequest{
			Filters: nil,
		})
		assert.NoError(tt, err)
		if !assert.NotEmpty(tt, got.Pagination.TotalResult) {
			return
		}
	}, retryDuration, tick)
	return token.GetToken()
}
