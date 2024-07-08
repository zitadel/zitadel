package command

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func instanceSetupZitadelIDs() ZitadelConfig {
	return ZitadelConfig{
		instanceID:   "INSTANCE",
		orgID:        "ORG",
		projectID:    "PROJECT",
		consoleAppID: "console-id",
		authAppID:    "auth-id",
		mgmtAppID:    "mgmt-id",
		adminAppID:   "admin-id",
	}
}

func projectAddedEvents(ctx context.Context, instanceID, orgID, id, owner string, externalSecure bool) []eventstore.Command {
	events := []eventstore.Command{
		project.NewProjectAddedEvent(ctx,
			&project.NewAggregate(id, orgID).Aggregate,
			"ZITADEL",
			false,
			false,
			false,
			domain.PrivateLabelingSettingUnspecified,
		),
		project.NewProjectMemberAddedEvent(ctx,
			&project.NewAggregate(id, orgID).Aggregate,
			owner,
			domain.RoleProjectOwner,
		),
		instance.NewIAMProjectSetEvent(ctx,
			&instance.NewAggregate(instanceID).Aggregate,
			id,
		),
	}
	events = append(events, apiAppEvents(ctx, orgID, id, "mgmt-id", "Management-API")...)
	events = append(events, apiAppEvents(ctx, orgID, id, "admin-id", "Admin-API")...)
	events = append(events, apiAppEvents(ctx, orgID, id, "auth-id", "Auth-API")...)

	consoleAppID := "console-id"
	consoleClientID := "clientID"
	events = append(events, oidcAppEvents(ctx, orgID, id, consoleAppID, "Console", consoleClientID, externalSecure)...)
	events = append(events,
		instance.NewIAMConsoleSetEvent(ctx,
			&instance.NewAggregate(instanceID).Aggregate,
			&consoleClientID,
			&consoleAppID,
		),
	)
	return events
}

func projectClientIDs() []string {
	return []string{"clientID", "clientID", "clientID", "clientID"}
}

func apiAppEvents(ctx context.Context, orgID, projectID, id, name string) []eventstore.Command {
	return []eventstore.Command{
		project.NewApplicationAddedEvent(
			ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			id,
			name,
		),
		project.NewAPIConfigAddedEvent(ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			id,
			"clientID",
			"",
			domain.APIAuthMethodTypePrivateKeyJWT,
		),
	}
}

func oidcAppEvents(ctx context.Context, orgID, projectID, id, name, clientID string, externalSecure bool) []eventstore.Command {
	return []eventstore.Command{
		project.NewApplicationAddedEvent(
			ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			id,
			name,
		),
		project.NewOIDCConfigAddedEvent(ctx,
			&project.NewAggregate(projectID, orgID).Aggregate,
			domain.OIDCVersionV1,
			id,
			clientID,
			"",
			[]string{},
			[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			domain.OIDCApplicationTypeUserAgent,
			domain.OIDCAuthMethodTypeNone,
			[]string{},
			!externalSecure,
			domain.OIDCTokenTypeBearer,
			false,
			false,
			false,
			0,
			nil,
			false,
		),
	}
}

func orgFilters(orgID string, machine, human bool) []expect {
	filters := []expect{
		expectFilter(),
		expectFilter(
			org.NewOrgAddedEvent(context.Background(), &org.NewAggregate(orgID).Aggregate, ""),
		),
	}
	if machine {
		filters = append(filters, machineFilters(orgID, true)...)
		filters = append(filters, adminMemberFilters(orgID, "USER-MACHINE")...)
	}
	if human {
		filters = append(filters, humanFilters(orgID)...)
		filters = append(filters, adminMemberFilters(orgID, "USER")...)
	}

	return append(filters,
		projectFilters()...,
	)
}

func orgEvents(ctx context.Context, instanceID, orgID, name, projectID, defaultDomain string, externalSecure bool, machine, human bool) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	domain := strings.ToLower(name + "." + defaultDomain)
	events := []eventstore.Command{
		org.NewOrgAddedEvent(ctx, &orgAgg.Aggregate, name),
		org.NewDomainAddedEvent(ctx, &orgAgg.Aggregate, domain),
		org.NewDomainVerifiedEvent(ctx, &orgAgg.Aggregate, domain),
		org.NewDomainPrimarySetEvent(ctx, &orgAgg.Aggregate, domain),
		instance.NewDefaultOrgSetEventEvent(ctx, &instanceAgg.Aggregate, orgID),
	}

	owner := ""
	if machine {
		machineID := "USER-MACHINE"
		events = append(events, machineEvents(ctx, instanceID, orgID, machineID, "PAT")...)
		owner = machineID
	}
	if human {
		userID := "USER"
		events = append(events, humanEvents(ctx, instanceID, orgID, userID)...)
		owner = userID
	}

	events = append(events, projectAddedEvents(ctx, instanceID, orgID, projectID, owner, externalSecure)...)
	return events
}

