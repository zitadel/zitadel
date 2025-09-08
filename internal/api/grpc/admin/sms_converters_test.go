package admin

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
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

func Test_listSMSConfigsToModel(t *testing.T) {
	type args struct {
		req *admin_pb.ListSMSProvidersRequest
	}
	tests := []struct {
		name string
		args args
		res  *query.SMSConfigsSearchQueries
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin_pb.ListSMSProvidersRequest{
					Query: &object_pb.ListQuery{
						Offset: 100,
						Limit:  100,
						Asc:    true,
					},
				},
			},
			res: &query.SMSConfigsSearchQueries{
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
			got, err := listSMSConfigsToModel(tt.args.req)
			require.NoError(t, err)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_SMSConfigsToPb(t *testing.T) {
	type args struct {
		req []*query.SMSConfig
	}
	tests := []struct {
		name string
		args args
		res  []*settings_pb.SMSProvider
	}{
		{
			name: "all fields filled",
			args: args{
				req: []*query.SMSConfig{
					{
						CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
						ResourceOwner: "resourceowner",
						AggregateID:   "agg",
						ID:            "id",
						Sequence:      1,
						Description:   "description",
						TwilioConfig: &query.Twilio{
							SID:              "sid",
							Token:            nil,
							SenderNumber:     "sender",
							VerifyServiceSID: "verify",
						},
						State: 1,
					},
					{
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
			},
			res: []*settings_pb.SMSProvider{
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
					Config: &settings_pb.SMSProvider_Twilio{
						Twilio: &settings_pb.TwilioConfig{
							Sid:              "sid",
							SenderNumber:     "sender",
							VerifyServiceSid: "verify",
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
					Config: &settings_pb.SMSProvider_Http{
						Http: &settings_pb.HTTPConfig{
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
			got := SMSConfigsToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_SMSConfigToProviderPb(t *testing.T) {
	type args struct {
		req *query.SMSConfig
	}
	tests := []struct {
		name string
		args args
		res  *settings_pb.SMSProvider
	}{
		{
			name: "all fields filled, twilio",
			args: args{
				req: &query.SMSConfig{

					CreationDate:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ChangeDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					ResourceOwner: "resourceowner",
					AggregateID:   "agg",
					ID:            "id",
					Sequence:      1,
					Description:   "description",
					TwilioConfig: &query.Twilio{
						SID:              "sid",
						Token:            nil,
						SenderNumber:     "sender",
						VerifyServiceSID: "verify",
					},
					State: 1,
				},
			},
			res: &settings_pb.SMSProvider{
				Details: &object_pb.ObjectDetails{
					Sequence:      1,
					CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ResourceOwner: "resourceowner",
				},
				Id:          "id",
				State:       1,
				Description: "description",
				Config: &settings_pb.SMSProvider_Twilio{
					Twilio: &settings_pb.TwilioConfig{
						Sid:              "sid",
						SenderNumber:     "sender",
						VerifyServiceSid: "verify",
					},
				},
			},
		},
		{
			name: "all fields filled, http",
			args: args{
				req: &query.SMSConfig{
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
			res: &settings_pb.SMSProvider{
				Details: &object_pb.ObjectDetails{
					Sequence:      1,
					CreationDate:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ChangeDate:    timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					ResourceOwner: "resourceowner",
				},
				Id:          "id",
				State:       1,
				Description: "description",
				Config: &settings_pb.SMSProvider_Http{
					Http: &settings_pb.HTTPConfig{
						Endpoint:   "endpoint",
						SigningKey: "key",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SMSConfigToProviderPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_smsStateToPb(t *testing.T) {
	type args struct {
		req domain.SMSConfigState
	}
	tests := []struct {
		name string
		args args
		res  settings_pb.SMSProviderConfigState
	}{
		{
			name: "unspecified",
			args: args{
				req: domain.SMSConfigStateUnspecified,
			},
			res: settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE,
		},
		{
			name: "removed",
			args: args{
				req: domain.SMSConfigStateRemoved,
			},
			res: settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE,
		},
		{
			name: "active",
			args: args{
				req: domain.SMSConfigStateActive,
			},
			res: settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_ACTIVE,
		},
		{
			name: "inactive",
			args: args{
				req: domain.SMSConfigStateInactive,
			},
			res: settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE,
		},
		{
			name: "default",
			args: args{
				req: 100,
			},
			res: settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smsStateToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_HTTPConfigToPb(t *testing.T) {
	type args struct {
		req *query.HTTP
	}
	tests := []struct {
		name string
		args args
		res  *settings_pb.SMSProvider_Http
	}{
		{
			name: "all fields filled",
			args: args{
				req: &query.HTTP{
					Endpoint:   "endpoint",
					SigningKey: "key",
				},
			},
			res: &settings_pb.SMSProvider_Http{
				Http: &settings_pb.HTTPConfig{
					Endpoint:   "endpoint",
					SigningKey: "key",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTTPConfigToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_TwilioConfigToPb(t *testing.T) {
	type args struct {
		req *query.Twilio
	}
	tests := []struct {
		name string
		args args
		res  *settings_pb.SMSProvider_Twilio
	}{
		{
			name: "all fields filled",
			args: args{
				req: &query.Twilio{
					SID:              "sid",
					SenderNumber:     "sender",
					VerifyServiceSID: "verify",
				},
			},
			res: &settings_pb.SMSProvider_Twilio{
				Twilio: &settings_pb.TwilioConfig{
					Sid:              "sid",
					SenderNumber:     "sender",
					VerifyServiceSid: "verify",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TwilioConfigToPb(tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_addSMSConfigTwilioToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.AddSMSProviderTwilioRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.AddTwilioConfig
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.AddSMSProviderTwilioRequest{
					Sid:              "sid",
					Token:            "token",
					SenderNumber:     "sender",
					Description:      "description",
					VerifyServiceSid: "verify",
				},
			},
			res: &command.AddTwilioConfig{
				ResourceOwner:    "instance",
				ID:               "",
				Description:      "description",
				SID:              "sid",
				Token:            "token",
				SenderNumber:     "sender",
				VerifyServiceSID: "verify",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addSMSConfigTwilioToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_updateSMSConfigTwilioToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.UpdateSMSProviderTwilioRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.ChangeTwilioConfig
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.UpdateSMSProviderTwilioRequest{
					Id:               "id",
					Sid:              "sid",
					SenderNumber:     "sender",
					Description:      "description",
					VerifyServiceSid: "verify",
				},
			},
			res: &command.ChangeTwilioConfig{
				ResourceOwner:    "instance",
				ID:               "id",
				Description:      gu.Ptr("description"),
				SID:              gu.Ptr("sid"),
				SenderNumber:     gu.Ptr("sender"),
				VerifyServiceSID: gu.Ptr("verify"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateSMSConfigTwilioToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_addSMSConfigHTTPToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.AddSMSProviderHTTPRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.AddSMSHTTP
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.AddSMSProviderHTTPRequest{
					Endpoint:    "endpoint",
					Description: "description",
				},
			},
			res: &command.AddSMSHTTP{
				ResourceOwner: "instance",
				ID:            "",
				Description:   "description",
				Endpoint:      "endpoint",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addSMSConfigHTTPToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}

func Test_updateSMSConfigHTTPToConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *admin_pb.UpdateSMSProviderHTTPRequest
	}
	tests := []struct {
		name string
		args args
		res  *command.ChangeSMSHTTP
	}{
		{
			name: "all fields filled",
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
				req: &admin_pb.UpdateSMSProviderHTTPRequest{
					Id:                   "id",
					Endpoint:             "endpoint",
					Description:          "description",
					ExpirationSigningKey: durationpb.New(time.Second),
				},
			},
			res: &command.ChangeSMSHTTP{
				ResourceOwner:        "instance",
				ID:                   "id",
				Description:          gu.Ptr("description"),
				Endpoint:             gu.Ptr("endpoint"),
				ExpirationSigningKey: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateSMSConfigHTTPToConfig(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.res, got)
		})
	}
}
