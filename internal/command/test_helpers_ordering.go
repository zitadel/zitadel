package command

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"time"
)

// generateTestKeyAndCert generates a real valid RSA key and self-signed cert for testing
func generateTestKeyAndCert() ([]byte, []byte, error) {
	// Generate RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// PEM encode the private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Create a simple self-signed certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	// Sign the certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// PEM encode the certificate
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return privateKeyPEM, certPEM, nil
}

// noOpEncryption implements generic encryption interface for testing
type noOpEncryption struct{}

func (n *noOpEncryption) Algorithm() string                                   { return "no-op" }
func (n *noOpEncryption) EncryptionKeyID() string                             { return "no-op-key" }
func (n *noOpEncryption) DecryptionKeyIDs() []string                          { return []string{"no-op-key"} }
func (n *noOpEncryption) Encrypt(value []byte) ([]byte, error)                { return value, nil }
func (n *noOpEncryption) Decrypt(hashed []byte, keyID string) ([]byte, error) { return hashed, nil }
func (n *noOpEncryption) DecryptString(hashed []byte, keyID string) (string, error) {
	return string(hashed), nil
}
