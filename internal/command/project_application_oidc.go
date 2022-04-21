package command

import (
	"context"
	"strings"
	"time"

	"github.com/caos/logging"

	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	project_repo "github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type addOIDCApp struct {
	AddApp
	Version                  domain.OIDCVersion
	RedirectUris             []string
	ResponseTypes            []domain.OIDCResponseType
	GrantTypes               []domain.OIDCGrantType
	ApplicationType          domain.OIDCApplicationType
	AuthMethodType           domain.OIDCAuthMethodType
	PostLogoutRedirectUris   []string
	DevMode                  bool
	AccessTokenType          domain.OIDCTokenType
	AccessTokenRoleAssertion bool
	IDTokenRoleAssertion     bool
	IDTokenUserinfoAssertion bool
	ClockSkew                time.Duration
	AdditionalOrigins        []string

	ClientID          string
	ClientSecret      *crypto.CryptoValue
	ClientSecretPlain string
}

//AddOIDCAppCommand prepares the commands to add an oidc app. The ClientID will be set during the CreateCommands
func AddOIDCAppCommand(app *addOIDCApp, clientSecretAlg crypto.HashAlgorithm) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if app.ID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument")
		}

		if app.Name = strings.TrimSpace(app.Name); app.Name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument")
		}

		if app.ClockSkew > time.Second*5 || app.ClockSkew < 0 {
			return nil, errors.ThrowInvalidArgument(nil, "V2-PnCMS", "Errors.Invalid.Argument")
		}

		for _, origin := range app.AdditionalOrigins {
			if !http_util.IsOrigin(origin) {
				return nil, errors.ThrowInvalidArgument(nil, "V2-DqWPX", "Errors.Invalid.Argument")
			}
		}

		if !domain.ContainsRequiredGrantTypes(app.ResponseTypes, app.GrantTypes) {
			return nil, errors.ThrowInvalidArgument(nil, "V2-sLpW1", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			project, err := projectWriteModel(ctx, filter, app.Aggregate.ID, app.Aggregate.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, errors.ThrowNotFound(err, "PROJE-6swVG", "Errors.Project.NotFound")
			}

			app.ClientID, err = domain.NewClientID(id.SonyFlakeGenerator, project.Name)
			if err != nil {
				return nil, errors.ThrowInternal(err, "V2-VMSQ1", "Errors.Internal")
			}

			if app.AuthMethodType == domain.OIDCAuthMethodTypeBasic || app.AuthMethodType == domain.OIDCAuthMethodTypePost {
				app.ClientSecret, app.ClientSecretPlain, err = newAppClientSecret(ctx, filter, clientSecretAlg)
				if err != nil {
					return nil, err
				}
			}

			return []eventstore.Command{
				project_repo.NewApplicationAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.ID,
					app.Name,
				),
				project_repo.NewOIDCConfigAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.Version,
					app.ID,
					app.ClientID,
					app.ClientSecret,
					app.RedirectUris,
					app.ResponseTypes,
					app.GrantTypes,
					app.ApplicationType,
					app.AuthMethodType,
					app.PostLogoutRedirectUris,
					app.DevMode,
					app.AccessTokenType,
					app.AccessTokenRoleAssertion,
					app.IDTokenRoleAssertion,
					app.IDTokenUserinfoAssertion,
					app.ClockSkew,
					app.AdditionalOrigins,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) AddOIDCApplication(ctx context.Context, application *domain.OIDCApp, resourceOwner string, appSecretGenerator crypto.Generator) (_ *domain.OIDCApp, err error) {
	if application == nil || application.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-34Fm0", "Errors.Application.Invalid")
	}
	project, err := c.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "PROJECT-3m9ss", "Errors.Project.NotFound")
	}
	addedApplication := NewOIDCApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	events, stringPw, err := c.addOIDCApplication(ctx, projectAgg, project, application, resourceOwner, appSecretGenerator)
	if err != nil {
		return nil, err
	}
	addedApplication.AppID = application.AppID
	pushedEvents, err := c.eventstore.Push(ctx, events...)
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

func (c *Commands) addOIDCApplication(ctx context.Context, projectAgg *eventstore.Aggregate, proj *domain.Project, oidcApp *domain.OIDCApp, resourceOwner string, appSecretGenerator crypto.Generator) (events []eventstore.Command, stringPW string, err error) {
	if oidcApp.AppName == "" || !oidcApp.IsValid() {
		return nil, "", caos_errs.ThrowInvalidArgument(nil, "PROJECT-1n8df", "Errors.Application.Invalid")
	}
	oidcApp.AppID, err = c.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}

	events = []eventstore.Command{
		project_repo.NewApplicationAddedEvent(ctx, projectAgg, oidcApp.AppID, oidcApp.AppName),
	}

	var stringPw string
	err = domain.SetNewClientID(oidcApp, c.idGenerator, proj)
	if err != nil {
		return nil, "", err
	}
	stringPw, err = domain.SetNewClientSecretIfNeeded(oidcApp, appSecretGenerator)
	if err != nil {
		return nil, "", err
	}
	events = append(events, project_repo.NewOIDCConfigAddedEvent(ctx,
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
		oidcApp.ClockSkew,
		oidcApp.AdditionalOrigins))

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
	if !existingOIDC.IsOIDC() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-GBr34", "Errors.Project.App.IsNotOIDC")
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
		oidc.ClockSkew,
		oidc.AdditionalOrigins)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1m88i", "Errors.NoChangesFound")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
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

func (c *Commands) ChangeOIDCApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string, appSecretGenerator crypto.Generator) (*domain.OIDCApp, error) {
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
	if !existingOIDC.IsOIDC() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Ghrh3", "Errors.Project.App.IsNotOIDC")
	}
	cryptoSecret, stringPW, err := domain.NewClientSecret(appSecretGenerator)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, project_repo.NewOIDCConfigSecretChangedEvent(ctx, projectAgg, appID, cryptoSecret))
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
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.NotExisting")
	}
	if !app.IsOIDC() {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-BHgn2", "Errors.Project.App.IsNotOIDC")
	}
	if app.ClientSecret == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.OIDCConfigInvalid")
	}

	projectAgg := ProjectAggregateFromWriteModel(&app.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(app.ClientSecret, []byte(secret), c.userPasswordAlg)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		_, err = c.eventstore.Push(ctx, project_repo.NewOIDCConfigSecretCheckSucceededEvent(ctx, projectAgg, app.AppID))
		return err
	}
	_, err = c.eventstore.Push(ctx, project_repo.NewOIDCConfigSecretCheckFailedEvent(ctx, projectAgg, app.AppID))
	logging.New().OnError(err).Error("could not push event OIDCClientSecretCheckFailed")
	return caos_errs.ThrowInvalidArgument(nil, "COMMAND-Bz542", "Errors.Project.App.ClientSecretInvalid")
}

func (c *Commands) getOIDCAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*OIDCApplicationWriteModel, error) {
	appWriteModel := NewOIDCApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
