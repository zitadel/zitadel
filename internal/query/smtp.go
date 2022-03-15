package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/query/projection"

	"github.com/caos/zitadel/internal/errors"
)

var (
	smtpConfigsTable = table{
		name: projection.SMTPConfigProjectionTable,
	}
	SMTPConfigColumnAggregateID = Column{
		name:  projection.SMTPConfigColumnAggregateID,
		table: smtpConfigsTable,
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
	SMTPConfigColumnTLS = Column{
		name:  projection.SMTPConfigColumnTLS,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnSenderAddress = Column{
		name:  projection.SMTPConfigColumnSenderAddress,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnSenderName = Column{
		name:  projection.SMTPConfigColumnSenderName,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnSMTPHost = Column{
		name:  projection.SMTPConfigColumnSMTPHost,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnSMTPUser = Column{
		name:  projection.SMTPConfigColumnSMTPUser,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnSMTPPassword = Column{
		name:  projection.SMTPConfigColumnSMTPPassword,
		table: smtpConfigsTable,
	}
)

type SMTPConfigs struct {
	SearchResponse
	SMTPConfigs []*SMTPConfig
}

type SMTPConfig struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	TLS           bool
	SenderAddress string
	SenderName    string
	Host          string
	User          string
	Password      *crypto.CryptoValue
}

func (q *Queries) SMTPConfigByAggregateID(ctx context.Context, aggregateID string) (*SMTPConfig, error) {
	stmt, scan := prepareSMTPConfigQuery()
	query, args, err := stmt.Where(sq.Eq{
		SMTPConfigColumnAggregateID.identifier(): aggregateID,
		SMTPConfigColumnInstanceID.identifier():  authz.GetCtxData(ctx).InstanceID,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3m9sl", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareSMTPConfigQuery() (sq.SelectBuilder, func(*sql.Row) (*SMTPConfig, error)) {
	password := new(crypto.CryptoValue)

	return sq.Select(
			SMTPConfigColumnAggregateID.identifier(),
			SMTPConfigColumnCreationDate.identifier(),
			SMTPConfigColumnChangeDate.identifier(),
			SMTPConfigColumnResourceOwner.identifier(),
			SMTPConfigColumnSequence.identifier(),
			SMTPConfigColumnTLS.identifier(),
			SMTPConfigColumnSenderAddress.identifier(),
			SMTPConfigColumnSenderName.identifier(),
			SMTPConfigColumnSMTPHost.identifier(),
			SMTPConfigColumnSMTPUser.identifier(),
			SMTPConfigColumnSMTPPassword.identifier()).
			From(smtpConfigsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*SMTPConfig, error) {
			config := new(SMTPConfig)
			err := row.Scan(
				&config.AggregateID,
				&config.CreationDate,
				&config.ChangeDate,
				&config.ResourceOwner,
				&config.Sequence,
				&config.TLS,
				&config.SenderAddress,
				&config.SenderName,
				&config.Host,
				&config.User,
				&password,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fwofw", "Errors.SMTPConfig.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-9k87F", "Errors.Internal")
			}
			config.Password = password
			return config, nil
		}
}
