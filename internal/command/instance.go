package command

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/ui/console"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	zitadelProjectName    = "ZITADEL"
	mgmtAppName           = "Management-API"
	adminAppName          = "Admin-API"
	authAppName           = "Auth-API"
	consoleAppName        = "Console"
	consoleRedirectPath   = console.HandlerPrefix + "/auth/callback"
	consolePostLogoutPath = console.HandlerPrefix + "/signedout"
)

type InstanceSetup struct {
	zitadel          ZitadelConfig
	idGenerator      id.Generator
	InstanceName     string
	CustomDomain     string
	DefaultLanguage  language.Tag
	Org              OrgSetup
	SecretGenerators struct {
		PasswordSaltCost         uint
		ClientSecret             *crypto.GeneratorConfig
		InitializeUserCode       *crypto.GeneratorConfig
		EmailVerificationCode    *crypto.GeneratorConfig
		PhoneVerificationCode    *crypto.GeneratorConfig
		PasswordVerificationCode *crypto.GeneratorConfig
		PasswordlessInitCode     *crypto.GeneratorConfig
		DomainVerification       *crypto.GeneratorConfig
	}
	PasswordComplexityPolicy struct {
		MinLength    uint64
		HasLowercase bool
		HasUppercase bool
		HasNumber    bool
		HasSymbol    bool
	}
	PasswordAgePolicy struct {
		ExpireWarnDays uint64
		MaxAgeDays     uint64
	}
	DomainPolicy struct {
		UserLoginMustBeDomain                  bool
		ValidateOrgDomains                     bool
		SMTPSenderAddressMatchesInstanceDomain bool
	}
	LoginPolicy struct {
		AllowUsernamePassword      bool
		AllowRegister              bool
		AllowExternalIDP           bool
		ForceMFA                   bool
		HidePasswordReset          bool
		IgnoreUnknownUsername      bool
		AllowDomainDiscovery       bool
		DisableLoginWithEmail      bool
		DisableLoginWithPhone      bool
		PasswordlessType           domain.PasswordlessType
		DefaultRedirectURI         string
		PasswordCheckLifetime      time.Duration
		ExternalLoginCheckLifetime time.Duration
		MfaInitSkipLifetime        time.Duration
		SecondFactorCheckLifetime  time.Duration
		MultiFactorCheckLifetime   time.Duration
	}
	NotificationPolicy struct {
		PasswordChange bool
	}
	PrivacyPolicy struct {
		TOSLink     string
		PrivacyLink string
		HelpLink    string
	}
	LabelPolicy struct {
		PrimaryColor        string
		BackgroundColor     string
		WarnColor           string
		FontColor           string
		PrimaryColorDark    string
		BackgroundColorDark string
		WarnColorDark       string
		FontColorDark       string
		HideLoginNameSuffix bool
		ErrorMsgPopup       bool
		DisableWatermark    bool
	}
	LockoutPolicy struct {
		MaxAttempts              uint64
		ShouldShowLockoutFailure bool
	}
	EmailTemplate     []byte
	MessageTexts      []*domain.CustomMessageText
	SMTPConfiguration *smtp.EmailConfig
	OIDCSettings      *struct {
		AccessTokenLifetime        time.Duration
		IdTokenLifetime            time.Duration
		RefreshTokenIdleExpiration time.Duration
		RefreshTokenExpiration     time.Duration
	}
	Quotas *struct {
		Items []*AddQuota
	}
}

type ZitadelConfig struct {
	projectID    string
	mgmtAppID    string
	adminAppID   string
	authAppID    string
	consoleAppID string
}

