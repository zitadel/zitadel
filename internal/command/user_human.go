package command

import (
	"context"
	"strings"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) getHuman(ctx context.Context, userID, resourceowner string) (*domain.Human, error) {
	human, err := c.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(human.UserState) {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-M9dsd", "Errors.User.NotFound")
	}
	return writeModelToHuman(human), nil
}

type AddHuman struct {
	// ID is optional, if empty it will be generated
	ID string
	// Username is required
	Username string
	// FirstName is required
	FirstName string
	// LastName is required
	LastName string
	// NickName is required
	NickName string
	// DisplayName is required
	DisplayName string
	// Email is required
	Email Email
	// PreferredLanguage is required
	PreferredLanguage language.Tag
	// Gender is required
	Gender domain.Gender
	// Phone represents an international phone number
	Phone Phone
	// Password is optional
	Password string
	// EncodedPasswordHash is optional
	EncodedPasswordHash string
	// PasswordChangeRequired is used if the `Password`-field is set
	PasswordChangeRequired bool
	Passwordless           bool
	ExternalIDP            bool
	Register               bool
	// UserAgentID is optional and can be passed in case the user registered themselves.
	// This will be used in the login UI to handle authentication automatically.
	UserAgentID string
	// AuthRequestID is optional and can be passed in case the user registered themselves.
	// This will be used to pass the information in notifications for links to the login UI.
	AuthRequestID string
	Metadata      []*AddMetadataEntry

	// Links are optional
	Links []*AddLink

	// TOTPSecret is optional
	TOTPSecret string

	// Details are set after a successful execution of the command
	Details *domain.ObjectDetails

	// EmailCode is set by the command
	EmailCode *string

	// PhoneCode is set by the command
	PhoneCode *string
}

type AddLink struct {
	IDPID         string
	DisplayName   string
	IDPExternalID string
}

func (h *AddHuman) Validate(hasher *crypto.Hasher) (err error) {
	if err := h.Email.Validate(); err != nil {
		return err
	}
	if h.Username = strings.TrimSpace(h.Username); h.Username == "" {
		return zerrors.ThrowInvalidArgument(nil, "V2-zzad3", "Errors.Invalid.Argument")
	}

	if h.FirstName = strings.TrimSpace(h.FirstName); h.FirstName == "" {
		return zerrors.ThrowInvalidArgument(nil, "USER-UCej2", "Errors.User.Profile.FirstNameEmpty")
	}
	if h.LastName = strings.TrimSpace(h.LastName); h.LastName == "" {
		return zerrors.ThrowInvalidArgument(nil, "USER-4hB7d", "Errors.User.Profile.LastNameEmpty")
	}
	h.ensureDisplayName()

	if h.Phone.Number != "" {
		if h.Phone.Number, err = h.Phone.Number.Normalize(); err != nil {
			return err
		}
	}

	for _, metadataEntry := range h.Metadata {
		if err := metadataEntry.Valid(); err != nil {
			return err
		}
	}
	if h.EncodedPasswordHash != "" {
		if !hasher.EncodingSupported(h.EncodedPasswordHash) {
			return zerrors.ThrowInvalidArgument(nil, "USER-JDk4t", "Errors.User.Password.NotSupported")
		}
	}
	return nil
}

type AddMetadataEntry struct {
	Key   string
	Value []byte
}

func (m *AddMetadataEntry) Valid() error {
	if m.Key = strings.TrimSpace(m.Key); m.Key == "" {
		return zerrors.ThrowInvalidArgument(nil, "USER-Drght", "Errors.User.Metadata.KeyEmpty")
	}
	if len(m.Value) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "USER-Dbgth", "Errors.User.Metadata.ValueEmpty")
	}
	return nil
}

// Deprecated: use commands.AddUserHuman
func (c *Commands) AddHuman(ctx context.Context, resourceOwner string, human *AddHuman, allowInitMail bool) (err error) {
	if resourceOwner == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter,
		c.AddHumanCommand(
			human,
			resourceOwner,
			c.userPasswordHasher,
			c.userEncryption,
			allowInitMail,
		))
	if err != nil {
		return err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return err
	}
	human.Details = &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}

	return nil
}

type humanCreationCommand interface {
	eventstore.Command
	AddPhoneData(phoneNumber domain.PhoneNumber)
	AddPasswordData(encoded string, changeRequired bool)
}

