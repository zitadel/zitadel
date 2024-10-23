package command

import (
	"context"
	"io"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// ChangeUserEmail sets a user's email address, generates a code
// and triggers a notification e-mail with the default confirmation URL format.
func (c *Commands) ChangeUserEmail(ctx context.Context, userID, email string, alg crypto.EncryptionAlgorithm) (*domain.Email, error) {
	return c.changeUserEmailWithCode(ctx, userID, email, alg, false, "")
}

// ChangeUserEmailURLTemplate sets a user's email address, generates a code
// and triggers a notification e-mail with the confirmation URL rendered from the passed urlTmpl.
// urlTmpl must be a valid [tmpl.Template].
func (c *Commands) ChangeUserEmailURLTemplate(ctx context.Context, userID, email string, alg crypto.EncryptionAlgorithm, urlTmpl string) (*domain.Email, error) {
	if err := domain.RenderConfirmURLTemplate(io.Discard, urlTmpl, userID, "code", "orgID"); err != nil {
		return nil, err
	}
	return c.changeUserEmailWithCode(ctx, userID, email, alg, false, urlTmpl)
}

// ChangeUserEmailReturnCode sets a user's email address, generates a code and does not send a notification email.
// The generated plain text code will be set in the returned Email object.
func (c *Commands) ChangeUserEmailReturnCode(ctx context.Context, userID, email string, alg crypto.EncryptionAlgorithm) (*domain.Email, error) {
	return c.changeUserEmailWithCode(ctx, userID, email, alg, true, "")
}

// ResendUserEmailCode generates a new code if there is a code existing
// and triggers a notification e-mail with the default confirmation URL format.
func (c *Commands) ResendUserEmailCode(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm) (*domain.Email, error) {
	return c.resendUserEmailCode(ctx, userID, alg, false, "")
}

// ResendUserEmailCodeURLTemplate generates a new code if there is a code existing
// and triggers a notification e-mail with the confirmation URL rendered from the passed urlTmpl.
// urlTmpl must be a valid [tmpl.Template].
func (c *Commands) ResendUserEmailCodeURLTemplate(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm, urlTmpl string) (*domain.Email, error) {
	if err := domain.RenderConfirmURLTemplate(io.Discard, urlTmpl, userID, "code", "orgID"); err != nil {
		return nil, err
	}
	return c.resendUserEmailCode(ctx, userID, alg, false, urlTmpl)
}

// ResendUserEmailReturnCode generates a new code if there is a code existing and does not send a notification email.
// The generated plain text code will be set in the returned Email object.
func (c *Commands) ResendUserEmailReturnCode(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm) (*domain.Email, error) {
	return c.resendUserEmailCode(ctx, userID, alg, true, "")
}

// ChangeUserEmailVerified sets a user's email address and marks it is verified.
// No code is generated and no confirmation e-mail is send.
func (c *Commands) ChangeUserEmailVerified(ctx context.Context, userID, email string) (*domain.Email, error) {
	cmd, err := c.NewUserEmailEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermission(ctx, domain.PermissionUserWrite, cmd.aggregate.ResourceOwner, userID); err != nil {
		return nil, err
	}
	if err = cmd.Change(ctx, domain.EmailAddress(email)); err != nil {
		return nil, err
	}
	cmd.SetVerified(ctx)
	return cmd.Push(ctx)
}

func (c *Commands) changeUserEmailWithCode(ctx context.Context, userID, email string, alg crypto.EncryptionAlgorithm, returnCode bool, urlTmpl string) (*domain.Email, error) {
	config, err := cryptoGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyEmailCode) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	gen := crypto.NewEncryptionGenerator(*config, alg)
	return c.changeUserEmailWithGenerator(ctx, userID, email, gen, returnCode, urlTmpl)
}

func (c *Commands) resendUserEmailCode(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm, returnCode bool, urlTmpl string) (*domain.Email, error) {
	config, err := cryptoGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyEmailCode) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	gen := crypto.NewEncryptionGenerator(*config, alg)
	return c.resendUserEmailCodeWithGenerator(ctx, userID, gen, returnCode, urlTmpl)
}

// changeUserEmailWithGenerator set a user's email address.
// returnCode controls if the plain text version of the code will be set in the return object.
// When the plain text code is returned, no notification e-mail will be send to the user.
// urlTmpl allows changing the target URL that is used by the e-mail and should be a validated Go template, if used.
func (c *Commands) changeUserEmailWithGenerator(ctx context.Context, userID, email string, gen crypto.Generator, returnCode bool, urlTmpl string) (*domain.Email, error) {
	cmd, err := c.changeUserEmailWithGeneratorEvents(ctx, userID, email, gen, returnCode, urlTmpl)
	if err != nil {
		return nil, err
	}
	return cmd.Push(ctx)
}

func (c *Commands) resendUserEmailCodeWithGenerator(ctx context.Context, userID string, gen crypto.Generator, returnCode bool, urlTmpl string) (*domain.Email, error) {
	cmd, err := c.resendUserEmailCodeWithGeneratorEvents(ctx, userID, gen, returnCode, urlTmpl)
	if err != nil {
		return nil, err
	}
	return cmd.Push(ctx)
}

