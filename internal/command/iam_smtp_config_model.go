package command

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMSMTPConfigWriteModel struct {
	eventstore.WriteModel

	SenderAddress string
	SenderName    string
	TLS           bool
	Host          string
	User          string
	Password      *crypto.CryptoValue
	State         domain.SMTPConfigState
}

func NewIAMSMTPConfigWriteModel() *IAMSMTPConfigWriteModel {
	return &IAMSMTPConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}

func (wm *IAMSMTPConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.SMTPConfigAddedEvent:
			wm.TLS = e.TLS
			wm.SenderAddress = e.SenderAddress
			wm.SenderName = e.SenderName
			wm.Host = e.Host
			wm.User = e.User
			wm.Password = e.Password
			wm.State = domain.SMTPConfigStateActive
		case *iam.SMTPConfigChangedEvent:
			if e.TLS != nil {
				wm.TLS = *e.TLS
			}
			if e.FromAddress != nil {
				wm.SenderAddress = *e.FromAddress
			}
			if e.FromName != nil {
				wm.SenderName = *e.FromName
			}
			if e.Host != nil {
				wm.Host = *e.Host
			}
			if e.User != nil {
				wm.User = *e.User
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IAMSMTPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.SMTPConfigAddedEventType,
			iam.SMTPConfigChangedEventType,
			iam.SMTPConfigPasswordChangedEventType).
		Builder()
}

func (wm *IAMSMTPConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, tls bool, fromAddress, fromName, smtpHost, smtpUser string) (*iam.SMTPConfigChangedEvent, bool, error) {
	changes := make([]iam.SMTPConfigChanges, 0)
	var err error

	if wm.TLS != tls {
		changes = append(changes, iam.ChangeSMTPConfigTLS(tls))
	}
	if wm.SenderAddress != fromAddress {
		changes = append(changes, iam.ChangeSMTPConfigFromAddress(fromAddress))
	}
	if wm.SenderName != fromName {
		changes = append(changes, iam.ChangeSMTPConfigFromName(fromName))
	}
	if wm.Host != smtpHost {
		changes = append(changes, iam.ChangeSMTPConfigSMTPHost(smtpHost))
	}
	if wm.User != smtpUser {
		changes = append(changes, iam.ChangeSMTPConfigSMTPUser(smtpUser))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewSMTPConfigChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
