package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	settings "github.com/zitadel/zitadel/internal/repository/organization_settings"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type relationalTablesProjection struct{}

// SkipV3ReducedEvents implements [handler.RelationalProjection]
func (relationalTablesProjection) SkipV3ReducedEvents() {}

func newRelationalTablesProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(relationalTablesProjection))
}

func (*relationalTablesProjection) Name() string {
	return "relational_tables"
}

func (p *relationalTablesProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceAddedEventType,
					Reduce: p.reduceInstanceAdded,
				},
				{
					Event:  instance.InstanceChangedEventType,
					Reduce: p.reduceInstanceChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: p.reduceInstanceDelete,
				},
				{
					Event:  instance.DefaultOrgSetEventType,
					Reduce: p.reduceDefaultOrgSet,
				},
				{
					Event:  instance.ProjectSetEventType,
					Reduce: p.reduceIAMProjectSet,
				},
				{
					Event:  instance.ManagementConsoleSetEventType,
					Reduce: p.reduceManagementConsoleSet,
				},
				{
					Event:  instance.DefaultLanguageSetEventType,
					Reduce: p.reduceDefaultLanguageSet,
				},

				// settings
				// Login
				{
					Event:  instance.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				{
					Event:  instance.LoginPolicyChangedEventType,
					Reduce: p.reduceLoginPolicyChanged,
				},
				{
					Event:  instance.LoginPolicyMultiFactorAddedEventType,
					Reduce: p.reduceLoginPolicyMFAAdded,
				},
				{
					Event:  instance.LoginPolicyMultiFactorRemovedEventType,
					Reduce: p.reduceLoginPolicyMFARemoved,
				},
				{
					Event:  instance.LoginPolicySecondFactorAddedEventType,
					Reduce: p.reduceLoginPolicySecondFactorAdded,
				},
				{
					Event:  instance.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduceLoginPolicySecondFactorRemoved,
				},
				// Label
				{
					Event:  instance.LabelPolicyAddedEventType,
					Reduce: p.reduceLabelPolicyAdded,
				},
				{
					Event:  instance.LabelPolicyChangedEventType,
					Reduce: p.reduceLabelPolicyChanged,
				},
				{
					Event:  instance.LabelPolicyActivatedEventType,
					Reduce: p.reduceLabelPolicyActivated,
				},
				{
					Event:  instance.LabelPolicyLogoAddedEventType,
					Reduce: p.reduceLabelPolicyLogoAdded,
				},
				{
					Event:  instance.LabelPolicyLogoRemovedEventType,
					Reduce: p.reduceLabelPolicyLogoRemoved,
				},
				{
					Event:  instance.LabelPolicyLogoDarkAddedEventType,
					Reduce: p.reduceLabelPolicyLogoDarkAdded,
				},
				{
					Event:  instance.LabelPolicyLogoDarkRemovedEventType,
					Reduce: p.reduceLabelPolicyLogoDarkRemoved,
				},
				{
					Event:  instance.LabelPolicyIconAddedEventType,
					Reduce: p.reduceLabelPolicyIconAdded,
				},
				{
					Event:  instance.LabelPolicyIconRemovedEventType,
					Reduce: p.reduceLabelPolicyIconRemoved,
				},
				{
					Event:  instance.LabelPolicyIconDarkAddedEventType,
					Reduce: p.reduceLabelPolicyIconDarkAdded,
				},
				{
					Event:  instance.LabelPolicyIconDarkRemovedEventType,
					Reduce: p.reduceLabelPolicyIconDarkRemoved,
				},
				{
					Event:  instance.LabelPolicyFontAddedEventType,
					Reduce: p.reduceLabelPolicyFontAdded,
				},
				{
					Event:  instance.LabelPolicyFontRemovedEventType,
					Reduce: p.reduceLabelPolicyFontRemoved,
				},
				// Password Complexity
				{
					Event:  instance.PasswordComplexityPolicyAddedEventType,
					Reduce: p.reducePasswordComplexityPolicyAdded,
				},
				{
					Event:  instance.PasswordComplexityPolicyChangedEventType,
					Reduce: p.reducePasswordComplexityPolicyChanged,
				},
				// Password Policy
				{
					Event:  instance.PasswordAgePolicyAddedEventType,
					Reduce: p.reducePasswordAgePolicyAdded,
				},
				{
					Event:  instance.PasswordAgePolicyChangedEventType,
					Reduce: p.reducePasswordAgePolicyChanged,
				},
				// Lockout Policy
				{
					Event:  instance.LockoutPolicyAddedEventType,
					Reduce: p.reduceLockoutPolicyAdded,
				},
				{
					Event:  instance.LockoutPolicyChangedEventType,
					Reduce: p.reduceLockoutPolicyChanged,
				},
				// Domain Policy
				{
					Event:  instance.DomainPolicyAddedEventType,
					Reduce: p.reduceDomainPolicyAdded,
				},
				{
					Event:  instance.DomainPolicyChangedEventType,
					Reduce: p.reduceDomainPolicyChanged,
				},
				// Security Policy
				{
					Event:  instance.SecurityPolicySetEventType,
					Reduce: p.reduceSecurityPolicySet,
				},
				// 	Notification
				{
					Event:  instance.NotificationPolicyAddedEventType,
					Reduce: p.reduceNotificationPolicyAdded,
				},
				{
					Event:  instance.NotificationPolicyChangedEventType,
					Reduce: p.reduceNotificationPolicyChanged,
				},
				// Privacy policy
				{
					Event:  instance.PrivacyPolicyAddedEventType,
					Reduce: p.reducePrivacyPolicyAdded,
				},
				{
					Event:  instance.PrivacyPolicyChangedEventType,
					Reduce: p.reducePrivacyPolicyChanged,
				},
				// Secret Generator
				{
					Event:  instance.SecretGeneratorAddedEventType,
					Reduce: p.reduceSecretGeneratorAdded,
				},
				{
					Event:  instance.SecretGeneratorChangedEventType,
					Reduce: p.reduceSecretGeneratorChanged,
				},
				{
					Event:  instance.SecretGeneratorRemovedEventType,
					Reduce: p.reduceSecretGeneratorRemoved,
				},

				// domains
				{
					Event:  instance.InstanceDomainAddedEventType,
					Reduce: p.reduceCustomInstanceDomainAdded,
				},
				{
					Event:  instance.InstanceDomainPrimarySetEventType,
					Reduce: p.reduceInstanceDomainPrimarySet,
				},
				{
					Event:  instance.InstanceDomainRemovedEventType,
					Reduce: p.reduceCustomInstanceDomainRemoved,
				},
				{
					Event:  instance.TrustedDomainAddedEventType,
					Reduce: p.reduceTrustedInstanceDomainAdded,
				},
				{
					Event:  instance.TrustedDomainRemovedEventType,
					Reduce: p.reduceTrustedInstanceDomainRemoved,
				},

				// IDP
				{
					Event:  instance.IDPConfigAddedEventType,
					Reduce: p.reduceIDPAdded,
				},
				{
					Event:  instance.IDPConfigChangedEventType,
					Reduce: p.reduceIDPChanged,
				},
				{
					Event:  instance.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDPDeactivated,
				},
				{
					Event:  instance.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPReactivated,
				},
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  instance.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  instance.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCConfigChanged,
				},
				{
					Event:  instance.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTConfigAdded,
				},
				{
					Event:  instance.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTConfigChanged,
				},
				{
					Event:  instance.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPAdded,
				},
				{
					Event:  instance.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPChanged,
				},
				{
					Event:  instance.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPAdded,
				},
				{
					Event:  instance.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPChanged,
				},
				{
					Event:  instance.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPMigratedAzureAD,
				},
				{
					Event:  instance.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPMigratedGoogle,
				},
				{
					Event:  instance.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPAdded,
				},
				{
					Event:  instance.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPChanged,
				},
				{
					Event:  instance.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPAdded,
				},
				{
					Event:  instance.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPChanged,
				},
				{
					Event:  instance.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPAdded,
				},
				{
					Event:  instance.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPChanged,
				},
				{
					Event:  instance.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPAdded,
				},
				{
					Event:  instance.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPChanged,
				},
				{
					Event:  instance.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPAdded,
				},
				{
					Event:  instance.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPChanged,
				},
				{
					Event:  instance.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPAdded,
				},
				{
					Event:  instance.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPChanged,
				},
				{
					Event:  instance.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPAdded,
				},
				{
					Event:  instance.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPChanged,
				},
				{
					Event:  instance.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  instance.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  instance.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  instance.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  instance.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  instance.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  instance.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgAddedEventType,
					Reduce: p.reduceOrgRelationalAdded,
				},
				{
					Event:  org.OrgChangedEventType,
					Reduce: p.reduceOrgRelationalChanged,
				},
				{
					Event:  org.OrgDeactivatedEventType,
					Reduce: p.reduceOrgRelationalDeactivated,
				},
				{
					Event:  org.OrgReactivatedEventType,
					Reduce: p.reduceOrgRelationalReactivated,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRelationalRemoved,
				},

				// metadata
				{
					Event:  org.MetadataSetType,
					Reduce: p.reduceOrganizationMetadataSet,
				},
				{
					Event:  org.MetadataRemovedType,
					Reduce: p.reduceOrganizationMetadataRemoved,
				},
				{
					Event:  org.MetadataRemovedAllType,
					Reduce: p.reduceOrganizationMetadataRemovedAll,
				},

				// domains
				{
					Event:  org.OrgDomainAddedEventType,
					Reduce: p.reduceOrganizationDomainAdded,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: p.reduceOrganizationDomainPrimarySet,
				},
				{
					Event:  org.OrgDomainRemovedEventType,
					Reduce: p.reduceOrganizationDomainRemoved,
				},
				{
					Event:  org.OrgDomainVerificationAddedEventType,
					Reduce: p.reduceOrganizationDomainVerificationAdded,
				},
				{
					Event:  org.OrgDomainVerifiedEventType,
					Reduce: p.reduceOrganizationDomainVerified,
				},

				// settings
				// Login
				{
					Event:  org.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				{
					Event:  org.LoginPolicyChangedEventType,
					Reduce: p.reduceLoginPolicyChanged,
				},
				{
					Event:  org.LoginPolicyRemovedEventType,
					Reduce: p.reduceLoginPolicyRemoved,
				},
				{
					Event:  org.LoginPolicyMultiFactorAddedEventType,
					Reduce: p.reduceLoginPolicyMFAAdded,
				},
				{
					Event:  org.LoginPolicyMultiFactorRemovedEventType,
					Reduce: p.reduceLoginPolicyMFARemoved,
				},
				{
					Event:  org.LoginPolicySecondFactorAddedEventType,
					Reduce: p.reduceLoginPolicySecondFactorAdded,
				},
				{
					Event:  org.LoginPolicySecondFactorRemovedEventType,
					Reduce: p.reduceLoginPolicySecondFactorRemoved,
				},
				// label
				{
					Event:  org.LabelPolicyAddedEventType,
					Reduce: p.reduceLabelPolicyAdded,
				},
				{
					Event:  org.LabelPolicyChangedEventType,
					Reduce: p.reduceLabelPolicyChanged,
				},
				{
					Event:  org.LabelPolicyRemovedEventType,
					Reduce: p.reduceLabelPolicyRemoved,
				},
				{
					Event:  org.LabelPolicyActivatedEventType,
					Reduce: p.reduceLabelPolicyActivated,
				},
				{
					Event:  org.LabelPolicyLogoAddedEventType,
					Reduce: p.reduceLabelPolicyLogoAdded,
				},
				{
					Event:  org.LabelPolicyLogoRemovedEventType,
					Reduce: p.reduceLabelPolicyLogoRemoved,
				},
				{
					Event:  org.LabelPolicyLogoDarkAddedEventType,
					Reduce: p.reduceLabelPolicyLogoDarkAdded,
				},
				{
					Event:  org.LabelPolicyLogoDarkRemovedEventType,
					Reduce: p.reduceLabelPolicyLogoDarkRemoved,
				},
				{
					Event:  org.LabelPolicyIconAddedEventType,
					Reduce: p.reduceLabelPolicyIconAdded,
				},
				{
					Event:  org.LabelPolicyIconRemovedEventType,
					Reduce: p.reduceLabelPolicyIconRemoved,
				},
				{
					Event:  org.LabelPolicyIconDarkAddedEventType,
					Reduce: p.reduceLabelPolicyIconDarkAdded,
				},
				{
					Event:  org.LabelPolicyIconDarkRemovedEventType,
					Reduce: p.reduceLabelPolicyIconDarkRemoved,
				},
				{
					Event:  org.LabelPolicyFontAddedEventType,
					Reduce: p.reduceLabelPolicyFontAdded,
				},
				{
					Event:  org.LabelPolicyFontRemovedEventType,
					Reduce: p.reduceLabelPolicyFontRemoved,
				},
				// Password Complexity
				{
					Event:  org.PasswordComplexityPolicyAddedEventType,
					Reduce: p.reducePasswordComplexityPolicyAdded,
				},
				{
					Event:  org.PasswordComplexityPolicyChangedEventType,
					Reduce: p.reducePasswordComplexityPolicyChanged,
				},
				{
					Event:  org.PasswordComplexityPolicyRemovedEventType,
					Reduce: p.reducePasswordComplexityPolicyRemoved,
				},
				// Password Policy
				{
					Event:  org.PasswordAgePolicyAddedEventType,
					Reduce: p.reducePasswordAgePolicyAdded,
				},
				{
					Event:  org.PasswordAgePolicyChangedEventType,
					Reduce: p.reducePasswordAgePolicyChanged,
				},
				{
					Event:  org.PasswordAgePolicyRemovedEventType,
					Reduce: p.reducePasswordAgePolicyRemoved,
				},
				// Lockout Policy
				{
					Event:  org.LockoutPolicyAddedEventType,
					Reduce: p.reduceLockoutPolicyAdded,
				},
				{
					Event:  org.LockoutPolicyChangedEventType,
					Reduce: p.reduceLockoutPolicyChanged,
				},
				{
					Event:  org.LockoutPolicyRemovedEventType,
					Reduce: p.reduceOrgLockoutPolicyRemoved,
				},
				// Domain Policy
				{
					Event:  org.DomainPolicyAddedEventType,
					Reduce: p.reduceDomainPolicyAdded,
				},
				{
					Event:  org.DomainPolicyChangedEventType,
					Reduce: p.reduceDomainPolicyChanged,
				},
				{
					Event:  org.DomainPolicyRemovedEventType,
					Reduce: p.reduceOrgDomainPolicyRemoved,
				},
				// Notification
				{
					Event:  org.NotificationPolicyAddedEventType,
					Reduce: p.reduceNotificationPolicyAdded,
				},
				{
					Event:  org.NotificationPolicyChangedEventType,
					Reduce: p.reduceNotificationPolicyChanged,
				},
				{
					Event:  org.NotificationPolicyRemovedEventType,
					Reduce: p.reduceOrgNotificationPolicyRemoved,
				},
				// Privacy
				{
					Event:  org.PrivacyPolicyAddedEventType,
					Reduce: p.reducePrivacyPolicyAdded,
				},
				{
					Event:  org.PrivacyPolicyChangedEventType,
					Reduce: p.reducePrivacyPolicyChanged,
				},
				{
					Event:  org.PrivacyPolicyRemovedEventType,
					Reduce: p.reduceOrgPrivacyPolicyRemoved,
				},

				// idps
				{
					Event:  org.IDPConfigAddedEventType,
					Reduce: p.reduceIDPAdded,
				},
				{
					Event:  org.IDPConfigChangedEventType,
					Reduce: p.reduceIDPChanged,
				},
				{
					Event:  org.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDPDeactivated,
				},
				{
					Event:  org.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPReactivated,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  org.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  org.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCConfigChanged,
				},
				{
					Event:  org.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTConfigAdded,
				},
				{
					Event:  org.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTConfigChanged,
				},
				{
					Event:  org.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPAdded,
				},
				{
					Event:  org.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPChanged,
				},
				{
					Event:  org.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPAdded,
				},
				{
					Event:  org.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPChanged,
				},
				{
					Event:  org.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPMigratedAzureAD,
				},
				{
					Event:  org.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPMigratedGoogle,
				},
				{
					Event:  org.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPAdded,
				},
				{
					Event:  org.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPChanged,
				},
				{
					Event:  org.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPAdded,
				},
				{
					Event:  org.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPChanged,
				},
				{
					Event:  org.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPAdded,
				},
				{
					Event:  org.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPChanged,
				},
				{
					Event:  org.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPAdded,
				},
				{
					Event:  org.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPChanged,
				},
				{
					Event:  org.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPAdded,
				},
				{
					Event:  org.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPChanged,
				},
				{
					Event:  org.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPAdded,
				},
				{
					Event:  org.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPChanged,
				},
				{
					Event:  org.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPAdded,
				},
				{
					Event:  org.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPChanged,
				},
				{
					Event:  org.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  org.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  org.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  org.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  org.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  org.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  org.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
			},
		},
		{
			Aggregate: settings.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  settings.OrganizationSettingsSetEventType,
					Reduce: p.reduceOrganizationSettingsSet,
				},
				{
					Event:  settings.OrganizationSettingsRemovedEventType,
					Reduce: p.reduceOrganizationSettingsRemoved,
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.reduceProjectAdded,
				},
				{
					Event:  project.ProjectChangedType,
					Reduce: p.reduceProjectChanged,
				},
				{
					Event:  project.ProjectDeactivatedType,
					Reduce: p.reduceProjectDeactivated,
				},
				{
					Event:  project.ProjectReactivatedType,
					Reduce: p.reduceProjectReactivated,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},

				// project grants
				{
					Event:  project.GrantAddedType,
					Reduce: p.reduceProjectGrantAdded,
				},
				{
					Event:  project.GrantChangedType,
					Reduce: p.reduceProjectGrantChanged,
				},
				{
					Event:  project.GrantCascadeChangedType,
					Reduce: p.reduceProjectGrantCascadeChanged,
				},
				{
					Event:  project.GrantDeactivatedType,
					Reduce: p.reduceProjectGrantDeactivated,
				},
				{
					Event:  project.GrantReactivatedType,
					Reduce: p.reduceProjectGrantReactivated,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: p.reduceProjectGrantRemoved,
				},

				// project roles
				{
					Event:  project.RoleAddedType,
					Reduce: p.reduceProjectRoleAdded,
				},
				{
					Event:  project.RoleChangedType,
					Reduce: p.reduceProjectRoleChanged,
				},
				{
					Event:  project.RoleRemovedType,
					Reduce: p.reduceProjectRoleRemoved,
				},
			},
		},
		{
			Aggregate: usergrant.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  usergrant.UserGrantAddedType,
					Reduce: p.reduceAuthorizationAdded,
				},
				{
					Event:  usergrant.UserGrantChangedType,
					Reduce: p.reduceAuthorizationChanged,
				},
				{
					Event:  usergrant.UserGrantCascadeChangedType,
					Reduce: p.reduceAuthorizationChanged,
				},
				{
					Event:  usergrant.UserGrantCascadeRemovedType,
					Reduce: p.reduceAuthorizationRemoved,
				},
				{
					Event:  usergrant.UserGrantRemovedType,
					Reduce: p.reduceAuthorizationRemoved,
				},
				{
					Event:  usergrant.UserGrantDeactivatedType,
					Reduce: p.reduceAuthorizationDeactivated,
				},
				{
					Event:  usergrant.UserGrantReactivatedType,
					Reduce: p.reduceAuthorizationReactivated,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserV1AddedType,
					Reduce: p.reduceHumanAdded,
				},
				{
					Event:  user.HumanAddedType,
					Reduce: p.reduceHumanAdded,
				},
				{
					Event:  user.UserV1RegisteredType,
					Reduce: p.reduceHumanRegistered,
				},
				{
					Event:  user.HumanRegisteredType,
					Reduce: p.reduceHumanRegistered,
				},
				{
					Event:  user.UserLockedType,
					Reduce: p.reduceUserLocked,
				},
				{
					Event:  user.UserUnlockedType,
					Reduce: p.reduceUserUnlocked,
				},
				{
					Event:  user.UserDeactivatedType,
					Reduce: p.reduceUserDeactivated,
				},
				{
					Event:  user.UserReactivatedType,
					Reduce: p.reduceUserReactivated,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},

				{
					Event:  user.UserUserNameChangedType,
					Reduce: p.reduceUsernameChanged,
				},
				{
					Event:  user.UserDomainClaimedType,
					Reduce: p.reduceUserDomainClaimed,
				},

				{
					Event:  user.HumanProfileChangedType,
					Reduce: p.reduceHumanProfileChanged,
				},
				{
					Event:  user.UserV1ProfileChangedType,
					Reduce: p.reduceHumanProfileChanged,
				},

				{
					Event:  user.HumanEmailChangedType,
					Reduce: p.reduceHumanEmailChanged,
				},
				{
					Event:  user.UserV1EmailChangedType,
					Reduce: p.reduceHumanEmailChanged,
				},
				{
					Event:  user.HumanEmailVerifiedType,
					Reduce: p.reduceHumanEmailVerified,
				},
				{
					Event:  user.UserV1EmailVerifiedType,
					Reduce: p.reduceHumanEmailVerified,
				},
				{
					Event:  user.HumanEmailCodeAddedType,
					Reduce: p.reduceHumanEmailCodeAdded,
				},
				{
					Event:  user.UserV1EmailCodeAddedType,
					Reduce: p.reduceHumanEmailCodeAdded,
				},
				{
					Event:  user.HumanEmailVerificationFailedType,
					Reduce: p.reduceHumanEmailVerificationFailed,
				},
				{
					Event:  user.UserV1EmailVerificationFailedType,
					Reduce: p.reduceHumanEmailVerificationFailed,
				},
				{
					Event:  user.HumanPhoneChangedType,
					Reduce: p.reduceHumanPhoneChanged,
				},
				{
					Event:  user.UserV1PhoneChangedType,
					Reduce: p.reduceHumanPhoneChanged,
				},
				{
					Event:  user.HumanPhoneRemovedType,
					Reduce: p.reduceHumanPhoneRemoved,
				},
				{
					Event:  user.UserV1PhoneRemovedType,
					Reduce: p.reduceHumanPhoneRemoved,
				},
				{
					Event:  user.HumanPhoneVerifiedType,
					Reduce: p.reduceHumanPhoneVerified,
				},
				{
					Event:  user.UserV1PhoneVerifiedType,
					Reduce: p.reduceHumanPhoneVerified,
				},
				{
					Event:  user.HumanPhoneCodeAddedType,
					Reduce: p.reduceHumanPhoneCodeAdded,
				},
				{
					Event:  user.UserV1PhoneCodeAddedType,
					Reduce: p.reduceHumanPhoneCodeAdded,
				},
				{
					Event:  user.HumanPhoneVerificationFailedType,
					Reduce: p.reduceHumanPhoneVerificationFailed,
				},
				{
					Event:  user.UserV1PhoneVerificationFailedType,
					Reduce: p.reduceHumanPhoneVerificationFailed,
				},

				{
					Event:  user.HumanAvatarAddedType,
					Reduce: p.reduceHumanAvatarAdded,
				},
				{
					Event:  user.HumanAvatarRemovedType,
					Reduce: p.reduceHumanAvatarRemoved,
				},

				{
					Event:  user.HumanPasswordChangedType,
					Reduce: p.reduceHumanPasswordChanged,
				},
				{
					Event:  user.UserV1PasswordChangedType,
					Reduce: p.reduceHumanPasswordChanged,
				},
				{
					Event:  user.HumanPasswordCodeAddedType,
					Reduce: p.reduceHumanPasswordCodeAdded,
				},
				{
					Event:  user.UserV1PasswordCodeAddedType,
					Reduce: p.reduceHumanPasswordCodeAdded,
				},
				{
					Event:  user.HumanPasswordCheckSucceededType,
					Reduce: p.reduceHumanPasswordCheckSucceeded,
				},
				{
					Event:  user.UserV1PasswordCheckSucceededType,
					Reduce: p.reduceHumanPasswordCheckSucceeded,
				},
				{
					Event:  user.HumanPasswordCheckFailedType,
					Reduce: p.reduceHumanPasswordCheckFailed,
				},
				{
					Event:  user.UserV1PasswordCheckFailedType,
					Reduce: p.reduceHumanPasswordCheckFailed,
				},
				{
					Event:  user.HumanPasswordHashUpdatedType,
					Reduce: p.reduceHumanPasswordHashUpdated,
				},

				{
					Event:  user.MachineAddedEventType,
					Reduce: p.reduceMachineAdded,
				},
				{
					Event:  user.MachineChangedEventType,
					Reduce: p.reduceMachineChanged,
				},

				{
					Event:  user.MachineSecretSetType,
					Reduce: p.reduceMachineSecretSet,
				},
				{
					Event:  user.MachineSecretHashUpdatedType,
					Reduce: p.reduceMachineSecretHashUpdated,
				},
				{
					Event:  user.MachineSecretRemovedType,
					Reduce: p.reduceMachineSecretRemoved,
				},

				{
					Event:  user.MachineKeyAddedEventType,
					Reduce: p.reduceMachineKeyAdded,
				},
				{
					Event:  user.MachineKeyRemovedEventType,
					Reduce: p.reduceMachineKeyRemoved,
				},
				{
					Event:  user.UserV1MFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},
				{
					Event:  user.HumanMFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},

				{
					Event:  user.PersonalAccessTokenAddedType,
					Reduce: p.reducePersonalAccessTokenAdded,
				},
				{
					Event:  user.PersonalAccessTokenRemovedType,
					Reduce: p.reducePersonalAccessTokenRemoved,
				},

				{
					Event:  user.MetadataSetType,
					Reduce: p.reduceUserMetadataSet,
				},
				{
					Event:  user.MetadataRemovedType,
					Reduce: p.reduceUserMetadataRemoved,
				},
				{
					Event:  user.MetadataRemovedAllType,
					Reduce: p.reduceUserMetadataRemovedAll,
				},
				{
					Event:  user.HumanPasswordlessTokenAddedType,
					Reduce: p.reducePasskeyAdded,
				},
				{
					Event:  user.HumanPasswordlessTokenVerifiedType,
					Reduce: p.reducePasskeyVerified,
				},
				{
					Event:  user.HumanPasswordlessTokenSignCountChangedType,
					Reduce: p.reducePasskeySignCountSet,
				},
				{
					Event:  user.HumanPasswordlessTokenRemovedType,
					Reduce: p.reducePasskeyRemoved,
				},

				{
					Event:  user.HumanU2FTokenAddedType,
					Reduce: p.reducePasskeyAdded,
				},
				{
					Event:  user.HumanU2FTokenVerifiedType,
					Reduce: p.reducePasskeyVerified,
				},
				{
					Event:  user.HumanU2FTokenSignCountChangedType,
					Reduce: p.reducePasskeySignCountSet,
				},
				{
					Event:  user.HumanU2FTokenRemovedType,
					Reduce: p.reducePasskeyRemoved,
				},

				{
					Event:  user.HumanPasswordlessInitCodeAddedType,
					Reduce: p.reducePasskeyInitCodeAdded,
				},
				{
					Event:  user.HumanPasswordlessInitCodeCheckFailedType,
					Reduce: p.reducePasskeyInitCodeCheckFailed,
				},
				{
					Event:  user.HumanPasswordlessInitCodeCheckSucceededType,
					Reduce: p.reducePasskeyInitCodeCheckSucceeded,
				},
				{
					Event:  user.HumanPasswordlessInitCodeRequestedType,
					Reduce: p.reducePasskeyInitCodeRequested,
				},
				{
					Event:  user.UserIDPLinkAddedType,
					Reduce: p.reduceIDPLinkAdded,
				},
				{
					Event:  user.UserIDPLinkCascadeRemovedType,
					Reduce: p.reduceIDPLinkCascadeRemoved,
				},
				{
					Event:  user.UserIDPLinkRemovedType,
					Reduce: p.reduceIDPLinkRemoved,
				},
				{
					Event:  user.UserIDPExternalIDMigratedType,
					Reduce: p.reduceIDPLinkUserIDMigrated,
				},
				{
					Event:  user.UserIDPExternalUsernameChangedType,
					Reduce: p.reduceIDPLinkUsernameChanged,
				},
				{
					Event:  user.HumanMFAOTPAddedType,
					Reduce: p.reduceTOTPAdded,
				},
				{
					Event:  user.UserV1MFAOTPAddedType,
					Reduce: p.reduceTOTPAdded,
				},
				{
					Event:  user.HumanMFAOTPVerifiedType,
					Reduce: p.reduceTOTPVerified,
				},
				{
					Event:  user.UserV1MFAOTPVerifiedType,
					Reduce: p.reduceTOTPVerified,
				},
				{
					Event:  user.HumanMFAOTPRemovedType,
					Reduce: p.reduceTOTPRemoved,
				},
				{
					Event:  user.UserV1MFAOTPRemovedType,
					Reduce: p.reduceTOTPRemoved,
				},
				{
					Event:  user.HumanMFAOTPCheckSucceededType,
					Reduce: p.reduceTOTPCheckSucceeded,
				},
				{
					Event:  user.UserV1MFAOTPCheckSucceededType,
					Reduce: p.reduceTOTPCheckSucceeded,
				},
				{
					Event:  user.HumanMFAOTPCheckFailedType,
					Reduce: p.reduceTOTPCheckFailed,
				},
				{
					Event:  user.UserV1MFAOTPCheckFailedType,
					Reduce: p.reduceTOTPCheckFailed,
				},
				{
					Event:  user.HumanOTPSMSAddedType,
					Reduce: p.reduceOTPSMSEnabled,
				},
				{
					Event:  user.HumanOTPSMSRemovedType,
					Reduce: p.reduceOTPSMSDisabled,
				},
				{
					Event:  user.HumanOTPSMSCheckSucceededType,
					Reduce: p.reduceOTPSMSCheckSucceeded,
				},
				{
					Event:  user.HumanOTPSMSCheckFailedType,
					Reduce: p.reduceOTPSMSCheckFailed,
				},
				{
					Event:  user.HumanOTPEmailAddedType,
					Reduce: p.reduceOTPEmailEnabled,
				},
				{
					Event:  user.HumanOTPEmailRemovedType,
					Reduce: p.reduceOTPEmailDisabled,
				},
				{
					Event:  user.HumanOTPEmailCheckSucceededType,
					Reduce: p.reduceOTPEmailCheckSucceeded,
				},
				{
					Event:  user.HumanOTPEmailCheckFailedType,
					Reduce: p.reduceOTPEmailCheckFailed,
				},
				{
					Event:  user.HumanInviteCodeAddedType,
					Reduce: p.reduceInviteCodeAdded,
				},
				{
					Event:  user.HumanInviteCheckSucceededType,
					Reduce: p.reduceInviteCheckSucceeded,
				},
				{
					Event:  user.HumanInviteCheckFailedType,
					Reduce: p.reduceInviteCheckFailed,
				},
				{
					Event:  user.HumanRecoveryCodesAddedType,
					Reduce: p.reduceRecoveryCodesAdded,
				},
				{
					Event:  user.HumanRecoveryCodesRemovedType,
					Reduce: p.reduceRecoveryCodesRemoved,
				},
				{
					Event:  user.HumanRecoveryCodeCheckSucceededType,
					Reduce: p.reduceRecoveryCodeCheckSucceeded,
				},
				{
					Event:  user.HumanRecoveryCodeCheckFailedType,
					Reduce: p.reduceRecoveryCodeCheckFailed,
				},
			},
		},
		{
			Aggregate: idpintent.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  idpintent.StartedEventType,
					Reduce: p.reduceIDPIntentStartedEvent,
				},
				{
					Event:  idpintent.SucceededEventType,
					Reduce: p.reduceIDPIntentSucceededEvent,
				},
				{
					Event:  idpintent.SAMLSucceededEventType,
					Reduce: p.reduceIDPIntentSAMLSucceededEvent,
				},
				{
					Event:  idpintent.SAMLRequestEventType,
					Reduce: p.reduceIDPIntentSAMLRequestEvent,
				},
				{
					Event:  idpintent.LDAPSucceededEventType,
					Reduce: p.reduceIDPIntentLDAPSucceededEvent,
				},
				{
					Event:  idpintent.FailedEventType,
					Reduce: p.reduceIDPIntentFailedEvent,
				},
				{
					Event:  idpintent.ConsumedEventType,
					Reduce: p.reduceIDPIntentConsumedEvent,
				},
			},
		},
		{
			Aggregate: session.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  session.AddedType,
					Reduce: p.reduceSessionAdded,
				},
				{
					Event:  session.UserCheckedType,
					Reduce: p.reduceSessionUserChecked,
				},
				{
					Event:  session.PasswordCheckedType,
					Reduce: p.reduceSessionPasswordChecked,
				},
				{
					Event:  session.IntentCheckedType,
					Reduce: p.reduceSessionIntentChecked,
				},
				{
					Event:  session.WebAuthNChallengedType,
					Reduce: p.reduceSessionWebAuthNChallenged,
				},
				{
					Event:  session.WebAuthNCheckedType,
					Reduce: p.reduceSessionWebAuthNChecked,
				},
				{
					Event:  session.TOTPCheckedType,
					Reduce: p.reduceSessionTOTPChecked,
				},
				{
					Event:  session.OTPSMSChallengedType,
					Reduce: p.reduceSessionOTPSMSChallenged,
				},
				{
					Event:  session.OTPSMSCheckedType,
					Reduce: p.reduceSessionOTPSMSChecked,
				},
				{
					Event:  session.OTPEmailChallengedType,
					Reduce: p.reduceSessionOTPEmailChallenged,
				},
				{
					Event:  session.OTPEmailCheckedType,
					Reduce: p.reduceSessionOTPEmailChecked,
				},
				{
					Event:  session.RecoveryCodeCheckedType,
					Reduce: p.reduceSessionRecoveryCodeChecked,
				},
				{
					Event:  session.TokenSetType,
					Reduce: p.reduceSessionTokenSet,
				},
				{
					Event:  session.MetadataSetType,
					Reduce: p.reduceSessionMetadataSet,
				},
				{
					Event:  session.LifetimeSetType,
					Reduce: p.reduceSessionLifetimeSet,
				},
				{
					Event:  session.TerminateType,
					Reduce: p.reduceSessionTerminated,
				},
			},
		},
	}
}
