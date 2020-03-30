package eventstore

import (
	"context"
)

func (app *app) Health(ctx context.Context) error {
	return app.repo.Health(ctx)
}
