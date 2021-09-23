package crtlcrd

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"

	databasev1 "github.com/caos/zitadel/operator/api/database/v1"
	zitadelv1 "github.com/caos/zitadel/operator/api/zitadel/v1"
	"github.com/caos/zitadel/operator/crtlcrd/database"
	"github.com/caos/zitadel/operator/crtlcrd/zitadel"
)

const (
	Database = "database"
	Zitadel  = "zitadel"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = databasev1.AddToScheme(scheme)
	_ = zitadelv1.AddToScheme(scheme)
}

func Start(monitor mntr.Monitor, version, metricsAddr string, features ...string) error {
	cfg := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "9adsd12l.caos.ch",
	})
	if err != nil {
		return fmt.Errorf("unable to start manager: %w", err)
	}

	k8sClient, err := kubernetes.NewK8sClientWithConfig(monitor, cfg)
	if err != nil {
		return err
	}

	for _, feature := range features {
		switch feature {
		case Database:
			if err = (&database.Reconciler{
				ClientInt: k8sClient,
				Monitor:   monitor,
				Scheme:    mgr.GetScheme(),
				Version:   version,
			}).SetupWithManager(mgr); err != nil {
				return fmt.Errorf("unable to create controller: %w", err)
			}
		case Zitadel:
			if err = (&zitadel.Reconciler{
				ClientInt: k8sClient,
				Monitor:   monitor,
				Scheme:    mgr.GetScheme(),
				Version:   version,
			}).SetupWithManager(mgr); err != nil {
				return fmt.Errorf("unable to create controller: %w", err)
			}
		}
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}
	return nil
}
