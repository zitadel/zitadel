package network

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateKeyPair(t *testing.T, commonName string) (certPEM, keyPEM []byte) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
}

// writeFile writes the file and moves its modification time forward,
// so the change is detected even on file systems with a coarse time resolution.
func writeFile(t *testing.T, path string, data []byte, modTime time.Time) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, data, 0o600))
	require.NoError(t, os.Chtimes(path, modTime, modTime))
}

func commonName(t *testing.T, cert *tls.Certificate) string {
	t.Helper()
	require.NotNil(t, cert)
	parsed, err := x509.ParseCertificate(cert.Certificate[0])
	require.NoError(t, err)
	return parsed.Subject.CommonName
}

func TestTLS_Config_disabled(t *testing.T) {
	config, err := (&TLS{Enabled: false}).Config()
	assert.NoError(t, err)
	assert.Nil(t, config)
}

func TestTLS_Config_missing(t *testing.T) {
	_, err := (&TLS{Enabled: true}).Config()
	assert.ErrorIs(t, err, ErrMissingConfig)
}

func TestTLS_Config_missingFile(t *testing.T) {
	_, err := (&TLS{Enabled: true, CertPath: filepath.Join(t.TempDir(), "missing.crt")}).Config()
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestTLS_Config_static(t *testing.T) {
	certPEM, keyPEM := generateKeyPair(t, "static")
	config, err := (&TLS{Enabled: true, Cert: certPEM, Key: keyPEM}).Config()
	require.NoError(t, err)
	assert.Nil(t, config.GetCertificate)
	require.Len(t, config.Certificates, 1)
	assert.Equal(t, "static", commonName(t, &config.Certificates[0]))
}

func TestTLS_Config_reload(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := filepath.Join(dir, "tls.crt"), filepath.Join(dir, "tls.key")
	start := time.Now().Add(-time.Hour)

	certPEM, keyPEM := generateKeyPair(t, "first")
	writeFile(t, certPath, certPEM, start)
	writeFile(t, keyPath, keyPEM, start)

	config, err := (&TLS{Enabled: true, CertPath: certPath, KeyPath: keyPath}).Config()
	require.NoError(t, err)
	require.NotNil(t, config.GetCertificate)

	cert, err := config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "first", commonName(t, cert))

	// both files replaced
	certPEM, keyPEM = generateKeyPair(t, "second")
	writeFile(t, certPath, certPEM, start.Add(time.Minute))
	writeFile(t, keyPath, keyPEM, start.Add(time.Minute))
	cert, err = config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "second", commonName(t, cert))

	// only the cert is replaced, the pair doesn't match: keep the current certificate
	certPEM, keyPEM = generateKeyPair(t, "third")
	writeFile(t, certPath, certPEM, start.Add(2*time.Minute))
	cert, err = config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "second", commonName(t, cert))

	// the key follows: the new pair is loaded
	writeFile(t, keyPath, keyPEM, start.Add(2*time.Minute))
	cert, err = config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "third", commonName(t, cert))

	// files removed: keep the current certificate
	require.NoError(t, os.Remove(certPath))
	cert, err = config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "third", commonName(t, cert))
}

func TestTLS_Config_reloadAfterReadError(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := filepath.Join(dir, "tls.crt"), filepath.Join(dir, "tls.key")
	start := time.Now().Add(-time.Hour)

	certPEM, keyPEM := generateKeyPair(t, "first")
	writeFile(t, certPath, certPEM, start)
	writeFile(t, keyPath, keyPEM, start)

	config, err := (&TLS{Enabled: true, CertPath: certPath, KeyPath: keyPath}).Config()
	require.NoError(t, err)

	// the cert path can be stat'ed but not read: keep the current certificate
	changed := start.Add(time.Minute)
	require.NoError(t, os.Remove(certPath))
	require.NoError(t, os.Mkdir(certPath, 0o700))
	require.NoError(t, os.Chtimes(certPath, changed, changed))
	cert, err := config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "first", commonName(t, cert))

	// the files are restored with the same modification time: they must still be loaded
	certPEM, keyPEM = generateKeyPair(t, "second")
	require.NoError(t, os.Remove(certPath))
	writeFile(t, certPath, certPEM, changed)
	writeFile(t, keyPath, keyPEM, start)
	cert, err = config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "second", commonName(t, cert))
}

func TestTLS_Config_reloadMixed(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	start := time.Now().Add(-time.Hour)

	certPEM, keyPEM := generateKeyPair(t, "first")
	writeFile(t, certPath, certPEM, start)

	config, err := (&TLS{Enabled: true, CertPath: certPath, Key: keyPEM}).Config()
	require.NoError(t, err)
	cert, err := config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "first", commonName(t, cert))

	// a new cert for a different key can't be used with the configured key
	otherCertPEM, _ := generateKeyPair(t, "other")
	writeFile(t, certPath, otherCertPEM, start.Add(time.Minute))
	cert, err = config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "first", commonName(t, cert))
}

func TestTLS_Config_reloadConcurrent(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := filepath.Join(dir, "tls.crt"), filepath.Join(dir, "tls.key")
	start := time.Now().Add(-time.Hour)

	certPEM, keyPEM := generateKeyPair(t, "first")
	writeFile(t, certPath, certPEM, start)
	writeFile(t, keyPath, keyPEM, start)

	config, err := (&TLS{Enabled: true, CertPath: certPath, KeyPath: keyPath}).Config()
	require.NoError(t, err)

	// The files are rotated before the workers are released,
	// so all of them see the change at the same time and race for the reload.
	certPEM, keyPEM = generateKeyPair(t, "second")
	writeFile(t, keyPath, keyPEM, start.Add(time.Minute))
	writeFile(t, certPath, certPEM, start.Add(time.Minute))

	release := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-release
			for j := 0; j < 50; j++ {
				// require must not be used outside of the test goroutine
				cert, err := config.GetCertificate(&tls.ClientHelloInfo{})
				if !assert.NoError(t, err) || !assert.NotNil(t, cert) {
					return
				}
				parsed, err := x509.ParseCertificate(cert.Certificate[0])
				if assert.NoError(t, err) {
					assert.Equal(t, "second", parsed.Subject.CommonName)
				}
			}
		}()
	}
	close(release)
	wg.Wait()

	cert, err := config.GetCertificate(&tls.ClientHelloInfo{})
	require.NoError(t, err)
	assert.Equal(t, "second", commonName(t, cert))
}
