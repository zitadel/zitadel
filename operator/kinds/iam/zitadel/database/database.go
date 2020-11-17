package database

import (
	"context"
	"errors"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
)

type Current struct {
	Host string
	Port string
}

func SetDatabaseInQueried(queried map[string]interface{}, current *Current) {
	queried["database"] = current
}

func GetDatabaseInQueried(queried map[string]interface{}) (*Current, error) {
	curr, ok := queried["database"].(*Current)
	if !ok {
		return nil, errors.New("Database current not in supported format")
	}

	return curr, nil
}

func newGit(monitor mntr.Monitor, repoURL string, repoKey string) (*git.Client, error) {
	gitClient := git.New(context.Background(), monitor, "orbos", "orbos@caos.ch")
	if err := gitClient.Configure(repoURL, []byte(repoKey)); err != nil {
		monitor.Error(err)
		return nil, err
	}

	if err := gitClient.Clone(); err != nil {
		return nil, err
	}
	return gitClient, nil
}

func GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, repoURL string, repoKey string) (string, string, error) {
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
