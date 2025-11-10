package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SessionsRelationalProjectionTable = "zitadel.sessions"
)

type sessionRelationalProjection struct{}

func (*sessionRelationalProjection) Name() string {
	return SessionsRelationalProjectionTable
}

func newSessionRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(sessionRelationalProjection))
}

func (p *sessionRelationalProjection) Reducers() []handler.AggregateReducer {
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
					Event:  session.WebAuthNChallengedType,
					Reduce: p.reduceWebAuthNChallenged,
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
					Event:  session.OTPSMSChallengedType,
					Reduce: p.reduceOTPSMSChecked,
				},
				{
					Event:  session.OTPSMSCheckedType,
					Reduce: p.reduceOTPSMSChecked,
				},
				{
					Event:  session.OTPEmailChallengedType,
					Reduce: p.reduceOTPEmailChecked,
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

func (p *sessionRelationalProjection) reduceSessionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.AddedEvent](event)
	if err != nil {
		return nil, err
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

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		return sessionRepo.Create(ctx, v3Tx, &domain.Session{
			InstanceID: e.Aggregate().InstanceID,
			ID:         e.Aggregate().ID,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
			UserAgent: &domain.SessionUserAgent{
				FingerprintID: e.UserAgent.FingerprintID,
				Description:   e.UserAgent.Description,
				IP:            e.UserAgent.IP,
				Header:        e.UserAgent.Header,
			},
		})
	}), nil
}

func (p *sessionRelationalProjection) reduceUserChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.UserCheckedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorUser{
				UserID:         e.UserID,
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reducePasswordChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.PasswordCheckedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorPassword{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceIntentChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.IntentCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorIDPIntent{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceWebAuthNChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.WebAuthNChallengedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetChallenge(&domain.SessionChallengePasskey{
				LastChallengedAt:     e.CreatedAt(),
				Challenge:            e.Challenge,
				AllowedCredentialIDs: e.AllowedCrentialIDs,
				UserVerification:     e.UserVerification,
				RPID:                 e.RPID,
			}),
		)
		return err
	}), nil
}
func (p *sessionRelationalProjection) reduceWebAuthNChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.WebAuthNCheckedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorPasskey{
				LastVerifiedAt: e.CreatedAt(),
				UserVerified:   e.UserVerified,
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceTOTPChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.TOTPCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorTOTP{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceOTPSMSChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPSMSChallengedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetChallenge(&domain.SessionChallengeOTPSMS{
				LastChallengedAt:  e.CreatedAt(),
				Code:              e.Code,
				Expiry:            e.Expiry,
				CodeReturned:      e.CodeReturned,
				GeneratorID:       e.GeneratorID,
				TriggeredAtOrigin: e.TriggeredAtOrigin,
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceOTPSMSChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPSMSCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorOTPSMS{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceOTPEmailChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPEmailChallengedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetChallenge(&domain.SessionChallengeOTPEmail{
				LastChallengedAt:  e.CreatedAt(),
				Code:              e.Code,
				Expiry:            e.Expiry,
				CodeReturned:      e.ReturnCode,
				URLTmpl:           e.URLTmpl,
				TriggeredAtOrigin: e.TriggeredAtOrigin,
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceOTPEmailChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPEmailCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorOTPEmail{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceTokenSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.TokenSetEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetToken(e.TokenID),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.MetadataSetEvent](event)
	if err != nil {
		return nil, err
	}

	metadataList := make([]domain.SessionMetadata, 0, len(e.Metadata))
	for key, value := range e.Metadata {
		metadataList = append(metadataList, domain.SessionMetadata{
			Metadata: domain.Metadata{
				Key:   key,
				Value: value,
			},
		})
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetMetadata(metadataList),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceLifetimeSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.LifetimeSetEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetLifetime(e.Lifetime),
		)
		return err
	}), nil
}

func (p *sessionRelationalProjection) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.TerminateEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		return sessionRepo.Delete(ctx, v3Tx, condition)
	}), nil
}

func (p *sessionRelationalProjection) reducePasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordChangedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := database.And(
			sessionRepo.InstanceIDCondition(e.Aggregate().InstanceID),
			sessionRepo.UserIDCondition(e.Aggregate().ID),
			sessionRepo.FactorConditions().FactorTypeCondition(domain.SessionFactorTypePassword),
			sessionRepo.FactorConditions().FactorLastVerifiedBeforeCondition(e.CreatedAt()),
		)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.ClearFactor(domain.SessionFactorTypePassword),
		)
		return err
	}), nil
}
