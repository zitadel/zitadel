package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AuthRequest struct {
	ID            string
	LoginClient   string
	ClientID      string
	RedirectURI   string
	State         string
	Nonce         string
	Scope         []string
	Audience      []string
	ResponseType  domain.OIDCResponseType
	CodeChallenge *domain.OIDCCodeChallenge
	Prompt        []domain.Prompt
	UILocales     []string
	MaxAge        *time.Duration
	LoginHint     string
	HintUserID    string
}

func (c *Commands) AddAuthRequest(ctx context.Context, authRequest *AuthRequest) (err error) {
	authRequest.ID, err = c.idGenerator.Next()
	if err != nil {
		return err
	}
	writeModel, err := c.getAuthRequestWriteModel(ctx, authRequest.ID)
	if err != nil {
		return err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateUnspecified {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Sf3gt", "Errors.AuthRequest.AlreadyExisting")
	}
	return c.pushAppendAndReduce(ctx, writeModel, authrequest.NewAddedEvent(
		ctx,
		&authrequest.NewAggregate(authRequest.ID, authz.GetInstance(ctx).InstanceID()).Aggregate,
		authRequest.LoginClient,
		authRequest.ClientID,
		authRequest.RedirectURI,
		authRequest.State,
		authRequest.Nonce,
		authRequest.Scope,
		authRequest.Audience,
		authRequest.ResponseType,
		authRequest.CodeChallenge,
		authRequest.Prompt,
		authRequest.UILocales,
		authRequest.MaxAge,
		authRequest.LoginHint,
		authRequest.HintUserID,
	))
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

func (c *Commands) LinkSessionToAuthRequest(ctx context.Context, id, sessionID, sessionToken string) error {
	writeModel, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled")
	}
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return err
	}
	if sessionWriteModel.State == domain.SessionStateUnspecified {
		return errors.ThrowNotFound(nil, "COMMAND-x0099887", "Errors.Session.NotExisting")
	}
	if err := c.sessionPermission(ctx, sessionWriteModel, sessionToken, domain.PermissionSessionWrite); err != nil {
		return err
	}
	return c.pushAppendAndReduce(ctx, writeModel, authrequest.NewSessionLinkedEvent(
		ctx,
		&authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		sessionID,
	))
}

func (c *Commands) FailAuthRequest(ctx context.Context, id string) error {
	writeModel, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Sx202nt", "Errors.AuthRequest.AlreadyHandled")
	}
	return c.pushAppendAndReduce(ctx, writeModel, authrequest.NewFailedEvent(
		ctx,
		&authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
	))
}
