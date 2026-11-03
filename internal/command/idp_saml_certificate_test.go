package command

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// newTestSAMLKeyAndCertificate returns a PEM encoded RSA key (PKCS#1) and a matching self-signed certificate,
// valid from one hour ago for one day.
func newTestSAMLKeyAndCertificate(t *testing.T) (key, certificate []byte) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	now := time.Now()
	certificate = newTestCertificate(t, privateKey, &privateKey.PublicKey, now.Add(-time.Hour), now.Add(24*time.Hour))
	key = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	return key, certificate
}

func newTestCertificate(t *testing.T, signer any, publicKey any, notBefore, notAfter time.Time) []byte {
	t.Helper()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "SP"},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, signer)
	require.NoError(t, err)
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
}

func Test_samlCertificateAndKey(t *testing.T) {
	now := time.Now()
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	otherRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	pkcs1Key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rsaKey)})
	pkcs8DER, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	require.NoError(t, err)
	pkcs8Key := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8DER})
	ecDER, err := x509.MarshalPKCS8PrivateKey(ecKey)
	require.NoError(t, err)
	ecPKCS8Key := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: ecDER})

	validCert := newTestCertificate(t, rsaKey, &rsaKey.PublicKey, now.Add(-time.Hour), now.Add(time.Hour))
	caCert := newTestCertificate(t, otherRSAKey, &otherRSAKey.PublicKey, now.Add(-time.Hour), now.Add(time.Hour))
	expiredCert := newTestCertificate(t, rsaKey, &rsaKey.PublicKey, now.Add(-2*time.Hour), now.Add(-time.Hour))
	notYetValidCert := newTestCertificate(t, rsaKey, &rsaKey.PublicKey, now.Add(time.Hour), now.Add(2*time.Hour))
	otherCert := newTestCertificate(t, otherRSAKey, &otherRSAKey.PublicKey, now.Add(-time.Hour), now.Add(time.Hour))
	ecCert := newTestCertificate(t, ecKey, &ecKey.PublicKey, now.Add(-time.Hour), now.Add(time.Hour))

	type args struct {
		certificate []byte
		key         []byte
	}
	tests := []struct {
		name     string
		args     args
		wantKey  []byte
		wantCert []byte
		wantErr  error
	}{
		{
			name:    "certificate missing",
			args:    args{key: pkcs1Key},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-s45bsj7m2u", "Errors.IDP.SAMLCertificateMissing"),
		},
		{
			name:    "key missing",
			args:    args{certificate: validCert},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-ru47reo7oq", "Errors.IDP.PrivateKeyMissing"),
		},
		{
			name:    "certificate not pem",
			args:    args{certificate: []byte("certificate"), key: pkcs1Key},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-khszemchha", "Errors.IDP.SAMLCertificateInvalid"),
		},
		{
			name: "certificate invalid",
			args: args{
				certificate: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("invalid")}),
				key:         pkcs1Key,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-ko5hxclbqm", "Errors.IDP.SAMLCertificateInvalid"),
		},
		{
			name:    "certificate expired",
			args:    args{certificate: expiredCert, key: pkcs1Key},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-3in4sinihu", "Errors.IDP.SAMLCertificateExpired"),
		},
		{
			name:    "certificate not yet valid",
			args:    args{certificate: notYetValidCert, key: pkcs1Key},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-3in4sinihu", "Errors.IDP.SAMLCertificateExpired"),
		},
		{
			name:    "key not pem",
			args:    args{certificate: validCert, key: []byte("key")},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-n6c1ub0yfk", "Errors.IDP.InvalidPrivateKey"),
		},
		{
			name: "key invalid",
			args: args{
				certificate: validCert,
				key:         pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: []byte("invalid")}),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-hadzypzq6y", "Errors.IDP.InvalidPrivateKey"),
		},
		{
			name:    "key not rsa",
			args:    args{certificate: ecCert, key: ecPKCS8Key},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-qv3lw8t2ze", "Errors.IDP.InvalidPrivateKey"),
		},
		{
			name:    "key does not match certificate",
			args:    args{certificate: otherCert, key: pkcs1Key},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-b3wng4a4hi", "Errors.IDP.SAMLCertificateKeyMismatch"),
		},
		{
			name:     "pkcs1 key",
			args:     args{certificate: validCert, key: pkcs1Key},
			wantKey:  pkcs1Key,
			wantCert: validCert,
		},
		{
			name:     "pkcs8 key",
			args:     args{certificate: validCert, key: pkcs8Key},
			wantKey:  pkcs1Key,
			wantCert: validCert,
		},
		{
			name: "certificate chain with additional text",
			args: args{
				certificate: append(append([]byte("subject=CN = SP\n"), validCert...), caCert...),
				key:         append([]byte("Bag Attributes\n"), pkcs8Key...),
			},
			wantKey:  pkcs1Key,
			wantCert: append(append([]byte{}, validCert...), caCert...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotCert, err := samlCertificateAndKey(tt.args.certificate, tt.args.key, now)
			if tt.wantErr != nil {
				assert.True(t, errors.Is(err, tt.wantErr), "got wrong err: %v", err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantKey, gotKey)
			assert.Equal(t, tt.wantCert, gotCert)
			// the SAML provider loads the key pair the same way
			_, err = tls.X509KeyPair(gotCert, gotKey)
			assert.NoError(t, err)
		})
	}
}
