package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceSessionAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.AddedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		var userAgent *domain.SessionUserAgent
		if e.UserAgent != nil {
			userAgent = &domain.SessionUserAgent{
				FingerprintID: e.UserAgent.FingerprintID,
				Description:   e.UserAgent.Description,
				IP:            e.UserAgent.IP,
				Header:        e.UserAgent.Header,
			}
		}

		sessionRepo := repository.SessionRepository()
		return sessionRepo.Create(ctx, v3Tx, &domain.Session{
			InstanceID: e.Aggregate().InstanceID,
			ID:         e.Aggregate().ID,
			CreatorID:  e.User,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
			UserAgent:  userAgent,
		})
	}), nil
}

func (p *relationalTablesProjection) reduceSessionUserChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.UserCheckedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-UHE92", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionPasswordChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.PasswordCheckedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-3krfa", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionIntentChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.IntentCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-ajkd2", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorIdentityProviderIntent{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceSessionWebAuthNChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.WebAuthNChallengedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-do35d", "reduce.wrong.db.pool %T", ex)
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
func (p *relationalTablesProjection) reduceSessionWebAuthNChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.WebAuthNCheckedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-djk22", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionTOTPChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.TOTPCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-cklr9", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionOTPSMSChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPSMSChallengedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-fk4f9", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionOTPSMSChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPSMSCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-dk39f", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionOTPEmailChallenged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPEmailChallengedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-fkf93", "reduce.wrong.db.pool %T", ex)
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
				URLTemplate:       e.URLTmpl,
				TriggeredAtOrigin: e.TriggeredAtOrigin,
			}),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceSessionOTPEmailChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.OTPEmailCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-3kll0", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionRecoveryCodeChecked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.RecoveryCodeCheckedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-fk45a", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		sessionRepo := repository.SessionRepository()
		condition := sessionRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)
		_, err = sessionRepo.Update(ctx, v3Tx, condition,
			sessionRepo.SetUpdatedAt(e.CreatedAt()),
			sessionRepo.SetFactor(&domain.SessionFactorRecoveryCode{
				LastVerifiedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceSessionTokenSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.TokenSetEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-asddt", "reduce.wrong.db.pool %T", ex)
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

func (p *relationalTablesProjection) reduceSessionMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*session.MetadataSetEvent](event)
	if err != nil {
		return nil, err
	}

	metadataList := make([]*domain.SessionMetadata, 0, len(e.Metadata))
	for key, value := range e.Metadata {
		metadataList = append(metadataList, &domain.SessionMetadata{
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

func (p *relationalTablesProjection) reduceSessionLifetimeSet(event eventstore.Event) (*handler.Statement, error) {
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

func (p *relationalTablesProjection) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
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
		_, err := sessionRepo.Delete(ctx, v3Tx, condition)
		return err
	}), nil
}
