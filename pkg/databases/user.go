package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	coredb "github.com/caos/zitadel/operator/database/kinds/databases/core"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func ListUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) ([]string, error) {
	desired, err := api.ReadDatabaseYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}
	current := &tree.Tree{}

	query, _, _, err := orbdb.AdaptFunc("", nil)(monitor, desired, current)
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

func AddUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	desired, err := api.ReadDatabaseYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return err
	}
	current := &tree.Tree{}

	query, _, _, err := orbdb.AdaptFunc("", nil)(monitor, desired, current)
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

func DeleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	desired, err := api.ReadDatabaseYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return err
	}
	current := &tree.Tree{}

	query, _, _, err := orbdb.AdaptFunc("", nil)(monitor, desired, current)
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
