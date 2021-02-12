package start

import (
	"context"
	"github.com/caos/zitadel/operator/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"github.com/caos/zitadel/operator/zitadel"
	"runtime/debug"
	"time"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/caos/zitadel/pkg/databases"
	kubernetes2 "github.com/caos/zitadel/pkg/kubernetes"
)

func Operator(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, version *string, gitops bool) error {
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

		takeoff := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbConfig, "ensure", version, gitops, []string{"operator", "iam"}), k8sClient)

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

func Restore(
	monitor mntr.Monitor,
	gitClient *git.Client,
	orbCfg *orbconfig.Orb,
	k8sClient *kubernetes.Client,
	backup string,
	gitops bool,
	version *string,
) error {
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

	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaledown", version, gitops, []string{"scaledown"}), k8sClient)(); err != nil {
		return err
	}

	if err := databases.Clear(monitor, k8sClient, gitClient, databasesList); err != nil {
		return err
	}

	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "migration", version, gitops, []string{"migration"}), k8sClient)(); err != nil {
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

	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaleup", version, gitops, []string{"scaleup"}), k8sClient)(); err != nil {
		return err
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
		return err
	}

	return nil
}

func Database(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, binaryVersion *string, gitops bool) error {
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

		takeoff := database.Takeoff(monitor, gitClient, orbdb.AdaptFunc("", binaryVersion, gitops, "operator", "database", "backup"), k8sClient)

		go func() {
			started := time.Now()
			takeoff()

			monitor.WithFields(map[string]interface{}{
				"took": time.Since(started),
			}).Info("Iteration done")

			takeoffChan <- struct{}{}
		}()
	}

	return nil
}

func Backup(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, backup string, binaryVersion *string) error {
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

	database.Takeoff(monitor, gitClient, orbdb.AdaptFunc(backup, binaryVersion, false, "instantbackup"), k8sClient)()
	return nil
}
