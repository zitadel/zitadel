package management

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func ctxWithOrg(orgID string) context.Context {
	return authz.SetCtxData(context.Background(), authz.CtxData{OrgID: orgID})
}

func Test_addOrgEmailProviderSMTPToConfig(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		req  *mgmt_pb.AddOrgEmailProviderSMTPRequest
		res  *command.AddOrgSMTPConfig
	}{
		{
			name: "all fields filled",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				SenderAddress:  "sender",
				SenderName:     "sendername",
				Tls:            true,
				Host:           "host",
				User:           "user",
				Password:       "password",
				ReplyToAddress: "address",
				Description:    "description",
			},
			res: &command.AddOrgSMTPConfig{
				ResourceOwner: "org1",
				Description:   "description",
				Host:          "host",
				User:          "user",
				PlainAuth: &command.PlainAuth{
					Password: "password",
				},
				Tls:            true,
				From:           "sender",
				FromName:       "sendername",
				ReplyToAddress: "address",
			},
		},
		{
			name: "legacy auth (username password) should map to plain",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				User:     "user",
				Password: "password",
			},
			res: &command.AddOrgSMTPConfig{
				ResourceOwner: "org1",
				User:          "user",
				PlainAuth: &command.PlainAuth{
					Password: "password",
				},
			},
		},
		{
			name: "plain auth should map to plain",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				User: "plain-user",
				Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
					Plain: &mgmt_pb.OrgSMTPPlainAuth{
						Password: "other_password",
					},
				},
			},
			res: &command.AddOrgSMTPConfig{
				ResourceOwner: "org1",
				User:          "plain-user",
				PlainAuth: &command.PlainAuth{
					Password: "other_password",
				},
			},
		},
		{
			name: "xoauth2 auth should map to xoauth2",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				User: "xoauth2-user",
				Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Xoauth2{
					Xoauth2: &mgmt_pb.OrgSMTPXOAuth2Auth{
						TokenEndpoint: "auth.example.com/token",
						Scopes:        []string{"scopes"},
						OAuth2Type: &mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials_{
							ClientCredentials: &mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials{
								ClientId:     "my-client",
								ClientSecret: "some-secret",
							},
						},
					},
				},
			},
			res: &command.AddOrgSMTPConfig{
				ResourceOwner: "org1",
				User:          "xoauth2-user",
				XOAuth2Auth: &command.XOAuth2Auth{
					TokenEndpoint: "auth.example.com/token",
					Scopes:        []string{"scopes"},
					ClientCredentialsAuth: &command.OAuth2ClientCredentials{
						ClientId:     "my-client",
						ClientSecret: "some-secret",
					},
				},
			},
		},
		{
			name: "none auth should not set any auth",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				Host: "host",
				Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_None{
					None: &mgmt_pb.OrgSMTPNoAuth{},
				},
			},
			res: &command.AddOrgSMTPConfig{
				ResourceOwner: "org1",
				Host:          "host",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addOrgEmailProviderSMTPToConfig(tt.ctx, tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_updateOrgEmailProviderSMTPToConfig(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		req  *mgmt_pb.UpdateOrgEmailProviderSMTPRequest
		res  *command.ChangeOrgSMTPConfig
	}{
		{
			name: "all fields filled",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				SenderAddress:  "sender",
				SenderName:     "sendername",
				Tls:            true,
				Host:           "host",
				User:           "user",
				ReplyToAddress: "address",
				Password:       "password",
				Description:    "description",
				Id:             "id",
			},
			res: &command.ChangeOrgSMTPConfig{
				ResourceOwner: "org1",
				ID:            "id",
				Description:   "description",
				Host:          "host",
				User:          "user",
				PlainAuth: &command.PlainAuth{
					Password: "password",
				},
				Tls:            true,
				From:           "sender",
				FromName:       "sendername",
				ReplyToAddress: "address",
			},
		},
		{
			name: "legacy auth (username password) should map to plain",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				User:     "user",
				Password: "password",
				Id:       "id",
			},
			res: &command.ChangeOrgSMTPConfig{
				ResourceOwner: "org1",
				ID:            "id",
				User:          "user",
				PlainAuth: &command.PlainAuth{
					Password: "password",
				},
			},
		},
		{
			name: "plain auth should map to plain",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				Id:   "id",
				User: "plain-user",
				Auth: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Plain{
					Plain: &mgmt_pb.OrgSMTPPlainAuth{
						Password: "other_password",
					},
				},
			},
			res: &command.ChangeOrgSMTPConfig{
				ResourceOwner: "org1",
				ID:            "id",
				User:          "plain-user",
				PlainAuth: &command.PlainAuth{
					Password: "other_password",
				},
			},
		},
		{
			name: "xoauth2 auth should map to xoauth2",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				Id:   "id",
				User: "xoauth2-user",
				Auth: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Xoauth2{
					Xoauth2: &mgmt_pb.OrgSMTPXOAuth2Auth{
						TokenEndpoint: "auth.example.com/token",
						Scopes:        []string{"scopes"},
						OAuth2Type: &mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials_{
							ClientCredentials: &mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials{
								ClientId:     "my-client",
								ClientSecret: "some-secret",
							},
						},
					},
				},
			},
			res: &command.ChangeOrgSMTPConfig{
				ResourceOwner: "org1",
				ID:            "id",
				User:          "xoauth2-user",
				XOAuth2Auth: &command.XOAuth2Auth{
					TokenEndpoint: "auth.example.com/token",
					Scopes:        []string{"scopes"},
					ClientCredentialsAuth: &command.OAuth2ClientCredentials{
						ClientId:     "my-client",
						ClientSecret: "some-secret",
					},
				},
			},
		},
		{
			name: "none auth should not set any auth",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				Id:   "id",
				Host: "host",
				Auth: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest_None{
					None: &mgmt_pb.OrgSMTPNoAuth{},
				},
			},
			res: &command.ChangeOrgSMTPConfig{
				ResourceOwner: "org1",
				ID:            "id",
				Host:          "host",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateOrgEmailProviderSMTPToConfig(tt.ctx, tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_addOrgEmailProviderHTTPToConfig(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		req  *mgmt_pb.AddOrgEmailProviderHTTPRequest
		res  *command.AddOrgSMTPConfigHTTP
	}{
		{
			name: "all fields filled",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.AddOrgEmailProviderHTTPRequest{
				Endpoint:    "endpoint",
				Description: "description",
			},
			res: &command.AddOrgSMTPConfigHTTP{
				ResourceOwner: "org1",
				Description:   "description",
				Endpoint:      "endpoint",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addOrgEmailProviderHTTPToConfig(tt.ctx, tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_updateOrgEmailProviderHTTPToConfig(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		req  *mgmt_pb.UpdateOrgEmailProviderHTTPRequest
		res  *command.ChangeOrgSMTPConfigHTTP
	}{
		{
			name: "all fields filled",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderHTTPRequest{
				Id:                   "id",
				Endpoint:             "endpoint",
				Description:          "description",
				ExpirationSigningKey: durationpb.New(time.Second),
			},
			res: &command.ChangeOrgSMTPConfigHTTP{
				ResourceOwner:        "org1",
				ID:                   "id",
				Description:          "description",
				Endpoint:             "endpoint",
				ExpirationSigningKey: true,
			},
		},
		{
			name: "no expiration signing key",
			ctx:  ctxWithOrg("org1"),
			req: &mgmt_pb.UpdateOrgEmailProviderHTTPRequest{
				Id:          "id",
				Endpoint:    "endpoint",
				Description: "description",
			},
			res: &command.ChangeOrgSMTPConfigHTTP{
				ResourceOwner:        "org1",
				ID:                   "id",
				Description:          "description",
				Endpoint:             "endpoint",
				ExpirationSigningKey: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateOrgEmailProviderHTTPToConfig(tt.ctx, tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_testOrgEmailProviderSMTPToConfig(t *testing.T) {
	tests := []struct {
		name string
		req  *mgmt_pb.TestOrgEmailProviderSMTPRequest
		res  *smtp.Config
	}{
		{
			name: "all fields filled",
			req: &mgmt_pb.TestOrgEmailProviderSMTPRequest{
				SenderAddress: "sender",
				SenderName:    "sendername",
				Tls:           true,
				Host:          "host",
				User:          "user",
				Password:      "password",
			},
			res: &smtp.Config{
				Tls:      true,
				From:     "sender",
				FromName: "sendername",
				SMTP: smtp.SMTP{
					Host: "host",
					PlainAuth: &smtp.PlainAuthConfig{
						User:     "user",
						Password: "password",
					},
				},
			},
		},
		{
			name: "legacy auth (username password) should map to plain",
			req: &mgmt_pb.TestOrgEmailProviderSMTPRequest{
				User:     "user",
				Password: "password",
			},
			res: &smtp.Config{
				SMTP: smtp.SMTP{
					PlainAuth: &smtp.PlainAuthConfig{
						User:     "user",
						Password: "password",
					},
				},
			},
		},
		{
			name: "plain auth should map to plain",
			req: &mgmt_pb.TestOrgEmailProviderSMTPRequest{
				User: "plain-user",
				Auth: &mgmt_pb.TestOrgEmailProviderSMTPRequest_Plain{
					Plain: &mgmt_pb.OrgSMTPPlainAuth{
						Password: "other_password",
					},
				},
			},
			res: &smtp.Config{
				SMTP: smtp.SMTP{
					PlainAuth: &smtp.PlainAuthConfig{
						User:     "plain-user",
						Password: "other_password",
					},
				},
			},
		},
		{
			name: "xoauth2 auth should map to xoauth2",
			req: &mgmt_pb.TestOrgEmailProviderSMTPRequest{
				User: "xoauth2-user",
				Auth: &mgmt_pb.TestOrgEmailProviderSMTPRequest_Xoauth2{
					Xoauth2: &mgmt_pb.OrgSMTPXOAuth2Auth{
						TokenEndpoint: "auth.example.com/token",
						Scopes:        []string{"scopes"},
						OAuth2Type: &mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials_{
							ClientCredentials: &mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials{
								ClientId:     "my-client",
								ClientSecret: "some-secret",
							},
						},
					},
				},
			},
			res: &smtp.Config{
				SMTP: smtp.SMTP{
					XOAuth2Auth: &smtp.XOAuth2AuthConfig{
						User:          "xoauth2-user",
						TokenEndpoint: "auth.example.com/token",
						Scopes:        []string{"scopes"},
						ClientCredentialsAuth: &smtp.OAuth2ClientCredentials{
							ClientId:     "my-client",
							ClientSecret: "some-secret",
						},
					},
				},
			},
		},
		{
			name: "none auth should not set any auth",
			req: &mgmt_pb.TestOrgEmailProviderSMTPRequest{
				Host: "host",
				Auth: &mgmt_pb.TestOrgEmailProviderSMTPRequest_None{
					None: &mgmt_pb.OrgSMTPNoAuth{},
				},
			},
			res: &smtp.Config{
				SMTP: smtp.SMTP{
					Host: "host",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testOrgEmailProviderSMTPToConfig(tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}
