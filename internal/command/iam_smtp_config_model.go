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

	FromAddress  string
	FromName     string
	TLS          bool
	SMTPHost     string
	SMTPUser     string
	SMTPPassword *crypto.CryptoValue
	State        domain.SMTPConfigState
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
			wm.FromAddress = e.FromAddress
			wm.FromName = e.FromName
			wm.SMTPHost = e.SMTPHost
			wm.SMTPUser = e.SMTPUser
			wm.SMTPPassword = e.SMTPPassword
			wm.State = domain.SMTPConfigStateActive
		case *iam.SMTPConfigChangedEvent:
			if e.TLS != nil {
				wm.TLS = *e.TLS
			}
			if e.FromAddress != nil {
				wm.FromAddress = *e.FromAddress
			}
			if e.FromName != nil {
				wm.FromName = *e.FromName
			}
			if e.SMTPHost != nil {
				wm.SMTPHost = *e.SMTPHost
			}
			if e.SMTPUser != nil {
				wm.SMTPUser = *e.SMTPUser
			}
			if e.SMTPPassword != nil {
				wm.SMTPPassword = e.SMTPPassword
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
			iam.SMTPConfigChangedEventType).
		Builder()
}

func (wm *IAMSMTPConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, tls bool, fromAddress, fromName, smtpHost, smtpUser, smtpPassword string, passwordCrypto crypto.EncryptionAlgorithm) (*iam.SMTPConfigChangedEvent, bool, error) {
	changes := make([]iam.SMTPConfigChanges, 0)
	var err error

	if wm.TLS != tls {
		changes = append(changes, iam.ChangeSMTPConfigTLS(tls))
	}
	if wm.FromAddress != fromAddress {
		changes = append(changes, iam.ChangeSMTPConfigFromAddress(fromAddress))
	}
	if wm.FromName != fromName {
		changes = append(changes, iam.ChangeSMTPConfigFromName(fromName))
	}
	if wm.SMTPHost != smtpHost {
		changes = append(changes, iam.ChangeSMTPConfigSMTPHost(smtpHost))
	}
	if wm.SMTPUser != smtpUser {
		changes = append(changes, iam.ChangeSMTPConfigSMTPUser(smtpUser))
	}

	existingPW, err := crypto.DecryptString(wm.SMTPPassword, passwordCrypto)
	if err != nil {
		return nil, false, err
	}
	if existingPW != smtpPassword {
		newPW, err := crypto.Encrypt([]byte(smtpPassword), passwordCrypto)
		if err != nil {
			return nil, false, err
		}
		changes = append(changes, iam.ChangeSMTPConfigSMTPPassword(newPW))
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
