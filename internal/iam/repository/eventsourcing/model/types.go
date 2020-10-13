package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	IAMAggregate models.AggregateType = "iam"

	IAMSetupStarted  models.EventType = "iam.setup.started"
	IAMSetupDone     models.EventType = "iam.setup.done"
	GlobalOrgSet     models.EventType = "iam.global.org.set"
	IAMProjectSet    models.EventType = "iam.project.iam.set"
	IAMMemberAdded   models.EventType = "iam.member.added"
	IAMMemberChanged models.EventType = "iam.member.changed"
	IAMMemberRemoved models.EventType = "iam.member.removed"

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
	LabelPolicyAdded                     models.EventType = "iam.policy.label.added"
	LabelPolicyChanged                   models.EventType = "iam.policy.label.changed"
)
