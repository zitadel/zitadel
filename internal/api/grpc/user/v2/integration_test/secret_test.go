//go:build integration

package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_AddSecret(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.AddSecretRequest
		prepare func(request *user.AddSecretRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "add secret, user not existing",
			args: args{
				CTX,
				&user.AddSecretRequest{
					UserId: "notexisting",
				},
				func(request *user.AddSecretRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "add secret, ok",
			args: args{
				CTX,
				&user.AddSecretRequest{},
				func(request *user.AddSecretRequest) error {
					resp := Instance.CreateUserTypeMachine(CTX, Instance.DefaultOrg.Id)
					request.UserId = resp.GetId()
					return nil
				},
			},
		},
		{
			name: "add secret human, not ok",
			args: args{
				CTX,
				&user.AddSecretRequest{},
				func(request *user.AddSecretRequest) error {
					resp := Instance.CreateUserTypeMachine(CTX, Instance.DefaultOrg.Id)
					request.UserId = resp.GetId()
					return nil
				},
			},
		},
		{
			name: "overwrite secret, ok",
			args: args{
				CTX,
				&user.AddSecretRequest{},
				func(request *user.AddSecretRequest) error {
					resp := Instance.CreateUserTypeMachine(CTX, Instance.DefaultOrg.Id)
					request.UserId = resp.GetId()
					_, err := Client.AddSecret(CTX, &user.AddSecretRequest{
						UserId: resp.GetId(),
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
			got, err := Client.AddSecret(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.ClientSecret, "client secret is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_AddSecret_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.AddSecretRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "system, ok",
			args: args{
				SystemCTX,
				&user.AddSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
			},
		},
		{
			name: "instance, ok",
			args: args{
				IamCTX,
				&user.AddSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
			},
		},
		{
			name: "org, error",
			args: args{
				CTX,
				&user.AddSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
			},
			wantErr: true,
		},
		{
			name: "user, error",
			args: args{
				UserCTX,
				&user.AddSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			require.NoError(t, err)
			got, err := Client.AddSecret(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.ClientSecret, "client secret is empty")
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}

func TestServer_RemoveSecret(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.RemoveSecretRequest
		prepare func(request *user.RemoveSecretRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove secret, user not existing",
			args: args{
				CTX,
				&user.RemoveSecretRequest{
					UserId: "notexisting",
				},
				func(request *user.RemoveSecretRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "remove secret, not existing",
			args: args{
				CTX,
				&user.RemoveSecretRequest{},
				func(request *user.RemoveSecretRequest) error {
					resp := Instance.CreateUserTypeMachine(CTX, Instance.DefaultOrg.Id)
					request.UserId = resp.GetId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "remove secret, ok",
			args: args{
				CTX,
				&user.RemoveSecretRequest{},
				func(request *user.RemoveSecretRequest) error {
					resp := Instance.CreateUserTypeMachine(CTX, Instance.DefaultOrg.Id)
					request.UserId = resp.GetId()
					_, err := Instance.Client.UserV2.AddSecret(CTX, &user.AddSecretRequest{
						UserId: resp.GetId(),
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
			got, err := Client.RemoveSecret(tt.args.ctx, tt.args.req)
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

func TestServer_RemoveSecret_Permission(t *testing.T) {
	otherOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	otherOrgUser, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: otherOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: integration.Username(),
			},
		},
	})
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		req     *user.RemoveSecretRequest
		prepare func(request *user.RemoveSecretRequest) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "system, ok",
			args: args{
				SystemCTX,
				&user.RemoveSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
				func(request *user.RemoveSecretRequest) error {
					_, err := Instance.Client.UserV2.AddSecret(IamCTX, &user.AddSecretRequest{
						UserId: otherOrgUser.GetId(),
					})
					return err
				},
			},
		},
		{
			name: "instance, ok",
			args: args{
				IamCTX,
				&user.RemoveSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
				func(request *user.RemoveSecretRequest) error {
					_, err := Instance.Client.UserV2.AddSecret(IamCTX, &user.AddSecretRequest{
						UserId: otherOrgUser.GetId(),
					})
					return err
				},
			},
		},
		{
			name: "org, error",
			args: args{
				CTX,
				&user.RemoveSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
				func(request *user.RemoveSecretRequest) error {
					_, err := Instance.Client.UserV2.AddSecret(IamCTX, &user.AddSecretRequest{
						UserId: otherOrgUser.GetId(),
					})
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "user, error",
			args: args{
				UserCTX,
				&user.RemoveSecretRequest{
					UserId: otherOrgUser.GetId(),
				},
				func(request *user.RemoveSecretRequest) error {
					_, err := Instance.Client.UserV2.AddSecret(IamCTX, &user.AddSecretRequest{
						UserId: otherOrgUser.GetId(),
					})
					return err
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			require.NoError(t, tt.args.prepare(tt.args.req))
			got, err := Client.RemoveSecret(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.DeletionDate, "client secret is empty")
			creationDate := got.DeletionDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
		})
	}
}