func orgIDs() []string {
	return slices.Concat([]string{"USER-MACHINE", "PAT", "USER"}, projectClientIDs())
}

func instancePoliciesFilters(instanceID string) []expect {
	return []expect{
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
	}
}

func instancePoliciesEvents(ctx context.Context, instanceID string) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	return []eventstore.Command{
		instance.NewPasswordComplexityPolicyAddedEvent(ctx, &instanceAgg.Aggregate, 8, true, true, true, true),
		instance.NewPasswordAgePolicyAddedEvent(ctx, &instanceAgg.Aggregate, 0, 0),
		instance.NewDomainPolicyAddedEvent(ctx, &instanceAgg.Aggregate, false, false, false),
		instance.NewLoginPolicyAddedEvent(ctx, &instanceAgg.Aggregate, true, true, true, false, false, false, false, true, false, false, domain.PasswordlessTypeAllowed, "", 240*time.Hour, 240*time.Hour, 720*time.Hour, 18*time.Hour, 12*time.Hour),
		instance.NewLoginPolicySecondFactorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecondFactorTypeTOTP),
		instance.NewLoginPolicySecondFactorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecondFactorTypeU2F),
		instance.NewLoginPolicyMultiFactorAddedEvent(ctx, &instanceAgg.Aggregate, domain.MultiFactorTypeU2FWithPIN),
		instance.NewPrivacyPolicyAddedEvent(ctx, &instanceAgg.Aggregate, "", "", "", "", "", "", ""),
		instance.NewNotificationPolicyAddedEvent(ctx, &instanceAgg.Aggregate, true),
		instance.NewLockoutPolicyAddedEvent(ctx, &instanceAgg.Aggregate, 0, 0, true),
		instance.NewLabelPolicyAddedEvent(ctx, &instanceAgg.Aggregate, "#5469d4", "#fafafa", "#cd3d56", "#000000", "#2073c4", "#111827", "#ff3b5b", "#ffffff", false, false, false, domain.LabelPolicyThemeAuto),
		instance.NewLabelPolicyActivatedEvent(ctx, &instanceAgg.Aggregate),
	}
}

