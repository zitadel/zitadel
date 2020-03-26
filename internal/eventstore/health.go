package eventstore

import (
	"context"
)

func (app *app) Health(ctx context.Context) error {
	return app.eventstore.Health(ctx)
}
