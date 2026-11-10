//go:build integration

package group_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func TestServer_CreateGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	alreadyExistingGroupName := integration.GroupName()
	resp := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), alreadyExistingGroupName)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *group_v2.CreateGroupRequest
		wantResp    bool
		wantGroupID string
		wantErrCode codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           " ",
				OrganizationId: "org1",
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "missing organization id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name: integration.GroupName(),
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization not found, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: "org1",
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "instance owner, already existing group (unique name constraint), error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           alreadyExistingGroupName,
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrCode: codes.AlreadyExists,
		},
		{
			name: "instance owner, already existing group ID, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Id:             gu.Ptr(resp.Id),
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrCode: codes.AlreadyExists,
		},
		{
			name: "organization owner, missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, with permission, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: instance.DefaultOrg.GetId(),
			},
			wantResp: true,
		},
		{
			name: "instance owner, create group with ID, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Id:             gu.Ptr("1234"),
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantResp:    true,
			wantGroupID: "1234",
		},
		{
			name: "instance owner, create group without user provided ID, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantResp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeCreationDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.CreateGroup(tt.ctx, tt.req)
			afterCreationDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetId())
				require.Empty(t, got.GetCreationDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			if tt.wantGroupID != "" {
				assert.Equal(t, tt.wantGroupID, got.Id, "want: %v, got: %v", tt.wantGroupID, got)
			}
			require.NotEmpty(t, got.GetId())
			assert.WithinRange(t, got.GetCreationDate().AsTime(), beforeCreationDate, afterCreationDate)
		})
	}
}

func TestServer_UpdateGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	// create a group in the default org
	groupDefOrg := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.GroupName())

	// create a new org
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	// create a group in the new org
	groupName := integration.GroupName()
	existingGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

	tests := []struct {
		name           string
		ctx            context.Context
		req            *group_v2.UpdateGroupRequest
		wantChangeDate bool
		wantErrCode    codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.UpdateGroupRequest{
				Name: gu.Ptr(integration.GroupName()),
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr(" "),
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr("updated group name"),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr("updated group name"),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, with permission, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.UpdateGroupRequest{
				Id:   groupDefOrg.GetId(),
				Name: gu.Ptr("updated group name 1"),
			},
			wantChangeDate: true,
		},
		{
			name: "instance owner, group not found, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   "12345",
				Name: gu.Ptr("updated group name 2"),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "instance owner, no change, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr(groupName),
			},
			wantChangeDate: false,
		},
		{
			name: "instance owner, change name, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr("groupXX"),
			},
			wantChangeDate: true,
		},
		{
			name: "instance owner, change description, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:          existingGroup.GetId(),
				Description: gu.Ptr("description of groupXX"),
			},
			wantChangeDate: true,
		},
		{
			name: "instance owner, full change, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:          existingGroup.GetId(),
				Name:        gu.Ptr("groupXXX"),
				Description: gu.Ptr("description of groupXXX"),
			},
			wantChangeDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeUpdateDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.UpdateGroup(tt.ctx, tt.req)
			afterUpdateDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetChangeDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			if tt.wantChangeDate {
				assert.WithinRange(t, got.GetChangeDate().AsTime(), beforeUpdateDate, afterUpdateDate)
			}
		})
	}
}

func TestServer_DeleteGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	// create a group in the default org
	groupDefOrg := instance.CreateGroup(iamOwnerCtx, t, instance.DefaultOrg.GetId(), integration.GroupName())

	// create a new org
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	// create a group in the new org
	groupName := integration.GroupName()
	existingGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

	// create a group in the new org to be deleted before the test
	deleteGroupName := integration.GroupName()
	deleteGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), deleteGroupName)
	deleteResp := instance.DeleteGroup(iamOwnerCtx, t, deleteGroup.GetId())

	tests := []struct {
		name             string
		ctx              context.Context
		req              *group_v2.DeleteGroupRequest
		wantErrCode      codes.Code
		wantDeletionDate bool
		deletionTime     *timestamp.Timestamp
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.DeleteGroupRequest{
				Id: "12345",
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: "",
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.DeleteGroupRequest{
				Id: existingGroup.GetId(),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.DeleteGroupRequest{
				Id: existingGroup.GetId(),
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "organization owner, with permission, ok",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeOrgOwner),
			req: &group_v2.DeleteGroupRequest{
				Id: groupDefOrg.GetId(),
			},
			wantDeletionDate: true,
		},
		{
			name: "group not found, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: "12345",
			},
		},
		{
			name: "delete, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: existingGroup.GetId(),
			},
			wantDeletionDate: true,
		},
		{
			name: "delete already deleted group, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: deleteGroup.GetId(),
			},
			deletionTime: deleteResp.GetDeletionDate(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeDeletionDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.DeleteGroup(tt.ctx, tt.req)
			afterDeletionDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetDeletionDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got.GetDeletionDate())
			if tt.deletionTime != nil {
				assert.Equal(t, tt.deletionTime, got.GetDeletionDate())
			}
			if tt.wantDeletionDate {
				assert.WithinRange(t, got.GetDeletionDate().AsTime(), beforeDeletionDate, afterDeletionDate)
			}
		})
	}
}
