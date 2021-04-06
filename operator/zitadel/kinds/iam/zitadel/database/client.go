package database

import (
	"context"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/pkg/databases"
)

var _ Client = (*GitOpsClient)(nil)
var _ Client = (*CrdClient)(nil)

type Client interface {
	GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error)
	DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error
	ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error)
}

type GitOpsClient struct {
	Monitor   mntr.Monitor
	gitClient *git.Client
}

func NewGitOpsClient(monitor mntr.Monitor, repoURL, repoKey string) (*GitOpsClient, error) {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return nil, err
	}

	return &GitOpsClient{
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

func (c *GitOpsClient) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error) {
	return databases.GitOpsGetConnectionInfo(
		monitor,
		k8sClient,
		c.gitClient,
	)
}

type CrdClient struct {
	Monitor mntr.Monitor
}

func NewCrdClient(monitor mntr.Monitor) *CrdClient {
	return &CrdClient{
		Monitor: monitor,
	}
}

func (c *CrdClient) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error) {
	return databases.CrdGetConnectionInfo(
		monitor,
		k8sClient,
	)
}
