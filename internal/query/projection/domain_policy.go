package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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

type domainPolicyProjection struct{}

func newDomainPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(domainPolicyProjection))
}

func (*domainPolicyProjection) Name() string {
	return DomainPolicyTable
}

func (*domainPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DomainPolicyIDCol, handler.ColumnTypeText),
			handler.NewColumn(DomainPolicyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DomainPolicyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(DomainPolicySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(DomainPolicyStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(DomainPolicyUserLoginMustBeDomainCol, handler.ColumnTypeBool),
			handler.NewColumn(DomainPolicyValidateOrgDomainsCol, handler.ColumnTypeBool),
			handler.NewColumn(DomainPolicySMTPSenderAddressMatchesInstanceDomainCol, handler.ColumnTypeBool),
			handler.NewColumn(DomainPolicyIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(DomainPolicyResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(DomainPolicyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(DomainPolicyOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(DomainPolicyInstanceIDCol, DomainPolicyIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{DomainPolicyOwnerRemovedCol})),
		),
	)
}

func (p *domainPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
	return handler.NewCreateStatement(
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
	return handler.NewUpdateStatement(
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
	return handler.NewDeleteStatement(
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

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(DomainPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(DomainPolicyResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}
