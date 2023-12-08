package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type ChangeHuman struct {
	ID       string
	Username *string
	Profile  *Profile
	Email    *Email
	Phone    *Phone

	Password *Password

	// Details are set after a successful execution of the command
	Details *domain.ObjectDetails

	// EmailCode is set by the command
	EmailCode *string

	// PhoneCode is set by the command
	PhoneCode *string
}

type Profile struct {
	FirstName         *string
	LastName          *string
	NickName          *string
	DisplayName       *string
	PreferredLanguage *language.Tag
	Gender            *domain.Gender
}

type Password struct {
	// Either you have to have permission, a password code or the old password to change
	PasswordCode        *string
	OldPassword         *string
	Password            *string
	EncodedPasswordHash *string

	ChangeRequired bool
}

func (h *ChangeHuman) Validate(hasher *crypto.PasswordHasher) (err error) {
	if h.Email != nil && h.Email.Address != "" {
		if err := h.Email.Validate(); err != nil {
			return err
		}
	}

	if h.Phone != nil && h.Phone.Number != "" {
		if h.Phone.Number, err = h.Phone.Number.Normalize(); err != nil {
			return err
		}
	}

	if h.Password != nil {
		if err := h.Password.Validate(hasher); err != nil {
			return err
		}
	}
	return nil
}

func (p *Password) Validate(hasher *crypto.PasswordHasher) error {
	if p.EncodedPasswordHash != nil {
		if !hasher.EncodingSupported(*p.EncodedPasswordHash) {
			return errors.ThrowInvalidArgument(nil, "USER-JDk4t", "Errors.User.Password.NotSupported")
		}
	}
	if p.Password == nil && p.EncodedPasswordHash == nil {
		return errors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Empty")
	}
	if p.Password != nil && p.EncodedPasswordHash != nil {
		return errors.ThrowInvalidArgument(nil, "COMMAND-3M0fsss", "Errors.User.Password.NotSupported")
	}
	return nil
}

