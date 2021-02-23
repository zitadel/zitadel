package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

func (r *CommandSide) AddAPIApplication(ctx context.Context, application *domain.APIApp, resourceOwner string) (_ *domain.APIApp, err error) {
	project, err := r.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication := NewAPIApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	events, stringPw, err := r.addAPIApplication(ctx, projectAgg, project, application, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication.AppID = application.AppID
	pushedEvents, err := r.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedApplication, pushedEvents...)
	if err != nil {
		return nil, err
	}
	result := apiWriteModelToAPIConfig(addedApplication)
	result.ClientSecretString = stringPw
	return result, nil
}

func (r *CommandSide) addAPIApplication(ctx context.Context, projectAgg *eventstore.Aggregate, proj *domain.Project, apiAppApp *domain.APIApp, resourceOwner string) (events []eventstore.EventPusher, stringPW string, err error) {
	if !apiAppApp.IsValid() {
		return nil, "", caos_errs.ThrowPreconditionFailed(nil, "PROJECT-Bff2g", "Errors.Application.Invalid")
	}
	apiAppApp.AppID, err = r.idGenerator.Next()
	if err != nil {
		return nil, "", err
	}

	events = []eventstore.EventPusher{
		project.NewApplicationAddedEvent(ctx, projectAgg, apiAppApp.AppID, apiAppApp.AppName, resourceOwner),
	}

	var stringPw string
	err = domain.SetNewClientID(apiAppApp, r.idGenerator, proj)
	if err != nil {
		return nil, "", err
	}
	stringPw, err = domain.SetNewClientSecretIfNeeded(apiAppApp, r.applicationSecretGenerator)
	if err != nil {
		return nil, "", err
	}
	events = append(events, project.NewAPIConfigAddedEvent(ctx,
		projectAgg,
		apiAppApp.AppID,
		apiAppApp.ClientID,
		apiAppApp.ClientSecret,
		apiAppApp.AuthMethodType))

	return events, stringPw, nil
}

func (r *CommandSide) ChangeAPIApplication(ctx context.Context, apiApp *domain.APIApp, resourceOwner string) (*domain.APIApp, error) {
	if !apiApp.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1m900", "Errors.Project.App.APIConfigInvalid")
	}

	existingAPI, err := r.getAPIAppWriteModel(ctx, apiApp.AggregateID, apiApp.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State == domain.AppStateUnspecified || existingAPI.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-1m88i", "Errors.NoChangesFound")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAPI, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := apiWriteModelToAPIConfig(existingAPI)
	return result, nil
}

func (r *CommandSide) ChangeAPIApplicationSecret(ctx context.Context, projectID, appID, resourceOwner string) (*domain.APIApp, error) {
	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-99i83", "Errors.IDMissing")
	}

	existingAPI, err := r.getAPIAppWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAPI.State == domain.AppStateUnspecified || existingAPI.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2g66f", "Errors.Project.App.NotExisting")
	}
	cryptoSecret, stringPW, err := domain.NewClientSecret(r.applicationSecretGenerator)
	if err != nil {
		return nil, err
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingAPI.WriteModel)

	pushedEvents, err := r.eventstore.PushEvents(ctx, project.NewAPIConfigSecretChangedEvent(ctx, projectAgg, appID, cryptoSecret))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingAPI, pushedEvents...)
	if err != nil {
		return nil, err
	}

	result := apiWriteModelToAPIConfig(existingAPI)
	result.ClientSecretString = stringPW
	return result, err
}
func (r *CommandSide) getAPIAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*APIApplicationWriteModel, error) {
	appWriteModel := NewAPIApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
