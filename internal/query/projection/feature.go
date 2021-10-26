package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/repository/features"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

type FeatureProjection struct {
	crdb.StatementHandler
}

const (
	FeatureTable = "zitadel.projections.features"
)

func NewFeatureProjection(ctx context.Context, config crdb.StatementHandlerConfig) *FeatureProjection {
	p := &FeatureProjection{}
	config.ProjectionName = FeatureTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *FeatureProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.FeaturesSetEventType,
					Reduce: p.reduceFeatureSet,
				},
				{
					Event:  org.FeaturesRemovedEventType,
					Reduce: p.reduceFeatureRemoved,
				},
			},
		},
		{
			Aggregate: iam.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.FeaturesSetEventType,
					Reduce: p.reduceFeatureSet,
				},
			},
		},
	}
}

const (
	FeatureAggregateIDCol              = "aggregate_id"
	FeatureCreationDateCol             = "creation_date"
	FeatureChangeDateCol               = "change_date"
	FeatureSequenceCol                 = "sequence"
	FeatureIsDefaultCol                = "is_default"
	FeatureTierNameCol                 = "tier_name"
	FeatureTierDescriptionCol          = "tier_description"
	FeatureStateCol                    = "state"
	FeatureStateDescriptionCol         = "state_description"
	FeatureAuditLogRetentionCol        = "audit_log_retention"
	FeatureLoginPolicyFactorsCol       = "login_policy_factors"
	FeatureLoginPolicyIDPCol           = "login_policy_idp"
	FeatureLoginPolicyPasswordlessCol  = "login_policy_passwordless"
	FeatureLoginPolicyRegistrationCol  = "login_policy_registration"
	FeatureLoginPolicyUsernameLoginCol = "login_policy_username_login"
	FeatureLoginPolicyPasswordResetCol = "login_policy_password_reset"
	FeaturePasswordComplexityPolicyCol = "password_complexity_policy"
	FeatureLabelPolicyPrivateLabelCol  = "label_policy_private_label"
	FeatureLabelPolicyWatermarkCol     = "label_policy_watermark"
	FeatureCustomDomainCol             = "custom_domain"
	FeaturePrivacyPolicyCol            = "privacy_policy"
	FeatureMetadataUserCol             = "meta_data_user"
	FeatureCustomTextMessageCol        = "custom_text_message"
	FeatureCustomTextLoginCol          = "custom_text_login"
	FeatureLockoutPolicyCol            = "lockout_policy"
	FeatureActionsCol                  = "actions"
)

func (p *FeatureProjection) reduceFeatureSet(event eventstore.EventReader) (*handler.Statement, error) {
	var featureEvent features.FeaturesSetEvent
	var isDefault bool
	switch e := event.(type) {
	case *iam.FeaturesSetEvent:
		featureEvent = e.FeaturesSetEvent
		isDefault = true
	case *org.FeaturesSetEvent:
		featureEvent = e.FeaturesSetEvent
		isDefault = false
	default:
		logging.LogWithFields("HANDL-M9ets", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{org.FeaturesSetEventType, iam.FeaturesSetEventType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-K0erf", "reduce.wrong.event.type")
	}

	return crdb.NewCreateStatement(&featureEvent, []handler.Column{
		handler.NewCol(FeatureAggregateIDCol, featureEvent.Aggregate().ID),
		handler.NewCol(FeatureCreationDateCol, featureEvent.CreationDate()),
		handler.NewCol(FeatureChangeDateCol, featureEvent.CreationDate()),
		handler.NewCol(FeatureSequenceCol, featureEvent.Sequence()),
		handler.NewCol(FeatureIsDefaultCol, isDefault),
		handler.NewCol(FeatureTierNameCol, featureEvent.TierName),
		handler.NewCol(FeatureTierDescriptionCol, featureEvent.TierDescription),
		handler.NewCol(FeatureStateCol, featureEvent.State),
		handler.NewCol(FeatureStateDescriptionCol, featureEvent.StateDescription),
		handler.NewCol(FeatureAuditLogRetentionCol, featureEvent.AuditLogRetention),
		handler.NewCol(FeatureLoginPolicyFactorsCol, featureEvent.LoginPolicyFactors),
		handler.NewCol(FeatureLoginPolicyIDPCol, featureEvent.LoginPolicyIDP),
		handler.NewCol(FeatureLoginPolicyPasswordlessCol, featureEvent.LoginPolicyPasswordless),
		handler.NewCol(FeatureLoginPolicyRegistrationCol, featureEvent.LoginPolicyRegistration),
		handler.NewCol(FeatureLoginPolicyUsernameLoginCol, featureEvent.LoginPolicyUsernameLogin),
		handler.NewCol(FeatureLoginPolicyPasswordResetCol, featureEvent.LoginPolicyPasswordReset),
		handler.NewCol(FeaturePasswordComplexityPolicyCol, featureEvent.PasswordComplexityPolicy),
		handler.NewCol(FeatureLabelPolicyPrivateLabelCol, featureEvent.LabelPolicyPrivateLabel),
		handler.NewCol(FeatureLabelPolicyWatermarkCol, featureEvent.LabelPolicyWatermark),
		handler.NewCol(FeatureCustomDomainCol, featureEvent.CustomDomain),
		handler.NewCol(FeaturePrivacyPolicyCol, featureEvent.PrivacyPolicy),
		handler.NewCol(FeatureMetadataUserCol, featureEvent.MetadataUser),
		handler.NewCol(FeatureCustomTextMessageCol, featureEvent.CustomTextMessage),
		handler.NewCol(FeatureCustomTextLoginCol, featureEvent.CustomTextLogin),
		handler.NewCol(FeatureLockoutPolicyCol, featureEvent.LockoutPolicy),
		handler.NewCol(FeatureActionsCol, featureEvent.Actions),
	}), nil
}

func (p *FeatureProjection) reduceFeatureRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*org.FeaturesRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-fN903", "seq", event.Sequence(), "expectedType", org.FeaturesRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-0p4rf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(FeatureAggregateIDCol, e.Aggregate().ID),
		},
	), nil
}
