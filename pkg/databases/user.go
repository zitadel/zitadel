package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	"github.com/caos/zitadel/operator/api/zitadel"
)

func CrdListUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) ([]string, error) {
	return listUsers(monitor, k8sClient, false,
		func() (*tree.Tree, error) {
			return database.ReadCrd(k8sClient)
		}, func() (*tree.Tree, error) {
			return zitadel.ReadCrd(k8sClient)
		})
}

func GitOpsListUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) ([]string, error) {
	return listUsers(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	})
}

func listUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	databaseTree func() (*tree.Tree, error),
	zitadelTree func() (*tree.Tree, error),
) ([]string, error) {

	currentDB, _, err := queryDatabase(monitor, k8sClient, gitOps, databaseTree, zitadelTree)
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
	return addUser(monitor, user, k8sClient, false,
		func() (*tree.Tree, error) {
			return database.ReadCrd(k8sClient)
		}, func() (*tree.Tree, error) {
			return zitadel.ReadCrd(k8sClient)
		})
}

func GitOpsAddUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	return addUser(monitor, user, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	})
}

func addUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	databaseTree func() (*tree.Tree, error),
	zitadelTree func() (*tree.Tree, error),
) error {
	currentDB, queried, err := queryDatabase(monitor, k8sClient, gitOps, databaseTree, zitadelTree)
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
	return deleteUser(monitor, user, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	})
}

func CrdDeleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
) error {
	return deleteUser(monitor, user, k8sClient, false,
		func() (*tree.Tree, error) {
			return database.ReadCrd(k8sClient)
		}, func() (*tree.Tree, error) {
			return zitadel.ReadCrd(k8sClient)
		})
}

func deleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	databaseTree func() (*tree.Tree, error),
	zitadelTree func() (*tree.Tree, error),
) error {
	currentDB, _, err := queryDatabase(monitor, k8sClient, gitOps, databaseTree, zitadelTree)
	if err != nil {
		return err
	}
	delUser, err := currentDB.GetDeleteUserFunc()(user)
	if err != nil {
		return err
	}
	return delUser(k8sClient)
}
