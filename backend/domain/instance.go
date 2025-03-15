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

	instance instanceOrchestrator
	user     userOrchestrator
}

type instanceOrchestrator interface {
	ByID(ctx context.Context, querier database.Querier, id string) (*repository.Instance, error)
	ByDomain(ctx context.Context, querier database.Querier, domain string) (*repository.Instance, error)
	SetUp(ctx context.Context, tx database.Transaction, instance *repository.Instance) (*repository.Instance, error)
}

func NewInstance(db database.Pool, tracer *tracing.Tracer, logger *logging.Logger) *Instance {
	b := &Instance{
		db:       db,
		instance: orchestrate.Instance(),
		user:     orchestrate.User(),
	}

	return b
}

func (b *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	return b.instance.ByID(ctx, b.db, id)
}

func (b *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	return b.instance.ByDomain(ctx, b.db, domain)
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
	_, err = b.instance.SetUp(ctx, tx, request.Instance)
	if err != nil {
		return err
	}
	_, err = b.user.Create(ctx, tx, request.User)
	return err
}