func instanceSetupPoliciesConfig() *InstanceSetup {
	return &InstanceSetup{
		PasswordComplexityPolicy: struct {
			MinLength    uint64
			HasLowercase bool
			HasUppercase bool
			HasNumber    bool
			HasSymbol    bool
		}{8, true, true, true, true},
		PasswordAgePolicy: struct {
			ExpireWarnDays uint64
			MaxAgeDays     uint64
		}{0, 0},
		DomainPolicy: struct {
			UserLoginMustBeDomain                  bool
			ValidateOrgDomains                     bool
			SMTPSenderAddressMatchesInstanceDomain bool
		}{false, false, false},
		LoginPolicy: struct {
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
		}{true, true, true, false, false, false, false, true, false, false, domain.PasswordlessTypeAllowed, "", 240 * time.Hour, 240 * time.Hour, 720 * time.Hour, 18 * time.Hour, 12 * time.Hour},
		NotificationPolicy: struct {
			PasswordChange bool
		}{true},
		PrivacyPolicy: struct {
			TOSLink        string
			PrivacyLink    string
			HelpLink       string
			SupportEmail   domain.EmailAddress
			DocsLink       string
			CustomLink     string
			CustomLinkText string
		}{"", "", "", "", "", "", ""},
		LabelPolicy: struct {
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
		}{"#5469d4", "#fafafa", "#cd3d56", "#000000", "#2073c4", "#111827", "#ff3b5b", "#ffffff", false, false, false, domain.LabelPolicyThemeAuto},
		LockoutPolicy: struct {
			MaxPasswordAttempts      uint64
			MaxOTPAttempts           uint64
			ShouldShowLockoutFailure bool
		}{0, 0, true},
	}
}

func setupInstanceElementsFilters(instanceID string) []expect {
	return slices.Concat(
		instanceElementsFilters(),
		instancePoliciesFilters(instanceID),
		// email template
		[]expect{expectFilter()},
	)
}

func setupInstanceElementsConfig() *InstanceSetup {
	conf := instanceSetupPoliciesConfig()
	conf.InstanceName = "ZITADEL"
	conf.DefaultLanguage = language.English
	conf.zitadel = instanceSetupZitadelIDs()
	conf.SecretGenerators = instanceElementsConfig()
	conf.EmailTemplate = []byte("something")
	return conf
}

func setupInstanceElementsEvents(ctx context.Context, instanceID, instanceName string, defaultLanguage language.Tag) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	return slices.Concat(
		instanceElementsEvents(ctx, instanceID, instanceName, defaultLanguage),
		instancePoliciesEvents(ctx, instanceID),
		[]eventstore.Command{instance.NewMailTemplateAddedEvent(ctx, &instanceAgg.Aggregate, []byte("something"))},
	)
}

func instanceElementsFilters() []expect {
	return []expect{
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
	}
}

func instanceElementsEvents(ctx context.Context, instanceID, instanceName string, defaultLanguage language.Tag) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	return []eventstore.Command{
		instance.NewInstanceAddedEvent(ctx, &instanceAgg.Aggregate, instanceName),
		instance.NewDefaultLanguageSetEvent(ctx, &instanceAgg.Aggregate, defaultLanguage),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeAppSecret, 64, 0, true, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeInitCode, 6, 72*time.Hour, false, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeVerifyEmailCode, 6, time.Hour, false, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeVerifyPhoneCode, 6, time.Hour, false, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypePasswordResetCode, 6, time.Hour, false, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypePasswordlessInitCode, 12, time.Hour, true, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeVerifyDomain, 32, 0, true, true, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeOTPSMS, 8, 5*time.Minute, false, false, true, false),
		instance.NewSecretGeneratorAddedEvent(ctx, &instanceAgg.Aggregate, domain.SecretGeneratorTypeOTPEmail, 8, 5*time.Minute, false, false, true, false),
	}
}
func instanceElementsConfig() *SecretGenerators {
	return &SecretGenerators{
		ClientSecret:             &crypto.GeneratorConfig{Length: 64, IncludeLowerLetters: true, IncludeUpperLetters: true, IncludeDigits: true},
		InitializeUserCode:       &crypto.GeneratorConfig{Length: 6, Expiry: 72 * time.Hour, IncludeUpperLetters: true, IncludeDigits: true},
		EmailVerificationCode:    &crypto.GeneratorConfig{Length: 6, Expiry: time.Hour, IncludeUpperLetters: true, IncludeDigits: true},
		PhoneVerificationCode:    &crypto.GeneratorConfig{Length: 6, Expiry: time.Hour, IncludeUpperLetters: true, IncludeDigits: true},
		PasswordVerificationCode: &crypto.GeneratorConfig{Length: 6, Expiry: time.Hour, IncludeUpperLetters: true, IncludeDigits: true},
		PasswordlessInitCode:     &crypto.GeneratorConfig{Length: 12, Expiry: time.Hour, IncludeLowerLetters: true, IncludeUpperLetters: true, IncludeDigits: true},
		DomainVerification:       &crypto.GeneratorConfig{Length: 32, IncludeLowerLetters: true, IncludeUpperLetters: true, IncludeDigits: true},
		OTPSMS:                   &crypto.GeneratorConfig{Length: 8, Expiry: 5 * time.Minute, IncludeDigits: true},
		OTPEmail:                 &crypto.GeneratorConfig{Length: 8, Expiry: 5 * time.Minute, IncludeDigits: true},
	}
}

