package admin

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
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
		Details:     object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:          config.ID,
		Description: config.Description,
		State:       smsStateToPb(config.State),
		Config:      SMSConfigToPb(config),
	}
}

func SMSConfigToPb(config *query.SMSConfig) settings_pb.SMSConfig {
	if config.TwilioConfig != nil {
		return TwilioConfigToPb(config.TwilioConfig)
	}
	if config.HTTPConfig != nil {
		return HTTPConfigToPb(config.HTTPConfig)
	}
	return nil
}

func HTTPConfigToPb(http *query.HTTP) *settings_pb.SMSProvider_Http {
	return &settings_pb.SMSProvider_Http{
		Http: &settings_pb.HTTPConfig{
			Endpoint:   http.Endpoint,
			SigningKey: http.SigningKey,
		},
	}
}

func TwilioConfigToPb(twilio *query.Twilio) *settings_pb.SMSProvider_Twilio {
	return &settings_pb.SMSProvider_Twilio{
		Twilio: &settings_pb.TwilioConfig{
			Sid:              twilio.SID,
			SenderNumber:     twilio.SenderNumber,
			VerifyServiceSid: twilio.VerifyServiceSID,
		},
	}
}

func smsStateToPb(state domain.SMSConfigState) settings_pb.SMSProviderConfigState {
	switch state {
	case domain.SMSConfigStateUnspecified, domain.SMSConfigStateRemoved:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	case domain.SMSConfigStateInactive:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	case domain.SMSConfigStateActive:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_ACTIVE
	default:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	}
}

func addSMSConfigTwilioToConfig(ctx context.Context, req *admin_pb.AddSMSProviderTwilioRequest) *command.AddTwilioConfig {
	return &command.AddTwilioConfig{
		ResourceOwner:    authz.GetInstance(ctx).InstanceID(),
		Description:      req.Description,
		SID:              req.Sid,
		SenderNumber:     req.SenderNumber,
		Token:            req.Token,
		VerifyServiceSID: req.VerifyServiceSid,
	}
}

func updateSMSConfigTwilioToConfig(ctx context.Context, req *admin_pb.UpdateSMSProviderTwilioRequest) *command.ChangeTwilioConfig {
	return &command.ChangeTwilioConfig{
		ResourceOwner:    authz.GetInstance(ctx).InstanceID(),
		ID:               req.Id,
		Description:      gu.Ptr(req.Description),
		SID:              gu.Ptr(req.Sid),
		SenderNumber:     gu.Ptr(req.SenderNumber),
		VerifyServiceSID: gu.Ptr(req.VerifyServiceSid),
	}
}

func addSMSConfigHTTPToConfig(ctx context.Context, req *admin_pb.AddSMSProviderHTTPRequest) *command.AddSMSHTTP {
	return &command.AddSMSHTTP{
		ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		Description:   req.GetDescription(),
		Endpoint:      req.GetEndpoint(),
	}
}

func updateSMSConfigHTTPToConfig(ctx context.Context, req *admin_pb.UpdateSMSProviderHTTPRequest) *command.ChangeSMSHTTP {
	// TODO handle expiration, currently only immediate expiration is supported
	expirationSigningKey := req.GetExpirationSigningKey() != nil
	return &command.ChangeSMSHTTP{
		ResourceOwner:        authz.GetInstance(ctx).InstanceID(),
		ID:                   req.Id,
		Description:          gu.Ptr(req.Description),
		Endpoint:             gu.Ptr(req.Endpoint),
		ExpirationSigningKey: expirationSigningKey,
	}
}
