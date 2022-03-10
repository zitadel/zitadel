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

const (
	SMSConfigProjectionTable = "zitadel.projections.sms_configs"
	SMSTwilioTable           = SMSConfigProjectionTable + "_" + smsTwilioTableSuffix

	SMSColumnID            = "id"
	SMSColumnAggregateID   = "aggregate_id"
	SMSColumnCreationDate  = "creation_date"
	SMSColumnChangeDate    = "change_date"
	SMSColumnSequence      = "sequence"
	SMSColumnState         = "state"
	SMSColumnResourceOwner = "resource_owner"

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
	config.InitChecks = []*handler.Check{
		crdb.NewMultiTableCheck(
			crdb.NewTable([]*crdb.Column{
				crdb.NewColumn(SMSColumnID, crdb.ColumnTypeText),
				crdb.NewColumn(SMSColumnAggregateID, crdb.ColumnTypeText),
				crdb.NewColumn(SMSColumnCreationDate, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(SMSColumnChangeDate, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(SMSColumnSequence, crdb.ColumnTypeInt64),
				crdb.NewColumn(SMSColumnState, crdb.ColumnTypeEnum),
				crdb.NewColumn(SMSColumnResourceOwner, crdb.ColumnTypeText),
			},
				crdb.NewPrimaryKey(SMSColumnID),
			),
			crdb.NewSecondaryTable([]*crdb.Column{
				crdb.NewColumn(SMSTwilioConfigColumnSMSID, crdb.ColumnTypeText, crdb.Default(SMSColumnID)),
				crdb.NewColumn(SMSTwilioConfigColumnSID, crdb.ColumnTypeText),
				crdb.NewColumn(SMSTwilioConfigColumnSenderNumber, crdb.ColumnTypeText),
				crdb.NewColumn(SMSTwilioConfigColumnToken, crdb.ColumnTypeJSONB),
			},
				crdb.NewPrimaryKey(SMSTwilioConfigColumnSMSID),
				smsTwilioTableSuffix,
			),
		),
	}
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
				handler.NewCol(SMSTwilioConfigColumnSenderNumber, e.SenderNumber),
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
