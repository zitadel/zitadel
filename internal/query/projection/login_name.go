package projection

import (
	"context"
	"fmt"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	LoginNameProjectionTable       = "projections.login_names"
	LoginNameUserProjectionTable   = LoginNameProjectionTable + "_" + loginNameUserSuffix
	LoginNamePolicyProjectionTable = LoginNameProjectionTable + "_" + loginNamePolicySuffix
	LoginNameDomainProjectionTable = LoginNameProjectionTable + "_" + loginNameDomainSuffix

	LoginNameCol = "login_name"

	loginNameUserSuffix           = "users"
	LoginNameUserIDCol            = "id"
	LoginNameUserUserNameCol      = "user_name"
	LoginNameUserResourceOwnerCol = "resource_owner"
	LoginNameUserInstanceIDCol    = "instance_id"

	loginNameDomainSuffix           = "domains"
	LoginNameDomainNameCol          = "name"
	LoginNameDomainIsPrimaryCol     = "is_primary"
	LoginNameDomainResourceOwnerCol = "resource_owner"
	LoginNameDomainInstanceIDCol    = "instance_id"

	loginNamePolicySuffix             = "policies"
	LoginNamePoliciesMustBeDomainCol  = "must_be_domain"
	LoginNamePoliciesIsDefaultCol     = "is_default"
	LoginNamePoliciesResourceOwnerCol = "resource_owner"
	LoginNamePoliciesInstanceIDCol    = "instance_id"
)

var (
	viewStmt = fmt.Sprintf("SELECT"+
		" user_id"+
		" , IF(%[1]s, CONCAT(%[2]s, '@', domain), %[2]s) AS login_name"+
		" , IFNULL(%[3]s, true) AS %[3]s"+
		" , %[4]s"+
		" FROM ("+
		" SELECT"+
		" policy_users.user_id"+
		" , policy_users.%[2]s"+
		" , policy_users.%[5]s"+
		" , policy_users.%[4]s"+
		" , policy_users.%[1]s"+
		" , domains.%[6]s AS domain"+
		" , domains.%[3]s"+
		" FROM ("+
		" SELECT"+
		" users.id as user_id"+
		" , users.%[2]s"+
		" , users.%[4]s"+
		" , users.%[5]s"+
		" , IFNULL(policy_custom.%[1]s, policy_default.%[1]s) AS %[1]s"+
		" FROM %[7]s users"+
		" LEFT JOIN %[8]s policy_custom on policy_custom.%[9]s = users.%[5]s AND policy_custom.%[10]s = users.%[4]s"+
		" LEFT JOIN %[8]s policy_default on policy_default.%[11]s = true) policy_users"+
		" LEFT JOIN %[12]s domains ON policy_users.%[1]s AND policy_users.%[5]s = domains.%[13]s AND policy_users.%[10]s = domains.%[14]s"+
		");",
		LoginNamePoliciesMustBeDomainCol,
		LoginNameUserUserNameCol,
		LoginNameDomainIsPrimaryCol,
		LoginNameUserInstanceIDCol,
		LoginNameUserResourceOwnerCol,
		LoginNameDomainNameCol,
		LoginNameUserProjectionTable,
		LoginNamePolicyProjectionTable,
		LoginNamePoliciesResourceOwnerCol,
		LoginNamePoliciesInstanceIDCol,
		LoginNamePoliciesIsDefaultCol,
		LoginNameDomainProjectionTable,
		LoginNameDomainResourceOwnerCol,
		LoginNameDomainInstanceIDCol,
	)
)

type LoginNameProjection struct {
	crdb.StatementHandler
}

