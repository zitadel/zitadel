package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	settings_pb "github.com/caos/zitadel/pkg/grpc/settings"
)

func listSMSConfigsToModel(req *admin_pb.ListSMSProvidersRequest) (*query.SMSConfigsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.SMSConfigsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}

func SMSConfigToPb(app *query.SMSConfig) settings_pb.SMSConfig {
	if app.TwilioConfig != nil {
		return TwilioConfigToPb(app.TwilioConfig)
	}
	return nil
}

func TwilioConfigToPb(twilio *query.Twilio) *settings_pb.SMSProvider_Twilio {
	return &settings_pb.SMSProvider_Twilio{
		Twilio: &settings_pb.TwilioConfig{
			Sid:          twilio.SID,
			SenderNumber: twilio.SenderNumber,
		},
	}
}

func smsStateToPb(state domain.SMSConfigState) settings_pb.SMSProviderConfigState {
	switch state {
	case domain.SMSConfigStateInactive:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	case domain.SMSConfigStateActive:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_ACTIVE
	default:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	}
}

func AddSMSConfigTwilioToConfig(req *admin_pb.AddSMSProviderTwilioRequest) *twilio.TwilioConfig {
	return &twilio.TwilioConfig{
		SID:          req.Sid,
		SenderNumber: req.SenderNumber,
		Token:        req.Token,
	}
}

func UpdateSMSConfigTwilioToConfig(req *admin_pb.UpdateSMSProviderTwilioRequest) *twilio.TwilioConfig {
	return &twilio.TwilioConfig{
		SID:          req.Sid,
		SenderNumber: req.SenderNumber,
	}
}
