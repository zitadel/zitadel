package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/api/database"
	coredb "github.com/caos/zitadel/operator/database/kinds/databases/core"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func CrdGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) (string, string, error) {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		return "", "", err
	}

	return getConnectionInfo(monitor, k8sClient, desired, false)
}

func GitOpsGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) (string, string, error) {
	desired, err := api.ReadDatabaseYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return "", "", err
	}

	return getConnectionInfo(monitor, k8sClient, desired, true)
}

func getConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
	gitOps bool,
) (string, string, error) {
	current := &tree.Tree{}

	query, _, _, _, _, err := orbdb.AdaptFunc("", nil, gitOps, "database")(monitor, desired, current)
	if err != nil {
		return "", "", err
	}

	queried := map[string]interface{}{}
	_, err = query(k8sClient, queried)
	if err != nil {
		return "", "", err
	}
	currentDB, err := coredb.ParseQueriedForDatabase(queried)
	if err != nil {
		return "", "", err
	}
	return currentDB.GetURL(), currentDB.GetPort(), nil
}
