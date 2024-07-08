package command

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

func (c *Commands) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	privateCrypto, publicCrypto, err := crypto.GenerateEncryptedKeyPair(c.keySize, c.keyAlgorithm)
	if err != nil {
		return err
	}
	keyID, err := id_generator.Next()
	if err != nil {
		return err
	}

	privateKeyExp := time.Now().UTC().Add(c.privateKeyLifetime)
	publicKeyExp := time.Now().UTC().Add(c.publicKeyLifetime)

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSigning,
		algorithm,
		privateCrypto, publicCrypto,
		privateKeyExp, publicKeyExp))
	return err
}

func (c *Commands) GenerateSAMLCACertificate(ctx context.Context, algorithm string) error {
	now := time.Now().UTC()
	after := now.Add(c.certificateLifetime)
	randInt, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return err
	}

	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCACertificate(c.certKeySize, c.keyAlgorithm, c.certificateAlgorithm, &crypto.CertificateInformations{
		SerialNumber: randInt,
		Organisation: []string{"ZITADEL"},
		CommonName:   "ZITADEL SAML CA",
		NotBefore:    now,
		NotAfter:     after,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
	})
	if err != nil {
		return err
	}
	keyID, err := id_generator.Next()
	if err != nil {
		return err
	}

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx,
		keypair.NewAddedEvent(
			ctx,
			keyAgg,
			domain.KeyUsageSAMLCA,
			algorithm,
			privateCrypto, publicCrypto,
			after, after,
		),
		keypair.NewAddedCertificateEvent(
			ctx,
			keyAgg,
			certificateCrypto,
			after,
		),
	)
	return err
}

func (c *Commands) GenerateSAMLResponseCertificate(ctx context.Context, algorithm string, caPrivateKey *rsa.PrivateKey, caCertificate []byte) error {
	now := time.Now().UTC()
	after := now.Add(c.certificateLifetime)
	randInt, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return err
	}

	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCertificate(c.certKeySize, c.keyAlgorithm, c.certificateAlgorithm, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: randInt,
		Organisation: []string{"ZITADEL"},
		CommonName:   "ZITADEL SAML response",
		NotBefore:    now,
		NotAfter:     after,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	})
	if err != nil {
		return err
	}
	keyID, err := id_generator.Next()
	if err != nil {
		return err
	}

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx,
		keypair.NewAddedEvent(
			ctx,
			keyAgg,
			domain.KeyUsageSAMLResponseSinging,
			algorithm,
			privateCrypto, publicCrypto,
			after, after,
		),
		keypair.NewAddedCertificateEvent(
			ctx,
			keyAgg,
			certificateCrypto,
			after,
		),
	)
	return err
}

func (c *Commands) GenerateSAMLMetadataCertificate(ctx context.Context, algorithm string, caPrivateKey *rsa.PrivateKey, caCertificate []byte) error {
	now := time.Now().UTC()
	after := now.Add(c.certificateLifetime)
	randInt, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return err
	}
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCertificate(c.certKeySize, c.keyAlgorithm, c.certificateAlgorithm, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: randInt,
		Organisation: []string{"ZITADEL"},
		CommonName:   "ZITADEL SAML metadata",
		NotBefore:    now,
		NotAfter:     after,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	})
	if err != nil {
		return err
	}
	keyID, err := id_generator.Next()
	if err != nil {
		return err
	}

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx,
		keypair.NewAddedEvent(
			ctx,
			keyAgg,
			domain.KeyUsageSAMLMetadataSigning,
			algorithm,
			privateCrypto, publicCrypto,
			after, after),
		keypair.NewAddedCertificateEvent(
			ctx,
			keyAgg,
			certificateCrypto,
			after,
		),
	)
	return err
}
