package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/ui/console"
	"github.com/caos/zitadel/internal/command/v2/instance"
	"github.com/caos/zitadel/internal/command/v2/org"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/command/v2/project"
	"github.com/caos/zitadel/internal/command/v2/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
	instance_repo "github.com/caos/zitadel/internal/repository/instance"
	org_repo "github.com/caos/zitadel/internal/repository/org"
	project_repo "github.com/caos/zitadel/internal/repository/project"
	user_repo "github.com/caos/zitadel/internal/repository/user"
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
	Org                      OrgSetup
	Zitadel                  ZitadelConfig
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

	projectID       string
	mgmtID          string
	mgmtClientID    string
	adminID         string
	adminClientID   string
	authID          string
	authClientID    string
	consoleID       string
	consoleClientID string
}

func (s *InstanceSetup) generateIDs() (err error) {
	s.Zitadel.projectID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	s.Zitadel.mgmtID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.mgmtClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}

	s.Zitadel.adminID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.adminClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}

	s.Zitadel.authID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.authClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}

	s.Zitadel.consoleID, err = id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}
	s.Zitadel.consoleClientID, err = domain.NewClientID(id.SonyFlakeGenerator, zitadelProjectName)
	if err != nil {
		return err
	}
	return nil
}

func (command *Command) SetUpInstance(ctx context.Context, setup *InstanceSetup) (*domain.ObjectDetails, error) {
	// TODO
	// instanceID, err := id.SonyFlakeGenerator.Next()
	// if err != nil {
	// 	return nil, err
	// }
	ctx = authz.SetCtxData(authz.WithInstance(ctx, authz.Instance{ID: "system"}), authz.CtxData{OrgID: domain.IAMID, ResourceOwner: domain.IAMID})

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

	instanceAgg := instance_repo.NewAggregate()
	orgAgg := org_repo.NewAggregate(orgID, orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)
	projectAgg := project_repo.NewAggregate(setup.Zitadel.projectID, orgID)

	validations := []preparation.Validation{
		instance.AddPasswordComplexityPolicy(
			instanceAgg,
			setup.PasswordComplexityPolicy.MinLength,
			setup.PasswordComplexityPolicy.HasLowercase,
			setup.PasswordComplexityPolicy.HasUppercase,
			setup.PasswordComplexityPolicy.HasNumber,
			setup.PasswordComplexityPolicy.HasSymbol,
		),
		instance.AddPasswordAgePolicy(
			instanceAgg,
			setup.PasswordAgePolicy.ExpireWarnDays,
			setup.PasswordAgePolicy.MaxAgeDays,
		),
		instance.AddDomainPolicy(
			instanceAgg,
			setup.DomainPolicy.UserLoginMustBeDomain,
		),
		instance.AddLoginPolicy(
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
		instance.AddSecondFactorToLoginPolicy(instanceAgg, domain.SecondFactorTypeOTP),
		instance.AddSecondFactorToLoginPolicy(instanceAgg, domain.SecondFactorTypeU2F),
		instance.AddMultiFactorToLoginPolicy(instanceAgg, domain.MultiFactorTypeU2FWithPIN),

		instance.AddPrivacyPolicy(instanceAgg, setup.PrivacyPolicy.TOSLink, setup.PrivacyPolicy.PrivacyLink, setup.PrivacyPolicy.HelpLink),
		instance.AddLockoutPolicy(instanceAgg, setup.LockoutPolicy.MaxAttempts, setup.LockoutPolicy.ShouldShowLockoutFailure),

		instance.AddEmailTemplate(instanceAgg, setup.EmailTemplate),
	}

	for _, msg := range setup.MessageTexts {
		validations = append(validations, instance.SetCustomTexts(instanceAgg, msg))
	}

	validations = append(validations,
		org.AddOrg(orgAgg, setup.Org.Name, command.iamDomain),
		user.AddHumanCommand(userAgg, &setup.Org.Human, command.userPasswordAlg),
		org.AddMember(orgAgg, userID, domain.RoleOrgOwner),

		project.AddProject(projectAgg, zitadelProjectName, userID, false, false, false, domain.PrivateLabelingSettingUnspecified),

		project.AddApp(
			projectAgg,
			setup.Zitadel.mgmtID,
			mgmtAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			setup.Zitadel.mgmtID,
			setup.Zitadel.mgmtClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			setup.Zitadel.adminID,
			adminAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			setup.Zitadel.adminID,
			setup.Zitadel.adminClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			setup.Zitadel.authID,
			authAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			setup.Zitadel.authID,
			setup.Zitadel.authClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			setup.Zitadel.consoleID,
			consoleAppName,
		),
		project.AddOIDCConfig(
			*projectAgg,
			domain.OIDCVersionV1,
			setup.Zitadel.consoleID,
			setup.Zitadel.consoleClientID,
			nil,
			[]string{setup.Zitadel.BaseURL + consoleRedirectPath},
			[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			domain.OIDCApplicationTypeUserAgent,
			domain.OIDCAuthMethodTypeNone,
			[]string{setup.Zitadel.BaseURL + consolePostLogoutPath},
			setup.Zitadel.IsDevMode,
			domain.OIDCTokenTypeBearer,
			false,
			false,
			false,
			0,
			nil,
		),
	)

	cmds, err := preparation.PrepareCommands(ctx, command.es.Filter, validations...)
	if err != nil {
		return nil, err
	}

	events, err := command.es.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
}
