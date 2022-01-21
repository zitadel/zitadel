package provided

import (
	"crypto/rsa"

	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
)

var _ core.DatabaseCurrent = (*Current)(nil)

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current struct {
		URL  string
		Port string
	}
}

func (c *Current) GetURL() string {
	return c.Current.URL
}

func (c *Current) GetPort() string {
	return c.Current.Port
}

func (c *Current) GetReadyQuery() operator.EnsureFunc {
	return nil
}

func (c *Current) GetCertificateKey() *rsa.PrivateKey {
	return nil

}
func (c *Current) SetCertificateKey(*rsa.PrivateKey) {
}

func (c *Current) GetAddUserFunc() func(user string) (operator.QueryFunc, error) {
	return nil
}
