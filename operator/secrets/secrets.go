package secrets

import (
	"errors"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	zitadelOrb "github.com/caos/zitadel/operator/kinds/orb"
	"strings"
)

const (
	zitadel = "zitadel"
)

func GetAllSecretsFunc(orb *orb.Orb) func(monitor mntr.Monitor, gitClient *git.Client) (map[string]*secret.Secret, map[string]*tree.Tree, error) {
	return func(monitor mntr.Monitor, gitClient *git.Client) (map[string]*secret.Secret, map[string]*tree.Tree, error) {
		allSecrets := make(map[string]*secret.Secret, 0)
		allTrees := make(map[string]*tree.Tree, 0)
		foundBoom, err := api.ExistsZitadelYml(gitClient)
		if err != nil {
			return nil, nil, err
		}

		if foundBoom {
			zitadelYML, err := api.ReadZitadelYml(gitClient)
			if err != nil {
				return nil, nil, err
			}
			allTrees[zitadel] = zitadelYML
			_, _, zitadelSecrets, err := zitadelOrb.AdaptFunc(orb, "secret", "", "", []string{})(monitor, zitadelYML, &tree.Tree{})
			if err != nil {
				return nil, nil, err
			}

			if zitadelSecrets != nil && len(zitadelSecrets) > 0 {
				secret.AppendSecrets(zitadel, allSecrets, zitadelSecrets)
			}
		}
		return allSecrets, allTrees, nil
	}
}

func PushFunc() func(monitor mntr.Monitor, gitClient *git.Client, trees map[string]*tree.Tree, path string) error {
	return func(monitor mntr.Monitor, gitClient *git.Client, trees map[string]*tree.Tree, path string) error {
		operator := ""
		if strings.HasPrefix(path, zitadel) {
			operator = zitadel

		} else {
			return errors.New("Operator unknown")
		}

		desired, found := trees[operator]
		if !found {
			return errors.New("Operator file not found")
		}

		if operator == zitadel {
			return api.PushZitadelDesiredFunc(gitClient, desired)(monitor)
		}
		return errors.New("Operator push function unknown")
	}
}
