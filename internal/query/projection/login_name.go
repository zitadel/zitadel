package projection

import (
	"context"
	"fmt"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
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

	loginNameDomainSuffix           = "domains"
	LoginNameDomainNameCol          = "name"
	LoginNameDomainIsPrimaryCol     = "is_primary"
	LoginNameDomainResourceOwnerCol = "resource_owner"

	loginNamePolicySuffix             = "policies"
	LoginNamePoliciesMustBeDomainCol  = "must_be_domain"
	LoginNamePoliciesIsDefaultCol     = "is_default"
	LoginNamePoliciesResourceOwnerCol = "resource_owner"
)

var (
	viewStmt = fmt.Sprintf("SELECT"+
		" user_id"+
		" , IF(%[1]s, CONCAT(%[2]s, '@', domain), %[2]s) AS login_name"+
		" , IFNULL(%[3]s, true) AS %[3]s"+
		" FROM ("+
		" SELECT"+
		" policy_users.user_id"+
		" , policy_users.%[2]s"+
		" , policy_users.%[4]s"+
		" , policy_users.%[1]s"+
		" , domains.%[5]s AS domain"+
		" , domains.%[3]s"+
		" FROM ("+
		" SELECT"+
		" users.id as user_id"+
		" , users.%[2]s"+
		" , users.%[4]s"+
		" , IFNULL(policy_custom.%[1]s, policy_default.%[1]s) AS %[1]s"+
		" FROM %[6]s users"+
		" LEFT JOIN %[7]s policy_custom on policy_custom.%[8]s = users.%[4]s"+
		" LEFT JOIN %[7]s policy_default on policy_default.%[9]s = true) policy_users"+
		" LEFT JOIN %[10]s domains ON policy_users.%[1]s AND policy_users.%[4]s = domains.%[11]s"+
		");",
		LoginNamePoliciesMustBeDomainCol,
		LoginNameUserUserNameCol,
		LoginNameDomainIsPrimaryCol,
		LoginNameUserResourceOwnerCol,
		LoginNameDomainNameCol,
		LoginNameUserProjectionTable,
		LoginNamePolicyProjectionTable,
		LoginNamePoliciesResourceOwnerCol,
		LoginNamePoliciesIsDefaultCol,
		LoginNameDomainProjectionTable,
		LoginNameDomainResourceOwnerCol,
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
		},
			crdb.NewPrimaryKey(LoginNameUserIDCol),
			loginNameUserSuffix,
			crdb.NewIndex("ro_idx", []string{LoginNameUserResourceOwnerCol}),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNameDomainNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameDomainIsPrimaryCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNameDomainResourceOwnerCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(LoginNameDomainResourceOwnerCol, LoginNameDomainNameCol),
			loginNameDomainSuffix,
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNamePoliciesMustBeDomainCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNamePoliciesIsDefaultCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNamePoliciesResourceOwnerCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(LoginNamePoliciesResourceOwnerCol),
			loginNamePolicySuffix,
			crdb.NewIndex("is_default_idx", []string{LoginNamePoliciesResourceOwnerCol, LoginNamePoliciesIsDefaultCol}),
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-yCV6S", "reduce.wrong.event.type %v", []eventstore.EventType{org.OrgIAMPolicyAddedEventType, iam.OrgIAMPolicyAddedEventType})
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

func (p *LoginNameProjection) reduceOrgIAMPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent *policy.OrgIAMPolicyChangedEvent

	switch e := event.(type) {
	case *org.OrgIAMPolicyChangedEvent:
		policyEvent = &e.OrgIAMPolicyChangedEvent
	case *iam.OrgIAMPolicyChangedEvent:
		policyEvent = &e.OrgIAMPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ArFDd", "reduce.wrong.event.type %v", []eventstore.EventType{org.OrgIAMPolicyChangedEventType, iam.OrgIAMPolicyChangedEventType})
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

func (p *LoginNameProjection) reduceOrgIAMPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgIAMPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ysEeB", "reduce.wrong.event.type %s", org.OrgIAMPolicyRemovedEventType)
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
