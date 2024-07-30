package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SMSConfigProjectionTable = "projections.sms_configs2"
	SMSTwilioTable           = SMSConfigProjectionTable + "_" + smsTwilioTableSuffix

	SMSColumnID            = "id"
	SMSColumnAggregateID   = "aggregate_id"
	SMSColumnCreationDate  = "creation_date"
	SMSColumnChangeDate    = "change_date"
	SMSColumnSequence      = "sequence"
	SMSColumnState         = "state"
	SMSColumnResourceOwner = "resource_owner"
	SMSColumnInstanceID    = "instance_id"

	smsTwilioTableSuffix                  = "twilio"
	SMSTwilioConfigColumnSMSID            = "sms_id"
	SMSTwilioColumnInstanceID             = "instance_id"
	SMSTwilioConfigColumnSID              = "sid"
	SMSTwilioConfigColumnSenderNumber     = "sender_number"
	SMSTwilioConfigColumnToken            = "token"
	SMSTwilioConfigColumnVerifyServiceSID = "verify_service_sid"
)

type smsConfigProjection struct{}

func newSMSConfigProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(smsConfigProjection))
}

func (*smsConfigProjection) Name() string {
	return SMSConfigProjectionTable
}

func (*smsConfigProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SMSColumnID, handler.ColumnTypeText),
			handler.NewColumn(SMSColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(SMSColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SMSColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SMSColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SMSColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(SMSColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(SMSColumnInstanceID, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(SMSColumnInstanceID, SMSColumnID),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(SMSTwilioConfigColumnSMSID, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioConfigColumnSID, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioConfigColumnSenderNumber, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioConfigColumnToken, handler.ColumnTypeJSONB),
			handler.NewColumn(SMSTwilioConfigColumnVerifyServiceSID, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(SMSTwilioColumnInstanceID, SMSTwilioConfigColumnSMSID),
			smsTwilioTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
	)
}

func (p *smsConfigProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.SMSConfigTwilioAddedEventType,
					Reduce: p.reduceSMSConfigTwilioAdded,
				},
				{
					Event:  instance.SMSConfigTwilioChangedEventType,
					Reduce: p.reduceSMSConfigTwilioChanged,
				},
				{
					Event:  instance.SMSConfigTwilioTokenChangedEventType,
					Reduce: p.reduceSMSConfigTwilioTokenChanged,
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
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SMSColumnInstanceID),
				},
			},
		},
	}
}

func (p *smsConfigProjection) reduceSMSConfigTwilioAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigTwilioAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-s8efs", "reduce.wrong.event.type %s", instance.SMSConfigTwilioAddedEventType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
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
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMSTwilioConfigColumnSMSID, e.ID),
				handler.NewCol(SMSTwilioColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMSTwilioConfigColumnSID, e.SID),
				handler.NewCol(SMSTwilioConfigColumnToken, e.Token),
				handler.NewCol(SMSTwilioConfigColumnSenderNumber, e.SenderNumber),
				handler.NewCol(SMSTwilioConfigColumnVerifyServiceSID, e.VerifyServiceSID),
			},
			handler.WithTableSuffix(smsTwilioTableSuffix),
		),
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigTwilioChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigTwilioChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fi99F", "reduce.wrong.event.type %s", instance.SMSConfigTwilioChangedEventType)
	}
	columns := make([]handler.Column, 0)
	if e.SID != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnSID, *e.SID))
	}
	if e.SenderNumber != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnSenderNumber, *e.SenderNumber))
	}
	if e.VerifyServiceSID != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnVerifyServiceSID, *e.VerifyServiceSID))
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMSTwilioConfigColumnSMSID, e.ID),
				handler.NewCond(SMSTwilioColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smsTwilioTableSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(SMSColumnID, e.ID),
				handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigTwilioTokenChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigTwilioTokenChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fi99F", "reduce.wrong.event.type %s", instance.SMSConfigTwilioTokenChangedEventType)
	}
	columns := make([]handler.Column, 0)
	if e.Token != nil {
		columns = append(columns, handler.NewCol(SMSTwilioConfigColumnToken, e.Token))
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMSTwilioConfigColumnSMSID, e.ID),
				handler.NewCond(SMSTwilioColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smsTwilioTableSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(SMSColumnID, e.ID),
				handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigActivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigActivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fj9Ef", "reduce.wrong.event.type %s", instance.SMSConfigActivatedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMSColumnState, domain.SMSConfigStateActive),
			handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMSColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
			handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dj9Js", "reduce.wrong.event.type %s", instance.SMSConfigDeactivatedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMSColumnState, domain.SMSConfigStateInactive),
			handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMSColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
			handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMSConfigRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-s9JJf", "reduce.wrong.event.type %s", instance.SMSConfigRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
			handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
