package network

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

var (
	ErrMissingConfig = errors.New("TLS is enabled: please specify a key (path) and a cert (path) or disable TLS if needed (e.g. by setting flag `--tlsMode external` or `--tlsMode disabled")
)

type TLS struct {
	//If enabled, ZITADEL will serve all traffic over TLS (HTTPS and gRPC)
	//you must then also provide a private key and certificate to be used for the connection
	//either directly or by a path to the corresponding file
	Enabled bool
	//Path to the private key of the TLS certificate, it will be used instead of Key.
	//The file is reloaded automatically when it changes.
	KeyPath string
	//Path to the certificate for the TLS connection, it will be used instead of Cert.
	//The file is reloaded automatically when it changes.
	CertPath string
	//Private key of the TLS certificate (KeyPath will this overwrite, if specified)
	Key []byte
	//Certificate for the TLS connection (CertPath will this overwrite, if specified)
	Cert []byte
}

// Config returns the TLS config for the server.
// If CertPath or KeyPath is set, the files are checked on every handshake
// and the certificate is reloaded when one of them changed.
func (t *TLS) Config() (*tls.Config, error) {
	if !t.Enabled {
		return nil, nil
	}
	reloader := &certReloader{
		certPath: t.CertPath,
		keyPath:  t.KeyPath,
		certPEM:  t.Cert,
		keyPEM:   t.Key,
	}
	if err := reloader.init(); err != nil {
		return nil, err
	}
	if t.CertPath == "" && t.KeyPath == "" {
		return &tls.Config{
			Certificates: []tls.Certificate{*reloader.cert},
		}, nil
	}
	return &tls.Config{
		GetCertificate: reloader.getCertificate,
	}, nil
}

// certReloader holds the parsed key pair and the modification times
// of the files it was loaded from.
// It is kept apart from [TLS] so the config struct can still be copied by value.
type certReloader struct {
	certPath, keyPath string
	// certPEM and keyPEM are used when the matching path is empty
	certPEM, keyPEM []byte

	mu          sync.RWMutex
	cert        *tls.Certificate
	certModTime time.Time
	keyModTime  time.Time
}

func (r *certReloader) init() error {
	certModTime, keyModTime, err := r.modTimes()
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.load(certModTime, keyModTime)
}

func (r *certReloader) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certModTime, keyModTime, err := r.modTimes()
	if err != nil {
		// The file might be replaced right now, keep serving the current certificate.
		logging.WithError(helloContext(hello), err).Warn("unable to check TLS certificate files, keeping current certificate")
		r.mu.RLock()
		defer r.mu.RUnlock()
		return r.cert, nil
	}

	r.mu.RLock()
	if r.unchanged(certModTime, keyModTime) {
		defer r.mu.RUnlock()
		return r.cert, nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()
	// Another handshake might have reloaded the files in the meantime.
	if r.unchanged(certModTime, keyModTime) {
		return r.cert, nil
	}
	if err := r.load(certModTime, keyModTime); err != nil {
		// The modification times are stored anyway,
		// so the next change of the files triggers a new attempt.
		r.certModTime, r.keyModTime = certModTime, keyModTime
		logging.WithError(helloContext(hello), err).Warn("unable to reload TLS certificate, keeping current certificate")
		return r.cert, nil
	}
	logging.Info(helloContext(hello), "TLS certificate reloaded")
	return r.cert, nil
}

// unchanged must be called while holding r.mu.
func (r *certReloader) unchanged(certModTime, keyModTime time.Time) bool {
	return certModTime.Equal(r.certModTime) && keyModTime.Equal(r.keyModTime)
}

// load reads and parses the key pair.
// It must be called while holding the write lock of r.mu.
func (r *certReloader) load(certModTime, keyModTime time.Time) error {
	certPEM, err := readOrDefault(r.certPath, r.certPEM)
	if err != nil {
		return err
	}
	keyPEM, err := readOrDefault(r.keyPath, r.keyPEM)
	if err != nil {
		return err
	}
	if len(certPEM) == 0 || len(keyPEM) == 0 {
		return ErrMissingConfig
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return err
	}
	r.cert = &cert
	r.certModTime, r.keyModTime = certModTime, keyModTime
	return nil
}

// modTimes returns the modification times of the cert and key files.
// The time is zero for a value that is not loaded from a file.
func (r *certReloader) modTimes() (certModTime, keyModTime time.Time, err error) {
	if certModTime, err = modTime(r.certPath); err != nil {
		return time.Time{}, time.Time{}, err
	}
	if keyModTime, err = modTime(r.keyPath); err != nil {
		return time.Time{}, time.Time{}, err
	}
	return certModTime, keyModTime, nil
}

func modTime(path string) (time.Time, error) {
	if path == "" {
		return time.Time{}, nil
	}
	// os.Stat follows symlinks, so updates of mounted Kubernetes secrets are detected as well.
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

func readOrDefault(path string, value []byte) ([]byte, error) {
	if path == "" {
		return value, nil
	}
	return os.ReadFile(path)
}

func helloContext(hello *tls.ClientHelloInfo) context.Context {
	if hello != nil {
		if ctx := hello.Context(); ctx != nil {
			return ctx
		}
	}
	return context.Background()
}
