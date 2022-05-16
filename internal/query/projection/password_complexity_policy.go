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

type passwordComplexityProjection struct {
	crdb.StatementHandler
}

const (
	PasswordComplexityTable = "zitadel.projections.password_complexity_policies"
)

func newPasswordComplexityProjection(ctx context.Context, config crdb.StatementHandlerConfig) *passwordComplexityProjection {
	p := &passwordComplexityProjection{}
	config.ProjectionName = PasswordComplexityTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *passwordComplexityProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.PasswordComplexityPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.PasswordComplexityPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.PasswordComplexityPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.PasswordComplexityPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.PasswordComplexityPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *passwordComplexityProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordComplexityPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PasswordComplexityPolicyAddedEvent:
		policyEvent = e.PasswordComplexityPolicyAddedEvent
		isDefault = false
	case *iam.PasswordComplexityPolicyAddedEvent:
		policyEvent = e.PasswordComplexityPolicyAddedEvent
		isDefault = true
	default:
		logging.LogWithFields("PROJE-mP8AR", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PasswordComplexityPolicyAddedEventType, iam.PasswordComplexityPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-KTHmJ", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(ComplexityPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(ComplexityPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(ComplexityPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(ComplexityPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(ComplexityPolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(ComplexityPolicyMinLengthCol, policyEvent.MinLength),
			handler.NewCol(ComplexityPolicyHasLowercaseCol, policyEvent.HasLowercase),
			handler.NewCol(ComplexityPolicyHasUppercaseCol, policyEvent.HasUppercase),
			handler.NewCol(ComplexityPolicyHasSymbolCol, policyEvent.HasSymbol),
			handler.NewCol(ComplexityPolicyHasNumberCol, policyEvent.HasNumber),
			handler.NewCol(ComplexityPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(ComplexityPolicyIsDefaultCol, isDefault),
		}), nil
}

func (p *passwordComplexityProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PasswordComplexityPolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordComplexityPolicyChangedEvent:
		policyEvent = e.PasswordComplexityPolicyChangedEvent
	case *iam.PasswordComplexityPolicyChangedEvent:
		policyEvent = e.PasswordComplexityPolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-L4UHn", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PasswordComplexityPolicyChangedEventType, iam.PasswordComplexityPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-cf3Xb", "reduce.wrong.event.type")
	}
	cols := []handler.Column{
		handler.NewCol(ComplexityPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(ComplexityPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.MinLength != nil {
		cols = append(cols, handler.NewCol(ComplexityPolicyMinLengthCol, *policyEvent.MinLength))
	}
	if policyEvent.HasLowercase != nil {
		cols = append(cols, handler.NewCol(ComplexityPolicyHasLowercaseCol, *policyEvent.HasLowercase))
	}
	if policyEvent.HasUppercase != nil {
		cols = append(cols, handler.NewCol(ComplexityPolicyHasUppercaseCol, *policyEvent.HasUppercase))
	}
	if policyEvent.HasSymbol != nil {
		cols = append(cols, handler.NewCol(ComplexityPolicyHasSymbolCol, *policyEvent.HasSymbol))
	}
	if policyEvent.HasNumber != nil {
		cols = append(cols, handler.NewCol(ComplexityPolicyHasNumberCol, *policyEvent.HasNumber))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(ComplexityPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *passwordComplexityProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordComplexityPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-ibd0c", "seq", event.Sequence(), "expectedType", org.PasswordComplexityPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-wttCd", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(ComplexityPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

const (
	ComplexityPolicyCreationDateCol  = "creation_date"
	ComplexityPolicyChangeDateCol    = "change_date"
	ComplexityPolicySequenceCol      = "sequence"
	ComplexityPolicyIDCol            = "id"
	ComplexityPolicyStateCol         = "state"
	ComplexityPolicyMinLengthCol     = "min_length"
	ComplexityPolicyHasLowercaseCol  = "has_lowercase"
	ComplexityPolicyHasUppercaseCol  = "has_uppercase"
	ComplexityPolicyHasSymbolCol     = "has_symbol"
	ComplexityPolicyHasNumberCol     = "has_number"
	ComplexityPolicyIsDefaultCol     = "is_default"
	ComplexityPolicyResourceOwnerCol = "resource_owner"
)
