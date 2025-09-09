package projection

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type orgDomainRelationalProjection struct{}

func newOrgDomainRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgDomainRelationalProjection))
}

func (*orgDomainRelationalProjection) Name() string {
	return "zitadel.org_domains"
}

func (p *orgDomainRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgDomainAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reducePrimarySet,
				},
				{
					Event:  org.OrgDomainRemovedEventType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.OrgDomainVerificationAddedEventType,
					Reduce: p.reduceVerificationAdded,
				},
				{
					Event:  org.OrgDomainVerifiedEventType,
					Reduce: p.reduceVerified,
				},
			},
		},
	}
}

func (p *orgDomainRelationalProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ZX9Fw", "reduce.wrong.event.type %s", org.OrgDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		return repository.OrganizationRepository(tx).Domains(false).Add(ctx, &domain.AddOrganizationDomain{
			InstanceID: e.Aggregate().InstanceID,
			OrgID:      e.Aggregate().ResourceOwner,
			Domain:     e.Domain,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
		})
	}), nil
}

func (p *orgDomainRelationalProjection) reducePrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dmFdb", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-h6xF0", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationRepository(tx).Domains(false)
		condition := database.And(
			domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
			domainRepo.OrgIDCondition(e.Aggregate().ResourceOwner),
			domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
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

func (p *orgDomainRelationalProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-MzC0n", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-X8oS8", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationRepository(tx).Domains(false)
		_, err := domainRepo.Remove(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.OrgIDCondition(e.Aggregate().ResourceOwner),
				domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
			),
		)
		return err
	}), nil
}

func (p *orgDomainRelationalProjection) reduceVerificationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerificationAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-oGzip", "reduce.wrong.event.type %s", org.OrgDomainVerificationAddedEventType)
	}
	var validationType domain.DomainValidationType
	switch e.ValidationType {
	case old_domain.OrgDomainValidationTypeDNS:
		validationType = domain.DomainValidationTypeDNS
	case old_domain.OrgDomainValidationTypeHTTP:
		validationType = domain.DomainValidationTypeHTTP
	case old_domain.OrgDomainValidationTypeUnspecified:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-FJfKB", "reduce.unsupported.validation.type %v", e.ValidationType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-yF03i", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationRepository(tx).Domains(false)
		condition := database.And(
			domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
			domainRepo.OrgIDCondition(e.Aggregate().ResourceOwner),
			domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
		)

		_, err := domainRepo.Update(ctx,
			condition,
			domainRepo.SetValidationType(validationType),
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

func (p *orgDomainRelationalProjection) reduceVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-7WrI2", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(database.Transaction)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-0ZGqC", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationRepository(tx).Domains(false)

		condition := database.And(
			domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
			domainRepo.OrgIDCondition(e.Aggregate().ResourceOwner),
			domainRepo.DomainCondition(database.TextOperationEqual, e.Domain),
		)

		_, err := domainRepo.Update(ctx,
			condition,
			domainRepo.SetVerified(),
			domainRepo.SetUpdatedAt(e.CreationDate()),
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
