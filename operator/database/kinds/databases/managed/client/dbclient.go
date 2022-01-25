package client

import (
	"crypto/rsa"

	"github.com/caos/zitadel/operator/database/kinds/databases/core"
)

type ManagedDatabase interface {
	core.SecureDatabase
	GetCertificateKey() *rsa.PrivateKey
	SetCertificateKey(*rsa.PrivateKey)
	GetCertificate() []byte
	SetCertificate([]byte)
}
