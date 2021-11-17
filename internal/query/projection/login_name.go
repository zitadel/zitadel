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
	"github.com/caos/zitadel/internal/repository/user"
)

type LoginNameProjection struct {
	crdb.StatementHandler
}

const (
	LoginNameProjectionTable       = "zitadel.projections.login_names"
	LoginNameUserProjectionTable   = "zitadel.projections.login_names" + "_" + loginNameUserSuffix
	LoginNamePolicyProjectionTable = "zitadel.projections.login_names" + "_" + loginNamePolicySuffix
	LoginNameDomainProjectionTable = "zitadel.projections.login_names" + "_" + loginNameDomainSuffix
)

func NewLoginNameProjection(ctx context.Context, config crdb.StatementHandlerConfig) *LoginNameProjection {
	p := &LoginNameProjection{}
	config.ProjectionName = LoginNameProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *LoginNameProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.HumanAddedType,
					Reduce: p.reduceHumanAdded,
				},
				{
					Event:  user.HumanRegisteredType,
					Reduce: p.reduceHumanRegistered,
				},
				{
					Event:  user.HumanEmailChangedType,
					Reduce: p.reduceEmailChanged,
				},
				// {
				// 	Event: user.HumanEmailVerifiedType,
				// 	// Reduce: p.reduceEmailVerified,
				// 	// email is changed as soon as the email changed
				// },
				{
					Event:  user.MachineAddedEventType,
					Reduce: p.reduceMachineAdded,
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
					Event:  org.OrgDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
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
	loginNameUserSuffix   = "users"
	loginNamePolicySuffix = "policies"
	loginNameDomainSuffix = "domains"

	LoginNameUserIDCol            = "id"
	LoginNameUserTypeCol          = "type"
	LoginNameUserUserNameCol      = "name"
	LoginNameUserEmailCol         = "email"
	LoginNameUserResourceOwnerCol = "resource_owner"

	LoginNameDomainNameCol          = "name"
	LoginNameDomainIsPrimaryCol     = "is_primary"
	LoginNameDomainIsVerifiedCol    = "is_verified"
	LoginNameDomainResourceOwnerCol = "resource_owner"

	LoginNamePoliciesMustBeDomainCol  = "must_be_domain"
	LoginNamePoliciesIsDefaultCol     = "is_default"
	LoginNamePoliciesResourceOwnerCol = "resource_owner"
)

func (p *LoginNameProjection) reduceHumanAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.HumanAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ayo69", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCol(LoginNameUserTypeCol, domain.UserTypeHuman),
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
			handler.NewCol(LoginNameUserEmailCol, e.EmailAddress),
			handler.NewCol(LoginNameUserResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *LoginNameProjection) reduceHumanRegistered(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.HumanRegisteredEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.HumanRegisteredType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-psrvi", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCol(LoginNameUserTypeCol, domain.UserTypeHuman),
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
			handler.NewCol(LoginNameUserEmailCol, e.EmailAddress),
			handler.NewCol(LoginNameUserResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *LoginNameProjection) reduceEmailChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.HumanEmailChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-kAeEo", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserEmailCol, e.EmailAddress),
		},
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *LoginNameProjection) reduceMachineAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.MachineAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.MachineAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-65IL8", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCol(LoginNameUserTypeCol, domain.UserTypeMachine),
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
			handler.NewCol(LoginNameUserResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *LoginNameProjection) reduceUserRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
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

func (p *LoginNameProjection) reduceUserNameChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.UserUserNameChangedType).Error("wrong event type")
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

func (p *LoginNameProjection) reduceUserDomainClaimed(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", user.UserDomainClaimedType).Error("wrong event type")
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

func (p *LoginNameProjection) reduceOrgIAMPolicyAdded(event eventstore.EventReader) (*handler.Statement, error) {
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
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyAddedEventType, iam.OrgIAMPolicyAddedEventType}).Error("wrong event type")
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

func (p *LoginNameProjection) reduceOrgIAMPolicyChanged(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent *policy.OrgIAMPolicyChangedEvent

	switch e := event.(type) {
	case *org.OrgIAMPolicyChangedEvent:
		policyEvent = &e.OrgIAMPolicyChangedEvent
	case *iam.OrgIAMPolicyChangedEvent:
		policyEvent = &e.OrgIAMPolicyChangedEvent
	default:
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyChangedEventType, iam.OrgIAMPolicyChangedEventType}).Error("wrong event type")
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

func (p *LoginNameProjection) reduceOrgIAMPolicyRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.OrgIAMPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgIAMPolicyRemovedEventType).Error("wrong event type")
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

func (p *LoginNameProjection) reduceDomainAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgDomainAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-JKUlP", "reduce.wrong.event.type")
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

func (p *LoginNameProjection) reducePrimaryDomainSet(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgDomainPrimarySetEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-t82ow", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(LoginNameDomainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNameUserSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(LoginNameDomainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(LoginNameDomainNameCol, e.Domain),
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNameDomainSuffix),
		),
	), nil
}

func (p *LoginNameProjection) reduceDomainRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgDomainRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-8S9aY", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameDomainNameCol, e.Domain),
			handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameDomainSuffix),
	), nil
}

func (p *LoginNameProjection) reduceDomainVerified(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		logging.LogWithFields("HANDL-zWCk3", "seq", event.Sequence(), "expectedType", org.OrgDomainVerifiedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-weGAh", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameDomainIsVerifiedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(LoginNameDomainNameCol, e.Domain),
			handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameDomainSuffix),
	), nil
}
