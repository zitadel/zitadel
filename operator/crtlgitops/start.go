package crtlgitops

import (
	"context"
	"time"

	"github.com/caos/zitadel/operator/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"github.com/caos/zitadel/operator/zitadel"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
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

			time.Sleep(time.Second * 10)
			takeoffChan <- struct{}{}
		}()
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

			time.Sleep(time.Second * 10)
			takeoffChan <- struct{}{}
		}()
	}

	return nil
}
