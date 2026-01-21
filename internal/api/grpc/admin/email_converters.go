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
			User:           config.User,
			SenderAddress:  config.SenderAddress,
			SenderName:     config.SenderName,
			ReplyToAddress: config.ReplyToAddress,
		},
	}

	if config.PlainAuth != nil {
		ret.Smtp.Auth = &settings_pb.EmailProviderSMTP_Plain{
			Plain: &settings_pb.SMTPPlainAuth{},
		}
	}
	if config.XOAuth2Auth != nil {
		xoauth2 := &settings_pb.EmailProviderSMTP_Xoauth2{
			Xoauth2: &settings_pb.SMTPXOAuth2Auth{
				TokenEndpoint: config.XOAuth2Auth.TokenEndpoint,
				Scopes:        config.XOAuth2Auth.Scopes,
			},
		}
		if config.XOAuth2Auth.ClientCredentials != nil {
			xoauth2.Xoauth2.OAuth2Type = &settings_pb.SMTPXOAuth2Auth_ClientCredentials_{
				ClientCredentials: &settings_pb.SMTPXOAuth2Auth_ClientCredentials{
					ClientId: config.XOAuth2Auth.ClientCredentials.ClientId,
				},
			}
		}
		ret.Smtp.Auth = xoauth2
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
		User:           req.User,
	}

	switch v := req.Auth.(type) {
	case *admin_pb.AddEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *admin_pb.AddEmailProviderSMTPRequest_Plain:
		cmd.PlainAuth = &command.PlainAuth{
			Password: v.Plain.Password,
		}
	case *admin_pb.AddEmailProviderSMTPRequest_Xoauth2:
		if xoauth2, ok := v.Xoauth2.OAuth2Type.(*admin_pb.SMTPXOAuth2Auth_ClientCredentials_); ok {
			cmd.XOAuth2Auth = &command.XOAuth2Auth{
				TokenEndpoint: v.Xoauth2.TokenEndpoint,
				Scopes:        v.Xoauth2.Scopes,
				ClientCredentialsAuth: &command.OAuth2ClientCredentials{
					ClientId:     xoauth2.ClientCredentials.ClientId,
					ClientSecret: xoauth2.ClientCredentials.ClientSecret,
				},
			}
		}
	default:
		// ensure backwards compatibility
		//nolint:staticcheck
		if req.User != "" || req.Password != "" {
			cmd.PlainAuth = &command.PlainAuth{
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
		User:           req.User,
	}

	switch v := req.Auth.(type) {
	case *admin_pb.UpdateEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *admin_pb.UpdateEmailProviderSMTPRequest_Plain:
		cmd.PlainAuth = &command.PlainAuth{
			Password: v.Plain.Password,
		}
	case *admin_pb.UpdateEmailProviderSMTPRequest_Xoauth2:
		cmd.XOAuth2Auth = xoauth2ToCmd(v)
	default:
		// ensure backwards compatibility
		if req.User != "" || req.Password != "" {
			cmd.PlainAuth = &command.PlainAuth{
				Password: req.Password,
			}
		}
	}

	return cmd
}

func xoauth2ToCmd(xoauth2 *admin_pb.UpdateEmailProviderSMTPRequest_Xoauth2) *command.XOAuth2Auth {
	cmd := &command.XOAuth2Auth{
		TokenEndpoint: xoauth2.Xoauth2.TokenEndpoint,
		Scopes:        xoauth2.Xoauth2.Scopes,
	}
	switch v := xoauth2.Xoauth2.OAuth2Type.(type) {
	case *admin_pb.SMTPXOAuth2Auth_ClientCredentials_:
		cmd.ClientCredentialsAuth = &command.OAuth2ClientCredentials{
			ClientId:     v.ClientCredentials.ClientId,
			ClientSecret: v.ClientCredentials.ClientSecret,
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
			User:     req.User,
			Password: v.Plain.Password,
		}
	case *admin_pb.TestEmailProviderSMTPRequest_Xoauth2:
		cfg.SMTP.XOAuth2Auth = xoauth2ToSmtp(v, req.User)
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

func xoauth2ToSmtp(xoauth2 *admin_pb.TestEmailProviderSMTPRequest_Xoauth2, user string) *smtp.XOAuth2AuthConfig {
	cfg := &smtp.XOAuth2AuthConfig{
		User:          user,
		TokenEndpoint: xoauth2.Xoauth2.TokenEndpoint,
		Scopes:        xoauth2.Xoauth2.Scopes,
	}
	switch v := xoauth2.Xoauth2.OAuth2Type.(type) {
	case *admin_pb.SMTPXOAuth2Auth_ClientCredentials_:
		cfg.ClientCredentialsAuth = &smtp.OAuth2ClientCredentials{
			ClientId:     v.ClientCredentials.ClientId,
			ClientSecret: v.ClientCredentials.ClientSecret,
		}
	}
	return cfg
}
