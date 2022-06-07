package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type mailTemplateProjection struct {
	crdb.StatementHandler
}

const (
	MailTemplateTable = "zitadel.projections.mail_templates"

	MailTemplateAggregateIDCol  = "aggregate_id"
	MailTemplateCreationDateCol = "creation_date"
	MailTemplateChangeDateCol   = "change_date"
	MailTemplateSequenceCol     = "sequence"
	MailTemplateStateCol        = "state"
	MailTemplateTemplateCol     = "template"
	MailTemplateIsDefaultCol    = "is_default"
)

func newMailTemplateProjection(ctx context.Context, config crdb.StatementHandlerConfig) *mailTemplateProjection {
	p := &mailTemplateProjection{}
	config.ProjectionName = MailTemplateTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *mailTemplateProjection) reducers() []handler.AggregateReducer {
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

func (p *mailTemplateProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
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
		logging.LogWithFields("PROJE-94jfG", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.MailTemplateAddedEventType, iam.MailTemplateAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-0pJ3f", "reduce.wrong.event.type")
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

func (p *mailTemplateProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.MailTemplateChangedEvent
	switch e := event.(type) {
	case *org.MailTemplateChangedEvent:
		policyEvent = e.MailTemplateChangedEvent
	case *iam.MailTemplateChangedEvent:
		policyEvent = e.MailTemplateChangedEvent
	default:
		logging.LogWithFields("PROJE-02J9f", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.MailTemplateChangedEventType, iam.MailTemplateChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-gJ03f", "reduce.wrong.event.type")
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

func (p *mailTemplateProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.MailTemplateRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-2m0fp", "seq", event.Sequence(), "expectedType", org.MailTemplateRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-3jJGs", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(MailTemplateAggregateIDCol, policyEvent.Aggregate().ID),
		}), nil
}
