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
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	zitadel                  ZitadelConfig
	InstanceName             string
	CustomDomain             string
	DefaultLanguage          language.Tag
	Org                      InstanceOrgSetup
	SecretGenerators         *SecretGenerators
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
		ForceMFALocalOnly          bool
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
		TOSLink        string
		PrivacyLink    string
		HelpLink       string
		SupportEmail   domain.EmailAddress
		DocsLink       string
		CustomLink     string
		CustomLinkText string
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
		ThemeMode           domain.LabelPolicyThemeMode
	}
	LockoutPolicy struct {
		MaxPasswordAttempts      uint64
		MaxOTPAttempts           uint64
		ShouldShowLockoutFailure bool
	}
	EmailTemplate     []byte
	MessageTexts      []*domain.CustomMessageText
	SMTPConfiguration *smtp.Config
	OIDCSettings      *OIDCSettings
	Quotas            *SetQuotas
	Features          *InstanceFeatures
	Limits            *SetLimits
	Restrictions      *SetRestrictions
}

type OIDCSettings struct {
	AccessTokenLifetime        time.Duration
	IdTokenLifetime            time.Duration
	RefreshTokenIdleExpiration time.Duration
	RefreshTokenExpiration     time.Duration
}

type SetQuotas struct {
	Items []*SetQuota
}

type SecretGenerators struct {
	ClientSecret             *crypto.GeneratorConfig
	InitializeUserCode       *crypto.GeneratorConfig
	EmailVerificationCode    *crypto.GeneratorConfig
	PhoneVerificationCode    *crypto.GeneratorConfig
	PasswordVerificationCode *crypto.GeneratorConfig
	PasswordlessInitCode     *crypto.GeneratorConfig
	DomainVerification       *crypto.GeneratorConfig
	OTPSMS                   *crypto.GeneratorConfig
	OTPEmail                 *crypto.GeneratorConfig
}

type ZitadelConfig struct {
	instanceID     string
	orgID          string
	projectID      string
	mgmtAppID      string
	adminAppID     string
	authAppID      string
	consoleAppID   string
	limitsID       string
	restrictionsID string
}

func (s *InstanceSetup) generateIDs() (err error) {
	s.zitadel.instanceID, err = id_generator.Next()
	if err != nil {
		return err
	}

	s.zitadel.orgID, err = id_generator.Next()
	if err != nil {
		return err
	}

	s.zitadel.projectID, err = id_generator.Next()
	if err != nil {
		return err
	}

	s.zitadel.mgmtAppID, err = id_generator.Next()
	if err != nil {
		return err
	}

	s.zitadel.adminAppID, err = id_generator.Next()
	if err != nil {
		return err
	}

	s.zitadel.authAppID, err = id_generator.Next()
	if err != nil {
		return err
	}

	s.zitadel.consoleAppID, err = id_generator.Next()
	if err != nil {
		return err
	}
	s.zitadel.limitsID, err = id_generator.Next()
	if err != nil {
		return err
	}
	s.zitadel.restrictionsID, err = id_generator.Next()
	return err
}

func (c *Commands) SetUpInstance(ctx context.Context, setup *InstanceSetup) (string, string, *MachineKey, *domain.ObjectDetails, error) {
	if err := setup.generateIDs(); err != nil {
		return "", "", nil, nil, err
	}
	ctx = contextWithInstanceSetupInfo(ctx, setup.zitadel.instanceID, setup.zitadel.projectID, setup.zitadel.consoleAppID, c.externalDomain)

	validations, pat, machineKey, err := setUpInstance(ctx, c, setup)
	if err != nil {
		return "", "", nil, nil, err
	}

	//nolint:staticcheck
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

	return setup.zitadel.instanceID, token, machineKey, &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: setup.zitadel.orgID,
	}, nil
}

