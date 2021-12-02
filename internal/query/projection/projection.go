package projection

import (
	"context"
	"database/sql"
	"time"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
)

const (
	CurrentSeqTable   = "projections.current_sequences"
	locksTable        = "projections.locks"
	failedEventsTable = "projections.failed_events"
)

func Start(ctx context.Context, sqlClient *sql.DB, es *eventstore.Eventstore, config Config, defaults systemdefaults.SystemDefaults) error {
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
		LockTable:         locksTable,
		FailedEventsTable: failedEventsTable,
		MaxFailureCount:   config.MaxFailureCount,
		BulkLimit:         config.BulkLimit,
	}

	NewOrgProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["orgs"]))
	NewActionProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["actions"]))
	NewFlowProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["flows"]))
	NewProjectProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["projects"]))
	NewPasswordComplexityProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["password_complexities"]))
	NewPasswordAgeProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["password_age_policy"]))
	NewLockoutPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["lockout_policy"]))
	NewPrivacyPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["privacy_policy"]))
	NewOrgIAMPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_iam_policy"]))
	NewLabelPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["label_policy"]))
	NewProjectGrantProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_grants"]))
	NewProjectRoleProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["project_roles"]))
	// owner.NewOrgOwnerProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_owners"]))
	NewOrgDomainProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["org_domains"]))
	NewLoginPolicyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_policies"]))
	NewIDPProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idps"]))
	NewAppProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["apps"]))
	NewIDPUserLinkProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_user_links"]))
	NewIDPLoginPolicyLinkProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["idp_login_policy_links"]))
	NewMailTemplateProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["mail_templates"]))
	NewMessageTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["message_texts"]))
	NewCustomTextProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["custom_texts"]))
	NewFeatureProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["features"]))
	NewUserProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["users"]))
	NewLoginNameProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["login_names"]))
	_, err := NewKeyProjection(ctx, applyCustomConfig(projectionConfig, config.Customizations["keys"]), defaults.KeyConfig)

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
