//go:build integration

package group_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func TestServer_CreateGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	alreadyExistingGroupName := "existingGroup1"
	resp := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), alreadyExistingGroupName)

	tests := []struct {
		name        string
		ctx         context.Context
		req         *group_v2.CreateGroupRequest
		wantResp    bool
		wantGroupID string
		wantErr     bool
	}{
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           " ",
				OrganizationId: "org1",
			},
			wantErr: true,
		},
		{
			name: "missing organization id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name: "example",
			},
			wantErr: true,
		},
		{
			name: "organization not found, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           "example",
				OrganizationId: "org1",
			},
			wantErr: true,
		},
		{
			name: "already existing group (unique name constraint), error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           alreadyExistingGroupName,
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "already existing group ID, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Id:             gu.Ptr(resp.Id),
				Name:           "groupY",
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantErr: true,
		},
		{
			name: "create group with ID, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Id:             gu.Ptr("1234"),
				Name:           "group1",
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantResp:    true,
			wantGroupID: "1234",
		},
		{
			name: "create group without user provided ID, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.CreateGroupRequest{
				Name:           "group2",
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantResp: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.CreateGroup(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantGroupID != "" {
				assert.Equal(t, tt.wantGroupID, got.Id, "want: %v, got: %v", tt.wantGroupID, got)
			}
			assertCreateGroupResponse(t, creationDate, changeDate, tt.wantResp, got)
		})
	}
}

func assertCreateGroupResponse(t *testing.T, creationDate, changeDate time.Time, expectResp bool, actualResp *group_v2.CreateGroupResponse) {
	if expectResp {
		assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, changeDate)
		assert.NotEmpty(t, actualResp.GetId())
	} else {
		assert.Nil(t, actualResp.CreationDate)
		assert.Nil(t, actualResp.Id)
	}
}

func TestServer_UpdateGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	groupName := "existingGroup2"
	existingGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

	tests := []struct {
		name           string
		ctx            context.Context
		req            *group_v2.UpdateGroupRequest
		wantChangeDate bool
		wantErr        bool
	}{
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr(" "),
			},
			wantErr: true,
		},
		{
			name: "group not found, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   "12345",
				Name: gu.Ptr("updated group name"),
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.UpdateGroupRequest{
				Id:   existingGroup.GetId(),
				Name: gu.Ptr("groupX"),
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
			creationDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.UpdateGroup(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			changeDate := time.Now().UTC()
			assertUpdateGroupResponse(t, creationDate, changeDate, tt.wantChangeDate, got)
		})
	}
}

func assertUpdateGroupResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate bool, actualResp *group_v2.UpdateGroupResponse) {
	if expectedChangeDate {
		assert.WithinRange(t, actualResp.GetChangeDate().AsTime(), creationDate, changeDate)
	} else {
		assert.Nil(t, actualResp.ChangeDate)
	}
}

func TestServer_DeleteGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	groupName := "existingGroup3"
	existingGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), groupName)

	startDate := time.Now().UTC()
	deleteGroupName := "existingGroup4"
	deleteGroup := instance.CreateGroup(iamOwnerCtx, t, orgResp.GetOrganizationId(), deleteGroupName)
	instance.DeleteGroup(iamOwnerCtx, t, deleteGroup.GetId())

	tests := []struct {
		name             string
		ctx              context.Context
		req              *group_v2.DeleteGroupRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "group not found, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupRequest{
				Id: "12345",
			},
			wantDeletionDate: false,
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
				Id: existingGroup.GetId(),
			},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.GroupV2.DeleteGroup(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assertDeleteProjectResponse(t, startDate, time.Now().UTC(), tt.wantDeletionDate, got)
		})
	}
}

func assertDeleteProjectResponse(t *testing.T, creationDate, deletionDate time.Time, wantDeletionDate bool, actualResp *group_v2.DeleteGroupResponse) {
	if wantDeletionDate {
		assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, deletionDate)
	} else {
		assert.Nil(t, actualResp.DeletionDate)
	}
}
