package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type Feature struct {
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
	Actions                  bool
}

var (
	feautureTable = table{
		name: projection.FeatureTable,
	}
	FeatureColumnAggregateID = Column{
		name:  projection.FeatureAggregateIDCol,
		table: feautureTable,
	}
	FeatureColumnChangeDate = Column{
		name:  projection.FeatureChangeDateCol,
		table: feautureTable,
	}
	FeatureColumnSequence = Column{
		name:  projection.FeatureSequenceCol,
		table: feautureTable,
	}
	FeatureColumnIsDefault = Column{
		name:  projection.FeatureIsDefaultCol,
		table: feautureTable,
	}
	FeatureTierName = Column{
		name:  projection.FeatureTierNameCol,
		table: feautureTable,
	}
	FeatureTierDescription = Column{
		name:  projection.FeatureTierDescriptionCol,
		table: feautureTable,
	}
	FeatureState = Column{
		name:  projection.FeatureStateCol,
		table: feautureTable,
	}
	FeatureStateDescription = Column{
		name:  projection.FeatureStateDescriptionCol,
		table: feautureTable,
	}
	FeatureAuditLogRetention = Column{
		name:  projection.FeatureAuditLogRetentionCol,
		table: feautureTable,
	}
	FeatureLoginPolicyFactors = Column{
		name:  projection.FeatureLoginPolicyFactorsCol,
		table: feautureTable,
	}
	FeatureLoginPolicyIDP = Column{
		name:  projection.FeatureLoginPolicyIDPCol,
		table: feautureTable,
	}
	FeatureLoginPolicyPasswordless = Column{
		name:  projection.FeatureLoginPolicyPasswordlessCol,
		table: feautureTable,
	}
	FeatureLoginPolicyRegistration = Column{
		name:  projection.FeatureLoginPolicyRegistrationCol,
		table: feautureTable,
	}
	FeatureLoginPolicyUsernameLogin = Column{
		name:  projection.FeatureLoginPolicyUsernameLoginCol,
		table: feautureTable,
	}
	FeatureLoginPolicyPasswordReset = Column{
		name:  projection.FeatureLoginPolicyPasswordResetCol,
		table: feautureTable,
	}
	FeaturePasswordComplexityPolicy = Column{
		name:  projection.FeaturePasswordComplexityPolicyCol,
		table: feautureTable,
	}
	FeatureLabelPolicyPrivateLabel = Column{
		name:  projection.FeatureLabelPolicyPrivateLabelCol,
		table: feautureTable,
	}
	FeatureLabelPolicyWatermark = Column{
		name:  projection.FeatureLabelPolicyWatermarkCol,
		table: feautureTable,
	}
	FeatureCustomDomain = Column{
		name:  projection.FeatureCustomDomainCol,
		table: feautureTable,
	}
	FeaturePrivacyPolicy = Column{
		name:  projection.FeaturePrivacyPolicyCol,
		table: feautureTable,
	}
	FeatureMetadataUser = Column{
		name:  projection.FeatureMetadataUserCol,
		table: feautureTable,
	}
	FeatureCustomTextMessage = Column{
		name:  projection.FeatureCustomTextMessageCol,
		table: feautureTable,
	}
	FeatureCustomTextLogin = Column{
		name:  projection.FeatureCustomTextLoginCol,
		table: feautureTable,
	}
	FeatureLockoutPolicy = Column{
		name:  projection.FeatureLockoutPolicyCol,
		table: feautureTable,
	}
	FeatureActions = Column{
		name:  projection.FeatureActionsCol,
		table: feautureTable,
	}
)

func (q *Queries) FeatureByOrgID(ctx context.Context, orgID string) (*Feature, error) {
	query, scan := prepareFeatureQuery()
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

func (q *Queries) DefaultFeature(ctx context.Context) (*Feature, error) {
	query, scan := prepareFeatureQuery()
	stmt, args, err := query.Where(sq.Eq{
		FeatureColumnAggregateID.identifier(): domain.IAMID,
	}).OrderBy(FeatureColumnIsDefault.identifier()).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-1Ndlg", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func prepareFeatureQuery() (sq.SelectBuilder, func(*sql.Row) (*Feature, error)) {
	tierName := sql.NullString{}
	tierDescription := sql.NullString{}
	stateDescription := sql.NullString{}
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
			FeatureActions.identifier(),
		).From(feautureTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Feature, error) {
			p := new(Feature)
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
				&p.Actions,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-M9fse", "Errors.Feature.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-3o9gd", "Errors.Internal")
			}
			p.TierName = tierName.String
			p.TierDescription = tierDescription.String
			p.StateDescription = stateDescription.String
			return p, nil
		}
}

func (f *Feature) FeatureList() []string {
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
	if f.Actions {
		list = append(list, domain.FeatureActions)
	}
	return list
}
