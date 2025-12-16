package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
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
	Description   string

	TwilioConfig *Twilio
	HTTPConfig   *HTTP
}

func (h *SMSConfig) decryptSigningKey(alg crypto.EncryptionAlgorithm) error {
	if h.HTTPConfig == nil {
		return nil
	}
	return h.HTTPConfig.decryptSigningKey(alg)
}

type Twilio struct {
	SID              string
	Token            *crypto.CryptoValue
	SenderNumber     string
	VerifyServiceSID string
}

type HTTP struct {
	Endpoint   string
	signingKey *crypto.CryptoValue
	SigningKey string
}

func (h *HTTP) decryptSigningKey(alg crypto.EncryptionAlgorithm) error {
	if h.signingKey == nil {
		return nil
	}
	keyValue, err := crypto.DecryptString(h.signingKey, alg)
	if err != nil {
		return zerrors.ThrowInternal(err, "QUERY-bxovy3YXwy", "Errors.Internal")
	}
	h.SigningKey = keyValue
	return nil
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
	SMSColumnID = Column{
		name:  projection.SMSColumnID,
		table: smsConfigsTable,
	}
	SMSColumnAggregateID = Column{
		name:  projection.SMSColumnAggregateID,
		table: smsConfigsTable,
	}
	SMSColumnCreationDate = Column{
		name:  projection.SMSColumnCreationDate,
		table: smsConfigsTable,
	}
	SMSColumnChangeDate = Column{
		name:  projection.SMSColumnChangeDate,
		table: smsConfigsTable,
	}
	SMSColumnResourceOwner = Column{
		name:  projection.SMSColumnResourceOwner,
		table: smsConfigsTable,
	}
	SMSColumnInstanceID = Column{
		name:  projection.SMSColumnInstanceID,
		table: smsConfigsTable,
	}
	SMSColumnState = Column{
		name:  projection.SMSColumnState,
		table: smsConfigsTable,
	}
	SMSColumnSequence = Column{
		name:  projection.SMSColumnSequence,
		table: smsConfigsTable,
	}
	SMSColumnDescription = Column{
		name:  projection.SMSColumnDescription,
		table: smsConfigsTable,
	}
)

var (
	smsTwilioTable = table{
		name:          projection.SMSTwilioTable,
		instanceIDCol: projection.SMSTwilioColumnInstanceID,
	}
	SMSTwilioColumnSMSID = Column{
		name:  projection.SMSTwilioColumnSMSID,
		table: smsTwilioTable,
	}
	SMSTwilioColumnSID = Column{
		name:  projection.SMSTwilioColumnSID,
		table: smsTwilioTable,
	}
	SMSTwilioColumnToken = Column{
		name:  projection.SMSTwilioColumnToken,
		table: smsTwilioTable,
	}
	SMSTwilioColumnSenderNumber = Column{
		name:  projection.SMSTwilioColumnSenderNumber,
		table: smsTwilioTable,
	}
	SMSTwilioColumnVerifyServiceSID = Column{
		name:  projection.SMSTwilioColumnVerifyServiceSID,
		table: smsTwilioTable,
	}
)

var (
	smsHTTPTable = table{
		name:          projection.SMSHTTPTable,
		instanceIDCol: projection.SMSHTTPColumnInstanceID,
	}
	SMSHTTPColumnSMSID = Column{
		name:  projection.SMSHTTPColumnSMSID,
		table: smsHTTPTable,
	}
	SMSHTTPColumnEndpoint = Column{
		name:  projection.SMSHTTPColumnEndpoint,
		table: smsHTTPTable,
	}
	SMSHTTPColumnSigningKey = Column{
		name:  projection.SMSHTTPColumnSigningKey,
		table: smsHTTPTable,
	}
)

func (q *Queries) SMSProviderConfigByID(ctx context.Context, id string) (config *SMSConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMSConfigQuery()
	stmt, args, err := query.Where(
		sq.Eq{
			SMSColumnID.identifier():         id,
			SMSColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	if err := config.decryptSigningKey(q.smsEncryptionAlgorithm); err != nil {
		return nil, err
	}
	return config, nil
}

func (q *Queries) SMSProviderConfigActive(ctx context.Context, instanceID string) (config *SMSConfig, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMSConfigQuery()
	stmt, args, err := query.Where(
		sq.Eq{
			SMSColumnInstanceID.identifier(): instanceID,
			SMSColumnState.identifier():      domain.SMSConfigStateActive,
		},
	).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-dn9JW", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		config, err = scan(row)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	if err := config.decryptSigningKey(q.smsEncryptionAlgorithm); err != nil {
		return nil, err
	}
	return config, nil
}

func (q *Queries) SearchSMSConfigs(ctx context.Context, queries *SMSConfigsSearchQueries) (configs *SMSConfigs, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareSMSConfigsQuery()
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			SMSColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-sn9Jf", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		configs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-l4bxm", "Errors.Internal")
	}
	for i := range configs.Configs {
		if err := configs.Configs[i].decryptSigningKey(q.smsEncryptionAlgorithm); err != nil {
			return nil, err
		}
	}
	configs.State, err = q.latestState(ctx, smsConfigsTable)
	return configs, nil
}

