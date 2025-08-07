package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SAMLRequest struct {
	ID          string
	LoginClient string

	ApplicationID  string
	ACSURL         string
	RelayState     string
	RequestID      string
	Binding        string
	Issuer         string
	Destination    string
	ResponseIssuer string
}

type CurrentSAMLRequest struct {
	*SAMLRequest
	SessionID   string
	UserID      string
	AuthMethods []domain.UserAuthMethodType
	AuthTime    time.Time
}

func (c *Commands) AddSAMLRequest(ctx context.Context, samlRequest *SAMLRequest) (_ *CurrentSAMLRequest, err error) {
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	samlRequest.ID = IDPrefixV2 + id
	writeModel, err := c.getSAMLRequestWriteModel(ctx, samlRequest.ID)
	if err != nil {
		return nil, err
	}
	if writeModel.SAMLRequestState != domain.SAMLRequestStateUnspecified {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-MO3vmsMLUt", "Errors.SAMLRequest.AlreadyExisting")
	}
	err = c.pushAppendAndReduce(ctx, writeModel, samlrequest.NewAddedEvent(
		ctx,
		&samlrequest.NewAggregate(samlRequest.ID, authz.GetInstance(ctx).InstanceID()).Aggregate,
		samlRequest.LoginClient,
		samlRequest.ApplicationID,
		samlRequest.ACSURL,
		samlRequest.RelayState,
		samlRequest.RequestID,
		samlRequest.Binding,
		samlRequest.Issuer,
		samlRequest.Destination,
		samlRequest.ResponseIssuer,
	))
	if err != nil {
		return nil, err
	}
	return samlRequestWriteModelToCurrentSAMLRequest(writeModel), nil
}

func (c *Commands) LinkSessionToSAMLRequest(ctx context.Context, id, sessionID, sessionToken string, checkLoginClient bool, projectPermissionCheck domain.ProjectPermissionCheck) (*domain.ObjectDetails, *CurrentSAMLRequest, error) {
	writeModel, err := c.getSAMLRequestWriteModel(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if writeModel.SAMLRequestState == domain.SAMLRequestStateUnspecified {
		return nil, nil, zerrors.ThrowNotFound(nil, "COMMAND-GH3PVLSfXC", "Errors.SAMLRequest.NotExisting")
	}
	if writeModel.SAMLRequestState != domain.SAMLRequestStateAdded {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-ttPKNdAIFT", "Errors.SAMLRequest.AlreadyHandled")
	}
	if checkLoginClient && authz.GetCtxData(ctx).UserID != writeModel.LoginClient {
		if err := c.checkPermission(ctx, domain.PermissionSessionLink, writeModel.ResourceOwner, ""); err != nil {
			return nil, nil, err
		}
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

	if projectPermissionCheck != nil {
		if err := projectPermissionCheck(ctx, writeModel.Issuer, sessionWriteModel.UserID); err != nil {
			return nil, nil, err
		}
	}

	if err := c.pushAppendAndReduce(ctx, writeModel, samlrequest.NewSessionLinkedEvent(
		ctx, &samlrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		sessionID,
		sessionWriteModel.UserID,
		sessionWriteModel.AuthenticationTime(),
		sessionWriteModel.AuthMethodTypes(),
	)); err != nil {
		return nil, nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), samlRequestWriteModelToCurrentSAMLRequest(writeModel), nil
}

func (c *Commands) FailSAMLRequest(ctx context.Context, id string, reason domain.SAMLErrorReason) (*domain.ObjectDetails, *CurrentSAMLRequest, error) {
	writeModel, err := c.getSAMLRequestWriteModel(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if writeModel.SAMLRequestState != domain.SAMLRequestStateAdded {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-32lGj1Fhjt", "Errors.SAMLRequest.AlreadyHandled")
	}
	if err := c.checkPermission(ctx, domain.PermissionSessionLink, writeModel.ResourceOwner, ""); err != nil {
		return nil, nil, err
	}
	err = c.pushAppendAndReduce(ctx, writeModel, samlrequest.NewFailedEvent(
		ctx,
		&samlrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
		reason,
	))
	if err != nil {
		return nil, nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), samlRequestWriteModelToCurrentSAMLRequest(writeModel), nil
}

func samlRequestWriteModelToCurrentSAMLRequest(writeModel *SAMLRequestWriteModel) (_ *CurrentSAMLRequest) {
	return &CurrentSAMLRequest{
		SAMLRequest: &SAMLRequest{
			ID:             writeModel.AggregateID,
			LoginClient:    writeModel.LoginClient,
			ApplicationID:  writeModel.ApplicationID,
			ACSURL:         writeModel.ACSURL,
			RelayState:     writeModel.RelayState,
			RequestID:      writeModel.RequestID,
			Binding:        writeModel.Binding,
			Issuer:         writeModel.Issuer,
			Destination:    writeModel.Destination,
			ResponseIssuer: writeModel.ResponseIssuer,
		},
		SessionID:   writeModel.SessionID,
		UserID:      writeModel.UserID,
		AuthMethods: writeModel.AuthMethods,
		AuthTime:    writeModel.AuthTime,
	}
}

func (c *Commands) GetCurrentSAMLRequest(ctx context.Context, id string) (_ *CurrentSAMLRequest, err error) {
	wm, err := c.getSAMLRequestWriteModel(ctx, id)
	if err != nil {
		return nil, err
	}
	return samlRequestWriteModelToCurrentSAMLRequest(wm), nil
}

func (c *Commands) getSAMLRequestWriteModel(ctx context.Context, id string) (writeModel *SAMLRequestWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewSAMLRequestWriteModel(ctx, id)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
