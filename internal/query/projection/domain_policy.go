package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	DomainPolicyTable = "projections.domain_policies2"

	DomainPolicyIDCol                                     = "id"
	DomainPolicyCreationDateCol                           = "creation_date"
	DomainPolicyChangeDateCol                             = "change_date"
	DomainPolicySequenceCol                               = "sequence"
	DomainPolicyStateCol                                  = "state"
	DomainPolicyUserLoginMustBeDomainCol                  = "user_login_must_be_domain"
	DomainPolicyValidateOrgDomainsCol                     = "validate_org_domains"
	DomainPolicySMTPSenderAddressMatchesInstanceDomainCol = "smtp_sender_address_matches_instance_domain"
	DomainPolicyIsDefaultCol                              = "is_default"
	DomainPolicyResourceOwnerCol                          = "resource_owner"
	DomainPolicyInstanceIDCol                             = "instance_id"
	DomainPolicyOwnerRemovedCol                           = "owner_removed"
)

type domainPolicyProjection struct {
	crdb.StatementHandler
}

func newDomainPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *domainPolicyProjection {
	p := new(domainPolicyProjection)
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
			crdb.NewColumn(DomainPolicyValidateOrgDomainsCol, crdb.ColumnTypeBool),
			crdb.NewColumn(DomainPolicySMTPSenderAddressMatchesInstanceDomainCol, crdb.ColumnTypeBool),
			crdb.NewColumn(DomainPolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(DomainPolicyResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(DomainPolicyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(DomainPolicyOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(DomainPolicyInstanceIDCol, DomainPolicyIDCol),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{DomainPolicyOwnerRemovedCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *domainPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.DomainPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.DomainPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  org.DomainPolicyRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.DomainPolicyAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(DomainPolicyInstanceIDCol),
				},
			},
		},
	}
}

func (p *domainPolicyProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.DomainPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *org.DomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		isDefault = false
	case *instance.DomainPolicyAddedEvent:
		policyEvent = e.DomainPolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-CSE7A", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyAddedEventType, instance.DomainPolicyAddedEventType})
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
			handler.NewCol(DomainPolicyValidateOrgDomainsCol, policyEvent.ValidateOrgDomains),
			handler.NewCol(DomainPolicySMTPSenderAddressMatchesInstanceDomainCol, policyEvent.SMTPSenderAddressMatchesInstanceDomain),
			handler.NewCol(DomainPolicyIsDefaultCol, isDefault),
			handler.NewCol(DomainPolicyResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(DomainPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *domainPolicyProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.DomainPolicyChangedEvent
	switch e := event.(type) {
	case *org.DomainPolicyChangedEvent:
		policyEvent = e.DomainPolicyChangedEvent
	case *instance.DomainPolicyChangedEvent:
		policyEvent = e.DomainPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-qgVug", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyChangedEventType, instance.DomainPolicyChangedEventType})
	}
	cols := []handler.Column{
		handler.NewCol(DomainPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(DomainPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.UserLoginMustBeDomain != nil {
		cols = append(cols, handler.NewCol(DomainPolicyUserLoginMustBeDomainCol, *policyEvent.UserLoginMustBeDomain))
	}
	if policyEvent.ValidateOrgDomains != nil {
		cols = append(cols, handler.NewCol(DomainPolicyValidateOrgDomainsCol, *policyEvent.ValidateOrgDomains))
	}
	if policyEvent.SMTPSenderAddressMatchesInstanceDomain != nil {
		cols = append(cols, handler.NewCol(DomainPolicySMTPSenderAddressMatchesInstanceDomainCol, *policyEvent.SMTPSenderAddressMatchesInstanceDomain))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(DomainPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(DomainPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *domainPolicyProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	policyEvent, ok := event.(*org.DomainPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-JAENd", "reduce.wrong.event.type %s", org.DomainPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		policyEvent,
		[]handler.Condition{
			handler.NewCond(DomainPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(DomainPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		}), nil
}

func (p *domainPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-JYD2K", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(DomainPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
