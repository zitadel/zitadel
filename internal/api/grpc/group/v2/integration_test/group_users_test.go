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
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func TestServer_AddUsersToGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	// group and user in the default org
	defOrgGroup := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.GroupName())
	defOrgUser := instance.CreateHumanUserVerified(iamOwnerCtx, instance.DefaultOrg.GetId(), integration.Email(), integration.Phone())

	// org1
	org1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	// group in org1
	group := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), integration.GroupName())
	// user1 in org1
	user1 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())
	// user2 in org1
	user2 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())
	// user3 in org1
	user3 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())
	// user4 in org1
	user4 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())

	tests := []struct {
		name           string
		ctx            context.Context
		req            *group_v2.AddUsersToGroupRequest
		wantChangeDate bool
		wantErrCode    codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "missing user ids, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id: group.GetId(),
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "group does not exist, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      "randomGroup",
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, with permission, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      defOrgGroup.GetId(),
				UserIds: []string{defOrgUser.GetUserId()},
			},
		},
		{
			name: "some users not found, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId(), "randomUser"},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "user in a different org not added, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId(), defOrgUser.GetUserId()},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "add all users to group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantChangeDate: true,
		},
		{
			name: "some users already in the group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId(), user3.GetUserId()},
			},
			wantChangeDate: true,
		},
		{
			name: "add all users to group (with duplicate user IDs), ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user4.GetUserId(), user4.GetUserId()},
			},
			wantChangeDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.AddUsersToGroup(tt.ctx, tt.req)
			afterDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetChangeDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			if tt.wantChangeDate {
				assert.WithinRange(t, got.GetChangeDate().AsTime(), beforeDate, afterDate)
			}
		})
	}
}

func TestServer_RemoveUsersFromGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	// group and user in the default org
	defOrgGroup := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.GroupName())
	defOrgUser := instance.CreateHumanUserVerified(iamOwnerCtx, instance.DefaultOrg.GetId(), integration.Email(), integration.Phone())

	// add user to the group in def org
	_, err := instance.Client.GroupV2.AddUsersToGroup(iamOwnerCtx, &group_v2.AddUsersToGroupRequest{
		Id:      defOrgGroup.GetId(),
		UserIds: []string{defOrgUser.GetUserId()},
	})
	require.NoError(t, err)

	// org1
	org1 := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	// group in org1
	group := instance.CreateGroup(iamOwnerCtx, t, org1.GetOrganizationId(), integration.GroupName())
	// user1 in org1
	user1 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())
	// user2 in org1
	user2 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())
	// user3 in org1
	user3 := instance.CreateHumanUserVerified(iamOwnerCtx, org1.OrganizationId, integration.Email(), integration.Phone())

	// add user1, user2, user3 to the group
	_, err = instance.Client.GroupV2.AddUsersToGroup(iamOwnerCtx, &group_v2.AddUsersToGroupRequest{
		Id:      group.GetId(),
		UserIds: []string{user1.GetUserId(), user2.GetUserId(), user3.GetUserId()},
	})
	require.NoError(t, err)

	tests := []struct {
		name           string
		ctx            context.Context
		req            *group_v2.RemoveUsersFromGroupRequest
		wantChangeDate bool
		wantErrCode    codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "missing user ids, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id: group.GetId(),
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "group does not exist, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      "randomGroup",
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, with permission, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      defOrgGroup.GetId(),
				UserIds: []string{defOrgUser.GetUserId()},
			},
			wantChangeDate: true,
		},
		{
			name: "users not in the group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{"user3", "user4"},
			},
			wantChangeDate: false,
		},
		{
			name: "some users not in the group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), defOrgUser.GetUserId(), "user3"},
			},
			wantChangeDate: true,
		},
		{
			name: "users removed, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user2.GetUserId()},
			},
			wantChangeDate: true,
		},
		{
			name: "users removed (with duplicate user IDs), ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.RemoveUsersFromGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user3.GetUserId(), user3.GetUserId()},
			},
			wantChangeDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.RemoveUsersFromGroup(tt.ctx, tt.req)
			afterDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetChangeDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			if tt.wantChangeDate {
				assert.WithinRange(t, got.GetChangeDate().AsTime(), beforeDate, afterDate)
			}
		})
	}
}
