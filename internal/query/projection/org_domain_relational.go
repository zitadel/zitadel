package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceOrganizationDomainAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ZX9Fw", "reduce.wrong.event.type %s", org.OrgDomainAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		return repository.OrganizationDomainRepository().Add(ctx, v3_sql.SQLTx(tx), &domain.AddOrganizationDomain{
			InstanceID: e.Aggregate().InstanceID,
			OrgID:      e.Aggregate().ResourceOwner,
			Domain:     e.Domain,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
		})
	}), nil
}

func (p *relationalTablesProjection) reduceOrganizationDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainPrimarySetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dmFdb", "reduce.wrong.event.type %s", org.OrgDomainPrimarySetEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-h6xF0", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationDomainRepository()
		_, err := domainRepo.Update(ctx, v3_sql.SQLTx(tx),
			domainRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ResourceOwner, e.Domain),
			domainRepo.SetPrimary(),
			domainRepo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceOrganizationDomainRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-MzC0n", "reduce.wrong.event.type %s", org.OrgDomainRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-X8oS8", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationDomainRepository()
		_, err := domainRepo.Remove(ctx, v3_sql.SQLTx(tx),
			domainRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ResourceOwner, e.Domain),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceOrganizationDomainVerificationAdded(event eventstore.Event) (*handler.Statement, error) {
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
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-yF03i", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationDomainRepository()

		_, err := domainRepo.Update(ctx, v3_sql.SQLTx(tx),
			domainRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ResourceOwner, e.Domain),
			domainRepo.SetValidationType(validationType),
			domainRepo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceOrganizationDomainVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.DomainVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-7WrI2", "reduce.wrong.event.type %s", org.OrgDomainVerifiedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-0ZGqC", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationDomainRepository()

		_, err := domainRepo.Update(ctx, v3_sql.SQLTx(tx),
			domainRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ResourceOwner, e.Domain),
			domainRepo.SetVerified(),
			domainRepo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}
