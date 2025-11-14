package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type IAMSMTPConfigWriteModel struct {
	eventstore.WriteModel

	ID          string
	Description string

	SMTPConfig *SMTPConfig
	HTTPConfig *HTTPConfig

	State domain.SMTPConfigState

	domain                                 string
	domainState                            domain.InstanceDomainState
	smtpSenderAddressMatchesInstanceDomain bool
}

type SMTPConfig struct {
	TLS            bool
	Host           string
	User           string
	Password       *crypto.CryptoValue
	SenderAddress  string
	SenderName     string
	ReplyToAddress string
}

func NewIAMSMTPConfigWriteModel(instanceID, id, domain string) *IAMSMTPConfigWriteModel {
	return &IAMSMTPConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
		ID:     id,
		domain: domain,
	}
}

func (wm *IAMSMTPConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.DomainAddedEvent:
			if e.Domain != wm.domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *instance.DomainRemovedEvent:
			if e.Domain != wm.domain {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		default:
			wm.WriteModel.AppendEvents(e)
		}

	}
}

func (wm *IAMSMTPConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.SMTPConfigAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigAddedEvent(e)
		case *instance.SMTPConfigChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigChangedEvent(e)
		case *instance.SMTPConfigPasswordChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.Password != nil {
				wm.SMTPConfig.Password = e.Password
			}
		case *instance.SMTPConfigHTTPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigHTTPAddedEvent(e)
		case *instance.SMTPConfigHTTPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigHTTPChangedEvent(e)
		case *instance.SMTPConfigRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceSMTPConfigRemovedEvent(e)
		case *instance.SMTPConfigActivatedEvent:
			if wm.ID != e.ID {
				wm.State = domain.SMTPConfigStateInactive
				continue
			}
			wm.State = domain.SMTPConfigStateActive
		case *instance.SMTPConfigDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMTPConfigStateInactive
		case *instance.DomainAddedEvent:
			wm.domainState = domain.InstanceDomainStateActive
		case *instance.DomainRemovedEvent:
			wm.domainState = domain.InstanceDomainStateRemoved
		case *instance.DomainPolicyAddedEvent:
			wm.smtpSenderAddressMatchesInstanceDomain = e.SMTPSenderAddressMatchesInstanceDomain
		case *instance.DomainPolicyChangedEvent:
			if e.SMTPSenderAddressMatchesInstanceDomain != nil {
				wm.smtpSenderAddressMatchesInstanceDomain = *e.SMTPSenderAddressMatchesInstanceDomain
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IAMSMTPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	// If ID equals ResourceOwner we're dealing with the old and unique smtp settings
	// Let's set the empty ID for the query
	if wm.ID == wm.ResourceOwner {
		wm.ID = ""
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.SMTPConfigAddedEventType,
			instance.SMTPConfigRemovedEventType,
			instance.SMTPConfigChangedEventType,
			instance.SMTPConfigPasswordChangedEventType,
			instance.SMTPConfigHTTPAddedEventType,
			instance.SMTPConfigHTTPChangedEventType,
			instance.SMTPConfigActivatedEventType,
			instance.SMTPConfigDeactivatedEventType,
			instance.SMTPConfigRemovedEventType,
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType,
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType).
		Builder()
}

func (wm *IAMSMTPConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id, description string, tls bool, fromAddress, fromName, replyToAddress, smtpHost, smtpUser string, smtpPassword *crypto.CryptoValue) (*instance.SMTPConfigChangedEvent, bool, error) {
	changes := make([]instance.SMTPConfigChanges, 0)
	var err error
	if wm.SMTPConfig == nil {
		return nil, false, nil
	}

	if wm.ID != id {
		changes = append(changes, instance.ChangeSMTPConfigID(id))
	}
	if wm.Description != description {
		changes = append(changes, instance.ChangeSMTPConfigDescription(description))
	}
	if wm.SMTPConfig.TLS != tls {
		changes = append(changes, instance.ChangeSMTPConfigTLS(tls))
	}
	if wm.SMTPConfig.SenderAddress != fromAddress {
		changes = append(changes, instance.ChangeSMTPConfigFromAddress(fromAddress))
	}
	if wm.SMTPConfig.SenderName != fromName {
		changes = append(changes, instance.ChangeSMTPConfigFromName(fromName))
	}
	if wm.SMTPConfig.ReplyToAddress != replyToAddress {
		changes = append(changes, instance.ChangeSMTPConfigReplyToAddress(replyToAddress))
	}
	if wm.SMTPConfig.Host != smtpHost {
		changes = append(changes, instance.ChangeSMTPConfigSMTPHost(smtpHost))
	}
	if wm.SMTPConfig.User != smtpUser {
		changes = append(changes, instance.ChangeSMTPConfigSMTPUser(smtpUser))
	}
	if smtpPassword != nil {
		changes = append(changes, instance.ChangeSMTPConfigSMTPPassword(smtpPassword))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSMTPConfigChangeEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *IAMSMTPConfigWriteModel) NewHTTPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id, description, endpoint string,
	signingKey *crypto.CryptoValue,
) (*instance.SMTPConfigHTTPChangedEvent, bool, error) {
	changes := make([]instance.SMTPConfigHTTPChanges, 0)
	var err error
	if wm.HTTPConfig == nil {
		return nil, false, nil
	}

	if wm.ID != id {
		changes = append(changes, instance.ChangeSMTPConfigHTTPID(id))
	}
	if wm.Description != description {
		changes = append(changes, instance.ChangeSMTPConfigHTTPDescription(description))
	}
	if wm.HTTPConfig.Endpoint != endpoint {
		changes = append(changes, instance.ChangeSMTPConfigHTTPEndpoint(endpoint))
	}
	// if signingkey is set, update it as it is encrypted
	if signingKey != nil {
		changes = append(changes, instance.ChangeSMTPConfigHTTPSigningKey(signingKey))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSMTPConfigHTTPChangeEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *IAMSMTPConfigWriteModel) reduceSMTPConfigAddedEvent(e *instance.SMTPConfigAddedEvent) {
	wm.Description = e.Description
	wm.SMTPConfig = &SMTPConfig{
		TLS:            e.TLS,
		Host:           e.Host,
		User:           e.User,
		Password:       e.Password,
		SenderName:     e.SenderName,
		SenderAddress:  e.SenderAddress,
		ReplyToAddress: e.ReplyToAddress,
	}
	wm.State = domain.SMTPConfigStateInactive
	// If ID has empty value we're dealing with the old and unique smtp settings
	// These would be the default values for ID and State
	if e.ID == "" {
		wm.Description = "generic"
		wm.ID = e.Aggregate().ResourceOwner
		wm.State = domain.SMTPConfigStateActive
	}
}

func (wm *IAMSMTPConfigWriteModel) reduceSMTPConfigHTTPAddedEvent(e *instance.SMTPConfigHTTPAddedEvent) {
	wm.Description = e.Description
	wm.HTTPConfig = &HTTPConfig{
		Endpoint:   e.Endpoint,
		SigningKey: e.SigningKey,
	}
	wm.State = domain.SMTPConfigStateInactive
	// If ID has empty value we're dealing with the old and unique smtp settings
	// These would be the default values for ID and State
	if e.ID == "" {
		wm.Description = "generic"
		wm.ID = e.Aggregate().ResourceOwner
		wm.State = domain.SMTPConfigStateActive
	}
}

func (wm *IAMSMTPConfigWriteModel) reduceSMTPConfigChangedEvent(e *instance.SMTPConfigChangedEvent) {
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
	if e.Password != nil {
		wm.SMTPConfig.Password = e.Password
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

	// If ID has empty value we're dealing with the old and unique smtp settings
	// These would be the default values for ID and State
	if e.ID == "" {
		wm.Description = "generic"
		wm.ID = e.Aggregate().ResourceOwner
		wm.State = domain.SMTPConfigStateActive
	}
}

func (wm *IAMSMTPConfigWriteModel) reduceSMTPConfigHTTPChangedEvent(e *instance.SMTPConfigHTTPChangedEvent) {
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

	// If ID has empty value we're dealing with the old and unique smtp settings
	// These would be the default values for ID and State
	if e.ID == "" {
		wm.Description = "generic"
		wm.ID = e.Aggregate().ResourceOwner
		wm.State = domain.SMTPConfigStateActive
	}
}

func (wm *IAMSMTPConfigWriteModel) reduceSMTPConfigRemovedEvent(e *instance.SMTPConfigRemovedEvent) {
	wm.Description = ""
	wm.HTTPConfig = nil
	wm.SMTPConfig = nil
	wm.State = domain.SMTPConfigStateRemoved

	// If ID has empty value we're dealing with the old and unique smtp settings
	// These would be the default values for ID and State
	if e.ID == "" {
		wm.ID = e.Aggregate().ResourceOwner
	}
}
