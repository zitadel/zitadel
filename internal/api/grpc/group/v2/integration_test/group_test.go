//go:build integration

package group_test

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func TestServer_CreateGroup(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	orgResp := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	tests := []struct {
		name        string
		ctx         context.Context
		req         *group.CreateGroupRequest
		wantResp    bool
		wantGroupID string
		wantErr     bool
	}{
		{
			name: "invalid name, error",
			ctx:  iamOwnerCtx,
			req: &group.CreateGroupRequest{
				Name:           " ",
				OrganizationId: "org1",
			},
			wantErr: true,
		},
		{
			name: "missing organization id, error",
			ctx:  iamOwnerCtx,
			req: &group.CreateGroupRequest{
				Name: "example",
			},
			wantErr: true,
		},
		{
			name: "organization not found, error",
			ctx:  iamOwnerCtx,
			req: &group.CreateGroupRequest{
				Name:           "example",
				OrganizationId: "org1",
			},
			wantErr: true,
		},
		{
			name: "create group with ID, ok",
			ctx:  iamOwnerCtx,
			req: &group.CreateGroupRequest{
				Id:             gu.Ptr("1234"),
				Name:           "example",
				OrganizationId: orgResp.GetOrganizationId(),
			},
			wantResp:    true,
			wantGroupID: "1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instance.Client.GroupV2.CreateGroup(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.wantGroupID != "" {
				assert.Equal(t, tt.wantGroupID, got.Id, "want: %v, got: %v", tt.wantGroupID, got)
			}

			if tt.wantResp {
				assert.NotEmpty(t, got.Id)
			} else {
				assert.Empty(t, got.Id)
			}
		})
	}

}