func (s *InstanceSetup) generateIDs(idGenerator id.Generator) (err error) {
	s.zitadel.projectID, err = idGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.mgmtAppID, err = idGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.adminAppID, err = idGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.authAppID, err = idGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.consoleAppID, err = idGenerator.Next()
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) SetUpInstance(ctx context.Context, setup *InstanceSetup) (string, string, *MachineKey, *domain.ObjectDetails, error) {
	instanceID, err := c.idGenerator.Next()
	if err != nil {
		return "", "", nil, nil, err
	}

	if err = c.eventstore.NewInstance(ctx, instanceID); err != nil {
		return "", "", nil, nil, err
	}

	ctx = authz.SetCtxData(authz.WithRequestedDomain(authz.WithInstanceID(ctx, instanceID), c.externalDomain), authz.CtxData{OrgID: instanceID, ResourceOwner: instanceID})

	orgID, err := c.idGenerator.Next()
	if err != nil {
		return "", "", nil, nil, err
	}

	userID, err := c.idGenerator.Next()
	if err != nil {
		return "", "", nil, nil, err
	}

	if err = setup.generateIDs(c.idGenerator); err != nil {
		return "", "", nil, nil, err
	}
	ctx = authz.WithConsole(ctx, setup.zitadel.projectID, setup.zitadel.consoleAppID)

	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	userAgg := user.NewAggregate(userID, orgID)
	projectAgg := project.NewAggregate(setup.zitadel.projectID, orgID)

	validations := []preparation.Validation{
		prepareAddInstance(instanceAgg, setup.InstanceName, setup.DefaultLanguage),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeAppSecret, setup.SecretGenerators.ClientSecret),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeInitCode, setup.SecretGenerators.InitializeUserCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyEmailCode, setup.SecretGenerators.EmailVerificationCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyPhoneCode, setup.SecretGenerators.PhoneVerificationCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypePasswordResetCode, setup.SecretGenerators.PasswordVerificationCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypePasswordlessInitCode, setup.SecretGenerators.PasswordlessInitCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyDomain, setup.SecretGenerators.DomainVerification),

		prepareAddDefaultPasswordComplexityPolicy(
			instanceAgg,
			setup.PasswordComplexityPolicy.MinLength,
			setup.PasswordComplexityPolicy.HasLowercase,
			setup.PasswordComplexityPolicy.HasUppercase,
			setup.PasswordComplexityPolicy.HasNumber,
			setup.PasswordComplexityPolicy.HasSymbol,
		),
		prepareAddDefaultPasswordAgePolicy(
			instanceAgg,
			setup.PasswordAgePolicy.ExpireWarnDays,
			setup.PasswordAgePolicy.MaxAgeDays,
		),
		prepareAddDefaultDomainPolicy(
			instanceAgg,
			setup.DomainPolicy.UserLoginMustBeDomain,
			setup.DomainPolicy.ValidateOrgDomains,
			setup.DomainPolicy.SMTPSenderAddressMatchesInstanceDomain,
		),
		prepareAddDefaultLoginPolicy(
			instanceAgg,
			setup.LoginPolicy.AllowUsernamePassword,
			setup.LoginPolicy.AllowRegister,
			setup.LoginPolicy.AllowExternalIDP,
			setup.LoginPolicy.ForceMFA,
			setup.LoginPolicy.HidePasswordReset,
			setup.LoginPolicy.IgnoreUnknownUsername,
			setup.LoginPolicy.AllowDomainDiscovery,
			setup.LoginPolicy.DisableLoginWithEmail,
			setup.LoginPolicy.DisableLoginWithPhone,
			setup.LoginPolicy.PasswordlessType,
			setup.LoginPolicy.DefaultRedirectURI,
			setup.LoginPolicy.PasswordCheckLifetime,
			setup.LoginPolicy.ExternalLoginCheckLifetime,
			setup.LoginPolicy.MfaInitSkipLifetime,
			setup.LoginPolicy.SecondFactorCheckLifetime,
			setup.LoginPolicy.MultiFactorCheckLifetime,
		),
		prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeOTP),
		prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeU2F),
		prepareAddMultiFactorToDefaultLoginPolicy(instanceAgg, domain.MultiFactorTypeU2FWithPIN),

		prepareAddDefaultPrivacyPolicy(instanceAgg, setup.PrivacyPolicy.TOSLink, setup.PrivacyPolicy.PrivacyLink, setup.PrivacyPolicy.HelpLink),
		prepareAddDefaultNotificationPolicy(instanceAgg, setup.NotificationPolicy.PasswordChange),
		prepareAddDefaultLockoutPolicy(instanceAgg, setup.LockoutPolicy.MaxAttempts, setup.LockoutPolicy.ShouldShowLockoutFailure),

		prepareAddDefaultLabelPolicy(
			instanceAgg,
			setup.LabelPolicy.PrimaryColor,
			setup.LabelPolicy.BackgroundColor,
			setup.LabelPolicy.WarnColor,
			setup.LabelPolicy.FontColor,
			setup.LabelPolicy.PrimaryColorDark,
			setup.LabelPolicy.BackgroundColorDark,
			setup.LabelPolicy.WarnColorDark,
			setup.LabelPolicy.FontColorDark,
			setup.LabelPolicy.HideLoginNameSuffix,
			setup.LabelPolicy.ErrorMsgPopup,
			setup.LabelPolicy.DisableWatermark,
		),
		prepareActivateDefaultLabelPolicy(instanceAgg),

		prepareAddDefaultEmailTemplate(instanceAgg, setup.EmailTemplate),
	}

	if setup.Quotas != nil {
		for _, q := range setup.Quotas.Items {
			quotaId, err := c.idGenerator.Next()
			if err != nil {
				return "", "", nil, nil, err
			}

			quotaAggregate := quota.NewAggregate(quotaId, instanceID, instanceID)

			validations = append(validations, c.AddQuotaCommand(quotaAggregate, q))
		}
	}

	for _, msg := range setup.MessageTexts {
		validations = append(validations, prepareSetInstanceCustomMessageTexts(instanceAgg, msg))
	}

	console := &addOIDCApp{
		AddApp: AddApp{
			Aggregate: *projectAgg,
			ID:        setup.zitadel.consoleAppID,
			Name:      consoleAppName,
		},
		Version:                  domain.OIDCVersionV1,
		RedirectUris:             []string{},
		ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
		GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
		ApplicationType:          domain.OIDCApplicationTypeUserAgent,
		AuthMethodType:           domain.OIDCAuthMethodTypeNone,
		PostLogoutRedirectUris:   []string{},
		DevMode:                  !c.externalSecure,
		AccessTokenType:          domain.OIDCTokenTypeBearer,
		AccessTokenRoleAssertion: false,
		IDTokenRoleAssertion:     false,
		IDTokenUserinfoAssertion: false,
		ClockSkew:                0,
	}

	validations = append(validations,
		AddOrgCommand(ctx, orgAgg, setup.Org.Name),
		c.prepareSetDefaultOrg(instanceAgg, orgAgg.ID),
	)

	var pat *PersonalAccessToken
	var machineKey *MachineKey
	// only a human or a machine user should be created as owner
	if setup.Org.Machine != nil && setup.Org.Machine.Machine != nil && !setup.Org.Machine.Machine.IsZero() {
		validations = append(validations,
			AddMachineCommand(userAgg, setup.Org.Machine.Machine),
		)
		if setup.Org.Machine.Pat != nil {
			pat = NewPersonalAccessToken(orgID, userID, setup.Org.Machine.Pat.ExpirationDate, setup.Org.Machine.Pat.Scopes, domain.UserTypeMachine)
			pat.TokenID, err = c.idGenerator.Next()
			if err != nil {
				return "", "", nil, nil, err
			}
			validations = append(validations, prepareAddPersonalAccessToken(pat, c.keyAlgorithm))
		}
		if setup.Org.Machine.MachineKey != nil {
			machineKey = NewMachineKey(orgID, userID, setup.Org.Machine.MachineKey.ExpirationDate, setup.Org.Machine.MachineKey.Type)
			machineKey.KeyID, err = c.idGenerator.Next()
			if err != nil {
				return "", "", nil, nil, err
			}
			validations = append(validations, prepareAddUserMachineKey(machineKey, c.machineKeySize))
		}
	} else if setup.Org.Human != nil {
		validations = append(validations,
			AddHumanCommand(userAgg, setup.Org.Human, c.userPasswordAlg, c.userEncryption),
		)
	}

	validations = append(validations,
		c.AddOrgMemberCommand(orgAgg, userID, domain.RoleOrgOwner),
		c.AddInstanceMemberCommand(instanceAgg, userID, domain.RoleIAMOwner),
		AddProjectCommand(projectAgg, zitadelProjectName, userID, false, false, false, domain.PrivateLabelingSettingUnspecified),
		SetIAMProject(instanceAgg, projectAgg.ID),

		c.AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        setup.zitadel.mgmtAppID,
					Name:      mgmtAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
			nil,
		),

		c.AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        setup.zitadel.adminAppID,
					Name:      adminAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
			nil,
		),

		c.AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        setup.zitadel.authAppID,
					Name:      authAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
			nil,
		),

		c.AddOIDCAppCommand(console, nil),
		SetIAMConsoleID(instanceAgg, &console.ClientID, &setup.zitadel.consoleAppID),
	)

	addGeneratedDomain, err := c.addGeneratedInstanceDomain(ctx, instanceAgg, setup.InstanceName)
	if err != nil {
		return "", "", nil, nil, err
	}
	validations = append(validations, addGeneratedDomain...)
	if setup.CustomDomain != "" {
		validations = append(validations,
			c.addInstanceDomain(instanceAgg, setup.CustomDomain, false),
			setPrimaryInstanceDomain(instanceAgg, setup.CustomDomain),
		)
	}

	if setup.SMTPConfiguration != nil {
		validations = append(validations,
			c.prepareAddSMTPConfig(
				instanceAgg,
				setup.SMTPConfiguration.From,
				setup.SMTPConfiguration.FromName,
				setup.SMTPConfiguration.SMTP.Host,
				setup.SMTPConfiguration.SMTP.User,
				[]byte(setup.SMTPConfiguration.SMTP.Password),
				setup.SMTPConfiguration.Tls,
			),
		)
	}

	if setup.OIDCSettings != nil {
		validations = append(validations,
			c.prepareAddOIDCSettings(
				instanceAgg,
				setup.OIDCSettings.AccessTokenLifetime,
				setup.OIDCSettings.IdTokenLifetime,
				setup.OIDCSettings.RefreshTokenIdleExpiration,
				setup.OIDCSettings.RefreshTokenExpiration,
			),
		)
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validations...)
	if err != nil {
		return "", "", nil, nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", "", nil, nil, err
	}

	var token string
	if pat != nil {
		token = pat.Token
	}

	return instanceID, token, machineKey, &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
}

