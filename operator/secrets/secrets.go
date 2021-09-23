package secrets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caos/orbos/pkg/kubernetes"

	crddb "github.com/caos/zitadel/operator/api/database"
	crdzit "github.com/caos/zitadel/operator/api/zitadel"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
)

const (
	zitadel  string = "zitadel"
	database string = "database"
)

func GetAllSecretsFunc(
	monitor mntr.Monitor,
	printLogs,
	gitops bool,
	gitClient *git.Client,
	k8sClient kubernetes.ClientInt,
	orb *orb.Orb,
) func() (
	map[string]*secret.Secret,
	map[string]*secret.Existing,
	map[string]*tree.Tree,
	error,
) {
	return func() (
		map[string]*secret.Secret,
		map[string]*secret.Existing,
		map[string]*tree.Tree,
		error,
	) {
		return getAllSecrets(monitor, printLogs, gitops, orb, gitClient, k8sClient)
	}
}

func getAllSecrets(
	monitor mntr.Monitor,
	printLogs,
	gitops bool,
	orb *orb.Orb,
	gitClient *git.Client,
	k8sClient kubernetes.ClientInt,
) (
	map[string]*secret.Secret,
	map[string]*secret.Existing,
	map[string]*tree.Tree,
	error,
) {
	allSecrets := make(map[string]*secret.Secret, 0)
	allExisting := make(map[string]*secret.Existing, 0)
	allTrees := make(map[string]*tree.Tree, 0)

	if err := secret.GetOperatorSecrets(
		monitor,
		printLogs,
		gitops,
		gitClient,
		git.ZitadelFile,
		allTrees,
		allSecrets,
		allExisting,
		func() (t *tree.Tree, err error) { return crdzit.ReadCrd(k8sClient) },
		func(t *tree.Tree) (map[string]*secret.Secret, map[string]*secret.Existing, bool, error) {
			_, _, _, secrets, existing, migrate, err := orbzit.AdaptFunc(orb, "secret", nil, gitops, []string{})(monitor, t, &tree.Tree{})
			return secrets, existing, migrate, err
		},
	); err != nil {
		return nil, nil, nil, err
	}

	if err := secret.GetOperatorSecrets(
		monitor,
		printLogs,
		gitops,
		gitClient,
		git.DatabaseFile,
		allTrees,
		allSecrets,
		allExisting,
		func() (t *tree.Tree, err error) { return crddb.ReadCrd(k8sClient) },
		func(t *tree.Tree) (map[string]*secret.Secret, map[string]*secret.Existing, bool, error) {
			_, _, _, secrets, existing, migrate, err := orbdb.AdaptFunc("", nil, gitops, "database", "backup")(monitor, t, nil)
			return secrets, existing, migrate, err
		},
	); err != nil {
		return nil, nil, nil, err
	}

	if k8sClient == nil {
		allExisting = nil
	}

	if len(allSecrets) == 0 && len(allExisting) == 0 {
		return nil, nil, nil, mntr.ToUserError(errors.New("couldn't find any secrets"))
	}

	return allSecrets, allExisting, allTrees, nil
}

func PushFunc(
	monitor mntr.Monitor,
	gitops bool,
	gitClient *git.Client,
	k8sClient kubernetes.ClientInt,
) func(
	trees map[string]*tree.Tree,
	path string,
) error {
	return func(
		trees map[string]*tree.Tree,
		path string,
	) error {
		return push(monitor, gitops, gitClient, k8sClient, trees, path)
	}
}

func push(
	monitor mntr.Monitor,
	gitops bool,
	gitClient *git.Client,
	k8sClient kubernetes.ClientInt,
	trees map[string]*tree.Tree,
	path string,
) error {

	var (
		applyCRDFunc func(*tree.Tree) error
		desiredFile  git.DesiredFile
	)

	if strings.HasPrefix(path, zitadel) {
		desiredFile = git.ZitadelFile
		applyCRDFunc = func(t *tree.Tree) error {
			return crdzit.WriteCrd(k8sClient, t)
		}
	} else if strings.HasPrefix(path, database) {
		desiredFile = git.DatabaseFile
		applyCRDFunc = func(t *tree.Tree) error {
			return crddb.WriteCrd(k8sClient, t)
		}
	} else {
		return errors.New("operator unknown")
	}

	desired, found := trees[desiredFile.WOExtension()]
	if !found {
		return mntr.ToUserError(fmt.Errorf("desired state not found for %s", desiredFile.WOExtension()))
	}

	if gitops {
		return gitClient.PushDesiredFunc(desiredFile, desired)(monitor)
	}
	return applyCRDFunc(desired)
}
