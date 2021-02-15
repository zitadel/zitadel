package orb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/orbos/pkg/treelabels"
	"github.com/caos/zitadel/operator"
	zitadelKubernetes "github.com/caos/zitadel/pkg/kubernetes"
	"github.com/pkg/errors"
)

func Reconcile(
	monitor mntr.Monitor,
	spec *Spec,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) (err error) {
		recMonitor := monitor.WithField("version", spec.Version)

		if spec.Version == "" {
			err := errors.New("No version provided for self-reconciling")
			recMonitor.Error(err)
			return err
		}

		imageRegistry := spec.CustomImageRegistry
		if imageRegistry == "" {
			imageRegistry = "ghcr.io"
		}

		if spec.SelfReconciling {
			desiredTree := &tree.Tree{
				Common: &tree.Common{
					Kind:    "databases.caos.ch/Orb",
					Version: "v0",
				},
			}

			if err := zitadelKubernetes.EnsureDatabaseArtifacts(monitor, treelabels.MustForAPI(desiredTree, mustDatabaseOperator(&spec.Version)), k8sClient, spec.Version, spec.NodeSelector, spec.Tolerations, imageRegistry, spec.GitOps); err != nil {
				recMonitor.Error(errors.Wrap(err, "Failed to deploy database-operator into k8s-cluster"))
				return err
			}
			recMonitor.Info("Applied database-operator")
		}
		return nil

	}
}