//nolint:gocognit
func (c *Commands) AddHumanCommand(human *AddHuman, orgID string, hasher *crypto.Hasher, codeAlg crypto.EncryptionAlgorithm, allowInitMail bool) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := human.Validate(hasher); err != nil {
			return nil, err
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if err := c.addHumanCommandCheckID(ctx, filter, human, orgID); err != nil {
				return nil, err
			}
			a := user.NewAggregate(human.ID, orgID)

			domainPolicy, err := domainPolicyWriteModel(ctx, filter, a.ResourceOwner)
			if err != nil {
				return nil, err
			}

			if err = c.userValidateDomain(ctx, a.ResourceOwner, human.Username, domainPolicy.UserLoginMustBeDomain); err != nil {
				return nil, err
			}

			var createCmd humanCreationCommand
			if human.Register {
				createCmd = user.NewHumanRegisteredEvent(
					ctx,
					&a.Aggregate,
					human.Username,
					human.FirstName,
					human.LastName,
					human.NickName,
					human.DisplayName,
					human.PreferredLanguage,
					human.Gender,
					human.Email.Address,
					domainPolicy.UserLoginMustBeDomain,
					"", // no user agent id available
				)
			} else {
				createCmd = user.NewHumanAddedEvent(
					ctx,
					&a.Aggregate,
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

			if err := addHumanCommandPassword(ctx, filter, createCmd, human, hasher); err != nil {
				return nil, err
			}

			cmds := make([]eventstore.Command, 0, 3)
			cmds = append(cmds, createCmd)

			cmds, err = c.addHumanCommandEmail(ctx, filter, cmds, a, human, codeAlg, allowInitMail)
			if err != nil {
				return nil, err
			}

			cmds, err = c.addHumanCommandPhone(ctx, filter, cmds, a, human, codeAlg)
			if err != nil {
				return nil, err
			}

			for _, metadataEntry := range human.Metadata {
				cmds = append(cmds, user.NewMetadataSetEvent(
					ctx,
					&a.Aggregate,
					metadataEntry.Key,
					metadataEntry.Value,
				))
			}
			for _, link := range human.Links {
				cmd, err := addLink(ctx, filter, a, link)
				if err != nil {
					return nil, err
				}
				cmds = append(cmds, cmd)
			}

			return cmds, nil
		}, nil
	}
}

func (c *Commands) addHumanCommandEmail(ctx context.Context, filter preparation.FilterToQueryReducer, cmds []eventstore.Command, a *user.Aggregate, human *AddHuman, codeAlg crypto.EncryptionAlgorithm, allowInitMail bool) ([]eventstore.Command, error) {
	if human.Email.Verified {
		cmds = append(cmds, user.NewHumanEmailVerifiedEvent(ctx, &a.Aggregate))
	}
	// if allowInitMail, used for v1 api (system, admin, mgmt, auth):
	// add init code if
	// email not verified or
	// user not registered and password set
	if allowInitMail && human.shouldAddInitCode() {
		initCode, err := c.newUserInitCode(ctx, filter, codeAlg)
		if err != nil {
			return nil, err
		}
		return append(cmds, user.NewHumanInitialCodeAddedEvent(ctx, &a.Aggregate, initCode.Crypted, initCode.Expiry, human.AuthRequestID)), nil
	}
	if !human.Email.Verified {
		emailCode, err := c.newEmailCode(ctx, filter, codeAlg)
		if err != nil {
			return nil, err
		}
		if human.Email.ReturnCode {
			human.EmailCode = &emailCode.Plain
		}
		return append(cmds, user.NewHumanEmailCodeAddedEventV2(ctx, &a.Aggregate, emailCode.Crypted, emailCode.Expiry, human.Email.URLTemplate, human.Email.ReturnCode, human.AuthRequestID)), nil
	}
	return cmds, nil
}

func addLink(ctx context.Context, filter preparation.FilterToQueryReducer, a *user.Aggregate, link *AddLink) (eventstore.Command, error) {
	exists, err := ExistsIDPOnOrgOrInstance(ctx, filter, authz.GetInstance(ctx).InstanceID(), a.ResourceOwner, link.IDPID)
	if !exists || err != nil {
		return nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-39nf2", "Errors.IDPConfig.NotExisting")
	}
	return user.NewUserIDPLinkAddedEvent(ctx, &a.Aggregate, link.IDPID, link.DisplayName, link.IDPExternalID), nil
}

