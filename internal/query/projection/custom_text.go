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

type CustomTextProjection struct {
	crdb.StatementHandler
}

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

func NewCustomTextProjection(ctx context.Context, config crdb.StatementHandlerConfig) *CustomTextProjection {
	p := &CustomTextProjection{}
	config.ProjectionName = CustomTextTable
	config.Reducers = p.reducers()
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
					Reduce: p.reduceAdded,
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
					Reduce: p.reduceAdded,
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

func (p *CustomTextProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var templateEvent policy.CustomTextSetEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.CustomTextSetEvent:
		templateEvent = e.CustomTextSetEvent
		isDefault = false
	case *iam.CustomTextSetEvent:
		templateEvent = e.CustomTextSetEvent
		isDefault = true
	default:
		logging.LogWithFields("PROJE-g0Jfs", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.CustomTextSetEventType, iam.CustomTextSetEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-KKfw4", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		&templateEvent,
		[]handler.Column{
			handler.NewCol(CustomTextAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCol(CustomTextCreationDateCol, templateEvent.CreationDate()),
			handler.NewCol(CustomTextChangeDateCol, templateEvent.CreationDate()),
			handler.NewCol(CustomTextSequenceCol, templateEvent.Sequence()),
			handler.NewCol(CustomTextIsDefaultCol, isDefault),
			handler.NewCol(CustomTextTemplateCol, templateEvent.Template),
			handler.NewCol(CustomTextLanguageCol, templateEvent.Language.String()),
			handler.NewCol(CustomTextKeyCol, templateEvent.Key),
			handler.NewCol(CustomTextTextCol, templateEvent.Text),
		}), nil
}

func (p *CustomTextProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.CustomTextRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-wm00r", "seq", event.Sequence(), "expectedType", org.CustomTextRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-sJ0gs", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(CustomTextAggregateIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(CustomTextTemplateCol, policyEvent.Template),
			handler.NewCond(CustomTextKeyCol, policyEvent.Key),
			handler.NewCond(CustomTextLanguageCol, policyEvent.Language),
		}), nil
}

func (p *CustomTextProjection) reduceTemplateRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.CustomTextRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-2j0gs", "seq", event.Sequence(), "expectedType", org.CustomTextRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-0gmeG", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(CustomTextAggregateIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(CustomTextTemplateCol, policyEvent.Template),
			handler.NewCond(CustomTextLanguageCol, policyEvent.Language),
		}), nil
}
