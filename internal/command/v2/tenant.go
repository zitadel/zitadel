package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command/v2/iam"
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
	zitadelProjectName = "ZITADEL"
	mgmtAppName        = "Management-API"
	adminAppName       = "Admin-API"
	authAppName        = "Auth-API"
	consoleAppName     = "Console"
)

type TenantSetup struct {
	Org                      OrgSetup
	Zitadel                  ZitadelSetup
	PasswordComplexityPolicy struct {
		MinLength    uint64
		HasLowercase bool
		HasUppercase bool
		HasNumber    bool
		HasSymbol    bool
	}
	PasswordAgePolicy struct {
		expireWarnDays uint64
		maxAgeDays     uint64
	}
	OrgIAMPolicy struct {
		UserLoginMustBeDomain bool
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

func (s *TenantSetup) generateIDs() (err error) {
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

type ZitadelSetup struct {
	projectID       string
	mgmtID          string
	mgmtClientID    string
	adminID         string
	adminClientID   string
	authID          string
	authClientID    string
	consoleID       string
	consoleClientID string
	IsDevMode       bool
}

func (command *Command) SetUpTenant(ctx context.Context, tenant *TenantSetup) (*domain.ObjectDetails, error) {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	if err = tenant.generateIDs(); err != nil {
		return nil, err
	}

	iamAgg := iam_repo.NewAggregate()
	orgAgg := org_repo.NewAggregate(orgID, orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)
	projectAgg := project_repo.NewAggregate(tenant.Zitadel.projectID, orgID)

	validations := []preparation.Validation{
		iam.AddPasswordComplexityPolicy(
			iamAgg,
			tenant.PasswordComplexityPolicy.MinLength,
			tenant.PasswordComplexityPolicy.HasLowercase,
			tenant.PasswordComplexityPolicy.HasUppercase,
			tenant.PasswordComplexityPolicy.HasNumber,
			tenant.PasswordComplexityPolicy.HasSymbol,
		),
		iam.AddPasswordAgePolicy(
			iamAgg,
			tenant.PasswordAgePolicy.expireWarnDays,
			tenant.PasswordAgePolicy.maxAgeDays,
		),
		iam.AddOrgIAMPolicy(
			iamAgg,
			tenant.OrgIAMPolicy.UserLoginMustBeDomain,
		),
		iam.AddLabelPolicy(
			iamAgg,
			tenant.LabelPolicy.PrimaryColor,
			tenant.LabelPolicy.BackgroundColor,
			tenant.LabelPolicy.WarnColor,
			tenant.LabelPolicy.FontColor,
			tenant.LabelPolicy.PrimaryColorDark,
			tenant.LabelPolicy.BackgroundColorDark,
			tenant.LabelPolicy.WarnColorDark,
			tenant.LabelPolicy.FontColorDark,
			tenant.LabelPolicy.HideLoginNameSuffix,
			tenant.LabelPolicy.ErrorMsgPopup,
			tenant.LabelPolicy.DisableWatermark,
		),
		iam.AddLoginPolicy(
			iamAgg,
			tenant.LoginPolicy.AllowUsernamePassword,
			tenant.LoginPolicy.AllowRegister,
			tenant.LoginPolicy.AllowExternalIDP,
			tenant.LoginPolicy.ForceMFA,
			tenant.LoginPolicy.HidePasswordReset,
			tenant.LoginPolicy.PasswordlessType,
			tenant.LoginPolicy.PasswordCheckLifetime,
			tenant.LoginPolicy.ExternalLoginCheckLifetime,
			tenant.LoginPolicy.MfaInitSkipLifetime,
			tenant.LoginPolicy.SecondFactorCheckLifetime,
			tenant.LoginPolicy.MultiFactorCheckLifetime,
		),
		iam.AddSecondFactorToLoginPolicy(iamAgg, domain.SecondFactorTypeOTP),
		iam.AddSecondFactorToLoginPolicy(iamAgg, domain.SecondFactorTypeU2F),
		iam.AddMultiFactorToLoginPolicy(iamAgg, domain.MultiFactorTypeU2FWithPIN),

		iam.AddPrivacyPolicy(iamAgg, tenant.PrivacyPolicy.TOSLink, tenant.PrivacyPolicy.PrivacyLink),
		iam.AddLockoutPolicy(iamAgg, tenant.LockoutPolicy.MaxAttempts, tenant.LockoutPolicy.ShouldShowLockoutFailure),

		iam.AddEmailTemplate(iamAgg, tenant.EmailTemplate),
	}

	for _, msg := range tenant.MessageTexts {
		validations = append(validations, iam.SetCustomTexts(iamAgg, msg))
	}

	validations = append(validations,
		org.AddOrg(orgAgg, tenant.Org.Name, command.iamDomain),
		org.AddDomain(orgAgg, tenant.Org.Domain),
		user.AddHumanCommand(userAgg, &tenant.Org.Human),
		org.AddMemberCommand(orgAgg, userID, domain.RoleOrgOwner),

		project.AddProject(projectAgg, zitadelProjectName, false, false, false, domain.PrivateLabelingSettingUnspecified),

		project.AddApp(
			projectAgg,
			tenant.Zitadel.mgmtID,
			mgmtAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			tenant.Zitadel.mgmtID,
			tenant.Zitadel.mgmtClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			tenant.Zitadel.adminID,
			adminAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			tenant.Zitadel.adminID,
			tenant.Zitadel.adminClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			tenant.Zitadel.authID,
			authAppName,
		),
		project.AddAPIConfig(
			*projectAgg,
			tenant.Zitadel.authID,
			tenant.Zitadel.authClientID,
			nil,
			domain.APIAuthMethodTypePrivateKeyJWT,
		),

		project.AddApp(
			projectAgg,
			tenant.Zitadel.consoleID,
			consoleAppName,
		),
		project.AddOIDCConfig(
			*projectAgg,
			domain.OIDCVersionV1,
			tenant.Zitadel.consoleID,
			tenant.Zitadel.consoleClientID,
			nil,
			[]string{"redirectUris"},
			[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			domain.OIDCApplicationTypeUserAgent,
			domain.OIDCAuthMethodTypeNone,
			[]string{"postLogoutRedirectUris "},
			tenant.Zitadel.IsDevMode,
			domain.OIDCTokenTypeBearer,
			false,
			false,
			false,
			0,
			[]string{"additionalOrigins"},
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