func (c *Commands) addHumanCommandPhone(ctx context.Context, filter preparation.FilterToQueryReducer, cmds []eventstore.Command, a *user.Aggregate, human *AddHuman, codeAlg crypto.EncryptionAlgorithm) ([]eventstore.Command, error) {
	if human.Phone.Number == "" {
		return cmds, nil
	}
	if human.Phone.Verified {
		return append(cmds, user.NewHumanPhoneVerifiedEvent(ctx, &a.Aggregate)), nil
	}
	phoneCode, err := c.newPhoneCode(ctx, filter, codeAlg)
	if err != nil {
		return nil, err
	}
	if human.Phone.ReturnCode {
		human.PhoneCode = &phoneCode.Plain
	}
	return append(cmds, user.NewHumanPhoneCodeAddedEventV2(ctx, &a.Aggregate, phoneCode.Crypted, phoneCode.Expiry, human.Phone.ReturnCode)), nil
}

// Deprecated: use commands.NewUserHumanWriteModel, to remove deprecated eventstore.Filter
func (c *Commands) addHumanCommandCheckID(ctx context.Context, filter preparation.FilterToQueryReducer, human *AddHuman, orgID string) (err error) {
	if human.ID == "" {
		human.ID, err = id_generator.Next()
		if err != nil {
			return err
		}
	}
	existingHuman, err := humanWriteModelByID(ctx, filter, human.ID, orgID)
	if err != nil {
		return err
	}
	if isUserStateExists(existingHuman.UserState) {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting")
	}
	return nil
}

func addHumanCommandPassword(ctx context.Context, filter preparation.FilterToQueryReducer, createCmd humanCreationCommand, human *AddHuman, hasher *crypto.Hasher) (err error) {
	if human.Password != "" {
		if err = humanValidatePassword(ctx, filter, human.Password); err != nil {
			return err
		}

		_, spanHash := tracing.NewNamedSpan(ctx, "passwap.Hash")
		encodedHash, err := hasher.Hash(human.Password)
		spanHash.EndWithError(err)
		if err != nil {
			return err
		}
		createCmd.AddPasswordData(encodedHash, human.PasswordChangeRequired)
		return nil
	}

	if human.EncodedPasswordHash != "" {
		createCmd.AddPasswordData(human.EncodedPasswordHash, human.PasswordChangeRequired)
	}
	return nil
}

func (c *Commands) userValidateDomain(ctx context.Context, resourceOwner string, username string, mustBeDomain bool) (err error) {
	if mustBeDomain {
		return nil
	}

	index := strings.LastIndex(username, "@")
	if index < 0 {
		return nil
	}

	domainCheck, err := c.searchOrgDomainVerifiedByDomain(ctx, username[index+1:])
	if err != nil {
		return err
	}

	if domainCheck.Verified && domainCheck.OrgID != resourceOwner {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername")
	}

	return nil
}

func humanValidatePassword(ctx context.Context, filter preparation.FilterToQueryReducer, password string) error {
	passwordComplexity, err := passwordComplexityPolicyWriteModel(ctx, filter)
	if err != nil {
		return err
	}

	return passwordComplexity.Validate(password)
}

func (h *AddHuman) ensureDisplayName() {
	if strings.TrimSpace(h.DisplayName) != "" {
		return
	}
	if strings.TrimSpace(h.FirstName) != "" && strings.TrimSpace(h.LastName) != "" {
		h.DisplayName = h.FirstName + " " + h.LastName
		return
	}
	if strings.TrimSpace(string(h.Email.Address)) != "" {
		h.DisplayName = string(h.Email.Address)
		return
	}
	h.DisplayName = h.Username
}

// shouldAddInitCode returns true for all added Humans which:
// - were not added from an external IDP
// - and either:
//   - have no verified email
//     and / or
//   - have no authentication method (password / passwordless)
func (h *AddHuman) shouldAddInitCode() bool {
	return len(h.Links) == 0 &&
		(!h.Email.Verified ||
			(!h.Passwordless && h.Password == ""))
}

// Deprecated: use commands.AddUserHuman
func (c *Commands) ImportHuman(ctx context.Context, orgID string, human *domain.Human, passwordless bool, links []*domain.UserIDPLink, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessCodeGenerator crypto.Generator) (_ *domain.Human, passwordlessCode *domain.PasswordlessInitCode, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if orgID == "" {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5N8fs", "Errors.ResourceOwnerMissing")
	}
	domainPolicy, err := c.getOrgDomainPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-2N9fs", "Errors.Org.DomainPolicy.NotFound")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, zerrors.ThrowPreconditionFailed(err, "COMMAND-4N8gs", "Errors.Org.PasswordComplexityPolicy.NotFound")
	}

	if human.AggregateID != "" {
		existing, err := c.getHumanWriteModelByID(ctx, human.AggregateID, human.ResourceOwner)
		if err != nil {
			return nil, nil, err
		}

		if existing.UserState != domain.UserStateUnspecified {
			return nil, nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-ziuna", "Errors.User.AlreadyExisting")
		}
	}

	events, addedHuman, addedCode, code, err := c.importHuman(ctx, orgID, human, passwordless, links, domainPolicy, pwPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessCodeGenerator)
	if err != nil {
		return nil, nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, nil, err
	}

	err = AppendAndReduce(addedHuman, pushedEvents...)
	if err != nil {
		return nil, nil, err
	}
	if addedCode != nil {
		err = AppendAndReduce(addedCode, pushedEvents...)
		if err != nil {
			return nil, nil, err
		}
		passwordlessCode = writeModelToPasswordlessInitCode(addedCode, code)
	}

	return writeModelToHuman(addedHuman), passwordlessCode, nil
}

