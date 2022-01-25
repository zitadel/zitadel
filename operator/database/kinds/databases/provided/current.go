package provided

import (
	"crypto/rsa"

	"github.com/caos/orbos/mntr"
	"github.com/caos/zitadel/operator/database/kinds/databases/core/certificate"

	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
)

var _ core.SecureDatabase = (*Current)(nil)

type Current struct {
	Common            *tree.Common `yaml:",inline"`
	URL, Port         string
	QueryParams       []string
	getAddUserFunc    func(user string) (operator.QueryFunc, error)
	getDeleteUserFunc func(user string) (operator.DestroyFunc, error)
	getListUserFunc   func(k8sClient kubernetes.ClientInt) ([]string, error)
	monitor           mntr.Monitor
	namespace         string
}

func (c *Current) GetURL() string { return c.URL }

func (c *Current) GetPort() string { return c.Port }

func (c *Current) GetQueryParams() []string { return c.QueryParams }

func (c *Current) GetCertificateKey() *rsa.PrivateKey { return nil }

func (c *Current) SetCertificateKey(*rsa.PrivateKey) {}
func (c *Current) GetCertificate() []byte            { return nil }
func (c *Current) SetCertificate([]byte)             {}
func (c *Current) GetAddUserFunc() func(_ string) (operator.QueryFunc, error) {

	if c.getAddUserFunc == nil {
		certificate.AdaptFunc(c.monitor, c.namespace)
	}

	return func(_ string) (operator.QueryFunc, error) {
		return func(_ kubernetes.ClientInt, _ map[string]interface{}) (operator.EnsureFunc, error) {
			return func(_ kubernetes.ClientInt) error { return nil }, nil
		}, nil
	}
}
func (c *Current) GetDeleteUserFunc() func(_ string) (operator.DestroyFunc, error) {
	return func(_ string) (operator.DestroyFunc, error) {
		return func(_ kubernetes.ClientInt) error { return nil }, nil
	}
}
func (c *Current) GetListUsersFunc() func(_ kubernetes.ClientInt) ([]string, error) {
	return func(_ kubernetes.ClientInt) ([]string, error) { return nil, nil }
}
