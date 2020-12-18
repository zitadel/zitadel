package start

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/databases"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/orb"
	kubernetes2 "github.com/caos/zitadel/pkg/kubernetes"
)

func Operator(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, migrationsPath string, version *string) error {
	takeoffChan := make(chan struct{})
	go func() {
		takeoffChan <- struct{}{}
	}()

	for range takeoffChan {
		orbConfig, err := orbconfig.ParseOrbConfig(orbConfigPath)
		if err != nil {
			monitor.Error(err)
			return err
		}

		gitClient := git.New(context.Background(), monitor, "orbos", "orbos@caos.ch")
		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			monitor.Error(err)
			return err
		}

		takeoff := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(orbConfig, "ensure", migrationsPath, version, []string{"operator", "iam"}), k8sClient)

		go func() {
			started := time.Now()
			takeoff()

			monitor.WithFields(map[string]interface{}{
				"took": time.Since(started),
			}).Info("Iteration done")
			debug.FreeOSMemory()

			takeoffChan <- struct{}{}
		}()
	}

	return nil
}

func Restore(monitor mntr.Monitor, gitClient *git.Client, orbCfg *orbconfig.Orb, k8sClient *kubernetes.Client, backup, migrationsPath string, version *string) error {
	databasesList := []string{
		"notification",
		"adminapi",
		"auth",
		"authz",
		"eventstore",
		"management",
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 0); err != nil {
		return err
	}

	if err := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaledown", migrationsPath, version, []string{"scaledown"}), k8sClient)(); err != nil {
		return err
	}

	if err := databases.Clear(monitor, k8sClient, gitClient, databasesList); err != nil {
		return err
	}

	if err := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "migration", migrationsPath, version, []string{"migration"}), k8sClient)(); err != nil {
		return err
	}

	if err := databases.Restore(
		monitor,
		k8sClient,
		gitClient,
		backup,
		databasesList,
	); err != nil {
		return err
	}

	if err := operator.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaleup", migrationsPath, version, []string{"scaleup"}), k8sClient)(); err != nil {
		return err
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
		return err
	}

	return nil
}
