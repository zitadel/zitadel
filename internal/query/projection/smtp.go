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
	SMTPConfigProjectionTable = "projections.smtp_configs5"
	SMTPConfigTable           = SMTPConfigProjectionTable + "_" + smtpConfigSMTPTableSuffix
	SMTPConfigHTTPTable       = SMTPConfigProjectionTable + "_" + smtpConfigHTTPTableSuffix

	SMTPConfigColumnInstanceID    = "instance_id"
	SMTPConfigColumnResourceOwner = "resource_owner"
	SMTPConfigColumnAggregateID   = "aggregate_id"
	SMTPConfigColumnID            = "id"
	SMTPConfigColumnCreationDate  = "creation_date"
	SMTPConfigColumnChangeDate    = "change_date"
	SMTPConfigColumnSequence      = "sequence"
	SMTPConfigColumnState         = "state"
	SMTPConfigColumnDescription   = "description"

	smtpConfigSMTPTableSuffix          = "smtp"
	SMTPConfigSMTPColumnInstanceID     = "instance_id"
	SMTPConfigSMTPColumnID             = "id"
	SMTPConfigSMTPColumnTLS            = "tls"
	SMTPConfigSMTPColumnSenderAddress  = "sender_address"
	SMTPConfigSMTPColumnSenderName     = "sender_name"
	SMTPConfigSMTPColumnReplyToAddress = "reply_to_address"
	SMTPConfigSMTPColumnHost           = "host"
	SMTPConfigSMTPColumnUser           = "username"
	SMTPConfigSMTPColumnPassword       = "password"

	smtpConfigHTTPTableSuffix      = "http"
	SMTPConfigHTTPColumnInstanceID = "instance_id"
	SMTPConfigHTTPColumnID         = "id"
	SMTPConfigHTTPColumnEndpoint   = "endpoint"
	SMTPConfigHTTPColumnSigningKey = "signing_key"
)

type smtpConfigProjection struct{}

func newSMTPConfigProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(smtpConfigProjection))
}

func (*smtpConfigProjection) Name() string {
	return SMTPConfigProjectionTable
}

