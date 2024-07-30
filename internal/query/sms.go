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
	SID              string
	Token            *crypto.CryptoValue
	SenderNumber     string
	VerifyServiceSID string
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
		name:          projection.SMSConfigProjectionTable,
		instanceIDCol: projection.SMSColumnInstanceID,
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
	SMSConfigColumnInstanceID = Column{
		name:  projection.SMSColumnInstanceID,
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
		name:          projection.SMSTwilioTable,
		instanceIDCol: projection.SMSTwilioColumnInstanceID,
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
	SMSTwilioConfigColumnVerifyServiceSID = Column{
		name:  projection.SMSTwilioConfigColumnVerifyServiceSID,
		table: smsTwilioConfigsTable,
	}
)

func (q *Queries) SMSProviderConfigByID(ctx context.Context, id string) (config *SMSConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMSConfigQuery(ctx, q.client)
	stmt, args, err := query.Where(
		sq.Eq{
			SMSConfigColumnID.identifier():         id,
			SMSConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, stmt, args...)
	return config, err
}

func (q *Queries) SMSProviderConfig(ctx context.Context, queries ...SearchQuery) (config *SMSConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMSConfigQuery(ctx, q.client)
	for _, searchQuery := range queries {
		query = searchQuery.toQuery(query)
	}
	stmt, args, err := query.Where(
		sq.Eq{
			SMSConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, stmt, args...)
	return config, err
}

func (q *Queries) SearchSMSConfigs(ctx context.Context, queries *SMSConfigsSearchQueries) (configs *SMSConfigs, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMSConfigsQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			SMSConfigColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-sn9Jf", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		configs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-aJnZL", "Errors.Internal")
	}
	configs.State, err = q.latestState(ctx, smsConfigsTable)
	return configs, err
}

func NewSMSProviderStateQuery(state domain.SMSConfigState) (SearchQuery, error) {
	return NewNumberQuery(SMSConfigColumnState, state, NumberEquals)
}

func prepareSMSConfigQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*SMSConfig, error)) {
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
			SMSTwilioConfigColumnVerifyServiceSID.identifier(),
		).From(smsConfigsTable.identifier()).
			LeftJoin(join(SMSTwilioConfigColumnSMSID, SMSConfigColumnID) + db.Timetravel(call.Took(ctx))).
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
				&twilioConfig.verifyServiceSid,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-fn99w", "Errors.SMSConfig.NotExisting")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-3n9Js", "Errors.Internal")
			}

			twilioConfig.set(config)

			return config, nil
		}
}

func prepareSMSConfigsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*SMSConfigs, error)) {
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
			SMSTwilioConfigColumnVerifyServiceSID.identifier(),
			countColumn.identifier(),
		).From(smsConfigsTable.identifier()).
			LeftJoin(join(SMSTwilioConfigColumnSMSID, SMSConfigColumnID) + db.Timetravel(call.Took(ctx))).
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
					&twilioConfig.verifyServiceSid,
					&configs.Count,
				)

				if err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-d9jJd", "Errors.Internal")
				}

				twilioConfig.set(config)

				configs.Configs = append(configs.Configs, config)
			}

			return configs, nil
		}
}

type sqlTwilioConfig struct {
	smsID            sql.NullString
	sid              sql.NullString
	token            *crypto.CryptoValue
	senderNumber     sql.NullString
	verifyServiceSid sql.NullString
}

func (c sqlTwilioConfig) set(smsConfig *SMSConfig) {
	if !c.smsID.Valid {
		return
	}
	smsConfig.TwilioConfig = &Twilio{
		SID:              c.sid.String,
		Token:            c.token,
		SenderNumber:     c.senderNumber.String,
		VerifyServiceSID: c.verifyServiceSid.String,
	}
}
