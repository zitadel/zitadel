package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AuthRequest struct {
	ID               string
	LoginClient      string
	ClientID         string
	RedirectURI      string
	State            string
	Nonce            string
	Scope            []string
	Audience         []string
	ResponseType     domain.OIDCResponseType
	ResponseMode     domain.OIDCResponseMode
	CodeChallenge    *domain.OIDCCodeChallenge
	Prompt           []domain.Prompt
	UILocales        []string
	MaxAge           *time.Duration
	LoginHint        *string
	HintUserID       *string
	NeedRefreshToken bool
}

type CurrentAuthRequest struct {
	*AuthRequest
	SessionID   string
	UserID      string
	AuthMethods []domain.UserAuthMethodType
	AuthTime    time.Time
}

const IDPrefixV2 = "V2_"

func (c *Commands) AddAuthRequest(ctx context.Context, authRequest *AuthRequest) (_ *CurrentAuthRequest, err error) {
	authRequestID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}
	authRequest.ID = IDPrefixV2 + authRequestID
	writeModel, err := c.getAuthRequestWriteModel(ctx, authRequest.ID)
	if err != nil {
		return nil, err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateUnspecified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sf3gt", "Errors.AuthRequest.AlreadyExisting")
	}
	err = c.pushAppendAndReduce(ctx, writeModel, authrequest.NewAddedEvent(
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
		authRequest.ResponseMode,
		authRequest.CodeChallenge,
		authRequest.Prompt,
		authRequest.UILocales,
		authRequest.MaxAge,
		authRequest.LoginHint,
		authRequest.HintUserID,
		authRequest.NeedRefreshToken,
	))
	if err != nil {
		return nil, err
	}
	return authRequestWriteModelToCurrentAuthRequest(writeModel), nil
}

func (c *Commands) LinkSessionToAuthRequest(ctx context.Context, id, sessionID, sessionToken string, checkLoginClient bool) (*domain.ObjectDetails, *CurrentAuthRequest, error) {
	writeModel, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if writeModel.AuthRequestState == domain.AuthRequestStateUnspecified {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-jae5P", "Errors.AuthRequest.NotExisting")
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled")
	}
	if checkLoginClient && authz.GetCtxData(ctx).UserID != writeModel.LoginClient {
		return nil, nil, zerrors.ThrowPermissionDenied(nil, "COMMAND-rai9Y", "Errors.AuthRequest.WrongLoginClient")
	}
	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetInstance(ctx).InstanceID())
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, nil, err
	}
	if err = sessionWriteModel.CheckIsActive(); err != nil {
		return nil, nil, err
	}
	if err := c.sessionTokenVerifier(ctx, sessionToken, sessionWriteModel.AggregateID, sessionWriteModel.TokenID); err != nil {
		return nil, nil, err
	}

	if err := c.pushAppendAndReduce(ctx, writeModel, authrequest.NewSessionLinkedEvent(
		ctx, &authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		sessionID,
		sessionWriteModel.UserID,
		sessionWriteModel.AuthenticationTime(),
		sessionWriteModel.AuthMethodTypes(),
	)); err != nil {
		return nil, nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), authRequestWriteModelToCurrentAuthRequest(writeModel), nil
}

func (c *Commands) FailAuthRequest(ctx context.Context, id string, reason domain.OIDCErrorReason) (*domain.ObjectDetails, *CurrentAuthRequest, error) {
	writeModel, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sx202nt", "Errors.AuthRequest.AlreadyHandled")
	}
	err = c.pushAppendAndReduce(ctx, writeModel, authrequest.NewFailedEvent(
		ctx,
		&authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		reason,
	))
	if err != nil {
		return nil, nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), authRequestWriteModelToCurrentAuthRequest(writeModel), nil
}

func (c *Commands) AddAuthRequestCode(ctx context.Context, authRequestID, code string) (err error) {
	if code == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Ht52d", "Errors.AuthRequest.InvalidCode")
	}
	writeModel, err := c.getAuthRequestWriteModel(ctx, authRequestID)
	if err != nil {
		return err
	}
	if writeModel.AuthRequestState != domain.AuthRequestStateAdded || writeModel.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.AlreadyHandled")
	}
	return c.pushAppendAndReduce(ctx, writeModel, authrequest.NewCodeAddedEvent(ctx,
		&authrequest.NewAggregate(writeModel.AggregateID, authz.GetInstance(ctx).InstanceID()).Aggregate))
}

func authRequestWriteModelToCurrentAuthRequest(writeModel *AuthRequestWriteModel) (_ *CurrentAuthRequest) {
	return &CurrentAuthRequest{
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
			ResponseMode:  writeModel.ResponseMode,
			CodeChallenge: writeModel.CodeChallenge,
			Prompt:        writeModel.Prompt,
			UILocales:     writeModel.UILocales,
			MaxAge:        writeModel.MaxAge,
			LoginHint:     writeModel.LoginHint,
			HintUserID:    writeModel.HintUserID,
		},
		SessionID:   writeModel.SessionID,
		UserID:      writeModel.UserID,
		AuthMethods: writeModel.AuthMethods,
		AuthTime:    writeModel.AuthTime,
	}
}

func (c *Commands) GetCurrentAuthRequest(ctx context.Context, id string) (_ *CurrentAuthRequest, err error) {
	wm, err := c.getAuthRequestWriteModel(ctx, id)
	if err != nil {
		return nil, err
	}
	return authRequestWriteModelToCurrentAuthRequest(wm), nil
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
