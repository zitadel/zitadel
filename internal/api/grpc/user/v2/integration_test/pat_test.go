//go:build integration

package user_test

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_AddPersonalAccessToken(t *testing.T) {
	resp := Instance.CreateUserTypeMachine(CTX)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
	type args struct {
		ctx     context.Context
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
				CTX,
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
				CTX,
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
				CTX,
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
			name: "add another pat, ok",
			args: args{
				CTX,
				&user.AddPersonalAccessTokenRequest{
					ExpirationDate: expirationDate,
				},
				func(request *user.AddPersonalAccessTokenRequest) error {
					request.UserId = userId
					_, err := Client.AddPersonalAccessToken(CTX, &user.AddPersonalAccessTokenRequest{
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

func TestServer_AddPersonalAccessToken_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, fmt.Sprintf("AddPersonalAccessToken-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
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
	resp := Instance.CreateUserTypeMachine(CTX)
	userId := resp.GetId()
	expirationDate := timestamppb.New(time.Now().Add(time.Hour * 24))
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
			name: "remove pat, user not existing",
			args: args{
				CTX,
				&user.RemovePersonalAccessTokenRequest{
					UserId: "notexisting",
				},
				func(request *user.RemovePersonalAccessTokenRequest) error {
					pat, err := Instance.Client.UserV2.AddPersonalAccessToken(CTX, &user.AddPersonalAccessTokenRequest{
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
				CTX,
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
				CTX,
				&user.RemovePersonalAccessTokenRequest{},
				func(request *user.RemovePersonalAccessTokenRequest) error {
					pat, err := Instance.Client.UserV2.AddPersonalAccessToken(CTX, &user.AddPersonalAccessTokenRequest{
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
			got, err := Client.RemovePersonalAccessToken(tt.args.ctx, tt.args.req)
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
	otherOrg := Instance.CreateOrganization(IamCTX, fmt.Sprintf("RemovePersonalAccessToken-%s", gofakeit.AppName()), gofakeit.Email())
	otherOrgUser, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: gofakeit.Name(),
			},
		},
	})
	request := &user.RemovePersonalAccessTokenRequest{
		UserId: otherOrgUser.GetId(),
	}
	prepare := func(request *user.RemovePersonalAccessTokenRequest) error {
		pat, err := Instance.Client.UserV2.AddPersonalAccessToken(IamCTX, &user.AddPersonalAccessTokenRequest{
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
