package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	EmailWebhookConfigProjectionTable = "projections.email_webhook_configs"

	EmailWebhookConfigColumnAggregateID          = "aggregate_id"
	EmailWebhookConfigColumnCreationDate         = "creation_date"
	EmailWebhookConfigColumnChangeDate           = "change_date"
	EmailWebhookConfigColumnSequence             = "sequence"
	EmailWebhookConfigColumnResourceOwner        = "resource_owner"
	EmailWebhookConfigColumnInstanceID           = "instance_id"
	EmailWebhookConfigColumnTLS                  = "tls"
	EmailWebhookConfigColumnSenderAddress        = "sender_address"
	EmailWebhookConfigColumnSenderName           = "sender_name"
	EmailWebhookConfigColumnReplyToAddress       = "reply_to_address"
	EmailWebhookConfigColumnEmailWebhookHost     = "host"
	EmailWebhookConfigColumnEmailWebhookUser     = "username"
	EmailWebhookConfigColumnEmailWebhookPassword = "password"
)

type emailWebhookConfigProjection struct {
	crdb.StatementHandler
}

func newEmailWebhookConfigProjection(ctx context.Context, config crdb.StatementHandlerConfig) *emailWebhookConfigProjection {
	p := new(emailWebhookConfigProjection)
	config.ProjectionName = EmailWebhookConfigProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(EmailWebhookConfigColumnAggregateID, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(EmailWebhookConfigColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(EmailWebhookConfigColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(EmailWebhookConfigColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnTLS, crdb.ColumnTypeBool),
			crdb.NewColumn(EmailWebhookConfigColumnSenderAddress, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnSenderName, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnReplyToAddress, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnEmailWebhookHost, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnEmailWebhookUser, crdb.ColumnTypeText),
			crdb.NewColumn(EmailWebhookConfigColumnEmailWebhookPassword, crdb.ColumnTypeJSONB, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(EmailWebhookConfigColumnInstanceID, EmailWebhookConfigColumnAggregateID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *emailWebhookConfigProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.EmailWebhookConfigAddedEventType,
					Reduce: p.reduceEmailWebhookConfigAdded,
				},
				{
					Event:  instance.EmailWebhookConfigChangedEventType,
					Reduce: p.reduceEmailWebhookConfigChanged,
				},
				{
					Event:  instance.EmailWebhookConfigPasswordChangedEventType,
					Reduce: p.reduceEmailWebhookConfigPasswordChanged,
				},
				{
					Event:  instance.EmailWebhookConfigRemovedEventType,
					Reduce: p.reduceEmailWebhookConfigRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(EmailWebhookConfigColumnInstanceID),
				},
			},
		},
	}
}

func (p *emailWebhookConfigProjection) reduceEmailWebhookConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.EmailWebhookConfigAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-sk99F", "reduce.wrong.event.type %s", instance.EmailWebhookConfigAddedEventType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(EmailWebhookConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(EmailWebhookConfigColumnCreationDate, e.CreationDate()),
			handler.NewCol(EmailWebhookConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(EmailWebhookConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(EmailWebhookConfigColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(EmailWebhookConfigColumnSequence, e.Sequence()),
			handler.NewCol(EmailWebhookConfigColumnTLS, e.TLS),
			handler.NewCol(EmailWebhookConfigColumnSenderAddress, e.SenderAddress),
			handler.NewCol(EmailWebhookConfigColumnSenderName, e.SenderName),
			handler.NewCol(EmailWebhookConfigColumnReplyToAddress, e.ReplyToAddress),
			handler.NewCol(EmailWebhookConfigColumnEmailWebhookHost, e.Host),
			handler.NewCol(EmailWebhookConfigColumnEmailWebhookUser, e.User),
			handler.NewCol(EmailWebhookConfigColumnEmailWebhookPassword, e.Password),
		},
	), nil
}

func (p *emailWebhookConfigProjection) reduceEmailWebhookConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.EmailWebhookConfigChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-wl0wd", "reduce.wrong.event.type %s", instance.EmailWebhookConfigChangedEventType)
	}

	columns := make([]handler.Column, 0, 8)
	columns = append(columns, handler.NewCol(EmailWebhookConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(EmailWebhookConfigColumnSequence, e.Sequence()))
	if e.TLS != nil {
		columns = append(columns, handler.NewCol(EmailWebhookConfigColumnTLS, *e.TLS))
	}
	if e.FromAddress != nil {
		columns = append(columns, handler.NewCol(EmailWebhookConfigColumnSenderAddress, *e.FromAddress))
	}
	if e.FromName != nil {
		columns = append(columns, handler.NewCol(EmailWebhookConfigColumnSenderName, *e.FromName))
	}
	if e.ReplyToAddress != nil {
		columns = append(columns, handler.NewCol(EmailWebhookConfigColumnReplyToAddress, *e.ReplyToAddress))
	}
	if e.Host != nil {
		columns = append(columns, handler.NewCol(EmailWebhookConfigColumnEmailWebhookHost, *e.Host))
	}
	if e.User != nil {
		columns = append(columns, handler.NewCol(EmailWebhookConfigColumnEmailWebhookUser, *e.User))
	}
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(EmailWebhookConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(EmailWebhookConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *emailWebhookConfigProjection) reduceEmailWebhookConfigPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.EmailWebhookConfigPasswordChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fk02f", "reduce.wrong.event.type %s", instance.EmailWebhookConfigChangedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(EmailWebhookConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(EmailWebhookConfigColumnSequence, e.Sequence()),
			handler.NewCol(EmailWebhookConfigColumnEmailWebhookPassword, e.Password),
		},
		[]handler.Condition{
			handler.NewCond(EmailWebhookConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(EmailWebhookConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *emailWebhookConfigProjection) reduceEmailWebhookConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.EmailWebhookConfigRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(EmailWebhookConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(EmailWebhookConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
