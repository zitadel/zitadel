package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	project_repo "github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeApplicationSecret(ctx context.Context, projectID, applicationID, resourceOwner string) (secret string, changeDate time.Time, err error) {
	if projectID == "" || applicationID == "" {
		return "", time.Time{}, zerrors.ThrowInvalidArgument(nil, "COMMAND-KJ29c", "Errors.IDMissing")
	}

	existingApplication, err := c.getApplicationSecretWriteModel(ctx, projectID, applicationID, resourceOwner)
	if err != nil {
		return "", time.Time{}, err
	}
	if !existingApplication.State.Exists() {
		return "", time.Time{}, zerrors.ThrowNotFound(nil, "COMMAND-Kd92s", "Errors.Project.App.NotExisting")
	}

	if err := c.checkPermissionUpdateApplication(ctx, existingApplication.ResourceOwner, existingApplication.AggregateID); err != nil {
		return "", time.Time{}, err
	}

	encodedHash, plain, err := c.newHashedSecret(ctx, c.eventstore.Filter) //nolint:staticcheck
	if err != nil {
		return "", time.Time{}, err
	}

	projectAgg := ProjectAggregateFromWriteModelWithCTX(ctx, &existingApplication.WriteModel)

	var command eventstore.Command
	command = project_repo.NewOIDCConfigSecretChangedEvent(ctx, projectAgg, applicationID, encodedHash)
	if existingApplication.IsAPI {
		command = project_repo.NewAPIConfigSecretChangedEvent(ctx, projectAgg, applicationID, encodedHash)
	}
	if err = c.pushAppendAndReduce(ctx, existingApplication, command); err != nil {
		return "", time.Time{}, err
	}

	return plain, existingApplication.ChangeDate, nil
}

func (c *Commands) getApplicationSecretWriteModel(ctx context.Context, projectID, applicationID, resourceOwner string) (_ *ApplicationSecretWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	appWriteModel := NewApplicationSecretWriteModel(projectID, applicationID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, appWriteModel)
	if err != nil {
		return nil, err
	}
	return appWriteModel, nil
}
