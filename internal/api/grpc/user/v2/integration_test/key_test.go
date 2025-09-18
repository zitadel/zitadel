//go:build integration

package user_test

import (
	"context"
	"slices"
	"testing"
	"time"

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
	resp := Instance.CreateUserTypeMachine(IamCTX, Instance.DefaultOrg.Id)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		req     *user.AddKeyRequest
		prepare func(request *user.AddKeyRequest) error
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantEmtpyKey bool
	}{
		{
			name: "add key, user not existing",
			args: args{
				&user.AddKeyRequest{
					UserId:         "notexisting",
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "generate key pair, ok",
			args: args{
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
			name: "add valid public key, ok",
			args: args{
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
					// This is the public key of the tester system user. This must be valid.
					PublicKey: []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzi+FFSJL7f5yw4KTwzgM
P34ePGycm/M+kT0M7V4Cgx5V3EaDIvTQKTLfBaEB45zb9LtjIXzDw0rXRoS2hO6t
h+CYQCz3KCvh09C0IzxZiB2IS3H/aT+5Bx9EFY+vnAkZjccbyG5YNRvmtOlnvIeI
H7qZ0tEwkPfF5GEZNPJPtmy3UGV7iofdVQS1xRj73+aMw5rvH4D8IdyiAC3VekIb
pt0Vj0SUX3DwKtog337BzTiPk3aXRF0sbFhQoqdJRI8NqgZjCwjq9yfI5tyxYswn
+JGzHGdHvW3idODlmwEt5K2pasiRIWK2OGfq+w0EcltQHabuqEPgZlmhCkRdNfix
BwIDAQAB
-----END PUBLIC KEY-----
`),
				},
				func(request *user.AddKeyRequest) error {
					request.UserId = userId
					return nil
				},
			},
			wantEmtpyKey: true,
		},
		{
			name: "add invalid public key, error",
			args: args{
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
					PublicKey: []byte(`
-----BEGIN PUBLIC KEY-----
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789
-----END PUBLIC KEY-----
`),
				},
				func(request *user.AddKeyRequest) error {
					request.UserId = userId
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "add key human, error",
			args: args{
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error {
					resp := Instance.CreateUserTypeHuman(IamCTX, integration.Email())
					request.UserId = resp.Id
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "add another key, ok",
			args: args{
				&user.AddKeyRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddKeyRequest) error {
					request.UserId = userId
					_, err := Client.AddKey(IamCTX, &user.AddKeyRequest{
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
			got, err := Client.AddKey(OrgCTX, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.KeyId, "key id is empty")
			if tt.wantEmtpyKey {
				assert.Empty(t, got.KeyContent, "key content is not empty")
			} else {
				assert.NotEmpty(t, got.KeyContent, "key content is empty")
			}
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_AddKey_Permission(t *testing.T) {
	OrgCTX := OrgCTX
	otherOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Client.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
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
			args:    args{OrgCTX, request},
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
	resp := Instance.CreateUserTypeMachine(IamCTX, Instance.DefaultOrg.Id)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
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
				&user.RemoveKeyRequest{
					UserId: "notexisting",
				},
				func(request *user.RemoveKeyRequest) error {
					key, err := Client.AddKey(IamCTX, &user.AddKeyRequest{
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
				&user.RemoveKeyRequest{},
				func(request *user.RemoveKeyRequest) error {
					key, err := Client.AddKey(IamCTX, &user.AddKeyRequest{
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
			got, err := Client.RemoveKey(OrgCTX, tt.args.req)
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
	OrgCTX := OrgCTX
	otherOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Client.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	request := &user.RemoveKeyRequest{
		UserId: otherOrgUser.GetId(),
	}
	prepare := func(request *user.RemoveKeyRequest) error {
		key, err := Client.AddKey(IamCTX, &user.AddKeyRequest{
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
			args:    args{OrgCTX, request, prepare},
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
	type args struct {
		ctx context.Context
		req *user.ListKeysRequest
	}
	type testCase struct {
		name string
		args args
		want *user.ListKeysResponse
	}
	OrgCTX := OrgCTX
	otherOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Client.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	require.NoError(t, err)
	otherOrgUserId := otherOrgUser.GetId()
	otherUserId := Instance.CreateUserTypeMachine(IamCTX, Instance.DefaultOrg.Id).GetId()
	onlySinceTestStartFilter := &user.KeysSearchFilter{Filter: &user.KeysSearchFilter_CreatedDateFilter{CreatedDateFilter: &filter.TimestampFilter{
		Timestamp: timestamppb.Now(),
		Method:    filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS,
	}}}
	myOrgId := Instance.DefaultOrg.GetId()
	myUserId := Instance.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Hour).Add(time.Hour * 24)
	myDataPoint := setupKeyDataPoint(IamCTX, t, Instance, myUserId, myOrgId, expiresInADay)
	otherUserDataPoint := setupKeyDataPoint(IamCTX, t, Instance, otherUserId, myOrgId, expiresInADay)
	otherOrgDataPointExpiringSoon := setupKeyDataPoint(IamCTX, t, Instance, otherOrgUserId, otherOrg.OrganizationId, time.Now().Truncate(time.Hour).Add(time.Hour))
	otherOrgDataPointExpiringLate := setupKeyDataPoint(IamCTX, t, Instance, otherOrgUserId, otherOrg.OrganizationId, expiresInADay.Add(time.Hour*24*30))
	sortingColumnExpirationDate := user.KeyFieldName_KEY_FIELD_NAME_KEY_EXPIRATION_DATE
	awaitKeys(IamCTX, t, Instance, onlySinceTestStartFilter,
		otherOrgDataPointExpiringSoon.GetId(),
		otherOrgDataPointExpiringLate.GetId(),
		otherUserDataPoint.GetId(),
		myDataPoint.GetId(),
	)
	tests := []testCase{
		{
			name: "list all, instance",
			args: args{
				IamCTX,
				&user.ListKeysRequest{Filters: []*user.KeysSearchFilter{onlySinceTestStartFilter}},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringLate,
					otherOrgDataPointExpiringSoon,
					otherUserDataPoint,
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, org",
			args: args{
				OrgCTX,
				&user.ListKeysRequest{Filters: []*user.KeysSearchFilter{onlySinceTestStartFilter}},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherUserDataPoint,
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, user",
			args: args{
				UserCTX,
				&user.ListKeysRequest{Filters: []*user.KeysSearchFilter{onlySinceTestStartFilter}},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list by id",
			args: args{
				IamCTX,
				&user.ListKeysRequest{
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.KeysSearchFilter_KeyIdFilter{
								KeyIdFilter: &filter.IDFilter{Id: otherOrgDataPointExpiringSoon.Id},
							},
						},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringSoon,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all from other org",
			args: args{
				IamCTX,
				&user.ListKeysRequest{
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.KeysSearchFilter_OrganizationIdFilter{
								OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId},
							},
						},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringLate,
					otherOrgDataPointExpiringSoon,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "sort by next expiration dates",
			args: args{
				IamCTX,
				&user.ListKeysRequest{
					Pagination: &filter.PaginationRequest{
						Asc: true,
					},
					SortingColumn: &sortingColumnExpirationDate,
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
						{Filter: &user.KeysSearchFilter_OrganizationIdFilter{OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId}}},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringSoon,
					otherOrgDataPointExpiringLate,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "get page",
			args: args{
				IamCTX,
				&user.ListKeysRequest{
					Pagination: &filter.PaginationRequest{
						Offset: 2,
						Limit:  2,
						Asc:    true,
					},
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringSoon,
					otherOrgDataPointExpiringLate,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 2,
				},
			},
		},
		{
			name: "empty list",
			args: args{
				UserCTX,
				&user.ListKeysRequest{
					Filters: []*user.KeysSearchFilter{
						{
							Filter: &user.KeysSearchFilter_KeyIdFilter{
								KeyIdFilter: &filter.IDFilter{Id: otherUserDataPoint.Id},
							},
						},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{},
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, 20*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListKeys(tt.args.ctx, tt.args.req)
				require.NoError(ttt, err)
				if !assert.Len(ttt, got.Result, len(tt.want.Result)) {
					return
				}
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					ttt.Errorf("ListKeys() mismatch (-want +got):\n%s", diff)
				}
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func TestServer_ListKeys_PermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, InstancePermissionV2)
	iamOwnerCtx := InstancePermissionV2.WithAuthorizationToken(OrgCTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		req *user.ListKeysRequest
	}
	type testCase struct {
		name string
		args args
		want *user.ListKeysResponse
	}
	otherOrg := InstancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := InstancePermissionV2.Client.UserV2.CreateUser(iamOwnerCtx, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	require.NoError(t, err)
	otherOrgUserId := otherOrgUser.GetId()
	otherUserId := InstancePermissionV2.CreateUserTypeMachine(iamOwnerCtx, InstancePermissionV2.DefaultOrg.Id).GetId()
	onlySinceTestStartFilter := &user.KeysSearchFilter{Filter: &user.KeysSearchFilter_CreatedDateFilter{CreatedDateFilter: &filter.TimestampFilter{
		Timestamp: timestamppb.Now(),
		Method:    filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS,
	}}}
	myOrgId := InstancePermissionV2.DefaultOrg.GetId()
	myUserId := InstancePermissionV2.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Hour).Add(time.Hour * 24)
	myDataPoint := setupKeyDataPoint(iamOwnerCtx, t, InstancePermissionV2, myUserId, myOrgId, expiresInADay)
	otherUserDataPoint := setupKeyDataPoint(iamOwnerCtx, t, InstancePermissionV2, otherUserId, myOrgId, expiresInADay)
	otherOrgDataPointExpiringSoon := setupKeyDataPoint(iamOwnerCtx, t, InstancePermissionV2, otherOrgUserId, otherOrg.OrganizationId, time.Now().Truncate(time.Hour).Add(time.Hour))
	otherOrgDataPointExpiringLate := setupKeyDataPoint(iamOwnerCtx, t, InstancePermissionV2, otherOrgUserId, otherOrg.OrganizationId, expiresInADay.Add(time.Hour*24*30))
	sortingColumnExpirationDate := user.KeyFieldName_KEY_FIELD_NAME_KEY_EXPIRATION_DATE
	awaitKeys(iamOwnerCtx, t, InstancePermissionV2, onlySinceTestStartFilter,
		otherOrgDataPointExpiringSoon.GetId(),
		otherOrgDataPointExpiringLate.GetId(),
		otherUserDataPoint.GetId(),
		myDataPoint.GetId(),
	)
	tests := []testCase{
		{
			name: "list all, InstancePermissionV2",
			args: args{
				iamOwnerCtx,
				&user.ListKeysRequest{Filters: []*user.KeysSearchFilter{onlySinceTestStartFilter}},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringLate,
					otherOrgDataPointExpiringSoon,
					otherUserDataPoint,
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, org",
			args: args{
				InstancePermissionV2.WithAuthorizationToken(iamOwnerCtx, integration.UserTypeOrgOwner),
				&user.ListKeysRequest{Filters: []*user.KeysSearchFilter{onlySinceTestStartFilter}},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherUserDataPoint,
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all, user",
			args: args{
				InstancePermissionV2.WithAuthorizationToken(iamOwnerCtx, integration.UserTypeNoPermission),
				&user.ListKeysRequest{Filters: []*user.KeysSearchFilter{onlySinceTestStartFilter}},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					myDataPoint,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list by id",
			args: args{
				iamOwnerCtx,
				&user.ListKeysRequest{
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.KeysSearchFilter_KeyIdFilter{
								KeyIdFilter: &filter.IDFilter{Id: otherOrgDataPointExpiringSoon.Id},
							},
						},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringSoon,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  1,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "list all from other org",
			args: args{
				iamOwnerCtx,
				&user.ListKeysRequest{
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.KeysSearchFilter_OrganizationIdFilter{
								OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId},
							},
						},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringLate,
					otherOrgDataPointExpiringSoon,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "sort by next expiration dates",
			args: args{
				iamOwnerCtx,
				&user.ListKeysRequest{
					Pagination: &filter.PaginationRequest{
						Asc: true,
					},
					SortingColumn: &sortingColumnExpirationDate,
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
						{Filter: &user.KeysSearchFilter_OrganizationIdFilter{OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId}}},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringSoon,
					otherOrgDataPointExpiringLate,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  2,
					AppliedLimit: 100,
				},
			},
		},
		{
			name: "get page",
			args: args{
				iamOwnerCtx,
				&user.ListKeysRequest{
					Pagination: &filter.PaginationRequest{
						Offset: 2,
						Limit:  2,
						Asc:    true,
					},
					Filters: []*user.KeysSearchFilter{
						onlySinceTestStartFilter,
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{
					otherOrgDataPointExpiringSoon,
					otherOrgDataPointExpiringLate,
				},
				Pagination: &filter.PaginationResponse{
					TotalResult:  4,
					AppliedLimit: 2,
				},
			},
		},
		{
			name: "empty list",
			args: args{
				InstancePermissionV2.WithAuthorizationToken(iamOwnerCtx, integration.UserTypeNoPermission),
				&user.ListKeysRequest{
					Filters: []*user.KeysSearchFilter{
						{
							Filter: &user.KeysSearchFilter_KeyIdFilter{
								KeyIdFilter: &filter.IDFilter{Id: otherUserDataPoint.Id},
							},
						},
					},
				},
			},
			want: &user.ListKeysResponse{
				Result: []*user.Key{},
				Pagination: &filter.PaginationResponse{
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, 20*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := InstancePermissionV2.Client.UserV2.ListKeys(tt.args.ctx, tt.args.req)
				require.NoError(ttt, err)
				assert.Len(ttt, got.Result, len(tt.want.Result))
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					ttt.Errorf("ListKeys() mismatch (-want +got):\n%s", diff)
				}
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func setupKeyDataPoint(ctx context.Context, t *testing.T, instance *integration.Instance, userId, orgId string, expirationDate time.Time) *user.Key {
	expirationDatePb := timestamppb.New(expirationDate)
	newKey, err := instance.Client.UserV2.AddKey(ctx, &user.AddKeyRequest{
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

func awaitKeys(ctx context.Context, t *testing.T, instance *integration.Instance, sinceTestStartFilter *user.KeysSearchFilter, keyIds ...string) {
	sortingColumn := user.KeyFieldName_KEY_FIELD_NAME_ID
	slices.Sort(keyIds)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		result, err := instance.Client.UserV2.ListKeys(ctx, &user.ListKeysRequest{
			Filters:       []*user.KeysSearchFilter{sinceTestStartFilter},
			SortingColumn: &sortingColumn,
			Pagination: &filter.PaginationRequest{
				Asc: true,
			},
		})
		require.NoError(collect, err)
		if !assert.Len(collect, result.Result, len(keyIds)) {
			return
		}
		for i := range keyIds {
			keyId := keyIds[i]
			require.Equal(collect, keyId, result.Result[i].GetId())
		}
	}, retryDuration, tick, "key not created in time")
}
