package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PasswordComplexityProjection struct {
	crdb.StatementHandler
}

const (
	PasswordComplexityTable = "zitadel.projections.password_complexity_policies"
)

func NewPasswordComplexityProjection(ctx context.Context, config crdb.StatementHandlerConfig) *PasswordComplexityProjection {
	p := &PasswordComplexityProjection{}
	config.ProjectionName = PasswordComplexityTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *PasswordComplexityProjection) reducers() []handler.AggregateReducer {
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

func (p *PasswordComplexityProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
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
		logging.LogWithFields("PROJE-mP8AR", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PasswordComplexityPolicyAddedEventType, iam.PasswordComplexityPolicyAddedEventType}).Error("was not an  event")
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

func (p *PasswordComplexityProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.PasswordComplexityPolicyChangedEvent
	switch e := event.(type) {
	case *org.PasswordComplexityPolicyChangedEvent:
		policyEvent = e.PasswordComplexityPolicyChangedEvent
	case *iam.PasswordComplexityPolicyChangedEvent:
		policyEvent = e.PasswordComplexityPolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-L4UHn", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PasswordComplexityPolicyChangedEventType, iam.PasswordComplexityPolicyChangedEventType}).Error("was not an  event")
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

func (p *PasswordComplexityProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PasswordComplexityPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-ibd0c", "seq", event.Sequence(), "expectedType", org.PasswordComplexityPolicyRemovedEventType).Error("was not an  event")
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
	ComplexityPolicyIsDefaultCol     = "is_default"     //TODO
	ComplexityPolicyResourceOwnerCol = "resource_owner" //TODO

)
