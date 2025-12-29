package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
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

func emailProvidersToPb(configs []*query.SMTPConfig) []*settings_pb.EmailProvider {
	c := make([]*settings_pb.EmailProvider, len(configs))
	for i, config := range configs {
		c[i] = emailProviderToProviderPb(config)
	}
	return c
}

func emailProviderToProviderPb(config *query.SMTPConfig) *settings_pb.EmailProvider {
	return &settings_pb.EmailProvider{
		Details:     object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:          config.ID,
		Description: config.Description,
		State:       emailProviderStateToPb(config.State),
		Config:      emailProviderToPb(config),
	}
}

func emailProviderStateToPb(state domain.SMTPConfigState) settings_pb.EmailProviderState {
	switch state {
	case domain.SMTPConfigStateUnspecified, domain.SMTPConfigStateRemoved:
		return settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED
	case domain.SMTPConfigStateActive:
		return settings_pb.EmailProviderState_EMAIL_PROVIDER_ACTIVE
	case domain.SMTPConfigStateInactive:
		return settings_pb.EmailProviderState_EMAIL_PROVIDER_INACTIVE
	default:
		return settings_pb.EmailProviderState_EMAIL_PROVIDER_STATE_UNSPECIFIED
	}
}

func emailProviderToPb(config *query.SMTPConfig) settings_pb.EmailConfig {
	if config.SMTPConfig != nil {
		return smtpToPb(config.SMTPConfig)
	}
	if config.HTTPConfig != nil {
		return httpToPb(config.HTTPConfig)
	}
	return nil
}

func httpToPb(http *query.HTTP) *settings_pb.EmailProvider_Http {
	return &settings_pb.EmailProvider_Http{
		Http: &settings_pb.EmailProviderHTTP{
			Endpoint:   http.Endpoint,
			SigningKey: http.SigningKey,
		},
	}
}

func smtpToPb(config *query.SMTP) *settings_pb.EmailProvider_Smtp {
	ret := &settings_pb.EmailProvider_Smtp{
		Smtp: &settings_pb.EmailProviderSMTP{
			Tls:            config.TLS,
			Host:           config.Host,
			SenderAddress:  config.SenderAddress,
			SenderName:     config.SenderName,
			ReplyToAddress: config.ReplyToAddress,
		},
	}

	if config.PlainAuth != nil {
		ret.Smtp.User = config.PlainAuth.User
		ret.Smtp.Auth = &settings_pb.EmailProviderSMTP_Plain{
			Plain: &settings_pb.SMTPPlainAuth{
				User: config.PlainAuth.User,
			},
		}
	}
	if config.XOAuth2Auth != nil {
		ret.Smtp.Auth = &settings_pb.EmailProviderSMTP_Xoauth2{
			Xoauth2: &settings_pb.SMTPXOAuth2Auth{
				User:          config.XOAuth2Auth.User,
				ClientId:      config.XOAuth2Auth.ClientId,
				TokenEndpoint: config.XOAuth2Auth.TokenEndpoint,
				Scopes:        config.XOAuth2Auth.Scopes,
			},
		}
	}

	return ret
}

func addEmailProviderSMTPToConfig(ctx context.Context, req *admin_pb.AddEmailProviderSMTPRequest) *command.AddSMTPConfig {
	cmd := &command.AddSMTPConfig{
		ResourceOwner:  authz.GetInstance(ctx).InstanceID(),
		Description:    req.Description,
		Tls:            req.Tls,
		From:           req.SenderAddress,
		FromName:       req.SenderName,
		ReplyToAddress: req.ReplyToAddress,
		Host:           req.Host,
	}

	switch v := req.Auth.(type) {
	case *admin_pb.AddEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *admin_pb.AddEmailProviderSMTPRequest_Plain:
		cmd.PlainAuth = &smtp.PlainAuthConfig{
			User:     v.Plain.User,
			Password: v.Plain.Password,
		}
	case *admin_pb.AddEmailProviderSMTPRequest_Xoauth2:
		cmd.XOAuth2Auth = &smtp.XOAuth2AuthConfig{
			User:          v.Xoauth2.User,
			ClientId:      v.Xoauth2.ClientId,
			ClientSecret:  v.Xoauth2.ClientSecret,
			TokenEndpoint: v.Xoauth2.TokenEndpoint,
			Scopes:        v.Xoauth2.Scopes,
		}
	default:
		// ensure backwards compatibility
		if req.User != "" || req.Password != "" {
			cmd.PlainAuth = &smtp.PlainAuthConfig{
				User:     req.User,
				Password: req.Password,
			}
		}
	}

	return cmd
}

