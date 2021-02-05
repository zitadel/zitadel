package controller

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	databasev1 "github.com/caos/zitadel/operator/api/database/v1"
	zitadelv1 "github.com/caos/zitadel/operator/api/zitadel/v1"
	"github.com/caos/zitadel/operator/controller/database"
	"github.com/caos/zitadel/operator/controller/zitadel"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
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
		return errors.Wrap(err, "unable to start manager")
	}

	for _, feature := range features {
		switch feature {
		case Database:
			if err = (&database.Reconciler{
				ClientInt: kubernetes.NewK8sClientWithConfig(monitor, cfg),
				Monitor:   monitor,
				Scheme:    mgr.GetScheme(),
				Version:   version,
			}).SetupWithManager(mgr); err != nil {
				return errors.Wrap(err, "unable to create controller")
			}
		case Zitadel:
			if err = (&zitadel.Reconciler{
				ClientInt: kubernetes.NewK8sClientWithConfig(monitor, cfg),
				Monitor:   monitor,
				Scheme:    mgr.GetScheme(),
				Version:   version,
			}).SetupWithManager(mgr); err != nil {
				return errors.Wrap(err, "unable to create controller")
			}
		}
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return errors.Wrap(err, "problem running manager")
	}
	return nil
}
