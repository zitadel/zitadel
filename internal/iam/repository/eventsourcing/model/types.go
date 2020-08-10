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

	IdpConfigAdded       models.EventType = "iam.idp.config.added"
	IdpConfigChanged     models.EventType = "iam.idp.config.changed"
	IdpConfigRemoved     models.EventType = "iam.idp.config.removed"
	IdpConfigDeactivated models.EventType = "iam.idp.config.deactivated"
	IdpConfigReactivated models.EventType = "iam.idp.config.reactivated"

	OidcIdpConfigAdded   models.EventType = "iam.idp.oidc.config.added"
	OidcIdpConfigChanged models.EventType = "iam.idp.oidc.config.changed"

	SamlIdpConfigAdded   models.EventType = "iam.idp.saml.config.added"
	SamlIdpConfigChanged models.EventType = "iam.idp.saml.config.changed"

	LoginPolicyAdded              models.EventType = "iam.policy.login.added"
	LoginPolicyChanged            models.EventType = "iam.policy.login.changed"
	LoginPolicyIdpProviderAdded   models.EventType = "iam.policy.login.idpprovider.added"
	LoginPolicyIdpProviderRemoved models.EventType = "iam.policy.login.idpprovider.removed"
)