func setupInstanceFilters(instanceID, orgID, projectID, appID, domain string) []expect {
	return slices.Concat(
		setupInstanceElementsFilters(instanceID),
		orgFilters(orgID, true, true),
		generatedDomainFilters(instanceID, orgID, projectID, appID, domain),
	)
}

func setupInstanceEvents(ctx context.Context, instanceID, orgID, projectID, appID, instanceName, orgName string, defaultLanguage language.Tag, domain string, externalSecure bool) []eventstore.Command {
	return slices.Concat(
		setupInstanceElementsEvents(ctx, instanceID, instanceName, defaultLanguage),
		orgEvents(ctx, instanceID, orgID, orgName, projectID, domain, externalSecure, true, true),
		generatedDomainEvents(ctx, instanceID, orgID, projectID, appID, domain),
	)
}

func setupInstanceConfig() *InstanceSetup {
	conf := setupInstanceElementsConfig()
	conf.Org = InstanceOrgSetup{
		Name:    "ZITADEL",
		Machine: instanceSetupMachineConfig(),
		Human:   instanceSetupHumanConfig(),
	}
	conf.CustomDomain = ""
	return conf
}

func generatedDomainEvents(ctx context.Context, instanceID, orgID, projectID, appID, defaultDomain string) []eventstore.Command {
	instanceAgg := instance.NewAggregate(instanceID)
	changed, _ := project.NewOIDCConfigChangedEvent(ctx, &project.NewAggregate(projectID, orgID).Aggregate, appID,
		[]project.OIDCConfigChanges{
			project.ChangeRedirectURIs([]string{"http://" + defaultDomain + "/ui/console/auth/callback"}),
			project.ChangePostLogoutRedirectURIs([]string{"http://" + defaultDomain + "/ui/console/signedout"}),
		},
	)
	return []eventstore.Command{
		instance.NewDomainAddedEvent(ctx, &instanceAgg.Aggregate, defaultDomain, true),
		changed,
		instance.NewDomainPrimarySetEvent(ctx, &instanceAgg.Aggregate, defaultDomain),
	}
}

