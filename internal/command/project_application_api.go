package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	project_repo "github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type addAPIApp struct {
	AddApp
	AuthMethodType domain.APIAuthMethodType

	ClientID          string
	EncodedHash       string
	ClientSecretPlain string
}

func (c *Commands) AddAPIAppCommand(app *addAPIApp) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if app.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument")
		}
		if app.Name = strings.TrimSpace(app.Name); app.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "PROJE-F7g21", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			project, err := projectWriteModel(ctx, filter, app.Aggregate.ID, app.Aggregate.ResourceOwner)
			if err != nil || !project.State.Valid() {
				return nil, zerrors.ThrowNotFound(err, "PROJE-Sf2gb", "Errors.Project.NotFound")
			}

			app.ClientID, err = id_generator.Next()
			if err != nil {
				return nil, zerrors.ThrowInternal(err, "V2-f0pgP", "Errors.Internal")
			}

			if app.AuthMethodType == domain.APIAuthMethodTypeBasic {
				app.EncodedHash, app.ClientSecretPlain, err = c.newHashedSecret(ctx, filter)
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
				project_repo.NewAPIConfigAddedEvent(
					ctx,
					&app.Aggregate.Aggregate,
					app.ID,
					app.ClientID,
					app.EncodedHash,
					app.AuthMethodType,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) AddAPIApplicationWithID(ctx context.Context, apiApp *domain.APIApp, resourceOwner, appID string) (_ *domain.APIApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	existingAPI, err := c.getAPIAppWriteModel(ctx, apiApp.AggregateID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State != domain.AppStateUnspecified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "PROJECT-mabu12", "Errors.Project.App.AlreadyExisting")
	}
	_, err = c.getProjectByID(ctx, apiApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "PROJECT-9fnsa", "Errors.Project.NotFound")
	}

	return c.addAPIApplicationWithID(ctx, apiApp, resourceOwner, appID)
}

func (c *Commands) AddAPIApplication(ctx context.Context, apiApp *domain.APIApp, resourceOwner string) (_ *domain.APIApp, err error) {
	if apiApp == nil || apiApp.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-5m9E", "Errors.Project.App.Invalid")
	}
	_, err = c.getProjectByID(ctx, apiApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "PROJECT-9fnsf", "Errors.Project.NotFound")
	}

	if !apiApp.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-Bff2g", "Errors.Project.App.Invalid")
	}

	appID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}

	return c.addAPIApplicationWithID(ctx, apiApp, resourceOwner, appID)
}

func (c *Commands) addAPIApplicationWithID(ctx context.Context, apiApp *domain.APIApp, resourceOwner string, appID string) (_ *domain.APIApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	apiApp.AppID = appID

	addedApplication := NewAPIApplicationWriteModel(apiApp.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)

	events := []eventstore.Command{
		project_repo.NewApplicationAddedEvent(ctx, projectAgg, apiApp.AppID, apiApp.AppName),
	}

	var plain string
	err = domain.SetNewClientID(apiApp)
	if err != nil {
		return nil, err
	}
	plain, err = domain.SetNewClientSecretIfNeeded(apiApp, func() (string, string, error) {
		return c.newHashedSecret(ctx, c.eventstore.Filter) //nolint:staticcheck
	})
	if err != nil {
		return nil, err
	}
	events = append(events, project_repo.NewAPIConfigAddedEvent(ctx,
		projectAgg,
		apiApp.AppID,
		apiApp.ClientID,
		apiApp.EncodedHash,
		apiApp.AuthMethodType))

	addedApplication.AppID = apiApp.AppID
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedApplication, pushedEvents...)
	if err != nil {
		return nil, err
	}
	result := apiWriteModelToAPIConfig(addedApplication)
	result.ClientSecretString = plain
	return result, nil
}