func NewLoginNameProjection(ctx context.Context, config crdb.StatementHandlerConfig) *LoginNameProjection {
	p := new(LoginNameProjection)
	config.ProjectionName = LoginNameProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewViewCheck(
		viewStmt,
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNameUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserUserNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserInstanceIDCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(LoginNameUserInstanceIDCol, LoginNameUserIDCol),
			loginNameUserSuffix,
			crdb.WithIndex(crdb.NewIndex("ro_idx", []string{LoginNameUserResourceOwnerCol})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNameDomainNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameDomainIsPrimaryCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNameDomainResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameDomainInstanceIDCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(LoginNameDomainInstanceIDCol, LoginNameDomainResourceOwnerCol, LoginNameDomainNameCol),
			loginNameDomainSuffix,
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNamePoliciesMustBeDomainCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNamePoliciesIsDefaultCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNamePoliciesResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNamePoliciesInstanceIDCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(LoginNamePoliciesInstanceIDCol, LoginNamePoliciesResourceOwnerCol),
			loginNamePolicySuffix,
			crdb.WithIndex(crdb.NewIndex("is_default_idx", []string{LoginNamePoliciesResourceOwnerCol, LoginNamePoliciesIsDefaultCol})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *LoginNameProjection) reducers() []handler.AggregateReducer {
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
					Event:  org.DomainPolicyAddedEventType,
					Reduce: p.reduceOrgIAMPolicyAdded,
				},
				{
					Event:  org.DomainPolicyChangedEventType,
					Reduce: p.reduceDomainPolicyChanged,
				},
				{
					Event:  org.DomainPolicyRemovedEventType,
					Reduce: p.reduceDomainPolicyRemoved,
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
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.DomainPolicyAddedEventType,
					Reduce: p.reduceOrgIAMPolicyAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: p.reduceDomainPolicyChanged,
				},
			},
		},
	}
}

func (p *LoginNameProjection) reduceUserCreated(event eventstore.Event) (*handler.Statement, error) {
	var userName string

	switch e := event.(type) {
	case *user.HumanAddedEvent:
		userName = e.UserName
	case *user.HumanRegisteredEvent:
		userName = e.UserName
	case *user.MachineAddedEvent:
		userName = e.UserName
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ayo69", "reduce.wrong.event.type %v", []eventstore.EventType{user.UserV1AddedType, user.HumanAddedType, user.UserV1RegisteredType, user.HumanRegisteredType, user.MachineAddedEventType})
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserIDCol, event.Aggregate().ID),
			handler.NewCol(LoginNameUserUserNameCol, userName),
			handler.NewCol(LoginNameUserResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCol(LoginNameUserInstanceIDCol, event.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *LoginNameProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-QIe3C", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *LoginNameProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-QlwjC", "reduce.wrong.event.type %s", user.UserUserNameChangedType)
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

func (p *LoginNameProjection) reduceUserDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-AQMBY", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
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

func (p *LoginNameProjection) reduceOrgIAMPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var (
		policyEvent *policy.DomainPolicyAddedEvent
		isDefault   bool
	)

	switch e := event.(type) {
	case *org.DomainPolicyAddedEvent:
		policyEvent = &e.DomainPolicyAddedEvent
		isDefault = false
	case *instance.DomainPolicyAddedEvent:
		policyEvent = &e.DomainPolicyAddedEvent
		isDefault = true
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-yCV6S", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyAddedEventType, instance.DomainPolicyAddedEventType})
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNamePoliciesMustBeDomainCol, policyEvent.UserLoginMustBeDomain),
			handler.NewCol(LoginNamePoliciesIsDefaultCol, isDefault),
			handler.NewCol(LoginNamePoliciesResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(LoginNamePoliciesInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *LoginNameProjection) reduceDomainPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent *policy.DomainPolicyChangedEvent

	switch e := event.(type) {
	case *org.DomainPolicyChangedEvent:
		policyEvent = &e.DomainPolicyChangedEvent
	case *instance.DomainPolicyChangedEvent:
		policyEvent = &e.DomainPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ArFDd", "reduce.wrong.event.type %v", []eventstore.EventType{org.DomainPolicyChangedEventType, instance.DomainPolicyChangedEventType})
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

func (p *LoginNameProjection) reduceDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ysEeB", "reduce.wrong.event.type %s", org.DomainPolicyRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNamePoliciesResourceOwnerCol, e.Aggregate().ResourceOwner),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *LoginNameProjection) reduceDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-weGAh", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameDomainNameCol, e.Domain),
			handler.NewCol(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNameDomainSuffix),
	), nil
}

func (p *LoginNameProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eOXPN", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
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

func (p *LoginNameProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-4RHYq", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
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
