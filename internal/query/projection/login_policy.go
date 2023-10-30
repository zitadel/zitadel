package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	LoginPolicyTable = "projections.login_policies5"

	LoginPolicyIDCol                    = "aggregate_id"
	LoginPolicyInstanceIDCol            = "instance_id"
	LoginPolicyCreationDateCol          = "creation_date"
	LoginPolicyChangeDateCol            = "change_date"
	LoginPolicySequenceCol              = "sequence"
	LoginPolicyIsDefaultCol             = "is_default"
	LoginPolicyAllowRegisterCol         = "allow_register"
	LoginPolicyAllowUsernamePasswordCol = "allow_username_password"
	LoginPolicyAllowExternalIDPsCol     = "allow_external_idps"
	LoginPolicyForceMFACol              = "force_mfa"
	LoginPolicyForceMFALocalOnlyCol     = "force_mfa_local_only"
	LoginPolicy2FAsCol                  = "second_factors"
	LoginPolicyMFAsCol                  = "multi_factors"
	LoginPolicyPasswordlessTypeCol      = "passwordless_type"
	LoginPolicyHidePWResetCol           = "hide_password_reset"
	IgnoreUnknownUsernames              = "ignore_unknown_usernames"
	AllowDomainDiscovery                = "allow_domain_discovery"
	DisableLoginWithEmail               = "disable_login_with_email"
	DisableLoginWithPhone               = "disable_login_with_phone"
	DefaultRedirectURI                  = "default_redirect_uri"
	PasswordCheckLifetimeCol            = "password_check_lifetime"
	ExternalLoginCheckLifetimeCol       = "external_login_check_lifetime"
	MFAInitSkipLifetimeCol              = "mfa_init_skip_lifetime"
	SecondFactorCheckLifetimeCol        = "second_factor_check_lifetime"
	MultiFactorCheckLifetimeCol         = "multi_factor_check_lifetime"
	LoginPolicyOwnerRemovedCol          = "owner_removed"
)

type loginPolicyProjection struct{}

func newLoginPolicyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(loginPolicyProjection))
}

func (*loginPolicyProjection) Name() string {
	return LoginPolicyTable
}

