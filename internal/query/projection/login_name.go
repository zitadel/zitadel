package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type loginNameProjection struct {
	crdb.StatementHandler
}

const (
	LoginNameProjectionTable       = "zitadel.projections.login_names"
	LoginNameUserProjectionTable   = LoginNameProjectionTable + "_" + loginNameUserSuffix
	LoginNamePolicyProjectionTable = LoginNameProjectionTable + "_" + loginNamePolicySuffix
	LoginNameDomainProjectionTable = LoginNameProjectionTable + "_" + loginNameDomainSuffix
)

func newLoginNameProjection(ctx context.Context, config crdb.StatementHandlerConfig) *loginNameProjection {
	p := &loginNameProjection{}
	config.ProjectionName = LoginNameProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *loginNameProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserV1AddedType,
					Reduce: p.reduceUserCreated,
				},
				{
					Event:  user.HumanAddedType,
					Reduce: p.reduceUserCreated,
				},
				{
					Event:  user.HumanRegisteredType,
					Reduce: p.reduceUserCreated,
				},
				{
					Event:  user.UserV1RegisteredType,
					Reduce: p.reduceUserCreated,
				},
				{
					Event:  user.MachineAddedEventType,
					Reduce: p.reduceUserCreated,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
				{
					Event:  user.UserUserNameChangedType,
					Reduce: p.reduceUserNameChanged,
				},
				{
					// changes the username of the user
					// this event occures in orgs
					// where policy.must_be_domain=false
					Event:  user.UserDomainClaimedType,
					Reduce: p.reduceUserDomainClaimed,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgIAMPolicyAddedEventType,
					Reduce: p.reduceOrgIAMPolicyAdded,
				},
				{
					Event:  org.OrgIAMPolicyChangedEventType,
					Reduce: p.reduceOrgIAMPolicyChanged,
				},
				{
					Event:  org.OrgIAMPolicyRemovedEventType,
					Reduce: p.reduceOrgIAMPolicyRemoved,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reducePrimaryDomainSet,
				},
				{
					Event:  org.OrgDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
				{
					Event:  org.OrgDomainVerifiedEventType,
					Reduce: p.reduceDomainVerified,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.OrgIAMPolicyAddedEventType,
					Reduce: p.reduceOrgIAMPolicyAdded,
				},
				{
					Event:  iam.OrgIAMPolicyChangedEventType,
					Reduce: p.reduceOrgIAMPolicyChanged,
				},
			},
		},
	}
}

const (
	LoginNameCol = "login_name"

	loginNameUserSuffix           = "users"
	LoginNameUserIDCol            = "id"
	LoginNameUserUserNameCol      = "user_name"
	LoginNameUserResourceOwnerCol = "resource_owner"

	loginNameDomainSuffix           = "domains"
	LoginNameDomainNameCol          = "name"
	LoginNameDomainIsPrimaryCol     = "is_primary"
	LoginNameDomainResourceOwnerCol = "resource_owner"

	loginNamePolicySuffix             = "policies"
	LoginNamePoliciesMustBeDomainCol  = "must_be_domain"
	LoginNamePoliciesIsDefaultCol     = "is_default"
	LoginNamePoliciesResourceOwnerCol = "resource_owner"
)

func (p *loginNameProjection) reduceUserCreated(event eventstore.Event) (*handler.Statement, error) {
	var userName string

	switch e := event.(type) {
	case *user.HumanAddedEvent:
		userName = e.UserName
	case *user.HumanRegisteredEvent:
		userName = e.UserName
	case *user.MachineAddedEvent:
		userName = e.UserName
	default:
		logging.LogWithFields("HANDL-tDUx3", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{user.UserV1AddedType, user.HumanAddedType, user.UserV1RegisteredType, user.HumanRegisteredType, user.MachineAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ayo69", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserIDCol, event.Aggregate().ID),
			handler.NewCol(LoginNameUserUserNameCol, userName),
			handler.NewCol(LoginNameUserResourceOwnerCol, event.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-8XEdC", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-QIe3C", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-UGo7U", "seq", event.Sequence(), "expectedType", user.UserUserNameChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-QlwjC", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
		},
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zIbyU", "seq", event.Sequence(), "expectedType", user.UserDomainClaimedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-AQMBY", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
		},
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceOrgIAMPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var (
		policyEvent *policy.OrgIAMPolicyAddedEvent
		isDefault   bool
	)

	switch e := event.(type) {
	case *org.OrgIAMPolicyAddedEvent:
		policyEvent = &e.OrgIAMPolicyAddedEvent
		isDefault = false
	case *iam.OrgIAMPolicyAddedEvent:
		policyEvent = &e.OrgIAMPolicyAddedEvent
		isDefault = true
	default:
		logging.LogWithFields("HANDL-PQluH", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyAddedEventType, iam.OrgIAMPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-yCV6S", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNamePoliciesMustBeDomainCol, policyEvent.UserLoginMustBeDomain),
			handler.NewCol(LoginNamePoliciesIsDefaultCol, isDefault),
			handler.NewCol(LoginNamePoliciesResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceOrgIAMPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent *policy.OrgIAMPolicyChangedEvent

	switch e := event.(type) {
	case *org.OrgIAMPolicyChangedEvent:
		policyEvent = &e.OrgIAMPolicyChangedEvent
	case *iam.OrgIAMPolicyChangedEvent:
		policyEvent = &e.OrgIAMPolicyChangedEvent
	default:
		logging.LogWithFields("HANDL-Z27QN", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyChangedEventType, iam.OrgIAMPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ArFDd", "reduce.wrong.event.type")
	}

	if policyEvent.UserLoginMustBeDomain == nil {
		return crdb.NewNoOpStatement(event), nil
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNamePoliciesMustBeDomainCol, *policyEvent.UserLoginMustBeDomain),
		},
		[]handler.Condition{
			handler.NewCond(LoginNamePoliciesResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceOrgIAMPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgIAMPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-1ZFHL", "seq", event.Sequence(), "expectedType", org.OrgIAMPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ysEeB", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNamePoliciesResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Rr7Tq", "seq", event.Sequence(), "expectedType", org.OrgDomainVerifiedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-weGAh", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameDomainNameCol, e.Domain),
			handler.NewCol(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNameDomainSuffix),
	), nil
}

func (p *loginNameProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("HANDL-0L5tW", "seq", event.Sequence(), "expectedType", org.OrgDomainPrimarySetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-eOXPN", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(LoginNameDomainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCond(LoginNameDomainIsPrimaryCol, true),
			},
			crdb.WithTableSuffix(loginNameDomainSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(LoginNameDomainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(LoginNameDomainNameCol, e.Domain),
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
			},
			crdb.WithTableSuffix(loginNameDomainSuffix),
		),
	), nil
}

func (p *loginNameProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-reP2u", "seq", event.Sequence(), "expectedType", org.OrgDomainRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-4RHYq", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameDomainNameCol, e.Domain),
			handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNameDomainSuffix),
	), nil
}
