package projection

import (
	"context"

	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/migration"
)

const (
	CurrentStateTable = "projections.current_states"
	LocksTable        = "projections.locks"
	FailedEventsTable = "projections.failed_events2"
)

var (
	projectionConfig                    handler.Config
	OrgProjection                       *handler.Handler
	OrgMetadataProjection               *handler.Handler
	ActionProjection                    *handler.Handler
	FlowProjection                      *handler.Handler
	ProjectProjection                   *handler.Handler
	PasswordComplexityProjection        *handler.Handler
	PasswordAgeProjection               *handler.Handler
	LockoutPolicyProjection             *handler.Handler
	PrivacyPolicyProjection             *handler.Handler
	DomainPolicyProjection              *handler.Handler
	LabelPolicyProjection               *handler.Handler
	ProjectGrantProjection              *handler.Handler
	ProjectRoleProjection               *handler.Handler
	OrgDomainProjection                 *handler.Handler
	LoginPolicyProjection               *handler.Handler
	IDPProjection                       *handler.Handler
	AppProjection                       *handler.Handler
	IDPUserLinkProjection               *handler.Handler
	IDPLoginPolicyLinkProjection        *handler.Handler
	IDPTemplateProjection               *handler.Handler
	MailTemplateProjection              *handler.Handler
	MessageTextProjection               *handler.Handler
	CustomTextProjection                *handler.Handler
	UserProjection                      *handler.Handler
	LoginNameProjection                 *handler.Handler
	OrgMemberProjection                 *handler.Handler
	InstanceDomainProjection            *handler.Handler
	InstanceTrustedDomainProjection     *handler.Handler
	InstanceMemberProjection            *handler.Handler
	ProjectMemberProjection             *handler.Handler
	ProjectGrantMemberProjection        *handler.Handler
	AuthNKeyProjection                  *handler.Handler
	PersonalAccessTokenProjection       *handler.Handler
	UserGrantProjection                 *handler.Handler
	UserMetadataProjection              *handler.Handler
	UserAuthMethodProjection            *handler.Handler
	InstanceProjection                  *handler.Handler
	SecretGeneratorProjection           *handler.Handler
	SMTPConfigProjection                *handler.Handler
	SMSConfigProjection                 *handler.Handler
	OIDCSettingsProjection              *handler.Handler
	DebugNotificationProviderProjection *handler.Handler
	KeyProjection                       *handler.Handler
	SecurityPolicyProjection            *handler.Handler
	NotificationPolicyProjection        *handler.Handler
	NotificationsProjection             interface{}
	NotificationsQuotaProjection        interface{}
	TelemetryPusherProjection           interface{}
	DeviceAuthProjection                *handler.Handler
	SessionProjection                   *handler.Handler
	AuthRequestProjection               *handler.Handler
	MilestoneProjection                 *handler.Handler
	QuotaProjection                     *quotaProjection
	LimitsProjection                    *handler.Handler
	RestrictionsProjection              *handler.Handler
	SystemFeatureProjection             *handler.Handler
	InstanceFeatureProjection           *handler.Handler
	TargetProjection                    *handler.Handler
	ExecutionProjection                 *handler.Handler
	UserSchemaProjection                *handler.Handler
	WebKeyProjection                    *handler.Handler
	DebugEventsProjection               *handler.Handler

	ProjectGrantFields      *handler.FieldHandler
	OrgDomainVerifiedFields *handler.FieldHandler
)

type projection interface {
	Start(ctx context.Context)
	Init(ctx context.Context) error
	Trigger(ctx context.Context, opts ...handler.TriggerOpt) (_ context.Context, err error)
	migration.Migration
}

var (
	projections []projection
)

