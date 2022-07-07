package network

import (
	"crypto/tls"
	"errors"
	"os"
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
}

func (t *TLS) Config() (_ *tls.Config, err error) {
	if !t.Enabled {
		return nil, nil
	}
	if t.KeyPath != "" {
		t.Key, err = os.ReadFile(t.KeyPath)
		if err != nil {
			return nil, err
		}
	}
	if t.CertPath != "" {
		t.Cert, err = os.ReadFile(t.CertPath)
		if err != nil {
			return nil, err
		}
	}
	if t.Key == nil || t.Cert == nil {
		return nil, ErrMissingConfig
	}
	tlsCert, err := tls.X509KeyPair(t.Cert, t.Key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}, nil
}
