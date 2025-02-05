package network

import (
	"crypto/tls"
	"errors"
	"os"
	"time"
)

var (
	ErrMissingConfig = errors.New("TLS is enabled: please specify a key (path) and a cert (path) or disable TLS if needed (e.g. by setting flag `--tlsMode external` or `--tlsMode disabled")
)

type TLS struct {
	//If enabled, ZITADEL will serve all traffic over TLS (HTTPS and gRPC)
	//you must then also provide a private key and certificate to be used for the connection
	//either directly or by a path to the corresponding file
	Enabled bool
	//Path to the private key of the TLS certificate, it will be loaded into the Key
	//and overwrite any exising value
	KeyPath string
	//Path to the certificate for the TLS connection, it will be loaded into the Cert
	//and overwrite any exising value
	CertPath string
	//Private key of the TLS certificate (KeyPath will this overwrite, if specified)
	Key []byte
	//Certificate for the TLS connection (CertPath will this overwrite, if specified)
	Cert []byte

	// cachedCert and cachedCertModTime are used to cache the parsed certificate
	cachedCert        *tls.Certificate
	cachedCertModTime time.Time
	cachedKeyModTime  time.Time
}

func (t *TLS) getCert() ([]byte, error) {
	if t.CertPath != "" {
		return os.ReadFile(t.CertPath)
	}
	return t.Cert, nil
}

func (t *TLS) getKey() ([]byte, error) {
	if t.KeyPath != "" {
		return os.ReadFile(t.KeyPath)
	}
	return t.Key, nil
}

func (t *TLS) updateCachedKeyPair() error {
	cert, err := t.getCert()
	if err != nil {
		return err
	}
	key, err := t.getKey()
	if err != nil {
		return err
	}
	if t.Key == nil || t.Cert == nil {
		return ErrMissingConfig
	}
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}
	t.cachedCert = &tlsCert
	return nil
}

func (t *TLS) Config() (_ *tls.Config, err error) {
	if !t.Enabled {
		return nil, nil
	}
	if err := t.updateCachedKeyPair(); err != nil {
		return nil, err
	}
	return &tls.Config{
		GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			updated := false
			if t.CertPath != "" {
				info, err := os.Stat(t.CertPath)
				if err != nil {
					return nil, err
				}
				if info.ModTime() != t.cachedCertModTime {
					updated = true
					t.cachedCertModTime = info.ModTime()
				}
			}
			if t.KeyPath != "" {
				info, err := os.Stat(t.KeyPath)
				if err != nil {
					return nil, err
				}
				if info.ModTime() != t.cachedKeyModTime {
					updated = true
					t.cachedKeyModTime = info.ModTime()
				}
			}
			if updated {
				if err := t.updateCachedKeyPair(); err != nil {
					return nil, err
				}
			}
			return t.cachedCert, nil
		},
	}, nil
}
