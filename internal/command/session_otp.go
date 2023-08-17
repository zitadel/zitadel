package command

import (
	"context"
	"io"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func (c *Commands) CreateOTPSMSChallengeReturnCode(dst *string) SessionCommand {
	return c.createOTPSMSChallenge(true, dst)
}

func (c *Commands) CreateOTPSMSChallenge() SessionCommand {
	return c.createOTPSMSChallenge(false, nil)
}

func (c *Commands) createOTPSMSChallenge(returnCode bool, dst *string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		writeModel := NewHumanOTPSMSWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
			return err
		}
		if !writeModel.OTPAdded() {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady")
		}
		code, plain, expiry, err := cmd.generate(ctx, domain.SecretGeneratorTypeOTPSMS)
		if err != nil {
			return err
		}
		if returnCode {
			*dst = plain
		}
		cmd.OTPSMSChallenged(ctx, code, expiry, returnCode)
		return nil
	}
}

func (c *Commands) CreateOTPEmailChallengeURLTemplate(urlTmpl string) (SessionCommand, error) {
	if err := domain.RenderOTPEmailURLTemplate(io.Discard, urlTmpl, "userID", "code"); err != nil {
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
	return func(ctx context.Context, cmd *SessionCommands) error {
		writeModel := NewHumanOTPEmailWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
			return err
		}
		if !writeModel.OTPAdded() {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady")
		}
		code, plain, expiry, err := cmd.generate(ctx, domain.SecretGeneratorTypeOTPEmail)
		if err != nil {
			return err
		}
		if returnCode {
			*dst = plain
		}
		cmd.OTPEmailChallenged(ctx, code, expiry, returnCode, urlTmpl)
		return nil
	}
}

func CheckOTPSMS(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (err error) {
		if cmd.sessionWriteModel.UserID == "" {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-VDrh3", "Errors.User.UserIDMissing")
		}
		challenge := cmd.sessionWriteModel.OTPSMSCodeChallenge
		if challenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SF3tv", "Errors.User.Code.NotFound")
		}
		err = crypto.VerifyCodeWithAlgorithm(challenge.CreationDate, challenge.Expiry, challenge.Code, code, cmd.otpAlg)
		if err != nil {
			return err
		}
		cmd.OTPSMSChecked(ctx, cmd.now())
		return nil
	}
}

func CheckOTPEmail(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) (err error) {
		if cmd.sessionWriteModel.UserID == "" {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-ejo2w", "Errors.User.UserIDMissing")
		}
		challenge := cmd.sessionWriteModel.OTPEmailCodeChallenge
		if challenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-zF3g3", "Errors.User.Code.NotFound")
		}
		err = crypto.VerifyCodeWithAlgorithm(challenge.CreationDate, challenge.Expiry, challenge.Code, code, cmd.otpAlg)
		if err != nil {
			return err
		}
		cmd.OTPEmailChecked(ctx, cmd.now())
		return nil
	}
}
