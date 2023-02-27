package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func (c *Commands) AddApplicationIcon(ctx context.Context, projectID string, appID string, dark bool, upload *AssetUpload) (*domain.ObjectDetails, error) {
	var pushedEvents []eventstore.Event

	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-88fi0", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, upload.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}

	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "COMMAND-1Xyud", "Errors.Assets.Object.PutFailed")
	}

	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	if dark {
		pushedEvents, err = c.eventstore.Push(ctx, project.NewApplicationDarkIconAddedEvent(ctx, projectAgg, asset.Name))
		if err != nil {
			return nil, err
		}
	} else {
		pushedEvents, err = c.eventstore.Push(ctx, project.NewApplicationLightIconAddedEvent(ctx, projectAgg, asset.Name))
		if err != nil {
			return nil, err
		}
	}

	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}

func (c *Commands) RemoveApplicationIcon(ctx context.Context, projectID, appID string, resourceOwner string, dark bool) (*domain.ObjectDetails, error) {
	var pushedEvents []eventstore.Event

	if projectID == "" || appID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-88fi0", "Errors.IDMissing")
	}

	existingApp, err := c.getApplicationWriteModel(ctx, projectID, appID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingApp.State == domain.AppStateUnspecified || existingApp.State == domain.AppStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-ov9d3", "Errors.Project.App.NotExisting")
	}

	orgID := authz.GetCtxData(ctx).OrgID
	if dark {
		err = c.removeAsset(ctx, orgID, existingApp.DarkIconURL)
		if err != nil {
			return nil, err
		}
	} else {
		err = c.removeAsset(ctx, orgID, existingApp.LightIconURL)
		if err != nil {
			return nil, err
		}
	}
	

	projectAgg := ProjectAggregateFromWriteModel(&existingApp.WriteModel)
	if dark {
		pushedEvents, err = c.eventstore.Push(ctx, project.NewApplicationDarkIconRemovedEvent(ctx, projectAgg, existingApp.LightIconURL))
		if err != nil {
			return nil, err
		}
	} else {
		pushedEvents, err = c.eventstore.Push(ctx, project.NewApplicationLightIconRemovedEvent(ctx, projectAgg, existingApp.DarkIconURL))
		if err != nil {
			return nil, err
		}
	}

	err = AppendAndReduce(existingApp, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingApp.WriteModel), nil
}
