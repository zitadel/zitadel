package projection

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
)

const (
	SMSConfigProjectionTable = "projections.sms_configs"
	SMSTwilioTable           = SMSConfigProjectionTable + "_" + smsTwilioTableSuffix

	SMSColumnID            = "id"
	SMSColumnAggregateID   = "aggregate_id"
	SMSColumnCreationDate  = "creation_date"
	SMSColumnChangeDate    = "change_date"
	SMSColumnSequence      = "sequence"
	SMSColumnState         = "state"
	SMSColumnResourceOwner = "resource_owner"
	SMSColumnInstanceID    = "instance_id"

	smsTwilioTableSuffix              = "twilio"
	SMSTwilioConfigColumnSMSID        = "sms_id"
	SMSTwilioConfigColumnSID          = "sid"
	SMSTwilioConfigColumnSenderNumber = "sender_number"
	SMSTwilioConfigColumnToken        = "token"
)

type SMSConfigProjection struct {
	crdb.StatementHandler
}

func NewSMSConfigProjection(ctx context.Context, config crdb.StatementHandlerConfig) *SMSConfigProjection {
	p := new(SMSConfigProjection)
	config.ProjectionName = SMSConfigProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(SMSColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(SMSColumnAggregateID, crdb.ColumnTypeText),
			crdb.NewColumn(SMSColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SMSColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SMSColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(SMSColumnState, crdb.ColumnTypeEnum),
			crdb.NewColumn(SMSColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(SMSColumnInstanceID, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(SMSColumnInstanceID, SMSColumnID),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(SMSTwilioConfigColumnSMSID, crdb.ColumnTypeText, crdb.Default(SMSColumnID)),
			crdb.NewColumn(SMSTwilioConfigColumnSID, crdb.ColumnTypeText),
			crdb.NewColumn(SMSTwilioConfigColumnSenderNumber, crdb.ColumnTypeText),
			crdb.NewColumn(SMSTwilioConfigColumnToken, crdb.ColumnTypeJSONB),
		},
			crdb.NewPrimaryKey(SMSTwilioConfigColumnSMSID),
			smsTwilioTableSuffix,
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *SMSConfigProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.SMSConfigTwilioAddedEventType,
					Reduce: p.reduceSMSConfigTwilioAdded,
				},
				{
					Event:  instance.SMSConfigTwilioChangedEventType,
					Reduce: p.reduceSMSConfigTwilioChanged,
				},
				{
					Event:  instance.SMSConfigActivatedEventType,
					Reduce: p.reduceSMSConfigActivated,
				},
				{
					Event:  instance.SMSConfigDeactivatedEventType,
					Reduce: p.reduceSMSConfigDeactivated,
				},
				{
					Event:  instance.SMSConfigRemovedEventType,
					Reduce: p.reduceSMSConfigRemoved,
				},
			},
		},
	}
}

func (p *SMSConfigProjection) reduceSMSConfigTwilioAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigTwilioAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-s8efs", "reduce.wrong.event.type %s", instance.SMSConfigTwilioAddedEventType)
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
				handler.NewCol(SMSColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMSColumnState, domain.SMSConfigStateInactive),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMSTwilioConfigColumnSMSID, e.ID),
				handler.NewCol(SMSTwilioConfigColumnSID, e.SID),
				handler.NewCol(SMSTwilioConfigColumnToken, e.Token),
				handler.NewCol(SMSTwilioConfigColumnSenderNumber, e.SenderNumber),
			},
			crdb.WithTableSuffix(smsTwilioTableSuffix),
		),
	), nil
}

func (p *SMSConfigProjection) reduceSMSConfigTwilioChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigTwilioChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fi99F", "reduce.wrong.event.type %s", instance.SMSConfigTwilioChangedEventType)
	}
	columns := make([]handler.Column, 0)
	if e.SID != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnSID, e.SID))
	}
	if e.SenderNumber != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnSenderNumber, e.SenderNumber))
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
	e, ok := event.(*instance.SMSConfigActivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fj9Ef", "reduce.wrong.event.type %s", instance.SMSConfigActivatedEventType)
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
	e, ok := event.(*instance.SMSConfigDeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-dj9Js", "reduce.wrong.event.type %s", instance.SMSConfigDeactivatedEventType)
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
	e, ok := event.(*instance.SMSConfigRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-s9JJf", "reduce.wrong.event.type %s", instance.SMSConfigRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
		},
	), nil
}
