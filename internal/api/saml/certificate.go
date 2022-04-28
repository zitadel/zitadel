package saml

import (
	"context"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/keypair"
	"github.com/zitadel/saml/pkg/provider/key"
	"gopkg.in/square/go-jose.v2"
	"time"
)

func (p *Storage) resetTimer(timer *time.Timer, shortRefresh bool) (nextCheck time.Duration) {
	nextCheck = p.certificateRotationCheck
	defer func() { timer.Reset(nextCheck) }()
	if shortRefresh || p.currentCACertificate == nil || p.currentResponseCertificate == nil || p.currentMetadataCertificate == nil {
		return nextCheck
	}
	maxLifetime := time.Until(p.currentCACertificate.Expiry())
	if maxLifetime < p.certificateGracefulPeriod+2*p.certificateRotationCheck {
		return nextCheck
	}
	return maxLifetime - p.certificateGracefulPeriod - p.certificateRotationCheck
}

func (p *Storage) GetCertificateAndKey(ctx context.Context, certAndKeyCh chan<- key.CertificateAndKey, usage model.KeyUsage) {
	renewTimer := time.NewTimer(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-p.certChan:
				if !renewTimer.Stop() {
					<-renewTimer.C
				}
				checkAfter := p.resetTimer(renewTimer, true)
				logging.Log("SAML-dK432").Infof("requested next signing key check in %s", checkAfter)
			case <-renewTimer.C:
				p.getCertificateAndKey(ctx, renewTimer, certAndKeyCh, usage)
			}
		}
	}()
}

func (p *Storage) getCertificateAndKey(ctx context.Context, renewTimer *time.Timer, certAndKeyCh chan<- key.CertificateAndKey, usage model.KeyUsage) {
	certs, err := p.query.ActiveCertificates(ctx, time.Now().Add(p.certificateGracefulPeriod), usage)
	if err != nil {
		checkAfter := p.resetTimer(renewTimer, true)
		logging.Log("SAML-ASff").Infof("next signing key check in %s", checkAfter)
		return
	}

	if len(certs.Certificates) == 0 {
		var sequence uint64
		if certs.LatestSequence != nil {
			sequence = certs.LatestSequence.Sequence
		}
		p.refreshCertificate(ctx, certAndKeyCh, usage, sequence)
		checkAfter := p.resetTimer(renewTimer, true)
		logging.Log("SAML-ASDf3").Infof("next signing key check in %s", checkAfter)
		return
	}
	err = p.exchangeCertificate(selectCertificate(certs.Certificates), certAndKeyCh, usage)
	logging.Log("SAML-aDfg3").OnError(err).Error("could not exchange signing key")
	checkAfter := p.resetTimer(renewTimer, err != nil)
	logging.Log("SAML-dK432").Infof("next signing key check in %s", checkAfter)
}

func (p *Storage) refreshCertificate(
	ctx context.Context,
	certAndKeyCh chan<- key.CertificateAndKey,
	usage model.KeyUsage,
	sequence uint64,
) {
	current := p.getCurrent(usage)
	currentCert := *current

	if currentCert != nil && currentCert.Expiry().Before(time.Now().UTC()) {
		logging.Log("SAML-ADg26").Info("unset current signing key")
		certAndKeyCh <- key.CertificateAndKey{}
	}
	ok, err := p.ensureIsLatestCertificate(ctx, sequence)
	if err != nil {
		logging.Log("SAML-sdz53").WithError(err).Error("could not ensure latest key")
		return
	}
	if !ok {
		logging.Log("EVENT-GBD13").Warn("view not up to date, retrying later")
		return
	}
	err = p.lockAndGenerateCertificateAndKey(ctx, usage, sequence)
	logging.Log("EVENT-B3d21").OnError(err).Warn("could not create signing key")
}

