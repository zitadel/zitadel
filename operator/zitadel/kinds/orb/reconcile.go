package orb

import (
	"github.com/caos/orbos/mntr"
	kubernetes2 "github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/orbos/pkg/treelabels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/pkg/kubernetes"
	"github.com/pkg/errors"
)

func Reconcile(
	monitor mntr.Monitor,
	spec *Spec,
	gitops bool,
) (
	operator.EnsureFunc,
	operator.DestroyFunc,
) {
	return func(k8sClient kubernetes2.ClientInt) (err error) {
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
						Kind:    "zitadel.caos.ch/Orb",
						Version: "v0",
					},
				}

				if err := kubernetes.EnsureZitadelOperatorArtifacts(monitor, treelabels.MustForAPI(desiredTree, mustZITADELOperator(&spec.Version)), k8sClient, spec.Version, spec.NodeSelector, spec.Tolerations, imageRegistry, gitops); err != nil {
					recMonitor.Error(errors.Wrap(err, "Failed to deploy zitadel-operator into k8s-cluster"))
					return err
				}
				recMonitor.Info("Applied zitadel-operator")
			}
			return nil
		}, func(k8sClient kubernetes2.ClientInt) error {
			if err := kubernetes.DestroyZitadelOperator(monitor, labels.MustForAPI(labels.NoopOperator("zitadel-operator"), "zitadel", "v0"), k8sClient, gitops); err != nil {
				monitor.Error(errors.Wrap(err, "Failed to destroy zitadel-operator in k8s-cluster"))
				return err
			}
			monitor.Info("Destroyed zitadel-operator")
			return nil
		}
}
