package projection

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
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
					Reduce: p.reduceCustomDomainAdded,
				},
				{
					Event:  instance.InstanceDomainPrimarySetEventType,
					Reduce: p.reduceDomainPrimarySet,
				},
				{
					Event:  instance.InstanceDomainRemovedEventType,
					Reduce: p.reduceCustomDomainRemoved,
				},
				{
					Event:  instance.TrustedDomainAddedEventType,
					Reduce: p.reduceTrustedDomainAdded,
				},
				{
					Event:  instance.TrustedDomainRemovedEventType,
					Reduce: p.reduceTrustedDomainRemoved,
				},
			},
		},
	}
}

func (p *instanceDomainRelationalProjection) reduceCustomDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-DU0xF", "reduce.wrong.event.type %s", instance.InstanceDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-bXCa6", "reduce.wrong.db.pool %T", ex)
		}
		return repository.InstanceRepository(tx).Domains(false).Add(ctx, &domain.AddInstanceDomain{
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

func (p *instanceDomainRelationalProjection) reduceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-TdEWA", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-QnjHo", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceRepository(tx).Domains(false)

		condition := database.And(
			domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
			domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
			domainRepo.TypeCondition(domain.DomainTypeCustom),
		)

		_, err := domainRepo.Update(ctx,
			condition,
			domainRepo.SetPrimary(),
		)
		if err != nil {
			return err
		}
		// we need to split the update into two statements because multiple events can have the same creation date
		// therefore we first do not set the updated_at timestamp
		_, err = domainRepo.Update(ctx,
			condition,
			domainRepo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *instanceDomainRelationalProjection) reduceCustomDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Hhcdl", "reduce.wrong.event.type %s", instance.InstanceDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-58ghE", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceRepository(tx).Domains(false)
		_, err := domainRepo.Remove(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
				domainRepo.TypeCondition(domain.DomainTypeCustom),
			),
		)
		return err
	}), nil
}

func (p *instanceDomainRelationalProjection) reduceTrustedDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.TrustedDomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-svHDh", "reduce.wrong.event.type %s", instance.TrustedDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-gx7tQ", "reduce.wrong.db.pool %T", ex)
		}
		return repository.InstanceRepository(tx).Domains(false).Add(ctx, &domain.AddInstanceDomain{
			InstanceID: e.Aggregate().InstanceID,
			Domain:     e.Domain,
			Type:       domain.DomainTypeTrusted,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
		})
	}), nil
}

func (p *instanceDomainRelationalProjection) reduceTrustedDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.TrustedDomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-4K74E", "reduce.wrong.event.type %s", instance.TrustedDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-D68ap", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.InstanceRepository(tx).Domains(false)
		_, err := domainRepo.Remove(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
				domainRepo.TypeCondition(domain.DomainTypeTrusted),
			),
		)
		return err
	}), nil
}
