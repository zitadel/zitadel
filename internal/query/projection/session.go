package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SessionsProjectionTable = "projections.sessions8"

	SessionColumnID                     = "id"
	SessionColumnCreationDate           = "creation_date"
	SessionColumnChangeDate             = "change_date"
	SessionColumnSequence               = "sequence"
	SessionColumnState                  = "state"
	SessionColumnResourceOwner          = "resource_owner"
	SessionColumnInstanceID             = "instance_id"
	SessionColumnCreator                = "creator"
	SessionColumnUserID                 = "user_id"
	SessionColumnUserResourceOwner      = "user_resource_owner"
	SessionColumnUserCheckedAt          = "user_checked_at"
	SessionColumnPasswordCheckedAt      = "password_checked_at"
	SessionColumnIntentCheckedAt        = "intent_checked_at"
	SessionColumnWebAuthNCheckedAt      = "webauthn_checked_at"
	SessionColumnWebAuthNUserVerified   = "webauthn_user_verified"
	SessionColumnTOTPCheckedAt          = "totp_checked_at"
	SessionColumnOTPSMSCheckedAt        = "otp_sms_checked_at"
	SessionColumnOTPEmailCheckedAt      = "otp_email_checked_at"
	SessionColumnMetadata               = "metadata"
	SessionColumnTokenID                = "token_id"
	SessionColumnUserAgentFingerprintID = "user_agent_fingerprint_id"
	SessionColumnUserAgentIP            = "user_agent_ip"
	SessionColumnUserAgentDescription   = "user_agent_description"
	SessionColumnUserAgentHeader        = "user_agent_header"
	SessionColumnExpiration             = "expiration"
)

type sessionProjection struct{}

func newSessionProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(sessionProjection))
}

func (*sessionProjection) Name() string {
	return SessionsProjectionTable
}

func (*sessionProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SessionColumnID, handler.ColumnTypeText),
			handler.NewColumn(SessionColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SessionColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SessionColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SessionColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(SessionColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(SessionColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SessionColumnCreator, handler.ColumnTypeText),
			handler.NewColumn(SessionColumnUserID, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SessionColumnUserResourceOwner, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SessionColumnUserCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnPasswordCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnIntentCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnWebAuthNCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnWebAuthNUserVerified, handler.ColumnTypeBool, handler.Nullable()),
			handler.NewColumn(SessionColumnTOTPCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnOTPSMSCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnOTPEmailCheckedAt, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(SessionColumnMetadata, handler.ColumnTypeJSONB, handler.Nullable()),
			handler.NewColumn(SessionColumnTokenID, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SessionColumnUserAgentFingerprintID, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SessionColumnUserAgentIP, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SessionColumnUserAgentDescription, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SessionColumnUserAgentHeader, handler.ColumnTypeJSONB, handler.Nullable()),
			handler.NewColumn(SessionColumnExpiration, handler.ColumnTypeTimestamp, handler.Nullable()),
		},
			handler.NewPrimaryKey(SessionColumnInstanceID, SessionColumnID),
			handler.WithIndex(handler.NewIndex(
				SessionColumnUserAgentFingerprintID+"_idx",
				[]string{SessionColumnUserAgentFingerprintID},
			)),
			handler.WithIndex(handler.NewIndex(SessionColumnUserID+"_idx", []string{SessionColumnUserID})),
		),
	)
}

func (p *sessionProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: session.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  session.AddedType,
					Reduce: p.reduceSessionAdded,
				},
				{
					Event:  session.UserCheckedType,
					Reduce: p.reduceUserChecked,
				},
				{
					Event:  session.PasswordCheckedType,
					Reduce: p.reducePasswordChecked,
				},
				{
					Event:  session.IntentCheckedType,
					Reduce: p.reduceIntentChecked,
				},
				{
					Event:  session.WebAuthNCheckedType,
					Reduce: p.reduceWebAuthNChecked,
				},
				{
					Event:  session.TOTPCheckedType,
					Reduce: p.reduceTOTPChecked,
				},
				{
					Event:  session.OTPSMSCheckedType,
					Reduce: p.reduceOTPSMSChecked,
				},
				{
					Event:  session.OTPEmailCheckedType,
					Reduce: p.reduceOTPEmailChecked,
				},
				{
					Event:  session.TokenSetType,
					Reduce: p.reduceTokenSet,
				},
				{
					Event:  session.MetadataSetType,
					Reduce: p.reduceMetadataSet,
				},
				{
					Event:  session.LifetimeSetType,
					Reduce: p.reduceLifetimeSet,
				},
				{
					Event:  session.TerminateType,
					Reduce: p.reduceSessionTerminated,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SMSColumnInstanceID),
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.HumanPasswordChangedType,
					Reduce: p.reducePasswordChanged,
				},
			},
		},
	}
}

