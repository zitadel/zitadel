package projection

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPIntentRelationalProjectionTable = "zitadel.identity_provider_intents"
)

type idpIntentRelationalProjection struct{}

func (*idpIntentRelationalProjection) Name() string {
	return IDPIntentRelationalProjectionTable
}

func newIDPIntentRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(idpIntentRelationalProjection))
}

// Reducers implements [handler.Projection].
func (i *idpIntentRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: idpintent.AggregateType,
			EventReducers: []handler.EventReducer{
				{Event: idpintent.StartedEventType, Reduce: i.reduceStartedEvent},
				{Event: idpintent.SucceededEventType, Reduce: i.reduceSucceededEvent},
				{Event: idpintent.SAMLSucceededEventType, Reduce: i.reduceSAMLSucceededEvent},
				{Event: idpintent.SAMLRequestEventType, Reduce: i.reduceSAMLRequestEvent},
				{Event: idpintent.LDAPSucceededEventType, Reduce: i.reduceLDAPSucceededEvent},
				{Event: idpintent.FailedEventType, Reduce: i.reduceFailedEvent},
				{Event: idpintent.ConsumedEventType, Reduce: i.reduceConsumedEvent},
			},
		},
	}

}

func (i *idpIntentRelationalProjection) reduceStartedEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.StartedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-Lqj3HB", "reduce.wrong.event.type %s", idpintent.StartedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rirFQm", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDPIntentRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IDPIntent{
			ID:           e.Aggregate().ID,
			InstanceID:   e.Aggregate().InstanceID,
			State:        domain.IDPIntentStateStarted,
			SuccessURL:   e.SuccessURL,
			FailureURL:   e.FailureURL,
			CreatedAt:    e.CreationDate(),
			UpdatedAt:    e.CreationDate(),
			IDPID:        e.IDPID,
			IDPArguments: e.IDPArguments,
		})
	}), nil
}

func (i *idpIntentRelationalProjection) reduceSucceededEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.SucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-hQSxj0", "reduce.wrong.event.type %s", idpintent.SucceededEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-ry4tFH", "reduce.wrong.db.pool %T", ex)
		}
		idpAccessToken, err := e.IDPAccessToken.Value()
		if err != nil {
			return err
		}
		jsonIDPAccessToken, ok := idpAccessToken.([]byte)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-32XE1L", "reduce.convertion.json.failed")
		}
		repo := repository.IDPIntentRepository()

		changes := []database.Change{
			repo.SetState(domain.IDPIntentStateSucceeded),
			repo.SetIDPUser(e.IDPUser),
			repo.SetIDPUserID(e.IDPUserID),
			repo.SetIDPUsername(e.IDPUserName),
			repo.SetUserID(e.UserID),
			repo.SetIDPAccessToken(jsonIDPAccessToken),
			repo.SetIDPIDToken(e.IDPIDToken),
			repo.SetSucceededAt(e.CreationDate()),
			repo.SetExpiresAt(e.ExpiresAt),
		}

		_, err = repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			changes...,
		)
		return err
	}), nil
}

func (i *idpIntentRelationalProjection) reduceSAMLSucceededEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.SAMLSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-t2uaiv", "reduce.wrong.event.type %s", idpintent.SAMLSucceededEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-0QGPXm", "reduce.wrong.db.pool %T", ex)
		}
		assertion, err := e.Assertion.Value()
		if err != nil {
			return err
		}
		jsonAssertion, ok := assertion.([]byte)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-1TFjTi", "reduce.convertion.json.failed")
		}

		repo := repository.IDPIntentRepository()
		_, err = repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetState(domain.IDPIntentStateSucceeded),
			repo.SetIDPUser(e.IDPUser),
			repo.SetIDPUserID(e.IDPUserID),
			repo.SetIDPUsername(e.IDPUserName),
			repo.SetUserID(e.UserID),
			repo.SetAssertion(jsonAssertion),
			repo.SetSucceededAt(e.CreationDate()),
			repo.SetExpiresAt(e.ExpiresAt),
		)
		return err
	}), nil
}

func (i *idpIntentRelationalProjection) reduceSAMLRequestEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.SAMLRequestEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-f7O6K2", "reduce.wrong.event.type %s", idpintent.SAMLRequestEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-ijgAPJ", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDPIntentRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetRequestID(e.RequestID),
		)
		return err
	}), nil
}

func (i *idpIntentRelationalProjection) reduceLDAPSucceededEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.LDAPSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-Gw8ddW", "reduce.wrong.event.type %s", idpintent.LDAPSucceededEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-tNp8aa", "reduce.wrong.db.pool %T", ex)
		}

		jsonAttributes, err := json.Marshal(e.EntryAttributes)
		if err != nil {
			return zerrors.ThrowInternal(nil, "HANDL-LuvCvD", "reduce.convertion.json.failed")
		}

		repo := repository.IDPIntentRepository()
		_, err = repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetState(domain.IDPIntentStateSucceeded),
			repo.SetIDPUser(e.IDPUser),
			repo.SetIDPUserID(e.IDPUserID),
			repo.SetIDPUsername(e.IDPUserName),
			repo.SetUserID(e.UserID),
			repo.SetIDPEntryAttributes(jsonAttributes),
			repo.SetSucceededAt(e.CreationDate()),
			repo.SetExpiresAt(e.ExpiresAt),
		)
		return err
	}), nil
}

func (i *idpIntentRelationalProjection) reduceFailedEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.FailedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-x3o9LZ", "reduce.wrong.event.type %s", idpintent.FailedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iqmR8O", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDPIntentRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetState(domain.IDPIntentStateFailed),
			repo.SetFailReason(e.Reason),
			repo.SetFailedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (i *idpIntentRelationalProjection) reduceConsumedEvent(evt eventstore.Event) (*handler.Statement, error) {
	e, ok := evt.(*idpintent.ConsumedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-jHZs6T", "reduce.wrong.event.type %s", idpintent.ConsumedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-8X9U9w", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDPIntentRepository()
		_, err := repo.Delete(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
		)
		return err
	}), nil
}
