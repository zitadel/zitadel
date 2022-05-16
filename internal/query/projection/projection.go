package projection

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
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
	OrgProjection                 *orgProjection
	LoginNameProjection           *loginNameProjection
	ActionProjection              *actionProjection
	FlowProjection                *flowProjection
	ProjectProjection             *projectProjection
	PasswordComplexityProjection  *passwordComplexityProjection
	PasswordAgeProjection         *passwordAgeProjection
	LockoutPolicyProjection       *lockoutPolicyProjection
	PrivacyPolicyProjection       *privacyPolicyProjection
	OrgIAMPolicyProjection        *orgIAMPolicyProjection
	LabelPolicyProjection         *labelPolicyProjection
	ProjectGrantProjection        *projectGrantProjection
	ProjectRoleProjection         *projectRoleProjection
	OrgDomainProjection           *orgDomainProjection
	LoginPolicyProjection         *loginPolicyProjection
	IDPProjection                 *idpProjection
	AppProjection                 *appProjection
	IDPUserLinkProjection         *idpUserLinkProjection
	IDPLoginPolicyLinkProjection  *idpLoginPolicyLinkProjection
	MailTemplateProjection        *mailTemplateProjection
	MessageTextProjection         *messageTextProjection
	CustomTextProjection          *customTextProjection
	FeatureProjection             *featureProjection
	UserProjection                *userProjection
	OrgMemberProjection           *orgMemberProjection
	IAMMemberProjection           *iamMemberProjection
	ProjectMemberProjection       *projectMemberProjection
	ProjectGrantMemberProjection  *projectGrantMemberProjection
	AuthNKeyProjection            *authNKeyProjection
	PersonalAccessTokenProjection *personalAccessTokenProjection
	UserGrantProjection           *userGrantProjection
	UserMetadataProjection        *userMetadataProjection
	UserAuthMethodProjection      *userAuthMethodProjection
	IAMProjection                 *iamProjection
	KeyProjection                 *keyProjection
)

func Start(ctx context.Context, sqlClient *sql.DB, es *eventstore.Eventstore, config Config, defaults systemdefaults.SystemDefaults, keyChan chan<- interface{}) (err error) {
	projectionConfig := crdb.StatementHandlerConfig{
		ProjectionHandlerConfig: handler.ProjectionHandlerConfig{
			HandlerConfig: handler.HandlerConfig{
				Eventstore: es,
			},
			RequeueEvery:     config.RequeueEvery.Duration,
			RetryFailedAfter: config.RetryFailedAfter.Duration,
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
	OrgIAMPolicyProjection = newOrgIAMPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_iam_policy"]))
	LabelPolicyProjection = newLabelPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["label_policy"]))
	ProjectGrantProjection = newProjectGrantProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_grants"]))
	ProjectRoleProjection = newProjectRoleProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_roles"]))
	// owner.NewOrgOwnerProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_owners"]))
	OrgDomainProjection = newOrgDomainProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_domains"]))
	LoginPolicyProjection = newLoginPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_policies"]))
	IDPProjection = newIDPProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idps"]))
	AppProjection = newAppProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["apps"]))
	IDPUserLinkProjection = newIDPUserLinkProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_user_links"]))
	IDPLoginPolicyLinkProjection = newIDPLoginPolicyLinkProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_login_policy_links"]))
	MailTemplateProjection = newMailTemplateProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["mail_templates"]))
	MessageTextProjection = newMessageTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["message_texts"]))
	CustomTextProjection = newCustomTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["custom_texts"]))
	FeatureProjection = newFeatureProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["features"]))
	UserProjection = newUserProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["users"]))
	LoginNameProjection = newLoginNameProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_names"]))
	OrgMemberProjection = newOrgMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_members"]))
	IAMMemberProjection = newIAMMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["iam_members"]))
	ProjectMemberProjection = newProjectMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_members"]))
	ProjectGrantMemberProjection = newProjectGrantMemberProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_grant_members"]))
	AuthNKeyProjection = newAuthNKeyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["authn_keys"]))
	PersonalAccessTokenProjection = newPersonalAccessTokenProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["personal_access_tokens"]))
	UserGrantProjection = newUserGrantProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_grants"]))
	UserMetadataProjection = newUserMetadataProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_metadata"]))
	UserAuthMethodProjection = newUserAuthMethodProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["user_auth_method"]))
	IAMProjection = newIAMProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["iam"]))
	KeyProjection, err = newKeyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["keys"]), defaults.KeyConfig, keyChan)

	return err
}

func applyCustomConfig(config crdb.StatementHandlerConfig, customConfig CustomConfig) crdb.StatementHandlerConfig {
	if customConfig.BulkLimit != nil {
		config.BulkLimit = *customConfig.BulkLimit
	}
	if customConfig.MaxFailureCount != nil {
		config.MaxFailureCount = *customConfig.MaxFailureCount
	}
	if customConfig.RequeueEvery != nil {
		config.RequeueEvery = customConfig.RequeueEvery.Duration
	}
	if customConfig.RetryFailedAfter != nil {
		config.RetryFailedAfter = customConfig.RetryFailedAfter.Duration
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
