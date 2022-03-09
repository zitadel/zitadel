package cockroachdb

import (
	"fmt"
	"github.com/caos/orbos/pkg/labels"
	"strings"

	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/pkg/databases/db"
)

var _ db.Connection = (*Current)(nil)

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
	}
}

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
