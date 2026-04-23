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
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
)

func TestServer_SetGroupMetadata(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	org := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	group := instance.CreateGroup(iamOwnerCtx, t, org.GetOrganizationId(), integration.GroupName())

	tests := []struct {
		name        string
		ctx         context.Context
		req         *group_v2.SetGroupMetadataRequest
		wantSetDate bool
		wantErrCode codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.SetGroupMetadataRequest{
				Id: group.GetId(),
				Metadata: []*group_v2.Metadata{
					{Key: "role", Value: []byte("admin")},
				},
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.SetGroupMetadataRequest{
				Metadata: []*group_v2.Metadata{
					{Key: "role", Value: []byte("admin")},
				},
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "empty metadata, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.SetGroupMetadataRequest{
				Id: group.GetId(),
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "group does not exist, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.SetGroupMetadataRequest{
				Id: "randomGroup",
				Metadata: []*group_v2.Metadata{
					{Key: "role", Value: []byte("admin")},
				},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.SetGroupMetadataRequest{
				Id: group.GetId(),
				Metadata: []*group_v2.Metadata{
					{Key: "role", Value: []byte("admin")},
				},
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "single metadata entry, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.SetGroupMetadataRequest{
				Id: group.GetId(),
				Metadata: []*group_v2.Metadata{
					{Key: "role", Value: []byte("admin")},
				},
			},
			wantSetDate: true,
		},
		{
			name: "bulk metadata entries, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.SetGroupMetadataRequest{
				Id: group.GetId(),
				Metadata: []*group_v2.Metadata{
					{Key: "label", Value: []byte("priority")},
					{Key: "owner", Value: []byte("platform")},
					{Key: "region", Value: []byte("eu-west")},
				},
			},
			wantSetDate: true,
		},
		{
			name: "overwrite existing key, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.SetGroupMetadataRequest{
				Id: group.GetId(),
				Metadata: []*group_v2.Metadata{
					{Key: "role", Value: []byte("viewer")},
				},
			},
			wantSetDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.SetGroupMetadata(tt.ctx, tt.req)
			afterDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetSetDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			if tt.wantSetDate {
				assert.WithinRange(t, got.GetSetDate().AsTime(), beforeDate, afterDate)
			}
		})
	}
}

