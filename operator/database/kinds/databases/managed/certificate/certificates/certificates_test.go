package certificates

import (
	"crypto/x509"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/pem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCertificates_CAE(t *testing.T) {
	priv, rootCa, err := NewCA()
	assert.NoError(t, err)
	assert.NotNil(t, priv)

	pemCa, err := pem.EncodeCertificate(rootCa)
	pemkey, err := pem.EncodeKey(priv)
	assert.NotNil(t, pemCa)
	assert.NotNil(t, pemkey)

	_, err = x509.ParseCertificate(rootCa)
	assert.NoError(t, err)
}

func TestCertificates_CA(t *testing.T) {
	_, rootCa, err := NewCA()
	assert.NoError(t, err)

	_, err = x509.ParseCertificate(rootCa)
	assert.NoError(t, err)
}

func TestCertificates_Chain(t *testing.T) {
	rootKey, rootCert, err := NewCA()
	assert.NoError(t, err)
	rootPem, err := pem.EncodeCertificate(rootCert)
	assert.NoError(t, err)

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPem)
	assert.Equal(t, ok, true)

	_, clientCert, err := NewClient(rootKey, rootCert, "test")

	cert, err := x509.ParseCertificate(clientCert)
	assert.NoError(t, err)

	opts := x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	_, err = cert.Verify(opts)
	assert.NoError(t, err)
}
