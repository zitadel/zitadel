package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddOIDCApplication(ctx context.Context, application *domain.OIDCApp, resourceOwner string) (_ *domain.OIDCApp, err error) {
	if application == nil || application.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-34Fm0", "Errors.Application.Invalid")
	}
	project, err := c.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "PROJECT-3m9ss", "Errors.Project.NotFound")
	}
	addedApplication := NewOIDCApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	events, stringPw, err := c.addOIDCApplication(ctx, projectAgg, project, application, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication.AppID = application.AppID
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedApplication, pushedEvents...)
	if err != nil {
		return nil, err
	}
	result := oidcWriteModelToOIDCConfig(addedApplication)
	result.ClientSecretString = stringPw
	result.FillCompliance()
	return result, nil
}

func (c *Commands) addOIDCApplication(ctx context.Context, projectAgg *eventstore.Aggregate, proj *domain.Project, oidcApp *domain.OIDCApp, resourceOwner string) (events []eventstore.EventPusher, stringPW string, err error) {
	if oidcApp.AppName == "" || !oidcApp.IsValid() {
		return nil, "", caos_errs.ThrowInvalidArgument(nil, "PROJECT-1n8df", "Errors.Application.Invalid")
	}
	oidcApp.AppID, err = c.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}

	events = []eventstore.EventPusher{
		project.NewApplicationAddedEvent(ctx, projectAgg, oidcApp.AppID, oidcApp.AppName),
	}

	var stringPw string
	err = domain.SetNewClientID(oidcApp, c.idGenerator, proj)
	if err != nil {
		return nil, "", err
	}
	stringPw, err = domain.SetNewClientSecretIfNeeded(oidcApp, c.applicationSecretGenerator)
	if err != nil {
		return nil, "", err
	}
	events = append(events, project.NewOIDCConfigAddedEvent(ctx,
		projectAgg,
		oidcApp.OIDCVersion,
		oidcApp.AppID,
		oidcApp.ClientID,
		oidcApp.ClientSecret,
		oidcApp.RedirectUris,
		oidcApp.ResponseTypes,
		oidcApp.GrantTypes,
		oidcApp.ApplicationType,
		oidcApp.AuthMethodType,
		oidcApp.PostLogoutRedirectUris,
		oidcApp.DevMode,
		oidcApp.AccessTokenType,
		oidcApp.AccessTokenRoleAssertion,
		oidcApp.IDTokenRoleAssertion,
		oidcApp.IDTokenUserinfoAssertion,
		oidcApp.ClockSkew))

	return events, stringPw, nil
}

func (c *Commands) ChangeOIDCApplication(ctx context.Context, oidc *domain.OIDCApp, resourceOwner string) (*domain.OIDCApp, error) {
	if !oidc.IsValid() || oidc.AppID == "" || oidc.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5m9fs", "Errors.Project.App.OIDCConfigInvalid")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, oidc.AggregateID, oidc.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)
	changedEvent, hasChanged, err := existingOIDC.NewChangedEvent(
		ctx,
		projectAgg,
		oidc.AppID,
		oidc.RedirectUris,
		oidc.PostLogoutRedirectUris,
		oidc.ResponseTypes,
		oidc.GrantTypes,
		oidc.ApplicationType,
		oidc.AuthMethodType,
		oidc.OIDCVersion,
		oidc.AccessTokenType,
		oidc.DevMode,
		oidc.AccessTokenRoleAssertion,
		oidc.IDTokenRoleAssertion,
		oidc.IDTokenUserinfoAssertion,
		oidc.ClockSkew)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1m88i", "Errors.NoChangesFound")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingOIDC, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := oidcWriteModelToOIDCConfig(existingOIDC)
	result.FillCompliance()
	return result, nil
}

func (c *Commands) ChangeOIDCApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string) (*domain.OIDCApp, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	cryptoSecret, stringPW, err := domain.NewClientSecret(c.applicationSecretGenerator)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)

	pushedEvents, err := c.eventstore.PushEvents(ctx, project.NewOIDCConfigSecretChangedEvent(ctx, projectAgg, appID, cryptoSecret))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingOIDC, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := oidcWriteModelToOIDCConfig(existingOIDC)
	result.ClientSecretString = stringPW
	return result, err
}

func (c *Commands) VerifyOIDCClientSecret(ctx context.Context, projectID, appID, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := c.getOIDCAppWriteModel(ctx, projectID, appID, "")
	if err != nil {
		return err
	}
	if !app.State.Exists() {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.NoExisting")
	}
	if app.ClientSecret == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.OIDCConfigInvalid")
	}

	projectAgg := ProjectAggregateFromWriteModel(&app.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(app.ClientSecret, []byte(secret), c.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		_, err = c.eventstore.PushEvents(ctx, project.NewOIDCConfigSecretCheckSucceededEvent(ctx, projectAgg, app.AppID))
		return err
	}
	_, err = c.eventstore.PushEvents(ctx, project.NewOIDCConfigSecretCheckFailedEvent(ctx, projectAgg, app.AppID))
	logging.Log("COMMAND-ADfhz").OnError(err).Error("could not push event OIDCClientSecretCheckFailed")
	return caos_errs.ThrowInvalidArgument(nil, "COMMAND-Bz542", "Errors.Project.App.OIDCSecretInvalid")
}

func (c *Commands) getOIDCAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*OIDCApplicationWriteModel, error) {
	appWriteModel := NewOIDCApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
