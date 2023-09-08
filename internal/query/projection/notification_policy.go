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
	NotificationPolicyProjectionTable = "projections.notification_policies"

	NotificationPolicyColumnID             = "id"
	NotificationPolicyColumnCreationDate   = "creation_date"
	NotificationPolicyColumnChangeDate     = "change_date"
	NotificationPolicyColumnResourceOwner  = "resource_owner"
	NotificationPolicyColumnInstanceID     = "instance_id"
	NotificationPolicyColumnSequence       = "sequence"
	NotificationPolicyColumnStateCol       = "state"
	NotificationPolicyColumnIsDefault      = "is_default"
	NotificationPolicyColumnPasswordChange = "password_change"
	NotificationPolicyColumnOwnerRemoved   = "owner_removed"
)

type notificationPolicyProjection struct {
	crdb.StatementHandler
}

func newNotificationPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *notificationPolicyProjection {
	p := new(notificationPolicyProjection)
	config.ProjectionName = NotificationPolicyProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(NotificationPolicyColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(NotificationPolicyColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(NotificationPolicyColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(NotificationPolicyColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(NotificationPolicyColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(NotificationPolicyColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(NotificationPolicyColumnStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(NotificationPolicyColumnIsDefault, crdb.ColumnTypeBool),
			crdb.NewColumn(NotificationPolicyColumnPasswordChange, crdb.ColumnTypeBool),
			crdb.NewColumn(NotificationPolicyColumnOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(NotificationPolicyColumnInstanceID, NotificationPolicyColumnID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *notificationPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.NotificationPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.NotificationPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.NotificationPolicyRemovedEventType,
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
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(NotificationPolicyColumnInstanceID),
				},
				{
					Event:  instance.NotificationPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.NotificationPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *notificationPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.NotificationPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.NotificationPolicyAddedEvent:
		policyEvent = e.NotificationPolicyAddedEvent
		isDefault = false
	case *instance.NotificationPolicyAddedEvent:
		policyEvent = e.NotificationPolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-x02s1m", "reduce.wrong.event.type %v", []eventstore.EventType{org.NotificationPolicyAddedEventType, instance.NotificationPolicyAddedEventType})
	}
	return crdb.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(NotificationPolicyColumnCreationDate, policyEvent.CreationDate()),
			handler.NewCol(NotificationPolicyColumnChangeDate, policyEvent.CreationDate()),
			handler.NewCol(NotificationPolicyColumnSequence, policyEvent.Sequence()),
			handler.NewCol(NotificationPolicyColumnID, policyEvent.Aggregate().ID),
			handler.NewCol(NotificationPolicyColumnStateCol, domain.PolicyStateActive),
			handler.NewCol(NotificationPolicyColumnPasswordChange, policyEvent.PasswordChange),
			handler.NewCol(NotificationPolicyColumnIsDefault, isDefault),
			handler.NewCol(NotificationPolicyColumnResourceOwner, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(NotificationPolicyColumnInstanceID, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *notificationPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.NotificationPolicyChangedEvent
	switch e := event.(type) {
	case *org.NotificationPolicyChangedEvent:
		policyEvent = e.NotificationPolicyChangedEvent
	case *instance.NotificationPolicyChangedEvent:
		policyEvent = e.NotificationPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-psom2h19", "reduce.wrong.event.type %v", []eventstore.EventType{org.NotificationPolicyChangedEventType, instance.NotificationPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(NotificationPolicyColumnChangeDate, policyEvent.CreationDate()),
		handler.NewCol(NotificationPolicyColumnSequence, policyEvent.Sequence()),
	}
	if policyEvent.PasswordChange != nil {
		cols = append(cols, handler.NewCol(NotificationPolicyColumnPasswordChange, *policyEvent.PasswordChange))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(NotificationPolicyColumnID, policyEvent.Aggregate().ID),
			handler.NewCond(NotificationPolicyColumnInstanceID, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *notificationPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.NotificationPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Po2iso2", "reduce.wrong.event.type %s", org.NotificationPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(NotificationPolicyColumnID, policyEvent.Aggregate().ID),
			handler.NewCond(NotificationPolicyColumnInstanceID, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *notificationPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-poxi9a", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(DomainPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
