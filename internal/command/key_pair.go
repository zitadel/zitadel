package command

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"math/rand"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

func (c *Commands) GenerateSigningKeyPair(ctx context.Context, algorithm string) error {
	privateCrypto, publicCrypto, err := crypto.GenerateEncryptedKeyPair(c.keySize, c.keyAlgorithm)
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	privateKeyExp := time.Now().UTC().Add(c.privateKeyLifetime)
	publicKeyExp := time.Now().UTC().Add(c.publicKeyLifetime)
	certificateExp := time.Now().UTC().Add(c.publicKeyLifetime)

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSigning,
		algorithm,
		privateCrypto, publicCrypto, nil,
		privateKeyExp, publicKeyExp, certificateExp))
	return err
}

func (c *Commands) GenerateSAMLCACertificate(ctx context.Context) error {
	now := time.Now().UTC()
	after := now.Add(c.certificateLifetime)
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCACertificate(c.certKeySize, c.certificateAlgorithm, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"ZITADEL"},
		CommonName:   "ZITADEL SAML CA",
		NotBefore:    &now,
		NotAfter:     &after,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
	})
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSAMLCA,
		"", //TODO do we need this information for SAML?
		privateCrypto, publicCrypto, certificateCrypto,
		after, after, after))
	return err
}

func (c *Commands) GenerateSAMLResponseCertificate(ctx context.Context, algorithm string, caPrivateKey *rsa.PrivateKey, caCertificate []byte) error {
	now := time.Now().UTC()
	after := now.Add(c.certificateLifetime)
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCertificate(c.certKeySize, c.certificateAlgorithm, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"ZITADEL"},
		CommonName:   "ZITADEL SAML response",
		NotBefore:    &now,
		NotAfter:     &after,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	})
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSAMLResponseSinging,
		algorithm,
		privateCrypto, publicCrypto, certificateCrypto,
		after, after, after))
	return err
}

func (c *Commands) GenerateSAMLMetadataCertificate(ctx context.Context, algorithm string, caPrivateKey *rsa.PrivateKey, caCertificate []byte) error {
	now := time.Now().UTC()
	after := now.Add(c.certificateLifetime)
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCertificate(c.certKeySize, c.certificateAlgorithm, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"ZITADEL"},
		CommonName:   "ZITADEL SAML metadata",
		NotBefore:    &now,
		NotAfter:     &after,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	})
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	keyPairWriteModel := NewKeyPairWriteModel(keyID, authz.GetInstance(ctx).InstanceID())
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.Push(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSAMLMetadataSigning,
		algorithm,
		privateCrypto, publicCrypto, certificateCrypto,
		after, after, after))
	return err
}
