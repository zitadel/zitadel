package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type defaults struct {
	db database.Pool
}

type clientSetter interface {
	setClient(database.QueryExecutor)
}

func (d *defaults) acquire(ctx context.Context, setter clientSetter) {
	d.db.Acquire(ctx)
	setter.setClient(d.db)
}