func (*loginPolicyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(LoginPolicyIDCol, handler.ColumnTypeText),
			handler.NewColumn(LoginPolicyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(LoginPolicyCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(LoginPolicyChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(LoginPolicySequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(LoginPolicyIsDefaultCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LoginPolicyAllowRegisterCol, handler.ColumnTypeBool),
			handler.NewColumn(LoginPolicyAllowUsernamePasswordCol, handler.ColumnTypeBool),
			handler.NewColumn(LoginPolicyAllowExternalIDPsCol, handler.ColumnTypeBool),
			handler.NewColumn(LoginPolicyForceMFACol, handler.ColumnTypeBool),
			handler.NewColumn(LoginPolicyForceMFALocalOnlyCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(LoginPolicy2FAsCol, handler.ColumnTypeEnumArray, handler.Nullable()),
			handler.NewColumn(LoginPolicyMFAsCol, handler.ColumnTypeEnumArray, handler.Nullable()),
			handler.NewColumn(LoginPolicyPasswordlessTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(LoginPolicyHidePWResetCol, handler.ColumnTypeBool),
			handler.NewColumn(IgnoreUnknownUsernames, handler.ColumnTypeBool),
			handler.NewColumn(AllowDomainDiscovery, handler.ColumnTypeBool),
			handler.NewColumn(DisableLoginWithEmail, handler.ColumnTypeBool),
			handler.NewColumn(DisableLoginWithPhone, handler.ColumnTypeBool),
			handler.NewColumn(DefaultRedirectURI, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(PasswordCheckLifetimeCol, handler.ColumnTypeInt64),
			handler.NewColumn(ExternalLoginCheckLifetimeCol, handler.ColumnTypeInt64),
			handler.NewColumn(MFAInitSkipLifetimeCol, handler.ColumnTypeInt64),
			handler.NewColumn(SecondFactorCheckLifetimeCol, handler.ColumnTypeInt64),
			handler.NewColumn(MultiFactorCheckLifetimeCol, handler.ColumnTypeInt64),
			handler.NewColumn(LoginPolicyOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(LoginPolicyInstanceIDCol, LoginPolicyIDCol),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{LoginPolicyOwnerRemovedCol})),
		),
	)
}

func (p *loginPolicyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				{
					Event:  org.LoginPolicyChangedEventType,
					Reduce: p.reduceLoginPolicyChanged,
				},
				{
					Event:  org.LoginPolicyMultiFactorAddedEventType,
					Reduce: p.reduceMFAAdded,
				},
				{
					Event:  org.LoginPolicyMultiFactorRemovedEventType,
					Reduce: p.reduceMFARemoved,
				},
				{
					Event:  org.LoginPolicyRemovedEventType,
					Reduce: p.reduceLoginPolicyRemoved,
				},
				{
					Event:  org.LoginPolicySecondFactorAddedEventType,
					Reduce: p.reduceSecondFactorAdded,
				},
				{
					Event:  org.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduceSecondFactorRemoved,
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
					Event:  instance.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				{
					Event:  instance.LoginPolicyChangedEventType,
					Reduce: p.reduceLoginPolicyChanged,
				},
				{
					Event:  instance.LoginPolicyMultiFactorAddedEventType,
					Reduce: p.reduceMFAAdded,
				},
				{
					Event:  instance.LoginPolicyMultiFactorRemovedEventType,
					Reduce: p.reduceMFARemoved,
				},
				{
					Event:  instance.LoginPolicySecondFactorAddedEventType,
					Reduce: p.reduceSecondFactorAdded,
				},
				{
					Event:  instance.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduceSecondFactorRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(LoginPolicyInstanceIDCol),
				},
			},
		},
	}
}

func (p *loginPolicyProjection) reduceLoginPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LoginPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *instance.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = true
	case *org.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = false
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-pYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyAddedEventType, instance.LoginPolicyAddedEventType})
	}

	return handler.NewCreateStatement(&policyEvent, []handler.Column{
		handler.NewCol(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		handler.NewCol(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		handler.NewCol(LoginPolicyCreationDateCol, policyEvent.CreationDate()),
		handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
		handler.NewCol(LoginPolicyAllowRegisterCol, policyEvent.AllowRegister),
		handler.NewCol(LoginPolicyAllowUsernamePasswordCol, policyEvent.AllowUserNamePassword),
		handler.NewCol(LoginPolicyAllowExternalIDPsCol, policyEvent.AllowExternalIDP),
		handler.NewCol(LoginPolicyForceMFACol, policyEvent.ForceMFA),
		handler.NewCol(LoginPolicyForceMFALocalOnlyCol, policyEvent.ForceMFALocalOnly),
		handler.NewCol(LoginPolicyPasswordlessTypeCol, policyEvent.PasswordlessType),
		handler.NewCol(LoginPolicyIsDefaultCol, isDefault),
		handler.NewCol(LoginPolicyHidePWResetCol, policyEvent.HidePasswordReset),
		handler.NewCol(IgnoreUnknownUsernames, policyEvent.IgnoreUnknownUsernames),
		handler.NewCol(AllowDomainDiscovery, policyEvent.AllowDomainDiscovery),
		handler.NewCol(DisableLoginWithEmail, policyEvent.DisableLoginWithEmail),
		handler.NewCol(DisableLoginWithPhone, policyEvent.DisableLoginWithPhone),
		handler.NewCol(DefaultRedirectURI, policyEvent.DefaultRedirectURI),
		handler.NewCol(PasswordCheckLifetimeCol, policyEvent.PasswordCheckLifetime),
		handler.NewCol(ExternalLoginCheckLifetimeCol, policyEvent.ExternalLoginCheckLifetime),
		handler.NewCol(MFAInitSkipLifetimeCol, policyEvent.MFAInitSkipLifetime),
		handler.NewCol(SecondFactorCheckLifetimeCol, policyEvent.SecondFactorCheckLifetime),
		handler.NewCol(MultiFactorCheckLifetimeCol, policyEvent.MultiFactorCheckLifetime),
	}), nil
}

func (p *loginPolicyProjection) reduceLoginPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LoginPolicyChangedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
	case *org.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-BpaO6", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyChangedEventType, instance.LoginPolicyChangedEventType})
	}

	cols := []handler.Column{
		handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.AllowRegister != nil {
		cols = append(cols, handler.NewCol(LoginPolicyAllowRegisterCol, *policyEvent.AllowRegister))
	}
	if policyEvent.AllowUserNamePassword != nil {
		cols = append(cols, handler.NewCol(LoginPolicyAllowUsernamePasswordCol, *policyEvent.AllowUserNamePassword))
	}
	if policyEvent.AllowExternalIDP != nil {
		cols = append(cols, handler.NewCol(LoginPolicyAllowExternalIDPsCol, *policyEvent.AllowExternalIDP))
	}
	if policyEvent.ForceMFA != nil {
		cols = append(cols, handler.NewCol(LoginPolicyForceMFACol, *policyEvent.ForceMFA))
	}
	if policyEvent.ForceMFALocalOnly != nil {
		cols = append(cols, handler.NewCol(LoginPolicyForceMFALocalOnlyCol, *policyEvent.ForceMFALocalOnly))
	}
	if policyEvent.PasswordlessType != nil {
		cols = append(cols, handler.NewCol(LoginPolicyPasswordlessTypeCol, *policyEvent.PasswordlessType))
	}
	if policyEvent.HidePasswordReset != nil {
		cols = append(cols, handler.NewCol(LoginPolicyHidePWResetCol, *policyEvent.HidePasswordReset))
	}
	if policyEvent.IgnoreUnknownUsernames != nil {
		cols = append(cols, handler.NewCol(IgnoreUnknownUsernames, *policyEvent.IgnoreUnknownUsernames))
	}
	if policyEvent.AllowDomainDiscovery != nil {
		cols = append(cols, handler.NewCol(AllowDomainDiscovery, *policyEvent.AllowDomainDiscovery))
	}
	if policyEvent.DisableLoginWithEmail != nil {
		cols = append(cols, handler.NewCol(DisableLoginWithEmail, *policyEvent.DisableLoginWithEmail))
	}
	if policyEvent.DisableLoginWithPhone != nil {
		cols = append(cols, handler.NewCol(DisableLoginWithPhone, *policyEvent.DisableLoginWithPhone))
	}
	if policyEvent.DefaultRedirectURI != nil {
		cols = append(cols, handler.NewCol(DefaultRedirectURI, *policyEvent.DefaultRedirectURI))
	}
	if policyEvent.PasswordCheckLifetime != nil {
		cols = append(cols, handler.NewCol(PasswordCheckLifetimeCol, *policyEvent.PasswordCheckLifetime))
	}
	if policyEvent.ExternalLoginCheckLifetime != nil {
		cols = append(cols, handler.NewCol(ExternalLoginCheckLifetimeCol, *policyEvent.ExternalLoginCheckLifetime))
	}
	if policyEvent.MFAInitSkipLifetime != nil {
		cols = append(cols, handler.NewCol(MFAInitSkipLifetimeCol, *policyEvent.MFAInitSkipLifetime))
	}
	if policyEvent.SecondFactorCheckLifetime != nil {
		cols = append(cols, handler.NewCol(SecondFactorCheckLifetimeCol, *policyEvent.SecondFactorCheckLifetime))
	}
	if policyEvent.MultiFactorCheckLifetime != nil {
		cols = append(cols, handler.NewCol(MultiFactorCheckLifetimeCol, *policyEvent.MultiFactorCheckLifetime))
	}

	return handler.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *loginPolicyProjection) reduceMFAAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.MultiFactorAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	case *org.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-WMhAV", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorAddedEventType, instance.LoginPolicyMultiFactorAddedEventType})
	}

	return handler.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			handler.NewArrayAppendCol(LoginPolicyMFAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *loginPolicyProjection) reduceMFARemoved(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.MultiFactorRemovedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	case *org.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-czU7n", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorRemovedEventType, instance.LoginPolicyMultiFactorRemovedEventType})
	}

	return handler.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			handler.NewArrayRemoveCol(LoginPolicyMFAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *loginPolicyProjection) reduceLoginPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.LoginPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-oRSvD", "reduce.wrong.event.type %s", org.LoginPolicyRemovedEventType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, e.Aggregate().ID),
			handler.NewCond(LoginPolicyInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *loginPolicyProjection) reduceSecondFactorAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.SecondFactorAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	case *org.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-agB2E", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorAddedEventType, instance.LoginPolicySecondFactorAddedEventType})
	}

	return handler.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			handler.NewArrayAppendCol(LoginPolicy2FAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *loginPolicyProjection) reduceSecondFactorRemoved(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.SecondFactorRemovedEvent
	switch e := event.(type) {
	case *instance.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	case *org.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-KYJvA", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorRemovedEventType, instance.LoginPolicySecondFactorRemovedEventType})
	}

	return handler.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			handler.NewArrayRemoveCol(LoginPolicy2FAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
			handler.NewCond(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *loginPolicyProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-B8NZW", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LoginPolicyInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(LoginPolicyIDCol, e.Aggregate().ID),
		},
	), nil
}
