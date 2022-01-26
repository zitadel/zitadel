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
	return listUsers(monitor, k8sClient, false, func() (*tree.Tree, error) {
		return zitadel.ReadCrd(k8sClient)
	}, func() (*tree.Tree, error) {
		return database.ReadCrd(k8sClient)
	})
}

func GitOpsListUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) ([]string, error) {
	return listUsers(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	})
}

func listUsers(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	desiredZitadel func() (*tree.Tree, error),
	desiredDatabase func() (*tree.Tree, error),
) ([]string, error) {
	queriedClient, err := client(monitor, k8sClient, gitOps, desiredZitadel, desiredDatabase)
	if err != nil {
		return nil, err
	}

	list, err := queriedClient.ListUsers(monitor, k8sClient)
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
	return addUser(monitor, k8sClient, false, func() (*tree.Tree, error) {
		return zitadel.ReadCrd(k8sClient)
	}, func() (*tree.Tree, error) {
		return database.ReadCrd(k8sClient)
	}, user)
}

func GitOpsAddUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	return addUser(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, user)
}

func addUser(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	desiredZitadel func() (*tree.Tree, error),
	desiredDatabase func() (*tree.Tree, error),
	user string,
) error {

	queriedClient, err := client(monitor, k8sClient, gitOps, desiredZitadel, desiredDatabase)
	if err != nil {
		return err
	}

	return queriedClient.AddUser(monitor, user, k8sClient)

}

func GitOpsDeleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	return deleteUser(monitor, k8sClient, true, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, user)
}

func CrdDeleteUser(
	monitor mntr.Monitor,
	user string,
	k8sClient kubernetes.ClientInt,
) error {
	return deleteUser(monitor, k8sClient, false, func() (*tree.Tree, error) {
		return zitadel.ReadCrd(k8sClient)
	}, func() (*tree.Tree, error) {
		return database.ReadCrd(k8sClient)
	}, user)
}

func deleteUser(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	desiredZitadel func() (*tree.Tree, error),
	desiredDatabase func() (*tree.Tree, error),
	user string,
) error {

	queriedClient, err := client(monitor, k8sClient, gitOps, desiredZitadel, desiredDatabase)
	if err != nil {
		return err
	}

	return queriedClient.DeleteUser(monitor, user, k8sClient)
}
