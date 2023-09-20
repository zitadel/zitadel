package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	SessionsProjectionTable = "projections.sessions5"

	SessionColumnID                   = "id"
	SessionColumnCreationDate         = "creation_date"
	SessionColumnChangeDate           = "change_date"
	SessionColumnSequence             = "sequence"
	SessionColumnState                = "state"
	SessionColumnResourceOwner        = "resource_owner"
	SessionColumnInstanceID           = "instance_id"
	SessionColumnCreator              = "creator"
	SessionColumnUserID               = "user_id"
	SessionColumnUserCheckedAt        = "user_checked_at"
	SessionColumnPasswordCheckedAt    = "password_checked_at"
	SessionColumnIntentCheckedAt      = "intent_checked_at"
	SessionColumnWebAuthNCheckedAt    = "webauthn_checked_at"
	SessionColumnWebAuthNUserVerified = "webauthn_user_verified"
	SessionColumnTOTPCheckedAt        = "totp_checked_at"
	SessionColumnOTPSMSCheckedAt      = "otp_sms_checked_at"
	SessionColumnOTPEmailCheckedAt    = "otp_email_checked_at"
	SessionColumnMetadata             = "metadata"
	SessionColumnTokenID              = "token_id"
)

type sessionProjection struct {
	crdb.StatementHandler
}

func newSessionProjection(ctx context.Context, config crdb.StatementHandlerConfig) *sessionProjection {
	p := new(sessionProjection)
	config.ProjectionName = SessionsProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(SessionColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SessionColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(SessionColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(SessionColumnState, crdb.ColumnTypeEnum),
			crdb.NewColumn(SessionColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnCreator, crdb.ColumnTypeText),
			crdb.NewColumn(SessionColumnUserID, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(SessionColumnUserCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnPasswordCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnIntentCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnWebAuthNCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnWebAuthNUserVerified, crdb.ColumnTypeBool, crdb.Nullable()),
			crdb.NewColumn(SessionColumnTOTPCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnOTPSMSCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnOTPEmailCheckedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(SessionColumnMetadata, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(SessionColumnTokenID, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(SessionColumnInstanceID, SessionColumnID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *sessionProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: session.AggregateType,
			EventRedusers: []handler.EventReducer{
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
					Event:  session.TerminateType,
					Reduce: p.reduceSessionTerminated,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SMSColumnInstanceID),
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sfrgf", "reduce.wrong.event.type %s", session.AddedType)
	}

	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnID, e.Aggregate().ID),
			handler.NewCol(SessionColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(SessionColumnCreationDate, e.CreationDate()),
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(SessionColumnState, domain.SessionStateActive),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnCreator, e.User),
		},
	), nil
}

func (p *sessionProjection) reduceUserChecked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.UserCheckedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-saDg5", "reduce.wrong.event.type %s", session.UserCheckedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnChangeDate, e.CreationDate()),
			handler.NewCol(SessionColumnSequence, e.Sequence()),
			handler.NewCol(SessionColumnUserID, e.UserID),
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SDgrb", "reduce.wrong.event.type %s", session.PasswordCheckedType)
	}

	return crdb.NewUpdateStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SDgr2", "reduce.wrong.event.type %s", session.IntentCheckedType)
	}

	return crdb.NewUpdateStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-WieM4", "reduce.wrong.event.type %s", session.WebAuthNCheckedType)
	}
	return crdb.NewUpdateStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Oqu8i", "reduce.wrong.event.type %s", session.TOTPCheckedType)
	}

	return crdb.NewUpdateStatement(
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

	return crdb.NewUpdateStatement(
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

	return crdb.NewUpdateStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAfd3", "reduce.wrong.event.type %s", session.TokenSetType)
	}

	return crdb.NewUpdateStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAfd3", "reduce.wrong.event.type %s", session.MetadataSetType)
	}

	return crdb.NewUpdateStatement(
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

func (p *sessionProjection) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TerminateEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAftn", "reduce.wrong.event.type %s", session.TerminateType)
	}

	return crdb.NewDeleteStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Deg3d", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SessionColumnPasswordCheckedAt, nil),
		},
		[]handler.Condition{
			handler.NewCond(SessionColumnUserID, e.Aggregate().ID),
			crdb.NewLessThanCond(SessionColumnPasswordCheckedAt, e.CreationDate()),
		},
	), nil
}