func (p *sessionProjection) reduceSessionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Sfrgf", "reduce.wrong.event.type %s", session.AddedType)
	}

	cols := make([]handler.Column, 0, 12)
	cols = append(cols,
		handler.NewCol(SessionColumnID, e.Aggregate().ID),
		handler.NewCol(SessionColumnInstanceID, e.Aggregate().InstanceID),
		handler.NewCol(SessionColumnCreationDate, e.CreationDate()),
		handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
		handler.NewCol(SessionColumnResourceOwner, e.Aggregate().ResourceOwner),
		handler.NewCol(SessionColumnState, domain.SessionStateActive),
		handler.NewCol(SessionColumnSequence, e.Sequence()),
		handler.NewCol(SessionColumnCreator, e.User),
	)
	if e.UserAgent != nil {
		cols = append(cols,
			handler.NewCol(SessionColumnUserAgentFingerprintID, e.UserAgent.FingerprintID),
			handler.NewCol(SessionColumnUserAgentDescription, e.UserAgent.Description),
		)
		if e.UserAgent.IP != nil {
			cols = append(cols,
				handler.NewCol(SessionColumnUserAgentIP, e.UserAgent.IP.String()),
			)
		}
		if e.UserAgent.Header != nil {
			cols = append(cols,
				handler.NewJSONCol(SessionColumnUserAgentHeader, e.UserAgent.Header),
			)
		}
	}

	return handler.NewCreateStatement(e, cols), nil
}

func (p *sessionProjection) reduceUserChecked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.UserCheckedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-saDg5", "reduce.wrong.event.type %s", session.UserCheckedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnUserID, e.UserID),
			handler.NewCol(SessionColumnUserResourceOwner, e.UserResourceOwner),
			handler.NewCol(SessionColumnUserCheckedAt, e.CheckedAt),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reducePasswordChecked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.PasswordCheckedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SDgrb", "reduce.wrong.event.type %s", session.PasswordCheckedType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnPasswordCheckedAt, e.CheckedAt),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceIntentChecked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.IntentCheckedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SDgr2", "reduce.wrong.event.type %s", session.IntentCheckedType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnIntentCheckedAt, e.CheckedAt),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceWebAuthNChecked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.WebAuthNCheckedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-WieM4", "reduce.wrong.event.type %s", session.WebAuthNCheckedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnWebAuthNCheckedAt, e.CheckedAt),
			handler.NewCol(SessionColumnWebAuthNUserVerified, e.UserVerified),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceTOTPChecked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TOTPCheckedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oqu8i", "reduce.wrong.event.type %s", session.TOTPCheckedType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnTOTPCheckedAt, e.CheckedAt),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceOTPSMSChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPSMSCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnOTPSMSCheckedAt, e.CheckedAt),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceOTPEmailChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPEmailCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnOTPEmailCheckedAt, e.CheckedAt),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceTokenSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TokenSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAfd3", "reduce.wrong.event.type %s", session.TokenSetType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnTokenID, e.TokenID),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.MetadataSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAfd3", "reduce.wrong.event.type %s", session.MetadataSetType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnMetadata, e.Metadata),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceLifetimeSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.LifetimeSetEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnExpiration, e.CreationDate().Add(e.Lifetime)),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TerminateEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAftn", "reduce.wrong.event.type %s", session.TerminateType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(SessionColumnID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *sessionProjection) reducePasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Deg3d", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnPasswordCheckedAt, nil),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnUserID, e.Aggregate().ID),
			handler.NewCond(SessionColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewLessThanCond(SessionColumnPasswordCheckedAt, e.CreationDate()),
		},
	), nil
}
