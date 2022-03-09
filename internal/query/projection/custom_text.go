package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

const (
	CustomTextTable = "zitadel.projections.custom_texts"

	CustomTextAggregateIDCol  = "aggregate_id"
	CustomTextCreationDateCol = "creation_date"
	CustomTextChangeDateCol   = "change_date"
	CustomTextSequenceCol     = "sequence"
	CustomTextIsDefaultCol    = "is_default"
	CustomTextTemplateCol     = "template"
	CustomTextLanguageCol     = "language"
	CustomTextKeyCol          = "key"
	CustomTextTextCol         = "text"
)

type CustomTextProjection struct {
	crdb.StatementHandler
}

func NewCustomTextProjection(ctx context.Context, config crdb.StatementHandlerConfig) *CustomTextProjection {
	p := new(CustomTextProjection)
	config.ProjectionName = CustomTextTable
	config.Reducers = p.reducers()
	config.InitChecks = []*handler.Check{
		crdb.NewTableCheck(
			crdb.NewTable([]*crdb.Column{
				crdb.NewColumn(CustomTextAggregateIDCol, crdb.ColumnTypeText),
				crdb.NewColumn(CustomTextCreationDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(CustomTextChangeDateCol, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(CustomTextSequenceCol, crdb.ColumnTypeInt64),
				crdb.NewColumn(CustomTextIsDefaultCol, crdb.ColumnTypeBool),
				crdb.NewColumn(CustomTextTemplateCol, crdb.ColumnTypeText),
				crdb.NewColumn(CustomTextLanguageCol, crdb.ColumnTypeText),
				crdb.NewColumn(CustomTextKeyCol, crdb.ColumnTypeText),
				crdb.NewColumn(CustomTextTextCol, crdb.ColumnTypeText),
			},
				crdb.NewPrimaryKey(CustomTextAggregateIDCol, CustomTextTemplateCol, CustomTextKeyCol, CustomTextLanguageCol),
			),
		),
	}
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *CustomTextProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.CustomTextSetEventType,
					Reduce: p.reduceSet,
				},
				{
					Event:  org.CustomTextRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.CustomTextTemplateRemovedEventType,
					Reduce: p.reduceTemplateRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.CustomTextSetEventType,
					Reduce: p.reduceSet,
				},
				{
					Event:  iam.CustomTextRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  iam.CustomTextTemplateRemovedEventType,
					Reduce: p.reduceTemplateRemoved,
				},
			},
		},
	}
}

func (p *CustomTextProjection) reduceSet(event eventstore.Event) (*handler.Statement, error) {
	var customTextEvent policy.CustomTextSetEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.CustomTextSetEvent:
		customTextEvent = e.CustomTextSetEvent
		isDefault = false
	case *iam.CustomTextSetEvent:
		customTextEvent = e.CustomTextSetEvent
		isDefault = true
	default:
		logging.LogWithFields("PROJE-g0Jfs", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.CustomTextSetEventType, iam.CustomTextSetEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-KKfw4", "reduce.wrong.event.type")
	}
	return crdb.NewUpsertStatement(
		&customTextEvent,
		[]handler.Column{
			handler.NewCol(CustomTextAggregateIDCol, customTextEvent.Aggregate().ID),
			handler.NewCol(CustomTextCreationDateCol, customTextEvent.CreationDate()),
			handler.NewCol(CustomTextChangeDateCol, customTextEvent.CreationDate()),
			handler.NewCol(CustomTextSequenceCol, customTextEvent.Sequence()),
			handler.NewCol(CustomTextIsDefaultCol, isDefault),
			handler.NewCol(CustomTextTemplateCol, customTextEvent.Template),
			handler.NewCol(CustomTextLanguageCol, customTextEvent.Language.String()),
			handler.NewCol(CustomTextKeyCol, customTextEvent.Key),
			handler.NewCol(CustomTextTextCol, customTextEvent.Text),
		}), nil
}

func (p *CustomTextProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	var customTextEvent policy.CustomTextRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextRemovedEvent:
		customTextEvent = e.CustomTextRemovedEvent
	case *iam.CustomTextRemovedEvent:
		customTextEvent = e.CustomTextRemovedEvent
	default:
		logging.LogWithFields("PROJE-2Nigw", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.CustomTextRemovedEventType, iam.CustomTextRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-n9wJg", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		&customTextEvent,
		[]handler.Condition{
			handler.NewCond(CustomTextAggregateIDCol, customTextEvent.Aggregate().ID),
			handler.NewCond(CustomTextTemplateCol, customTextEvent.Template),
			handler.NewCond(CustomTextKeyCol, customTextEvent.Key),
			handler.NewCond(CustomTextLanguageCol, customTextEvent.Language.String()),
		}), nil
}

func (p *CustomTextProjection) reduceTemplateRemoved(event eventstore.Event) (*handler.Statement, error) {
	var customTextEvent policy.CustomTextTemplateRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextTemplateRemovedEvent:
		customTextEvent = e.CustomTextTemplateRemovedEvent
	case *iam.CustomTextTemplateRemovedEvent:
		customTextEvent = e.CustomTextTemplateRemovedEvent
	default:
		logging.LogWithFields("PROJE-J9wfg", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.CustomTextTemplateRemovedEventType, iam.CustomTextTemplateRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-29iPf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		&customTextEvent,
		[]handler.Condition{
			handler.NewCond(CustomTextAggregateIDCol, customTextEvent.Aggregate().ID),
			handler.NewCond(CustomTextTemplateCol, customTextEvent.Template),
			handler.NewCond(CustomTextLanguageCol, customTextEvent.Language.String()),
		}), nil
}
