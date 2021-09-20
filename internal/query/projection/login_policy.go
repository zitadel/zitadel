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
				// {
				// 	Event:  org.LoginPolicyIDPProviderAddedEventType,
				// 	Reduce: p.reduceIDPAddedEvent,
				// },
				// {
				// 	Event:  org.LoginPolicyIDPProviderCascadeRemovedEventType,
				// 	Reduce: p.reduceIDPCascadeRemoved,
				// },
				// {
				// 	Event:  org.LoginPolicyIDPProviderRemovedEventType,
				// 	Reduce: p.reduceIDPRemoved,
				// },
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
				// {
				// 	Event:  iam.LoginPolicyIDPProviderAddedEventType,
				// 	Reduce: p.reduceIDPAddedEvent,
				// },
				// {
				// 	Event:  iam.LoginPolicyIDPProviderCascadeRemovedEventType,
				// 	Reduce: p.reduceIDPCascadeRemoved,
				// },
				// {
				// 	Event:  iam.LoginPolicyIDPProviderRemovedEventType,
				// 	Reduce: p.reduceIDPRemoved,
				// },
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
	loginPolicyUserLoginMustBeDomainCol = "user_login_must_be_domain"
)

func (p *LoginPolicyProjection) reduceLoginPolicyAdded(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *iam.LoginPolicyAddedEvent:
		return crdb.NewCreateStatement(e, []handler.Column{
			handler.NewCol(loginPolicyIDCol, e.Aggregate().ID),
			handler.NewCol(loginPolicyCreationDateCol, e.CreationDate()),
			handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, e.Sequence()),
			handler.NewCol(loginPolicyAllowRegisterCol, e.AllowRegister),
			handler.NewCol(loginPolicyAllowUserNamePasswordCol, e.AllowUserNamePassword),
			handler.NewCol(loginPolicyAllowExternalIDPsCol, e.AllowExternalIDP),
			handler.NewCol(loginPolicyForceMFACol, e.ForceMFA),
			handler.NewCol(loginPolicyPasswordlessTypeCol, e.PasswordlessType),
			handler.NewCol(loginPolicyIsDefaultCol, true),
			handler.NewCol(loginPolicyHidePWResetCol, e.HidePasswordReset),
		}), nil
	case *org.OrgIAMPolicyAddedEvent:
		return crdb.NewCreateStatement(e, []handler.Column{
			handler.NewCol(loginPolicyIDCol, e.Aggregate().ID),
			handler.NewCol(loginPolicyCreationDateCol, e.CreationDate()),
			handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, e.Sequence()),
			handler.NewCol(loginPolicyIsDefaultCol, false),
			handler.NewCol(loginPolicyUserLoginMustBeDomainCol, e.UserLoginMustBeDomain),
		}), nil
	default:
		logging.LogWithFields("HANDL-IW6So", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.OrgIAMPolicyAddedEventType, iam.OrgIAMPolicyAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-pYPxS", "reduce.wrong.event.type")
	}
}

func (p *LoginPolicyProjection) reduceLoginPolicyChanged(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *iam.LoginPolicyChangedEvent:
		cols := []handler.Column{
			handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
			handler.NewCol(loginPolicySequenceCol, e.Sequence()),
		}
		if e.AllowRegister != nil {
			cols = append(cols, handler.NewCol(loginPolicyAllowRegisterCol, *e.AllowRegister))
		}
		if e.AllowUserNamePassword != nil {
			cols = append(cols, handler.NewCol(loginPolicyAllowUserNamePasswordCol, *e.AllowUserNamePassword))
		}
		if e.AllowExternalIDP != nil {
			cols = append(cols, handler.NewCol(loginPolicyAllowExternalIDPsCol, *e.AllowExternalIDP))
		}
		if e.ForceMFA != nil {
			cols = append(cols, handler.NewCol(loginPolicyForceMFACol, *e.ForceMFA))
		}
		if e.PasswordlessType != nil {
			cols = append(cols, handler.NewCol(loginPolicyPasswordlessTypeCol, *e.PasswordlessType))
		}
		if e.HidePasswordReset != nil {
			cols = append(cols, handler.NewCol(loginPolicyHidePWResetCol, *e.HidePasswordReset))
		}
		return crdb.NewUpdateStatement(
			e,
			cols,
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, true),
			},
		), nil
	case *org.OrgIAMPolicyChangedEvent:
		if e.UserLoginMustBeDomain == nil {
			return crdb.NewNoOpStatement(e), nil
		}
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				handler.NewCol(loginPolicyAllowRegisterCol, *e.UserLoginMustBeDomain),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, false),
			},
		), nil
	default:
		logging.LogWithFields("HANDL-NIvFo", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyChangedEventType, iam.LoginPolicyChangedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BpaO6", "reduce.wrong.event.type")
	}
}

