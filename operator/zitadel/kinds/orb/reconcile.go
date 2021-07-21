package orb

import (
	"errors"
	"fmt"

	"github.com/caos/orbos/mntr"
	kubernetes2 "github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/orbos/pkg/treelabels"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/pkg/kubernetes"
)

func Reconcile(
	monitor mntr.Monitor,
	spec *Spec,
	gitops bool,
) operator.EnsureFunc {
	return func(k8sClient kubernetes2.ClientInt) (err error) {
		recMonitor := monitor.WithField("version", spec.Version)

		if spec.Version == "" {
			return errors.New("no version provided for self-reconciling")
		}

		if spec.SelfReconciling {
			desiredTree := &tree.Tree{
				Common: tree.NewCommon("zitadel.caos.ch/Orb", "v0", false),
			}

			if err := kubernetes.EnsureZitadelOperatorArtifacts(monitor, treelabels.MustForAPI(desiredTree, mustZITADELOperator(&spec.Version)), k8sClient, spec.Version, spec.NodeSelector, spec.Tolerations, spec.CustomImageRegistry, gitops); err != nil {
				return fmt.Errorf("failed to deploy zitadel-operator into k8s-cluster: %w", err)
			}
			recMonitor.Info("Applied zitadel-operator")
		}
		return nil

	}
}
