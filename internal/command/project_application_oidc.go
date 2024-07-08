package command

import (
	"context"
	"strings"
	"time"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	project_repo "github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	ClientSecret      string
	ClientSecretPlain string
}

// AddOIDCAppCommand prepares the commands to add an oidc app. The ClientID will be set during the CreateCommands
func (c *Commands) AddOIDCAppCommand(app *addOIDCApp) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if app.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument")
		}

		if app.Name = strings.TrimSpace(app.Name); app.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument")
		}

		if app.ClockSkew > time.Second*5 || app.ClockSkew < 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "V2-PnCMS", "Errors.Invalid.Argument")
		}

		for _, origin := range app.AdditionalOrigins {
			if !http_util.IsOrigin(strings.TrimSpace(origin)) {
				return nil, zerrors.ThrowInvalidArgument(nil, "V2-DqWPX", "Errors.Invalid.Argument")
			}
		}

		if !domain.ContainsRequiredGrantTypes(app.ResponseTypes, app.GrantTypes) {
			return nil, zerrors.ThrowInvalidArgument(nil, "V2-sLpW1", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (_ []eventstore.Command, err error) {
			project, err := projectWriteModel(ctx, filter, app.Aggregate.ID, app.Aggregate.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, zerrors.ThrowNotFound(err, "PROJE-6swVG", "Errors.Project.NotFound")
			}

			app.ClientID, err = id_generator.Next()
			if err != nil {
				return nil, zerrors.ThrowInternal(err, "V2-VMSQ1", "Errors.Internal")
			}

			if app.AuthMethodType == domain.OIDCAuthMethodTypeBasic || app.AuthMethodType == domain.OIDCAuthMethodTypePost {
				app.ClientSecret, app.ClientSecretPlain, err = c.newHashedSecret(ctx, filter)
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
					trimStringSliceWhiteSpaces(app.RedirectUris),
					app.ResponseTypes,
					app.GrantTypes,
					app.ApplicationType,
					app.AuthMethodType,
					trimStringSliceWhiteSpaces(app.PostLogoutRedirectUris),
					app.DevMode,
					app.AccessTokenType,
					app.AccessTokenRoleAssertion,
					app.IDTokenRoleAssertion,
					app.IDTokenUserinfoAssertion,
					app.ClockSkew,
					trimStringSliceWhiteSpaces(app.AdditionalOrigins),
					app.SkipSuccessPageForNativeApp,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) AddOIDCApplicationWithID(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner, appID string) (_ *domain.OIDCApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingApp, err := c.getOIDCAppWriteModel(ctx, oidcApp.AggregateID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State != domain.AppStateUnspecified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-lxowmp", "Errors.Project.App.AlreadyExisting")
	}

	_, err = c.getProjectByID(ctx, oidcApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "PROJECT-3m9s2", "Errors.Project.NotFound")
	}

	return c.addOIDCApplicationWithID(ctx, oidcApp, resourceOwner, appID)
}

func (c *Commands) AddOIDCApplication(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner string) (_ *domain.OIDCApp, err error) {
	if oidcApp == nil || oidcApp.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-34Fm0", "Errors.Project.App.Invalid")
	}
	_, err = c.getProjectByID(ctx, oidcApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "PROJECT-3m9ss", "Errors.Project.NotFound")
	}

	if oidcApp.AppName == "" || !oidcApp.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-1n8df", "Errors.Project.App.Invalid")
	}

	appID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.addOIDCApplicationWithID(ctx, oidcApp, resourceOwner, appID)
}

