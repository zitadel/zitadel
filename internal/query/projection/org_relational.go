package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	OrgRelationProjectionTable = "zitadel.organizations"
)

type orgRelationalProjection struct{}

func (*orgRelationalProjection) Name() string {
	return OrgRelationProjectionTable
}

func newOrgRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgRelationalProjection))
}

func (p *orgRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgAddedEventType,
					Reduce: p.reduceOrgRelationalAdded,
				},
				{
					Event:  org.OrgChangedEventType,
					Reduce: p.reduceOrgRelationalChanged,
				},
				{
					Event:  org.OrgDeactivatedEventType,
					Reduce: p.reduceOrgRelationalDeactivated,
				},
				{
					Event:  org.OrgReactivatedEventType,
					Reduce: p.reduceOrgRelationalReactivated,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRelationalRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(OrgColumnInstanceID),
				},
			},
		},
	}
}

func (p *orgRelationalProjection) reduceOrgRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uYq5R", "reduce.wrong.event.type %s", org.OrgAddedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.OrganizationRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.Organization{
			ID:         e.Aggregate().ID,
			Name:       e.Name,
			InstanceID: e.Aggregate().InstanceID,
			State:      domain.OrgStateActive,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreatedAt(),
		})
	}), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Bg9om", "reduce.wrong.event.type %s", org.OrgChangedEventType)
	}
	if e.Name == "" {
		return handler.NewNoOpStatement(e), nil
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.OrganizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetName(e.Name),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BApK5", "reduce.wrong.event.type %s", org.OrgDeactivatedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.OrganizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetState(domain.OrgStateInactive),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-o38DE", "reduce.wrong.event.type %s", org.OrgReactivatedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.OrganizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetState(domain.OrgStateInactive),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *orgRelationalProjection) reduceOrgRelationalRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-DGm9g", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.OrganizationRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx), repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID))
		return err
	}), nil
}
