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
	CreationDate             time.Time
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
		name: projection.FeatureAggregateIDCol,
	}
	FeatureColumnCreationDate = Column{
		name: projection.FeatureCreationDateCol,
	}
	FeatureColumnChangeDate = Column{
		name: projection.FeatureChangeDateCol,
	}
	FeatureColumnSequence = Column{
		name: projection.FeatureSequenceCol,
	}
	FeatureColumnIsDefault = Column{
		name: projection.FeatureIsDefaultCol,
	}
	FeatureTierName = Column{
		name: projection.FeatureTierNameCol,
	}
	FeatureTierDescription = Column{
		name: projection.FeatureTierDescriptionCol,
	}
	FeatureState = Column{
		name: projection.FeatureStateCol,
	}
	FeatureStateDescription = Column{
		name: projection.FeatureStateDescriptionCol,
	}
	FeatureAuditLogRetention = Column{
		name: projection.FeatureAuditLogRetentionCol,
	}
	FeatureLoginPolicyFactors = Column{
		name: projection.FeatureLoginPolicyFactorsCol,
	}
	FeatureLoginPolicyIDP = Column{
		name: projection.FeatureLoginPolicyIDPCol,
	}
	FeatureLoginPolicyPasswordless = Column{
		name: projection.FeatureLoginPolicyPasswordlessCol,
	}
	FeatureLoginPolicyRegistration = Column{
		name: projection.FeatureLoginPolicyRegistrationCol,
	}
	FeatureLoginPolicyUsernameLogin = Column{
		name: projection.FeatureLoginPolicyUsernameLoginCol,
	}
	FeatureLoginPolicyPasswordReset = Column{
		name: projection.FeatureLoginPolicyPasswordResetCol,
	}
	FeaturePasswordComplexityPolicy = Column{
		name: projection.FeaturePasswordComplexityPolicyCol,
	}
	FeatureLabelPolicyPrivateLabel = Column{
		name: projection.FeatureLabelPolicyPrivateLabelCol,
	}
	FeatureLabelPolicyWatermark = Column{
		name: projection.FeatureLabelPolicyWatermarkCol,
	}
	FeatureCustomDomain = Column{
		name: projection.FeatureCustomDomainCol,
	}
	FeaturePrivacyPolicy = Column{
		name: projection.FeaturePrivacyPolicyCol,
	}
	FeatureMetadataUser = Column{
		name: projection.FeatureMetadataUserCol,
	}
	FeatureCustomTextMessage = Column{
		name: projection.FeatureCustomTextMessageCol,
	}
	FeatureCustomTextLogin = Column{
		name: projection.FeatureCustomTextLoginCol,
	}
	FeatureLockoutPolicy = Column{
		name: projection.FeatureLockoutPolicyCol,
	}
	FeatureActions = Column{
		name: projection.FeatureActionsCol,
	}
)

func (q *Queries) FeatureByID(ctx context.Context, orgID string) (*Feature, error) {
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
	return sq.Select(
			FeatureColumnAggregateID.identifier(),
			FeatureColumnCreationDate.identifier(),
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
		).From(loginPolicyTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Feature, error) {
			p := new(Feature)
			err := row.Scan(
				&p.AggregateID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.Sequence,
				&p.IsDefault,
				&p.TierName,
				&p.TierDescription,
				&p.State,
				&p.StateDescription,
				&p.AuditLogRetention,
				&p.LoginPolicyFactors,
				&p.LoginPolicyIDP,
				&p.LoginPolicyPasswordless,
				&p.LoginPolicyRegistration,
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
			return p, nil
		}
}
