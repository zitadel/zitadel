package projection

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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

type loginNameProjection struct{}

func newLoginNameProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(loginNameProjection))
}

func (*loginNameProjection) Name() string {
	return LoginNameProjectionTable
}

func (*loginNameProjection) Init() *old_handler.Check {
	return handler.NewViewCheck(
		viewStmt,
		handler.NewSuffixedTable(
			[]*handler.InitColumn{
				handler.NewColumn(LoginNameUserIDCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameUserUserNameCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameUserResourceOwnerCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameUserInstanceIDCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameUserOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
			},
			handler.NewPrimaryKey(LoginNameUserInstanceIDCol, LoginNameUserIDCol),
			loginNameUserSuffix,
			handler.WithIndex(handler.NewIndex("resource_owner", []string{LoginNameUserResourceOwnerCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{LoginNameUserOwnerRemovedCol})),
			handler.WithIndex(
				handler.NewIndex("lnu_instance_ro_id", []string{LoginNameUserInstanceIDCol, LoginNameUserResourceOwnerCol, LoginNameUserIDCol},
					handler.WithInclude(
						LoginNameUserUserNameCol,
						LoginNameUserOwnerRemovedCol,
					),
				),
			),
		),
		handler.NewSuffixedTable(
			[]*handler.InitColumn{
				handler.NewColumn(LoginNameDomainNameCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameDomainIsPrimaryCol, handler.ColumnTypeBool, handler.Default(false)),
				handler.NewColumn(LoginNameDomainResourceOwnerCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameDomainInstanceIDCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNameDomainOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
			},
			handler.NewPrimaryKey(LoginNameDomainInstanceIDCol, LoginNameDomainResourceOwnerCol, LoginNameDomainNameCol),
			loginNameDomainSuffix,
			handler.WithIndex(handler.NewIndex("owner_removed", []string{LoginNameDomainOwnerRemovedCol})),
		),
		handler.NewSuffixedTable(
			[]*handler.InitColumn{
				handler.NewColumn(LoginNamePoliciesMustBeDomainCol, handler.ColumnTypeBool),
				handler.NewColumn(LoginNamePoliciesIsDefaultCol, handler.ColumnTypeBool),
				handler.NewColumn(LoginNamePoliciesResourceOwnerCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNamePoliciesInstanceIDCol, handler.ColumnTypeText),
				handler.NewColumn(LoginNamePoliciesOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
			},
			handler.NewPrimaryKey(LoginNamePoliciesInstanceIDCol, LoginNamePoliciesResourceOwnerCol),
			loginNamePolicySuffix,
			handler.WithIndex(handler.NewIndex("is_default", []string{LoginNamePoliciesResourceOwnerCol, LoginNamePoliciesIsDefaultCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{LoginNamePoliciesOwnerRemovedCol})),
		),
	)
}

func (p *loginNameProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
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

	return handler.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserIDCol, event.Aggregate().ID),
			handler.NewCol(LoginNameUserUserNameCol, userName),
			handler.NewCol(LoginNameUserResourceOwnerCol, event.Aggregate().ResourceOwner),
			handler.NewCol(LoginNameUserInstanceIDCol, event.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-QIe3C", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-QlwjC", "reduce.wrong.event.type %s", user.UserUserNameChangedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
		},
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNameUserSuffix),
	), nil
}

func (p *loginNameProjection) reduceUserDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-AQMBY", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameUserUserNameCol, e.UserName),
		},
		[]handler.Condition{
			handler.NewCond(LoginNameUserIDCol, e.Aggregate().ID),
			handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNameUserSuffix),
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

	return handler.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNamePoliciesMustBeDomainCol, policyEvent.UserLoginMustBeDomain),
			handler.NewCol(LoginNamePoliciesIsDefaultCol, isDefault),
			handler.NewCol(LoginNamePoliciesResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCol(LoginNamePoliciesInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNamePolicySuffix),
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
		return handler.NewNoOpStatement(event), nil
	}

	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNamePoliciesMustBeDomainCol, *policyEvent.UserLoginMustBeDomain),
		},
		[]handler.Condition{
			handler.NewCond(LoginNamePoliciesResourceOwnerCol, policyEvent.Aggregate().ResourceOwner),
			handler.NewCond(LoginNamePoliciesInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceDomainPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ysEeB", "reduce.wrong.event.type %s", org.DomainPolicyRemovedEventType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNamePoliciesResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(LoginNamePoliciesInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNamePolicySuffix),
	), nil
}

func (p *loginNameProjection) reduceDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-weGAh", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
	}

	return handler.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(LoginNameDomainNameCol, e.Domain),
			handler.NewCol(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNameDomainSuffix),
	), nil
}

func (p *loginNameProjection) reducePrimaryDomainSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eOXPN", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(LoginNameDomainIsPrimaryCol, false),
			},
			[]handler.Condition{
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCond(LoginNameDomainIsPrimaryCol, true),
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(loginNameDomainSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(LoginNameDomainIsPrimaryCol, true),
			},
			[]handler.Condition{
				handler.NewCond(LoginNameDomainNameCol, e.Domain),
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(loginNameDomainSuffix),
		),
	), nil
}

func (p *loginNameProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-4RHYq", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(LoginNameDomainNameCol, e.Domain),
			handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(loginNameDomainSuffix),
	), nil
}

func (p *loginNameProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASeg3", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}

	return handler.NewMultiStatement(
		event,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(loginNameDomainSuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNamePoliciesInstanceIDCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(loginNamePolicySuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(loginNameUserSuffix),
		),
	), nil
}

func (p *loginNameProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-px02mo", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewMultiStatement(
		event,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameDomainInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(LoginNameDomainResourceOwnerCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(loginNameDomainSuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNamePoliciesInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(LoginNamePoliciesResourceOwnerCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(loginNamePolicySuffix),
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(LoginNameUserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCond(LoginNameUserResourceOwnerCol, e.Aggregate().ID),
			},
			handler.WithTableSuffix(loginNameUserSuffix),
		),
	), nil
}
