package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type Features struct {
	AggregateID              string
	ChangeDate               time.Time
	Sequence                 uint64
	IsDefault                bool
	TierName                 string
	TierDescription          string
	State                    domain.FeaturesState
	StateDescription         string
	AuditLogRetention        time.Duration
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
	LoginPolicyPasswordReset bool
	PasswordComplexityPolicy bool
	LabelPolicyPrivateLabel  bool
	LabelPolicyWatermark     bool
	CustomDomain             bool
	PrivacyPolicy            bool
	MetadataUser             bool
	CustomTextMessage        bool
	CustomTextLogin          bool
	LockoutPolicy            bool
	ActionsAllowed           domain.ActionsAllowed
	MaxActions               int32
}

var (
	featureTable = table{
		name: projection.FeatureTable,
	}
	FeatureColumnAggregateID = Column{
		name:  projection.FeatureAggregateIDCol,
		table: featureTable,
	}
	FeatureColumnChangeDate = Column{
		name:  projection.FeatureChangeDateCol,
		table: featureTable,
	}
	FeatureColumnSequence = Column{
		name:  projection.FeatureSequenceCol,
		table: featureTable,
	}
	FeatureColumnIsDefault = Column{
		name:  projection.FeatureIsDefaultCol,
		table: featureTable,
	}
	FeatureTierName = Column{
		name:  projection.FeatureTierNameCol,
		table: featureTable,
	}
	FeatureTierDescription = Column{
		name:  projection.FeatureTierDescriptionCol,
		table: featureTable,
	}
	FeatureState = Column{
		name:  projection.FeatureStateCol,
		table: featureTable,
	}
	FeatureStateDescription = Column{
		name:  projection.FeatureStateDescriptionCol,
		table: featureTable,
	}
	FeatureAuditLogRetention = Column{
		name:  projection.FeatureAuditLogRetentionCol,
		table: featureTable,
	}
	FeatureLoginPolicyFactors = Column{
		name:  projection.FeatureLoginPolicyFactorsCol,
		table: featureTable,
	}
	FeatureLoginPolicyIDP = Column{
		name:  projection.FeatureLoginPolicyIDPCol,
		table: featureTable,
	}
	FeatureLoginPolicyPasswordless = Column{
		name:  projection.FeatureLoginPolicyPasswordlessCol,
		table: featureTable,
	}
	FeatureLoginPolicyRegistration = Column{
		name:  projection.FeatureLoginPolicyRegistrationCol,
		table: featureTable,
	}
	FeatureLoginPolicyUsernameLogin = Column{
		name:  projection.FeatureLoginPolicyUsernameLoginCol,
		table: featureTable,
	}
	FeatureLoginPolicyPasswordReset = Column{
		name:  projection.FeatureLoginPolicyPasswordResetCol,
		table: featureTable,
	}
	FeaturePasswordComplexityPolicy = Column{
		name:  projection.FeaturePasswordComplexityPolicyCol,
		table: featureTable,
	}
	FeatureLabelPolicyPrivateLabel = Column{
		name:  projection.FeatureLabelPolicyPrivateLabelCol,
		table: featureTable,
	}
	FeatureLabelPolicyWatermark = Column{
		name:  projection.FeatureLabelPolicyWatermarkCol,
		table: featureTable,
	}
	FeatureCustomDomain = Column{
		name:  projection.FeatureCustomDomainCol,
		table: featureTable,
	}
	FeaturePrivacyPolicy = Column{
		name:  projection.FeaturePrivacyPolicyCol,
		table: featureTable,
	}
	FeatureMetadataUser = Column{
		name:  projection.FeatureMetadataUserCol,
		table: featureTable,
	}
	FeatureCustomTextMessage = Column{
		name:  projection.FeatureCustomTextMessageCol,
		table: featureTable,
	}
	FeatureCustomTextLogin = Column{
		name:  projection.FeatureCustomTextLoginCol,
		table: featureTable,
	}
	FeatureLockoutPolicy = Column{
		name:  projection.FeatureLockoutPolicyCol,
		table: featureTable,
	}
	FeatureActionsAllowed = Column{
		name:  projection.FeatureActionsAllowedCol,
		table: featureTable,
	}
	FeatureMaxActions = Column{
		name:  projection.FeatureMaxActionsCol,
		table: featureTable,
	}
)

