//go:build integration

package group_test

import (
	"context"
	"testing"

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
		name          string
		ctx           context.Context
		req           *group_v2.CreateGroupRequest
		wantResp      bool
		wantGroupID   string
		wantErrorCode codes.Code
		wantErrMsg    string
	}{
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           " ",
				OrganizationId: "org1",
			},
			wantErrorCode: codes.InvalidArgument,
			wantErrMsg:    "Errors.Group.InvalidName (GROUP-m177lN)",
		},
		{
			name: "missing organization id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name: integration.GroupName(),
			},
			wantErrorCode: codes.InvalidArgument,
			wantErrMsg:    "invalid CreateGroupRequest.OrganizationId: value length must be between 1 and 200 runes, inclusive",
		},
		{
			name: "organization not found, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           integration.GroupName(),
				OrganizationId: "org1",
			},
			wantErrorCode: codes.FailedPrecondition,
			wantErrMsg:    "Organisation not found (CMDGRP-j1mH8l)",
		},
		{
			name: "already existing group (unique name constraint), error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           alreadyExistingGroupName,
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrorCode: codes.AlreadyExists,
			wantErrMsg:    "Errors.Group.AlreadyExists (V3-DKcYh)",
		},
		{
			name: "already existing group ID, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Id:             gu.Ptr(resp.Id),
				Name:           integration.GroupName(),
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErrorCode: codes.AlreadyExists,
			wantErrMsg:    "Errors.Group.AlreadyExists (CMDGRP-shRut3)",
		},
		{
			name: "create group with ID, ok",
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
			name: "create group without user provided ID, ok",
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
			got, err := instance.Client.GroupV2.CreateGroup(tt.ctx, tt.req)
			if tt.wantErrorCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetId())
				require.Empty(t, got.GetCreationDate())
				assert.Equal(t, tt.wantErrorCode, status.Code(err))
				assert.Equal(t, tt.wantErrMsg, status.Convert(err).Message())
				return
			}
			require.NoError(t, err)
			if tt.wantGroupID != "" {
				assert.Equal(t, tt.wantGroupID, got.Id, "want: %v, got: %v", tt.wantGroupID, got)
			}
			require.NotEmpty(t, got.GetId())
			require.NotEmpty(t, got.GetCreationDate())
		})
	}
}

func TestServer_UpdateGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	groupName := integration.GroupName()
	existingGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

	tests := []struct {
		name           string
		ctx            context.Context
		req            *group_v2.UpdateGroupRequest
		wantChangeDate bool
		wantErrCode    codes.Code
		wantErrMsg     string
	}{
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr(" "),
			},
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "Errors.Group.InvalidName (GROUP-dUNd3r)",
		},
		{
			name: "group not found, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   "12345",
				Name: gu.Ptr("updated group name"),
			},
			wantErrCode: codes.NotFound,
			wantErrMsg:  "Errors.Group.NotFound (CMDGRP-b33zly)",
		},
		{
			name: "no change, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr(groupName),
			},
			wantChangeDate: true,
		},
		{
			name: "change name, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr("groupXX"),
			},
			wantChangeDate: true,
		},
		{
			name: "change description, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:          existingGroup.GetId(),
				Description: gu.Ptr("description of groupXX"),
			},
			wantChangeDate: true,
		},
		{
			name: "full change, ok",
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
			got, err := instance.Client.GroupV2.UpdateGroup(tt.ctx, tt.req)
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetChangeDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				assert.Equal(t, tt.wantErrMsg, status.Convert(err).Message())
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got.GetChangeDate())
		})
	}
}

func TestServer_DeleteGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	groupName := integration.GroupName()
	existingGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

	deleteGroupName := integration.GroupName()
	deleteGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), deleteGroupName)
	deleteResp := instance.DeleteGroup(iamOwnerCtx, t, deleteGroup.GetId())

	tests := []struct {
		name         string
		ctx          context.Context
		req          *group_v2.DeleteGroupRequest
		wantErrCode  codes.Code
		wantErrMsg   string
		deletionTime *timestamp.Timestamp
	}{
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: "",
			},
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "invalid DeleteGroupRequest.Id: value length must be between 1 and 200 runes, inclusive",
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
			got, err := instance.Client.GroupV2.DeleteGroup(tt.ctx, tt.req)
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetDeletionDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				assert.Equal(t, tt.wantErrMsg, status.Convert(err).Message())
				return
			}
			require.NoError(t, err)
			require.NotEmpty(t, got.GetDeletionDate())
			if tt.deletionTime != nil {
				assert.Equal(t, tt.deletionTime, got.GetDeletionDate())
			}
		})
	}
}