func generatedDomainFilters(instanceID, orgID, projectID, appID, generatedDomain string) []expect {
	return []expect{
		expectFilter(),
		expectFilter(
			project.NewApplicationAddedEvent(context.Background(),
				&project.NewAggregate(projectID, orgID).Aggregate,
				appID,
				"console",
			),
			project.NewOIDCConfigAddedEvent(context.Background(),
				&project.NewAggregate(projectID, orgID).Aggregate,
				domain.OIDCVersionV1,
				appID,
				"clientID",
				"",
				[]string{},
				[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				domain.OIDCApplicationTypeUserAgent,
				domain.OIDCAuthMethodTypeNone,
				[]string{},
				true,
				domain.OIDCTokenTypeBearer,
				false,
				false,
				false,
				0,
				nil,
				false,
			),
		),
		expectFilter(
			func() eventstore.Event {
				event := instance.NewDomainAddedEvent(context.Background(),
					&instance.NewAggregate(instanceID).Aggregate,
					generatedDomain,
					true,
				)
				event.Data, _ = json.Marshal(event)
				return event
			}(),
		),
	}
}

func humanFilters(orgID string) []expect {
	return []expect{
		expectFilter(),
		expectFilter(
			org.NewDomainPolicyAddedEvent(
				context.Background(),
				&org.NewAggregate(orgID).Aggregate,
				true,
				true,
				true,
			),
		),
		expectFilter(
			org.NewPasswordComplexityPolicyAddedEvent(
				context.Background(),
				&org.NewAggregate(orgID).Aggregate,
				2,
				false,
				false,
				false,
				false,
			),
		),
	}
}

func instanceSetupHumanConfig() *AddHuman {
	return &AddHuman{
		Username:  "zitadel-admin",
		FirstName: "ZITADEL",
		LastName:  "Admin",
		Email: Email{
			Address:  domain.EmailAddress("admin@zitadel.test"),
			Verified: true,
		},
		PreferredLanguage:      language.English,
		Password:               "password",
		PasswordChangeRequired: false,
	}
}

func machineFilters(orgID string, pat bool) []expect {
	filters := []expect{
		expectFilter(),
		expectFilter(
			org.NewDomainPolicyAddedEvent(
				context.Background(),
				&org.NewAggregate(orgID).Aggregate,
				true,
				true,
				true,
			),
		),
	}
	if pat {
		filters = append(filters,
			expectFilter(),
			expectFilter(),
		)
	}
	return filters
}

func instanceSetupMachineConfig() *AddMachine {
	return &AddMachine{
		Machine: &Machine{
			Username:        "zitadel-admin-machine",
			Name:            "ZITADEL-machine",
			Description:     "Admin",
			AccessTokenType: domain.OIDCTokenTypeBearer,
		},
		Pat: &AddPat{
			ExpirationDate: time.Time{},
			Scopes:         nil,
		},
		/* not predictable with the key value in the events
		MachineKey: &AddMachineKey{
			Type:           domain.AuthNKeyTypeJSON,
			ExpirationDate: time.Time{},
		},
		*/
	}
}

func projectFilters() []expect {
	return []expect{
		expectFilter(),
		expectFilter(),
		expectFilter(),
		expectFilter(),
	}
}

func adminMemberFilters(orgID, userID string) []expect {
	return []expect{
		expectFilter(
			addHumanEvent(context.Background(), orgID, userID),
		),
		expectFilter(),
		expectFilter(
			addHumanEvent(context.Background(), orgID, userID),
		),
		expectFilter(),
	}
}

func humanEvents(ctx context.Context, instanceID, orgID, userID string) []eventstore.Command {
	agg := user.NewAggregate(userID, orgID)
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	return []eventstore.Command{
		addHumanEvent(ctx, orgID, userID),
		user.NewHumanEmailVerifiedEvent(ctx, &agg.Aggregate),
		org.NewMemberAddedEvent(ctx, &orgAgg.Aggregate, userID, domain.RoleOrgOwner),
		instance.NewMemberAddedEvent(ctx, &instanceAgg.Aggregate, userID, domain.RoleIAMOwner),
	}
}

func addHumanEvent(ctx context.Context, orgID, userID string) *user.HumanAddedEvent {
	agg := user.NewAggregate(userID, orgID)
	return func() *user.HumanAddedEvent {
		event := user.NewHumanAddedEvent(
			ctx,
			&agg.Aggregate,
			"zitadel-admin",
			"ZITADEL",
			"Admin",
			"",
			"ZITADEL Admin",
			language.English,
			0,
			"admin@zitadel.test",
			false,
		)
		event.AddPasswordData("$plain$x$password", false)
		return event
	}()
}

// machineEvents all events from setup to create the machine user, machinekey can't be tested here, as the public key is not provided and as such the value in the event can't be expected
func machineEvents(ctx context.Context, instanceID, orgID, userID, patID string) []eventstore.Command {
	agg := user.NewAggregate(userID, orgID)
	instanceAgg := instance.NewAggregate(instanceID)
	orgAgg := org.NewAggregate(orgID)
	events := []eventstore.Command{addMachineEvent(ctx, orgID, userID)}
	if patID != "" {
		events = append(events,
			user.NewPersonalAccessTokenAddedEvent(
				ctx,
				&agg.Aggregate,
				patID,
				time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC),
				nil,
			),
		)
	}
	return append(events,
		org.NewMemberAddedEvent(ctx, &orgAgg.Aggregate, userID, domain.RoleOrgOwner),
		instance.NewMemberAddedEvent(ctx, &instanceAgg.Aggregate, userID, domain.RoleIAMOwner),
	)
}