func TestServer_ListGroupMetadata(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	org := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	group := instance.CreateGroup(iamOwnerCtx, t, org.GetOrganizationId(), integration.GroupName())

	setResp, err := instance.Client.GroupV2.SetGroupMetadata(iamOwnerCtx, &group_v2.SetGroupMetadataRequest{
		Id: group.GetId(),
		Metadata: []*group_v2.Metadata{
			{Key: "key1", Value: []byte("value1")},
			{Key: "key2", Value: []byte("value2")},
			{Key: "key2.1", Value: []byte("value3")},
			{Key: "key2.2", Value: []byte("value4")},
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name string
		ctx  context.Context
		req  *group_v2.ListGroupMetadataRequest
		want *group_v2.ListGroupMetadataResponse
	}{
		{
			name: "list all, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.ListGroupMetadataRequest{
				Id: group.GetId(),
			},
			want: &group_v2.ListGroupMetadataResponse{
				Pagination: &filter.PaginationResponse{TotalResult: 4},
				Metadata: []*metadata.Metadata{
					{Key: "key1", Value: []byte("value1"), CreationDate: setResp.GetSetDate(), ChangeDate: setResp.GetSetDate()},
					{Key: "key2", Value: []byte("value2"), CreationDate: setResp.GetSetDate(), ChangeDate: setResp.GetSetDate()},
					{Key: "key2.1", Value: []byte("value3"), CreationDate: setResp.GetSetDate(), ChangeDate: setResp.GetSetDate()},
					{Key: "key2.2", Value: []byte("value4"), CreationDate: setResp.GetSetDate(), ChangeDate: setResp.GetSetDate()},
				},
			},
		},
		{
			name: "filter by key prefix with pagination, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.ListGroupMetadataRequest{
				Id: group.GetId(),
				Pagination: &filter.PaginationRequest{
					Offset: 1,
					Limit:  2,
				},
				Filters: []*metadata.MetadataSearchFilter{
					{
						Filter: &metadata.MetadataSearchFilter_KeyFilter{
							KeyFilter: &metadata.MetadataKeyFilter{
								Key:    "key2",
								Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH,
							},
						},
					},
				},
			},
			want: &group_v2.ListGroupMetadataResponse{
				Pagination: &filter.PaginationResponse{TotalResult: 3, AppliedLimit: 2},
				Metadata: []*metadata.Metadata{
					{Key: "key2.1", Value: []byte("value3"), CreationDate: setResp.GetSetDate(), ChangeDate: setResp.GetSetDate()},
					{Key: "key2.2", Value: []byte("value4"), CreationDate: setResp.GetSetDate(), ChangeDate: setResp.GetSetDate()},
				},
			},
		},
		{
			name: "non-existent group, empty result",
			ctx:  iamOwnerCtx,
			req: &group_v2.ListGroupMetadataRequest{
				Id: "nonexistent-group",
			},
			want: &group_v2.ListGroupMetadataResponse{
				Pagination: &filter.PaginationResponse{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 1*time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := instance.Client.GroupV2.ListGroupMetadata(tt.ctx, tt.req)
				require.NoError(ttt, err)
				assert.EqualExportedValues(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for group metadata projection")
		})
	}
}

func TestServer_DeleteGroupMetadata(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	org := instance.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	group := instance.CreateGroup(iamOwnerCtx, t, org.GetOrganizationId(), integration.GroupName())

	// Seed metadata so Delete has something to remove.
	_, err := instance.Client.GroupV2.SetGroupMetadata(iamOwnerCtx, &group_v2.SetGroupMetadataRequest{
		Id: group.GetId(),
		Metadata: []*group_v2.Metadata{
			{Key: "role", Value: []byte("admin")},
			{Key: "label", Value: []byte("priority")},
			{Key: "region", Value: []byte("eu-west")},
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name             string
		ctx              context.Context
		req              *group_v2.DeleteGroupMetadataRequest
		wantDeletionDate bool
		wantErrCode      codes.Code
	}{
		{
			name: "unauthenticated, error",
			ctx:  context.Background(),
			req: &group_v2.DeleteGroupMetadataRequest{
				Id:   group.GetId(),
				Keys: []string{"role"},
			},
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "missing id, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupMetadataRequest{
				Keys: []string{"role"},
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "empty keys, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupMetadataRequest{
				Id: group.GetId(),
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "group does not exist, error",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupMetadataRequest{
				Id:   "randomGroup",
				Keys: []string{"role"},
			},
			wantErrCode: codes.FailedPrecondition,
		},
		{
			name: "missing permission, error",
			ctx:  instance.WithAuthorizationToken(CTX, integration.UserTypeNoPermission),
			req: &group_v2.DeleteGroupMetadataRequest{
				Id:   group.GetId(),
				Keys: []string{"role"},
			},
			wantErrCode: codes.NotFound,
		},
		{
			name: "delete single key, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupMetadataRequest{
				Id:   group.GetId(),
				Keys: []string{"role"},
			},
			wantDeletionDate: true,
		},
		{
			name: "delete multiple keys, ok",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupMetadataRequest{
				Id:   group.GetId(),
				Keys: []string{"label", "region"},
			},
			wantDeletionDate: true,
		},
		{
			name: "delete unknown key, not found",
			ctx:  iamOwnerCtx,
			req: &group_v2.DeleteGroupMetadataRequest{
				Id:   group.GetId(),
				Keys: []string{"does-not-exist"},
			},
			wantErrCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeDate := time.Now().UTC()
			got, err := instance.Client.GroupV2.DeleteGroupMetadata(tt.ctx, tt.req)
			afterDate := time.Now().UTC()
			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				require.Empty(t, got.GetDeletionDate())
				assert.Equal(t, tt.wantErrCode, status.Code(err))
				return
			}
			require.NoError(t, err)
			if tt.wantDeletionDate {
				assert.WithinRange(t, got.GetDeletionDate().AsTime(), beforeDate, afterDate)
			}
		})
	}
}
