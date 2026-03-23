package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func listOrgEmailProvidersToModel(req *mgmt_pb.ListOrgEmailProvidersRequest) (*query.SMTPConfigsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.SMTPConfigsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}

func orgEmailProvidersToPb(configs []*query.SMTPConfig) []*settings_pb.EmailProvider {
	c := make([]*settings_pb.EmailProvider, len(configs))
	for i, config := range configs {
		c[i] = orgEmailProviderToProviderPb(config)
	}
	return c
}

func orgEmailProviderToProviderPb(config *query.SMTPConfig) *settings_pb.EmailProvider {
	return &settings_pb.EmailProvider{
		Details:     object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:          config.ID,
		Description: config.Description,
		State:       orgEmailProviderStateToPb(config.State),
		Config:      orgEmailProviderToPb(config),
	}
}

func orgEmailProviderStateToPb(state domain.SMTPConfigState) settings_pb.EmailProviderState {
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

func orgEmailProviderToPb(config *query.SMTPConfig) settings_pb.EmailConfig {
	if config.SMTPConfig != nil {
		return orgSmtpToPb(config.SMTPConfig)
	}
	if config.HTTPConfig != nil {
		return orgHttpToPb(config.HTTPConfig)
	}
	return nil
}

func orgHttpToPb(http *query.HTTP) *settings_pb.EmailProvider_Http {
	return &settings_pb.EmailProvider_Http{
		Http: &settings_pb.EmailProviderHTTP{
			Endpoint:   http.Endpoint,
			SigningKey: http.SigningKey,
		},
	}
}

func orgSmtpToPb(config *query.SMTP) *settings_pb.EmailProvider_Smtp {
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

	if config.XOAuth2Auth == nil && config.PlainAuth == nil {
		ret.Smtp.Auth = &settings_pb.EmailProviderSMTP_None{
			None: &settings_pb.SMTPNoAuth{},
		}
	}

	return ret
}

func addOrgEmailProviderSMTPToConfig(ctx context.Context, req *mgmt_pb.AddOrgEmailProviderSMTPRequest) *command.AddOrgSMTPConfig {
	cmd := &command.AddOrgSMTPConfig{
		ResourceOwner:  authz.GetCtxData(ctx).OrgID,
		Description:    req.Description,
		Tls:            req.Tls,
		From:           req.SenderAddress,
		FromName:       req.SenderName,
		ReplyToAddress: req.ReplyToAddress,
		Host:           req.Host,
		User:           req.User,
	}

	switch v := req.Auth.(type) {
	case *mgmt_pb.AddOrgEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain:
		cmd.PlainAuth = &command.PlainAuth{
			Password: v.Plain.Password,
		}
	case *mgmt_pb.AddOrgEmailProviderSMTPRequest_Xoauth2:
		if xoauth2, ok := v.Xoauth2.OAuth2Type.(*mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials_); ok {
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

func updateOrgEmailProviderSMTPToConfig(ctx context.Context, req *mgmt_pb.UpdateOrgEmailProviderSMTPRequest) *command.ChangeOrgSMTPConfig {
	cmd := &command.ChangeOrgSMTPConfig{
		ResourceOwner:  authz.GetCtxData(ctx).OrgID,
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
	case *mgmt_pb.UpdateOrgEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Plain:
		cmd.PlainAuth = &command.PlainAuth{
			Password: v.Plain.Password,
		}
	case *mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Xoauth2:
		cmd.XOAuth2Auth = orgXoauth2ToCmd(v)
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

func orgXoauth2ToCmd(xoauth2 *mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Xoauth2) *command.XOAuth2Auth {
	cmd := &command.XOAuth2Auth{
		TokenEndpoint: xoauth2.Xoauth2.TokenEndpoint,
		Scopes:        xoauth2.Xoauth2.Scopes,
	}
	if v, ok := xoauth2.Xoauth2.OAuth2Type.(*mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials_); ok {
		cmd.ClientCredentialsAuth = &command.OAuth2ClientCredentials{
			ClientId:     v.ClientCredentials.ClientId,
			ClientSecret: v.ClientCredentials.ClientSecret,
		}
	}
	return cmd
}

func addOrgEmailProviderHTTPToConfig(ctx context.Context, req *mgmt_pb.AddOrgEmailProviderHTTPRequest) *command.AddOrgSMTPConfigHTTP {
	return &command.AddOrgSMTPConfigHTTP{
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
		Description:   req.Description,
		Endpoint:      req.Endpoint,
	}
}

func updateOrgEmailProviderHTTPToConfig(ctx context.Context, req *mgmt_pb.UpdateOrgEmailProviderHTTPRequest) *command.ChangeOrgSMTPConfigHTTP {
	expirationSigningKey := req.GetExpirationSigningKey() != nil
	return &command.ChangeOrgSMTPConfigHTTP{
		ResourceOwner:        authz.GetCtxData(ctx).OrgID,
		ID:                   req.Id,
		Description:          req.Description,
		Endpoint:             req.Endpoint,
		ExpirationSigningKey: expirationSigningKey,
	}
}

func testOrgEmailProviderSMTPToConfig(req *mgmt_pb.TestOrgEmailProviderSMTPRequest) *smtp.Config {
	cfg := &smtp.Config{
		Tls:      req.Tls,
		From:     req.SenderAddress,
		FromName: req.SenderName,
		SMTP: smtp.SMTP{
			Host: req.Host,
		},
	}

	switch v := req.Auth.(type) {
	case *mgmt_pb.TestOrgEmailProviderSMTPRequest_None:
		// Nothing to do, no auth is required
	case *mgmt_pb.TestOrgEmailProviderSMTPRequest_Plain:
		cfg.SMTP.PlainAuth = &smtp.PlainAuthConfig{
			User:     req.User,
			Password: v.Plain.Password,
		}
	case *mgmt_pb.TestOrgEmailProviderSMTPRequest_Xoauth2:
		cfg.SMTP.XOAuth2Auth = orgXoauth2ToSmtp(v, req.User)
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

func orgXoauth2ToSmtp(xoauth2 *mgmt_pb.TestOrgEmailProviderSMTPRequest_Xoauth2, user string) *smtp.XOAuth2AuthConfig {
	cfg := &smtp.XOAuth2AuthConfig{
		User:          user,
		TokenEndpoint: xoauth2.Xoauth2.TokenEndpoint,
		Scopes:        xoauth2.Xoauth2.Scopes,
	}
	if v, ok := xoauth2.Xoauth2.OAuth2Type.(*mgmt_pb.OrgSMTPXOAuth2Auth_ClientCredentials_); ok {
		cfg.ClientCredentialsAuth = &smtp.OAuth2ClientCredentials{
			ClientId:     v.ClientCredentials.ClientId,
			ClientSecret: v.ClientCredentials.ClientSecret,
		}
	}
	return cfg
}
