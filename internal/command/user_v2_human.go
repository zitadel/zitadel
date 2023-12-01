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
	if h.Email.Address != "" {
		if err := h.Email.Validate(); err != nil {
			return err
		}
	}

	if h.Phone.Number != "" {
		if h.Phone.Number, err = h.Phone.Number.Normalize(); err != nil {
			return err
		}
	}

	if h.Password != nil {
		if h.Password.Validate(hasher); err != nil {
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
	if p.Password == nil {
		return errors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Empty")
	}
	if p.OldPassword == nil && p.PasswordCode == nil {
		return errors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Empty")
	}

	return nil
}

func (c *Commands) AddUserHuman(ctx context.Context, resourceOwner string, human *AddHuman, allowInitMail bool) (err error) {
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

	if err := c.addHumanCommandPassword(ctx, createCmd, human, c.userPasswordHasher); err != nil {
		return err
	}

	cmds := make([]eventstore.Command, 0, 3)
	cmds = append(cmds, createCmd)
	filter := c.eventstore.Filter

	cmds, err = c.addHumanCommandEmail(ctx, filter, cmds, existingHuman.Aggregate(), human, c.userEncryption, allowInitMail)
	if err != nil {
		return err
	}

	cmds, err = c.addHumanCommandPhone(ctx, filter, cmds, existingHuman.Aggregate(), human, c.userEncryption)
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

	err = c.pushAppendAndReduce(ctx, existingHuman, cmds...)
	if err != nil {
		return err
	}
	human.Details = writeModelToObjectDetails(&existingHuman.WriteModel)
	return nil
}

func (c *Commands) ChangeUserHuman(ctx context.Context, resourceOwner string, human *ChangeHuman) error {
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
	)
	if err != nil {
		return err
	}
	if isUserStateExists(existingHuman.UserState) {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting")
	}

	cmds := make([]eventstore.Command, 0, 4)
	if human.Username != nil {
		if err := c.changeUsername(ctx, cmds, existingHuman, *human.Username); err != nil {
			return err
		}
	}
	if human.Profile != nil {
		changeUserProfile(ctx, cmds, existingHuman, human.Profile)
	}
	if human.Email != nil {
		changeUserEmail(ctx, cmds, existingHuman, human.Email)
	}
	if human.Phone != nil {
		changeUserPhone(ctx, cmds, existingHuman, human.Phone)
	}
	if human.Password != nil {
		// password changes can only be handled if user is active
		if existingHuman.UserState == domain.UserStateInitial {
			return errors.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.NotInitialised")
		}
		c.changeUserPassword(ctx, cmds, existingHuman, human.Password)
	}

	err = c.pushAppendAndReduce(ctx, existingHuman, cmds...)
	if err != nil {
		return err
	}
	human.Details = writeModelToObjectDetails(&existingHuman.WriteModel)
	return nil
}

func changeUserEmail(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, email *Email) {
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
		//TODO email code generate
	}
}

func changeUserPhone(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, phone *Phone) {
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
		//TODO phone code generate
	}
}

func changeUserProfile(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, profile *Profile) error {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	cmd, err := wm.NewProfileChangedEvent(ctx, profile.FirstName, profile.LastName, profile.NickName, profile.DisplayName, profile.PreferredLanguage, profile.Gender)
	if err != nil {
		return err
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return nil
}

func (c *Commands) changeUserPassword(ctx context.Context, cmds []eventstore.Command, wm *UserHumanWriteModel, password *Password) error {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	// Either have a code to set the password
	if password.PasswordCode != nil {
		if err := crypto.VerifyCodeWithAlgorithm(wm.PasswordCodeCreationDate, wm.PasswordCodeExpiry, wm.PasswordCode, *password.PasswordCode, c.userEncryption); err != nil {
			return err
		}
	}
	// or have the old password to change it
	if password.OldPassword != nil {
		if _, err := c.verifyPassword(ctx, wm.PasswordEncodedHash, *password.OldPassword); err != nil {
			return err
		}
	}
	// or if neither, have the permission to do so
	if password.PasswordCode == nil && password.OldPassword == nil {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, wm.ResourceOwner, wm.AggregateID); err != nil {
			return err
		}
	}

	cmd, err := c.setPasswordCommand(ctx, &wm.Aggregate().Aggregate, wm.UserState, *password.Password, password.ChangeRequired)
	if err != nil {
		return err
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return nil
}

func (c *Commands) userHumanWriteModel(ctx context.Context, userID, resourceOwner string, profileWM, emailWM, phoneWM, passwordWM bool) (writeModel *UserHumanWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserHumanWriteModel(userID, resourceOwner, profileWM, emailWM, phoneWM, passwordWM)
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
