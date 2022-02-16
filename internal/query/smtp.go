package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	SMTPConfigColumnSequence = Column{
		name:  projection.SMTPConfigColumnSequence,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnTls = Column{
		name:  projection.SMTPConfigColumnTls,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnFromAddress = Column{
		name:  projection.SMTPConfigColumnFromAddress,
		table: smtpConfigsTable,
	}
	SMTPConfigColumnFromName = Column{
		name:  projection.SMTPConfigColumnFromName,
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

	TLS          bool
	FromAddress  string
	FromName     string
	SMTPHost     string
	SMTPUser     string
	SMTPPassword *crypto.CryptoValue
}

func (q *Queries) SMTPConfigByAggregateID(ctx context.Context, aggregateID string) (*SMTPConfig, error) {
	stmt, scan := prepareSMTPConfigQuery()
	query, args, err := stmt.Where(sq.Eq{
		SMTPConfigColumnAggregateID.identifier(): aggregateID,
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
			SMTPConfigColumnTls.identifier(),
			SMTPConfigColumnFromAddress.identifier(),
			SMTPConfigColumnFromName.identifier(),
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
				&config.FromAddress,
				&config.FromName,
				&config.SMTPHost,
				&config.SMTPUser,
				&password,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fwofw", "Errors.SMTPConfig.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-9k87F", "Errors.Internal")
			}
			config.SMTPPassword = password
			return config, nil
		}
}
