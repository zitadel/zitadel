package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/kubernetes"
)

func (c *Client) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return databases.DeleteUser(
		monitor,
		user,
		k8sClient,
		c.gitClient,
	)
}

func (c *Client) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return databases.AddUser(
		monitor,
		user,
		k8sClient,
		c.gitClient,
	)
}

func (c *Client) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return databases.ListUsers(
		monitor,
		k8sClient,
		c.gitClient,
	)
}
