package command

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"math/rand"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	keypair "github.com/caos/zitadel/internal/repository/keypair"
)

const (
	oidcUser = "OIDC"
	samlUser = "SAML"
)

func (c *Commands) GenerateSAMLCACertificate(ctx context.Context) error {
	ctx = setSAMLCtx(ctx)
	now := time.Now().UTC()
	after := now.Add(c.privateKeyLifetime)
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCACertificate(c.certKeySize, c.keyAlgorithm, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"caos AG"},
		CommonName:   "Zitadel SAML CA",
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

	keyPairWriteModel := NewKeyPairWriteModel(keyID, domain.IAMID)
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSAMLCA,
		"", //TODO do we need this information for SAML?
		privateCrypto, publicCrypto, certificateCrypto,
		after, after, after))
	return err
}

func (c *Commands) GenerateSAMLResponseCertificate(ctx context.Context, reason string, algorithm string, caPrivateKey *rsa.PrivateKey, caCertificate []byte) error {
	ctx = setSAMLCtx(ctx)
	now := time.Now().UTC()
	after := now.Add(c.privateKeyLifetime)
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCertificate(c.certKeySize, c.keyAlgorithm, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"caos AG"},
		CommonName:   reason,
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

	keyPairWriteModel := NewKeyPairWriteModel(keyID, domain.IAMID)
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSAMLResponseSinging,
		algorithm,
		privateCrypto, publicCrypto, certificateCrypto,
		after, after, after))
	return err
}

func (c *Commands) GenerateSAMLMetadataCertificate(ctx context.Context, reason string, algorithm string, caPrivateKey *rsa.PrivateKey, caCertificate []byte) error {
	ctx = setSAMLCtx(ctx)
	now := time.Now().UTC()
	after := now.Add(c.privateKeyLifetime)
	privateCrypto, publicCrypto, certificateCrypto, err := crypto.GenerateEncryptedKeyPairWithCertificate(c.certKeySize, c.keyAlgorithm, caPrivateKey, caCertificate, &crypto.CertificateInformations{
		SerialNumber: big.NewInt(int64(rand.Intn(50000))),
		Organisation: []string{"caos AG"},
		CommonName:   reason,
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

	keyPairWriteModel := NewKeyPairWriteModel(keyID, domain.IAMID)
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSAMLMetadataSigning,
		algorithm,
		privateCrypto, publicCrypto, certificateCrypto,
		after, after, after))
	return err
}

func (c *Commands) GenerateOIDCSigningPair(ctx context.Context, algorithm string) error {
	ctx = setOIDCCtx(ctx)
	privateCrypto, publicCrypto, err := crypto.GenerateEncryptedKeyPair(c.keySize, c.keyAlgorithm)
	if err != nil {
		return err
	}
	keyID, err := c.idGenerator.Next()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	privateKeyExp := now.Add(c.privateKeyLifetime)
	publicKeyExp := now.Add(c.publicKeyLifetime)
	certificateExp := now.Add(c.publicKeyLifetime)

	keyPairWriteModel := NewKeyPairWriteModel(keyID, domain.IAMID)
	keyAgg := KeyPairAggregateFromWriteModel(&keyPairWriteModel.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, keypair.NewAddedEvent(
		ctx,
		keyAgg,
		domain.KeyUsageSigning,
		algorithm,
		privateCrypto, publicCrypto, nil,
		privateKeyExp, publicKeyExp, certificateExp))
	return err
}

func setOIDCCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: oidcUser, OrgID: domain.IAMID})
}

func setSAMLCtx(ctx context.Context) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: samlUser, OrgID: domain.IAMID})
}
