package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeHuman struct {
	ID            string
	ResourceOwner string
	State         *domain.UserState
	Username      *string
	Profile       *Profile
	Email         *Email
	Phone         *Phone

	Metadata             []*domain.Metadata
	MetadataKeysToRemove []string

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

	if h.Phone != nil {
		if err := h.Phone.Validate(); err != nil {
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
	if h.State != nil {
		return true
	}
	if len(h.Metadata) > 0 {
		return true
	}
	if len(h.MetadataKeysToRemove) > 0 {
		return true
	}
	return false
}

func (c *Commands) AddUserHuman(ctx context.Context, resourceOwner string, human *AddHuman, allowInitMail bool, alg crypto.EncryptionAlgorithm) (err error) {
	if resourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-095xh8fll1", "Errors.Internal")
	}
	if human.Details == nil {
		human.Details = &domain.ObjectDetails{}
	}
	human.Details.ResourceOwner = resourceOwner
	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	if human.ID == "" {
		human.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}
	// check for permission to create user on resourceOwner
	if !human.Register {
		if err := c.checkPermissionUpdateUser(ctx, resourceOwner, human.ID); err != nil {
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

	cmds, err := c.addUserHumanCommands(ctx, filter, existingHuman, human, allowInitMail, alg, createCmd)
	if err != nil {
		return err
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

func (c *Commands) addUserHumanCommands(ctx context.Context, filter preparation.FilterToQueryReducer, existingHuman *UserV2WriteModel, human *AddHuman, allowInitMail bool, alg crypto.EncryptionAlgorithm, addUserCommand eventstore.Command) ([]eventstore.Command, error) {
	cmds := []eventstore.Command{addUserCommand}
	var err error
	cmds, err = c.addHumanCommandEmail(ctx, filter, cmds, existingHuman.Aggregate(), human, alg, allowInitMail)
	if err != nil {
		return nil, err
	}

	cmds, err = c.addHumanCommandPhone(ctx, filter, cmds, existingHuman.Aggregate(), human, alg)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		cmds = append(cmds, cmd)
	}

	if human.TOTPSecret != "" {
		encryptedSecret, err := crypto.Encrypt([]byte(human.TOTPSecret), c.multifactors.OTP.CryptoMFA)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds,
			user.NewHumanOTPAddedEvent(ctx, &existingHuman.Aggregate().Aggregate, encryptedSecret),
			user.NewHumanOTPVerifiedEvent(ctx, &existingHuman.Aggregate().Aggregate, ""),
		)
	}

	if human.SetInactive {
		cmds = append(cmds, user.NewUserDeactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate))
	}
	return cmds, nil
}

func (c *Commands) ChangeUserHuman(ctx context.Context, human *ChangeHuman, alg crypto.EncryptionAlgorithm) (err error) {
	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	existingHuman, err := c.UserHumanWriteModel(
		ctx,
		human.ID,
		human.ResourceOwner,
		human.Profile != nil,
		human.Email != nil,
		human.Phone != nil,
		human.Password != nil,
		false, // avatar not updateable
		false, // IDPLinks not updateable
		len(human.Metadata) > 0 || len(human.MetadataKeysToRemove) > 0,
	)
	if err != nil {
		return err
	}

	if human.Changed() {
		if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
			return err
		}
	}

	userAgg := UserAggregateFromWriteModelCtx(ctx, &existingHuman.WriteModel)
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

	for _, md := range human.Metadata {
		cmd, err := c.setUserMetadata(ctx, userAgg, md)
		if err != nil {
			return err
		}

		cmds = append(cmds, cmd)
	}

	for _, mdKey := range human.MetadataKeysToRemove {
		cmd, err := c.removeUserMetadata(ctx, userAgg, mdKey)
		if err != nil {
			return err
		}

		cmds = append(cmds, cmd)
	}

	if human.State != nil {
		// only allow toggling between active and inactive
		// any other target state is not supported
		switch {
		case isUserStateActive(*human.State):
			if isUserStateActive(existingHuman.UserState) {
				// user is already active => no change needed
				break
			}

			// do not allow switching from other states than active (e.g. locked)
			if !isUserStateInactive(existingHuman.UserState) {
				return zerrors.ThrowInvalidArgumentf(nil, "USER2-statex1", "Errors.User.State.Invalid")
			}

			cmds = append(cmds, user.NewUserReactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate))
		case isUserStateInactive(*human.State):
			if isUserStateInactive(existingHuman.UserState) {
				// user is already inactive => no change needed
				break
			}

			// do not allow switching from other states than active (e.g. locked)
			if !isUserStateActive(existingHuman.UserState) {
				return zerrors.ThrowInvalidArgumentf(nil, "USER2-statex2", "Errors.User.State.Invalid")
			}

			cmds = append(cmds, user.NewUserDeactivatedEvent(ctx, &existingHuman.Aggregate().Aggregate))
		default:
			return zerrors.ThrowInvalidArgumentf(nil, "USER2-statex3", "Errors.User.State.Invalid")
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

	if phone.Remove {
		return append(cmds, user.NewHumanPhoneRemovedEvent(ctx, &wm.Aggregate().Aggregate)), nil, nil
	}

	if phone.Number != "" && phone.Number != wm.Phone {
		cmds = append(cmds, user.NewHumanPhoneChangedEvent(ctx, &wm.Aggregate().Aggregate, phone.Number))

		if phone.Verified {
			return append(cmds, user.NewHumanPhoneVerifiedEvent(ctx, &wm.Aggregate().Aggregate)), code, nil
		} else {
			cryptoCode, generatorID, err := c.newPhoneCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypeVerifyPhoneCode, alg, c.defaultSecretGenerators.PhoneVerificationCode) //nolint:staticcheck
			if err != nil {
				return cmds, code, err
			}
			cmds = append(cmds, user.NewHumanPhoneCodeAddedEventV2(ctx, &wm.Aggregate().Aggregate, cryptoCode.CryptedCode(), cryptoCode.CodeExpiry(), phone.ReturnCode, generatorID))
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
		verification = c.setPasswordWithVerifyCode(
			wm.PasswordCodeCreationDate,
			wm.PasswordCodeExpiry,
			wm.PasswordCode,
			wm.PasswordCodeGeneratorID,
			wm.PasswordCodeVerificationID,
			password.PasswordCode,
		)
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

func (c *Commands) HumanMFAInitSkippedV2(ctx context.Context, userID string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Wei5kooz1i", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userStateWriteModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-auj6jeBei4", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUser(ctx, existingHuman.ResourceOwner, existingHuman.AggregateID); err != nil {
		return nil, err
	}

	if err := c.pushAppendAndReduce(ctx, existingHuman, user.NewHumanMFAInitSkippedEvent(ctx, &existingHuman.Aggregate().Aggregate)); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
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

func (c *Commands) UserHumanWriteModel(ctx context.Context, userID, resourceOwner string, profileWM, emailWM, phoneWM, passwordWM, avatarWM, idpLinksWM, metadataWM bool) (writeModel *UserV2WriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserHumanWriteModel(userID, resourceOwner, profileWM, emailWM, phoneWM, passwordWM, avatarWM, idpLinksWM, metadataWM)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if !isUserStateExists(writeModel.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-ugjs0upun6", "Errors.User.NotFound")
	}

	return writeModel, nil
}
