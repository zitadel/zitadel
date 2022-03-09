package cockroachdb

import (
	"github.com/caos/orbos/pkg/labels"
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
		CertsSecret       string
	}
}

func (c *Current) Host() string { return c.Current.Host }
func (c *Current) Port() string { return c.Current.Port }
func (c *Current) User() string { return c.Current.User }
func (c *Current) PasswordSecret() (*labels.Selectable, string) {
	return c.Current.PasswordSecret, c.Current.PasswordSecretKey
}

func (c *Current) SSL() *db.SSL {
	if c.Current.CertsSecret == "" {
		return nil
	}
	return &db.SSL{
		CertsSecret:    c.Current.CertsSecret,
		RootCert:       true,
		UserCertAndKey: false,
	}
}

func (c *Current) Options() string {
	if c.Current.Cluster != "" {
		return "--cluster=" + c.Current.Cluster
	}
	return ""
}
