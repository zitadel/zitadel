package domain

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
	"github.com/zitadel/zitadel/backend/storage/repository"
	"github.com/zitadel/zitadel/backend/storage/repository/cache"
	"github.com/zitadel/zitadel/backend/storage/repository/event"
	"github.com/zitadel/zitadel/backend/storage/repository/sql"
	"github.com/zitadel/zitadel/backend/storage/repository/telemetry/logged"
	"github.com/zitadel/zitadel/backend/storage/repository/telemetry/traced"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type Instance struct {
	db     database.Pool
	tracer *tracing.Tracer
	logger *slog.Logger
	cache  *cache.Instance
}

func NewInstance(db database.Pool, tracer *tracing.Tracer, logger *slog.Logger) *Instance {
	b := &Instance{
		db:     db,
		tracer: tracer,
		logger: logger,

		cache: &cache.Instance{},
	}

	return b
}

func (b *Instance) instanceCommandRepo(tx database.Transaction) repository.InstanceRepository {
	return logged.NewInstance(
		b.logger,
		traced.NewInstance(
			b.tracer,
			event.NewInstance(
				eventstore.New(tx),
				b.cache.SetNext(
					sql.NewInstance(tx),
				),
			),
		),
	)
}

func (b *Instance) instanceQueryRepo(tx database.QueryExecutor) repository.InstanceRepository {
	return logged.NewInstance(
		b.logger,
		traced.NewInstance(
			b.tracer,
			b.cache.SetNext(
				sql.NewInstance(tx),
			),
		),
	)
}

func (b *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	return b.instanceQueryRepo(b.db).ByID(ctx, id)
}

func (b *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	return b.instanceQueryRepo(b.db).ByDomain(ctx, domain)
}

type SetUpInstance struct {
	Instance *repository.Instance
	User     *repository.User
}

func (b *Instance) SetUp(ctx context.Context, request *SetUpInstance) (err error) {
	tx, err := b.db.Begin(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.End(ctx, err)
	}()
	err = b.instanceCommandRepo(tx).SetUp(ctx, request.Instance)
	if err != nil {
		return err
	}
	return b.userCommandRepo(tx).Create(ctx, request.User)
}
