package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/pkg/databases"
	kubernetes2 "github.com/caos/zitadel/pkg/kubernetes"
	"time"
)

var (
	databasesList = []string{
		"notification",
		"adminapi",
		"auth",
		"authz",
		"eventstore",
		"management",
	}
	userList = []string{
		"notification",
		"adminapi",
		"auth",
		"authz",
		"eventstore",
		"management",
		"queries",
	}
)

func GitOpsClearMigrateRestore(
	monitor mntr.Monitor,
	gitClient *git.Client,
	orbCfg *orbconfig.Orb,
	k8sClient *kubernetes.Client,
	backup string,
	version *string,
) error {

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 0); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	if err := GitOpsScaleDown(monitor, orbCfg, gitClient, k8sClient, version); err != nil {
		return err
	}

	if err := databases.GitOpsClear(monitor, k8sClient, gitClient, databasesList, userList); err != nil {
		return err
	}

	if err := GitOpsMigrations(monitor, orbCfg, gitClient, k8sClient, version); err != nil {
		return err
	}

	if err := databases.GitOpsRestore(monitor, k8sClient, gitClient, backup, databasesList); err != nil {
		return err
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
		return err
	}

	return nil
}

func CrdClearMigrateRestore(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	backup string,
	version *string,
) error {

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 0); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	if err := CrdScaleDown(monitor, k8sClient, version); err != nil {
		return err
	}

	if err := databases.CrdClear(monitor, k8sClient, databasesList, userList); err != nil {
		return err
	}

	if err := CrdMigrations(monitor, k8sClient, version); err != nil {
		return err
	}

	if err := databases.CrdRestore(monitor, k8sClient, backup, databasesList); err != nil {
		return err
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
		return err
	}

	return nil
}