func (p *Storage) ensureIsLatestCertificate(ctx context.Context, sequence uint64) (bool, error) {
	maxSequence, err := p.getMaxKeySequence(ctx)
	if err != nil {
		return false, fmt.Errorf("error retrieving new events: %w", err)
	}
	return sequence == maxSequence, nil
}

func (p *Storage) lockAndGenerateCertificateAndKey(ctx context.Context, usage model.KeyUsage, sequence uint64) error {
	errs := p.locker.Lock(ctx, p.certificateRotationCheck*2)
	err, ok := <-errs
	if err != nil || !ok {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		logging.Log("OIDC-Dfg32").OnError(err).Warn("initial lock failed")
		return err
	}

	switch usage {
	case model.KeyUsageSAMLMetadataSigning, model.KeyUsageSAMLResponseSinging:

		certs, err := p.query.ActiveCertificates(ctx, time.Now().Add(p.certificateGracefulPeriod), model.KeyUsageSAMLCA)
		if err != nil {
			return fmt.Errorf("error while reading ca certificate: %w", err)
		}
		certAndKey := selectCertificate(certs.Certificates)

		keyData, err := crypto.Decrypt(certAndKey.Key(), p.encAlg)
		if err != nil {
			return err
		}
		privateKey, err := crypto.BytesToPrivateKey(keyData)
		if err != nil {
			return err
		}

		certData, err := crypto.Decrypt(certAndKey.Certificate(), p.encAlg)
		if err != nil {
			return err
		}
		cert, err := crypto.BytesToCertificate(certData)
		if err != nil {
			return err
		}

		switch usage {
		case model.KeyUsageSAMLMetadataSigning:
			return p.command.GenerateSAMLMetadataCertificate(ctx, p.certificateAlgorithm, privateKey, cert)
		case model.KeyUsageSAMLResponseSinging:
			return p.command.GenerateSAMLResponseCertificate(ctx, p.certificateAlgorithm, privateKey, cert)
		default:
			return fmt.Errorf("unknown usage")
		}
	case model.KeyUsageSAMLCA:
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

func (p *Storage) getCurrent(usage model.KeyUsage) *query.Certificate {
	switch usage {
	case model.KeyUsageSAMLResponseSinging:
		return &p.currentResponseCertificate
	case model.KeyUsageSAMLMetadataSigning:
		return &p.currentMetadataCertificate
	case model.KeyUsageSAMLCA:
		return &p.currentCACertificate
	}

	return nil
}

func (p *Storage) setCurrent(usage model.KeyUsage, current *query.Certificate) {
	switch usage {
	case model.KeyUsageSAMLResponseSinging:
		p.currentResponseCertificate = *current
	case model.KeyUsageSAMLMetadataSigning:
		p.currentMetadataCertificate = *current
	case model.KeyUsageSAMLCA:
		p.currentCACertificate = *current
	}
}

func (p *Storage) exchangeCertificate(certificate query.Certificate, certAndKeyCh chan<- key.CertificateAndKey, usage model.KeyUsage) (err error) {
	current := p.getCurrent(usage)
	currentCert := *current

	if currentCert != nil && currentCert.ID() == certificate.ID() {
		logging.Log("OIDC-Abb3e").Info("no new signing key")
		return nil
	}

	keyData, err := crypto.Decrypt(certificate.Key(), p.encAlg)
	if err != nil {
		return err
	}
	privateKey, err := crypto.BytesToPrivateKey(keyData)
	if err != nil {
		return err
	}

	certData, err := crypto.Decrypt(certificate.Certificate(), p.encAlg)
	if err != nil {
		return err
	}
	cert, err := crypto.BytesToCertificate(certData)
	if err != nil {
		return err
	}

	certAndKeyCh <- key.CertificateAndKey{
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
	}

	p.setCurrent(usage, &certificate)
	logging.LogWithFields("OIDC-dsg54", "keyID", certificate.ID()).Info("exchanged signing key")
	return nil
}

func selectCertificate(certs []query.Certificate) query.Certificate {
	return certs[len(certs)-1]
}
