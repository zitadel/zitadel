package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/project/model"
	proj_view_model "github.com/caos/zitadel/internal/project/repository/view/model"
)

type ApplicationRepo struct {
	View        *view.View
	PasswordAlg crypto.HashAlgorithm
}

func (a *ApplicationRepo) ApplicationByClientID(ctx context.Context, clientID string) (*model.ApplicationView, error) {
	app, err := a.View.ApplicationByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return proj_view_model.ApplicationViewToModel(app), nil
}

func (a *ApplicationRepo) AuthorizeOIDCApplication(ctx context.Context, clientID, secret string) error {
	app, err := a.View.ApplicationByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	app.
	app, err := a.View.ApplicationByClientID(ctx, clientID)
	if err != nil {

	}
	crypto.CompareHash(app.oidc)
}
