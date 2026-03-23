package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgSMTPConfigWriteModel struct {
	eventstore.WriteModel

	ID          string
	Description string

	SMTPConfig *SMTPConfig
	HTTPConfig *HTTPConfig

	State domain.SMTPConfigState
}

func NewOrgSMTPConfigWriteModel(orgID, id string) *OrgSMTPConfigWriteModel {
	return &OrgSMTPConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   orgID,
			ResourceOwner: orgID,
		},
		ID: id,
	}
}

func (wm *OrgSMTPConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch event.(type) {
		case *org.OrgSMTPConfigAddedEvent,
			*org.OrgSMTPConfigChangedEvent,
			*org.OrgSMTPConfigPasswordChangedEvent,
			*org.OrgSMTPConfigHTTPAddedEvent,
			*org.OrgSMTPConfigHTTPChangedEvent,
			*org.OrgSMTPConfigRemovedEvent,
			*org.OrgSMTPConfigActivatedEvent,
			*org.OrgSMTPConfigDeactivatedEvent:
			wm.WriteModel.AppendEvents(event)
		}
	}
}

func (wm *OrgSMTPConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.OrgSMTPConfigAddedEvent:
			if wm.ID != "" && wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigAddedEvent(e)
		case *org.OrgSMTPConfigChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigChangedEvent(e)
		case *org.OrgSMTPConfigPasswordChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.Password != nil {
				if wm.SMTPConfig.PlainAuth == nil {
					wm.SMTPConfig.PlainAuth = &instance.PlainAuth{Password: e.Password}
				} else {
					wm.SMTPConfig.PlainAuth.Password = e.Password
				}
			}
		case *org.OrgSMTPConfigHTTPAddedEvent:
			if wm.ID != "" && wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigHTTPAddedEvent(e)
		case *org.OrgSMTPConfigHTTPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigHTTPChangedEvent(e)
		case *org.OrgSMTPConfigRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigRemovedEvent()
		case *org.OrgSMTPConfigActivatedEvent:
			if wm.ID != e.ID {
				wm.State = domain.SMTPConfigStateInactive
				continue
			}
			wm.State = domain.SMTPConfigStateActive
		case *org.OrgSMTPConfigDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMTPConfigStateInactive
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrgSMTPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.OrgSMTPConfigAddedEventType,
			org.OrgSMTPConfigChangedEventType,
			org.OrgSMTPConfigPasswordChangedEventType,
			org.OrgSMTPConfigHTTPAddedEventType,
			org.OrgSMTPConfigHTTPChangedEventType,
			org.OrgSMTPConfigActivatedEventType,
			org.OrgSMTPConfigDeactivatedEventType,
			org.OrgSMTPConfigRemovedEventType,
		).
		Builder()
}

