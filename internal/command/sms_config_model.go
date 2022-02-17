package command

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMSMSConfigWriteModel struct {
	eventstore.WriteModel

	ID     string
	Twilio *TwilioConfig
	State  domain.SMSConfigState
}

type TwilioConfig struct {
	SID        string
	Token      *crypto.CryptoValue
	SenderName string
}

func NewIAMSMSConfigWriteModel(id string) *IAMSMSConfigWriteModel {
	return &IAMSMSConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   domain.IAMID,
			ResourceOwner: domain.IAMID,
		},
		ID: id,
	}
}

func (wm *IAMSMSConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *iam.SMSConfigTwilioAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio = &TwilioConfig{
				SID:        e.SID,
				Token:      e.Token,
				SenderName: e.SenderName,
			}
			wm.State = domain.SMSConfigStateInactive
		case *iam.SMSConfigTwilioChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.SID != nil {
				wm.Twilio.SID = *e.SID
			}
			if e.SenderName != nil {
				wm.Twilio.SenderName = *e.SenderName
			}
		case *iam.SMSConfigTwilioTokenChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio.Token = e.Token
		case *iam.SMSConfigActivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateActive
		case *iam.SMSConfigDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateInactive
		case *iam.SMSConfigRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio = nil
			wm.State = domain.SMSConfigStateRemoved
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
			iam.SMSConfigTwilioAddedEventType,
			iam.SMSConfigTwilioChangedEventType,
			iam.SMSConfigTwilioTokenChangedEventType,
			iam.SMSConfigActivatedEventType,
			iam.SMSConfigDeactivatedEventType,
			iam.SMSConfigRemovedEventType).
		Builder()
}

func (wm *IAMSMSConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id, sid, senderName string) (*iam.SMSConfigTwilioChangedEvent, bool, error) {
	changes := make([]iam.SMSConfigTwilioChanges, 0)
	var err error

	if wm.Twilio.SID != sid {
		changes = append(changes, iam.ChangeSMSConfigTwilioSID(sid))
	}
	if wm.Twilio.SenderName != senderName {
		changes = append(changes, iam.ChangeSMSConfigTwilioSenderName(senderName))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewSMSConfigTwilioChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
