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
	LockoutPolicyTable = "projections.lockout_policies2"

	LockoutPolicyIDCol                  = "id"
	LockoutPolicyCreationDateCol        = "creation_date"
	LockoutPolicyChangeDateCol          = "change_date"
	LockoutPolicySequenceCol            = "sequence"
	LockoutPolicyStateCol               = "state"
	LockoutPolicyIsDefaultCol           = "is_default"
	LockoutPolicyResourceOwnerCol       = "resource_owner"
	LockoutPolicyInstanceIDCol          = "instance_id"
	LockoutPolicyMaxPasswordAttemptsCol = "max_password_attempts"
	LockoutPolicyShowLockOutFailuresCol = "show_failure"
	LockoutPolicyOwnerRemovedCol        = "owner_removed"
)

type lockoutPolicyProjection struct{}

func newLockoutPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(lockoutPolicyProjection))
}

func (*lockoutPolicyProjection) Name() string {
	return LockoutPolicyTable
}

func (*lockoutPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(LockoutPolicyIDCol, handler.ColumnTypeText),
			handler.NewColumn(LockoutPolicyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(LockoutPolicyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(LockoutPolicySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(LockoutPolicyStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(LockoutPolicyIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LockoutPolicyResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(LockoutPolicyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(LockoutPolicyMaxPasswordAttemptsCol, handler.ColumnTypeInt64),
			handler.NewColumn(LockoutPolicyShowLockOutFailuresCol, handler.ColumnTypeBool),
			handler.NewColumn(LockoutPolicyOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(LockoutPolicyInstanceIDCol, LockoutPolicyIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{LockoutPolicyOwnerRemovedCol})),
		),
	)
}

func (p *lockoutPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.LockoutPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.LockoutPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.LockoutPolicyRemovedEventType,
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
					Event:  instance.LockoutPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.LockoutPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(LockoutPolicyInstanceIDCol),
				},
			},
		},
	}
}

func (p *lockoutPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LockoutPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
		isDefault = false
	case *instance.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-d8mZO", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyAddedEventType, instance.LockoutPolicyAddedEventType})
	}
	return handler.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LockoutPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(LockoutPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LockoutPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(LockoutPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(LockoutPolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(LockoutPolicyMaxPasswordAttemptsCol, policyEvent.MaxPasswordAttempts),
			handler.NewCol(LockoutPolicyShowLockOutFailuresCol, policyEvent.ShowLockOutFailures),
			handler.NewCol(LockoutPolicyIsDefaultCol, isDefault),
			handler.NewCol(LockoutPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(LockoutPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *lockoutPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LockoutPolicyChangedEvent
	switch e := event.(type) {
	case *org.LockoutPolicyChangedEvent:
		policyEvent = e.LockoutPolicyChangedEvent
	case *instance.LockoutPolicyChangedEvent:
		policyEvent = e.LockoutPolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-pT3mQ", "reduce.wrong.event.type, %v", []eventstore.EventType{org.LockoutPolicyChangedEventType, instance.LockoutPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(LockoutPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(LockoutPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.MaxPasswordAttempts != nil {
		cols = append(cols, handler.NewCol(LockoutPolicyMaxPasswordAttemptsCol, *policyEvent.MaxPasswordAttempts))
	}
	if policyEvent.ShowLockOutFailures != nil {
		cols = append(cols, handler.NewCol(LockoutPolicyShowLockOutFailuresCol, *policyEvent.ShowLockOutFailures))
	}
	return handler.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(LockoutPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *lockoutPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LockoutPolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Bqut9", "reduce.wrong.event.type %s", org.LockoutPolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(LockoutPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LabelPolicyInstanceIDCol, event.Aggregate().InstanceID),
		}), nil
}

func (p *lockoutPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-IoW0x", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LockoutPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(LockoutPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
