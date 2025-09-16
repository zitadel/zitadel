package projection

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type orgMetadataRelationalProjection struct{}

func newOrgMetadataRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(orgMetadataRelationalProjection))
}

func (*orgMetadataRelationalProjection) Name() string {
	return "zitadel.org_metadata"
}

func (p *orgMetadataRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.MetadataSetType,
					Reduce: p.reduceSet,
				},
				{
					Event:  org.MetadataRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  org.MetadataRemovedAllType,
					Reduce: p.reduceRemovedAll,
				},
			},
		},
	}
}

func (p *orgMetadataRelationalProjection) reduceSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MetadataSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOO4e", "reduce.wrong.event.type %s", org.MetadataSetType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-xg4IJ", "reduce.wrong.db.pool %T", ex)
		}
		if !json.Valid(e.Value) {
			logging.Info("invalid")
		} else {
			logging.Info("valid")
		}
		return repository.OrganizationRepository(v3_sql.SQLTx(tx)).Metadata(false).Set(ctx, &domain.OrganizationMetadata{
			Metadata: domain.Metadata{
				InstanceID: e.Aggregate().InstanceID,
				Key:        e.Key,
				Value:      (*json.RawMessage)(&e.Value),
				CreatedAt:  e.CreationDate(),
				UpdatedAt:  e.CreationDate(),
			},
			OrgID: e.Aggregate().ResourceOwner,
		})
	}), nil
}

func (p *orgMetadataRelationalProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MetadataRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-XE6TF", "reduce.wrong.event.type %s", org.MetadataRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-QKMlz", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationRepository(v3_sql.SQLTx(tx)).Metadata(false)
		_, err := domainRepo.Remove(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.OrgIDCondition(e.Aggregate().ResourceOwner),
				domainRepo.KeyCondition(database.TextOperationEqual, e.Key),
			),
		)
		return err
	}), nil
}

func (p *orgMetadataRelationalProjection) reduceRemovedAll(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.MetadataRemovedAllEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-EmyAe", "reduce.wrong.event.type %s", org.MetadataRemovedAllType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-hjEHg", "reduce.wrong.db.pool %T", ex)
		}
		domainRepo := repository.OrganizationRepository(v3_sql.SQLTx(tx)).Metadata(false)
		_, err := domainRepo.Remove(ctx,
			database.And(
				domainRepo.InstanceIDCondition(e.Aggregate().InstanceID),
				domainRepo.OrgIDCondition(e.Aggregate().ResourceOwner),
			),
		)
		return err
	}), nil
}
