package database

import (
	"context"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
)

var _ ClientInt = (*Client)(nil)

type ClientInt interface {
	GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error)
	DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error)
}

type Client struct {
	Monitor   mntr.Monitor
	gitClient *git.Client
}

func NewClient(monitor mntr.Monitor, repoURL, repoKey string) (*Client, error) {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return nil, err
	}

	return &Client{
		Monitor:   monitor,
		gitClient: gitClient,
	}, nil
}

func newGit(monitor mntr.Monitor, repoURL string, repoKey string) (*git.Client, error) {
	gitClient := git.New(context.Background(), monitor, "orbos", "orbos@caos.ch")
	if err := gitClient.Configure(repoURL, []byte(repoKey)); err != nil {
		return nil, err
	}

	if err := gitClient.Clone(); err != nil {
		return nil, err
	}
	return gitClient, nil
}

func (c *Client) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error) {
	return databases.GetConnectionInfo(
		monitor,
		k8sClient,
		c.gitClient,
	)
}
