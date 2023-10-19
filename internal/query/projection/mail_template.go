package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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

type mailTemplateProjection struct{}

func newMailTemplateProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(mailTemplateProjection))
}

func (*mailTemplateProjection) Name() string {
	return MailTemplateTable
}

func (*mailTemplateProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(MailTemplateAggregateIDCol, handler.ColumnTypeText),
			handler.NewColumn(MailTemplateInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(MailTemplateCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(MailTemplateChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(MailTemplateSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(MailTemplateStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(MailTemplateIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(MailTemplateTemplateCol, handler.ColumnTypeBytes),
			handler.NewColumn(MailTemplateOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(MailTemplateInstanceIDCol, MailTemplateAggregateIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{MailTemplateOwnerRemovedCol})),
		),
	)
}

func (p *mailTemplateProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
	return handler.NewCreateStatement(
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
	return handler.NewUpdateStatement(
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
	return handler.NewDeleteStatement(
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

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(MailTemplateInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(MailTemplateAggregateIDCol, e.Aggregate().ID),
		},
	), nil
}
