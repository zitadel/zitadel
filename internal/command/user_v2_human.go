package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	PasswordCode        string
	OldPassword         string
	Password            string
	EncodedPasswordHash string

	ChangeRequired bool
}

func (h *ChangeHuman) Validate(hasher *crypto.Hasher) (err error) {
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

func (p *Password) Validate(hasher *crypto.Hasher) error {
	if p.EncodedPasswordHash != "" {
		if !hasher.EncodingSupported(p.EncodedPasswordHash) {
			return zerrors.ThrowInvalidArgument(nil, "USER-oz74onzvqr", "Errors.User.Password.NotSupported")
		}
	}
	if p.Password == "" && p.EncodedPasswordHash == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty")
	}
	return nil
}

func (h *ChangeHuman) Changed() bool {
	if h.Username != nil {
		return true
	}
	if h.Profile != nil {
		return true
	}
	if h.Email != nil {
		return true
	}
	if h.Phone != nil {
		return true
	}
	if h.Password != nil {
		return true
	}
	return false
}

func (c *Commands) AddUserHuman(ctx context.Context, resourceOwner string, human *AddHuman, allowInitMail bool, alg crypto.EncryptionAlgorithm) (err error) {
	if resourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-095xh8fll1", "Errors.Internal")
	}

	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	if human.ID == "" {
		human.ID, err = id_generator.Next()
		if err != nil {
			return err
		}
	}

	// only check if user is already existing
	existingHuman, err := c.userExistsWriteModel(
		ctx,
		human.ID,
	)
	if err != nil {
		return err
	}
	if isUserStateExists(existingHuman.UserState) {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-7yiox1isql", "Errors.User.AlreadyExisting")
	}
	// check for permission to create user on resourceOwner
	if !human.Register {
		if err := c.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, human.ID); err != nil {
			return err
		}
	}
	// add resourceowner for the events with the aggregate
	existingHuman.ResourceOwner = resourceOwner

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
			&existingHuman.Aggregate().Aggregate,
			human.Username,
			human.FirstName,
			human.LastName,
			human.NickName,
			human.DisplayName,
			human.PreferredLanguage,
			human.Gender,
			human.Email.Address,
			domainPolicy.UserLoginMustBeDomain,
			human.UserAgentID,
		)
	} else {
		createCmd = user.NewHumanAddedEvent(
			ctx,
			&existingHuman.Aggregate().Aggregate,
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

	// separated to change when old user logic is not used anymore
	filter := c.eventstore.Filter //nolint:staticcheck
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
			&existingHuman.Aggregate().Aggregate,
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

	if human.TOTPSecret != "" {
		encryptedSecret, err := crypto.Encrypt([]byte(human.TOTPSecret), c.multifactors.OTP.CryptoMFA)
		if err != nil {
			return err
		}
		cmds = append(cmds,
			user.NewHumanOTPAddedEvent(ctx, &existingHuman.Aggregate().Aggregate, encryptedSecret),
			user.NewHumanOTPVerifiedEvent(ctx, &existingHuman.Aggregate().Aggregate, ""),
		)
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

func (c *Commands) ChangeUserHuman(ctx context.Context, human *ChangeHuman, alg crypto.EncryptionAlgorithm) (err error) {
	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	existingHuman, err := c.userHumanWriteModel(
		ctx,
		human.ID,
		human.Profile != nil,
		human.Email != nil,
		human.Phone != nil,
		human.Password != nil,
		false, // avatar not updateable
		false, // IDPLinks not updateable
	)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return zerrors.ThrowNotFound(nil, "COMMAND-ugjs0upun6", "Errors.User.NotFound")
	}

	if human.Changed() {
		if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
			return err
		}
	}

	cmds := make([]eventstore.Command, 0)
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
		cmds, err = c.changeUserPassword(ctx, cmds, existingHuman, human.Password)
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

func (c *Commands) changeUserEmail(ctx context.Context, cmds []eventstore.Command, wm *UserV2WriteModel, email *Email, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, code *string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if email.Address != "" && email.Address != wm.Email {
		cmds = append(cmds, user.NewHumanEmailChangedEvent(ctx, &wm.Aggregate().Aggregate, email.Address))

		if email.Verified {
			return append(cmds, user.NewHumanEmailVerifiedEvent(ctx, &wm.Aggregate().Aggregate)), code, nil
		} else {
			cryptoCode, err := c.newEmailCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
			if err != nil {
				return cmds, code, err
			}
			cmds = append(cmds, user.NewHumanEmailCodeAddedEventV2(ctx, &wm.Aggregate().Aggregate, cryptoCode.Crypted, cryptoCode.Expiry, email.URLTemplate, email.ReturnCode, ""))
			if email.ReturnCode {
				code = &cryptoCode.Plain
			}
			return cmds, code, nil
		}
	}
	// only create separate event of verified if email was not changed
	if email.Verified && wm.IsEmailVerified != email.Verified {
		return append(cmds, user.NewHumanEmailVerifiedEvent(ctx, &wm.Aggregate().Aggregate)), nil, nil
	}
	return cmds, code, nil
}

func (c *Commands) changeUserPhone(ctx context.Context, cmds []eventstore.Command, wm *UserV2WriteModel, phone *Phone, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, code *string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	if phone.Number != "" && phone.Number != wm.Phone {
		cmds = append(cmds, user.NewHumanPhoneChangedEvent(ctx, &wm.Aggregate().Aggregate, phone.Number))

		if phone.Verified {
			return append(cmds, user.NewHumanPhoneVerifiedEvent(ctx, &wm.Aggregate().Aggregate)), code, nil
		} else {
			cryptoCode, err := c.newPhoneCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
			if err != nil {
				return cmds, code, err
			}
			cmds = append(cmds, user.NewHumanPhoneCodeAddedEventV2(ctx, &wm.Aggregate().Aggregate, cryptoCode.Crypted, cryptoCode.Expiry, phone.ReturnCode))
			if phone.ReturnCode {
				code = &cryptoCode.Plain
			}
			return cmds, code, nil
		}
	}
	// only create separate event of verified if email was not changed
	if phone.Verified && wm.IsPhoneVerified != phone.Verified {
		return append(cmds, user.NewHumanPhoneVerifiedEvent(ctx, &wm.Aggregate().Aggregate)), code, nil
	}
	return cmds, code, nil
}

func changeUserProfile(ctx context.Context, cmds []eventstore.Command, wm *UserV2WriteModel, profile *Profile) ([]eventstore.Command, error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	cmd, err := wm.NewProfileChangedEvent(ctx, profile.FirstName, profile.LastName, profile.NickName, profile.DisplayName, profile.PreferredLanguage, profile.Gender)
	if cmd != nil {
		return append(cmds, cmd), err
	}
	return cmds, err
}

func (c *Commands) changeUserPassword(ctx context.Context, cmds []eventstore.Command, wm *UserV2WriteModel, password *Password) ([]eventstore.Command, error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	// if no verification is set, the user must have the permission to change the password
	verification := c.setPasswordWithPermission(wm.AggregateID, wm.ResourceOwner)
	// otherwise check the password code...
	if password.PasswordCode != "" {
		verification = c.setPasswordWithVerifyCode(wm.PasswordCodeCreationDate, wm.PasswordCodeExpiry, wm.PasswordCode, password.PasswordCode)
	}
	// ...or old password
	if password.OldPassword != "" {
		verification = c.checkCurrentPassword(password.Password, password.EncodedPasswordHash, password.OldPassword, wm.PasswordEncodedHash)
	}
	cmd, err := c.setPasswordCommand(
		ctx,
		&wm.Aggregate().Aggregate,
		wm.UserState,
		password.Password,
		password.EncodedPasswordHash,
		"",
		password.ChangeRequired,
		verification,
	)
	if cmd != nil {
		return append(cmds, cmd), err
	}
	return cmds, err
}

func (c *Commands) userExistsWriteModel(ctx context.Context, userID string) (writeModel *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserExistsWriteModel(userID, "")
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) userHumanWriteModel(ctx context.Context, userID string, profileWM, emailWM, phoneWM, passwordWM, avatarWM, idpLinksWM bool) (writeModel *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserHumanWriteModel(userID, "", profileWM, emailWM, phoneWM, passwordWM, avatarWM, idpLinksWM)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
