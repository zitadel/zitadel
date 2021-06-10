package model

import "github.com/caos/zitadel/internal/eventstore/v1/models"

const (
	IAMAggregate models.AggregateType = "iam"

	IAMSetupStarted         models.EventType = "iam.setup.started"
	IAMSetupDone            models.EventType = "iam.setup.done"
	GlobalOrgSet            models.EventType = "iam.global.org.set"
	IAMProjectSet           models.EventType = "iam.project.iam.set"
	IAMMemberAdded          models.EventType = "iam.member.added"
	IAMMemberChanged        models.EventType = "iam.member.changed"
	IAMMemberRemoved        models.EventType = "iam.member.removed"
	IAMMemberCascadeRemoved models.EventType = "iam.member.cascade.removed"

	IDPConfigAdded       models.EventType = "iam.idp.config.added"
	IDPConfigChanged     models.EventType = "iam.idp.config.changed"
	IDPConfigRemoved     models.EventType = "iam.idp.config.removed"
	IDPConfigDeactivated models.EventType = "iam.idp.config.deactivated"
	IDPConfigReactivated models.EventType = "iam.idp.config.reactivated"

	OIDCIDPConfigAdded   models.EventType = "iam.idp.oidc.config.added"
	OIDCIDPConfigChanged models.EventType = "iam.idp.oidc.config.changed"

	SAMLIDPConfigAdded   models.EventType = "iam.idp.saml.config.added"
	SAMLIDPConfigChanged models.EventType = "iam.idp.saml.config.changed"

	LoginPolicyAdded                     models.EventType = "iam.policy.login.added"
	LoginPolicyChanged                   models.EventType = "iam.policy.login.changed"
	LoginPolicyIDPProviderAdded          models.EventType = "iam.policy.login.idpprovider.added"
	LoginPolicyIDPProviderRemoved        models.EventType = "iam.policy.login.idpprovider.removed"
	LoginPolicyIDPProviderCascadeRemoved models.EventType = "iam.policy.login.idpprovider.cascade.removed"
	LoginPolicySecondFactorAdded         models.EventType = "iam.policy.login.secondfactor.added"
	LoginPolicySecondFactorRemoved       models.EventType = "iam.policy.login.secondfactor.removed"
	LoginPolicyMultiFactorAdded          models.EventType = "iam.policy.login.multifactor.added"
	LoginPolicyMultiFactorRemoved        models.EventType = "iam.policy.login.multifactor.removed"

	LabelPolicyAdded     models.EventType = "iam.policy.label.added"
	LabelPolicyChanged   models.EventType = "iam.policy.label.changed"
	LabelPolicyActivated models.EventType = "iam.policy.label.activated"

	LabelPolicyLogoAdded       models.EventType = "iam.policy.label.logo.added"
	LabelPolicyLogoRemoved     models.EventType = "iam.policy.label.logo.removed"
	LabelPolicyIconAdded       models.EventType = "iam.policy.label.icon.added"
	LabelPolicyIconRemoved     models.EventType = "iam.policy.label.icon.removed"
	LabelPolicyLogoDarkAdded   models.EventType = "iam.policy.label.logo.dark.added"
	LabelPolicyLogoDarkRemoved models.EventType = "iam.policy.label.logo.dark.removed"
	LabelPolicyIconDarkAdded   models.EventType = "iam.policy.label.icon.dark.added"
	LabelPolicyIconDarkRemoved models.EventType = "iam.policy.label.icon.dark.removed"
	LabelPolicyFontAdded       models.EventType = "iam.policy.label.font.added"
	LabelPolicyFontRemoved     models.EventType = "iam.policy.label.font.removed"
	LabelPolicyAssetsRemoved   models.EventType = "iam.policy.label.assets.removed"

	MailTemplateAdded   models.EventType = "iam.mail.template.added"
	MailTemplateChanged models.EventType = "iam.mail.template.changed"

	CustomTextSet     models.EventType = "iam.customtext.set"
	CustomTextRemoved models.EventType = "iam.customtext.removed"

	CustomTextSet     models.EventType = "iam.customtext.set"
	CustomTextRemoved models.EventType = "iam.customtext.removed"

	PasswordComplexityPolicyAdded   models.EventType = "iam.policy.password.complexity.added"
	PasswordComplexityPolicyChanged models.EventType = "iam.policy.password.complexity.changed"

	PasswordAgePolicyAdded   models.EventType = "iam.policy.password.age.added"
	PasswordAgePolicyChanged models.EventType = "iam.policy.password.age.changed"

	PasswordLockoutPolicyAdded   models.EventType = "iam.policy.password.lockout.added"
	PasswordLockoutPolicyChanged models.EventType = "iam.policy.password.lockout.changed"

	OrgIAMPolicyAdded   models.EventType = "iam.policy.org.iam.added"
	OrgIAMPolicyChanged models.EventType = "iam.policy.org.iam.changed"
)
