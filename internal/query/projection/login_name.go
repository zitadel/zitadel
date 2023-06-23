package projection

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	LoginNameTableAlias            = "login_names2"
	LoginNameProjectionTable       = "projections." + LoginNameTableAlias
	LoginNameUserProjectionTable   = LoginNameProjectionTable + "_" + loginNameUserSuffix
	LoginNamePolicyProjectionTable = LoginNameProjectionTable + "_" + loginNamePolicySuffix
	LoginNameDomainProjectionTable = LoginNameProjectionTable + "_" + loginNameDomainSuffix

	LoginNameCol                   = "login_name"
	LoginNameUserCol               = "user_id"
	LoginNameIsPrimaryCol          = "is_primary"
	LoginNameInstanceIDCol         = "instance_id"
	LoginNameOwnerRemovedUserCol   = "user_owner_removed"
	LoginNameOwnerRemovedPolicyCol = "policy_owner_removed"
	LoginNameOwnerRemovedDomainCol = "domain_owner_removed"

	usersAlias         = "users"
	policyCustomAlias  = "policy_custom"
	policyDefaultAlias = "policy_default"
	policyUsersAlias   = "policy_users"
	domainsAlias       = "domains"
	domainAlias        = "domain"

	loginNameUserSuffix           = "users"
	LoginNameUserIDCol            = "id"
	LoginNameUserUserNameCol      = "user_name"
	LoginNameUserResourceOwnerCol = "resource_owner"
	LoginNameUserInstanceIDCol    = "instance_id"
	LoginNameUserOwnerRemovedCol  = "owner_removed"

	loginNameDomainSuffix           = "domains"
	LoginNameDomainNameCol          = "name"
	LoginNameDomainIsPrimaryCol     = "is_primary"
	LoginNameDomainResourceOwnerCol = "resource_owner"
	LoginNameDomainInstanceIDCol    = "instance_id"
	LoginNameDomainOwnerRemovedCol  = "owner_removed"

	loginNamePolicySuffix             = "policies"
	LoginNamePoliciesMustBeDomainCol  = "must_be_domain"
	LoginNamePoliciesIsDefaultCol     = "is_default"
	LoginNamePoliciesResourceOwnerCol = "resource_owner"
	LoginNamePoliciesInstanceIDCol    = "instance_id"
	LoginNamePoliciesOwnerRemovedCol  = "owner_removed"
)

var (
	policyUsers = sq.Select(
		alias(
			col(usersAlias, LoginNameUserIDCol),
			LoginNameUserCol,
		),
		col(usersAlias, LoginNameUserUserNameCol),
		col(usersAlias, LoginNameUserInstanceIDCol),
		col(usersAlias, LoginNameUserResourceOwnerCol),
		alias(
			coalesce(col(policyCustomAlias, LoginNamePoliciesMustBeDomainCol), col(policyDefaultAlias, LoginNamePoliciesMustBeDomainCol)),
			LoginNamePoliciesMustBeDomainCol,
		),
		alias(col(usersAlias, LoginNameUserOwnerRemovedCol),
			LoginNameOwnerRemovedUserCol),
		alias(coalesce(col(policyCustomAlias, LoginNamePoliciesOwnerRemovedCol), "false"),
			LoginNameOwnerRemovedPolicyCol),
	).From(alias(LoginNameUserProjectionTable, usersAlias)).
		LeftJoin(
			leftJoin(LoginNamePolicyProjectionTable, policyCustomAlias,
				eq(col(policyCustomAlias, LoginNamePoliciesResourceOwnerCol), col(usersAlias, LoginNameUserResourceOwnerCol)),
				eq(col(policyCustomAlias, LoginNamePoliciesInstanceIDCol), col(usersAlias, LoginNameUserInstanceIDCol)),
			),
		).
		LeftJoin(
			leftJoin(LoginNamePolicyProjectionTable, policyDefaultAlias,
				eq(col(policyDefaultAlias, LoginNamePoliciesIsDefaultCol), "true"),
				eq(col(policyDefaultAlias, LoginNamePoliciesInstanceIDCol), col(usersAlias, LoginNameUserInstanceIDCol)),
			),
		)

	loginNamesTable = sq.Select(
		col(policyUsersAlias, LoginNameUserCol),
		col(policyUsersAlias, LoginNameUserUserNameCol),
		col(policyUsersAlias, LoginNameUserResourceOwnerCol),
		alias(col(policyUsersAlias, LoginNameUserInstanceIDCol),
			LoginNameInstanceIDCol),
		col(policyUsersAlias, LoginNamePoliciesMustBeDomainCol),
		alias(col(domainsAlias, LoginNameDomainNameCol),
			domainAlias),
		col(domainsAlias, LoginNameDomainIsPrimaryCol),
		col(policyUsersAlias, LoginNameOwnerRemovedUserCol),
		col(policyUsersAlias, LoginNameOwnerRemovedPolicyCol),
		alias(coalesce(col(domainsAlias, LoginNameDomainOwnerRemovedCol), "false"),
			LoginNameOwnerRemovedDomainCol),
	).FromSelect(policyUsers, policyUsersAlias).
		LeftJoin(
			leftJoin(LoginNameDomainProjectionTable, domainsAlias,
				col(policyUsersAlias, LoginNamePoliciesMustBeDomainCol),
				eq(col(policyUsersAlias, LoginNameUserResourceOwnerCol), col(domainsAlias, LoginNameDomainResourceOwnerCol)),
				eq(col(policyUsersAlias, LoginNamePoliciesInstanceIDCol), col(domainsAlias, LoginNameDomainInstanceIDCol)),
			),
		)

	viewStmt, _ = sq.Select(
		LoginNameUserCol,
		alias(
			whenThenElse(
				LoginNamePoliciesMustBeDomainCol,
				concat(LoginNameUserUserNameCol, "'@'", domainAlias),
				LoginNameUserUserNameCol),
			LoginNameCol),
		alias(coalesce(LoginNameDomainIsPrimaryCol, "true"),
			LoginNameIsPrimaryCol),
		LoginNameInstanceIDCol,
		LoginNameOwnerRemovedUserCol,
		LoginNameOwnerRemovedPolicyCol,
		LoginNameOwnerRemovedDomainCol,
	).FromSelect(loginNamesTable, LoginNameTableAlias).MustSql()
)

