package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// ChangeUserPhone sets a user's phone number, generates a code
// and triggers a notification sms.
func (c *Commands) ChangeUserPhone(ctx context.Context, userID, phone string, alg crypto.EncryptionAlgorithm) (*domain.Phone, error) {
	return c.changeUserPhoneWithCode(ctx, userID, phone, alg, false)
}

// ChangeUserPhoneReturnCode sets a user's phone number, generates a code and does not send a notification sms.
// The generated plain text code will be set in the returned Phone object.
func (c *Commands) ChangeUserPhoneReturnCode(ctx context.Context, userID, phone string, alg crypto.EncryptionAlgorithm) (*domain.Phone, error) {
	return c.changeUserPhoneWithCode(ctx, userID, phone, alg, true)
}

// ChangeUserPhoneVerified sets a user's phone number and marks it is verified.
// No code is generated and no confirmation sms is send.
func (c *Commands) ChangeUserPhoneVerified(ctx context.Context, userID, phone string) (*domain.Phone, error) {
	cmd, err := c.NewUserPhoneEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermission(ctx, domain.PermissionUserWrite, cmd.aggregate.ResourceOwner, userID); err != nil {
		return nil, err
	}
	if err = cmd.Change(ctx, domain.PhoneNumber(phone)); err != nil {
		return nil, err
	}
	cmd.SetVerified(ctx)
	return cmd.Push(ctx)
}

// ResendUserPhoneCode generates a code
// and triggers a notification sms.
func (c *Commands) ResendUserPhoneCode(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm) (*domain.Phone, error) {
	return c.resendUserPhoneCode(ctx, userID, alg, false)
}

// ResendUserPhoneCodeReturnCode generates a code and does not send a notification sms.
// The generated plain text code will be set in the returned Phone object.
func (c *Commands) ResendUserPhoneCodeReturnCode(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm) (*domain.Phone, error) {
	return c.resendUserPhoneCode(ctx, userID, alg, true)
}

func (c *Commands) changeUserPhoneWithCode(ctx context.Context, userID, phone string, alg crypto.EncryptionAlgorithm, returnCode bool) (*domain.Phone, error) {
	config, err := cryptoGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	gen := crypto.NewEncryptionGenerator(*config, alg)
	return c.changeUserPhoneWithGenerator(ctx, userID, phone, gen, returnCode)
}

func (c *Commands) resendUserPhoneCode(ctx context.Context, userID string, alg crypto.EncryptionAlgorithm, returnCode bool) (*domain.Phone, error) {
	config, err := cryptoGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	gen := crypto.NewEncryptionGenerator(*config, alg)
	return c.resendUserPhoneCodeWithGenerator(ctx, userID, gen, returnCode)
}

// changeUserPhoneWithGenerator set a user's phone number.
// returnCode controls if the plain text version of the code will be set in the return object.
// When the plain text code is returned, no notification sms will be send to the user.
func (c *Commands) changeUserPhoneWithGenerator(ctx context.Context, userID, phone string, gen crypto.Generator, returnCode bool) (*domain.Phone, error) {
	cmd, err := c.NewUserPhoneEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if authz.GetCtxData(ctx).UserID != userID {
		if err = c.checkPermission(ctx, domain.PermissionUserWrite, cmd.aggregate.ResourceOwner, userID); err != nil {
			return nil, err
		}
	}
	if err = cmd.Change(ctx, domain.PhoneNumber(phone)); err != nil {
		return nil, err
	}
	if err = cmd.AddGeneratedCode(ctx, gen, returnCode); err != nil {
		return nil, err
	}
	return cmd.Push(ctx)
}

// resendUserPhoneCodeWithGenerator generates a new code.
// returnCode controls if the plain text version of the code will be set in the return object.
// When the plain text code is returned, no notification sms will be send to the user.
func (c *Commands) resendUserPhoneCodeWithGenerator(ctx context.Context, userID string, gen crypto.Generator, returnCode bool) (*domain.Phone, error) {
	cmd, err := c.NewUserPhoneEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if authz.GetCtxData(ctx).UserID != userID {
		if err = c.checkPermission(ctx, domain.PermissionUserWrite, cmd.aggregate.ResourceOwner, userID); err != nil {
			return nil, err
		}
	}
	if cmd.model.Code == nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "PHONE-5xrra88eq8", "Errors.User.Code.Empty")
	}
	if err = cmd.AddGeneratedCode(ctx, gen, returnCode); err != nil {
		return nil, err
	}
	return cmd.Push(ctx)
}

