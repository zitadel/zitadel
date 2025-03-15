package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/orchestrate"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type Instance struct {
	db database.Pool

	orchestrator instanceOrchestrator
}

type instanceOrchestrator interface {
	ByID(ctx context.Context, querier database.Querier, id string) (*repository.Instance, error)
	ByDomain(ctx context.Context, querier database.Querier, domain string) (*repository.Instance, error)
	SetUp(ctx context.Context, tx database.Transaction, instance *repository.Instance) (*repository.Instance, error)
}

func NewInstance(db database.Pool, tracer *tracing.Tracer, logger *logging.Logger) *Instance {
	b := &Instance{
		db:           db,
		orchestrator: orchestrate.Instance(),
	}

	return b
}

func (b *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	return b.orchestrator.ByID(ctx, b.db, id)
}

func (b *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	return b.orchestrator.ByDomain(ctx, b.db, domain)
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
	_, err = b.orchestrator.SetUp(ctx, tx, request.Instance)
	if err != nil {
		return err
	}
	return b.userCommandRepo(tx).Create(ctx, request.User)
}
