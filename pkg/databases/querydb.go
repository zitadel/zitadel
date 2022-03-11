package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	coredb "github.com/caos/zitadel/operator/database/kinds/databases/core"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

func queryDatabase(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitOps bool,
	databaseTree func() (*tree.Tree, error),
	zitadelTree func() (*tree.Tree, error),
) (coredb.DatabaseCurrent, map[string]interface{}, error) {
	current := &tree.Tree{}

	var query operator.QueryFunc
	dbTree, err := databaseTree()
	if err != nil {
		return nil, nil, err
	}
	if dbTree != nil && dbTree.Original != nil {
		query, _, _, _, _, _, err = orbdb.AdaptFunc("", nil, gitOps, "database")(monitor, dbTree, current)
		if err != nil {
			return nil, nil, err
		}
	} else {
		zitTree, err := zitadelTree()
		if err != nil {
			return nil, nil, err
		}

		query, _, _, _, _, _, err = orbzit.AdaptFunc("", nil, gitOps, []string{"dbconnection"}, nil)(monitor, zitTree, current)
		if err != nil {
			return nil, nil, err
		}
	}

	queried := map[string]interface{}{}
	_, err = query(k8sClient, queried)
	if err != nil {
		return nil, nil, err
	}
	client, err := coredb.ParseQueriedForDatabase(queried)
	return client, queried, err
}
