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
	"github.com/zitadel/zitadel/internal/zerrors"
)

const InstanceRelationalProjectionTable = "zitadel.instances"

type instanceRelationalProjection struct{}

func newInstanceRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceRelationalProjection))
}

func (*instanceRelationalProjection) Name() string {
	return InstanceRelationalProjectionTable
}

func (p *instanceRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceAddedEventType,
					Reduce: p.reduceInstanceAdded,
				},
				{
					Event:  instance.InstanceChangedEventType,
					Reduce: p.reduceInstanceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: p.reduceInstanceDelete,
				},
				{
					Event:  instance.DefaultOrgSetEventType,
					Reduce: p.reduceDefaultOrgSet,
				},
				{
					Event:  instance.ProjectSetEventType,
					Reduce: p.reduceIAMProjectSet,
				},
				{
					Event:  instance.ConsoleSetEventType,
					Reduce: p.reduceConsoleSet,
				},
				{
					Event:  instance.DefaultLanguageSetEventType,
					Reduce: p.reduceDefaultLanguageSet,
				},
			},
		},
	}
}

func (p *instanceRelationalProjection) reduceInstanceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-29nRr", "reduce.wrong.event.type %s", instance.InstanceAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		return repository.InstanceRepository(v3_sql.SQLTx(tx)).Create(ctx, &domain.Instance{
			ID:        e.Aggregate().ID,
			Name:      e.Name,
			CreatedAt: e.CreationDate(),
			UpdatedAt: e.CreationDate(),
		})
	}), nil
}

func (p *instanceRelationalProjection) reduceInstanceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-so2am1", "reduce.wrong.event.type %s", instance.InstanceChangedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.InstanceRepository(v3_sql.SQLTx(tx))
		return p.updateInstance(ctx, event, repo, repo.SetName(e.Name))
	}), nil
}

func (p *instanceRelationalProjection) reduceInstanceDelete(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-so2am1", "reduce.wrong.event.type %s", instance.InstanceChangedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repository.InstanceRepository(v3_sql.SQLTx(tx)).Delete(ctx, e.Aggregate().ID)
		return err
	}), nil
}

func (p *instanceRelationalProjection) reduceDefaultOrgSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DefaultOrgSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-2n9f2", "reduce.wrong.event.type %s", instance.DefaultOrgSetEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.InstanceRepository(v3_sql.SQLTx(tx))
		return p.updateInstance(ctx, event, repo, repo.SetDefaultOrg(e.OrgID))
	}), nil
}

func (p *instanceRelationalProjection) reduceIAMProjectSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.ProjectSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-30o0e", "reduce.wrong.event.type %s", instance.ProjectSetEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.InstanceRepository(v3_sql.SQLTx(tx))
		return p.updateInstance(ctx, event, repo, repo.SetIAMProject(e.ProjectID))
	}), nil
}

func (p *instanceRelationalProjection) reduceConsoleSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.ConsoleSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Dgf11", "reduce.wrong.event.type %s", instance.ConsoleSetEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.InstanceRepository(v3_sql.SQLTx(tx))
		return p.updateInstance(ctx, event, repo, repo.SetConsoleClientID(e.ClientID), repo.SetConsoleAppID(e.AppID))
	}), nil
}

func (p *instanceRelationalProjection) reduceDefaultLanguageSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DefaultLanguageSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-30o0e", "reduce.wrong.event.type %s", instance.DefaultLanguageSetEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-rVUyy", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.InstanceRepository(v3_sql.SQLTx(tx))
		return p.updateInstance(ctx, event, repo, repo.SetDefaultLanguage(e.Language))
	}), nil
}

func (p *instanceRelationalProjection) updateInstance(ctx context.Context, event eventstore.Event, repo domain.InstanceRepository, changes ...database.Change) error {
	_, err := repo.Update(ctx, event.Aggregate().ID, changes...)
	if err != nil {
		return err
	}

	instance, err := repo.Get(ctx, database.WithCondition(repo.IDCondition(event.Aggregate().ID)))
	if err != nil {
		return err
	}
	if instance.UpdatedAt.Equal(event.CreatedAt()) {
		return nil
	}
	// we need to split the update into two statements because multiple events can have the same creation date
	// therefore we first do not set the updated_at timestamp
	_, err = repo.Update(ctx,
		event.Aggregate().ID,
		repo.SetUpdatedAt(event.CreatedAt()),
	)
	return err
}
