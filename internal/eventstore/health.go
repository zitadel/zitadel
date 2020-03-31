package eventstore

import (
	"context"
)

func (app *eventstore) Health(ctx context.Context) error {
	return app.repo.Health(ctx)
}
