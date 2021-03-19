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
	"github.com/caos/zitadel/operator/api"
)

const (
	zitadel  string = "zitadel"
	database string = "database"
)

func GetAllSecretsFunc(
	monitor mntr.Monitor,
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
		return getAllSecrets(monitor, gitops, orb, gitClient, k8sClient)
	}
}

func getAllSecrets(
	monitor mntr.Monitor,
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
		gitops,
		allTrees,
		allSecrets,
		allExisting,
		zitadel,
		func() (bool, error) { return api.ExistsZitadelYml(gitClient) },
		func() (t *tree.Tree, err error) { return api.ReadZitadelYml(gitClient) },
		func() (t *tree.Tree, err error) { return crdzit.ReadCrd(k8sClient) },
		func(t *tree.Tree) (map[string]*secret.Secret, map[string]*secret.Existing, bool, error) {
			_, _, secrets, existing, migrate, err := orbzit.AdaptFunc(orb, "secret", nil, gitops, []string{})(monitor, t, &tree.Tree{})
			return secrets, existing, migrate, err
		},
	); err != nil {
		return nil, nil, nil, err
	}

	if err := secret.GetOperatorSecrets(
		monitor,
		gitops,
		allTrees,
		allSecrets,
		allExisting,
		database,
		func() (bool, error) { return api.ExistsDatabaseYml(gitClient) },
		func() (t *tree.Tree, err error) { return api.ReadDatabaseYml(gitClient) },
		func() (t *tree.Tree, err error) { return crddb.ReadCrd(k8sClient) },
		func(t *tree.Tree) (map[string]*secret.Secret, map[string]*secret.Existing, bool, error) {
			_, _, secrets, existing, migrate, err := orbdb.AdaptFunc("", nil, gitops, "database", "backup")(monitor, t, nil)
			return secrets, existing, migrate, err
		},
	); err != nil {
		return nil, nil, nil, err
	}

	if k8sClient == nil {
		allExisting = nil
	}

	if len(allSecrets) == 0 && len(allExisting) == 0 {
		return nil, nil, nil, errors.New("couldn't find any secrets")
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
		pushGitFunc  func(*tree.Tree) error
		applyCRDFunc func(*tree.Tree) error
		operator     string
	)

	if strings.HasPrefix(path, zitadel) {
		operator = zitadel
		pushGitFunc = func(desired *tree.Tree) error {
			return api.PushZitadelDesiredFunc(gitClient, desired)(monitor)
		}
		applyCRDFunc = func(t *tree.Tree) error {
			return crdzit.WriteCrd(k8sClient, t)
		}
	} else if strings.HasPrefix(path, database) {
		operator = database
		pushGitFunc = func(desired *tree.Tree) error {
			return api.PushDatabaseDesiredFunc(gitClient, desired)(monitor)
		}
		applyCRDFunc = func(t *tree.Tree) error {
			return crddb.WriteCrd(k8sClient, t)
		}
	} else {
		return errors.New("operator unknown")
	}

	desired, found := trees[operator]
	if !found {
		return fmt.Errorf("desired state for %s not found", operator)
	}

	if gitops {
		return pushGitFunc(desired)
	}
	return applyCRDFunc(desired)
}