func addMachineEvent(ctx context.Context, orgID, userID string) *user.MachineAddedEvent {
	agg := user.NewAggregate(userID, orgID)
	return user.NewMachineAddedEvent(ctx,
		&agg.Aggregate,
		"zitadel-admin-machine",
		"ZITADEL-machine",
		"Admin",
		false,
		domain.OIDCTokenTypeBearer,
	)
}

func testSetup(ctx context.Context, c *Commands, validations []preparation.Validation) error {
	//nolint:staticcheck
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validations...)
	if err != nil {
		return err
	}

	_, err = c.eventstore.Push(ctx, cmds...)
	return err
}

func TestCommandSide_setupMinimalInterfaces(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		ctx         context.Context
		instanceAgg *instance.Aggregate
		orgAgg      *org.Aggregate
		owner       string
		ids         ZitadelConfig
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "create, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						projectFilters(),
						[]expect{expectPush(
							projectAddedEvents(context.Background(),
								"INSTANCE",
								"ORG",
								"PROJECT",
								"owner",
								false,
							)...,
						),
						},
					)...,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, projectClientIDs()...),
			},
			args: args{
				ctx:         contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				owner:       "owner",
				ids:         instanceSetupZitadelIDs(),
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			validations := make([]preparation.Validation, 0)
			setupMinimalInterfaces(r, &validations, tt.args.instanceAgg, tt.args.orgAgg, tt.args.owner, tt.args.ids)

			err := testSetup(tt.args.ctx, r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_setupAdmins(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id_generator.Generator
		userPasswordHasher *crypto.Hasher
		roles              []authz.RoleMapping
		keyAlgorithm       crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		instanceAgg *instance.Aggregate
		orgAgg      *org.Aggregate
		machine     *AddMachine
		human       *AddHuman
	}
	type res struct {
		owner      string
		pat        bool
		machineKey bool
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "human, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						humanFilters("ORG"),
						adminMemberFilters("ORG", "USER"),
						[]expect{
							expectPush(
								humanEvents(context.Background(),
									"INSTANCE",
									"ORG",
									"USER",
								)...,
							),
						},
					)...,
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "USER"),
				userPasswordHasher: mockPasswordHasher("x"),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
			},
			args: args{
				ctx:         contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				human:       instanceSetupHumanConfig(),
			},
			res: res{
				owner:      "USER",
				pat:        false,
				machineKey: false,
				err:        nil,
			},
		},
		{
			name: "machine, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						machineFilters("ORG", true),
						adminMemberFilters("ORG", "USER-MACHINE"),
						[]expect{
							expectPush(
								machineEvents(context.Background(),
									"INSTANCE",
									"ORG",
									"USER-MACHINE",
									"PAT",
								)...,
							),
						},
					)...,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "USER-MACHINE", "PAT"),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:         contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				machine:     instanceSetupMachineConfig(),
			},
			res: res{
				owner:      "USER-MACHINE",
				pat:        true,
				machineKey: false,
				err:        nil,
			},
		},
		{
			name: "human and machine, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						machineFilters("ORG", true),
						adminMemberFilters("ORG", "USER-MACHINE"),
						humanFilters("ORG"),
						adminMemberFilters("ORG", "USER"),
						[]expect{
							expectPush(
								slices.Concat(
									machineEvents(context.Background(),
										"INSTANCE",
										"ORG",
										"USER-MACHINE",
										"PAT",
									),
									humanEvents(context.Background(),
										"INSTANCE",
										"ORG",
										"USER",
									),
								)...,
							),
						},
					)...,
				),
				userPasswordHasher: mockPasswordHasher("x"),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "USER-MACHINE", "PAT", "USER"),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:         contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgAgg:      org.NewAggregate("ORG"),
				machine:     instanceSetupMachineConfig(),
				human:       instanceSetupHumanConfig(),
			},
			res: res{
				owner:      "USER",
				pat:        true,
				machineKey: false,
				err:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				zitadelRoles:       tt.fields.roles,
				userPasswordHasher: tt.fields.userPasswordHasher,
				keyAlgorithm:       tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			validations := make([]preparation.Validation, 0)
			owner, pat, mk, err := setupAdmins(r, &validations, tt.args.instanceAgg, tt.args.orgAgg, tt.args.machine, tt.args.human)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			err = testSetup(tt.args.ctx, r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			if tt.res.err == nil {
				assert.Equal(t, owner, tt.res.owner)
				if tt.res.pat {
					assert.NotNil(t, pat)
				}
				if tt.res.machineKey {
					assert.NotNil(t, mk)
				}
			}
		})
	}
}

