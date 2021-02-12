package secrets

import (
	"errors"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"strings"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	zitadelOrb "github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

const (
	zitadel  string = "zitadel"
	database string = "database"
)

func GetAllSecretsFunc(orb *orb.Orb) func(monitor mntr.Monitor, gitClient *git.Client) (map[string]*secret.Secret, map[string]*tree.Tree, error) {
	return func(monitor mntr.Monitor, gitClient *git.Client) (map[string]*secret.Secret, map[string]*tree.Tree, error) {
		allSecrets := make(map[string]*secret.Secret, 0)
		allTrees := make(map[string]*tree.Tree, 0)
		foundZitadel, err := api.ExistsZitadelYml(gitClient)
		if err != nil {
			return nil, nil, err
		}

		if foundZitadel {
			zitadelYML, err := api.ReadZitadelYml(gitClient)
			if err != nil {
				return nil, nil, err
			}
			allTrees[zitadel] = zitadelYML
			_, _, zitadelSecrets, err := zitadelOrb.AdaptFunc(orb, "secret", nil, []string{})(monitor, zitadelYML, &tree.Tree{})
			if err != nil {
				return nil, nil, err
			}

			if zitadelSecrets != nil && len(zitadelSecrets) > 0 {
				secret.AppendSecrets(zitadel, allSecrets, zitadelSecrets)
			}
		} else {
			monitor.Info("no file for zitadel found")
		}

		foundDB, err := api.ExistsDatabaseYml(gitClient)
		if err != nil {
			return nil, nil, err
		}
		if foundDB {
			dbYML, err := api.ReadDatabaseYml(gitClient)
			if err != nil {
				return nil, nil, err
			}
			allTrees[database] = dbYML

			_, _, dbSecrets, err := orbdb.AdaptFunc("", nil, "database", "backup")(monitor, dbYML, nil)
			if err != nil {
				return nil, nil, err
			}
			if dbSecrets != nil && len(dbSecrets) > 0 {
				secret.AppendSecrets(database, allSecrets, dbSecrets)
			}
		} else {
			monitor.Info("no file for database found")
		}
		return allSecrets, allTrees, nil
	}
}

func PushFunc() func(monitor mntr.Monitor, gitClient *git.Client, trees map[string]*tree.Tree, path string) error {
	return func(monitor mntr.Monitor, gitClient *git.Client, trees map[string]*tree.Tree, path string) error {
		operator := ""
		if strings.HasPrefix(path, zitadel) {
			operator = zitadel
		} else if strings.HasPrefix(path, database) {
			operator = database
		} else {
			return errors.New("Operator unknown")
		}

		desired, found := trees[operator]
		if !found {
			return errors.New("Operator file not found")
		}

		if operator == zitadel {
			return api.PushZitadelDesiredFunc(gitClient, desired)(monitor)
		} else if operator == database {
			return api.PushDatabaseDesiredFunc(gitClient, desired)(monitor)
		}
		return errors.New("Operator push function unknown")
	}
}