func (c *Commands) ChangeAPIApplication(ctx context.Context, apiApp *domain.APIApp, resourceOwner string) (*domain.APIApp, error) {
	if apiApp.AppID == "" || apiApp.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-1m900", "Errors.Project.App.APIConfigInvalid")
	}

	existingAPI, err := c.getAPIAppWriteModel(ctx, apiApp.AggregateID, apiApp.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State == domain.AppStateUnspecified || existingAPI.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	if !existingAPI.IsAPI() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Gnwt3", "Errors.Project.App.IsNotAPI")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingAPI.WriteModel)
	changedEvent, hasChanged, err := existingAPI.NewChangedEvent(
		ctx,
		projectAgg,
		apiApp.AppID,
		apiApp.AuthMethodType)
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
	err = AppendAndReduce(existingAPI, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return apiWriteModelToAPIConfig(existingAPI), nil
}

func (c *Commands) ChangeAPIApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string) (*domain.APIApp, error) {
	if projectID == "" || appID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingAPI, err := c.getAPIAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State == domain.AppStateUnspecified || existingAPI.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	if !existingAPI.IsAPI() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-aeH4", "Errors.Project.App.IsNotAPI")
	}
	encodedHash, plain, err := c.newHashedSecret(ctx, c.eventstore.Filter) //nolint:staticcheck
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingAPI.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, project_repo.NewAPIConfigSecretChangedEvent(ctx, projectAgg, appID, encodedHash))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAPI, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := apiWriteModelToAPIConfig(existingAPI)
	result.ClientSecretString = plain
	return result, err
}

func (c *Commands) VerifyAPIClientSecret(ctx context.Context, projectID, appID, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	app, err := c.getAPIAppWriteModel(ctx, projectID, appID, "")
	if err != nil {
		return err
	}
	if !app.State.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-DFnbf", "Errors.Project.App.NotExisting")
	}
	if !app.IsAPI() {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Bf3fw", "Errors.Project.App.IsNotAPI")
	}
	if app.HashedSecret == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-D3t5g", "Errors.Project.App.APIConfigInvalid")
	}

	projectAgg := ProjectAggregateFromWriteModel(&app.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.secretHasher.Verify(app.HashedSecret, secret)
	spanPasswordComparison.EndWithError(err)
	if err == nil {
		c.apiSecretCheckSucceeded(ctx, projectAgg, app.AppID, updated)
		return err
	}
	c.apiSecretCheckFailed(ctx, projectAgg, app.AppID)
	return zerrors.ThrowInvalidArgument(err, "COMMAND-SADfg", "Errors.Project.App.ClientSecretInvalid")
}

func (c *Commands) APISecretCheckSucceeded(ctx context.Context, appID, projectID, resourceOwner, updated string) {
	agg := project_repo.NewAggregate(projectID, resourceOwner)
	c.apiSecretCheckSucceeded(ctx, &agg.Aggregate, appID, updated)
}

func (c *Commands) APISecretCheckFailed(ctx context.Context, appID, projectID, resourceOwner string) {
	agg := project_repo.NewAggregate(projectID, resourceOwner)
	c.apiSecretCheckFailed(ctx, &agg.Aggregate, appID)
}

func (c *Commands) getAPIAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (_ *APIApplicationWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	appWriteModel := NewAPIApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}

func (c *Commands) apiSecretCheckSucceeded(ctx context.Context, agg *eventstore.Aggregate, appID, updated string) {
	cmds := append(
		make([]eventstore.Command, 0, 2),
		project_repo.NewAPIConfigSecretCheckSucceededEvent(ctx, agg, appID),
	)
	if updated != "" {
		cmds = append(cmds, project_repo.NewAPIConfigSecretHashUpdatedEvent(ctx, agg, appID, updated))
	}
	c.asyncPush(ctx, cmds...)
}

func (c *Commands) apiSecretCheckFailed(ctx context.Context, agg *eventstore.Aggregate, appID string) {
	c.asyncPush(ctx, project_repo.NewAPIConfigSecretCheckFailedEvent(ctx, agg, appID))
}