func (c *Commands) AddUserHuman(ctx context.Context, resourceOwner string, human *AddHuman, allowInitMail bool, alg crypto.EncryptionAlgorithm) (err error) {
	if resourceOwner == "" {
		return errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}

	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	if human.ID == "" {
		human.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	existingHuman, err := c.userHumanWriteModel(
		ctx,
		human.ID,
		resourceOwner,
		true, // user profile is always necessary at creation
		human.Email.Address != "",
		human.Phone.Number != "",
		human.Password != "" || human.EncodedPasswordHash != "",
		true,  // state is always necessary at creation
		false, // avatar not updateable
	)
	if err != nil {
		return err
	}
	if isUserStateExists(existingHuman.UserState) {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting")
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, resourceOwner)
	if err != nil {
		return err
	}

	if err = c.userValidateDomain(ctx, resourceOwner, human.Username, domainPolicy.UserLoginMustBeDomain); err != nil {
		return err
	}

	var createCmd humanCreationCommand
	if human.Register {
		createCmd = user.NewHumanRegisteredEvent(
			ctx,
			UserAggregateFromWriteModel(&existingHuman.WriteModel),
			human.Username,
			human.FirstName,
			human.LastName,
			human.NickName,
			human.DisplayName,
			human.PreferredLanguage,
			human.Gender,
			human.Email.Address,
			domainPolicy.UserLoginMustBeDomain,
		)
	} else {
		createCmd = user.NewHumanAddedEvent(
			ctx,
			UserAggregateFromWriteModel(&existingHuman.WriteModel),
			human.Username,
			human.FirstName,
			human.LastName,
			human.NickName,
			human.DisplayName,
			human.PreferredLanguage,
			human.Gender,
			human.Email.Address,
			domainPolicy.UserLoginMustBeDomain,
		)
	}

	if human.Phone.Number != "" {
		createCmd.AddPhoneData(human.Phone.Number)
	}

	filter := c.eventstore.Filter
	if err := addHumanCommandPassword(ctx, filter, createCmd, human, c.userPasswordHasher); err != nil {
		return err
	}

	cmds := make([]eventstore.Command, 0, 3)
	cmds = append(cmds, createCmd)

	cmds, err = c.addHumanCommandEmail(ctx, filter, cmds, existingHuman.Aggregate(), human, alg, allowInitMail)
	if err != nil {
		return err
	}

	cmds, err = c.addHumanCommandPhone(ctx, filter, cmds, existingHuman.Aggregate(), human, alg)
	if err != nil {
		return err
	}

	for _, metadataEntry := range human.Metadata {
		cmds = append(cmds, user.NewMetadataSetEvent(
			ctx,
			UserAggregateFromWriteModel(&existingHuman.WriteModel),
			metadataEntry.Key,
			metadataEntry.Value,
		))
	}
	for _, link := range human.Links {
		cmd, err := addLink(ctx, filter, existingHuman.Aggregate(), link)
		if err != nil {
			return err
		}
		cmds = append(cmds, cmd)
	}

	if len(cmds) == 0 {
		human.Details = writeModelToObjectDetails(&existingHuman.WriteModel)
		return nil
	}

	err = c.pushAppendAndReduce(ctx, existingHuman, cmds...)
	if err != nil {
		return err
	}
	human.Details = writeModelToObjectDetails(&existingHuman.WriteModel)
	return nil
}

func (c *Commands) LockUserHuman(ctx context.Context, resourceOwner string, userID string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userHumanStateWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, errors.ThrowNotFound(nil, "COMMAND-k2unb", "Errors.User.NotFound")
	}
	if !hasUserState(existingHuman.UserState, domain.UserStateActive, domain.UserStateInitial) {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-3NN8v", "Errors.User.ShouldBeActiveOrInitial")
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserLockedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) UnlockUserHuman(ctx context.Context, resourceOwner string, userID string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userHumanStateWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, errors.ThrowNotFound(nil, "COMMAND-k2unb", "Errors.User.NotFound")
	}
	if !hasUserState(existingHuman.UserState, domain.UserStateLocked) {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-4M0ds", "Errors.User.NotLocked")
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserUnlockedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) DeactivateUserHuman(ctx context.Context, resourceOwner string, userID string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userHumanStateWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, errors.ThrowNotFound(nil, "COMMAND-k2unb", "Errors.User.NotFound")
	}
	if isUserStateInitial(existingHuman.UserState) {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-ke0fw", "Errors.User.CantDeactivateInitial")
	}
	if isUserStateInactive(existingHuman.UserState) {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-5M0sf", "Errors.User.AlreadyInactive")
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserDeactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) ReactivateUserHuman(ctx context.Context, resourceOwner string, userID string) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	if userID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-2M0sd", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userHumanStateWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, errors.ThrowNotFound(nil, "COMMAND-k2unb", "Errors.User.NotFound")
	}
	if !isUserStateInactive(existingHuman.UserState) {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-6M0sf", "Errors.User.NotInactive")
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewUserReactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

func (c *Commands) ChangeUserHuman(ctx context.Context, resourceOwner string, human *ChangeHuman, alg crypto.EncryptionAlgorithm) (err error) {
	if resourceOwner == "" {
		return errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	existingHuman, err := c.userHumanWriteModel(
		ctx,
		human.ID,
		resourceOwner,
		human.Profile != nil,
		human.Email != nil,
		human.Phone != nil,
		human.Password != nil,
		true,  // state is always necessary at update
		false, // avatar not updateable
	)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return errors.ThrowNotFound(nil, "COMMAND-k2unb", "Errors.User.NotFound")
	}

	cmds := make([]eventstore.Command, 0)
	// if user change has not all necessary elements check permission
	if (human.Email != nil && human.Email.Verified) ||
		(human.Phone != nil && human.Phone.Verified) ||
		(human.Password != nil && human.Password.PasswordCode == nil && human.Password.OldPassword == nil) {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
			return err
		}
	}

	if human.Username != nil {
		cmds, err = c.changeUsername(ctx, cmds, existingHuman, *human.Username)
		if err != nil {
			return err
		}
	}
	if human.Profile != nil {
		cmds, err = changeUserProfile(ctx, cmds, existingHuman, human.Profile)
		if err != nil {
			return err
		}
	}
	if human.Email != nil {
		cmds, human.EmailCode, err = c.changeUserEmail(ctx, cmds, existingHuman, human.Email, alg)
		if err != nil {
			return err
		}
	}
	if human.Phone != nil {
		cmds, human.PhoneCode, err = c.changeUserPhone(ctx, cmds, existingHuman, human.Phone, alg)
		if err != nil {
			return err
		}
	}
	if human.Password != nil {
		cmds, err = c.changeUserPassword(ctx, cmds, existingHuman, human.Password, alg)
		if err != nil {
			return err
		}
	}

	if len(cmds) == 0 {
		human.Details = writeModelToObjectDetails(&existingHuman.WriteModel)
		return nil
	}
	err = c.pushAppendAndReduce(ctx, existingHuman, cmds...)
	if err != nil {
		return err
	}
	human.Details = writeModelToObjectDetails(&existingHuman.WriteModel)
	return nil
}

