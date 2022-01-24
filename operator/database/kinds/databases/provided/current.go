package provided

import (
	"crypto/rsa"

	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
)

var _ core.DatabaseCurrent = (*Current)(nil)

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current struct {
		URL            string
		Port           string
		QueryParams    []string
		CertificateKey *rsa.PrivateKey
	}
}

func (c *Current) GetURL() string { return c.Current.URL }

func (c *Current) GetPort() string { return c.Current.Port }

func (c *Current) GetQueryParams() []string { return c.Current.QueryParams }

func (c *Current) GetAddUserFunc() func(user string) (operator.QueryFunc, error) {
	return func(_ string) (operator.QueryFunc, error) {
		return func(_ kubernetes.ClientInt, _ map[string]interface{}) (operator.EnsureFunc, error) {
			return func(_ kubernetes.ClientInt) error { return nil }, nil
		}, nil
	}
}

func (c *Current) GetDeleteUserFunc() func(user string) (operator.DestroyFunc, error) {
	return func(_ string) (operator.DestroyFunc, error) {
		return func(k8sClient kubernetes.ClientInt) error {
			return nil
		}, nil
	}
}

func (c *Current) GetListUsersFunc() func(k8sClient kubernetes.ClientInt) ([]string, error) {
	return func(_ kubernetes.ClientInt) ([]string, error) { return nil, nil }
}