func col(table, name string) string {
	return table + "." + name
}

func alias(col, alias string) string {
	return col + " AS " + alias
}

func coalesce(values ...string) string {
	str := "COALESCE("
	for i, value := range values {
		if i > 0 {
			str += ", "
		}
		str += value
	}
	str += ")"
	return str
}

func eq(first, second string) string {
	return first + " = " + second
}

func leftJoin(table, alias, on string, and ...string) string {
	st := table + " " + alias + " ON " + on
	for _, a := range and {
		st += " AND " + a
	}
	return st
}

func concat(strs ...string) string {
	return "CONCAT(" + strings.Join(strs, ", ") + ")"
}

func whenThenElse(when, then, el string) string {
	return "(CASE WHEN " + when + " THEN " + then + " ELSE " + el + " END)"
}

type loginNameProjection struct {
	crdb.StatementHandler
}

func newLoginNameProjection(ctx context.Context, config crdb.StatementHandlerConfig) *loginNameProjection {
	p := new(loginNameProjection)
	config.ProjectionName = LoginNameProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewViewCheck(
		viewStmt,
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNameUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserUserNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameUserOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(LoginNameUserInstanceIDCol, LoginNameUserIDCol),
			loginNameUserSuffix,
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{LoginNameUserResourceOwnerCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{LoginNameUserOwnerRemovedCol})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNameDomainNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameDomainIsPrimaryCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(LoginNameDomainResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameDomainInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNameDomainOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(LoginNameDomainInstanceIDCol, LoginNameDomainResourceOwnerCol, LoginNameDomainNameCol),
			loginNameDomainSuffix,
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{LoginNameDomainOwnerRemovedCol})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LoginNamePoliciesMustBeDomainCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNamePoliciesIsDefaultCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginNamePoliciesResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNamePoliciesInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginNamePoliciesOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(LoginNamePoliciesInstanceIDCol, LoginNamePoliciesResourceOwnerCol),
			loginNamePolicySuffix,
			crdb.WithIndex(crdb.NewIndex("is_default", []string{LoginNamePoliciesResourceOwnerCol, LoginNamePoliciesIsDefaultCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{LoginNamePoliciesOwnerRemovedCol})),
		),
	)
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
					Reduce: p.reduceOrgIAMPolicyAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: p.reduceDomainPolicyChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: p.reduceInstanceRemoved,
				},
			},
		},
	}
}

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

func (p *loginNameProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-QIe3C", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceOrgIAMPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *loginNameProjection) reduceDomainPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
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
			handler.NewCond(LoginNamePoliciesInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ysEeB", "reduce.wrong.event.type %s", org.DomainPolicyRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNamePoliciesResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(LoginNamePoliciesInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceDomainVerified(event eventstore.Event) (*handler.Statement, error) {
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

func (p *loginNameProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
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
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
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
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(loginNameDomainSuffix),
		),
	), nil
}

func (p *loginNameProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-4RHYq", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
	}

	return crdb.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameDomainNameCol, e.Domain),
			handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(loginNameDomainSuffix),
	), nil
}

func (p *loginNameProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASeg3", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return crdb.NewMultiStatement(
		event,
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNameDomainSuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNamePoliciesInstanceIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNamePolicySuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNameUserSuffix),
		),
	), nil
}

func (p *loginNameProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-px02mo", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewMultiStatement(
		event,
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNameDomainSuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNamePoliciesInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(LoginNamePoliciesResourceOwnerCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNamePolicySuffix),
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(LoginNameUserResourceOwnerCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(loginNameUserSuffix),
		),
	), nil
}
