package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	SMTPConfigProjectionTable = "projections.smtp_configs1"

	SMTPConfigColumnAggregateID    = "aggregate_id"
	SMTPConfigColumnCreationDate   = "creation_date"
	SMTPConfigColumnChangeDate     = "change_date"
	SMTPConfigColumnSequence       = "sequence"
	SMTPConfigColumnResourceOwner  = "resource_owner"
	SMTPConfigColumnInstanceID     = "instance_id"
	SMTPConfigColumnTLS            = "tls"
	SMTPConfigColumnSenderAddress  = "sender_address"
	SMTPConfigColumnSenderName     = "sender_name"
	SMTPConfigColumnReplyToAddress = "reply_to_address"
	SMTPConfigColumnSMTPHost       = "host"
	SMTPConfigColumnSMTPUser       = "username"
	SMTPConfigColumnSMTPPassword   = "password"
)

type smtpConfigProjection struct{}

func newSMTPConfigProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(smtpConfigProjection))
}

func (*smtpConfigProjection) Name() string {
	return SMTPConfigProjectionTable
}

func (*smtpConfigProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SMTPConfigColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SMTPConfigColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SMTPConfigColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SMTPConfigColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnTLS, handler.ColumnTypeBool),
			handler.NewColumn(SMTPConfigColumnSenderAddress, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnSenderName, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnReplyToAddress, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnSMTPHost, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnSMTPUser, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnSMTPPassword, handler.ColumnTypeJSONB, handler.Nullable()),
		},
			handler.NewPrimaryKey(SMTPConfigColumnInstanceID, SMTPConfigColumnAggregateID),
		),
	)
}

func (p *smtpConfigProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.SMTPConfigAddedEventType,
					Reduce: p.reduceSMTPConfigAdded,
				},
				{
					Event:  instance.SMTPConfigChangedEventType,
					Reduce: p.reduceSMTPConfigChanged,
				},
				{
					Event:  instance.SMTPConfigPasswordChangedEventType,
					Reduce: p.reduceSMTPConfigPasswordChanged,
				},
				{
					Event:  instance.SMTPConfigRemovedEventType,
					Reduce: p.reduceSMTPConfigRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SMTPConfigColumnInstanceID),
				},
			},
		},
	}
}

func (p *smtpConfigProjection) reduceSMTPConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMTPConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-sk99F", "reduce.wrong.event.type %s", instance.SMTPConfigAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(SMTPConfigColumnCreationDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			handler.NewCol(SMTPConfigColumnTLS, e.TLS),
			handler.NewCol(SMTPConfigColumnSenderAddress, e.SenderAddress),
			handler.NewCol(SMTPConfigColumnSenderName, e.SenderName),
			handler.NewCol(SMTPConfigColumnReplyToAddress, e.ReplyToAddress),
			handler.NewCol(SMTPConfigColumnSMTPHost, e.Host),
			handler.NewCol(SMTPConfigColumnSMTPUser, e.User),
			handler.NewCol(SMTPConfigColumnSMTPPassword, e.Password),
		},
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMTPConfigChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-wl0wd", "reduce.wrong.event.type %s", instance.SMTPConfigChangedEventType)
	}

	columns := make([]handler.Column, 0, 8)
	columns = append(columns, handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(SMTPConfigColumnSequence, e.Sequence()))
	if e.TLS != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnTLS, *e.TLS))
	}
	if e.FromAddress != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnSenderAddress, *e.FromAddress))
	}
	if e.FromName != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnSenderName, *e.FromName))
	}
	if e.ReplyToAddress != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnReplyToAddress, *e.ReplyToAddress))
	}
	if e.Host != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnSMTPHost, *e.Host))
	}
	if e.User != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnSMTPUser, *e.User))
	}
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SMTPConfigPasswordChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fk02f", "reduce.wrong.event.type %s", instance.SMTPConfigChangedEventType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			handler.NewCol(SMTPConfigColumnSMTPPassword, e.Password),
		},
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
