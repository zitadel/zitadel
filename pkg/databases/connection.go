package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/api/database"
	"github.com/caos/zitadel/operator/api/zitadel"
	coredb "github.com/caos/zitadel/operator/database/kinds/databases/core"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/pkg/errors"
)

func CrdGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) (string, string, error) {

	return getConnectionInfo(monitor, k8sClient, false, func() (*tree.Tree, error) {
		return zitadel.ReadCrd(k8sClient)
	}, func() (*tree.Tree, error) {
		return database.ReadCrd(k8sClient)
	})
}

func GitOpsGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) (string, string, error) {

	return getConnectionInfo(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	})
}

func getConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	desiredZitadel func() (*tree.Tree, error),
	desiredDatabase func() (*tree.Tree, error),
) (string, string, error) {
	current := &tree.Tree{}

	zitadelTree, err := desiredZitadel()
	if err != nil {
		return "", "", err
	}

	query, _, _, _, _, _, err := orbzit.AdaptFunc("", nil, gitOps, []string{"dbconnection"}, nil)(monitor, zitadelTree, current)
	if err != nil {
		return "", "", err
	}

	currentDB, err := parse(k8sClient, query)
	noCurrentState := errors.Is(err, coredb.ErrNoCurrentState)
	if err != nil && !noCurrentState {
		return "", "", err
	}

	if noCurrentState {
		databaseTree, err := desiredDatabase()
		if err != nil {
			return "", "", err
		}

		query, _, _, _, _, _, err = orbdb.AdaptFunc("", nil, gitOps, "database")(monitor, databaseTree, current)
		if err != nil {
			return "", "", err
		}

		currentDB, err = parse(k8sClient, query)
		if err != nil {
			return "", "", err
		}
	}

	return currentDB.GetURL(), currentDB.GetPort(), nil
}

func parse(k8sClient kubernetes.ClientInt, query operator.QueryFunc) (coredb.DatabaseCurrent, error) {
	queried := map[string]interface{}{}
	_, err := query(k8sClient, queried)
	if err != nil {
		return nil, err
	}
	return coredb.ParseQueriedForDatabase(queried)
}
