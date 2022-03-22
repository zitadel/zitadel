package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/ui/console"
	iam "github.com/caos/zitadel/internal/command/v2/instance"
	"github.com/caos/zitadel/internal/command/v2/org"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/command/v2/project"
	"github.com/caos/zitadel/internal/command/v2/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
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
	OrgIAMPolicy struct {
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

func (command *Command) SetUpTenant(ctx context.Context, instance *InstanceSetup) (*domain.ObjectDetails, error) {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	if err = instance.generateIDs(); err != nil {
		return nil, err
	}

	instance.Org.Human.PasswordChangeRequired = true

	iamAgg := iam_repo.NewAggregate()
	orgAgg := org_repo.NewAggregate(orgID, orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)
	projectAgg := project_repo.NewAggregate(instance.Zitadel.projectID, orgID)

	validations := []preparation.Validation{
		iam.AddPasswordComplexityPolicy(
			iamAgg,
			instance.PasswordComplexityPolicy.MinLength,
			instance.PasswordComplexityPolicy.HasLowercase,
			instance.PasswordComplexityPolicy.HasUppercase,
			instance.PasswordComplexityPolicy.HasNumber,
			instance.PasswordComplexityPolicy.HasSymbol,
		),
		iam.AddPasswordAgePolicy(
			iamAgg,
			instance.PasswordAgePolicy.ExpireWarnDays,
			instance.PasswordAgePolicy.MaxAgeDays,
		),
		iam.AddOrgIAMPolicy(
			iamAgg,
			instance.OrgIAMPolicy.UserLoginMustBeDomain,
		),
		iam.AddLoginPolicy(
			iamAgg,
			instance.LoginPolicy.AllowUsernamePassword,
			instance.LoginPolicy.AllowRegister,
			instance.LoginPolicy.AllowExternalIDP,
			instance.LoginPolicy.ForceMFA,
			instance.LoginPolicy.HidePasswordReset,
			instance.LoginPolicy.PasswordlessType,
			instance.LoginPolicy.PasswordCheckLifetime,
			instance.LoginPolicy.ExternalLoginCheckLifetime,
			instance.LoginPolicy.MfaInitSkipLifetime,
			instance.LoginPolicy.SecondFactorCheckLifetime,
			instance.LoginPolicy.MultiFactorCheckLifetime,
		),
		iam.AddSecondFactorToLoginPolicy(iamAgg, domain.SecondFactorTypeOTP),
		iam.AddSecondFactorToLoginPolicy(iamAgg, domain.SecondFactorTypeU2F),
		iam.AddMultiFactorToLoginPolicy(iamAgg, domain.MultiFactorTypeU2FWithPIN),

		iam.AddPrivacyPolicy(iamAgg, instance.PrivacyPolicy.TOSLink, instance.PrivacyPolicy.PrivacyLink),
		iam.AddLockoutPolicy(iamAgg, instance.LockoutPolicy.MaxAttempts, instance.LockoutPolicy.ShouldShowLockoutFailure),

		iam.AddEmailTemplate(iamAgg, instance.EmailTemplate),
	}

	for _, msg := range instance.MessageTexts {
		validations = append(validations, iam.SetCustomTexts(iamAgg, msg))
	}

	validations = append(validations,
		org.AddOrg(orgAgg, instance.Org.Name, command.iamDomain),
		org.AddDomain(orgAgg, instance.Org.Domain),
		user.AddHumanCommand(userAgg, &instance.Org.Human, command.userPasswordAlg),
		org.AddMember(orgAgg, userID, domain.RoleOrgOwner),

		project.AddProject(projectAgg, zitadelProjectName, userID, false, false, false, domain.PrivateLabelingSettingUnspecified),

		project.AddApp(
			projectAgg,
			instance.Zitadel.mgmtID,
			mgmtAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			instance.Zitadel.mgmtID,
			instance.Zitadel.mgmtClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			instance.Zitadel.adminID,
			adminAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			instance.Zitadel.adminID,
			instance.Zitadel.adminClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			instance.Zitadel.authID,
			authAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			instance.Zitadel.authID,
			instance.Zitadel.authClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			instance.Zitadel.consoleID,
			consoleAppName,
		),
		project.AddOIDCConfig(
			*projectAgg,
			domain.OIDCVersionV1,
			instance.Zitadel.consoleID,
			instance.Zitadel.consoleClientID,
			nil,
			[]string{instance.Zitadel.BaseURL + consoleRedirectPath},
			[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			domain.OIDCApplicationTypeUserAgent,
			domain.OIDCAuthMethodTypeNone,
			[]string{instance.Zitadel.BaseURL + consolePostLogoutPath},
			instance.Zitadel.IsDevMode,
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
