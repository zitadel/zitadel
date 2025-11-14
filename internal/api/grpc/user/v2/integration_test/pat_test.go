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

func TestServer_AddPersonalAccessToken(t *testing.T) {
	resp := Instance.CreateUserTypeMachine(IamCTX, Instance.DefaultOrg.Id)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		req     *user.AddPersonalAccessTokenRequest
		prepare func(request *user.AddPersonalAccessTokenRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add pat, user not existing",
			args: args{
				&user.AddPersonalAccessTokenRequest{
					UserId:         "notexisting",
					ExpirationDate: expirationDate,
				},
				func(request *user.AddPersonalAccessTokenRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "add pat, ok",
			args: args{
				&user.AddPersonalAccessTokenRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddPersonalAccessTokenRequest) error {
					request.UserId = userId
					return nil
				},
			},
		},
		{
			name: "add pat human, not ok",
			args: args{
				&user.AddPersonalAccessTokenRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddPersonalAccessTokenRequest) error {
					resp := Instance.CreateUserTypeHuman(IamCTX, integration.Email())
					request.UserId = resp.Id
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "add another pat, ok",
			args: args{
				&user.AddPersonalAccessTokenRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddPersonalAccessTokenRequest) error {
					request.UserId = userId
					_, err := Client.AddPersonalAccessToken(IamCTX, &user.AddPersonalAccessTokenRequest{
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
			got, err := Client.AddPersonalAccessToken(OrgCTX, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.TokenId, "id is empty")
			assert.NotEmpty(t, got.Token, "token is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_AddPersonalAccessToken_Permission(t *testing.T) {
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
	request := &user.AddPersonalAccessTokenRequest{
		ExpirationDate: timestamppb.New(time.Now().Add(time.Hour * 24)),
		UserId:         otherOrgUser.GetId(),
	}
	type args struct {
		ctx context.Context
		req *user.AddPersonalAccessTokenRequest
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
			got, err := Client.AddPersonalAccessToken(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.TokenId, "id is empty")
			assert.NotEmpty(t, got.Token, "token is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_RemovePersonalAccessToken(t *testing.T) {
	resp := Instance.CreateUserTypeMachine(IamCTX, Instance.DefaultOrg.Id)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		req     *user.RemovePersonalAccessTokenRequest
		prepare func(request *user.RemovePersonalAccessTokenRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove pat, user not existing",
			args: args{
				&user.RemovePersonalAccessTokenRequest{
					UserId: "notexisting",
				},
				func(request *user.RemovePersonalAccessTokenRequest) error {
					pat, err := Client.AddPersonalAccessToken(OrgCTX, &user.AddPersonalAccessTokenRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					request.TokenId = pat.GetTokenId()
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "remove pat, not existing",
			args: args{
				&user.RemovePersonalAccessTokenRequest{
					TokenId: "notexisting",
				},
				func(request *user.RemovePersonalAccessTokenRequest) error {
					request.UserId = userId
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "remove pat, ok",
			args: args{
				&user.RemovePersonalAccessTokenRequest{},
				func(request *user.RemovePersonalAccessTokenRequest) error {
					pat, err := Client.AddPersonalAccessToken(OrgCTX, &user.AddPersonalAccessTokenRequest{
						ExpirationDate: expirationDate,
						UserId:         userId,
					})
					request.TokenId = pat.GetTokenId()
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
			got, err := Client.RemovePersonalAccessToken(OrgCTX, tt.args.req)
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

func TestServer_RemovePersonalAccessToken_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Client.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	request := &user.RemovePersonalAccessTokenRequest{
		UserId: otherOrgUser.GetId(),
	}
	prepare := func(request *user.RemovePersonalAccessTokenRequest) error {
		pat, err := Client.AddPersonalAccessToken(IamCTX, &user.AddPersonalAccessTokenRequest{
			ExpirationDate: timestamppb.New(time.Now().Add(time.Hour * 24)),
			UserId:         otherOrgUser.GetId(),
		})
		request.TokenId = pat.GetTokenId()
		return err
	}
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *user.RemovePersonalAccessTokenRequest
		prepare func(request *user.RemovePersonalAccessTokenRequest) error
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
			got, err := Client.RemovePersonalAccessToken(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.DeletionDate, "client pat is empty")
			creationDate := got.DeletionDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_ListPersonalAccessTokens(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.ListPersonalAccessTokensRequest
	}
	type testCase struct {
		name string
		args args
		want *user.ListPersonalAccessTokensResponse
	}
	OrgCTX := OrgCTX
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
	onlySinceTestStartFilter := &user.PersonalAccessTokensSearchFilter{Filter: &user.PersonalAccessTokensSearchFilter_CreatedDateFilter{CreatedDateFilter: &filter.TimestampFilter{
		Timestamp: timestamppb.Now(),
		Method:    filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS,
	}}}
	myOrgId := Instance.DefaultOrg.GetId()
	myUserId := Instance.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Hour).Add(time.Hour * 24)
	myDataPoint := setupPATDataPoint(SystemCTX, t, Instance, myUserId, myOrgId, expiresInADay)
	otherUserDataPoint := setupPATDataPoint(SystemCTX, t, Instance, otherUserId, myOrgId, expiresInADay)
	otherOrgDataPointExpiringSoon := setupPATDataPoint(SystemCTX, t, Instance, otherOrgUserId, otherOrg.OrganizationId, time.Now().Truncate(time.Hour).Add(time.Hour))
	otherOrgDataPointExpiringLate := setupPATDataPoint(SystemCTX, t, Instance, otherOrgUserId, otherOrg.OrganizationId, expiresInADay.Add(time.Hour*24*30))
	sortingColumnExpirationDate := user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_EXPIRATION_DATE
	awaitPersonalAccessTokens(SystemCTX, t, Instance,
		onlySinceTestStartFilter,
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.PersonalAccessTokensSearchFilter_TokenIdFilter{
								TokenIdFilter: &filter.IDFilter{Id: otherOrgDataPointExpiringSoon.Id},
							},
						},
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.PersonalAccessTokensSearchFilter_OrganizationIdFilter{
								OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId},
							},
						}},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Pagination: &filter.PaginationRequest{
						Asc: true,
					},
					SortingColumn: &sortingColumnExpirationDate,
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
						{Filter: &user.PersonalAccessTokensSearchFilter_OrganizationIdFilter{OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId}}},
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Pagination: &filter.PaginationRequest{
						Offset: 2,
						Limit:  2,
						Asc:    true,
					},
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{
						{
							Filter: &user.PersonalAccessTokensSearchFilter_TokenIdFilter{
								TokenIdFilter: &filter.IDFilter{Id: otherUserDataPoint.Id},
							},
						},
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{},
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
				got, err := Client.ListPersonalAccessTokens(tt.args.ctx, tt.args.req)
				require.NoError(ttt, err)
				assert.Len(ttt, got.Result, len(tt.want.Result))
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					ttt.Errorf("ListPersonalAccessTokens() mismatch (-want +got):\n%s", diff)
				}
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func TestServer_ListPersonalAccessTokens_PermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, InstancePermissionV2)
	iamOwnerCtx := InstancePermissionV2.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)

	type args struct {
		ctx context.Context
		req *user.ListPersonalAccessTokensRequest
	}
	type testCase struct {
		name string
		args args
		want *user.ListPersonalAccessTokensResponse
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
	otherUserId := InstancePermissionV2.CreateUserTypeMachine(SystemCTX, InstancePermissionV2.DefaultOrg.Id).GetId()
	onlySinceTestStartFilter := &user.PersonalAccessTokensSearchFilter{Filter: &user.PersonalAccessTokensSearchFilter_CreatedDateFilter{CreatedDateFilter: &filter.TimestampFilter{
		Timestamp: timestamppb.Now(),
		Method:    filter.TimestampFilterMethod_TIMESTAMP_FILTER_METHOD_AFTER_OR_EQUALS,
	}}}
	myOrgId := InstancePermissionV2.DefaultOrg.GetId()
	myUserId := InstancePermissionV2.Users.Get(integration.UserTypeNoPermission).ID
	expiresInADay := time.Now().Truncate(time.Hour).Add(time.Hour * 24)
	myDataPoint := setupPATDataPoint(iamOwnerCtx, t, InstancePermissionV2, myUserId, myOrgId, expiresInADay)
	otherUserDataPoint := setupPATDataPoint(iamOwnerCtx, t, InstancePermissionV2, otherUserId, myOrgId, expiresInADay)
	otherOrgDataPointExpiringSoon := setupPATDataPoint(iamOwnerCtx, t, InstancePermissionV2, otherOrgUserId, otherOrg.OrganizationId, time.Now().Truncate(time.Hour).Add(time.Hour))
	otherOrgDataPointExpiringLate := setupPATDataPoint(iamOwnerCtx, t, InstancePermissionV2, otherOrgUserId, otherOrg.OrganizationId, expiresInADay.Add(time.Hour*24*30))
	sortingColumnExpirationDate := user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_EXPIRATION_DATE
	awaitPersonalAccessTokens(iamOwnerCtx, t, InstancePermissionV2,
		onlySinceTestStartFilter,
		otherOrgDataPointExpiringSoon.GetId(),
		otherOrgDataPointExpiringLate.GetId(),
		otherUserDataPoint.GetId(),
		myDataPoint.GetId(),
	)
	tests := []testCase{
		{
			name: "list all, instance",
			args: args{
				iamOwnerCtx,
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{onlySinceTestStartFilter},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.PersonalAccessTokensSearchFilter_TokenIdFilter{
								TokenIdFilter: &filter.IDFilter{Id: otherOrgDataPointExpiringSoon.Id},
							},
						},
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
						{
							Filter: &user.PersonalAccessTokensSearchFilter_OrganizationIdFilter{
								OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId},
							},
						}},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Pagination: &filter.PaginationRequest{
						Asc: true,
					},
					SortingColumn: &sortingColumnExpirationDate,
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
						{Filter: &user.PersonalAccessTokensSearchFilter_OrganizationIdFilter{OrganizationIdFilter: &filter.IDFilter{Id: otherOrg.OrganizationId}}},
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Pagination: &filter.PaginationRequest{
						Offset: 2,
						Limit:  2,
						Asc:    true,
					},
					Filters: []*user.PersonalAccessTokensSearchFilter{
						onlySinceTestStartFilter,
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{
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
				&user.ListPersonalAccessTokensRequest{
					Filters: []*user.PersonalAccessTokensSearchFilter{
						{
							Filter: &user.PersonalAccessTokensSearchFilter_TokenIdFilter{
								TokenIdFilter: &filter.IDFilter{Id: otherUserDataPoint.Id},
							},
						},
					},
				},
			},
			want: &user.ListPersonalAccessTokensResponse{
				Result: []*user.PersonalAccessToken{},
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
				got, err := InstancePermissionV2.Client.UserV2.ListPersonalAccessTokens(tt.args.ctx, tt.args.req)
				require.NoError(ttt, err)
				assert.Len(ttt, got.Result, len(tt.want.Result))
				if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
					ttt.Errorf("ListPersonalAccessTokens() mismatch (-want +got):\n%s", diff)
				}
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func setupPATDataPoint(ctx context.Context, t *testing.T, instance *integration.Instance, userId, orgId string, expirationDate time.Time) *user.PersonalAccessToken {
	expirationDatePb := timestamppb.New(expirationDate)
	newPersonalAccessToken, err := instance.Client.UserV2.AddPersonalAccessToken(ctx, &user.AddPersonalAccessTokenRequest{
		UserId:         userId,
		ExpirationDate: expirationDatePb,
	})
	require.NoError(t, err)
	return &user.PersonalAccessToken{
		CreationDate:   newPersonalAccessToken.CreationDate,
		ChangeDate:     newPersonalAccessToken.CreationDate,
		Id:             newPersonalAccessToken.GetTokenId(),
		UserId:         userId,
		OrganizationId: orgId,
		ExpirationDate: expirationDatePb,
	}
}

func awaitPersonalAccessTokens(ctx context.Context, t *testing.T, instance *integration.Instance, sinceTestStartFilter *user.PersonalAccessTokensSearchFilter, patIds ...string) {
	sortingColumn := user.PersonalAccessTokenFieldName_PERSONAL_ACCESS_TOKEN_FIELD_NAME_ID
	slices.Sort(patIds)
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		result, err := instance.Client.UserV2.ListPersonalAccessTokens(ctx, &user.ListPersonalAccessTokensRequest{
			Filters:       []*user.PersonalAccessTokensSearchFilter{sinceTestStartFilter},
			SortingColumn: &sortingColumn,
			Pagination: &filter.PaginationRequest{
				Asc: true,
			},
		})
		require.NoError(collect, err)
		if !assert.Len(collect, result.Result, len(patIds)) {
			return
		}
		for i := range patIds {
			patId := patIds[i]
			require.Equal(collect, patId, result.Result[i].GetId())
		}
	}, retryDuration, tick, "pat not created in time")
}
