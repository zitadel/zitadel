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

type instanceDomainRelationalProjection struct{}

func newInstanceDomainRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(instanceDomainRelationalProjection))
}

func (*instanceDomainRelationalProjection) Name() string {
	return "zitadel.instance_domains"
}

func (p *instanceDomainRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceDomainAddedEventType,
					Reduce: p.reduceDomainAdded,
				},
				{
					Event:  instance.InstanceDomainPrimarySetEventType,
					Reduce: p.reduceDomainPrimarySet,
				},
				{
					Event:  instance.InstanceDomainRemovedEventType,
					Reduce: p.reduceDomainRemoved,
				},
			},
		},
	}
}

func (p *instanceDomainRelationalProjection) reduceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-DU0xF", "reduce.wrong.event.type %s", instance.InstanceDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-bXCa6", "reduce.wrong.db.pool %T", ex)
		}
		return repository.InstanceRepository(v3_sql.SQLTx(tx)).Domains(false).Add(ctx, &domain.AddInstanceDomain{
			InstanceID:  e.Aggregate().InstanceID,
			Domain:      e.Domain,
			IsGenerated: e.Generated,
			CreatedAt:   e.CreationDate(),
			UpdatedAt:   e.CreationDate(),
		})
	}), nil
}

func (p *instanceDomainRelationalProjection) reduceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-TdEWA", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-QnjHo", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceRepository(v3_sql.SQLTx(tx)).Domains(false)
		_, err := domainRepo.Update(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
			),
			domainRepo.SetPrimary(),
			domainRepo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *instanceDomainRelationalProjection) reduceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Hhcdl", "reduce.wrong.event.type %s", instance.InstanceDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-58ghE", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceRepository(v3_sql.SQLTx(tx)).Domains(false)
		_, err := domainRepo.Remove(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
			),
		)
		return err
	}), nil
}
