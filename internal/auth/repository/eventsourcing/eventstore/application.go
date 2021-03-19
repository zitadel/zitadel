package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/project/model"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type ApplicationRepo struct {
	View          *view.View
	ProjectEvents *proj_event.ProjectEventstore
}

func (a *ApplicationRepo) ApplicationByClientID(ctx context.Context, clientID string) (*model.ApplicationView, error) {
	app, err := a.View.ApplicationByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return proj_view_model.ApplicationViewToModel(app), nil
}

func (a *ApplicationRepo) AuthorizeClientIDSecret(ctx context.Context, clientID, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := a.View.ApplicationByClientID(ctx, clientID)
	if err != nil {
		return err
	}
	return a.ProjectEvents.VerifyClientSecret(ctx, app.ProjectID, app.ID, secret)
}
