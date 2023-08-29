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

	SenderAddress  string
	SenderName     string
	ReplyToAddress string
	TLS            bool
	Host           string
	User           string
	Password       *crypto.CryptoValue
	State          domain.SMTPConfigState

	domain                                 string
	domainState                            domain.InstanceDomainState
	smtpSenderAddressMatchesInstanceDomain bool
}

func NewInstanceSMTPConfigWriteModel(instanceID, domain string) *InstanceSMTPConfigWriteModel {
	return &InstanceSMTPConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
		},
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
			wm.TLS = e.TLS
			wm.SenderAddress = e.SenderAddress
			wm.SenderName = e.SenderName
			wm.ReplyToAddress = e.ReplyToAddress
			wm.Host = e.Host
			wm.User = e.User
			wm.Password = e.Password
			wm.State = domain.SMTPConfigStateActive
		case *instance.SMTPConfigChangedEvent:
			if e.TLS != nil {
				wm.TLS = *e.TLS
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
			if e.Host != nil {
				wm.Host = *e.Host
			}
			if e.User != nil {
				wm.User = *e.User
			}
		case *instance.SMTPConfigRemovedEvent:
			wm.State = domain.SMTPConfigStateRemoved
			wm.TLS = false
			wm.SenderName = ""
			wm.SenderAddress = ""
			wm.ReplyToAddress = ""
			wm.Host = ""
			wm.User = ""
			wm.Password = nil
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
			instance.SMTPConfigChangedEventType,
			instance.SMTPConfigPasswordChangedEventType,
			instance.InstanceDomainAddedEventType,
			instance.InstanceDomainRemovedEventType,
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType).
		Builder()
}

func (wm *InstanceSMTPConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, tls bool, fromAddress, fromName, replyToAddress, smtpHost, smtpUser string) (*instance.SMTPConfigChangedEvent, bool, error) {
	changes := make([]instance.SMTPConfigChanges, 0)
	var err error

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

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSMTPConfigChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
