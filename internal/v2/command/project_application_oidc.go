package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddOIDCApplication(ctx context.Context, application *domain.OIDCApp, resourceOwner string) (_ *domain.OIDCApp, err error) {
	project, err := r.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication := NewOIDCApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	stringPw, err := r.addOIDCApplication(ctx, projectAgg, project, application, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication.AppID = application.AppID
	err = r.eventstore.PushAggregate(ctx, addedApplication, projectAgg)
	if err != nil {
		return nil, err
	}

	result := oidcWriteModelToOIDCConfig(addedApplication)
	result.ClientSecretString = stringPw
	result.FillCompliance()
	return result, nil
}

func (r *CommandSide) addOIDCApplication(ctx context.Context, projectAgg *project.Aggregate, proj *domain.Project, oidcApp *domain.OIDCApp, resourceOwner string) (stringPW string, err error) {
	if !oidcApp.IsValid() {
		return "", caos_errs.ThrowPreconditionFailed(nil, "PROJECT-Bff2g", "Errors.Application.Invalid")
	}
	oidcApp.AppID, err = r.idGenerator.Next()
	if err != nil {
		return "", err
	}

	projectAgg.PushEvents(project.NewApplicationAddedEvent(ctx, oidcApp.AppID, oidcApp.AppName, resourceOwner))

	var stringPw string
	err = oidcApp.GenerateNewClientID(r.idGenerator, proj)
	if err != nil {
		return "", err
	}
	stringPw, err = oidcApp.GenerateClientSecretIfNeeded(r.applicationSecretGenerator)
	if err != nil {
		return "", err
	}
	projectAgg.PushEvents(project.NewOIDCConfigAddedEvent(ctx,
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

	return stringPw, nil
}

func (r *CommandSide) ChangeOIDCApplication(ctx context.Context, oidc *domain.OIDCApp, resourceOwner string) (*domain.OIDCApp, error) {
	if !oidc.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1m900", "Errors.Project.App.OIDCConfigInvalid")
	}

	existingOIDC, err := r.getOIDCAppWriteModel(ctx, oidc.AggregateID, oidc.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	changedEvent, hasChanged, err := existingOIDC.NewChangedEvent(
		ctx,
		oidc.AppID,
		oidc.ClientID,
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
	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)
	projectAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingOIDC, projectAgg)
	if err != nil {
		return nil, err
	}
	result := oidcWriteModelToOIDCConfig(existingOIDC)
	result.FillCompliance()
	return result, nil
}

func (r *CommandSide) ChangeOIDCApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string) (*domain.OIDCApp, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingOIDC, err := r.getOIDCAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	cryptoSecret, stringPW, err := domain.NewClientSecret(r.applicationSecretGenerator)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)
	projectAgg.PushEvents(project.NewOIDCConfigSecretChangedEvent(ctx, appID, cryptoSecret))

	err = r.eventstore.PushAggregate(ctx, existingOIDC, projectAgg)
	if err != nil {
		return nil, err
	}

	result := oidcWriteModelToOIDCConfig(existingOIDC)
	result.ClientSecretString = stringPW
	return result, err
}
func (r *CommandSide) getOIDCAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*OIDCApplicationWriteModel, error) {
	appWriteModel := NewOIDCApplicationWriteModelWithAppIDC(projectID, appID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
