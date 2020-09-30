package database

import (
	"context"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
)

func newGit(monitor mntr.Monitor, repoURL string, repoKey string) (*git.Client, error) {
	gitClient := git.New(context.Background(), monitor, "orbos", "orbos@caos.ch")
	if err := gitClient.Configure(repoURL, []byte(repoKey)); err != nil {
		monitor.Error(err)
		return nil, err
	}
	return gitClient, nil
}

func GetConnectionInfo(monitor mntr.Monitor, k8sClient *kubernetes.Client, repoURL string, repoKey string) (string, string, error) {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return "", "", err
	}

	return databases.GetConnectionInfo(
		monitor,
		k8sClient,
		gitClient,
	)
}

func DeleteUser(monitor mntr.Monitor, user string, k8sClient *kubernetes.Client, repoURL string, repoKey string) error {
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

func AddUser(monitor mntr.Monitor, user string, k8sClient *kubernetes.Client, repoURL string, repoKey string) error {
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

func ListUsers(monitor mntr.Monitor, k8sClient *kubernetes.Client, repoURL string, repoKey string) ([]string, error) {
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
