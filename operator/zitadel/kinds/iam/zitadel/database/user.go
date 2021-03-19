package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/pkg/databases"
)

func (c *GitOpsClient) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return databases.GitOpsDeleteUser(
		monitor,
		user,
		k8sClient,
		c.gitClient,
	)
}

func (c *GitOpsClient) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return databases.GitOpsAddUser(
		monitor,
		user,
		k8sClient,
		c.gitClient,
	)
}

func (c *GitOpsClient) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return databases.GitOpsListUsers(
		monitor,
		k8sClient,
		c.gitClient,
	)
}

func (c *CrdClient) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return databases.CrdDeleteUser(
		monitor,
		user,
		k8sClient,
	)
}

func (c *CrdClient) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return databases.CrdAddUser(
		monitor,
		user,
		k8sClient,
	)
}

func (c *CrdClient) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return databases.CrdListUsers(
		monitor,
		k8sClient,
	)
}
