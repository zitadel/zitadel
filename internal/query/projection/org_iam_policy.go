package projection

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

const (
	DomainPolicyTable = "projections.domain_policies"

	DomainPolicyIDCol                    = "id"
	DomainPolicyCreationDateCol          = "creation_date"
	DomainPolicyChangeDateCol            = "change_date"
	DomainPolicySequenceCol              = "sequence"
	DomainPolicyStateCol                 = "state"
	DomainPolicyUserLoginMustBeDomainCol = "user_login_must_be_domain"
	DomainPolicyIsDefaultCol             = "is_default"
	DomainPolicyResourceOwnerCol         = "resource_owner"
	DomainPolicyInstanceIDCol            = "instance_id"
)

type DomainPolicyProjection struct {
	crdb.StatementHandler
}

func NewDomainPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *DomainPolicyProjection {
	p := new(DomainPolicyProjection)
	config.ProjectionName = DomainPolicyTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(DomainPolicyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(DomainPolicyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DomainPolicyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DomainPolicySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(DomainPolicyStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(DomainPolicyUserLoginMustBeDomainCol, crdb.ColumnTypeBool),
			crdb.NewColumn(DomainPolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(DomainPolicyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(DomainPolicyInstanceIDCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(DomainPolicyInstanceIDCol, DomainPolicyIDCol),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *DomainPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgDomainPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.OrgDomainPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.OrgDomainPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceDomainPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.InstanceDomainPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *DomainPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.DomainPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.OrgDomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		isDefault = false
	case *instance.InstanceDomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.OrgDomainPolicyAddedEventType, instance.InstanceDomainPolicyAddedEventType})
	}
	return crdb.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(DomainPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(DomainPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(DomainPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(DomainPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(DomainPolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(DomainPolicyUserLoginMustBeDomainCol, policyEvent.UserLoginMustBeDomain),
			handler.NewCol(DomainPolicyIsDefaultCol, isDefault),
			handler.NewCol(DomainPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(DomainPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *DomainPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.DomainPolicyChangedEvent
	switch e := event.(type) {
	case *org.OrgDomainPolicyChangedEvent:
		policyEvent = e.DomainPolicyChangedEvent
	case *instance.InstanceDomainPolicyChangedEvent:
		policyEvent = e.DomainPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.OrgDomainPolicyChangedEventType, instance.InstanceDomainPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(DomainPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(DomainPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.UserLoginMustBeDomain != nil {
		cols = append(cols, handler.NewCol(DomainPolicyUserLoginMustBeDomainCol, *policyEvent.UserLoginMustBeDomain))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(DomainPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *DomainPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.OrgDomainPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-JAENd", "reduce.wrong.event.type %s", org.OrgDomainPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(DomainPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
