package command

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	project_repo "github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type addOIDCApp struct {
	AddApp
	Version                     domain.OIDCVersion
	RedirectUris                []string
	ResponseTypes               []domain.OIDCResponseType
	GrantTypes                  []domain.OIDCGrantType
	ApplicationType             domain.OIDCApplicationType
	AuthMethodType              domain.OIDCAuthMethodType
	PostLogoutRedirectUris      []string
	DevMode                     bool
	AccessTokenType             domain.OIDCTokenType
	AccessTokenRoleAssertion    bool
	IDTokenRoleAssertion        bool
	IDTokenUserinfoAssertion    bool
	ClockSkew                   time.Duration
	AdditionalOrigins           []string
	SkipSuccessPageForNativeApp bool

	ClientID          string
	ClientSecret      *crypto.CryptoValue
	ClientSecretPlain string
}

// AddOIDCAppCommand prepares the commands to add an oidc app. The ClientID will be set during the CreateCommands
func (c *Commands) AddOIDCAppCommand(app *addOIDCApp, clientSecretAlg crypto.HashAlgorithm) preparation.Validation {
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

			app.ClientID, err = domain.NewClientID(c.idGenerator, project.Name)
			if err != nil {
				return nil, errors.ThrowInternal(err, "V2-VMSQ1", "Errors.Internal")
			}

			if app.AuthMethodType == domain.OIDCAuthMethodTypeBasic || app.AuthMethodType == domain.OIDCAuthMethodTypePost {
				code, err := c.newAppClientSecret(ctx, filter, clientSecretAlg)
				if err != nil {
					return nil, err
				}
				app.ClientSecret, app.ClientSecretPlain = code.Crypted, code.Plain
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
					app.SkipSuccessPageForNativeApp,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) AddOIDCApplicationWithID(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner, appID string, appSecretGenerator crypto.Generator) (_ *domain.OIDCApp, err error) {
	existingApp, err := c.getOIDCAppWriteModel(ctx, oidcApp.AggregateID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State != domain.AppStateUnspecified {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-lxowmp", "Errors.Project.App.AlreadyExisting")
	}

	project, err := c.getProjectByID(ctx, oidcApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "PROJECT-3m9s2", "Errors.Project.NotFound")
	}

	return c.addOIDCApplicationWithID(ctx, oidcApp, resourceOwner, project, appID, appSecretGenerator)
}

func (c *Commands) AddOIDCApplication(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner string, appSecretGenerator crypto.Generator) (_ *domain.OIDCApp, err error) {
	if oidcApp == nil || oidcApp.AggregateID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "PROJECT-34Fm0", "Errors.Project.App.Invalid")
	}
	project, err := c.getProjectByID(ctx, oidcApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "PROJECT-3m9ss", "Errors.Project.NotFound")
	}

	if oidcApp.AppName == "" || !oidcApp.IsValid() {
		return nil, errors.ThrowInvalidArgument(nil, "PROJECT-1n8df", "Errors.Project.App.Invalid")
	}

	appID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	return c.addOIDCApplicationWithID(ctx, oidcApp, resourceOwner, project, appID, appSecretGenerator)
}

func (c *Commands) addOIDCApplicationWithID(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner string, project *domain.Project, appID string, appSecretGenerator crypto.Generator) (_ *domain.OIDCApp, err error) {

	addedApplication := NewOIDCApplicationWriteModel(oidcApp.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)

	oidcApp.AppID = appID

	events := []eventstore.Command{
		project_repo.NewApplicationAddedEvent(ctx, projectAgg, oidcApp.AppID, oidcApp.AppName),
	}

	var stringPw string
	err = domain.SetNewClientID(oidcApp, c.idGenerator, project)
	if err != nil {
		return nil, err
	}
	stringPw, err = domain.SetNewClientSecretIfNeeded(oidcApp, appSecretGenerator)
	if err != nil {
		return nil, err
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
		oidcApp.AdditionalOrigins,
		oidcApp.SkipNativeAppSuccessPage,
	))

	addedApplication.AppID = oidcApp.AppID
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

func (c *Commands) ChangeOIDCApplication(ctx context.Context, oidc *domain.OIDCApp, resourceOwner string) (*domain.OIDCApp, error) {
	if !oidc.IsValid() || oidc.AppID == "" || oidc.AggregateID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-5m9fs", "Errors.Project.App.OIDCConfigInvalid")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, oidc.AggregateID, oidc.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, errors.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	if !existingOIDC.IsOIDC() {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-GBr34", "Errors.Project.App.IsNotOIDC")
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
		oidc.AdditionalOrigins,
		oidc.SkipNativeAppSuccessPage,
	)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-1m88i", "Errors.NoChangesFound")
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
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, errors.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	if !existingOIDC.IsOIDC() {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-Ghrh3", "Errors.Project.App.IsNotOIDC")
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
		return errors.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.NotExisting")
	}
	if !app.IsOIDC() {
		return errors.ThrowInvalidArgument(nil, "COMMAND-BHgn2", "Errors.Project.App.IsNotOIDC")
	}
	if app.ClientSecret == nil {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.OIDCConfigInvalid")
	}

	projectAgg := ProjectAggregateFromWriteModel(&app.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
	err = crypto.CompareHash(app.ClientSecret, []byte(secret), c.codeAlg)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		_, err = c.eventstore.Push(ctx, project_repo.NewOIDCConfigSecretCheckSucceededEvent(ctx, projectAgg, app.AppID))
		return err
	}
	_, err = c.eventstore.Push(ctx, project_repo.NewOIDCConfigSecretCheckFailedEvent(ctx, projectAgg, app.AppID))
	logging.OnError(err).Error("could not push event OIDCClientSecretCheckFailed")
	return errors.ThrowInvalidArgument(nil, "COMMAND-Bz542", "Errors.Project.App.ClientSecretInvalid")
}

func (c *Commands) getOIDCAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*OIDCApplicationWriteModel, error) {
	appWriteModel := NewOIDCApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}

func getOIDCAppWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, projectID, appID, resourceOwner string) (*OIDCApplicationWriteModel, error) {
	appWriteModel := NewOIDCApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	events, err := filter(ctx, appWriteModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return appWriteModel, nil
	}
	appWriteModel.AppendEvents(events...)
	err = appWriteModel.Reduce()
	return appWriteModel, err
}
