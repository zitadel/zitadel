package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PrivacyPolicyProjection struct {
	crdb.StatementHandler
}

const (
	PrivacyPolicyTable = "zitadel.projections.privacy_policies"

	PrivacyPolicyCreationDateCol  = "creation_date"
	PrivacyPolicyChangeDateCol    = "change_date"
	PrivacyPolicySequenceCol      = "sequence"
	PrivacyPolicyIDCol            = "id"
	PrivacyPolicyStateCol         = "state"
	PrivacyPolicyPrivacyLinkCol   = "privacy_link"
	PrivacyPolicyTOSLinkCol       = "tos_link"
	PrivacyPolicyIsDefaultCol     = "is_default"
	PrivacyPolicyResourceOwnerCol = "resource_owner"
)

func NewPrivacyPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *PrivacyPolicyProjection {
	p := new(PrivacyPolicyProjection)
	config.ProjectionName = PrivacyPolicyTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *PrivacyPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.PrivacyPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.PrivacyPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.PrivacyPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.PrivacyPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.PrivacyPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *PrivacyPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PrivacyPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
		isDefault = false
	case *instance.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
		isDefault = true
	default:
		logging.LogWithFields("PROJE-BrdLn", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PrivacyPolicyAddedEventType, instance.PrivacyPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-kRNh8", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(PrivacyPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(PrivacyPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(PrivacyPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(PrivacyPolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(PrivacyPolicyPrivacyLinkCol, policyEvent.PrivacyLink),
			handler.NewCol(PrivacyPolicyTOSLinkCol, policyEvent.TOSLink),
			handler.NewCol(PrivacyPolicyIsDefaultCol, isDefault),
			handler.NewCol(PrivacyPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
		}), nil
}

func (p *PrivacyPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PrivacyPolicyChangedEvent
	switch e := event.(type) {
	case *org.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	case *instance.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-1nQWm", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.PrivacyPolicyChangedEventType, instance.PrivacyPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-91weZ", "reduce.wrong.event.type")
	}
	cols := []handler.Column{
		handler.NewCol(PrivacyPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(PrivacyPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.PrivacyLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyPrivacyLinkCol, *policyEvent.PrivacyLink))
	}
	if policyEvent.TOSLink != nil {
		cols = append(cols, handler.NewCol(PrivacyPolicyTOSLinkCol, *policyEvent.TOSLink))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *PrivacyPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.PrivacyPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-hN5Ip", "seq", event.Sequence(), "expectedType", org.PrivacyPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-FvtGO", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
