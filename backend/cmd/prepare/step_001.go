package prepare

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
)

type Step001 struct {
	Database database.Pool
}

func (v *Step001) Migrate(ctx context.Context) error {
	conn, err := v.Database.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release(ctx)

	eventstore.New(conn).
	return nil
}
