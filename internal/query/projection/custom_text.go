package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	CustomTextTable = "projections.custom_texts2"

	CustomTextAggregateIDCol  = "aggregate_id"
	CustomTextInstanceIDCol   = "instance_id"
	CustomTextCreationDateCol = "creation_date"
	CustomTextChangeDateCol   = "change_date"
	CustomTextSequenceCol     = "sequence"
	CustomTextIsDefaultCol    = "is_default"
	CustomTextTemplateCol     = "template"
	CustomTextLanguageCol     = "language"
	CustomTextKeyCol          = "key"
	CustomTextTextCol         = "text"
	CustomTextOwnerRemovedCol = "owner_removed"
)

type customTextProjection struct{}

func newCustomTextProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(customTextProjection))
}

func (*customTextProjection) Name() string {
	return CustomTextTable
}

func (*customTextProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(CustomTextAggregateIDCol, handler.ColumnTypeText),
			handler.NewColumn(CustomTextInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(CustomTextCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(CustomTextChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(CustomTextSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(CustomTextIsDefaultCol, handler.ColumnTypeBool),
			handler.NewColumn(CustomTextTemplateCol, handler.ColumnTypeText),
			handler.NewColumn(CustomTextLanguageCol, handler.ColumnTypeText),
			handler.NewColumn(CustomTextKeyCol, handler.ColumnTypeText),
			handler.NewColumn(CustomTextTextCol, handler.ColumnTypeText),
			handler.NewColumn(CustomTextOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(CustomTextInstanceIDCol, CustomTextAggregateIDCol, CustomTextTemplateCol, CustomTextKeyCol, CustomTextLanguageCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{CustomTextOwnerRemovedCol})),
		),
	)
}

func (p *customTextProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.CustomTextSetEventType,
					Reduce: p.reduceSet,
				},
				{
					Event:  instance.CustomTextRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  instance.CustomTextTemplateRemovedEventType,
					Reduce: p.reduceTemplateRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(CustomTextInstanceIDCol),
				},
			},
		},
	}
}

func (p *customTextProjection) reduceSet(event eventstore.Event) (*handler.Statement, error) {
	var customTextEvent policy.CustomTextSetEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.CustomTextSetEvent:
		customTextEvent = e.CustomTextSetEvent
		isDefault = false
	case *instance.CustomTextSetEvent:
		customTextEvent = e.CustomTextSetEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-KKfw4", "reduce.wrong.event.type %v", []eventstore.EventType{org.CustomTextSetEventType, instance.CustomTextSetEventType})
	}
	return handler.NewUpsertStatement(
		&customTextEvent,
		[]handler.Column{
			handler.NewCol(CustomTextInstanceIDCol, nil),
			handler.NewCol(CustomTextAggregateIDCol, nil),
			handler.NewCol(CustomTextTemplateCol, nil),
			handler.NewCol(CustomTextKeyCol, nil),
			handler.NewCol(CustomTextLanguageCol, nil),
		},
		[]handler.Column{
			handler.NewCol(CustomTextAggregateIDCol, customTextEvent.Aggregate().ID),
			handler.NewCol(CustomTextInstanceIDCol, customTextEvent.Aggregate().InstanceID),
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

func (p *customTextProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	var customTextEvent policy.CustomTextRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextRemovedEvent:
		customTextEvent = e.CustomTextRemovedEvent
	case *instance.CustomTextRemovedEvent:
		customTextEvent = e.CustomTextRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-n9wJg", "reduce.wrong.event.type %v", []eventstore.EventType{org.CustomTextRemovedEventType, instance.CustomTextRemovedEventType})
	}
	return handler.NewDeleteStatement(
		&customTextEvent,
		[]handler.Condition{
			handler.NewCond(CustomTextAggregateIDCol, customTextEvent.Aggregate().ID),
			handler.NewCond(CustomTextTemplateCol, customTextEvent.Template),
			handler.NewCond(CustomTextKeyCol, customTextEvent.Key),
			handler.NewCond(CustomTextLanguageCol, customTextEvent.Language.String()),
			handler.NewCond(CustomTextInstanceIDCol, customTextEvent.Aggregate().InstanceID),
		}), nil
}

func (p *customTextProjection) reduceTemplateRemoved(event eventstore.Event) (*handler.Statement, error) {
	var customTextEvent policy.CustomTextTemplateRemovedEvent
	switch e := event.(type) {
	case *org.CustomTextTemplateRemovedEvent:
		customTextEvent = e.CustomTextTemplateRemovedEvent
	case *instance.CustomTextTemplateRemovedEvent:
		customTextEvent = e.CustomTextTemplateRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-29iPf", "reduce.wrong.event.type %v", []eventstore.EventType{org.CustomTextTemplateRemovedEventType, instance.CustomTextTemplateRemovedEventType})
	}
	return handler.NewDeleteStatement(
		&customTextEvent,
		[]handler.Condition{
			handler.NewCond(CustomTextAggregateIDCol, customTextEvent.Aggregate().ID),
			handler.NewCond(CustomTextTemplateCol, customTextEvent.Template),
			handler.NewCond(CustomTextLanguageCol, customTextEvent.Language.String()),
			handler.NewCond(CustomTextInstanceIDCol, customTextEvent.Aggregate().InstanceID),
		}), nil
}

func (p *customTextProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-V2T3z", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(CustomTextInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(CustomTextAggregateIDCol, e.Aggregate().ID),
		},
	), nil
}
