package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	IamAggregate models.AggregateType = "iam"

	IamSetupStarted  models.EventType = "iam.setup.started"
	IamSetupDone     models.EventType = "iam.setup.done"
	GlobalOrgSet     models.EventType = "iam.global.org.set"
	IamProjectSet    models.EventType = "iam.project.iam.set"
	IamMemberAdded   models.EventType = "iam.member.added"
	IamMemberChanged models.EventType = "iam.member.changed"
	IamMemberRemoved models.EventType = "iam.member.removed"

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
)