func (q *Queries) FeaturesByOrgID(ctx context.Context, orgID string) (*Features, error) {
	query, scan := prepareFeaturesQuery()
	stmt, args, err := query.Where(
		sq.Or{
			sq.Eq{
				FeatureColumnAggregateID.identifier(): orgID,
			},
			sq.Eq{
				FeatureColumnAggregateID.identifier(): domain.IAMID,
			},
		}).
		OrderBy(FeatureColumnIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-P9gwg", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) DefaultFeatures(ctx context.Context) (*Features, error) {
	query, scan := prepareFeaturesQuery()
	stmt, args, err := query.Where(sq.Eq{
		FeatureColumnAggregateID.identifier(): domain.IAMID,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-1Ndlg", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func prepareFeaturesQuery() (sq.SelectBuilder, func(*sql.Row) (*Features, error)) {
	return sq.Select(
			FeatureColumnAggregateID.identifier(),
			FeatureColumnChangeDate.identifier(),
			FeatureColumnSequence.identifier(),
			FeatureColumnIsDefault.identifier(),
			FeatureTierName.identifier(),
			FeatureTierDescription.identifier(),
			FeatureState.identifier(),
			FeatureStateDescription.identifier(),
			FeatureAuditLogRetention.identifier(),
			FeatureLoginPolicyFactors.identifier(),
			FeatureLoginPolicyIDP.identifier(),
			FeatureLoginPolicyPasswordless.identifier(),
			FeatureLoginPolicyRegistration.identifier(),
			FeatureLoginPolicyUsernameLogin.identifier(),
			FeatureLoginPolicyPasswordReset.identifier(),
			FeaturePasswordComplexityPolicy.identifier(),
			FeatureLabelPolicyPrivateLabel.identifier(),
			FeatureLabelPolicyWatermark.identifier(),
			FeatureCustomDomain.identifier(),
			FeaturePrivacyPolicy.identifier(),
			FeatureMetadataUser.identifier(),
			FeatureCustomTextMessage.identifier(),
			FeatureCustomTextLogin.identifier(),
			FeatureLockoutPolicy.identifier(),
			FeatureActionsAllowed.identifier(),
			FeatureMaxActions.identifier(),
		).From(featureTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Features, error) {
			p := new(Features)
			tierName := sql.NullString{}
			tierDescription := sql.NullString{}
			stateDescription := sql.NullString{}
			err := row.Scan(
				&p.AggregateID,
				&p.ChangeDate,
				&p.Sequence,
				&p.IsDefault,
				&tierName,
				&tierDescription,
				&p.State,
				&stateDescription,
				&p.AuditLogRetention,
				&p.LoginPolicyFactors,
				&p.LoginPolicyIDP,
				&p.LoginPolicyPasswordless,
				&p.LoginPolicyRegistration,
				&p.LoginPolicyUsernameLogin,
				&p.LoginPolicyPasswordReset,
				&p.PasswordComplexityPolicy,
				&p.LabelPolicyPrivateLabel,
				&p.LabelPolicyWatermark,
				&p.CustomDomain,
				&p.PrivacyPolicy,
				&p.MetadataUser,
				&p.CustomTextMessage,
				&p.CustomTextLogin,
				&p.LockoutPolicy,
				&p.ActionsAllowed,
				&p.MaxActions,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-M9fse", "Errors.Features.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-3o9gd", "Errors.Internal")
			}
			p.TierName = tierName.String
			p.TierDescription = tierDescription.String
			p.StateDescription = stateDescription.String
			return p, nil
		}
}

func (f *Features) EnabledFeatureTypes() []string {
	list := make([]string, 0)
	if f.LoginPolicyFactors {
		list = append(list, domain.FeatureLoginPolicyFactors)
	}
	if f.LoginPolicyIDP {
		list = append(list, domain.FeatureLoginPolicyIDP)
	}
	if f.LoginPolicyPasswordless {
		list = append(list, domain.FeatureLoginPolicyPasswordless)
	}
	if f.LoginPolicyRegistration {
		list = append(list, domain.FeatureLoginPolicyRegistration)
	}
	if f.LoginPolicyUsernameLogin {
		list = append(list, domain.FeatureLoginPolicyUsernameLogin)
	}
	if f.LoginPolicyPasswordReset {
		list = append(list, domain.FeatureLoginPolicyPasswordReset)
	}
	if f.PasswordComplexityPolicy {
		list = append(list, domain.FeaturePasswordComplexityPolicy)
	}
	if f.LabelPolicyPrivateLabel {
		list = append(list, domain.FeatureLabelPolicyPrivateLabel)
	}
	if f.LabelPolicyWatermark {
		list = append(list, domain.FeatureLabelPolicyWatermark)
	}
	if f.CustomDomain {
		list = append(list, domain.FeatureCustomDomain)
	}
	if f.PrivacyPolicy {
		list = append(list, domain.FeaturePrivacyPolicy)
	}
	if f.MetadataUser {
		list = append(list, domain.FeatureMetadataUser)
	}
	if f.CustomTextMessage {
		list = append(list, domain.FeatureCustomTextMessage)
	}
	if f.CustomTextLogin {
		list = append(list, domain.FeatureCustomTextLogin)
	}
	if f.LockoutPolicy {
		list = append(list, domain.FeatureLockoutPolicy)
	}
	if f.ActionsAllowed != domain.ActionsNotAllowed {
		list = append(list, domain.FeatureActions)
	}
	return list
}