func (c *Commands) changeUserEmail(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, email *Email, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, code *string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if email.Address != "" {
		cmd := wm.NewEmailAddressChangedEvent(ctx, email.Address)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if email.Verified {
		cmd := wm.NewEmailIsVerifiedEvent(ctx, email.Verified)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else {
		// if email not changed or empty, send no code
		if email.Address == "" || email.Address == wm.Email {
			return cmds, code, err
		}
		c, err := c.newEmailCode(ctx, c.eventstore.Filter, alg)
		if err != nil {
			return cmds, code, err
		}
		cmds = append(cmds, user.NewHumanEmailCodeAddedEventV2(ctx, &wm.Aggregate().Aggregate, c.Crypted, c.Expiry, email.URLTemplate, email.ReturnCode))
		if email.ReturnCode {
			code = &c.Plain
		}
	}
	return cmds, code, nil
}

func (c *Commands) changeUserPhone(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, phone *Phone, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, code *string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if phone.Number != "" {
		cmd := wm.NewPhoneNumberChangedEvent(ctx, phone.Number)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if phone.Verified {
		cmd := wm.NewPhoneIsVerifiedEvent(ctx, phone.Verified)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else {
		// if phone not changed or empty, send no code
		if phone.Number == "" || phone.Number == wm.Phone {
			return cmds, code, err
		}
		c, err := c.newPhoneCode(ctx, c.eventstore.Filter, alg)
		if err != nil {
			return cmds, code, err
		}
		cmds = append(cmds, user.NewHumanPhoneCodeAddedEventV2(ctx, &wm.Aggregate().Aggregate, c.Crypted, c.Expiry, phone.ReturnCode))
		if phone.ReturnCode {
			code = &c.Plain
		}
	}
	return cmds, code, nil
}

func changeUserProfile(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, profile *Profile) ([]eventstore.Command, error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	cmd, err := wm.NewProfileChangedEvent(ctx, profile.FirstName, profile.LastName, profile.NickName, profile.DisplayName, profile.PreferredLanguage, profile.Gender)
	if cmd != nil {
		return append(cmds, cmd), err
	}
	return cmds, err
}

func (c *Commands) changeUserPassword(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, password *Password, alg crypto.EncryptionAlgorithm) ([]eventstore.Command, error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	// Either have a code to set the password
	if password.PasswordCode != nil {
		if err := crypto.VerifyCodeWithAlgorithm(wm.PasswordCodeCreationDate, wm.PasswordCodeExpiry, wm.PasswordCode, *password.PasswordCode, alg); err != nil {
			return cmds, err
		}
	}
	// or have the old password to change it
	if password.OldPassword != nil {
		if _, err := c.verifyPassword(ctx, wm.PasswordEncodedHash, *password.OldPassword); err != nil {
			return cmds, err
		}
	}

	if password.EncodedPasswordHash != nil {
		cmd, err := c.setPasswordCommand(ctx, &wm.Aggregate().Aggregate, wm.UserState, *password.EncodedPasswordHash, password.ChangeRequired, true)
		if cmd != nil {
			return append(cmds, cmd), err
		}
		return cmds, err
	}
	if password.Password != nil {
		cmd, err := c.setPasswordCommand(ctx, &wm.Aggregate().Aggregate, wm.UserState, *password.Password, password.ChangeRequired, false)
		if cmd != nil {
			return append(cmds, cmd), err
		}
		return cmds, err
	}
	return cmds, nil
}

func (c *Commands) userHumanWriteModel(ctx context.Context, userID, resourceOwner string, profileWM, emailWM, phoneWM, passwordWM, stateWM, avatarWM bool) (writeModel *UserHumanWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserHumanWriteModel(userID, resourceOwner, profileWM, emailWM, phoneWM, passwordWM, stateWM, avatarWM)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) userHumanStateWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *UserHumanWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserHumanStateWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) orgDomainVerifiedWriteModel(ctx context.Context, domain string) (writeModel *OrgDomainVerifiedWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewOrgDomainVerifiedWriteModel(domain)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
