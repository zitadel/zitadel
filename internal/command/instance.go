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
		UserLoginMustBeDomain bool
		ValidateOrgDomains    bool
	}
	LoginPolicy struct {
		AllowUsernamePassword      bool
		AllowRegister              bool
		AllowExternalIDP           bool
		ForceMFA                   bool
		HidePasswordReset          bool
		PasswordlessType           domain.PasswordlessType
		PasswordCheckLifetime      time.Duration
		ExternalLoginCheckLifetime time.Duration
		MfaInitSkipLifetime        time.Duration
		SecondFactorCheckLifetime  time.Duration
		MultiFactorCheckLifetime   time.Duration
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
}

type ZitadelConfig struct {
	projectID    string
	mgmtAppID    string
	adminAppID   string
	authAppID    string
	consoleAppID string
}

func (s *InstanceSetup) generateIDs() (err error) {
	s.zitadel.projectID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.mgmtAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.adminAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.authAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.zitadel.consoleAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) SetUpInstance(ctx context.Context, setup *InstanceSetup) (string, *domain.ObjectDetails, error) {
	instanceID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	if err = c.eventstore.NewInstance(ctx, instanceID); err != nil {
		return "", nil, err
	}

	ctx = authz.SetCtxData(authz.WithRequestedDomain(authz.WithInstanceID(ctx, instanceID), c.externalDomain), authz.CtxData{OrgID: instanceID, ResourceOwner: instanceID})

	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return "", nil, err
	}

	if err = setup.generateIDs(); err != nil {
		return "", nil, err
	}
	ctx = authz.WithConsole(ctx, setup.zitadel.projectID, setup.zitadel.consoleAppID)

	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	userAgg := user.NewAggregate(userID, orgID)
	projectAgg := project.NewAggregate(setup.zitadel.projectID, orgID)

	validations := []preparation.Validation{
		addInstance(instanceAgg, setup.InstanceName, setup.DefaultLanguage),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeAppSecret, setup.SecretGenerators.ClientSecret),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeInitCode, setup.SecretGenerators.InitializeUserCode),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyEmailCode, setup.SecretGenerators.EmailVerificationCode),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyPhoneCode, setup.SecretGenerators.PhoneVerificationCode),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypePasswordResetCode, setup.SecretGenerators.PasswordVerificationCode),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypePasswordlessInitCode, setup.SecretGenerators.PasswordlessInitCode),
		addSecretGeneratorConfig(instanceAgg, domain.SecretGeneratorTypeVerifyDomain, setup.SecretGenerators.DomainVerification),

		AddPasswordComplexityPolicy(
			instanceAgg,
			setup.PasswordComplexityPolicy.MinLength,
			setup.PasswordComplexityPolicy.HasLowercase,
			setup.PasswordComplexityPolicy.HasUppercase,
			setup.PasswordComplexityPolicy.HasNumber,
			setup.PasswordComplexityPolicy.HasSymbol,
		),
		AddPasswordAgePolicy(
			instanceAgg,
			setup.PasswordAgePolicy.ExpireWarnDays,
			setup.PasswordAgePolicy.MaxAgeDays,
		),
		AddDefaultDomainPolicy(
			instanceAgg,
			setup.DomainPolicy.UserLoginMustBeDomain,
			setup.DomainPolicy.ValidateOrgDomains,
		),
		AddDefaultLoginPolicy(
			instanceAgg,
			setup.LoginPolicy.AllowUsernamePassword,
			setup.LoginPolicy.AllowRegister,
			setup.LoginPolicy.AllowExternalIDP,
			setup.LoginPolicy.ForceMFA,
			setup.LoginPolicy.HidePasswordReset,
			setup.LoginPolicy.PasswordlessType,
			setup.LoginPolicy.PasswordCheckLifetime,
			setup.LoginPolicy.ExternalLoginCheckLifetime,
			setup.LoginPolicy.MfaInitSkipLifetime,
			setup.LoginPolicy.SecondFactorCheckLifetime,
			setup.LoginPolicy.MultiFactorCheckLifetime,
		),
		AddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeOTP),
		AddSecondFactorToDefaultLoginPolicy(instanceAgg, domain.SecondFactorTypeU2F),
		AddMultiFactorToDefaultLoginPolicy(instanceAgg, domain.MultiFactorTypeU2FWithPIN),

		AddPrivacyPolicy(instanceAgg, setup.PrivacyPolicy.TOSLink, setup.PrivacyPolicy.PrivacyLink, setup.PrivacyPolicy.HelpLink),
		AddDefaultLockoutPolicy(instanceAgg, setup.LockoutPolicy.MaxAttempts, setup.LockoutPolicy.ShouldShowLockoutFailure),

		AddDefaultLabelPolicy(
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
		ActivateDefaultLabelPolicy(instanceAgg),

		AddEmailTemplate(instanceAgg, setup.EmailTemplate),
	}

	for _, msg := range setup.MessageTexts {
		validations = append(validations, SetInstanceCustomTexts(instanceAgg, msg))
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

	validations = append(validations,
		AddOrgCommand(ctx, orgAgg, setup.Org.Name),
		AddHumanCommand(userAgg, &setup.Org.Human, c.userPasswordAlg, c.userEncryption),
		c.AddOrgMemberCommand(orgAgg, userID, domain.RoleOrgOwner),
		c.AddInstanceMemberCommand(instanceAgg, userID, domain.RoleIAMOwner),

		AddProjectCommand(projectAgg, zitadelProjectName, userID, false, false, false, domain.PrivateLabelingSettingUnspecified),
		SetIAMProject(instanceAgg, projectAgg.ID),

		AddAPIAppCommand(
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

		AddAPIAppCommand(
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

		AddAPIAppCommand(
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

		AddOIDCAppCommand(console, nil),
		SetIAMConsoleID(instanceAgg, &console.ClientID, &setup.zitadel.consoleAppID),
	)
	addGenerateddDomain, err := c.addGeneratedInstanceDomain(ctx, instanceAgg, setup.InstanceName)
	if err != nil {
		return "", nil, err
	}
	validations = append(validations, addGenerateddDomain...)
	if setup.CustomDomain != "" {
		validations = append(validations,
			c.addInstanceDomain(instanceAgg, setup.CustomDomain, false),
			setPrimaryInstanceDomain(instanceAgg, setup.CustomDomain),
		)
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validations...)
	if err != nil {
		return "", nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return instanceID, &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
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
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func addInstance(a *instance.Aggregate, instanceName string, defaultLanguage language.Tag) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewInstanceAddedEvent(ctx, &a.Aggregate, instanceName),
				instance.NewDefaultLanguageSetEvent(ctx, &a.Aggregate, defaultLanguage),
			}, nil
		}, nil
	}
}

//SetIAMProject defines the command to set the id of the IAM project onto the instance
func SetIAMProject(a *instance.Aggregate, projectID string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMProjectSetEvent(ctx, &a.Aggregate, projectID),
			}, nil
		}, nil
	}
}

//SetIAMConsoleID defines the command to set the clientID of the Console App onto the instance
func SetIAMConsoleID(a *instance.Aggregate, clientID, appID *string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMConsoleSetEvent(ctx, &a.Aggregate, clientID, appID),
			}, nil
		}, nil
	}
}

func (c *Commands) setGlobalOrg(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *InstanceWriteModel, orgID string) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.GlobalOrgID != "" {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-HGG24", "Errors.IAM.GlobalOrgAlreadySet")
	}
	return instance.NewGlobalOrgSetEventEvent(ctx, iamAgg, orgID), nil
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
