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
	PasswordAgeTable = "projections.password_age_policies"

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
)

type PasswordAgeProjection struct {
	crdb.StatementHandler
}

func NewPasswordAgeProjection(ctx context.Context, config crdb.StatementHandlerConfig) *PasswordAgeProjection {
	p := new(PasswordAgeProjection)
	config.ProjectionName = PasswordAgeTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(AgePolicyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AgePolicyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AgePolicyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(AgePolicySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(AgePolicyStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(AgePolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(AgePolicyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(AgePolicyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AgePolicyExpireWarnDaysCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(AgePolicyMaxAgeDaysCol, crdb.ColumnTypeInt64),
		},
			crdb.NewPrimaryKey(AgePolicyInstanceIDCol, AgePolicyIDCol),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *PasswordAgeProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
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
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.PasswordAgePolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.PasswordAgePolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *PasswordAgeProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordAgePolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		isDefault = false
	case *iam.PasswordAgePolicyAddedEvent:
		policyEvent = e.PasswordAgePolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-CJqF0", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyAddedEventType, iam.PasswordAgePolicyAddedEventType})
	}
	return crdb.NewCreateStatement(
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

func (p *PasswordAgeProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordAgePolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	case *iam.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-i7FZt", "reduce.wrong.event.type %v", []eventstore.EventType{org.PasswordAgePolicyChangedEventType, iam.PasswordAgePolicyChangedEventType})
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
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(AgePolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *PasswordAgeProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordAgePolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-EtHWB", "reduce.wrong.event.type %s", org.PasswordAgePolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(AgePolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
