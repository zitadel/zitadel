package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type IAMSMSConfigWriteModel struct {
	eventstore.WriteModel

	ID     string
	Twilio *TwilioConfig
	State  domain.SMSConfigState
}

type TwilioConfig struct {
	SID              string
	Token            *crypto.CryptoValue
	SenderNumber     string
	VerifyServiceSID string
}

func NewIAMSMSConfigWriteModel(instanceID, id string) *IAMSMSConfigWriteModel {
	return &IAMSMSConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
		ID: id,
	}
}

func (wm *IAMSMSConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.SMSConfigTwilioAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio = &TwilioConfig{
				SID:              e.SID,
				Token:            e.Token,
				SenderNumber:     e.SenderNumber,
				VerifyServiceSID: e.VerifyServiceSID,
			}
			wm.State = domain.SMSConfigStateInactive
		case *instance.SMSConfigTwilioChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.SID != nil {
				wm.Twilio.SID = *e.SID
			}
			if e.SenderNumber != nil {
				wm.Twilio.SenderNumber = *e.SenderNumber
			}
			if e.VerifyServiceSID != nil {
				wm.Twilio.VerifyServiceSID = *e.VerifyServiceSID
			}
		case *instance.SMSConfigTwilioTokenChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio.Token = e.Token
		case *instance.SMSConfigActivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateActive
		case *instance.SMSConfigDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateInactive
		case *instance.SMSConfigRemovedEvent:
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
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.SMSConfigTwilioAddedEventType,
			instance.SMSConfigTwilioChangedEventType,
			instance.SMSConfigTwilioTokenChangedEventType,
			instance.SMSConfigActivatedEventType,
			instance.SMSConfigDeactivatedEventType,
			instance.SMSConfigRemovedEventType).
		Builder()
}

func (wm *IAMSMSConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id, sid, senderNumber string, verifyServiceSID string) (*instance.SMSConfigTwilioChangedEvent, bool, error) {
	changes := make([]instance.SMSConfigTwilioChanges, 0)
	var err error

	if wm.Twilio.SID != sid {
		changes = append(changes, instance.ChangeSMSConfigTwilioSID(sid))
	}
	if wm.Twilio.SenderNumber != senderNumber {
		changes = append(changes, instance.ChangeSMSConfigTwilioSenderNumber(senderNumber))
	}
	if wm.Twilio.VerifyServiceSID != verifyServiceSID {
		changes = append(changes, instance.ChangeSMSConfigTwilioVerifyServiceSID(verifyServiceSID))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSMSConfigTwilioChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
