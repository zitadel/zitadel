package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/kubernetes"
)

func DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt, repoURL string, repoKey string) error {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return err
	}

	return databases.DeleteUser(
		monitor,
		user,
		k8sClient,
		gitClient,
	)
}

func AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt, repoURL string, repoKey string) error {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return err
	}

	return databases.AddUser(
		monitor,
		user,
		k8sClient,
		gitClient,
	)
}

func ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, repoURL string, repoKey string) ([]string, error) {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return nil, err
	}

	return databases.ListUsers(
		monitor,
		k8sClient,
		gitClient,
	)
}