func contextWithInstanceSetupInfo(ctx context.Context, instanceID, projectID, consoleAppID, externalDomain string) context.Context {
	return authz.WithConsole(
		authz.SetCtxData(
			authz.WithRequestedDomain(
				authz.WithInstanceID(
					ctx,
					instanceID),
				externalDomain,
			),
			authz.CtxData{ResourceOwner: instanceID},
		),
		projectID,
		consoleAppID,
	)
}

func setUpInstance(ctx context.Context, c *Commands, setup *InstanceSetup) (validations []preparation.Validation, pat *PersonalAccessToken, machineKey *MachineKey, err error) {
	instanceAgg := instance.NewAggregate(setup.zitadel.instanceID)

	validations = setupInstanceElements(instanceAgg, setup)

	// default organization on setup'd instance
	pat, machineKey, err = setupDefaultOrg(ctx, c, &validations, instanceAgg, setup.Org.Name, setup.Org.Machine, setup.Org.Human, setup.zitadel)
	if err != nil {
		return nil, nil, nil, err
	}

	// domains
	if err := setupGeneratedDomain(ctx, c, &validations, instanceAgg, setup.InstanceName); err != nil {
		return nil, nil, nil, err
	}
	setupCustomDomain(c, &validations, instanceAgg, setup.CustomDomain)

	// optional setting if set
	setupMessageTexts(&validations, setup.MessageTexts, instanceAgg)
	if err := setupQuotas(c, &validations, setup.Quotas, setup.zitadel.instanceID); err != nil {
		return nil, nil, nil, err
	}
	setupSMTPSettings(c, &validations, setup.SMTPConfiguration, instanceAgg)
	setupOIDCSettings(c, &validations, setup.OIDCSettings, instanceAgg)
	setupFeatures(&validations, setup.Features, setup.zitadel.instanceID)
	setupLimits(c, &validations, limits.NewAggregate(setup.zitadel.limitsID, setup.zitadel.instanceID), setup.Limits)
	setupRestrictions(c, &validations, restrictions.NewAggregate(setup.zitadel.restrictionsID, setup.zitadel.instanceID, setup.zitadel.instanceID), setup.Restrictions)

	return validations, pat, machineKey, nil
}

func setupInstanceElements(instanceAgg *instance.Aggregate, setup *InstanceSetup) []preparation.Validation {
	return []preparation.Validation{
		prepareAddInstance(instanceAgg, setup.InstanceName, setup.DefaultLanguage),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeAppSecret, setup.SecretGenerators.ClientSecret),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeInitCode, setup.SecretGenerators.InitializeUserCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyEmailCode, setup.SecretGenerators.EmailVerificationCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyPhoneCode, setup.SecretGenerators.PhoneVerificationCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypePasswordResetCode, setup.SecretGenerators.PasswordVerificationCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypePasswordlessInitCode, setup.SecretGenerators.PasswordlessInitCode),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyDomain, setup.SecretGenerators.DomainVerification),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeOTPSMS, setup.SecretGenerators.OTPSMS),
		prepareAddSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeOTPEmail, setup.SecretGenerators.OTPEmail),

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
			setup.LoginPolicy.ForceMFALocalOnly,
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
		prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeTOTP),
		prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeU2F),
		/* TODO: incomment when usable
		prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeOTPEmail),
		prepareAddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeOTPSMS),
		*/
		prepareAddMultiFactorToDefaultLoginPolicy(instanceAgg, domain.MultiFactorTypeU2FWithPIN),

		prepareAddDefaultPrivacyPolicy(instanceAgg, setup.PrivacyPolicy.TOSLink, setup.PrivacyPolicy.PrivacyLink, setup.PrivacyPolicy.HelpLink, setup.PrivacyPolicy.SupportEmail, setup.PrivacyPolicy.DocsLink, setup.PrivacyPolicy.CustomLink, setup.PrivacyPolicy.CustomLinkText),
		prepareAddDefaultNotificationPolicy(instanceAgg, setup.NotificationPolicy.PasswordChange),
		prepareAddDefaultLockoutPolicy(instanceAgg, setup.LockoutPolicy.MaxPasswordAttempts, setup.LockoutPolicy.MaxOTPAttempts, setup.LockoutPolicy.ShouldShowLockoutFailure),

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
			setup.LabelPolicy.ThemeMode,
		),
		prepareAddDefaultEmailTemplate(instanceAgg, setup.EmailTemplate),
	}
}

