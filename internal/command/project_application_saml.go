package command

import (
	"context"

	"github.com/zitadel/saml/pkg/provider/xml"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddSAMLApplication(ctx context.Context, application *domain.SAMLApp, resourceOwner string) (_ *domain.SAMLApp, err error) {
	if application == nil || application.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-35Fn0", "Errors.Project.App.Invalid")
	}

	_, err = c.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "PROJECT-3p9ss", "Errors.Project.NotFound")
	}

	addedApplication := NewSAMLApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	events, err := c.addSAMLApplication(ctx, projectAgg, application)
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
	result := samlWriteModelToSAMLConfig(addedApplication)
	return result, nil
}

func (c *Commands) addSAMLApplication(ctx context.Context, projectAgg *eventstore.Aggregate, samlApp *domain.SAMLApp) (events []eventstore.Command, err error) {

	if samlApp.AppName == "" || !samlApp.IsValid() {
		return nil, zerrors.ThrowInvalidArgument(nil, "PROJECT-1n9df", "Errors.Project.App.Invalid")
	}

	if samlApp.Metadata == nil && samlApp.MetadataURL == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "SAML-podix9", "Errors.Project.App.SAMLMetadataMissing")
	}

	if samlApp.MetadataURL != "" {
		data, err := xml.ReadMetadataFromURL(c.httpClient, samlApp.MetadataURL)
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "SAML-wmqlo1", "Errors.Project.App.SAMLMetadataMissing")
		}
		samlApp.Metadata = data
	}

	entity, err := xml.ParseMetadataXmlIntoStruct(samlApp.Metadata)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "SAML-bquso", "Errors.Project.App.SAMLMetadataFormat")
	}

	samlApp.AppID, err = id_generator.Next()
	if err != nil {
		return nil, err
	}

	return []eventstore.Command{
		project.NewApplicationAddedEvent(ctx, projectAgg, samlApp.AppID, samlApp.AppName),
		project.NewSAMLConfigAddedEvent(ctx,
			projectAgg,
			samlApp.AppID,
			string(entity.EntityID),
			samlApp.Metadata,
			samlApp.MetadataURL,
		),
	}, nil
}

func (c *Commands) ChangeSAMLApplication(ctx context.Context, samlApp *domain.SAMLApp, resourceOwner string) (*domain.SAMLApp, error) {
	if !samlApp.IsValid() || samlApp.AppID == "" || samlApp.AggregateID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5n9fs", "Errors.Project.App.SAMLConfigInvalid")
	}

	existingSAML, err := c.getSAMLAppWriteModel(ctx, samlApp.AggregateID, samlApp.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingSAML.State == domain.AppStateUnspecified || existingSAML.State == domain.AppStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	if !existingSAML.IsSAML() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-GBr35", "Errors.Project.App.IsNotSAML")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingSAML.WriteModel)

	if samlApp.MetadataURL != "" {
		data, err := xml.ReadMetadataFromURL(c.httpClient, samlApp.MetadataURL)
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "SAML-J3kg3", "Errors.Project.App.SAMLMetadataMissing")
		}
		samlApp.Metadata = data
	}

	entity, err := xml.ParseMetadataXmlIntoStruct(samlApp.Metadata)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "SAML-3fk2b", "Errors.Project.App.SAMLMetadataFormat")
	}

	changedEvent, hasChanged, err := existingSAML.NewChangedEvent(
		ctx,
		projectAgg,
		samlApp.AppID,
		string(entity.EntityID),
		samlApp.Metadata,
		samlApp.MetadataURL)
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
	err = AppendAndReduce(existingSAML, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return samlWriteModelToSAMLConfig(existingSAML), nil
}

func (c *Commands) getSAMLAppWriteModel(ctx context.Context, projectID, appID, resourceOwner string) (*SAMLApplicationWriteModel, error) {
	appWriteModel := NewSAMLApplicationWriteModelWithAppID(projectID, appID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