func TestCommandSide_setupDefaultOrg(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id_generator.Generator
		userPasswordHasher *crypto.Hasher
		roles              []authz.RoleMapping
		keyAlgorithm       crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		instanceAgg *instance.Aggregate
		orgName     string
		machine     *AddMachine
		human       *AddHuman
		ids         ZitadelConfig
	}
	type res struct {
		pat        bool
		machineKey bool
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "human and machine, ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						orgFilters(
							"ORG",
							true,
							true,
						),
						[]expect{
							expectPush(
								slices.Concat(
									orgEvents(context.Background(),
										"INSTANCE",
										"ORG",
										"ZITADEL",
										"PROJECT",
										"DOMAIN",
										false,
										true,
										true,
									),
								)...,
							),
						},
					)...,
				),
				userPasswordHasher: mockPasswordHasher("x"),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, orgIDs()...),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:         contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				orgName:     "ZITADEL",
				machine: &AddMachine{
					Machine: &Machine{
						Username:        "zitadel-admin-machine",
						Name:            "ZITADEL-machine",
						Description:     "Admin",
						AccessTokenType: domain.OIDCTokenTypeBearer,
					},
					Pat: &AddPat{
						ExpirationDate: time.Time{},
						Scopes:         nil,
					},
					/* not predictable with the key value in the events
					MachineKey: &AddMachineKey{
						Type:           domain.AuthNKeyTypeJSON,
						ExpirationDate: time.Time{},
					},
					*/
				},
				human: &AddHuman{
					Username:  "zitadel-admin",
					FirstName: "ZITADEL",
					LastName:  "Admin",
					Email: Email{
						Address:  domain.EmailAddress("admin@zitadel.test"),
						Verified: true,
					},
					PreferredLanguage:      language.English,
					Password:               "password",
					PasswordChangeRequired: false,
				},
				ids: ZitadelConfig{
					instanceID:   "INSTANCE",
					orgID:        "ORG",
					projectID:    "PROJECT",
					consoleAppID: "console-id",
					authAppID:    "auth-id",
					mgmtAppID:    "mgmt-id",
					adminAppID:   "admin-id",
				},
			},
			res: res{
				pat:        true,
				machineKey: false,
				err:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				zitadelRoles:       tt.fields.roles,
				userPasswordHasher: tt.fields.userPasswordHasher,
				keyAlgorithm:       tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			validations := make([]preparation.Validation, 0)
			pat, mk, err := setupDefaultOrg(tt.args.ctx, r, &validations, tt.args.instanceAgg, tt.args.orgName, tt.args.machine, tt.args.human, tt.args.ids)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			err = testSetup(context.Background(), r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			if tt.res.err == nil {
				if tt.res.pat {
					assert.NotNil(t, pat)
				}
				if tt.res.machineKey {
					assert.NotNil(t, mk)
				}
			}
		})
	}
}

