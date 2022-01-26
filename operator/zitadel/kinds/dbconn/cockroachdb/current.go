package cockroachdb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/pkg/databases/db"
)

var _ db.Client = (*Current)(nil)

type Current struct {
	Common  *tree.Common `yaml:",inline"`
	Current struct {
		URL  string
		Port string
	}
}

func (c *Current) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error) {
	return c.Current.URL, c.Current.Port, nil
}
func (c *Current) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return nil
}
func (c *Current) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return nil
}

func (c *Current) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return nil, nil
}