func (c *Commands) changeUserEmailWithGeneratorEvents(ctx context.Context, userID, email string, gen crypto.Generator, returnCode bool, urlTmpl string) (*UserEmailEvents, error) {
	cmd, err := c.NewUserEmailEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermissionUpdateUser(ctx, cmd.aggregate.ResourceOwner, userID); err != nil {
		return nil, err
	}
	if err = cmd.Change(ctx, domain.EmailAddress(email)); err != nil {
		return nil, err
	}
	if err = cmd.AddGeneratedCode(ctx, gen, urlTmpl, returnCode); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Commands) resendUserEmailCodeWithGeneratorEvents(ctx context.Context, userID string, gen crypto.Generator, returnCode bool, urlTmpl string) (*UserEmailEvents, error) {
	cmd, err := c.NewUserEmailEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermissionUpdateUser(ctx, cmd.aggregate.ResourceOwner, userID); err != nil {
		return nil, err
	}
	if cmd.model.Code == nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "EMAIL-5w5ilin4yt", "Errors.User.Code.Empty")
	}
	if err = cmd.AddGeneratedCode(ctx, gen, urlTmpl, returnCode); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Commands) VerifyUserEmail(ctx context.Context, userID, code string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	config, err := cryptoGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyEmailCode) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	gen := crypto.NewEncryptionGenerator(*config, alg)
	return c.verifyUserEmailWithGenerator(ctx, userID, code, gen)
}

func (c *Commands) verifyUserEmailWithGenerator(ctx context.Context, userID, code string, gen crypto.Generator) (*domain.ObjectDetails, error) {
	cmd, err := c.NewUserEmailEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	err = cmd.VerifyCode(ctx, code, gen)
	if err != nil {
		return nil, err
	}
	if _, err = cmd.Push(ctx); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&cmd.model.WriteModel), nil
}

// UserEmailEvents allows step-by-step additions of events,
// operating on the Human Email Model.
type UserEmailEvents struct {
	eventstore *eventstore.Eventstore
	aggregate  *eventstore.Aggregate
	events     []eventstore.Command
	model      *HumanEmailWriteModel

	plainCode *string
}

// NewUserEmailEvents constructs a UserEmailEvents with a Human Email Write Model,
// filtered by userID and resourceOwner.
// If a model cannot be found, or it's state is invalid and error is returned.
func (c *Commands) NewUserEmailEvents(ctx context.Context, userID string) (*UserEmailEvents, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0Gzs3", "Errors.User.Email.IDMissing")
	}

	model, err := c.emailWriteModel(ctx, userID, "")
	if err != nil {
		return nil, err
	}
	if model.UserState == domain.UserStateUnspecified || model.UserState == domain.UserStateDeleted {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ieJ2e", "Errors.User.Email.NotFound")
	}
	if model.UserState == domain.UserStateInitial {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-uz0Uu", "Errors.User.NotInitialised")
	}
	return &UserEmailEvents{
		eventstore: c.eventstore,
		aggregate:  UserAggregateFromWriteModel(&model.WriteModel),
		model:      model,
	}, nil
}

// Change sets a new email address.
// The generated event unsets any previously generated code and verified flag.
func (c *UserEmailEvents) Change(ctx context.Context, email domain.EmailAddress) error {
	if err := email.Validate(); err != nil {
		return err
	}
	event, hasChanged := c.model.NewChangedEvent(ctx, c.aggregate, email)
	if !hasChanged {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Email.NotChanged")
	}
	c.events = append(c.events, event)
	return nil
}

// SetVerified sets the email address to verified.
func (c *UserEmailEvents) SetVerified(ctx context.Context) {
	c.events = append(c.events, user.NewHumanEmailVerifiedEvent(ctx, c.aggregate))
}

// AddGeneratedCode generates a new encrypted code and sets it to the email address.
// When returnCode a plain text of the code will be returned from Push.
func (c *UserEmailEvents) AddGeneratedCode(ctx context.Context, gen crypto.Generator, urlTmpl string, returnCode bool) error {
	cmd, code, err := generateCodeCommand(ctx, c.aggregate, gen, urlTmpl, returnCode)
	if err != nil {
		return err
	}
	c.events = append(c.events, cmd)
	if returnCode {
		c.plainCode = &code
	}
	return nil
}

func generateCodeCommand(ctx context.Context, agg *eventstore.Aggregate, gen crypto.Generator, urlTmpl string, returnCode bool) (eventstore.Command, string, error) {
	value, plain, err := crypto.NewCode(gen)
	if err != nil {
		return nil, "", err
	}

	cmd := user.NewHumanEmailCodeAddedEventV2(ctx, agg, value, gen.Expiry(), urlTmpl, returnCode, "")
	if returnCode {
		return cmd, plain, nil
	}
	return cmd, "", nil
}

func (c *UserEmailEvents) VerifyCode(ctx context.Context, code string, gen crypto.Generator) error {
	if code == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Fia4a", "Errors.User.Code.Empty")
	}

	err := crypto.VerifyCode(c.model.CodeCreationDate, c.model.CodeExpiry, c.model.Code, code, gen.Alg())
	if err == nil {
		c.events = append(c.events, user.NewHumanEmailVerifiedEvent(ctx, c.aggregate))
		return nil
	}
	_, err = c.eventstore.Push(ctx, user.NewHumanEmailVerificationFailedEvent(ctx, c.aggregate))
	logging.WithFields("id", "COMMAND-Zoo6b", "userID", c.aggregate.ID).OnError(err).Error("NewHumanEmailVerificationFailedEvent push failed")
	return zerrors.ThrowInvalidArgument(err, "COMMAND-eis9R", "Errors.User.Code.Invalid")
}

// Push all events to the eventstore and Reduce them into the Model.
func (c *UserEmailEvents) Push(ctx context.Context) (*domain.Email, error) {
	pushedEvents, err := c.eventstore.Push(ctx, c.events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(c.model, pushedEvents...)
	if err != nil {
		return nil, err
	}
	email := writeModelToEmail(c.model)
	email.PlainCode = c.plainCode

	return email, nil
}