func (c *Commands) addOIDCApplicationWithID(ctx context.Context, oidcApp *domain.OIDCApp, resourceOwner string, appID string) (_ *domain.OIDCApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	addedApplication := NewOIDCApplicationWriteModel(oidcApp.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)

	oidcApp.AppID = appID

	events := []eventstore.Command{
		project_repo.NewApplicationAddedEvent(ctx, projectAgg, oidcApp.AppID, oidcApp.AppName),
	}

	var plain string
	err = domain.SetNewClientID(oidcApp)
	if err != nil {
		return nil, err
	}
	plain, err = domain.SetNewClientSecretIfNeeded(oidcApp, func() (string, string, error) {
		return c.newHashedSecret(ctx, c.eventstore.Filter) //nolint:staticcheck
	})
	if err != nil {
		return nil, err
	}
	events = append(events, project_repo.NewOIDCConfigAddedEvent(ctx,
		projectAgg,
		oidcApp.OIDCVersion,
		oidcApp.AppID,
		oidcApp.ClientID,
		oidcApp.EncodedHash,
		trimStringSliceWhiteSpaces(oidcApp.RedirectUris),
		oidcApp.ResponseTypes,
		oidcApp.GrantTypes,
		oidcApp.ApplicationType,
		oidcApp.AuthMethodType,
		trimStringSliceWhiteSpaces(oidcApp.PostLogoutRedirectUris),
		oidcApp.DevMode,
		oidcApp.AccessTokenType,
		oidcApp.AccessTokenRoleAssertion,
		oidcApp.IDTokenRoleAssertion,
		oidcApp.IDTokenUserinfoAssertion,
		oidcApp.ClockSkew,
		trimStringSliceWhiteSpaces(oidcApp.AdditionalOrigins),
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
	result.ClientSecretString = plain
	result.FillCompliance()
	return result, nil
}

func (c *Commands) ChangeOIDCApplication(ctx context.Context, oidc *domain.OIDCApp, resourceOwner string) (*domain.OIDCApp, error) {
	if !oidc.IsValid() || oidc.AppID == "" || oidc.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5m9fs", "Errors.Project.App.OIDCConfigInvalid")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, oidc.AggregateID, oidc.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	if !existingOIDC.IsOIDC() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-GBr34", "Errors.Project.App.IsNotOIDC")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)
	changedEvent, hasChanged, err := existingOIDC.NewChangedEvent(
		ctx,
		projectAgg,
		oidc.AppID,
		trimStringSliceWhiteSpaces(oidc.RedirectUris),
		trimStringSliceWhiteSpaces(oidc.PostLogoutRedirectUris),
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
		trimStringSliceWhiteSpaces(oidc.AdditionalOrigins),
		oidc.SkipNativeAppSuccessPage,
	)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-1m88i", "Errors.NoChangesFound")
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

func (c *Commands) ChangeOIDCApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string) (*domain.OIDCApp, error) {
	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingOIDC, err := c.getOIDCAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingOIDC.State == domain.AppStateUnspecified || existingOIDC.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	if !existingOIDC.IsOIDC() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Ghrh3", "Errors.Project.App.IsNotOIDC")
	}
	encodedHash, plain, err := c.newHashedSecret(ctx, c.eventstore.Filter) //nolint:staticcheck
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingOIDC.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, project_repo.NewOIDCConfigSecretChangedEvent(ctx, projectAgg, appID, encodedHash))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingOIDC, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := oidcWriteModelToOIDCConfig(existingOIDC)
	result.ClientSecretString = plain
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
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-D8hba", "Errors.Project.App.NotExisting")
	}
	if !app.IsOIDC() {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-BHgn2", "Errors.Project.App.IsNotOIDC")
	}
	if app.HashedSecret == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-D6hba", "Errors.Project.App.OIDCConfigInvalid")
	}

	projectAgg := ProjectAggregateFromWriteModel(&app.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.secretHasher.Verify(app.HashedSecret, secret)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		c.oidcSecretCheckSucceeded(ctx, projectAgg, appID, updated)
		return nil
	}
	c.oidcSecretCheckFailed(ctx, projectAgg, appID)
	return zerrors.ThrowInvalidArgument(err, "COMMAND-Bz542", "Errors.Project.App.ClientSecretInvalid")
}

func (c *Commands) OIDCSecretCheckSucceeded(ctx context.Context, appID, projectID, resourceOwner, updated string) {
	agg := project_repo.NewAggregate(projectID, resourceOwner)
	c.oidcSecretCheckSucceeded(ctx, &agg.Aggregate, appID, updated)
}

func (c *Commands) OIDCSecretCheckFailed(ctx context.Context, appID, projectID, resourceOwner string) {
	agg := project_repo.NewAggregate(projectID, resourceOwner)
	c.oidcSecretCheckFailed(ctx, &agg.Aggregate, appID)
}

func (c *Commands) getOIDCAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (_ *OIDCApplicationWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	appWriteModel := NewOIDCApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
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

func trimStringSliceWhiteSpaces(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}

func (c *Commands) oidcSecretCheckSucceeded(ctx context.Context, agg *eventstore.Aggregate, appID, updated string) {
	cmds := append(
		make([]eventstore.Command, 0, 2),
		project_repo.NewOIDCConfigSecretCheckSucceededEvent(ctx, agg, appID),
	)
	if updated != "" {
		cmds = append(cmds, project_repo.NewOIDCConfigSecretHashUpdatedEvent(ctx, agg, appID, updated))
	}
	c.asyncPush(ctx, cmds...)
}

func (c *Commands) oidcSecretCheckFailed(ctx context.Context, agg *eventstore.Aggregate, appID string) {
	c.asyncPush(ctx, project_repo.NewOIDCConfigSecretCheckFailedEvent(ctx, agg, appID))
}
