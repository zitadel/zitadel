package managed

import (
	"crypto/rsa"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate"
)

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current *CurrentDB
}

type CurrentDB struct {
	URL               string
	Port              string
	ReadyFunc         operator.EnsureFunc
	CA                *certificate.Current
	AddUserFunc       func(user string) (operator.QueryFunc, error)
	DeleteUserFunc    func(user string) (operator.DestroyFunc, error)
	ListUsersFunc     func(k8sClient kubernetes.ClientInt) ([]string, error)
	ListDatabasesFunc func(k8sClient kubernetes.ClientInt) ([]string, error)
}

func (c *Current) GetURL() string {
	return c.Current.URL
}

func (c *Current) GetPort() string {
	return c.Current.Port
}

func (c *Current) GetReadyQuery() operator.EnsureFunc {
	return c.Current.ReadyFunc
}

func (c *Current) GetCA() *certificate.Current {
	return c.Current.CA
}

func (c *Current) GetCertificateKey() *rsa.PrivateKey {
	return c.Current.CA.CertificateKey
}

func (c *Current) SetCertificateKey(key *rsa.PrivateKey) {
	c.Current.CA.CertificateKey = key
}

func (c *Current) GetCertificate() []byte {
	return c.Current.CA.Certificate
}

func (c *Current) SetCertificate(cert []byte) {
	c.Current.CA.Certificate = cert
}

func (c *Current) GetListDatabasesFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return c.Current.ListDatabasesFunc
}

func (c *Current) GetListUsersFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return c.Current.ListUsersFunc
}

func (c *Current) GetAddUserFunc() func(user string) (operator.QueryFunc, error) {
	return c.Current.AddUserFunc
}

func (c *Current) GetDeleteUserFunc() func(user string) (operator.DestroyFunc, error) {
	return c.Current.DeleteUserFunc
}
