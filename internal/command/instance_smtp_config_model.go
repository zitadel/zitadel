package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceSMTPConfigWriteModel struct {
	eventstore.WriteModel
	ID             string
	SenderAddress  string
	SenderName     string
	ReplyToAddress string
	TLS            bool
	Host           string
	User           string
	Password       *crypto.CryptoValue
	State          domain.SMTPConfigState
	ProviderType   uint32

	domain                                 string
	domainState                            domain.InstanceDomainState
	smtpSenderAddressMatchesInstanceDomain bool
}

func NewIAMSMTPConfigWriteModel(instanceID, id, domain string) *InstanceSMTPConfigWriteModel {
	return &InstanceSMTPConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
		ID:     id,
		domain: domain,
	}
}

func (wm *InstanceSMTPConfigWriteModel) AppendEvents(events ...eventstore.Event) {
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

func (wm *InstanceSMTPConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.SMTPConfigAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.TLS = e.TLS
			wm.Host = e.Host
			wm.User = e.User
			wm.Password = e.Password
			wm.SenderAddress = e.SenderAddress
			wm.SenderName = e.SenderName
			wm.ReplyToAddress = e.ReplyToAddress
			wm.ProviderType = e.ProviderType
			wm.State = domain.SMTPConfigStateInactive
		case *instance.SMTPConfigChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.TLS != nil {
				wm.TLS = *e.TLS
			}
			if e.Host != nil {
				wm.Host = *e.Host
			}
			if e.User != nil {
				wm.User = *e.User
			}
			if e.Password != nil {
				wm.Password = e.Password
			}
			if e.FromAddress != nil {
				wm.SenderAddress = *e.FromAddress
			}
			if e.FromName != nil {
				wm.SenderName = *e.FromName
			}
			if e.ReplyToAddress != nil {
				wm.ReplyToAddress = *e.ReplyToAddress
			}
			if e.ProviderType != nil {
				wm.ProviderType = *e.ProviderType
			}
		case *instance.SMTPConfigRemovedEvent:
			if wm.ID != e.ID {
				continue
			}

			wm.TLS = false
			wm.SenderName = ""
			wm.SenderAddress = ""
			wm.ReplyToAddress = ""
			wm.Host = ""
			wm.User = ""
			wm.Password = nil
			wm.ProviderType = 0
			wm.State = domain.SMTPConfigStateRemoved
		case *instance.SMTPConfigActivatedEvent:
			if wm.ID != e.ID {
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

func (wm *InstanceSMTPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
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
			instance.SMTPConfigActivatedEventType,
			instance.SMTPConfigDeactivatedEventType,
			instance.SMTPConfigRemovedEventType,
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType,
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType).
		Builder()
}

func (wm *InstanceSMTPConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id string, tls bool, fromAddress, fromName, replyToAddress, smtpHost, smtpUser string, smtpPassword *crypto.CryptoValue, providerType uint32) (*instance.SMTPConfigChangedEvent, bool, error) {
	changes := make([]instance.SMTPConfigChanges, 0)
	var err error

	if wm.ID != id {
		changes = append(changes, instance.ChangeSMTPConfigID(id))
	}

	if wm.TLS != tls {
		changes = append(changes, instance.ChangeSMTPConfigTLS(tls))
	}
	if wm.SenderAddress != fromAddress {
		changes = append(changes, instance.ChangeSMTPConfigFromAddress(fromAddress))
	}
	if wm.SenderName != fromName {
		changes = append(changes, instance.ChangeSMTPConfigFromName(fromName))
	}
	if wm.ReplyToAddress != replyToAddress {
		changes = append(changes, instance.ChangeSMTPConfigReplyToAddress(replyToAddress))
	}
	if wm.Host != smtpHost {
		changes = append(changes, instance.ChangeSMTPConfigSMTPHost(smtpHost))
	}
	if wm.User != smtpUser {
		changes = append(changes, instance.ChangeSMTPConfigSMTPUser(smtpUser))
	}
	if smtpPassword != nil {
		changes = append(changes, instance.ChangeSMTPConfigSMTPPassword(smtpPassword))
	}
	if wm.ProviderType != providerType {
		changes = append(changes, instance.ChangeSMTPConfigProviderType(providerType))
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
