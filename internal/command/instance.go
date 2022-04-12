package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/ui/console"
	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
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
	Org      OrgSetup
	Zitadel  ZitadelConfig
	Features struct {
		TierName                 string
		TierDescription          string
		Retention                time.Duration
		State                    domain.FeaturesState
		StateDescription         string
		LoginPolicyFactors       bool
		LoginPolicyIDP           bool
		LoginPolicyPasswordless  bool
		LoginPolicyRegistration  bool
		LoginPolicyUsernameLogin bool
		LoginPolicyPasswordReset bool
		PasswordComplexityPolicy bool
		LabelPolicyPrivateLabel  bool
		LabelPolicyWatermark     bool
		CustomDomain             bool
		PrivacyPolicy            bool
		MetadataUser             bool
		CustomTextMessage        bool
		CustomTextLogin          bool
		LockoutPolicy            bool
		ActionsAllowed           domain.ActionsAllowed
		MaxActions               int
	}
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
	EmailTemplate []byte
	MessageTexts  []*domain.CustomMessageText
}

type ZitadelConfig struct {
	IsDevMode bool
	BaseURL   string

	projectID    string
	mgmtAppID    string
	adminAppID   string
	authAppID    string
	consoleAppID string
}

func (s *InstanceSetup) generateIDs() (err error) {
	s.Zitadel.projectID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.Zitadel.mgmtAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.Zitadel.adminAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.Zitadel.authAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.Zitadel.consoleAppID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	return nil
}

func (c *commandNew) SetUpInstance(ctx context.Context, setup *InstanceSetup) (*domain.ObjectDetails, error) {
	instanceID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}
	requestedDomain := authz.GetInstance(ctx).RequestedDomain()
	ctx = authz.SetCtxData(authz.WithRequestedDomain(authz.WithInstanceID(ctx, instanceID), requestedDomain), authz.CtxData{OrgID: instanceID, ResourceOwner: instanceID})

	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	if err = setup.generateIDs(); err != nil {
		return nil, err
	}

	setup.Org.Human.PasswordChangeRequired = true

	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID, orgID)
	userAgg := user.NewAggregate(userID, orgID)
	projectAgg := project.NewAggregate(setup.Zitadel.projectID, orgID)

	validations := []preparation.Validation{
		SetDefaultFeatures(
			instanceAgg,
			setup.Features.TierName,
			setup.Features.TierDescription,
			setup.Features.State,
			setup.Features.StateDescription,
			setup.Features.Retention,
			setup.Features.LoginPolicyFactors,
			setup.Features.LoginPolicyIDP,
			setup.Features.LoginPolicyPasswordless,
			setup.Features.LoginPolicyRegistration,
			setup.Features.LoginPolicyUsernameLogin,
			setup.Features.LoginPolicyPasswordReset,
			setup.Features.PasswordComplexityPolicy,
			setup.Features.LabelPolicyPrivateLabel,
			setup.Features.LabelPolicyWatermark,
			setup.Features.CustomDomain,
			setup.Features.PrivacyPolicy,
			setup.Features.MetadataUser,
			setup.Features.CustomTextMessage,
			setup.Features.CustomTextLogin,
			setup.Features.LockoutPolicy,
			setup.Features.ActionsAllowed,
			setup.Features.MaxActions,
		),
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
			ID:        setup.Zitadel.consoleAppID,
			Name:      consoleAppName,
		},
		Version:                  domain.OIDCVersionV1,
		RedirectUris:             []string{setup.Zitadel.BaseURL + consoleRedirectPath},
		ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
		GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
		ApplicationType:          domain.OIDCApplicationTypeUserAgent,
		AuthMethodType:           domain.OIDCAuthMethodTypeNone,
		PostLogoutRedirectUris:   []string{setup.Zitadel.BaseURL + consolePostLogoutPath},
		DevMode:                  setup.Zitadel.IsDevMode,
		AccessTokenType:          domain.OIDCTokenTypeBearer,
		AccessTokenRoleAssertion: false,
		IDTokenRoleAssertion:     false,
		IDTokenUserinfoAssertion: false,
		ClockSkew:                0,
	}

	validations = append(validations,
		AddOrgCommand(ctx, orgAgg, setup.Org.Name),
		addHumanCommand(userAgg, &setup.Org.Human, c.userPasswordAlg, c.phoneAlg, c.emailAlg, c.initCodeAlg),
		c.AddOrgMember(orgAgg, userID, domain.RoleOrgOwner),
		c.AddInstanceMember(instanceAgg, userID, domain.RoleIAMOwner),

		AddProjectCommand(projectAgg, zitadelProjectName, userID, false, false, false, domain.PrivateLabelingSettingUnspecified),
		SetIAMProject(instanceAgg, projectAgg.ID),

		AddAPIAppCommand(
			&addAPIApp{
				AddApp: AddApp{
					Aggregate: *projectAgg,
					ID:        setup.Zitadel.mgmtAppID,
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
					ID:        setup.Zitadel.adminAppID,
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
					ID:        setup.Zitadel.authAppID,
					Name:      authAppName,
				},
				AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
			},
			nil,
		),

		AddOIDCAppCommand(console, nil),
		SetIAMConsoleID(instanceAgg, &console.ClientID),
	)

	cmds, err := preparation.PrepareCommands(ctx, c.es.Filter, validations...)
	if err != nil {
		return nil, err
	}

	events, err := c.es.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
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
func SetIAMConsoleID(a *instance.Aggregate, clientID *string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				instance.NewIAMConsoleSetEvent(ctx, &a.Aggregate, clientID),
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
