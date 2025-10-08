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
			got, err := Client.AddKey(CTX, tt.args.req)
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
	OrgCTX := CTX
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
			got, err := Client.RemoveKey(CTX, tt.args.req)
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
	OrgCTX := CTX
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
	OrgCTX := CTX
	otherOrg := Instance.CreateOrganization(SystemCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Client.CreateUser(SystemCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	require.NoError(t, err)
	otherOrgUserId := otherOrgUser.GetId()
	otherUserId := Instance.CreateUserTypeMachine(SystemCTX, Instance.DefaultOrg.Id).GetId()
	onlySinceTestStartFilter := &user.KeysSearchFilter{Filter: &user.KeysSearchFilter_CreatedDateFilter{CreatedDateFilter: &filter.TimestampFilter{
		Timestamp: timestamppb.Now(),
		Method:    filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS,
	}}}
	myOrgId := Instance.DefaultOrg.GetId()
	myUserId := Instance.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Hour).Add(time.Hour * 24)
	myDataPoint := setupKeyDataPoint(t, myUserId, myOrgId, expiresInADay)
	otherUserDataPoint := setupKeyDataPoint(t, otherUserId, myOrgId, expiresInADay)
	otherOrgDataPointExpiringSoon := setupKeyDataPoint(t, otherOrgUserId, otherOrg.OrganizationId, time.Now().Truncate(time.Hour).Add(time.Hour))
	otherOrgDataPointExpiringLate := setupKeyDataPoint(t, otherOrgUserId, otherOrg.OrganizationId, expiresInADay.Add(time.Hour*24*30))
	sortingColumnExpirationDate := user.KeyFieldName_KEY_FIELD_NAME_KEY_EXPIRATION_DATE
	awaitKeys(t, onlySinceTestStartFilter,
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
					TotalResult:  2,
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
					TotalResult:  1,
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
					TotalResult:  0,
					AppliedLimit: 100,
				},
			},
		},
	}
	t.Run("with permission flag v2", func(t *testing.T) {
		setPermissionCheckV2Flag(t, true)
		defer setPermissionCheckV2Flag(t, false)
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Client.ListKeys(tt.args.ctx, tt.args.req)
				require.NoError(t, err)
				assert.Len(t, got.Result, len(tt.want.Result))
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					t.Errorf("ListKeys() mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})
	t.Run("without permission flag v2", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := Client.ListKeys(tt.args.ctx, tt.args.req)
				require.NoError(t, err)
				assert.Len(t, got.Result, len(tt.want.Result))
				// ignore the total result, as this is a known bug with the in-memory permission checks.
				// The command can't know how many keys exist in the system if the SQL statement has a limit.
				// This is fixed, once the in-memory permission checks are removed with https://github.com/zitadel/zitadel/issues/9188
				tt.want.Pagination.TotalResult = got.Pagination.TotalResult
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					t.Errorf("ListKeys() mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})
}

func setupKeyDataPoint(t *testing.T, userId, orgId string, expirationDate time.Time) *user.Key {
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

func awaitKeys(t *testing.T, sinceTestStartFilter *user.KeysSearchFilter, keyIds ...string) {
	sortingColumn := user.KeyFieldName_KEY_FIELD_NAME_ID
	slices.Sort(keyIds)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		result, err := Client.ListKeys(SystemCTX, &user.ListKeysRequest{
			Filters:       []*user.KeysSearchFilter{sinceTestStartFilter},
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
