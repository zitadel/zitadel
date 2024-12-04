package command

import (
	"context"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateUserInvite struct {
	UserID          string
	URLTemplate     string
	ReturnCode      bool
	ApplicationName string
}

func (c *Commands) CreateInviteCode(ctx context.Context, invite *CreateUserInvite) (details *domain.ObjectDetails, returnCode *string, err error) {
	invite.UserID = strings.TrimSpace(invite.UserID)
	if invite.UserID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4jio3", "Errors.User.UserIDMissing")
	}
	wm, err := c.userInviteCodeWriteModel(ctx, invite.UserID, "")
	if err != nil {
		return nil, nil, err
	}
	if err := c.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, nil, err
	}
	if !wm.UserState.Exists() {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wgvn4", "Errors.User.NotFound")
	}
	if !wm.CreationAllowed() {
		return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-EF34g", "Errors.User.AlreadyInitialised")
	}
	code, err := c.newUserInviteCode(ctx, c.eventstore.Filter, c.userEncryption) //nolint
	if err != nil {
		return nil, nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, user.NewHumanInviteCodeAddedEvent(
		ctx,
		UserAggregateFromWriteModelCtx(ctx, &wm.WriteModel),
		code.Crypted,
		code.Expiry,
		invite.URLTemplate,
		invite.ReturnCode,
		invite.ApplicationName,
		"",
	))
	if err != nil {
		return nil, nil, err
	}
	if invite.ReturnCode {
		returnCode = &code.Plain
	}
	return writeModelToObjectDetails(&wm.WriteModel), returnCode, nil
}

// ResendInviteCode resends the invite mail with a new code and an optional authRequestID.
// It will reuse the applicationName from the previous code.
func (c *Commands) ResendInviteCode(ctx context.Context, userID, resourceOwner, authRequestID string) (objectDetails *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2n8vs", "Errors.User.UserIDMissing")
	}

	existingCode, err := c.userInviteCodeWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if err := c.checkPermissionUpdateUser(ctx, existingCode.ResourceOwner, userID); err != nil {
		return nil, err
	}
	if !existingCode.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-H3b2a", "Errors.User.NotFound")
	}
	if !existingCode.CreationAllowed() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Gg42s", "Errors.User.AlreadyInitialised")
	}
	if existingCode.InviteCode == nil || existingCode.CodeReturned {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wr3gq", "Errors.User.Code.NotFound")
	}
	code, err := c.newUserInviteCode(ctx, c.eventstore.Filter, c.userEncryption) //nolint
	if err != nil {
		return nil, err
	}
	if authRequestID == "" {
		authRequestID = existingCode.AuthRequestID
	}
	err = c.pushAppendAndReduce(ctx, existingCode,
		user.NewHumanInviteCodeAddedEvent(
			ctx,
			UserAggregateFromWriteModelCtx(ctx, &existingCode.WriteModel),
			code.Crypted,
			code.Expiry,
			existingCode.URLTemplate,
			false,
			existingCode.ApplicationName,
			authRequestID,
		))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingCode.WriteModel), nil
}

func (c *Commands) InviteCodeSent(ctx context.Context, userID, orgID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Sgf31", "Errors.User.UserIDMissing")
	}
	existingCode, err := c.userInviteCodeWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !existingCode.UserState.Exists() {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-HN34a", "Errors.User.NotFound")
	}
	if existingCode.InviteCode == nil || existingCode.CodeReturned {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wr3gq", "Errors.User.Code.NotFound")
	}
	userAgg := UserAggregateFromWriteModelCtx(ctx, &existingCode.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanInviteCodeSentEvent(ctx, userAgg))
	return err
}

func (c *Commands) VerifyInviteCode(ctx context.Context, userID, code string) (details *domain.ObjectDetails, err error) {
	return c.VerifyInviteCodeSetPassword(ctx, userID, code, "", "")
}

func (c *Commands) VerifyInviteCodeSetPassword(ctx context.Context, userID, code, password, userAgentID string) (details *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Gk3f2", "Errors.User.UserIDMissing")
	}
	wm, err := c.userInviteCodeWriteModel(ctx, userID, "")
	if err != nil {
		return nil, err
	}
	if !wm.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-F5g2h", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModelCtx(ctx, &wm.WriteModel)
	err = crypto.VerifyCode(wm.InviteCodeCreationDate, wm.InviteCodeExpiry, wm.InviteCode, code, c.userEncryption)
	if err != nil {
		_, err = c.eventstore.Push(ctx, user.NewHumanInviteCheckFailedEvent(ctx, userAgg))
		logging.WithFields("userID", userAgg.ID).OnError(err).Error("NewHumanInviteCheckFailedEvent push failed")
		return nil, zerrors.ThrowInvalidArgument(err, "COMMAND-Wgn4q", "Errors.User.Code.Invalid")
	}
	commands := []eventstore.Command{
		user.NewHumanInviteCheckSucceededEvent(ctx, userAgg),
		user.NewHumanEmailVerifiedEvent(ctx, userAgg),
	}
	if password != "" {
		passwordCommand, err := c.setPasswordCommand(
			ctx,
			userAgg,
			wm.UserState,
			password,
			"",
			userAgentID,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
		commands = append(commands, passwordCommand)
	}
	err = c.pushAppendAndReduce(ctx, wm, commands...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) userInviteCodeWriteModel(ctx context.Context, userID, orgID string) (writeModel *UserV2InviteWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = newUserV2InviteWriteModel(userID, orgID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
