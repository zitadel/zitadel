package client

import (
	"crypto/rsa"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
)

type ManagedDatabase interface {
	core.DatabaseCurrent
	GetReadyQuery() operator.EnsureFunc
	GetCertificateKey() *rsa.PrivateKey
	SetCertificateKey(*rsa.PrivateKey)
	GetCertificate() []byte
	SetCertificate([]byte)
	//	GetListDatabasesFunc() func(k8sClient kubernetes.ClientInt) ([]string, error)
}