func (wm *OrgSMTPConfigWriteModel) NewChangedEvent(
	ctx context.Context, aggregate *eventstore.Aggregate,
	id,
	description string,
	tls bool,
	fromAddress,
	fromName,
	replyToAddress,
	smtpHost string,
	smtpUser string,
	plainAuth *instance.PlainAuth,
	xoauth2Auth *instance.XOAuth2Auth,
) (*org.OrgSMTPConfigChangedEvent, bool, error) {
	changes := make([]org.OrgSMTPConfigChanges, 0)
	if wm.SMTPConfig == nil {
		return nil, false, nil
	}

	if wm.ID != id {
		changes = append(changes, org.ChangeOrgSMTPConfigID(id))
	}
	if wm.Description != description {
		changes = append(changes, org.ChangeOrgSMTPConfigDescription(description))
	}
	if wm.SMTPConfig.TLS != tls {
		changes = append(changes, org.ChangeOrgSMTPConfigTLS(tls))
	}
	if wm.SMTPConfig.SenderAddress != fromAddress {
		changes = append(changes, org.ChangeOrgSMTPConfigFromAddress(fromAddress))
	}
	if wm.SMTPConfig.SenderName != fromName {
		changes = append(changes, org.ChangeOrgSMTPConfigFromName(fromName))
	}
	if wm.SMTPConfig.ReplyToAddress != replyToAddress {
		changes = append(changes, org.ChangeOrgSMTPConfigReplyToAddress(replyToAddress))
	}
	if wm.SMTPConfig.Host != smtpHost {
		changes = append(changes, org.ChangeOrgSMTPConfigSMTPHost(smtpHost))
	}
	if wm.SMTPConfig.User != smtpUser {
		changes = append(changes, org.ChangeOrgSMTPConfigSMTPUser(smtpUser))
	}
	if plainAuth != nil {
		changes = append(changes, orgSmtpPlainAuthChanges(wm.SMTPConfig.PlainAuth, *plainAuth)...)
	}
	if xoauth2Auth != nil {
		changes = append(changes, orgSmtpXOAuthChanges(wm.SMTPConfig.XOAuth2Auth, *xoauth2Auth)...)
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := org.NewOrgSMTPConfigChangeEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *OrgSMTPConfigWriteModel) NewHTTPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id, description, endpoint string,
	signingKey *crypto.CryptoValue,
) (*org.OrgSMTPConfigHTTPChangedEvent, bool, error) {
	changes := make([]org.OrgSMTPConfigHTTPChanges, 0)
	if wm.HTTPConfig == nil {
		return nil, false, nil
	}

	if wm.ID != id {
		changes = append(changes, org.ChangeOrgSMTPConfigHTTPID(id))
	}
	if wm.Description != description {
		changes = append(changes, org.ChangeOrgSMTPConfigHTTPDescription(description))
	}
	if wm.HTTPConfig.Endpoint != endpoint {
		changes = append(changes, org.ChangeOrgSMTPConfigHTTPEndpoint(endpoint))
	}
	if signingKey != nil {
		changes = append(changes, org.ChangeOrgSMTPConfigHTTPSigningKey(signingKey))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := org.NewOrgSMTPConfigHTTPChangeEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *OrgSMTPConfigWriteModel) reduceSMTPConfigAddedEvent(e *org.OrgSMTPConfigAddedEvent) {
	wm.ID = e.ID
	wm.Description = e.Description
	wm.SMTPConfig = &SMTPConfig{
		TLS:            e.TLS,
		Host:           e.Host,
		User:           e.User,
		SenderName:     e.SenderName,
		SenderAddress:  e.SenderAddress,
		ReplyToAddress: e.ReplyToAddress,
	}
	if e.PlainAuth != nil {
		wm.SMTPConfig.PlainAuth = &instance.PlainAuth{
			Password: e.PlainAuth.Password,
		}
	}
	if e.XOAuth2Auth != nil {
		wm.SMTPConfig.XOAuth2Auth = &instance.XOAuth2Auth{
			TokenEndpoint: e.XOAuth2Auth.TokenEndpoint,
			Scopes:        e.XOAuth2Auth.Scopes,
		}
		if e.XOAuth2Auth.ClientCredentials != nil {
			wm.SMTPConfig.XOAuth2Auth.ClientCredentials = &instance.XOAuth2ClientCredentials{
				ClientId:     e.XOAuth2Auth.ClientCredentials.ClientId,
				ClientSecret: e.XOAuth2Auth.ClientCredentials.ClientSecret,
			}
		}
	}
	wm.State = domain.SMTPConfigStateInactive
}

func (wm *OrgSMTPConfigWriteModel) reduceSMTPConfigHTTPAddedEvent(e *org.OrgSMTPConfigHTTPAddedEvent) {
	wm.ID = e.ID
	wm.Description = e.Description
	wm.HTTPConfig = &HTTPConfig{
		Endpoint:  e.Endpoint,
		SigningKey: e.SigningKey,
	}
	wm.State = domain.SMTPConfigStateInactive
}

func (wm *OrgSMTPConfigWriteModel) reduceSMTPConfigChangedEvent(e *org.OrgSMTPConfigChangedEvent) {
	if wm.SMTPConfig == nil {
		return
	}
	if e.Description != nil {
		wm.Description = *e.Description
	}
	if e.TLS != nil {
		wm.SMTPConfig.TLS = *e.TLS
	}
	if e.Host != nil {
		wm.SMTPConfig.Host = *e.Host
	}
	if e.User != nil {
		wm.SMTPConfig.User = *e.User
	}
	if e.FromAddress != nil {
		wm.SMTPConfig.SenderAddress = *e.FromAddress
	}
	if e.FromName != nil {
		wm.SMTPConfig.SenderName = *e.FromName
	}
	if e.ReplyToAddress != nil {
		wm.SMTPConfig.ReplyToAddress = *e.ReplyToAddress
	}
	if !e.PlainAuth.IsEmpty() {
		if wm.SMTPConfig.PlainAuth == nil {
			wm.SMTPConfig.PlainAuth = &instance.PlainAuth{}
			wm.SMTPConfig.XOAuth2Auth = nil
		}
		if e.PlainAuth.Password != nil {
			wm.SMTPConfig.PlainAuth.Password = e.PlainAuth.Password
		} else if e.Password != nil {
			wm.SMTPConfig.PlainAuth.Password = e.Password
		}
	}
	if !e.XOAuth2Auth.IsEmpty() {
		if wm.SMTPConfig.XOAuth2Auth == nil {
			wm.SMTPConfig.XOAuth2Auth = &instance.XOAuth2Auth{}
			wm.SMTPConfig.PlainAuth = nil
		}
		if e.XOAuth2Auth.TokenEndpoint != nil {
			wm.SMTPConfig.XOAuth2Auth.TokenEndpoint = *e.XOAuth2Auth.TokenEndpoint
		}
		if e.XOAuth2Auth.Scopes != nil {
			wm.SMTPConfig.XOAuth2Auth.Scopes = e.XOAuth2Auth.Scopes
		}
		if wm.SMTPConfig.XOAuth2Auth.ClientCredentials != nil && !e.XOAuth2Auth.ClientCredentials.IsEmpty() {
			if e.XOAuth2Auth.ClientCredentials.ClientId != nil {
				wm.SMTPConfig.XOAuth2Auth.ClientCredentials.ClientId = *e.XOAuth2Auth.ClientCredentials.ClientId
			}
			if e.XOAuth2Auth.ClientCredentials.ClientSecret != nil {
				wm.SMTPConfig.XOAuth2Auth.ClientCredentials.ClientSecret = e.XOAuth2Auth.ClientCredentials.ClientSecret
			}
		}
	}
}

func (wm *OrgSMTPConfigWriteModel) reduceSMTPConfigHTTPChangedEvent(e *org.OrgSMTPConfigHTTPChangedEvent) {
	if wm.HTTPConfig == nil {
		return
	}
	if e.Description != nil {
		wm.Description = *e.Description
	}
	if e.Endpoint != nil {
		wm.HTTPConfig.Endpoint = *e.Endpoint
	}
	if e.SigningKey != nil {
		wm.HTTPConfig.SigningKey = e.SigningKey
	}
}

func (wm *OrgSMTPConfigWriteModel) reduceSMTPConfigRemovedEvent() {
	wm.Description = ""
	wm.HTTPConfig = nil
	wm.SMTPConfig = nil
	wm.State = domain.SMTPConfigStateRemoved
}

func orgSmtpPlainAuthChanges(wm *instance.PlainAuth, auth instance.PlainAuth) []org.OrgSMTPConfigChanges {
	if wm == nil {
		return []org.OrgSMTPConfigChanges{
			org.ChangeOrgSMTPConfigSMTPPassword(auth.Password),
		}
	}
	var changes []org.OrgSMTPConfigChanges
	if auth.Password != nil {
		changes = append(changes, org.ChangeOrgSMTPConfigSMTPPassword(auth.Password))
	}
	return changes
}

func orgSmtpXOAuthChanges(wm *instance.XOAuth2Auth, auth instance.XOAuth2Auth) []org.OrgSMTPConfigChanges {
	if wm == nil {
		return []org.OrgSMTPConfigChanges{
			org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientId(auth.ClientCredentials.ClientId),
			org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientSecret(auth.ClientCredentials.ClientSecret),
			org.ChangeOrgSMTPConfigXOAuth2TokenEndpoint(auth.TokenEndpoint),
			org.ChangeOrgSMTPConfigXOAuth2Scopes(auth.Scopes),
		}
	}
	var changes []org.OrgSMTPConfigChanges
	if wm.TokenEndpoint != auth.TokenEndpoint {
		changes = append(changes, org.ChangeOrgSMTPConfigXOAuth2TokenEndpoint(auth.TokenEndpoint))
	}
	if len(wm.Scopes) != len(auth.Scopes) {
		changes = append(changes, org.ChangeOrgSMTPConfigXOAuth2Scopes(auth.Scopes))
	} else {
		for _, s := range auth.Scopes {
			if !slices.Contains(wm.Scopes, s) {
				changes = append(changes, org.ChangeOrgSMTPConfigXOAuth2Scopes(auth.Scopes))
				break
			}
		}
	}
	if auth.ClientCredentials != nil {
		changes = append(changes, orgSmtpXOAuthClientCredentialChanges(auth.ClientCredentials, *auth.ClientCredentials)...)
	}
	return changes
}

func orgSmtpXOAuthClientCredentialChanges(wm *instance.XOAuth2ClientCredentials, cc instance.XOAuth2ClientCredentials) []org.OrgSMTPConfigChanges {
	if wm == nil {
		return []org.OrgSMTPConfigChanges{
			org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientId(cc.ClientId),
			org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientSecret(cc.ClientSecret),
		}
	}
	var changes []org.OrgSMTPConfigChanges
	if wm.ClientId != cc.ClientId {
		changes = append(changes, org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientId(cc.ClientId))
	}
	if wm.ClientSecret != cc.ClientSecret {
		changes = append(changes, org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientSecret(cc.ClientSecret))
	}
	return changes
}
