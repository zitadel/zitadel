package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
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

			app.ClientID, err = c.idGenerator.Next()
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

	if _, err := c.checkProjectExists(ctx, apiApp.AggregateID, resourceOwner); err != nil {
		return nil, err
	}
	return c.addAPIApplicationWithID(ctx, apiApp, resourceOwner, appID)
}

func (c *Commands) AddAPIApplication(ctx context.Context, apiApp *domain.APIApp, resourceOwner string) (_ *domain.APIApp, err error) {
	if apiApp == nil || apiApp.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-5m9E", "Errors.Project.App.Invalid")
	}

	projectResOwner, err := c.checkProjectExists(ctx, apiApp.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if resourceOwner == "" {
		resourceOwner = projectResOwner
	}

	if !apiApp.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-Bff2g", "Errors.Project.App.Invalid")
	}

	appID := apiApp.AppID
	if appID == "" {
		appID, err = c.idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}

	return c.addAPIApplicationWithID(ctx, apiApp, resourceOwner, appID)
}

func (c *Commands) addAPIApplicationWithID(ctx context.Context, apiApp *domain.APIApp, resourceOwner string, appID string) (_ *domain.APIApp, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	apiApp.AppID = appID

	addedApplication := NewAPIApplicationWriteModel(apiApp.AggregateID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, addedApplication); err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateApplication(ctx, addedApplication.ResourceOwner, addedApplication.AggregateID); err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)

	events := []eventstore.Command{
		project_repo.NewApplicationAddedEvent(ctx, projectAgg, apiApp.AppID, apiApp.AppName),
	}

	var plain string
	err = domain.SetNewClientID(apiApp, c.idGenerator)
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

func (c *Commands) UpdateAPIApplication(ctx context.Context, apiApp *domain.APIApp, resourceOwner string) (*domain.APIApp, error) {
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
	if err := c.eventstore.FilterToQueryReducer(ctx, existingAPI); err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateApplication(ctx, existingAPI.ResourceOwner, existingAPI.AggregateID); err != nil {
		return nil, err
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

	if err := c.checkPermissionUpdateApplication(ctx, existingAPI.ResourceOwner, existingAPI.AggregateID); err != nil {
		return nil, err
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

func (c *Commands) APIUpdateSecret(ctx context.Context, appID, projectID, resourceOwner, updated string) {
	agg := project_repo.NewAggregate(projectID, resourceOwner)
	c.apiUpdateSecret(ctx, &agg.Aggregate, appID, updated)
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

func (c *Commands) apiUpdateSecret(ctx context.Context, agg *eventstore.Aggregate, appID, updated string) {
	c.asyncPush(ctx, project_repo.NewAPIConfigSecretHashUpdatedEvent(ctx, agg, appID, updated))
}
