package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc/amr"
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
	LoginHint     *string
	HintUserID    *string
}

type AuthenticatedAuthRequest struct {
	*AuthRequest
	SessionID string
	UserID    string
	AMR       []string
	AuthTime  time.Time
}

const IDPrefixV2 = "V2_"

func (c *Commands) AddAuthRequest(ctx context.Context, authRequest *AuthRequest) (err error) {
	authRequestID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	authRequest.ID = IDPrefixV2 + authRequestID
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

func (c *Commands) LinkSessionToAuthRequest(ctx context.Context, id, sessionID, sessionToken string, checkLoginClient bool) (*domain.ObjectDetails, error) {
	writeModel, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if checkLoginClient && authz.GetCtxData(ctx).UserID != writeModel.LoginClient {
		return nil, errors.ThrowPermissionDenied(nil, "COMMAND-rai9Y", "Errors.AuthRequest.WrongLoginClient")
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled")
	}
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetCtxData(ctx).OrgID)
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	if sessionWriteModel.State == domain.SessionStateUnspecified {
		return nil, errors.ThrowNotFound(nil, "COMMAND-x0099887", "Errors.Session.NotExisting")
	}
	if err := c.sessionPermission(ctx, sessionWriteModel, sessionToken, domain.PermissionSessionWrite); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, writeModel, authrequest.NewSessionLinkedEvent(
		ctx, &authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		sessionID,
		sessionWriteModel.UserID,
		sessionWriteModel.AuthenticationTime(),
		amr.List(sessionWriteModel),
	)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
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

func (c *Commands) AddAuthRequestCode(ctx context.Context, authRequestID, code string) (err error) {
	if code == "" {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-Ht52d", "Errors.AuthRequest.InvalidCode")
	}
	writeModel, err := c.getAuthRequestWriteModel(ctx, authRequestID)
	if err != nil {
		return err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded || writeModel.SessionID == "" {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.AlreadyHandled") //TODO: key
	}
	return c.pushAppendAndReduce(ctx, writeModel, authrequest.NewCodeAddedEvent(ctx,
		&authrequest.NewAggregate(writeModel.AggregateID, authz.GetInstance(ctx).InstanceID()).Aggregate,
		code,
	))
}

func (c *Commands) ExchangeAuthCode(ctx context.Context, code string) (authRequest *AuthenticatedAuthRequest, err error) {
	//split := strings.Split(code, ":")
	//if len(split) != 2 {
	//	return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-Sfr3s", "Errors.AuthRequest.InvalidCode")
	//}
	writeModel, err := c.getAuthRequestWriteModel(ctx, code)
	if err != nil {
		return nil, err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateCodeAdded {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.NoCode")
	}
	//if writeModel.ExchangeCode != split[1] {
	//	return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-DBNqz", "Errors.AuthRequest.InvalidCode")
	//}
	err = c.pushAppendAndReduce(ctx, writeModel, authrequest.NewCodeExchangedEvent(ctx,
		&authrequest.NewAggregate(writeModel.AggregateID, authz.GetInstance(ctx).InstanceID()).Aggregate))
	if err != nil {
		return nil, err
	}
	return authRequestWriteModelToAuthenticatedAuthRequest(writeModel), nil
}

func authRequestWriteModelToAuthenticatedAuthRequest(writeModel *AuthRequestWriteModel) (_ *AuthenticatedAuthRequest) {
	return &AuthenticatedAuthRequest{
		AuthRequest: &AuthRequest{
			ID:            writeModel.AggregateID,
			LoginClient:   writeModel.LoginClient,
			ClientID:      writeModel.ClientID,
			RedirectURI:   writeModel.RedirectURI,
			State:         writeModel.State,
			Nonce:         writeModel.Nonce,
			Scope:         writeModel.Scope,
			Audience:      writeModel.Audience,
			ResponseType:  writeModel.ResponseType,
			CodeChallenge: writeModel.CodeChallenge,
			Prompt:        writeModel.Prompt,
			UILocales:     writeModel.UILocales,
			MaxAge:        writeModel.MaxAge,
			LoginHint:     writeModel.LoginHint,
			HintUserID:    writeModel.HintUserID,
		},
		SessionID: writeModel.SessionID,
		UserID:    writeModel.UserID,
		AMR:       writeModel.AMR,
		AuthTime:  writeModel.AuthTime,
	}
}

func (c *Commands) GetAuthRequestWriteModel(ctx context.Context, id string) (writeModel *AuthRequestWriteModel, err error) {
	return c.getAuthRequestWriteModel(ctx, id)
}

func (c *Commands) getAuthRequestWriteModel(ctx context.Context, id string) (writeModel *AuthRequestWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewAuthRequestWriteModel(ctx, id)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
