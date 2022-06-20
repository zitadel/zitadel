package projection

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
)

const (
	CurrentSeqTable   = "projections.current_sequences"
	LocksTable        = "projections.locks"
	FailedEventsTable = "projections.failed_events"
)

var (
	OrgProjection                       *orgProjection
	ActionProjection                    *actionProjection
	FlowProjection                      *flowProjection
	ProjectProjection                   *projectProjection
	PasswordComplexityProjection        *passwordComplexityProjection
	PasswordAgeProjection               *passwordAgeProjection
	LockoutPolicyProjection             *lockoutPolicyProjection
	PrivacyPolicyProjection             *privacyPolicyProjection
	DomainPolicyProjection              *domainPolicyProjection
	LabelPolicyProjection               *labelPolicyProjection
	ProjectGrantProjection              *projectGrantProjection
	ProjectRoleProjection               *projectRoleProjection
	OrgDomainProjection                 *orgDomainProjection
	LoginPolicyProjection               *loginPolicyProjection
	IDPProjection                       *idpProjection
	AppProjection                       *appProjection
	IDPUserLinkProjection               *idpUserLinkProjection
	IDPLoginPolicyLinkProjection        *idpLoginPolicyLinkProjection
	MailTemplateProjection              *mailTemplateProjection
	MessageTextProjection               *messageTextProjection
	CustomTextProjection                *customTextProjection
	UserProjection                      *userProjection
	LoginNameProjection                 *loginNameProjection
	OrgMemberProjection                 *orgMemberProjection
	InstanceDomainProjection            *instanceDomainProjection
	InstanceMemberProjection            *instanceMemberProjection
	ProjectMemberProjection             *projectMemberProjection
	ProjectGrantMemberProjection        *projectGrantMemberProjection
	AuthNKeyProjection                  *authNKeyProjection
	PersonalAccessTokenProjection       *personalAccessTokenProjection
	UserGrantProjection                 *userGrantProjection
	UserMetadataProjection              *userMetadataProjection
	UserAuthMethodProjection            *userAuthMethodProjection
	InstanceProjection                  *instanceProjection
	SecretGeneratorProjection           *secretGeneratorProjection
	SMTPConfigProjection                *smtpConfigProjection
	SMSConfigProjection                 *smsConfigProjection
	OIDCSettingsProjection              *oidcSettingsProjection
	DebugNotificationProviderProjection *debugNotificationProviderProjection
	KeyProjection                       *keyProjection
)

func Start(ctx context.Context, sqlClient *sql.DB, es *eventstore.Eventstore, config Config, keyEncryptionAlgorithm crypto.EncryptionAlgorithm, certEncryptionAlgorithm crypto.EncryptionAlgorithm) error {
	projectionConfig := crdb.StatementHandlerConfig{
		ProjectionHandlerConfig: handler.ProjectionHandlerConfig{
			HandlerConfig: handler.HandlerConfig{
				Eventstore: es,
			},
			RequeueEvery:     config.RequeueEvery,
			RetryFailedAfter: config.RetryFailedAfter,
		},
		Client:            sqlClient,
		SequenceTable:     CurrentSeqTable,
		LockTable:         LocksTable,
		FailedEventsTable: FailedEventsTable,
		MaxFailureCount:   config.MaxFailureCount,
		BulkLimit:         config.BulkLimit,
	}

	OrgProjection = newOrgProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["orgs"]))
	ActionProjection = newActionProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["actions"]))
	FlowProjection = newFlowProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["flows"]))
	ProjectProjection = newProjectProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["projects"]))
	PasswordComplexityProjection = newPasswordComplexityProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["password_complexities"]))
	PasswordAgeProjection = newPasswordAgeProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["password_age_policy"]))
	LockoutPolicyProjection = newLockoutPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["lockout_policy"]))
	PrivacyPolicyProjection = newPrivacyPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["privacy_policy"]))
	DomainPolicyProjection = newDomainPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_iam_policy"]))
	LabelPolicyProjection = newLabelPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["label_policy"]))
	ProjectGrantProjection = newProjectGrantProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_grants"]))
	ProjectRoleProjection = newProjectRoleProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_roles"]))
	OrgDomainProjection = newOrgDomainProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_domains"]))
	LoginPolicyProjection = newLoginPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_policies"]))
	IDPProjection = newIDPProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idps"]))
	AppProjection = newAppProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["apps"]))
	IDPUserLinkProjection = newIDPUserLinkProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_user_links"]))
	IDPLoginPolicyLinkProjection = newIDPLoginPolicyLinkProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_login_policy_links"]))
	MailTemplateProjection = newMailTemplateProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["mail_templates"]))
	MessageTextProjection = newMessageTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["message_texts"]))
	CustomTextProjection = newCustomTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["custom_texts"]))
	UserProjection = newUserProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["users"]))
	LoginNameProjection = newLoginNameProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_names"]))
	OrgMemberProjection = newOrgMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_members"]))
	InstanceDomainProjection = newInstanceDomainProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["instance_domains"]))
	InstanceMemberProjection = newInstanceMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["iam_members"]))
	ProjectMemberProjection = newProjectMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_members"]))
	ProjectGrantMemberProjection = newProjectGrantMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_grant_members"]))
	AuthNKeyProjection = newAuthNKeyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["authn_keys"]))
	PersonalAccessTokenProjection = newPersonalAccessTokenProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["personal_access_tokens"]))
	UserGrantProjection = newUserGrantProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_grants"]))
	UserMetadataProjection = newUserMetadataProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_metadata"]))
	UserAuthMethodProjection = newUserAuthMethodProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_auth_method"]))
	InstanceProjection = newInstanceProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["instances"]))
	SecretGeneratorProjection = newSecretGeneratorProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["secret_generators"]))
	SMTPConfigProjection = newSMTPConfigProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["smtp_configs"]))
	SMSConfigProjection = newSMSConfigProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["sms_config"]))
	OIDCSettingsProjection = newOIDCSettingsProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["oidc_settings"]))
	DebugNotificationProviderProjection = newDebugNotificationProviderProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["debug_notification_provider"]))
	KeyProjection = newKeyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["keys"]), keyEncryptionAlgorithm, certEncryptionAlgorithm)
	return nil
}

func applyCustomConfig(config crdb.StatementHandlerConfig, customConfig CustomConfig) crdb.StatementHandlerConfig {
	if customConfig.BulkLimit != nil {
		config.BulkLimit = *customConfig.BulkLimit
	}
	if customConfig.MaxFailureCount != nil {
		config.MaxFailureCount = *customConfig.MaxFailureCount
	}
	if customConfig.RequeueEvery != nil {
		config.RequeueEvery = *customConfig.RequeueEvery
	}
	if customConfig.RetryFailedAfter != nil {
		config.RetryFailedAfter = *customConfig.RetryFailedAfter
	}

	return config
}

func iteratorPool(workerCount int) chan func() {
	if workerCount <= 0 {
		return nil
	}

	queue := make(chan func())
	for i := 0; i < workerCount; i++ {
		go func() {
			for iteration := range queue {
				iteration()
				time.Sleep(2 * time.Second)
			}
		}()
	}
	return queue
}
