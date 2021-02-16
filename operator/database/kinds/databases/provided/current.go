package provided

import (
	"github.com/caos/orbos/pkg/tree"
)

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
