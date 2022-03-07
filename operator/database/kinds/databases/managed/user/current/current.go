package current

import (
	"crypto/rsa"
)

type Current struct {
	CertificateKey *rsa.PrivateKey
	Certificate    []byte
}
