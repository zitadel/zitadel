package projection

import (
	"context"
	"database/sql"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceCustomInstanceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-DU0xF", "reduce.wrong.event.type %s", instance.InstanceDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-bXCa6", "reduce.wrong.db.pool %T", ex)
		}
		return repository.InstanceDomainRepository().Add(ctx, v3_sql.SQLTx(tx), &domain.AddInstanceDomain{
			InstanceID:  e.Aggregate().InstanceID,
			Domain:      e.Domain,
			IsPrimary:   gu.Ptr(false),
			IsGenerated: &e.Generated,
			Type:        domain.DomainTypeCustom,
			CreatedAt:   e.CreationDate(),
			UpdatedAt:   e.CreationDate(),
		})
	}), nil
}

func (p *relationalTablesProjection) reduceInstanceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-TdEWA", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-QnjHo", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceDomainRepository()

		_, err := domainRepo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				domainRepo.PrimaryKeyCondition(e.Domain),
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.TypeCondition(domain.DomainTypeCustom),
			),
			domainRepo.SetPrimary(),
			domainRepo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceCustomInstanceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Hhcdl", "reduce.wrong.event.type %s", instance.InstanceDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-58ghE", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceDomainRepository()
		_, err := domainRepo.Remove(ctx, v3_sql.SQLTx(tx),
			database.And(
				domainRepo.PrimaryKeyCondition(e.Domain),
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.TypeCondition(domain.DomainTypeCustom),
			),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceTrustedInstanceDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.TrustedDomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-svHDh", "reduce.wrong.event.type %s", instance.TrustedDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-gx7tQ", "reduce.wrong.db.pool %T", ex)
		}
		return repository.InstanceDomainRepository().Add(ctx, v3_sql.SQLTx(tx), &domain.AddInstanceDomain{
			InstanceID: e.Aggregate().InstanceID,
			Domain:     e.Domain,
			Type:       domain.DomainTypeTrusted,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
		})
	}), nil
}

func (p *relationalTablesProjection) reduceTrustedInstanceDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.TrustedDomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-4K74E", "reduce.wrong.event.type %s", instance.TrustedDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-D68ap", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceDomainRepository()
		_, err := domainRepo.Remove(ctx, v3_sql.SQLTx(tx),
			database.And(
				domainRepo.PrimaryKeyCondition(e.Domain),
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.TypeCondition(domain.DomainTypeTrusted),
			),
		)
		return err
	}), nil
}
