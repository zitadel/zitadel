package saml

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/zitadel/logging"
	"github.com/zitadel/saml/pkg/provider/key"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"gopkg.in/square/go-jose.v2"
	"time"
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
			return errors.ThrowInternal(err, "SAML-8u01nks", "no certificate found")
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

	var sequence uint64
	if certs.LatestSequence != nil {
		sequence = certs.LatestSequence.Sequence
	}

	return nil, p.refreshCertificate(ctx, usage, sequence)
}

func (p *Storage) refreshCertificate(
	ctx context.Context,
	usage domain.KeyUsage,
	sequence uint64,
) error {
	current := p.getCurrent(usage)
	currentCert := *current

	if currentCert != nil && currentCert.Expiry().Before(time.Now().UTC()) {
		logging.Log("SAML-ADg26").Info("unset current signing key")
		return fmt.Errorf("unset current signing key")
	}
	ok, err := p.ensureIsLatestCertificate(ctx, sequence)
	if err != nil {
		logging.Log("SAML-sdz53").WithError(err).Error("could not ensure latest key")
		return err
	}
	if !ok {
		logging.Log("EVENT-GBD13").Warn("view not up to date, retrying later")
		return err
	}
	err = p.lockAndGenerateCertificateAndKey(ctx, usage, sequence)
	logging.Log("EVENT-B3d21").OnError(err).Warn("could not create signing key")
	return nil
}

func (p *Storage) ensureIsLatestCertificate(ctx context.Context, sequence uint64) (bool, error) {
	maxSequence, err := p.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return sequence == maxSequence, nil
}

func (p *Storage) lockAndGenerateCertificateAndKey(ctx context.Context, usage domain.KeyUsage, sequence uint64) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = setSAMLCtx(ctx)

	errs := p.locker.Lock(ctx, lockDuration, authz.GetInstance(ctx).InstanceID())
	err, ok := <-errs
	if err != nil || !ok {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.OnError(err).Warn("initial lock failed")
		return err
	}

	switch usage {
	case domain.KeyUsageSAMLMetadataSigning, domain.KeyUsageSAMLResponseSinging:
		certAndKey, err := p.GetCertificateAndKey(ctx, domain.KeyUsageSAMLCA)
		if err != nil {
			return fmt.Errorf("error while reading ca certificate: %w", err)
		}
		if certAndKey.Key.Key == nil || certAndKey.Certificate.Key == nil {
			return fmt.Errorf("has no ca certificate")
		}
		certWebKey := certAndKey.Certificate.Key.(jose.JSONWebKey)
		keyWebKey := certAndKey.Key.Key.(jose.JSONWebKey)

		switch usage {
		case domain.KeyUsageSAMLMetadataSigning:
			return p.command.GenerateSAMLMetadataCertificate(ctx, p.certificateAlgorithm, keyWebKey.Key.(*rsa.PrivateKey), certWebKey.Key.([]byte))
		case domain.KeyUsageSAMLResponseSinging:
			return p.command.GenerateSAMLResponseCertificate(ctx, p.certificateAlgorithm, keyWebKey.Key.(*rsa.PrivateKey), certWebKey.Key.([]byte))
		default:
			return fmt.Errorf("unknown usage")
		}
	case domain.KeyUsageSAMLCA:
		return p.command.GenerateSAMLCACertificate(ctx)
	default:
		return fmt.Errorf("unknown certificate usage")
	}
}

func (p *Storage) getMaxKeySequence(ctx context.Context) (uint64, error) {
	return p.eventstore.LatestSequence(ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsMaxSequence).
			ResourceOwner(domain.IAMID).
			AddQuery().
			AggregateTypes(keypair.AggregateType).
			Builder(),
	)
}

func (p *Storage) getCurrent(usage domain.KeyUsage) *query.Certificate {
	switch usage {
	case domain.KeyUsageSAMLResponseSinging:
		return &p.currentResponseCertificate
	case domain.KeyUsageSAMLMetadataSigning:
		return &p.currentMetadataCertificate
	case domain.KeyUsageSAMLCA:
		return &p.currentCACertificate
	}

	return nil
}

func (p *Storage) setCurrent(usage domain.KeyUsage, current *query.Certificate) {
	switch usage {
	case domain.KeyUsageSAMLResponseSinging:
		p.currentResponseCertificate = *current
	case domain.KeyUsageSAMLMetadataSigning:
		p.currentMetadataCertificate = *current
	case domain.KeyUsageSAMLCA:
		p.currentCACertificate = *current
	}
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

	certData, err := crypto.Decrypt(certificate.Certificate(), p.encAlg)
	if err != nil {
		return nil, err
	}
	cert, err := crypto.BytesToCertificate(certData)
	if err != nil {
		return nil, err
	}

	return &key.CertificateAndKey{
		Key: &jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(p.certificateAlgorithm),
			Key: jose.JSONWebKey{
				KeyID: certificate.ID(),
				Key:   privateKey,
			},
		},
		Certificate: &jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(p.certificateAlgorithm),
			Key: jose.JSONWebKey{
				KeyID: certificate.ID(),
				Key:   cert,
			},
		},
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
		time.Sleep(retryBackoff)
		err = retryable()
		if err == nil {
			return nil
		}
	}
	return err
}
