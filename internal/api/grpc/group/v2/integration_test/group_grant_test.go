//go:build integration

package group_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func createProjectWithRole(ctx context.Context, t *testing.T, orgID, roleKey string) string {
	project := instance.CreateProject(ctx, t, orgID, integration.ProjectName(), false, false)
	instance.AddProjectRole(ctx, t, project.GetId(), roleKey, roleKey, "")
	return project.GetId()
}

func TestServer_CreateGroupGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	groupResp := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.GroupName())
	projectID := createProjectWithRole(iamOwnerCtx, t, orgResp.GetOrganizationId(), "role1")

	tests := []struct {
		name        string
		ctx         context.Context
		req         *group_v2.CreateGroupGrantRequest
		wantErrCode codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.CreateGroupGrantRequest{
				GroupId:   groupResp.GetId(),
				ProjectId: projectID,
				RoleKeys:  []string{"role1"},
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "missing role keys, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupGrantRequest{
				GroupId:   groupResp.GetId(),
				ProjectId: projectID,
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "role does not exist, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupGrantRequest{
				GroupId:   groupResp.GetId(),
				ProjectId: projectID,
				RoleKeys:  []string{"does-not-exist"},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "group grant created, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupGrantRequest{
				GroupId:   groupResp.GetId(),
				ProjectId: projectID,
				RoleKeys:  []string{"role1"},
			},
		},
		{
			name: "duplicate group grant, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupGrantRequest{
				GroupId:   groupResp.GetId(),
				ProjectId: projectID,
				RoleKeys:  []string{"role1"},
			},
			wantErrCode: codes.AlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := instance.Client.GroupV2.CreateGroupGrant(tt.ctx, tt.req)
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, resp.GetId())
			assert.NotNil(t, resp.GetCreationDate())
		})
	}
}

func TestServer_ListGroupGrants(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	groupName := integration.GroupName()
	groupResp := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)
	projectID := createProjectWithRole(iamOwnerCtx, t, orgResp.GetOrganizationId(), "role1")

	grantResp, err := instance.Client.GroupV2.CreateGroupGrant(iamOwnerCtx, &group_v2.CreateGroupGrantRequest{
		GroupId:   groupResp.GetId(),
		ProjectId: projectID,
		RoleKeys:  []string{"role1"},
	})
	require.NoError(t, err)

	req := &group_v2.ListGroupGrantsRequest{
		Filters: []*group_v2.GroupGrantsSearchFilter{
			{
				Filter: &group_v2.GroupGrantsSearchFilter_GroupIds{
					GroupIds: &filter.InIDsFilter{Ids: []string{groupResp.GetId()}},
				},
			},
		},
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := instance.Client.GroupV2.ListGroupGrants(iamOwnerCtx, req)
		require.NoError(ttt, err)
		require.Len(ttt, got.GetGroupGrants(), 1)
		grant := got.GetGroupGrants()[0]
		// provenance: the supplying group is named on the grant
		assert.Equal(ttt, grantResp.GetId(), grant.GetId())
		assert.Equal(ttt, groupResp.GetId(), grant.GetGroupId())
		assert.Equal(ttt, groupName, grant.GetGroupName())
		assert.Equal(ttt, orgResp.GetOrganizationId(), grant.GetOrganizationId())
		assert.Equal(ttt, projectID, grant.GetProjectId())
		assert.Equal(ttt, []string{"role1"}, grant.GetRoleKeys())
	}, retryDuration, tick, "timeout waiting for expected result")
}

func TestServer_DeleteGroupGrant(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	groupResp := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), integration.GroupName())
	projectID := createProjectWithRole(iamOwnerCtx, t, orgResp.GetOrganizationId(), "role1")

	grantResp, err := instance.Client.GroupV2.CreateGroupGrant(iamOwnerCtx, &group_v2.CreateGroupGrantRequest{
		GroupId:   groupResp.GetId(),
		ProjectId: projectID,
		RoleKeys:  []string{"role1"},
	})
	require.NoError(t, err)

	deleteResp, err := instance.Client.GroupV2.DeleteGroupGrant(iamOwnerCtx, &group_v2.DeleteGroupGrantRequest{
		Id: grantResp.GetId(),
	})
	require.NoError(t, err)
	assert.NotNil(t, deleteResp.GetDeletionDate())

	// deleting again is idempotent
	_, err = instance.Client.GroupV2.DeleteGroupGrant(iamOwnerCtx, &group_v2.DeleteGroupGrantRequest{
		Id: grantResp.GetId(),
	})
	require.NoError(t, err)

	// after deletion, the grant can be created again (unique constraint released)
	_, err = instance.Client.GroupV2.CreateGroupGrant(iamOwnerCtx, &group_v2.CreateGroupGrantRequest{
		GroupId:   groupResp.GetId(),
		ProjectId: projectID,
		RoleKeys:  []string{"role1"},
	})
	require.NoError(t, err)
}
