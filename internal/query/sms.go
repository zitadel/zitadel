package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

type SMSConfigs struct {
	SearchResponse
	Configs []*SMSConfig
}

type SMSConfig struct {
	AggregateID   string
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.SMSConfigState
	Sequence      uint64

	TwilioConfig *Twilio
}

type Twilio struct {
	SID          string
	Token        *crypto.CryptoValue
	SenderNumber string
}

type SMSConfigsSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *SMSConfigsSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

var (
	smsConfigsTable = table{
		name: projection.SMSConfigProjectionTable,
	}
	SMSConfigColumnID = Column{
		name:  projection.SMSColumnID,
		table: smsConfigsTable,
	}
	SMSConfigColumnAggregateID = Column{
		name:  projection.SMSColumnAggregateID,
		table: smsConfigsTable,
	}
	SMSConfigColumnCreationDate = Column{
		name:  projection.SMSColumnCreationDate,
		table: smsConfigsTable,
	}
	SMSConfigColumnChangeDate = Column{
		name:  projection.SMSColumnChangeDate,
		table: smsConfigsTable,
	}
	SMSConfigColumnResourceOwner = Column{
		name:  projection.SMSColumnResourceOwner,
		table: smsConfigsTable,
	}
	SMSConfigColumnState = Column{
		name:  projection.SMSColumnState,
		table: smsConfigsTable,
	}
	SMSConfigColumnSequence = Column{
		name:  projection.SMSColumnSequence,
		table: smsConfigsTable,
	}
)

var (
	smsTwilioConfigsTable = table{
		name: projection.SMSTwilioTable,
	}
	SMSTwilioConfigColumnSMSID = Column{
		name:  projection.SMSTwilioConfigColumnSMSID,
		table: smsTwilioConfigsTable,
	}
	SMSTwilioConfigColumnSID = Column{
		name:  projection.SMSTwilioConfigColumnSID,
		table: smsTwilioConfigsTable,
	}
	SMSTwilioConfigColumnToken = Column{
		name:  projection.SMSTwilioConfigColumnToken,
		table: smsTwilioConfigsTable,
	}
	SMSTwilioConfigColumnSenderNumber = Column{
		name:  projection.SMSTwilioConfigColumnSenderNumber,
		table: smsTwilioConfigsTable,
	}
)

func (q *Queries) SMSProviderConfigByID(ctx context.Context, id string) (*SMSConfig, error) {
	stmt, scan := prepareSMSConfigQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			SMSConfigColumnID.identifier(): id,
		},
	).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) SearchSMSConfigs(ctx context.Context, queries *SMSConfigsSearchQueries) (*SMSConfigs, error) {
	query, scan := prepareSMSConfigsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-sn9Jf", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-aJnZL", "Errors.Internal")
	}
	apps, err := scan(rows)
	if err != nil {
		return nil, err
	}
	apps.LatestSequence, err = q.latestSequence(ctx, smsConfigsTable)
	return apps, err
}

func prepareSMSConfigQuery() (sq.SelectBuilder, func(*sql.Row) (*SMSConfig, error)) {
	return sq.Select(
			SMSConfigColumnID.identifier(),
			SMSConfigColumnAggregateID.identifier(),
			SMSConfigColumnCreationDate.identifier(),
			SMSConfigColumnChangeDate.identifier(),
			SMSConfigColumnResourceOwner.identifier(),
			SMSConfigColumnState.identifier(),
			SMSConfigColumnSequence.identifier(),

			SMSTwilioConfigColumnSMSID.identifier(),
			SMSTwilioConfigColumnSID.identifier(),
			SMSTwilioConfigColumnToken.identifier(),
			SMSTwilioConfigColumnSenderNumber.identifier(),
		).From(smsConfigsTable.identifier()).
			LeftJoin(join(SMSTwilioConfigColumnSMSID, SMSConfigColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*SMSConfig, error) {
			config := new(SMSConfig)

			var (
				twilioConfig = sqlTwilioConfig{}
			)

			err := row.Scan(
				&config.ID,
				&config.AggregateID,
				&config.CreationDate,
				&config.ChangeDate,
				&config.ResourceOwner,
				&config.State,
				&config.Sequence,

				&twilioConfig.smsID,
				&twilioConfig.sid,
				&twilioConfig.token,
				&twilioConfig.senderNumber,
			)

			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fn99w", "Errors.SMSConfig.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-3n9Js", "Errors.Internal")
			}

			twilioConfig.set(config)

			return config, nil
		}
}

func prepareSMSConfigsQuery() (sq.SelectBuilder, func(*sql.Rows) (*SMSConfigs, error)) {
	return sq.Select(
			SMSConfigColumnID.identifier(),
			SMSConfigColumnAggregateID.identifier(),
			SMSConfigColumnCreationDate.identifier(),
			SMSConfigColumnChangeDate.identifier(),
			SMSConfigColumnResourceOwner.identifier(),
			SMSConfigColumnState.identifier(),
			SMSConfigColumnSequence.identifier(),

			SMSTwilioConfigColumnSMSID.identifier(),
			SMSTwilioConfigColumnSID.identifier(),
			SMSTwilioConfigColumnToken.identifier(),
			SMSTwilioConfigColumnSenderNumber.identifier(),
			countColumn.identifier(),
		).From(smsConfigsTable.identifier()).
			LeftJoin(join(SMSTwilioConfigColumnSMSID, SMSConfigColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Rows) (*SMSConfigs, error) {
			configs := &SMSConfigs{Configs: []*SMSConfig{}}

			for row.Next() {
				config := new(SMSConfig)
				var (
					twilioConfig = sqlTwilioConfig{}
				)

				err := row.Scan(
					&config.ID,
					&config.AggregateID,
					&config.CreationDate,
					&config.ChangeDate,
					&config.ResourceOwner,
					&config.State,
					&config.Sequence,

					&twilioConfig.smsID,
					&twilioConfig.sid,
					&twilioConfig.token,
					&twilioConfig.senderNumber,
					&configs.Count,
				)

				if err != nil {
					return nil, errors.ThrowInternal(err, "QUERY-d9jJd", "Errors.Internal")
				}

				twilioConfig.set(config)

				configs.Configs = append(configs.Configs, config)
			}

			return configs, nil
		}
}

type sqlTwilioConfig struct {
	smsID        sql.NullString
	sid          sql.NullString
	token        *crypto.CryptoValue
	senderNumber sql.NullString
}

func (c sqlTwilioConfig) set(smsConfig *SMSConfig) {
	if !c.smsID.Valid {
		return
	}
	smsConfig.TwilioConfig = &Twilio{
		SID:          c.sid.String,
		Token:        c.token,
		SenderNumber: c.senderNumber.String,
	}
}
