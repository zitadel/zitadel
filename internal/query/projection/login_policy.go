package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

const (
	LoginPolicyTable = "projections.login_policies"

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
	LoginPolicy2FAsCol                  = "second_factors"
	LoginPolicyMFAsCol                  = "multi_factors"
	LoginPolicyPasswordlessTypeCol      = "passwordless_type"
	LoginPolicyHidePWResetCol           = "hide_password_reset"
	PasswordCheckLifetimeCol            = "password_check_lifetime"
	ExternalLoginCheckLifetimeCol       = "external_login_check_lifetime"
	MFAInitSkipLifetimeCol              = "mfa_init_skip_lifetime"
	SecondFactorCheckLifetimeCol        = "second_factor_check_lifetime"
	MultiFactorCheckLifetimeCol         = "multi_factor_check_lifetime"
)

type LoginPolicyProjection struct {
	crdb.StatementHandler
}

func NewLoginPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *LoginPolicyProjection {
	p := new(LoginPolicyProjection)
	config.ProjectionName = LoginPolicyTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(LoginPolicyIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginPolicyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LoginPolicyCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(LoginPolicyChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(LoginPolicySequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(LoginPolicyIsDefaultCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(LoginPolicyAllowRegisterCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginPolicyAllowUsernamePasswordCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginPolicyAllowExternalIDPsCol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginPolicyForceMFACol, crdb.ColumnTypeBool),
			crdb.NewColumn(LoginPolicy2FAsCol, crdb.ColumnTypeEnumArray, crdb.Nullable()),
			crdb.NewColumn(LoginPolicyMFAsCol, crdb.ColumnTypeEnumArray, crdb.Nullable()),
			crdb.NewColumn(LoginPolicyPasswordlessTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(LoginPolicyHidePWResetCol, crdb.ColumnTypeBool),
			crdb.NewColumn(PasswordCheckLifetimeCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(ExternalLoginCheckLifetimeCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(MFAInitSkipLifetimeCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(SecondFactorCheckLifetimeCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(MultiFactorCheckLifetimeCol, crdb.ColumnTypeInt64),
		},
			crdb.NewPrimaryKey(LoginPolicyInstanceIDCol, LoginPolicyIDCol),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *LoginPolicyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
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
					Reduce: p.reduce2FAAdded,
				},
				{
					Event:  org.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduce2FARemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
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
					Reduce: p.reduce2FAAdded,
				},
				{
					Event:  instance.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduce2FARemoved,
				},
			},
		},
	}
}

func (p *LoginPolicyProjection) reduceLoginPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return crdb.NewCreateStatement(&policyEvent, []handler.Column{
		handler.NewCol(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		handler.NewCol(LoginPolicyInstanceIDCol, policyEvent.Aggregate().InstanceID),
		handler.NewCol(LoginPolicyCreationDateCol, policyEvent.CreationDate()),
		handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
		handler.NewCol(LoginPolicyAllowRegisterCol, policyEvent.AllowRegister),
		handler.NewCol(LoginPolicyAllowUsernamePasswordCol, policyEvent.AllowUserNamePassword),
		handler.NewCol(LoginPolicyAllowExternalIDPsCol, policyEvent.AllowExternalIDP),
		handler.NewCol(LoginPolicyForceMFACol, policyEvent.ForceMFA),
		handler.NewCol(LoginPolicyPasswordlessTypeCol, policyEvent.PasswordlessType),
		handler.NewCol(LoginPolicyIsDefaultCol, isDefault),
		handler.NewCol(LoginPolicyHidePWResetCol, policyEvent.HidePasswordReset),
		handler.NewCol(PasswordCheckLifetimeCol, policyEvent.PasswordCheckLifetime),
		handler.NewCol(ExternalLoginCheckLifetimeCol, policyEvent.ExternalLoginCheckLifetime),
		handler.NewCol(MFAInitSkipLifetimeCol, policyEvent.MFAInitSkipLifetime),
		handler.NewCol(SecondFactorCheckLifetimeCol, policyEvent.SecondFactorCheckLifetime),
		handler.NewCol(MultiFactorCheckLifetimeCol, policyEvent.MultiFactorCheckLifetime),
	}), nil
}

func (p *LoginPolicyProjection) reduceLoginPolicyChanged(event eventstore.Event) (*handler.Statement, error) {
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
	if policyEvent.PasswordlessType != nil {
		cols = append(cols, handler.NewCol(LoginPolicyPasswordlessTypeCol, *policyEvent.PasswordlessType))
	}
	if policyEvent.HidePasswordReset != nil {
		cols = append(cols, handler.NewCol(LoginPolicyHidePWResetCol, *policyEvent.HidePasswordReset))
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

	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduceMFAAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.MultiFactorAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	case *org.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-WMhAV", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorAddedEventType, instance.LoginPolicyMultiFactorAddedEventType})
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayAppendCol(LoginPolicyMFAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduceMFARemoved(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.MultiFactorRemovedEvent
	switch e := event.(type) {
	case *instance.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	case *org.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-czU7n", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyMultiFactorRemovedEventType, instance.LoginPolicyMultiFactorRemovedEventType})
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayRemoveCol(LoginPolicyMFAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduceLoginPolicyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.LoginPolicyRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-oRSvD", "reduce.wrong.event.type %s", org.LoginPolicyRemovedEventType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduce2FAAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.SecondFactorAddedEvent
	switch e := event.(type) {
	case *instance.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	case *org.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-agB2E", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorAddedEventType, instance.LoginPolicySecondFactorAddedEventType})
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayAppendCol(LoginPolicy2FAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduce2FARemoved(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.SecondFactorRemovedEvent
	switch e := event.(type) {
	case *instance.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	case *org.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-KYJvA", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicySecondFactorRemovedEventType, instance.LoginPolicySecondFactorRemovedEventType})
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(LoginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(LoginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayRemoveCol(LoginPolicy2FAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(LoginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}