func (c *Commands) importHuman(ctx context.Context, orgID string, human *domain.Human, passwordless bool, links []*domain.UserIDPLink, domainPolicy *domain.DomainPolicy, pwPolicy *domain.PasswordComplexityPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessCodeGenerator crypto.Generator) (events []eventstore.Command, humanWriteModel *HumanWriteModel, passwordlessCodeWriteModel *HumanPasswordlessInitCodeWriteModel, code string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if orgID == "" {
		return nil, nil, nil, "", zerrors.ThrowInvalidArgument(nil, "COMMAND-00p2b", "Errors.Org.Empty")
	}
	if err = human.Normalize(); err != nil {
		return nil, nil, nil, "", err
	}
	events, humanWriteModel, err = c.createHuman(ctx, orgID, human, links, passwordless, domainPolicy, pwPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator)
	if err != nil {
		return nil, nil, nil, "", err
	}
	if passwordless {
		var codeEvent eventstore.Command
		codeEvent, passwordlessCodeWriteModel, code, err = c.humanAddPasswordlessInitCode(ctx, human.AggregateID, orgID, true, passwordlessCodeGenerator)
		if err != nil {
			return nil, nil, nil, "", err
		}
		events = append(events, codeEvent)
	}
	return events, humanWriteModel, passwordlessCodeWriteModel, code, nil
}

func (c *Commands) createHuman(ctx context.Context, orgID string, human *domain.Human, links []*domain.UserIDPLink, passwordless bool, domainPolicy *domain.DomainPolicy, pwPolicy *domain.PasswordComplexityPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator crypto.Generator) (events []eventstore.Command, addedHuman *HumanWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = human.CheckDomainPolicy(domainPolicy); err != nil {
		return nil, nil, err
	}
	human.Username = strings.TrimSpace(human.Username)
	human.EmailAddress = human.EmailAddress.Normalize()
	if err = c.userValidateDomain(ctx, orgID, human.Username, domainPolicy.UserLoginMustBeDomain); err != nil {
		return nil, nil, err
	}

	if human.AggregateID == "" {
		userID, err := id_generator.Next()
		if err != nil {
			return nil, nil, err
		}
		human.AggregateID = userID
	}

	human.EnsureDisplayName()
	if human.Password != nil {
		if err := human.HashPasswordIfExisting(ctx, pwPolicy, c.userPasswordHasher, human.Password.ChangeRequired); err != nil {
			return nil, nil, err
		}
	}

	addedHuman = NewHumanWriteModel(human.AggregateID, orgID)
	//TODO: adlerhurst maybe we could simplify the code below
	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)

	events = append(events, createAddHumanEvent(ctx, userAgg, human, domainPolicy.UserLoginMustBeDomain))

	for _, link := range links {
		event, err := c.addUserIDPLink(ctx, userAgg, link, false)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, event)
	}

	if human.IsInitialState(passwordless, len(links) > 0) {
		initCode, err := domain.NewInitUserCode(initCodeGenerator)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, user.NewHumanInitialCodeAddedEvent(ctx, userAgg, initCode.Code, initCode.Expiry, ""))
	} else {
		if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
			events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
		} else {
			emailCode, _, err := domain.NewEmailCode(emailCodeGenerator)
			if err != nil {
				return nil, nil, err
			}
			events = append(events, user.NewHumanEmailCodeAddedEvent(ctx, userAgg, emailCode.Code, emailCode.Expiry, ""))
		}
	}

	if human.Phone != nil && human.PhoneNumber != "" && !human.IsPhoneVerified {
		phoneCode, err := domain.NewPhoneCode(phoneCodeGenerator)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	} else if human.Phone != nil && human.PhoneNumber != "" && human.IsPhoneVerified {
		events = append(events, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
	}

	return events, addedHuman, nil
}