func setupLimits(commands *Commands, validations *[]preparation.Validation, limitsAgg *limits.Aggregate, setLimits *SetLimits) {
	if setLimits != nil {
		*validations = append(*validations, commands.SetLimitsCommand(limitsAgg, &limitsWriteModel{}, setLimits))
	}
}

func setupRestrictions(commands *Commands, validations *[]preparation.Validation, restrictionsAgg *restrictions.Aggregate, setRestrictions *SetRestrictions) {
	if setRestrictions != nil {
		*validations = append(*validations, commands.SetRestrictionsCommand(restrictionsAgg, &restrictionsWriteModel{}, setRestrictions))
	}
}

func setupQuotas(commands *Commands, validations *[]preparation.Validation, setQuotas *SetQuotas, instanceID string) error {
	if setQuotas == nil {
		return nil
	}
	for _, q := range setQuotas.Items {
		quotaId, err := id_generator.Next()
		if err != nil {
			return err
		}
		*validations = append(*validations, commands.SetQuotaCommand(quota.NewAggregate(quotaId, instanceID), nil, true, q))
	}
	return nil
}

func setupFeatures(validations *[]preparation.Validation, features *InstanceFeatures, instanceID string) {
	if features != nil {
		*validations = append(*validations, prepareSetFeatures(instanceID, features))
	}
}

func setupOIDCSettings(commands *Commands, validations *[]preparation.Validation, oidcSettings *OIDCSettings, instanceAgg *instance.Aggregate) {
	if oidcSettings == nil {
		return
	}
	*validations = append(*validations,
		commands.prepareAddOIDCSettings(
			instanceAgg,
			oidcSettings.AccessTokenLifetime,
			oidcSettings.IdTokenLifetime,
			oidcSettings.RefreshTokenIdleExpiration,
			oidcSettings.RefreshTokenExpiration,
		),
	)
}

func setupSMTPSettings(commands *Commands, validations *[]preparation.Validation, smtpConfig *smtp.Config, instanceAgg *instance.Aggregate) {
	if smtpConfig == nil {
		return
	}
	*validations = append(*validations,
		commands.prepareAddAndActivateSMTPConfig(
			instanceAgg,
			smtpConfig.Description,
			smtpConfig.From,
			smtpConfig.FromName,
			smtpConfig.ReplyToAddress,
			smtpConfig.SMTP.Host,
			smtpConfig.SMTP.User,
			[]byte(smtpConfig.SMTP.Password),
			smtpConfig.Tls,
		),
	)
}

func setupCustomDomain(commands *Commands, validations *[]preparation.Validation, instanceAgg *instance.Aggregate, customDomain string) {
	if customDomain == "" {
		return
	}
	*validations = append(*validations,
		commands.addInstanceDomain(instanceAgg, customDomain, false),
		setPrimaryInstanceDomain(instanceAgg, customDomain),
	)
}

func setupGeneratedDomain(ctx context.Context, commands *Commands, validations *[]preparation.Validation, instanceAgg *instance.Aggregate, instanceName string) error {
	addGeneratedDomain, err := commands.addGeneratedInstanceDomain(ctx, instanceAgg, instanceName)
	if err != nil {
		return err
	}
	*validations = append(*validations, addGeneratedDomain...)
	return nil
}

