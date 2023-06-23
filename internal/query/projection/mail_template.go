package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	MailTemplateTable = "projections.mail_templates2"

	MailTemplateAggregateIDCol  = "aggregate_id"
	MailTemplateInstanceIDCol   = "instance_id"
	MailTemplateCreationDateCol = "creation_date"
	MailTemplateChangeDateCol   = "change_date"
	MailTemplateSequenceCol     = "sequence"
	MailTemplateStateCol        = "state"
	MailTemplateIsDefaultCol    = "is_default"
	MailTemplateTemplateCol     = "template"
	MailTemplateOwnerRemovedCol = "owner_removed"
)

type mailTemplateProjection struct {
	crdb.StatementHandler
}

func newMailTemplateProjection(ctx context.Context, config crdb.StatementHandlerConfig) *mailTemplateProjection {
	p := new(mailTemplateProjection)
	config.ProjectionName = MailTemplateTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(MailTemplateAggregateIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MailTemplateInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MailTemplateCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(MailTemplateChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(MailTemplateSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(MailTemplateStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(MailTemplateIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(MailTemplateTemplateCol, crdb.ColumnTypeBytes),
			crdb.NewColumn(MailTemplateOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(MailTemplateInstanceIDCol, MailTemplateAggregateIDCol),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{MailTemplateOwnerRemovedCol})),
		),
	)
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
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.MailTemplateAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.MailTemplateChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(MailTemplateInstanceIDCol),
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
	case *instance.MailTemplateAddedEvent:
		templateEvent = e.MailTemplateAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-0pJ3f", "reduce.wrong.event.type, %v", []eventstore.EventType{org.MailTemplateAddedEventType, instance.MailTemplateAddedEventType})
	}
	return crdb.NewCreateStatement(
		&templateEvent,
		[]handler.Column{
			handler.NewCol(MailTemplateAggregateIDCol, templateEvent.Aggregate().ID),
			handler.NewCol(MailTemplateInstanceIDCol, templateEvent.Aggregate().InstanceID),
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
	case *instance.MailTemplateChangedEvent:
		policyEvent = e.MailTemplateChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-gJ03f", "reduce.wrong.event.type, %v", []eventstore.EventType{org.MailTemplateChangedEventType, instance.MailTemplateChangedEventType})
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
			handler.NewCond(MailTemplateInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *mailTemplateProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.MailTemplateRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-3jJGs", "reduce.wrong.event.type %s", org.MailTemplateRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(MailTemplateAggregateIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(MailTemplateInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *mailTemplateProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-CThXR", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(MailTemplateInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(MailTemplateAggregateIDCol, e.Aggregate().ID),
		},
	), nil
}
