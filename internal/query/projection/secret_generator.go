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
	SecretGeneratorProjectionTable = "projections.secret_generators2"

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

type secretGeneratorProjection struct{}

func newSecretGeneratorProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(secretGeneratorProjection))
}

func (*secretGeneratorProjection) Name() string {
	return SecretGeneratorProjectionTable
}

func (*secretGeneratorProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SecretGeneratorColumnGeneratorType, handler.ColumnTypeEnum),
			handler.NewColumn(SecretGeneratorColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(SecretGeneratorColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SecretGeneratorColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SecretGeneratorColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SecretGeneratorColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(SecretGeneratorColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SecretGeneratorColumnLength, handler.ColumnTypeInt64),
			handler.NewColumn(SecretGeneratorColumnExpiry, handler.ColumnTypeInt64),
			handler.NewColumn(SecretGeneratorColumnIncludeLowerLetters, handler.ColumnTypeBool),
			handler.NewColumn(SecretGeneratorColumnIncludeUpperLetters, handler.ColumnTypeBool),
			handler.NewColumn(SecretGeneratorColumnIncludeDigits, handler.ColumnTypeBool),
			handler.NewColumn(SecretGeneratorColumnIncludeSymbols, handler.ColumnTypeBool),
		},
			handler.NewPrimaryKey(SecretGeneratorColumnInstanceID, SecretGeneratorColumnGeneratorType, SecretGeneratorColumnAggregateID),
		),
	)
}

func (p *secretGeneratorProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.SecretGeneratorAddedEventType,
					Reduce: p.reduceSecretGeneratorAdded,
				},
				{
					Event:  instance.SecretGeneratorChangedEventType,
					Reduce: p.reduceSecretGeneratorChanged,
				},
				{
					Event:  instance.SecretGeneratorRemovedEventType,
					Reduce: p.reduceSecretGeneratorRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SecretGeneratorColumnInstanceID),
				},
			},
		},
	}
}

func (p *secretGeneratorProjection) reduceSecretGeneratorAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecretGeneratorAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-sk99F", "reduce.wrong.event.type %s", instance.SecretGeneratorAddedEventType)
	}
	return handler.NewCreateStatement(
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

func (p *secretGeneratorProjection) reduceSecretGeneratorChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecretGeneratorChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-s00Fs", "reduce.wrong.event.type %s", instance.SecretGeneratorChangedEventType)
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
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(SecretGeneratorColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SecretGeneratorColumnGeneratorType, e.GeneratorType),
			handler.NewCond(SecretGeneratorColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *secretGeneratorProjection) reduceSecretGeneratorRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.SecretGeneratorRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-fmiIf", "reduce.wrong.event.type %s", instance.SecretGeneratorRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SecretGeneratorColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(SecretGeneratorColumnGeneratorType, e.GeneratorType),
			handler.NewCond(SecretGeneratorColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
