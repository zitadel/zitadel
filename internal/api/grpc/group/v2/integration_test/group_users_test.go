//go:build integration

package group_test

import (
	"context"
	"testing"

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

	tests := []struct {
		name              string
		ctx               context.Context
		req               *group_v2.AddUsersToGroupRequest
		wantResp          bool
		wantFailedUserIDs []string
		wantErrCode       codes.Code
		wantErrMsg        string
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.Unauthenticated,
			wantErrMsg:  "auth header missing",
		},
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "invalid AddUsersToGroupRequest.Id: value length must be between 1 and 200 runes, inclusive",
		},
		{
			name: "missing user ids, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id: group.GetId(),
			},
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "invalid AddUsersToGroupRequest.UserIds: value must contain at least 1 item(s)",
		},
		{
			name: "group does not exist, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      "randomGroup",
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.FailedPrecondition,
			wantErrMsg:  "Errors.Group.NotFound (CMDGRP-eQfeur)",
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.NotFound,
			wantErrMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			name: "organization owner, missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantErrCode: codes.NotFound,
			wantErrMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			name: "organization owner, with permission, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.AddUsersToGroupRequest{
				Id:      defOrgGroup.GetId(),
				UserIds: []string{defOrgUser.GetUserId()},
			},
			wantResp: true,
		},
		{
			name: "some users already in the group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId(), user3.GetUserId()},
			},
			wantResp: true,
		},
		{
			name: "some users not found, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId(), "randomUser"},
			},
			wantResp:          true,
			wantFailedUserIDs: []string{"randomUser"},
		},
		{
			name: "no users were found, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{"randomUser1", "randomUser2"},
			},
			wantResp:          true,
			wantFailedUserIDs: []string{"randomUser1", "randomUser2"},
		},
		{
			name: "user in a different org not added, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId(), defOrgUser.GetUserId()},
			},
			wantResp:          true,
			wantFailedUserIDs: []string{defOrgUser.GetUserId()},
		},
		{
			name: "add all users to group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.AddUsersToGroupRequest{
				Id:      group.GetId(),
				UserIds: []string{user1.GetUserId(), user2.GetUserId()},
			},
			wantResp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.GroupV2.AddUsersToGroup(tt.ctx, tt.req)
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetChangeDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				assert.Equal(t, tt.wantErrMsg, status.Convert(err).Message())
				return
			}
			require.NoError(t, err)
			require.True(t, tt.wantResp)
			require.NotEmpty(t, got.GetChangeDate())
			assert.Equal(t, tt.wantFailedUserIDs, got.FailedUserIds)
		})
	}
}