func (c *Commands) UpdateInstance(ctx context.Context, name string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareUpdateInstance(instanceAgg, name)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) SetDefaultLanguage(ctx context.Context, defaultLanguage language.Tag) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareSetDefaultLanguage(instanceAgg, defaultLanguage)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) SetDefaultOrg(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareSetDefaultOrg(instanceAgg, orgID)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) ChangeSystemConfig(ctx context.Context, externalDomain string, externalPort uint16, externalSecure bool) error {
	validations, err := c.prepareChangeSystemConfig(externalDomain, externalPort, externalSecure)(ctx, c.eventstore.Filter)
	if err != nil {
		return err
	}
	for instanceID, instanceValidations := range validations {
		if len(instanceValidations.Validations) == 0 {
			continue
		}
		ctx := authz.WithConsole(authz.WithInstanceID(ctx, instanceID), instanceValidations.ProjectID, instanceValidations.ConsoleAppID)
		cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, instanceValidations.Validations...)
		if err != nil {
			return err
		}
		_, err = c.eventstore.Push(ctx, cmds...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) prepareChangeSystemConfig(externalDomain string, externalPort uint16, externalSecure bool) func(ctx context.Context, filter preparation.FilterToQueryReducer) (map[string]*SystemConfigChangesValidation, error) {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer) (map[string]*SystemConfigChangesValidation, error) {
		if externalDomain == "" || externalPort == 0 {
			return nil, nil
		}
		writeModel, err := getSystemConfigWriteModel(ctx, filter, externalDomain, c.externalDomain, externalPort, c.externalPort, externalSecure, c.externalSecure)
		if err != nil {
			return nil, err
		}
		return writeModel.NewChangedEvents(c), nil
	}
}

