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
	OrgIAMPolicyTable = "projections.org_iam_policies"

	OrgIAMPolicyIDCol                    = "id"
	OrgIAMPolicyCreationDateCol          = "creation_date"
	OrgIAMPolicyChangeDateCol            = "change_date"
	OrgIAMPolicySequenceCol              = "sequence"
	OrgIAMPolicyStateCol                 = "state"
	OrgIAMPolicyUserLoginMustBeDomainCol = "user_login_must_be_domain"
	OrgIAMPolicyIsDefaultCol             = "is_default"
	OrgIAMPolicyResourceOwnerCol         = "resource_owner"
	OrgIAMPolicyInstanceCol              = "instance_id"
)

type OrgIAMPolicyProjection struct {
	crdb.StatementHandler
}

func NewOrgIAMPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OrgIAMPolicyProjection {
	p := new(OrgIAMPolicyProjection)
	config.ProjectionName = OrgIAMPolicyTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(OrgIAMPolicyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OrgIAMPolicyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(OrgIAMPolicyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(OrgIAMPolicySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(OrgIAMPolicyStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(OrgIAMPolicyUserLoginMustBeDomainCol, crdb.ColumnTypeBool),
			crdb.NewColumn(OrgIAMPolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(OrgIAMPolicyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(OrgIAMPolicyInstanceCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(OrgIAMPolicyIDCol),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OrgIAMPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgIAMPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.OrgIAMPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.OrgIAMPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.OrgIAMPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  iam.OrgIAMPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
			},
		},
	}
}

func (p *OrgIAMPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.OrgIAMPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.OrgIAMPolicyAddedEvent:
		policyEvent = e.OrgIAMPolicyAddedEvent
		isDefault = false
	case *iam.OrgIAMPolicyAddedEvent:
		policyEvent = e.OrgIAMPolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.OrgIAMPolicyAddedEventType, iam.OrgIAMPolicyAddedEventType})
	}
	return crdb.NewCreateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(OrgIAMPolicyCreationDateCol, policyEvent.CreationDate()),
			handler.NewCol(OrgIAMPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(OrgIAMPolicySequenceCol, policyEvent.Sequence()),
			handler.NewCol(OrgIAMPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCol(OrgIAMPolicyStateCol, domain.PolicyStateActive),
			handler.NewCol(OrgIAMPolicyUserLoginMustBeDomainCol, policyEvent.UserLoginMustBeDomain),
			handler.NewCol(OrgIAMPolicyIsDefaultCol, isDefault),
			handler.NewCol(OrgIAMPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(OrgIAMPolicyInstanceCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *OrgIAMPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.OrgIAMPolicyChangedEvent
	switch e := event.(type) {
	case *org.OrgIAMPolicyChangedEvent:
		policyEvent = e.OrgIAMPolicyChangedEvent
	case *iam.OrgIAMPolicyChangedEvent:
		policyEvent = e.OrgIAMPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.OrgIAMPolicyChangedEventType, iam.OrgIAMPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(OrgIAMPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(OrgIAMPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.UserLoginMustBeDomain != nil {
		cols = append(cols, handler.NewCol(OrgIAMPolicyUserLoginMustBeDomainCol, *policyEvent.UserLoginMustBeDomain))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(OrgIAMPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}

func (p *OrgIAMPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.OrgIAMPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-JAENd", "reduce.wrong.event.type %s", org.OrgIAMPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(OrgIAMPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