func (*smtpConfigProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SMTPConfigColumnID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SMTPConfigColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SMTPConfigColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SMTPConfigColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnDescription, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigColumnState, handler.ColumnTypeEnum),
		},
			handler.NewPrimaryKey(SMTPConfigColumnInstanceID, SMTPConfigColumnID),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(SMTPConfigSMTPColumnID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnTLS, handler.ColumnTypeBool),
			handler.NewColumn(SMTPConfigSMTPColumnSenderAddress, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnSenderName, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnReplyToAddress, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnHost, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnUser, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigSMTPColumnPassword, handler.ColumnTypeJSONB, handler.Nullable()),
		},
			handler.NewPrimaryKey(SMTPConfigSMTPColumnInstanceID, SMTPConfigSMTPColumnID),
			smtpConfigSMTPTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(SMTPConfigHTTPColumnID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigHTTPColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigHTTPColumnEndpoint, handler.ColumnTypeText),
			handler.NewColumn(SMTPConfigHTTPColumnSigningKey, handler.ColumnTypeJSONB, handler.Nullable()),
		},
			handler.NewPrimaryKey(SMTPConfigHTTPColumnInstanceID, SMTPConfigHTTPColumnID),
			smtpConfigHTTPTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
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
					Event:  instance.SMTPConfigHTTPAddedEventType,
					Reduce: p.reduceSMTPConfigHTTPAdded,
				},
				{
					Event:  instance.SMTPConfigHTTPChangedEventType,
					Reduce: p.reduceSMTPConfigHTTPChanged,
				},
				{
					Event:  instance.SMTPConfigActivatedEventType,
					Reduce: p.reduceSMTPConfigActivated,
				},
				{
					Event:  instance.SMTPConfigDeactivatedEventType,
					Reduce: p.reduceSMTPConfigDeactivated,
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
	e, err := assertEvent[*instance.SMTPConfigAddedEvent](event)
	if err != nil {
		return nil, err
	}

	description := e.Description
	state := domain.SMTPConfigStateInactive
	if e.ID == "" {
		description = "generic"
		state = domain.SMTPConfigStateActive
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigColumnCreationDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, state),
				handler.NewCol(SMTPConfigColumnDescription, description),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigSMTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigSMTPColumnTLS, e.TLS),
				handler.NewCol(SMTPConfigSMTPColumnSenderAddress, e.SenderAddress),
				handler.NewCol(SMTPConfigSMTPColumnSenderName, e.SenderName),
				handler.NewCol(SMTPConfigSMTPColumnReplyToAddress, e.ReplyToAddress),
				handler.NewCol(SMTPConfigSMTPColumnHost, e.Host),
				handler.NewCol(SMTPConfigSMTPColumnUser, e.User),
				handler.NewCol(SMTPConfigSMTPColumnPassword, e.Password),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		),
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigHTTPAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigHTTPAddedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigColumnCreationDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(SMTPConfigColumnAggregateID, e.Aggregate().ID),
				handler.NewCol(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
				handler.NewCol(SMTPConfigColumnDescription, e.Description),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(SMTPConfigHTTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCol(SMTPConfigHTTPColumnEndpoint, e.Endpoint),
				handler.NewCol(SMTPConfigHTTPColumnSigningKey, e.SigningKey),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		),
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigHTTPChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigHTTPChangedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnDescription, *e.Description))
	}
	stmts = append(stmts, handler.AddUpdateStatement(
		columns,
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	))

	smtpColumns := make([]handler.Column, 0, 1)
	if e.Endpoint != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigHTTPColumnEndpoint, *e.Endpoint))
	}
	if e.SigningKey != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigHTTPColumnSigningKey, e.SigningKey))
	}
	if len(smtpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			smtpColumns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigHTTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigHTTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigHTTPTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigChangedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, 3)
	columns := []handler.Column{
		handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
	}
	if e.Description != nil {
		columns = append(columns, handler.NewCol(SMTPConfigColumnDescription, *e.Description))
	}
	if len(columns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			columns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
	}

	smtpColumns := make([]handler.Column, 0, 7)
	if e.TLS != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnTLS, *e.TLS))
	}
	if e.FromAddress != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnSenderAddress, *e.FromAddress))
	}
	if e.FromName != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnSenderName, *e.FromName))
	}
	if e.ReplyToAddress != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnReplyToAddress, *e.ReplyToAddress))
	}
	if e.Host != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnHost, *e.Host))
	}
	if e.User != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnUser, *e.User))
	}
	if e.Password != nil {
		smtpColumns = append(smtpColumns, handler.NewCol(SMTPConfigSMTPColumnPassword, *e.Password))
	}
	if len(smtpColumns) > 0 {
		stmts = append(stmts, handler.AddUpdateStatement(
			smtpColumns,
			[]handler.Condition{
				handler.NewCond(SMTPConfigSMTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, stmts...), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigPasswordChangedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigSMTPColumnPassword, e.Password),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigSMTPColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigSMTPColumnInstanceID, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(smtpConfigSMTPTableSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigActivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigActivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
			},
			[]handler.Condition{
				handler.Not(handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate()))),
				handler.NewCond(SMTPConfigColumnState, domain.SMTPConfigStateActive),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
				handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
				handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateActive),
			},
			[]handler.Condition{
				handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
				handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
			},
		),
	), nil
}

func (p *smtpConfigProjection) reduceSMTPConfigDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.SMTPConfigDeactivatedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SMTPConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(SMTPConfigColumnSequence, e.Sequence()),
			handler.NewCol(SMTPConfigColumnState, domain.SMTPConfigStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
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
			handler.NewCond(SMTPConfigColumnID, getSMTPConfigID(e.ID, e.Aggregate())),
			handler.NewCond(SMTPConfigColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func getSMTPConfigID(id string, aggregate *eventstore.Aggregate) string {
	if id != "" {
		return id
	}
	// Deal with old and unique SMTP settings (empty ID)
	return aggregate.ResourceOwner
}
