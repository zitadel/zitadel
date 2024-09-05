package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func listEmailProvidersToModel(req *admin_pb.ListEmailProvidersRequest) (*query.SMTPConfigsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.SMTPConfigsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}

func EmailProvidersToPb(configs []*query.SMTPConfig) []*settings_pb.EmailProvider {
	c := make([]*settings_pb.EmailProvider, len(configs))
	for i, config := range configs {
		c[i] = EmailProviderToProviderPb(config)
	}
	return c
}

func EmailProviderToProviderPb(config *query.SMTPConfig) *settings_pb.EmailProvider {
	return &settings_pb.EmailProvider{
		Details:     object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:          config.ID,
		Description: config.Description,
		State:       settings_pb.EmailProviderState(config.State),
		Config:      EmailProviderToPb(config),
	}
}

func EmailProviderToPb(config *query.SMTPConfig) settings_pb.EmailConfig {
	if config.SMTPConfig != nil {
		return SMTPToPb(config.SMTPConfig)
	}
	if config.HTTPConfig != nil {
		return HTTPToPb(config.HTTPConfig)
	}
	return nil
}

func HTTPToPb(http *query.HTTP) *settings_pb.EmailProvider_Http {
	return &settings_pb.EmailProvider_Http{
		Http: &settings_pb.EmailProviderHTTP{
			Endpoint: http.Endpoint,
		},
	}
}

func SMTPToPb(config *query.SMTP) *settings_pb.EmailProvider_Smtp {
	return &settings_pb.EmailProvider_Smtp{
		Smtp: &settings_pb.EmailProviderSMTP{
			Tls:           config.TLS,
			Host:          config.Host,
			User:          config.User,
			SenderAddress: config.SenderAddress,
			SenderName:    config.SenderName,
		},
	}
}

func addEmailProviderSMTPToConfig(ctx context.Context, req *admin_pb.AddEmailProviderSMTPRequest) *command.AddSMTPConfig {
	return &command.AddSMTPConfig{
		ResourceOwner:  authz.GetInstance(ctx).InstanceID(),
		Description:    req.Description,
		Tls:            req.Tls,
		From:           req.SenderAddress,
		FromName:       req.SenderName,
		ReplyToAddress: req.ReplyToAddress,
		Host:           req.Host,
		User:           req.User,
		Password:       req.Password,
	}
}

func updateEmailProviderSMTPToConfig(ctx context.Context, req *admin_pb.UpdateEmailProviderSMTPRequest) *command.ChangeSMTPConfig {
	return &command.ChangeSMTPConfig{
		ResourceOwner:  authz.GetInstance(ctx).InstanceID(),
		ID:             req.Id,
		Description:    req.Description,
		Tls:            req.Tls,
		From:           req.SenderAddress,
		FromName:       req.SenderName,
		ReplyToAddress: req.ReplyToAddress,
		Host:           req.Host,
		User:           req.User,
		Password:       req.Password,
	}
}

func addEmailProviderHTTPToConfig(ctx context.Context, req *admin_pb.AddEmailProviderHTTPRequest) *command.AddSMTPConfigHTTP {
	return &command.AddSMTPConfigHTTP{
		ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		Description:   req.Description,
		Endpoint:      req.Endpoint,
	}
}

func updateEmailProviderHTTPToConfig(ctx context.Context, req *admin_pb.UpdateEmailProviderHTTPRequest) *command.ChangeSMTPConfigHTTP {
	return &command.ChangeSMTPConfigHTTP{
		ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		ID:            req.Id,
		Description:   req.Description,
		Endpoint:      req.Endpoint,
	}
}
