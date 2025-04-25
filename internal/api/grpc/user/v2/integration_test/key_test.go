//go:build integration

package user_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_AddKey(t *testing.T) {
	resp := Instance.CreateUserTypeMachine(CTX)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		ctx     context.Context
		req     *user.AddKeyRequest
		prepare func(request *user.AddKeyRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add key, user not existing",
			args: args{
				CTX,
				&user.AddKeyRequest{
					UserId:         "notexisting",
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "add key, ok",
			args: args{
				CTX,
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error {
					request.UserId = userId
					return nil
				},
			},
		},
		{
			name: "add key human, not ok",
			args: args{
				CTX,
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error {
					request.UserId = userId
					return nil
				},
			},
		},
		{
			name: "add another key, ok",
			args: args{
				CTX,
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error {
					request.UserId = userId
					_, err := Client.AddKey(CTX, &user.AddKeyRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					return err
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)
			got, err := Client.AddKey(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.KeyId, "key id is empty")
			assert.NotEmpty(t, got.KeyContent, "key content is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_AddKey_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, fmt.Sprintf("AddKey-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
			},
		},
	})
	require.NoError(t, err)
	request := &user.AddKeyRequest{
		ExpirationDate: timestamppb.New(time.Now().Add(time.Hour * 24)),
		UserId:         otherOrgUser.GetId(),
	}
	type args struct {
		ctx context.Context
		req *user.AddKeyRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "system, ok",
			args: args{SystemCTX, request},
		},
		{
			name: "instance, ok",
			args: args{IamCTX, request},
		},
		{
			name:    "org, error",
			args:    args{CTX, request},
			wantErr: true,
		},
		{
			name:    "user, error",
			args:    args{UserCTX, request},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			require.NoError(t, err)
			got, err := Client.AddKey(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.KeyId, "key id is empty")
			assert.NotEmpty(t, got.KeyContent, "key content is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_RemoveKey(t *testing.T) {
	resp := Instance.CreateUserTypeMachine(CTX)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		ctx     context.Context
		req     *user.RemoveKeyRequest
		prepare func(request *user.RemoveKeyRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove key, user not existing",
			args: args{
				CTX,
				&user.RemoveKeyRequest{
					UserId: "notexisting",
				},
				func(request *user.RemoveKeyRequest) error {
					key, err := Instance.Client.UserV2.AddKey(CTX, &user.AddKeyRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					request.KeyId = key.GetKeyId()
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "remove key, not existing",
			args: args{
				CTX,
				&user.RemoveKeyRequest{
					KeyId: "notexisting",
				},
				func(request *user.RemoveKeyRequest) error {
					request.UserId = userId
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "remove key, ok",
			args: args{
				CTX,
				&user.RemoveKeyRequest{},
				func(request *user.RemoveKeyRequest) error {
					key, err := Instance.Client.UserV2.AddKey(CTX, &user.AddKeyRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					request.KeyId = key.GetKeyId()
					request.UserId = userId
					return err
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)
			got, err := Client.RemoveKey(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			deletionDate := got.DeletionDate.AsTime()
			assert.Greater(t, deletionDate, now, "creation date is before the test started")
			assert.Less(t, deletionDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_RemoveKey_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, fmt.Sprintf("RemoveKey-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
			},
		},
	})
	request := &user.RemoveKeyRequest{
		UserId: otherOrgUser.GetId(),
	}
	prepare := func(request *user.RemoveKeyRequest) error {
		key, err := Instance.Client.UserV2.AddKey(IamCTX, &user.AddKeyRequest{
			ExpirationDate: timestamppb.New(time.Now().Add(time.Hour * 24)),
			UserId:         otherOrgUser.GetId(),
		})
		request.KeyId = key.GetKeyId()
		return err
	}
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *user.RemoveKeyRequest
		prepare func(request *user.RemoveKeyRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "system, ok",
			args: args{SystemCTX, request, prepare},
		},
		{
			name: "instance, ok",
			args: args{IamCTX, request, prepare},
		},
		{
			name:    "org, error",
			args:    args{CTX, request, prepare},
			wantErr: true,
		},
		{
			name:    "user, error",
			args:    args{UserCTX, request, prepare},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			require.NoError(t, tt.args.prepare(tt.args.req))
			got, err := Client.RemoveKey(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.DeletionDate, "client key is empty")
			creationDate := got.DeletionDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_ListKeys(t *testing.T) {
	/*	otherInstance := integration.NewInstance(SystemCTX)
		otherInstanceUserId := otherInstance.CreateUserTypeMachine(SystemCTX).GetId()*/
	otherOrg := Instance.CreateOrganization(SystemCTX, fmt.Sprintf("ListKeys-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Client.CreateUser(SystemCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
			},
		},
	})
	require.NoError(t, err)
	otherOrgUserId := otherOrgUser.GetId()
	otherUserId := Instance.CreateUserTypeMachine(SystemCTX).GetId()
	myOrgId := Instance.DefaultOrg.GetId()
	myUserId := Instance.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Second).Add(time.Hour * 24)
	//	otherInstanceDataPoint := setupDataPoint(t, otherInstance.Client.UserV2, otherInstanceUserId, otherInstance.DefaultOrg.GetId(), expiresInADay)
	otherOrgDataPoint := setupDataPoint(t, otherOrgUserId, otherOrg.OrganizationId, expiresInADay)
	otherUserDataPoint := setupDataPoint(t, otherUserId, myOrgId, expiresInADay)
	myDataPoint := setupDataPoint(t, myUserId, myOrgId, expiresInADay)
	awaitKeys(t, otherOrgDataPoint.GetId(), otherUserDataPoint.GetId(), myDataPoint.GetId())
	type args struct {
		ctx context.Context
		req *user.ListKeysRequest
	}
	tests := []struct {
		name string
		args args
		want *user.ListKeysResponse
	}{
		{
			name: "list all, instance",
			args: args{
				IamCTX,
				&user.ListKeysRequest{},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					//					otherInstanceDataPoint,
					myDataPoint,
					otherUserDataPoint,
					otherOrgDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  3,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, org",
			args: args{
				CTX,
				&user.ListKeysRequest{},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					//					otherInstanceDataPoint,
					myDataPoint,
					otherUserDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		/*
			{
				name: "list all, user",
			},
			{
				name: "list by id",
			},
			{
				name: "list by multiple ids",
			},
			{
				name: "list all from other instance",
			},
			{
				name: "list all from other instance and org",
			},
			{
				name: "sort by descending expiration date",
			},
			{
				name: "get page",
			},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.ListKeys(tt.args.ctx, tt.args.req)
			require.NoError(t, err)
			assert.Len(t, got.Result, len(tt.want.Result))
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("ListKeys() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func setupDataPoint(t *testing.T, userId, orgId string, expirationDate time.Time) *user.Key {
	expirationDatePb := timestamppb.New(expirationDate)
	newKey, err := Client.AddKey(SystemCTX, &user.AddKeyRequest{
		UserId:         userId,
		ExpirationDate: expirationDatePb,
		PublicKey:      nil,
	})
	require.NoError(t, err)
	return &user.Key{
		CreationDate:   newKey.CreationDate,
		ChangeDate:     newKey.CreationDate,
		Id:             newKey.GetKeyId(),
		UserId:         userId,
		OrganizationId: orgId,
		ExpirationDate: expirationDatePb,
	}
}

func awaitKeys(t *testing.T, keyIds ...string) {
	sortingColumn := user.KeyFieldName_KEY_FIELD_NAME_ID
	slices.Sort(keyIds)
	var filters []*user.KeysSearchFilter
	for _, keyId := range keyIds {
		filters = append(filters, &user.KeysSearchFilter{
			Filter: &user.KeysSearchFilter_KeyIdFilter{
				KeyIdFilter: &user.IDFilter{Id: keyId},
			},
		})
	}
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		result, err := Client.ListKeys(SystemCTX, &user.ListKeysRequest{
			Filters: []*user.KeysSearchFilter{{
				Filter: &user.KeysSearchFilter_OrFilter{
					OrFilter: &user.KeysOrFilter{Filters: filters},
				},
			}},
			SortingColumn: &sortingColumn,
			Pagination: &filter.PaginationRequest{
				Asc: true,
			},
		})
		require.NoError(t, err)
		if !assert.Len(collect, result.Result, len(keyIds)) {
			return
		}
		for i := range keyIds {
			keyId := keyIds[i]
			require.Equal(collect, keyId, result.Result[i].GetId())
		}
	}, 5*time.Second, time.Second, "key not created in time")
}
