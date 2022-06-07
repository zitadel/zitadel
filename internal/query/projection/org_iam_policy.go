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

type orgIAMPolicyProjection struct {
	crdb.StatementHandler
}

const (
	OrgIAMPolicyTable = "zitadel.projections.org_iam_policies"

	OrgIAMPolicyCreationDateCol          = "creation_date"
	OrgIAMPolicyChangeDateCol            = "change_date"
	OrgIAMPolicySequenceCol              = "sequence"
	OrgIAMPolicyIDCol                    = "id"
	OrgIAMPolicyStateCol                 = "state"
	OrgIAMPolicyUserLoginMustBeDomainCol = "user_login_must_be_domain"
	OrgIAMPolicyIsDefaultCol             = "is_default"
	OrgIAMPolicyResourceOwnerCol         = "resource_owner"
)

func newOrgIAMPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *orgIAMPolicyProjection {
	p := &orgIAMPolicyProjection{}
	config.ProjectionName = OrgIAMPolicyTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *orgIAMPolicyProjection) reducers() []handler.AggregateReducer {
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

func (p *orgIAMPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
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
		logging.LogWithFields("PROJE-XakxJ", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyAddedEventType, iam.OrgIAMPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-CSE7A", "reduce.wrong.event.type")
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
		}), nil
}

func (p *orgIAMPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.OrgIAMPolicyChangedEvent
	switch e := event.(type) {
	case *org.OrgIAMPolicyChangedEvent:
		policyEvent = e.OrgIAMPolicyChangedEvent
	case *iam.OrgIAMPolicyChangedEvent:
		policyEvent = e.OrgIAMPolicyChangedEvent
	default:
		logging.LogWithFields("PROJE-SvTK0", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyChangedEventType, iam.OrgIAMPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-qgVug", "reduce.wrong.event.type")
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

func (p *orgIAMPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.OrgIAMPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("PROJE-ovQya", "seq", event.Sequence(), "expectedType", org.OrgIAMPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "PROJE-JAENd", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(OrgIAMPolicyIDCol, policyEvent.Aggregate().ID),
		}), nil
}
