package cockroachdb

import (
	"crypto/rsa"

	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate"
)

var _ core.DatabaseCurrent = (*Current)(nil)

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current struct {
		Host              string
		Port              string
		Cluster           string
		User              string
		PasswordSecret    *labels.Selectable
		PasswordSecretKey string
		Secure            bool
		AddUserFunc       func(user string) (operator.QueryFunc, error)
		DeleteUserFunc    func(user string) (operator.DestroyFunc, error)
		ListUsersFunc     func(k8sClient kubernetes.ClientInt) ([]string, error)
		CA                certificate.Current
	}
}

func (c *Current) GetURL() string { return c.Current.Host }

func (c *Current) GetPort() string { return c.Current.Port }

func (c *Current) GetReadyQuery() operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error { return nil }
}

func (c *Current) GetCertificateKey() *rsa.PrivateKey { return c.Current.CA.CertificateKey }

func (c *Current) SetCertificateKey(key *rsa.PrivateKey) { c.Current.CA.CertificateKey = key }

func (c *Current) GetCertificate() []byte { return c.Current.CA.Certificate }

func (c *Current) SetCertificate(cert []byte) { c.Current.CA.Certificate = cert }

func (c *Current) GetAddUserFunc() func(user string) (operator.QueryFunc, error) {
	return c.Current.AddUserFunc
}

func (c *Current) GetDeleteUserFunc() func(user string) (operator.DestroyFunc, error) {
	return c.Current.DeleteUserFunc
}

func (c *Current) GetListUsersFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return c.Current.ListUsersFunc
}

func (c *Current) GetListDatabasesFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return nil
}

/*

func (c *Current) Host() string { return c.Current.Host }
func (c *Current) Port() string { return c.Current.Port }
func (c *Current) User() string { return c.Current.User }
func (c *Current) PasswordSecret() (*labels.Selectable, string) {
	return c.Current.PasswordSecret, c.Current.PasswordSecretKey
}

func (c *Current) SSL() *db.SSL {
	return &db.SSL{
		RootCert:       c.Current.Secure,
		UserCertAndKey: false,
	}
}

func (c *Current) Options() string {
	if c.Current.Cluster != "" {
		return "--cluster=" + c.Current.Cluster
	}
	return ""
}

func (c *Current) ConnectionParams(certsDir string) string {

	var params []string
	certs := fmt.Sprintf("sslmode=verify-full&sslrootcert=%s/client.%s.crt", certsDir, c.Current.User)
	if !c.Current.Secure {
		certs = "sslmode=disable"
	}
	params = append(params, certs)

	if c.Current.Cluster != "" {
		params = append(params, "options=--cluster%3D"+c.Current.Cluster)
	}

	return strings.Join(params, "&")
}


*/