func setupMinimalInterfaces(commands *Commands, validations *[]preparation.Validation, instanceAgg *instance.Aggregate, orgAgg *org.Aggregate, projectOwner string, ids ZitadelConfig) {
	projectAgg := project.NewAggregate(ids.projectID, orgAgg.ID)

	cnsl := &addOIDCApp{
		AddApp: AddApp{
			Aggregate: *projectAgg,
			ID:        ids.consoleAppID,
			Name:      consoleAppName,
		},
		Version:                  domain.OIDCVersionV1,
		RedirectUris:             []string{},
		ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
		GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
		ApplicationType:          domain.OIDCApplicationTypeUserAgent,
		AuthMethodType:           domain.OIDCAuthMethodTypeNone,
		PostLogoutRedirectUris:   []string{},
		DevMode:                  !commands.externalSecure,
		AccessTokenType:          domain.OIDCTokenTypeBearer,
		AccessTokenRoleAssertion: false,
		IDTokenRoleAssertion:     false,
		IDTokenUserinfoAssertion: false,
		ClockSkew:                0,
	}

	*validations = append(*validations,
		AddProjectCommand(projectAgg, zitadelProjectName, projectOwner, false, false, false, domain.PrivateLabelingSettingUnspecified),
		SetIAMProject(instanceAgg, projectAgg.ID),

		commands.AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        ids.mgmtAppID,
					Name:      mgmtAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
		),

		commands.AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        ids.adminAppID,
					Name:      adminAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
		),

		commands.AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        ids.authAppID,
					Name:      authAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
		),

		commands.AddOIDCAppCommand(cnsl),
		SetIAMConsoleID(instanceAgg, &cnsl.ClientID, &ids.consoleAppID),
	)
}

func setupDefaultOrg(ctx context.Context,
	commands *Commands,
	validations *[]preparation.Validation,
	instanceAgg *instance.Aggregate,
	name string,
	machine *AddMachine,
	human *AddHuman,
	ids ZitadelConfig,
) (pat *PersonalAccessToken, machineKey *MachineKey, err error) {
	orgAgg := org.NewAggregate(ids.orgID)

	*validations = append(
		*validations,
		AddOrgCommand(ctx, orgAgg, name),
		commands.prepareSetDefaultOrg(instanceAgg, ids.orgID),
	)

	projectOwner, pat, machineKey, err := setupAdmins(commands, validations, instanceAgg, orgAgg, machine, human)
	if err != nil {
		return nil, nil, err
	}
	setupMinimalInterfaces(commands, validations, instanceAgg, orgAgg, projectOwner, ids)
	return pat, machineKey, nil
}

func setupAdmins(commands *Commands,
	validations *[]preparation.Validation,
	instanceAgg *instance.Aggregate,
	orgAgg *org.Aggregate,
	machine *AddMachine,
	human *AddHuman,
) (owner string, pat *PersonalAccessToken, machineKey *MachineKey, err error) {
	if human == nil && machine == nil {
		return "", nil, nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-z1yi2q2ot7", "Error.Instance.NoAdmin")
	}

	if machine != nil && machine.Machine != nil && !machine.Machine.IsZero() {
		machineUserID, err := id_generator.Next()
		if err != nil {
			return "", nil, nil, err
		}
		owner = machineUserID

		pat, machineKey, err = setupMachineAdmin(commands, validations, machine, orgAgg.ID, machineUserID)
		if err != nil {
			return "", nil, nil, err
		}

		setupAdminMembers(commands, validations, instanceAgg, orgAgg, machineUserID)
	}
	if human != nil {
		humanUserID, err := id_generator.Next()
		if err != nil {
			return "", nil, nil, err
		}
		owner = humanUserID
		human.ID = humanUserID

		*validations = append(*validations,
			commands.AddHumanCommand(human, orgAgg.ID, commands.userPasswordHasher, commands.userEncryption, true),
		)

		setupAdminMembers(commands, validations, instanceAgg, orgAgg, humanUserID)
	}
	return owner, pat, machineKey, nil
}