func NewSMSProviderStateQuery(state domain.SMSConfigState) (SearchQuery, error) {
	return NewNumberQuery(SMSColumnState, state, NumberEquals)
}

func prepareSMSConfigQuery() (sq.SelectBuilder, func(*sql.Row) (*SMSConfig, error)) {
	return sq.Select(
			SMSColumnID.identifier(),
			SMSColumnAggregateID.identifier(),
			SMSColumnCreationDate.identifier(),
			SMSColumnChangeDate.identifier(),
			SMSColumnResourceOwner.identifier(),
			SMSColumnState.identifier(),
			SMSColumnSequence.identifier(),
			SMSColumnDescription.identifier(),

			SMSTwilioColumnSMSID.identifier(),
			SMSTwilioColumnSID.identifier(),
			SMSTwilioColumnToken.identifier(),
			SMSTwilioColumnSenderNumber.identifier(),
			SMSTwilioColumnVerifyServiceSID.identifier(),

			SMSHTTPColumnSMSID.identifier(),
			SMSHTTPColumnEndpoint.identifier(),
			SMSHTTPColumnSigningKey.identifier(),
		).From(smsConfigsTable.identifier()).
			LeftJoin(join(SMSTwilioColumnSMSID, SMSColumnID)).
			LeftJoin(join(SMSHTTPColumnSMSID, SMSColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*SMSConfig, error) {
			config := new(SMSConfig)

			var (
				twilioConfig = sqlTwilioConfig{}
				httpConfig   = sqlHTTPConfig{}
			)

			err := row.Scan(
				&config.ID,
				&config.AggregateID,
				&config.CreationDate,
				&config.ChangeDate,
				&config.ResourceOwner,
				&config.State,
				&config.Sequence,
				&config.Description,

				&twilioConfig.smsID,
				&twilioConfig.sid,
				&twilioConfig.token,
				&twilioConfig.senderNumber,
				&twilioConfig.verifyServiceSid,

				&httpConfig.id,
				&httpConfig.endpoint,
				&httpConfig.signingKey,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-fn99w", "Errors.SMSConfig.NotExisting")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-3n9Js", "Errors.Internal")
			}

			twilioConfig.set(config)
			httpConfig.setSMS(config)

			return config, nil
		}
}

func prepareSMSConfigsQuery() (sq.SelectBuilder, func(*sql.Rows) (*SMSConfigs, error)) {
	return sq.Select(
			SMSColumnID.identifier(),
			SMSColumnAggregateID.identifier(),
			SMSColumnCreationDate.identifier(),
			SMSColumnChangeDate.identifier(),
			SMSColumnResourceOwner.identifier(),
			SMSColumnState.identifier(),
			SMSColumnSequence.identifier(),
			SMSColumnDescription.identifier(),

			SMSTwilioColumnSMSID.identifier(),
			SMSTwilioColumnSID.identifier(),
			SMSTwilioColumnToken.identifier(),
			SMSTwilioColumnSenderNumber.identifier(),
			SMSTwilioColumnVerifyServiceSID.identifier(),

			SMSHTTPColumnSMSID.identifier(),
			SMSHTTPColumnEndpoint.identifier(),
			SMSHTTPColumnSigningKey.identifier(),

			countColumn.identifier(),
		).From(smsConfigsTable.identifier()).
			LeftJoin(join(SMSTwilioColumnSMSID, SMSColumnID)).
			LeftJoin(join(SMSHTTPColumnSMSID, SMSColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Rows) (*SMSConfigs, error) {
			configs := &SMSConfigs{Configs: []*SMSConfig{}}

			for row.Next() {
				config := new(SMSConfig)
				var (
					twilioConfig = sqlTwilioConfig{}
					httpConfig   = sqlHTTPConfig{}
				)

				err := row.Scan(
					&config.ID,
					&config.AggregateID,
					&config.CreationDate,
					&config.ChangeDate,
					&config.ResourceOwner,
					&config.State,
					&config.Sequence,
					&config.Description,

					&twilioConfig.smsID,
					&twilioConfig.sid,
					&twilioConfig.token,
					&twilioConfig.senderNumber,
					&twilioConfig.verifyServiceSid,

					&httpConfig.id,
					&httpConfig.endpoint,
					&httpConfig.signingKey,

					&configs.Count,
				)

				if err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-d9jJd", "Errors.Internal")
				}

				twilioConfig.set(config)
				httpConfig.setSMS(config)

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

type sqlHTTPConfig struct {
	id         sql.NullString
	endpoint   sql.NullString
	signingKey *crypto.CryptoValue
}

func (c sqlHTTPConfig) setSMS(smsConfig *SMSConfig) {
	if !c.id.Valid {
		return
	}
	smsConfig.HTTPConfig = &HTTP{
		Endpoint:   c.endpoint.String,
		signingKey: c.signingKey,
	}
}
