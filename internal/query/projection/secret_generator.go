package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/project"
)

const (
	SecretGeneratorProjectionTable = "projections.secret_generators"

	SecretGeneratorColumnGeneratorType       = "generator_type"
	SecretGeneratorColumnAggregateID         = "aggregate_id"
	SecretGeneratorColumnCreationDate        = "creation_date"
	SecretGeneratorColumnChangeDate          = "change_date"
	SecretGeneratorColumnSequence            = "sequence"
	SecretGeneratorColumnResourceOwner       = "resource_owner"
	SecretGeneratorColumnInstanceID          = "instance_id"
	SecretGeneratorColumnLength              = "length"
	SecretGeneratorColumnExpiry              = "expiry"
	SecretGeneratorColumnIncludeLowerLetters = "include_lower_letters"
	SecretGeneratorColumnIncludeUpperLetters = "include_upper_letters"
	SecretGeneratorColumnIncludeDigits       = "include_digits"
	SecretGeneratorColumnIncludeSymbols      = "include_symbols"
)

type SecretGeneratorProjection struct {
	crdb.StatementHandler
}

func NewSecretGeneratorProjection(ctx context.Context, config crdb.StatementHandlerConfig) *SecretGeneratorProjection {
	p := new(SecretGeneratorProjection)
	config.ProjectionName = SecretGeneratorProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(SecretGeneratorColumnGeneratorType, crdb.ColumnTypeText),
			crdb.NewColumn(SecretGeneratorColumnAggregateID, crdb.ColumnTypeText),
			crdb.NewColumn(SecretGeneratorColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SecretGeneratorColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SecretGeneratorColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(SecretGeneratorColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(SecretGeneratorColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(SecretGeneratorColumnLength, crdb.ColumnTypeInt64),
			crdb.NewColumn(SecretGeneratorColumnExpiry, crdb.ColumnTypeInt64),
			crdb.NewColumn(SecretGeneratorColumnIncludeLowerLetters, crdb.ColumnTypeBool),
			crdb.NewColumn(SecretGeneratorColumnIncludeUpperLetters, crdb.ColumnTypeBool),
			crdb.NewColumn(SecretGeneratorColumnIncludeDigits, crdb.ColumnTypeBool),
			crdb.NewColumn(SecretGeneratorColumnIncludeSymbols, crdb.ColumnTypeBool),
		},
			crdb.NewPrimaryKey(SecretGeneratorColumnInstanceID, SecretGeneratorColumnGeneratorType, SecretGeneratorColumnAggregateID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *SecretGeneratorProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.SecretGeneratorAddedEventType,
					Reduce: p.reduceSecretGeneratorAdded,
				},
				{
					Event:  iam.SecretGeneratorChangedEventType,
					Reduce: p.reduceSecretGeneratorChanged,
				},
				{
					Event:  iam.SecretGeneratorRemovedEventType,
					Reduce: p.reduceSecretGeneratorRemoved,
				},
			},
		},
	}
}

func (p *SecretGeneratorProjection) reduceSecretGeneratorAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SecretGeneratorAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-sk99F", "reduce.wrong.event.type %s", iam.SecretGeneratorAddedEventType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SecretGeneratorColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(SecretGeneratorColumnGeneratorType, e.GeneratorType),
			handler.NewCol(SecretGeneratorColumnCreationDate, e.CreationDate()),
			handler.NewCol(SecretGeneratorColumnChangeDate, e.CreationDate()),
			handler.NewCol(SecretGeneratorColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(SecretGeneratorColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(SecretGeneratorColumnSequence, e.Sequence()),
			handler.NewCol(SecretGeneratorColumnLength, e.Length),
			handler.NewCol(SecretGeneratorColumnExpiry, e.Expiry),
			handler.NewCol(SecretGeneratorColumnIncludeLowerLetters, e.IncludeLowerLetters),
			handler.NewCol(SecretGeneratorColumnIncludeUpperLetters, e.IncludeUpperLetters),
			handler.NewCol(SecretGeneratorColumnIncludeDigits, e.IncludeDigits),
			handler.NewCol(SecretGeneratorColumnIncludeSymbols, e.IncludeSymbols),
		},
	), nil
}

func (p *SecretGeneratorProjection) reduceSecretGeneratorChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SecretGeneratorChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-s00Fs", "reduce.wrong.event.type %s", iam.SecretGeneratorChangedEventType)
	}

	columns := make([]handler.Column, 0, 7)
	columns = append(columns, handler.NewCol(SecretGeneratorColumnChangeDate, e.CreationDate()),
		handler.NewCol(SecretGeneratorColumnSequence, e.Sequence()))
	if e.Length != nil {
		columns = append(columns, handler.NewCol(SecretGeneratorColumnLength, *e.Length))
	}
	if e.Expiry != nil {
		columns = append(columns, handler.NewCol(SecretGeneratorColumnExpiry, *e.Expiry))
	}
	if e.IncludeLowerLetters != nil {
		columns = append(columns, handler.NewCol(SecretGeneratorColumnIncludeLowerLetters, *e.IncludeLowerLetters))
	}
	if e.IncludeUpperLetters != nil {
		columns = append(columns, handler.NewCol(SecretGeneratorColumnIncludeUpperLetters, *e.IncludeUpperLetters))
	}
	if e.IncludeDigits != nil {
		columns = append(columns, handler.NewCol(SecretGeneratorColumnIncludeDigits, *e.IncludeDigits))
	}
	if e.IncludeSymbols != nil {
		columns = append(columns, handler.NewCol(SecretGeneratorColumnIncludeSymbols, *e.IncludeSymbols))
	}
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(SecretGeneratorColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SecretGeneratorColumnGeneratorType, e.GeneratorType),
		},
	), nil
}

func (p *SecretGeneratorProjection) reduceSecretGeneratorRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.SecretGeneratorRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fmiIf", "reduce.wrong.event.type %s", iam.SecretGeneratorRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SecretGeneratorColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SecretGeneratorColumnGeneratorType, e.GeneratorType),
		},
	), nil
}
