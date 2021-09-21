package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type LoginPolicyProjection struct {
	crdb.StatementHandler
}

const (
	loginPolicyProjection = "zitadel.projections.login_policies"
)

func NewLoginPolicyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *LoginPolicyProjection {
	p := &LoginPolicyProjection{}
	config.ProjectionName = loginPolicyProjection
	config.Reducers = p.reducers()
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
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				{
					Event:  iam.LoginPolicyChangedEventType,
					Reduce: p.reduceLoginPolicyChanged,
				},
				{
					Event:  iam.LoginPolicyMultiFactorAddedEventType,
					Reduce: p.reduceMFAAdded,
				},
				{
					Event:  iam.LoginPolicyMultiFactorRemovedEventType,
					Reduce: p.reduceMFARemoved,
				},
				{
					Event:  iam.LoginPolicySecondFactorAddedEventType,
					Reduce: p.reduce2FAAdded,
				},
				{
					Event:  iam.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduce2FARemoved,
				},
			},
		},
	}
}

const (
	loginPolicyIDCol                    = "aggregate_id"
	loginPolicyCreationDateCol          = "creation_date"
	loginPolicyChangeDateCol            = "change_date"
	loginPolicySequenceCol              = "sequence"
	loginPolicyAllowRegisterCol         = "allow_register"
	loginPolicyAllowUserNamePasswordCol = "allow_username_password"
	loginPolicyAllowExternalIDPsCol     = "allow_external_idps"
	loginPolicyForceMFACol              = "force_mfa"
	loginPolicy2FAsCol                  = "second_factors"
	loginPolicyMFAsCol                  = "multi_factors"
	loginPolicyPasswordlessTypeCol      = "passwordless_type"
	loginPolicyIsDefaultCol             = "is_default"
	loginPolicyHidePWResetCol           = "hide_password_reset"
)

func (p *LoginPolicyProjection) reduceLoginPolicyAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.LoginPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *iam.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = true
	case *org.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = false
	default:
		logging.LogWithFields("HANDL-IW6So", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyAddedEventType, iam.LoginPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pYPxS", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(&policyEvent, []handler.Column{
		handler.NewCol(loginPolicyIDCol, policyEvent.Aggregate().ID),
		handler.NewCol(loginPolicyCreationDateCol, policyEvent.CreationDate()),
		handler.NewCol(loginPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(loginPolicySequenceCol, policyEvent.Sequence()),
		handler.NewCol(loginPolicyAllowRegisterCol, policyEvent.AllowRegister),
		handler.NewCol(loginPolicyAllowUserNamePasswordCol, policyEvent.AllowUserNamePassword),
		handler.NewCol(loginPolicyAllowExternalIDPsCol, policyEvent.AllowExternalIDP),
		handler.NewCol(loginPolicyForceMFACol, policyEvent.ForceMFA),
		handler.NewCol(loginPolicyPasswordlessTypeCol, policyEvent.PasswordlessType),
		handler.NewCol(loginPolicyIsDefaultCol, isDefault),
		handler.NewCol(loginPolicyHidePWResetCol, policyEvent.HidePasswordReset),
	}), nil
}

func (p *LoginPolicyProjection) reduceLoginPolicyChanged(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.LoginPolicyChangedEvent
	switch e := event.(type) {
	case *iam.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
	case *org.LoginPolicyChangedEvent:
		policyEvent = e.LoginPolicyChangedEvent
	default:
		logging.LogWithFields("HANDL-NIvFo", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyChangedEventType, iam.LoginPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BpaO6", "reduce.wrong.event.type")
	}

	cols := []handler.Column{
		handler.NewCol(loginPolicyChangeDateCol, policyEvent.CreationDate()),
		handler.NewCol(loginPolicySequenceCol, policyEvent.Sequence()),
	}
	if policyEvent.AllowRegister != nil {
		cols = append(cols, handler.NewCol(loginPolicyAllowRegisterCol, *policyEvent.AllowRegister))
	}
	if policyEvent.AllowUserNamePassword != nil {
		cols = append(cols, handler.NewCol(loginPolicyAllowUserNamePasswordCol, *policyEvent.AllowUserNamePassword))
	}
	if policyEvent.AllowExternalIDP != nil {
		cols = append(cols, handler.NewCol(loginPolicyAllowExternalIDPsCol, *policyEvent.AllowExternalIDP))
	}
	if policyEvent.ForceMFA != nil {
		cols = append(cols, handler.NewCol(loginPolicyForceMFACol, *policyEvent.ForceMFA))
	}
	if policyEvent.PasswordlessType != nil {
		cols = append(cols, handler.NewCol(loginPolicyPasswordlessTypeCol, *policyEvent.PasswordlessType))
	}
	if policyEvent.HidePasswordReset != nil {
		cols = append(cols, handler.NewCol(loginPolicyHidePWResetCol, *policyEvent.HidePasswordReset))
	}
	return crdb.NewUpdateStatement(
		&policyEvent,
		cols,
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduceMFAAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.MultiFactorAddedEvent
	switch e := event.(type) {
	case *iam.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	case *org.LoginPolicyMultiFactorAddedEvent:
		policyEvent = e.MultiFactorAddedEvent
	default:
		logging.LogWithFields("HANDL-fYAHO", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyMultiFactorAddedEventType, iam.LoginPolicyMultiFactorAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-WMhAV", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(loginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayAppendCol(loginPolicyMFAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduceMFARemoved(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.MultiFactorRemovedEvent
	switch e := event.(type) {
	case *iam.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	case *org.LoginPolicyMultiFactorRemovedEvent:
		policyEvent = e.MultiFactorRemovedEvent
	default:
		logging.LogWithFields("HANDL-vtC31", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyMultiFactorRemovedEventType, iam.LoginPolicyMultiFactorRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-czU7n", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(loginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayRemoveCol(loginPolicyMFAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduceLoginPolicyRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.LoginPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-gF5q6", "seq", event.Sequence(), "expectedType", org.LoginPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-oRSvD", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduce2FAAdded(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.SecondFactorAddedEvent
	switch e := event.(type) {
	case *iam.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	case *org.LoginPolicySecondFactorAddedEvent:
		policyEvent = e.SecondFactorAddedEvent
	default:
		logging.LogWithFields("HANDL-dwadQ", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicySecondFactorAddedEventType, iam.LoginPolicySecondFactorAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-agB2E", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(loginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayAppendCol(loginPolicy2FAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduce2FARemoved(event eventstore.EventReader) (*handler.Statement, error) {
	var policyEvent policy.SecondFactorRemovedEvent
	switch e := event.(type) {
	case *iam.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	case *org.LoginPolicySecondFactorRemovedEvent:
		policyEvent = e.SecondFactorRemovedEvent
	default:
		logging.LogWithFields("HANDL-2IE8Y", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicySecondFactorRemovedEventType, iam.LoginPolicySecondFactorRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-KYJvA", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		&policyEvent,
		[]handler.Column{
			handler.NewCol(loginPolicyChangeDateCol, policyEvent.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, policyEvent.Sequence()),
			crdb.NewArrayRemoveCol(loginPolicy2FAsCol, policyEvent.MFAType),
		},
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, policyEvent.Aggregate().ID),
		},
	), nil
}
