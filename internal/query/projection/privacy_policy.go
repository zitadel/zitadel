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
	PrivacyPolicyTable = "projections.privacy_policies"

	PrivacyPolicyIDCol            = "id"
	PrivacyPolicyCreationDateCol  = "creation_date"
	PrivacyPolicyChangeDateCol    = "change_date"
	PrivacyPolicySequenceCol      = "sequence"
	PrivacyPolicyStateCol         = "state"
	PrivacyPolicyIsDefaultCol     = "is_default"
	PrivacyPolicyResourceOwnerCol = "resource_owner"
	PrivacyPolicyInstanceIDCol    = "instance_id"
	PrivacyPolicyPrivacyLinkCol   = "privacy_link"
	PrivacyPolicyTOSLinkCol       = "tos_link"
)

type PrivacyPolicyProjection struct {
	crdb.StatementHandler
}

func NewPrivacyPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *PrivacyPolicyProjection {
	p := new(PrivacyPolicyProjection)
	config.ProjectionName = PrivacyPolicyTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(PrivacyPolicyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(PrivacyPolicyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PrivacyPolicyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(PrivacyPolicySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(PrivacyPolicyStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(PrivacyPolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(PrivacyPolicyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(PrivacyPolicyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(PrivacyPolicyPrivacyLinkCol, crdb.ColumnTypeText),
			crdb.NewColumn(PrivacyPolicyTOSLinkCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(PrivacyPolicyInstanceIDCol, PrivacyPolicyIDCol),
		),
	)
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
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.PrivacyPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.PrivacyPolicyChangedEventType,
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
	case *iam.PrivacyPolicyAddedEvent:
		policyEvent = e.PrivacyPolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-kRNh8", "reduce.wrong.event.type %v", []eventstore.EventType{org.PrivacyPolicyAddedEventType, iam.PrivacyPolicyAddedEventType})
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
			handler.NewCol(PrivacyPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *PrivacyPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.PrivacyPolicyChangedEvent
	switch e := event.(type) {
	case *org.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	case *iam.PrivacyPolicyChangedEvent:
		policyEvent = e.PrivacyPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-91weZ", "reduce.wrong.event.type %v", []eventstore.EventType{org.PrivacyPolicyChangedEventType, iam.PrivacyPolicyChangedEventType})
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
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-FvtGO", "reduce.wrong.event.type %s", org.PrivacyPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(PrivacyPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
