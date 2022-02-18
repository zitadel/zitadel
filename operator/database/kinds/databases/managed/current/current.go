package current

import (
	"crypto/rsa"
	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/zitadel/pkg/databases/db"

	cacurr "github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/current"

	"github.com/caos/orbos/pkg/tree"
)

var _ db.Connection = (*Current)(nil)

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current *CurrentDB
}

type CurrentDB struct {
	URL               string
	Port              string
	User              string
	PasswordSecret    *labels.Selectable
	PasswordSecretKey string
	CA                *cacurr.Current

	/*
		AddUserFunc    func(user string) (operator.QueryFunc, error)
		DeleteUserFunc func(user string) (operator.DestroyFunc, error)
		ListUsersFunc  func(k8sClient kubernetes.ClientInt) ([]string, error)

	*/
}

func (c *Current) GetCA() *cacurr.Current {
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

func (c *Current) Host() string { return "cockroachdb-public" }
func (c *Current) Port() string { return "26257" }
func (c *Current) User() string { return c.Current.User }
func (c *Current) PasswordSecret() (*labels.Selectable, string) {
	return c.Current.PasswordSecret, c.Current.PasswordSecretKey
}

func (c *Current) SSL() *db.SSL {
	return &db.SSL{
		RootCert:       true,
		UserCertAndKey: true,
	}
}
func (c *Current) Options() string { return "" }

/*
func (c *Current) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	destroy, err := c.Current.DeleteUserFunc(user)
	if err != nil {
		return err
	}
	return destroy(k8sClient)
}
func (c *Current) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	query, err := c.Current.AddUserFunc(user)
	if err != nil {
		return err
	}

	ensure, err := query(k8sClient, nil)
	if err != nil {
		return err
	}
	return ensure(k8sClient)
}
func (c *Current) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return c.Current.ListUsersFunc(k8sClient)
}
*/