func (c *Commands) HumanSkipMFAInit(ctx context.Context, userID, resourceowner string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-2xpX9", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return zerrors.ThrowNotFound(nil, "COMMAND-m9cV8", "Errors.User.NotFound")
	}

	_, err = c.eventstore.Push(ctx,
		user.NewHumanMFAInitSkippedEvent(ctx, UserAggregateFromWriteModel(&existingHuman.WriteModel)))
	return err
}

// TODO: adlerhurst maybe we can simplify createAddHumanEvent and createRegisterHumanEvent
func createAddHumanEvent(ctx context.Context, aggregate *eventstore.Aggregate, human *domain.Human, userLoginMustBeDomain bool) *user.HumanAddedEvent {
	addEvent := user.NewHumanAddedEvent(
		ctx,
		aggregate,
		human.Username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
		userLoginMustBeDomain,
	)
	if human.Phone != nil {
		addEvent.AddPhoneData(human.PhoneNumber)
	}
	if human.Address != nil {
		addEvent.AddAddressData(
			human.Country,
			human.Locality,
			human.PostalCode,
			human.Region,
			human.StreetAddress)
	}
	if human.Password != nil {
		addEvent.AddPasswordData(human.Password.EncodedSecret, human.Password.ChangeRequired)
	}
	if human.HashedPassword != "" {
		addEvent.AddPasswordData(human.HashedPassword, false)
	}
	return addEvent
}

func (c *Commands) HumansSignOut(ctx context.Context, agentID string, userIDs []string) error {
	if agentID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}
	if len(userIDs) == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-M0od3", "Errors.User.UserIDMissing")
	}
	events := make([]eventstore.Command, 0)
	for _, userID := range userIDs {
		existingUser, err := c.getHumanWriteModelByID(ctx, userID, "")
		if err != nil {
			return err
		}
		if !isUserStateExists(existingUser.UserState) {
			continue
		}
		events = append(events, user.NewHumanSignedOutEvent(
			ctx,
			UserAggregateFromWriteModel(&existingUser.WriteModel),
			agentID))
	}
	if len(events) == 0 {
		return nil
	}
	_, err := c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) getHumanWriteModelByID(ctx context.Context, userID, resourceowner string) (_ *HumanWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	humanWriteModel := NewHumanWriteModel(userID, resourceowner)
	err = c.eventstore.FilterToQueryReducer(ctx, humanWriteModel)
	if err != nil {
		return nil, err
	}
	return humanWriteModel, nil
}

func humanWriteModelByID(ctx context.Context, filter preparation.FilterToQueryReducer, userID, resourceowner string) (*HumanWriteModel, error) {
	humanWriteModel := NewHumanWriteModel(userID, resourceowner)
	events, err := filter(ctx, humanWriteModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return humanWriteModel, nil
	}
	humanWriteModel.AppendEvents(events...)
	err = humanWriteModel.Reduce()
	return humanWriteModel, err
}

func AddHumanFromDomain(user *domain.Human, metadataList []*domain.Metadata, authRequest *domain.AuthRequest, idp *domain.UserIDPLink) *AddHuman {
	addMetadata := make([]*AddMetadataEntry, len(metadataList))
	for i, metadata := range metadataList {
		addMetadata[i] = &AddMetadataEntry{
			Key:   metadata.Key,
			Value: metadata.Value,
		}
	}
	human := new(AddHuman)
	if user.Profile != nil {
		human.Username = user.Username
		human.FirstName = user.FirstName
		human.LastName = user.LastName
		human.NickName = user.NickName
		human.DisplayName = user.DisplayName
		human.PreferredLanguage = user.PreferredLanguage
		human.Gender = user.Gender
		human.Register = true
		human.Metadata = addMetadata
	}
	if authRequest != nil {
		human.UserAgentID = authRequest.AgentID
		human.AuthRequestID = authRequest.ID
	}
	if user.Email != nil {
		human.Email = Email{
			Address:  user.Email.EmailAddress,
			Verified: user.Email.IsEmailVerified,
		}
	}
	if user.Phone != nil {
		human.Phone = Phone{
			Number:   user.Phone.PhoneNumber,
			Verified: user.Phone.IsPhoneVerified,
		}
	}
	if user.Password != nil {
		human.Password = user.Password.SecretString
	}
	if idp != nil {
		human.Links = []*AddLink{
			{
				IDPID:         idp.IDPConfigID,
				DisplayName:   idp.DisplayName,
				IDPExternalID: idp.ExternalUserID,
			},
		}
	}
	if human.Username = strings.TrimSpace(human.Username); human.Username == "" {
		human.Username = string(human.Email.Address)
	}
	return human
}
