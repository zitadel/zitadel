package command

import (
	"context"
	"io"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) CreateOTPSMSChallengeReturnCode(dst *string) SessionCommand {
	return c.createOTPSMSChallenge(true, dst)
}

func (c *Commands) CreateOTPSMSChallenge() SessionCommand {
	return c.createOTPSMSChallenge(false, nil)
}

func (c *Commands) createOTPSMSChallenge(returnCode bool, dst *string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		if cmd.sessionWriteModel.UserID == "" {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKL3g", "Errors.User.UserIDMissing")
		}
		writeModel := NewHumanOTPSMSWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
			return nil, err
		}
		if !writeModel.OTPAdded() {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady")
		}
		code, generatorID, err := cmd.createPhoneCode(ctx, cmd.eventstore.Filter, domain.SecretGeneratorTypeOTPSMS, cmd.otpAlg, c.defaultSecretGenerators.OTPSMS) //nolint:staticcheck
		if err != nil {
			return nil, err
		}
		if returnCode {
			*dst = code.Plain
		}
		cmd.OTPSMSChallenged(ctx, code.CryptedCode(), code.CodeExpiry(), returnCode, generatorID)
		return nil, nil
	}
}

func (c *Commands) OTPSMSSent(ctx context.Context, sessionID, resourceOwner string, generatorInfo *senders.CodeGeneratorInfo) error {
	sessionWriteModel := NewSessionWriteModel(sessionID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return err
	}
	if sessionWriteModel.OTPSMSCodeChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-G3t31", "Errors.User.Code.NotFound")
	}
	return c.pushAppendAndReduce(ctx, sessionWriteModel,
		session.NewOTPSMSSentEvent(ctx, &session.NewAggregate(sessionID, sessionWriteModel.ResourceOwner).Aggregate, generatorInfo),
	)
}

func (c *Commands) CreateOTPEmailChallengeURLTemplate(urlTmpl string) (SessionCommand, error) {
	if err := domain.RenderOTPEmailURLTemplate(io.Discard, urlTmpl, "code", "userID", "loginName", "displayName", "sessionID", language.English); err != nil {
		return nil, err
	}
	return c.createOTPEmailChallenge(false, urlTmpl, nil), nil
}

func (c *Commands) CreateOTPEmailChallengeReturnCode(dst *string) SessionCommand {
	return c.createOTPEmailChallenge(true, "", dst)
}

func (c *Commands) CreateOTPEmailChallenge() SessionCommand {
	return c.createOTPEmailChallenge(false, "", nil)
}

func (c *Commands) createOTPEmailChallenge(returnCode bool, urlTmpl string, dst *string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		if cmd.sessionWriteModel.UserID == "" {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-JK3gp", "Errors.User.UserIDMissing")
		}
		writeModel := NewHumanOTPEmailWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
			return nil, err
		}
		if !writeModel.OTPAdded() {
			return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady")
		}
		code, err := cmd.createCode(ctx, cmd.eventstore.Filter, domain.SecretGeneratorTypeOTPEmail, cmd.otpAlg, c.defaultSecretGenerators.OTPEmail) //nolint:staticcheck
		if err != nil {
			return nil, err
		}
		if returnCode {
			*dst = code.Plain
		}
		cmd.OTPEmailChallenged(ctx, code.Crypted, code.Expiry, returnCode, urlTmpl)
		return nil, nil
	}
}

func (c *Commands) OTPEmailSent(ctx context.Context, sessionID, resourceOwner string) error {
	sessionWriteModel := NewSessionWriteModel(sessionID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, sessionWriteModel)
	if err != nil {
		return err
	}
	if sessionWriteModel.OTPEmailCodeChallenge == nil {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-SLr02", "Errors.User.Code.NotFound")
	}
	return c.pushAppendAndReduce(ctx, sessionWriteModel,
		session.NewOTPEmailSentEvent(ctx, &session.NewAggregate(sessionID, sessionWriteModel.ResourceOwner).Aggregate),
	)
}

func CheckOTPSMS(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (_ []eventstore.Command, err error) {
		writeModel := func(ctx context.Context, userID string, resourceOwner string) (OTPCodeWriteModel, error) {
			otpWriteModel := NewHumanOTPSMSCodeWriteModel(cmd.sessionWriteModel.UserID, "")
			err := cmd.eventstore.FilterToQueryReducer(ctx, otpWriteModel)
			if err != nil {
				return nil, err
			}
			// explicitly set the challenge from the session write model since the code write model will only check user events
			otpWriteModel.otpCode = cmd.sessionWriteModel.OTPSMSCodeChallenge
			return otpWriteModel, nil
		}
		succeededEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
			return user.NewHumanOTPSMSCheckSucceededEvent(ctx, aggregate, nil)
		}
		failedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
			return user.NewHumanOTPSMSCheckFailedEvent(ctx, aggregate, nil)
		}
		commands, err := checkOTP(
			ctx,
			cmd.sessionWriteModel.UserID,
			code,
			"",
			nil,
			writeModel,
			cmd.eventstore.FilterToQueryReducer,
			cmd.otpAlg,
			cmd.getCodeVerifier,
			succeededEvent,
			failedEvent,
		)
		if err != nil {
			return commands, err
		}
		cmd.eventCommands = append(cmd.eventCommands, commands...)
		cmd.OTPSMSChecked(ctx, cmd.now())
		return nil, nil
	}
}

func CheckOTPEmail(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (_ []eventstore.Command, err error) {
		writeModel := func(ctx context.Context, userID string, resourceOwner string) (OTPCodeWriteModel, error) {
			otpWriteModel := NewHumanOTPEmailCodeWriteModel(cmd.sessionWriteModel.UserID, "")
			err := cmd.eventstore.FilterToQueryReducer(ctx, otpWriteModel)
			if err != nil {
				return nil, err
			}
			// explicitly set the challenge from the session write model since the code write model will only check user events
			otpWriteModel.otpCode = cmd.sessionWriteModel.OTPEmailCodeChallenge
			return otpWriteModel, nil
		}
		succeededEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
			return user.NewHumanOTPEmailCheckSucceededEvent(ctx, aggregate, nil)
		}
		failedEvent := func(ctx context.Context, aggregate *eventstore.Aggregate, info *user.AuthRequestInfo) eventstore.Command {
			return user.NewHumanOTPEmailCheckFailedEvent(ctx, aggregate, nil)
		}
		commands, err := checkOTP(
			ctx,
			cmd.sessionWriteModel.UserID,
			code,
			"",
			nil,
			writeModel,
			cmd.eventstore.FilterToQueryReducer,
			cmd.otpAlg,
			nil, // email currently always uses local code checks
			succeededEvent,
			failedEvent,
		)
		if err != nil {
			return commands, err
		}
		cmd.eventCommands = append(cmd.eventCommands, commands...)
		cmd.OTPEmailChecked(ctx, cmd.now())
		return nil, nil
	}
}
