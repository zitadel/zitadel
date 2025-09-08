package admin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func Test_listEmailProvidersToModel(t *testing.T) {
	type args struct {
		req *admin_pb.ListEmailProvidersRequest
	}
	tests := []struct {
		name string
		args args
		res  *query.SMTPConfigsSearchQueries
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin_pb.ListEmailProvidersRequest{
					Query: &object_pb.ListQuery{
						Offset: 100,
						Limit:  100,
						Asc:    true,
					},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listEmailProvidersToModel(tt.args.req)
			require.NoError(t, err)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_emailProvidersToPb(t *testing.T) {
	type args struct {
		req []*query.SMTPConfig
	}
	tests := []struct {
		name string
		args args
		res  []*settings_pb.EmailProvider
	}{
		{
			name: "all fields filled",
			args: args{
				req: []*query.SMTPConfig{
					{
						CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						ResourceOwner: "resourceowner",
						AggregateID:   "agg",
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
						},
						HTTPConfig: nil,
						State:      1,
					},
					{
						CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						ResourceOwner: "resourceowner",
						AggregateID:   "agg",
						ID:            "id",
						Sequence:      1,
						Description:   "description",
						SMTPConfig:    nil,
						HTTPConfig: &query.HTTP{
							Endpoint:   "endpoint",
							SigningKey: "key",
						},
						State: 1,
					},
				},
			},
			res: []*settings_pb.EmailProvider{
				{
					Details: &object_pb.ObjectDetails{
						Sequence:      1,
						CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ResourceOwner: "resourceowner",
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
						},
					},
				},
				{
					Details: &object_pb.ObjectDetails{
						Sequence:      1,
						CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
						ResourceOwner: "resourceowner",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := emailProvidersToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_emailProviderToProviderPb(t *testing.T) {
	type args struct {
		req *query.SMTPConfig
	}
	tests := []struct {
		name string
		args args
		res  *settings_pb.EmailProvider
	}{
		{
			name: "all fields filled, smtp",
			args: args{
				req: &query.SMTPConfig{

					CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ResourceOwner: "resourceowner",
					AggregateID:   "agg",
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
					},
					HTTPConfig: &query.HTTP{
						Endpoint:   "endpoint",
						SigningKey: "key",
					},
					State: 1,
				},
			},
			res: &settings_pb.EmailProvider{
				Details: &object_pb.ObjectDetails{
					Sequence:      1,
					CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ResourceOwner: "resourceowner",
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
					},
				},
			},
		},
		{
			name: "all fields filled, http",
			args: args{
				req: &query.SMTPConfig{
					CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ResourceOwner: "resourceowner",
					AggregateID:   "agg",
					ID:            "id",
					Sequence:      1,
					Description:   "description",
					HTTPConfig: &query.HTTP{
						Endpoint:   "endpoint",
						SigningKey: "key",
					},
					State: 1,
				},
			},
			res: &settings_pb.EmailProvider{
				Details: &object_pb.ObjectDetails{
					Sequence:      1,
					CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ResourceOwner: "resourceowner",
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
			got := emailProviderToProviderPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_emailProviderStateToPb(t *testing.T) {
	type args struct {
		req domain.SMTPConfigState
	}
	tests := []struct {
		name string
		args args
		res  settings_pb.EmailProviderState
	}{
		{
			name: "unspecified",
			args: args{
				req: domain.SMTPConfigStateUnspecified,
			},
			res: settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED,
		},
		{
			name: "removed",
			args: args{
				req: domain.SMTPConfigStateRemoved,
			},
			res: settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED,
		},
		{
			name: "active",
			args: args{
				req: domain.SMTPConfigStateActive,
			},
			res: settings_pb.EmailProviderState_EMAIL_PROVIDER_ACTIVE,
		},
		{
			name: "inactive",
			args: args{
				req: domain.SMTPConfigStateInactive,
			},
			res: settings_pb.EmailProviderState_EMAIL_PROVIDER_INACTIVE,
		},
		{
			name: "default",
			args: args{
				req: 100,
			},
			res: settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := emailProviderStateToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_httpToPb(t *testing.T) {
	type args struct {
		req *query.HTTP
	}
	tests := []struct {
		name string
		args args
		res  *settings_pb.EmailProvider_Http
	}{
		{
			name: "all fields filled",
			args: args{
				req: &query.HTTP{
					Endpoint:   "endpoint",
					SigningKey: "key",
				},
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
			got := httpToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_smtpToPb(t *testing.T) {
	type args struct {
		req *query.SMTP
	}
	tests := []struct {
		name string
		args args
		res  *settings_pb.EmailProvider_Smtp
	}{
		{
			name: "all fields filled",
			args: args{
				req: &query.SMTP{
					SenderAddress:  "sender",
					SenderName:     "sendername",
					TLS:            true,
					Host:           "host",
					User:           "user",
					ReplyToAddress: "address",
				},
			},
			res: &settings_pb.EmailProvider_Smtp{
				Smtp: &settings_pb.EmailProviderSMTP{
					SenderAddress:  "sender",
					SenderName:     "sendername",
					Tls:            true,
					Host:           "host",
					User:           "user",
					ReplyToAddress: "address",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smtpToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_addEmailProviderSMTPToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.AddEmailProviderSMTPRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.AddSMTPConfig
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.AddEmailProviderSMTPRequest{
					SenderAddress:  "sender",
					SenderName:     "sendername",
					Tls:            true,
					Host:           "host",
					User:           "user",
					Password:       "password",
					ReplyToAddress: "address",
					Description:    "description",
				},
			},
			res: &command.AddSMTPConfig{
				ResourceOwner:  "instance",
				Description:    "description",
				Host:           "host",
				User:           "user",
				Password:       "password",
				Tls:            true,
				From:           "sender",
				FromName:       "sendername",
				ReplyToAddress: "address",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addEmailProviderSMTPToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_updateEmailProviderSMTPToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.UpdateEmailProviderSMTPRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.ChangeSMTPConfig
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.UpdateEmailProviderSMTPRequest{
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
			},
			res: &command.ChangeSMTPConfig{
				ResourceOwner:  "instance",
				ID:             "id",
				Description:    "description",
				Host:           "host",
				User:           "user",
				Password:       "password",
				Tls:            true,
				From:           "sender",
				FromName:       "sendername",
				ReplyToAddress: "address",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateEmailProviderSMTPToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_addEmailProviderHTTPToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.AddEmailProviderHTTPRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.AddSMTPConfigHTTP
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.AddEmailProviderHTTPRequest{
					Endpoint:    "endpoint",
					Description: "description",
				},
			},
			res: &command.AddSMTPConfigHTTP{
				ResourceOwner: "instance",
				ID:            "",
				Description:   "description",
				Endpoint:      "endpoint",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addEmailProviderHTTPToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_updateEmailProviderHTTPToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.UpdateEmailProviderHTTPRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.ChangeSMTPConfigHTTP
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.UpdateEmailProviderHTTPRequest{
					Id:                   "id",
					Endpoint:             "endpoint",
					Description:          "description",
					ExpirationSigningKey: durationpb.New(time.Second),
				},
			},
			res: &command.ChangeSMTPConfigHTTP{
				ResourceOwner:        "instance",
				ID:                   "id",
				Description:          "description",
				Endpoint:             "endpoint",
				ExpirationSigningKey: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateEmailProviderHTTPToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}