func (c *Commands) VerifyUserPhone(ctx context.Context, userID, code string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	config, err := cryptoGeneratorConfig(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	gen := crypto.NewEncryptionGenerator(*config, alg)
	return c.verifyUserPhoneWithGenerator(ctx, userID, code, gen)
}

func (c *Commands) verifyUserPhoneWithGenerator(ctx context.Context, userID, code string, gen crypto.Generator) (*domain.ObjectDetails, error) {
	cmd, err := c.NewUserPhoneEvents(ctx, userID)
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

func (c *Commands) RemoveUserPhone(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	return c.removeUserPhone(ctx, userID)
}

func (c *Commands) removeUserPhone(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	cmd, err := c.NewUserPhoneEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = c.checkPermission(ctx, domain.PermissionUserWrite, cmd.aggregate.ResourceOwner, userID); err != nil {
		return nil, err
	}
	cmd.Remove(ctx)
	if _, err = cmd.Push(ctx); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&cmd.model.WriteModel), nil
}

// UserPhoneEvents allows step-by-step additions of events,
// operating on the Human Phone Model.
type UserPhoneEvents struct {
	eventstore *eventstore.Eventstore
	aggregate  *eventstore.Aggregate
	events     []eventstore.Command
	model      *HumanPhoneWriteModel

	plainCode *string
}

// NewUserPhoneEvents constructs a UserPhoneEvents with a Human Phone Write Model,
// filtered by userID and resourceOwner.
// If a model cannot be found, or it's state is invalid and error is returned.
func (c *Commands) NewUserPhoneEvents(ctx context.Context, userID string) (*UserPhoneEvents, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-xP292j", "Errors.User.Phone.IDMissing")
	}

	model, err := c.phoneWriteModelByID(ctx, userID, "")
	if err != nil {
		return nil, err
	}
	if model.UserState == domain.UserStateUnspecified || model.UserState == domain.UserStateDeleted {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ieJ2e", "Errors.User.Phone.NotFound")
	}
	if model.UserState == domain.UserStateInitial {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-uz0Uu", "Errors.User.NotInitialised")
	}
	return &UserPhoneEvents{
		eventstore: c.eventstore,
		aggregate:  UserAggregateFromWriteModel(&model.WriteModel),
		model:      model,
	}, nil
}

// Change sets a new phone number.
// The generated event unsets any previously generated code and verified flag.
func (c *UserPhoneEvents) Change(ctx context.Context, phone domain.PhoneNumber) error {
	phone, err := phone.Normalize()
	if err != nil {
		return err
	}
	event, hasChanged := c.model.NewChangedEvent(ctx, c.aggregate, phone)
	if !hasChanged {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-Uch5e", "Errors.User.Phone.NotChanged")
	}
	c.events = append(c.events, event)
	return nil
}

func (c *UserPhoneEvents) Remove(ctx context.Context) {
	c.events = append(c.events, user.NewHumanPhoneRemovedEvent(ctx, c.aggregate))
}

// SetVerified sets the phone number to verified.
func (c *UserPhoneEvents) SetVerified(ctx context.Context) {
	c.events = append(c.events, user.NewHumanPhoneVerifiedEvent(ctx, c.aggregate))
}

// AddGeneratedCode generates a new encrypted code and sets it to the phone number.
// When returnCode a plain text of the code will be returned from Push.
func (c *UserPhoneEvents) AddGeneratedCode(ctx context.Context, gen crypto.Generator, returnCode bool) error {
	value, plain, err := crypto.NewCode(gen)
	if err != nil {
		return err
	}

	c.events = append(c.events, user.NewHumanPhoneCodeAddedEventV2(ctx, c.aggregate, value, gen.Expiry(), returnCode))
	if returnCode {
		c.plainCode = &plain
	}
	return nil
}

func (c *UserPhoneEvents) VerifyCode(ctx context.Context, code string, gen crypto.Generator) error {
	if code == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Fia4a", "Errors.User.Code.Empty")
	}

	err := crypto.VerifyCode(c.model.CodeCreationDate, c.model.CodeExpiry, c.model.Code, code, gen.Alg())
	if err == nil {
		c.events = append(c.events, user.NewHumanPhoneVerifiedEvent(ctx, c.aggregate))
		return nil
	}
	_, err = c.eventstore.Push(ctx, user.NewHumanPhoneVerificationFailedEvent(ctx, c.aggregate))
	logging.WithFields("id", "COMMAND-Zoo6b", "userID", c.aggregate.ID).OnError(err).Error("NewHumanPhoneVerificationFailedEvent push failed")
	return zerrors.ThrowInvalidArgument(err, "COMMAND-eis9R", "Errors.User.Code.Invalid")
}

// Push all events to the eventstore and Reduce them into the Model.
func (c *UserPhoneEvents) Push(ctx context.Context) (*domain.Phone, error) {
	pushedEvents, err := c.eventstore.Push(ctx, c.events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(c.model, pushedEvents...)
	if err != nil {
		return nil, err
	}
	phone := writeModelToPhone(c.model)
	phone.PlainCode = c.plainCode

	return phone, nil
}