func Create(ctx context.Context, sqlClient *database.DB, es handler.EventStore, config Config, keyEncryptionAlgorithm crypto.EncryptionAlgorithm, certEncryptionAlgorithm crypto.EncryptionAlgorithm, systemUsers map[string]*internal_authz.SystemAPIUser) error {
	projectionConfig = handler.Config{
		Client:                sqlClient,
		Eventstore:            es,
		BulkLimit:             uint16(config.BulkLimit),
		RequeueEvery:          config.RequeueEvery,
		HandleActiveInstances: config.HandleActiveInstances,
		MaxFailureCount:       config.MaxFailureCount,
		RetryFailedAfter:      config.RetryFailedAfter,
		TransactionDuration:   config.TransactionDuration,
	}

	OrgProjection = newOrgProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["orgs"]))
	OrgMetadataProjection = newOrgMetadataProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_metadata"]))
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
	IDPTemplateProjection = newIDPTemplateProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_templates"]))
	MailTemplateProjection = newMailTemplateProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["mail_templates"]))
	MessageTextProjection = newMessageTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["message_texts"]))
	CustomTextProjection = newCustomTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["custom_texts"]))
	UserProjection = newUserProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["users"]))
	LoginNameProjection = newLoginNameProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_names"]))
	OrgMemberProjection = newOrgMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_members"]))
	InstanceDomainProjection = newInstanceDomainProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["instance_domains"]))
	InstanceTrustedDomainProjection = newInstanceTrustedDomainProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["instance_trusted_domains"]))
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
	SecurityPolicyProjection = newSecurityPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["security_policies"]))
	NotificationPolicyProjection = newNotificationPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["notification_policies"]))
	DeviceAuthProjection = newDeviceAuthProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["device_auth"]))
	SessionProjection = newSessionProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["sessions"]))
	AuthRequestProjection = newAuthRequestProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["auth_requests"]))
	MilestoneProjection = newMilestoneProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["milestones"]))
	QuotaProjection = newQuotaProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["quotas"]))
	LimitsProjection = newLimitsProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["limits"]))
	RestrictionsProjection = newRestrictionsProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["restrictions"]))
	SystemFeatureProjection = newSystemFeatureProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["system_features"]))
	InstanceFeatureProjection = newInstanceFeatureProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["instance_features"]))
	TargetProjection = newTargetProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["targets"]))
	ExecutionProjection = newExecutionProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["executions"]))
	UserSchemaProjection = newUserSchemaProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_schemas"]))
	WebKeyProjection = newWebKeyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["web_keys"]))
	DebugEventsProjection = newDebugEventsProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["debug_events"]))

	ProjectGrantFields = newFillProjectGrantFields(applyCustomConfig(projectionConfig, config.Customizations[fieldsProjectGrant]))
	OrgDomainVerifiedFields = newFillOrgDomainVerifiedFields(applyCustomConfig(projectionConfig, config.Customizations[fieldsOrgDomainVerified]))

	newProjectionsList()
	return nil
}

func Projections() []projection {
	return projections
}

func Init(ctx context.Context) error {
	for _, p := range projections {
		if err := p.Init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}

func ProjectInstance(ctx context.Context) error {
	for _, projection := range projections {
		_, err := projection.Trigger(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ApplyCustomConfig(customConfig CustomConfig) handler.Config {
	return applyCustomConfig(projectionConfig, customConfig)
}

func applyCustomConfig(config handler.Config, customConfig CustomConfig) handler.Config {
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
	if customConfig.HandleActiveInstances != nil {
		config.HandleActiveInstances = *customConfig.HandleActiveInstances
	}
	if customConfig.TransactionDuration != nil {
		config.TransactionDuration = *customConfig.TransactionDuration
	}

	return config
}

// we know this is ugly, but we need to have a singleton slice of all projections
// and are only able to initialize it after all projections are created
// as setup and start currently create them individually, we make sure we get the right one
// will be refactored when changing to new id based projections
//
// Event handlers NotificationsProjection, NotificationsQuotaProjection and NotificationsProjection are not added here, because they do not reduce to database statements
func newProjectionsList() {
	projections = []projection{
		OrgProjection,
		OrgMetadataProjection,
		ActionProjection,
		FlowProjection,
		ProjectProjection,
		PasswordComplexityProjection,
		PasswordAgeProjection,
		LockoutPolicyProjection,
		PrivacyPolicyProjection,
		DomainPolicyProjection,
		LabelPolicyProjection,
		ProjectGrantProjection,
		ProjectRoleProjection,
		OrgDomainProjection,
		LoginPolicyProjection,
		IDPProjection,
		IDPTemplateProjection,
		AppProjection,
		IDPUserLinkProjection,
		IDPLoginPolicyLinkProjection,
		MailTemplateProjection,
		MessageTextProjection,
		CustomTextProjection,
		UserProjection,
		LoginNameProjection,
		OrgMemberProjection,
		InstanceDomainProjection,
		InstanceTrustedDomainProjection,
		InstanceMemberProjection,
		ProjectMemberProjection,
		ProjectGrantMemberProjection,
		AuthNKeyProjection,
		PersonalAccessTokenProjection,
		UserGrantProjection,
		UserMetadataProjection,
		UserAuthMethodProjection,
		InstanceProjection,
		SecretGeneratorProjection,
		SMTPConfigProjection,
		SMSConfigProjection,
		OIDCSettingsProjection,
		DebugNotificationProviderProjection,
		KeyProjection,
		SecurityPolicyProjection,
		NotificationPolicyProjection,
		DeviceAuthProjection,
		SessionProjection,
		AuthRequestProjection,
		MilestoneProjection,
		QuotaProjection.handler,
		LimitsProjection,
		RestrictionsProjection,
		SystemFeatureProjection,
		InstanceFeatureProjection,
		TargetProjection,
		ExecutionProjection,
		UserSchemaProjection,
		WebKeyProjection,
		DebugEventsProjection,
	}
}
