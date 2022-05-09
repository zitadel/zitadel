package command

import (
	"context"
	"fmt"
	"github.com/zitadel/saml/pkg/provider/xml"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func (c *Commands) AddSAMLApplication(ctx context.Context, application *domain.SAMLApp, resourceOwner string) (_ *domain.SAMLApp, err error) {
	if application == nil || application.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-35Fn0", "Errors.Application.Invalid")
	}

	addedApplication := NewSAMLApplicationWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	events, err := c.addSAMLApplication(ctx, projectAgg, application, resourceOwner)
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

func (c *Commands) addSAMLApplication(ctx context.Context, projectAgg *eventstore.Aggregate, samlApp *domain.SAMLApp, resourceOwner string) (events []eventstore.Command, err error) {
	if samlApp.AppName == "" || !samlApp.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "PROJECT-1n9df", "Errors.Application.Invalid")
	}
	samlApp.AppID, err = c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	events = []eventstore.Command{
		project.NewApplicationAddedEvent(ctx, projectAgg, samlApp.AppID, samlApp.AppName),
	}

	if samlApp.Metadata == "" && samlApp.MetadataURL == "" {
		return nil, fmt.Errorf("no metadata provided")
	}

	var metadata []byte
	if samlApp.MetadataURL != "" {
		data, err := xml.ReadMetadataFromURL(samlApp.MetadataURL)
		if err != nil {
			return nil, err
		}
		metadata = data
	} else {
		metadata = []byte(samlApp.Metadata)
	}

	entity, err := xml.ParseMetadataXmlIntoStruct(metadata)
	if err != nil {
		return nil, err
	}

	events = append(events, project.NewSAMLConfigAddedEvent(ctx,
		projectAgg,
		samlApp.AppID,
		string(entity.EntityID),
		samlApp.Metadata,
		samlApp.MetadataURL,
	))

	return events, nil
}

func (c *Commands) ChangeSAMLApplication(ctx context.Context, saml *domain.SAMLApp, resourceOwner string) (*domain.SAMLApp, error) {
	if !saml.IsValid() || saml.AppID == "" || saml.AggregateID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5n9fs", "Errors.Project.App.SAMLConfigInvalid")
	}

	existingSAML, err := c.getSAMLAppWriteModel(ctx, saml.AggregateID, saml.AppID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingSAML.State == domain.AppStateUnspecified || existingSAML.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2n8uU", "Errors.Project.App.NotExisting")
	}
	if !existingSAML.IsSAML() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-GBr34", "Errors.Project.App.IsNotOIDC")
	}
	projectAgg := ProjectAggregateFromWriteModel(&existingSAML.WriteModel)

	var metadata []byte
	if saml.MetadataURL != "" {
		data, err := xml.ReadMetadataFromURL(saml.MetadataURL)
		if err != nil {
			return nil, err
		}
		metadata = data
	}

	entity, err := xml.ParseMetadataXmlIntoStruct(metadata)
	if err != nil {
		return nil, err
	}

	changedEvent, hasChanged, err := existingSAML.NewChangedEvent(
		ctx,
		projectAgg,
		saml.AppID,
		string(entity.EntityID),
		saml.Metadata,
		saml.MetadataURL)
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