func updateEmailProviderSMTPToConfig(ctx context.Context, req *admin_pb.UpdateEmailProviderSMTPRequest) *command.ChangeSMTPConfig {
	cmd := &command.ChangeSMTPConfig{
		ResourceOwner:  authz.GetInstance(ctx).InstanceID(),
		ID:             req.Id,
		Description:    req.Description,
		Tls:            req.Tls,
		From:           req.SenderAddress,
		FromName:       req.SenderName,
		ReplyToAddress: req.ReplyToAddress,
		Host:           req.Host,
	}

	switch v := req.Auth.(type) {
	case *admin_pb.UpdateEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *admin_pb.UpdateEmailProviderSMTPRequest_Plain:
		cmd.PlainAuth = &smtp.PlainAuthConfig{
			User:     v.Plain.User,
			Password: v.Plain.Password,
		}
	case *admin_pb.UpdateEmailProviderSMTPRequest_Xoauth2:
		cmd.XOAuth2Auth = &smtp.XOAuth2AuthConfig{
			User:          v.Xoauth2.User,
			ClientId:      v.Xoauth2.ClientId,
			ClientSecret:  v.Xoauth2.ClientSecret,
			TokenEndpoint: v.Xoauth2.TokenEndpoint,
			Scopes:        v.Xoauth2.Scopes,
		}
	default:
		// ensure backwards compatibility
		if req.User != "" || req.Password != "" {
			cmd.PlainAuth = &smtp.PlainAuthConfig{
				User:     req.User,
				Password: req.Password,
			}
		}
	}

	return cmd
}

func addEmailProviderHTTPToConfig(ctx context.Context, req *admin_pb.AddEmailProviderHTTPRequest) *command.AddSMTPConfigHTTP {
	return &command.AddSMTPConfigHTTP{
		ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		Description:   req.Description,
		Endpoint:      req.Endpoint,
	}
}

func updateEmailProviderHTTPToConfig(ctx context.Context, req *admin_pb.UpdateEmailProviderHTTPRequest) *command.ChangeSMTPConfigHTTP {
	// TODO handle expiration, currently only immediate expiration is supported
	expirationSigningKey := req.GetExpirationSigningKey() != nil
	return &command.ChangeSMTPConfigHTTP{
		ResourceOwner:        authz.GetInstance(ctx).InstanceID(),
		ID:                   req.Id,
		Description:          req.Description,
		Endpoint:             req.Endpoint,
		ExpirationSigningKey: expirationSigningKey,
	}
}

func testEmailProviderSMTPToConfig(req *admin_pb.TestEmailProviderSMTPRequest) *smtp.Config {
	cfg := &smtp.Config{
		Tls:      req.Tls,
		From:     req.SenderAddress,
		FromName: req.SenderName,
		SMTP: smtp.SMTP{
			Host: req.Host,
		},
	}

	switch v := req.Auth.(type) {
	case *admin_pb.TestEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *admin_pb.TestEmailProviderSMTPRequest_Plain:
		cfg.SMTP.PlainAuth = &smtp.PlainAuthConfig{
			User:     v.Plain.User,
			Password: v.Plain.Password,
		}
	case *admin_pb.TestEmailProviderSMTPRequest_Xoauth2:
		cfg.SMTP.XOAuth2Auth = &smtp.XOAuth2AuthConfig{
			User:          v.Xoauth2.User,
			ClientId:      v.Xoauth2.ClientId,
			ClientSecret:  v.Xoauth2.ClientSecret,
			TokenEndpoint: v.Xoauth2.TokenEndpoint,
			Scopes:        v.Xoauth2.Scopes,
		}
	default:
		// ensure backwards compatibility
		if req.User != "" || req.Password != "" {
			cfg.SMTP.PlainAuth = &smtp.PlainAuthConfig{
				User:     req.User,
				Password: req.Password,
			}
		}
	}

	return cfg
}
