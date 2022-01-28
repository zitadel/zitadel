package databases

import (
	"context"

	"github.com/caos/zitadel/operator/api/database"
	"github.com/caos/zitadel/operator/api/zitadel"

	"github.com/caos/orbos/pkg/tree"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/pkg/errors"

	"github.com/caos/zitadel/operator"

	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/orbos/pkg/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
)

/*
var _ db.Connection = (*GitOpsClient)(nil)
var _ db.Connection = (*CrdClient)(nil)

type GitOpsClient struct {
	Monitor   mntr.Monitor
	gitClient *git.Client
}

*/

func NewConnection(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, gitops bool, orbconfig *orb.Orb) (db.Connection, error) {
	if gitops {
		return newGitOpsConnection(monitor, k8sClient, orbconfig.URL, orbconfig.Repokey)
	}
	return newCrdConnection(monitor, k8sClient)
}

func newGitOpsConnection(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, repoURL, repoKey string) (db.Connection, error) {
	gitClient, err := newGit(monitor, repoURL, repoKey)
	if err != nil {
		return nil, err
	}

	return connection(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	})
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

/*
func (c *GitOpsClient) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, func(string) string, error) {
	return gitOpsGetConnectionInfo(
		monitor,
		k8sClient,
		c.gitClient,
	)
}

*/

type CrdClient struct {
	Monitor mntr.Monitor
}

func newCrdConnection(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (db.Connection, error) {
	return connection(monitor, k8sClient, false, func() (*tree.Tree, error) {
		return zitadel.ReadCrd(k8sClient)
	}, func() (*tree.Tree, error) {
		return database.ReadCrd(k8sClient)
	})
}

/*
func (c *CrdClient) GetConnectionInfo(monitor mntr.Monitor, k8sClient kubernetes.ClientInt) (string, string, func(string) string, error) {
	return crdGetConnectionInfo(
		monitor,
		k8sClient,
	)
}

func crdGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) (db.Connection, error) {

	return connection(monitor, k8sClient, false, func() (*tree.Tree, error) {
		return zitadel.ReadCrd(k8sClient)
	}, func() (*tree.Tree, error) {
		return database.ReadCrd(k8sClient)
	})
}

func gitOpsGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) (db.Connection, error) {

	return connection(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	})
}
*/

/*
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
*/

func connection(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	desiredZitadel func() (*tree.Tree, error),
	desiredDatabase func() (*tree.Tree, error),
) (db.Connection, error) {
	current := &tree.Tree{}

	zitadelTree, err := desiredZitadel()
	if err != nil {
		return nil, err
	}

	query, _, _, _, _, _, err := orbzit.AdaptFunc("", nil, gitOps, []string{"dbconnection"}, nil)(monitor, zitadelTree, current)
	if err != nil {
		return nil, err
	}

	queriedConn, err := parse(k8sClient, query)
	noCurrentState := errors.Is(err, db.ErrNoCurrentState)
	if err != nil && !noCurrentState {
		return nil, err
	}

	if noCurrentState {
		databaseTree, err := desiredDatabase()
		if err != nil {
			return nil, err
		}

		query, _, _, _, _, _, err = orbdb.AdaptFunc("", nil, gitOps, "database")(monitor, databaseTree, current)
		if err != nil {
			return nil, err
		}

		queriedConn, err = parse(k8sClient, query)
		if err != nil {
			return nil, err
		}
	}
	return queriedConn, nil
}

func parse(k8sClient kubernetes.ClientInt, query operator.QueryFunc) (db.Connection, error) {
	queried := map[string]interface{}{}
	_, err := query(k8sClient, queried)
	if err != nil {
		return nil, err
	}
	return db.ParseQueriedForDatabase(queried)
}