func prepareAddInstance(a *instance.Aggregate, instanceName string, defaultLanguage language.Tag) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewInstanceAddedEvent(ctx, &a.Aggregate, instanceName),
				instance.NewDefaultLanguageSetEvent(ctx, &a.Aggregate, defaultLanguage),
			}, nil
		}, nil
	}
}

// SetIAMProject defines the command to set the id of the IAM project onto the instance
func SetIAMProject(a *instance.Aggregate, projectID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMProjectSetEvent(ctx, &a.Aggregate, projectID),
			}, nil
		}, nil
	}
}

// SetIAMConsoleID defines the command to set the clientID of the Console App onto the instance
func SetIAMConsoleID(a *instance.Aggregate, clientID, appID *string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMConsoleSetEvent(ctx, &a.Aggregate, clientID, appID),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareSetDefaultOrg(a *instance.Aggregate, orgID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if orgID == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-SWffe", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getInstanceWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.DefaultOrgID == orgID {
				return nil, errors.ThrowPreconditionFailed(nil, "INST-SDfw2", "Errors.Instance.NotChanged")
			}
			if exists, err := ExistsOrg(ctx, filter, orgID); err != nil || !exists {
				return nil, errors.ThrowPreconditionFailed(err, "INSTA-Wfe21", "Errors.Org.NotFound")
			}
			return []eventstore.Command{instance.NewDefaultOrgSetEventEvent(ctx, &a.Aggregate, orgID)}, nil
		}, nil
	}
}

