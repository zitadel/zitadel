package statefulset

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/labels"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	"strings"
	"time"
)

func CleanPVCs(
	monitor mntr.Monitor,
	namespace string,
	sfsSelectable *labels.Selectable,
	replicaCount int,
) resources.QueryFunc {
	name := sfsSelectable.Name()
	return func(k8sClient kubernetes.ClientInt) (resources.EnsureFunc, error) {
		pvcs, err := k8sClient.ListPersistentVolumeClaims(namespace)
		if err != nil {
			return nil, err
		}
		internalPvcs := []string{}
		for _, pvc := range pvcs.Items {
			if strings.HasPrefix(pvc.Name, datadirInternal+"-"+name) {
				internalPvcs = append(internalPvcs, pvc.Name)
			}
		}
		return func(k8sClient kubernetes.ClientInt) error {
			noSFS := false
			monitor.Info("Scale down statefulset")
			if err := k8sClient.ScaleStatefulset(namespace, name, 0); err != nil {
				if macherrs.IsNotFound(err) {
					noSFS = true
				} else {
					return err
				}
			}
			time.Sleep(2 * time.Second)

			monitor.Info("Delete persistent volume claims")
			for _, pvcName := range internalPvcs {
				if err := k8sClient.DeletePersistentVolumeClaim(namespace, pvcName, cleanTimeout); err != nil {
					return err
				}
			}
			time.Sleep(2 * time.Second)

			if !noSFS {
				monitor.Info("Scale up statefulset")
				if err := k8sClient.ScaleStatefulset(namespace, name, replicaCount); err != nil {
					return err
				}
			}
			return nil
		}, nil
	}
}
