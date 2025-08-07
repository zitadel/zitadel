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
	PasswordAgeTable = "projections.password_age_policies2"

	AgePolicyIDCol             = "id"
	AgePolicyCreationDateCol   = "creation_date"
	AgePolicyChangeDateCol     = "change_date"
	AgePolicySequenceCol       = "sequence"
	AgePolicyStateCol          = "state"
	AgePolicyIsDefaultCol      = "is_default"
	AgePolicyResourceOwnerCol  = "resource_owner"
	AgePolicyInstanceIDCol     = "instance_id"
	AgePolicyExpireWarnDaysCol = "expire_warn_days"
	AgePolicyMaxAgeDaysCol     = "max_age_days"
	AgePolicyOwnerRemovedCol   = "owner_removed"
)

type passwordAgeProjection struct{}

func newPasswordAgeProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(passwordAgeProjection))
}

func (*passwordAgeProjection) Name() string {
	return PasswordAgeTable
}

func (*passwordAgeProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(AgePolicyIDCol, handler.ColumnTypeText),
			handler.NewColumn(AgePolicyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(AgePolicyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(AgePolicySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(AgePolicyStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(AgePolicyIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(AgePolicyResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(AgePolicyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(AgePolicyExpireWarnDaysCol, handler.ColumnTypeInt64),
			handler.NewColumn(AgePolicyMaxAgeDaysCol, handler.ColumnTypeInt64),
			handler.NewColumn(AgePolicyOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(AgePolicyInstanceIDCol, AgePolicyIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{AgePolicyOwnerRemovedCol})),
		),
	)
}

func (p *passwordAgeProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.PasswordAgePolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.PasswordAgePolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.PasswordAgePolicyRemovedEventType,
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
					Event:  instance.PasswordAgePolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.PasswordAgePolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(AgePolicyInstanceIDCol),
				},
			},
		},
	}
}

func (p *passwordAgeProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordAgePolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		isDefault = false
	case *instance.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		isDefault = true
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-CJqF0", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyAddedEventType, instance.PasswordAgePolicyAddedEventType})
	}
	return handler.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(AgePolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(AgePolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(AgePolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(AgePolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(AgePolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(AgePolicyExpireWarnDaysCol, policyEvent.ExpireWarnDays),
			handler.NewCol(AgePolicyMaxAgeDaysCol, policyEvent.MaxAgeDays),
			handler.NewCol(AgePolicyIsDefaultCol, isDefault),
			handler.NewCol(AgePolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(AgePolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *passwordAgeProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordAgePolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	case *instance.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-i7FZt", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyChangedEventType, instance.PasswordAgePolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(AgePolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(AgePolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.ExpireWarnDays != nil {
		cols = append(cols, handler.NewCol(AgePolicyExpireWarnDaysCol, *policyEvent.ExpireWarnDays))
	}
	if policyEvent.MaxAgeDays != nil {
		cols = append(cols, handler.NewCol(AgePolicyMaxAgeDaysCol, *policyEvent.MaxAgeDays))
	}
	return handler.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(AgePolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(AgePolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *passwordAgeProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordAgePolicyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-EtHWB", "reduce.wrong.event.type %s", org.PasswordAgePolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(AgePolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(AgePolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *passwordAgeProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-edLs2", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(AgePolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(AgePolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
