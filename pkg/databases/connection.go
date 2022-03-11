package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	"github.com/caos/zitadel/operator/api/zitadel"
)

func CrdGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) (string, string, error) {
	return getConnectionInfo(monitor, k8sClient, false,
		func() (*tree.Tree, error) {
			return database.ReadCrd(k8sClient)
		}, func() (*tree.Tree, error) {
			return zitadel.ReadCrd(k8sClient)
		})
}

func GitOpsGetConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) (string, string, error) {
	return getConnectionInfo(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	})
}

func getConnectionInfo(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	databaseTree func() (*tree.Tree, error),
	zitadelTree func() (*tree.Tree, error),
) (url string, port string, err error) {

	currentDB, _, err := queryDatabase(monitor, k8sClient, gitOps, databaseTree, zitadelTree)
	if err != nil {
		return "", "", err
	}

	return currentDB.GetURL(), currentDB.GetPort(), nil
}