func (c *Commands) setIAMProject(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *InstanceWriteModel, projectID string) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.ProjectID != "" {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-EGbw2", "Errors.IAM.IAMProjectAlreadySet")
	}
	return instance.NewIAMProjectSetEvent(ctx, iamAgg, projectID), nil
}

func (c *Commands) prepareUpdateInstance(a *instance.Aggregate, name string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "INST-092mid", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getInstanceWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, errors.ThrowNotFound(nil, "INST-nuso2m", "Errors.Instance.NotFound")
			}
			if writeModel.Name == name {
				return nil, errors.ThrowPreconditionFailed(nil, "INST-alpxism", "Errors.Instance.NotChanged")
			}
			return []eventstore.Command{instance.NewInstanceChangedEvent(ctx, &a.Aggregate, name)}, nil
		}, nil
	}
}

func (c *Commands) prepareSetDefaultLanguage(a *instance.Aggregate, defaultLanguage language.Tag) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if defaultLanguage == language.Und {
			return nil, errors.ThrowInvalidArgument(nil, "INST-28nlD", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getInstanceWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.DefaultLanguage == defaultLanguage {
				return nil, errors.ThrowPreconditionFailed(nil, "INST-DS3rq", "Errors.Instance.NotChanged")
			}
			return []eventstore.Command{instance.NewDefaultLanguageSetEvent(ctx, &a.Aggregate, defaultLanguage)}, nil
		}, nil
	}
}

func getInstanceWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (*InstanceWriteModel, error) {
	writeModel := NewInstanceWriteModel(authz.GetInstance(ctx).InstanceID())
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}

func getSystemConfigWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, externalDomain, newExternalDomain string, externalPort, newExternalPort uint16, externalSecure, newExternalSecure bool) (*SystemConfigWriteModel, error) {
	writeModel := NewSystemConfigWriteModel(externalDomain, newExternalDomain, externalPort, newExternalPort, externalSecure, newExternalSecure)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}

func (c *Commands) RemoveInstance(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareRemoveInstance(instanceAgg))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareRemoveInstance(a *instance.Aggregate) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getInstanceWriteModelByID(ctx, a.ID)
			if err != nil {
				return nil, errors.ThrowNotFound(err, "COMMA-pax9m3", "Errors.Instance.NotFound")
			}
			if !writeModel.State.Exists() {
				return nil, errors.ThrowNotFound(err, "COMMA-AE3GS", "Errors.Instance.NotFound")
			}
			return []eventstore.Command{instance.NewInstanceRemovedEvent(ctx,
					&a.Aggregate,
					writeModel.Name,
					writeModel.Domains)},
				nil
		}, nil
	}
}

func (c *Commands) getInstanceWriteModelByID(ctx context.Context, orgID string) (*InstanceWriteModel, error) {
	instanceWriteModel := NewInstanceWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, instanceWriteModel)
	if err != nil {
		return nil, err
	}
	return instanceWriteModel, nil
}
