package orb

import (
	"errors"
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/orbos/pkg/treelabels"

	"github.com/caos/zitadel/operator"
	zitadelKubernetes "github.com/caos/zitadel/pkg/kubernetes"
)

func Reconcile(
	monitor mntr.Monitor,
	spec *Spec,
	gitops bool,
) (
	operator.EnsureFunc,
	operator.DestroyFunc,
) {
	return func(k8sClient kubernetes.ClientInt) (err error) {
			recMonitor := monitor.WithField("version", spec.Version)

			if spec.Version == "" {
				return errors.New("no version provided for self-reconciling")
			}

			if spec.SelfReconciling {
				desiredTree := &tree.Tree{
					Common: tree.NewCommon("databases.caos.ch/Orb", "v0", false),
				}

				if err := zitadelKubernetes.EnsureDatabaseArtifacts(monitor, treelabels.MustForAPI(desiredTree, mustDatabaseOperator(&spec.Version)), k8sClient, spec.Version, spec.NodeSelector, spec.Tolerations, spec.CustomImageRegistry, gitops); err != nil {
					return fmt.Errorf("failed to deploy database-operator into k8s-cluster: %w", err)
				}
				recMonitor.Info("Applied database-operator")
			}
			return nil
		}, func(k8sClient kubernetes.ClientInt) error {
			if err := zitadelKubernetes.DestroyDatabaseOperator(monitor, labels.MustForAPI(labels.NoopOperator("database-operator"), "database", "v0"), k8sClient, gitops); err != nil {
				return fmt.Errorf("failed to destroy database-operator in k8s-cluster: %w", err)
			}
			monitor.Info("Destroyed database-operator")
			return nil
		}
}