func setupMachineAdmin(commands *Commands, validations *[]preparation.Validation, machine *AddMachine, orgID, userID string) (pat *PersonalAccessToken, machineKey *MachineKey, err error) {
	*validations = append(*validations,
		AddMachineCommand(user.NewAggregate(userID, orgID), machine.Machine),
	)
	if machine.Pat != nil {
		pat = NewPersonalAccessToken(orgID, userID, machine.Pat.ExpirationDate, machine.Pat.Scopes, domain.UserTypeMachine)
		pat.TokenID, err = id_generator.Next()
		if err != nil {
			return nil, nil, err
		}
		*validations = append(*validations, prepareAddPersonalAccessToken(pat, commands.keyAlgorithm))
	}
	if machine.MachineKey != nil {
		machineKey = NewMachineKey(orgID, userID, machine.MachineKey.ExpirationDate, machine.MachineKey.Type)
		machineKey.KeyID, err = id_generator.Next()
		if err != nil {
			return nil, nil, err
		}
		*validations = append(*validations, prepareAddUserMachineKey(machineKey, commands.machineKeySize))
	}
	return pat, machineKey, nil
}

func setupAdminMembers(commands *Commands, validations *[]preparation.Validation, instanceAgg *instance.Aggregate, orgAgg *org.Aggregate, userID string) {
	*validations = append(*validations,
		commands.AddOrgMemberCommand(orgAgg, userID, domain.RoleOrgOwner),
		commands.AddInstanceMemberCommand(instanceAgg, userID, domain.RoleIAMOwner),
	)
}

func setupMessageTexts(validations *[]preparation.Validation, setupMessageTexts []*domain.CustomMessageText, instanceAgg *instance.Aggregate) {
	for _, msg := range setupMessageTexts {
		*validations = append(*validations, prepareSetInstanceCustomMessageTexts(instanceAgg, msg))
	}
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
		if len(cmds) == 0 {
			continue
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
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SWffe", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getInstanceWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if writeModel.DefaultOrgID == orgID {
				return nil, zerrors.ThrowPreconditionFailed(nil, "INST-SDfw2", "Errors.Instance.NotChanged")
			}
			if exists, err := ExistsOrg(ctx, filter, orgID); err != nil || !exists {
				return nil, zerrors.ThrowPreconditionFailed(err, "INSTA-Wfe21", "Errors.Org.NotFound")
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
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-EGbw2", "Errors.IAM.IAMProjectAlreadySet")
	}
	return instance.NewIAMProjectSetEvent(ctx, iamAgg, projectID), nil
}

func (c *Commands) prepareUpdateInstance(a *instance.Aggregate, name string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-092mid", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getInstanceWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INST-nuso2m", "Errors.Instance.NotFound")
			}
			if writeModel.Name == name {
				return nil, zerrors.ThrowPreconditionFailed(nil, "INST-alpxism", "Errors.Instance.NotChanged")
			}
			return []eventstore.Command{instance.NewInstanceChangedEvent(ctx, &a.Aggregate, name)}, nil
		}, nil
	}
}

func (c *Commands) prepareSetDefaultLanguage(a *instance.Aggregate, defaultLanguage language.Tag) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if err := domain.LanguageIsDefined(defaultLanguage); err != nil {
			return nil, err
		}
		if err := domain.LanguagesAreSupported(i18n.SupportedLanguages(), defaultLanguage); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := getInstanceWriteModel(ctx, filter)
			if writeModel.DefaultLanguage == defaultLanguage {
				return nil, zerrors.ThrowPreconditionFailed(nil, "INST-DS3rq", "Errors.Instance.NotChanged")
			}
			instanceID := authz.GetInstance(ctx).InstanceID()
			restrictionsWM, err := c.getRestrictionsWriteModel(ctx, instanceID, instanceID)
			if err != nil {
				return nil, err
			}
			if err := domain.LanguageIsAllowed(false, restrictionsWM.allowedLanguages, defaultLanguage); err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
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
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareRemoveInstance(a *instance.Aggregate) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getInstanceWriteModelByID(ctx, a.ID)
			if err != nil {
				return nil, zerrors.ThrowNotFound(err, "COMMA-pax9m3", "Errors.Instance.NotFound")
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(err, "COMMA-AE3GS", "Errors.Instance.NotFound")
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
