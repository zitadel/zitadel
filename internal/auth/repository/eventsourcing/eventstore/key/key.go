package key

import "gopkg.in/square/go-jose.v2"

type CertificateAndKey struct {
	Certificate *jose.SigningKey
	Key         *jose.SigningKey
}
