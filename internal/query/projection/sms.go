package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	SMSConfigProjectionTable = "projections.sms_configs3"
	SMSTwilioTable           = SMSConfigProjectionTable + "_" + smsTwilioTableSuffix
	SMSHTTPTable             = SMSConfigProjectionTable + "_" + smsHTTPTableSuffix

	SMSColumnID            = "id"
	SMSColumnAggregateID   = "aggregate_id"
	SMSColumnCreationDate  = "creation_date"
	SMSColumnChangeDate    = "change_date"
	SMSColumnSequence      = "sequence"
	SMSColumnState         = "state"
	SMSColumnResourceOwner = "resource_owner"
	SMSColumnInstanceID    = "instance_id"
	SMSColumnDescription   = "description"

	smsTwilioTableSuffix            = "twilio"
	SMSTwilioColumnSMSID            = "sms_id"
	SMSTwilioColumnInstanceID       = "instance_id"
	SMSTwilioColumnSID              = "sid"
	SMSTwilioColumnSenderNumber     = "sender_number"
	SMSTwilioColumnToken            = "token"
	SMSTwilioColumnVerifyServiceSID = "verify_service_sid"

	smsHTTPTableSuffix      = "http"
	SMSHTTPColumnSMSID      = "sms_id"
	SMSHTTPColumnInstanceID = "instance_id"
	SMSHTTPColumnEndpoint   = "endpoint"
	SMSHTTPColumnSigningKey = "signing_key"
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
			handler.NewColumn(SMSColumnDescription, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(SMSColumnInstanceID, SMSColumnID),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(SMSTwilioColumnSMSID, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioColumnSID, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioColumnSenderNumber, handler.ColumnTypeText),
			handler.NewColumn(SMSTwilioColumnToken, handler.ColumnTypeJSONB),
			handler.NewColumn(SMSTwilioColumnVerifyServiceSID, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(SMSTwilioColumnInstanceID, SMSTwilioColumnSMSID),
			smsTwilioTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(SMSHTTPColumnSMSID, handler.ColumnTypeText),
			handler.NewColumn(SMSHTTPColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMSHTTPColumnEndpoint, handler.ColumnTypeText),
			handler.NewColumn(SMSHTTPColumnSigningKey, handler.ColumnTypeJSONB, handler.Nullable()),
		},
			handler.NewPrimaryKey(SMSHTTPColumnInstanceID, SMSHTTPColumnSMSID),
			smsHTTPTableSuffix,
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
					Event:  instance.SMSConfigHTTPAddedEventType,
					Reduce: p.reduceSMSConfigHTTPAdded,
				},
				{
					Event:  instance.SMSConfigHTTPChangedEventType,
					Reduce: p.reduceSMSConfigHTTPChanged,
				},
				{
					Event:  instance.SMSConfigTwilioActivatedEventType,
					Reduce: p.reduceSMSConfigTwilioActivated,
				},
				{
					Event:  instance.SMSConfigTwilioDeactivatedEventType,
					Reduce: p.reduceSMSConfigTwilioDeactivated,
				},
				{
					Event:  instance.SMSConfigTwilioRemovedEventType,
					Reduce: p.reduceSMSConfigTwilioRemoved,
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
	e, err := assertEvent[*instance.SMSConfigTwilioAddedEvent](event)
	if err != nil {
		return nil, err
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
				handler.NewCol(SMSColumnDescription, e.Description),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMSTwilioColumnSMSID, e.ID),
				handler.NewCol(SMSTwilioColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMSTwilioColumnSID, e.SID),
				handler.NewCol(SMSTwilioColumnToken, e.Token),
				handler.NewCol(SMSTwilioColumnSenderNumber, e.SenderNumber),
				handler.NewCol(SMSTwilioColumnVerifyServiceSID, e.VerifyServiceSID),
			},
			handler.WithTableSuffix(smsTwilioTableSuffix),
		),
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigTwilioChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigTwilioChangedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
		handler.NewCol(SMSColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(SMSColumnDescription, *e.Description))
	}
	if len(columns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMSColumnID, e.ID),
				handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
	}

	twilioColumns := make([]handler.Column, 0, 3)
	if e.SID != nil {
		twilioColumns = append(twilioColumns, handler.NewCol(SMSTwilioColumnSID, *e.SID))
	}
	if e.SenderNumber != nil {
		twilioColumns = append(twilioColumns, handler.NewCol(SMSTwilioColumnSenderNumber, *e.SenderNumber))
	}
	if e.VerifyServiceSID != nil {
		twilioColumns = append(twilioColumns, handler.NewCol(SMSTwilioColumnVerifyServiceSID, *e.VerifyServiceSID))
	}
	if len(twilioColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			twilioColumns,
			[]handler.Condition{
				handler.NewCond(SMSTwilioColumnSMSID, e.ID),
				handler.NewCond(SMSTwilioColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smsTwilioTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *smsConfigProjection) reduceSMSConfigTwilioTokenChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigTwilioTokenChangedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSTwilioColumnToken, e.Token),
			},
			[]handler.Condition{
				handler.NewCond(SMSTwilioColumnSMSID, e.ID),
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

func (p *smsConfigProjection) reduceSMSConfigHTTPAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigHTTPAddedEvent](event)
	if err != nil {
		return nil, err
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
				handler.NewCol(SMSColumnDescription, e.Description),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMSHTTPColumnSMSID, e.ID),
				handler.NewCol(SMSHTTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMSHTTPColumnEndpoint, e.Endpoint),
				handler.NewCol(SMSHTTPColumnSigningKey, e.SigningKey),
			},
			handler.WithTableSuffix(smsHTTPTableSuffix),
		),
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigHTTPChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigHTTPChangedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
		handler.NewCol(SMSColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(SMSColumnDescription, *e.Description))
	}
	stmts = append(stmts, handler.AddUpdateStatement(
		columns,
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
			handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
		},
	))

	httpColumns := make([]handler.Column, 0)
	if e.SigningKey != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMSHTTPColumnSigningKey, e.SigningKey))
	}
	if e.Endpoint != nil {
		httpColumns = append(httpColumns, handler.NewCol(SMSHTTPColumnEndpoint, *e.Endpoint))
	}
	if len(httpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			httpColumns,
			[]handler.Condition{
				handler.NewCond(SMSHTTPColumnSMSID, e.ID),
				handler.NewCond(SMSHTTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smsHTTPTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *smsConfigProjection) reduceSMSConfigTwilioActivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigTwilioActivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnState, domain.SMSConfigStateInactive),
				handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.Not(handler.NewCond(SMSColumnID, e.ID)),
				handler.NewCond(SMSColumnState, domain.SMSConfigStateActive),
				handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnState, domain.SMSConfigStateActive),
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

func (p *smsConfigProjection) reduceSMSConfigTwilioDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigTwilioDeactivatedEvent](event)
	if err != nil {
		return nil, err
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

func (p *smsConfigProjection) reduceSMSConfigTwilioRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigTwilioRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
			handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *smsConfigProjection) reduceSMSConfigActivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigActivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnState, domain.SMSConfigStateInactive),
				handler.NewCol(SMSColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMSColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.Not(handler.NewCond(SMSColumnID, e.ID)),
				handler.NewCond(SMSColumnState, domain.SMSConfigStateActive),
				handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMSColumnState, domain.SMSConfigStateActive),
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

func (p *smsConfigProjection) reduceSMSConfigDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMSConfigDeactivatedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*instance.SMSConfigRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMSColumnID, e.ID),
			handler.NewCond(SMSColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
