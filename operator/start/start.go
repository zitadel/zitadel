package start

import (
	"context"
	"github.com/caos/zitadel/operator/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"github.com/caos/zitadel/operator/zitadel"
	"time"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/caos/zitadel/pkg/databases"
	kubernetes2 "github.com/caos/zitadel/pkg/kubernetes"
)

func Operator(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, version *string) error {
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

		takeoff := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbConfig, "ensure", version, []string{"operator", "iam"}), k8sClient)

		go func() {
			started := time.Now()
			takeoff()

			monitor.WithFields(map[string]interface{}{
				"took": time.Since(started),
			}).Info("Iteration done")

			time.Sleep(time.Second * 10)
			takeoffChan <- struct{}{}
		}()
	}

	return nil
}

func Restore(monitor mntr.Monitor, gitClient *git.Client, orbCfg *orbconfig.Orb, k8sClient *kubernetes.Client, backup string, version *string) error {
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

	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaledown", version, []string{"scaledown"}), k8sClient)(); err != nil {
		return err
	}

	if err := databases.Clear(monitor, k8sClient, gitClient, databasesList); err != nil {
		return err
	}

	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "migration", version, []string{"migration"}), k8sClient)(); err != nil {
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

	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaleup", version, []string{"scaleup"}), k8sClient)(); err != nil {
		return err
	}

	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
		return err
	}

	return nil
}

func Database(monitor mntr.Monitor, orbConfigPath string, k8sClient *kubernetes.Client, binaryVersion *string) error {
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

		takeoff := database.Takeoff(monitor, gitClient, orbdb.AdaptFunc("", binaryVersion, "database", "backup"), k8sClient)

		go func() {
			started := time.Now()
			takeoff()

			monitor.WithFields(map[string]interface{}{
				"took": time.Since(started),
			}).Info("Iteration done")

			time.Sleep(time.Second * 10)
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

	database.Takeoff(monitor, gitClient, orbdb.AdaptFunc(backup, binaryVersion, "instantbackup"), k8sClient)()
	return nil
}
