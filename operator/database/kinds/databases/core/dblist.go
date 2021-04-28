package core

import (
	"crypto/rsa"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
)

var current DatabaseCurrent = &CurrentDBList{}

type CurrentDBList struct {
	Common  *tree.Common `yaml:",inline"`
	Current *DatabaseCurrentDBList
}

type DatabaseCurrentDBList struct {
	Databases []string
	Users     []string
}

func (c *CurrentDBList) GetURL() string {
	return ""
}

func (c *CurrentDBList) GetPort() string {
	return ""
}

func (c *CurrentDBList) GetReadyQuery() operator.EnsureFunc {
	return nil
}

func (c *CurrentDBList) GetCertificateKey() *rsa.PrivateKey {
	return nil
}

func (c *CurrentDBList) SetCertificateKey(key *rsa.PrivateKey) {
	return
}

func (c *CurrentDBList) GetCertificate() []byte {
	return nil
}

func (c *CurrentDBList) SetCertificate(cert []byte) {
	return
}

func (c *CurrentDBList) GetListDatabasesFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return func(k8sClient kubernetes.ClientInt) ([]string, error) {
		return c.Current.Databases, nil
	}
}

func (c *CurrentDBList) GetListUsersFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return func(k8sClient kubernetes.ClientInt) ([]string, error) {
		return c.Current.Users, nil
	}
}

func (c *CurrentDBList) GetAddUserFunc() func(user string) (operator.QueryFunc, error) {
	return nil
}

func (c *CurrentDBList) GetDeleteUserFunc() func(user string) (operator.DestroyFunc, error) {
	return nil
}
