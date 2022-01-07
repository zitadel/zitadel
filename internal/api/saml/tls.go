package saml

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

// ConfigureTLS not requiring users to present client certificates.
func ConfigureTLS(certpath, keypath, capath string) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		//Some but not all operations will require a client cert
		ClientAuth: tls.VerifyClientCertIfGiven,
		MinVersion: tls.VersionTLS12,
	}
	if capath != "" {
		caCert, err := ioutil.ReadFile(capath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		tlsConfig.ClientCAs = caCertPool
	}

	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}
