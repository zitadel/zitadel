//go:build integration

package user_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_SetUserMetadata(t *testing.T) {
	iamOwnerCTX := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(request *user.SetUserMetadataRequest)
		req     *user.SetUserMetadataRequest
		setDate bool
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Instance.WithAuthorizationToken(context.Background(), integration.UserTypeNoPermission),
			dep: func(req *user.SetUserMetadataRequest) {
				req.UserId = Instance.CreateUserTypeHuman(CTX, integration.Email()).GetId()
			},
			req: &user.SetUserMetadataRequest{
				Metadata: []*user.Metadata{{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value1")))}},
			},
			wantErr: true,
		},
		{
			name: "set user metadata",
			ctx:  iamOwnerCTX,
			dep: func(req *user.SetUserMetadataRequest) {
				req.UserId = Instance.CreateUserTypeHuman(CTX, integration.Email()).GetId()
			},
			req: &user.SetUserMetadataRequest{
				Metadata: []*user.Metadata{{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value1")))}},
			},
			setDate: true,
		},
		{
			name: "set user metadata, multiple",
			ctx:  iamOwnerCTX,
			dep: func(req *user.SetUserMetadataRequest) {
				req.UserId = Instance.CreateUserTypeHuman(CTX, integration.Email()).GetId()
			},
			req: &user.SetUserMetadataRequest{
				Metadata: []*user.Metadata{
					{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value1")))},
					{Key: "key2", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value2")))},
					{Key: "key3", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value3")))},
				},
			},
			setDate: true,
		},
		{
			name: "set user metadata on non existent user",
			ctx:  iamOwnerCTX,
			req: &user.SetUserMetadataRequest{
				UserId:   "notexisting",
				Metadata: []*user.Metadata{{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value1")))}},
			},
			wantErr: true,
		},
		{
			name: "update user metadata",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			dep: func(req *user.SetUserMetadataRequest) {
				req.UserId = Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
				Instance.SetUserMetadata(iamOwnerCTX, req.UserId, "key1", "value1")
			},
			req: &user.SetUserMetadataRequest{
				Metadata: []*user.Metadata{{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value2")))}},
			},
			setDate: true,
		},
		{
			name: "update user metadata with same value",
			ctx:  Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner),
			dep: func(req *user.SetUserMetadataRequest) {
				req.UserId = Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
				Instance.SetUserMetadata(iamOwnerCTX, req.UserId, "key1", "value1")
			},
			req: &user.SetUserMetadataRequest{
				Metadata: []*user.Metadata{{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value1")))}},
			},
			setDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			if tt.dep != nil {
				tt.dep(tt.req)
			}
			got, err := Client.SetUserMetadata(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertSetUserMetadataResponse(t, creationDate, changeDate, tt.setDate, got)
		})
	}
}

func assertSetUserMetadataResponse(t *testing.T, creationDate, changeDate time.Time, expectedSetDat bool, actualResp *user.SetUserMetadataResponse) {
	if expectedSetDat {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetSetDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetSetDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.SetDate)
	}
}

func TestServer_ListUserMetadata(t *testing.T) {
	iamOwnerCTX := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		dep func(context.Context, *user.ListUserMetadataRequest, *user.ListUserMetadataResponse)
		req *user.ListUserMetadataRequest
	}

	tests := []struct {
		name    string
		args    args
		want    *user.ListUserMetadataResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			args: args{
				ctx: Instance.WithAuthorizationToken(context.Background(), integration.UserTypeNoPermission),
				dep: func(ctx context.Context, request *user.ListUserMetadataRequest, response *user.ListUserMetadataResponse) {
					userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
					request.UserId = userID
					Instance.SetUserMetadata(iamOwnerCTX, userID, "key1", "value1")
				},
				req: &user.ListUserMetadataRequest{},
			},
			wantErr: true,
		},
		{
			name: "list request",
			args: args{
				ctx: iamOwnerCTX,
				dep: func(ctx context.Context, request *user.ListUserMetadataRequest, response *user.ListUserMetadataResponse) {
					userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
					request.UserId = userID
					metadataResp := Instance.SetUserMetadata(iamOwnerCTX, userID, "key1", "value1")

					response.Metadata[0] = &metadata.Metadata{
						CreationDate: metadataResp.GetSetDate(),
						ChangeDate:   metadataResp.GetSetDate(),
						Key:          "key1",
						Value:        []byte(base64.StdEncoding.EncodeToString([]byte("value1"))),
					}
				},
				req: &user.ListUserMetadataRequest{},
			},
			want: &user.ListUserMetadataResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Metadata: []*metadata.Metadata{
					{},
				},
			},
		},
		{
			name: "list request single key",
			args: args{
				ctx: iamOwnerCTX,
				dep: func(ctx context.Context, request *user.ListUserMetadataRequest, response *user.ListUserMetadataResponse) {
					userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
					request.UserId = userID
					key := "key1"
					response.Metadata[0] = setUserMetadata(iamOwnerCTX, userID, key, "value1")
					Instance.SetUserMetadata(iamOwnerCTX, userID, "key2", "value2")
					Instance.SetUserMetadata(iamOwnerCTX, userID, "key3", "value3")
					request.Filters[0] = &metadata.MetadataSearchFilter{
						Filter: &metadata.MetadataSearchFilter_KeyFilter{KeyFilter: &metadata.MetadataKeyFilter{Key: key}},
					}
				},
				req: &user.ListUserMetadataRequest{
					Filters: []*metadata.MetadataSearchFilter{{}},
				},
			},
			want: &user.ListUserMetadataResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
				Metadata: []*metadata.Metadata{
					{},
				},
			},
		},
		{
			name: "list multiple keys",
			args: args{
				ctx: iamOwnerCTX,
				dep: func(ctx context.Context, request *user.ListUserMetadataRequest, response *user.ListUserMetadataResponse) {
					userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
					request.UserId = userID

					response.Metadata[2] = setUserMetadata(iamOwnerCTX, userID, "key1", "value1")
					response.Metadata[1] = setUserMetadata(iamOwnerCTX, userID, "key2", "value2")
					response.Metadata[0] = setUserMetadata(iamOwnerCTX, userID, "key3", "value3")
				},
				req: &user.ListUserMetadataRequest{},
			},
			want: &user.ListUserMetadataResponse{
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
				Metadata: []*metadata.Metadata{
					{}, {}, {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.dep != nil {
				tt.args.dep(tt.args.ctx, tt.args.req, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, listErr := Instance.Client.UserV2.ListUserMetadata(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, listErr)
					return
				}
				require.NoError(ttt, listErr)
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Metadata, len(tt.want.Metadata)) {
					assert.EqualExportedValues(ttt, got.Metadata, tt.want.Metadata)
				}
				assertPaginationResponse(ttt, tt.want.Pagination, got.Pagination)
			}, retryDuration, tick, "timeout waiting for expected execution result")
		})
	}
}

func setUserMetadata(ctx context.Context, userID, key, value string) *metadata.Metadata {
	metadataResp := Instance.SetUserMetadata(ctx, userID, key, value)
	return &metadata.Metadata{
		CreationDate: metadataResp.GetSetDate(),
		ChangeDate:   metadataResp.GetSetDate(),
		Key:          key,
		Value:        []byte(base64.StdEncoding.EncodeToString([]byte(value))),
	}
}

func assertPaginationResponse(t *assert.CollectT, expected *filter.PaginationResponse, actual *filter.PaginationResponse) {
	assert.Equal(t, expected.AppliedLimit, actual.AppliedLimit)
	assert.Equal(t, expected.TotalResult, actual.TotalResult)
}

func TestServer_DeleteUserMetadata(t *testing.T) {
	iamOwnerCTX := Instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *user.DeleteUserMetadataRequest) (time.Time, time.Time)
		req              *user.DeleteUserMetadataRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "empty id",
			ctx:  iamOwnerCTX,
			req: &user.DeleteUserMetadataRequest{
				UserId: "",
			},
			wantErr: true,
		},
		{
			name: "delete, user not existing",
			ctx:  iamOwnerCTX,
			req: &user.DeleteUserMetadataRequest{
				UserId: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "delete",
			ctx:  iamOwnerCTX,
			prepare: func(request *user.DeleteUserMetadataRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
				request.UserId = userID
				key := "key1"
				Instance.SetUserMetadata(iamOwnerCTX, userID, key, "value1")
				request.Keys = []string{key}
				return creationDate, time.Time{}
			},
			req:              &user.DeleteUserMetadataRequest{},
			wantDeletionDate: true,
		},
		{
			name: "delete, empty list",
			ctx:  iamOwnerCTX,
			prepare: func(request *user.DeleteUserMetadataRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
				request.UserId = userID
				key := "key1"
				Instance.SetUserMetadata(iamOwnerCTX, userID, key, "value1")
				Instance.DeleteUserMetadata(iamOwnerCTX, userID, key)
				return creationDate, time.Now().UTC()
			},
			req:     &user.DeleteUserMetadataRequest{},
			wantErr: true,
		},
		{
			name: "delete, already removed",
			ctx:  iamOwnerCTX,
			prepare: func(request *user.DeleteUserMetadataRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
				request.UserId = userID
				key := "key1"
				Instance.SetUserMetadata(iamOwnerCTX, userID, key, "value1")
				Instance.DeleteUserMetadata(iamOwnerCTX, userID, key)
				request.Keys = []string{key}
				return creationDate, time.Now().UTC()
			},
			req:     &user.DeleteUserMetadataRequest{},
			wantErr: true,
		},
		{
			name: "delete, multiple",
			ctx:  iamOwnerCTX,
			prepare: func(request *user.DeleteUserMetadataRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				userID := Instance.CreateUserTypeHuman(iamOwnerCTX, integration.Email()).GetId()
				request.UserId = userID
				key1 := "key1"
				Instance.SetUserMetadata(iamOwnerCTX, userID, key1, "value1")
				key2 := "key2"
				Instance.SetUserMetadata(iamOwnerCTX, userID, key2, "value1")
				key3 := "key3"
				Instance.SetUserMetadata(iamOwnerCTX, userID, key3, "value1")
				request.Keys = []string{key1, key2, key3}
				return creationDate, time.Time{}
			},
			req:              &user.DeleteUserMetadataRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := Instance.Client.UserV2.DeleteUserMetadata(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteProjectResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func assertDeleteProjectResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *user.DeleteUserMetadataResponse) {
	if expectedDeletionDate {
		if !deletionDate.IsZero() {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, deletionDate)
		} else {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.DeletionDate)
	}
}
