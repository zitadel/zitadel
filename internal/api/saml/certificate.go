package saml

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/zitadel/logging"
	"github.com/zitadel/saml/pkg/provider/key"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	locksTable = "projections.locks"
	signingKey = "signing_key"
	samlUser   = "SAML"

	retryBackoff   = 500 * time.Millisecond
	retryCount     = 3
	lockDuration   = retryCount * retryBackoff * 5
	gracefulPeriod = 10 * time.Minute
)

type CertificateAndKey struct {
	algorithm   jose.SignatureAlgorithm
	id          string
	key         interface{}
	certificate interface{}
}

func (c *CertificateAndKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return c.algorithm
}

func (c *CertificateAndKey) Key() interface{} {
	return c.key
}

func (c *CertificateAndKey) Certificate() interface{} {
	return c.certificate
}

func (c *CertificateAndKey) ID() string {
	return c.id
}

func (p *Storage) GetCertificateAndKey(ctx context.Context, usage domain.KeyUsage) (certAndKey *key.CertificateAndKey, err error) {
	err = retry(func() error {
		certAndKey, err = p.getCertificateAndKey(ctx, usage)
		if err != nil {
			return err
		}
		if certAndKey == nil {
			return zerrors.ThrowInternal(err, "SAML-8u01nks", "no certificate found")
		}
		return nil
	})
	return certAndKey, err
}

func (p *Storage) getCertificateAndKey(ctx context.Context, usage domain.KeyUsage) (*key.CertificateAndKey, error) {
	certs, err := p.query.ActiveCertificates(ctx, time.Now().Add(gracefulPeriod), usage)
	if err != nil {
		return nil, err
	}

	if len(certs.Certificates) > 0 {
		return p.certificateToCertificateAndKey(selectCertificate(certs.Certificates))
	}

	var position float64
	if certs.State != nil {
		position = certs.State.Position
	}

	return nil, p.refreshCertificate(ctx, usage, position)
}

func (p *Storage) refreshCertificate(
	ctx context.Context,
	usage domain.KeyUsage,
	position float64,
) error {
	ok, err := p.ensureIsLatestCertificate(ctx, position)
	if err != nil {
		logging.WithError(err).Error("could not ensure latest key")
		return err
	}
	if !ok {
		logging.Warn("view not up to date, retrying later")
		return err
	}
	err = p.lockAndGenerateCertificateAndKey(ctx, usage)
	logging.OnError(err).Warn("could not create signing key")
	return nil
}

func (p *Storage) ensureIsLatestCertificate(ctx context.Context, position float64) (bool, error) {
	maxSequence, err := p.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return position >= maxSequence, nil
}

func (p *Storage) lockAndGenerateCertificateAndKey(ctx context.Context, usage domain.KeyUsage) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = setSAMLCtx(ctx)

	errs := p.locker.Lock(ctx, lockDuration, authz.GetInstance(ctx).InstanceID())
	err, ok := <-errs
	if err != nil || !ok {
		if zerrors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.OnError(err).Debug("initial lock failed")
		return err
	}

	switch usage {
	case domain.KeyUsageSAMLMetadataSigning, domain.KeyUsageSAMLResponseSinging:
		certAndKey, err := p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLCA)
		if err != nil {
			return fmt.Errorf("error while reading ca certificate: %w", err)
		}
		if certAndKey.Key == nil || certAndKey.Certificate == nil {
			return fmt.Errorf("has no ca certificate")
		}

		switch usage {
		case domain.KeyUsageSAMLMetadataSigning:
			return p.command.GenerateSAMLMetadataCertificate(setSAMLCtx(ctx), p.certificateAlgorithm, certAndKey.Key, certAndKey.Certificate)
		case domain.KeyUsageSAMLResponseSinging:
			return p.command.GenerateSAMLResponseCertificate(setSAMLCtx(ctx), p.certificateAlgorithm, certAndKey.Key, certAndKey.Certificate)
		default:
			return fmt.Errorf("unknown usage")
		}
	case domain.KeyUsageSAMLCA:
		return p.command.GenerateSAMLCACertificate(setSAMLCtx(ctx), p.certificateAlgorithm)
	default:
		return fmt.Errorf("unknown certificate usage")
	}
}

func (p *Storage) getMaxKeySequence(ctx context.Context) (float64, error) {
	return p.eventstore.LatestSequence(ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxSequence).
			ResourceOwner(authz.GetInstance(ctx).InstanceID()).
			AwaitOpenTransactions().
			AddQuery().
			AggregateTypes(keypair.AggregateType).
			EventTypes(
				keypair.AddedEventType,
				keypair.AddedCertificateEventType,
			).
			Or().
			AggregateTypes(instance.AggregateType).
			EventTypes(instance.InstanceRemovedEventType).
			Builder(),
	)
}

func (p *Storage) certificateToCertificateAndKey(certificate query.Certificate) (_ *key.CertificateAndKey, err error) {
	keyData, err := crypto.Decrypt(certificate.Key(), p.encAlg)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.BytesToPrivateKey(keyData)
	if err != nil {
		return nil, err
	}

	cert, err := crypto.BytesToCertificate(certificate.Certificate())
	if err != nil {
		return nil, err
	}

	return &key.CertificateAndKey{
		Key:         privateKey,
		Certificate: cert,
	}, nil
}

func selectCertificate(certs []query.Certificate) query.Certificate {
	return certs[len(certs)-1]
}

func setSAMLCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: samlUser, OrgID: authz.GetInstance(ctx).InstanceID()})
}

func retry(retryable func() error) (err error) {
	for i := 0; i < retryCount; i++ {
		err = retryable()
		if err == nil {
			return nil
		}
		time.Sleep(retryBackoff)
	}
	return err
}
