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

type passwordAgeProjection struct {
	crdb.StatementHandler
}

const (
	PasswordAgeTable = "zitadel.projections.password_age_policies"
)

func newPasswordAgeProjection(ctx context.Context, config crdb.StatementHandlerConfig) *passwordAgeProjection {
	p := &passwordAgeProjection{}
	config.ProjectionName = PasswordAgeTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *passwordAgeProjection) reducers() []handler.AggregateReducer {
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

func (p *passwordAgeProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
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
		logging.LogWithFields("PROJE-stxcL", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PasswordAgePolicyAddedEventType, iam.PasswordAgePolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-CJqF0", "reduce.wrong.event.type")
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
		}), nil
}

func (p *passwordAgeProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordAgePolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	case *iam.PasswordAgePolicyChangedEvent:
		policyEvent = e.PasswordAgePolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-EZ53p", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PasswordAgePolicyChangedEventType, iam.PasswordAgePolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-i7FZt", "reduce.wrong.event.type")
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

func (p *passwordAgeProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordAgePolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-iwqfN", "seq", event.Sequence(), "expectedType", org.PasswordAgePolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-EtHWB", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(AgePolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

const (
	AgePolicyCreationDateCol   = "creation_date"
	AgePolicyChangeDateCol     = "change_date"
	AgePolicySequenceCol       = "sequence"
	AgePolicyIDCol             = "id"
	AgePolicyStateCol          = "state"
	AgePolicyExpireWarnDaysCol = "expire_warn_days"
	AgePolicyMaxAgeDaysCol     = "max_age_days"
	AgePolicyIsDefaultCol      = "is_default"
	AgePolicyResourceOwnerCol  = "resource_owner"
)
