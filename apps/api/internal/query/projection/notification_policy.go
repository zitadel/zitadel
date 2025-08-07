package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
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

type notificationPolicyProjection struct{}

func newNotificationPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(notificationPolicyProjection))
}

func (*notificationPolicyProjection) Name() string {
	return NotificationPolicyProjectionTable
}

func (*notificationPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(NotificationPolicyColumnID, handler.ColumnTypeText),
			handler.NewColumn(NotificationPolicyColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(NotificationPolicyColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(NotificationPolicyColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(NotificationPolicyColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(NotificationPolicyColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(NotificationPolicyColumnStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(NotificationPolicyColumnIsDefault, handler.ColumnTypeBool),
			handler.NewColumn(NotificationPolicyColumnPasswordChange, handler.ColumnTypeBool),
			handler.NewColumn(NotificationPolicyColumnOwnerRemoved, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(NotificationPolicyColumnInstanceID, NotificationPolicyColumnID),
		),
	)
}

func (p *notificationPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-x02s1m", "reduce.wrong.event.type %v", []eventstore.EventType{org.NotificationPolicyAddedEventType, instance.NotificationPolicyAddedEventType})
	}
	return handler.NewCreateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-psom2h19", "reduce.wrong.event.type %v", []eventstore.EventType{org.NotificationPolicyChangedEventType, instance.NotificationPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(NotificationPolicyColumnChangeDate, policyEvent.CreationDate()),
		handler.NewCol(NotificationPolicyColumnSequence, policyEvent.Sequence()),
	}
	if policyEvent.PasswordChange != nil {
		cols = append(cols, handler.NewCol(NotificationPolicyColumnPasswordChange, *policyEvent.PasswordChange))
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Po2iso2", "reduce.wrong.event.type %s", org.NotificationPolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(NotificationPolicyColumnID, policyEvent.Aggregate().ID),
			handler.NewCond(NotificationPolicyColumnInstanceID, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *notificationPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-poxi9a", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(DomainPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
