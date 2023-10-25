package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	emailWebhookConfigsTable = table{
		name:          projection.EmailWebhookConfigProjectionTable,
		instanceIDCol: projection.EmailWebhookConfigColumnInstanceID,
	}
	EmailWebhookConfigColumnAggregateID = Column{
		name:  projection.EmailWebhookConfigColumnAggregateID,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnCreationDate = Column{
		name:  projection.EmailWebhookConfigColumnCreationDate,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnChangeDate = Column{
		name:  projection.EmailWebhookConfigColumnChangeDate,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnResourceOwner = Column{
		name:  projection.EmailWebhookConfigColumnResourceOwner,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnInstanceID = Column{
		name:  projection.EmailWebhookConfigColumnInstanceID,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnSequence = Column{
		name:  projection.EmailWebhookConfigColumnSequence,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnTLS = Column{
		name:  projection.EmailWebhookConfigColumnTLS,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnSenderAddress = Column{
		name:  projection.EmailWebhookConfigColumnSenderAddress,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnSenderName = Column{
		name:  projection.EmailWebhookConfigColumnSenderName,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnReplyToAddress = Column{
		name:  projection.EmailWebhookConfigColumnReplyToAddress,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnEmailWebhookHost = Column{
		name:  projection.EmailWebhookConfigColumnEmailWebhookHost,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnEmailWebhookUser = Column{
		name:  projection.EmailWebhookConfigColumnEmailWebhookUser,
		table: emailWebhookConfigsTable,
	}
	EmailWebhookConfigColumnEmailWebhookPassword = Column{
		name:  projection.EmailWebhookConfigColumnEmailWebhookPassword,
		table: emailWebhookConfigsTable,
	}
)

type EmailWebhookConfigs struct {
	SearchResponse
	EmailWebhookConfigs []*EmailWebhookConfig
}

type EmailWebhookConfig struct {
	AggregateID   string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	TLS            bool
	SenderAddress  string
	SenderName     string
	ReplyToAddress string
	Host           string
	User           string
	Password       *crypto.CryptoValue
}

func (q *Queries) EmailWebhookConfigByAggregateID(ctx context.Context, aggregateID string) (config *EmailWebhookConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareEmailWebhookConfigQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		EmailWebhookConfigColumnAggregateID.identifier(): aggregateID,
		EmailWebhookConfigColumnInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3m9sl", "Errors.Query.SQLStatment")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, query, args...)
	return config, err
}

func prepareEmailWebhookConfigQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*EmailWebhookConfig, error)) {
	password := new(crypto.CryptoValue)

	return sq.Select(
			EmailWebhookConfigColumnAggregateID.identifier(),
			EmailWebhookConfigColumnCreationDate.identifier(),
			EmailWebhookConfigColumnChangeDate.identifier(),
			EmailWebhookConfigColumnResourceOwner.identifier(),
			EmailWebhookConfigColumnSequence.identifier(),
			EmailWebhookConfigColumnTLS.identifier(),
			EmailWebhookConfigColumnSenderAddress.identifier(),
			EmailWebhookConfigColumnSenderName.identifier(),
			EmailWebhookConfigColumnReplyToAddress.identifier(),
			EmailWebhookConfigColumnEmailWebhookHost.identifier(),
			EmailWebhookConfigColumnEmailWebhookUser.identifier(),
			EmailWebhookConfigColumnEmailWebhookPassword.identifier()).
			From(emailWebhookConfigsTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*EmailWebhookConfig, error) {
			config := new(EmailWebhookConfig)
			err := row.Scan(
				&config.AggregateID,
				&config.CreationDate,
				&config.ChangeDate,
				&config.ResourceOwner,
				&config.Sequence,
				&config.TLS,
				&config.SenderAddress,
				&config.SenderName,
				&config.ReplyToAddress,
				&config.Host,
				&config.User,
				&password,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fwofw", "Errors.EmailWebhookConfig.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-9k87F", "Errors.Internal")
			}
			config.Password = password
			return config, nil
		}
}
