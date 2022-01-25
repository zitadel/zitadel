package databases

import (
	"context"

	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/orbos/pkg/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
)

var _ db.Client = (*GitOpsClient)(nil)
var _ db.Client = (*CrdClient)(nil)

type GitOpsClient struct {
	Monitor   mntr.Monitor
	gitClient *git.Client
}

func NewClient(monitor mntr.Monitor, gitops bool, orbconfig *orb.Orb) (db.Client, error) {
	if gitops {
		return newGitOpsClient(monitor, orbconfig.URL, orbconfig.Repokey)
	}
	return newCrdClient(monitor), nil
}

func newGitOpsClient(monitor mntr.Monitor, repoURL, repoKey string) (*GitOpsClient, error) {
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
	return GitOpsGetConnectionInfo(
		monitor,
		k8sClient,
		c.gitClient,
	)
}

type CrdClient struct {
	Monitor mntr.Monitor
}

func newCrdClient(monitor mntr.Monitor) *CrdClient {
	return &CrdClient{
		Monitor: monitor,
	}
}

func (c *CrdClient) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, error) {
	return CrdGetConnectionInfo(
		monitor,
		k8sClient,
	)
}

func (c *GitOpsClient) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return GitOpsDeleteUser(
		monitor,
		user,
		k8sClient,
		c.gitClient,
	)
}

func (c *GitOpsClient) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return GitOpsAddUser(
		monitor,
		user,
		k8sClient,
		c.gitClient,
	)
}

func (c *GitOpsClient) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return GitOpsListUsers(
		monitor,
		k8sClient,
		c.gitClient,
	)
}

func (c *CrdClient) DeleteUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return CrdDeleteUser(
		monitor,
		user,
		k8sClient,
	)
}

func (c *CrdClient) AddUser(monitor mntr.Monitor, user string, k8sClient kubernetes.ClientInt) error {
	return CrdAddUser(
		monitor,
		user,
		k8sClient,
	)
}

func (c *CrdClient) ListUsers(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) ([]string, error) {
	return CrdListUsers(
		monitor,
		k8sClient,
	)
}
