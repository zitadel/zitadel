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

type lockoutPolicyProjection struct {
	crdb.StatementHandler
}

const (
	LockoutPolicyTable = "zitadel.projections.lockout_policies"

	LockoutPolicyCreationDateCol        = "creation_date"
	LockoutPolicyChangeDateCol          = "change_date"
	LockoutPolicySequenceCol            = "sequence"
	LockoutPolicyIDCol                  = "id"
	LockoutPolicyStateCol               = "state"
	LockoutPolicyMaxPasswordAttemptsCol = "max_password_attempts"
	LockoutPolicyShowLockOutFailuresCol = "show_failure"
	LockoutPolicyIsDefaultCol           = "is_default"
	LockoutPolicyResourceOwnerCol       = "resource_owner"
)

func newLockoutPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *lockoutPolicyProjection {
	p := &lockoutPolicyProjection{}
	config.ProjectionName = LockoutPolicyTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *lockoutPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
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
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.LockoutPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.LockoutPolicyChangedEventType,
					Reduce: p.reduceChanged,
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
	case *iam.LockoutPolicyAddedEvent:
		policyEvent = e.LockoutPolicyAddedEvent
		isDefault = true
	default:
		logging.LogWithFields("PROJE-uFqFM", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LockoutPolicyAddedEventType, iam.LockoutPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-d8mZO", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
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
		}), nil
}

func (p *lockoutPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LockoutPolicyChangedEvent
	switch e := event.(type) {
	case *org.LockoutPolicyChangedEvent:
		policyEvent = e.LockoutPolicyChangedEvent
	case *iam.LockoutPolicyChangedEvent:
		policyEvent = e.LockoutPolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-iIkej", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LockoutPolicyChangedEventType, iam.LockoutPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-pT3mQ", "reduce.wrong.event.type")
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
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(LockoutPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *lockoutPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.LockoutPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-U5cys", "seq", event.Sequence(), "expectedType", org.LockoutPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-Bqut9", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(LockoutPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
