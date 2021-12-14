package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type ApplicationRepo struct {
	Commands *command.Commands
	Query    *query.Queries
}

func (a *ApplicationRepo) AuthorizeClientIDSecret(ctx context.Context, clientID, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := a.Query.AppByClientID(ctx, clientID)
	if err != nil {
		return err
	}
	if app.OIDCConfig != nil {
		return a.Commands.VerifyOIDCClientSecret(ctx, app.ProjectID, app.ID, secret)
	}
	return a.Commands.VerifyAPIClientSecret(ctx, app.ProjectID, app.ID, secret)
}
