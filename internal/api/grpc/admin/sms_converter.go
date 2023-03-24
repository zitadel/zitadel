package admin

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
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

func SMSConfigsToPb(configs []*query.SMSConfig) []*settings_pb.SMSProvider {
	c := make([]*settings_pb.SMSProvider, len(configs))
	for i, config := range configs {
		c[i] = SMSConfigToProviderPb(config)
	}
	return c
}

func SMSConfigToProviderPb(config *query.SMSConfig) *settings_pb.SMSProvider {
	return &settings_pb.SMSProvider{
		Details: object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:      config.ID,
		State:   smsStateToPb(config.State),
		Config:  SMSConfigToPb(config),
	}
}

func SMSConfigToPb(config *query.SMSConfig) settings_pb.SMSConfig {
	if config.TwilioConfig != nil {
		return TwilioConfigToPb(config.TwilioConfig)
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

func AddSMSConfigTwilioToConfig(req *admin_pb.AddSMSProviderTwilioRequest) *twilio.Config {
	return &twilio.Config{
		SID:          req.Sid,
		SenderNumber: req.SenderNumber,
		Token:        req.Token,
	}
}

func UpdateSMSConfigTwilioToConfig(req *admin_pb.UpdateSMSProviderTwilioRequest) *twilio.Config {
	return &twilio.Config{
		SID:          req.Sid,
		SenderNumber: req.SenderNumber,
	}
}
