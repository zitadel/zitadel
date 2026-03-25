package management

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func Test_listOrgEmailProvidersToModel(t *testing.T) {
	tests := []struct {
		name string
		req  *mgmt_pb.ListOrgEmailProvidersRequest
		res  *query.SMTPConfigsSearchQueries
	}{
		{
			name: "all fields filled",
			req: &mgmt_pb.ListOrgEmailProvidersRequest{
				Query: &object_pb.ListQuery{
					Offset: 100,
					Limit:  100,
					Asc:    true,
				},
			},
			res: &query.SMTPConfigsSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset: 100,
					Limit:  100,
					Asc:    true,
				},
			},
		},
		{
			name: "nil query",
			req:  &mgmt_pb.ListOrgEmailProvidersRequest{},
			res: &query.SMTPConfigsSearchQueries{
				SearchRequest: query.SearchRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listOrgEmailProvidersToModel(tt.req)
			require.NoError(t, err)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_orgEmailProvidersToPb(t *testing.T) {
	tests := []struct {
		name string
		req  []*query.SMTPConfig
		res  []*settings_pb.EmailProvider
	}{
		{
			name: "smtp and http providers",
			req: []*query.SMTPConfig{
				{
					CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ResourceOwner: "org1",
					AggregateID:   "agg",
					ID:            "id1",
					Sequence:      1,
					Description:   "smtp provider",
					SMTPConfig: &query.SMTP{
						TLS:            true,
						SenderAddress:  "sender@example.com",
						SenderName:     "Sender",
						ReplyToAddress: "reply@example.com",
						Host:           "smtp.example.com:587",
						User:           "user",
						PlainAuth:      &query.PlainAuth{},
					},
					State: 1,
				},
				{
					CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ResourceOwner: "org1",
					AggregateID:   "agg",
					ID:            "id2",
					Sequence:      2,
					Description:   "http provider",
					HTTPConfig: &query.HTTP{
						Endpoint:   "https://email.example.com",
						SigningKey: "key123",
					},
					State: 1,
				},
			},
			res: []*settings_pb.EmailProvider{
				{
					Details: &object_pb.ObjectDetails{
						Sequence:      1,
						CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ResourceOwner: "org1",
					},
					Id:          "id1",
					State:       1,
					Description: "smtp provider",
					Config: &settings_pb.EmailProvider_Smtp{
						Smtp: &settings_pb.EmailProviderSMTP{
							SenderAddress:  "sender@example.com",
							SenderName:     "Sender",
							Tls:            true,
							Host:           "smtp.example.com:587",
							User:           "user",
							ReplyToAddress: "reply@example.com",
							Auth: &settings_pb.EmailProviderSMTP_Plain{
								Plain: &settings_pb.SMTPPlainAuth{},
							},
						},
					},
				},
				{
					Details: &object_pb.ObjectDetails{
						Sequence:      2,
						CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ResourceOwner: "org1",
					},
					Id:          "id2",
					State:       1,
					Description: "http provider",
					Config: &settings_pb.EmailProvider_Http{
						Http: &settings_pb.EmailProviderHTTP{
							Endpoint:   "https://email.example.com",
							SigningKey: "key123",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orgEmailProvidersToPb(tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_orgEmailProviderToProviderPb(t *testing.T) {
	tests := []struct {
		name string
		req  *query.SMTPConfig
		res  *settings_pb.EmailProvider
	}{
		{
			name: "smtp config",
			req: &query.SMTPConfig{
				CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				ResourceOwner: "org1",
				ID:            "id",
				Sequence:      1,
				Description:   "description",
				SMTPConfig: &query.SMTP{
					TLS:            true,
					SenderAddress:  "sender",
					SenderName:     "sendername",
					ReplyToAddress: "address",
					Host:           "host",
					User:           "user",
					PlainAuth:      &query.PlainAuth{},
				},
				State: 1,
			},
			res: &settings_pb.EmailProvider{
				Details: &object_pb.ObjectDetails{
					Sequence:      1,
					CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ResourceOwner: "org1",
				},
				Id:          "id",
				State:       1,
				Description: "description",
				Config: &settings_pb.EmailProvider_Smtp{
					Smtp: &settings_pb.EmailProviderSMTP{
						SenderAddress:  "sender",
						SenderName:     "sendername",
						Tls:            true,
						Host:           "host",
						User:           "user",
						ReplyToAddress: "address",
						Auth: &settings_pb.EmailProviderSMTP_Plain{
							Plain: &settings_pb.SMTPPlainAuth{},
						},
					},
				},
			},
		},
		{
			name: "http config",
			req: &query.SMTPConfig{
				CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				ResourceOwner: "org1",
				ID:            "id",
				Sequence:      1,
				Description:   "description",
				HTTPConfig: &query.HTTP{
					Endpoint:   "endpoint",
					SigningKey: "key",
				},
				State: 1,
			},
			res: &settings_pb.EmailProvider{
				Details: &object_pb.ObjectDetails{
					Sequence:      1,
					CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ResourceOwner: "org1",
				},
				Id:          "id",
				State:       1,
				Description: "description",
				Config: &settings_pb.EmailProvider_Http{
					Http: &settings_pb.EmailProviderHTTP{
						Endpoint:   "endpoint",
						SigningKey: "key",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orgEmailProviderToProviderPb(tt.req)
			assert.EqualValues(t, tt.res, got)
		})
	}
}

func Test_orgEmailProviderStateToPb(t *testing.T) {
	tests := []struct {
		name string
		req  domain.SMTPConfigState
		res  settings_pb.EmailProviderState
	}{
		{
			name: "unspecified",
			req:  domain.SMTPConfigStateUnspecified,
			res:  settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED,
		},
		{
			name: "removed",
			req:  domain.SMTPConfigStateRemoved,
			res:  settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED,
		},
		{
			name: "active",
			req:  domain.SMTPConfigStateActive,
			res:  settings_pb.EmailProviderState_EMAIL_PROVIDER_ACTIVE,
		},
		{
			name: "inactive",
			req:  domain.SMTPConfigStateInactive,
			res:  settings_pb.EmailProviderState_EMAIL_PROVIDER_INACTIVE,
		},
		{
			name: "default",
			req:  100,
			res:  settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orgEmailProviderStateToPb(tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_orgHttpToPb(t *testing.T) {
	tests := []struct {
		name string
		req  *query.HTTP
		res  *settings_pb.EmailProvider_Http
	}{
		{
			name: "all fields filled",
			req: &query.HTTP{
				Endpoint:   "endpoint",
				SigningKey: "key",
			},
			res: &settings_pb.EmailProvider_Http{
				Http: &settings_pb.EmailProviderHTTP{
					Endpoint:   "endpoint",
					SigningKey: "key",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orgHttpToPb(tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_orgSmtpToPb(t *testing.T) {
	tests := []struct {
		name string
		req  *query.SMTP
		res  *settings_pb.EmailProvider_Smtp
	}{
		{
			name: "all fields filled with plain auth",
			req: &query.SMTP{
				SenderAddress:  "sender",
				SenderName:     "sendername",
				TLS:            true,
				Host:           "host",
				User:           "user",
				PlainAuth:      &query.PlainAuth{},
				ReplyToAddress: "address",
			},
			res: &settings_pb.EmailProvider_Smtp{
				Smtp: &settings_pb.EmailProviderSMTP{
					SenderAddress:  "sender",
					SenderName:     "sendername",
					Tls:            true,
					Host:           "host",
					User:           "user",
					ReplyToAddress: "address",
					Auth: &settings_pb.EmailProviderSMTP_Plain{
						Plain: &settings_pb.SMTPPlainAuth{},
					},
				},
			},
		},
		{
			name: "no auth maps to none",
			req:  &query.SMTP{},
			res: &settings_pb.EmailProvider_Smtp{
				Smtp: &settings_pb.EmailProviderSMTP{
					Auth: &settings_pb.EmailProviderSMTP_None{
						None: &settings_pb.SMTPNoAuth{},
					},
				},
			},
		},
		{
			name: "xoauth2 auth with client credentials",
			req: &query.SMTP{
				User: "xoauth2-user",
				XOAuth2Auth: &query.XOAuth2Auth{
					TokenEndpoint: "auth.example.com/token",
					Scopes:        []string{"scopes"},
					ClientCredentials: &query.XOAuthClientCredentials{
						ClientId: "my-client",
					},
				},
			},
			res: &settings_pb.EmailProvider_Smtp{
				Smtp: &settings_pb.EmailProviderSMTP{
					User: "xoauth2-user",
					Auth: &settings_pb.EmailProviderSMTP_Xoauth2{
						Xoauth2: &settings_pb.SMTPXOAuth2Auth{
							TokenEndpoint: "auth.example.com/token",
							Scopes:        []string{"scopes"},
							OAuth2Type: &settings_pb.SMTPXOAuth2Auth_ClientCredentials_{
								ClientCredentials: &settings_pb.SMTPXOAuth2Auth_ClientCredentials{
									ClientId: "my-client",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orgSmtpToPb(tt.req)
			assert.Equal(t, tt.res, got)
		})
	}
}
