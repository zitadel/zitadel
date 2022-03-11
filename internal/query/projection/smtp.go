package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/project"
)

const (
	SMTPConfigProjectionTable = "projections.smtp_configs"

	SMTPConfigColumnAggregateID   = "aggregate_id"
	SMTPConfigColumnCreationDate  = "creation_date"
	SMTPConfigColumnChangeDate    = "change_date"
	SMTPConfigColumnSequence      = "sequence"
	SMTPConfigColumnResourceOwner = "resource_owner"
	SMTPConfigColumnTLS           = "tls"
	SMTPConfigColumnSenderAddress = "sender_address"
	SMTPConfigColumnSenderName    = "sender_name"
	SMTPConfigColumnSMTPHost      = "host"
	SMTPConfigColumnSMTPUser      = "username"
	SMTPConfigColumnSMTPPassword  = "password"
)

type SMTPConfigProjection struct {
	crdb.StatementHandler
}

func NewSMTPConfigProjection(ctx context.Context, config crdb.StatementHandlerConfig) *SMTPConfigProjection {
	p := new(SMTPConfigProjection)
	config.ProjectionName = SMTPConfigProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(SMTPConfigColumnAggregateID, crdb.ColumnTypeText),
			crdb.NewColumn(SMTPConfigColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SMTPConfigColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SMTPConfigColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(SMTPConfigColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(SMTPConfigColumnTLS, crdb.ColumnTypeBool),
			crdb.NewColumn(SMTPConfigColumnSenderAddress, crdb.ColumnTypeText),
			crdb.NewColumn(SMTPConfigColumnSenderName, crdb.ColumnTypeText),
			crdb.NewColumn(SMTPConfigColumnSMTPHost, crdb.ColumnTypeText),
			crdb.NewColumn(SMTPConfigColumnSMTPUser, crdb.ColumnTypeText),
			crdb.NewColumn(SMTPConfigColumnSMTPPassword, crdb.ColumnTypeJSONB, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(SMTPConfigColumnAggregateID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *SMTPConfigProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.SMTPConfigAddedEventType,
					Reduce: p.reduceSMTPConfigAdded,
				},
				{
					Event:  iam.SMTPConfigChangedEventType,
					Reduce: p.reduceSMTPConfigChanged,
				},
				{
					Event:  iam.SMTPConfigPasswordChangedEventType,
					Reduce: p.reduceSMTPConfigPasswordChanged,
				},
			},
		},
	}
}

func (p *SMTPConfigProjection) reduceSMTPConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMTPConfigAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-wkofs", "seq", event.Sequence(), "expectedType", iam.SMTPConfigAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-sk99F", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(SMTPConfigColumnCreationDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			handler.NewCol(SMTPConfigColumnTLS, e.TLS),
			handler.NewCol(SMTPConfigColumnSenderAddress, e.SenderAddress),
			handler.NewCol(SMTPConfigColumnSenderName, e.SenderName),
			handler.NewCol(SMTPConfigColumnSMTPHost, e.Host),
			handler.NewCol(SMTPConfigColumnSMTPUser, e.User),
			handler.NewCol(SMTPConfigColumnSMTPPassword, e.Password),
		},
	), nil
}

func (p *SMTPConfigProjection) reduceSMTPConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMTPConfigChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-wo00f", "seq", event.Sequence(), "expected", iam.SMTPConfigChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-wl0wd", "reduce.wrong.event.type")
	}

	columns := make([]handler.Column, 0, 7)
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
	if e.Host != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnSMTPHost, *e.Host))
	}
	if e.User != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnSMTPUser, *e.User))
	}
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnAggregateID, e.Aggregate().ID),
		},
	), nil
}

func (p *SMTPConfigProjection) reduceSMTPConfigPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SMTPConfigPasswordChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-f92sf", "seq", event.Sequence(), "expected", iam.SMTPConfigChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-fk02f", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			handler.NewCol(SMTPConfigColumnSMTPPassword, e.Password),
		},
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnAggregateID, e.Aggregate().ID),
		},
	), nil
}
