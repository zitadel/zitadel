package query

import (
	"context"
	"database/sql"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// OIDCAppLinkConfig holds iOS/Android association fields for one active OIDC app.
type OIDCAppLinkConfig struct {
	AppID                         string
	IOSTeamID                     string
	IOSBundleID                   string
	AndroidPackageName            string
	AndroidSHA256CertFingerprints []string
}

// SearchOIDCAppLinkConfigs returns association fields for active OIDC apps in the
// current instance that have at least one iOS or Android app-link field set.
func (q *Queries) SearchOIDCAppLinkConfigs(ctx context.Context) (configs []*OIDCAppLinkConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareOIDCAppLinkConfigsQuery()
	iosTeamIDCol := AppOIDCConfigColumnIOSTeamID.identifier()
	iosBundleIDCol := AppOIDCConfigColumnIOSBundleID.identifier()
	androidPackageCol := AppOIDCConfigColumnAndroidPackageName.identifier()
	stmt, args, err := query.Where(sq.And{
		sq.Eq{
			AppOIDCConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			AppColumnState.identifier():                domain.AppStateActive,
		},
		sq.Or{
			sq.Expr("COALESCE(" + iosTeamIDCol + ", '') <> ''"),
			sq.Expr("COALESCE(" + iosBundleIDCol + ", '') <> ''"),
			sq.Expr("COALESCE(" + androidPackageCol + ", '') <> ''"),
		},
	}).OrderBy(AppOIDCConfigColumnAppID.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-ApLn1", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		configs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-ApLn2", "Errors.Internal")
	}
	return configs, nil
}

func prepareOIDCAppLinkConfigsQuery() (sq.SelectBuilder, func(*sql.Rows) ([]*OIDCAppLinkConfig, error)) {
	return sq.Select(
			AppOIDCConfigColumnAppID.identifier(),
			AppOIDCConfigColumnIOSTeamID.identifier(),
			AppOIDCConfigColumnIOSBundleID.identifier(),
			AppOIDCConfigColumnAndroidPackageName.identifier(),
			AppOIDCConfigColumnAndroidSHA256CertFingerprints.identifier(),
		).From(appOIDCConfigsTable.identifier()).
			Join(join(AppColumnID, AppOIDCConfigColumnAppID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) ([]*OIDCAppLinkConfig, error) {
			configs := make([]*OIDCAppLinkConfig, 0)
			for rows.Next() {
				var (
					appID          sql.NullString
					iosTeamID      sql.NullString
					iosBundleID    sql.NullString
					androidPackage sql.NullString
					fingerprints   database.TextArray[string]
				)
				if err := rows.Scan(
					&appID,
					&iosTeamID,
					&iosBundleID,
					&androidPackage,
					&fingerprints,
				); err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-ApLn3", "Errors.Internal")
				}
				configs = append(configs, &OIDCAppLinkConfig{
					AppID:                         appID.String,
					IOSTeamID:                     strings.TrimSpace(iosTeamID.String),
					IOSBundleID:                   strings.TrimSpace(iosBundleID.String),
					AndroidPackageName:            strings.TrimSpace(androidPackage.String),
					AndroidSHA256CertFingerprints: fingerprints,
				})
			}
			return configs, nil
		}
}
