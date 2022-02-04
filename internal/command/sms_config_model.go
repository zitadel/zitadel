package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMSMSConfigWriteModel struct {
	eventstore.WriteModel

	Twilio TwilioConfig
	State  domain.SMSConfigState
}

type TwilioConfig struct {
	SID   string
	Token string
	From  string
}

func NewIAMSMSConfigWriteModel() *IAMSMSConfigWriteModel {
	return &IAMSMSConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
	}
}

func (wm *IAMSMSConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.SMSConfigAddedEvent:
			wm.TLS = e.TLS
			wm.FromAddress = e.FromAddress
			wm.FromName = e.FromName
			wm.SMSHost = e.SMSHost
			wm.SMSUser = e.SMSUser
			wm.SMSPassword = e.SMSPassword
			wm.State = domain.SMSConfigStateActive
		case *iam.SMSConfigChangedEvent:
			if e.TLS != nil {
				wm.TLS = *e.TLS
			}
			if e.FromAddress != nil {
				wm.FromAddress = *e.FromAddress
			}
			if e.FromName != nil {
				wm.FromName = *e.FromName
			}
			if e.SMSHost != nil {
				wm.SMSHost = *e.SMSHost
			}
			if e.SMSUser != nil {
				wm.SMSUser = *e.SMSUser
			}
		}
	}
	return wm.WriteModel.Reduce()
}
func (wm *IAMSMSConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.SMSConfigAddedEventType,
			iam.SMSConfigChangedEventType,
			iam.SMSConfigPasswordChangedEventType).
		Builder()
}

func (wm *IAMSMSConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, tls bool, fromAddress, fromName, smtpHost, smtpUser string) (*iam.SMSConfigChangedEvent, bool, error) {
	changes := make([]iam.SMSConfigChanges, 0)
	var err error

	if wm.TLS != tls {
		changes = append(changes, iam.ChangeSMSConfigTLS(tls))
	}
	if wm.FromAddress != fromAddress {
		changes = append(changes, iam.ChangeSMSConfigFromAddress(fromAddress))
	}
	if wm.FromName != fromName {
		changes = append(changes, iam.ChangeSMSConfigFromName(fromName))
	}
	if wm.SMSHost != smtpHost {
		changes = append(changes, iam.ChangeSMSConfigSMSHost(smtpHost))
	}
	if wm.SMSUser != smtpUser {
		changes = append(changes, iam.ChangeSMSConfigSMSUser(smtpUser))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewSMSConfigChangeEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
