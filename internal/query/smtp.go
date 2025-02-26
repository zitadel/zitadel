package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SMTPConfigsSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

type SMTPConfigs struct {
	SearchResponse
	Configs []*SMTPConfig
}

var (
	smtpConfigsTable = table{
		name:          projection.SMTPConfigProjectionTable,
		instanceIDCol: projection.SMTPConfigColumnInstanceID,
	}
	SMTPConfigColumnCreationDate = Column{
		name:  projection.SMTPConfigColumnCreationDate,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnChangeDate = Column{
		name:  projection.SMTPConfigColumnChangeDate,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnResourceOwner = Column{
		name:  projection.SMTPConfigColumnResourceOwner,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnInstanceID = Column{
		name:  projection.SMTPConfigColumnInstanceID,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnSequence = Column{
		name:  projection.SMTPConfigColumnSequence,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnID = Column{
		name:  projection.SMTPConfigColumnID,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnState = Column{
		name:  projection.SMTPConfigColumnState,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnDescription = Column{
		name:  projection.SMTPConfigColumnDescription,
		table: smtpConfigsTable,
	}

	smtpConfigsSMTPTable = table{
		name:          projection.SMTPConfigTable,
		instanceIDCol: projection.SMTPConfigColumnInstanceID,
	}
	SMTPConfigSMTPColumnInstanceID = Column{
		name:  projection.SMTPConfigColumnInstanceID,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnID = Column{
		name:  projection.SMTPConfigColumnID,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnTLS = Column{
		name:  projection.SMTPConfigSMTPColumnTLS,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnSenderAddress = Column{
		name:  projection.SMTPConfigSMTPColumnSenderAddress,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnSenderName = Column{
		name:  projection.SMTPConfigSMTPColumnSenderName,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnReplyToAddress = Column{
		name:  projection.SMTPConfigSMTPColumnReplyToAddress,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnHost = Column{
		name:  projection.SMTPConfigSMTPColumnHost,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnUser = Column{
		name:  projection.SMTPConfigSMTPColumnUser,
		table: smtpConfigsSMTPTable,
	}
	SMTPConfigSMTPColumnPassword = Column{
		name:  projection.SMTPConfigSMTPColumnPassword,
		table: smtpConfigsSMTPTable,
	}

	smtpConfigsHTTPTable = table{
		name:          projection.SMTPConfigHTTPTable,
		instanceIDCol: projection.SMTPConfigHTTPColumnInstanceID,
	}
	SMTPConfigHTTPColumnInstanceID = Column{
		name:  projection.SMTPConfigHTTPColumnInstanceID,
		table: smtpConfigsHTTPTable,
	}
	SMTPConfigHTTPColumnID = Column{
		name:  projection.SMTPConfigHTTPColumnID,
		table: smtpConfigsHTTPTable,
	}
	SMTPConfigHTTPColumnEndpoint = Column{
		name:  projection.SMTPConfigHTTPColumnEndpoint,
		table: smtpConfigsHTTPTable,
	}
)

type SMTPConfig struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	AggregateID   string
	ID            string
	Sequence      uint64
	Description   string

	SMTPConfig *SMTP
	HTTPConfig *HTTP

	State domain.SMTPConfigState
}

type SMTP struct {
	TLS            bool
	SenderAddress  string
	SenderName     string
	ReplyToAddress string
	Host           string
	User           string
	Password       *crypto.CryptoValue
}

func (q *Queries) SMTPConfigActive(ctx context.Context, resourceOwner string) (config *SMTPConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareSMTPConfigQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		SMTPConfigColumnResourceOwner.identifier(): resourceOwner,
		SMTPConfigColumnInstanceID.identifier():    resourceOwner,
		SMTPConfigColumnState.identifier():         domain.SMTPConfigStateActive,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-3m9sl", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, query, args...)
	return config, err
}

func (q *Queries) SMTPConfigByID(ctx context.Context, instanceID, id string) (config *SMTPConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareSMTPConfigQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		SMTPConfigColumnInstanceID.identifier(): instanceID,
		SMTPConfigColumnID.identifier():         id,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-8f8gw", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, query, args...)
	return config, err
}

func prepareSMTPConfigQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*SMTPConfig, error)) {
	password := new(crypto.CryptoValue)

	return sq.Select(
			SMTPConfigColumnCreationDate.identifier(),
			SMTPConfigColumnChangeDate.identifier(),
			SMTPConfigColumnResourceOwner.identifier(),
			SMTPConfigColumnSequence.identifier(),
			SMTPConfigColumnID.identifier(),
			SMTPConfigColumnState.identifier(),
			SMTPConfigColumnDescription.identifier(),

			SMTPConfigSMTPColumnID.identifier(),
			SMTPConfigSMTPColumnTLS.identifier(),
			SMTPConfigSMTPColumnSenderAddress.identifier(),
			SMTPConfigSMTPColumnSenderName.identifier(),
			SMTPConfigSMTPColumnReplyToAddress.identifier(),
			SMTPConfigSMTPColumnHost.identifier(),
			SMTPConfigSMTPColumnUser.identifier(),
			SMTPConfigSMTPColumnPassword.identifier(),

			SMTPConfigHTTPColumnID.identifier(),
			SMTPConfigHTTPColumnEndpoint.identifier()).
			From(smtpConfigsTable.identifier()).
			LeftJoin(join(SMTPConfigSMTPColumnID, SMTPConfigColumnID)).
			LeftJoin(join(SMTPConfigHTTPColumnID, SMTPConfigColumnID) + db.Timetravel(call.Took(ctx))).
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
				&password,
				&httpConfig.id,
				&httpConfig.endpoint,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-fwofw", "Errors.SMTPConfig.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-9k87F", "Errors.Internal")
			}
			smtpConfig.password = password
			smtpConfig.set(config)
			httpConfig.setSMTP(config)
			return config, nil
		}
}

func prepareSMTPConfigsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*SMTPConfigs, error)) {
	return sq.Select(
			SMTPConfigColumnCreationDate.identifier(),
			SMTPConfigColumnChangeDate.identifier(),
			SMTPConfigColumnResourceOwner.identifier(),
			SMTPConfigColumnSequence.identifier(),
			SMTPConfigColumnID.identifier(),
			SMTPConfigColumnState.identifier(),
			SMTPConfigColumnDescription.identifier(),

			SMTPConfigSMTPColumnID.identifier(),
			SMTPConfigSMTPColumnTLS.identifier(),
			SMTPConfigSMTPColumnSenderAddress.identifier(),
			SMTPConfigSMTPColumnSenderName.identifier(),
			SMTPConfigSMTPColumnReplyToAddress.identifier(),
			SMTPConfigSMTPColumnHost.identifier(),
			SMTPConfigSMTPColumnUser.identifier(),
			SMTPConfigSMTPColumnPassword.identifier(),

			SMTPConfigHTTPColumnID.identifier(),
			SMTPConfigHTTPColumnEndpoint.identifier(),
			countColumn.identifier(),
		).From(smtpConfigsTable.identifier()).
			LeftJoin(join(SMTPConfigSMTPColumnID, SMTPConfigColumnID)).
			LeftJoin(join(SMTPConfigHTTPColumnID, SMTPConfigColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*SMTPConfigs, error) {
			configs := &SMTPConfigs{Configs: []*SMTPConfig{}}
			for rows.Next() {
				config := new(SMTPConfig)
				password := new(crypto.CryptoValue)
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
					&password,
					&httpConfig.id,
					&httpConfig.endpoint,
					&configs.Count,
				)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return nil, zerrors.ThrowNotFound(err, "QUERY-fwofw", "Errors.SMTPConfig.NotFound")
					}
					return nil, zerrors.ThrowInternal(err, "QUERY-9k87F", "Errors.Internal")
				}
				smtpConfig.password = password
				smtpConfig.set(config)
				httpConfig.setSMTP(config)
				configs.Configs = append(configs.Configs, config)
			}
			return configs, nil
		}
}

func (q *Queries) SearchSMTPConfigs(ctx context.Context, queries *SMTPConfigsSearchQueries) (configs *SMTPConfigs, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMTPConfigsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			SMTPConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-sZ7Cx", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		configs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-tOpKN", "Errors.Internal")
	}
	configs.State, err = q.latestState(ctx, smsConfigsTable)
	return configs, err
}

type sqlSmtpConfig struct {
	id             sql.NullString
	tls            sql.NullBool
	senderAddress  sql.NullString
	senderName     sql.NullString
	replyToAddress sql.NullString
	host           sql.NullString
	user           sql.NullString
	password       *crypto.CryptoValue
}

func (c sqlSmtpConfig) set(smtpConfig *SMTPConfig) {
	if !c.id.Valid {
		return
	}
	smtpConfig.SMTPConfig = &SMTP{
		TLS:            c.tls.Bool,
		SenderAddress:  c.senderAddress.String,
		SenderName:     c.senderName.String,
		ReplyToAddress: c.replyToAddress.String,
		Host:           c.host.String,
		User:           c.user.String,
		Password:       c.password,
	}
}

func (c sqlHTTPConfig) setSMTP(smtpConfig *SMTPConfig) {
	if !c.id.Valid {
		return
	}
	smtpConfig.HTTPConfig = &HTTP{
		Endpoint: c.endpoint.String,
	}
}