// func (p *LoginPolicyProjection) reduceIDPAddedEvent(event eventstore.EventReader) (*handler.Statement, error) {
// 	switch e := event.(type) {
// 	case *iam.IdentityProviderAddedEvent:
// 		_ = e
// 		return nil, nil
// 	case *org.IdentityProviderAddedEvent:
// 		return nil, nil
// 	default:
// 		logging.LogWithFields("HANDL-CtSYI", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyIDPProviderAddedEventType, iam.LoginPolicyIDPProviderAddedEventType}).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "HANDL-DV2XO", "reduce.wrong.event.type")
// 	}
// }

// func (p *LoginPolicyProjection) reduceIDPRemoved(event eventstore.EventReader) (*handler.Statement, error) {
// 	switch e := event.(type) {
// 	case *iam.IdentityProviderRemovedEvent:
// 		_ = e
// 		return nil, nil
// 	case *org.IdentityProviderRemovedEvent:
// 		return nil, nil
// 	default:
// 		logging.LogWithFields("HANDL-I4RWX", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyIDPProviderRemovedEventType, iam.LoginPolicyIDPProviderRemovedEventType}).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "HANDL-PD2l4", "reduce.wrong.event.type")
// 	}
// }

// func (p *LoginPolicyProjection) reduceIDPCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
// 	switch e := event.(type) {
// 	case *iam.IdentityProviderCascadeRemovedEvent:
// 		_ = e
// 		return nil, nil
// 	case *org.IdentityProviderCascadeRemovedEvent:
// 		return nil, nil
// 	default:
// 		logging.LogWithFields("HANDL-IKmF9", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyIDPProviderCascadeRemovedEventType, iam.LoginPolicyIDPProviderCascadeRemovedEventType}).Error("wrong event type")
// 		return nil, errors.ThrowInvalidArgument(nil, "HANDL-GYuEp", "reduce.wrong.event.type")
// 	}
// }

func (p *LoginPolicyProjection) reduceMFAAdded(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *iam.LoginPolicyMultiFactorAddedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayAppendCol(loginPolicyMFAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, true),
			},
		), nil
	case *org.LoginPolicyMultiFactorAddedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayAppendCol(loginPolicyMFAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, false),
			},
		), nil
	default:
		logging.LogWithFields("HANDL-fYAHO", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyMultiFactorAddedEventType, iam.LoginPolicyMultiFactorAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-WMhAV", "reduce.wrong.event.type")
	}
}

func (p *LoginPolicyProjection) reduceMFARemoved(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *iam.LoginPolicyMultiFactorRemovedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayRemoveCol(loginPolicyMFAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, true),
			},
		), nil
	case *org.LoginPolicyMultiFactorRemovedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayRemoveCol(loginPolicyMFAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, false),
			},
		), nil
	default:
		logging.LogWithFields("HANDL-fYAHO", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicyMultiFactorRemovedEventType, iam.LoginPolicyMultiFactorRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-WMhAV", "reduce.wrong.event.type")
	}
}

func (p *LoginPolicyProjection) reduceLoginPolicyRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.LoginPolicyRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-fYAHO", "seq", event.Sequence(), "expectedType", org.LoginPolicyRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-WMhAV", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *LoginPolicyProjection) reduce2FAAdded(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *iam.LoginPolicySecondFactorAddedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayAppendCol(loginPolicy2FAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, true),
			},
		), nil
	case *org.LoginPolicySecondFactorAddedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayAppendCol(loginPolicy2FAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, false),
			},
		), nil
	default:
		logging.LogWithFields("HANDL-fYAHO", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicySecondFactorAddedEventType, iam.LoginPolicySecondFactorAddedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-WMhAV", "reduce.wrong.event.type")
	}
}

func (p *LoginPolicyProjection) reduce2FARemoved(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *iam.LoginPolicySecondFactorRemovedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayRemoveCol(loginPolicy2FAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, true),
			},
		), nil
	case *org.LoginPolicySecondFactorRemovedEvent:
		return crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(loginPolicyChangeDateCol, e.CreationDate()),
				handler.NewCol(loginPolicySequenceCol, e.Sequence()),
				crdb.NewArrayRemoveCol(loginPolicy2FAsCol, e.MFAType),
			},
			[]handler.Condition{
				handler.NewCond(loginPolicyIDCol, e.Aggregate().ID),
				// handler.NewCond(loginPolicyIsDefaultCol, true),
			},
		), nil
	default:
		logging.LogWithFields("HANDL-fYAHO", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.LoginPolicySecondFactorRemovedEventType, iam.LoginPolicySecondFactorRemovedEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-WMhAV", "reduce.wrong.event.type")
	}
}
