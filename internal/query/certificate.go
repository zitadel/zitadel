package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Certificate interface {
	Key
	Expiry() time.Time
	Key() *crypto.CryptoValue
	Certificate() []byte
}

type Certificates struct {
	SearchResponse
	Certificates []Certificate
}

type rsaCertificate struct {
	key
	expiry      time.Time
	privateKey  *crypto.CryptoValue
	certificate []byte
}

func (c *rsaCertificate) Expiry() time.Time {
	return c.expiry
}

func (c *rsaCertificate) Key() *crypto.CryptoValue {
	return c.privateKey
}

func (c *rsaCertificate) Certificate() []byte {
	return c.certificate
}

var (
	certificateTable = table{
		name:          projection.CertificateTable,
		instanceIDCol: projection.CertificateColumnInstanceID,
	}
	CertificateColID = Column{
		name:  projection.CertificateColumnID,
		table: certificateTable,
	}
	CertificateColExpiry = Column{
		name:  projection.CertificateColumnExpiry,
		table: certificateTable,
	}
	CertificateColCertificate = Column{
		name:  projection.CertificateColumnCertificate,
		table: certificateTable,
	}
)

func (q *Queries) ActiveCertificates(ctx context.Context, t time.Time, usage crypto.KeyUsage) (certs *Certificates, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareCertificateQuery()
	if t.IsZero() {
		t = time.Now()
	}
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{
				KeyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
				KeyColUse.identifier():        usage,
			},
			sq.Gt{CertificateColExpiry.identifier(): t},
			sq.Gt{KeyPrivateColExpiry.identifier(): t},
		},
	).OrderBy(KeyPrivateColExpiry.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SDfkg", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		certs, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Sgan4", "Errors.Internal")
	}

	certs.State, err = q.latestState(ctx, keyTable)
	if !zerrors.IsNotFound(err) {
		return certs, err
	}
	return certs, nil
}

func prepareCertificateQuery() (sq.SelectBuilder, func(*sql.Rows) (*Certificates, error)) {
	return sq.Select(
			KeyColID.identifier(),
			KeyColCreationDate.identifier(),
			KeyColChangeDate.identifier(),
			KeyColSequence.identifier(),
			KeyColResourceOwner.identifier(),
			KeyColAlgorithm.identifier(),
			KeyColUse.identifier(),
			CertificateColExpiry.identifier(),
			CertificateColCertificate.identifier(),
			KeyPrivateColKey.identifier(),
			countColumn.identifier(),
		).From(keyTable.identifier()).
			LeftJoin(join(CertificateColID, KeyColID)).
			LeftJoin(join(KeyPrivateColID, KeyColID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Certificates, error) {
			certificates := make([]Certificate, 0)
			var count uint64
			for rows.Next() {
				k := new(rsaCertificate)
				err := rows.Scan(
					&k.id,
					&k.creationDate,
					&k.changeDate,
					&k.sequence,
					&k.resourceOwner,
					&k.algorithm,
					&k.use,
					&k.expiry,
					&k.certificate,
					&k.privateKey,
					&count,
				)
				if err != nil {
					return nil, err
				}

				certificates = append(certificates, k)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-rKd6k", "Errors.Query.CloseRows")
			}

			return &Certificates{
				Certificates: certificates,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
