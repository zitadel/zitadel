package projection

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

const (
	MailTemplateTable = "projections.mail_templates"

	MailTemplateAggregateIDCol  = "aggregate_id"
	MailTemplateCreationDateCol = "creation_date"
	MailTemplateChangeDateCol   = "change_date"
	MailTemplateSequenceCol     = "sequence"
	MailTemplateStateCol        = "state"
	MailTemplateIsDefaultCol    = "is_default"
	MailTemplateTemplateCol     = "template"
)

type MailTemplateProjection struct {
	crdb.StatementHandler
}

func NewMailTemplateProjection(ctx context.Context, config crdb.StatementHandlerConfig) *MailTemplateProjection {
	p := new(MailTemplateProjection)
	config.ProjectionName = MailTemplateTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(MailTemplateAggregateIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MailTemplateCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(MailTemplateChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(MailTemplateSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(MailTemplateStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(MailTemplateIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(MailTemplateTemplateCol, crdb.ColumnTypeBytes),
		},
			crdb.NewPrimaryKey(MailTemplateAggregateIDCol),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *MailTemplateProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.MailTemplateAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.MailTemplateChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.MailTemplateRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.MailTemplateAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.MailTemplateChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *MailTemplateProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var templateEvent policy.MailTemplateAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.MailTemplateAddedEvent:
		templateEvent = e.MailTemplateAddedEvent
		isDefault = false
	case *iam.MailTemplateAddedEvent:
		templateEvent = e.MailTemplateAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-0pJ3f", "reduce.wrong.event.type, %v", []eventstore.EventType{org.MailTemplateAddedEventType, iam.MailTemplateAddedEventType})
	}
	return crdb.NewCreateStatement(
		&templateEvent,
		[]handler.Column{
			handler.NewCol(MailTemplateAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCol(MailTemplateCreationDateCol, templateEvent.CreationDate()),
			handler.NewCol(MailTemplateChangeDateCol, templateEvent.CreationDate()),
			handler.NewCol(MailTemplateSequenceCol, templateEvent.Sequence()),
			handler.NewCol(MailTemplateStateCol, domain.PolicyStateActive),
			handler.NewCol(MailTemplateIsDefaultCol, isDefault),
			handler.NewCol(MailTemplateTemplateCol, templateEvent.Template),
		}), nil
}

func (p *MailTemplateProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.MailTemplateChangedEvent
	switch e := event.(type) {
	case *org.MailTemplateChangedEvent:
		policyEvent = e.MailTemplateChangedEvent
	case *iam.MailTemplateChangedEvent:
		policyEvent = e.MailTemplateChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-gJ03f", "reduce.wrong.event.type, %v", []eventstore.EventType{org.MailTemplateChangedEventType, iam.MailTemplateChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(MailTemplateChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(MailTemplateSequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.Template != nil {
		cols = append(cols, handler.NewCol(MailTemplateTemplateCol, *policyEvent.Template))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(MailTemplateAggregateIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *MailTemplateProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.MailTemplateRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-3jJGs", "reduce.wrong.event.type %s", org.MailTemplateRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(MailTemplateAggregateIDCol, policyEvent.Aggregate().ID),
		}), nil
}
