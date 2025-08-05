package command

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddDeviceAuth(ctx context.Context, clientID, deviceCode, userCode string, expires time.Time, scopes, audience []string, needRefreshToken bool) (*domain.ObjectDetails, error) {
	aggr := deviceauth.NewAggregate(deviceCode, authz.GetInstance(ctx).InstanceID())
	model := NewDeviceAuthWriteModel(deviceCode, aggr.ResourceOwner)

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewAddedEvent(
		ctx,
		aggr,
		clientID,
		deviceCode,
		userCode,
		expires,
		scopes,
		audience,
		needRefreshToken,
	))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) ApproveDeviceAuth(
	ctx context.Context,
	deviceCode,
	userID,
	userOrgID string,
	authMethods []domain.UserAuthMethodType,
	authTime time.Time,
	preferredLanguage *language.Tag,
	userAgent *domain.UserAgent,
	sessionID string,
) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-Hief9", "Errors.DeviceAuth.NotFound")
	}
	if model.State != domain.DeviceAuthStateInitiated {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-GEJL3", "Errors.DeviceAuth.AlreadyHandled")
	}
	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewApprovedEvent(ctx, model.aggregate, userID, userOrgID, authMethods, authTime, preferredLanguage, userAgent, sessionID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) ApproveDeviceAuthWithSession(
	ctx context.Context,
	deviceCode,
	sessionID,
	sessionToken string,
) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-D2hf2", "Errors.DeviceAuth.NotFound")
	}
	if model.State != domain.DeviceAuthStateInitiated {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-D30Jf", "Errors.DeviceAuth.AlreadyHandled")
	}
	if err := c.checkPermission(ctx, domain.PermissionSessionLink, model.ResourceOwner, ""); err != nil {
		return nil, err
	}

	sessionWriteModel := NewSessionWriteModel(sessionID, authz.GetInstance(ctx).InstanceID())
	err = c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return nil, err
	}
	if err = sessionWriteModel.CheckIsActive(); err != nil {
		return nil, err
	}
	if err := c.sessionTokenVerifier(ctx, sessionToken, sessionWriteModel.AggregateID, sessionWriteModel.TokenID); err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewApprovedEvent(
		ctx,
		model.aggregate,
		sessionWriteModel.UserID,
		sessionWriteModel.UserResourceOwner,
		sessionWriteModel.AuthMethodTypes(),
		sessionWriteModel.AuthenticationTime(),
		sessionWriteModel.PreferredLanguage,
		sessionWriteModel.UserAgent,
		sessionID,
	))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) CancelDeviceAuth(ctx context.Context, id string, reason domain.DeviceAuthCanceled) (*domain.ObjectDetails, error) {
	model, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, id)
	if err != nil {
		return nil, err
	}
	if !model.State.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-gee5A", "Errors.DeviceAuth.NotFound")
	}
	if err := c.checkPermission(ctx, domain.PermissionSessionLink, model.ResourceOwner, ""); err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, deviceauth.NewCanceledEvent(ctx, model.aggregate, reason))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(model, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) getDeviceAuthWriteModelByDeviceCode(ctx context.Context, deviceCode string) (*DeviceAuthWriteModel, error) {
	model := &DeviceAuthWriteModel{
		WriteModel: eventstore.WriteModel{AggregateID: deviceCode},
	}
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	model.aggregate = deviceauth.NewAggregate(model.AggregateID, model.InstanceID)
	return model, nil
}

type DeviceAuthStateError domain.DeviceAuthState

func (e DeviceAuthStateError) Error() string {
	return fmt.Sprintf("device auth state not approved: %s", domain.DeviceAuthState(e).String())
}

// CreateOIDCSessionFromDeviceAuth creates a new OIDC session if the device authorization
// flow is completed (user logged in).
// A [DeviceAuthStateError] is returned if the device authorization was not approved,
// containing a [domain.DeviceAuthState] which can be used to inform the client about the state.
//
// As devices can poll at various intervals, an explicit state takes precedence over expiry.
// This is to prevent cases where users might approve or deny the authorization on time, but the next poll
// happens after expiry.
func (c *Commands) CreateOIDCSessionFromDeviceAuth(ctx context.Context, deviceCode, backChannelLogoutURI string) (_ *OIDCSession, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deviceAuthModel, err := c.getDeviceAuthWriteModelByDeviceCode(ctx, deviceCode)
	if err != nil {
		return nil, err
	}

	switch deviceAuthModel.State {
	case domain.DeviceAuthStateApproved:
		break
	case domain.DeviceAuthStateUndefined:
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ua1Vo", "Errors.DeviceAuth.NotFound")

	case domain.DeviceAuthStateInitiated:
		if deviceAuthModel.Expires.Before(time.Now()) {
			c.asyncPush(ctx, deviceauth.NewCanceledEvent(ctx, deviceAuthModel.aggregate, domain.DeviceAuthCanceledExpired))
			return nil, DeviceAuthStateError(domain.DeviceAuthStateExpired)
		}
		fallthrough
	case domain.DeviceAuthStateDenied, domain.DeviceAuthStateExpired, domain.DeviceAuthStateDone:
		fallthrough
	default:
		return nil, DeviceAuthStateError(deviceAuthModel.State)
	}

	cmd, err := c.newOIDCSessionAddEvents(ctx, deviceAuthModel.UserID, deviceAuthModel.UserOrgID)
	if err != nil {
		return nil, err
	}

	cmd.AddSession(ctx,
		deviceAuthModel.UserID,
		deviceAuthModel.UserOrgID,
		deviceAuthModel.SessionID,
		deviceAuthModel.ClientID,
		deviceAuthModel.Audience,
		deviceAuthModel.Scopes,
		deviceAuthModel.UserAuthMethods,
		deviceAuthModel.AuthTime,
		"",
		deviceAuthModel.PreferredLanguage,
		deviceAuthModel.UserAgent,
	)
	cmd.RegisterLogout(ctx, deviceAuthModel.SessionID, deviceAuthModel.UserID, deviceAuthModel.ClientID, backChannelLogoutURI)
	if err = cmd.AddAccessToken(ctx, deviceAuthModel.Scopes, deviceAuthModel.UserID, deviceAuthModel.UserOrgID, domain.TokenReasonAuthRequest, nil); err != nil {
		return nil, err
	}

	if deviceAuthModel.NeedRefreshToken {
		if err = cmd.AddRefreshToken(ctx, deviceAuthModel.UserID); err != nil {
			return nil, err
		}
	}
	cmd.DeviceAuthRequestDone(ctx, deviceAuthModel.aggregate)
	return cmd.PushEvents(ctx)
}

func (cmd *OIDCSessionEvents) DeviceAuthRequestDone(ctx context.Context, deviceAuthAggregate *eventstore.Aggregate) {
	cmd.events = append(cmd.events, deviceauth.NewDoneEvent(ctx, deviceAuthAggregate))
}
