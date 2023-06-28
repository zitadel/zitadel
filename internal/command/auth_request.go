package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddAuthRequest(ctx context.Context, request *domain.AuthRequest) error {
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	oidcRequest := request.Request.(*domain.AuthRequestOIDC)
	writeModel, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateUnspecified {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Sf3gt", "Errors.AuthRequest.AlreadyExisting")
	}
	return c.pushAppendAndReduce(ctx, writeModel, authrequest.NewAddedEvent(
		ctx,
		&authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		request.LoginClient,
		request.ApplicationID,
		request.CallbackURI,
		request.TransferState,
		oidcRequest.Nonce,
		oidcRequest.Scopes,
		oidcRequest.ResponseType,
		oidcRequest.CodeChallenge,
		request.Prompt,
		request.UiLocales,
		request.MaxAuthAge,
		request.LoginHint,
		request.UserID))
}

func (c *Commands) getAuthRequestWriteModel(ctx context.Context, id string) (writeModel *AuthRequestWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewAuthRequestWriteModel(id)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
