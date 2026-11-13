package query

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	orgSMTPConfigsTable = table{
		name:          projection.OrgSMTPConfigProjectionTable,
		instanceIDCol: projection.OrgSMTPConfigColumnInstanceID,
	}
	OrgSMTPConfigColumnCreationDate = Column{
		name:  projection.OrgSMTPConfigColumnCreationDate,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnChangeDate = Column{
		name:  projection.OrgSMTPConfigColumnChangeDate,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnResourceOwner = Column{
		name:  projection.OrgSMTPConfigColumnResourceOwner,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnInstanceID = Column{
		name:  projection.OrgSMTPConfigColumnInstanceID,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnSequence = Column{
		name:  projection.OrgSMTPConfigColumnSequence,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnID = Column{
		name:  projection.OrgSMTPConfigColumnID,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnState = Column{
		name:  projection.OrgSMTPConfigColumnState,
		table: orgSMTPConfigsTable,
	}
	OrgSMTPConfigColumnDescription = Column{
		name:  projection.OrgSMTPConfigColumnDescription,
		table: orgSMTPConfigsTable,
	}

	orgSMTPConfigsSMTPTable = table{
		name:          projection.OrgSMTPConfigTable,
		instanceIDCol: projection.OrgSMTPConfigSMTPColumnInstanceID,
	}
	OrgSMTPConfigSMTPColumnInstanceID = Column{
		name:  projection.OrgSMTPConfigSMTPColumnInstanceID,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnID = Column{
		name:  projection.OrgSMTPConfigSMTPColumnID,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnTLS = Column{
		name:  projection.OrgSMTPConfigSMTPColumnTLS,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnSenderAddress = Column{
		name:  projection.OrgSMTPConfigSMTPColumnSenderAddress,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnSenderName = Column{
		name:  projection.OrgSMTPConfigSMTPColumnSenderName,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnReplyToAddress = Column{
		name:  projection.OrgSMTPConfigSMTPColumnReplyToAddress,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnHost = Column{
		name:  projection.OrgSMTPConfigSMTPColumnHost,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnUser = Column{
		name:  projection.OrgSMTPConfigSMTPColumnUser,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnPlainAuthPassword = Column{
		name:  projection.OrgSMTPConfigSMTPColumnPlainAuthPassword,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId = Column{
		name:  projection.OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret = Column{
		name:  projection.OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint = Column{
		name:  projection.OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint,
		table: orgSMTPConfigsSMTPTable,
	}
	OrgSMTPConfigSMTPColumnXOAuth2AuthScope = Column{
		name:  projection.OrgSMTPConfigSMTPColumnXOAuth2AuthScope,
		table: orgSMTPConfigsSMTPTable,
	}

	orgSMTPConfigsHTTPTable = table{
		name:          projection.OrgSMTPConfigHTTPTable,
		instanceIDCol: projection.OrgSMTPConfigHTTPColumnInstanceID,
	}
	OrgSMTPConfigHTTPColumnInstanceID = Column{
		name:  projection.OrgSMTPConfigHTTPColumnInstanceID,
		table: orgSMTPConfigsHTTPTable,
	}
	OrgSMTPConfigHTTPColumnID = Column{
		name:  projection.OrgSMTPConfigHTTPColumnID,
		table: orgSMTPConfigsHTTPTable,
	}
	OrgSMTPConfigHTTPColumnEndpoint = Column{
		name:  projection.OrgSMTPConfigHTTPColumnEndpoint,
		table: orgSMTPConfigsHTTPTable,
	}
	OrgSMTPConfigHTTPColumnSigningKey = Column{
		name:  projection.OrgSMTPConfigHTTPColumnSigningKey,
		table: orgSMTPConfigsHTTPTable,
	}
)

func (q *Queries) OrgSMTPConfigActive(ctx context.Context, orgID string) (config *SMTPConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareOrgSMTPConfigQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgSMTPConfigColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
		OrgSMTPConfigColumnResourceOwner.identifier(): orgID,
		OrgSMTPConfigColumnState.identifier():         domain.SMTPConfigStateActive,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-org3m9", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, query, args...)
	if err != nil {
		return nil, err
	}
	if err := config.decryptSigningKey(q.smtpEncryptionAlgorithm); err != nil {
		return nil, err
	}
	return config, nil
}

func (q *Queries) OrgSMTPConfigByID(ctx context.Context, orgID, id string) (config *SMTPConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareOrgSMTPConfigQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgSMTPConfigColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
		OrgSMTPConfigColumnResourceOwner.identifier(): orgID,
		OrgSMTPConfigColumnID.identifier():            id,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-org8f9", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, query, args...)
	if err != nil {
		return nil, err
	}
	if err := config.decryptSigningKey(q.smtpEncryptionAlgorithm); err != nil {
		return nil, err
	}
	return config, nil
}

func (q *Queries) SearchOrgSMTPConfigs(ctx context.Context, orgID string, queries *SMTPConfigsSearchQueries) (configs *SMTPConfigs, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareOrgSMTPConfigsQuery()
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			OrgSMTPConfigColumnInstanceID.identifier():    authz.GetInstance(ctx).InstanceID(),
			OrgSMTPConfigColumnResourceOwner.identifier(): orgID,
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-orgSz7", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		configs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-orgTpK", "Errors.Internal")
	}
	for i := range configs.Configs {
		if err := configs.Configs[i].decryptSigningKey(q.smtpEncryptionAlgorithm); err != nil {
			return nil, err
		}
	}
	configs.State, err = q.latestState(ctx, orgSMTPConfigsTable)
	return configs, err
}

func prepareOrgSMTPConfigQuery() (sq.SelectBuilder, func(*sql.Row) (*SMTPConfig, error)) {
	return sq.Select(
			OrgSMTPConfigColumnCreationDate.identifier(),
			OrgSMTPConfigColumnChangeDate.identifier(),
			OrgSMTPConfigColumnResourceOwner.identifier(),
			OrgSMTPConfigColumnSequence.identifier(),
			OrgSMTPConfigColumnID.identifier(),
			OrgSMTPConfigColumnState.identifier(),
			OrgSMTPConfigColumnDescription.identifier(),

			OrgSMTPConfigSMTPColumnID.identifier(),
			OrgSMTPConfigSMTPColumnTLS.identifier(),
			OrgSMTPConfigSMTPColumnSenderAddress.identifier(),
			OrgSMTPConfigSMTPColumnSenderName.identifier(),
			OrgSMTPConfigSMTPColumnReplyToAddress.identifier(),
			OrgSMTPConfigSMTPColumnHost.identifier(),
			OrgSMTPConfigSMTPColumnUser.identifier(),
			OrgSMTPConfigSMTPColumnPlainAuthPassword.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthScope.identifier(),

			OrgSMTPConfigHTTPColumnID.identifier(),
			OrgSMTPConfigHTTPColumnEndpoint.identifier(),
			OrgSMTPConfigHTTPColumnSigningKey.identifier(),
		).
			From(orgSMTPConfigsTable.identifier()).
			LeftJoin(join(OrgSMTPConfigSMTPColumnID, OrgSMTPConfigColumnID)).
			LeftJoin(join(OrgSMTPConfigHTTPColumnID, OrgSMTPConfigColumnID)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SMTPConfig, error) {
			config := new(SMTPConfig)
			var (
				smtpConfig = sqlSmtpConfig{}
				httpConfig = sqlHTTPConfig{}
			)
			err := row.Scan(
				&config.CreationDate,
				&config.ChangeDate,
				&config.ResourceOwner,
				&config.Sequence,
				&config.ID,
				&config.State,
				&config.Description,
				&smtpConfig.id,
				&smtpConfig.tls,
				&smtpConfig.senderAddress,
				&smtpConfig.senderName,
				&smtpConfig.replyToAddress,
				&smtpConfig.host,
				&smtpConfig.user,
				&smtpConfig.plainAuthPassword,
				&smtpConfig.xoauth2AuthClientCredentialsClientId,
				&smtpConfig.xoauth2AuthClientClientCredentialsSecret,
				&smtpConfig.xoauth2AuthTokenEndpoint,
				&smtpConfig.xoauth2AuthScope,
				&httpConfig.id,
				&httpConfig.endpoint,
				&httpConfig.signingKey,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-orgFw", "Errors.SMTPConfig.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-org9kF", "Errors.Internal")
			}

			smtpConfig.set(config)
			httpConfig.setSMTP(config)

			return config, nil
		}
}

func prepareOrgSMTPConfigsQuery() (sq.SelectBuilder, func(*sql.Rows) (*SMTPConfigs, error)) {
	return sq.Select(
			OrgSMTPConfigColumnCreationDate.identifier(),
			OrgSMTPConfigColumnChangeDate.identifier(),
			OrgSMTPConfigColumnResourceOwner.identifier(),
			OrgSMTPConfigColumnSequence.identifier(),
			OrgSMTPConfigColumnID.identifier(),
			OrgSMTPConfigColumnState.identifier(),
			OrgSMTPConfigColumnDescription.identifier(),

			OrgSMTPConfigSMTPColumnID.identifier(),
			OrgSMTPConfigSMTPColumnTLS.identifier(),
			OrgSMTPConfigSMTPColumnSenderAddress.identifier(),
			OrgSMTPConfigSMTPColumnSenderName.identifier(),
			OrgSMTPConfigSMTPColumnReplyToAddress.identifier(),
			OrgSMTPConfigSMTPColumnHost.identifier(),
			OrgSMTPConfigSMTPColumnUser.identifier(),
			OrgSMTPConfigSMTPColumnPlainAuthPassword.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientId.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthClientCredentialsClientSecret.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthTokenEndpoint.identifier(),
			OrgSMTPConfigSMTPColumnXOAuth2AuthScope.identifier(),

			OrgSMTPConfigHTTPColumnID.identifier(),
			OrgSMTPConfigHTTPColumnEndpoint.identifier(),
			OrgSMTPConfigHTTPColumnSigningKey.identifier(),

			countColumn.identifier(),
		).From(orgSMTPConfigsTable.identifier()).
			LeftJoin(join(OrgSMTPConfigSMTPColumnID, OrgSMTPConfigColumnID)).
			LeftJoin(join(OrgSMTPConfigHTTPColumnID, OrgSMTPConfigColumnID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*SMTPConfigs, error) {
			configs := &SMTPConfigs{Configs: []*SMTPConfig{}}
			for rows.Next() {
				config := new(SMTPConfig)
				var (
					smtpConfig = sqlSmtpConfig{}
					httpConfig = sqlHTTPConfig{}
				)
				err := rows.Scan(
					&config.CreationDate,
					&config.ChangeDate,
					&config.ResourceOwner,
					&config.Sequence,
					&config.ID,
					&config.State,
					&config.Description,
					&smtpConfig.id,
					&smtpConfig.tls,
					&smtpConfig.senderAddress,
					&smtpConfig.senderName,
					&smtpConfig.replyToAddress,
					&smtpConfig.host,
					&smtpConfig.user,
					&smtpConfig.plainAuthPassword,
					&smtpConfig.xoauth2AuthClientCredentialsClientId,
					&smtpConfig.xoauth2AuthClientClientCredentialsSecret,
					&smtpConfig.xoauth2AuthTokenEndpoint,
					&smtpConfig.xoauth2AuthScope,
					&httpConfig.id,
					&httpConfig.endpoint,
					&httpConfig.signingKey,
					&configs.Count,
				)
				if err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-org9kR", "Errors.Internal")
				}

				smtpConfig.set(config)
				httpConfig.setSMTP(config)

				configs.Configs = append(configs.Configs, config)
			}
			return configs, nil
		}
}
