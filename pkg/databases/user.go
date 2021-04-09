package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	coredb "github.com/caos/zitadel/operator/database/kinds/databases/core"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func CrdListUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) ([]string, error) {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return listUsers(monitor, k8sClient, desired)
}

func GitOpsListUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) ([]string, error) {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return listUsers(monitor, k8sClient, desired)
}

func listUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
) ([]string, error) {
	current := &tree.Tree{}

	query, _, _, _, _, err := orbdb.AdaptFunc("", nil, false, "database")(monitor, desired, current)
	if err != nil {
		return nil, err
	}

	queried := map[string]interface{}{}
	_, err = query(k8sClient, queried)
	if err != nil {
		return nil, err
	}
	currentDB, err := coredb.ParseQueriedForDatabase(queried)
	if err != nil {
		return nil, err
	}

	list, err := currentDB.GetListUsersFunc()(k8sClient)
	if err != nil {
		return nil, err
	}

	users := []string{}
	for _, listedUser := range list {
		if listedUser != "root" {
			users = append(users, listedUser)
		}
	}

	return users, nil
}

func CrdAddUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return err
	}
	return addUser(monitor, user, k8sClient, desired)
}

func GitOpsAddUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		monitor.Error(err)
		return err
	}
	return addUser(monitor, user, k8sClient, desired)
}

func addUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
) error {
	current := &tree.Tree{}

	query, _, _, _, _, err := orbdb.AdaptFunc("", nil, false, "database")(monitor, desired, current)
	if err != nil {
		return err
	}

	queried := map[string]interface{}{}
	_, err = query(k8sClient, queried)
	if err != nil {
		return err
	}
	currentDB, err := coredb.ParseQueriedForDatabase(queried)
	if err != nil {
		return err
	}

	queryUser, err := currentDB.GetAddUserFunc()(user)
	if err != nil {
		return err
	}
	ensureUser, err := queryUser(k8sClient, queried)
	if err != nil {
		return err
	}
	return ensureUser(k8sClient)
}

func GitOpsDeleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		monitor.Error(err)
		return err
	}

	return deleteUser(monitor, user, k8sClient, desired)
}

func CrdDeleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return err
	}

	return deleteUser(monitor, user, k8sClient, desired)
}

func deleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
) error {
	current := &tree.Tree{}

	query, _, _, _, _, err := orbdb.AdaptFunc("", nil, false, "database")(monitor, desired, current)
	if err != nil {
		return err
	}

	queried := map[string]interface{}{}
	_, err = query(k8sClient, queried)
	if err != nil {
		return err
	}
	currentDB, err := coredb.ParseQueriedForDatabase(queried)
	if err != nil {
		return err
	}

	deleteUser, err := currentDB.GetDeleteUserFunc()(user)
	if err != nil {
		return err
	}
	return deleteUser(k8sClient)
}
