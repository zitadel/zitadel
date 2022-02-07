package command

import (
	"context"

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
	SID   string
	Token string
	From  string
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
				SID:   e.SID,
				Token: e.Token,
				From:  e.From,
			}
			wm.State = domain.SMSConfigStateInactive
		case *iam.SMSConfigTwilioChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.SID != nil {
				wm.Twilio.SID = *e.SID
			}
			if e.Token != nil {
				wm.Twilio.Token = *e.Token
			}
			if e.From != nil {
				wm.Twilio.From = *e.From
			}
		case *iam.SMSConfigTwilioActivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateActive
		case *iam.SMSConfigTwilioDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateInactive
		case *iam.SMSConfigTwilioRemovedEvent:
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
			iam.SMSConfigTwilioActivatedEventType,
			iam.SMSConfigTwilioDeactivatedEventType,
			iam.SMSConfigTwilioRemovedEventType).
		Builder()
}

func (wm *IAMSMSConfigWriteModel) NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id, sid, token, from string) (*iam.SMSConfigTwilioChangedEvent, bool, error) {
	changes := make([]iam.SMSConfigTwilioChanges, 0)
	var err error

	if wm.Twilio.SID != sid {
		changes = append(changes, iam.ChangeSMSConfigTwilioSID(sid))
	}
	if wm.Twilio.Token != token {
		changes = append(changes, iam.ChangeSMSConfigTwilioToken(token))
	}
	if wm.Twilio.From != from {
		changes = append(changes, iam.ChangeSMSConfigTwilioFrom(from))
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