func TestCommandSide_setupInstanceElements(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx         context.Context
		instanceAgg *instance.Aggregate
		setup       *InstanceSetup
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						setupInstanceElementsFilters("INSTANCE"),
						[]expect{
							expectPush(
								setupInstanceElementsEvents(context.Background(),
									"INSTANCE",
									"ZITADEL",
									language.English,
								)...,
							),
						},
					)...,
				),
			},
			args: args{
				ctx:         contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				instanceAgg: instance.NewAggregate("INSTANCE"),
				setup:       setupInstanceElementsConfig(),
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			validations := setupInstanceElements(tt.args.instanceAgg, tt.args.setup)

			err := testSetup(context.Background(), r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_setUpInstance(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id_generator.Generator
		userPasswordHasher *crypto.Hasher
		roles              []authz.RoleMapping
		keyAlgorithm       crypto.EncryptionAlgorithm
		generateDomain     func(string, string) (string, error)
	}
	type args struct {
		ctx   context.Context
		setup *InstanceSetup
	}
	type res struct {
		pat        bool
		machineKey bool
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					slices.Concat(
						setupInstanceFilters("INSTANCE", "ORG", "PROJECT", "console-id", "DOMAIN"),
						[]expect{
							expectPush(
								setupInstanceEvents(context.Background(),
									"INSTANCE",
									"ORG",
									"PROJECT",
									"console-id",
									"ZITADEL",
									"ZITADEL",
									language.English,
									"DOMAIN",
									false,
								)...,
							),
						},
					)...,
				),
				userPasswordHasher: mockPasswordHasher("x"),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, orgIDs()...),
				roles: []authz.RoleMapping{
					{Role: domain.RoleOrgOwner, Permissions: []string{""}},
					{Role: domain.RoleIAMOwner, Permissions: []string{""}},
				},
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				generateDomain: func(string, string) (string, error) {
					return "DOMAIN", nil
				},
			},
			args: args{
				ctx:   contextWithInstanceSetupInfo(context.Background(), "INSTANCE", "PROJECT", "console-id", "DOMAIN"),
				setup: setupInstanceConfig(),
			},
			res: res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				zitadelRoles:       tt.fields.roles,
				userPasswordHasher: tt.fields.userPasswordHasher,
				keyAlgorithm:       tt.fields.keyAlgorithm,
				GenerateDomain:     tt.fields.generateDomain,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			validations, pat, mk, err := setUpInstance(tt.args.ctx, r, tt.args.setup)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			err = testSetup(tt.args.ctx, r, validations)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			if tt.res.err == nil {
				if tt.res.pat {
					assert.NotNil(t, pat)
				}
				if tt.res.machineKey {
					assert.NotNil(t, mk)
				}
			}
		})
	}
}

func TestCommandSide_UpdateInstance(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		name string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty name, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "instance not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "instance removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusher(
							instance.NewInstanceRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "instance change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewInstanceAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
					),
					expectPush(
						instance.NewInstanceChangedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"INSTANCE_CHANGED",
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.UpdateInstance(tt.args.ctx, tt.args.name)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveInstance(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "instance not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "instance removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusher(
							instance.NewInstanceRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "instance remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewInstanceAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"instance.domain",
								true,
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"custom.domain",
								false,
							),
						),
					),
					expectPush(
						instance.NewInstanceRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"INSTANCE",
							[]string{
								"instance.domain",
								"custom.domain",
							},
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveInstance(tt.args.ctx, tt.args.instanceID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
