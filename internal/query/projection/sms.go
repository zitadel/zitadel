package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
)

type SMSConfigProjection struct {
	crdb.StatementHandler
}

const (
	SMSConfigProjectionTable = "zitadel.projections.sms_configs"
	SMSTwilioTable           = SMSConfigProjectionTable + "_" + smsTwilioTableSuffix
)

func NewSMSConfigProjection(ctx context.Context, config crdb.StatementHandlerConfig) *SMSConfigProjection {
	p := &SMSConfigProjection{}
	config.ProjectionName = SMSConfigProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *SMSConfigProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.SMSConfigTwilioAddedEventType,
					Reduce: p.reduceSMSConfigTwilioAdded,
				},
				{
					Event:  iam.SMSConfigTwilioChangedEventType,
					Reduce: p.reduceSMSConfigTwilioChanged,
				},
				{
					Event:  iam.SMSConfigActivatedEventType,
					Reduce: p.reduceSMSConfigActivated,
				},
				{
					Event:  iam.SMSConfigDeactivatedEventType,
					Reduce: p.reduceSMSConfigDeactivated,
				},
				{
					Event:  iam.SMSConfigRemovedEventType,
					Reduce: p.reduceSMSConfigRemoved,
				},
			},
		},
	}
}

const (
	SMSColumnID            = "id"
	SMSColumnAggregateID   = "aggregate_id"
	SMSColumnCreationDate  = "creation_date"
	SMSColumnChangeDate    = "change_date"
	SMSColumnResourceOwner = "resource_owner"
	SMSColumnState         = "state"
	SMSColumnSequence      = "sequence"

	smsTwilioTableSuffix            = "twilio"
	SMSTwilioConfigColumnSMSID      = "sms_id"
	SMSTwilioConfigColumnSID        = "sid"
	SMSTwilioConfigColumnToken      = "token"
	SMSTwilioConfigColumnSenderName = "sender_name"
)

func (p *SMSConfigProjection) reduceSMSConfigTwilioAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMSConfigTwilioAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-9jiWf", "seq", event.Sequence(), "expectedType", iam.SMSConfigTwilioAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-s8efs", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnID, e.ID),
				handler.NewCol(SMSColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMSColumnCreationDate, e.CreationDate()),
				handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMSColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMSColumnState, domain.SMSConfigStateInactive),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMSTwilioConfigColumnSMSID, e.ID),
				handler.NewCol(SMSTwilioConfigColumnSID, e.SID),
				handler.NewCol(SMSTwilioConfigColumnToken, e.Token),
				handler.NewCol(SMSTwilioConfigColumnSenderName, e.SenderName),
			},
			crdb.WithTableSuffix(smsTwilioTableSuffix),
		),
	), nil
}

func (p *SMSConfigProjection) reduceSMSConfigTwilioChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMSConfigTwilioChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-fm9el", "seq", event.Sequence(), "expectedType", iam.SMSConfigTwilioChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-fi99F", "reduce.wrong.event.type")
	}
	columns := make([]handler.Column, 0)
	if e.SID != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnSID, e.SID))
	}
	if e.SenderName != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnSenderName, e.SenderName))
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMSTwilioConfigColumnSMSID, e.ID),
			},
			crdb.WithTableSuffix(smsTwilioTableSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(SMSColumnID, e.ID),
			},
		),
	), nil
}

func (p *SMSConfigProjection) reduceSMSConfigActivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMSConfigActivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-fm03F", "seq", event.Sequence(), "expectedType", iam.SMSConfigActivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-fj9Ef", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMSColumnState, domain.SMSConfigStateActive),
			handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMSColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
		},
	), nil
}

func (p *SMSConfigProjection) reduceSMSConfigDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMSConfigDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-9fnHS", "seq", event.Sequence(), "expectedType", iam.SMSConfigDeactivatedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-dj9Js", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMSColumnState, domain.SMSConfigStateInactive),
			handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMSColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
		},
	), nil
}

func (p *SMSConfigProjection) reduceSMSConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMSConfigRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-0Opew", "seq", event.Sequence(), "expectedType", iam.SMSConfigRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-s9JJf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
		},
	), nil
}
